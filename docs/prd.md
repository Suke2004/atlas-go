# Atlas — Product Requirements Document (PRD)

> **Version**: 1.1 (Decisions Locked)
> **Date**: 2026-07-26
> **Status**: Approved

---

## 1. Executive Summary

**Atlas** is a self-hosted personal operating system — a single, unified web application that replaces the 10–15 browser tabs a knowledge worker opens every day. Rather than context-switching between notes apps, task managers, calendar tools, finance trackers, and AI assistants, Atlas centralises everything into one coherent, fast, locally-owned interface.

> **Tagline**: *Your personal operating system.*

Atlas is not productivity software. It is the operating system layer above your software — a central brain that connects tasks, knowledge, projects, journal, finance, documents, learning, and AI into one living system.

---

## 2. Problem Statement

### 2.1 The Fragmentation Problem

The average knowledge worker uses **12–15 separate tools** daily:

| Tool Category     | Example Products                |
|-------------------|---------------------------------|
| Tasks             | Todoist, Linear, Notion         |
| Notes             | Obsidian, Notion, Roam          |
| Calendar          | Google Calendar                 |
| Projects          | GitHub, Linear, Jira            |
| Finance           | Splitwise, YNAB, spreadsheets   |
| Documents         | Google Drive, Notion, Dropbox   |
| AI Chat           | ChatGPT, Claude, Gemini         |
| Learning Tracker  | Spreadsheets, Notion            |
| Journal           | Day One, Notion, physical book  |
| Analytics         | RescueTime, manual logs         |

**Problems caused by fragmentation**:
- Context switching destroys focus (average cost: 23 minutes per switch)
- Data lives in silos — you cannot ask "what did I work on last month?" across all tools
- Subscriptions add up ($150–$300/year for a typical stack)
- No single source of truth for personal productivity

### 2.2 Privacy & Ownership Problem

Most tools store data on third-party servers. Users have no guarantees about:
- Data longevity (service shutdowns)
- Data privacy (training AI models)
- Data portability (proprietary formats)

### 2.3 The AI Integration Problem

AI tools have no context about your life. They cannot answer:
- "When did I first start learning Go?"
- "Summarize what I worked on last week"
- "Find all my notes about system design"

---

## 3. Vision & Goals

### 3.1 Product Vision

Atlas is the single tab you keep permanently open. It knows your tasks, projects, notes, journal entries, finances, documents, and learning progress — with an AI that can reason across all of it.

### 3.2 Strategic Goals

| Goal | Metric |
|------|--------|
| Replace 10+ tools with one | User closes ≥10 apps after adoption |
| Sub-second UI response | Every interaction < 300ms |
| Full local ownership | 100% self-hostable, no vendor lock-in |
| AI that knows your life | Natural language queries over all personal data |
| Zero context-switch cost | Any module accessible in ≤2 clicks |

### 3.3 Non-Goals (v1)

- Multi-user / team collaboration — Atlas is **single-user only**, forever intentionally simple
- Native mobile app (responsive web only; mobile-friendly polish → v4)
- Real-time sync across multiple devices (v3)
- Public sharing / publishing
- External REST API (v4)
- All external integrations: GitHub, Gmail, Google Calendar, Weather, LeetCode, Discord — all **v2**
- Sub-tasks and task checklists (v2); task dependencies **are** supported in v1
- Finance CSV import (v2)
- Document OCR (v3)
- Semantic / vector search (v3)
- Knowledge graph visualisation (v2)

### 3.4 Locked Architectural Decisions

| Decision | Answer |
|----------|--------|
| Database | SQLite (WAL mode); PostgreSQL → v3 |
| Users | Single-user only |
| Offline | Local-first; everything works with zero internet |
| Search | SQLite FTS5 only; no Elasticsearch/Meilisearch/Typesense |
| Theme | Light / Dark / System (stored in settings; default: System) |
| First-run | Setup wizard (name, username, password, timezone, theme) + optional demo data |
| AI | Provider interface pattern; OpenAI + Ollama in v2; never coupled to one SDK |
| Form save | On blur (fields); every 30 seconds (note/journal editors) |
| File storage | `/data/db/`, `/data/uploads/`, `/data/backups/`, `/data/logs/` |

