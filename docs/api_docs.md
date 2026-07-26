# Atlas — System API Reference & Endpoint Registry

> **Purpose**: Authoritative single source of truth for all HTTP endpoints, routes, request formats, response types, and HTMX swap fragments across the Atlas personal operating system.
> **Rule**: Every phase MUST update this file whenever new routes or endpoints are added.

---

## Global System & Health Endpoints

### 1. Health Check
- **Route**: `GET /health`
- **Auth Required**: No (Public)
- **Description**: Returns system health, version, uptime, and database connection status.
- **Response Format**: `application/json`
- **Sample Response**:
  ```json
  {
    "status": "ok",
    "version": "dev",
    "uptime_seconds": 1420
  }
  ```

---

## System Static Assets

### 2. Static Asset Delivery
- **Route**: `GET /static/*`
- **Auth Required**: No (Public)
- **Description**: Serves local bundled CSS, JS, and font assets (`/static/css/app.css`, `/static/js/htmx.min.js`, `/static/js/alpine.min.js`, `/static/js/lucide.min.js`).
- **Response Format**: CSS / JavaScript / Asset binary.

---

## First-Run Setup & Onboarding Wizard (`internal/setup`)

### 3. Show Setup Wizard
- **Route**: `GET /setup`
- **Auth Required**: No (Intercepted by `FirstRunGate` middleware when `CountUsers() == 0`)
- **Description**: Renders multi-step account creation wizard HTML view.
- **Response Format**: `text/html; charset=utf-8`

### 4. Process First-Run Setup
- **Route**: `POST /setup`
- **Auth Required**: No (First-Run Gate)
- **Parameters** (`application/x-www-form-urlencoded`):
  - `display_name` (string, required) — e.g. "Alex Mercer"
  - `username` (string, required) — e.g. "alex"
  - `email` (string, optional) — e.g. "alex@local"
  - `password` (string, required) — bcrypt hashed (cost 12)
  - `timezone` (string, required) — e.g. "UTC", "Asia/Kolkata"
  - `theme` (string, required) — "system" | "dark" | "light"
- **Behavior**: Creates primary owner account in SQLite, creates session cookie, and redirects to `/setup/demo-choice`.
- **Response**: `302 Found` → `/setup/demo-choice` or `200 OK` with error form.

### 5. Show Demo Choice Screen
- **Route**: `GET /setup/demo-choice`
- **Auth Required**: No (First-Run Gate)
- **Description**: Offers option to seed sample projects, tasks, and notes into database.
- **Response Format**: `text/html; charset=utf-8`

### 6. Process Demo Seed
- **Route**: `POST /setup/seed`
- **Auth Required**: No (First-Run Gate)
- **Behavior**: Executes `setup.SeedDemoData()`, populates sample project, tasks, notes, and journal entry, then redirects to `/login`.
- **Response**: `302 Found` → `/login`

---

## Authentication & Session Management (`internal/auth`)

### 7. Show Login Page
- **Route**: `GET /login`
- **Auth Required**: No
- **Description**: Renders login credential form view.
- **Response Format**: `text/html; charset=utf-8`

### 8. Process Login
- **Route**: `POST /login`
- **Auth Required**: No (Rate limited: ≤5 attempts/minute)
- **Parameters** (`application/x-www-form-urlencoded`):
  - `username` (string, required)
  - `password` (string, required)
- **Behavior**: Validates credentials against bcrypt hash, generates 32-byte session token, stores session token in SQLite `sessions` table, sets HTTP-only `SameSite=Lax` cookie, and redirects to `/`.
- **Response**: `302 Found` → `/` or `200 OK` with error form.

### 9. Logout
- **Route**: `POST /logout`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Behavior**: Destroys session token in SQLite, clears session cookie, and redirects to `/login`.
- **Response**: `302 Found` → `/login`

---

## Protected Application Routes (`layout.Base` Shell)

### 10. Dashboard / Layout Shell Root
- **Route**: `GET /`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Description**: Renders primary layout shell (`base.templ`) with sidebar, topbar, theme switcher, and welcome dashboard widget.
- **Response Format**: `text/html; charset=utf-8`

