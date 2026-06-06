---
phase: quick
plan: 260606-rqq
type: execute
wave: 1
depends_on: []
files_modified:
  - CLAUDE.md
  - docs/ARCHITECTURE.md
  - README.md
  - docs/DEVELOPMENT.md
  - .planning/codebase/ARCHITECTURE.md
  - .planning/codebase/STRUCTURE.md
  - .planning/codebase/CONVENTIONS.md
  - .planning/research/ARCHITECTURE.md
  - .planning/seeds/SEED-001-golang-based-project-template.md
autonomous: true
requirements: []

must_haves:
  truths:
    - "No file in the repo contains a reference to module/sample/controller.go"
    - "All docs reflect the current file name: module/sample/handlers.go"
  artifacts:
    - path: "CLAUDE.md"
      provides: "Updated module design and architecture references"
      contains: "handlers.go"
    - path: "docs/ARCHITECTURE.md"
      provides: "Accurate handler location documentation"
      contains: "handlers.go"
    - path: ".planning/codebase/CONVENTIONS.md"
      provides: "Accurate module design convention"
      contains: "handlers.go"
  key_links:
    - from: "CLAUDE.md"
      to: "module/sample/handlers.go"
      via: "Module Design section listing"
      pattern: "handlers\\.go"
---

<objective>
Update all documentation references from `controller.go` to `handlers.go` across the codebase after the file rename in `module/sample/`.

Purpose: Keep docs accurate and consistent with the actual file structure so developers cloning and reading this template are not confused by stale filenames.
Output: All doc files reference `handlers.go` instead of `controller.go` for the sample module handler file.
</objective>

<execution_context>
@$HOME/.claude/get-shit-done/workflows/execute-plan.md
@$HOME/.claude/get-shit-done/templates/summary.md
</execution_context>

<context>
@.planning/STATE.md
@.planning/codebase/STRUCTURE.md
@.planning/codebase/CONVENTIONS.md
</context>

<tasks>

<task type="auto">
  <name>Task 1: Update primary documentation files</name>
  <files>CLAUDE.md, docs/ARCHITECTURE.md, README.md, docs/DEVELOPMENT.md</files>
  <action>
    In each file, replace every occurrence of `controller.go` that refers to the sample module handler file with `handlers.go`. The rename is specific to `module/sample/controller.go` — do not change any other use of the word "controller" (e.g., in prose describing the controller pattern, or references to `controller.go` in other modules if any exist).

    Specific changes per file:

    **CLAUDE.md:**
    - Line 87 (Module Design section): `controller.go` → `handlers.go`
    - Line 104 (Architecture Layers, Location field): `module/sample/controller.go` → `module/sample/handlers.go`
    - Line 132 (Key Abstractions, Pattern field): `routes.go, controller.go, service.go` → `routes.go, handlers.go, service.go`
    - Line 146 (Entry Points, Location field): `module/sample/controller.go` → `module/sample/handlers.go`

    **docs/ARCHITECTURE.md:**
    - Line 24: `module/sample/controller.go` → `module/sample/handlers.go`
    - Line 42: any reference → `handlers.go`
    - Line 74: any reference → `handlers.go`
    - Line 90: any reference → `handlers.go`

    **README.md:**
    - Line 391: `controller.go` → `handlers.go`

    **docs/DEVELOPMENT.md:**
    - Line 132: `controller.go` → `handlers.go`
    - Line 170: `controller.go` → `handlers.go`

    Use `sed -i '' 's/controller\.go/handlers.go/g'` on each file, then verify the result with grep to confirm no `controller.go` references remain in these files. This sed pattern is safe because no other file in the module uses `controller.go`.
  </action>
  <verify>
    <automated>grep -n "controller\.go" "/Users/nazwanibrahim/Documents/MRSB/SWDEV/BASED PROJECT/go-project/CLAUDE.md" "/Users/nazwanibrahim/Documents/MRSB/SWDEV/BASED PROJECT/go-project/docs/ARCHITECTURE.md" "/Users/nazwanibrahim/Documents/MRSB/SWDEV/BASED PROJECT/go-project/README.md" "/Users/nazwanibrahim/Documents/MRSB/SWDEV/BASED PROJECT/go-project/docs/DEVELOPMENT.md" 2>/dev/null; echo "Exit: $?"</automated>
  </verify>
  <done>grep returns no matches in any of the four primary doc files. All occurrences have been replaced with handlers.go.</done>
