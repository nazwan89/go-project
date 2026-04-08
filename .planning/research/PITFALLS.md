# Pitfalls: Go/Fiber API Baseline

**Date:** 2026-04-08
**Confidence:** HIGH (direct codebase analysis + production patterns)

---

## Critical

### 1. No Graceful Shutdown (Data Loss on Deploy)

**Problem:** `app.Listen()` has no SIGTERM handler. SIGTERM kills in-flight requests mid-operation — DB transactions don't complete, connection pool not closed.

**Fix:**
```go
go func() {
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
    <-sigChan
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    app.ShutdownWithContext(ctx)
}()
```

### 2. Error Responses Expose Internals (OWASP A01)

**Problem:** `InternalServerErrorHandler` returns `err.Error()` to clients — leaks stack traces, file paths, DB schema.

**Fix:** Log full error server-side, return generic message + request_id to client.

### 3. Fiber Context Reuse (Race Condition)

**Problem:** Fiber pools `*fiber.Ctx`. Storing context references or spawning goroutines that outlive the handler causes one request to read another's data.

**Fix:** Extract all values from `c` before going async. Never store `*fiber.Ctx` outside the handler function.

---

## High Priority

### 4. No Input Validation (Injection Risk)

`c.Params("name")` returns raw strings — no length limit, no format check. Echoed back in response = XSS risk. Will be SQL injection risk when GORM is added.

**Fix:** Validate length + format on all params/query strings before use.

### 5. Body Parser Silent 200 on Failure

`c.BodyParser()` errors currently return HTTP 200. Invalid JSON accepted silently.

**Fix:** Always return `c.Status(400)` when `BodyParser` errors.

### 6. No Request Size Limit (DoS)

Fiber accepts unlimited body size. Multi-GB payloads exhaust memory.

**Fix:**
```go
app := fiber.New(fiber.Config{
    BodyLimit: 1 * 1024 * 1024, // 1MB
})
```

### 7. No Rate Limiting (Brute Force)

No rate limiter. Any endpoint can be hammered unlimited times.

**Fix:** Add `limiter.New()` middleware. Note: in-memory limiter doesn't work across multiple instances — document this limitation; teams needing multi-instance must swap to Redis-backed limiter.

### 8. Hardcoded Config (Inflexible Deployments)

Timezone `"Asia/Kuala_Lumpur"` and `AppName "Project Name"` are hardcoded. Same container can't run in multiple environments.

**Fix:** Load all config from `.env` / environment variables.

---

## Medium Priority

### 9. Middleware Order Undocumented

Recovery must be first. Reordering silently breaks things. Document the required order with comments.

### 10. Error Suppression Without Context

```go
loc, _ := time.LoadLocation("Asia/Kuala_Lumpur")  // Silent failure → UTC fallback
```

Log suppressed errors so failures don't go unnoticed in production.

### 11. No CORS Configuration

Browser requests from any origin are blocked. Add `cors.New()` with env-driven allow list.

### 12. Timestamp Has No Timezone Indicator

`"2006-01-02 15:04:05"` format doesn't say which timezone. Use `time.RFC3339` and always output UTC.

---

## Minor

- String concatenation `"Hello, " + name` → use `fmt.Sprintf`
- No request ID middleware → logs can't be correlated
- No response compression → wasted bandwidth

---

## Impact on Roadmap

| Phase | Pitfalls to Address |
|-------|-------------------|
| Phase 1 | Graceful shutdown, error response sanitization, request size limit, hardcoded config |
| Phase 2 | Structured logging, request ID, CORS, timestamp normalization |
| Phase 3 | Rate limiting, input validation, body parser fixes |
| Future | GORM-specific pitfalls (N+1, connection pooling) — research when GORM added |
