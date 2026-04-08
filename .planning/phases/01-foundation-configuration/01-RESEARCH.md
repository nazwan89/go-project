# Phase 1: Foundation & Configuration - Research

**Researched:** 2026-04-08
**Domain:** Go/Fiber v2 configuration, graceful shutdown, input validation, error sanitization
**Confidence:** HIGH

---

## Summary

Phase 1 establishes the foundational robustness of the Go/Fiber starter template. The codebase already has a working skeleton: Fiber v2.52.12, godotenv, a centralized `utils/error_handler.go`, and a sample module. What is missing is: externalized config beyond PORT (APP_NAME, TIMEZONE are not read from env), graceful shutdown (app.Listen blocks forever with no signal handling), body size limiting (no BodyLimit set in fiber.Config), input validation helpers (path params and query strings are used raw), body parser error handling (the sample controller currently returns a silent 200 on parse failure), and error sanitization (InternalServerErrorHandler leaks `err.Error()` directly to the client).

All seven requirements in this phase are achievable with the existing dependencies plus one addition (`github.com/go-playground/validator/v10`). No framework changes are needed. The changes are localized: `main.go`, `utils/error_handler.go`, a new `utils/validate.go`, a new `config/app_config.go`, and `.env`/`.env.example` files.

**Primary recommendation:** Add `go-playground/validator/v10` for input validation. All other requirements are solvable with Go stdlib (`os/signal`, `context`) and the Fiber v2 API already present in the project.

---

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| FOUND-01 | App config loaded from env vars (APP_NAME, TIMEZONE, PORT) — no hardcoded values | Config struct pattern; godotenv already present |
| FOUND-02 | Graceful shutdown on SIGTERM — drain in-flight requests before exit (30s timeout) | Fiber `ShutdownWithTimeout()` + `os/signal` stdlib |
| FOUND-03 | Request body size limit (1MB) — prevent memory exhaustion DoS | `fiber.Config{BodyLimit: 1 * 1024 * 1024}` |
| FOUND-04 | Input validation helpers for path params and query strings (length, format) | `go-playground/validator/v10` v10.30.2 |
| FOUND-05 | Body parser always returns 400 on parse failure (not silent 200) | Explicit `c.BodyParser` error handling pattern |
| SEC-03 | Error response sanitization — internal errors return generic message + request_id only (no stack traces) | `SafeError` pattern in `InternalServerErrorHandler` |
| DOC-01 | `.env.example` file with all required environment variables and descriptions | Documented convention — static file |
</phase_requirements>

---

## Project Constraints (from CLAUDE.md)

- **Framework:** Fiber v2 only — no alternatives
- **Language:** Go 1.25.4
- **Simplicity:** Baseline must remain easy to understand; avoid over-engineering
- **Module naming:** Use generic `sample` module only; no project-specific code
- **File naming:** lowercase, underscore-separated (`error_handler.go`, `date_time.go`)
- **Function naming:** camelCase; handlers prefixed with `response`; public functions uppercase
- **Error handling:** Centralized via `utils.ErrorHandler` passed to `fiber.Config`
- **Error format:** Consistent JSON with timestamp and path
- **Logging:** Standard `log` package for startup; Fiber's `logger.New()` for HTTP (Zap comes in Phase 2)
- **Module design:** `controller.go`, `service.go`, `routes.go`, `form.go` per module
- **Struct tags:** JSON first, then form — `json:"name" form:"name"`
- **Indentation:** Standard Go `gofmt` — tab indentation

---

## Existing Codebase — What Is Already In Place

Understanding what exists prevents duplicating work or introducing conflicts.

### `main.go` — Current state

- `godotenv.Load()` is called (ignoring error) — correct pattern for optional `.env` [VERIFIED: codebase read]
- `fiber.New(fiber.Config{AppName: "Project Name", ErrorHandler: utils.ErrorHandler})` — AppName is hardcoded string [VERIFIED: codebase read]
- `recover.New()`, `logger.New()`, `healthcheck.New()` middleware are registered [VERIFIED: codebase read]
- `PORT` is read from env with default `"8080"` — correct [VERIFIED: codebase read]
- `app.Listen(":" + port)` blocks; no signal handling or graceful shutdown [VERIFIED: codebase read]
- No `BodyLimit` set in `fiber.Config` — defaults to 4MB [VERIFIED: source code]

