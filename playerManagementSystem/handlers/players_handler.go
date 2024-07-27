package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/playerManagementSystem/databases"
	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/playerManagementSystem/models"

	"github.com/gin-gonic/gin"
)

// @Summary      List players
// @Description  Retrieve a list of players and their ranks from the database.
// @Tags         players
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.PlayerRank  "A list of players with their ranks"
// @Failure      500  {object}  error  "Internal server error"
// @Router       /players [get]
func GetPlayers(c *gin.Context, db *sql.DB) {
	playerRanks, err := databases.GetPlayersData(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, playerRanks)
}

// @Summary      Create a new player
// @Description  Create a new player in the database using the provided player details.
// @Tags         players
// @Accept       json
// @Produce      json
// @Param        player  body  models.PlayerRank  true  "Player details to be created"
// @Success      201  {object}  number  "Player created successfully with the generated ID"
// @Failure      400  {object}  error  "Bad request due to invalid input"
// @Failure      500  {object}  error  "Internal server error"
// @Router       /players [post]
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

// @Summary      Retrieve a player by ID
// @Description  Get details of a specific player identified by their ID from the database.
// @Tags         players
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Player ID"
// @Success      200  {object}  models.PlayerRank  "Player details"
// @Failure      400  {object}  error  "Invalid ID supplied"
// @Failure      500  {object}  error  "Internal server error"
// @Router       /players/{id} [get]
func GetPlayer(c *gin.Context, db *sql.DB) {
	id, _ := strconv.Atoi(c.Param("id"))
	playerRank, err := databases.GetPlayer(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, playerRank)
}

// @Summary      Update player details
// @Description  Update the details of an existing player in the database using the provided player information.
// @Tags         players
// @Accept       json
// @Produce      json
// @Param        player  body  models.PlayerRank  true  "Player details to be updated"
// @Success      200  {object}  string  "Player updated successfully"
// @Failure      400  {object}  error  "Bad request due to invalid input"
// @Failure      500  {object}  error  "Internal server error"
// @Router       /players [put]
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

// @Summary      Delete a player
// @Description  Remove a player from the database using the provided player ID.
// @Tags         players
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Player ID to be deleted"
// @Success      200  {object}  string  "Player deleted successfully"
// @Failure      400  {object}  error  "Invalid ID supplied"
// @Failure      500  {object}  error  "Internal server error"
// @Router       /players/{id} [delete]
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
