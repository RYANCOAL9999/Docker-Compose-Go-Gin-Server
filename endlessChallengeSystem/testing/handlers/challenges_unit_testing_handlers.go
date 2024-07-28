package handlers

import (
	"bytes"
	"encoding/json"
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
