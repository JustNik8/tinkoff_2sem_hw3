package dto

import "time"

type MessageInfoRequest struct {
	Nickname string `json:"nickname" validate:"required"`
	Message  string `json:"message" validate:"required"`
}

type MessageInfoResponse struct {
	Nickname string    `json:"nickname"`
	Message  string    `json:"message"`
	Time     time.Time `json:"time"`
}
