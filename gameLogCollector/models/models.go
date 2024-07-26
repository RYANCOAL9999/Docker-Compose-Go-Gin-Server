package models

import "time"

type Log struct {
	ID        int64     `json:"id"`
	PlayerID  string    `json:"player_id" binding:"required"`
	Action    string    `json:"action" binding:"required"`
	Timestamp time.Time `json:"timestamp" binding:"required"`
	Details   string    `json:"details" binding:"required"`
}
