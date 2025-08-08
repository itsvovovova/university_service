package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"university_bot/src/types"
)

func AddLoginUser(chatID int64, login string) error {
	query := `
		INSERT INTO users (chat_id, login, state, last_update) 
		VALUES ($1, $2, 'login', CURRENT_TIMESTAMP)
		ON CONFLICT (chat_id) DO UPDATE 
		SET login = EXCLUDED.login, state = 'login', last_update = CURRENT_TIMESTAMP;
	`
	_, err := DB.Exec(query, chatID, login)
	return err
}

func AddPasswordLkUser(chatID int64, password string) error {
	query := `
		UPDATE users 
		SET password_lk = $1, state = 'passwordLk', last_update = CURRENT_TIMESTAMP
		WHERE chat_id = $2;
	`
	_, err := DB.Exec(query, password, chatID)
	return err
}

func AddPasswordEduUser(chatID int64, password string) error {
	query := `
		UPDATE users 
		SET password_edu = $1, state = 'passwordEdu', last_update = CURRENT_TIMESTAMP
		WHERE chat_id = $2;
	`
	_, err := DB.Exec(query, password, chatID)
	return err
}

func GetActiveUsers() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE last_update > CURRENT_DATE - INTERVAL '30 days';`
	err := DB.QueryRow(query).Scan(&count)
	return count, err
}

func AddUserSchedule(schedules []types.Schedule) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			fmt.Println(err)
		}
	}(tx)

	stmt, err := tx.Prepare(`
		INSERT INTO schedule (
			login, schedule_date, day_of_week, week_number, 
			pair_number, subject, start_time, end_time, auditory, teacher
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (login, schedule_date, pair_number) DO UPDATE
		SET subject = EXCLUDED.subject,
			start_time = EXCLUDED.start_time,
			end_time = EXCLUDED.end_time,
			auditory = EXCLUDED.auditory,
			teacher = EXCLUDED.teacher;
	`)
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(stmt)

	for _, s := range schedules {
		_, err := stmt.Exec(
			s.Login,
			s.ScheduleDate,
			s.DayOfWeek,
			s.WeekNumber,
			s.PairNumber,
			s.Subject,
			s.StartTime,
			s.EndTime,
			s.Auditory,
			s.Teacher,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func AddUserDeadlines(deadlines []types.Deadline) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			fmt.Println(err)
		}
	}(tx)

	stmt, err := tx.Prepare(`
	INSERT INTO deadlines (
		login, task_name, deadline_date, subject
	) VALUES ($1, $2, $3, $4)
	ON CONFLICT (login, task_name, deadline_date) DO UPDATE
	SET subject = EXCLUDED.subject;
`)
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(stmt)

	for _, d := range deadlines {
		_, err := stmt.Exec(
			d.Login,
			d.TaskName,
			d.DeadlineDate,
			d.Subject,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func GetUserSchedule(login string, startDate, endDate time.Time) ([]types.Schedule, error) {
	query := `
		SELECT 
			schedule_date, day_of_week, week_number, pair_number,
			subject, start_time, end_time, auditory, teacher
		FROM schedule 
		WHERE login = $1 
		AND schedule_date BETWEEN $2 AND $3
		ORDER BY schedule_date, pair_number;
	`

	rows, err := DB.Query(query, login, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(rows)

	var schedules []types.Schedule
	for rows.Next() {
		var s types.Schedule
		err := rows.Scan(
			&s.ScheduleDate,
			&s.DayOfWeek,
			&s.WeekNumber,
			&s.PairNumber,
			&s.Subject,
			&s.StartTime,
			&s.EndTime,
			&s.Auditory,
			&s.Teacher,
		)
		if err != nil {
			return nil, err
		}
		s.Login = login
		schedules = append(schedules, s)
	}

	return schedules, nil
}

func GetUpcomingDeadlines(login string, limit int) ([]types.Deadline, error) {
	query := `
	SELECT 
		id, task_name, deadline_date, subject, is_completed
	FROM deadlines
	WHERE login = $1 
	AND deadline_date >= CURRENT_DATE
	AND is_completed = FALSE
	ORDER BY deadline_date
	LIMIT $2;
`

	rows, err := DB.Query(query, login, limit)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(rows)

	var deadlines []types.Deadline
	for rows.Next() {
		var d types.Deadline
		err := rows.Scan(
			&d.ID,
			&d.TaskName,
			&d.DeadlineDate,
			&d.Subject,
			&d.IsCompleted,
		)
		if err != nil {
			return nil, err
		}
		d.Login = login
		deadlines = append(deadlines, d)
	}

	return deadlines, nil
}

func UpdateUserState(chatID int64, state string) error {
	query := `
		UPDATE users 
		SET state = $1, last_update = CURRENT_TIMESTAMP
		WHERE chat_id = $2;
	`
	_, err := DB.Exec(query, state, chatID)
	return err
}

func GetUserState(chatID int64) (string, error) {
	var state string
	query := `SELECT state FROM users WHERE chat_id = $1;`
	err := DB.QueryRow(query, chatID).Scan(&state)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	return state, err
}

func GetUserLogin(chatID int64) (string, error) {
	var login string
	query := `SELECT login FROM users WHERE chat_id = $1;`
	err := DB.QueryRow(query, chatID).Scan(&login)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	return login, err
}

func GetUserPasswordLk(chatID int64) (string, error) {
	var password string
	query := `SELECT password_lk FROM users WHERE chat_id = $1;`
	err := DB.QueryRow(query, chatID).Scan(&password)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	return password, err
}

func AddRememberMeToken(chatID int64, token string) error {
	query := `
		UPDATE users 
		SET remember_me_token = $1, last_update = CURRENT_TIMESTAMP
		WHERE chat_id = $2;
	`
	_, err := DB.Exec(query, token, chatID)
	return err
}

func AddPhpSessionToken(chatID int64, token string) error {
	query := `
		UPDATE users 
		SET php_session_token = $1, last_update = CURRENT_TIMESTAMP
		WHERE chat_id = $2;
	`
	_, err := DB.Exec(query, token, chatID)
	return err
}

func GetAllUsers() ([]types.User, error) {
	query := `SELECT chat_id, login FROM users;`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var users []types.User
	for rows.Next() {
		var u types.User
		err := rows.Scan(&u.ChatID, &u.Login)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func ExistUser(chatID int64) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE chat_id = $1);`
	err := DB.QueryRow(query, chatID).Scan(&exists)
	return exists, err
}

