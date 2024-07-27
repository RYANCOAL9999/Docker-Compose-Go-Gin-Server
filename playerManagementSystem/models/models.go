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
