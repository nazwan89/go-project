package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	// Load environment variables from .env if present.
	// In production, set variables through the hosting provider or container orchestration.
	// ========================
	_ = godotenv.Load()

	// ========================
	// Load typed configuration from environment variables.
	// All config reads happen here — never call os.Getenv in handlers.
	// ========================
	cfg := utils.LoadConfig()

	// ========================
	// Fiber App Configuration
	// ========================
	app := fiber.New(fiber.Config{
		AppName:      cfg.AppName,
		ErrorHandler: utils.ErrorHandler,
		BodyLimit:    1 * 1024 * 1024, // 1MB — prevents memory exhaustion DoS (FOUND-03)
		ReadTimeout:  30 * time.Second, // bounds keepalive connections during graceful shutdown (FOUND-02)
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
			"timestamp": utils.CurrentTimestamp(cfg.Location),
		})
	})

	// ========================
	// Register Module Routes
	// ========================
	api := app.Group("/api")
	sample.RegisterRoutes(api)

	// ========================
	// Start Server (non-blocking) — FOUND-02
	// ========================
	go func() {
		log.Printf("Service Starting On Port %s", cfg.Port)
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Printf("Server stopped: %v", err)
		}
	}()

	// ========================
	// Graceful Shutdown — FOUND-02
	// Block until SIGTERM or SIGINT received.
	// ========================
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server (30s timeout)...")
	if err := app.ShutdownWithTimeout(30 * time.Second); err != nil {
		log.Fatalf("Server forced shutdown: %v", err)
	}
	log.Println("Server exited cleanly")
}