func AddNotification(chatID int64, notificationType, title, message string) error {
	query := `
		INSERT INTO notifications (
			user_id, type_notification, title, message
		) VALUES ($1, $2, $3, $4)
	`
	_, err := DB.Exec(query, chatID, notificationType, title, message)
	return err
}

func GetUnreadNotifications(chatID int64) ([]types.Notification, error) {
	query := `
	SELECT id, type_notification, title, message
	FROM notifications
	WHERE user_id = $1 
	AND read_at IS NULL
	ORDER BY created_at DESC;
`

	rows, err := DB.Query(query, chatID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(rows)

	var notifications []types.Notification
	for rows.Next() {
		var n types.Notification
		err := rows.Scan(
			&n.ID,
			&n.Type,
			&n.Title,
			&n.Message,
		)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}

	return notifications, nil
}

func MarkNotificationAsRead(notificationID int) error {
	query := `
	UPDATE notifications 
	SET is_sent = TRUE
	WHERE id = $1;
`
	_, err := DB.Exec(query, notificationID)
	return err
}

func GetUnsentNotifications() ([]types.Notification, error) {
	query := `
		SELECT 
			n.id, u.chat_id, n.type_notification, 
			n.title, n.message, n.created_at
		FROM notifications n
		JOIN users u ON n.user_id = u.chat_id
		WHERE n.is_sent = FALSE
		ORDER BY n.created_at;
	`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var notifications []types.Notification
	for rows.Next() {
		var n types.Notification
		err := rows.Scan(
			&n.ID,
			&n.UserID,
			&n.Type,
			&n.Title,
			&n.Message,
		)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}

	return notifications, nil
}

func MarkNotificationAsSent(notificationID int) error {
	query := `
		UPDATE notifications 
		SET is_sent = TRUE, sent_at = CURRENT_TIMESTAMP
		WHERE id = $1;
	`
	_, err := DB.Exec(query, notificationID)
	return err
}

func BroadcastNotification(notificationType, title, message string) error {
	query := `
		INSERT INTO notifications (user_id, type_notification, title, message)
		SELECT chat_id, $1, $2, $3 FROM users;
	`
	_, err := DB.Exec(query, notificationType, title, message)
	return err
}

func UpdateUserScores(scores []types.Score) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			fmt.Println(err)
		}
	}(tx)

	stmt, err := tx.Prepare(`
        INSERT INTO scores (
            login, chat_id, subject, score, max_score
        ) VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (login, subject) DO UPDATE
        SET score = EXCLUDED.score,
            max_score = EXCLUDED.max_score;
    `)
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(stmt)

	for _, s := range scores {
		_, err := stmt.Exec(
			s.Login,
			s.ChatID,
			s.Subject,
			s.Score,
			s.MaxScore,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func GetUserScores(login string) ([]types.Score, error) {
	query := `
        SELECT subject, score, max_score
        FROM scores
        WHERE login = $1;
    `

	rows, err := DB.Query(query, login)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(rows)

	var scores []types.Score
	for rows.Next() {
		var s types.Score
		err := rows.Scan(
			&s.Subject,
			&s.Score,
			&s.MaxScore,
		)
		if err != nil {
			return nil, err
		}
		s.Login = login
		scores = append(scores, s)
	}

	return scores, nil
}

func GetUserScoresByChatID(chatID int64) ([]types.Score, error) {
	query := `
        SELECT subject, score, max_score
        FROM scores
        WHERE chat_id = $1;
    `

	rows, err := DB.Query(query, chatID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(rows)

	var scores []types.Score
	for rows.Next() {
		var s types.Score
		err := rows.Scan(
			&s.Subject,
			&s.Score,
			&s.MaxScore,
		)
		if err != nil {
			return nil, err
		}
		s.ChatID = chatID
		scores = append(scores, s)
	}

	return scores, nil
}

func GetRecentScores(limit int) ([]types.Score, error) {
	query := `
        SELECT login, chat_id, subject, score, max_score
        FROM scores
        ORDER BY id DESC
        LIMIT $1;
    `

	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(rows)

	var scores []types.Score
	for rows.Next() {
		var s types.Score
		err := rows.Scan(
			&s.Login,
			&s.ChatID,
			&s.Subject,
			&s.Score,
			&s.MaxScore,
		)
		if err != nil {
			return nil, err
		}
		scores = append(scores, s)
	}

	return scores, nil
}
