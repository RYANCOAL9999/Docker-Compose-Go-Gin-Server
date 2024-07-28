package databases

import (
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

}
