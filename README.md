# go-project

## Install Go 
### Windows
1. Download from https://go.dev/dl/
2. Run the installer.
3. Verify installation using bash: "go version"

### MacOS
```sh
brew install go
go version
```

### Linux
```sh
sudo apt update
sudo apt install golang-go -y
go version
```

## Create Project
### Create Project Folder:
```sh
mkdir 'project name'
cd 'project name'
```

### Initialize a Go Module:
```sh
go mod init 'project name'
```

### Install Fiber:
```sh
go get github.com/gofiber/fiber/v2
```

## Create File
### Create `.env` file
 - Create an `.env` file (optional: use for local development). Go will not read this file automatically; you either need to load it in code with a package such as `github.com/joho/godotenv` or export the variable before starting the app.

```sh
touch .env
```
- Add code
```sh
PORT=8080
```

- To start the server with the value from `.env` you can run:
```sh
go run main.go
```

- or manually export before running:
```sh
export PORT=8080
go run main.go
```

- (see the `main.go` example which uses `godotenv.Load()`)

### Create `main.go` file
```sh
touch main.go
```
- Add Code
```sh
package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"

	"project/module/sample"
	"project/utils"
)

func main() {
	// ========================
	// load environment variables from .env if present
	// Go does _not_ automatically read the file; you must do this yourself
	// or export the variables before running.
	// In production, you should set environment variables through your hosting provider or container orchestration system.
	// ========================
	_ = godotenv.Load() // ignore error – file may not exist in production

	// ========================
	// Fiber App Configuration
	// ========================
	app := fiber.New(fiber.Config{
		AppName:      "Project Name",
		ErrorHandler: utils.ErrorHandler,
	})

	// ========================
	// Middleware
	// ========================
	app.Use(recover.New())
	app.Use(logger.New())

	// ========================
	// Health Check
	// ========================
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "OK",
			"message":   "Service is running",
			"timestamp": utils.CurrentTimestamp(),
		})
	})

	// ========================
	// Register Module Routes
	// ========================
	api := app.Group("/api")
	sample.RegisterRoutes(api)

	// ========================
	// Port Configuration
	// ========================
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf(
		"Service Starting On Port %s",
		port,
	)

	// ========================
	// Start Server
	// ========================
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf(
			"Failed To Start Server: %v",
			err,
		)
	}
}
```

### Create Utils file
1. Create `utils/date_time.go` file
```sh
mkdir -p utils
touch utils/date_time.go
```
- Add code
```sh
package utils

import "time"

const DefaultTimeFormat = "2006-01-02 15:04:05"

func CurrentTimestamp() string {
    loc, _ := time.LoadLocation("Asia/Kuala_Lumpur")
    return time.Now().In(loc).Format(DefaultTimeFormat)
}

func CurrentUTCTime() time.Time {
    return time.Now().UTC()
}
```

2. Create `utils/error_handler.go` file
```sh
touch utils/error_handler.go
```
- Add code
```sh
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
```

## Folder Structure
### Create the folders:
```sh
mkdir -p module/{module1,module2,module3}
```

## Run the Server
```sh
go run main.go
```

## Create Dockerfile
1. Create docker file
```sh
touch Dockerfile
```
2. Add code
```sh
FROM golang:1.25.4-alpine AS builder

RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Kuala_Lumpur /etc/localtime && \
    echo "Asia/Kuala_Lumpur" > /etc/timezone

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o project .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/project .

EXPOSE 8080
ENTRYPOINT ["./project"]
```
3. Build and run
- build images for linux & windows
```sh
docker build -t 'project-name':latest .
```
- build images for macos
```sh
docker buildx build --platform linux/amd64 -t 'project-name':latest .
```
- run images
```sh
docker run -p 8080:8080 'project-name':latest
```

## Deploy to Kubernetes
1. Create a file: k8s-deployment.yaml
```sh
apiVersion: apps/v1
kind: Deployment
metadata:
name: 'project-name'
spec:
replicas: 1
selector:
    matchLabels:
    app: 'project-name'
template:
    metadata:
    labels:
        app: 'project-name'
    spec:
    containers:
        - name: 'project-name'
        image: 'project-name':latest
        ports:
            - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
name: 'project-name'
spec:
type: ClusterIP
selector:
    app: 'project-name'
ports:
    - port: 80
    targetPort: 8080
```

2. Deploy
```sh
kubectl apply -f k8s-deployment.yaml
```

