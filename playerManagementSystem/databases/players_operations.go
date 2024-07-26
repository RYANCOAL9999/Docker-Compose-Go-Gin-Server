package databases

import (
	"database/sql"
	"errors"

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
		return nil, err
	}
	defer rows.Close()

	var playerRanks []models.PlayerRank
	for rows.Next() {
		var playerRank models.PlayerRank
		if err := rows.Scan(&playerRank.ID, &playerRank.Name, &playerRank.Rank); err != nil {
			return nil, err
		}
		playerRanks = append(playerRanks, playerRank)
	}
	return playerRanks, nil
}

func AddPlayer(db *sql.DB, name string, rank int) (*int64, error) {
	result, err := db.Exec(`
		INSERT INTO players (name, level_id) 
		SELECT ?, id 
		FROM levels 
		WHERE name = ?
	`, name, rank)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return &id, nil
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
		return nil, err
	} else if err != nil {
		return nil, err
	}
	return &playerRank, err
}

func UpdatePlayer(db *sql.DB, playerRank models.PlayerRank) error {
	var id int = playerRank.ID
	var rankNUll bool = playerRank.Rank == 0
	var query string
	var itemID int
	var playerRankScanID int
	var _ sql.Result

	if !rankNUll {
		query = `SELECT id FROM level WHERE rank = ?`
		itemID = playerRank.Rank

	} else {
		query = `SELECT id FROM player WHERE id = ?`
		itemID = id
	}

	err := db.QueryRow(query, itemID).Scan(&playerRankScanID)
	if err == sql.ErrNoRows {
		return err
	} else if err != nil {
		return err
	}

	query = `UPDATE players SET`
	args := []interface{}{}

	if !rankNUll {
		query += " level = ?"
		args = append(args, playerRankScanID)
	}

	if playerRank.Name != "" {
		query += " name = ?"
		args = append(args, playerRank.Name)
	}

	args = append(args, id)

	_, err = db.Exec(query, args...)

	if err != nil {
		return err
	}
	return nil
}

func DeletePlayer(db *sql.DB, id int) error {
	result, err := db.Exec("DELETE FROM players WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		err = errors.New("room not found")
		return err
	}
	return nil
}
