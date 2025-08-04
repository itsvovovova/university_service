package api

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"math"
	"strconv"
	"university_bot/config"
	"university_bot/src/db"
	"university_bot/src/service"
)

func Login(ctx context.Context, b *bot.Bot, update *models.Update) {
	number, err := strconv.Atoi(update.Message.Text)
	if err != nil || float64(number) < math.Pow10(5) || float64(number) > math.Pow10(6)-1 {
		service.SendMessageWithRetries(
			ctx, b, config.IncorrectLoginErrorText, update.Message.Chat.ID, config.MaxRetries)
		return
	}
	if err = db.AddLoginUser(number); err != nil {
		service.SendMessageWithRetries(
			ctx, b, config.ServerErrorText, update.Message.Chat.ID, config.MaxRetries)
		fmt.Println(config.LoginDbErrorText)
		return
	}
}

func PasswordLk(ctx context.Context, b *bot.Bot, update *models.Update) {
	text := update.Message.Text
	if err := db.AddPasswordLkUser(text); err != nil {
		service.SendMessageWithRetries(
			ctx, b, config.ServerErrorText, update.Message.Chat.ID, config.MaxRetries)
	}
	if err := db.UpdateUserState("PasswordEdu"); err != nil {
		service.SendMessageWithRetries(
			ctx, b, config.ServerErrorText, update.Message.Chat.ID, config.MaxRetries)
	}
}

func PasswordEdu(ctx context.Context, b *bot.Bot, update *models.Update) {
	text := update.Message.Text
	if err := db.AddPasswordEduUser(text); err != nil {
		service.SendMessageWithRetries(
			ctx, b, config.ServerErrorText, update.Message.Chat.ID, config.MaxRetries)
		fmt.Println(config.PasswordDbErrorText)
	}
	if err := db.UpdateUserState("success"); err != nil {
		service.SendMessageWithRetries(
			ctx, b, config.ServerErrorText, update.Message.Chat.ID, config.MaxRetries)
	}
}