### `utils/error_handler.go` — Current state

- `ErrorHandler` routes `*fiber.Error` by code (404, 405, others) to specific handlers [VERIFIED: codebase read]
- `InternalServerErrorHandler` exposes `err.Error()` directly as `"message"` field — **this leaks internal details** [VERIFIED: codebase read]
- `BadRequestHandler(c, message string)` exists and returns correct 400 structure [VERIFIED: codebase read]
- No `request_id` field in any response [VERIFIED: codebase read]

### `module/sample/controller.go` — Current state

- `responseHelloWithForm` calls `c.BodyParser(&req)` and on error returns `c.JSON(Response{Message: "Invalid request body"})` with **no status code override** — this returns 200 [VERIFIED: codebase read]
- Path params (`c.Params("name")`) and query strings (`c.Query("name", "World")`) are used with no length or format validation [VERIFIED: codebase read]

### `.env` — Current state

- Contains only `PORT = 8080` [VERIFIED: codebase read]
- No `.env.example` exists [VERIFIED: directory listing]

---

## Standard Stack

### Core

| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| `github.com/gofiber/fiber/v2` | v2.52.12 | HTTP framework + BodyLimit + ShutdownWithTimeout | Already in project; all Phase 1 features available |
| `github.com/joho/godotenv` | v1.5.1 | Load `.env` file in development | Already in project; standard Go env-file loader |
| `os/signal` | stdlib | SIGTERM/SIGINT signal capture | Go stdlib; no dependency needed |
| `context` | stdlib | 30-second shutdown timeout | Go stdlib; used with ShutdownWithContext |
| `github.com/go-playground/validator/v10` | v10.30.2 | Struct tag-based input validation | De facto standard for Go validation |

[VERIFIED: `go list -m` against project module cache and registry]

### Supporting

| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| `github.com/google/uuid` | v1.6.0 | Generate request_id for error responses | Already present as transitive dep; import directly for SEC-03 |

[VERIFIED: `go list -m all` against project]

### Alternatives Considered

| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| go-playground/validator | Manual regex validation | Validator is well-maintained, struct-tag driven, and already the ecosystem standard. Manual regex would require hand-rolling what validator provides out of the box. |
| `ShutdownWithTimeout` | `ShutdownWithContext` | Both achieve the same result. `ShutdownWithTimeout` is simpler API for this use case (no context to pass). Either is correct. |

### Installation

```bash
go get github.com/go-playground/validator/v10@latest
```

---

## Architecture Patterns

### Recommended Project Structure After Phase 1

```
.
├── main.go                    # Config load, app init, graceful shutdown
├── utils/
│   ├── (config moved to config/app_config.go)
│   ├── error_handler.go       # MODIFIED: sanitize InternalServerError, add request_id
│   ├── validate.go            # NEW: ValidatePathParam(), ValidateQueryParam() helpers
│   └── date_time.go           # unchanged
├── module/
│   └── sample/
│       └── controller.go      # MODIFIED: fix body parser error response
├── .env                       # MODIFIED: add APP_NAME, TIMEZONE
└── .env.example               # NEW: document all env vars
```

### Pattern 1: Centralized Config Struct (FOUND-01)

**What:** Load all environment variables once at startup into a typed struct. Pass the struct to components that need it.

**When to use:** Any config value read from the environment. Never call `os.Getenv` in handlers.

**Example:**
```go
// config/app_config.go
package config

import (
    "os"
    "time"
)

type AppConfig struct {
    AppName  string
    Port     string
    Timezone string
    Location *time.Location
}

func LoadConfig() AppConfig {
    appName := os.Getenv("APP_NAME")
    if appName == "" {
        appName = "go-project"
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    tz := os.Getenv("TIMEZONE")
    if tz == "" {
        tz = "UTC"
    }
    loc, err := time.LoadLocation(tz)
    if err != nil {
        loc = time.UTC
    }

    return AppConfig{
        AppName:  appName,
        Port:     port,
        Timezone: tz,
        Location: loc,
    }
}
```

