package db

func get_user_schedule(login int) (map[string]string, error) {
	return map[string]string{}, nil
}

func get_user_deadlines(login int) (map[string]string, error) {
	return map[string]string{}, nil
}

func get_active_users() (map[string]string, error) {
	return map[string]string{}, nil
}

func get_user_score() (map[string]string, error) {
	return map[string]string{}, nil
}

func update_user_schedule(login int, new_schedule map[string]string) error {
	return nil
}

func update_user_deadlines(login int, new_deadlines map[string]string) error {
	return nil
}

func update_active_users(users map[string]string) error {
	return nil
}
