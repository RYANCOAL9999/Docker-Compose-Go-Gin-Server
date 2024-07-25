package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/playerManagementSystem/handlers"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Database connection
	db, err := sql.Open("mysql", os.Getenv("DB_CONNECTION_STRING"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Setup Players routes
	handlers.SetupPlayersRoutes(r.Group("/players"), db)

	// Setup Levels routes
	handlers.SetupLevelsRoutes(r.Group("/levele"), db)

	r.Run(":8080")
}
