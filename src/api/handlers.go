package api

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"strings"
	"university_bot/config"
	"university_bot/src/db"
	"university_bot/src/service"
)

/*
users table
    chat_id SERIAL PRIMARY KEY,
    state = any([login, passwordLk, passwordEdu, success]) NOT NULL
    login NULL
    password_edu NULL
    password_lk NULL
    last_update NULL

schedule table
    login SERIAL PRIMARY KEY,
    day_of_week
    start_time
    end_time
    subject

deadlines table (
    login SERIAL PRIMARY KEY,
    task_name
    due_date
    subject
    is_completed
    notification_sent bool
);
*/

func Schedule(ctx context.Context, b *bot.Bot, update *models.Update) {
	schedule, err := db.GetUserSchedule(update.Message.Chat.ID)
	if err != nil {
		service.SendMessageWithRetries(ctx, b, config.ServerErrorText, update.Message.Chat.ID, config.MaxRetries)
		fmt.Println(config.GetScheduleErrorText)
	}
	var resultText strings.Builder
	for object := range schedule {
		resultText.WriteString(fmt.Sprintf("%s: %s баллов\n", object, schedule[object]))
	}
	service.SendMessageWithRetries(ctx, b, resultText.String(), update.Message.Chat.ID, config.MaxRetries)
}

func ListDeadlines(ctx context.Context, b *bot.Bot, update *models.Update) {

}

func ListActiveUsers(ctx context.Context, b *bot.Bot, update *models.Update) {

}

func ListScores(ctx context.Context, b *bot.Bot, update *models.Update) {

}
