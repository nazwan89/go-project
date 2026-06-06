# Conventions

## Naming Patterns

- **Files:** lowercase, underscore-separated (`error_handler.go`, `date_time.go`)
- **Functions:** camelCase; handlers prefixed with `response` (responseHello, responseHelloWithName); public functions uppercase (ErrorHandler, CurrentTimestamp)
- **Variables:** camelCase
- **Types:** PascalCase (`Request`, `Response`)

## Code Style

- Standard Go `gofmt` formatting, tab indentation
- Struct tags: JSON first, then form — `json:"name" form:"name"`

## Import Organization

1. Standard library (`time`, `fmt`, `log`)
2. Third-party (`github.com/gofiber/*`, `github.com/joho/*`)
3. Local project (`project/module/*`, `project/utils`)

## Error Handling

- Centralized via `utils.ErrorHandler` passed to Fiber config
- Specific handlers: `NotFoundHandler` (404), `MethodNotAllowedHandler` (405), `InternalServerErrorHandler` (500), `BadRequestHandler` (400)
- Type assertion pattern: `if e, ok := err.(*fiber.Error); ok`
- Consistent response format with timestamp and path

## Logging

- Standard Go `log` package — `log.Printf()` for startup, `log.Fatalf()` for fatal errors
- Fiber's `logger.New()` middleware for HTTP requests

## Module Design

Each feature module follows this layout:
- `handlers.go` — HTTP handlers (private functions)
- `service.go` — Business logic (private functions)
- `routes.go` — Route registration (public via `RegisterRoutes`)
- `form.go` — Request/Response structs (exported)
