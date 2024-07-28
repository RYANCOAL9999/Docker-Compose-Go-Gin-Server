package databases

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	object "github.com/RYANCOAL9999/SpinnrTechnologyInterview/paymentProcessingSystem/databases"
	object_models "github.com/RYANCOAL9999/SpinnrTechnologyInterview/paymentProcessingSystem/models"
	"github.com/stretchr/testify/assert"
)

func TestGetPayment(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	describle := object_models.Describle{
		CardNumber:     "1234567890123456",
		ExpirationDate: "12/25",
		Key:            "123",
		Currency:       "USD",
	}
	describleJSON, _ := json.Marshal(describle)

	rows := sqlmock.NewRows([]string{"ID", "Method", "Amount", "Describle", "Timestamp"}).
		AddRow(1, "CreditCardPayment", 100.00, string(describleJSON), time.Now())

	mock.ExpectQuery("SELECT (.+) FROM Payment").WithArgs(1).WillReturnRows(rows)

	payment, err := object.GetPayment(db, 1)
	assert.NoError(t, err)
	assert.NotNil(t, payment)
	assert.Equal(t, 1, payment.ID)
	assert.Equal(t, "CreditCardPayment", payment.Method)
	assert.Equal(t, 100.00, payment.Amount)
	assert.Equal(t, describle, payment.Describle)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddPayment(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	payment := object_models.Payment{
		Method: "CreditCardPayment",
		Amount: 100.00,
		Describle: object_models.Describle{
			CardNumber:     "1234567890123456",
			ExpirationDate: "12/25",
			Key:            "123",
			Currency:       "USD",
		},
	}

	mock.ExpectExec("INSERT INTO Payment").
		WithArgs(payment.Method, payment.Amount, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	id, err := object.AddPayment(db, payment)
	assert.NoError(t, err)
	assert.Equal(t, 1, id)

	assert.NoError(t, mock.ExpectationsWereMet())
}
