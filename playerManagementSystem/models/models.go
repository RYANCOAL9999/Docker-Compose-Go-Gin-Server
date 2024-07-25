package models

type Player struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Level int    `json:"level"`
}

type Level struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
