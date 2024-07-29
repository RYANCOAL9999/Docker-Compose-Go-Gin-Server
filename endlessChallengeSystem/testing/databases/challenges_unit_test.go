package databases

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"regexp"
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
	// Setup
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer db.Close()

	// Test cases
	t.Run("Database query error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM Challenge").
			WillReturnError(sql.ErrConnDone)

		challenges, err := object.ListChallenges(db, 0)
		assert.Error(t, err)
		assert.Nil(t, challenges)
		assert.Contains(t, err.Error(), "error querying database with ListChallenges")
	})

	t.Run("Row scan error", func(t *testing.T) {
		// Simulate a row with invalid data
		rows := sqlmock.NewRows([]string{"ID", "PlayerID", "Amount", "Status", "Won", "CreatedAt", "Probability"}).
			AddRow("invalid", "1001", 20.01, "Joined", false, "invalid_time", 0.5)

		mock.ExpectQuery("SELECT (.+) FROM Challenge").
			WillReturnRows(rows)

		challenges, err := object.ListChallenges(db, 0)
		assert.Error(t, err)
		assert.Nil(t, challenges)
		assert.Contains(t, err.Error(), "error scanning row with ListChallenges")
	})

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
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

	mock.ExpectQuery(regexp.QuoteMeta("SELECT CreatedAt, Probability FROM Challenge WHERE PlayerID = ? ORDER BY CreatedAt DESC LIMIT 1")).
		WithArgs(playerID).
		WillReturnRows(rows)

	resultTime, resultProb, err := object.GetLastChallenge(db, playerID)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if resultTime == nil || resultTime.Sub(lastChallengeTime).Abs() > time.Millisecond {
		t.Errorf("expected time close to %v, got %v", lastChallengeTime, resultTime)
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
	expectedSQL := "SELECT CreatedAt, Probability FROM Challenge WHERE PlayerID = ? ORDER BY CreatedAt DESC LIMIT 1"

	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WithArgs(playerID).
		WillReturnError(expectedError)

	fmt.Printf("Debug: Expected SQL query: %s\n", expectedSQL)

	resultTime, resultProb, err := object.GetLastChallenge(db, playerID)
	fmt.Printf("Debug: resultTime=%v, resultProb=%v, err=%v\n", resultTime, resultProb, err)

	if err == nil {
		t.Error("expected an error, but got nil")
	} else if !strings.Contains(err.Error(), expectedError.Error()) {
		t.Errorf("expected error containing '%v', got '%v'", expectedError, err)
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

	newChallengeNeed := object_models.NewChallengeNeed{
		PlayerID: 1,
		Amount:   100,
	}
	probability := 0.5
	status := object_models.Ready

	// Define the expected SQL for INSERT
	sql := regexp.QuoteMeta("INSERT INTO Challenge (PlayerID, Amount, Status, Won, CreatedAt, Probability) VALUES (?, ?, ?, false, NOW(), ?)")

	// Expect transaction to begin
	mock.ExpectBegin()

	// Expect the INSERT query
	mock.ExpectExec(sql).
		WithArgs(newChallengeNeed.PlayerID, newChallengeNeed.Amount, int(status), probability).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Expect transaction to be committed
	mock.ExpectCommit()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when beginning a transaction", err)
	}

	// Call the function
	_, err = object.AddNewChallenge(tx, newChallengeNeed, status, probability)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	err = tx.Commit() // Commit on success
	if err != nil {
		t.Errorf("unexpected error with commit on TestAddNewChallenge: %s", err)
	}

	// Check if all expectations were met
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

	// Expect transaction to fail when beginning
	mock.ExpectBegin().WillReturnError(errors.New("failed to begin transaction"))

	_, err = db.Begin()
	if err == nil {
		t.Errorf("expected an error when beginning a transaction, but got nil")
	}

	// Check if all expectations were met
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
	// Expect transaction to be committed

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

	// Define the amount for the prize pool update
	amount := 100.1

	// Define the expected SQL for UPDATE
	sql := regexp.QuoteMeta("UPDATE PrizePool SET Amount = Amount + ? WHERE ID = 1")

	// Expect transaction to begin
	mock.ExpectBegin()

	// Simulate an error during the UPDATE
	mock.ExpectExec(sql).
		WithArgs(amount).
		WillReturnError(errors.New("some error"))

	// Expect transaction to be rolled back due to the error
	mock.ExpectRollback()

	// Start a new transaction
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when beginning a transaction", err)
	}

	// Call the function
	err = object.UpdatePricePool(tx, amount)
	if err == nil {
		t.Errorf("expected an error when updating the prize pool, but got nil")
	}

	// Verify the error message
	if !strings.Contains(err.Error(), "error updating price pool") {
		t.Errorf("unexpected error message: %s", err)
	}

	// Rollback the transaction (simulate the error handling in the function)
	if err := tx.Rollback(); err != nil {
		t.Errorf("unexpected error with rollback: %s", err)
	}

	// Check if all expectations were met
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

	// Define the amount for updating the prize pool
	amount := 100.0

	// Define the expected SQL for UPDATE
	sql := regexp.QuoteMeta("UPDATE PrizePool SET Amount = Amount + ? WHERE ID = 1")

	mock.ExpectBegin()
	mock.ExpectExec(sql).
		WithArgs(amount).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when beginning a transaction", err)
	}

	err = object.UpdatePricePool(tx, 100.0)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdatePricePool_ErrorMalformedInput(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Define the amount for updating the prize pool
	amount := 3.141516192132132134

	// Define the expected SQL for UPDATE
	sql := regexp.QuoteMeta("UPDATE PrizePool SET Amount = Amount + ? WHERE ID = 1")

	// Expect the UPDATE operation to fail with a malformed input error
	mock.ExpectBegin()
	mock.ExpectExec(sql).
		WithArgs(amount).
		WillReturnError(fmt.Errorf("malformed input"))

	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when beginning a transaction", err)
	}

	// Call the function to test
	err = object.UpdatePricePool(tx, amount)
	if err == nil {
		t.Errorf("expected an error but got nil")
	}

	// Check if all expectations were met
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

	challengeID := 1
	playerID := 1
	probability := 0.75
	status := object_models.Status(1)

	// Mock Exec to return an error
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE Challenge").
		WithArgs(probability, int(status), challengeID, playerID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when beginning a transaction", err)
	}

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

	challengeID := 999
	playerID := 888
	probability := 0.5
	status := object_models.Ready

	// Mock Exec to return an error
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE Challenge").
		WithArgs(probability, status, challengeID, playerID).
		WillReturnError(errors.New("SQL execution error"))
	mock.ExpectRollback()

	// Assertion
	var logOutput bytes.Buffer
	log.SetOutput(&logOutput)

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when beginning a transaction", err)
	}

	// Call the function
	err = object.UpdateProbability(tx, challengeID, playerID, probability, status)
	if err == nil {
		t.Errorf("Expected an error due to SQL execution failure, but got nil")
	}

	// Ensure expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
