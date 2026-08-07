package domain

import (
	"errors"
	"time"
)

var (
	ErrEventNotFound = errors.New("event not found")
)

type Event struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	StartAt     time.Time  `json:"start_at"`
	EndAt       time.Time  `json:"end_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type UpdateEventInput struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	StartAt     *time.Time `json:"start_at"`
	EndAt       *time.Time `json:"end_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
