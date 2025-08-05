package db

func GetUserSchedule(chatID int64) (map[string]string, error) {
	return map[string]string{}, nil
}

func GetUserDeadlines(login int) (map[string]string, error) {
	return map[string]string{}, nil
}

func GetActiveUsers() (int, error) {
	return 0, nil
}

func GetUserScore() (map[string]string, error) {
	return map[string]string{}, nil
}

func AddUserSchedule(login int) error {
	return nil
}

func AddUserDeadlines(login int) error {
	return nil
}

func AddActiveUsers() error {
	return nil
}

func AddUserScore() error {
	return nil
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

func AddLoginUser(number int) error {
	return nil
}

func AddPasswordLkUser(password string) error {
	return nil
}

func AddPasswordEduUser(password string) error {
	return nil
}

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
