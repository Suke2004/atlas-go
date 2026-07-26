# Roadmap

> Atlas is a self-hosted personal operating system. This roadmap tracks planned releases.
> See uture_plans.md for full feature details.

---

## Current Status

| Release | Status |
|---------|--------|
| v0.1.0 — Foundation | ✅ Shipped (2026-07-26) |
| v0.2.0 — Projects + Tasks | ✅ Shipped (2026-07-26) |
| v0.3.0 — Knowledge Base | ✅ Shipped (2026-07-26) |
| v0.4.0 — Executive Mind-Sync Journal + Global Search | ✅ Shipped (2026-07-26) |
| v0.5.0 — Finance Engine + Tech Skill Roadmap | ✅ Shipped (2026-07-26) |
| v0.6.0 — Documents Engine & AI Provider Switching | ✅ Shipped (2026-07-26) |
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

## v0.2.0 — Projects + Tasks *(Shipped 2026-07-26)*

- [x] Projects: CRUD, milestones, tech stack tags, GitHub REST integration
- [x] Tasks: CRUD, List view, 3-column Kanban view
- [x] Task dependencies, priority, energy levels
- [x] Task → Project linking + progress auto-calculation
- [x] Linear/Vercel-inspired Slate Obsidian theme with zero emojis & Lucide SVG icons

---

## v0.3.0 — Knowledge Base + Command Center *(Shipped 2026-07-26)*

- [x] Knowledge Base: Markdown Editor with Live Preview
- [x] Note Template Library (ADR, Meeting Notes, Brainstorming)
- [x] One-Click Formatting Toolbar
- [x] Wiki-style internal links `[[Note Title]]` + Bidirectional Backlinks
- [x] Reading time & word count telemetry + 30s background autosave
- [x] Executive Command Center root `/` dashboard with active initiative progress bars

---

## v0.6.0 — Documents Engine & AI Provider Switching *(Shipped 2026-07-26)*

- [x] Documents: Upload, file storage in `/data/uploads`, inline preview for PDF, images, text
- [x] Document metadata editing (title, JSON tags) & FTS5 full-text search indexing
- [x] AI Provider Switching: Hot-swappable interface between local Ollama and Google Gemini
- [x] AI Document Summarisation endpoint using active AI Provider backend
- [x] System Settings: Theme options, AI provider configuration, and credentials management

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
