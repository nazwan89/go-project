<!-- GSD:project-start source:PROJECT.md -->
## Project

**Go Based Project**

A production-ready Go/Fiber starter template that teams can clone as the baseline for any new API project. It provides the structural patterns, tooling, and infrastructure concerns already solved — so new projects start with a solid foundation rather than from scratch.

**Core Value:** Every new API project should start with security, observability, and testability already in place — not bolted on later.

### Constraints

- **Tech Stack**: Go + Fiber — do not introduce alternative frameworks
- **Simplicity**: Baseline must remain easy to understand; avoid over-engineering abstractions
- **Clone-ready**: No project-specific code in the final template; all examples use generic `sample` module
<!-- GSD:project-end -->

<!-- GSD:stack-start source:codebase/STACK.md -->
## Technology Stack

## Languages
- Go 1.25.4 - Full backend application; all production code is written in Go
## Runtime
- Go 1.25.4 runtime
- Alpine Linux 3.x (containerized deployment)
- Go modules (`go.mod`, `go.sum`)
- Lockfile: Present (`go.sum`)
## Frameworks
- Fiber v2.52.12 - HTTP web framework for building REST APIs and handling HTTP requests/responses
- Location: Used throughout `cmd/main.go` and all route handlers in `module/sample/`
- Fiber's built-in middleware suite:
- Configuration in `cmd/main.go` lines 47-49
- Air v1.x - Live code reloading during development
- Config: `.air.toml`
- Purpose: Automatic rebuild and restart on file changes
## Key Dependencies
- `github.com/gofiber/fiber/v2` v2.52.12 - Core HTTP framework and request/response handling
- `github.com/joho/godotenv` v1.5.1 - Environment variable loading from `.env` files for local development
- `github.com/google/uuid` v1.6.0 - UUID generation (imported by Fiber)
- `github.com/valyala/fasthttp` v1.51.0 - High-performance HTTP client used by Fiber
- `github.com/valyala/bytebufferpool` v1.0.0 - Memory pooling for byte buffers (Fiber dependency)
- `github.com/klauspost/compress` v1.17.9 - HTTP compression support (Brotli and others)
- `github.com/andybalholm/brotli` v1.1.0 - Brotli compression algorithm
- `golang.org/x/sys` v0.42.0 - System-level Go libraries
- Standard library packages:
## Configuration
- `.env` file (optional) - Local development configuration
- Environment variable `PORT` - Server port configuration (default: 8080)
- Location: `config/app_config.go` via `config.LoadAppConfig()`, called at `cmd/main.go` line 32
- In production: Set environment variables through hosting provider or container orchestration
- `Dockerfile` - Multi-stage build for Alpine Linux container
- `.air.toml` - Hot reload configuration for development
- `go.mod` - Module definition and dependency management
## Platform Requirements
- Go 1.25.4 or compatible version
- Optional: Air for hot reloading (`go install github.com/air-verse/air@latest`)
- macOS, Linux, or Windows with Go installed
- Docker with Alpine Linux 3.x base image
- Container port: 8080 (configurable via `PORT` environment variable)
- Timezone: Set to Asia/Kuala_Lumpur in Dockerfile (lines 4-5, 19-20)
- Health check endpoint available at `/` (root)
- Minimal: Alpine Linux base (approximately 6-7 MB)
- No external dependencies beyond Go standard library and Fiber framework
- Single compiled binary: `project`
<!-- GSD:stack-end -->

<!-- GSD:conventions-start source:CONVENTIONS.md -->
## Conventions

## Naming Patterns
- **Files:** lowercase, underscore-separated (`error_handler.go`, `date_time.go`)
- **Functions:** camelCase; handlers prefixed with `response` (responseHello, responseHelloWithName); public functions uppercase (ErrorHandler, CurrentTimestamp)
- **Variables:** camelCase
- **Types:** PascalCase (`Request`, `Response`)
## Code Style
- Standard Go `gofmt` formatting, tab indentation
- Struct tags: JSON first, then form — `json:"name" form:"name"`
## Import Organization
## Error Handling
- Centralized via `utils.ErrorHandler` passed to Fiber config
- Specific handlers: `NotFoundHandler` (404), `MethodNotAllowedHandler` (405), `InternalServerErrorHandler` (500), `BadRequestHandler` (400)
- Type assertion pattern: `if e, ok := err.(*fiber.Error); ok`
- Consistent response format with timestamp and path
## Logging
- Standard Go `log` package — `log.Printf()` for startup, `log.Fatalf()` for fatal errors
- Fiber's `logger.New()` middleware for HTTP requests
## Module Design
- `handlers.go` — HTTP handlers (private functions)
- `service.go` — Business logic (private functions)
- `routes.go` — Route registration (public via `RegisterRoutes`)
- `form.go` — Request/Response structs (exported)
<!-- GSD:conventions-end -->

