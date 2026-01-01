package handlers

import (
	"net/http"

	"github.com/lopesmarcello/vitals/internal/analyzer"
	"github.com/lopesmarcello/vitals/views"
)

func AnalyzeURL(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error to parse form", http.StatusBadRequest)
		return
	}

	urlParam := r.FormValue("url")

	if urlParam == "" {
		http.Error(w, "Missing 'url' query parameter", http.StatusBadRequest)
		return
	}

	stats, err := analyzer.Analyze(r.Context(), urlParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	views.Results(stats).Render(r.Context(), w)
}
