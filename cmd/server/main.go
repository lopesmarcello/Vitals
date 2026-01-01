package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lopesmarcello/vitals/internal/handlers"
	"github.com/lopesmarcello/vitals/views"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		views.Home().Render(r.Context(), w)
	})

	// /check?url=https://google.com
	r.Post("/check", handlers.AnalyzeURL)

	fmt.Println("Starting server on port :8080")
	http.ListenAndServe(":8080", r)
}
