package services

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/rdmnl/nexora/shared"
)

func DetectAndMerge(existing []shared.ServerInfo, configNodes []shared.ServerInfo) []shared.ServerInfo {
	detected := detectRunningServers()
	return mergeWithConfig(detected, configNodes)
}

func detectRunningServers() []shared.ServerInfo {
	cmd := exec.Command("lsof", "-i", "-P", "-n")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		log.Printf("Error running lsof: %v", err)
		return nil
	}

	lines := strings.Split(out.String(), "\n")
	serverMap := make(map[int]shared.ServerInfo)
	currentPID := os.Getpid()

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}

		pid := toInt32(fields[1])
		port := extractPort(fields[8])

		if int(pid) == currentPID || port == 0 {
			continue
		}

		serverMap[port] = shared.ServerInfo{
			Name:      fields[0],
			Port:      port,
			ProcessID: pid,
		}
	}
	return mapToSlice(serverMap)
}

func extractPort(address string) int {
	parts := strings.Split(address, ":")
	if len(parts) < 2 {
		return 0
	}
	port, _ := strconv.Atoi(parts[1])
	return port
}

func mapToSlice(serverMap map[int]shared.ServerInfo) []shared.ServerInfo {
	serverList := make([]shared.ServerInfo, 0, len(serverMap))
	for _, server := range serverMap {
		serverList = append(serverList, server)
	}
	return serverList
}

func toInt32(s string) int32 {
	i, _ := strconv.Atoi(s)
	return int32(i)
}

func mergeWithConfig(detected, configNodes []shared.ServerInfo) []shared.ServerInfo {
	serverMap := make(map[int]shared.ServerInfo)

	for _, server := range detected {
		serverMap[server.Port] = server
	}

	for _, node := range configNodes {
		if server, exists := serverMap[node.Port]; exists {
			server.Name = node.Name
			serverMap[node.Port] = server
		} else {
			serverMap[node.Port] = node
		}
	}

	return mapToSlice(serverMap)
}
