# Atlas — Autonomous Agent Operating Directives

> **CRITICAL DIRECTIVE FOR ALL AI ASSISTANTS & MODELS**:
> This project operates under strict engineering guidelines and a predictable interaction protocol.
> Follow these instructions on EVERY prompt, in EVERY conversation, regardless of model or context reset.

---

## 1. Core Principles & Safeguards

- **Continuous Verification**: Continuously test and verify your work after every modification (`go build ./...` & `go test`). Catch and fix errors immediately to keep defects at zero.
- **Clarify Before Acting**: If you have ANY doubt, ambiguity, or design trade-off, **STOP AND ASK THE USER FIRST**. Never make unverified assumptions about business logic or architecture.
- **Documentation Sync**: Keep all documentation (`docs/`, `docs/api_docs.md`, `progress.md`, `CHANGELOG.md`, `ROADMAP.md`, `docs/retros/`) 100% updated in lockstep with code changes.
- **Branch & Push Discipline**: Commit logical steps incrementally and push commits to the respective feature branch (`feat/phase-N-...`) and `main` on remote origin.

---

## 2. Operating Protocol ("goahead" Directive)

When the user says **"goahead"** or prompts to continue/start a phase:

1. **Step 1: Inspect Project State**
   - Read `progress.md` to identify the current phase and next phase.
   - Read `guidelines.md` to refresh technical rules, conventions, and locked decisions.
   - Read `future_plans.md` for feature context.

2. **Step 2: Present Implementation Plan (Wait for Approval)**
   - Present a concise, structured implementation plan with:
     - **WHAT**: What components, routes, services, handlers, and templates will be built.
     - **WHY**: The purpose of these components in the Atlas architecture.
     - **HOW**: Step-by-step breakdown of schema, services, handlers, templates, and tests.
     - **USER REVIEW / OPEN QUESTIONS / DOUBTS**: Highlight any trade-offs or design choices requiring feedback.
   - **STOP AND WAIT**: Ask the user to review the plan. Do NOT start modifying files until the user approves or provides adjustments.

3. **Step 3: Incremental Execution with Continuous Verification**
   - Once approved, build the phase step-by-step.
   - **Verify continuously**: Run `go build ./...` and `go test ./...` after every edit to catch errors instantly.
   - **Commit after every logical unit**: Create conventional commits (`feat(...)`, `fix(...)`, `docs(...)`, `test(...)`).

4. **Step 4: Execute Milestone Ceremony, API Docs, Documentation & Git Push**
   - Run full test suite: `powershell -Command "$env:CGO_ENABLED='0'; go test ./tests/unit/... ./tests/integration/... -v"`
   - Write phase retrospective: `docs/retros/phase-{N}.md`
   - Update `docs/api_docs.md`: Record all new/modified HTTP API endpoints, request parameters, HTMX partial responses, and status codes.
   - Update `docs/` module documentation as per changes made.
   - Update `progress.md`: mark tasks completed `✅`, update dates and retros table.
   - Update `CHANGELOG.md` & `ROADMAP.md`.
   - Create git release tag if milestone complete: `git tag -a vX.Y.Z -m "..."`
   - **Push to Remote**: Push commits, feature branch, and tags to GitHub `origin`.

---

## 3. Mandatory Coding Guidelines Summary

- **Layering**: `Handler → Service → Repository → DB` (Strictly enforced. Handlers never touch DB directly).
- **SQLite**: Pure-Go driver (`modernc.org/sqlite`). Connection initialized with WAL mode, foreign keys enabled, and busy timeout 5000ms.
- **SQL Queries**: All queries in `queries/*.sql` and compiled via sqlc into `internal/db/`. No raw string concatenation.
- **Migrations**: All migrations in `migrations/00N_name.sql` and run via embedded Goose (`internal/db/migrate.go`).
- **Templates**: Compiled Templ templates (`web/templates/`). HTML fragments for HTMX responses.
- **Logging**: Zap structured logging (`internal/logger`). JSON in production, colored console in dev.
- **Single-User First-Run**: If `CountUsers() == 0`, `FirstRunGate` redirects all non-static requests to `/setup`.

---

## 4. Persistent Project Files Reference

- `guidelines.md` — Complete engineering law & conventions (26 sections)
- `progress.md` — Live task tracker & milestone retros table
- `future_plans.md` — Total roadmap (v0.1 → v1.0+)
- `ARCHITECTURE.md` — System architecture & request lifecycle
- `docs/database.md` — SQLite schema, indices & FTS5 triggers
- `docs/retros/` — Phase retrospectives

---

*This document is automatically loaded into context by the agent workspace runner.*
