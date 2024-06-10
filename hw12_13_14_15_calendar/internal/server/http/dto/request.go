package dto

import "time"

type CreateEvent struct {
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	UserID       string    `json:"userId"`
	RemindBefore int64     `json:"remindBefore"`
	StartTime    time.Time `json:"startTime"`
	EndTime      time.Time `json:"endTime"`
}

type UpdateEvent struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	UserID       string    `json:"userId"`
	RemindBefore int64     `json:"remindBefore"`
	StartTime    time.Time `json:"startTime"`
	EndTime      time.Time `json:"endTime"`
}
