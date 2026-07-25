# Roadmap

> Atlas is a self-hosted personal operating system. This roadmap tracks planned releases.
> See uture_plans.md for full feature details.

---

## Current Status

| Release | Status |
|---------|--------|
| v0.1.0 — Foundation | ✅ Shipped (2026-07-26) |
| v0.2.0 — Projects + Tasks | ⬜ Planned |
| v0.3.0 — Knowledge + Journal + Search | ⬜ Planned |
| v0.4.0 — Dashboard + Settings | ⬜ Planned |
| v0.5.0 — AI + Finance + Learning | ⬜ Planned |
| v1.0.0 — Stable | ⬜ Planned |

---

## v0.1.0 — Foundation *(Shipped 2026-07-26)*

- [x] PRD, TRD, and architecture documentation
- [x] Go server with Chi router and Templ templates
- [x] SQLite + WAL + all migrations
- [x] First-run setup wizard
- [x] Session-based authentication (login, logout, CSRF)
- [x] Light / Dark / System theme system
- [x] Docker + GitHub Actions CI

---

## v0.2.0 — Projects + Tasks

- [ ] Projects: CRUD, milestones, timeline
- [ ] Tasks: CRUD, List view, Kanban view
- [ ] Task dependencies
- [ ] Task → Project linking + progress auto-calculation

---

## v0.3.0 — Knowledge + Journal + Search

- [ ] Knowledge Base: Markdown editor, live preview, 30s autosave
- [ ] Note tagging and bidirectional linking
- [ ] Journal: daily entries, mood/energy/sleep, on-blur save
- [ ] Global search (Ctrl+K, SQLite FTS5)

---

## v0.4.0 — Dashboard + Settings

- [ ] Dashboard: live widgets (Focus, Projects, Journal reminder, Quote)
- [ ] Settings: theme toggle, widget visibility
- [ ] Skeleton loaders, empty states, toast notifications
- [ ] Polish: HTMX indicators, responsive layout

---

## v0.5.0 — AI + Finance + Learning

- [ ] AI Workspace: provider interface (OpenAI + Ollama)
- [ ] Chat, PDF Q&A, "Ask Atlas" (NL queries over personal data)
- [ ] Finance: manual income/expense, budgets, charts
- [ ] Learning Tracker: tracks, streaks, session logging

---

## v0.6.0 — Documents

- [ ] File upload (PDF, images, text)
- [ ] In-browser PDF preview
- [ ] AI-generated document summaries
- [ ] Full-text search of document contents

---

## v0.7.0 — Analytics + Intelligence

- [ ] Activity heatmap (GitHub-style)
- [ ] Weekly/monthly trend charts
- [ ] Document OCR
- [ ] Semantic/vector search (embeddings)
- [ ] Automated weekly AI review

---

## v0.8.0 — Integrations

- [ ] GitHub (commit feed, repo linking)
- [ ] Notification scheduler (deadline, streak, journal reminder)
- [ ] Calendar view for tasks

---

## v1.0.0 — Stable Release

- [ ] All v0.x features complete and tested
- [ ] Full documentation
- [ ] 80%+ test coverage
- [ ] Security hardening pass
- [ ] Performance audit (all responses < 300ms)

---

## Post v1.0 (Future)

- Browser extension for quick capture
- Mobile-responsive overhaul
- Plugin architecture
- REST API (/api/v1/)
- Multi-device sync

---

*See uture_plans.md for detailed feature descriptions and progress.md for current implementation status.*
