package handlers

import (
	"database/sql"
	"net/http"

	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/playerManagementSystem/databases"
	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/playerManagementSystem/models"

	"github.com/gin-gonic/gin"
)

func GetLevels(c *gin.Context, db *sql.DB) {
	levels, err := databases.GetLevelsData(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, levels)
}

func CreateLevel(c *gin.Context, db *sql.DB) {
	var newLevel models.Level
	if err := c.BindJSON(&newLevel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := databases.AddLevel(db, newLevel.Name, newLevel.Rank)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}
