---
id: SEED-001
status: dormant
planted: 2026-04-08
planted_during: initial setup
trigger_when: when starting a new backend service
scope: Medium
---

# SEED-001: Create based project using golang

## Why This Matters

This seed addresses all four pillars of a healthy Go backend foundation:

- **Standardize Go structure** — Establish a consistent, reusable project structure for all future Go services
- **Speed up bootstrapping** — Reduce time-to-start when creating new Go backend services
- **Enforce best practices** — Bake in patterns like layered architecture, logging, health checks, and error handling from day one
- **Team consistency** — Ensure all team members start from the same foundation

## When to Surface

**Trigger:** When starting a new backend service

This seed should be presented during `/gsd-new-milestone` when the milestone
scope matches any of these conditions:
- A new Go microservice or backend service is being created
- A new module or domain service is being added to the system
- The team decides to spin off a standalone Go API

## Scope Estimate

**Medium** — A phase or two of focused work. Involves scaffolding the project layout,
wiring up routing, middleware, logging, health checks, and documenting the template
for reuse across services.

## Breadcrumbs

Related code and decisions found in the current codebase:

- `main.go` — Entry point; current app bootstrap and router setup
- `go.mod` — Module definition; current dependencies (Echo, etc.)
- `module/sample/handlers.go` — Example handler pattern to replicate
- `module/sample/routes.go` — Route registration pattern to formalize
- `module/sample/service.go` — Service layer separation pattern
- `module/sample/form.go` — Request form/validation pattern
- `utils/date_time.go` — Utility layer structure
- `utils/error_handler.go` — Centralized error handling pattern
- `Dockerfile` — Container setup already established

## Notes

The current `go-project` repo already demonstrates the pattern informally via the
`module/sample` module. The based project template should codify this into a
reusable scaffold — possibly as a `go-template` repo or a `make new-module` generator.
Consider including: layered architecture (controller/service/repo), structured logging,
health check endpoint, environment config loading, and a standard Makefile.
