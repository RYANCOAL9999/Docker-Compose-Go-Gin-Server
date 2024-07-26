package handlers

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func SetupLogsRoutes(players *gin.RouterGroup, db *sql.DB) {
	// Player routes
	players.GET("/", func(c *gin.Context) {
		// GetPlayers(c, db)
	})
}
