package shared

type ServerInfo struct {
	Name        string  `json:"name"`
	Port        int     `json:"port"`
	ProcessID   int32   `json:"process_id"`
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage uint64  `json:"memory_usage"`
}

type Config struct {
	Nodes []ServerInfo `yaml:"nodes"`
}
