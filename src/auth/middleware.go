package auth

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"university_bot/src/db"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		r.Body = io.NopCloser(bytes.NewReader(data))
		var request struct {
			ChatId int `json:"chat_id"`
		}
		if err := json.Unmarshal(data, &request); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		flag, err := db.ExistUser(request.ChatId)
		if flag == false {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		}
		if err != nil {
			http.Error(w, http.StatusText(401), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
