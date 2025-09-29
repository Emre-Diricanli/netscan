package api

import (
	"encoding/json"
	"net/http"

	"netscan-server/internal/core"
)

type HTTPAPI struct {
	scanner *core.Scanner
}

func NewHTTPAPI(scanner *core.Scanner) *HTTPAPI {
	return &HTTPAPI{scanner: scanner}
}

func (h *HTTPAPI) Routes(mux *http.ServeMux) {
	mux.HandleFunc("/scan", h.handleScan)
	// health
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
}

func (h *HTTPAPI) handleScan(w http.ResponseWriter, r *http.Request) {
	h.allowCORS(&w, r)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req core.ScanRequest
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.Target == "" {
		http.Error(w, "target required", http.StatusBadRequest)
		return
	}

	results, err := h.scanner.Run(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(results)
}

func (h *HTTPAPI) allowCORS(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
