# Phase 1 Retrospective — Database Foundation

> **Phase**: 1
> **Started**: 2026-07-26
> **Completed**: 2026-07-26
> **Time Estimated**: 1 day
> **Time Actual**: 1 day

---

## What Was Built

Phase 1 implemented the complete database layer and search foundation for Atlas:
- Goose SQL Migrations (`migrations/001_users.sql` through `009_fts5.sql`) defining schemas for users, sessions, settings, projects, milestones, tasks, task_labels, task_dependencies, notes, note_tags, note_links, journal_entries, journal_items, transactions, budgets, learning_tracks, and learning_sessions.
- Global FTS5 search index (`search_index`) with automatic SQLite database triggers on INSERT, UPDATE, and DELETE for tasks, notes, projects, and journal entries.
- Database connection wrapper (`internal/db/db.go`) enforcing SQLite WAL mode (`_journal_mode=WAL`), foreign keys (`_foreign_keys=ON`), busy timeout (`_busy_timeout=5000`), normal sync (`_synchronous=NORMAL`), and transaction runner (`WithTx`).
- Embedded migrations runner (`internal/db/migrate.go` and `migrations/embed.go`) using `embed.FS` and `goose` library to apply migrations automatically at server boot without external CLI requirements.
- Compiled sqlc query functions (`internal/db/*.sql.go`) for users, sessions, settings, projects, tasks, notes, journal, and search.
- Automated unit tests (`tests/unit/db_test.go`) and integration tests (`tests/integration/db_migration_test.go`) verifying pragmas, transaction rollback, and table creation.
- Demo seed script (`scripts/seed.go`) and database documentation (`docs/database.md`).

---

## What Went Well

- Using Go `embed.FS` with the `pressly/goose/v3` library allows Atlas to run migrations automatically on boot with zero external tool dependencies in production.
- FTS5 virtual table synchronization handled completely via SQLite triggers keeps Go application code simple and guarantees strict search index parity.

---

## What Was Harder Than Expected

- Ensuring `go.sum` and module dependencies were properly tidied across scripts and internal packages.

---

## Decisions Made During Build

| Decision | Reason | Impact |
|----------|--------|--------|
| Embedded Goose migrations | Eliminates dependency on external Goose binary during runtime | Zero setup required when deploying Atlas binary |
| Automatic FTS5 Triggers | Prevents out-of-sync search indices | FTS5 stays 100% synced without manual Go code calls |

---

## Bugs Found and Fixed

| Bug | Root Cause | Fix |
|-----|-----------|-----|
| Missing `bcrypt` entry in `go.sum` for seed script | Seed script introduced `bcrypt` dependency | Ran `go mod tidy` to update `go.sum` |

---

## What to Carry Forward

- Keep all SQL query parameterization inside `queries/*.sql` to maintain type-safe generated methods.
- Maintain test coverage across both unit tests and integration tests for every database schema update.

---

## Links

- Phase 1 Commit: feat(db): implement SQLite WAL database foundation, goose migrations, and FTS5 search index
- Documentation: `docs/database.md`, `ARCHITECTURE.md`
