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

func ExistUser(chat_id int) (bool, error) { return true, nil }
