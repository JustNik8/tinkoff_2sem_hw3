package dto

import "time"

type MessageInfoDTO struct {
	Nickname string    `json:"nickname" validate:"required"`
	Message  string    `json:"message" validate:"required"`
	Time     time.Time `json:"time"`
}
