package main

import (
	"fmt"
	"os"
	"social_api/db"
	"social_api/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func main() {
	err := db.InitDB()

	if err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		os.Exit(1)
	}

	defer db.CloseDB()

	app := fiber.New()

	app.Use(requestid.New())

	router.SetupRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	app.Listen(":" + port)
}