</task>

<task type="auto">
  <name>Task 2: Update planning codebase map files and secondary references</name>
  <files>.planning/codebase/ARCHITECTURE.md, .planning/codebase/STRUCTURE.md, .planning/codebase/CONVENTIONS.md, .planning/research/ARCHITECTURE.md, .planning/seeds/SEED-001-golang-based-project-template.md</files>
  <action>
    Apply the same `controller.go` → `handlers.go` replacement to the planning codebase map files and lower-priority secondary files for full consistency.

    **Target files and known locations:**

    **.planning/codebase/ARCHITECTURE.md:**
    - Line 20: `controller.go` → `handlers.go`
    - Line 81: `controller.go` → `handlers.go`
    - Line 106: `controller.go` → `handlers.go`

    **.planning/codebase/STRUCTURE.md:**
    - Line 15: `controller.go` → `handlers.go`
    - Line 34: `controller.go` → `handlers.go`
    - Line 41: `controller.go` → `handlers.go`
    - Line 49: `controller.go` → `handlers.go`

    **.planning/codebase/CONVENTIONS.md:**
    - Line 36: `controller.go` → `handlers.go`

    **.planning/research/ARCHITECTURE.md** and **.planning/seeds/SEED-001-golang-based-project-template.md:**
    - Replace all occurrences of `controller.go` → `handlers.go` (these are secondary/historical but updating keeps the repo self-consistent)

    Do NOT touch `.planning/phases/` directories — those are historical execution records and should remain as committed.

    Use `sed -i '' 's/controller\.go/handlers.go/g'` on each file individually, then run a final grep across all updated files.
  </action>
  <verify>
    <automated>grep -rn "controller\.go" "/Users/nazwanibrahim/Documents/MRSB/SWDEV/BASED PROJECT/go-project/.planning/codebase/" "/Users/nazwanibrahim/Documents/MRSB/SWDEV/BASED PROJECT/go-project/.planning/research/" "/Users/nazwanibrahim/Documents/MRSB/SWDEV/BASED PROJECT/go-project/.planning/seeds/" 2>/dev/null; echo "Exit: $?"</automated>
  </verify>
  <done>grep returns no matches across all planning codebase map, research, and seed files. Historical phase records in .planning/phases/ are untouched.</done>
</task>

</tasks>

<threat_model>
## Trust Boundaries

| Boundary | Description |
|----------|-------------|
| Local filesystem | All edits are to local documentation files only; no external input or network calls |

## STRIDE Threat Register

| Threat ID | Category | Component | Disposition | Mitigation Plan |
|-----------|----------|-----------|-------------|-----------------|
| T-rqq-01 | Tampering | sed replace on doc files | accept | Replacement is scoped to `controller.go` literal; no logic or secrets involved |
</threat_model>

<verification>
Run a final repo-wide grep to confirm no stray `controller.go` references remain outside the intentionally preserved `.planning/phases/` historical records:

```bash
grep -rn "controller\.go" \
  "/Users/nazwanibrahim/Documents/MRSB/SWDEV/BASED PROJECT/go-project/" \
  --exclude-dir=".planning/phases" \
  --exclude-dir=".git" \
  2>/dev/null
```

Expected result: zero matches.
</verification>

<success_criteria>
- All 9 target documentation files have `handlers.go` where `controller.go` previously appeared
- A repo-wide grep (excluding `.planning/phases/` and `.git`) returns zero hits for `controller.go`
- Historical phase execution records in `.planning/phases/` are unmodified
- The project still builds: `go build ./...` passes without errors (file rename is pre-existing; docs-only change should not affect build)
</success_criteria>

<output>
After completion, create `.planning/quick/260606-rqq-update-all-documentation-references-from/260606-rqq-SUMMARY.md`
</output>
