package utils

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

// ErrorHandler handles all application errors including 404 and 405
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Check if it's a Fiber error
	if e, ok := err.(*fiber.Error); ok {
		switch e.Code {
		case 404:
			return NotFoundHandler(c)
		case 405:
			return MethodNotAllowedHandler(c)
		default:
			return c.Status(e.Code).JSON(fiber.Map{
				"error":   e.Error(),
				"message": e.Message,
				"code":    e.Code,
				"timestamp": CurrentTimestamp(),
			})
		}
	}
	
	// Generic error
	return InternalServerErrorHandler(c, err)
}

// NotFoundHandler handles 404 errors
func NotFoundHandler(c *fiber.Ctx) error {
	return c.Status(404).JSON(fiber.Map{
		"error":   "Endpoint not found",
		"message": "The requested endpoint does not exist",
		"path":    c.Path(),
		"timestamp": CurrentTimestamp(),
	})
}

// MethodNotAllowedHandler handles 405 errors
func MethodNotAllowedHandler(c *fiber.Ctx) error {
	return c.Status(405).JSON(fiber.Map{
		"error":   "Method Not Allowed",
		"message": fmt.Sprintf("%s method is not allowed for this endpoint", c.Method()),
		"path":    c.Path(),
	})
}

// InternalServerErrorHandler handles 500 errors
func InternalServerErrorHandler(c *fiber.Ctx, err error) error {
	return c.Status(500).JSON(fiber.Map{
		"error":   "Internal Server Error",
		"message": err.Error(),
	})
}

// BadRequestHandler handles 400 errors
func BadRequestHandler(c *fiber.Ctx, message string) error {
	return c.Status(400).JSON(fiber.Map{
		"error":   "Bad Request",
		"message": message,
	})
}