package models

import (
	"time"
)

// table for GameLog
type GameLog struct {
	ID        int64     `json:"id"`
	PlayerID  int       `json:"player_id" binding:"required"`
	Action    string    `json:"action" binding:"required"`
	Timestamp time.Time `json:"timestamp"`
	Details   string    `json:"details" binding:"required"`
}

// ErrorResponse represents an error response with a single error message.
type ErrorResponse struct {
	Error string `json:"error"`
}

// CreateResponse represents an id after created a item.
type CreateResponse struct {
	ID int `json:"id"`
}
