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

func LoginLk(ctx context.Context, b *bot.Bot, update *models.Update) {
	url := "https://lk.gubkin.ru/new/api/api.php?module=auth&method=login"
	request := types.LoginRequest{
		Login:      login,
		Password:   password,
		RememberMe: "1",
	}
	data, err := json.Marshal(request)
	if err != nil {
		return types.LoginResponse{}, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	if err != nil {
		return types.LoginResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	client := http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		return types.LoginResponse{}, err
	}
	if response.StatusCode != 200 {
		return types.LoginResponse{}, fmt.Errorf("login failed: %s", response.Status)
	}
	cookie := response.Cookies()
	var cookies []http.Cookie
	for _, cookie := range cookie {
		cookies = append(cookies, *cookie)
	}
	rememberMeToken, err := getTokenFromCookie(cookies, "rememberMe")
	if err != nil {
		return types.LoginResponse{}, err
	}
	phpSessionToken, err := getTokenFromCookie(cookies, "php_session")
	if err != nil {
		return types.LoginResponse{}, err
	}
	return types.LoginResponse{
		RememberMe: rememberMeToken,
		PhpSession: phpSessionToken,
	}, nil
}

func LoginEdu(login int, password string) error {
	// here we are trying to access the website
	return nil
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
