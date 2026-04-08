# Go Based Project

## What This Is

A production-ready Go/Fiber starter template that teams can clone as the baseline for any new API project. It provides the structural patterns, tooling, and infrastructure concerns already solved — so new projects start with a solid foundation rather than from scratch.

## Core Value

Every new API project should start with security, observability, and testability already in place — not bolted on later.

## Requirements

### Validated

<!-- Shipped and confirmed valuable. -->

- ✓ Modular route structure (`module/<name>/controller+service+routes+form`) — established pattern
- ✓ Centralized error handling (`utils/error_handler.go`)
- ✓ Hot reload with Air
- ✓ Docker containerization
- ✓ Health check endpoint

### Active

<!-- Current scope. Building toward these. -->

**Security**
- [ ] Rate limiting middleware
- [ ] CORS configuration
- [ ] Request size limits
- [ ] Input sanitization / validation layer

**Observability**
- [ ] Structured logging (replace standard `log` with zerolog or zap)
- [ ] Request ID / correlation ID middleware

**Configuration**
- [ ] All config loaded from `.env` / environment variables (port, DB URL, timezone, app name)
- [ ] Graceful shutdown on SIGTERM

**Database**
- [ ] GORM + Postgres connection setup
- [ ] Database migration scaffolding
- [ ] Sample module updated to demonstrate DB usage

**Testing**
- [ ] Unit test setup for utility functions
- [ ] Handler test setup using `net/http/httptest`
- [ ] Sample tests for the sample module

### Out of Scope

- Authentication / JWT — project-specific, not baseline
- Redis / caching — project-specific
- CI/CD pipelines — team-specific, not part of template
- gRPC / WebSocket — out of scope for v1 REST baseline

## Context

- Existing codebase: Go 1.25.4, Fiber v2.52.12, godotenv for env loading
- Pattern established: `module/<name>/` with controller/service/routes/form
- Current gaps: no tests, no DB layer, hardcoded config, no security middleware, standard `log` only
- Codebase map: `.planning/codebase/` (STACK, ARCHITECTURE, CONVENTIONS, CONCERNS, etc.)
- Team: MRSB / milradius internal projects

## Constraints

- **Tech Stack**: Go + Fiber — do not introduce alternative frameworks
- **Simplicity**: Baseline must remain easy to understand; avoid over-engineering abstractions
- **Clone-ready**: No project-specific code in the final template; all examples use generic `sample` module

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| GORM + Postgres | Standard ORM for most internal projects | — Pending |
| Zerolog or Zap for logging | Structured, high-performance logging | — Pending |
| Fiber built-in middleware for security | Consistent with existing Fiber usage | — Pending |

---
*Last updated: 2026-04-08 after initial project definition*
