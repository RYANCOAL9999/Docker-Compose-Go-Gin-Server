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
	var localprobability float64 = probability
	if localprobability == 0 {
		localprobability += rand.Float64()
	}

	won := localprobability < winProbability

	tx, err := db.Begin()
	if err != nil {
		log.Printf("Failed to start transaction: %v", err)
		return
	}
	defer tx.Rollback()

	if won {
		err = databases.DistributePrizePool(tx, challengeID, playerID)
	} else {
		err = databases.Updateprobability(tx, challengeID, playerID, localprobability)
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

func JoinChallenges(c *gin.Context, db *sql.DB) {
	var newChallengeNeed models.NewChallengeNeed

	if err := c.ShouldBindJSON(&newChallengeNeed); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lastChallengeTime, err := databases.GetLastChallengeTime(db, newChallengeNeed.PlayerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if time.Since(*lastChallengeTime) < time.Minute {
		c.JSON(http.StatusTooEarly, gin.H{"error": "You can only participate once per minute"})
		return
	}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	lastChallengeID, err := databases.AddNewChallenge(tx, newChallengeNeed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Add New Challenge "})
		return
	}

	// need to update PrizePool Value
	err = databases.UpdatePricePool(tx, newChallengeNeed.Amount)
	if err != nil {
		c.JSON(http.StatusTooEarly, gin.H{"error": "Failed to update price pool"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	challengeStatus, probability, err := databases.GetChallenge(db, lastChallengeID, newChallengeNeed.PlayerID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go func() {
		time.Sleep(30 * time.Second)
		CalculateChallengeResult(db, lastChallengeID, newChallengeNeed.PlayerID, *probability)
	}()

	c.JSON(http.StatusCreated, gin.H{"status": challengeStatus})
}

func ShowChallenges(c *gin.Context, db *sql.DB) {
	var args interface{}
	limit, _ := strconv.Atoi(c.Query("limit"))
	challenges, err := databases.ListChallenges(db, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(challenges) > 1 {
		args = challenges[0]
	} else {
		args = challenges
	}
	c.JSON(http.StatusOK, args)
}
