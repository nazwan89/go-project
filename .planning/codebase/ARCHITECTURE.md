# Architecture

**Analysis Date:** 2026-04-08

## Pattern Overview

**Overall:** Modular MVC-style architecture with Fiber web framework

**Key Characteristics:**
- Module-based organization with clear separation of concerns
- Request routing through Fiber framework with middleware pipeline
- Service layer abstraction for business logic
- Centralized error handling and utility functions
- Environment-based configuration

## Layers

**Presentation Layer (HTTP):**
- Purpose: Handle HTTP requests and responses via Fiber framework
- Location: `main.go` (route registration), `module/sample/controller.go` (request handlers)
- Contains: HTTP route definitions, request handling functions, middleware configuration
- Depends on: Fiber framework, service layer, utils
- Used by: Client applications making HTTP requests

**Service Layer:**
- Purpose: Encapsulate business logic and reusable application functionality
- Location: `module/sample/service.go`
- Contains: Pure business logic functions (generateGreeting)
- Depends on: No external dependencies except Go standard library
- Used by: Controllers to process data and execute business operations

**Route Registration Layer:**
- Purpose: Define and group API endpoints for modules
- Location: `module/sample/routes.go`
- Contains: Route definitions with HTTP methods and handler mappings
- Depends on: Fiber framework, controllers
- Used by: Main application to register module routes

**Model/Contract Layer:**
- Purpose: Define request/response data structures
- Location: `module/sample/form.go`
- Contains: Request and Response struct types with JSON/form tags
- Depends on: No external dependencies
- Used by: Controllers for request parsing and response serialization

**Utility/Infrastructure Layer:**
- Purpose: Provide cross-cutting concerns and infrastructure utilities
- Location: `utils/` directory
- Contains: Error handling, timestamp formatting, middleware setup
- Depends on: Fiber framework, Go standard library
- Used by: Main application, all controllers

## Data Flow

**Request Processing Flow:**

1. HTTP request arrives at Fiber application
2. Middleware pipeline processes request (recover, logger, healthcheck)
3. Router matches request to registered route handler (controller)
4. Controller receives Fiber context and processes request
5. Controller calls service layer function if business logic needed
6. Service layer returns processed data
7. Controller formats response using Response struct
8. ErrorHandler catches any errors and formats error response
9. HTTP response sent to client with timestamp via utils

**Example: GET /api/sample/hello-service?name=John**

1. Fiber routing matches to `responseHelloWithService` handler
2. Handler extracts query parameter "name"
3. Handler calls `generateGreeting(name)` from service layer
4. Service returns formatted greeting string
5. Handler wraps response in fiber.Map with JSON serialization
6. Response sent with status 200 and JSON body

## Key Abstractions

**Module Structure:**
- Purpose: Organize feature-specific code into self-contained packages
- Examples: `module/sample/` (greeting/sample module)
- Pattern: Each module contains routes.go, controller.go, service.go, form.go for feature isolation

**Request/Response Types:**
- Purpose: Provide strongly-typed data contracts for API endpoints
- Examples: `Request{Name string}`, `Response{Message string}` in `module/sample/form.go`
- Pattern: JSON struct tags for request body parsing, form tags for HTML form parsing

**Error Handler Abstraction:**
- Purpose: Centralize HTTP error response formatting across application
- Examples: NotFoundHandler, MethodNotAllowedHandler, InternalServerErrorHandler, BadRequestHandler
- Pattern: Specific handler functions for each HTTP error code, returns consistently formatted error JSON

## Entry Points

**Application Entry Point:**
- Location: `main.go`
- Triggers: `go run main.go` or compiled binary execution
- Responsibilities: Load environment, configure Fiber app, register middleware, register module routes, start server on configurable port

**Module Entry Point:**
- Location: `module/sample/routes.go` - RegisterRoutes() function
- Triggers: Called from main.go during app initialization
- Responsibilities: Register all routes for sample module under /api/sample namespace

**Request Entry Point:**
- Location: `module/sample/controller.go` - HTTP handler functions
- Triggers: When HTTP request matches route (e.g., GET /api/sample/hello)
- Responsibilities: Extract request data, call service layer, return JSON response

## Error Handling

**Strategy:** Centralized error handler with specific handlers for common HTTP errors

**Patterns:**
- Fiber error type checking using type assertion (if e, ok := err.(*fiber.Error))
- Dedicated handlers for 404, 405, 500, 400 errors
- All error responses include timestamp via `CurrentTimestamp()`
- All error responses return JSON with consistent structure: error, message, code/path, timestamp

**Error Response Format:**
```json
{
  "error": "Error type",
  "message": "Detailed message",
  "timestamp": "2006-01-02 15:04:05"
}
```

## Cross-Cutting Concerns

**Logging:** Fiber logger middleware (logger.New()) automatically logs HTTP requests; configured in main.go middleware stack

**Validation:** Request validation occurs in controller via `c.BodyParser()` for JSON bodies; returns BadRequest if parsing fails

**Authentication:** Not currently implemented; ready for middleware integration via `app.Use()` pattern

**Error Handling:** Centralized via `ErrorHandler` passed to Fiber config; catches all unhandled errors

**Timestamps:** All responses include Asia/Kuala_Lumpur timezone timestamp via `utils.CurrentTimestamp()`

---

*Architecture analysis: 2026-04-08*
