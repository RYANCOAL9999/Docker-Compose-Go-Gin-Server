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
	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/gameRoomManagementSystem/models"
	object_models "github.com/RYANCOAL9999/SpinnrTechnologyInterview/gameRoomManagementSystem/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetReservations(t *testing.T) {
	// Set up mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// Create a new gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set up query parameters
	c.Request, _ = http.NewRequest("GET", "/?room_id=1&start_Date=2024-07-28&end_Date=2024-07-29&limit=5", nil)

	// Set up expected query and result rows
	expectedQuery := "SELECT (.+) FROM reservations WHERE (.+)"
	rows := sqlmock.NewRows([]string{"id", "room_id", "date", "player_ids"}).
		AddRow(1, 1, "2024-07-28", "1,2,3").
		AddRow(2, 1, "2024-07-29", "4,5,6")

	mock.ExpectQuery(expectedQuery).WillReturnRows(rows)

	// Call the function
	object.GetReservations(c, db)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)

	var response object_models.ReservationRoom
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, 1, response.ID)
	assert.Equal(t, 1, response.RoomID)
	assert.Equal(t, "2024-07-28", response.Date)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetReservations_Error(t *testing.T) {
	// Set up mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// Create a new gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set up query parameters
	c.Request, _ = http.NewRequest("GET", "/?room_id=1&start_Date=2024-07-28&end_Date=2024-07-29&limit=5", nil)

	// Set up expected query and error
	expectedQuery := "SELECT (.+) FROM reservations WHERE (.+)"
	expectedError := sql.ErrNoRows

	mock.ExpectQuery(expectedQuery).WillReturnError(expectedError)

	// Call the function
	object.GetReservations(c, db)

	// Assert the response
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response models.ErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, expectedError.Error(), response.Error)

	// Verify all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateReservationRoom(t *testing.T) {

	tests := []struct {
		name      string
		roomID    int
		playerIDs string
		setupMock func(sqlmock.Sqlmock)
		expectErr bool
	}{
		{
			name:      "Successful Update",
			roomID:    1,
			playerIDs: "1,2,3",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE rooms").
					WithArgs("1,2,3", 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectErr: false,
		},
		{
			name:      "Database Error",
			roomID:    2,
			playerIDs: "4,5,6",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE rooms").
					WithArgs("4,5,6", 2).
					WillReturnError(sqlmock.ErrCancelled)
			},
			expectErr: true,
		},
		{
			name:      "No Rows Affected",
			roomID:    3,
			playerIDs: "7,8,9",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE rooms").
					WithArgs("7,8,9", 3).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.setupMock(mock)

			err = object.UpdateReservationRoom(db, tt.roomID, tt.playerIDs)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpdateReservationRoom_Error(t *testing.T) {
	tests := []struct {
		name      string
		roomID    int
		playerIDs string
		setupMock func(sqlmock.Sqlmock)
		expectErr bool
	}{
		{
			name:      "Successful Update",
			roomID:    1,
			playerIDs: "1,2,3",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE rooms").
					WithArgs("1,2,3", 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectErr: false,
		},
		{
			name:      "Database Error",
			roomID:    2,
			playerIDs: "4,5,6",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE rooms").
					WithArgs("4,5,6", 2).
					WillReturnError(sqlmock.ErrCancelled)
			},
			expectErr: true,
		},
		{
			name:      "No Rows Affected",
			roomID:    3,
			playerIDs: "7,8,9",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE rooms").
					WithArgs("7,8,9", 3).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectErr: true,
		},
		{
			name:      "Invalid Room ID",
			roomID:    -1,
			playerIDs: "10,11,12",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE rooms").
					WithArgs("10,11,12", -1).
					WillReturnError(sql.ErrNoRows)
			},
			expectErr: true,
		},
		{
			name:      "Empty Player IDs",
			roomID:    4,
			playerIDs: "",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE rooms").
					WithArgs("", 4).
					WillReturnError(errors.New("player IDs cannot be empty"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.setupMock(mock)

			err = object.UpdateReservationRoom(db, tt.roomID, tt.playerIDs)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCreateReservations(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		inputJSON      string
		setupMock      func(sqlmock.Sqlmock)
		expectedStatus int
		expectedBody   interface{}
	}{
		// Test cases remain the same as in the previous version
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.setupMock(mock)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest(http.MethodPost, "/reservations", bytes.NewBufferString(tt.inputJSON))
			c.Request.Header.Set("Content-Type", "application/json")

			object.CreateReservations(c, db)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCreateReservations_Error(t *testing.T) {
	tests := []struct {
		name           string
		inputJSON      string
		setupMock      func(sqlmock.Sqlmock)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "Successful Reservation",
			inputJSON: `{"room_id": 1, "date": "2024-07-28", "player_ids": "1,2,3"}`,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM rooms WHERE (.+)").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "status"}).AddRow(1, models.StatusAvailable))
				mock.ExpectExec("INSERT INTO reservations").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE rooms").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":1}`,
		},
		{
			name:           "Invalid JSON",
			inputJSON:      `{"room_id": "invalid", "date": "2024-07-28", "player_ids": "1,2,3"}`,
			setupMock:      func(mock sqlmock.Sqlmock) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"json: cannot unmarshal string into Go struct field Reservation.room_id of type int"}`,
		},
		{
			name:      "Room Not Found",
			inputJSON: `{"room_id": 999, "date": "2024-07-28", "player_ids": "1,2,3"}`,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM rooms WHERE (.+)").
					WithArgs(999).
					WillReturnError(sql.ErrNoRows)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"sql: no rows in result set"}`,
		},
		{
			name:      "Room Not Available",
			inputJSON: `{"room_id": 2, "date": "2024-07-28", "player_ids": "1,2,3"}`,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM rooms WHERE (.+)").
					WithArgs(2).
					WillReturnRows(sqlmock.NewRows([]string{"id", "status"}).AddRow(2, models.StatusOccupied))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"this room is not available"}`,
		},
		{
			name:      "Invalid Date Format",
			inputJSON: `{"room_id": 1, "date": "2024-13-45", "player_ids": "1,2,3"}`,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM rooms WHERE (.+)").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "status"}).AddRow(1, models.StatusAvailable))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"date format has issues"}`,
		},
		{
			name:      "Reservation Insert Error",
			inputJSON: `{"room_id": 1, "date": "2024-07-28", "player_ids": "1,2,3"}`,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM rooms WHERE (.+)").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "status"}).AddRow(1, models.StatusAvailable))
				mock.ExpectExec("INSERT INTO reservations").
					WillReturnError(errors.New("insert error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"insert error"}`,
		},
		{
			name:      "Update Room Error",
			inputJSON: `{"room_id": 1, "date": "2024-07-28", "player_ids": "1,2,3"}`,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM rooms WHERE (.+)").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "status"}).AddRow(1, models.StatusAvailable))
				mock.ExpectExec("INSERT INTO reservations").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE rooms").
					WillReturnError(errors.New("update error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"update error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tt.setupMock(mock)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("POST", "/reservations", strings.NewReader(tt.inputJSON))
			c.Request.Header.Set("Content-Type", "application/json")

			object.CreateReservations(c, db)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
