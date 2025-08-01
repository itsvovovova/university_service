package api

import "net/http"

/*
	r.Get("/schedule", get_schedule)
	r.Get("/deadlines", get_deadlines)
	r.Get("/users", get_active_users)
	r.Get("/score", get_score)
*/

func get_schedule(w http.ResponseWriter, r *http.Request) map[string]string {
	return map[string]string{}
}

func get_deadlines(w http.ResponseWriter, r *http.Request) map[string]string {
	return map[string]string{}
}

func get_active_users(w http.ResponseWriter, r *http.Request) map[string]string {
	return map[string]string{}
}

func get_scores(w http.ResponseWriter, r *http.Request) map[string]string {
	return map[string]string{}
}
