package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/lopesmarcello/vitals/internal/analyzer"
)

func AnalyzeURL(w http.ResponseWriter, r *http.Request) {
	urlParam := r.URL.Query().Get("url")
	if urlParam == "" {
		http.Error(w, "Missing 'url' query parameter", http.StatusBadRequest)
		return
	}

	stats, err := analyzer.AnalyzeNetwork(urlParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
