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

func getLogs(c *gin.Context, db *sql.DB) {

	var args interface{}
	playerID, _ := strconv.Atoi(c.Query("player_id"))
	action := c.Query("action")
	startTime, _ := time.Parse(time_format, c.Query("start_time"))
	endTime, _ := time.Parse(time_format, c.Query("end_time"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	logs, err := databases.ListLogs(db, playerID, action, startTime, endTime, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(logs) > 1 {
		args = logs[0]
	} else {
		args = logs
	}
	c.JSON(http.StatusOK, args)
}

func CreateLog(c *gin.Context, db *sql.DB) {
	var newLog models.Log
	if err := c.BindJSON(&newLog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := databases.AddLog(db, newLog)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}
