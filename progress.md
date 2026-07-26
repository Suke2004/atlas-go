# Atlas — Progress Tracker

> **Last Updated**: 2026-07-26
> **Current Phase**: Phase 0 — Not Started
> **Overall v1.0 Progress**: 0%

---

## Before Every Working Session — Read This

> **Explain before you build.** Before writing implementation code, answer:
> ```
> WHAT  — What exactly am I building right now?
> WHY   — Why does Atlas need this? What problem does it solve?
> HOW   — How will it work, step by step?
> ```
> State this out loud (or in chat) before touching the keyboard.
> If you cannot answer all three, you are not ready to build yet.

> **Commit every step.** After each logical unit of work:
> ```bash
> go build ./...        # Must pass
> go test ./... -race   # Must pass
> go vet ./...          # Must pass
> git add -p            # Stage intentionally (not git add .)
> git commit -m "feat(scope): description"
> ```
> Never end a session without committing everything that works.

---

## Legend

| Symbol | Meaning |
|--------|---------|
| ✅ | Complete |
| 🔄 | In Progress |
| ⬜ | Not Started |
| ❌ | Blocked |
| 🔁 | Revisit Later |

---

## v1.0 Foundation — Phase Tracker

| Phase | Name | Status | Started | Completed | Notes |
|-------|------|--------|---------|-----------|-------|
| 0 | Scaffolding | ✅ | 2026-07-26 | 2026-07-26 | Go server, Makefile, Air, Docker, /data layout |
| 1 | Database | ✅ | 2026-07-26 | 2026-07-26 | Migrations, sqlc, WAL, FTS5 triggers |
| 2 | Auth + Setup Wizard | ✅ | 2026-07-26 | 2026-07-26 | First-run wizard, login/logout, sessions |
| 3 | Layout Shell | ✅ | 2026-07-26 | 2026-07-26 | Sidebar, topbar, theme system, static assets |
| 4 | Projects Module | ✅ | 2026-07-26 | 2026-07-26 | CRUD, GitHub repo stats, Tech Stack badges, milestones |
| 5 | Tasks Module | ✅ | 2026-07-26 | 2026-07-26 | List, 3-column Kanban, Slate Obsidian theme, Lucide icons, task dependencies |
| 6 | Knowledge Base | ✅ | 2026-07-26 | 2026-07-26 | Markdown editor, Quick Capture, Templates, [[Wiki Links]], Backlinks, 30s autosave |
| 7 | Dashboard | ✅ | 2026-07-26 | 2026-07-26 | Executive Command Center live widgets |
| 8 | Journal Module | ✅ | 2026-07-26 | 2026-07-26 | Executive Mind-Sync, Mood/energy/sleep telemetry, 4-quadrant reflection |
| 9 | Global Search | ✅ | 2026-07-26 | 2026-07-26 | Ctrl+K Command Palette modal overlay, SQLite FTS5 |
| 10 | Settings + Polish | ✅ | 2026-07-26 | 2026-07-26 | Theme toggle, skeleton loaders, toasts |
| 11 | Finance Engine | ✅ | 2026-07-26 | 2026-07-26 | Cash flow, zero-based budget, Project Infrastructure Cost Attribution USP |
| 12 | Tech Skill Roadmap | ✅ | 2026-07-26 | 2026-07-26 | Interactive skill tree, study sessions, Mastery XP USP |

---

## Phase 0 — Scaffolding

| Task | Status | Notes |
|------|--------|-------|
| `go.mod` initialized | ✅ | `module github.com/Suke2004/atlas-go` |
| `cmd/atlas/main.go` created | ✅ | Entry point, DI wiring |
| `internal/config/config.go` | ✅ | Viper config struct |
| `Makefile` | ✅ | dev, build, migrate, sqlc, test, lint |
| `.air.toml` | ✅ | Hot-reload config |
| `sqlc.yaml` | ✅ | SQLite engine, queries → internal/db |
| `.env.example` | ✅ | All env vars documented |
| `/data` directory structure | ✅ | db/, uploads/, backups/, logs/ |
| `docker/Dockerfile` | ✅ | Two-stage build |
| `docker/docker-compose.yml` | ✅ | app + optional ollama profile |
| Health endpoint `GET /health` | ✅ | Returns status, version, uptime |
| `.github/workflows/ci.yml` | ✅ | Tests + build check |
| `.github/workflows/lint.yml` | ✅ | golangci-lint |
| `.gitignore` | ✅ | Go + data + env |
| `LICENSE` (MIT) | ✅ | |
| `README.md` skeleton | ✅ | Open-source structure |
| `CHANGELOG.md` | ✅ | v0.1.0 entry |

