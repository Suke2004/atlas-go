# Atlas — Technical Requirements Document (TRD)

> **Version**: 1.1 (Decisions Locked)
> **Date**: 2026-07-26
> **Status**: Approved

---

## 1. System Overview

Atlas is a self-hosted, single-user web application built with Go on the backend and HTMX-powered server-side rendering on the frontend. The architecture is deliberately minimal: one binary, one database file, one Docker container.

### 1.1 Design Philosophy

| Principle | Implementation |
|-----------|----------------|
| **Server-side simplicity** | HTML rendered by Templ templates on the server; no SPA framework |
| **Partial updates** | HTMX swaps DOM fragments — no full-page reloads |
| **Single database** | SQLite for all relational data + FTS5 for search |
| **Zero build complexity** | Go binary + static assets; Air for dev hot-reload |
| **Local-first** | All data stays on the host machine |

---

## 2. Tech Stack

### 2.1 Backend

| Technology | Version | Purpose |
|------------|---------|---------|
| Go | 1.23+ | Application language |
| Chi | v5 | HTTP router (lightweight, idiomatic) |
| Templ | latest | Typesafe Go HTML templates |
| SQLite (mattn/go-sqlite3) | latest | Primary database |
| sqlc | v1.x | SQL → Go code generation |
| Goose | v3 | Database migrations |
| Zap | v1.x | Structured logging |
| Viper | v2 | Configuration management |
| bcrypt | stdlib | Password hashing |
| sessions (gorilla/sessions) | v1 | Session management |
| Air | latest | Hot-reload in development |

### 2.2 Frontend

| Technology | Version | Purpose |
|------------|---------|---------|
| HTMX | 2.x | Partial HTML updates over HTTP |
| TailwindCSS | 3.x (CDN or CLI) | Utility-first styling |
| Alpine.js | 3.x | Minimal reactive UI state (modals, dropdowns) |
| Heroicons | 2.x | SVG icon library |
| Chart.js | 4.x | Data visualisation |
| CodeMirror | 6.x | Markdown editor with syntax highlighting |

### 2.3 Infrastructure

| Technology | Purpose |
|------------|---------|
| Docker | Container packaging |
| Docker Compose | Multi-service orchestration (app + optional Ollama) |
| Caddy / Nginx | Reverse proxy + HTTPS (production) |
| GitHub Actions | CI/CD pipeline |
| Makefile | Developer convenience commands |
| Air | Development hot-reload |

### 2.4 AI / Search

| Technology | Purpose |
|------------|---------|
| AI Provider Interface | Go interface; never coupled to one SDK |
| OpenAI API (v2) | Primary cloud AI backend |
| Ollama (v2) | Local LLM — Llama 3, Mistral, etc. |
| Anthropic / Gemini (v2+) | Additional providers via same interface |
| SQLite FTS5 | Full-text search across all modules (v1) |
| Vector search (v3) | Semantic embeddings; pgvector or sqlite-vec |

---

## 3. Architecture

### 3.1 High-Level Architecture

```
Browser
  │
  │  HTTP (HTML fragments via HTMX)
  ▼
┌──────────────────────────────────────────┐
│              Atlas Go Server             │
│                                          │
│  Chi Router                              │
│    │                                     │
│    ├── Middleware Chain                  │
│    │    ├── Logger (Zap)                 │
│    │    ├── Auth (session check)         │
│    │    ├── CSRF                         │
│    │    └── Request ID                   │
│    │                                     │
│    ├── Handlers (per module)             │
│    │    └── → Services → Repositories   │
│    │                                     │
│    └── Templ Template Renderer           │
│                                          │
│  SQLite Database                         │
│    ├── Application data (all modules)    │
│    └── FTS5 virtual tables (search)      │
│                                          │
│  Local Filesystem                        │
│    ├── /data/db/        ← SQLite database file
│    ├── /data/uploads/   ← Document files (v2)
│    ├── /data/backups/   ← Automated backup copies
│    └── /data/logs/      ← File logs (if enabled)
└──────────────────────────────────────────┘
  │
  │  HTTP (OpenAI-compatible API)
  ▼
┌──────────────────────────────────────────┐
│         AI Backend (optional)            │
│  OpenAI API  OR  Ollama (local)          │
└──────────────────────────────────────────┘
```

