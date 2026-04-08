<!-- generated-by: gsd-doc-writer -->
# Configuration

This project is configured exclusively through environment variables. There are no JSON/YAML config files — all settings are loaded once at startup via `config.LoadAppConfig()` in `config/app_config.go` and passed through the application as a typed `AppConfig` struct. Handlers never call `os.Getenv` directly.

## Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `APP_NAME` | Optional | `go-project` | Human-readable name for the application. Used in the Fiber app configuration and visible in logs. |
| `PORT` | Optional | `8080` | TCP port the HTTP server listens on. |
| `TIMEZONE` | Optional | `UTC` | IANA timezone string for timestamp formatting in API responses (e.g., `UTC`, `Asia/Kuala_Lumpur`, `America/New_York`). Invalid values fall back to `UTC` silently. |

All three variables are optional. The application starts successfully with no `.env` file present.

## Config File Format

There are no structured config files (JSON/YAML/TOML) beyond environment variables. The Fiber framework settings (body limit, read timeout, error handler) are hardcoded constants in `main.go` and are not environment-driven:

| Setting | Value | Location |
|---------|-------|----------|
| `BodyLimit` | `1 MB` | `main.go` line 40 |
| `ReadTimeout` | `30 seconds` | `main.go` line 41 |
| `ShutdownTimeout` | `30 seconds` | `main.go` line 87 |

## Required vs Optional Settings

No environment variable causes a startup failure if absent — all have safe defaults defined in `config/app_config.go`:

```go
// Defaults applied in LoadAppConfig()
appName = "go-project"   // if APP_NAME is empty
port    = "8080"         // if PORT is empty
tz      = "UTC"          // if TIMEZONE is empty or invalid
```

The application will always start. An invalid `TIMEZONE` value silently falls back to UTC without logging a warning.

## Defaults

| Variable | Default Value | Set In |
|----------|---------------|--------|
| `APP_NAME` | `go-project` | `config/app_config.go` line 23 |
| `PORT` | `8080` | `config/app_config.go` line 28 |
| `TIMEZONE` | `UTC` | `config/app_config.go` line 33 |

The `Dockerfile` also sets `ENV TZ=Asia/Kuala_Lumpur` at the OS level (line 26), which affects the container's system timezone independently of the `TIMEZONE` variable used for API response timestamps.

## Per-Environment Overrides

For local development, copy `.env.example` to `.env` and set values there:

```bash
cp .env.example .env
```

The `.env` file is loaded by `godotenv` at startup (`main.go` line 26) and is gitignored — never commit real values to source control.

For production, set environment variables through your hosting provider or container orchestration system. The `.env` file is not used in production containers; the `Dockerfile` does not copy it into the image.

**Docker example:**

```bash
docker run -e PORT=9000 -e APP_NAME=my-service -e TIMEZONE=Asia/Kuala_Lumpur -p 9000:9000 my-image
```

<!-- VERIFY: Kubernetes/Helm secret management or any container orchestration platform-specific override mechanism -->

## Hot Reload Configuration (Development Only)

Air's hot reload behaviour is controlled by `.air.toml` in the project root. Key settings:

| Setting | Value | Description |
|---------|-------|-------------|
| `build.cmd` | `go build -o ./tmp/main .` | Build command run on file change |
| `build.bin` | `./tmp/main` | Output binary path |
| `build.include_ext` | `go, tpl, tmpl, html` | File extensions that trigger a rebuild |
| `build.exclude_regex` | `_test.go` | Test files are excluded from watch |
| `build.delay` | `1000ms` | Debounce delay before rebuild starts |

Air does not load the `.env` file itself (`env_files = []`); `godotenv` handles that at application startup.
