package net

// package netutil

import (
	"fmt"
	"strconv"
	"strings"
)

func ParsePortSpec(spec string) ([]int, error) {
	s := strings.TrimSpace(spec)
	if s == "" {
		return nil, fmt.Errorf("empty port spec")
	}
	if s == "all" {
		ports := make([]int, 65535)
		for i := 1; i <= 65535; i++ {
			ports[i-1] = i
		}
		return ports, nil
	}
	if strings.Contains(s, "-") {
		r := strings.Split(s, "-")
		if len(r) != 2 {
			return nil, fmt.Errorf("invalid range %s", s)
		}
		a, err := strconv.Atoi(r[0])
		if err != nil {
			return nil, err
		}
		b, err := strconv.Atoi(r[1])
		if err != nil {
			return nil, err
		}
		if a < 1 {
			a = 1
		}
		if b > 65535 {
			b = 65535
		}
		if b < a {
			return nil, fmt.Errorf("invalid range %s", s)
		}
		out := make([]int, 0, b-a+1)
		for p := a; p <= b; p++ {
			out = append(out, p)
		}
		return out, nil
	}
	p, err := strconv.Atoi(s)
	if err != nil {
		return nil, err
	}
	return []int{p}, nil
}