### 3.2 Request Lifecycle

```
1. Browser sends HTTP request (GET/POST/PUT/DELETE)
2. Chi router matches route
3. Middleware chain runs (logger → auth → CSRF)
4. Handler called
5. Handler calls Service layer (business logic)
6. Service calls Repository layer (SQL via sqlc)
7. Repository executes query against SQLite
8. Data flows back: Repository → Service → Handler
9. Handler renders Templ template → HTML fragment
10. HTMX swaps target element in DOM
```

### 3.3 HTMX Pattern

Every interactive element uses HTMX attributes:
```html
<!-- Example: delete a task -->
<button
  hx-delete="/tasks/42"
  hx-confirm="Delete this task?"
  hx-target="#task-42"
  hx-swap="outerHTML"
>
  Delete
</button>
```

No JavaScript required. Server returns empty response or replacement HTML.

---

## 4. Directory Structure

```
atlas/
│
├── cmd/
│   └── server/
│       └── main.go              # Entry point, dependency wiring
│
├── internal/
│   ├── setup/
│   │   ├── handler.go           # First-run wizard HTTP handlers
│   │   ├── service.go           # IsFirstRun, CreateFirstUser, SeedDemoData
│   │   └── middleware.go        # FirstRunGate — redirect to /setup if no users
│   │
│   ├── auth/
│   │   ├── handler.go           # Login/logout HTTP handlers
│   │   ├── service.go           # Auth business logic
│   │   ├── middleware.go        # Session auth middleware
│   │   └── session.go           # Session store wrapper
│   │
│   ├── dashboard/
│   │   ├── handler.go
│   │   └── service.go           # Aggregates data from all modules
│   │
│   ├── tasks/
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── repository.go        # sqlc-generated queries wrapper
│   │
│   ├── projects/
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── repository.go
│   │
│   ├── notes/                   # Knowledge Base
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── repository.go
│   │
│   ├── journal/
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── repository.go
│   │
│   ├── search/
│   │   ├── handler.go
│   │   └── service.go           # FTS5 query across all modules
│   │
│   ├── settings/
│   │   ├── handler.go
│   │   └── service.go           # Key-value settings + theme management
│   │
│   ├── ai/                      # v2 — provider interface defined now
│   │   ├── provider.go          # Provider interface
│   │   ├── openai/provider.go   # OpenAI implementation
│   │   ├── ollama/provider.go   # Ollama implementation
│   │   └── service.go           # Chat, RAG, embeddings
│   │
│   ├── finance/                 # v2
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── repository.go
│   │
│   ├── documents/               # v2
│   │   ├── handler.go
│   │   ├── service.go           # Upload, OCR trigger, AI summary
│   │   └── repository.go
│   │
│   ├── learning/                # v2
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── repository.go
│   │
│   ├── analytics/               # v3
│   │   ├── handler.go
│   │   └── service.go
│   │
│   ├── notifications/           # v2
│   │   ├── handler.go
│   │   └── scheduler.go         # Background goroutine for checks
│   │
│   ├── db/
│   │   ├── db.go                # SQLite connection + pragmas
│   │   ├── query.sql.go         # sqlc generated
│   │   └── models.go            # sqlc generated
│   │
│   ├── config/
│   │   └── config.go            # Viper config struct
│   │
│   └── middleware/
│       ├── logger.go
│       ├── requestid.go
│       └── csrf.go
│
├── templates/
│   ├── layout/
│   │   ├── base.templ            # Root HTML shell, nav sidebar, theme attr
│   │   └── head.templ
│   ├── setup/
│   │   ├── wizard.templ          # Multi-step first-run wizard
│   │   └── demo_choice.templ     # Demo data offer page
│   ├── dashboard/
│   │   ├── page.templ
│   │   └── widgets/
│   │       ├── focus.templ
│   │       ├── projects.templ
│   │       ├── expenses.templ     # v2 (Finance module)
│   │       ├── streak.templ       # v2 (Learning module)
│   │       ├── journal_reminder.templ
│   │       └── quote.templ
│   ├── tasks/
│   │   ├── list.templ
│   │   ├── kanban.templ
│   │   ├── form.templ
│   │   └── deps.templ             # Task dependency list partial
│   ├── projects/
│   ├── notes/
│   ├── journal/
│   ├── search/
│   ├── settings/
│   └── auth/
│       └── login.templ
│
├── static/
│   ├── css/
│   │   └── app.css              # TailwindCSS output
│   ├── js/
│   │   ├── htmx.min.js
│   │   ├── alpine.min.js
│   │   └── charts.js            # Chart.js wrapper
│   └── icons/
│
├── migrations/
│   ├── 001_init.sql
│   ├── 002_tasks.sql
│   ├── 003_projects.sql
│   ├── 004_notes.sql
│   ├── 005_journal.sql
│   ├── 006_finance.sql
│   ├── 007_documents.sql
│   ├── 008_learning.sql
│   └── 009_fts5.sql
│
├── queries/                     # sqlc source SQL queries
│   ├── tasks.sql
│   ├── projects.sql
│   ├── notes.sql
│   ├── journal.sql
│   ├── finance.sql
│   ├── learning.sql
│   └── search.sql
│
├── scripts/
│   └── seed.go                  # Demo data seed (called by setup wizard)
│
├── docker/
│   ├── Dockerfile
│   └── docker-compose.yml
│
├── docs/
│   ├── prd.md
│   └── trd.md
│
├── tests/
│   ├── integration/
│   └── unit/
│
├── .air.toml                    # Air hot-reload config
├── .env.example
├── sqlc.yaml
├── Makefile
├── go.mod
└── go.sum
```

