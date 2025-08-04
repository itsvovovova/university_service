package auth

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"net/http"
	"university_bot/src/api"
	"university_bot/src/db"
	"university_bot/src/service"
)

func Middleware(ctx context.Context, b *bot.Bot, update *models.Update, next bot.HandlerFunc) {
		currentState, err := db.GetUserState()
		if err != nil {
			service.SendMessageWithRetries(ctx, b, "")
			return
		}
		switch currentState {
		case "login":
			// send_message
			if err := api.Login(w, r); err != nil {
				http.Error(w, http.StatusText(500), 500)
			} else {
				if err := db.UpdateUserState("passwordLk"); err != nil {
					http.Error(w, http.StatusText(500), 500)
				}
			}
			return
		case "passwordLk":
			// send_message
			if err := api.PasswordLk(w, r); err != nil {
				http.Error(w, http.StatusText(500), 500)
			} else if err := db.UpdateUserState("passwordEdu"); err != nil {
				http.Error(w, http.StatusText(500), 500)
			}

			if err := service.LoginLk(w, r); err != nil {
				http.Error(w, http.StatusText(401), 401)
				if err := db.UpdateUserState("login"); err != nil {
					http.Error(w, http.StatusText(500), 500)
				}
			}
			return
		case "passwordEdu":
			// send_message
			if err := api.PasswordEdu(w, r); err != nil {
				http.Error(w, http.StatusText(500), 500)
			} else if err := db.UpdateUserState("success"); err != nil {
				http.Error(w, http.StatusText(500), 500)
			}
			if err := service.LoginEdu(); err != nil {
				http.Error(w, http.StatusText(401), 401)
			}
			return
		default:
			next(ctx, b, update)
}
