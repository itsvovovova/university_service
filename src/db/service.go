package db

func AddLoginUser(chatID int64, login int) error {
	var query = "INSERT INTO users(chat_id, login) VALUES ($1, $2)"
	if _, err := db.Exec(query, chatID, login); err != nil {
		return err
	}
	return nil
}

func AddPasswordLkUser(password string) error {
	var query = "INSERT INTO users(password_lk, last_update) VALUES ($1, CURRENT_DATE)"
	if _, err := db.Exec(query, password); err != nil {
		return err
	}
	return nil
}

func AddPasswordEduUser(password string) error {
	var query = "INSERT INTO users(password_edu, last_update) VALUES ($1, CURRENT_DATE)"
	if _, err := db.Exec(query, password); err != nil {
		return err
	}
	return nil
}

func GetActiveUsers() (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func AddUserSchedule(login int, schedule map[string][7]string) error {
	return nil
}

func AddUserDeadlines(login int) error {
	return nil
}

func AddUserScore() error {
	return nil
}

func GetUserSchedule(chatID int64) (map[string]string, error) {
	return map[string]string{}, nil
}

func GetUserDeadlines(login int) (map[string]string, error) {
	return map[string]string{}, nil
}

func GetUserScore() (map[string]string, error) {
	return map[string]string{}, nil
}

func UpdateUserSchedule(login int, new_schedule map[string]string) error {
	return nil
}

func UpdateUserDeadlines(login int, new_deadlines map[string]string) error {
	return nil
}

func UpdateActiveUsers(users map[string]string) error {
	return nil
}

func GetUsers(login int) (map[string]string, error) { return map[string]string{}, nil }

func ExistUser(chatID int64) (bool, error) { return true, nil }

func AddUserInfo() {
}

func GetUserState() (string, error) {
	return "", nil
}

func UpdateUserState(state string) error {
	return nil
}

func GetUserPasswordLk(chatID int64) (string, error) {
	return "", nil
}

func GetUserLogin(chatID int64) (string, error) {
	return "", nil
}

func AddRememberMeToken(chatID int64, rememberMe string) error {
	return nil
}

func AddPhpSessionToken(chatID int64, phpSession string) error {
	return nil
}
