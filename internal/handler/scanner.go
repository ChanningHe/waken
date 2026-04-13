package handler

import (
	"log"
	"net/http"

	"github.com/channinghe/waken/internal/scanner"
)

func Scan(w http.ResponseWriter, r *http.Request) {
	hosts, err := scanner.Scan()
	if err != nil {
		log.Printf("scan failed: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to scan network")
		return
	}
	if hosts == nil {
		hosts = []scanner.Host{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"hosts": hosts})
}
