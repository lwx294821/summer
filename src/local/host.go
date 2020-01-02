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
	IP  string
	Flags string
	Ether string
}

//使用pcap包需要提前准备 yum install -y libpcap-devel
func FindAllNetWorkDevs() []Device{
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}
	var devs  []Device
	for _, d := range devices {
		var ip = findDevIpv(d)
		if ip == "" {
			continue
		}
		macaddr,flag := findDevMacAddrByIp(ip)
		if macaddr == nil{
			continue
		}
		devs=append(devs,Device{Name:d.Name,IP:ip,Flags:flag.String(),Ether:macaddr.String()})
	}
	return devs
}

//获取设备的IPv4或者IPv6,不包括loop设备
func findDevIpv(device pcap.Interface) string {
	for _, addr := range device.Addresses {
		if addr.IP.IsLoopback(){
			return ""
		}else {
			return addr.IP.String()
		}
	}
	return ""
}

//根据IP获取网卡的MAC地址
func findDevMacAddrByIp(ip string) (net.HardwareAddr,net.Flags) {
	interfaces, err := net.Interfaces()
	if err != nil {
       return nil,0
	}
	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Println(err)
			return nil,0
		}
		for _, addr := range addrs {
			if a, ok := addr.(*net.IPNet); ok {
				if ip == a.IP.String() {
					return i.HardwareAddr,i.Flags
				}
			}
		}
	}
	return nil,0
}
