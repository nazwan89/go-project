# Technology Stack

**Analysis Date:** 2026-04-08

## Languages

**Primary:**
- Go 1.25.4 - Full backend application; all production code is written in Go

## Runtime

**Environment:**
- Go 1.25.4 runtime
- Alpine Linux 3.x (containerized deployment)

**Package Manager:**
- Go modules (`go.mod`, `go.sum`)
- Lockfile: Present (`go.sum`)

## Frameworks

**Core:**
- Fiber v2.52.12 - HTTP web framework for building REST APIs and handling HTTP requests/responses
- Location: Used throughout `main.go` and all route handlers in `module/sample/`

**Middleware:**
- Fiber's built-in middleware suite:
  - `recover` - Panic recovery middleware for graceful error handling
  - `logger` - HTTP request logging middleware
  - `healthcheck` - Health check middleware for service monitoring
- Configuration in `main.go` lines 37-39

**Hot Reload (Development):**
- Air v1.x - Live code reloading during development
- Config: `.air.toml`
- Purpose: Automatic rebuild and restart on file changes

## Key Dependencies

**Critical:**
- `github.com/gofiber/fiber/v2` v2.52.12 - Core HTTP framework and request/response handling
- `github.com/joho/godotenv` v1.5.1 - Environment variable loading from `.env` files for local development

**Infrastructure:**
- `github.com/google/uuid` v1.6.0 - UUID generation (imported by Fiber)
- `github.com/valyala/fasthttp` v1.51.0 - High-performance HTTP client used by Fiber
- `github.com/valyala/bytebufferpool` v1.0.0 - Memory pooling for byte buffers (Fiber dependency)
- `github.com/klauspost/compress` v1.17.9 - HTTP compression support (Brotli and others)
- `github.com/andybalholm/brotli` v1.1.0 - Brotli compression algorithm
- `golang.org/x/sys` v0.28.0 - System-level Go libraries

**Build/Dev:**
- Standard library packages:
  - `log` - Application logging (used in `main.go`)
  - `os` - Environment variable access
  - `time` - Timestamp utilities (used in `utils/date_time.go`)
  - `fmt` - String formatting (used in error handlers)

## Configuration

**Environment:**
- `.env` file (optional) - Local development configuration
- Environment variable `PORT` - Server port configuration (default: 8080)
- Location: `main.go` lines 61-64
- In production: Set environment variables through hosting provider or container orchestration

**Build:**
- `Dockerfile` - Multi-stage build for Alpine Linux container
- `.air.toml` - Hot reload configuration for development
- `go.mod` - Module definition and dependency management

## Platform Requirements

**Development:**
- Go 1.25.4 or compatible version
- Optional: Air for hot reloading (`go install github.com/air-verse/air@latest`)
- macOS, Linux, or Windows with Go installed

**Production:**
- Docker with Alpine Linux 3.x base image
- Container port: 8080 (configurable via `PORT` environment variable)
- Timezone: Set to Asia/Kuala_Lumpur in Dockerfile (lines 4-5, 19-20)
- Health check endpoint available at `/` (root)

**Runtime Requirements:**
- Minimal: Alpine Linux base (approximately 6-7 MB)
- No external dependencies beyond Go standard library and Fiber framework
- Single compiled binary: `project`

---

*Stack analysis: 2026-04-08*
