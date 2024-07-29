package databases

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	object "github.com/RYANCOAL9999/SpinnrTechnologyInterview/gameLogCollector/databases"
	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/gameLogCollector/models"
	object_models "github.com/RYANCOAL9999/SpinnrTechnologyInterview/gameLogCollector/models"
	"github.com/stretchr/testify/assert"
)

func TestListLogs(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"ID", "PlayerID", "Action", "Timestamp", "Details"}).
		AddRow(1, 1001, "LOGIN", time.Now(), "User logged in").
		AddRow(2, 1001, "PURCHASE", time.Now(), "User made a purchase")

	mock.ExpectQuery("SELECT (.+) FROM GameLog").WillReturnRows(rows)

	logs, err := object.ListLogs(db, 1001, "", time.Time{}, time.Time{}, 0)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(logs) != 2 {
		t.Errorf("Expected 2 logs, got %d", len(logs))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestListLogs_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	tests := []struct {
		name        string
		setupMock   func(mock sqlmock.Sqlmock)
		playerID    int
		action      string
		startTime   time.Time
		endTime     time.Time
		limit       int
		expectedErr string
	}{
		{
			name: "Database query error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM GameLog").WillReturnError(sql.ErrConnDone)
			},
			playerID:    1,
			action:      "login",
			startTime:   time.Now().Add(-24 * time.Hour),
			endTime:     time.Now(),
			limit:       10,
			expectedErr: "error querying database with ListLogs: sql: connection is already closed",
		},
		{
			name: "Row scan error",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"ID", "PlayerID", "Action", "Timestamp", "Details"}).
					AddRow("invalid", 1, "login", time.Now(), "details") // ID should be int, not string
				mock.ExpectQuery("SELECT (.+) FROM GameLog").WillReturnRows(rows)
			},
			playerID:    1,
			action:      "login",
			startTime:   time.Now().Add(-24 * time.Hour),
			endTime:     time.Now(),
			limit:       10,
			expectedErr: "error scanning row with with ListLogs: sql: Scan error on column index 0, name \"ID\": converting driver.Value type string (\"invalid\") to a int: invalid syntax",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mock)

			_, err := object.ListLogs(db, tt.playerID, tt.action, tt.startTime, tt.endTime, tt.limit)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErr)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestAddLog(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO GameLog").
		WithArgs(1001, "LOGIN", sqlmock.AnyArg(), "User logged in").
		WillReturnResult(sqlmock.NewResult(1, 1))

	log := object_models.GameLog{
		PlayerID: 1001,
		Action:   "LOGIN",
		Details:  "User logged in",
	}

	id, err := object.AddLog(db, log)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if id != 1 {
		t.Errorf("Expected ID 1, got %d", id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestAddLog_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	tests := []struct {
		name        string
		setupMock   func(mock sqlmock.Sqlmock)
		input       object_models.GameLog
		expectedID  int
		expectedErr string
	}{
		{
			name: "Database execution error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO GameLog").WillReturnError(sql.ErrConnDone)
			},
			input: models.GameLog{
				PlayerID: 1,
				Action:   "login",
				Details:  "User logged in",
			},
			expectedID:  0,
			expectedErr: "error querying database with AddLog: sql: connection is already closed",
		},
		{
			name: "LastInsertId error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO GameLog").WillReturnResult(sqlmock.NewErrorResult(errors.New("LastInsertId error")))
			},
			input: models.GameLog{
				PlayerID: 1,
				Action:   "logout",
				Details:  "User logged out",
			},
			expectedID:  0,
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mock)

			id, err := object.AddLog(db, tt.input)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedID, id)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