[ASSUMED: struct design; pattern is idiomatic Go. The field set derives from FOUND-01 requirements.]

### Pattern 2: Graceful Shutdown with SIGTERM (FOUND-02)

**What:** Capture OS signals, call `ShutdownWithTimeout` with 30-second limit, then exit cleanly.

**When to use:** Always — any production server must drain in-flight requests.

**Key constraint from source docs:** `ShutdownWithTimeout` does not close keepalive connections. Set `ReadTimeout` in `fiber.Config` to something non-zero so keepalive connections are eventually closed. [VERIFIED: `go doc github.com/gofiber/fiber/v2 App.ShutdownWithTimeout`]

**Example:**
```go
// main.go (replace the blocking app.Listen call)
package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"
)

// Start server in goroutine
go func() {
    if err := app.Listen(":" + cfg.Port); err != nil {
        log.Printf("Server stopped: %v", err)
    }
}()

// Block until signal
quit := make(chan os.Signal, 1)
signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
<-quit

log.Println("Shutting down server...")
if err := app.ShutdownWithTimeout(30 * time.Second); err != nil {
    log.Fatalf("Server forced to shutdown: %v", err)
}
log.Println("Server exited cleanly")
```

[VERIFIED: `go doc github.com/gofiber/fiber/v2 App.ShutdownWithTimeout`; pattern confirmed by Fiber official recipe]

### Pattern 3: Body Size Limit (FOUND-03)

**What:** Set `BodyLimit` in `fiber.Config`. When exceeded, Fiber routes a `fiber.Error` with status 413 through the `ErrorHandler`.

**Critical finding:** The requirement says "oversized requests return 400 Bad Request". However, the HTTP standard uses 413 for payload too large, and Fiber's source routes body-too-large to `ErrRequestEntityTooLarge` (413). The `ErrorHandler` receives this as a `*fiber.Error` with code 413. The ErrorHandler can intercept it and return 400 if that is a firm requirement, or return 413 (semantically more correct). [VERIFIED: Fiber source code `/app.go` lines 1081-1096 and `/helpers.go` line 944]

**Recommended approach:** Intercept the 413 in `ErrorHandler` and return a 400 with a clear message, since the requirement is explicit.

```go
// fiber.Config in main.go
app := fiber.New(fiber.Config{
    AppName:      cfg.AppName,
    ErrorHandler: utils.ErrorHandler,
    BodyLimit:    1 * 1024 * 1024, // 1MB
    ReadTimeout:  30 * time.Second, // required for keepalive shutdown
})
```

```go
// utils/error_handler.go — intercept 413 → return 400
case fiber.StatusRequestEntityTooLarge:
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
        "error":     "Bad Request",
        "message":   "Request body exceeds maximum allowed size of 1MB",
        "timestamp": CurrentTimestamp(),
    })
```

[VERIFIED: Fiber source code; pattern is standard Fiber error handler customization]

### Pattern 4: Input Validation Helpers (FOUND-04)

**What:** A `utils/validate.go` package with helpers that validate path params and query strings for length and format before use. Backed by `go-playground/validator/v10`.

**When to use:** Any handler that uses `c.Params()` or `c.Query()` for values that feed business logic.

**Example:**
```go
// utils/validate.go
package utils

import (
    "github.com/go-playground/validator/v10"
)

var validate = validator.New()

type PathParam struct {
    Value string `validate:"required,min=1,max=100,alphanum"`
}

type QueryParam struct {
    Value string `validate:"omitempty,min=1,max=200"`
}

// ValidatePathParam validates a path parameter value.
// Returns an error with a human-readable message if invalid.
func ValidatePathParam(value string) error {
    p := PathParam{Value: value}
    return validate.Struct(p)
}

// ValidateQueryParam validates a query string parameter value.
func ValidateQueryParam(value string) error {
    if value == "" {
        return nil // optional query params are allowed to be empty
    }
    q := QueryParam{Value: value}
    return validate.Struct(q)
}
```

