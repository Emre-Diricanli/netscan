package net

// package netutil

import (
	"time"

	"github.com/go-ping/ping"
)

func PingRTT(ip string, timeout time.Duration) (time.Duration, bool) {
	p := ping.New(ip)
	p.Count = 1
	p.Timeout = timeout
	p.SetPrivileged(true)
	if err := p.Run(); err != nil {
		return 0, false
	}
	stats := p.Statistics()
	if stats.PacketsRecv < 1 {
		return 0, false
	}
	return stats.AvgRtt, true
}
