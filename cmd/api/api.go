package api

import (
	"github.com/gofiber/fiber/v2"
	"project/module/sample"
)

func RegisterRoutes(app *fiber.App) {
	api := app.Group("/api")
	sample.RegisterRoutes(api)
}
