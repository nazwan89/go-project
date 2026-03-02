package sample

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(api fiber.Router) {
	group := api.Group("/sample")

	group.Get("/hello", responseHello)
	group.Get("/hello/:name", responseHelloWithName)
	group.Get("/hello-query", responseHelloWithQuery)
	group.Get("/hello-service", responseHelloWithService)
	group.Post("/hello-form", responseHelloWithForm)
}
