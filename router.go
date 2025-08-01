package university_bot

import (
	"github.com/go-chi/chi/v5"
	"university_bot/src/api"
)

func New() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/schedule", api.Schedule)
	r.Get("/deadlines", api.ListDeadlines)
	r.Get("/users", api.ListActiveUsers)
	r.Get("/score", api.ListScores)
	return r
}
