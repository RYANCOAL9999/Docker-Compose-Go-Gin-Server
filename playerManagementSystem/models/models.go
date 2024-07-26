package models

// table and return struct for Level
type Level struct {
	ID   int    `json:"id"`
	Name string `json:"name" binding:"required"`
	Rank int    `json:"Rank" binding:"required"`
}

// table for Player
type Player struct {
	ID       int    `json:"id"`
	Name     string `json:"name" binding:"required"`
	Level_id int    `json:"Level_id" binding:"required"`
}

// return struct for player
type PlayerRank struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Rank int    `json:"rank"`
}
