<!-- generated-by: gsd-doc-writer -->
# Deployment

This project ships as a single statically compiled Go binary inside an Alpine Linux container. The only deployment artifact is the Docker image produced by the multi-stage `Dockerfile` at the project root. There is no CI/CD pipeline configured in the repository.

## Deployment Targets

| Target | Config File | Notes |
|--------|-------------|-------|
| Docker (self-hosted) | `Dockerfile` | Multi-stage Alpine build; primary deployment method |
| <!-- VERIFY: container registry (e.g., Docker Hub, GHCR, ECR) --> | — | No registry config detected in repo |
| <!-- VERIFY: container orchestration platform (e.g., Kubernetes, ECS, Fly.io) --> | — | No orchestration config detected in repo |

## Build Pipeline

No automated CI/CD workflow files are present in the repository. The build and deployment steps are performed manually.

### 1. Build the Docker image

```bash
docker build -t go-project:latest .
```

The `Dockerfile` uses a two-stage build:

1. **Builder stage** (`golang:1.25.4-alpine`) — installs timezone data, downloads Go modules, compiles the binary with `CGO_ENABLED=0 GOOS=linux go build -o project .`
2. **Runtime stage** (`alpine:latest`) — copies only the compiled `project` binary into a minimal Alpine image (~6–7 MB total). Sets timezone to `Asia/Kuala_Lumpur` at the OS level via `ENV TZ=Asia/Kuala_Lumpur`.

The final image exposes port `8080` and uses `ENTRYPOINT ["./project"]`.

### 2. Tag and push the image

<!-- VERIFY: target container registry URL and authentication method -->

```bash
docker tag go-project:latest <registry>/<org>/go-project:<version>
docker push <registry>/<org>/go-project:<version>
```

### 3. Run the container

```bash
docker run -d \
  --name go-project \
  -p 8080:8080 \
  -e APP_NAME=go-project \
  -e PORT=8080 \
  -e TIMEZONE=Asia/Kuala_Lumpur \
  go-project:latest
```

## Environment Setup

All configuration is supplied via environment variables at container runtime. The `.env` file is **not** included in the Docker image and is not used in production.

Set the following variables through your hosting provider, container orchestration secrets manager, or the `docker run -e` flag:

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `APP_NAME` | Optional | `go-project` | Human-readable application name shown in logs |
| `PORT` | Optional | `8080` | TCP port the server listens on; must match the container's exposed port |
| `TIMEZONE` | Optional | `UTC` | IANA timezone for API response timestamps (e.g., `Asia/Kuala_Lumpur`) |

All three variables are optional — the application starts cleanly with none set. See [CONFIGURATION.md](CONFIGURATION.md) for full details on defaults and config loading.

**Docker Compose example (development/staging):**

```yaml
services:
  api:
    image: go-project:latest
    ports:
      - "8080:8080"
    environment:
      - APP_NAME=go-project
      - PORT=8080
      - TIMEZONE=Asia/Kuala_Lumpur
```

<!-- VERIFY: production secret management (e.g., Kubernetes Secrets, AWS Secrets Manager, HashiCorp Vault) -->

## Rollback Procedure

No automated rollback is configured. To revert a deployment, redeploy the previous image tag:

```bash
docker stop go-project
docker rm go-project
docker run -d \
  --name go-project \
  -p 8080:8080 \
  -e APP_NAME=go-project \
  -e PORT=8080 \
  -e TIMEZONE=Asia/Kuala_Lumpur \
  go-project:<previous-version>
```

<!-- VERIFY: platform-specific rollback command if deployed to a managed container service -->

## Monitoring

No monitoring library is configured in `go.mod` or `main.go`. The application provides a built-in health check endpoint via Fiber's `healthcheck` middleware:

| Endpoint | Method | Response |
|----------|--------|----------|
| `/livez` | `GET` | `200 OK` when the process is running |
| `/readyz` | `GET` | `200 OK` when the process is ready to serve |
| `/` | `GET` | `200 OK` with `{"status":"OK","message":"Service is running","timestamp":"..."}` |

HTTP request logs are emitted to stdout via Fiber's `logger` middleware and are available through your container runtime's log collection mechanism (e.g., `docker logs go-project`).

<!-- VERIFY: external monitoring service (e.g., Sentry, Datadog, New Relic) if configured outside the repository -->
<!-- VERIFY: log aggregation platform (e.g., Loki, CloudWatch, Datadog Logs) -->
