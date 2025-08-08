package types

type User struct {
	Login    string
	Password string
}
type LoginRequest struct {
	User       User
	RememberMe string
}

type LoginResponse struct {
	RememberMe string
	PhpSession string
}

type Schedule struct {
	Login       int
	CurrentDate string
	PairNumber  int
	subject     string
}

type Deadline struct {
	Login        int
	CurrentDate  string
	Subject      string
	DeadlineDate string
}

type Score struct {
	Login    int
	Subject  string
	Score    int
	MaxScore int
}
