package databases

import (
	"database/sql"
	"errors"
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

func TestListChallenges_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	t.Run("Database query error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM Challenge").WillReturnError(sql.ErrConnDone)

		challenges, err := object.ListChallenges(db, 0)
		assert.Error(t, err)
		assert.Nil(t, challenges)
		assert.Contains(t, err.Error(), "error querying database with ListChallenges")
	})

	t.Run("Row scan error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"ID", "PlayerID", "Amount", "Status", "Won", "CreatedAt", "Probability"}).
			AddRow("invalid", "1001", 20.01, "Joined", false, "invalid_time", 0.5)

		mock.ExpectQuery("SELECT (.+) FROM Challenge").WillReturnRows(rows)

		challenges, err := object.ListChallenges(db, 0)
		assert.Error(t, err)
		assert.Nil(t, challenges)
		assert.Contains(t, err.Error(), "error scanning row with ListChallenges")
	})

	t.Run("Invalid limit parameter", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM Challenge LIMIT ?").
			WithArgs(-1).
			WillReturnError(errors.New("invalid LIMIT clause"))

		challenges, err := object.ListChallenges(db, -1)
		assert.Error(t, err)
		assert.Nil(t, challenges)
		assert.Contains(t, err.Error(), "error querying database with ListChallenges")
	})

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

func TestGetLastChallengeTime_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	t.Run("Database query error", func(t *testing.T) {
		mock.ExpectQuery("SELECT CreatedAt FROM Challenge").
			WithArgs(1001).
			WillReturnError(sql.ErrConnDone)

		lastTime, err := object.GetLastChallengeTime(db, 1001)
		assert.Error(t, err)
		assert.Nil(t, lastTime)
		assert.Contains(t, err.Error(), "error querying database with CreateChallenge")
	})

	t.Run("No rows returned", func(t *testing.T) {
		mock.ExpectQuery("SELECT CreatedAt FROM Challenge").
			WithArgs(1002).
			WillReturnError(sql.ErrNoRows)

		lastTime, err := object.GetLastChallengeTime(db, 1002)
		assert.NoError(t, err)
		assert.NotNil(t, lastTime)
		assert.True(t, lastTime.IsZero())
	})

	t.Run("Invalid data returned", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"CreatedAt"}).AddRow("invalid_time")

		mock.ExpectQuery("SELECT CreatedAt FROM Challenge").
			WithArgs(1003).
			WillReturnRows(rows)

		lastTime, err := object.GetLastChallengeTime(db, 1003)
		assert.Error(t, err)
		assert.Nil(t, lastTime)
		assert.Contains(t, err.Error(), "error querying database with CreateChallenge")
	})

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

func TestAddNewChallenge_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	t.Run("Database execution error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO Challenge").
			WithArgs(1001, 20.01).
			WillReturnError(errors.New("database error"))

		tx, _ := db.Begin()
		id, err := object.AddNewChallenge(tx, object_models.NewChallengeNeed{PlayerID: 1001, Amount: 20.01})

		assert.Error(t, err)
		assert.Equal(t, 0, id)
		assert.Contains(t, err.Error(), "error querying database with addNewChallenge")
	})

	t.Run("LastInsertId error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO Challenge").
			WithArgs(1002, 30.02).
			WillReturnResult(sqlmock.NewErrorResult(errors.New("LastInsertId error")))

		tx, _ := db.Begin()
		id, err := object.AddNewChallenge(tx, object_models.NewChallengeNeed{PlayerID: 1002, Amount: 30.02})

		assert.NoError(t, err)
		assert.Equal(t, 0, id)
	})

	t.Run("Invalid transaction", func(t *testing.T) {
		invalidTx := &sql.Tx{}
		id, err := object.AddNewChallenge(invalidTx, object_models.NewChallengeNeed{PlayerID: 1003, Amount: 40.03})

		assert.Error(t, err)
		assert.Equal(t, 0, id)
		assert.Contains(t, err.Error(), "error querying database with addNewChallenge")
	})

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

func TestUpdatePricePool_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	t.Run("Database execution error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE PrizePool").
			WithArgs(20.01).
			WillReturnError(errors.New("database error"))

		tx, _ := db.Begin()
		err := object.UpdatePricePool(tx, 20.01)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error updating price pool")
	})

	t.Run("No rows affected", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE PrizePool").
			WithArgs(30.02).
			WillReturnResult(sqlmock.NewResult(0, 0))

		tx, _ := db.Begin()
		err := object.UpdatePricePool(tx, 30.02)

		assert.NoError(t, err)
	})

	t.Run("Invalid transaction", func(t *testing.T) {
		invalidTx := &sql.Tx{}
		err := object.UpdatePricePool(invalidTx, 40.03)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error updating price pool")
	})

	t.Run("Negative amount", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE PrizePool").
			WithArgs(-50.04).
			WillReturnResult(sqlmock.NewResult(0, 1))

		tx, _ := db.Begin()
		err := object.UpdatePricePool(tx, -50.04)

		assert.NoError(t, err)
	})

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

func TestDistributePrizePool_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	tests := []struct {
		name        string
		setupMock   func(mock sqlmock.Sqlmock)
		expectedErr string
	}{
		{
			name: "Error fetching prize pool amount",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT Amount FROM PrizePool").WillReturnError(sql.ErrNoRows)
				mock.ExpectRollback()
			},
			expectedErr: "error fetching prize pool amount",
		},
		{
			name: "Error updating challenge",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT Amount FROM PrizePool").WillReturnRows(sqlmock.NewRows([]string{"Amount"}).AddRow(1000.0))
				mock.ExpectExec("UPDATE challenges").WillReturnError(sql.ErrTxDone)
				mock.ExpectRollback()
			},
			expectedErr: "error updating player's balance",
		},
		{
			name: "Error resetting prize pool",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery("SELECT Amount FROM PrizePool").WillReturnRows(sqlmock.NewRows([]string{"Amount"}).AddRow(1000.0))
				mock.ExpectExec("UPDATE challenges").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE PrizePool").WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			expectedErr: "error resetting prize pool",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mock)

			tx, err := db.Begin()
			if err != nil {
				t.Fatalf("Failed to begin transaction: %v", err)
			}

			err = object.DistributePrizePool(tx, 1, 1)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErr)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