Usage in controller:
```go
func responseHelloWithName(c *fiber.Ctx) error {
    name := c.Params("name")
    if err := utils.ValidatePathParam(name); err != nil {
        return utils.BadRequestHandler(c, "Invalid path parameter: name must be 1-100 alphanumeric characters")
    }
    // ...
}
```

[VERIFIED: `go-playground/validator/v10` v10.30.2 from registry; struct tag approach from official docs]

### Pattern 5: Body Parser Error Handling (FOUND-05)

**What:** `c.BodyParser` returns an error on parse failure. The handler must check this error and call the error handler — not silently return 200.

**Root cause in current code:** `responseHelloWithForm` calls `c.JSON(Response{...})` directly without `c.Status(400)`, so Fiber sends 200.

**Fix pattern:**
```go
func responseHelloWithForm(c *fiber.Ctx) error {
    var req Request

    if err := c.BodyParser(&req); err != nil {
        // Return the error so ErrorHandler processes it, OR call BadRequestHandler directly
        return fiber.NewError(fiber.StatusBadRequest, "Invalid or malformed request body")
    }

    // ... rest of handler
}
```

Alternatively, using `BadRequestHandler`:
```go
    if err := c.BodyParser(&req); err != nil {
        return utils.BadRequestHandler(c, "Invalid or malformed request body")
    }
```

**Note:** `c.BodyParser` returns a raw error, not a `*fiber.Error`. If returned directly to Fiber, the `ErrorHandler` receives it as a generic error (line 1095 of app.go: `err = NewError(StatusBadRequest, err.Error())`) — which would expose the raw error message. Therefore, always wrap with `fiber.NewError(400, "safe message")` or call `BadRequestHandler` directly. [VERIFIED: Fiber source code `ctx.go` line 454, `app.go` line 1095]

### Pattern 6: Error Sanitization (SEC-03)

**What:** Replace `err.Error()` in `InternalServerErrorHandler` with a generic message. Generate a `request_id` (UUID v4) and include it in the response so logs can be correlated without leaking internals.

