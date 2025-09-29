package core

import (
	"sync"
	"time"

	nu "netscan-server/internal/net"
)

type Scanner struct {
	oui nu.VendorDB
}

func NewScanner(oui nu.VendorDB) *Scanner {
	return &Scanner{oui: oui}
}

func (s *Scanner) Run(req ScanRequest) ([]HostResult, error) {
	if req.Ports == "" {
		req.Ports = "1-1024"
	}
	if req.Concurrency == 0 {
		req.Concurrency = 100
	}
	if req.TimeoutMs == 0 {
		req.TimeoutMs = 400
	}

	ports, err := nu.ParsePortSpec(req.Ports)
	if err != nil {
		return nil, err
	}

	// --- DISCOVERY (same as you currently have; shortened for brevity) ---
	hostsMap := make(map[string]string) // ip -> mac

	if req.Discovery == "arp-cache" {
		ips, err := nu.HostsFromCIDR(req.Target)
		if err != nil {
			return nil, err
		}
		for _, ip := range ips {
			go nu.Ping(ip)
		}
		time.Sleep(500 * time.Millisecond)
		arpCache := nu.ReadARP()
		for ip, mac := range arpCache {
			hostsMap[ip] = mac
		}
	} else {
		arpEntries, err := nu.ARPSweep(req.Target, 1200*time.Millisecond)
		if err != nil {
			return nil, err
		}
		for ip, e := range arpEntries {
			hostsMap[ip] = e.MAC.String()
		}
	}

	if len(hostsMap) == 0 {
		ips, err := nu.HostsFromCIDR(req.Target)
		if err != nil {
			return nil, err
		}
		for _, ip := range ips {
			if nu.TryConnect(ip, 80, 250*time.Millisecond) || nu.TryConnect(ip, 443, 250*time.Millisecond) ||
				(len(ports) > 0 && nu.TryConnect(ip, ports[0], 300*time.Millisecond)) {
				hostsMap[ip] = ""
			}
		}
	}

	// --- ENRICH + SCAN ---
	var (
		results []HostResult
		mu      sync.Mutex
		wg      sync.WaitGroup
	)
	hostSema := make(chan struct{}, 200)

	for ip, mac := range hostsMap {
		wg.Add(1)
		go func(ip, mac string) {
			defer wg.Done()
			hostSema <- struct{}{}
			defer func() { <-hostSema }()

			start := time.Now()
			hostname := nu.ReverseLookup(ip)
			open := nu.ScanPorts(ip, ports, time.Duration(req.TimeoutMs)*time.Millisecond, req.Concurrency)

			// vendor lookup
			vendor := ""
			if s.oui != nil {
				vendor = nu.VendorFromMAC(s.oui, mac)
			}

			mu.Lock()
			results = append(results, HostResult{
				IP:             ip,
				Hostname:       hostname,
				MAC:            mac,
				Vendor:         vendor,
				OpenPorts:      open,
				ScanStartedISO: start.Format(time.RFC3339),
				ScanDurationMS: time.Since(start).Milliseconds(),
			})
			mu.Unlock()
		}(ip, mac)
	}
	wg.Wait()
	return results, nil
}
