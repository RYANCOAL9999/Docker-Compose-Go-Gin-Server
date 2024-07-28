package handlers

import (
	"database/sql"
	"log"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"

	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/endlessChallengeSystem/databases"
	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/endlessChallengeSystem/models"
	"github.com/gin-gonic/gin"
)

// 1% chance of winning
const winProbability float64 = 0.01

func CalculateChallengeResult(db *sql.DB, challengeID int, playerID int, probability float64) {

	// Delay the calculation by 30 seconds
	time.Sleep(30 * time.Second)

	localProbability := winProbability + probability

	won := rand.Float64() < winProbability

	var joined models.Status = models.Joined

	tx, err := db.Begin()
	if err != nil {
		log.Printf("Failed to start transaction: %v", err)
		return
	}
	defer tx.Rollback()

	if won {
		err = databases.DistributePrizePool(tx, challengeID, playerID, joined)
	} else {
		err = databases.UpdateProbability(tx, challengeID, playerID, localProbability, joined)
	}

	if err != nil {
		log.Printf("Failed to distribute prize pool: %v", err)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return
	}

	log.Printf("Challenge result calculated for player %d. Won: %v", playerID, won)
}

// @Summary      Join a challenge
// @Description  Allows a player to join a new challenge, provided they haven't participated in the last minute. It processes the challenge creation within a transaction, updates the prize pool, and starts a background task to calculate the challenge result after 30 seconds. Returns the status of the challenge creation.
// @Tags         challenges
// @Accept       json
// @Produce      json
// @Param        challenge  body  models.NewChallengeNeed  true  "Details for joining the challenge"
// @Success      201  {object}  models.JoinChallengeResponse "Challenge joined successfully, returns the status of the challenge, it represent as number, 1 is joined, 0 is Ready"
// @Failure      400  {object}  models.ErrorResponse "Bad request due to invalid input data"
// @Failure      425  {object}  models.ErrorResponse "Too many requests if attempting to join within a minute"
// @Failure      500  {object}  models.ErrorResponse "Internal server error during challenge creation or transaction"
// @Router       /challenges/join [post]
func JoinChallenges(c *gin.Context, db *sql.DB) {
	var newChallengeNeed models.NewChallengeNeed

	if err := c.ShouldBindJSON(&newChallengeNeed); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	var probability float64 = 0

	lastChallengeTime, lastprobability, err := databases.GetLastChallenge(db, newChallengeNeed.PlayerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	if time.Since(*lastChallengeTime) < time.Minute {
		c.JSON(http.StatusTooEarly, models.ErrorResponse{Error: "You can only participate once per minute"})
		return
	}

	//No error means that player is ready to join Challenge
	const status models.Status = models.Ready

	//Having a last record
	if !lastChallengeTime.IsZero() {
		probability = lastprobability
	}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	lastChallengeID, err := databases.AddNewChallenge(tx, newChallengeNeed, status, probability)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to Add New Challenge "})
		return
	}

	// need to update PrizePool Value
	err = databases.UpdatePricePool(tx, newChallengeNeed.Amount)
	if err != nil {
		c.JSON(http.StatusTooEarly, models.ErrorResponse{Error: "Failed to update price pool"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to commit transaction"})
		return
	}

	go func() {
		time.Sleep(30 * time.Second)
		CalculateChallengeResult(db, lastChallengeID, newChallengeNeed.PlayerID, probability)
	}()

	c.JSON(http.StatusCreated, models.JoinChallengeResponse{Status: status})
}

// @Summary      List recent challenges
// @Description  Retrieves a list of recent challenges based on the provided limit. Returns the most recent challenge if there are multiple results.
// @Tags         challenges
// @Accept       json
// @Produce      json
// @Param        limit  query  int  true  "Maximum number of challenges to retrieve"
// @Success      200  {object}  []models.Challenge  "List of recent challenges or the most recent challenge"
// @Failure      400  {object}  models.ErrorResponse "Bad request due to invalid input data"
// @Failure      500  {object}  models.ErrorResponse "Internal server error during retrieval"
// @Router       /challenges [get]
func ShowChallenges(c *gin.Context, db *sql.DB) {
	var args interface{}
	limit, _ := strconv.Atoi(c.Query("limit"))
	challenges, err := databases.ListChallenges(db, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	if len(challenges) > 1 {
		args = challenges[0]
	} else {
		args = challenges
	}
	c.JSON(http.StatusOK, args)
}
