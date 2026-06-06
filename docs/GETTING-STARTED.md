<!-- generated-by: gsd-doc-writer -->
# Getting Started

This guide takes you from a fresh clone to a running local server and a passing test suite in a few minutes.

## Prerequisites

Ensure the following tools are installed before proceeding.

| Tool | Required Version | Notes |
|------|-----------------|-------|
| Go | `>= 1.25.4` | [https://go.dev/dl/](https://go.dev/dl/) |
| Git | Any recent version | For cloning the repository |
| Air (optional) | Latest | Live reload during development — `go install github.com/air-verse/air@latest` |
| Docker (optional) | Any recent version | For container-based local runs |

Verify your Go installation:

```bash
go version
```

Expected output: `go version go1.25.4 ...` (or higher).

## Installation Steps

**1. Clone the repository**

```bash
git clone <repository-url>
cd go-project
```

**2. Download Go module dependencies**

```bash
go mod download
```

This fetches all dependencies declared in `go.mod` into the local module cache. No further package manager commands are needed.

**3. Configure environment variables**

```bash
cp .env.example .env
```

The defaults in `.env.example` work for local development without any edits. See [CONFIGURATION.md](CONFIGURATION.md) for a full variable reference.

## Running the Server Locally

### Option A — standard `go run` (no hot reload)

```bash
go run ./cmd
```

Expected output:

```
Service Starting On Port 8080
```

The server is now listening at `http://localhost:8080`.

### Option B — Air (hot reload)

If Air is installed, run from the project root:

```bash
air
```

Air watches `.go`, `.tpl`, `.tmpl`, and `.html` files. Any saved change triggers an automatic rebuild and restart. Build output and errors are written to `tmp/build-errors.log`.

### Option C — Docker

```bash
docker build -t go-project .
docker run -p 8080:8080 go-project
```

The multi-stage `Dockerfile` compiles the binary on Alpine Linux and produces a minimal image. The container exposes port `8080`.

## First API Call

With the server running, verify it responds correctly.

**Health check (root route):**

```bash
curl http://localhost:8080/
```

Expected JSON response:

```json
{
  "message": "Service is running",
  "status": "OK",
  "timestamp": "..."
}
```

**Sample module — static greeting:**

```bash
curl http://localhost:8080/api/sample/hello
```

**Sample module — greeting with path parameter:**

```bash
curl http://localhost:8080/api/sample/hello/World
```

**Sample module — greeting with query parameter:**

```bash
curl "http://localhost:8080/api/sample/hello-query?name=World"
```

**Sample module — POST with form/JSON body:**

```bash
curl -X POST http://localhost:8080/api/sample/hello-form \
  -H "Content-Type: application/json" \
  -d '{"name": "World"}'
```

## Running Tests

Run the full test suite from the project root:

```bash
go test ./...
```

Run a specific package:

```bash
go test ./module/sample/...
go test ./test/config/...
go test ./test/utils/...
```

Run tests with verbose output:

```bash
go test -v ./...
```

No test database, external services, or additional setup is required — all tests use in-process Fiber `httptest` instances.

## Common Setup Issues

**Wrong Go version**
If `go run ./cmd` reports a module directive error, your installed Go version is below `1.25.4`. Upgrade at [https://go.dev/dl/](https://go.dev/dl/).

**Port 8080 already in use**
Set a different port before starting the server:

```bash
PORT=9090 go run ./cmd
```

Or update `PORT=9090` in your `.env` file.

**Air not found**
If `air` is not on your `PATH` after installing, ensure `$(go env GOPATH)/bin` is included in your shell's `PATH`:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

**Module download failures**
If `go mod download` fails due to network restrictions, set the Go module proxy:

```bash
GOPROXY=https://proxy.golang.org,direct go mod download
```

## Next Steps

- [ARCHITECTURE.md](ARCHITECTURE.md) — how the codebase is structured and why
- [CONFIGURATION.md](CONFIGURATION.md) — full environment variable reference
- [DEVELOPMENT.md](DEVELOPMENT.md) — build commands, code style, and PR process
