package domain

import "time"

type MessageInfo struct {
	ID       string
	Nickname string
	Message  string
	Time     time.Time
}
