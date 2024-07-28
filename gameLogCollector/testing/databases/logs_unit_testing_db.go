package databases

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	object "github.com/RYANCOAL9999/SpinnrTechnologyInterview/gameLogCollector/databases"
	object_models "github.com/RYANCOAL9999/SpinnrTechnologyInterview/gameLogCollector/models"
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
