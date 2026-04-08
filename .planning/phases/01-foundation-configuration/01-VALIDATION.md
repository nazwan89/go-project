---
phase: 1
slug: foundation-configuration
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-04-08
---

# Phase 1 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | go test |
| **Config file** | none — standard Go test tooling |
| **Quick run command** | `go build ./...` |
| **Full suite command** | `go test ./...` |
| **Estimated runtime** | ~5 seconds |

---

## Sampling Rate

- **After every task commit:** Run `go build ./...`
- **After every plan wave:** Run `go test ./...`
- **Before `/gsd-verify-work`:** Full suite must be green
- **Max feedback latency:** 10 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Threat Ref | Secure Behavior | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|------------|-----------------|-----------|-------------------|-------------|--------|
| 1-01-01 | 01 | 1 | FOUND-01 | — | APP_NAME/TIMEZONE/PORT read from env | unit | `go test ./... -run TestConfig` | ❌ W0 | ⬜ pending |
| 1-01-02 | 01 | 1 | FOUND-02 | — | Graceful shutdown within 30s | unit | `go test ./... -run TestGracefulShutdown` | ❌ W0 | ⬜ pending |
| 1-01-03 | 01 | 2 | FOUND-03 | — | 1MB body limit returns 400 | unit | `go test ./... -run TestBodyLimit` | ❌ W0 | ⬜ pending |
| 1-01-04 | 01 | 2 | FOUND-04 | — | Path/query validation rejects bad input | unit | `go test ./... -run TestInputValidation` | ❌ W0 | ⬜ pending |
| 1-01-05 | 01 | 2 | FOUND-05 | — | BodyParser failure returns 400 structured | unit | `go test ./... -run TestBodyParserError` | ❌ W0 | ⬜ pending |
| 1-01-06 | 01 | 3 | SEC-03 | — | No stack traces in error responses | unit | `go test ./... -run TestErrorSanitization` | ❌ W0 | ⬜ pending |
| 1-01-07 | 01 | 3 | DOC-01 | — | .env.example documents all vars | manual | verify file exists and is complete | N/A | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `main_test.go` — integration test stubs for config, shutdown, body limit
- [ ] `utils/error_handler_test.go` — stubs for error sanitization (SEC-03)
- [ ] `module/sample/controller_test.go` — stubs for input validation (FOUND-04, FOUND-05)

*Existing infrastructure: Go stdlib + fiber/v2 test helpers available.*

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| .env.example completeness | DOC-01 | File content review needed | Check all env vars documented with descriptions |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 10s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
