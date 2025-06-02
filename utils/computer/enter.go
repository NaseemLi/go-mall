package computer

import (
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

func GetCpuPercent() float64 {
	cpuPercent, _ := cpu.Percent(time.Second, false)
	return cpuPercent[0]
}

func GetMemPercent() float64 {
	memInfo, _ := mem.VirtualMemory()
	return memInfo.UsedPercent
}

func GetDiskPercent() (percent float64) {
	diskList, _ := disk.Partitions(false)
	var used uint64
	var total uint64
	for _, stat := range diskList {
		usage, _ := disk.Usage(stat.Mountpoint)
		used += usage.Used
		total += usage.Total
	}
	percent = (float64(used) / float64(total)) * 100
	return
}