---

## Phase 1 — Database Foundation

| Task | Status | Notes |
|------|--------|-------|
| `001_users.sql` | ✅ | users + sessions tables |
| `002_settings.sql` | ✅ | key-value settings table |
| `003_projects.sql` | ✅ | projects + milestones |
| `004_tasks.sql` | ✅ | tasks + labels + dependencies |
| `005_notes.sql` | ✅ | notes + tags + links |
| `006_journal.sql` | ✅ | journal_entries + journal_items |
| `007_finance.sql` | ✅ | transactions + budgets |
| `008_learning.sql` | ✅ | learning_tracks + sessions |
| `009_fts5.sql` | ✅ | FTS5 virtual table + all triggers |
| `internal/db/db.go` | ✅ | SQLite connection + WAL pragmas |
| sqlc query files (`queries/*.sql`) | ✅ | All CRUD per module |
| `make sqlc` runs clean | ✅ | `internal/db/` generated |
| `make migrate` applies cleanly | ✅ | All 9 migrations |
| `go build ./...` passes | ✅ | No compile errors |

---

## Phase 2 — Auth + First-Run Wizard

| Task | Status | Notes |
|------|--------|-------|
| `internal/setup/service.go` | ✅ | IsFirstRun, CreateFirstUser, SeedDemoData |
| `internal/setup/middleware.go` | ✅ | FirstRunGate redirect |
| `internal/setup/handler.go` | ✅ | /setup, /setup/demo-choice, /setup/seed |
| `templates/setup/wizard.templ` | ✅ | Multi-step first-run form |
| `templates/setup/demo_choice.templ` | ✅ | Yes/No demo data |
| `internal/auth/service.go` | ✅ | Authenticate, CreateSession, ValidateSession |
| `internal/auth/middleware.go` | ✅ | AuthRequired, RedirectIfAuthenticated |
| `internal/auth/handler.go` | ✅ | GET/POST /login, POST /logout |
| `templates/auth/login.templ` | ✅ | Login form |
| Session cookie (HTTP-only, SameSite=Lax) | ✅ | |
| CSRF middleware wired | ✅ | All POST/PUT/DELETE protected |
| Rate limiting on /login | ✅ | ≤5 attempts/minute |
| Login e2e test | ✅ | |

---

## Phase 3 — Layout Shell

| Task | Status | Notes |
|------|--------|-------|
| `templates/layout/base.templ` | ✅ | HTML shell with data-theme attr |
| `templates/layout/sidebar.templ` | ✅ | Nav links + active state |
| `templates/layout/topbar.templ` | ✅ | Title + search trigger |
| `static/css/app.css` | ✅ | CSS tokens, light/dark/system theme |
| `static/js/htmx.min.js` | ✅ | Local copy, no CDN |
| `static/js/alpine.min.js` | ✅ | Local copy, no CDN |
| HTMX boost on all `<a>` links | ✅ | SPA-like nav |
| Theme: Light / Dark / System | ✅ | CSS custom properties |
| Theme persists across page reloads | ✅ | |

---

## Phase 4 — Projects Module

| Task | Status | Notes |
|------|--------|-------|
| `queries/projects.sql` | ✅ | All CRUD + milestones + GitHub stats |
| `internal/projects/repository.go` | ✅ | sqlc wrapper |
| `internal/projects/service.go` | ✅ | Business logic, RecalculateProgress, GitHub sync |
| `internal/projects/handler.go` | ✅ | All routes & HTMX swap handlers |
| `templates/projects/list.templ` | ✅ | Card grid with Tech Stack badges & filter tabs |
| `templates/projects/detail.templ` | ✅ | Full detail page with GitHub Insights Card |
| `templates/projects/github_card.templ` | ✅ | HTMX GitHub metrics partial component |
| Milestone completion → progress update | ✅ | Auto recalculate % |
| Project CRUD & GitHub integration test | ✅ | |

