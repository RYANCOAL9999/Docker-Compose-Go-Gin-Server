package models

import (
	"time"
)

type Status int

const (
	Ready Status = iota
	joined
)

// table for Prize Pool
type PrizePool struct {
	Amount float64 `json:"amount"`
}

// table for Challenge
type Challenge struct {
	ID        int       `json:"id"`
	PlayerID  string    `json:"player_id"`
	Amount    float64   `json:"amount"`
	Status    Status    `json:"status"`
	Won       bool      `json:"won"`
	CreatedAt time.Time `json:"created_at"`
}

// New Challenge Struct for request
type NewChallengeNeed struct {
	PlayerID int     `json:"player_id" binding:"required"`
	Amount   float64 `json:"amount" binding:"required,eq=20.01"`
}
