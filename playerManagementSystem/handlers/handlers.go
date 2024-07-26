package handlers

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

type RouterFunctionHeader struct {
	context *gin.Context
	db      *sql.DB
}

func NewRouterFunctionHeader(c *gin.Context, db *sql.DB) *RouterFunctionHeader {
	return &RouterFunctionHeader{
		context: c,
		db:      db,
	}
}

func SetupPlayersRoutes(players *gin.RouterGroup, db *sql.DB) {
	// Player routes
	players.GET("/", func(c *gin.Context) { GetPlayers(c, db) })
	players.POST("/", func(c *gin.Context) { CreatePlayer(c, db) })
	players.GET("/:id", func(c *gin.Context) { GetPlayer(c, db) })
	players.PUT("/:id", func(c *gin.Context) { UpdatePlayer(c, db) })
	players.DELETE("/:id", func(c *gin.Context) { DeletePlayer(c, db) })
}

func SetupLevelsRoutes(levels *gin.RouterGroup, db *sql.DB) {
	// Level routes
	levels.GET("/", func(c *gin.Context) { GetLevels(c, db) })
	levels.POST("/", func(c *gin.Context) { CreateLevel(c, db) })
}