---

## 5. Database Schema

### 5.1 SQLite Configuration

```sql
PRAGMA journal_mode = WAL;       -- Write-Ahead Logging for concurrent reads
PRAGMA foreign_keys = ON;        -- Enforce FK constraints
PRAGMA synchronous = NORMAL;     -- Balance durability vs speed
PRAGMA cache_size = -64000;      -- 64MB cache
PRAGMA temp_store = MEMORY;      -- Temp tables in memory
```

### 5.2 Core Tables

```sql
-- Users (single-user, but extensible)
CREATE TABLE users (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    username    TEXT NOT NULL UNIQUE,
    email       TEXT NOT NULL UNIQUE,
    password    TEXT NOT NULL,           -- bcrypt hash
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Sessions
CREATE TABLE sessions (
    id          TEXT PRIMARY KEY,
    user_id     INTEGER NOT NULL REFERENCES users(id),
    expires_at  DATETIME NOT NULL,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Projects
CREATE TABLE projects (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id         INTEGER NOT NULL REFERENCES users(id),
    name            TEXT NOT NULL,
    description     TEXT,
    status          TEXT NOT NULL DEFAULT 'active'
                        CHECK(status IN ('active','paused','completed','archived')),
    priority        TEXT NOT NULL DEFAULT 'medium'
                        CHECK(priority IN ('critical','high','medium','low')),
    deadline        DATE,
    progress        INTEGER NOT NULL DEFAULT 0 CHECK(progress BETWEEN 0 AND 100),
    repository_url  TEXT,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Milestones
CREATE TABLE milestones (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id  INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    title       TEXT NOT NULL,
    due_date    DATE,
    completed   BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Tasks
CREATE TABLE tasks (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id             INTEGER NOT NULL REFERENCES users(id),
    project_id          INTEGER REFERENCES projects(id) ON DELETE SET NULL,
    title               TEXT NOT NULL,
    description         TEXT,
    status              TEXT NOT NULL DEFAULT 'todo'
                            CHECK(status IN ('todo','in_progress','done','cancelled')),
    priority            TEXT NOT NULL DEFAULT 'medium'
                            CHECK(priority IN ('critical','high','medium','low')),
    energy_required     TEXT NOT NULL DEFAULT 'shallow'
                            CHECK(energy_required IN ('deep','shallow','admin')),
    estimated_minutes   INTEGER,
    actual_minutes      INTEGER,
    deadline            DATETIME,
    recurrence          TEXT DEFAULT 'none'
                            CHECK(recurrence IN ('daily','weekly','monthly','none')),
    created_at          DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Task Labels
CREATE TABLE task_labels (
    task_id     INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    label       TEXT NOT NULL,
    PRIMARY KEY (task_id, label)
);

-- Task Dependencies
CREATE TABLE task_dependencies (
    task_id         INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    depends_on_id   INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    PRIMARY KEY (task_id, depends_on_id)
);

-- Notes (Knowledge Base)
CREATE TABLE notes (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id     INTEGER NOT NULL REFERENCES users(id),
    project_id  INTEGER REFERENCES projects(id) ON DELETE SET NULL,
    title       TEXT NOT NULL,
    content     TEXT,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Note Tags
CREATE TABLE note_tags (
    note_id     INTEGER NOT NULL REFERENCES notes(id) ON DELETE CASCADE,
    tag         TEXT NOT NULL,
    PRIMARY KEY (note_id, tag)
);

-- Note Links (bidirectional via two rows)
CREATE TABLE note_links (
    source_id   INTEGER NOT NULL REFERENCES notes(id) ON DELETE CASCADE,
    target_id   INTEGER NOT NULL REFERENCES notes(id) ON DELETE CASCADE,
    PRIMARY KEY (source_id, target_id)
);

-- Journal Entries
CREATE TABLE journal_entries (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id         INTEGER NOT NULL REFERENCES users(id),
    date            DATE NOT NULL,
    mood            INTEGER CHECK(mood BETWEEN 1 AND 5),
    energy          INTEGER CHECK(energy BETWEEN 1 AND 5),
    sleep_hours     REAL,
    content         TEXT,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, date)
);

-- Journal Sections (wins, problems, ideas, tomorrow)
CREATE TABLE journal_items (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    entry_id        INTEGER NOT NULL REFERENCES journal_entries(id) ON DELETE CASCADE,
    type            TEXT NOT NULL CHECK(type IN ('win','problem','idea','tomorrow')),
    content         TEXT NOT NULL
);

-- Finance Transactions
CREATE TABLE transactions (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id     INTEGER NOT NULL REFERENCES users(id),
    amount      REAL NOT NULL,
    type        TEXT NOT NULL CHECK(type IN ('income','expense')),
    category    TEXT NOT NULL,
    description TEXT,
    date        DATE NOT NULL DEFAULT (DATE('now')),
    recurring   BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Finance Budgets
CREATE TABLE budgets (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id     INTEGER NOT NULL REFERENCES users(id),
    category    TEXT NOT NULL,
    amount      REAL NOT NULL,
    period      TEXT NOT NULL DEFAULT 'monthly',
    UNIQUE(user_id, category)
);

-- Documents
CREATE TABLE documents (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id         INTEGER NOT NULL REFERENCES users(id),
    filename        TEXT NOT NULL,
    original_name   TEXT NOT NULL,
    mime_type       TEXT NOT NULL,
    file_size       INTEGER NOT NULL,
    file_path       TEXT NOT NULL,           -- relative to /data/uploads/
    ocr_content     TEXT,                    -- extracted text
    ai_summary      TEXT,                    -- AI-generated summary
    ocr_status      TEXT DEFAULT 'pending'
                        CHECK(ocr_status IN ('pending','processing','done','failed')),
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Document Tags
CREATE TABLE document_tags (
    document_id     INTEGER NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    tag             TEXT NOT NULL,
    PRIMARY KEY (document_id, tag)
);

-- Learning Tracks
CREATE TABLE learning_tracks (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id         INTEGER NOT NULL REFERENCES users(id),
    name            TEXT NOT NULL,
    type            TEXT NOT NULL
                        CHECK(type IN ('dsa','course','book','paper','language','framework','other')),
    progress        INTEGER NOT NULL DEFAULT 0 CHECK(progress BETWEEN 0 AND 100),
    current_streak  INTEGER NOT NULL DEFAULT 0,
    longest_streak  INTEGER NOT NULL DEFAULT 0,
    total_sessions  INTEGER NOT NULL DEFAULT 0,
    last_session    DATE,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Learning Sessions
CREATE TABLE learning_sessions (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    track_id    INTEGER NOT NULL REFERENCES learning_tracks(id) ON DELETE CASCADE,
    date        DATE NOT NULL DEFAULT (DATE('now')),
    duration_minutes INTEGER,
    notes       TEXT,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Settings (key-value store)
CREATE TABLE settings (
    user_id     INTEGER NOT NULL REFERENCES users(id),
    key         TEXT NOT NULL,
    value       TEXT NOT NULL,
    PRIMARY KEY (user_id, key)
);
```

