package db

func GetUserSchedule(login int) (map[string]string, error) {
	return map[string]string{}, nil
}

func GetUserDeadlines(login int) (map[string]string, error) {
	return map[string]string{}, nil
}

func GetActiveUsers() (map[string]string, error) {
	return map[string]string{}, nil
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

func ExistUser(chatID int) (bool, error) { return true, nil }

func AddLoginUser(number int) error {
	return nil
}

func AddPasswordLkUser(password string) error {
	return nil
}

func AddPasswordEduUser(password string) error {
	return nil
}

func AddUser() error {
	return nil
}

func GetUserState() (string, error) {
	return "", nil
}

func UpdateUserState(state string) error {
	return nil
}

func GetUserPasswordLk(chatID int) (string, error) {
	return "", nil
}

func GetUserLogin(chatID int) (string, error) {
	return "", nil
}
