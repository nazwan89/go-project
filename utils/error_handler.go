package utils

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
			return BadRequestHandler(c, "Request body exceeds the maximum allowed size of 1MB")
		default:
			if e.Code == fiber.StatusInternalServerError {
				return InternalServerErrorHandler(c, err)
			}
			return c.Status(e.Code).JSON(fiber.Map{
				"error":     "Bad Request",
				"message":   e.Message,
				"code":      e.Code,
				"timestamp": CurrentTimestamp(time.UTC),
			})
		}
	}

	// Generic error
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
// SEC-03: Returns generic message + request_id only. Real error is logged internally.
// request_id is a UUID v4; Phase 2 (OBS-02) will replace this with middleware-injected correlation ID.
func InternalServerErrorHandler(c *fiber.Ctx, err error) error {
	requestID := uuid.New().String()
	log.Printf("[ERROR] request_id=%s path=%s error=%v", requestID, c.Path(), err)
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error":      "Internal Server Error",
		"message":    "An unexpected error occurred. Please contact support if the issue persists.",
		"request_id": requestID,
		"timestamp":  CurrentTimestamp(time.UTC),
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
