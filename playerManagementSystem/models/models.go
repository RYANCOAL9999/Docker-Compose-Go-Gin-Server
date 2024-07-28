package models

// table and return struct for Level
type Level struct {
	ID   int    `json:"id"`
	Name string `json:"name" binding:"required"`
	LV   int    `json:"lv" binding:"required"`
}

// table for Player
type Player struct {
	ID      int    `json:"id"`
	Name    string `json:"name" binding:"required"`
	LevelID int    `json:"level_id" binding:"required"`
}

// return struct for player
type PlayerRank struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	LV   int    `json:"lv"`
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
