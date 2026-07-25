# Atlas — Autonomous Agent Operating Directives

> **CRITICAL DIRECTIVE FOR ALL AI ASSISTANTS & MODELS**:
> This project operates under strict engineering guidelines and a predictable interaction protocol.
> Follow these instructions on EVERY prompt, in EVERY conversation, regardless of model or context reset.

---

## 1. Operating Protocol ("goahead" Directive)

When the user says **"goahead"** or prompts to continue/start a phase:

1. **Step 1: Inspect Project State**
   - Read `progress.md` to identify the current phase and next phase.
   - Read `guidelines.md` to refresh all technical rules, conventions, and locked decisions.
   - Read `future_plans.md` for feature context.

2. **Step 2: Present Implementation Plan (Wait for Approval)**
   - Present a concise, structured implementation plan with:
     - **WHAT**: What components, routes, services, handlers, and templates will be built.
     - **WHY**: The purpose of these components in the Atlas architecture.
     - **HOW**: Step-by-step breakdown of schema, services, handlers, templates, and tests.
     - **USER REVIEW / OPEN QUESTIONS**: Any design choices or trade-offs requiring feedback.
   - **STOP AND WAIT**: Ask the user to review the plan. Do NOT start modifying files until the user approves or provides adjustments.

3. **Step 3: Execute with Commit Discipline**
   - Once approved, build the phase incrementally.
   - **Commit after every self-contained, logical step** that leaves the codebase passing `go build ./...` and `go test ./...`.
   - Use Conventional Commits (`feat(...)`, `fix(...)`, `docs(...)`, `test(...)`, `build(...)`).

4. **Step 4: Execute Milestone Ceremony & Documentation**
   - Run tests: `powershell -Command "$env:CGO_ENABLED='0'; go test ./tests/unit/... ./tests/integration/... -v"`
   - Write phase retrospective: `docs/retros/phase-{N}.md`
   - Update `progress.md`: mark tasks completed `✅`, update dates and retros table.
   - Update `CHANGELOG.md` with new features.
   - Update `ROADMAP.md` if milestone shipped.
   - Tag release in git if milestone complete: `git tag -a vX.Y.Z -m "..."`

---

## 2. Mandatory Coding Guidelines Summary

- **Layering**: `Handler → Service → Repository → DB` (Strictly enforced. Handlers never touch DB directly).
- **SQLite**: Pure-Go driver (`modernc.org/sqlite`). Connection initialized with WAL mode, foreign keys enabled, and busy timeout 5000ms.
- **SQL Queries**: All queries in `queries/*.sql` and compiled via sqlc into `internal/db/`. No raw string concatenation.
- **Migrations**: All migrations in `migrations/00N_name.sql` and run via embedded Goose (`internal/db/migrate.go`).
- **Templates**: Compiled Templ templates (`web/templates/`). HTML fragments for HTMX responses.
- **Logging**: Zap structured logging (`internal/logger`). JSON in production, colored console in dev.
- **Single-User First-Run**: If `CountUsers() == 0`, `FirstRunGate` redirects all non-static requests to `/setup`.

---

## 3. Persistent Project Files Reference

- `guidelines.md` — Complete engineering law & conventions (25 sections)
- `progress.md` — Live task tracker & milestone retros table
- `future_plans.md` — Total roadmap (v0.1 → v1.0+)
- `ARCHITECTURE.md` — System architecture & request lifecycle
- `docs/database.md` — SQLite schema, indices & FTS5 triggers
- `docs/retros/` — Phase retrospectives

---

*This document is automatically loaded into context by the agent workspace runner.*
