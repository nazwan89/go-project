<!-- generated-by: gsd-doc-writer -->
# Architecture

## System Overview

`go-project` is a production-ready Go/Fiber REST API starter template. It accepts HTTP requests through the Fiber v2 framework, routes them through a module-based handler and service pipeline, and returns JSON responses with consistent error formatting. The primary architectural style is a layered, module-scoped design: each feature domain lives in its own package under `module/`, with a strict separation between routing, HTTP handling, business logic, and request/response types. Cross-cutting concerns (error handling, input validation, timestamp formatting, configuration) are centralised in `utils/` and `config/` packages that all modules share.

## Component Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                          main.go                            │
│  LoadAppConfig → Fiber App → Middleware → Routes → Listen   │
└────────────┬──────────────────────────────────┬────────────┘
             │                                  │
             ▼                                  ▼
    ┌─────────────────┐               ┌──────────────────┐
    │   config/       │               │   module/sample/ │
    │  AppConfig      │               │  RegisterRoutes  │
    └─────────────────┘               └────────┬─────────┘
                                               │
                          ┌────────────────────┼────────────────────┐
                          ▼                    ▼                     ▼
                   controller.go          service.go             form.go
                  (HTTP handlers)      (business logic)    (Request/Response)
                          │
                          ▼
                   ┌─────────────┐
                   │   utils/    │
                   │ ErrorHandler│
                   │ Validate    │
                   │ CurrentTime │
                   └─────────────┘
```

## Data Flow

A typical request moves through the system in the following sequence:

1. **Entry** — `main.go` starts Fiber on the configured port and registers three middleware layers: `recover.New()` (panic recovery), `logger.New()` (HTTP access logging), and `healthcheck.New()` (built-in liveness endpoint at `GET /livez` and `GET /readyz`).
2. **Routing** — The request is matched against routes registered under the `/api` group. `sample.RegisterRoutes(api)` registers all `/api/sample/*` routes during startup.
3. **Handler** — The matched handler function in `controller.go` extracts path params, query params, or body data from the `fiber.Ctx`. Input is validated with `utils.ValidatePathParam` or `utils.ValidateQueryParam` before any business logic runs.
4. **Service** — For routes that require non-trivial logic, the handler calls a function in `service.go` (e.g., `generateGreeting`). Service functions are pure — they take plain values and return plain values with no Fiber dependency.
5. **Response** — The handler serialises the result to JSON via `c.JSON()`. For typed responses, a `Response` struct from `form.go` is used directly.
6. **Error path** — Any `error` returned from a handler reaches `utils.ErrorHandler`, which is set as `fiber.Config.ErrorHandler`. It inspects the `*fiber.Error` type and dispatches to the appropriate specific handler (`NotFoundHandler`, `MethodNotAllowedHandler`, `BadRequestHandler`, or `InternalServerErrorHandler`). The 500 path logs the real error with a UUID `request_id` internally and returns only a generic message to the caller.

## Key Abstractions

| Abstraction | File | Description |
|---|---|---|
| `AppConfig` | `config/app_config.go` | Typed configuration struct; loaded once in `main.go` from env vars with safe defaults |
| `LoadAppConfig()` | `config/app_config.go` | Single entry point for all environment-driven configuration |
| `ErrorHandler` | `utils/error_handler.go` | Fiber-compatible top-level error handler; dispatches to specific 4xx/5xx handlers |
| `NotFoundHandler` | `utils/error_handler.go` | Returns a consistent 404 JSON body with `path` and `timestamp` |
| `MethodNotAllowedHandler` | `utils/error_handler.go` | Returns a 405 JSON body with the attempted method name |
| `InternalServerErrorHandler` | `utils/error_handler.go` | Logs real error with UUID `request_id`; returns sanitised 500 JSON to client |
| `BadRequestHandler` | `utils/error_handler.go` | Returns a 400 JSON body with a caller-supplied message string |
| `CurrentTimestamp()` | `utils/date_time.go` | Returns formatted timestamp string using a `*time.Location` (timezone-aware) |
| `ValidatePathParam()` | `utils/validate.go` | Enforces alphanumeric, 1-100 char constraint on path parameters |
| `ValidateQueryParam()` | `utils/validate.go` | Enforces max 200 char constraint on query string parameters |
| `RegisterRoutes()` | `module/sample/routes.go` | Public function that wires all sample module routes onto a Fiber router group |
| `Request` / `Response` | `module/sample/form.go` | Typed request body and response envelope for the sample module |

## Directory Structure Rationale

```
go-project/
├── main.go                  # Application entry point: config, Fiber setup, middleware, server lifecycle
├── config/
│   └── app_config.go        # Typed AppConfig struct and LoadAppConfig(); all env-var reads live here
├── module/
│   └── sample/              # Self-contained feature module; clone this directory to add new domains
│       ├── routes.go        # RegisterRoutes() — route-to-handler mapping; the only public function
│       ├── controller.go    # Private HTTP handler functions; call service layer and return JSON
│       ├── service.go       # Private business logic functions; no Fiber dependency
│       └── form.go          # Exported Request and Response struct types with JSON/form tags
├── utils/
│   ├── error_handler.go     # Centralised error dispatch and per-code JSON response helpers
│   ├── date_time.go         # Timezone-aware timestamp formatting
│   └── validate.go          # Input validation helpers (path params, query params)
├── test/
│   ├── config/              # Black-box tests for config package
│   └── utils/               # Black-box tests for utils package
├── Dockerfile               # Multi-stage build: golang:alpine builder → alpine:latest runtime
├── .air.toml                # Air hot-reload configuration for local development
├── go.mod                   # Module definition: module path "project", Go 1.25.4
└── go.sum                   # Dependency lockfile
```

Each module under `module/` follows a fixed four-file layout (`routes.go`, `controller.go`, `service.go`, `form.go`). This convention makes it immediately obvious where to find routing, handling, logic, and types for any feature domain without reading file contents first. The `utils/` and `config/` packages are intentionally flat — they hold only cross-cutting infrastructure concerns and must not contain feature-specific business logic.

## Middleware Pipeline

Requests pass through the following middleware stack in order before reaching a route handler:

1. `recover.New()` — catches any panic in a handler and converts it to a 500 error routed through `ErrorHandler`
2. `logger.New()` — writes a structured HTTP access log line per request
3. `healthcheck.New()` — intercepts `GET /livez` and `GET /readyz` without reaching the route tree

## Security Constraints Built In

- **Body size limit** — `BodyLimit: 1MB` on the Fiber app prevents memory exhaustion from oversized request bodies. Fiber routes 413 through `ErrorHandler`, which remaps it to `BadRequestHandler` returning HTTP 400.
- **Read timeout** — `ReadTimeout: 30s` bounds the time a slow client can hold a keepalive connection open.
- **Error sanitisation** — `InternalServerErrorHandler` never forwards the real error message to the client. Instead it logs the error with a UUID `request_id` and returns only a generic message.
- **Input validation** — All path and query parameters pass through `ValidatePathParam` / `ValidateQueryParam` before being used in logic, protecting against oversized or malformed input.

## Deployment Model

The project is shipped as a single statically-linked binary built inside a multi-stage Docker image. The builder stage (`golang:1.25.4-alpine`) compiles the binary with `CGO_ENABLED=0`; the runtime stage (`alpine:latest`) copies only the compiled binary. The resulting image has no Go toolchain, no source code, and a minimal attack surface. The timezone is set to `Asia/Kuala_Lumpur` in the Dockerfile and can be overridden via the `TIMEZONE` environment variable at runtime.

<!-- VERIFY: Production deployment target, orchestration platform, and registry details -->
