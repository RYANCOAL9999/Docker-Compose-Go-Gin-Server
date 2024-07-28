package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	object "github.com/RYANCOAL9999/SpinnrTechnologyInterview/paymentProcessingSystem/handlers"
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
