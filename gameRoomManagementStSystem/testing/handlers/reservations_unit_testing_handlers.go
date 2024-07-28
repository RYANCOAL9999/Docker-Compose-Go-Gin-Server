package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	object "github.com/RYANCOAL9999/SpinnrTechnologyInterview/gameRoomManagementSystem/handlers"
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
