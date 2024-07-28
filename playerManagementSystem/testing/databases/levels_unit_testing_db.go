package databases

import (
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
