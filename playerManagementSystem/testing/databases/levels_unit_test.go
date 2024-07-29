package databases

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	object "github.com/RYANCOAL9999/SpinnrTechnologyInterview/playerManagementSystem/databases"
	object_models "github.com/RYANCOAL9999/SpinnrTechnologyInterview/playerManagementSystem/models"
	"github.com/stretchr/testify/assert"
)

func TestGetLevelsData(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"ID", "Name", "LV"}).
		AddRow(1, "Novice", 1).
		AddRow(2, "Expert", 10)

	mock.ExpectQuery("SELECT ID, Name, LV FROM Level").WillReturnRows(rows)

	levels, err := object.GetLevelsData(db)

	assert.NoError(t, err)
	assert.Len(t, levels, 2)
	assert.Equal(t, object_models.Level{ID: 1, Name: "Novice", LV: 1}, levels[0])
	assert.Equal(t, object_models.Level{ID: 2, Name: "Expert", LV: 10}, levels[1])

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestGetLevelsData_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT ID, Name, LV FROM Level").WillReturnError(sql.ErrConnDone)

	levels, err := object.GetLevelsData(db)

	assert.Error(t, err)
	assert.Nil(t, levels)
	assert.Contains(t, err.Error(), "error querying database with GetLevelsData")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestAddLevel(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO Level").
		WithArgs("Intermediate", 5).
		WillReturnResult(sqlmock.NewResult(3, 1))

	id, err := object.AddLevel(db, "Intermediate", 5)

	assert.NoError(t, err)
	assert.Equal(t, 3, id)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestAddLevel_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO Level").
		WithArgs("Intermediate", 5).
		WillReturnError(sql.ErrTxDone)

	id, err := object.AddLevel(db, "Intermediate", 5)

	assert.Error(t, err)
	assert.Equal(t, 0, id)
	assert.Contains(t, err.Error(), "error querying database with AddLevel")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}
