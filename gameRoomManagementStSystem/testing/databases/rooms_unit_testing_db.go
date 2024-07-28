package databases

import (
	"database/sql"
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	object "github.com/RYANCOAL9999/SpinnrTechnologyInterview/gameRoomManagementSystem/databases"
	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/gameRoomManagementSystem/models"
	object_models "github.com/RYANCOAL9999/SpinnrTechnologyInterview/gameRoomManagementSystem/models"
	"github.com/stretchr/testify/assert"
)

func TestListRooms(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"ID", "Name", "Status"}).
		AddRow(1, "Room 1", 0).
		AddRow(2, "Room 2", 1).
		AddRow(3, "Room 3", 2)

	mock.ExpectQuery("SELECT ID, Name, Status FROM Room").WillReturnRows(rows)

	rooms, err := object.ListRooms(db)

	assert.NoError(t, err)
	assert.Len(t, rooms, 3)
	assert.Equal(t, 1, rooms[0].ID)
	assert.Equal(t, "Room 1", rooms[0].Name)
	assert.Equal(t, 0, int(rooms[0].Status))
	assert.Equal(t, 2, rooms[1].ID)
	assert.Equal(t, "Room 2", rooms[1].Name)
	assert.Equal(t, 1, int(rooms[1].Status))
	assert.Equal(t, 3, rooms[2].ID)
	assert.Equal(t, "Room 3", rooms[2].Name)
	assert.Equal(t, 2, int(rooms[2].Status))

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestListRooms_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT ID, Name, Status FROM Room").WillReturnError(sqlmock.ErrCancelled)

	rooms, err := object.ListRooms(db)

	assert.Error(t, err)
	assert.Nil(t, rooms)
	assert.Contains(t, err.Error(), "error querying database with ListRooms")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestShowRoom(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	tests := []struct {
		name      string
		roomID    int
		mockRow   *sqlmock.Rows
		expectErr bool
	}{
		{
			name:   "Existing room",
			roomID: 1,
			mockRow: sqlmock.NewRows([]string{"ID", "Name", "Status"}).
				AddRow(1, "Test Room", 0),
			expectErr: false,
		},
		{
			name:      "Non-existent room",
			roomID:    999,
			mockRow:   sqlmock.NewRows([]string{"ID", "Name", "Status"}),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery("SELECT ID, Name, Status FROM Room WHERE ID = ?").
				WithArgs(tt.roomID).
				WillReturnRows(tt.mockRow)

			room, err := object.ShowRoom(db, tt.roomID)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, room)
				assert.Contains(t, err.Error(), "error querying database with Show Room")
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, room)
				assert.Equal(t, tt.roomID, room.ID)
				assert.Equal(t, "Test Room", room.Name)
				assert.Equal(t, 0, int(room.Status))
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestShowRoom_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT ID, Name, Status FROM Room WHERE ID = ?").
		WithArgs(1).
		WillReturnError(sql.ErrConnDone)

	room, err := object.ShowRoom(db, 1)

	assert.Error(t, err)
	assert.Nil(t, room)
	assert.Contains(t, err.Error(), "error scanning row with Show Room")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestAddRoom(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	tests := []struct {
		name        string
		roomName    string
		description string
		mockResult  sql.Result
		expectErr   bool
	}{
		{
			name:        "Successful room addition",
			roomName:    "New Room",
			description: "A brand new room",
			mockResult:  sqlmock.NewResult(1, 1),
			expectErr:   false,
		},
		{
			name:        "Database error",
			roomName:    "Error Room",
			description: "Room causing error",
			mockResult:  nil,
			expectErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectErr {
				mock.ExpectExec("INSERT INTO Room").
					WithArgs(tt.roomName, tt.description).
					WillReturnError(sql.ErrConnDone)
			} else {
				mock.ExpectExec("INSERT INTO Room").
					WithArgs(tt.roomName, tt.description).
					WillReturnResult(tt.mockResult)
			}

			id, err := object.AddRoom(db, tt.roomName, tt.description)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, 0, id)
				assert.Contains(t, err.Error(), "error querying database with AddRoom")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, 1, id)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestAddRoomWithInvalidInput(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	tests := []struct {
		name        string
		roomName    string
		description string
	}{
		{
			name:        "Empty room name",
			roomName:    "",
			description: "Valid description",
		},
		{
			name:        "Empty description",
			roomName:    "Valid Name",
			description: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectExec("INSERT INTO Room").
				WithArgs(tt.roomName, tt.description).
				WillReturnResult(sqlmock.NewResult(0, 0))

			id, err := object.AddRoom(db, tt.roomName, tt.description)

			assert.Equal(t, 0, id)
			assert.NoError(t, err, "AddRoom should not return an error for invalid input")

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpdateRoomData(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	testCases := []struct {
		name        string
		room        object_models.Room
		expectQuery string
		expectArgs  []interface{}
		mockResult  sql.Result
		expectError bool
	}{
		{
			name: "Update all fields",
			room: object_models.Room{
				ID:          1,
				Name:        "Updated Room",
				Status:      object_models.StatusOccupied,
				Description: "New description",
				PlayerIDs:   "1,2,3",
			},
			expectQuery: "UPDATE Room SET Name = ?, Status = ?, Description = ?, PlayerIDs = ? WHERE id = ?",
			expectArgs:  []interface{}{"Updated Room", object_models.StatusOccupied, "New description", "1,2,3", 1},
			mockResult:  sqlmock.NewResult(1, 1),
			expectError: false,
		},
		{
			name: "Update partial fields",
			room: object_models.Room{
				ID:     2,
				Name:   "Partial Update",
				Status: object_models.StatusMaintenance,
			},
			expectQuery: "UPDATE Room SET Name = ?, Status = ? WHERE id = ?",
			expectArgs:  []interface{}{"Partial Update", object_models.StatusMaintenance, 2},
			mockResult:  sqlmock.NewResult(2, 1),
			expectError: false,
		},
		{
			name:        "No fields to update",
			room:        object_models.Room{ID: 3},
			expectError: true,
		},
		{
			name: "Non-existent room",
			room: object_models.Room{
				ID:   4,
				Name: "Non-existent Room",
			},
			expectQuery: "UPDATE Room SET Name = ? WHERE id = ?",
			expectArgs:  []interface{}{"Non-existent Room", 4},
			mockResult:  sqlmock.NewResult(0, 0),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.expectError {

				driverArgs := make([]driver.Value, len(tc.expectArgs))
				for i, arg := range tc.expectArgs {
					driverArgs[i] = arg.(driver.Value)
				}

				mock.ExpectExec(tc.expectQuery).
					WithArgs(driverArgs...).
					WillReturnResult(tc.mockResult)
			}

			err := object.UpdateRoomData(db, tc.room)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpdateRoomData_Error(t *testing.T) {

}

func TestDeleteRoom(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	tests := []struct {
		name           string
		roomID         int
		mockResult     sql.Result
		expectError    bool
		expectedErrMsg string
	}{
		{
			name:           "Successful deletion",
			roomID:         1,
			mockResult:     sqlmock.NewResult(0, 1),
			expectError:    false,
			expectedErrMsg: "",
		},
		{
			name:           "Room not found",
			roomID:         999,
			mockResult:     sqlmock.NewResult(0, 0),
			expectError:    true,
			expectedErrMsg: "room not found",
		},
		{
			name:           "Database error",
			roomID:         2,
			mockResult:     nil,
			expectError:    true,
			expectedErrMsg: "error querying database with DeleteRoom",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockResult == nil {
				mock.ExpectExec("DELETE FROM Room WHERE ID = ?").
					WithArgs(tt.roomID).
					WillReturnError(sql.ErrConnDone)
			} else {
				mock.ExpectExec("DELETE FROM Room WHERE ID = ?").
					WithArgs(tt.roomID).
					WillReturnResult(tt.mockResult)
			}

			err := object.DeleteRoom(db, tt.roomID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrMsg)
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDeleteRoom_Error(t *testing.T) {

}

func TestSearchPlayerInRoom(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	tests := []struct {
		name            string
		playerIDs       []int
		mockRows        *sqlmock.Rows
		expectError     bool
		expectedErrMsg  string
		expectedPlayers []models.PlayerRank
	}{
		{
			name:      "All players found",
			playerIDs: []int{1, 2, 3},
			mockRows: sqlmock.NewRows([]string{"ID", "Name", "LV"}).
				AddRow(1, "Player1", 5).
				AddRow(2, "Player2", 3).
				AddRow(3, "Player3", 7),
			expectError: false,
			expectedPlayers: []models.PlayerRank{
				{ID: 1, Name: "Player1", LV: 5},
				{ID: 2, Name: "Player2", LV: 3},
				{ID: 3, Name: "Player3", LV: 7},
			},
		},
		{
			name:      "Some players not found",
			playerIDs: []int{1, 2, 3},
			mockRows: sqlmock.NewRows([]string{"ID", "Name", "LV"}).
				AddRow(1, "Player1", 5).
				AddRow(3, "Player3", 7),
			expectError:    true,
			expectedErrMsg: "some players were not found",
			expectedPlayers: []models.PlayerRank{
				{ID: 1, Name: "Player1", LV: 5},
				{ID: 3, Name: "Player3", LV: 7},
			},
		},
		{
			name:           "Database error",
			playerIDs:      []int{1, 2},
			mockRows:       nil,
			expectError:    true,
			expectedErrMsg: "error querying database with searchPlayerInRoom",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queryRegex := "SELECT P.ID, P.Name, L.LV FROM Player P INNER JOIN Level L ON P.LevelID = L.ID WHERE P.ID IN \\(\\?,\\?.*\\) ORDER BY P.ID"
			if tt.mockRows == nil {
				mock.ExpectQuery(queryRegex).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(sql.ErrConnDone)
			} else {
				mock.ExpectQuery(queryRegex).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(tt.mockRows)
			}

			players, err := object.SearchPlayerInRoom(db, tt.playerIDs)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrMsg)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedPlayers, players)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestSearchPlayerInRoomWithEmptyInput(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	players, err := object.SearchPlayerInRoom(db, []int{})

	assert.NoError(t, err)
	assert.Empty(t, players)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}