### 5.3 FTS5 Virtual Tables

```sql
-- Search index (synced via triggers)
CREATE VIRTUAL TABLE search_index USING fts5(
    entity_type,        -- 'task' | 'note' | 'project' | 'journal' | 'document' | 'transaction'
    entity_id,          -- FK to source table
    title,
    content,
    tags
);

-- Auto-sync trigger example (tasks):
CREATE TRIGGER tasks_fts_insert AFTER INSERT ON tasks BEGIN
    INSERT INTO search_index(entity_type, entity_id, title, content)
    VALUES ('task', NEW.id, NEW.title, NEW.description);
END;

CREATE TRIGGER tasks_fts_update AFTER UPDATE ON tasks BEGIN
    DELETE FROM search_index WHERE entity_id = OLD.id AND entity_type = 'task';
    INSERT INTO search_index(entity_type, entity_id, title, content)
    VALUES ('task', NEW.id, NEW.title, NEW.description);
END;

CREATE TRIGGER tasks_fts_delete AFTER DELETE ON tasks BEGIN
    DELETE FROM search_index WHERE entity_id = OLD.id AND entity_type = 'task';
END;
```

---

## 6. API Design

### 6.1 Route Conventions

All routes follow RESTful conventions with HTMX-friendly responses.

```
GET    /                       → Dashboard (full page or HTMX partial)
GET    /login                  → Login page
POST   /login                  → Authenticate
POST   /logout                 → Destroy session

GET    /tasks                  → Task list
GET    /tasks/new              → New task form (partial)
POST   /tasks                  → Create task
GET    /tasks/{id}/edit        → Edit form (partial)
PUT    /tasks/{id}             → Update task
DELETE /tasks/{id}             → Delete task
PUT    /tasks/{id}/status      → Toggle status (Kanban drag)

GET    /projects               → Project list
GET    /projects/new           → New project form
POST   /projects               → Create project
GET    /projects/{id}          → Project detail
PUT    /projects/{id}          → Update project
DELETE /projects/{id}          → Delete project

GET    /notes                  → Note list
GET    /notes/new              → New note form
POST   /notes                  → Create note
GET    /notes/{id}             → Note detail
PUT    /notes/{id}             → Update note
DELETE /notes/{id}             → Delete note

GET    /journal                → Journal calendar
GET    /journal/today          → Today's entry (or new)
PUT    /journal/{date}         → Save entry

GET    /finance                → Finance overview
POST   /finance/transactions   → Add transaction
DELETE /finance/transactions/{id} → Delete transaction

POST   /documents/upload       → Upload document
GET    /documents/{id}/preview → Preview document

GET    /learning               → Learning tracks
POST   /learning               → Create track
POST   /learning/{id}/session  → Log session

GET    /analytics              → Analytics dashboard

GET    /ai/chat                → AI workspace
POST   /ai/chat                → Send message (streaming SSE)

GET    /search                 → Search page
GET    /search?q={query}       → Search results (partial)

GET    /settings               → Settings page
PUT    /settings               → Update settings
```

