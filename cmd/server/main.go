package main

import (
	"fmt"
	"net/http"
	"os"

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
		if r.Header.Get("HX-Request") == "true" {
			views.HomeContent().Render(r.Context(), w)
		} else {
			views.Home().Render(r.Context(), w)
		}
	})

	// /check?url=https://google.com
	r.Post("/check", handlers.AnalyzeURL)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	fmt.Println("Starting server on port :" + port)
	http.ListenAndServe(":"+port, r)
}
