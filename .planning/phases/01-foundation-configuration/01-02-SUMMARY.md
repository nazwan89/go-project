---
phase: 01-foundation-configuration
plan: 02
subsystem: utils, module/sample
tags: [validation, security, error-handling, input-sanitization]
dependency_graph:
  requires: []
  provides: [utils/validate.go, ValidatePathParam, ValidateQueryParam, request_id-in-500-responses]
  affects: [utils/error_handler.go, module/sample/controller.go, utils/date_time.go, main.go]
tech_stack:
  added: [github.com/go-playground/validator/v10 v10.30.2]
  patterns: [struct-tag validation, fiber.NewError for safe 400s, uuid request_id for 500s]
key_files:
  created: [utils/validate.go, utils/validate_test.go, module/sample/controller_test.go]
  modified: [utils/error_handler.go, utils/date_time.go, module/sample/controller.go, main.go, go.mod, go.sum]
decisions:
  - "Used go-playground/validator/v10 struct-tag approach for ValidatePathParam/ValidateQueryParam to keep validation declarative"
  - "Used fiber.NewError(400) in responseHelloWithForm instead of return err directly, preventing raw parser error exposure"
  - "Updated CurrentTimestamp() signature to CurrentTimestamp(loc *time.Location) in this worktree to unblock compilation (mirrors Plan 01 change)"
  - "413 oversized body remapped to 400 Bad Request in ErrorHandler (aligns with T-02-05 threat mitigation)"
metrics:
  duration: "~15 minutes"
  completed: "2026-04-08"
  tasks_completed: 2
  files_changed: 7
---

# Phase 01 Plan 02: Input Protection and Error Sanitization Summary

Input protection and error sanitization implemented: ValidatePathParam/ValidateQueryParam helpers via go-playground/validator/v10, UUID-tagged internal error responses with log-only real errors, and 400 returns on body parse failure.

## What Was Built

### Task 1: Validator, utils/validate.go, request_id in InternalServerErrorHandler

**utils/validate.go** — New file providing two public validation helpers:
- `ValidatePathParam(value string) error` — rejects empty, >100 chars, non-alphanumeric inputs using `required,min=1,max=100,alphanum` tag
- `ValidateQueryParam(value string) error` — allows empty (optional), rejects >200 chars using `omitempty,min=1,max=200` tag

**utils/error_handler.go** — Updated `InternalServerErrorHandler`:
- Generates UUID v4 `request_id` via `uuid.New().String()`
- Logs real error internally: `log.Printf("[ERROR] request_id=%s path=%s error=%v", ...)`
- Returns only generic message to client — no `err.Error()` content exposed
- Added 413→400 remapping in `ErrorHandler` (T-02-05 mitigation)
- Updated all `CurrentTimestamp()` calls to `CurrentTimestamp(time.UTC)`

**utils/date_time.go** — Updated `CurrentTimestamp()` signature to `CurrentTimestamp(loc *time.Location)` (mirrors Plan 01 change; required to unblock compilation in this worktree)

**Commits:** `57beade`

### Task 2: Body parser error handling and path/query param validation in controller

**module/sample/controller.go** — Three fixes applied:
- `responseHelloWithForm`: replaced silent `return c.JSON(Response{Message: "Invalid request body"})` (HTTP 200) with `return fiber.NewError(fiber.StatusBadRequest, "Invalid or malformed request body")` (HTTP 400)
- `responseHelloWithName`: added `utils.ValidatePathParam(name)` check before using path param
- `responseHelloWithQuery` + `responseHelloWithService`: added `utils.ValidateQueryParam(name)` check before using query param

**module/sample/controller_test.go** — New test file with:
- `TestBodyParserError_Returns400` — verifies 400 on malformed body
- `TestBodyParserError_StructuredResponse` — verifies JSON error structure with "error" and "Bad Request"
- `TestPathParamValidation_Valid` — verifies 200 for valid path param
- `TestPathParamValidation_TooLong` — verifies 400 for 101-char path param
- `TestErrorSanitization_NoErrMessage` — verifies sensitive error content does not leak
- `TestErrorSanitization_HasRequestID` — verifies request_id present in 500 response

**Commits:** `5a2100b`

## Verification Results

```
go build ./...                                    PASS
go test ./utils/ -run "TestValidatePath|TestValidateQuery"   PASS (7/7)
go test ./module/sample/ -run "TestBodyParser|TestPathParam|TestErrorSanitization"  PASS (6/6)
go test ./...                                     PASS (all packages)
```

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 3 - Blocking] Updated CurrentTimestamp() signature to accept loc parameter**
- **Found during:** Task 1 — error_handler.go needed updating
- **Issue:** Plan 01 (running in parallel) changes `CurrentTimestamp()` to `CurrentTimestamp(loc *time.Location)`. This worktree needed to match for compilation.
- **Fix:** Updated `utils/date_time.go` and all callers (`error_handler.go`, `main.go`) to use the new signature. Passes `time.UTC` in error handlers and `nil` in `main.go` root route.
- **Files modified:** `utils/date_time.go`, `utils/error_handler.go`, `main.go`
- **Commit:** `57beade`

**2. [Rule 2 - Missing critical functionality] Remap 413 → 400 in ErrorHandler**
- **Found during:** Task 1 — reviewing T-02-05 threat: oversized body returning 413 instead of 400
- **Issue:** Plan mentions "413 remapped to 400 in ErrorHandler" (T-02-05) but `ErrorHandler` had no explicit 413 case
- **Fix:** Added `case 413: return BadRequestHandler(c, "Request body too large")` to ErrorHandler switch
- **Files modified:** `utils/error_handler.go`
- **Commit:** `57beade`

**3. [Rule 1 - Bug] Corrected test route path from /api/sample/form to /api/sample/hello-form**
- **Found during:** Task 2 — writing controller_test.go
- **Issue:** Plan's test snippet used `/api/sample/form` but the actual route registered in `routes.go` is `/hello-form` → `/api/sample/hello-form`
- **Fix:** Used correct route path `/api/sample/hello-form` in all controller tests
- **Files modified:** `module/sample/controller_test.go`
- **Commit:** `5a2100b`

## Known Stubs

None. All validation and error handling is fully wired.

## Threat Surface Scan

No new network endpoints, auth paths, or file access patterns introduced. All changes harden existing surfaces per the plan's threat model.

## Self-Check: PASSED

- `utils/validate.go` — EXISTS
- `utils/validate_test.go` — EXISTS
- `utils/error_handler.go` — EXISTS (contains request_id, uuid.New, no err.Error() in responses)
- `module/sample/controller.go` — EXISTS (contains ValidatePathParam, ValidateQueryParam, fiber.NewError 400)
- `module/sample/controller_test.go` — EXISTS
- Commit `57beade` — EXISTS
- Commit `5a2100b` — EXISTS
- `go test ./...` — PASS
