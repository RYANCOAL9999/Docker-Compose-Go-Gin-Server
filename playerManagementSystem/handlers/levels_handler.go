package handlers

import (
	"database/sql"
	"net/http"

	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/playerManagementSystem/models"

	"github.com/gin-gonic/gin"
)

func GetLevels(c *gin.Context, db *sql.DB) {
	rows, err := db.Query("SELECT id, name FROM levels")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	var levels []models.Level
	for rows.Next() {
		var l models.Level
		if err := rows.Scan(&l.ID, &l.Name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		levels = append(levels, l)
	}

	c.JSON(http.StatusOK, levels)
}

func CreateLevel(c *gin.Context, db *sql.DB) {
	var newLevel models.Level
	if err := c.BindJSON(&newLevel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("INSERT INTO levels (name) VALUES (?)", newLevel.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{"id": id})
}
