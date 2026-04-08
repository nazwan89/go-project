---
phase: 01-foundation-configuration
plan: "01"
subsystem: config
tags: [config, shutdown, security, env]
dependency_graph:
  requires: []
  provides: [AppConfig, LoadConfig, CurrentTimestamp(loc), graceful-shutdown]
  affects: [main.go, utils/date_time.go, utils/error_handler.go]
tech_stack:
  added: []
  patterns: [env-driven config, graceful SIGTERM shutdown, location-aware timestamps]
key_files:
  created:
    - utils/config.go
    - utils/config_test.go
  modified:
    - utils/date_time.go
    - utils/error_handler.go
    - main.go
    - .env
decisions:
  - "CurrentTimestamp signature changed to accept *time.Location; error handlers pass time.UTC for consistency and Phase 2 UTC normalization preparation"
  - "error_handler.go updated in Task 1 (not Task 2) to unblock compilation of test suite"
  - "413 StatusRequestEntityTooLarge remapped to 400 Bad Request in ErrorHandler per FOUND-03"
  - "InternalServerErrorHandler no longer exposes err.Error() in response body (SEC-03 partial)"
metrics:
  duration: "~2 minutes"
  completed: "2026-04-08T14:21:42Z"
  tasks_completed: 2
  files_changed: 6
---

# Phase 1 Plan 01: Externalize Config and Graceful Shutdown Summary

**One-liner:** Env-driven AppConfig struct with LoadConfig(), location-aware CurrentTimestamp, and SIGTERM graceful shutdown with 30s timeout via goroutine + signal channel pattern.

---

## Tasks Completed

| # | Task | Commit | Files |
|---|------|--------|-------|
| 1 (RED) | Add failing tests for LoadConfig and CurrentTimestamp(loc) | 12881f7 | utils/config_test.go |
| 1 (GREEN) | Create AppConfig/LoadConfig, env-driven CurrentTimestamp, update error handlers | 81e1cd1 | utils/config.go, utils/date_time.go, utils/error_handler.go |
| 2 | Config-driven main.go with graceful SIGTERM shutdown | cb32d8e | main.go, .env |

---

## What Changed

### utils/config.go (new)
- `AppConfig` struct: `AppName`, `Port`, `Timezone`, `Location *time.Location`
- `LoadConfig()` reads `APP_NAME`, `PORT`, `TIMEZONE` from env with defaults: `"go-project"`, `"8080"`, `"UTC"`
- Invalid `TIMEZONE` silently falls back to `time.UTC` — startup never panics (T-01-04 mitigated)

### utils/date_time.go (modified)
- `CurrentTimestamp()` zero-arg → `CurrentTimestamp(loc *time.Location)` — removes hardcoded `"Asia/Kuala_Lumpur"`
- Nil-safe: falls back to `time.UTC` if loc is nil
- Added comment noting Phase 2 will normalize to RFC3339 UTC

### utils/error_handler.go (modified)
- All `CurrentTimestamp()` calls updated to `CurrentTimestamp(time.UTC)`
- Added `fiber.StatusRequestEntityTooLarge` (413) case → remapped to 400 Bad Request (T-01-01 mitigated)
- `InternalServerErrorHandler` no longer exposes `err.Error()` in response body (T-01-02 mitigated)
- Added `"time"` import

### main.go (modified)
- Calls `utils.LoadConfig()` immediately after `godotenv.Load()`
- `fiber.Config.AppName` now uses `cfg.AppName` (removes hardcoded `"Project Name"`)
- Added `BodyLimit: 1 * 1024 * 1024` (1MB DoS protection, T-01-01 / FOUND-03)
- Added `ReadTimeout: 30 * time.Second` (bounds keepalive connections, T-01-03 mitigated)
- Server starts in goroutine (non-blocking)
- Blocks on signal channel; calls `app.ShutdownWithTimeout(30 * time.Second)` on SIGTERM/SIGINT (FOUND-02)
- `cfg.Port` used instead of `os.Getenv("PORT")`

### .env (modified)
- Added `APP_NAME=go-project` and `TIMEZONE=UTC` entries alongside existing `PORT=8080`

---

## Verification Results

```
go test ./utils/ -run TestLoadConfig -count=1 -v
=== RUN   TestLoadConfig_Defaults    --- PASS
=== RUN   TestLoadConfig_FromEnv     --- PASS
=== RUN   TestLoadConfig_InvalidTimezone --- PASS
=== RUN   TestCurrentTimestamp_UsesConfig --- PASS
PASS ok  project/utils  0.491s

go build ./...   # exits 0
go vet ./...     # exits 0
grep "Asia/Kuala_Lumpur" utils/date_time.go utils/config.go utils/error_handler.go  # no matches
grep '"Project Name"' main.go  # no matches
grep "ShutdownWithTimeout" main.go  # match found
```

---

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 3 - Blocking] Updated error_handler.go in Task 1 rather than waiting for Task 2**
- **Found during:** Task 1 GREEN phase
- **Issue:** Changing `CurrentTimestamp()` signature broke compilation of `error_handler.go`, preventing the test suite from running
- **Fix:** Updated all `CurrentTimestamp()` call sites in `error_handler.go` to `CurrentTimestamp(time.UTC)` immediately; also added 413->400 remapping and removed `err.Error()` exposure (items that were planned for Task 2 anyway)
- **Files modified:** `utils/error_handler.go`
- **Commit:** 81e1cd1

---

## Known Stubs

None. All data flows are wired. `CurrentTimestamp` uses the real system clock with the configured location.

---

## Threat Flags

No new security-relevant surface introduced beyond what the plan's threat model already covers.

All four T-01-xx threat mitigations are in place:
- T-01-01: `BodyLimit: 1MB` + 413→400 remapping in ErrorHandler
- T-01-02: `InternalServerErrorHandler` no longer returns `err.Error()`
- T-01-03: `ReadTimeout: 30s` bounds keepalive connections during shutdown
- T-01-04: `LoadConfig()` falls back to `time.UTC` on invalid TIMEZONE

---

## Self-Check: PASSED

All files exist. All commits verified.

| Item | Status |
|------|--------|
| utils/config.go | FOUND |
| utils/config_test.go | FOUND |
| utils/date_time.go | FOUND |
| utils/error_handler.go | FOUND |
| main.go | FOUND |
| .env | FOUND |
| 01-01-SUMMARY.md | FOUND |
| commit 12881f7 | FOUND |
| commit 81e1cd1 | FOUND |
| commit cb32d8e | FOUND |
