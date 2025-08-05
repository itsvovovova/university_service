package api

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"strconv"
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
    subjects

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
	// TODO: Тут тоже нужно будет вывод поправить, как только придумаю, как это будет в бд выглядеть
	for object := range schedule {
		resultText.WriteString(fmt.Sprintf("%s: %s", object, schedule[object]))
	}
	service.SendMessageWithRetries(ctx, b, resultText.String(), update.Message.Chat.ID, config.MaxRetries)
}

func ListDeadlines(ctx context.Context, b *bot.Bot, update *models.Update) {
	deadlines, err := db.GetUserSchedule(update.Message.Chat.ID)
	if err != nil {
		service.SendMessageWithRetries(ctx, b, config.ServerErrorText, update.Message.Chat.ID, config.MaxRetries)
		fmt.Println(config.GetDeadlinesErrorText)
	}
	var resultText strings.Builder
	// TODO: нужно поправить дизайн ответа: очевидно, что пользователю нужно выдавать дедлайны по-другому
	for object := range deadlines {
		resultText.WriteString(fmt.Sprintf("%s: %s", object, deadlines[object]))
	}
	service.SendMessageWithRetries(ctx, b, resultText.String(), update.Message.Chat.ID, config.MaxRetries)
}

func ListActiveUsers(ctx context.Context, b *bot.Bot, update *models.Update) {
	activeUsers, err := db.GetActiveUsers()
	if err != nil {
		service.SendMessageWithRetries(ctx, b, config.ServerErrorText, update.Message.Chat.ID, config.MaxRetries)
		fmt.Println(config.GetActiveUsersErrorText)
	}
	service.SendMessageWithRetries(
		ctx, b, "Количество пользователей: "+strconv.Itoa(activeUsers), update.Message.Chat.ID, config.MaxRetries)
}

func ListScores(ctx context.Context, b *bot.Bot, update *models.Update) {
	scores, err := db.GetUserSchedule(update.Message.Chat.ID)
	if err != nil {
		service.SendMessageWithRetries(ctx, b, config.ServerErrorText, update.Message.Chat.ID, config.MaxRetries)
		fmt.Println(config.GetScheduleErrorText)
	}
	var resultText strings.Builder
	for object := range scores {
		resultText.WriteString(fmt.Sprintf("%s: %s баллов\n", object, scores[object]))
	}
	service.SendMessageWithRetries(ctx, b, resultText.String(), update.Message.Chat.ID, config.MaxRetries)
}