**Key principle from JetBrains Go blog (2026-03):** "When crossing public API boundaries, replace all generated error messages with static, pre-defined strings." [CITED: https://blog.jetbrains.com/go/2026/03/02/secure-go-error-handling-best-practices/]

```go
// utils/error_handler.go — updated InternalServerErrorHandler
import "github.com/google/uuid"

func InternalServerErrorHandler(c *fiber.Ctx, err error) error {
    requestID := uuid.New().String()

    // Log the real error internally — never expose it
    log.Printf("[ERROR] request_id=%s error=%v path=%s", requestID, err, c.Path())

    return c.Status(500).JSON(fiber.Map{
        "error":      "Internal Server Error",
        "message":    "An unexpected error occurred. Please contact support if the issue persists.",
        "request_id": requestID,
        "timestamp":  CurrentTimestamp(),
    })
}
```

[ASSUMED: exact log format. The UUID approach is idiomatic and `github.com/google/uuid` is already a transitive dependency in the project.]

### Anti-Patterns to Avoid

- **Calling `os.Getenv` in handlers or service functions:** All config reads must happen at startup via `LoadConfig()`.
- **Returning `c.JSON(...)` directly on body parser error without setting status:** This produces a 200 with an "error" payload — misleading to clients and API consumers.
- **Passing raw `err.Error()` to clients:** Database errors, JSON decode failures, and internal errors contain implementation details. Always use pre-defined messages.
- **Not setting `ReadTimeout` when using `ShutdownWithTimeout`:** The Fiber docs explicitly warn that keepalive connections are not closed by shutdown — a `ReadTimeout` is required to bound how long they linger. [VERIFIED: `go doc` output]
- **Panicking on invalid TIMEZONE:** `time.LoadLocation` can fail. Always fall back to `time.UTC` to prevent startup crash.

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Struct field validation (length, format, required) | Custom string length checks with `if len(s) > 100` | `go-playground/validator/v10` struct tags | Handles nested structs, cross-field validation, i18n error messages; battle-tested in production |
| UUID generation for request_id | Custom random string generators | `github.com/google/uuid` (already a dep) | UUID v4 is cryptographically random; custom generators often have collision or bias issues |
| Signal handling | Custom signal loop | `os/signal` + `signal.Notify` | stdlib; well-defined behaviour for SIGTERM and SIGINT |

**Key insight:** The validator library covers an enormous surface area (50+ built-in validators, custom validators, cross-field rules). Any hand-rolled solution will miss edge cases that validator already handles.

---

## Common Pitfalls

### Pitfall 1: BodyLimit Exceeded Closes Connection Instead of Returning 413

**What goes wrong:** On very large uploads, fasthttp closes the TCP connection rather than sending a proper HTTP response. This was a known Fiber v2 issue (GitHub #1940). The error IS routed through `ErrorHandler` in the normal HTTP body parsing path, but at the transport level for extremely large payloads the connection may close.

**Why it happens:** Fasthttp's MaxRequestBodySize enforcement happens at the transport layer. Once the limit is hit, the connection is dropped before a response is fully formed.

**How to avoid:** Set `BodyLimit` in `fiber.Config` (1MB per FOUND-03). For the template, document that clients should check for connection resets as well as 400/413 responses. This is a known limitation of fasthttp-based servers.

**Warning signs:** Clients reporting connection reset errors rather than HTTP error responses for large payloads.

[VERIFIED: Fiber GitHub issue #1940 and source code `app.go` lines 1081-1096]

### Pitfall 2: Body Parser Returns Raw Error Message Through Default Handler

**What goes wrong:** If a handler does `return err` after `c.BodyParser` fails, the Fiber `ErrorHandler` receives a raw Go error. The default `serverErrorHandler` (in `app.go` line 1095) wraps unknown errors as `NewError(StatusBadRequest, err.Error())` — which exposes the raw JSON decode error message to the client.

**Why it happens:** `BodyParser` returns standard Go errors, not `*fiber.Error`. The wrapping in app.go exposes the message.

**How to avoid:** Always return `fiber.NewError(400, "safe message")` or call `BadRequestHandler(c, "safe message")` explicitly. Never `return err` after BodyParser.

[VERIFIED: Fiber source code `app.go` line 1095, `ctx.go` line 454]

### Pitfall 3: Graceful Shutdown Without ReadTimeout Leaves Keepalive Connections Open

**What goes wrong:** `ShutdownWithTimeout(30s)` times out but keepalive connections are still alive, so the process never exits cleanly during the 30-second window.

**Why it happens:** Fiber's `Shutdown*` methods explicitly do not close keepalive connections (documented in godoc).

**How to avoid:** Always set `ReadTimeout` in `fiber.Config` when using graceful shutdown. A value like `30 * time.Second` bounds keepalive connection lifetime.

[VERIFIED: `go doc github.com/gofiber/fiber/v2 App.ShutdownWithTimeout`]

### Pitfall 4: TIMEZONE Env Var Invalid Location Panics on Startup

**What goes wrong:** `time.LoadLocation(os.Getenv("TIMEZONE"))` with an invalid timezone string returns an error. If the error is not handled, startup proceeds with a nil Location and panics on first use.

**Why it happens:** `time.LoadLocation` depends on the system's timezone database. Invalid strings or missing `tzdata` in Alpine Docker builds cause failures.

**How to avoid:** Always fall back to `time.UTC` if `time.LoadLocation` returns an error. In the Dockerfile, `tzdata` is already installed via `apk add --no-cache tzdata`. [VERIFIED: Dockerfile lines 3-5, 18-20]

### Pitfall 5: `AppName` Hardcoded Breaks Template Clone-Readiness

**What goes wrong:** `AppName: "Project Name"` in `fiber.Config` and `"Asia/Kuala_Lumpur"` hardcoded in `date_time.go` mean every team cloning the template must find and replace hardcoded values.

**Why it happens:** Initial skeleton was written before the config externalization requirement.

**How to avoid:** `AppName` must come from `cfg.AppName` (FOUND-01). The timezone in `date_time.go` should use the loaded `Location` from `AppConfig`. Note: Phase 2 (OBS-04) will normalize all timestamps to UTC RFC3339 — so `date_time.go` will be further revised in Phase 2. For Phase 1, externalizing the location read is sufficient.

[VERIFIED: codebase read of `main.go` line 30 and `utils/date_time.go` line 6]

---

## Code Examples

### Graceful Shutdown in main.go (complete pattern)

```go
// Source: go doc github.com/gofiber/fiber/v2 App.ShutdownWithTimeout + stdlib os/signal
package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gofiber/fiber/v2"
    "github.com/joho/godotenv"

    "project/utils"
    "project/module/sample"
)

func main() {
    _ = godotenv.Load()

    cfg := config.LoadConfig()

    app := fiber.New(fiber.Config{
        AppName:      cfg.AppName,
        ErrorHandler: utils.ErrorHandler,
        BodyLimit:    1 * 1024 * 1024, // 1MB (FOUND-03)
        ReadTimeout:  30 * time.Second, // required for graceful shutdown of keepalive
    })

    // middleware ...
    // routes ...

    // Start server non-blocking (FOUND-02)
    go func() {
        log.Printf("Service Starting On Port %s", cfg.Port)
        if err := app.Listen(":" + cfg.Port); err != nil {
            log.Printf("Server stopped: %v", err)
        }
    }()

    // Wait for shutdown signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
    <-quit

    log.Println("Shutting down server (30s timeout)...")
    if err := app.ShutdownWithTimeout(30 * time.Second); err != nil {
        log.Fatalf("Server forced shutdown: %v", err)
    }
    log.Println("Server exited cleanly")
}
```

### Error Sanitization in error_handler.go

```go
// Source: JetBrains Go blog 2026-03; go-playground/validator pattern
import (
    "log"
    "github.com/google/uuid"
)

func InternalServerErrorHandler(c *fiber.Ctx, err error) error {
    requestID := uuid.New().String()
    log.Printf("[ERROR] request_id=%s path=%s error=%v", requestID, c.Path(), err)
    return c.Status(500).JSON(fiber.Map{
        "error":      "Internal Server Error",
        "message":    "An unexpected error occurred. Please contact support if the issue persists.",
        "request_id": requestID,
        "timestamp":  CurrentTimestamp(),
    })
}
```

### .env.example

```bash
# Application
APP_NAME=go-project
TIMEZONE=UTC

# Server
PORT=8080
```

---

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| Hand-rolled env parsing | godotenv + AppConfig struct | Established pattern | Avoids scattered `os.Getenv` calls |
| Blocking `app.Listen` | Goroutine + `ShutdownWithTimeout` | Fiber v2.x | Required for graceful SIGTERM drain |
| No body limit | `fiber.Config{BodyLimit: N}` | Fiber v2 | Prevents memory exhaustion DoS |

**Deprecated/outdated:**
- `io/ioutil.ReadFile` for env files: replaced by `godotenv` (already in use)
- Manual signal handling with `os.Exit(1)` in signal handler: replaced by `ShutdownWithTimeout` + clean exit

---

## Assumptions Log

| # | Claim | Section | Risk if Wrong |
|---|-------|---------|---------------|
| A1 | `AppConfig` struct design (fields: AppName, Port, Timezone, Location) | Architecture Patterns | Planner may add/remove fields; low risk — struct is internal |
| A2 | `request_id` field name in error response (vs `trace_id`, `correlation_id`) | Pattern 6 / SEC-03 | Naming inconsistency with Phase 2 (OBS-02 will add X-Request-ID header); Phase 2 should unify |
| A3 | `uuid.New().String()` acceptable for request_id in Phase 1 before Phase 2 introduces dedicated request ID middleware | Pattern 6 / SEC-03 | Phase 2 OBS-02 will add proper middleware; Phase 1 implementation is temporary until middleware takes over |
| A4 | `ValidatePathParam` uses `alphanum` tag as default path param constraint | Pattern 4 / FOUND-04 | Path params containing hyphens, underscores, or other chars would fail; planner should decide exact allowed character set |

---

## Open Questions

1. **Body size: 400 vs 413 for oversized bodies**
   - What we know: HTTP standard says 413; the requirement says 400; Fiber routes body-too-large through ErrorHandler as a `*fiber.Error` with code 413
   - What's unclear: Whether the team intends to override semantic correctness for consistency, or whether 413 is acceptable
   - Recommendation: Implement as 400 per FOUND-03 requirement; intercept 413 in ErrorHandler and remap to 400

2. **`request_id` in Phase 1 vs Phase 2 correlation ID**
   - What we know: SEC-03 requires request_id in error responses; OBS-02 (Phase 2) will add a proper Request ID middleware that injects correlation ID into every request
   - What's unclear: Should Phase 1 generate a new UUID per error call, or should Phase 2's middleware retroactively replace this with context-stored IDs?
   - Recommendation: Phase 1 generates UUID at the error handler level (simple); Phase 2 refactors to read from Fiber's `Locals` set by the middleware

3. **Timezone in `date_time.go`**
   - What we know: `date_time.go` hardcodes `"Asia/Kuala_Lumpur"`; FOUND-01 requires no hardcoded values; Phase 2 (OBS-04) will normalize to UTC RFC3339
   - What's unclear: Whether to fully fix `date_time.go` in Phase 1 or partially (just read from env) and fully in Phase 2
   - Recommendation: Phase 1 externalizes the location (reads from AppConfig); Phase 2 converts to UTC. Planner should create two tasks if needed.

---

## Environment Availability

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| Go runtime | All | Yes | go1.25.4 darwin/arm64 | — |
| `github.com/gofiber/fiber/v2` | All | Yes | v2.52.12 | — |
| `github.com/joho/godotenv` | FOUND-01 | Yes | v1.5.1 | — |
| `github.com/google/uuid` | SEC-03 | Yes (transitive dep) | v1.6.0 | — |
| `github.com/go-playground/validator/v10` | FOUND-04 | Not yet installed | v10.30.2 latest | None — must `go get` |
| `os/signal`, `syscall` | FOUND-02 | Yes (stdlib) | — | — |

**Missing dependencies with no fallback:**
- `github.com/go-playground/validator/v10` — must be installed before FOUND-04 can be implemented

[VERIFIED: `go list -m all`, `go list -m github.com/go-playground/validator/v10@latest`, `go version`]

---

## Validation Architecture

### Test Framework

| Property | Value |
|----------|-------|
| Framework | None installed yet — Wave 0 gap |
| Config file | None found |
| Quick run command | `go test ./... -run TestFoundation -count=1` |
| Full suite command | `go test ./...` |

### Phase Requirements → Test Map

| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| FOUND-01 | Config loads APP_NAME, TIMEZONE, PORT from env | unit | `go test ./utils/ -run TestLoadConfig` | Wave 0 |
| FOUND-02 | Graceful shutdown drains within 30s | integration | Manual — requires process signal test | Manual only |
| FOUND-03 | Requests > 1MB return 400 | unit (httptest) | `go test ./... -run TestBodyLimit` | Wave 0 |
| FOUND-04 | ValidatePathParam rejects invalid input | unit | `go test ./utils/ -run TestValidatePathParam` | Wave 0 |
| FOUND-05 | Body parser failure returns 400 | unit (httptest) | `go test ./... -run TestBodyParserError` | Wave 0 |
| SEC-03 | Error response has no stack trace, has request_id | unit (httptest) | `go test ./... -run TestErrorSanitization` | Wave 0 |
| DOC-01 | `.env.example` exists with all vars | smoke | Manual file check or `go test ./... -run TestEnvExample` | Wave 0 |

### Sampling Rate

- **Per task commit:** `go test ./...`
- **Per wave merge:** `go test ./...`
- **Phase gate:** Full suite green before `/gsd-verify-work`

### Wave 0 Gaps

- [ ] `config/app_config_test.go` — covers FOUND-01
- [ ] `utils/validate_test.go` — covers FOUND-04
- [ ] `module/sample/controller_test.go` — covers FOUND-03, FOUND-05, SEC-03
- [ ] No test framework config needed (Go stdlib `testing` + `net/http/httptest`)

---

## Security Domain

### Applicable ASVS Categories

| ASVS Category | Applies | Standard Control |
|---------------|---------|-----------------|
| V2 Authentication | No | Out of scope Phase 1 |
| V3 Session Management | No | Out of scope Phase 1 |
| V4 Access Control | No | Out of scope Phase 1 |
| V5 Input Validation | Yes | `go-playground/validator/v10` struct tags |
| V6 Cryptography | No | No crypto in Phase 1 |
| V12 File/Resource | Yes (partial) | BodyLimit = 1MB (prevents DoS via large payload) |

### Known Threat Patterns for Go/Fiber

| Pattern | STRIDE | Standard Mitigation |
|---------|--------|---------------------|
| Large payload DoS | Denial of Service | `fiber.Config{BodyLimit: 1 * 1024 * 1024}` |
| Internal error leakage | Information Disclosure | Sanitize `InternalServerErrorHandler`; never return `err.Error()` |
| Unvalidated path param injection | Tampering | `go-playground/validator/v10` on all path params |
| Stack trace exposure | Information Disclosure | Generic message + request_id only; log internally |
| Uncontrolled JSON decode errors leaking to client | Information Disclosure | Wrap BodyParser errors with `fiber.NewError(400, "safe message")` |

---

## Sources

### Primary (HIGH confidence)

- Fiber v2.52.12 source code (`/app.go`, `/ctx.go`, `/helpers.go`) — BodyLimit behavior, error routing, shutdown methods [VERIFIED: read from `$GOPATH/pkg/mod/github.com/gofiber/fiber/v2@v2.52.12/`]
- `go doc github.com/gofiber/fiber/v2 App.ShutdownWithTimeout` — graceful shutdown API and keepalive warning [VERIFIED: live go doc output]
- `go doc github.com/gofiber/fiber/v2 Config` — BodyLimit default value (4MB), AppName, ErrorHandler fields [VERIFIED: live go doc output]
- Project codebase (`main.go`, `utils/error_handler.go`, `module/sample/controller.go`) — existing state audit [VERIFIED: Read tool]

### Secondary (MEDIUM confidence)

- [Fiber Graceful Shutdown Recipe](https://docs.gofiber.io/recipes/graceful-shutdown/) — signal channel pattern [CITED]
- [JetBrains Go Blog: Secure Error Handling Best Practices (2026-03-02)](https://blog.jetbrains.com/go/2026/03/02/secure-go-error-handling-best-practices/) — SafeError pattern, sanitization at API boundaries [CITED]
- [go-playground/validator v10.30.2 on pkg.go.dev](https://pkg.go.dev/github.com/go-playground/validator/v10) — current version, struct tag validation [VERIFIED: registry lookup]

### Tertiary (LOW confidence)

- GitHub issue #1940 (gofiber/fiber) — connection close behavior on BodyLimit exceed [CITED; fixed in later PR but behavior may still occur in edge cases]

---

## Metadata

**Confidence breakdown:**

- Standard stack: HIGH — all versions verified against live registry and `go list -m`
- Architecture patterns: HIGH for Fiber-specific patterns (verified from source); MEDIUM for struct designs (assumed but idiomatic)
- Pitfalls: HIGH — derived from reading Fiber source code directly, not just documentation

**Research date:** 2026-04-08
**Valid until:** 2026-10-08 (stable stack; Fiber v2 is in maintenance mode, go-playground/validator is stable)
