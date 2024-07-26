package models

type Level struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Rank int    `json:"Rank"`
}

type PlayerRank struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Rank int    `json:"rank"`
}
