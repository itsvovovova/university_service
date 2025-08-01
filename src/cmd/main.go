package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	if err := http.ListenAndServe(":8000", router); err != nil {
		log.Fatalf("Server not started, error: %v", err)
	}
}