---

---

## Projects Module (`internal/projects`) — Phase 4

### 11. List Projects & Filter Tabs
- **Route**: `GET /projects`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Query Parameters**:
  - `status` (string, optional) — "all" | "active" | "completed" | "archived"
- **Behavior**: Returns Project Cards grid with Tech Stack badges, GitHub star/fork metrics, progress bars, and status filters. Supports HTMX partial swap (`HX-Target: projects-grid`).
- **Response Format**: `text/html; charset=utf-8`

### 12. Create Project
- **Route**: `POST /projects`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Form Parameters** (`application/x-www-form-urlencoded`):
  - `name` (string, required)
  - `description` (string, optional)
  - `status` (string, optional) — "active" | "completed" | "archived" | "on_hold"
  - `color` (string, optional) — hex color e.g. "#6366f1"
  - `target_date` (string, optional)
  - `github_url` (string, optional) — e.g. "https://github.com/owner/repo"
  - `tech_stack` (string, optional) — comma-separated e.g. "Go, SQLite, HTMX, TailwindCSS"
- **Behavior**: Auto-queries GitHub API for repo stats and tech stack if `github_url` provided, inserts project record into SQLite, and redirects to `/projects`.
- **Response**: `303 See Other` → `/projects`

### 13. Project Detail View
- **Route**: `GET /projects/{id}`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Description**: Renders Project Detail page with header, Tech Stack pills, GitHub Insights Card, Milestone Checklist, and edit/delete modals.
- **Response Format**: `text/html; charset=utf-8`

### 14. Update Project
- **Route**: `POST /projects/{id}/edit`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Form Parameters**: `name`, `description`, `status`, `color`, `target_date`, `github_url`, `tech_stack`.
- **Response**: `303 See Other` → `/projects/{id}`

### 15. Delete Project
- **Route**: `POST /projects/{id}/delete`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Behavior**: Deletes project and associated milestones from SQLite.
- **Response**: `303 See Other` → `/projects`

### 16. Sync GitHub Stats
- **Route**: `POST /projects/{id}/sync-github`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Behavior**: Fetches latest repo stars, forks, open issues, primary language, and tech stack from GitHub API and updates SQLite.
- **Response Format**: `text/html; charset=utf-8` (HTMX `GitHubCard` fragment)

### 17. Import GitHub Issues as Milestones
- **Route**: `POST /projects/{id}/import-issues`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Behavior**: Fetches open GitHub issues and converts them into project milestones.
- **Response**: `303 See Other` → `/projects/{id}`

### 18. Add Milestone
- **Route**: `POST /projects/{id}/milestones`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Form Parameters**: `title` (string, required), `due_date` (string, optional).
- **Behavior**: Inserts milestone, recalculates project progress %, and redirects.
- **Response**: `303 See Other` → `/projects/{id}`

### 19. Toggle Milestone Completion
- **Route**: `POST /projects/{id}/milestones/{milestoneID}/toggle`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Form Parameters**: `completed` ("true" | "false")
- **Behavior**: Toggles milestone completion state, records completion timestamp, recalculates project progress %, and redirects.
- **Response**: `303 See Other` → `/projects/{id}`

### 20. Delete Milestone
- **Route**: `POST /projects/{id}/milestones/{milestoneID}/delete`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Behavior**: Deletes milestone, recalculates project progress %, and redirects.
- **Response**: `303 See Other` → `/projects/{id}`

---

## Planned API Endpoints (Upcoming Phases)

| Phase | Method | Route | Description | Response Type |
|-------|--------|-------|-------------|---------------|
| **Phase 5** | `GET` | `/tasks` | Task list & Kanban view | HTML Page / Partial |
| **Phase 6** | `GET` | `/notes` | Knowledge Base editor & wiki links | HTML Page |
| **Phase 8** | `GET` | `/journal` | Daily journal mood/energy entry | HTML Page |
| **Phase 9** | `GET` | `/search` | Global search query | HTMX Search JSON/HTML |

