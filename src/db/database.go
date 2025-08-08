package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
	"university_bot/config"
)

var DB *sql.DB

func InitDB() error {
	connectionStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.CurrentConfig.Database.Host,
		config.CurrentConfig.Database.Port,
		config.CurrentConfig.Database.Username,
		config.CurrentConfig.Database.Password,
		config.CurrentConfig.Database.DatabaseName,
		config.CurrentConfig.Database.SSLMode)

	var err error
	DB, err = sql.Open("postgres", connectionStr)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxLifetime(5 * time.Minute)

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("не удалось пингануть базу данных: %w", err)
	}

	log.Println("Successfully connected to database")

	if err := createUsersTable(); err != nil {
		return err
	}

	if err := createScheduleTable(); err != nil {
		return err
	}

	if err := createDeadlinesTable(); err != nil {
		return err
	}

	if err := createNotificationsTable(); err != nil {
		return err
	}

	return nil
}

func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

func createUsersTable() error {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			chat_id BIGINT PRIMARY KEY,
			state VARCHAR(20) NOT NULL CHECK (state IN ('login', 'passwordLk', 'passwordEdu', 'success')),
			login VARCHAR(50) UNIQUE,
			password_edu VARCHAR(255),
			password_lk VARCHAR(255),
			last_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return fmt.Errorf("не удалось создать таблицу users: %w", err)
	}
	return nil
}

func createScheduleTable() error {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS schedule (
			id SERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL REFERENCES users(chat_id),
			schedule_date DATE NOT NULL,
			day_of_week SMALLINT CHECK (day_of_week BETWEEN 1 AND 7),
			week_number SMALLINT,
			pair_number SMALLINT NOT NULL CHECK (pair_number BETWEEN 1 AND 7),
			subject VARCHAR(255) NOT NULL,
			start_time TIME,
			end_time TIME,
			auditory VARCHAR(50),
			teacher VARCHAR(100),
			UNIQUE(user_id, schedule_date, pair_number)
		);
	`)
	if err != nil {
		return fmt.Errorf("не удалось создать таблицу schedule: %w", err)
	}

	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_schedule_user ON schedule(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_schedule_date ON schedule(schedule_date)",
		"CREATE INDEX IF NOT EXISTS idx_schedule_user_week ON schedule(user_id, week_number)",
	}

	for _, query := range indexes {
		if _, err := DB.Exec(query); err != nil {
			return fmt.Errorf("не удалось создать индекс: %w", err)
		}
	}

	return nil
}

func createDeadlinesTable() error {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS deadlines (
			id SERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL REFERENCES users(chat_id),
			task_name VARCHAR(255) NOT NULL,
			deadline_date DATE NOT NULL,
			subject VARCHAR(255),
			description TEXT,
			priority SMALLINT DEFAULT 3 CHECK (priority BETWEEN 1 AND 5),
			is_completed BOOLEAN DEFAULT FALSE,
			notification_sent BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return fmt.Errorf("не удалось создать таблицу deadlines: %w", err)
	}

	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_deadlines_user ON deadlines(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_deadlines_date ON deadlines(deadline_date)",
		"CREATE INDEX IF NOT EXISTS idx_deadlines_user_incomplete ON deadlines(user_id) WHERE is_completed = FALSE",
	}

	for _, query := range indexes {
		if _, err := DB.Exec(query); err != nil {
			return fmt.Errorf("не удалось создать индекс в таблице deadlines: %w", err)
		}
	}

	return nil
}

func createNotificationsTable() error {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS notifications (
			id SERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL REFERENCES users(chat_id),
			type_notification VARCHAR(50) NOT NULL CHECK (
				type_notification IN ('update_schedule', 'update_scores', 'update_deadlines')
			),
			title VARCHAR(255),
			message TEXT NOT NULL,
			is_sent BOOLEAN DEFAULT FALSE,
			sent_at TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			read_at TIMESTAMP
		);
	`)
	if err != nil {
		return fmt.Errorf("не удалось создать таблицу notification: %w", err)
	}

	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_notifications_user ON notifications(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_notifications_unsent ON notifications(user_id) WHERE is_sent = FALSE",
	}

	for _, query := range indexes {
		if _, err := DB.Exec(query); err != nil {
			return fmt.Errorf("не удалось создать индекс: %w", err)
		}
	}

	return nil
}
