package handlers

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func SetupLogsRoutes(Logs *gin.RouterGroup, db *sql.DB) {
	// Player routes
	Logs.POST("/", func(c *gin.Context) {})
	Logs.GET("/results", func(c *gin.Context) {})
}
