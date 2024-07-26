package models

import "time"

type PlayerRank struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Rank int    `json:"rank"`
}

type Status int

const (
	StatusAvailable Status = iota
	StatusOccupied
	StatusMaintenance
	StatusClosed
)

type Room struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Status      Status `json:"status"`
	Description string `json:"description"`
	Player_ids  string `json:"player_ids"`
}

type ReservationCreate struct {
	ID     int       `json:"id"`
	RoomID int       `json:"room_id"`
	Date   time.Time `json:"date"`
}

type Reservation struct {
	ID     int          `json:"id"`
	RoomID int          `json:"room_id"`
	Date   time.Time    `json:"date"`
	Player []PlayerRank `json:"player_ids"`
}
