package packet

import (
"context"
"errors"
"fmt"
"github.com/google/gopacket"
"github.com/google/gopacket/layers"
"github.com/google/gopacket/pcap"
"log"
"net"
"time"
)

type NetPkt struct {
	LocalIP    net.IP
	RemoteIP   net.IP
	LocalPort  layers.TCPPort
	RemotePort layers.TCPPort
	Direct     string
	Len        uint16
}

var (
	snapshot_len int32 = 65536
	promiscuous  bool  = false
	err          error
	timeout      time.Duration = 30 * time.Second
	handle       *pcap.Handle

	ethLayer layers.Ethernet
	ipLayer  layers.IPv4
	tcpLayer layers.TCP


	LocalIP, RemoteIP     net.IP
	LocalPort, RemotePort layers.TCPPort
	Dir                   string
	pLen                  uint16

	PkgAcc  map[string]map[string]int64
	AccIntv int64 = 3
	Ctx     context.Context
	Cancel  context.CancelFunc

	Start bool
)

func init() {
	PkgAcc = make(map[string]map[string]int64)
}

func StartNetSniff(ipAddr string,ctx context.Context) {
	if Start {
		return
	}
	Start = true
	PkgAcc = make(map[string]map[string]int64)

	ctx, cancel := context.WithTimeout(ctx,time.Duration(150)*time.Second)
	Ctx = ctx
	Cancel = cancel
	go startNetSniff(ctx, ipAddr)

}

func StopNetSniff() {
	if !Start {
		return
	}

	Cancel()
	PkgAcc = make(map[string]map[string]int64)
}

func formatNetPck(p map[string]map[string]int64) {
	for itemId, info := range p{
		log.Printf(" %-40s %-6d %-6d %-6d %-6d \n",
			itemId, info["in"],info["out"], info["intRate"], info["outRate"])
	}
}

func startNetSniff(ctx context.Context, ipAddr string) {
	handle, err := getPcapHandle(ipAddr)
	if err != nil {
		log.Println(err.Error())
		return
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	go accumulator(ctx)

	for packet := range packetSource.Packets() {
		select {
		case <-ctx.Done():
			log.Println("⇨ Packet sniff Stop")
			handle.Close()
			return
		default:
			parser := gopacket.NewDecodingLayerParser(
				layers.LayerTypeEthernet,
				&ethLayer,
				&ipLayer,
				&tcpLayer,
			)
			foundLayerTypes := []gopacket.LayerType{}

			err := parser.DecodeLayers(packet.Data(), &foundLayerTypes)
			if err != nil {
				//log.Println("Trouble decoding layers: ", err)
			}
			for _, layerType := range foundLayerTypes {
				if layerType == layers.LayerTypeIPv4 {
					if ipLayer.SrcIP.String() != ipAddr {
						LocalIP = ipLayer.DstIP
						RemoteIP = ipLayer.SrcIP
						LocalPort = tcpLayer.DstPort
						RemotePort = tcpLayer.SrcPort
						pLen = ipLayer.Length

						Dir = "in"
						itemId := fmt.Sprintf("%s:%d-%s:%d", LocalIP, LocalPort, RemoteIP, RemotePort)

						if _, ok := PkgAcc[itemId]; !ok {
							PkgAcc[itemId] = map[string]int64{
								"in":          0,
								"out":         0,
								"inRate":      0,
								"outRate":     0,
								"lastAccTime": 0,
								"lastAccIn":   0,
								"lastAccOut":  0,
							}
						}

						PkgAcc[itemId][Dir] = PkgAcc[itemId][Dir] + int64(pLen)
					} else {
						LocalIP = ipLayer.SrcIP
						RemoteIP = ipLayer.DstIP
						LocalPort = tcpLayer.SrcPort
						RemotePort = tcpLayer.DstPort
						pLen = ipLayer.Length

						Dir = "out"

						itemId := fmt.Sprintf("%s:%d-%s:%d", LocalIP, LocalPort, RemoteIP, RemotePort)

						if _, ok := PkgAcc[itemId]; !ok {
							PkgAcc[itemId] = map[string]int64{
								"in":          0,
								"out":         0,
								"inRate":      0,
								"outRate":     0,
								"lastAccTime": 0,
								"lastAccIn":   0,
								"lastAccOut":  0,
							}
						}
						PkgAcc[itemId][Dir] = PkgAcc[itemId][Dir] + int64(pLen)
					}
				}
			}
		}
	}
}

func getPcapHandle(ip string) (*pcap.Handle, error) {
	devs, err := pcap.FindAllDevs()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	var device string
	for _, dev := range devs {
		for _, v := range dev.Addresses {
			if v.IP.String() == ip {
				device = dev.Name
				break
			}
		}
	}
	if device == "" {
		return nil, errors.New("find device error")
	}

	h, err := pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
	if err != nil {
		return nil, err
	}
	log.Println("⇨ Starting net traffic monitor...")
	var filter string = "tcp and (not broadcast and not multicast)"
	err = h.SetBPFFilter(filter)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func accumulator(ctx context.Context) {
	log.Println("⇨ Accumulator thread is starting...")
	for {
		select {
		case <-ctx.Done():
			log.Println("⇨ Accumulator thread is stop.")
			Start = false
			return
		default:
			for _, pkgMap := range PkgAcc {
				start := time.Now().Unix()
				in := pkgMap["in"]
				out := pkgMap["out"]

				if pkgMap["lastAccTime"] == 0 {
					pkgMap["lastAccTime"] = start - AccIntv
				}

				last := pkgMap["lastAccTime"]
				pkgMap["lastAccTime"] = start

				durSec := start - last

				if in == 0 {
					pkgMap["inRate"] = 0
				} else {
					pkgMap["inRate"] = (in - pkgMap["lastAccIn"]) / durSec
				}

				if out == 0 {
					pkgMap["outRate"] = 0
				} else {
					pkgMap["outRate"] = (out - pkgMap["lastAccOut"]) / durSec

				}

				pkgMap["lastAccIn"] = in
				pkgMap["lastAccOut"] = out
			}
			formatNetPck(PkgAcc)
			time.Sleep(time.Duration(AccIntv) * time.Second)

		}

	}
}
