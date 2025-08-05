package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"io"
	"log"
	"net/http"
	"time"
	"university_bot/config"
	"university_bot/src/db"
	"university_bot/src/types"
)

func ReadText(w http.ResponseWriter, r *http.Request) string {
	text, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Couldn't open the text", http.StatusBadRequest)
		return ""
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("Error when closing a file: %v", err)
		}
	}()
	return string(text)
}

func getTokenFromCookie(cookies []http.Cookie, token string) (string, error) {
	for _, cookie := range cookies {
		if token == cookie.Name {
			return cookie.Value, nil
		}
	}
	return "", errors.New("cookie not found")
}

func getUserInfo(ctx context.Context, b *bot.Bot, update *models.Update) (string, string, error) {
	chatID := update.Message.Chat.ID
	login, err := db.GetUserLogin(chatID)
	if err != nil {
		SendMessageWithRetries(ctx, b, config.ServerErrorText, chatID, config.MaxRetries)
		fmt.Println(config.GetUserLoginErrorText)
		return "", "", err
	}
	password, err := db.GetUserPasswordLk(chatID)
	if err != nil {
		SendMessageWithRetries(ctx, b, config.ServerErrorText, chatID, config.MaxRetries)
		fmt.Println(config.GetUserPasswordLkErrorText)
		return "", "", err
	}
	return login, password, nil
}

func sendPostRequest(ctx context.Context, b *bot.Bot, update *models.Update, request types.LoginRequest) (*http.Response, error) {
	chatID := update.Message.Chat.ID
	url := "https://lk.gubkin.ru/new/api/api.php?module=auth&method=login"
	data, err := json.Marshal(request)
	if err != nil {
		SendMessageWithRetries(ctx, b, config.ServerErrorText, chatID, config.MaxRetries)
		fmt.Println(config.PostRequestConvertErrorText)
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	if err != nil {
		SendMessageWithRetries(ctx, b, config.ServerErrorText, chatID, config.MaxRetries)
		fmt.Print(config.PostRequestConvertErrorText)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	client := http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		SendMessageWithRetries(ctx, b, config.ServerErrorText, chatID, config.MaxRetries)
		fmt.Println(config.PostRequestServerErrorText)
		return nil, err
	}
	if response.StatusCode != 200 {
		SendMessageWithRetries(ctx, b, config.ServerErrorText, chatID, config.MaxRetries)
		fmt.Println(config.StatusRequestErrorText)
		return nil, err
	}
	return response, nil
}

func getUserTokens(ctx context.Context, b *bot.Bot, update *models.Update, response http.Response) (string, string, error) {
	cookie, chatID := response.Cookies(), update.Message.Chat.ID
	var cookies []http.Cookie
	for _, cookie := range cookie {
		cookies = append(cookies, *cookie)
	}
	rememberMeToken, err := getTokenFromCookie(cookies, "rememberMe")
	if err != nil {
		SendMessageWithRetries(ctx, b, config.ServerErrorText, chatID, config.MaxRetries)
		fmt.Println(config.GetRememberMeErrorText)
		return "", "", err
	}
	phpSessionToken, err := getTokenFromCookie(cookies, "php_session")
	if err != nil {
		SendMessageWithRetries(ctx, b, config.ServerErrorText, chatID, config.MaxRetries)
		fmt.Print(config.GetPhpSessionErrorText)
		return "", "", err
	}
	return rememberMeToken, phpSessionToken, nil
}

func LoginLk(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID
	login, password, err := getUserInfo(ctx, b, update)
	if err != nil {
		return
	}
	request := types.LoginRequest{
		Login:      login,
		Password:   password,
		RememberMe: "1",
	}
	response, err := sendPostRequest(ctx, b, update, request)
	if err != nil {
		return
	}
	rememberMe, phpSession, err := getUserTokens(ctx, b, update, *response)
	if err != nil {
		return
	}
	if err := db.AddRememberMeToken(chatID, rememberMe); err != nil {
		SendMessageWithRetries(ctx, b, config.ServerErrorText, chatID, config.MaxRetries)
		fmt.Println(config.AddRememberMeDbErrorText)
		return
	}

	if err := db.AddPhpSessionToken(chatID, phpSession); err != nil {
		SendMessageWithRetries(ctx, b, config.ServerErrorText, chatID, config.MaxRetries)
		fmt.Println(config.AddPhpSessionDbErrorText)
		return
	}
}

func LoginEdu(login int, password string) error {
	// here we are trying to access the website
	return nil // TODO: доделать..........
}
func SendMessageWithRetries(ctx context.Context, b *bot.Bot, message string, chatID int64, retries int) {
	for i := 0; i < retries; i++ {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   message,
		}); err != nil {
			log.Printf("Error sending message for chatID: %v", err)
		} else {
			return
		}
	}
}
