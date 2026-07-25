# Architecture

> A technical overview of Atlas's system design, layers, and request lifecycle.

---

## System Overview

Atlas is a **server-rendered web application**. All HTML is generated on the server by Templ templates. The browser receives HTML fragments via HTMX — there is no JavaScript SPA framework, no virtual DOM, no frontend build pipeline.

`
Browser
  │
  │  HTTP (GET/POST/PUT/DELETE)
  │  Response: HTML fragments via HTMX
  ▼
┌──────────────────────────────────────────┐
│              Atlas Go Server             │
│                                          │
│  Chi Router                              │
│    │                                     │
│    ├── Middleware Chain                  │
│    │    ├── RequestID                    │
│    │    ├── Logger (Zap)                 │
│    │    ├── Recoverer                    │
│    │    ├── CSRF                         │
│    │    ├── FirstRunGate                 │
│    │    └── AuthRequired                 │
│    │                                     │
│    └── Module Handlers                   │
│         └── Services → Repositories     │
│                                          │
│  Templ Template Renderer                 │
│  SQLite (WAL mode)                       │
│    ├── Relational tables                 │
│    └── FTS5 virtual tables               │
│                                          │
│  /data/                                  │
│    ├── db/atlas.db                       │
│    ├── uploads/                          │
│    └── backups/                          │
└──────────────────────────────────────────┘
  │
  │  HTTP (OpenAI-compatible API) — v2
  ▼
AI Backend: OpenAI / Ollama / Anthropic / Gemini
`

---

## Layered Architecture

Every request follows this **strict pipeline**:

`
Presentation Layer     templates/*.templ
        ↑
Handler Layer          internal/*/handler.go
        │
Service Layer          internal/*/service.go
        │
Repository Layer       internal/*/repository.go
        │
Database Layer         internal/db/  (sqlc generated)
        │
        ▼
SQLite Database        /data/db/atlas.db
`

### Layer Responsibilities

| Layer | File | Responsibility |
|-------|------|----------------|
| **Template** | *.templ | Render data → HTML. No logic. |
| **Handler** | handler.go | Parse HTTP request → call service → render template. Log errors. |
| **Service** | service.go | Business logic, validation, orchestration. |
| **Repository** | epository.go | DB access only. Thin wrapper around sqlc. |
| **DB** | internal/db/ | sqlc-generated type-safe query functions. |

### Rules

- Handlers never touch the DB directly.
- Repositories contain zero business logic.
- Services never render HTML or write to http.ResponseWriter.
- Modules do not import each other's packages — they communicate via service interfaces passed through constructors.

---

## Request Lifecycle

`
1. Browser sends HTTP request
2. Chi router matches URL → selects handler
3. Middleware chain runs:
   a. RequestID — generate/attach unique ID
   b. Logger — log method, path, remote addr
   c. Recoverer — catch panics
   d. CSRF — validate token on mutations
   e. FirstRunGate — redirect to /setup if no users
   f. AuthRequired — validate session cookie
4. Handler called:
   a. Parse and validate request (form values, path params)
   b. Call service method
   c. Handle error → return 422 + error partial
5. Service executes business logic:
   a. Call repository methods
   b. Coordinate across modules if needed
   c. Return domain struct or error
6. Repository executes sqlc-generated query against SQLite
7. Data flows back: Repository → Service → Handler
8. Handler renders Templ template → HTML fragment
9. HTTP response written with appropriate status code
10. HTMX swaps target DOM element in browser
`

---

## HTMX Pattern

Standard approach for all interactions:

`html
<!-- Create: response prepended to list -->
<form hx-post="/tasks" hx-target="#task-list" hx-swap="afterbegin"
      hx-on::after-request="this.reset()">

<!-- Update: response replaces the item -->
<form hx-put="/tasks/42" hx-target="#task-42" hx-swap="outerHTML">

<!-- Delete: item removed -->
<button hx-delete="/tasks/42" hx-target="#task-42" hx-swap="outerHTML"
        hx-confirm="Delete this task?">

<!-- Dashboard widget: loads independently on page ready -->
<div hx-get="/dashboard/widgets/focus" hx-trigger="load"
     hx-indicator="#focus-loader">

<!-- Search: fires after 200ms of input idle -->
<input hx-get="/search" hx-trigger="input delay:200ms"
       hx-target="#search-results">
`

---

## Database Design

### SQLite Configuration

`sql
PRAGMA journal_mode = WAL;       -- Concurrent reads + single writer
PRAGMA foreign_keys = ON;        -- FK enforcement
PRAGMA synchronous = NORMAL;     -- Durability/speed balance
PRAGMA cache_size = -64000;      -- 64MB page cache
PRAGMA temp_store = MEMORY;
PRAGMA busy_timeout = 5000;      -- 5s wait on write lock
`

### Module Tables

| Module | Tables |
|--------|--------|
| Auth | users, sessions |
| Settings | settings |
| Projects | projects, milestones |
| Tasks | 	asks, 	ask_labels, 	ask_dependencies |
| Notes | 
otes, 
ote_tags, 
ote_links |
| Journal | journal_entries, journal_items |
| Finance | 	ransactions, udgets |
| Learning | learning_tracks, learning_sessions |
| Documents | documents, document_tags |
| Search | search_index (FTS5 virtual table) |

### FTS5 Search Index

A single FTS5 virtual table indexes all searchable content:

`sql
CREATE VIRTUAL TABLE search_index USING fts5(
    entity_type,          -- 'task' | 'note' | 'project' | ...
    entity_id UNINDEXED,  -- FK to source table
    user_id UNINDEXED,    -- scoped per user
    title,
    content,
    tags
);
`

Kept in sync via **database triggers** on INSERT, UPDATE, DELETE — never synced in Go code.

---

## AI Provider Interface

All AI features are mediated through a Go interface. No service ever imports an AI SDK directly.

`go
// internal/ai/provider.go
type Provider interface {
    Complete(ctx context.Context, messages []Message) (string, error)
    Stream(ctx context.Context, messages []Message) (<-chan string, error)
    Embed(ctx context.Context, text string) ([]float32, error)
}

// Implementations (v2):
// internal/ai/openai/provider.go
// internal/ai/ollama/provider.go
`

Provider selected at startup from ATLAS_AI_PROVIDER env var.

---

## Module Dependency Graph

`
cmd/atlas
    ├── config
    ├── db
    ├── middleware
    ├── setup          (no module deps)
    ├── auth           (no module deps)
    ├── settings       (no module deps)
    ├── projects
    ├── tasks          → projects (RecalculateProgress)
    ├── notes          → projects (optional link)
    ├── journal        (no module deps)
    ├── search         (queries FTS5 — no module deps)
    ├── dashboard      → tasks, projects, notes, journal
    └── ai             (interface only in v1)
`

---

## File Storage

`
/data/
├── db/atlas.db         Single SQLite file. The entire database.
├── uploads/            User-uploaded documents (v2)
│   └── {uuid}.{ext}    Random filename; original name stored in DB
├── backups/            Point-in-time SQLite copies
└── logs/               File-based logs (if enabled)
`

---

## Deployment

### Development

`ash
make dev    # Air watches .go, .templ, .sql — hot-reloads on change
`

### Production (Docker)

`ash
docker compose up -d
# Atlas available at :8080
# Data persisted in named Docker volume (atlas_data)
`

### Reverse Proxy (Caddy)

`
your-domain.com {
    reverse_proxy localhost:8080
}
`

Caddy handles TLS automatically via Let's Encrypt.

---

*For implementation details, see each module's docs/{module}/architecture.md.*
