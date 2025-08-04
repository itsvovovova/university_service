package auth

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"university_bot/config"
	"university_bot/src/api"
	"university_bot/src/db"
	"university_bot/src/service"
)

func Middleware(ctx context.Context, b *bot.Bot, update *models.Update, next bot.HandlerFunc) {
	currentState, err := db.GetUserState()
	if err != nil {
		service.SendMessageWithRetries(ctx, b, config.ServerErrorText, update.Message.Chat.ID, 100)
		return
	}
	switch currentState {
	case "login":
		api.Login(ctx, b, update)
	case "passwordLk":
		api.PasswordLk(ctx, b, update)
		service.LoginLk(ctx, b, update)
	case "passwordEdu":
		api.PasswordEdu(ctx, b, update)
	default:
		next(ctx, b, update)
	}
}
