# Feature Landscape

**Domain:** Production-ready Go/Fiber REST API baseline template
**Researched:** 2026-04-08

---

## Table Stakes

| Feature | Status | Complexity |
|---------|--------|-----------|
| Error Handling (centralized, consistent) | ✓ Shipped | Low |
| Health Check Endpoint | ✓ Shipped | Low |
| Panic Recovery | ✓ Shipped | Low |
| Modular Route Organization | ✓ Shipped | Low |
| Request/Response Logging | ✓ Shipped (Fiber logger) | Low |
| Environment Configuration (no hardcoding) | Partial — port works, app name/timezone hardcoded | Low |
| Input Validation Helpers | Partial — binding exists, no validation helpers | Low |
| Graceful Shutdown (SIGTERM, drain requests) | Missing | Low |
| Database Connection Setup (GORM + Postgres) | Missing | Medium |
| Database Migrations | Missing | Medium |
| HTTP Status Code Consistency (201, 404, etc.) | Partial — error handlers exist, no POST/201 examples | Low |

---

## Differentiators

| Feature | Value | Complexity |
|---------|-------|-----------|
| Structured Logging (Zerolog) | Machine-readable logs, log aggregation | Medium |
| Request ID / Correlation ID Middleware | Trace requests across logs | Medium |
| Rate Limiting Middleware | Prevent abuse, 429 responses | Medium |
| CORS Configuration (env-aware) | Cross-origin access control | Low |
| Request Size Limits | Prevent memory exhaustion, 413 response | Low |
| Unit Test Utilities | Test helpers, example handler tests | Medium |
| Integration Test Scaffold | DB fixtures, transaction rollback | Medium |
| Request Timeout Management | Context-based, prevent hanging | Low |
| Dependency Injection Pattern | Simple constructor injection for testability | Medium |

---

## Anti-Features (Deliberately Excluded)

| Feature | Why |
|---------|-----|
| Authentication (JWT/OAuth) | Project-specific — different every time |
| Authorization / RBAC | Project-specific permission models |
| Caching (Redis) | Optional — not all projects need it |
| API Versioning | Teams choose their own pattern |
| Swagger / OpenAPI | Can drift from code; document as optional |
| Message Queues | Optional, async patterns vary |
| gRPC / WebSocket | Out of scope for v1 REST baseline |
| CI/CD Pipelines | Team-specific |
| Kubernetes configs | Team-specific |

---

## Feature Dependencies

```
Environment Config
  ↓
Structured Logging → Request ID Middleware
  ↓
Database Connection → Database Migrations → Sample DB Module

Rate Limiting, CORS, Request Size Limits (independent)

Unit Test Utilities → Integration Test Scaffold
```

---

## Recommended Phase Order

| Phase | Focus | Key Additions |
|-------|-------|---------------|
| 1 | Foundation fixes | Graceful shutdown, env config, input validation helpers |
| 2 | Observability + DB | Structured logging, request ID, GORM + migrations, sample DB usage |
| 3 | Security + Testing | Rate limiting, CORS, size limits, unit + integration tests |
| 4 | Polish | Timeouts, DI pattern, optional Swagger docs |
