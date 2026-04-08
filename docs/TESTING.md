<!-- generated-by: gsd-doc-writer -->
# Testing

This document describes the test structure, how to run tests, what each test file covers, and how to write new tests for this project.

## Test Framework and Setup

The project uses the Go standard library `testing` package — no third-party test framework is required. Tests run with the standard `go test` toolchain included with Go 1.25.4.

**No special setup is required before running tests.** Dependency installation is managed by Go modules; running any `go test` command will resolve dependencies automatically via `go.sum`.

For controller tests, the project uses Go's built-in `net/http/httptest` package alongside Fiber's `app.Test()` method to spin up an in-memory HTTP server — no real network port is bound during testing.

## Running Tests

**Run all tests across the entire project:**

```bash
go test ./...
```

**Run tests for a specific package:**

```bash
# Controller tests (module/sample package)
go test ./module/sample/

# Config tests
go test ./test/config/

# Utils tests
go test ./test/utils/
```

**Run a single test function:**

```bash
go test ./module/sample/ -run TestBodyParserError_Returns400
```

**Run tests with verbose output:**

```bash
go test ./... -v
```

**Run tests with the race detector:**

```bash
go test -race ./...
```

## Test File Structure

Tests are organized in two locations:

```
module/
  sample/
    controller_test.go   # HTTP handler integration tests (same package: sample)
test/
  config/
    app_config_test.go   # Config loading unit tests (package: config_test)
  utils/
    date_time_test.go    # Timestamp utility unit tests (package: utils_test)
    validate_test.go     # Input validation unit tests (package: utils_test)
```

### module/sample/controller_test.go

**Package:** `sample` (white-box — same package as the handlers under test)

This file tests the HTTP layer end-to-end using an in-memory Fiber app created by `setupApp()`. The helper mirrors the real `main.go` configuration (same `ErrorHandler`, same `BodyLimit`, same route registration) without starting a real server.

| Test | What it verifies |
|---|---|
| `TestBodyParserError_Returns400` | Malformed JSON body to `POST /api/sample/hello-form` returns HTTP 400 |
| `TestBodyParserError_StructuredResponse` | Bad JSON body produces a JSON response containing `"error"` and `"Bad Request"` |
| `TestPathParamValidation_Valid` | A valid name in `GET /api/sample/hello/:name` returns HTTP 200 |
| `TestPathParamValidation_TooLong` | A 101-character path parameter returns HTTP 400 |
| `TestErrorSanitization_NoErrMessage` | A 500 error containing sensitive text does not expose that text in the response body |
| `TestErrorSanitization_HasRequestID` | A 500 error response includes a `request_id` field |

### test/config/app_config_test.go

**Package:** `config_test` (black-box — external package testing the `config` package's public API)

Unit tests for `config.LoadAppConfig()`. Each test manipulates environment variables using `os.Setenv` / `os.Unsetenv` and defers cleanup to avoid cross-test contamination.

| Test | What it verifies |
|---|---|
| `TestLoadAppConfig_Defaults` | `APP_NAME`, `PORT`, and `TIMEZONE` all fall back to expected defaults when unset |
| `TestLoadAppConfig_FromEnv` | Environment variables are correctly read and the `Location` field is populated |
| `TestLoadAppConfig_InvalidTimezone` | An unrecognisable timezone string causes `Location` to fall back to `time.UTC` |

### test/utils/validate_test.go

**Package:** `utils_test` (black-box — tests the exported `ValidatePathParam` and `ValidateQueryParam` functions)

Unit tests for the input-validation utilities. Uses dot-import (`. "project/utils"`) so function names are unqualified.

| Test | What it verifies |
|---|---|
| `TestValidatePathParam_Valid` | Alphanumeric strings of acceptable length return no error |
| `TestValidatePathParam_TooLong` | A 101-character string returns an error |
| `TestValidatePathParam_Empty` | An empty string returns an error |
| `TestValidatePathParam_InvalidChars` | Strings with spaces, slashes, `@`, hyphens, and underscores return errors |
| `TestValidateQueryParam_Empty` | An empty string returns no error (query params are optional) |
| `TestValidateQueryParam_Valid` | Alphanumeric strings of acceptable length return no error |
| `TestValidateQueryParam_TooLong` | A 201-character string returns an error |

### test/utils/date_time_test.go

**Package:** `utils_test` (black-box)

Unit test for the `CurrentTimestamp` utility function.

| Test | What it verifies |
|---|---|
| `TestCurrentTimestamp_UsesConfig` | `CurrentTimestamp` returns a non-empty string for both `Asia/Kuala_Lumpur` and `time.UTC` locations. The test is skipped if the `Asia/Kuala_Lumpur` timezone data is unavailable on the host. |

## Coverage Requirements

No coverage threshold is configured. There is no `.nycrc`, coverage section in `go.mod`, or CI enforcement of a minimum percentage. To generate a coverage report manually:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## CI Integration

No CI pipeline is currently configured (no `.github/workflows/` directory is present in the repository). Tests are run locally by developers before committing.

## Writing New Tests

### Controller tests (HTTP handlers)

Place new handler tests in `module/sample/controller_test.go` (or a new `_test.go` file inside the same `module/<name>/` package). Use the `setupApp()` helper pattern to get a pre-configured Fiber app:

```go
func TestMyHandler_ReturnsExpected(t *testing.T) {
    app := setupApp()

    req := httptest.NewRequest(http.MethodGet, "/api/sample/hello", nil)
    resp, err := app.Test(req, -1)
    if err != nil {
        t.Fatalf("request failed: %v", err)
    }

    if resp.StatusCode != http.StatusOK {
        t.Errorf("expected 200, got %d", resp.StatusCode)
    }
}
```

The `-1` timeout passed to `app.Test` disables the default 200 ms timeout — use a positive millisecond value when testing endpoints with intentional delays.

### Unit tests (utils / config)

Place utility and config unit tests under `test/<package>/` using the `<package>_test` package name (black-box testing). This matches the existing convention in `test/config/` and `test/utils/`.

```go
package utils_test

import (
    "testing"
    . "project/utils"
)

func TestMyUtility_DoesExpected(t *testing.T) {
    result := MyUtility("input")
    if result != "expected" {
        t.Errorf("expected %q, got %q", "expected", result)
    }
}
```

### Naming convention

Follow the existing pattern: `Test<Subject>_<Condition>` (e.g., `TestValidatePathParam_TooLong`, `TestBodyParserError_Returns400`). This makes it easy to target individual tests with `-run` and to read test output at a glance.

### Service-layer tests

Business logic in `service.go` files can be tested as white-box tests in the same package (no HTTP setup needed):

```go
package sample

import "testing"

func TestGenerateGreeting_IncludesName(t *testing.T) {
    got := generateGreeting("Alice")
    if got != "Hello, Alice! This message was generated by the service layer." {
        t.Errorf("unexpected greeting: %s", got)
    }
}
```
