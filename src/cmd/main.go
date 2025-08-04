package main

import (
	"context"
	"github.com/go-telegram/bot"
	"log"
	"net/http"
	"os"
	"os/signal"
	"university_bot/src/api"
	"university_bot/src/auth"
)

func main() {
	if err := http.ListenAndServe("PORT", router); err != nil {
		log.Fatalf("Server not started, error: %v", err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	CurrentBot, err := bot.New("YOUR_BOT_TOKEN_FROM_BOTFATHER")
	if err != nil {
		panic(err)
	}
	CurrentBot.Start(ctx)
}
