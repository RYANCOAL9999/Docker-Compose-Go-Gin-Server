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
		ID, PlayerID, Amount, Status, Won, CreatedAt, Probability 
		FROM Challenge
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
		err := rows.Scan(
			&challenge.ID,
			&challenge.PlayerID,
			&challenge.Amount,
			&challenge.Status,
			&challenge.Won,
			&challenge.CreatedAt,
			&challenge.Probability,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row with ListChallenges: %w", err)
		}
		challenges = append(challenges, challenge)
	}
	return challenges, nil
}

func GetLastChallengeTime(db *sql.DB, playerID int) (*time.Time, error) {
	var lastChallengeTime time.Time
	err := db.QueryRow(`
		SELECT 
		CreatedAt 
		FROM Challenge 
		WHERE PlayerID = ? 
		ORDER BY CreatedAt DESC 
		LIMIT 1
	`, playerID).Scan(
		&lastChallengeTime,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error querying database with CreateChallenge: %w", err)
	}
	return &lastChallengeTime, nil
}

func GetChallenge(db *sql.DB, challengeID int, playerID int) (*models.Status, *float64, error) {
	var status models.Status
	var probability float64
	err := db.QueryRow(`
		SELECT 
		Status, Probability 
		FROM Challenge 
		WHERE ID = ? AND PlayerID = ? 
	`, challengeID, playerID).Scan(
		&status,
		&probability,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, nil, fmt.Errorf("error querying database with CreateChallenge: %w", err)
	}
	return &status, &probability, nil
}

func AddNewChallenge(tx *sql.Tx, newChallengeNeed models.NewChallengeNeed) (int, error) {
	result, err := tx.Exec(`
		INSERT INTO Challenge (PlayerID, Amount, Status, Won, CreatedAt, Probability) 
		VALUES (?, ?, 1, false, NOW(), 0)
	`, newChallengeNeed.PlayerID, newChallengeNeed.Amount)
	if err != nil {
		return 0, fmt.Errorf("error querying database with addNewChallenge: %w", err)
	}
	challengeID, _ := result.LastInsertId()
	return int(challengeID), nil
}

func UpdatePricePool(tx *sql.Tx, amount float64) error {
	_, err := tx.Exec(`
		UPDATE PrizePool 
		SET Amount = Amount + ? 
		WHERE ID = 1
	`, amount)
	if err != nil {
		return fmt.Errorf("error updating price pool: %w", err)
	}
	return nil
}

func UpdateChallenge(tx *sql.Tx, status models.Status, won bool, playerID int) error {
	_, err := tx.Exec(`
		UPDATE Challenge 
		SET Status = ?, Won = ? WHERE PlayerID = ?
	`, status, won, playerID)
	if err != nil {
		return fmt.Errorf("error updating challenge: %w", err)
	}
	return nil
}

func DistributePrizePool(tx *sql.Tx, challengeID int, playerID int) error {
	// know the player last win how much money
	var prize float64
	err := tx.QueryRow(`
		SELECT 
		Amount 
		FROM PrizePool 
		WHERE ID = 1
	`).Scan(
		&prize,
	)
	if err != nil {
		return fmt.Errorf("error fetching prize pool amount: %w", err)
	}

	// Update last challenges's won
	_, err = tx.Exec(`
		UPDATE challenges 
		SET Won = 1, Probability = 0 
		WHERE ID = ? AND PlayerID = ?
	`, challengeID, playerID)
	if err != nil {
		return fmt.Errorf("error updating player's balance: %w", err)
	}

	// Reset prize pool
	_, err = tx.Exec(
		`UPDATE PrizePool 
		SET Amount = 0 
		WHERE ID = 1
	`)
	if err != nil {
		return fmt.Errorf("error resetting prize pool: %w", err)
	}

	log.Printf("Prize pool of %f distributed to player %d and challenge %d", prize, playerID, challengeID)
	return nil
}

func Updateprobability(tx *sql.Tx, challengeID int, playerID int, probability float64) error {

	// Update last challenge's won
	_, err := tx.Exec(`
		UPDATE Challenge 
		SET Probability = 0 
		WHERE ID = ? AND PlayerID = ?
	`, challengeID, playerID)
	if err != nil {
		return fmt.Errorf("error updating player's balance: %w", err)
	}

	log.Printf("probability of %f with player %d and challenge %d", probability, playerID, challengeID)
	return nil
}