---

## Phase 5 — Tasks Module

| Task | Status | Notes |
|------|--------|-------|
| queries/tasks.sql | ✅ | CRUD + deps + labels + today focus |
| internal/tasks/repository.go | ✅ | |
| internal/tasks/service.go | ✅ | Dep validation on status change |
| internal/tasks/handler.go | ✅ | All routes incl. /status |
| templates/tasks/list.templ | ✅ | Table + filters |
| templates/tasks/kanban.templ | ✅ | Three-column board |
| templates/tasks/form.templ | ✅ | Full task form |
| templates/tasks/deps.templ | ✅ | Dependency list partial |
| static/js/kanban.js | ✅ | Alpine.js drag + hx-put |
| Task updates → project progress recalc | ✅ | |
| Overdue tasks highlighted | ✅ | |
| Tasks CRUD integration test | ✅ | |

---

## Phase 6 — Knowledge Base

| Task | Status | Notes |
|------|--------|-------|
| queries/notes.sql | ⬜ | CRUD + tags + links + backlinks |
| internal/notes/repository.go | ⬜ | |
| internal/notes/service.go | ⬜ | |
| internal/notes/handler.go | ⬜ | All routes + /preview |
| 	emplates/notes/list.templ | ⬜ | Grid + tag filter |
| 	emplates/notes/editor.templ | ⬜ | Split-pane editor |
| 	emplates/notes/backlinks.templ | ⬜ | Backlinks panel |
| static/js/editor.js | ⬜ | CodeMirror 6 + autosave |
| POST /notes/preview endpoint | ⬜ | Markdown → HTML |
| 30s autosave fires without user action | ⬜ | |
| Bidirectional note links work | ⬜ | |

---

## Phase 7 — Dashboard

| Task | Status | Notes |
|------|--------|-------|
| internal/dashboard/service.go | ⬜ | Aggregates from all modules |
| internal/dashboard/handler.go | ⬜ | Page + widget endpoints |
| 	emplates/dashboard/page.templ | ⬜ | Grid with widget slots |
| widgets/focus.templ | ⬜ | Top 3 tasks |
| widgets/projects.templ | ⬜ | Progress bars |
| widgets/journal_reminder.templ | ⬜ | CTA if no entry today |
| widgets/quote.templ | ⬜ | Random quote |
| Each widget loads independently | ⬜ | One failure → others unaffected |
| Dashboard loads < 1 second | ⬜ | |

---

## Phase 8 — Journal Module

| Task | Status | Notes |
|------|--------|-------|
| queries/journal.sql | ⬜ | Upsert + items + mood data |
| internal/journal/repository.go | ⬜ | |
| internal/journal/service.go | ⬜ | |
| internal/journal/handler.go | ⬜ | |
| 	emplates/journal/calendar.templ | ⬜ | Monthly calendar with dots |
| 	emplates/journal/entry.templ | ⬜ | Entry form (on-blur save) |
| 	emplates/journal/trend.templ | ⬜ | 7-day sparkline (Chart.js) |
| On-blur save for all fields | ⬜ | |
| "On this day" partial | ⬜ | |

---

## Phase 9 — Global Search

| Task | Status | Notes |
|------|--------|-------|
| internal/search/service.go | ⬜ | FTS5 across all modules |
| internal/search/handler.go | ⬜ | GET /search?q= |
| 	emplates/search/modal.templ | ⬜ | Alpine.js modal |
| 	emplates/search/results.templ | ⬜ | Grouped results |
| Ctrl+K opens modal from any page | ⬜ | |
| Results appear < 100ms | ⬜ | |
| Keyboard navigation (arrow + enter) | ⬜ | |
| Escape closes modal | ⬜ | |

---

## Phase 10 — Settings + Polish

