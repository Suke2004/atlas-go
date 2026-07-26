# Atlas HTTP API Documentation

This document records all HTTP API routes, request parameters, responses, HTMX partial fragments, and status codes implemented across Atlas phases.

---

## Registered API Endpoints

### 1. Health Check
- **Route**: `GET /health`
- **Auth Required**: No
- **Response**: `200 OK` → `{"status": "ok", "app": "Atlas"}`

### 2. First-Run Onboarding Wizard Page
- **Route**: `GET /setup`
- **Auth Required**: No (`FirstRunGate` bypasses if `CountUsers() == 0`)
- **Response Format**: `text/html; charset=utf-8` (Renders Setup Templ template)

### 3. Process Onboarding Account Creation
- **Route**: `POST /setup`
- **Auth Required**: No
- **Form Parameters**: `username` (required), `display_name` (required), `password` (required), `timezone` (optional)
- **Response**: `303 See Other` → `/setup/demo-choice` (Creates session cookie `atlas_session`)

### 4. Show Seed Demo Choice Page
- **Route**: `GET /setup/demo-choice`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Response Format**: `text/html; charset=utf-8`

### 5. Process Demo Data Seeding
- **Route**: `POST /setup/seed`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Form Parameters**: `action` ("seed" | "skip")
- **Response**: `303 See Other` → `/` (Dashboard root)

### 6. Show Login Page
- **Route**: `GET /login`
- **Auth Required**: No
- **Response Format**: `text/html; charset=utf-8`

### 7. Process User Authentication
- **Route**: `POST /login`
- **Auth Required**: No
- **Form Parameters**: `username` (required), `password` (required)
- **Response**: `303 See Other` → `/` on success, or `401 Unauthorized` HTML on failure.

### 8. User Logout
- **Route**: `POST /logout`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Response**: `303 See Other` → `/login` (Clears session cookie)

### 9. Application Dashboard (Root Layout Shell)
- **Route**: `GET /`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Response Format**: `text/html; charset=utf-8`

### 10. Serve Static Assets
- **Route**: `GET /static/*`
- **Auth Required**: No
- **Response**: Static file assets (CSS, JS, Lucide icons, web fonts).

### 11. List Projects & Filter Tabs
- **Route**: `GET /projects`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Query Parameters**:
  - `status` (string, optional) — "all" | "active" | "completed" | "archived"
  - `tag` (string, optional) — Tech stack tag filter e.g. "Go"
  - `view` (string, optional) — "grid" | "table" | "roadmap"
  - `search` (string, optional) — Live query string
- **Response Format**: `text/html; charset=utf-8`

### 12. Create Project
- **Route**: `POST /projects`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Form Parameters**: `name`, `description`, `status`, `color`, `target_date`, `github_url`, `tech_stack`.
- **Response**: `303 See Other` → `/projects`

### 13. Project Detail View
- **Route**: `GET /projects/{id}`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Response Format**: `text/html; charset=utf-8`

### 14. Update Project
- **Route**: `POST /projects/{id}/edit`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Form Parameters**: `name`, `description`, `status`, `color`, `target_date`, `github_url`, `tech_stack`.
- **Response**: `303 See Other` → `/projects/{id}`

### 15. Delete Project
- **Route**: `POST /projects/{id}/delete`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Response**: `303 See Other` → `/projects`

### 16. Sync GitHub Stats
- **Route**: `POST /projects/{id}/sync-github`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Response Format**: `text/html; charset=utf-8` (HTMX `GitHubCard` fragment)

### 17. Import GitHub Issues as Milestones
- **Route**: `POST /projects/{id}/import-issues`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Response**: `303 See Other` → `/projects/{id}`

### 18. Add Milestone
- **Route**: `POST /projects/{id}/milestones`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Form Parameters**: `title`, `due_date`.
- **Response**: `303 See Other` → `/projects/{id}`

### 19. Toggle Milestone Completion
- **Route**: `POST /projects/{id}/milestones/{milestoneID}/toggle`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Response**: `303 See Other` → `/projects/{id}`

