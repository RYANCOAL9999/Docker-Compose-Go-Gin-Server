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
const winProbability float32 = 0.01

// time format
const time_format string = "2006-01-02 00:00:00"

func CalculateChallengeResult(db *sql.DB, playerID int) {
	// Delay the calculation by 30 seconds
	time.Sleep(30 * time.Second)

	won := rand.Float32() < winProbability

	tx, err := db.Begin()
	if err != nil {
		log.Printf("Failed to start transaction: %v", err)
		return
	}
	defer tx.Rollback()

	err = databases.UpdateChallenge(tx, models.Ready, won, playerID)
	if err != nil {
		log.Printf("Failed to update challenge: %v", err)
		return
	}

	if won {
		err = databases.DistributePrizePool(tx, playerID)
		if err != nil {
			log.Printf("Failed to distribute prize pool: %v", err)
			return
		}
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

	err = databases.AddNewChallenge(tx, newChallengeNeed)
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

	go func() {
		time.Sleep(30 * time.Second)
		CalculateChallengeResult(db, newChallengeNeed.PlayerID)
	}()

	challengeStatus, err := databases.GetChallengeStatus(db, newChallengeNeed.PlayerID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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
