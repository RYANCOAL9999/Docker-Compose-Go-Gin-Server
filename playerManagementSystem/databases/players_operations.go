package databases

import (
	"database/sql"
	"fmt"

	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/playerManagementSystem/models"
)

func GetPlayersData(db *sql.DB) ([]models.PlayerRank, error) {
	rows, err := db.Query(`
		SELECT p.id, p.name, l.rank 
		FROM players p
		INNER JOIN 
		level l 
		ON p.level_id = l.id
		ORDER BY p.id
	`)
	if err != nil {
		return nil, fmt.Errorf("error querying database with GetPlayersData: %w", err)
	}
	defer rows.Close()

	var playerRanks []models.PlayerRank
	for rows.Next() {
		var playerRank models.PlayerRank
		if err := rows.Scan(&playerRank.ID, &playerRank.Name, &playerRank.Rank); err != nil {
			return nil, fmt.Errorf("error scanning row with GetPlayersData: %w", err)
		}
		playerRanks = append(playerRanks, playerRank)
	}
	return playerRanks, nil
}

func AddPlayer(db *sql.DB, name string, rank int) (int, error) {
	result, err := db.Exec(`
		INSERT INTO players (name, level_id) 
		SELECT ?, id 
		FROM levels 
		WHERE name = ?
	`, name, rank)
	if err != nil {
		return 0, fmt.Errorf("error querying database with AddPlayer: %w", err)
	}

	id, _ := result.LastInsertId()
	return int(id), nil
}

func GetPlayer(db *sql.DB, id int) (*models.PlayerRank, error) {
	var playerRank models.PlayerRank
	err := db.QueryRow(`
		SELECT p.id, p.name, l.rank 
		FROM players p
		INNER JOIN 
		level l 
		ON p.level_id = l.id
		WHERE id = ?
	`, id).Scan(&playerRank.ID, &playerRank.Name, &playerRank.Rank)
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

	if playerRank.Rank != 0 {
		var levelID int
		err := db.QueryRow("SELECT id FROM level WHERE rank = ?", playerRank.Rank).Scan(&levelID)
		if err != nil {
			return fmt.Errorf("error querying database with UpdatePlayer: %w", err)
		}
		query += " level_id = ?"
		args = append(args, levelID)
	}

	if playerRank.Name != "" {
		if len(args) > 0 {
			query += ","
		}
		query += " name = ?"
		args = append(args, playerRank.Name)
	}

	if len(args) == 0 {
		return fmt.Errorf("no fields update for player with id: %d", playerRank.ID)
	}

	query += " WHERE id = ?"
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
	result, err := db.Exec("DELETE FROM players WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("error querying database with DeletePlayer: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no rows were deleted, player with id %d may not exist", id)
	}
	return nil
}
