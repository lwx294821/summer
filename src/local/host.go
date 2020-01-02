package local

import (
	"github.com/google/gopacket/pcap"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"log"
	"net"
)

func NodeInfo() {
	h, _ := host.Info()
	c, _ := cpu.Info()
	m, _ := mem.VirtualMemory()
	d, _ := disk.Partitions(true)
	dio, _ := disk.IOCounters()
	log.Println(h.Hostname, c, m, d, dio)
}

func Ifconfig() bool {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		log.Println("net.Interfaces failed, err:", err.Error())
		return false
	}
	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			log.Println(netInterfaces[i].Name,netInterfaces[i].Flags.String())
			addrs, _ := netInterfaces[i].Addrs()
			for _, address := range addrs {
				//address.(*net.IPNet) 类型断言,即判断该变量是否为指定类型同时对变量进行类型转换
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						log.Println(ipnet.IP.String())
						return true
					}
				}
			}
		}
	}
	return false
}

type Device struct {
	Name  string
	IPv4  string
	IPv6  string
	Flags string
}

//使用pcap包需要提前准备 yum install -y libpcap-devel
func FindAllNetWorkDevs() {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}
	for _, d := range devices {
		var ip = findDevIpv(d)
		if ip == "" {
			continue
		}
		var macaddr = findDevMacAddrByIp(ip)
		log.Println("name:"+d.Name, "ip:"+ip, "macaddr:"+macaddr)
		for _, addr := range d.Addresses {
			if !addr.IP.IsLoopback() {

			}
		}
	}
}

//获取设备的IPv4地址
func findDevIpv(device pcap.Interface) string {
	for _, addr := range device.Addresses {
		if ipv4 := addr.IP.To4(); ipv4 != nil {
			if addr.IP.IsLoopback() {
				return ""
			}
			return ipv4.String()
		} else if ipv6 := addr.IP.To16(); ipv6 != nil {
			return ipv6.String()
		}
	}
	return ""
}

//根据IP获取网卡的MAC地址
func findDevMacAddrByIp(ip string) string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Println(err)
			return ""
		}
		for _, addr := range addrs {
			if a, ok := addr.(*net.IPNet); ok {
				if ip == a.IP.String() {
					return i.HardwareAddr.String()
				}
			}
		}
	}
	return ""
}
