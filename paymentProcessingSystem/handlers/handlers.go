package handlers

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func SetupPaymentsRoutes(players *gin.RouterGroup, db *sql.DB) {
	// Player routes
	players.GET("/:id", func(c *gin.Context) {
		// GetPlayers(c, db)
	})
	players.POST("/", func(c *gin.Context) {
		// CreatePlayer(c, db)
	})
}
