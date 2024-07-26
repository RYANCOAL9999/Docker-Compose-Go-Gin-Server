package handlers

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func SetupPaymentsRoutes(payments *gin.RouterGroup, db *sql.DB) {
	// Player routes
	payments.GET("/:id", func(c *gin.Context) { ShowPayment(c, db) })
	payments.POST("/", func(c *gin.Context) { CreatePayment(c, db) })
}
