package local

import (
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
