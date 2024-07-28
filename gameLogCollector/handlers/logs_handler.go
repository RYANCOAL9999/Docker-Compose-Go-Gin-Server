package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/gameLogCollector/databases"
	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/gameLogCollector/models"
	"github.com/gin-gonic/gin"
)

const time_format string = "2006-01-02 00:00:00"

// @Summary      Retrieve game logs
// @Description  Fetches a list of game logs, allowing optional filtering by player ID, action, start time, end time, and limit. If more than one log is found, returns the first log. Returns a list of logs otherwise.
// @Tags         game_logs
// @Accept       json
// @Produce      json
// @Param        player_id  query  int     false  "Filter logs by player ID"
// @Param        action     query  string  false  "Filter logs by action"
// @Param        start_time query  string  false  "Start time for filtering logs (format: YYYY-MM-DDTHH:MM:SSZ)"
// @Param        end_time   query  string  false  "End time for filtering logs (format: YYYY-MM-DDTHH:MM:SSZ)"
// @Param        limit      query  int     false  "Limit the number of logs returned"
// @Success      200  {object}  []models.GameLog  "List of game logs matching the criteria"
// @Failure      400  {object}  models.ErrorResponse "Bad request due to invalid query parameters"
// @Failure      500  {object}  models.ErrorResponse "Internal server error"
// @Router       /game_logs [get]
func getLogs(c *gin.Context, db *sql.DB) {
	var args interface{}
	playerID, _ := strconv.Atoi(c.Query("player_id"))
	action := c.Query("action")
	startTime, _ := time.Parse(time_format, c.Query("start_time"))
	endTime, _ := time.Parse(time_format, c.Query("end_time"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	logs, err := databases.ListLogs(db, playerID, action, startTime, endTime, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	if len(logs) <= 1 {
		args = logs[0]
	} else {
		args = logs
	}
	c.JSON(http.StatusOK, args)
}

// @Summary      Create a game log
// @Description  Adds a new game log entry with the provided details. The request body must contain the player ID, action, timestamp, and details. Returns the ID of the newly created log entry if successful.
// @Tags         game_logs
// @Accept       json
// @Produce      json
// @Param        game_log  body models.GameLog  true  "Details of the game log to be created"
// @Success      201  {object}  models.CreateResponse "Game log created successfully, returns the ID of the new game log"
// @Failure      400  {object}  models.ErrorResponse "Bad request due to invalid input data"
// @Failure      500  {object}  models.ErrorResponse "Internal server error"
// @Router       /game_logs [post]
func CreateLog(c *gin.Context, db *sql.DB) {
	var newLog models.GameLog
	if err := c.BindJSON(&newLog); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	id, err := databases.AddLog(db, newLog)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, models.CreateResponse{ID: id})
}
