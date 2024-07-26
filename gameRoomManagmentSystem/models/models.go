package models

import "time"

type Status int

const (
	StatusAvailable Status = iota
	StatusOccupied
	StatusMaintenance
	StatusClosed
)

// table for Room
type Room struct {
	ID          int    `json:"id"`
	Name        string `json:"name" binding:"required"`
	Status      Status `json:"status" binding:"required"`
	Description string `json:"description" binding:"required"`
	PlayerIDs   string `json:"player_ids" binding:"required"`
}

// struct for Completed the ReservationRoom
type PlayerRank struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Rank int    `json:"rank"`
}

// table for Reservation
type Reservation struct {
	ID     int       `json:"id"`
	RoomID int       `json:"room_id" binding:"required"`
	Date   time.Time `json:"date" binding:"required"`
}

// return struct for Reservation
type ReservationRoom struct {
	ID     int          `json:"id"`
	RoomID int          `json:"room_id"`
	Date   time.Time    `json:"date"`
	Player []PlayerRank `json:"player"`
}
