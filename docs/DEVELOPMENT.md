<!-- generated-by: gsd-doc-writer -->
# Development Guide

This guide covers everything you need to extend the template: adding new modules, following naming and code style conventions, applying the error handling patterns, and understanding what to keep and what to replace when customising for a real project.

---

## Local Setup

### Prerequisites

- Go 1.25.4 or later (`go version`)
- Optional: Air for hot reloading

```bash
go install github.com/air-verse/air@latest
```

### Clone and install

```bash
git clone <your-repo-url>
cd go-project
go mod download
```

### Configure the environment

```bash
cp .env.example .env
# Edit .env as needed — see docs/CONFIGURATION.md for all variables
```

### Run the development server

With hot reload (Air):

```bash
air
```

Without hot reload:

```bash
go run ./cmd
```

The server starts on port `8080` by default. Override with the `PORT` environment variable.

### Verify the server is running

```bash
curl http://localhost:8080/
```

Expected response:

```json
{
  "message": "Service is running",
  "status": "OK",
  "timestamp": "2024-01-01 12:00:00"
}
```

---

## Build Commands

| Command | Description |
|---|---|
| `go run ./cmd` | Run the server without compiling a binary |
| `go build -o ./tmp/main ./cmd` | Compile to `./tmp/main` (same as Air uses) |
| `go test ./...` | Run all tests |
| `go test ./module/sample/...` | Run tests for the sample module only |
| `go test ./test/...` | Run package-level tests in the `test/` directory |
| `go vet ./...` | Static analysis — run before committing |
| `gofmt -w .` | Format all Go files in place |
| `air` | Start dev server with hot reload (requires Air installed) |

---

## Code Style

The project uses standard Go formatting with no additional linter configuration files. Apply these rules consistently:

- **Formatter:** `gofmt` — run `gofmt -w .` before committing. Tab indentation, no trailing spaces.
- **Vet:** run `go vet ./...` to catch common mistakes (incorrect printf formats, unreachable code, etc.).
- **Import organisation:** standard library first, then third-party, then internal packages. Separate each group with a blank line.

```go
import (
    "os"
    "time"

    "github.com/gofiber/fiber/v2"

    "project/utils"
)
```

- **Struct tags:** JSON tag first, then form tag.

```go
type Request struct {
    Name string `json:"name" form:"name"`
}
```

- **Naming:**
  - Files: lowercase, underscore-separated — `error_handler.go`, `date_time.go`
  - Functions: camelCase for unexported, PascalCase for exported
  - HTTP handler functions: prefix with `response` — `responseHello`, `responseHelloWithName`
  - Types: PascalCase — `Request`, `Response`, `AppConfig`
  - Variables: camelCase

- **Logging:** use the standard `log` package. `log.Printf()` for startup messages; `log.Fatalf()` for unrecoverable errors. Do not use third-party logging libraries in the baseline.

---

## Adding a New Module

Each feature domain is a self-contained package under `module/`. Every module contains exactly four files. Use the `sample` module as the reference implementation.

### The 4-file pattern

```
module/
└── yourfeature/
    ├── form.go        # Request and Response structs
    ├── service.go     # Business logic (pure functions)
    ├── handlers.go    # HTTP handlers (calls service, returns JSON)
    └── routes.go      # Route registration (one exported function)
```

### Step 1 — Create the package directory

```bash
mkdir module/yourfeature
```

### Step 2 — Define request/response types (`form.go`)

```go
package yourfeature

type Request struct {
    Name string `json:"name" form:"name"`
}

type Response struct {
    Message string `json:"message"`
}
```

Export all types. Use JSON struct tags on every field. Add `form` tags for endpoints that accept `application/x-www-form-urlencoded` bodies.

### Step 3 — Write business logic (`service.go`)

```go
package yourfeature

func generateMessage(name string) string {
    return "Hello, " + name + "!"
}
```

Service functions are unexported by convention. They must not import Fiber or any HTTP-layer package. Dependencies are passed in as arguments — no global state.

### Step 4 — Write HTTP handlers (`handlers.go`)

```go
package yourfeature

import (
    "github.com/gofiber/fiber/v2"

    "project/utils"
)

func responseIndex(c *fiber.Ctx) error {
    name := c.Query("name", "World")

    if err := utils.ValidateQueryParam(name); err != nil {
        return utils.BadRequestHandler(c, "Invalid query parameter: name must not exceed 200 characters")
    }

    return c.JSON(Response{
        Message: generateMessage(name),
    })
}
```

Handler functions are unexported. The naming convention is `response<Action>`. Always validate path parameters with `utils.ValidatePathParam` and query parameters with `utils.ValidateQueryParam` before using their values. Return errors using the utility handlers from `utils/` — never construct error responses inline.

### Step 5 — Register routes (`routes.go`)

```go
package yourfeature

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(api fiber.Router) {
    group := api.Group("/yourfeature")

    group.Get("/", responseIndex)
}
```

`RegisterRoutes` is the only exported function in a module. The `api` argument is the `/api` group created in `cmd/main.go`.

### Step 6 — Register the module in `cmd/main.go`

