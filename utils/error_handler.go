package utils

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ErrorHandler handles all application errors including 404 and 405.
func ErrorHandler(c *fiber.Ctx, err error) error {
	if e, ok := err.(*fiber.Error); ok {
		switch e.Code {
		case fiber.StatusNotFound:
			return NotFoundHandler(c)
		case fiber.StatusMethodNotAllowed:
			return MethodNotAllowedHandler(c)
		case fiber.StatusRequestEntityTooLarge:
			// Fiber routes 413 (body too large) through ErrorHandler.
			// Requirement FOUND-03 mandates 400 Bad Request for oversized bodies.
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":     "Bad Request",
				"message":   "Request body exceeds the maximum allowed size of 1MB",
				"timestamp": CurrentTimestamp(time.UTC),
			})
		default:
			return c.Status(e.Code).JSON(fiber.Map{
				"error":     e.Error(),
				"message":   e.Message,
				"code":      e.Code,
				"timestamp": CurrentTimestamp(time.UTC),
			})
		}
	}

	return InternalServerErrorHandler(c, err)
}

// NotFoundHandler handles 404 errors.
func NotFoundHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error":     "Endpoint not found",
		"message":   "The requested endpoint does not exist",
		"path":      c.Path(),
		"timestamp": CurrentTimestamp(time.UTC),
	})
}

// MethodNotAllowedHandler handles 405 errors.
func MethodNotAllowedHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusMethodNotAllowed).JSON(fiber.Map{
		"error":     "Method Not Allowed",
		"message":   fmt.Sprintf("%s method is not allowed for this endpoint", c.Method()),
		"path":      c.Path(),
		"timestamp": CurrentTimestamp(time.UTC),
	})
}

// InternalServerErrorHandler handles 500 errors.
// SEC-03: Never expose err.Error() to the client. Log internally; return generic message + request_id.
// Note: request_id injection is added in Task 1 of Plan 02-02.
func InternalServerErrorHandler(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error":     "Internal Server Error",
		"message":   "An unexpected error occurred. Please contact support if the issue persists.",
		"timestamp": CurrentTimestamp(time.UTC),
	})
}

// BadRequestHandler handles 400 errors.
func BadRequestHandler(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error":     "Bad Request",
		"message":   message,
		"timestamp": CurrentTimestamp(time.UTC),
	})
}
