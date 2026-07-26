# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

### Added
- Application Layout Shell (`web/templates/layout/base.templ`, `sidebar.templ`, `topbar.templ`, `toast.templ`)
- Local vendor static asset delivery (`web/static/js/htmx.min.js`, `alpine.min.js`, `lucide.min.js`, `app.css`)
- Light / Dark / System theme switcher with instant `localStorage` head listener
- HTMX SPA-like boosted navigation links (`hx-boost="true"`)
- Global Search trigger button (`Ctrl+K`)
- Reusable toast notification system with color variants and auto-dismiss
- Layout component test suite (`tests/unit/layout_test.go`)


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
