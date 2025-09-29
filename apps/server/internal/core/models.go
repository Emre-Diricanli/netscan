package core

type ScanRequest struct {
	Target      string `json:"target"`
	Ports       string `json:"ports"`
	Concurrency int    `json:"concurrency,omitempty"`
	TimeoutMs   int    `json:"timeout_ms,omitempty"`
	Discovery   string `json:"discovery,omitempty"` // "arp-raw" | "arp-cache"
}

type HostResult struct {
	IP             string `json:"ip"`
	Hostname       string `json:"hostname,omitempty"`
	MAC            string `json:"mac,omitempty"`
	Vendor         string `json:"vendor,omitempty"` // <--- make sure this exists
	PingMs         int64  `json:"ping_ms,omitempty"`
	OpenPorts      []int  `json:"open_ports"`
	ScanStartedISO string `json:"scan_started"`
	ScanDurationMS int64  `json:"scan_duration_ms"`
}