## Setting Up Hot Reload with Air
### Install Air
```sh
# Install Air globally (recommended)
go install github.com/air-verse/air@latest

# Verify installation
air -v
```
- If air command not found, add GOPATH/bin to your PATH:
```sh
export PATH=$PATH:$(go env GOPATH)/bin
```
### Initialize Air Configuration
- Navigate to your project root and initialize Air:
```sh
cd /path/to/your/booking-service
air init
```
- This creates a .air.toml configuration file

### Run with hot reload
```sh
air
```

## Project Overview

A production-ready Go/Fiber starter template that provides the structural patterns, tooling, and infrastructure concerns already solved — so new API projects start with security, observability, and testability already in place.

**Tech Stack:**
- Go 1.25.4
- Fiber v2.52.12 (HTTP framework)
- godotenv v1.5.1 (env loading)
- go-playground/validator v10 (input validation)
- Air (hot reload for development)

**Key features built in:**
- Centralized error handling with sanitized responses (no internal detail leakage)
- Input validation on path and query parameters
- Graceful shutdown with 30-second timeout
- 1 MB request body limit
- Panic recovery middleware
- Structured module layout ready to clone and extend

## Module Structure

Each feature module lives under `module/` and follows a consistent four-file layout:

```
module/
└── sample/
    ├── routes.go      — Route registration (public via RegisterRoutes)
    ├── controller.go  — HTTP handlers (private functions)
    ├── service.go     — Business logic (private functions)
    └── form.go        — Request/Response structs with JSON and form tags
```

Shared infrastructure lives in:

```
utils/
├── date_time.go       — Timestamp formatting helpers
├── error_handler.go   — Centralized HTTP error response functions
└── validate.go        — Path and query parameter validators

config/
└── app_config.go      — Typed configuration loaded from environment variables
```

To add a new module, create a new directory under `module/`, add the four standard files, and call `RegisterRoutes` from `main.go`.

## Configuration

All configuration is loaded from environment variables via `config.LoadAppConfig()`. An optional `.env` file is supported for local development (loaded by `godotenv.Load()` at startup).

| Variable   | Required | Default      | Description                                    |
|------------|----------|--------------|------------------------------------------------|
| `PORT`     | No       | `8080`       | TCP port the server listens on                 |
| `APP_NAME` | No       | `go-project` | Application name reported in Fiber config      |
| `TIMEZONE` | No       | `UTC`        | IANA timezone name used for timestamp output   |

If `TIMEZONE` is set to an invalid location string, `UTC` is used silently.

## API Endpoints

The sample module is registered under `/api/sample` and demonstrates the handler patterns to copy when building real modules.

| Method | Path                        | Description                                      |
|--------|-----------------------------|--------------------------------------------------|
| `GET`  | `/`                         | Health check — returns service status and timestamp |
| `GET`  | `/livez`                    | Liveness probe (Fiber healthcheck middleware)    |
| `GET`  | `/readyz`                   | Readiness probe (Fiber healthcheck middleware)   |
| `GET`  | `/api/sample/hello`         | Returns a static greeting message                |
| `GET`  | `/api/sample/hello/:name`   | Returns a greeting for the given path parameter  |
| `GET`  | `/api/sample/hello-query`   | Returns a greeting using `?name=` query param    |
| `GET`  | `/api/sample/hello-service` | Returns a greeting generated by the service layer |
| `POST` | `/api/sample/hello-form`    | Accepts a JSON or form body with a `name` field  |

**Request body for `POST /api/sample/hello-form`:**
```json
{ "name": "Alice" }
```

**Response shape:**
```json
{ "message": "Hello, Alice! This message was generated from the form handler." }
```

## Testing

Tests use the Go standard `testing` package and Fiber's built-in `app.Test()` helper — no real server is started.

**Run all tests:**
```sh
go test ./...
```

**Run tests for a specific package:**
```sh
go test ./module/sample/...
go test ./test/config/...
go test ./test/utils/...
```

**Run with verbose output:**
```sh
go test -v ./...
```

Test files are co-located with their packages or placed under `test/` for cross-package concerns:

```
module/sample/controller_test.go   — HTTP handler tests (body parsing, param validation, error sanitization)
test/config/app_config_test.go     — Config loading tests
test/utils/date_time_test.go       — Timestamp helper tests
test/utils/validate_test.go        — Input validation tests
```

No coverage threshold is configured. Tests are run with `go test ./...`.
