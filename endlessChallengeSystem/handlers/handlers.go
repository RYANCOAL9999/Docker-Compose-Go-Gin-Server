package handlers

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func SetupChallengeRoutes(challenges *gin.RouterGroup, db *sql.DB) {
	challenges.POST("/", func(c *gin.Context) { JoinChallenges(c, db) })
	challenges.GET("/results", func(c *gin.Context) { ShowChallenges(c, db) })
}
