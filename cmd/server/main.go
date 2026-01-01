package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lopesmarcello/vitals/internal/handlers"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// /check?url=https://google.com
	r.Get("/check", handlers.AnalyzeURL)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Vitals"))
	})

	fmt.Println("Starting server on port :8080")
	http.ListenAndServe(":8080", r)
}
