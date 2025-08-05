package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"university_bot/config"
)

var db *sql.DB

func InitDb() {
	var connectionStr = fmt.Sprintf(
		"port=%s host=%s user=%s password=%s dbname=%s sslmode=%s",
		config.CurrentConfig.Database.Port,
		config.CurrentConfig.Database.Host,
		config.CurrentConfig.Database.Username,
		config.CurrentConfig.Database.Password,
		config.CurrentConfig.Database.DatabaseName,
		config.CurrentConfig.Database.SSLMode)

	var err error
	db, err = sql.Open("postgres", connectionStr)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("Connected to database")

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS users (
        chat_id BIGINT PRIMARY KEY,
        state VARCHAR(20) NOT NULL CHECK (state IN ('login', 'passwordLk', 'passwordEdu', 'success')),
        login INT,
        password_edu VARCHAR(255),
        password_lk VARCHAR(255),
        last_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
`)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Create schedule table
	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS schedule (
        id SERIAL PRIMARY KEY,
        chat_id BIGINT REFERENCES users(chat_id) ON DELETE CASCADE,
        week_start DATE NOT NULL,
        day_of_week INTEGER NOT NULL CHECK (day_of_week BETWEEN 1 AND 7),
        pair_number INTEGER NOT NULL CHECK (pair_number BETWEEN 1 AND 7),
        subject VARCHAR(255) NOT NULL,
        UNIQUE(chat_id, week_start, day_of_week, pair_number)
    );
`)
	if err != nil {
		log.Fatal(err)
		return
	}

	_, err = db.Exec(`
    CREATE INDEX IF NOT EXISTS idx_schedule_user_week 
    ON schedule (chat_id, week_start);
`)
	if err != nil {
		log.Fatal(err)
		return
	}

	_, err = db.Exec(`
    CREATE INDEX IF NOT EXISTS idx_schedule_user_day 
    ON schedule (chat_id, week_start, day_of_week);
`)
	if err != nil {
		log.Fatal(err)
		return
	}

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS deadlines (
        id SERIAL PRIMARY KEY,
        chat_id BIGINT REFERENCES users(chat_id) ON DELETE CASCADE,
        task_name VARCHAR(255) NOT NULL,
        date_value DATE NOT NULL,
        subject VARCHAR(255),
        is_completed BOOLEAN DEFAULT FALSE,
        notification_sent BOOLEAN DEFAULT FALSE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
`)
	if err != nil {
		log.Fatal(err)
		return
	}

	_, err = db.Exec(`
    CREATE INDEX IF NOT EXISTS idx_deadlines_user_date 
    ON deadlines (chat_id, date_value);
`)
	if err != nil {
		log.Fatal(err)
		return
	}

	_, err = db.Exec(`
    CREATE INDEX IF NOT EXISTS idx_deadlines_incomplete 
    ON deadlines (chat_id, is_completed, date_value);
`)
	if err != nil {
		log.Fatal(err)
		return
	}
	return
}
