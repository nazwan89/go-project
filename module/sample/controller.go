package sample

import (
	"github.com/gofiber/fiber/v2"
)

func responseHello(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Hello from the Greeting Module!",
	})
}

func responseHelloWithName(c *fiber.Ctx) error {
	name := c.Params("name")
	if name == "" {
		name = "World"
	}
	return c.JSON(fiber.Map{
		"message": "Hello, " + name + "!",
	})
}

func responseHelloWithQuery(c *fiber.Ctx) error {
	name := c.Query("name", "World")
	return c.JSON(fiber.Map{
		"message": "Hello, " + name + "!",
	})
}

func responseHelloWithService(c *fiber.Ctx) error {
	name := c.Query("name", "World")
	greetingMessage := generateGreeting(name)
	return c.JSON(fiber.Map{
		"message": greetingMessage,
	})
}

func responseHelloWithForm(c *fiber.Ctx) error {
	var req Request

	if err := c.BodyParser(&req); err != nil {
		return c.JSON(Response{
			Message: "Invalid request body",
		})
	}

	name := req.Name
	return c.JSON(Response{
		Message: "Hello, " + name + "! This message was generated from the form handler.",
	})
}
