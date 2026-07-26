# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

### Added
- Projects Module (`internal/projects`): Project CRUD operations, status filters (*All*, *Active*, *Completed*, *Archived*), and color accent pills.
- GitHub Integration & Live Metrics: GitHub repository URL import, live stars (⭐), forks (🍴), open issues (🐛), primary language, and last pushed date tracking.
- Tech Stack Breakdown & Badges: Automated GitHub languages API detection and custom tech stack badge tags.
- GitHub Action Triggers: One-click HTMX stats sync (`POST /projects/{id}/sync-github`) and GitHub open issue import as project milestones (`POST /projects/{id}/import-issues`).
- Milestone Checklist & Progress Engine: Interactive milestone checkboxes with automatic project completion percentage recalculation (`Completed / Total * 100`).
- API Documentation update (`docs/api_docs.md`): Registered 10 new HTTP endpoints.
- Projects unit and integration test suite (`tests/unit/projects_test.go` & `tests/integration/projects_flow_test.go`).



---

## [v0.1.0] - TBD

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
