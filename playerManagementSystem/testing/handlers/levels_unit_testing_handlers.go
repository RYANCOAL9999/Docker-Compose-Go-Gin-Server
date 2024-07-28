package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	object "github.com/RYANCOAL9999/SpinnrTechnologyInterview/playerManagementSystem/handlers"
	object_models "github.com/RYANCOAL9999/SpinnrTechnologyInterview/playerManagementSystem/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetLevels(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"ID", "Name", "LV"}).
		AddRow(1, "Novice", 1).
		AddRow(2, "Expert", 10)

	mock.ExpectQuery("SELECT ID, Name, LV FROM Level").WillReturnRows(rows)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	object.GetLevels(c, db)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []object_models.Level
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, object_models.Level{ID: 1, Name: "Novice", LV: 1}, response[0])
	assert.Equal(t, object_models.Level{ID: 2, Name: "Expert", LV: 10}, response[1])
}

func TestCreateLevel(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO Level").
		WithArgs("Intermediate", 5).
		WillReturnResult(sqlmock.NewResult(3, 1))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	newLevel := object_models.Level{Name: "Intermediate", LV: 5}
	body, _ := json.Marshal(newLevel)
	c.Request, _ = http.NewRequest("POST", "/levels", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	object.CreateLevel(c, db)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response object_models.CreateResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 3, response.ID)
}