---

## 4. Target Users

### 4.1 Primary Persona — "The Builder"

Software engineers, CS students, technical founders, and researchers who:
- Have a high volume of concurrent projects and learning tracks
- Already use productivity systems but feel the friction of fragmentation
- Are comfortable self-hosting (Docker, basic CLI)
- Value privacy and data ownership over SaaS convenience
- Spend 8–12 hours/day in front of a computer

### 4.2 Secondary Persona — "The Learner"

Graduate students and self-taught developers who:
- Track multiple learning paths (DSA, ML, languages, courses)
- Need to connect learning notes to projects
- Want to journal and reflect on progress

---

## 5. Product Modules

### 5.1 Module Overview

```
Atlas
├── Dashboard        — Morning briefing, daily overview
├── Projects         — Project lifecycle management
├── Tasks            — Advanced task manager (Kanban, Calendar, Timeline)
├── Knowledge Base   — Second brain, linked notes
├── Journal          — Daily reflection, mood & energy tracking
├── Documents        — File store with OCR & AI summaries
├── Finance          — Income, expenses, subscriptions, predictions
├── Learning         — Skill tracking (LeetCode, books, courses, papers)
├── AI Workspace     — Chat, PDF Q&A, code explanation
├── Analytics        — Heatmaps, weekly/monthly/yearly reviews
├── Search           — Global Ctrl+K search across all modules
├── Notifications    — Smart reminders and nudges
└── Settings         — Preferences, integrations, themes
```

---

### 5.2 Module 1 — Dashboard

**Purpose**: The default view opened on login. A personalised morning briefing that aggregates live data from all modules.

#### User Stories

| ID | Story |
|----|-------|
| D-01 | See a greeting with the current date/day |
| D-02 | See "Today's Focus" — top 3 priority tasks |
| D-03 | See today's calendar events in chronological order |
| D-04 | See active project progress bars |
| D-05 | See GitHub contribution count for yesterday |
| D-06 | See current LeetCode / DSA learning streak |
| D-07 | See today's and this-month's expenses |
| D-08 | See current weather |
| D-09 | See a motivational quote |
| D-10 | See a journal reminder if no entry written today |

#### Functional Requirements

- Dashboard widgets load independently (HTMX partial updates)
- Each widget auto-refreshes at configurable intervals
- Widgets are re-orderable and hideable via Settings
- Clicking any widget navigates to the corresponding module
- Clock updates every second client-side (no server round-trips)

#### Acceptance Criteria

- Dashboard fully loads in < 1 second on localhost
- Each widget renders independently — one failure does not block others
- "Today's Focus" shows top 3 tasks by priority + deadline

---

### 5.3 Module 2 — Project Manager

**Purpose**: Track every personal and professional project from inception to completion.

#### Data Model

```
Project
├── name            string
├── description     markdown
├── status          [active | paused | completed | archived]
├── priority        [critical | high | medium | low]
├── deadline        date
├── progress        0–100%
├── repository_url  string (optional)
├── documentation   markdown
├── notes           linked to Knowledge Base
├── files           linked to Documents
├── milestones      []Milestone
├── people          []Person
└── timeline        []TimelineEvent
```

#### User Stories

| ID | Story |
|----|-------|
| P-01 | Create a project with name, description, deadline, priority |
| P-02 | Add milestones with dates and completion status |
| P-03 | See a vertical timeline of project events |
| P-04 | Track progress percentage (manual or auto from tasks) |
| P-05 | Link a project to a GitHub repository |
| P-06 | Link knowledge base notes to a project |
| P-07 | Archive completed projects |
| P-08 | See all projects in a card grid with status badges |

