package sample

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"

	"project/utils"
)

// setupApp creates a minimal Fiber app wired the same way as main.go
// but without starting a real server (uses httptest).
func setupApp() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: utils.ErrorHandler,
		BodyLimit:    1 * 1024 * 1024,
	})
	api := app.Group("/api")
	RegisterRoutes(api)
	return app
}

func TestBodyParserError_Returns400(t *testing.T) {
	app := setupApp()

	req := httptest.NewRequest(http.MethodPost, "/api/sample/hello-form", strings.NewReader("not-valid-json-or-form"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestBodyParserError_StructuredResponse(t *testing.T) {
	app := setupApp()

	req := httptest.NewRequest(http.MethodPost, "/api/sample/hello-form", strings.NewReader("{bad json"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, `"error"`) {
		t.Errorf("expected JSON body with 'error' key, got: %s", bodyStr)
	}
	if !strings.Contains(bodyStr, "Bad Request") {
		t.Errorf("expected 'Bad Request' in response body, got: %s", bodyStr)
	}
}

func TestPathParamValidation_Valid(t *testing.T) {
	app := setupApp()

	req := httptest.NewRequest(http.MethodGet, "/api/sample/hello/World", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 for valid path param, got %d", resp.StatusCode)
	}
}

func TestPathParamValidation_TooLong(t *testing.T) {
	app := setupApp()

	longName := strings.Repeat("a", 101)
	req := httptest.NewRequest(http.MethodGet, "/api/sample/hello/"+longName, nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for 101-char path param, got %d", resp.StatusCode)
	}
}

func TestErrorSanitization_NoErrMessage(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: utils.ErrorHandler,
	})

	// Route that triggers an internal server error with a sensitive message
	app.Get("/test-error", func(c *fiber.Ctx) error {
		return fiber.NewError(500, "sensitive internal detail: db password=hunter2")
	})

	req := httptest.NewRequest(http.MethodGet, "/test-error", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)

	if strings.Contains(bodyStr, "hunter2") {
		t.Errorf("sensitive error content leaked in response: %s", bodyStr)
	}
	if strings.Contains(bodyStr, "db password") {
		t.Errorf("sensitive error content leaked in response: %s", bodyStr)
	}
}

func TestErrorSanitization_HasRequestID(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: utils.ErrorHandler,
	})

	app.Get("/test-500", func(c *fiber.Ctx) error {
		return fiber.NewError(500, "internal error")
	})

	req := httptest.NewRequest(http.MethodGet, "/test-500", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "request_id") {
		t.Errorf("expected 'request_id' in 500 response body, got: %s", bodyStr)
	}
}