<!-- GSD:architecture-start source:ARCHITECTURE.md -->
## Architecture

## Pattern Overview
- Module-based organization with clear separation of concerns
- Request routing through Fiber framework with middleware pipeline
- Service layer abstraction for business logic
- Centralized error handling and utility functions
- Environment-based configuration
## Layers
- Purpose: Handle HTTP requests and responses via Fiber framework
- Location: `cmd/main.go` (route registration), `module/sample/handlers.go` (request handlers)
- Contains: HTTP route definitions, request handling functions, middleware configuration
- Depends on: Fiber framework, service layer, utils
- Used by: Client applications making HTTP requests
- Purpose: Encapsulate business logic and reusable application functionality
- Location: `module/sample/service.go`
- Contains: Pure business logic functions (generateGreeting)
- Depends on: No external dependencies except Go standard library
- Used by: Controllers to process data and execute business operations
- Purpose: Define and group API endpoints for modules
- Location: `module/sample/routes.go`
- Contains: Route definitions with HTTP methods and handler mappings
- Depends on: Fiber framework, controllers
- Used by: Main application to register module routes
- Purpose: Define request/response data structures
- Location: `module/sample/form.go`
- Contains: Request and Response struct types with JSON/form tags
- Depends on: No external dependencies
- Used by: Controllers for request parsing and response serialization
- Purpose: Provide cross-cutting concerns and infrastructure utilities
- Location: `utils/` directory
- Contains: Error handling, timestamp formatting, middleware setup
- Depends on: Fiber framework, Go standard library
- Used by: Main application, all controllers
## Data Flow
## Key Abstractions
- Purpose: Organize feature-specific code into self-contained packages
- Examples: `module/sample/` (greeting/sample module)
- Pattern: Each module contains routes.go, handlers.go, service.go, form.go for feature isolation
- Purpose: Provide strongly-typed data contracts for API endpoints
- Examples: `Request{Name string}`, `Response{Message string}` in `module/sample/form.go`
- Pattern: JSON struct tags for request body parsing, form tags for HTML form parsing
- Purpose: Centralize HTTP error response formatting across application
- Examples: NotFoundHandler, MethodNotAllowedHandler, InternalServerErrorHandler, BadRequestHandler
- Pattern: Specific handler functions for each HTTP error code, returns consistently formatted error JSON
## Entry Points
- Location: `cmd/main.go`
- Triggers: `go run ./cmd` or compiled binary execution
- Responsibilities: Load environment, configure Fiber app, register middleware, register module routes, start server on configurable port
- Location: `module/sample/routes.go` - RegisterRoutes() function
- Triggers: Called from cmd/main.go during app initialization
- Responsibilities: Register all routes for sample module under /api/sample namespace
- Location: `module/sample/handlers.go` - HTTP handler functions
- Triggers: When HTTP request matches route (e.g., GET /api/sample/hello)
- Responsibilities: Extract request data, call service layer, return JSON response
## Error Handling
- Fiber error type checking using type assertion (if e, ok := err.(*fiber.Error))
- Dedicated handlers for 404, 405, 500, 400 errors
- All error responses include timestamp via `CurrentTimestamp()`
- All error responses return JSON with consistent structure: error, message, code/path, timestamp
```json
```
## Cross-Cutting Concerns
<!-- GSD:architecture-end -->

<!-- GSD:skills-start source:skills/ -->
## Project Skills

No project skills found. Add skills to any of: `.claude/skills/`, `.agents/skills/`, `.cursor/skills/`, or `.github/skills/` with a `SKILL.md` index file.
<!-- GSD:skills-end -->

<!-- GSD:workflow-start source:GSD defaults -->
## GSD Workflow Enforcement

Before using Edit, Write, or other file-changing tools, start work through a GSD command so planning artifacts and execution context stay in sync.

Use these entry points:
- `/gsd-quick` for small fixes, doc updates, and ad-hoc tasks
- `/gsd-debug` for investigation and bug fixing
- `/gsd-execute-phase` for planned phase work

Do not make direct repo edits outside a GSD workflow unless the user explicitly asks to bypass it.
<!-- GSD:workflow-end -->



<!-- GSD:profile-start -->
## Developer Profile

> Profile not yet configured. Run `/gsd-profile-user` to generate your developer profile.
> This section is managed by `generate-claude-profile` -- do not edit manually.
<!-- GSD:profile-end -->
