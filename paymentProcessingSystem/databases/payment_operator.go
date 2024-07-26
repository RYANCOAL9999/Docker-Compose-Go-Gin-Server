package databases

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/paymentProcessingSystem/models"
)

func GetPayment(db *sql.DB, id int) (*models.Payment, error) {
	var payment models.Payment
	err := db.QueryRow(`
		SELECT id, method, amount, describle FROM payment WHERE id = ?
	`, id).Scan(&payment.ID, &payment.Method, &payment.Amount, &payment.Describle)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("error querying database with GetPlayer: %w", err)
	} else if err != nil {
		return nil, fmt.Errorf("error scanning row with GetPlayer: %w", err)
	}
	return &payment, err
}

func AddPayment(db *sql.DB, payment models.Payment) (int, error) {
	result, err := db.Exec(`
		INSERT INTO payment (method, amount, describle, timestamp) VALUES (?, ?, ?, ?)
	`, payment.Method, payment.Amount, payment.Describle, time.Now())
	if err != nil {
		return 0, fmt.Errorf("error querying database with AddPayment: %w", err)
	}
	id, _ := result.LastInsertId()
	return int(id), err
}
