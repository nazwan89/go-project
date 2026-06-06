---
phase: quick
plan: 260606-sgw
type: execute
wave: 1
depends_on: []
files_modified:
  - cmd/api/api.go
  - cmd/main.go
autonomous: true
requirements: []
must_haves:
  truths:
    - "cmd/api/api.go exists and exports RegisterRoutes(app *fiber.App)"
    - "cmd/main.go no longer imports project/module/sample directly"
    - "go build ./cmd/... passes with no errors"
  artifacts:
    - path: "cmd/api/api.go"
      provides: "Centralized API route registration"
      exports: ["RegisterRoutes"]
    - path: "cmd/main.go"
      provides: "Updated entry point using api.RegisterRoutes"
  key_links:
    - from: "cmd/main.go"
      to: "cmd/api/api.go"
      via: "api.RegisterRoutes(app)"
      pattern: "api\\.RegisterRoutes"
    - from: "cmd/api/api.go"
      to: "module/sample"
      via: "sample.RegisterRoutes(api)"
      pattern: "sample\\.RegisterRoutes"
---

<objective>
Extract API route registration from cmd/main.go into a dedicated cmd/api/api.go package.

Purpose: Keeps main.go focused on app bootstrap; all module route wiring lives in one place (cmd/api/).
Output: cmd/api/api.go with exported RegisterRoutes, cmd/main.go updated to call it.
</objective>

<execution_context>
@$HOME/.claude/get-shit-done/workflows/execute-plan.md
@$HOME/.claude/get-shit-done/templates/summary.md
</execution_context>

<context>
@.planning/PROJECT.md
@.planning/STATE.md
@cmd/main.go
</context>

<tasks>

<task type="auto">
  <name>Task 1: Create cmd/api/api.go and update cmd/main.go</name>
  <files>cmd/api/api.go, cmd/main.go</files>
  <action>
Create the new file cmd/api/api.go with package api:

```go
package api

import (
	"github.com/gofiber/fiber/v2"
	"project/module/sample"
)

func RegisterRoutes(app *fiber.App) {
	api := app.Group("/api")
	sample.RegisterRoutes(api)
}
```

Then update cmd/main.go:
1. Remove import "project/module/sample" from the import block.
2. Add import "project/cmd/api" to the import block.
3. Replace the route registration block (lines 62-66):
   ```go
   // ========================
   // Register Module Routes
   // ========================
   api := app.Group("/api")
   sample.RegisterRoutes(api)
   ```
   with:
   ```go
   // ========================
   // Register Module Routes
   // ========================
   api.RegisterRoutes(app)
   ```

Keep all other code in main.go exactly as-is, including comments, spacing, and comment block style.

Note: the local variable `api` (from `app.Group("/api")`) is removed, so there is no name collision with the imported `api` package.
  </action>
  <verify>
    <automated>cd /Users/nazwanibrahim/Documents/MRSB/SWDEV/BASED PROJECT/go-project && go build ./cmd/...</automated>
  </verify>
  <done>go build ./cmd/... exits 0 with no errors. cmd/api/api.go exists. cmd/main.go imports project/cmd/api and calls api.RegisterRoutes(app). cmd/main.go no longer imports project/module/sample.</done>
</task>

</tasks>

<threat_model>
## Trust Boundaries

| Boundary | Description |
|----------|-------------|
| N/A | Pure internal refactor — no new trust boundaries introduced |

## STRIDE Threat Register

| Threat ID | Category | Component | Disposition | Mitigation Plan |
|-----------|----------|-----------|-------------|-----------------|
| N/A | — | — | accept | No new attack surface; pure code reorganization within cmd/ |
</threat_model>

<verification>
Run `go build ./cmd/...` from the project root. Build must succeed with exit code 0.
</verification>

<success_criteria>
- cmd/api/api.go exists with package api and exported RegisterRoutes function
- cmd/main.go imports project/cmd/api and does NOT import project/module/sample
- go build ./cmd/... passes
- Runtime behavior is identical (routes still registered under /api)
</success_criteria>

<output>
After completion, create `.planning/quick/260606-sgw-create-cmd-api-api-go-with-a-registerrou/260606-sgw-SUMMARY.md`
</output>
