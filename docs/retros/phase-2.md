# Phase 2 Retrospective — Auth + Setup Wizard

> **Phase**: 2
> **Started**: 2026-07-26
> **Completed**: 2026-07-26
> **Time Estimated**: 1 day
> **Time Actual**: 1 day

---

## What Was Built

Phase 2 implemented the first-run onboarding wizard and session-based authentication for Atlas:
- First-Run Setup Service (`internal/setup/service.go`): `IsFirstRun`, `CreateFirstUser` (bcrypt cost 12), and `SeedDemoData` (populates demo project, tasks, notes, journal entry).
- FirstRunGate Middleware (`internal/setup/middleware.go`): Intercepts incoming HTTP requests and redirects any non-static request to `/setup` if no user accounts exist in SQLite.
- First-Run Wizard Templates (`web/templates/setup/wizard.templ` & `wizard_templ.go`): Multi-step setup form (display name, username, email, password, timezone, theme) and Demo Data choice screen.
- First-Run Setup Handler (`internal/setup/handler.go`): Handlers for `GET /setup`, `POST /setup`, `GET /setup/demo-choice`, and `POST /setup/seed`.
- Session Store & Auth Service (`internal/auth/session.go` & `service.go`): Gorilla session store wrapper with HTTP-only and SameSite=Lax cookies, bcrypt authentication, and SQLite session persistence.
- AuthRequired Middleware (`internal/auth/middleware.go`): Protects application routes and populates request context with authenticated `*db.User`.
- Login Handlers & Template (`internal/auth/handler.go` & `web/templates/auth/login.templ`): Handlers for `GET /login`, `POST /login`, and `POST /logout`.
- Automated Unit & Integration Tests (`tests/unit/auth_test.go` & `tests/integration/setup_flow_test.go`).

---

## What Went Well

- `FirstRunGate` middleware creates an effortless first-run onboarding experience: zero-config start → `/setup` redirect → account creation → demo choice → dashboard.
- Session cookie handling coupled with SQLite session token persistence guarantees secure, resilient user authentication.

---

## What Was Harder Than Expected

- Ensuring Templ template generation compatibility across standard Go `io.Writer` interfaces.

---

## Decisions Made During Build

| Decision | Reason | Impact |
|----------|--------|--------|
| `FirstRunGate` middleware | Avoids manual setup flag checks in every handler | Automatic setup redirection across all endpoints |
| 1-Year Session for Owner | Single-user personal OS should not constantly ask owner to re-login | Seamless long-term usability |

---

## Bugs Found and Fixed

| Bug | Root Cause | Fix |
|-----|-----------|-----|
| Missing `templ` dependency in `go.mod` for compiled template files | New Templ templates added to `web/templates/` | Ran `go mod tidy` to add `github.com/a-h/templ` |

---

## What to Carry Forward

- Keep authentication middleware light by relying on request context values.
- Continue writing integration tests that exercise full HTTP request → response → database verification.

---

## Links

- Phase 2 Commit: feat(auth): implement first-run setup wizard, session authentication, and middleware
- Documentation: `docs/trd.md`, `guidelines.md`
