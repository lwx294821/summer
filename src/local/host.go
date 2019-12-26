package local

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"log"
)

func Machine(){
	h, _ := host.Info()
	c, _ := cpu.Info()
	m, _ := mem.VirtualMemory()
	d, _ := disk.Partitions(true)
	dio, _ := disk.IOCounters()
	log.Println(h.Hostname,c,m,d,dio)



}
