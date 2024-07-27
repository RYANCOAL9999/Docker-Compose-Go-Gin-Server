package handlers

import (
	"database/sql"
	"net/http"

	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/playerManagementSystem/databases"
	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/playerManagementSystem/models"

	"github.com/gin-gonic/gin"
)

// @Summary      List levels
// @Description  Retrieve a list of levels from the database.
// @Tags         levels
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Level  "A list of levels"
// @Failure      500  {object}  error  "Internal server error"
// @Router       /levels [get]
func GetLevels(c *gin.Context, db *sql.DB) {
	levels, err := databases.GetLevelsData(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, levels)
}

// @Summary      Create a new level
// @Description  Create a new level in the database using the provided level details.
// @Tags         levels
// @Accept       json
// @Produce      json
// @Param        level  body  models.Level  true  "Level details to be created"
// @Success      201  {object}  number "Level created successfully with the generated ID"
// @Failure      400  {object}  error  "Bad request due to invalid input"
// @Failure      500  {object}  error  "Internal server error"
// @Router       /levels [post]
func CreateLevel(c *gin.Context, db *sql.DB) {
	var newLevel models.Level
	if err := c.BindJSON(&newLevel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := databases.AddLevel(db, newLevel.Name, newLevel.LV)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}