### 6.2 HTMX Response Rules

| Trigger | Response type | HTMX target |
|---------|--------------|-------------|
| Form submit (create) | 201 + new item HTML | Prepend to list |
| Form submit (update) | 200 + updated item HTML | Replace item |
| Delete button | 200 empty or 204 | Remove item (outerHTML swap) |
| Error | 422 + error partial | Error container |
| Search | 200 + results partial | Results container |

---

## 7. Authentication & Security

### 7.1 Authentication Flow

```
1. POST /login {username, password}
2. Server: bcrypt.CompareHashAndPassword(stored_hash, password)
3. On success: create sessions row, set session cookie (HTTP-only, SameSite=Lax)
4. On all subsequent requests: auth middleware reads cookie → validates session → sets user in context
5. POST /logout: delete session row, clear cookie
```

### 7.2 Middleware Stack

```go
r.Use(middleware.RequestID)
r.Use(middleware.RealIP)
r.Use(zapLogger)
r.Use(middleware.Recoverer)
r.Use(csrfMiddleware)

r.Route("/", func(r chi.Router) {
    r.Use(authRequired)  // all authenticated routes
    // ... module routes
})
```

### 7.3 CSRF Protection

- Use `gorilla/csrf` or custom double-submit cookie pattern
- CSRF token embedded in all forms via Templ helper
- Token validated on all POST/PUT/DELETE requests

