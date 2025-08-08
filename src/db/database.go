package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
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

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("не удалось пингануть базу данных: %w", err)
	}

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

	if err := createScoresTable(); err != nil {
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
			login BIGINT NOT NULL REFERENCES users(chat_id),
			schedule_date DATE NOT NULL,
			day_of_week SMALLINT CHECK (day_of_week BETWEEN 1 AND 7),
			week_number SMALLINT,
			pair_number SMALLINT NOT NULL CHECK (pair_number BETWEEN 1 AND 7),
			subject VARCHAR(255) NOT NULL,
			start_time TIME,
			end_time TIME,
			auditory VARCHAR(5),
			teacher VARCHAR(100),
			UNIQUE(login, schedule_date, pair_number)
		);
	`)
	if err != nil {
		return fmt.Errorf("не удалось создать таблицу schedule: %w", err)
	}

	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_schedule_user ON schedule(login)",
		"CREATE INDEX IF NOT EXISTS idx_schedule_date ON schedule(schedule_date)",
		"CREATE INDEX IF NOT EXISTS idx_schedule_user_week ON schedule(login, week_number)",
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
		login BIGINT NOT NULL REFERENCES users(chat_id),
		task_name VARCHAR(255) NOT NULL,
		deadline_date DATE NOT NULL,
		subject VARCHAR(255),
		is_completed BOOLEAN DEFAULT FALSE
	);
`)
	if err != nil {
		return fmt.Errorf("не удалось создать таблицу deadlines: %w", err)
	}

	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_deadlines_user ON deadlines(login)",
		"CREATE INDEX IF NOT EXISTS idx_deadlines_date ON deadlines(deadline_date)",
		"CREATE INDEX IF NOT EXISTS idx_deadlines_user_incomplete ON deadlines(login) WHERE is_completed = FALSE",
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
			type_notification IN ('update_schedule', 'update_scores', 'update_deadlines')),
		title VARCHAR(255),
		message TEXT NOT NULL,
		is_sent BOOLEAN DEFAULT FALSE
	);
`)
	if err != nil {
		return fmt.Errorf("не удалось создать таблицу notification: %w", err)
	}

	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_notifications_user ON notifications(login)",
		"CREATE INDEX IF NOT EXISTS idx_notifications_unsent ON notifications(login) WHERE is_sent = FALSE",
	}

	for _, query := range indexes {
		if _, err := DB.Exec(query); err != nil {
			return fmt.Errorf("не удалось создать индекс: %w", err)
		}
	}

	return nil
}

func createScoresTable() error {
	_, err := DB.Exec(`
        CREATE TABLE IF NOT EXISTS scores (
            id SERIAL PRIMARY KEY,
            login VARCHAR(50) NOT NULL REFERENCES users(login),
            chat_id BIGINT NOT NULL REFERENCES users(chat_id),
            subject VARCHAR(255) NOT NULL,
            score SMALLINT NOT NULL,
            max_score SMALLINT NOT NULL,
            UNIQUE(login, subject)  // Уникальность по логину и предмету
        );
    `)
	if err != nil {
		return fmt.Errorf("не удалось создать таблицу scores: %w", err)
	}

	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_scores_login ON scores(login)",
		"CREATE INDEX IF NOT EXISTS idx_scores_chat_id ON scores(chat_id)",
		"CREATE INDEX IF NOT EXISTS idx_scores_subject ON scores(subject)",
	}

	for _, query := range indexes {
		if _, err := DB.Exec(query); err != nil {
			return fmt.Errorf("не удалось создать индекс: %w", err)
		}
	}

	return nil
}