#### Functional Requirements

- Progress auto-calculates as: `completed tasks / total tasks × 100`
- Timeline events added automatically when milestones are completed
- Projects filterable by status, priority, and deadline

---

### 5.4 Module 3 — Task Manager

**Purpose**: Advanced personal task management with context about time, energy, project, and dependencies.

#### Data Model

```
Task
├── title               string
├── description         markdown
├── status              [todo | in_progress | done | cancelled]
├── priority            [critical | high | medium | low]
├── energy_required     [deep | shallow | admin]
├── estimated_minutes   int
├── actual_minutes      int
├── deadline            datetime (optional)
├── project_id          FK → Project
├── dependencies        []Task.id
├── labels              []string
└── recurrence          (daily | weekly | monthly | none)
```

#### Views

- **List View** — Default, filterable, sortable
- **Kanban View** — Columns: Todo → In Progress → Done
- **Calendar View** — Tasks plotted by deadline
- **Timeline View** — Gantt-style with dependencies

#### User Stories

| ID | Story |
|----|-------|
| T-01 | Create a task with all fields |
| T-02 | Mark tasks complete with actual time logged |
| T-03 | View tasks in Kanban, List, Calendar, and Timeline views |
| T-04 | Filter tasks by project, label, priority, energy level |
| T-05 | Set task dependencies (Task B blocks Task A) |
| T-06 | Create recurring tasks |
| T-07 | Quick-add a task from any page via keyboard shortcut |
| T-08 | See overdue tasks highlighted |

---

### 5.5 Module 4 — Knowledge Base

**Purpose**: A second brain for capturing, organising, and connecting ideas, notes, and research.

#### Data Model

```
Note
├── title           string
├── content         markdown
├── tags            []string
├── linked_notes    []Note.id (bidirectional)
├── related_project FK → Project (optional)
├── files           []Attachment
├── created_at      datetime
└── updated_at      datetime
```

#### User Stories

| ID | Story |
|----|-------|
| K-01 | Create and edit notes in Markdown with live preview |
| K-02 | Tag notes with multiple labels |
| K-03 | Link notes to other notes (bidirectional) |
| K-04 | Link notes to projects |
| K-05 | See all notes linked to a given note ("backlinks") |
| K-06 | View a visual graph of note connections (v2) |
| K-07 | Search notes by title, content, or tag |
| K-08 | Upload images and files into notes |

#### Note Editor Requirements

- Full Markdown (GFM + math blocks)
- Live preview split-pane
- Auto-save every 30 seconds
- Syntax highlighting in code blocks
- Slash commands for quick formatting

---

### 5.6 Module 5 — Journal

**Purpose**: Daily reflection and mood tracking. Builds a queryable personal history.

#### Data Model

```
JournalEntry
├── date            date (unique per day)
├── mood            [1–5]
├── energy          [1–5]
├── sleep_hours     float
├── wins            []string
├── problems        []string
├── ideas           []string
├── tomorrow_focus  []string
└── content         markdown (free-form)
```

#### User Stories

| ID | Story |
|----|-------|
| J-01 | Create or update today's journal entry |
| J-02 | See a calendar view of journal history |
| J-03 | Track mood, energy, and sleep trends over time |
| J-04 | Ask AI: "What did I work on last month?" |
| J-05 | Receive a daily reminder if no entry written |
| J-06 | See "On this day" — entries from 1 year ago |

---

### 5.7 Module 6 — Documents

**Purpose**: A personal file system with AI-powered search and summarisation.

#### User Stories

| ID | Story |
|----|-------|
| Doc-01 | Upload PDF, images, and text documents |
| Doc-02 | Search across all document contents (OCR full-text) |
| Doc-03 | View AI-generated summaries per document |
| Doc-04 | Tag and categorise documents |
| Doc-05 | Find a specific phrase inside a PDF |
| Doc-06 | Preview PDFs in-browser |

