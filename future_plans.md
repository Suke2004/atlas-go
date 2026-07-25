# Atlas — Future Plans

> A complete vision document for Atlas across all versions.
> This is the source of truth for where Atlas is going, not just where it is.

---

## Philosophy

Atlas is not a productivity app. It is a **personal operating system** — a single, permanently open tab that replaces the 10–15 tools a developer uses daily. Every feature decision must serve this vision.

> *One tab replaces 10–15 websites you open every day.*

---

## Versioning Map

| Version | Name | Theme | Status |
|---------|------|-------|--------|
| v0.1.0 | Scaffold | Project setup + auth | ⬜ Not started |
| v0.2.0 | Projects | Projects + Tasks | ⬜ Not started |
| v0.3.0 | Knowledge | Notes + Journal + Search | ⬜ Not started |
| v0.4.0 | Interface | Dashboard + Settings + Polish | ⬜ Not started |
| v0.5.0 | Intelligence | AI Workspace + Finance + Learning | ⬜ Not started |
| v0.6.0 | Documents | Files + OCR + Documents | ⬜ Not started |
| v0.7.0 | Analytics | Heatmaps + Reports + Insights | ⬜ Not started |
| v0.8.0 | Connected | GitHub + Calendar integrations | ⬜ Not started |
| v0.9.0 | Hardening | Security + Performance + Testing | ⬜ Not started |
| v1.0.0 | Stable | Public release | ⬜ Not started |

---

## Version 1 — Foundation

### Milestone M0 — Scaffold (v0.1.0)

**Goal**: A working Go server with full toolchain, first-run wizard, and auth.

`
✦ Go 1.23+ server with Chi router
✦ Templ for HTML templates
✦ SQLite + WAL mode
✦ sqlc for type-safe queries
✦ Goose for migrations
✦ Air for hot-reload
✦ Docker + Docker Compose
✦ GitHub Actions CI (test + lint)
✦ First-run setup wizard (name, username, password, timezone, theme)
✦ Session-based authentication
✦ CSRF protection
✦ /health endpoint
✦ /data directory structure
`

### Milestone M1 — Projects + Tasks (v0.2.0)

`
✦ Projects: CRUD, milestones, timeline, progress auto-calculation
✦ Tasks: CRUD, List view, Kanban view
✦ Task dependencies (replaces sub-tasks for v1)
✦ Task labels, priority, energy level, estimated/actual time
✦ Task → Project linking
✦ Project progress auto-recalculates from task completion
`

### Milestone M2 — Knowledge + Journal + Search (v0.3.0)

`
✦ Knowledge Base: Markdown editor (CodeMirror 6), live preview
✦ Note tagging, bidirectional note linking, backlinks panel
✦ Note autosave (every 30 seconds)
✦ Journal: daily entries, mood/energy/sleep tracking
✦ Journal: wins, problems, ideas, tomorrow sections
✦ Journal: on-blur save, "On this day" partial
✦ Global Search: Ctrl+K modal, SQLite FTS5
✦ Search across Tasks, Notes, Projects, Journal entries
✦ Keyboard-navigable search results
`

### Milestone M3 — Dashboard + Settings (v0.4.0)

`
✦ Dashboard: live widgets pulling real data (no mocks)
  - Today's Focus (top 3 tasks)
  - Active Projects (progress bars)
  - Journal Reminder (if no entry today)
  - Motivational Quote
✦ HTMX independent widget loading (failure isolation)
✦ Settings: Light / Dark / System theme (persisted)
✦ Settings: dashboard widget visibility
✦ Skeleton loaders, empty states, toast notifications
✦ HTMX indicators on all mutations
`

---

## Version 2 — Integration

### Milestone M4 — AI Workspace (v0.5.0)

**AI Provider Interface** (never tied to one SDK):

`go
type Provider interface {
    Complete(ctx, messages) (string, error)
    Stream(ctx, messages) (<-chan string, error)
    Embed(ctx, text) ([]float32, error)
}
`

**Supported backends**:
- OpenAI (GPT-4o, o1)
- Ollama (Llama 3.2, Mistral, Gemma)
- Anthropic Claude (v2+)
- Google Gemini (v2+)

**Features**:
`
✦ Chat workspace (general AI conversation)
✦ PDF Q&A (upload PDF → ask questions about it)
✦ Code Explanation (paste code → explain it)
✦ "Ask Atlas" — NL queries over your personal data
  - "When did I first start learning Go?"
  - "Find all my notes about authentication"
  - "Summarise last week's journal entries"
✦ Save AI responses as Knowledge Base notes
✦ Switch AI provider in Settings
`

### Milestone M5 — Finance (v0.5.0)

`
✦ Manual income/expense logging
✦ Category grouping (Food, Transport, College, etc.)
✦ Monthly budget per category
✦ Budget exceeded alert
✦ Monthly spending chart (Chart.js)
✦ Subscription tracking (next billing date)
✦ Dashboard expense widget (enabled after Finance module)
✦ AI insight: "You spent 35% more on food this month"
✦ CSV export of transactions
`

### Milestone M6 — Learning Tracker (v0.5.0)

