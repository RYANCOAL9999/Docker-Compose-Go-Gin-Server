package databases

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	object "github.com/RYANCOAL9999/SpinnrTechnologyInterview/playerManagementSystem/databases"
	object_models "github.com/RYANCOAL9999/SpinnrTechnologyInterview/playerManagementSystem/models"
	"github.com/stretchr/testify/assert"
)

func TestGetPlayersData(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"ID", "Name", "LV"}).
		AddRow(1, "Alice", 5).
		AddRow(2, "Bob", 10)

	mock.ExpectQuery("SELECT P.ID as ID, P.Name as Name, L.LV as LV FROM Player P").
		WillReturnRows(rows)

	players, err := object.GetPlayersData(db)

	assert.NoError(t, err)
	assert.Len(t, players, 2)
	assert.Equal(t, object_models.PlayerRank{ID: 1, Name: "Alice", LV: 5}, players[0])
	assert.Equal(t, object_models.PlayerRank{ID: 2, Name: "Bob", LV: 10}, players[1])

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestGetPlayersData_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	t.Run("Database query error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM Player P INNER JOIN Level L (.+)").
			WillReturnError(sql.ErrConnDone)

		playerRanks, err := object.GetPlayersData(db)

		assert.Nil(t, playerRanks)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error querying database with GetPlayersData")
	})

	t.Run("Row scan error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"ID", "Name", "LV"}).
			AddRow("invalid", "Player1", 10) // ID should be an integer, not a string

		mock.ExpectQuery("SELECT (.+) FROM Player P INNER JOIN Level L (.+)").
			WillReturnRows(rows)

		playerRanks, err := object.GetPlayersData(db)

		assert.Nil(t, playerRanks)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error scanning row with GetPlayersData")
	})

	t.Run("No rows returned", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"ID", "Name", "LV"})

		mock.ExpectQuery("SELECT (.+) FROM Player P INNER JOIN Level L (.+)").
			WillReturnRows(rows)

		playerRanks, err := object.GetPlayersData(db)

		assert.NoError(t, err)
		assert.Empty(t, playerRanks)
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}

}

func TestAddPlayer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO Player").
		WithArgs("Charlie", 3).
		WillReturnResult(sqlmock.NewResult(3, 1))

	id, err := object.AddPlayer(db, "Charlie", 3)

	assert.NoError(t, err)
	assert.Equal(t, 3, id)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestAddPlayer_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	name := "JohnDoe"
	level := 0

	mock.ExpectExec("INSERT INTO Player").
		WithArgs(name, level).
		WillReturnResult(sqlmock.NewResult(1, 1))

	id, err := object.AddPlayer(db, name, level)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if id != 0 {
		t.Errorf("expected id to be 0, got %d", id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestGetPlayer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"ID", "Name", "LV"}).
		AddRow(1, "Alice", 5)

	mock.ExpectQuery("SELECT P.ID as ID, P.Name as Name, L.LV as LV FROM Player P").
		WithArgs(1).
		WillReturnRows(rows)

	player, err := object.GetPlayer(db, 1)

	assert.NoError(t, err)
	assert.NotNil(t, player)
	assert.Equal(t, object_models.PlayerRank{ID: 1, Name: "Alice", LV: 5}, *player)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestGetPlayer_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT P.ID as ID, P.Name as Name, L.LV as LV FROM Player P INNER JOIN Levels L ON P.LevelID = L.ID WHERE P.ID = ?").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)

	playerRank, err := object.GetPlayer(db, 1)
	if playerRank != nil {
		t.Errorf("expected nil playerRank, got %v", playerRank)
	}

	if err == nil || !strings.Contains(err.Error(), "error querying database with GetPlayer") {
		t.Errorf("expected error querying database with GetPlayer, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdatePlayer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT ID FROM Level WHERE LV = ?").
		WithArgs(6).
		WillReturnRows(sqlmock.NewRows([]string{"ID"}).AddRow(2))

	mock.ExpectExec("UPDATE players SET").
		WithArgs(2, "Alice Updated", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = object.UpdatePlayer(db, object_models.PlayerRank{ID: 1, Name: "Alice Updated", LV: 6})

	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestUpdatePlayer_Error(t *testing.T) {
	// Mock SQL database
	mockDB := &sql.DB{}

	// Prepare test data
	player := object_models.PlayerRank{
		ID: 1,
	}

	// Call the function
	err := object.UpdatePlayer(mockDB, player)

	// Check for error
	if err == nil {
		t.Error("Expected an error but got nil")
	}
}

func TestDeletePlayer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectExec("DELETE FROM players WHERE ID = ?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = object.DeletePlayer(db, 1)

	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestDeletePlayer_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("DELETE FROM players WHERE ID = ?").
		WithArgs(1).
		WillReturnError(fmt.Errorf("query failed"))

	err = object.DeletePlayer(db, 1)
	if err == nil || err.Error() != "error querying database with DeletePlayer: query failed" {
		t.Errorf("expected error 'error querying database with DeletePlayer: query failed', but got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