#### Functional Requirements

- OCR runs asynchronously after upload
- Full-text search via SQLite FTS5
- AI summarisation runs once and is cached
- Documents stored on local filesystem (not DB blobs)

---

### 5.8 Module 7 — Finance

**Purpose**: Personal finance tracking with AI insights.

#### User Stories

| ID | Story |
|----|-------|
| F-01 | Log an expense or income manually |
| F-02 | View expenses grouped by category |
| F-03 | See monthly spending charts |
| F-04 | Track subscriptions and see next billing dates |
| F-05 | Set monthly budget per category |
| F-06 | Receive alert when budget is exceeded |
| F-07 | AI insight: "You spent 35% more on food this month" |
| F-08 | Export transactions as CSV |

---

### 5.9 Module 8 — Learning Tracker

**Purpose**: Track progress across all learning tracks.

#### Data Model

```
LearningTrack
├── name            string
├── type            [dsa | course | book | paper | language | framework]
├── progress        0–100%
├── current_streak  int (days)
├── longest_streak  int (days)
├── total_sessions  int
└── notes           linked to Knowledge Base
```

#### User Stories

| ID | Story |
|----|-------|
| L-01 | Create a learning track with type and progress |
| L-02 | Log a session (increments streak) |
| L-03 | See current and longest streak per track |
| L-04 | Link a learning track to Knowledge Base notes |
| L-05 | View a heatmap calendar of learning activity |
| L-06 | See all tracks sorted by recent activity |

---

### 5.10 Module 9 — AI Workspace

**Purpose**: A privacy-first AI assistant with context over your Atlas data.

#### Features

| Feature | Description |
|---------|-------------|
| Chat | General-purpose AI conversation |
| PDF Q&A | Upload a PDF, ask questions about it |
| Ask Atlas | Natural language queries over your own data |
| Note Generation | Generate structured notes from chat |
| Code Explanation | Explain code snippets |
| Journal Summary | Summarise past journal entries |
| Weekly Review | AI-generated weekly productivity report |

#### AI Backends

- OpenAI API (GPT-4o, o1)
- Ollama (local LLM — Llama 3, Mistral, etc.)
- Any OpenAI-compatible endpoint

#### User Stories

| ID | Story |
|----|-------|
| AI-01 | Chat with an AI assistant |
| AI-02 | Ask: "When did I first start learning Go?" |
| AI-03 | Upload a PDF and ask questions about it |
| AI-04 | Ask: "Summarise last week's journal entries" |
| AI-05 | Switch between OpenAI and Ollama in Settings |
| AI-06 | Save AI chat responses as Knowledge Base notes |

---

### 5.11 Module 10 — Analytics

**Purpose**: Visualise productivity, learning, and finance trends.

#### Metrics Tracked

- Tasks completed (daily/weekly/monthly)
- Projects finished
- Hours worked (logged via task actual time)
- Journal entries (streak and total)
- LeetCode problems solved
- Books / courses completed
- Total expenses
- GitHub commits (if integrated)

#### Visualisations

- **Activity Heatmap** — GitHub-style contribution graph
- **Weekly Bar Charts** — Tasks, study sessions, expenses
- **Monthly Trend Lines** — Mood, energy, expenses
- **Yearly Wrapped Report** — Annual summary (v2)

---

### 5.12 Module 11 — Global Search

**Purpose**: One keyboard shortcut to search everything in Atlas.

#### Behaviour

- Triggered by `Ctrl+K` from any page
- Searches across: Tasks, Notes, Projects, Journal, Documents, Transactions, Learning Tracks
- Results grouped by type with icons
- Fuzzy full-text search via SQLite FTS5
- Keyboard navigable (arrow keys + Enter)
- Results appear within 100ms

---

### 5.13 Module 12 — Notifications

**Purpose**: Smart, non-intrusive reminders.

