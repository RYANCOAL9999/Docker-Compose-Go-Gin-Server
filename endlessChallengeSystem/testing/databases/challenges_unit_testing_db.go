package databases

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	object "github.com/RYANCOAL9999/SpinnrTechnologyInterview/endlessChallengeSystem/databases"
	object_models "github.com/RYANCOAL9999/SpinnrTechnologyInterview/endlessChallengeSystem/models"
	"github.com/stretchr/testify/assert"
)

func TestListChallenges(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"ID", "PlayerID", "Amount", "Status", "Won", "CreatedAt", "Probability"}).
		AddRow(1, "1001", 20.01, object_models.Joined, false, time.Now(), 0.5).
		AddRow(2, "1002", 20.01, object_models.Ready, false, time.Now(), 0.7)

	mock.ExpectQuery("SELECT (.+) FROM Challenge").WillReturnRows(rows)

	challenges, err := object.ListChallenges(db, 0)
	assert.NoError(t, err)
	assert.Len(t, challenges, 2)
	assert.Equal(t, int64(1), challenges[0].ID)
	assert.Equal(t, "1001", challenges[0].PlayerID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetLastChallengeTime(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	expectedTime := time.Now()
	rows := sqlmock.NewRows([]string{"CreatedAt"}).AddRow(expectedTime)

	mock.ExpectQuery("SELECT CreatedAt FROM Challenge").WithArgs(1001).WillReturnRows(rows)

	lastTime, err := object.GetLastChallengeTime(db, 1001)
	assert.NoError(t, err)
	assert.Equal(t, expectedTime, *lastTime)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddNewChallenge(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO Challenge").WithArgs(1001, 20.01).WillReturnResult(sqlmock.NewResult(1, 1))

	tx, _ := db.Begin()
	id, err := object.AddNewChallenge(tx, object_models.NewChallengeNeed{PlayerID: 1001, Amount: 20.01})
	assert.NoError(t, err)
	assert.Equal(t, 1, id)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdatePricePool(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE PrizePool").WithArgs(20.01).WillReturnResult(sqlmock.NewResult(1, 1))

	tx, _ := db.Begin()
	err = object.UpdatePricePool(tx, 20.01)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDistributePrizePool(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT Amount FROM PrizePool").WillReturnRows(sqlmock.NewRows([]string{"Amount"}).AddRow(100.0))
	mock.ExpectExec("UPDATE challenges").WithArgs(1, 1001).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("UPDATE PrizePool").WillReturnResult(sqlmock.NewResult(1, 1))

	tx, _ := db.Begin()
	err = object.DistributePrizePool(tx, 1, 1001)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
