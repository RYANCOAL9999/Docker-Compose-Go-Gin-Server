package databases

import (
	"database/sql"

	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/playerManagementSystem/models"
)

func GetLevelsData(db *sql.DB) ([]models.Level, error) {
	rows, err := db.Query("SELECT id, name FROM levels")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var levels []models.Level
	for rows.Next() {
		var l models.Level
		if err := rows.Scan(&l.ID, &l.Name); err != nil {
			return nil, err
		}
		levels = append(levels, l)
	}
	return levels, nil
}

func AddLevel(db *sql.DB, name string, rank int) (*int64, error) {
	result, err := db.Exec("INSERT INTO levels (name, rank) VALUES (?)", name, rank)
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	return &id, nil
}
