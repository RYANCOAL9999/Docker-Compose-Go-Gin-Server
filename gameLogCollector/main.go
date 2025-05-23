package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/gameLogCollector/docs"
	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/gameLogCollector/handlers"
	"github.com/joho/godotenv"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// @title Game Log Collector API
// @version 1.0
// @description This is a game log collector server.
// @termsOfService http://swagger.io/terms/

// @contact.name Steven Poon
// @contact.url  https://github.com/RYANCOAL9999
// @contact.email lmf242003@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host :8084
// @BasePath /v2
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
	handlers.SetupLogsRoutes(r.Group("/logs"), db)

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
