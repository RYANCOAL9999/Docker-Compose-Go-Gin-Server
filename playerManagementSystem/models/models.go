package models

// table and return struct for Level
type Level struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Rank int    `json:"Rank"`
}

// table for Player
type Player struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Level_id int    `json:"Level_id"`
}

// return struct for player
type PlayerRank struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Rank int    `json:"rank"`
}
