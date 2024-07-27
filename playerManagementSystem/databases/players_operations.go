package databases

import (
	"database/sql"
	"fmt"

	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/playerManagementSystem/models"
)

func GetPlayersData(db *sql.DB) ([]models.PlayerRank, error) {
	rows, err := db.Query(`
		SELECT 
		P.ID as ID, 
		P.Name as Name,
		L.LV as LV
		FROM Player P
		INNER JOIN 
		Level L 
		ON P.LevelID = L.ID
		ORDER BY ID
	`)
	if err != nil {
		return nil, fmt.Errorf("error querying database with GetPlayersData: %w", err)
	}
	defer rows.Close()

	var playerRanks []models.PlayerRank
	for rows.Next() {
		var playerRank models.PlayerRank
		err := rows.Scan(
			&playerRank.ID,
			&playerRank.Name,
			&playerRank.LV,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row with GetPlayersData: %w", err)
		}
		playerRanks = append(playerRanks, playerRank)
	}
	return playerRanks, nil
}

func AddPlayer(db *sql.DB, name string, lv int) (int, error) {
	result, err := db.Exec(`
		INSERT INTO Player (Name, LevelID) 
		SELECT 
		?, ID 
		FROM Level 
		WHERE LV = ?
	`, name, lv)
	if err != nil {
		return 0, fmt.Errorf("error querying database with AddPlayer: %w", err)
	}

	id, _ := result.LastInsertId()
	return int(id), nil
}

func GetPlayer(db *sql.DB, id int) (*models.PlayerRank, error) {
	var playerRank models.PlayerRank
	err := db.QueryRow(`
		SELECT 
		P.ID as ID, 
		P.Name as Name,
		L.LV as LV
		FROM Player P
		INNER JOIN 
		Levels L 
		ON P.LevelID = L.ID
		WHERE P.ID = ?
	`, id).Scan(&playerRank.ID, &playerRank.Name, &playerRank.LV)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("error querying database with GetPlayer: %w", err)
	} else if err != nil {
		return nil, fmt.Errorf("error scanning row with GetPlayer: %w", err)
	}
	return &playerRank, err
}

func UpdatePlayer(db *sql.DB, playerRank models.PlayerRank) error {
	query := "UPDATE players SET"
	args := []interface{}{}

	if playerRank.LV != 0 {
		var levelID int
		err := db.QueryRow(`
			SELECT 
			ID
			FROM Level 
			WHERE LV = ?
		`, playerRank.LV).Scan(&levelID)
		if err != nil {
			return fmt.Errorf("error querying database with UpdatePlayer: %w", err)
		}
		query += " LevelID = ?"
		args = append(args, levelID)
	}

	if playerRank.Name != "" {
		if len(args) > 0 {
			query += ","
		}
		query += " Name = ?"
		args = append(args, playerRank.Name)
	}

	if len(args) == 0 {
		return fmt.Errorf("no fields update for player with id: %d", playerRank.ID)
	}

	query += " WHERE ID = ?"
	args = append(args, playerRank.ID)

	// Execute the update query
	result, err := db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error updating player: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows were updated, player with id %d may not exist", playerRank.ID)
	}

	return nil
}

func DeletePlayer(db *sql.DB, id int) error {
	result, err := db.Exec(`
		DELETE FROM players 
		WHERE ID = ?
	`, id)
	if err != nil {
		return fmt.Errorf("error querying database with DeletePlayer: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no rows were deleted, player with id %d may not exist", id)
	}
	return nil
}
