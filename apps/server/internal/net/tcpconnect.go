package net

// package netutil

import (
	"net"
	"sort"
	"strconv"
	"sync"
	"time"
)

func TryConnect(ip string, port int, timeout time.Duration) bool {
	addr := net.JoinHostPort(ip, strconv.Itoa(port))
	c, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return false
	}
	_ = c.Close()
	return true
}

func ScanPorts(ip string, ports []int, timeout time.Duration, workers int) []int {
	open := make([]int, 0)
	portCh := make(chan int, len(ports))
	resCh := make(chan int, len(ports))
	var wg sync.WaitGroup

	go func() {
		for _, p := range ports {
			portCh <- p
		}
		close(portCh)
	}()

	if workers <= 0 {
		workers = 50
	}
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			for p := range portCh {
				if TryConnect(ip, p, timeout) {
					resCh <- p
				} else {
					resCh <- 0
				}
			}
			wg.Done()
		}()
	}

	go func() { wg.Wait(); close(resCh) }()

	for r := range resCh {
		if r != 0 {
			open = append(open, r)
		}
	}
	sort.Ints(open)
	return open
}
