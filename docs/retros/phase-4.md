# Phase 4 Retrospective — Projects Module

> **Phase**: 4
> **Started**: 2026-07-26
> **Completed**: 2026-07-26
> **Time Estimated**: 1 day
> **Time Actual**: 1 day

---

## What Was Built

Phase 4 implemented the Projects Module with GitHub integration, Tech Stack detection, and Milestone progress tracking:
- Migration `010_github_projects.sql`: Added `github_url`, `github_stars`, `github_forks`, `github_open_issues`, `github_language`, `github_last_pushed_at`, and `tech_stack` columns to `projects` table.
- GitHub Client & Tech Stack Parser (`internal/projects/github.go`): URL parser, repo stats fetcher (stars, forks, open issues, primary language), language breakdown fetcher (`/repos/owner/repo/languages`), and open issue importer.
- Repository & Service Layer (`internal/projects/repository.go` & `service.go`): Project CRUD operations, milestone checklist management, `RecalculateProgress` math (`Completed / Total * 100`), `SyncGitHubStats`, and `ImportGitHubIssues`.
- HTTP Handlers (`internal/projects/handler.go`): Handlers for `/projects` (list with status filters `all`, `active`, `completed`, `archived`), `/projects/{id}` (detail view), `/projects/{id}/edit`, `/projects/{id}/delete`, `/projects/{id}/sync-github`, `/projects/{id}/import-issues`, and milestone CRUD/toggles.
- Templ UI Templates (`web/templates/projects/`):
  - `list.templ`: Project card grid with status filter tabs, Tech Stack badge pills, GitHub star/fork metrics, progress bars, and creation modal.
  - `detail.templ`: Detailed project page with Tech Stack breakdown header, GitHub Insights Card, Milestone Checklist with HTMX toggles, and edit/delete forms.
  - `github_card.templ`: HTMX partial swap component for live GitHub metrics re-rendering.
- API Documentation (`docs/api_docs.md`): Registered 10 new HTTP API endpoints (`GET/POST /projects`, `POST /sync-github`, `POST /import-issues`, `POST /milestones`).
- Automated Unit & Integration Tests (`tests/unit/projects_test.go` & `tests/integration/projects_flow_test.go`).

---

## What Went Well

- Automated Tech Stack detection via GitHub Languages API provides instant insight into project technology compositions.
- Dynamic milestone progress recalculation gives immediate visual feedback on project completion status.

---

## What Was Harder Than Expected

- Hand-crafting `sql.NullString` scanning for optional SQLite `target_date` and `github_last_pushed_at` fields across Go database drivers.

---

## Decisions Made During Build

| Decision | Reason | Impact |
|----------|--------|--------|
| Tech Stack Badges | Clear visual distinction of project technologies | Rendered on both List cards and Detail header |
| GitHub Issue Import | Quick conversion of repository issues to project milestones | Seamless workflow integration with GitHub repos |

---

## Bugs Found and Fixed

| Bug | Root Cause | Fix |
|-----|-----------|-----|
| Scan error on `target_date` in database models | `target_date` was scanned into `sql.NullTime` when empty strings were present | Updated `TargetDate` and `GithubLastPushedAt` fields in `models.go` to `sql.NullString` |

---

## What to Carry Forward

- Reuse `RecalculateProgress` logic when Phase 5 Tasks Module integrates task completion counts into project progress.

---

## Links

- Phase 4 Commit: feat(projects): implement Projects module with GitHub integration, tech stack detection, and milestone progress tracking
- Documentation: `docs/api_docs.md`, `progress.md`
