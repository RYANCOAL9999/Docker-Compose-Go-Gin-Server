package databases

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	object "github.com/RYANCOAL9999/SpinnrTechnologyInterview/paymentProcessingSystem/databases"
	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/paymentProcessingSystem/models"
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

func TestGetPayment_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	t.Run("No rows error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM Payment WHERE ID = ?").
			WithArgs(1).
			WillReturnError(sql.ErrNoRows)

		payment, err := object.GetPayment(db, 1)

		assert.Nil(t, payment)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error querying database with GetPayment")
	})

	t.Run("Query execution error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM Payment WHERE ID = ?").
			WithArgs(2).
			WillReturnError(sql.ErrConnDone)

		payment, err := object.GetPayment(db, 2)

		assert.Nil(t, payment)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error scanning row with GetPayment")
	})

	t.Run("JSON Unmarshal error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"ID", "Method", "Amount", "Describle", "Timestamp"}).
			AddRow(3, "Credit Card", 100.00, "{invalid json}", time.Now())

		mock.ExpectQuery("SELECT (.+) FROM Payment WHERE ID = ?").
			WithArgs(3).
			WillReturnRows(rows)

		payment, err := object.GetPayment(db, 3)

		assert.Nil(t, payment)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error json Unmarshal with GetPayment")
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
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

func TestAddPayment_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	t.Run("JSON required error", func(t *testing.T) {
		payment := object_models.Payment{
			Method: "Credit Card",
			Amount: 100.00,
		}

		id, err := object.AddPayment(db, payment)

		assert.Equal(t, 0, id)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error missing param on AddPayment")
	})

	t.Run("Database execution error", func(t *testing.T) {
		payment := models.Payment{
			Method:    "Credit Card",
			Amount:    100.00,
			Describle: object_models.Describle{},
		}

		mock.ExpectExec("INSERT INTO Payment").
			WithArgs(payment.Method, payment.Amount, sqlmock.AnyArg()).
			WillReturnError(sql.ErrConnDone)

		id, err := object.AddPayment(db, payment)

		assert.Equal(t, 0, id)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error querying database with AddPayment")
	})

	t.Run("LastInsertId error", func(t *testing.T) {
		payment := models.Payment{
			Method:    "Credit Card",
			Amount:    100.00,
			Describle: object_models.Describle{},
		}

		mock.ExpectExec("INSERT INTO Payment").
			WithArgs(payment.Method, payment.Amount, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewErrorResult(sql.ErrNoRows))

		id, err := object.AddPayment(db, payment)

		assert.Equal(t, 0, id)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error querying database with AddPayment")
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}
