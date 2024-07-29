package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	object "github.com/RYANCOAL9999/SpinnrTechnologyInterview/gameRoomManagementSystem/handlers"
	object_models "github.com/RYANCOAL9999/SpinnrTechnologyInterview/gameRoomManagementSystem/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetRooms(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"ID", "Name", "Status", "Description", "PlayerIDs"}).
		AddRow(1, "Room 1", object_models.StatusAvailable, "Description 1", "1,2,3").
		AddRow(2, "Room 2", object_models.StatusOccupied, "Description 2", "4,5,6")

	mock.ExpectQuery("SELECT (.+) FROM Room").WillReturnRows(rows)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	object.GetRooms(c, db)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []object_models.Room
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, "Room 1", response[0].Name)
	assert.Equal(t, "Room 2", response[1].Name)
}

func TestGetRooms_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectQuery("SELECT (.+) FROM rooms").WillReturnError(errors.New("database error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	object.GetRooms(c, db)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "database error")
}

func TestCreateRoom(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO Room").WithArgs("New Room", "New Description").WillReturnResult(sqlmock.NewResult(1, 1))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	room := object_models.Room{Name: "New Room", Description: "New Description"}
	body, _ := json.Marshal(room)
	c.Request, _ = http.NewRequest("POST", "/rooms", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	object.CreateRoom(c, db)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response object_models.CreateResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 1, response.ID)
}

func TestCreateRoom_Error(t *testing.T) {
	tests := []struct {
		name           string
		inputJSON      string
		mockSetup      func(sqlmock.Sqlmock)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Invalid JSON",
			inputJSON:      `{"name": 123, "description": "test"}`,
			mockSetup:      func(mock sqlmock.Sqlmock) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"json: cannot unmarshal number into Go struct field Room.name of type string"}`,
		},
		{
			name:      "Database Error",
			inputJSON: `{"name": "Test Room", "description": "A test room"}`,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO rooms").WillReturnError(errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"database error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			tt.mockSetup(mock)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("POST", "/rooms", strings.NewReader(tt.inputJSON))
			c.Request.Header.Set("Content-Type", "application/json")

			object.CreateRoom(c, db)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestGetRoom(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"ID", "Name", "Status", "Description", "PlayerIDs"}).
		AddRow(1, "Room 1", object_models.StatusAvailable, "Description 1", "1,2,3")

	mock.ExpectQuery("SELECT (.+) FROM Room WHERE ID = ?").WithArgs(1).WillReturnRows(rows)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	object.GetRoom(c, db)

	assert.Equal(t, http.StatusOK, w.Code)

	var response object_models.Room
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Room 1", response.Name)
}

func TestGetRoom_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectQuery("SELECT (.+) FROM rooms WHERE (.+)").WillReturnError(sql.ErrNoRows)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	object.GetRoom(c, db)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "sql: no rows in result set")
}

func TestUpdateRoom(t *testing.T) {

}

func TestUpdateRoom_Error(t *testing.T) {

}

func TestDeleteRoom(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("UPDATE Room SET").WithArgs("Updated Room", object_models.StatusOccupied, "Updated Description", "1,2,3", 1).WillReturnResult(sqlmock.NewResult(1, 1))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	room := object_models.Room{ID: 1, Name: "Updated Room", Status: object_models.StatusOccupied, Description: "Updated Description", PlayerIDs: "1,2,3"}
	body, _ := json.Marshal(room)
	c.Request, _ = http.NewRequest("PUT", "/rooms", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	object.UpdateRoom(c, db)

	assert.Equal(t, http.StatusOK, w.Code)

	var response object_models.SuccessResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
}

func TestDeleteRoom_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectExec("DELETE FROM rooms WHERE (.+)").WillReturnError(errors.New("database error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	object.DeleteRoom(c, db)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "database error")
}
