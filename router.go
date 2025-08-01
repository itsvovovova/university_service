package api

import (
	"github.com/go-chi/chi/v5"
)

func New() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/schedule", get_schedule)
	r.Get("/deadlines", get_deadlines)
	r.Get("/users", get_active_users)
}
