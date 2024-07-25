package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/playerManagementSystem/models"

	"github.com/gin-gonic/gin"
)

func GetPlayers(c *gin.Context, db *sql.DB) {
	rows, err := db.Query("SELECT id, name, level FROM players")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var players []models.Player
	for rows.Next() {
		var p models.Player
		if err := rows.Scan(&p.ID, &p.Name, &p.Level); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		players = append(players, p)
	}

	c.JSON(http.StatusOK, players)
}

func CreatePlayer(c *gin.Context, db *sql.DB) {
	var newPlayer models.Player
	if err := c.BindJSON(&newPlayer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("INSERT INTO players (name, level) VALUES (?, ?)", newPlayer.Name, newPlayer.Level)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func GetPlayer(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	var p models.Player
	err := db.QueryRow("SELECT id, name, level FROM players WHERE id = ?", id).Scan(&p.ID, &p.Name, &p.Level)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)

}

func UpdatePlayer(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	var player models.Player
	if err := c.BindJSON(&player); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("UPDATE players SET name = ?, level = ? WHERE id = ?", player.Name, player.Level, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	player.ID, _ = strconv.Atoi(id)
	c.JSON(http.StatusOK, player)
}

func DeletePlayer(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	result, err := db.Exec("DELETE FROM players WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
