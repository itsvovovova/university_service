package main

import (
	"context"
	"github.com/go-telegram/bot"
	"log"
	"os"
	"os/signal"
	"university_bot/src/consumer"
	"university_bot/src/db"
)

func main() {
	if err := db.InitDB(); err != nil {
		log.Fatal("Не удалось инициализировать базу данных: ", err)
	}
	defer func() {
		err := db.CloseDB()
		if err != nil {
			panic(err)
		}
	}()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	CurrentBot, err := bot.New("YOUR_BOT_TOKEN_FROM_BOTFATHER")
	if err != nil {
		panic(err)
	}
	CurrentBot.Start(ctx)
	dataChan := consumer.Parser()
	consumer.ProcessData(dataChan)
}
