package databases

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
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
	assert.Equal(t, int(1), challenges[0].ID)
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

func TestGetLastChallenge(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	playerID := 1
	lastChallengeTime := time.Now()
	lastProbability := 0.75

	rows := sqlmock.NewRows([]string{"CreatedAt", "Probability"}).
		AddRow(lastChallengeTime, lastProbability)

	mock.ExpectQuery("SELECT CreatedAt, Probability FROM Challenge WHERE PlayerID = ? ORDER BY CreatedAt DESC LIMIT 1").
		WithArgs(playerID).
		WillReturnRows(rows)

	resultTime, resultProb, err := object.GetLastChallenge(db, playerID)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if resultTime == nil || !resultTime.Equal(lastChallengeTime) {
		t.Errorf("expected %v, got %v", lastChallengeTime, resultTime)
	}

	if resultProb != lastProbability {
		t.Errorf("expected %v, got %v", lastProbability, resultProb)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetLastChallenge_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	playerID := 1
	expectedError := fmt.Errorf("some SQL error")

	mock.ExpectQuery("SELECT CreatedAt, Probability FROM Challenge WHERE PlayerID = ? ORDER BY CreatedAt DESC LIMIT 1").
		WithArgs(playerID).
		WillReturnError(expectedError)

	resultTime, resultProb, err := object.GetLastChallenge(db, playerID)
	if err == nil || !errors.Is(err, expectedError) {
		t.Errorf("expected error %v, got %v", expectedError, err)
	}

	if resultTime != nil {
		t.Errorf("expected nil time, got %v", resultTime)
	}

	if resultProb != 0 {
		t.Errorf("expected probability 0, got %v", resultProb)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAddNewChallenge(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when beginning a transaction", err)
	}

	newChallengeNeed := object_models.NewChallengeNeed{
		PlayerID: 1,
		Amount:   100,
	}
	status := object_models.Status(1)
	probability := 0.5

	mock.ExpectExec("INSERT INTO Challenge").
		WithArgs(newChallengeNeed.PlayerID, newChallengeNeed.Amount, int(status), probability).
		WillReturnResult(sqlmock.NewResult(1, 1))

	challengeID, err := object.AddNewChallenge(tx, newChallengeNeed, status, probability)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if challengeID != 1 {
		t.Errorf("expected challengeID to be 1, got %d", challengeID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAddNewChallenge_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when beginning a transaction", err)
	}

	newChallengeNeed := object_models.NewChallengeNeed{
		PlayerID: 1,
		Amount:   100,
	}
	status := object_models.Status(1)
	probability := 0.5

	mock.ExpectExec("INSERT INTO Challenge").
		WithArgs(newChallengeNeed.PlayerID, newChallengeNeed.Amount, int(status), probability).
		WillReturnError(fmt.Errorf("some SQL error"))

	_, err = object.AddNewChallenge(tx, newChallengeNeed, status, probability)
	if err == nil {
		t.Error("expected an error but got none")
	}

	if !strings.Contains(err.Error(), "error querying database with addNewChallenge") {
		t.Errorf("unexpected error message: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
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
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when beginning a transaction", err)
	}

	mock.ExpectExec("UPDATE PrizePool").
		WithArgs(100.0).
		WillReturnError(fmt.Errorf("some error"))

	err = object.UpdatePricePool(tx, 100.0)
	if err == nil {
		t.Errorf("expected an error but got none")
	}

	if rollbackErr := tx.Rollback(); rollbackErr != nil {
		t.Errorf("unexpected error during rollback: %s", rollbackErr)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDistributePrizePool(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when beginning a transaction", err)
	}

	mock.ExpectExec("UPDATE PrizePool").
		WithArgs(100.0).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = object.UpdatePricePool(tx, 100.0)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDistributePrizePool_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when beginning a transaction", err)
	}

	injectionString := "-3.14159265"
	mock.ExpectExec("UPDATE PrizePool").
		WithArgs(injectionString).
		WillReturnError(fmt.Errorf("malformed input"))
	str, _ := strconv.ParseFloat(injectionString, 64)
	err = object.UpdatePricePool(tx, str)
	if err == nil {
		t.Errorf("expected an error but got none")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateProbability(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when beginning a transaction", err)
	}

	challengeID := 1
	playerID := 1
	probability := 0.75
	status := object_models.Status(1)

	mock.ExpectExec("UPDATE Challenge").
		WithArgs(probability, int(status), challengeID, playerID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = object.UpdateProbability(tx, challengeID, playerID, probability, status)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateProbability_Error(t *testing.T) {
	// Setup
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when beginning a transaction", err)
	}
	challengeID := 999
	playerID := 888
	probability := 0.5
	status := object_models.Ready

	// Mock Exec to return an error
	mock.ExpectExec("UPDATE Challenge").
		WithArgs(probability, status, challengeID, playerID).
		WillReturnError(errors.New("SQL execution error"))

	// Assertion
	var logOutput bytes.Buffer
	log.SetOutput(&logOutput)

	// Call the function
	err = object.UpdateProbability(tx, challengeID, playerID, probability, status)
	if err == nil {
		t.Errorf("Expected an error due to SQL execution failure, but got nil")
	} else if err.Error() != "error updating probability: SQL execution error" {
		t.Errorf("Expected error message 'error updating probability: SQL execution error', but got '%s'", err)
	}

	// Ensure expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// Rollback the transaction in the test
	if err := tx.Rollback(); err != nil {
		t.Errorf("Failed to rollback transaction: %v", err)
	}
}
