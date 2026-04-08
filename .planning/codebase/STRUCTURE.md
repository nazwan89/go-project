# Structure

## Directory Layout

```
go-project/
├── main.go                     # Entry point — Fiber app setup, middleware, module registration
├── go.mod                      # Module definition: based-project/go-project
├── go.sum                      # Dependency lock file
├── Dockerfile                  # Container build (with Air hot reload)
├── README.md                   # Project overview
│
├── module/                     # Feature modules (one dir per domain)
│   └── sample/                 # Example module
│       ├── controller.go       # HTTP handlers (private)
│       ├── service.go          # Business logic (private)
│       ├── routes.go           # Route registration (public RegisterRoutes)
│       └── form.go             # Request/Response structs (exported)
│
├── utils/                      # Shared utilities
│   ├── error_handler.go        # Centralized error handling
│   └── date_time.go            # Timestamp helpers
│
└── tmp/                        # Air hot-reload build output (gitignored)
    └── main
```

## Key Locations

| Location | Purpose |
|----------|---------|
| `main.go` | App bootstrap, middleware wiring, module registration |
| `module/<name>/routes.go` | Call `RegisterRoutes(app)` to mount a feature |
| `module/<name>/controller.go` | HTTP request handlers |
| `module/<name>/service.go` | Business logic separated from transport |
| `utils/error_handler.go` | All error response formatting |

## Naming Conventions

- Module directories: lowercase, singular (`sample`, not `samples`)
- Go files within modules: role-based (`controller.go`, `service.go`, `routes.go`, `form.go`)
- Utility files: snake_case describing function (`error_handler.go`, `date_time.go`)

## Adding a New Module

1. Create `module/<name>/` directory
2. Add `form.go` — request/response types
3. Add `service.go` — business logic
4. Add `controller.go` — HTTP handlers calling service
5. Add `routes.go` — `RegisterRoutes(app *fiber.App)` function
6. Call `<name>.RegisterRoutes(app)` in `main.go`