---

## 8. Configuration

### 8.1 Environment Variables (.env)

```env
# Server
ATLAS_PORT=8080
ATLAS_ENV=development              # development | production

# Database
ATLAS_DB_PATH=./data/atlas.db

# Session
ATLAS_SESSION_SECRET=<32-byte-random-hex>
ATLAS_SESSION_MAX_AGE=86400        # 24 hours

# File Storage
ATLAS_UPLOAD_PATH=./data/uploads
ATLAS_MAX_UPLOAD_MB=50

# AI
ATLAS_AI_PROVIDER=openai           # openai | ollama
ATLAS_AI_API_KEY=sk-...
ATLAS_AI_BASE_URL=https://api.openai.com/v1
ATLAS_AI_MODEL=gpt-4o

# Ollama (if provider=ollama)
ATLAS_OLLAMA_URL=http://localhost:11434
ATLAS_OLLAMA_MODEL=llama3.2

# External Integrations (v2)
ATLAS_GITHUB_TOKEN=ghp_...
ATLAS_WEATHER_API_KEY=...
```

### 8.2 sqlc Configuration (sqlc.yaml)

```yaml
version: "2"
sql:
  - engine: "sqlite"
    queries: "./queries/"
    schema: "./migrations/"
    gen:
      go:
        package: "db"
        out: "internal/db"
        emit_json_tags: true
        emit_prepared_queries: false
```

---

## 9. Development Workflow

### 9.1 Makefile Commands

```makefile
.PHONY: dev build test migrate seed lint tidy templ

dev:
	air

build:
	templ generate
	go build -o bin/atlas ./cmd/server

templ:
	templ generate

migrate:
	goose -dir migrations sqlite ./data/atlas.db up

migrate-down:
	goose -dir migrations sqlite ./data/atlas.db down

migrate-create:
	goose -dir migrations create $(name) sql

sqlc:
	sqlc generate

seed:
	go run ./scripts/seed.go

test:
	go test ./... -v

lint:
	golangci-lint run

tidy:
	go mod tidy

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-build:
	docker compose build
```

### 9.2 Air Configuration (.air.toml)

```toml
root = "."
tmp_dir = "tmp"

[build]
  cmd = "templ generate && go build -o ./tmp/main ./cmd/server"
  bin = "./tmp/main"
  include_ext = ["go", "templ", "sql"]
  exclude_dir = ["tmp", "vendor", "node_modules"]
  delay = 500

[log]
  time = true
```

### 9.3 Development Setup

```bash
# 1. Clone and enter repo
git clone <repo>
cd atlas

# 2. Install tools
go install github.com/a-h/templ/cmd/templ@latest
go install github.com/pressly/goose/v3/cmd/goose@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/air-verse/air@latest

# 3. Copy and fill env
cp .env.example .env

# 4. Run migrations
make migrate

# 5. Start dev server
make dev
```

---

## 10. Docker Deployment

### 10.1 Dockerfile

