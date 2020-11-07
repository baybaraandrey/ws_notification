package rest

type VirtualMemory struct {
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"used_percent"`
}

type Metrics struct {
	Hostname          string        `json:"hostname"`
	Time              string        `json:"time"`
	VirtualMemoryStat VirtualMemory `json:"virtual_memory_stat"`
}
