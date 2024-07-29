package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	object "github.com/RYANCOAL9999/SpinnrTechnologyInterview/endlessChallengeSystem/handlers"
	object_models "github.com/RYANCOAL9999/SpinnrTechnologyInterview/endlessChallengeSystem/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestJoinChallenges(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/challenges", func(c *gin.Context) {
		object.JoinChallenges(c, db)
	})

	// Mock expectations
	mock.ExpectQuery(`SELECT CreatedAt FROM Challenge WHERE PlayerID = \?`).
		WithArgs(1001). // Assuming PlayerID is the argument
		WillReturnRows(sqlmock.NewRows([]string{"CreatedAt"}).
			AddRow(time.Now().Add(-2 * time.Minute)))

	mock.ExpectBegin()

	mock.ExpectExec(`INSERT INTO Challenge \(PlayerID, Amount, Status, Won, CreatedAt, Probability\)`).
		WithArgs(1001, 20.01, object_models.Ready, 0, sqlmock.AnyArg(), 0.5).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(`UPDATE PrizePool SET Amount = \? WHERE ID = \?`).
		WithArgs(20.01, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	mock.ExpectQuery(`SELECT Status, Probability FROM Challenge WHERE ID = \?`).
		WithArgs(1). // Assuming Challenge ID is 1
		WillReturnRows(sqlmock.NewRows([]string{"Status", "Probability"}).
			AddRow(object_models.Joined, 0.5))

	// Create the request
	newChallenge := object_models.NewChallengeNeed{PlayerID: 1001, Amount: 20.01}
	jsonValue, _ := json.Marshal(newChallenge)
	req, _ := http.NewRequest("POST", "/challenges", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()

	// Serve the request
	r.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusCreated, w.Code)

	var response object_models.JoinChallengeResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, object_models.Joined, response.Status)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestJoinChallenges_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Register the route with the correct path
	router.POST("/challenges", func(c *gin.Context) {
		object.JoinChallenges(c, db)
	})

	// Invalid JSON data
	invalidJSON := `{"PlayerID": "player123", "Amount": "invalid_amount"}`

	req, _ := http.NewRequest(http.MethodPost, "/challenges", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert that the response status code is 400 Bad Request
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestShowChallenges(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/challenges/results", func(c *gin.Context) {
		object.ShowChallenges(c, db)
	})

	// Mocking the SQL query and response
	rows := sqlmock.NewRows([]string{"ID", "PlayerID", "Amount", "Status", "Won", "CreatedAt", "Probability"}).
		AddRow(1, 1001, 100.01, object_models.Joined, false, time.Now().Format(time.RFC3339), 0.5)

	mock.ExpectQuery("SELECT (.+) FROM Challenge LIMIT ?").WithArgs(1).WillReturnRows(rows)

	// Making the HTTP request
	req, _ := http.NewRequest("GET", "/challenges/results?limit=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Checking the response status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Unmarshalling the response into a slice of Challenge
	var response []object_models.Challenge
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verifying the contents of the response
	assert.Equal(t, 1, len(response))
	assert.Equal(t, 1, response[0].ID)
	assert.Equal(t, 1001, response[0].PlayerID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestShowChallenges_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		setupMock          func(mock sqlmock.Sqlmock)
		queryParams        string
		expectedStatusCode int
		expectedErrorMsg   string
	}{
		{
			name: "Database error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnError(sql.ErrConnDone)
			},
			queryParams:        "limit=10",
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorMsg:   "sql: connection is already closed",
		},
		{
			name: "Invalid limit parameter",
			setupMock: func(mock sqlmock.Sqlmock) {
				// No expectations set as the error occurs before DB interaction
			},
			queryParams:        "limit=invalid",
			expectedStatusCode: http.StatusInternalServerError, // Note: The function doesn't handle this error explicitly
			expectedErrorMsg:   "",                             // No error message expected in this case
		},
		{
			name: "No challenges found",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "player_id", "amount"}))
			},
			queryParams:        "limit=10",
			expectedStatusCode: http.StatusOK,
			expectedErrorMsg:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mock)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(http.MethodGet, "/challenges?"+tt.queryParams, nil)

			object.ShowChallenges(c, db)

			assert.Equal(t, tt.expectedStatusCode, w.Code)

			if tt.expectedErrorMsg != "" {
				var response object_models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response.Error, tt.expectedErrorMsg)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestCalculateChallengeResult(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	challengeID := 1
	playerID := 1
	probability := 0.5

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE Challenge SET").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	object.CalculateChallengeResult(db, challengeID, playerID, probability)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCalculateChallengeResult_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	challengeID := 1
	playerID := 1
	probability := 0.5

	mock.ExpectBegin().WillReturnError(fmt.Errorf("transaction start error"))

	object.CalculateChallengeResult(db, challengeID, playerID, probability)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
