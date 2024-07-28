package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
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
	r.POST("/challenges/join", func(c *gin.Context) {
		object.JoinChallenges(c, db)
	})

	mock.ExpectQuery("SELECT CreatedAt FROM Challenge").WillReturnRows(sqlmock.NewRows([]string{"CreatedAt"}).AddRow(time.Now().Add(-2 * time.Minute)))
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO Challenge").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("UPDATE PrizePool").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectQuery("SELECT Status, Probability FROM Challenge").WillReturnRows(sqlmock.NewRows([]string{"Status", "Probability"}).AddRow(object_models.Joined, 0.5))

	newChallenge := object_models.NewChallengeNeed{PlayerID: 1001, Amount: 20.01}
	jsonValue, _ := json.Marshal(newChallenge)

	req, _ := http.NewRequest("POST", "/challenges/join", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response object_models.JoinChallengeResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, object_models.Joined, response.Status)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestJoinChallenges_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		setupMock          func(mock sqlmock.Sqlmock)
		inputJSON          string
		expectedStatusCode int
		expectedErrorMsg   string
	}{
		{
			name:               "Invalid JSON input",
			setupMock:          func(mock sqlmock.Sqlmock) {},
			inputJSON:          `{"playerID": "invalid"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorMsg:   "json: cannot unmarshal string into Go struct field",
		},
		{
			name: "GetLastChallengeTime error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnError(sql.ErrConnDone)
			},
			inputJSON:          `{"playerID": 1, "amount": 100}`,
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorMsg:   "sql: connection is already closed",
		},
		{
			name: "Too early for new challenge",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"last_challenge_time"}).AddRow(time.Now()))
			},
			inputJSON:          `{"playerID": 1, "amount": 100}`,
			expectedStatusCode: http.StatusTooEarly,
			expectedErrorMsg:   "You can only participate once per minute",
		},
		{
			name: "Failed to start transaction",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"last_challenge_time"}).AddRow(time.Now().Add(-2 * time.Minute)))
				mock.ExpectBegin().WillReturnError(sql.ErrConnDone)
			},
			inputJSON:          `{"playerID": 1, "amount": 100}`,
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorMsg:   "Failed to start transaction",
		},
		{
			name: "Failed to add new challenge",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"last_challenge_time"}).AddRow(time.Now().Add(-2 * time.Minute)))
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO").WillReturnError(sql.ErrTxDone)
				mock.ExpectRollback()
			},
			inputJSON:          `{"playerID": 1, "amount": 100}`,
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorMsg:   "Failed to Add New Challenge",
		},
		{
			name: "Failed to update price pool",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"last_challenge_time"}).AddRow(time.Now().Add(-2 * time.Minute)))
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE PrizePool").WillReturnError(sql.ErrTxDone)
				mock.ExpectRollback()
			},
			inputJSON:          `{"playerID": 1, "amount": 100}`,
			expectedStatusCode: http.StatusTooEarly,
			expectedErrorMsg:   "Failed to update price pool",
		},
		{
			name: "Failed to commit transaction",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"last_challenge_time"}).AddRow(time.Now().Add(-2 * time.Minute)))
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE PrizePool").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit().WillReturnError(sql.ErrTxDone)
			},
			inputJSON:          `{"playerID": 1, "amount": 100}`,
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorMsg:   "Failed to commit transaction",
		},
		{
			name: "Failed to get challenge",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"last_challenge_time"}).AddRow(time.Now().Add(-2 * time.Minute)))
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE PrizePool").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				mock.ExpectQuery("SELECT").WillReturnError(sql.ErrNoRows)
			},
			inputJSON:          `{"playerID": 1, "amount": 100}`,
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorMsg:   "sql: no rows in result set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mock)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(http.MethodPost, "/join", strings.NewReader(tt.inputJSON))
			c.Request.Header.Set("Content-Type", "application/json")

			object.JoinChallenges(c, db)

			assert.Equal(t, tt.expectedStatusCode, w.Code)

			var response object_models.ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Contains(t, response.Error, tt.expectedErrorMsg)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestShowChallenges(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/challenges", func(c *gin.Context) {
		object.ShowChallenges(c, db)
	})

	rows := sqlmock.NewRows([]string{"ID", "PlayerID", "Amount", "Status", "Won", "CreatedAt", "Probability"}).
		AddRow(1, "1001", 20.01, object_models.Joined, false, time.Now(), 0.5)

	mock.ExpectQuery("SELECT (.+) FROM Challenge").WillReturnRows(rows)

	req, _ := http.NewRequest("GET", "/challenges?limit=1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response object_models.Challenge
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), response.ID)
	assert.Equal(t, "1001", response.PlayerID)

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
			expectedStatusCode: http.StatusOK, // Note: The function doesn't handle this error explicitly
			expectedErrorMsg:   "",            // No error message expected in this case
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
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT Amount FROM PrizePool").WillReturnRows(sqlmock.NewRows([]string{"Amount"}).AddRow(100.0))
	mock.ExpectExec("UPDATE challenges").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("UPDATE PrizePool").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	object.CalculateChallengeResult(db, 1, 1001, 0.005)

	time.Sleep(35 * time.Second) // Wait for the goroutine to complete

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCalculateChallengeResult_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	tests := []struct {
		name          string
		setupMock     func(mock sqlmock.Sqlmock)
		expectedError string
	}{
		{
			name: "Failed to start transaction",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(sql.ErrConnDone)
			},
			expectedError: "Failed to start transaction",
		},
		{
			name: "Failed to distribute prize pool",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE challenges").WillReturnError(sql.ErrTxDone)
				mock.ExpectRollback()
			},
			expectedError: "Failed to distribute prize pool",
		},
		{
			name: "Failed to update probability",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE challenges").WillReturnError(sql.ErrTxDone)
				mock.ExpectRollback()
			},
			expectedError: "Failed to distribute prize pool",
		},
		{
			name: "Failed to commit transaction",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE challenges").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit().WillReturnError(sql.ErrTxDone)
				mock.ExpectRollback()
			},
			expectedError: "Failed to commit transaction",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mock)

			// Capture log output
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer func() { log.SetOutput(os.Stderr) }()

			object.CalculateChallengeResult(db, 1, 1, 0.5)

			logOutput := buf.String()
			assert.Contains(t, logOutput, tt.expectedError)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
