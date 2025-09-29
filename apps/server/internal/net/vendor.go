package net

// package netutil

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

// ---------- Interfaces & types ----------

type VendorDB interface {
	Lookup(oui string) string
}

type MemoryOUI map[string]string

func (m MemoryOUI) Lookup(oui string) string {
	key := normalizeOUI(oui)
	if key == "" {
		return ""
	}
	// Look up first 6 hex (24-bit OUI)
	if v, ok := m[key[:6]]; ok {
		return v
	}
	return ""
}

// Default tiny seed (useful if CSV missing)
func DefaultOUI() MemoryOUI {
	return MemoryOUI{
		"BC2411": "Proxmox Server Solutions GmbH",
		"A08CFD": "Hewlett Packard",
		"B06CB1": "Google LLC",
	}
}

// ---------- Loader ----------

// LoadOUIFromCSV reads the official IEEE CSV and builds a MemoryOUI map.
//
// Get file from: https://standards-oui.ieee.org/oui/oui.csv
//
// CSV header usually has columns like:
//
//	"Registry","Assignment","Organization Name","Organization Address"
//
// Example Assignment formats:
//
//	"FC-FB-FB" or "FC:FB:FB" or "FCF BFB" (variations). We normalize to hex.
func LoadOUIFromCSV(path string) (MemoryOUI, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open OUI csv: %w", err)
	}
	defer f.Close()

	cr := csv.NewReader(f)
	cr.FieldsPerRecord = -1 // organization address may contain commas

	header, err := cr.Read()
	if err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}
	colAssign := indexOf(header, "Assignment")
	colOrg := indexOf(header, "Organization Name")
	if colAssign < 0 || colOrg < 0 {
		return nil, fmt.Errorf("OUI csv missing required columns (Assignment/Organization Name)")
	}

	db := MemoryOUI{}
	for {
		rec, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			// Skip malformed lines rather than failing the whole load
			continue
		}
		if colAssign >= len(rec) || colOrg >= len(rec) {
			continue
		}
		assign := rec[colAssign]
		org := strings.TrimSpace(rec[colOrg])
		key := normalizeOUI(assign)
		if key == "" {
			continue
		}
		// store only first 6 hex characters (24-bit OUI)
		db[key[:6]] = org
	}
	return db, nil
}

func indexOf(arr []string, name string) int {
	for i, v := range arr {
		if strings.EqualFold(strings.TrimSpace(v), name) {
			return i
		}
	}
	return -1
}

func normalizeOUI(s string) string {
	if s == "" {
		return ""
	}
	// remove separators and spaces, uppercase
	up := strings.ToUpper(s)
	up = strings.ReplaceAll(up, ":", "")
	up = strings.ReplaceAll(up, "-", "")
	up = strings.ReplaceAll(up, ".", "")
	up = strings.ReplaceAll(up, " ", "")
	// Some CSVs might include extra text; keep hex only
	var hexOnly strings.Builder
	for _, r := range up {
		if (r >= '0' && r <= '9') || (r >= 'A' && r <= 'F') {
			hexOnly.WriteRune(r)
		}
	}
	h := hexOnly.String()
	if len(h) < 6 {
		return ""
	}
	return h
}

// VendorFromMAC: convenience for looking up a full MAC like "aa:bb:cc:dd:ee:ff".
func VendorFromMAC(db VendorDB, mac string) string {
	if mac == "" || db == nil {
		return ""
	}
	// reuse normalize, it will uppercase and strip separators
	n := normalizeOUI(mac)
	if len(n) < 6 {
		return ""
	}
	return db.Lookup(n[:6])
}
