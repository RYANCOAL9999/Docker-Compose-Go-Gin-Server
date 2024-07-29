package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	object "github.com/RYANCOAL9999/SpinnrTechnologyInterview/paymentProcessingSystem/handlers"
	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/paymentProcessingSystem/models"
	object_models "github.com/RYANCOAL9999/SpinnrTechnologyInterview/paymentProcessingSystem/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestShowPayment(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/payments/:id", func(c *gin.Context) {
		object.ShowPayment(c, db)
	})

	describle := object_models.Describle{
		CardNumber:     "1234567890123456",
		ExpirationDate: "12/25",
		Key:            "123",
		Currency:       "USD",
	}
	describleJSON, _ := json.Marshal(describle)

	rows := sqlmock.NewRows([]string{"ID", "Method", "Amount", "Describle", "Timestamp"}).
		AddRow(1, "CreditCardPayment", 100.00, string(describleJSON), "2023-01-01T00:00:00Z")

	mock.ExpectQuery("SELECT (.+) FROM Payment").WithArgs(1).WillReturnRows(rows)

	req, _ := http.NewRequest("GET", "/payments/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response object_models.Payment
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 1, response.ID)
	assert.Equal(t, "CreditCardPayment", response.Method)
	assert.Equal(t, 100.00, response.Amount)
	assert.Equal(t, describle, response.Describle)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestShowPayment_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	paymentID := 1
	payment := object_models.Payment{ID: paymentID, Amount: 100.0}
	mock.ExpectQuery("SELECT (.+) FROM payments WHERE id = ?").
		WithArgs(paymentID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "amount"}).AddRow(payment.ID, payment.Amount))

	router.GET("/payments/:id", func(c *gin.Context) {
		object.ShowPayment(c, db)
	})

}

func TestCreatePayment(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/payments", func(c *gin.Context) {
		object.CreatePayment(c, db)
	})

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

	jsonValue, _ := json.Marshal(payment)

	mock.ExpectExec("INSERT INTO Payment").
		WithArgs(payment.Method, payment.Amount, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	req, _ := http.NewRequest("POST", "/payments", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response object_models.PaymentResult
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 1, response.ID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreatePayment_Error(t *testing.T) {

	gin.SetMode(gin.TestMode)

	t.Run("Invalid JSON payload", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		invalidJSON := []byte(`{"method": "CreditCardPayment", "amount": "invalid"}`)
		c.Request, _ = http.NewRequest("POST", "/payments", bytes.NewBuffer(invalidJSON))
		c.Request.Header.Set("Content-Type", "application/json")

		db, _, _ := sqlmock.New()
		defer db.Close()

		object.CreatePayment(c, db)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "json: cannot unmarshal string into Go struct field")
	})

	t.Run("Unsupported payment method", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		payment := models.Payment{Method: "UnsupportedMethod", Amount: 100.00}
		payloadBytes, _ := json.Marshal(payment)
		c.Request, _ = http.NewRequest("POST", "/payments", bytes.NewBuffer(payloadBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		db, mock, _ := sqlmock.New()
		defer db.Close()

		mock.ExpectExec("INSERT INTO Payment").
			WithArgs(payment.Method, payment.Amount, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		object.CreatePayment(c, db)

		assert.Equal(t, http.StatusCreated, w.Code)
		var result models.PaymentResult
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Equal(t, 1, result.ID)
		assert.Empty(t, result.Status) // Status should be empty for unsupported method
	})

	t.Run("Database error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		payment := models.Payment{Method: "CreditCardPayment", Amount: 100.00}
		payloadBytes, _ := json.Marshal(payment)
		c.Request, _ = http.NewRequest("POST", "/payments", bytes.NewBuffer(payloadBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		db, mock, _ := sqlmock.New()
		defer db.Close()

		mock.ExpectExec("INSERT INTO Payment").
			WithArgs(payment.Method, payment.Amount, sqlmock.AnyArg()).
			WillReturnError(sql.ErrConnDone)

		object.CreatePayment(c, db)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "error querying database with AddPayment")
	})

}