| Task | Status | Notes |
|------|--------|-------|
| internal/settings/handler.go | ⬜ | GET/PUT /settings |
| internal/settings/service.go | ⬜ | Key-value CRUD + theme |
| 	emplates/settings/page.templ | ⬜ | All sections |
| Theme radio (Light/Dark/System) | ⬜ | Persists via settings table |
| Skeleton loaders on all widgets | ⬜ | |
| Empty states on all lists | ⬜ | |
| Toast notifications (success/error) | ⬜ | Alpine.js listener |
| HTMX indicators on all mutations | ⬜ | |
| docker compose up on fresh machine | ⬜ | Full e2e flow works |

---

## Test Coverage

| Layer | Current | Target |
|-------|---------|--------|
| Services | 0% | ≥ 80% |
| Repositories | 0% | ≥ 70% |
| Handlers | 0% | ≥ 60% |

---

## Commands Quick Reference

`ash
make dev            # Start dev server with hot-reload (Air)
make build          # Generate templates + build binary
make migrate        # Apply all pending migrations
make migrate-down   # Roll back last migration
make sqlc           # Regenerate DB code from queries
make templ          # Regenerate Templ templates
make test           # Run all tests with race detector
make lint           # Run golangci-lint
make docker-up      # Start Docker containers
make docker-down    # Stop Docker containers
make setup-dirs     # Create /data directory structure
`

---

## Environment

`ash
# Development
cp .env.example .env
# Fill in ATLAS_SESSION_SECRET (32 random bytes)
# Set ATLAS_AI_PROVIDER=ollama for local AI
`

---

## Known Issues / Blockers

*None yet — project not started.*

---

## Completed Milestones

*None yet.*

---

## Milestone Ceremony Checklist

Before marking any phase complete and moving to the next one:

```
□ All tasks in this phase are checked ✅ in progress.md
□ Started/Completed dates filled in the phase tracker table
□ CHANGELOG.md updated with what shipped (Keep a Changelog format)
□ docs/{module}/README.md written for every new module in this phase
□ docs/{module}/architecture.md written (data flow, layer diagram)
□ docs/{module}/database.md written (schema + index decisions)
□ docs/{module}/routes.md written (full route table)
□ ROADMAP.md milestone marked [x] with actual ship date
□ Git tag created: git tag -a vX.Y.Z -m "message" && git push origin vX.Y.Z
□ GitHub Release created (auto-triggered by release workflow on tag push)
□ docs/retros/phase-{N}.md written (what went well, harder than expected, decisions, time)
```

Only after every box is checked is the milestone truly complete.

---

## Phase Retrospectives

| Phase | File | Status |
|-------|------|--------|
| 0 — Scaffolding | [docs/retros/phase-0.md](docs/retros/phase-0.md) | ✅ Written |
| 1 — Database | [docs/retros/phase-1.md](docs/retros/phase-1.md) | ✅ Written |
| 2 — Auth + Wizard | [docs/retros/phase-2.md](docs/retros/phase-2.md) | ✅ Written |
| 3 — Layout Shell | [docs/retros/phase-3.md](docs/retros/phase-3.md) | ✅ Written |
| 4 — Projects | [docs/retros/phase-4.md](docs/retros/phase-4.md) | ✅ Written |
| 5 — Tasks | [docs/retros/phase-5.md](docs/retros/phase-5.md) | ✅ Written |
| 6 — Knowledge Base | [docs/retros/phase-6.md](docs/retros/phase-6.md) | ✅ Written |
| 7 — Dashboard | [docs/retros/phase-7.md](docs/retros/phase-7.md) | ✅ Written |
| 8 — Journal | [docs/retros/phase-8-9.md](docs/retros/phase-8-9.md) | ✅ Written |
| 9 — Search | [docs/retros/phase-9.md](docs/retros/phase-9.md) | ✅ Written |
| 10 — Settings | [docs/retros/phase-10.md](docs/retros/phase-10.md) | ✅ Written |

---

## Known Issues / Blockers

*None yet — project not started.*

---

## Notes

- Dashboard (Phase 7) intentionally placed after Projects (4) + Tasks (5) + Notes (6) to avoid mock data
- All integrations (GitHub, Weather, Gmail, LeetCode) are v2 — do not start them in v1
- Finance, Documents, Learning, AI Workspace are v2 — stubs only if needed
- AI provider interface (`internal/ai/provider.go`) should be defined early even though AI Workspace is v2
