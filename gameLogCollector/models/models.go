package models

import "time"

// table for GameLog
type GameLog struct {
	ID        int64     `json:"id"`
	PlayerID  int       `json:"player_id" binding:"required"`
	Action    string    `json:"action" binding:"required"`
	Timestamp time.Time `json:"timestamp"`
	Details   string    `json:"details" binding:"required"`
}
