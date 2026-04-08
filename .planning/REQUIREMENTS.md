# Requirements: Go Based Project

**Defined:** 2026-04-08
**Core Value:** Every new API project should start with security, observability, and testability already in place — not bolted on later.

---

## v1 Requirements

### Foundation

- [ ] **FOUND-01**: App config loaded from environment variables (APP_NAME, TIMEZONE, PORT) — no hardcoded values
- [ ] **FOUND-02**: Graceful shutdown on SIGTERM — drain in-flight requests before exit (30s timeout)
- [ ] **FOUND-03**: Request body size limit (1MB) — prevent memory exhaustion DoS
- [ ] **FOUND-04**: Input validation helpers for path params and query strings (length, format)
- [ ] **FOUND-05**: Body parser always returns 400 on parse failure (not silent 200)

### Observability

- [ ] **OBS-01**: Structured logging with Zap — JSON output in production, readable in development
- [ ] **OBS-02**: Request ID middleware — correlation ID injected into every request, logged and returned in response headers
- [ ] **OBS-03**: CORS middleware — env-driven allow list (origins, methods, headers)
- [ ] **OBS-04**: Timestamps normalized to UTC RFC3339 format across all responses

### Database

- [ ] **DB-01**: GORM + Postgres connection with connection pooling (MaxOpenConns, MaxIdleConns, MaxLifetime)
- [ ] **DB-02**: Migration scaffolding using golang-migrate with example up/down SQL files under `migrations/`
- [ ] **DB-03**: Sample module updated to demonstrate CRUD operations (User model with GORM)
- [ ] **DB-04**: Repository layer pattern established in sample module (controller → service → repository → GORM)

### Security

- [ ] **SEC-01**: Rate limiting middleware — 100 req/min per IP, returns 429 with error body
- [ ] **SEC-02**: Security headers via gofiber/helmet (HSTS, CSP, X-Frame-Options, X-Content-Type-Options)
- [ ] **SEC-03**: Error response sanitization — internal errors return generic message + request_id only (no stack traces)

### Testing

- [ ] **TEST-01**: Unit test setup with httptest + Testify — example tests for sample module handlers
- [ ] **TEST-02**: Service layer unit tests with mock repository — example mock pattern
- [ ] **TEST-03**: `.env.test` example and test helper for spinning up Fiber app in tests

### Documentation

- [ ] **DOC-01**: `.env.example` file with all required environment variables and descriptions
- [ ] **DOC-02**: Updated README with setup, run, test, and "how to add a new module" instructions

---

## v2 Requirements

### Advanced Observability

- **ADV-01**: Prometheus metrics endpoint (`/metrics`)
- **ADV-02**: Sentry / error tracking integration
- **ADV-03**: OpenTelemetry tracing

### Auth Layer (Project-Specific Pattern)

- **AUTH-01**: Example JWT middleware (template only — teams swap their own secret/claims)

### Developer Experience

- **DX-01**: Swagger/OpenAPI generation via swag CLI (optional, documented)
- **DX-02**: Makefile with common commands (run, test, migrate, build)

---

## Out of Scope

| Feature | Reason |
|---------|--------|
| Authentication (JWT, OAuth) | Project-specific — different secret, claims, and provider per project |
| Authorization / RBAC | Project-specific permission models |
| Caching (Redis) | Optional — not all projects need it; adds operational overhead |
| API versioning | Teams choose their own pattern |
| gRPC / WebSocket | Out of scope for v1 REST baseline |
| Message queues | Optional, async patterns vary per project |
| CI/CD pipelines | Team-specific; not part of template |
| Kubernetes configs | Team-specific deployment configs |
| Multi-instance rate limiting (Redis-backed) | Projects that need this will swap the limiter; baseline uses in-memory |

---

## Traceability

| Requirement | Phase | Status |
|-------------|-------|--------|
| FOUND-01 | Phase 1 | Pending |
| FOUND-02 | Phase 1 | Pending |
| FOUND-03 | Phase 1 | Pending |
| FOUND-04 | Phase 1 | Pending |
| FOUND-05 | Phase 1 | Pending |
| OBS-01 | Phase 2 | Pending |
| OBS-02 | Phase 2 | Pending |
| OBS-03 | Phase 2 | Pending |
| OBS-04 | Phase 2 | Pending |
| DB-01 | Phase 3 | Pending |
| DB-02 | Phase 3 | Pending |
| DB-03 | Phase 3 | Pending |
| DB-04 | Phase 3 | Pending |
| SEC-01 | Phase 4 | Pending |
| SEC-02 | Phase 4 | Pending |
| SEC-03 | Phase 1 | Pending |
| TEST-01 | Phase 5 | Pending |
| TEST-02 | Phase 5 | Pending |
| TEST-03 | Phase 5 | Pending |
| DOC-01 | Phase 1 | Pending |
| DOC-02 | Phase 5 | Pending |

**Coverage:**
- v1 requirements: 21 total
- Mapped to phases: 21
- Unmapped: 0 ✓

---
*Requirements defined: 2026-04-08*
*Last updated: 2026-04-08 after initial definition*
