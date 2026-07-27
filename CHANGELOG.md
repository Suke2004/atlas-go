# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.7.0] - 2026-07-27

### Added
- **Analytics & Insights Engine (`/analytics`)**:
  - 365-day GitHub-style contribution activity heatmap scaling by activity level across all modules.
  - Telemetry counters tracking Total Contributions, Current Active Streak, and Longest Streak Record.
  - Automated AI Weekly Review Synthesis endpoint powered by active AI provider (`Gemini` or `Ollama`).
  - Current month category expense breakdown distribution cards.

---

## [v0.6.0] - 2026-07-26

### Added
- **Documents Engine Module (`/documents`)**:
  - File upload engine with support for all file types, saving to `/data/uploads/`.
  - In-browser inline previews for PDF files, images, and plain text code/documents.
  - Document metadata editor (title, comma-separated tags) and download endpoint.
  - Automatic FTS5 search indexing across document names and text content.
- **AI Provider Engine & Settings (`/settings`)**:
  - Hot-swappable AI provider architecture supporting **Google Gemini** (Gemini 2.0 Flash / Pro) and local **Ollama** (Llama 3.2).
  - One-click AI document summarisation endpoint.
  - Full system settings panel for appearance theme selection and API credential configuration.

---

## [v0.5.0] - 2026-07-26

### Added
- **Finance Engine Module (`/finance`)**:
  - Monthly cash flow banner tracking Total Income, Total Expenses, Net Savings, and Savings Rate percentage.
  - Zero-based budget allocation and full transaction ledger.
  - 🚀 **Atlas USP: *"Project-Linked Infrastructure Cost Attribution"***: Automatically attributes server hosting, domain fees, and SaaS API expenses directly to active **Atlas Projects**.
- **Tech Skill Roadmap Engine (`/learning`)**:
  - Interactive skill roadmap tracks grid with domain categories (`language`, `framework`, `dsa`, `course`).
  - Active study streak counter and session logger.
  - 🚀 **Atlas USP: *"Code-Verified Proof of Mastery"***: Mastery XP engine combining total study hours and active streaks with linked project proof.

---

## [v0.4.0] - 2026-07-26

### Added
- **Executive Mind-Sync & Journal Module**:
  - Industry-first Executive Mind-Sync velocity correlation engine connecting daily tasks completed and notes created to mood & energy telemetry.
  - Morning Prep vs. Evening Review dual-mode reflection routines.
  - 4-Quadrant reflection cards: **Wins & Achievements**, **Challenges & Blockers**, **Ideas & Breakthroughs**, and **Tomorrow's Focus**.
  - Mood score (1-5), Energy rating (1-5), Sleep duration tracking, and on-blur background auto-save.
- **Global `Ctrl+K` Command Palette Search**:
  - Global `Ctrl+K` / `Cmd+K` keyboard modal overlay available anywhere in Atlas.
  - Powered by SQLite FTS5 full-text search across Projects, Tasks, and Notes.

---

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