```dockerfile
# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/a-h/templ/cmd/templ@latest
RUN templ generate
RUN CGO_ENABLED=1 GOOS=linux go build -o atlas ./cmd/server

# Run stage
FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/atlas .
COPY --from=builder /app/static ./static
COPY --from=builder /app/migrations ./migrations

RUN mkdir -p /data/uploads

EXPOSE 8080

CMD ["./atlas"]
```

### 10.2 docker-compose.yml

```yaml
version: '3.9'

services:
  atlas:
    build:
      context: .
      dockerfile: docker/Dockerfile
    restart: unless-stopped
    ports:
      - "8080:8080"
    volumes:
      - atlas_data:/data
    env_file:
      - .env
    environment:
      ATLAS_DB_PATH: /data/atlas.db
      ATLAS_UPLOAD_PATH: /data/uploads

  # Optional: Local LLM
  ollama:
    image: ollama/ollama:latest
    restart: unless-stopped
    ports:
      - "11434:11434"
    volumes:
      - ollama_data:/root/.ollama
    profiles:
      - ollama   # Only starts with: docker compose --profile ollama up

volumes:
  atlas_data:
  ollama_data:
```

---

## 11. Testing Strategy

### 11.1 Unit Tests

- Service layer: pure function tests with mocked repositories
- Repository layer: tests against in-memory SQLite (`:memory:`)
- Handlers: `httptest` with mock services

### 11.2 Integration Tests

- Full request/response cycle tests
- Test DB seeded and torn down per test run
- Focus: auth flows, HTMX partial responses, search accuracy

### 11.3 Coverage Targets

| Layer | Target Coverage |
|-------|----------------|
| Services | ≥ 80% |
| Repositories | ≥ 70% |
| Handlers | ≥ 60% |

---

## 12. Performance Considerations

### 12.1 SQLite Optimisations

- WAL mode enabled at startup (concurrent readers + one writer)
- Indexes on all FK columns and frequently filtered columns
- FTS5 for all text search (avoids LIKE '%query%' scans)
- Connection pool: single write connection + multiple read connections via `database/sql`

### 12.2 HTMX Optimisations

- Dashboard widgets load in parallel via separate HTMX requests
- `hx-trigger="load"` on each widget — independent failure isolation
- Responses cached at handler level where data is stable (quotes, weather)
- `hx-boost` on standard `<a>` links for SPA-like navigation

### 12.3 Static Assets

- Tailwind CSS purged for production build (only used classes)
- HTMX and Alpine.js served from local static/ (no CDN dependency)
- `Cache-Control: max-age=31536000` on versioned static assets

---

## 13. CI/CD Pipeline (GitHub Actions)

```yaml
# .github/workflows/ci.yml
name: CI

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      - name: Install templ
        run: go install github.com/a-h/templ/cmd/templ@latest
      - name: Generate templates
        run: templ generate
      - name: Run tests
        run: go test ./... -race -coverprofile=coverage.out
      - name: Build
        run: CGO_ENABLED=1 go build ./cmd/server

  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: golangci/golangci-lint-action@v4
```

---

## 14. Security Checklist

- [ ] All user input sanitised before DB insertion (sqlc parameterised queries)
- [ ] No raw SQL string formatting with user data
- [ ] File upload: MIME type validation, max size enforced, random filename on disk
- [ ] Session tokens are cryptographically random (≥32 bytes)
- [ ] Passwords hashed with bcrypt (cost factor ≥ 12)
- [ ] CSRF token on all state-changing requests
- [ ] HTTP-only cookies (no JS access to session)
- [ ] SameSite=Lax cookie attribute
- [ ] Rate limiting on login endpoint (≤5 attempts per minute)
- [ ] Uploaded files served from non-web-accessible path or with strict MIME headers

---

## 15. Monitoring & Observability

### 15.1 Structured Logging (Zap)

Every request logged with:
- Request ID
- Method + path
- Duration (ms)
- Status code
- User ID (if authenticated)

### 15.2 Health Endpoint

```
GET /health

Response 200:
{
  "status": "ok",
  "db": "connected",
  "version": "1.0.0",
  "uptime": "4h32m"
}
```

---

*End of TRD v1.0*
