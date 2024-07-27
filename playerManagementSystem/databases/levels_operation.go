package databases

import (
	"database/sql"
	"fmt"

	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/playerManagementSystem/models"
)

func GetLevelsData(db *sql.DB) ([]models.Level, error) {
	rows, err := db.Query(`
		SELECT 
		ID, Name, LV 
		FROM Level
	`)
	if err != nil {
		return nil, fmt.Errorf("error querying database with GetLevelsData: %w", err)
	}
	defer rows.Close()
	var levels []models.Level
	for rows.Next() {
		var l models.Level
		err := rows.Scan(
			&l.ID,
			&l.Name,
			&l.LV,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row with GetLevelsData: %w", err)
		}
		levels = append(levels, l)
	}
	return levels, nil
}

func AddLevel(db *sql.DB, name string, lv int) (int, error) {
	result, err := db.Exec(`
		INSERT INTO Level (Name, LV) 
		VALUES (?, ?)
	`, name, lv)
	if err != nil {
		return 0, fmt.Errorf("error querying database with AddLevel: %w", err)
	}
	id, _ := result.LastInsertId()
	return int(id), nil
}