Import the new package and call `RegisterRoutes` in the "Register Module Routes" block:

```go
import (
    // ...
    "project/module/yourfeature"
)

// in main():
api := app.Group("/api")
sample.RegisterRoutes(api)
yourfeature.RegisterRoutes(api)   // add this line
```

---

## Error Handling Patterns

All error responses go through `utils.ErrorHandler`, which is registered as Fiber's global error handler in `cmd/main.go`. Do not construct raw JSON error responses in handlers.

### Return a 400 Bad Request

Use `utils.BadRequestHandler` when input validation fails:

```go
if err := utils.ValidatePathParam(name); err != nil {
    return utils.BadRequestHandler(c, "Invalid path parameter: name must be 1-100 alphanumeric characters")
}
```

### Return a 400 on body parse failure

Wrap body parse errors with `fiber.NewError` — do not pass the raw BodyParser error to the client:

```go
if err := c.BodyParser(&req); err != nil {
    return fiber.NewError(fiber.StatusBadRequest, "Invalid or malformed request body")
}
```

### Return a 500 Internal Server Error

Return any non-fiber error from a handler. `utils.ErrorHandler` catches it, logs the real error with a `request_id`, and returns a sanitised response to the client:

```go
result, err := doSomething()
if err != nil {
    return err  // logged internally; client receives a generic 500 with request_id
}
```

Never include internal error details (database messages, stack traces, credentials) in the JSON response. The `InternalServerErrorHandler` in `utils/error_handler.go` enforces this by design.

### Error response shapes

All error responses follow a consistent structure:

**400 Bad Request:**
```json
{
  "error": "Bad Request",
  "message": "Human-readable description of what was wrong",
  "timestamp": "2024-01-01 12:00:00"
}
```

**404 Not Found:**
```json
{
  "error": "Endpoint not found",
  "message": "The requested endpoint does not exist",
  "path": "/api/nonexistent",
  "timestamp": "2024-01-01 12:00:00"
}
```

**500 Internal Server Error:**
```json
{
  "error": "Internal Server Error",
  "message": "An unexpected error occurred. Please contact support if the issue persists.",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2024-01-01 12:00:00"
}
```

---

## Input Validation

Two validation helpers are provided in `utils/validate.go`. Call them in every handler before using parameter values in business logic.

| Helper | Use for | Rules |
|---|---|---|
| `utils.ValidatePathParam(value)` | `c.Params()` values | Required, 1–100 chars, alphanumeric only |
| `utils.ValidateQueryParam(value)` | `c.Query()` values | Optional, max 200 chars if non-empty |

For request bodies, use `c.BodyParser` and handle its error explicitly (see the error handling section above). Add struct-level validation with `github.com/go-playground/validator/v10` for complex body validation.

---

## Configuration

All environment variable reads happen in `config/LoadAppConfig()`. Never call `os.Getenv` directly in handlers or service functions. Pass the `AppConfig` struct as a parameter if a handler needs a config value.

To add a new configuration value:

1. Add a field to `AppConfig` in `config/app_config.go`.
2. Read the environment variable inside `LoadAppConfig()` with a safe default.
3. Add the variable to `.env.example` with a comment.
4. Pass `cfg` into the handler or middleware that needs it.

---

## Testing

### Writing tests

- Module-level tests live alongside the source files: `module/yourfeature/controller_test.go`.
- Package-level tests (config, utils) live under `test/`: `test/config/app_config_test.go`, `test/utils/validate_test.go`.
- Use `package yourfeature` for white-box tests (access to unexported identifiers). Use `package yourfeature_test` or `package config_test` for black-box tests.

### Setting up a test app

Use Fiber's `app.Test()` with `httptest` — no real server required:

```go
func setupApp() *fiber.App {
    app := fiber.New(fiber.Config{
        ErrorHandler: utils.ErrorHandler,
        BodyLimit:    1 * 1024 * 1024,
    })
    api := app.Group("/api")
    RegisterRoutes(api)
    return app
}

func TestMyHandler(t *testing.T) {
    app := setupApp()
    req := httptest.NewRequest(http.MethodGet, "/api/yourfeature/", nil)
    resp, err := app.Test(req, -1)
    // ...
}
```

### Running tests

```bash
go test ./...                          # all tests
go test ./module/sample/...            # sample module only
go test ./test/...                     # package-level tests only
go test -v ./module/sample/...         # verbose output
```

---

## Extending the Template

When using this template as the baseline for a real project:

1. **Rename the module.** Update the module name in `go.mod` from `project` to your actual module path (e.g., `github.com/your-org/your-service`). Update all internal import paths accordingly.
2. **Replace the `sample` module.** The `module/sample/` package is the reference implementation. Delete it or rename it once you have added your own modules.
3. **Keep `utils/` and `config/` intact.** These packages provide the cross-cutting concerns (error handling, validation, configuration, timestamps) that every module depends on. Extend them; do not bypass them.
4. **Do not call `os.Getenv` outside `config/`.** All environment access belongs in `config/app_config.go`.
5. **Do not import Fiber in `service.go`.** The service layer must remain free of HTTP dependencies so it can be tested and reused independently.
