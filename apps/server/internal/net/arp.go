package net

// package netutil

import (
	"bufio"
	"context"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func Ping(ip string) { // populate ARP using system ping
	ctx, cancel := context.WithTimeout(context.Background(), 1200*time.Millisecond)
	defer cancel()
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "ping", "-n", "1", "-w", "1000", ip)
	} else {
		cmd = exec.CommandContext(ctx, "ping", "-c", "1", "-W", "1", ip)
	}
	_ = cmd.Run()
}

func ReadARP() map[string]string {
	out := make(map[string]string)
	if runtime.GOOS == "linux" {
		if b, err := exec.Command("ip", "neigh").Output(); err == nil {
			sc := bufio.NewScanner(strings.NewReader(string(b)))
			for sc.Scan() {
				parts := strings.Fields(sc.Text())
				if len(parts) >= 5 {
					ip := parts[0]
					for i := 0; i < len(parts); i++ {
						if parts[i] == "lladdr" && i+1 < len(parts) {
							out[ip] = parts[i+1]
							break
						}
					}
				}
			}
			return out
		}
	}
	// fallback: arp -a
	b, err := exec.Command("arp", "-a").Output()
	if err != nil {
		return out
	}
	sc := bufio.NewScanner(strings.NewReader(string(b)))
	for sc.Scan() {
		line := sc.Text()
		if strings.Contains(line, "(") && strings.Contains(line, ")") && strings.Contains(line, " at ") {
			parts := strings.Split(line, ")")
			ip := strings.TrimPrefix(parts[0], "? (")
			after := parts[1]
			idx := strings.Index(after, " at ")
			if idx >= 0 {
				mac := strings.Fields(after[idx+4:])[0]
				out[ip] = mac
			}
		} else {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				ip := fields[0]
				out[ip] = strings.ReplaceAll(fields[1], "-", ":")
			}
		}
	}
	return out
}