`
✦ Learning tracks (DSA, courses, books, papers, languages, frameworks)
✦ Session logging (increments streak)
✦ Current streak + longest streak per track
✦ Heatmap calendar of learning activity
✦ Link tracks to Knowledge Base notes
✦ Dashboard streak widget (enabled after Learning module)
`

### Milestone M7 — Documents (v0.6.0)

`
✦ Upload PDFs, images, text files
✦ In-browser PDF preview
✦ Document tagging and categorisation
✦ AI-generated summary per document (cached)
✦ Full-text indexing (FTS5)
✦ File storage at /data/uploads/ (not DB blobs)
`

### Milestone M8 — Integrations (v0.8.0)

**All integrations are strictly v2. None touch v1.**

`
✦ GitHub integration
  - Yesterday's commit count (dashboard widget)
  - Repository link in Projects
  - GitHub webhook for issue notifications
✦ Notification system (background goroutine)
  - Deadline today / tomorrow
  - Budget exceeded
  - Journal reminder (8 PM)
  - Project inactive for 7 days
  - Learning streak at risk (9 PM)
✦ Calendar view for tasks
✦ Weather widget (dashboard)
`

---

## Version 3 — Intelligence

### Milestone M9 — Semantic Search (v0.7.0)

`
✦ Document OCR (async, Tesseract or cloud)
  - Text extracted after upload
  - Content indexed into FTS5
✦ Semantic / vector search
  - Embeddings stored per note, journal entry, document
  - sqlite-vec or pgvector for similarity search
  - "Find notes similar to this idea"
✦ Enhanced "Ask Atlas" with semantic retrieval
`

### Milestone M10 — Analytics (v0.7.0)

`
✦ GitHub-style activity heatmap (all module activity)
✦ Weekly bar charts: tasks, study sessions, expenses
✦ Monthly trend lines: mood, energy, spending
✦ Tasks completed: daily / weekly / monthly
✦ Hours worked (via task actual_minutes)
✦ Journal entry streak and total count
✦ Learning track completion history
✦ LeetCode problems solved (manual logging)
✦ Yearly wrapped report (v3.1)
`

### Milestone M11 — Automated Reviews (v0.7.0)

`
✦ AI weekly review generation
  "Productive Days: 5/7"
  "Most active project: Atlas"
  "Hours worked: 32"
  "Learning progress: Go 72%"
  "Spending analysis: 12% over budget"
✦ Monthly AI summary
✦ Journal summarisation by time range
✦ "On this day" AI reflection
`

---

## Version 4 — Platform

### Milestone M12 — Extensibility (v1.0.0+)

`
✦ Browser extension for quick capture
  - Quick-add task from any page
  - Clip web content to Knowledge Base
  - Log expense from anywhere
✦ Mobile-responsive UI overhaul
  - Fully usable on phone/tablet
  - Touch-friendly Kanban
✦ Plugin architecture
  - Plugin manifest (YAML)
  - Plugin API (Go interface)
  - First-party plugins: Pomodoro timer, Habit tracker
✦ Multi-device sync
  - Atlas server as sync hub
  - Conflict resolution strategy
✦ External REST API (/api/v1/)
  - Authentication via API keys
  - Full CRUD on all modules
  - Webhook support
✦ PostgreSQL support (as alternative to SQLite)
`

---

## Long-Term Vision (Post v1.0)

These are directional — not committed. They guide architecture decisions now.

### Graph Knowledge Layer
Every note, task, project, journal entry, and document should eventually be a node in a graph.
Edges: "references", "depends on", "inspired by", "completed during".
Visualise as a force-directed graph (D3.js or similar).

### Local LLM Deep Integration
- Ollama as the default AI backend for offline-first use
- Embeddings computed locally (no data leaves the machine)
- Model switching: task-specific (code vs. summarisation vs. QA)
- Fine-tuning on personal data (opt-in)

### Collaboration Mode (Optional)
- Share a read-only view of a project or note
- Guest access with scoped permissions
- Real-time presence indicators

### Atlas Mobile App
- React Native or Flutter wrapper around Atlas server
- Offline read + write with background sync
- Widget for home screen (Today's Focus, streaks)

### Atlas CLI
`ash
atlas task add "Implement auth" --project atlas --priority high
atlas note new "Go interfaces"
atlas journal today --mood 4 --energy 5
atlas search "JWT authentication"
`

---

## Features Intentionally Never in Atlas

These are explicit **non-features** — not deferred, never built:

| Feature | Why |
|---------|-----|
| Social / sharing feed | Atlas is personal, not social |
| Public profiles | Privacy-first |
| Team / org features | Single-user by design |
| Gamification / badges | Not the identity of the product |
| Time tracking SaaS | Atlas logs time, not sells it |
| Recurring billing | Self-hosted = no subscriptions |

---

## Open Research Areas

Problems to solve eventually, not yet decided:

1. **Conflict resolution** — when multi-device sync arrives, how to handle concurrent writes to the same note?
2. **Backup strategy** — automatic /data/backups/ snapshots? Encrypted? Frequency?
3. **Migration path** — how does a user import data from Notion, Obsidian, Todoist?
4. **Search ranking** — FTS5 BM25 is good; when to switch to vector search?
5. **AI context window** — how to efficiently pass personal data to LLM without hitting limits?

---

*This document is updated with every major version milestone.*
