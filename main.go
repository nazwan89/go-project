    package main

    import (
        "log"
        "os"

        "github.com/joho/godotenv"
        "github.com/gofiber/fiber/v2"
        "github.com/gofiber/fiber/v2/middleware/logger"
        "github.com/gofiber/fiber/v2/middleware/recover"
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
            AppName: "Project Name",
            // ErrorHandler: handlers.ErrorHandler,
        })

		// ========================
		// Middleware
		// ========================
		app.Use(recover.New())

		app.Use(logger.New(
			logger.Config{
				Format: "[${time}] ${status} - ${method} ${path}\n",
			},
		))

		// ========================
		// Health Check
		// ========================
		app.Get("/", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"status":  "OK",
				"message": "Service is running",
			})
		})

		// ========================
		// Port Configuration
		// ========================
		port := os.Getenv("PORT")
		if port == "" {port = "8080"}

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