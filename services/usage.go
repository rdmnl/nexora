package services

import (
	"log"
	"time"

	"github.com/rdmnl/nexora/shared"
	"github.com/shirou/gopsutil/process"
)

func UpdateServerUsage(servers []shared.ServerInfo) {
	for i := range servers {
		server := &servers[i]

		proc, err := process.NewProcess(server.ProcessID)
		if err != nil {
			// log.Printf("Error accessing process %d: %v", server.ProcessID, err)
			continue
		}

		time.Sleep(100 * time.Millisecond)
		cpuPercent, err := proc.CPUPercent()
		if err != nil {
			log.Printf("Error getting CPU usage for process %d: %v", server.ProcessID, err)
			continue
		}

		memInfo, err := proc.MemoryInfo()
		if err != nil {
			log.Printf("Error getting memory info for process %d: %v", server.ProcessID, err)
			continue
		}

		server.CPUUsage = cpuPercent
		server.MemoryUsage = memInfo.RSS / (1024 * 1024) // Convert to MB for readability
	}
}
