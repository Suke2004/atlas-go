# Contributing to Atlas

Thank you for your interest in contributing to Atlas. This document covers everything you need to know.

---

## Philosophy

Atlas is a **personal operating system** — a single application that replaces a developer's entire tool stack. Every contribution should serve that vision: clarity, performance, local ownership, and long-term maintainability over clever abstractions.

---

## Before You Start

1. Read guidelines.md — the engineering law of this project.
2. Read docs/architecture.md — understand the system before changing it.
3. Check open issues and the project board before starting new work.
4. For large features, open an issue first and describe the approach.

---

## Development Setup

### Prerequisites

`ash
go 1.23+
templ (go install github.com/a-h/templ/cmd/templ@latest)
goose (go install github.com/pressly/goose/v3/cmd/goose@latest)
sqlc (go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest)
air  (go install github.com/air-verse/air@latest)
golangci-lint
`

### Local Setup

`ash
git clone https://github.com/Suke2004/atlas-go
cd atlas-go

cp .env.example .env
# Fill ATLAS_SESSION_SECRET with 32 random bytes

make setup-dirs    # Create /data/db, /data/uploads, etc.
make migrate       # Apply all migrations
make dev           # Start dev server with hot-reload
`

Open http://localhost:8080 — the setup wizard will appear on first run.

---

## Workflow

### Branch Strategy

`
main       → Protected. Stable releases only.
develop    → Integration branch. Target your PRs here.
feature/*  → Feature branches.
fix/*      → Bug fix branches.
docs/*     → Documentation branches.
`

### Creating a Feature Branch

`ash
git checkout develop
git pull origin develop
git checkout -b feature/your-feature-name
`

### Commit Convention

Use **Conventional Commits**:

`
feat(scope): add X
fix(scope): resolve Y
refactor(scope): simplify Z
docs(scope): update W
test(scope): add tests for V
`

Valid scopes: uth, dashboard, 	asks, projects, 
otes, journal, search, settings, i, inance, db, setup, docs, ci

### Opening a Pull Request

- Target develop (not main)
- PR title follows commit convention: eat(tasks): add dependency validation
- Fill in the PR template
- Ensure all CI checks pass
- Request review

---

## Code Standards

### Must Pass Before PR Merge

`ash
go build ./...                     # No compile errors
go test ./... -race -cover         # All tests pass, no race conditions
go vet ./...                       # No vet warnings
golangci-lint run                  # No lint errors
templ generate                     # Templates regenerate cleanly
`

### Coverage Requirements

| Layer | Minimum |
|-------|---------|
| Services | 80% |
| Repositories | 70% |
| Handlers | 60% |

### Key Rules

- Follow the Handler → Service → Repository → DB layering — strictly.
- No raw SQL string building — all queries through sqlc.
- No CDN dependencies — all JS/CSS served locally.
- No package-level global variables.
- Always wrap errors with context: mt.Errorf("tasks: create: %w", err).

---

## Issue Labels

| Label | Use |
|-------|-----|
| ug | Something is broken |
| enhancement | Improvement to existing feature |
| eature | New feature request |
| documentation | Docs improvements |
| good first issue | Easy entry point for new contributors |
| help wanted | We need help here |
| performance | Speed / efficiency issue |
| security | Security concern |
| 	esting | Test coverage gaps |
| efactor | Code quality improvement |

---

## Questions?

Open a Discussion (not an Issue) for questions and ideas.

---

*Atlas is a personal operating system. Keep that identity in every change.*