| Type | Trigger |
|------|---------|
| Deadline Today | Task/Project deadline = today |
| Interview Tomorrow | Event title contains "interview" |
| Budget Exceeded | Monthly category spend > budget |
| Journal Reminder | 8 PM, no entry written today |
| Project Inactive | No activity on active project for 7 days |
| Learning Streak Risk | No session logged today by 9 PM |

---

## 6. Non-Functional Requirements

### 6.1 Performance

| Requirement | Target |
|-------------|--------|
| Dashboard initial load | < 1 second |
| Any navigation action | < 300ms |
| Search results | < 100ms |
| Task/note create | < 200ms |
| AI first token (streaming) | < 500ms |

### 6.2 Security

- Single-user auth (username + password, bcrypt, cost ≥ 12)
- First-run wizard creates the only user account; no public registration endpoint
- Session-based auth with HTTP-only, SameSite=Lax cookies
- All routes protected by auth middleware (except `/setup`, `/login`, `/static`)
- CSRF protection on all POST/PUT/DELETE requests
- Rate limiting on `/login` (≤ 5 attempts/minute)
- File uploads: MIME validation, max size enforced, random filename on disk

### 6.3 Privacy

- All data stored locally (SQLite + filesystem)
- No telemetry
- AI calls routable to local Ollama
- Full data export as JSON and CSV

### 6.4 Portability

- Single `docker compose up` to deploy
- Database is a single SQLite file — trivial to backup
- No cloud provider dependency

---

## 7. Roadmap

### Version 1.0 — Foundation (Weeks 1–10)

- [ ] First-run setup wizard (name, username, password, timezone, theme, optional demo data)
- [ ] Authentication (login/logout, session management, rate limiting)
- [ ] Projects module (CRUD, milestones, timeline)
- [ ] Tasks module (List + Kanban views, task dependencies)
- [ ] Knowledge Base (Markdown editor, note linking, 30s autosave)
- [ ] Dashboard (live widgets pulling real data from Projects + Tasks + Notes)
- [ ] Journal module (mood/energy/sleep, on-blur save)
- [ ] Global search (Ctrl+K, SQLite FTS5)
- [ ] Settings page (Light / Dark / System theme, widget preferences)

### Version 2.0 — Integration (Weeks 11–18)

- [ ] AI Workspace (provider interface: OpenAI + Ollama, PDF Q&A, "Ask Atlas")
- [ ] Finance module (manual income/expense tracking, budgets, charts)
- [ ] Documents module (upload, preview, tagging)
- [ ] Learning Tracker module
- [ ] Calendar view for tasks
- [ ] Sub-tasks and task checklists
- [ ] Finance CSV import (rule engine, auto-categorisation)
- [ ] GitHub integration (commit feed widget)
- [ ] Knowledge graph visualisation
- [ ] Notification system (deadline, journal reminder, streak risk)

### Version 3.0 — Intelligence (Weeks 19–26)

- [ ] Document OCR (async, indexed via FTS5)
- [ ] Semantic / vector search (embeddings)
- [ ] Analytics module (heatmaps, weekly/monthly/yearly charts)
- [ ] "Ask Atlas" AI over all personal data
- [ ] Automated weekly AI productivity review
- [ ] PostgreSQL support option

### Version 4.0 — Platform (Weeks 27+)

- [ ] Browser extension for quick capture
- [ ] Mobile-responsive UI overhaul
- [ ] Plugin architecture
- [ ] Multi-device sync
- [ ] External REST API (`/api/v1/`)

---

## 8. Success Metrics

| Metric | Target |
|--------|--------|
| Daily active use | User opens Atlas every day |
| Modules used | ≥ 5 modules used weekly |
| Tools replaced | ≥ 8 external tools closed |
| Journal entries | ≥ 20 consecutive days |
| Search usage | Ctrl+K used ≥ 5× per day |

---

*End of PRD v1.0*
