package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	object "github.com/RYANCOAL9999/SpinnrTechnologyInterview/gameLogCollector/handlers"
	object_models "github.com/RYANCOAL9999/SpinnrTechnologyInterview/gameLogCollector/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetLogs(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/game_logs", func(c *gin.Context) {
		object.GetLogs(c, db)
	})

	rows := sqlmock.NewRows([]string{"ID", "PlayerID", "Action", "Timestamp", "Details"}).
		AddRow(1, 1001, "LOGIN", time.Now(), "User logged in")

	mock.ExpectQuery("SELECT (.+) FROM GameLog").WillReturnRows(rows)

	req, _ := http.NewRequest("GET", "/game_logs?player_id=1001", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response object_models.GameLog
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), response.ID)
	assert.Equal(t, 1001, response.PlayerID)
	assert.Equal(t, "LOGIN", response.Action)
}

func TestGetLogs_Error(t *testing.T) {

}

func TestCreateLog(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/game_logs", func(c *gin.Context) {
		object.CreateLog(c, db)
	})

	mock.ExpectExec("INSERT INTO GameLog").
		WithArgs(1001, "LOGIN", sqlmock.AnyArg(), "User logged in").
		WillReturnResult(sqlmock.NewResult(1, 1))

	newLog := object_models.GameLog{
		PlayerID: 1001,
		Action:   "LOGIN",
		Details:  "User logged in",
	}
	jsonValue, _ := json.Marshal(newLog)

	req, _ := http.NewRequest("POST", "/game_logs", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response object_models.CreateResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 1, response.ID)
}

func TestCreateLog_Error(t *testing.T) {

}
