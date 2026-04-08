# Stack Research: Go/Fiber API Baseline

**Date:** 2026-04-08
**Milestone:** Subsequent — upgrading existing Go/Fiber skeleton to production baseline

---

## Structured Logging: Zerolog (RECOMMENDED)

**`github.com/rs/zerolog` v1.32.0+**

- Lower memory allocation than Zap, simpler API than slog, good Fiber middleware integration
- **Confidence:** HIGH

**Do NOT use:**
- Zap — heavier, more complex than needed
- slog alone — still emerging, less Fiber integration

```go
if os.Getenv("ENV") == "production" {
    log.Logger = zerolog.New(os.Stderr)
}
```

---

## Database: GORM v1.25 + Postgres

- `gorm.io/gorm` v1.25.x
- `gorm.io/driver/postgres` v1.5.4
- `github.com/lib/pq` v1.10.9

**Why GORM:** Standard Go ORM, Postgres first-class, built-in migrations, connection pooling via `sql.DB`.

**Do NOT use:** sqlc, Ent, sqlx — overkill or too low-level for a baseline template.

```go
func InitDB(dsn string) (*gorm.DB, error) {
    return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
```

---

## Security Middleware: Fiber Native + Helmet

- Fiber built-in: `cors`, `limiter`, `recover` (already in v2.52.12)
- `github.com/gofiber/helmet/v2` v2.2.22+

```go
app.Use(cors.New(cors.Config{
    AllowOrigins: os.Getenv("ALLOWED_ORIGINS"),
    AllowMethods: "GET,POST,PUT,DELETE",
}))
app.Use(limiter.New(limiter.Config{
    Max:        100,
    Expiration: 1 * time.Minute,
}))
app.Use(helmet.New())
```

---

## Configuration: Viper + godotenv (keep existing)

- Keep: `github.com/joho/godotenv` v1.5.1 (already in go.mod)
- Add: `github.com/spf13/viper` v1.17.0+

godotenv loads `.env` files → Viper reads env vars, sets defaults, validates.

**Do NOT use:** `os.Getenv()` directly scattered throughout code.

```go
func LoadAppConfig() (*Config, error) {
    _ = godotenv.Load()
    viper.SetDefault("PORT", 8080)
    viper.SetDefault("ENV", "development")
    viper.AutomaticEnv()
    return &Config{
        Port:        viper.GetInt("PORT"),
        DatabaseURL: viper.GetString("DATABASE_URL"),
    }, nil
}
```

---

## Testing: httptest + Testify

- Go built-in: `net/http/httptest`
- Add: `github.com/stretchr/testify` v1.8.4+

**Do NOT use:** Ginkgo/Gomega — overkill for baseline.

```go
func TestHandler(t *testing.T) {
    app := fiber.New()
    app.Get("/users/:id", GetUser)
    req := httptest.NewRequest("GET", "/users/1", nil)
    resp, _ := app.Test(req)
    assert.Equal(t, http.StatusOK, resp.StatusCode)
}
```

---

## Installation

```bash
go get github.com/rs/zerolog@v1.32.0
go get github.com/spf13/viper@v1.17.0
go get gorm.io/gorm@v1.25.0
go get gorm.io/driver/postgres@v1.5.4
go get github.com/lib/pq@v1.10.9
go get github.com/gofiber/helmet/v2@v2.2.22
go get -t github.com/stretchr/testify@v1.8.4
go mod tidy
```

---

## Decision Matrix

| Area | Chosen | Confidence |
|------|--------|-----------|
| Logging | Zerolog | HIGH |
| ORM | GORM + Postgres | HIGH |
| Security | Fiber native + Helmet | HIGH |
| Config | Viper + godotenv | MEDIUM-HIGH |
| Testing | httptest + Testify | HIGH |
