package main

import (
	"context"
	"github.com/go-telegram/bot"
	"os"
	"os/signal"
	"university_bot/src/db"
)

func main() {
	db.InitDb()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	CurrentBot, err := bot.New("YOUR_BOT_TOKEN_FROM_BOTFATHER")
	if err != nil {
		panic(err)
	}
	CurrentBot.Start(ctx)
}
