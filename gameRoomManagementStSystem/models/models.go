package models

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
	Status      Status `json:"status"`
	Description string `json:"description" binding:"required"`
	PlayerIDs   string `json:"player_ids"`
}

// struct for Completed the ReservationRoom
type PlayerRank struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	LV   int    `json:"lv"`
}

// struct for Reservation request
type Reservation struct {
	ID        int    `json:"id"`
	RoomID    int    `json:"room_id" binding:"required"`
	Date      string `json:"date" binding:"required"`
	PlayerIDs string `json:"player_ids" binding:"required"`
}

// return struct for Reservation
type ReservationRoom struct {
	ID     int          `json:"id"`
	RoomID int          `json:"room_id"`
	Date   string       `json:"date"`
	Player []PlayerRank `json:"player"`
}

// ErrorResponse represents an error response with a single error message.
type ErrorResponse struct {
	Error string `json:"error"`
}

// CreateResponse represents an id after created a item.
type CreateResponse struct {
	ID int `json:"id"`
}

// SuccessResponse represents an any after update or delete item.
type SuccessResponse struct {
}
