package databases

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/paymentProcessingSystem/models"
)

func GetPayment(db *sql.DB, id int) (*models.Payment, error) {
	var payment models.Payment
	var describleString string
	err := db.QueryRow(`
		SELECT 
		ID, Method, Amount, Describle, Timestamp 
		FROM Payment 
		WHERE ID = ?
	`, id).Scan(
		&payment.ID,
		&payment.Method,
		&payment.Amount,
		&describleString,
		&payment.Timestamp,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("error querying database with GetPayment: %w", err)
	} else if err != nil {
		return nil, fmt.Errorf("error scanning row with GetPayment: %w", err)
	}

	err = json.Unmarshal([]byte(describleString), &payment.Describle)
	if err != nil {
		return nil, fmt.Errorf("error json Unmarshal with GetPayment: %w", err)
	}

	return &payment, err
}

func AddPayment(db *sql.DB, payment models.Payment) (int, error) {
	jsonBytes, err := json.Marshal(payment.Describle)
	if err != nil {
		return 0, fmt.Errorf("error marshal on AddPayment: %w", err)
	}
	result, err := db.Exec(`
		INSERT INTO Payment (Method, Amount, Describle, Timestamp) 
		VALUES (?, ?, ?, Now())
	`, payment.Method, payment.Amount, jsonBytes)
	if err != nil {
		return 0, fmt.Errorf("error querying database with AddPayment: %w", err)
	}
	id, _ := result.LastInsertId()
	return int(id), err
}
