package handlers

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func SetupLogsRoutes(logs *gin.RouterGroup, db *sql.DB) {
	// Player routes
	logs.POST("/", func(c *gin.Context) { CreateLog(c, db) })
	logs.GET("/", func(c *gin.Context) { getLogs(c, db) })
}
