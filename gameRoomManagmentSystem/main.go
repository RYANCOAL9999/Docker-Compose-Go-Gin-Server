package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/gameRoomManagmentSystem/handlers"

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

	//Using the Default setting
	var r *gin.Engine = gin.Default()

	//write the logs to gin.DefaultWriter
	r.Use(gin.Logger())

	//Recovery returns a middleware if server is panics
	r.Use(gin.Recovery())

	// Setup Rooms routes
	handlers.SetupRoomsRoutes(r.Group("/rooms"), db)

	// Setup Reservations routes
	handlers.SetupReservationsRoutes(r.Group("/reservations"), db)

	fmt.Println("Starting Go API service...")

	fmt.Println(`
	______     ______        ______     ______   __    
   /\  ___\   /\  __ \      /\  __ \   /\  == \ /\ \   
   \ \ \__ \  \ \ \/\ \     \ \  __ \  \ \  _-/ \ \ \  
	\ \_____\  \ \_____\     \ \_\ \_\  \ \_\    \ \_\ 
	 \/_____/   \/_____/      \/_/\/_/   \/_/     \/_/ `)

	// Run with 8080 port
	r.Run(":8080")
}
