# ROADMAP: Go Based Project

**Project:** Production-ready Go/Fiber starter template
**Created:** 2026-04-08
**Granularity:** Standard
**Total Phases:** 5
**Coverage:** 21/21 v1 requirements mapped

---

## Phases

- [ ] **Phase 1: Foundation & Configuration** - Core config, graceful shutdown, input protection, error sanitization, environment setup
- [ ] **Phase 2: Observability & Middleware** - Structured logging, request correlation, CORS, UTC timestamps  
- [ ] **Phase 3: Database Layer** - GORM+Postgres integration, migrations, CRUD sample, repository pattern
- [ ] **Phase 4: Security Hardening** - Rate limiting, security headers, and comprehensive error responses
- [ ] **Phase 5: Testing & Documentation** - Unit tests, mocks, test helpers, README updates

---

## Phase Details

### Phase 1: Foundation & Configuration

**Goal:** All application configuration is externalized, startup is robust, and user inputs are validated and limited.

**Depends on:** None (foundation phase)

**Requirements:** FOUND-01, FOUND-02, FOUND-03, FOUND-04, FOUND-05, SEC-03, DOC-01

**Success Criteria** (what must be TRUE):
1. Application reads APP_NAME, TIMEZONE, PORT from environment variables with no hardcoded values
2. Application gracefully shuts down on SIGTERM, draining in-flight requests with 30-second timeout
3. All HTTP request bodies are limited to 1MB; oversized requests return 400 Bad Request
4. All path parameters and query strings are validated for length and format before use
5. Body parser failures return 400 Bad Request with structured error response (not silent 200)
6. All errors return sanitized responses (no stack traces); only generic message + request_id exposed
7. `.env.example` file documents all required environment variables with descriptions

**Plans:** 3 plans

Plans:
- [x] 01-01-PLAN.md — Externalize config (APP_NAME, TIMEZONE, PORT), update date_time.go, implement graceful SIGTERM shutdown in main.go
- [x] 01-02-PLAN.md — Install validator, create input validation helpers, fix body parser 400 errors, add request_id to InternalServerErrorHandler
- [x] 01-03-PLAN.md — Create .env.example, run full test suite phase gate, human verification checkpoint

---

### Phase 2: Observability & Middleware

**Goal:** All events are logged structurally, requests are correlatable end-to-end, and responses include consistent metadata.

**Depends on:** Phase 1

**Requirements:** OBS-01, OBS-02, OBS-03, OBS-04

**Success Criteria** (what must be TRUE):
1. All logs are output as JSON in production, human-readable in development via Zap configuration
2. Every request is assigned a unique correlation ID, logged at entry, and returned in X-Request-ID response header
3. CORS requests are evaluated against env-driven allow list (origins, methods, headers); blocked requests return 403
4. All timestamps in responses are RFC3339 format in UTC (e.g., 2026-04-08T15:30:45Z)

**Plans:** TBD

**UI hint:** yes

---

### Phase 3: Database Layer

**Goal:** Applications have a persistent database foundation with an established ORM, migration tooling, and a clear repository pattern to follow.

**Depends on:** Phase 1

**Requirements:** DB-01, DB-02, DB-03, DB-04

**Success Criteria** (what must be TRUE):
1. GORM+Postgres connection is established with pooling configured (MaxOpenConns=25, MaxIdleConns=5, MaxLifetime=5m)
2. golang-migrate is scaffolded in `migrations/` with example up/down SQL files for creating a Users table
3. Sample module (User) demonstrates full CRUD operations via GORM with create, read, update, delete handlers working end-to-end
4. Service layer calls repository layer, which calls GORM; this three-tier pattern is enforced in sample module and documented

**Plans:** TBD

**UI hint:** yes

---

### Phase 4: Security Hardening

**Goal:** All applications are protected against common attack vectors via middleware that enforces rate limits and critical security headers.

**Depends on:** Phase 1, Phase 2

**Requirements:** SEC-01, SEC-02

**Success Criteria** (what must be TRUE):
1. Rate limiter middleware returns 429 Too Many Requests for IPs exceeding 100 req/min; includes X-Retry-After header
2. gofiber/helmet middleware is integrated, returning HSTS, CSP, X-Frame-Options, X-Content-Type-Options headers on all responses
3. All security middleware responses include request_id for debugging

**Plans:** TBD

---

### Phase 5: Testing & Documentation

**Goal:** The template provides working examples for testing all layers and clear guidance for teams to understand and extend the starter.

**Depends on:** Phase 3, Phase 4

**Requirements:** TEST-01, TEST-02, TEST-03, DOC-02

**Success Criteria** (what must be TRUE):
1. Handler tests for sample module routes exist using httptest + Testify assertions; tests run via `go test ./...`
2. Service layer tests mock the repository; tests demonstrate the mock pattern with example assertions
3. `.env.test` example and test helper exist, allowing tests to spin up the Fiber app with test configuration in-memory
4. README includes "How to add a new module" section with step-by-step instructions for controller/service/repository/routes pattern
5. README documents setup (Go + Postgres install), how to run locally, how to run tests, how to run migrations

**Plans:** TBD

**UI hint:** yes

---

## Progress Table

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Foundation & Configuration | 0/3 | Planned | - |
| 2. Observability & Middleware | 0/? | Not started | - |
| 3. Database Layer | 0/? | Not started | - |
| 4. Security Hardening | 0/? | Not started | - |
| 5. Testing & Documentation | 0/? | Not started | - |

---

## Coverage Summary

**Requirements Mapped:** 21/21 ✓

| Phase | Count | Requirements |
|-------|-------|--------------|
| 1 | 7 | FOUND-01, FOUND-02, FOUND-03, FOUND-04, FOUND-05, SEC-03, DOC-01 |
| 2 | 4 | OBS-01, OBS-02, OBS-03, OBS-04 |
| 3 | 4 | DB-01, DB-02, DB-03, DB-04 |
| 4 | 2 | SEC-01, SEC-02 |
| 5 | 4 | TEST-01, TEST-02, TEST-03, DOC-02 |

**No orphaned requirements.** All v1 requirements mapped to exactly one phase.

---

*Roadmap created: 2026-04-08*
*Phase 1 planned: 2026-04-08 — 3 plans, 2 waves*
