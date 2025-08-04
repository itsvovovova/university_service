package api

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
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

func Schedule(ctx context.Context, b *bot.Bot, update *models.Update) {}

func ListDeadlines(ctx context.Context, b *bot.Bot, update *models.Update) {

}

func ListActiveUsers(ctx context.Context, b *bot.Bot, update *models.Update) {

}

func ListScores(ctx context.Context, b *bot.Bot, update *models.Update) {

}
