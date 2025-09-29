package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"netscan-server/internal/api"
	"netscan-server/internal/core"
	nu "netscan-server/internal/net"
)

func main() {
	var (
		ouiPath string
		addr    string
	)
	flag.StringVar(&ouiPath, "oui", "./oui.csv", "Path to IEEE OUI CSV (download from https://standards-oui.ieee.org/oui/oui.csv)")
	flag.StringVar(&addr, "addr", ":8080", "Listen address")
	flag.Parse()

	var db nu.VendorDB
	if _, err := os.Stat(ouiPath); err == nil {
		if loaded, err := nu.LoadOUIFromCSV(ouiPath); err == nil {
			fmt.Printf("Loaded OUI entries: %d from %s\n", len(loaded), ouiPath)
			db = loaded
		} else {
			fmt.Printf("WARN: failed to load OUI CSV: %v (falling back to default)\n", err)
			db = nu.DefaultOUI()
		}
	} else {
		fmt.Printf("OUI CSV not found at %s — using default minimal map.\n", ouiPath)
		db = nu.DefaultOUI()
	}

	sc := core.NewScanner(db)
	mux := http.NewServeMux()
	api.NewHTTPAPI(sc).Routes(mux)

	fmt.Printf("netscan server -> http://%s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		panic(err)
	}
}
