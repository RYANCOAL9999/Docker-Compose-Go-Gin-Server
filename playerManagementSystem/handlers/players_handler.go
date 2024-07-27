package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/playerManagementSystem/databases"
	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/playerManagementSystem/models"

	"github.com/gin-gonic/gin"
)

func GetPlayers(c *gin.Context, db *sql.DB) {
	playerRanks, err := databases.GetPlayersData(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, playerRanks)
}

func CreatePlayer(c *gin.Context, db *sql.DB) {
	var newPlayerRank models.PlayerRank
	if err := c.BindJSON(&newPlayerRank); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := databases.AddPlayer(db, newPlayerRank.Name, newPlayerRank.LV)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func GetPlayer(c *gin.Context, db *sql.DB) {
	id, _ := strconv.Atoi(c.Param("id"))
	playerRank, err := databases.GetPlayer(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, playerRank)
}

func UpdatePlayer(c *gin.Context, db *sql.DB) {
	var playerRank models.PlayerRank
	if err := c.BindJSON(&playerRank); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := databases.UpdatePlayer(db, playerRank)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var args interface{}
	c.JSON(http.StatusOK, args)
}

func DeletePlayer(c *gin.Context, db *sql.DB) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := databases.DeletePlayer(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	var args interface{}
	c.JSON(http.StatusOK, args)
}
