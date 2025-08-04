package api

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"math"
	"strconv"
	"university_bot/src/db"
	"university_bot/src/service"
)

func Login(ctx context.Context, b *bot.Bot, update *models.Update) {
	number, err := strconv.Atoi(update.Message.Text)
	if err != nil || float64(number) < math.Pow10(5) || float64(number) > math.Pow10(6)-1 {
		service.SendMessageWithRetries(
			ctx, b, "Невалидный логин, попробуй еще", update.Message.Text, 100)
		return
	}
	if err = db.AddLoginUser(number); err != nil {
		service.SendMessageWithRetries(
			ctx, b, "Проблемы с сервером, приносим извинения, попробуйте позже", update.Message.Text, 100)
		fmt.Println("Возникли трудности с добавлением логина пользователя в базу данных")
		return
	}
}

func PasswordLk(ctx context.Context, b *bot.Bot, update *models.Update) error {
	text := update.Message.Text
	if err := db.AddPasswordLkUser(text); err != nil {
		service.SendMessageWithRetries(
			ctx, b, "Возникли трудности с добавлением пароля пользователя из ЛК в базу данных", update.Message.Text, 100)
		return err
	}
	return nil
}

func PasswordEdu(ctx context.Context, b *bot.Bot, update *models.Update) error {
	text := update.Message.Text
	if err := db.AddPasswordEduUser(text); err != nil {
		service.SendMessageWithRetries(
			ctx, b, "Возникли трудности с добавлением пароля пользователя из edu в базу данных", update.Message.Text, 100)
		return err
	}
	return nil
}
