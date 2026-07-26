# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.3.0] - 2026-07-26

### Added
- **Knowledge Base Module & Productivity Wiki**:
  - Markdown Editor with Live HTML Preview and one-click formatting toolbar (**Bold**, *Italic*, `Code Block`, Heading, Quote, Bullet List, Task List, `[[Wiki Link]]`).
  - Pre-built note templates: **ADR (Architecture Decision Record)**, **Meeting Notes**, and **Feature Brainstorm**.
  - Bidirectional Wiki Links (`[[Note Title]]`) and automatic "Linked References" backlinks network panel.
  - Live Reading Time & Word Count telemetry bar with background 30s HTMX autosave.
  - Tag Indexing cloud and search filter engine.
- **Executive Command Center Dashboard (`/`)**:
  - Live system telemetry cards (Active Initiatives, Today's Focus Queue, GitHub Stars, Task Completion %).
  - Real-time active initiatives progress bars and direct action links.

---

## [v0.2.0] - 2026-07-26

### Added
- **Tasks Module & Product-Grade Slate Interface**:
  - Full Task CRUD engine with List & 3-Column Kanban Board views (`To Do`, `In Progress`, `Completed`).
  - Linear/Vercel-inspired Slate Obsidian design system (`#090d16`) with zero emojis and crisp Lucide SVG vector iconography.
  - Multi-axis task filtering: Status tabs, Priority (`Critical`, `High`, `Medium`, `Low`), Energy levels (`High`, `Medium`, `Low`), and Project linking.
  - Task Dependencies, estimated/actual minutes, due date tracking, and right slide-over inspector drawer.
  - Automatic project completion % recalculation whenever linked tasks are updated.

---

## [v0.1.0] - 2026-07-26

### Added
- Go server with Chi router and Templ templates
- SQLite database with WAL mode and all migrations
- First-run setup wizard (name, username, timezone, theme, optional demo data)
- Session-based authentication (login, logout)
- CSRF protection on all state-changing requests
- Rate limiting on login endpoint
- Light / Dark / System theme system
- /health endpoint
- Docker and Docker Compose setup
- GitHub Actions CI (test + lint)

---

*Entries are added as features ship.*
