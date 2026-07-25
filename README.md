<div align="center">

# Atlas

**Your personal operating system.**

A self-hosted, server-rendered application that replaces your notes app, task manager, journal, AI assistant, finance tracker, and knowledge base with a single, permanently open tab.

[![Go Version](https://img.shields.io/badge/go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue?style=flat)](LICENSE)
[![Build](https://github.com/Suke2004/atlas-go/actions/workflows/ci.yml/badge.svg)](https://github.com/Suke2004/atlas-go/actions)
[![Release](https://img.shields.io/github/v/release/Suke2004/atlas-go?style=flat)](https://github.com/Suke2004/atlas-go/releases)

</div>

---

> **Atlas is a self-hosted personal operating system designed to consolidate the tools, knowledge, and workflows of daily life into a single, server-rendered application. Rather than functioning as a collection of isolated productivity tools, Atlas models relationships between projects, tasks, notes, documents, events, and AI-assisted knowledge to create a cohesive workspace. Built with Go and HTMX, it emphasizes simplicity, performance, maintainability, and long-term ownership over frontend complexity.**

---

## Why Atlas?

The average developer uses 12–15 separate tools daily: tasks in Todoist, notes in Obsidian, journal in Notion, AI in ChatGPT, finance in spreadsheets. Each context switch costs 23 minutes of focus. None of them talk to each other.

Atlas is one tab that replaces all of them. Everything is connected. Everything is local. Everything is yours.

---

## Features

### v1.0 (Foundation)
- **First-run Setup Wizard** — guided onboarding with optional demo data
- **Projects** — full lifecycle tracking with milestones, timeline, and auto-progress
- **Tasks** — List and Kanban views, task dependencies, priority, energy level, labels
- **Knowledge Base** — Markdown editor with live preview, bidirectional note linking, 30s autosave
- **Journal** — daily entries, mood/energy/sleep tracking, trend charts
- **Global Search** — Ctrl+K across all modules, SQLite FTS5, < 100ms results
- **Dashboard** — personalised morning briefing with live widgets
- **Themes** — Light / Dark / System preference
- **Local-first** — SQLite database, all data stays on your machine

### v2.0 (Planned)
- **AI Workspace** — chat, PDF Q&A, "Ask Atlas" natural language queries over your data
- **Finance** — income/expense tracking, budgets, charts, AI insights
- **Learning Tracker** — streaks, progress, session logging
- **Documents** — file upload, AI summaries, full-text search
- **GitHub Integration** — commit feed, repository linking

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Language | Go 1.23+ |
| Router | Chi v5 |
| Templates | Templ (typesafe, compiled) |
| Frontend | HTMX 2.x + Alpine.js 3.x |
| Styling | TailwindCSS 3.x |
| Database | SQLite (WAL mode) |
| Search | SQLite FTS5 |
| Queries | sqlc (type-safe, generated) |
| Migrations | Goose v3 |
| Logging | Zap |
| Dev server | Air (hot-reload) |

No SPA framework. No virtual DOM. No Node.js in production. One binary, one database file.

---

## Quick Start

### Prerequisites

- Go 1.23+
- Docker + Docker Compose (for containerised deployment)

### With Docker (recommended)

```bash
git clone https://github.com/Suke2004/atlas-go
cd atlas-go

cp .env.example .env
# Edit .env — set ATLAS_SESSION_SECRET to a random 32-byte hex string

docker compose up -d
```

Open http://localhost:8080 — the setup wizard will appear.

### From Source

```bash
git clone https://github.com/Suke2004/atlas-go
cd atlas-go

# Install tools
go install github.com/a-h/templ/cmd/templ@latest
go install github.com/pressly/goose/v3/cmd/goose@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/air-verse/air@latest

# Setup
cp .env.example .env
make setup-dirs   # Creates /data/db, /data/uploads, etc.
make migrate      # Apply database migrations
make dev          # Start with hot-reload
```

Open http://localhost:8080

---

## Configuration

Copy `.env.example` to `.env` and fill in the values:

```env
ATLAS_PORT=8080
ATLAS_ENV=development
ATLAS_DB_PATH=./data/db/atlas.db
ATLAS_SESSION_SECRET=<32-byte-random-hex>
ATLAS_THEME_DEFAULT=system
```

See `docs/configuration.md` for the full reference.

---

## Project Structure

```
atlas/
├── cmd/atlas/          Entry point
├── internal/           Application code
│   ├── auth/           Authentication
│   ├── setup/          First-run wizard
│   ├── dashboard/      Dashboard aggregation
│   ├── projects/       Projects module
│   ├── tasks/          Tasks module
│   ├── notes/          Knowledge Base
│   ├── journal/        Journal module
│   ├── search/         Global search (FTS5)
│   └── settings/       User settings
├── web/
│   ├── templates/      Templ templates
│   └── static/         CSS, JS, icons
├── migrations/         Database migrations
├── queries/            sqlc SQL queries
└── docs/               Documentation
```

---

## Development

```bash
make dev          # Hot-reload dev server
make test         # Run all tests
make lint         # Lint check
make migrate      # Apply migrations
make sqlc         # Regenerate DB code
make build        # Production build
```

See `guidelines.md` for the complete engineering guide.

---

## Roadmap

| Release | Focus | Status |
|---------|-------|--------|
| v0.1.0 | Scaffold + Auth | 🔄 In progress |
| v0.2.0 | Projects + Tasks | ⬜ Planned |
| v0.3.0 | Knowledge + Journal + Search | ⬜ Planned |
| v0.4.0 | Dashboard + Settings | ⬜ Planned |
| v0.5.0 | AI + Finance + Learning | ⬜ Planned |
| v1.0.0 | Stable release | ⬜ Planned |

Full roadmap: [`ROADMAP.md`](ROADMAP.md) · [`future_plans.md`](future_plans.md)

---

## Documentation

| Document | Description |
|----------|-------------|
| [`docs/prd.md`](docs/prd.md) | Product Requirements |
| [`docs/trd.md`](docs/trd.md) | Technical Requirements |
| [`ARCHITECTURE.md`](ARCHITECTURE.md) | System architecture |
| [`guidelines.md`](guidelines.md) | Engineering guidelines |
| [`progress.md`](progress.md) | Implementation progress |
| [`future_plans.md`](future_plans.md) | Full feature roadmap |
| [`docs/getting-started.md`](docs/getting-started.md) | Setup guide |
| [`docs/configuration.md`](docs/configuration.md) | Configuration reference |

---

## Contributing

Atlas is open to contributions. Read [`CONTRIBUTING.md`](CONTRIBUTING.md) and [`guidelines.md`](guidelines.md) before starting.

Key rules:
- Follow the Handler → Service → Repository → DB layering
- Use Conventional Commits (`feat(tasks): add dependency validation`)
- Target the `develop` branch — never commit to `main`
- All queries through sqlc — no raw SQL string building

---

## License

[MIT](LICENSE) — free to use, modify, and self-host.

---

<div align="center">

Built with Go and HTMX · Self-hosted · Local-first · Single-user

</div>