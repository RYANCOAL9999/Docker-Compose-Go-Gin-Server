package handlers

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func SetupChallengeRoutes(players *gin.RouterGroup, db *sql.DB) {
	// Player routes
	players.GET("/results", func(c *gin.Context) {
		// GetPlayers(c, db)
	})
	players.POST("/", func(c *gin.Context) {
		// CreatePlayer(c, db)
	})
}
