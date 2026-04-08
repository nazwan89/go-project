package sample

import (
	"github.com/gofiber/fiber/v2"

	"project/utils"
)

func responseHello(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Hello from the Greeting Module!",
	})
}

func responseHelloWithName(c *fiber.Ctx) error {
	name := c.Params("name")

	// FOUND-04: Validate path parameter before use
	if err := utils.ValidatePathParam(name); err != nil {
		return utils.BadRequestHandler(c, "Invalid path parameter: name must be 1-100 alphanumeric characters")
	}

	return c.JSON(fiber.Map{
		"message": "Hello, " + name + "!",
	})
}

func responseHelloWithQuery(c *fiber.Ctx) error {
	name := c.Query("name", "World")

	// FOUND-04: Validate query string parameter before use.
	// Empty string is allowed (default "World" is set above only when query is absent).
	// Explicit empty ?name= would also get "World" from the default, so validation
	// runs on the actual resolved value.
	if err := utils.ValidateQueryParam(name); err != nil {
		return utils.BadRequestHandler(c, "Invalid query parameter: name must not exceed 200 characters")
	}

	return c.JSON(fiber.Map{
		"message": "Hello, " + name + "!",
	})
}

func responseHelloWithService(c *fiber.Ctx) error {
	name := c.Query("name", "World")

	// FOUND-04: Validate query string parameter before use.
	if err := utils.ValidateQueryParam(name); err != nil {
		return utils.BadRequestHandler(c, "Invalid query parameter: name must not exceed 200 characters")
	}

	greetingMessage := generateGreeting(name)
	return c.JSON(fiber.Map{
		"message": greetingMessage,
	})
}

func responseHelloWithForm(c *fiber.Ctx) error {
	var req Request

	// FOUND-05: Return 400 on body parse failure — never return the raw error or a silent 200.
	// Using fiber.NewError prevents the raw BodyParser error message from reaching the client
	// (Fiber would otherwise wrap it as-is at app.go:1095).
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid or malformed request body")
	}

	name := req.Name
	if name == "" {
		name = "World"
	}

	return c.JSON(Response{
		Message: "Hello, " + name + "! This message was generated from the form handler.",
	})
}
