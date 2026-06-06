# Architecture Research: Go/Fiber API Baseline

**Date:** 2026-04-08
**Confidence:** HIGH

---

## System Architecture

```
Client HTTP Request
    ↓
Fiber Router + Middleware Stack
  (recovery → request-id → logging → rate-limit → CORS → security → size-limit)
    ↓
Module Layer (controller — extracts & validates request)
    ↓
Service Layer (business logic)
    ↓
Repository Layer (NEW — data access)
    ↓
GORM Models
    ↓
PostgreSQL
```

---

## Critical Components

### 1. Configuration Layer (NEW) — `config/app_config.go`

Single `Config` struct, not package-level globals. Explicit dependencies, testable, validated at startup.

```go
type Config struct {
    AppName       string
    Port          string
    DatabaseURL   string
    LogLevel      string
    RateLimitReqs int
}

func Load() *Config {
    _ = godotenv.Load()
    return &Config{
        AppName:     getEnv("APP_NAME", "go-project"),
        Port:        getEnv("PORT", "8080"),
        DatabaseURL: os.Getenv("DATABASE_URL"),
    }
}
```

### 2. Dependency Injection — `di/container.go`

Manual DI (NOT wire, NOT fx). Explicit, zero-overhead, no code generation.

```go
type Container struct {
    DB     *gorm.DB
    Config *config.Config
    Logger *zap.Logger
}

// RegisterRoutes signature evolves to accept container:
func RegisterRoutes(g fiber.Router, c *di.Container) { ... }
```

### 3. Module Evolution — Add `model.go` + `repository.go`

Keep existing: `form.go`, `handlers.go`, `service.go`, `routes.go`

Add per module:
- `model.go` — GORM model definitions
- `repository.go` — DB access (service calls repo, repo calls GORM)

This keeps service layer DB-agnostic and easy to mock/test.

### 4. Database Layer (NEW) — `database/postgres.go`

```go
func Connect(databaseURL string) (*gorm.DB, error) {
    db, _ := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
    sqlDB, _ := db.DB()
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetConnMaxLifetime(time.Hour)
    return db, nil
}
```

**Migrations:** Use `golang-migrate` with SQL files under `migrations/` — NOT AutoMigrate in production. AutoMigrate only for dev.

### 5. Middleware Ordering (CRITICAL)

```go
app.Use(recover.New())           // 1. MUST be first — catches panics
app.Use(middleware.RequestID())  // 2. Correlation ID — all logs include it
app.Use(logger.New())            // 3. Request logging
app.Use(middleware.RateLimit())  // 4. Reject bad IPs early
app.Use(cors.New(...))           // 5. Reject bad origins
app.Use(helmet.New())            // 6. Security headers
app.Use(middleware.SizeLimit())  // 7. Before parsing large bodies
```

### 6. Structured Logging — `logging/logger.go`

Use **Zap** (high performance, structured JSON, correlation ID propagation).

```go
logger.Info("user created",
    zap.String("request_id", reqID),
    zap.Uint("user_id", userID),
    zap.Int64("duration_ms", ms),
)
```

### 7. Graceful Shutdown — `main.go`

```go
go func() {
    sigChan := make(chan os.Signal, 1) // buffered — critical
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    <-sigChan
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    app.ShutdownWithContext(ctx) // drain in-flight requests
}()
```

---

## File Structure (Post Phase 2+3)

```
go-project/
├── config/
│   └── app_config.go      # NEW — AppConfig struct + LoadConfig()
├── di/
│   └── container.go       # NEW — DI container
├── logging/
│   └── logger.go          # NEW — Zap setup
├── database/
│   ├── postgres.go        # NEW — GORM connection + pool
│   └── migration.go       # NEW — migration runner
├── middleware/
│   ├── request_id.go      # NEW
│   ├── security.go        # NEW — security headers
│   ├── rate_limit.go      # NEW
│   └── cors.go            # NEW
├── module/
│   └── sample/
│       ├── model.go       # NEW — GORM model
│       ├── repository.go  # NEW — DB access
│       ├── service.go     # Enhanced — repo injection
│       ├── handlers.go    # Enhanced — service injection
│       ├── routes.go      # Enhanced — accept container
│       └── form.go        # Existing
├── migrations/            # NEW
│   ├── 001_initial.up.sql
│   └── 001_initial.down.sql
├── utils/                 # Existing
├── main.go
└── .env.example           # NEW
```

---

## Key Decisions

| Decision | Why | Alternative Not Chosen |
|----------|-----|----------------------|
| Manual DI | Explicit, zero overhead, easy to understand | wire/fx — overkill for this size |
| Single Config struct | Testable, all config visible | Package globals — hard to test |
| golang-migrate | Version control + rollback | AutoMigrate — unsafe for production |
| Zap logging | Structured, high-perf, correlation IDs | zerolog — also good, Zap more widely used |
| Graceful shutdown | Drain requests, no data loss | Hard shutdown — client errors |

---

## Anti-Patterns to Avoid

- Global database variable — pass via DI
- Package-level logger — inject via DI
- Wrong middleware order — recovery MUST be first
- AutoMigrate in production — use golang-migrate
- No connection pooling — always set MaxOpenConns/MaxIdleConns
- Hardcoded config — load from environment always
