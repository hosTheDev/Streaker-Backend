package main

import (
	"streaker-backend/config"
	"streaker-backend/routes"

	"github.com/gofiber/fiber/v2"

	//To access the .env file.
	"os"
	"github.com/joho/godotenv"

	//for logging
	"log"
)

func main() {
	err := godotenv.Load()
	if(err != nil){
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")

	// Initialize Fiber
	app := fiber.New()

	// Connect to MongoDB
	config.ConnectDB(os.Getenv("CONNECTION_STRING"))

	// Setup Routes
	routes.SetupRoutes(app)

	// Start Server
	app.Listen(port)
}
