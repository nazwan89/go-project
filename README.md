# go-project
Based project using golang fiber

### Install Go 
1. Windows
	- Download from https://go.dev/dl/
	- Run the installer.
	- Verify installation using bash: "go version"

2. MacOS
```sh
    brew install go
    go version
```

3. Linux
```sh
    sudo apt update
    sudo apt install golang-go -y
    go version
```

### Create Project
1. Create Project Folder:
```sh
    mkdir 'project name'
    cd 'project name'
```

2. Initialize a Go Module:
```sh
    go mod init 'project name'
```

3. Install Fiber:
```sh
    go get github.com/gofiber/fiber/v2
```

4. Create an `.env` file (optional: use for local development). Go will not read this file automatically; you either need to load it in code with a package such as `github.com/joho/godotenv` or export the variable before starting the app.

```sh
    touch .env
```
- .env
```sh
    PORT=8080
```

To start the server with the value from `.env` you can run:

```sh
    go run main.go
```

or manually export before running:

```sh
    export PORT=8080
    go run main.go
```

(see the `main.go` example which uses `godotenv.Load()`)

5. Create main.go file
```sh
    touch main.go
```
Add Code Files
- main.go
```sh
    package main

    import (
        "log"
        "os"

        "github.com/joho/godotenv"
        "github.com/gofiber/fiber/v2"
        "github.com/gofiber/fiber/v2/middleware/logger"
        "github.com/gofiber/fiber/v2/middleware/recover"
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
            AppName: "Project Name",
            // ErrorHandler: handlers.ErrorHandler,
        })

		// ========================
		// Middleware
		// ========================
		app.Use(recover.New())

		app.Use(logger.New(
			logger.Config{
				Format: "[${time}] ${status} - ${method} ${path}\n",
			},
		))

		// ========================
		// Health Check
		// ========================
		app.Get("/", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"status":  "OK",
				"message": "Service is running",
			})
		})

		// ========================
		// Port Configuration
		// ========================
		port := os.Getenv("PORT")
		if port == "" {port = "8080"}

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

6. Create date_time utils
```sh
    mkdir -p utils
    touch utils/date_time.go
```
Add code
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

7. Create error utils
```sh
    touch utils/error_handler.go
```
Add code
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

### Folder Structure
1. Create the folders:
```sh
    mkdir -p module/{module1,module2,module3}
```