package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/playerManagementSystem/docs"
	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/playerManagementSystem/handlers"
	"github.com/joho/godotenv"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	// Database connection
	db, err := sql.Open("mysql", os.Getenv("DB_CONNECTION_STRING")+"?parseTime=true")
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

	docs.SwaggerInfo.BasePath = "/api/v1"

	// Setup Levels routes
	handlers.SetupLevelsRoutes(r.Group("/levels"), db)

	// Setup Players routes
	handlers.SetupPlayersRoutes(r.Group("/players"), db)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	fmt.Println("Starting Go API service...")

	fmt.Println(`
	______     ______        ______     ______   __    
   /\  ___\   /\  __ \      /\  __ \   /\  == \ /\ \   
   \ \ \__ \  \ \ \/\ \     \ \  __ \  \ \  _-/ \ \ \  
	\ \_____\  \ \_____\     \ \_\ \_\  \ \_\    \ \_\ 
	 \/_____/   \/_____/      \/_/\/_/   \/_/     \/_/ `)

	// Run with port
	r.Run(os.Getenv("PORT"))
}