### 20. Delete Milestone
- **Route**: `POST /projects/{id}/milestones/{milestoneID}/delete`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Response**: `303 See Other` → `/projects/{id}`

### 21. List Tasks & Kanban View
- **Route**: `GET /tasks`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Query Parameters**: `status`, `priority`, `energy`, `view` ("list" | "kanban"), `search`.
- **Response Format**: `text/html; charset=utf-8`

### 22. Create Task
- **Route**: `POST /tasks`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Form Parameters**: `title`, `description`, `status`, `priority`, `energy_level`, `project_id`, `due_date`, `estimated_minutes`.
- **Response**: `303 See Other` → `/tasks`

### 23. Update Task Status
- **Route**: `POST /tasks/{id}/status`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Form Parameters**: `status` ("todo" | "in_progress" | "done")
- **Response**: `303 See Other` → `/tasks`

### 24. Delete Task
- **Route**: `POST /tasks/{id}/delete`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Response**: `303 See Other` → `/tasks`

### 25. List Notes & Knowledge Base View
- **Route**: `GET /notes`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Query Parameters**: `tag`, `search`, `view` ("grid" | "list").
- **Response Format**: `text/html; charset=utf-8`

### 26. Render New Note Editor
- **Route**: `GET /notes/new`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Query Parameters**: `template` ("adr" | "meeting" | "brainstorm").
- **Response Format**: `text/html; charset=utf-8`

### 27. Create Note
- **Route**: `POST /notes`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Form Parameters**: `title`, `content`, `project_id`, `tags`, `is_pinned`.
- **Response**: `303 See Other` → `/notes/{id}`

### 28. Render Note Detail / Editor
- **Route**: `GET /notes/{id}`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Response Format**: `text/html; charset=utf-8`

### 29. Update Note Content & Wiki Links
- **Route**: `POST /notes/{id}/edit`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Form Parameters**: `title`, `content`, `project_id`, `tags`, `is_pinned`.
- **Response**: `303 See Other` → `/notes/{id}`

### 30. Background Autosave Note
- **Route**: `POST /notes/{id}/autosave`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Form Parameters**: `title`, `content`.
- **Response Format**: `text/html; charset=utf-8` (`200 OK` fragment)

### 31. Daily Reflection & Mind-Sync Dashboard
- **Route**: `GET /journal`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Query Parameters**: `date` (format: `YYYY-MM-DD`).
- **Response Format**: `text/html; charset=utf-8`

### 32. Auto-Save Daily Telemetry & Summary
- **Route**: `POST /journal/save`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Form Parameters**: `entry_date`, `mood_rating`, `energy_rating`, `sleep_hours`, `summary`.
- **Response Format**: `text/html; charset=utf-8` (`200 OK` fragment)

### 33. Add 4-Quadrant Reflection Item
- **Route**: `POST /journal/items`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Form Parameters**: `entry_id`, `entry_date`, `category` ("win" | "problem" | "idea" | "tomorrow"), `content`.
- **Response**: `303 See Other` → `/journal?date={entry_date}`

### 34. Delete Reflection Item
- **Route**: `POST /journal/items/{id}/delete`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Query Parameters**: `date`.
- **Response**: `303 See Other` → `/journal?date={date}`

### 35. Global Command Palette FTS5 Search
- **Route**: `GET /api/search`
- **Auth Required**: Yes (`AuthRequired` middleware)
- **Query Parameters**: `q`.
- **Response Format**: `text/html; charset=utf-8` (Search Results Fragment)



---

## Planned API Endpoints (Upcoming Phases)

| Phase | Method | Route | Description | Response Type |
|-------|--------|-------|-------------|---------------|
| **Phase 6** | `GET` | `/notes` | Knowledge Base editor & wiki links | HTML Page |
| **Phase 8** | `GET` | `/journal` | Daily journal mood/energy entry | HTML Page |
| **Phase 9** | `GET` | `/search` | Global search query | HTMX Search JSON/HTML |
