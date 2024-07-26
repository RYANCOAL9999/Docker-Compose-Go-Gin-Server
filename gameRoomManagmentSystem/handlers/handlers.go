package handlers

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func SetupRoomsRoutes(rooms *gin.RouterGroup, db *sql.DB) {
	// Player routes
	rooms.GET("/", func(c *gin.Context) { GetRooms(c, db) })
	rooms.POST("/", func(c *gin.Context) { CreateRoom(c, db) })
	rooms.GET("/:id", func(c *gin.Context) { GetRoom(c, db) })
	rooms.PUT("/:id", func(c *gin.Context) { UpdateRoom(c, db) })
	rooms.DELETE("/:id", func(c *gin.Context) { DeleteRoom(c, db) })
}

func SetupReservationsRoutes(reservations *gin.RouterGroup, db *sql.DB) {
	// Level routes
	reservations.GET("/", func(c *gin.Context) { GetReservations(c, db) })
	reservations.POST("/", func(c *gin.Context) { CreateReservations(c, db) })
}
