package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"

	"project/module/sample"
	"project/utils"
)

func main() {
	// ========================
	// load environment variables from .env if present
	// Go does _not_ automatically read the file; you must do this yourself
	// or export the variables before running.
	// In production, you should set environment variables through your hosting provider or container orchestration system.
	// ========================
	_ = godotenv.Load() // ignore error – file may not exist in production

	// ========================
	// Fiber App Configuration
	// ========================
	app := fiber.New(fiber.Config{
		AppName:      "Project Name",
		ErrorHandler: utils.ErrorHandler,
	})

	// ========================
	// Middleware
	// ========================
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(healthcheck.New())

	// ========================
	// Root Route
	// ========================
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "OK",
			"message":   "Service is running",
			"timestamp": utils.CurrentTimestamp(),
		})
	})

	// ========================
	// Register Module Routes
	// ========================
	api := app.Group("/api")
	sample.RegisterRoutes(api)

	// ========================
	// Port Configuration
	// ========================
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf(
		"Service Starting On Port %s",
		port,
	)

	// ========================
	// Start Server
	// ========================
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf(
			"Failed To Start Server: %v",
			err,
		)
	}
}
