package databases

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/endlessChallengeSystem/models"
)

func ListChallenges(db *sql.DB, limit int) ([]models.Challenge, error) {
	var query string = `
		SELECT 
		id, player_id, amount, won, created_at 
		FROM challenges
		ORDER BY id
	`
	args := []interface{}{}

	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)
	}

	rows, err := db.Query(query, args...)

	if err != nil {
		return nil, fmt.Errorf("error querying database with ListChallenges: %w", err)
	}

	defer rows.Close()

	var challenges []models.Challenge
	for rows.Next() {
		var challenge models.Challenge
		err := rows.Scan(&challenge.ID, &challenge.PlayerID, &challenge.Amount, &challenge.Won, &challenge.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning row with ListChallenges: %w", err)
		}
		challenges = append(challenges, challenge)
	}
	return challenges, nil
}

func GetLastChallengeTime(db *sql.DB, playerIDs int) (*time.Time, error) {
	var lastChallengeTime time.Time
	err := db.QueryRow("SELECT created_at FROM challenges WHERE player_id = ? ORDER BY created_at DESC LIMIT 1", playerIDs).Scan(&lastChallengeTime)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error querying database with CreateChallenge: %w", err)
	}
	return &lastChallengeTime, nil
}

func GetChallengeStatus(db *sql.DB, playerID int) (*int, error) {
	var status int
	err := db.QueryRow("SELECT status FROM challenges WHERE player_id = ? ", playerID).Scan(&status)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error querying database with CreateChallenge: %w", err)
	}
	return &status, nil
}

func AddNewChallenge(tx *sql.Tx, newChallengeNeed models.NewChallengeNeed) error {
	_, err := tx.Exec(`
		INSERT INTO challenges (player_id, amount, status, won, created_at) 
		VALUES (?, ?, 1, false, NOW())
		ON DUPLICATE KEY UPDATE 
		amount = VALUES(amount),
		status = VALUES(status),
		won = VALUES(won),
		created_at = VALUES(created_at)
	`, newChallengeNeed.PlayerID, newChallengeNeed.Amount)
	if err != nil {
		return fmt.Errorf("error querying database with addNewChallenge: %w", err)
	}
	return nil
}

func UpdatePricePool(tx *sql.Tx, amount float64) error {
	_, err := tx.Exec("UPDATE prize_pool SET amount = amount + ?", amount)
	if err != nil {
		return fmt.Errorf("error updating price pool: %w", err)
	}
	return nil
}

func UpdateChallenge(tx *sql.Tx, status models.Status, won bool, playerID int) error {
	_, err := tx.Exec(`
		UPDATE challenges SET status = ?, won = ? WHERE player_id = ?
	`, status, won, playerID)
	if err != nil {
		return fmt.Errorf("error updating challenge: %w", err)
	}
	return nil
}

func DistributePrizePool(tx *sql.Tx, playerID int) error {

	var prize float64
	err := tx.QueryRow("SELECT amount FROM prize_pool").Scan(&prize)
	if err != nil {
		return fmt.Errorf("error fetching prize pool amount: %w", err)
	}

	// Update challenges's amount
	_, err = tx.Exec("UPDATE challenges SET amount = amount + ? WHERE player_id = ? AND status = ?", prize, playerID, models.Ready)
	if err != nil {
		return fmt.Errorf("error updating player's balance: %w", err)
	}

	// Reset prize pool
	_, err = tx.Exec("UPDATE prize_pool SET amount = 0")
	if err != nil {
		return fmt.Errorf("error resetting prize pool: %w", err)
	}

	log.Printf("Prize pool of %f distributed to player %d", prize, playerID)
	return nil
}
