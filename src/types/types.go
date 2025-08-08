package types

import "time"

type LoginRequest struct {
	Login      string
	Password   string
	RememberMe string
}

type User struct {
	ChatID int64
	Login  string
}

type Schedule struct {
	Login        string
	ScheduleDate time.Time
	DayOfWeek    int
	WeekNumber   int
	PairNumber   int
	Subject      string
	StartTime    string
	EndTime      string
	Auditory     string
	Teacher      string
}

type Deadline struct {
	ID           int
	Login        string
	TaskName     string
	DeadlineDate time.Time
	Subject      string
	IsCompleted  bool
}

type Notification struct {
	ID      int
	UserID  int64
	Type    string
	Title   string
	Message string
	IsSent  bool
}

type Score struct {
	Login    string
	ChatID   int64
	Subject  string
	Score    int
	MaxScore int
}
