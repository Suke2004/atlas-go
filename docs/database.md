# Atlas — Database Documentation

> Single-user SQLite database schema, WAL configuration, indices, and FTS5 triggers.

---

## Configuration & Pragmas

Atlas configures SQLite for high concurrent performance and strict durability using WAL (Write-Ahead Logging) mode:

```sql
PRAGMA journal_mode = WAL;       -- Concurrent readers + single writer
PRAGMA foreign_keys = ON;        -- Enforce foreign key constraints
PRAGMA busy_timeout = 5000;      -- 5-second wait on table locks
PRAGMA synchronous = NORMAL;     -- Balance between performance & durability
PRAGMA cache_size = -64000;      -- 64MB memory page cache
```

---

## Schema Overview & Entity Relationship

```
 ┌──────────────┐       ┌────────────────────┐       ┌────────────────────┐
 │    users     │──────<│      sessions      │       │      settings      │
 └──────────────┘       └────────────────────┘       └────────────────────┘
        │
        ├──────────────<┌────────────────────┐
        │               │      projects      │
        │               └────────────────────┘
        │                          │
        │                          ├────────────────<┌────────────────────┐
        │                          │                 │     milestones     │
        │                          │                 └────────────────────┘
        ├──────────────<┌──────────┴─────────┐
        │               │       tasks        │──────<┌────────────────────┐
        │               └────────────────────┘       │    task_labels     │
        │                          │                 └────────────────────┘
        │                          ├────────────────<┌────────────────────┐
        │                          │                 │ task_dependencies  │
        │                          │                 └────────────────────┘
        ├──────────────<┌──────────┴─────────┐
        │               │       notes        │──────<┌────────────────────┐
        │               └────────────────────┘       │     note_tags      │
        │                          │                 └────────────────────┘
        │                          ├────────────────<┌────────────────────┐
        │                          │                 │     note_links     │
        │                          │                 └────────────────────┘
        ├──────────────<┌──────────┴─────────┐
        │               │  journal_entries   │──────<┌────────────────────┐
        │               └────────────────────┘       │   journal_items    │
        │                                            └────────────────────┘
        ├──────────────<┌────────────────────┐
        │               │    transactions    │
        │               └────────────────────┘
        └──────────────<┌────────────────────┐
                        │  learning_tracks   │──────<┌────────────────────┐
                        └────────────────────┘       │ learning_sessions  │
                                                     └────────────────────┘
```

---

## Table Definitions

### 1. `users`
Core user account (single-user model).
- `id`: INTEGER PRIMARY KEY AUTOINCREMENT
- `username`: TEXT UNIQUE NOT NULL
- `display_name`: TEXT NOT NULL
- `email`: TEXT UNIQUE
- `password_hash`: TEXT NOT NULL (bcrypt cost ≥ 12)
- `timezone`: TEXT NOT NULL DEFAULT 'UTC'
- `created_at`, `updated_at`: DATETIME

### 2. `sessions`
Authentication sessions.
- `id`: TEXT PRIMARY KEY (UUID / hex)
- `user_id`: INTEGER FK → `users(id)` ON DELETE CASCADE
- `expires_at`: DATETIME NOT NULL

### 3. `settings`
User preference key-value store.
- `user_id`: INTEGER FK → `users(id)` ON DELETE CASCADE
- `key`: TEXT NOT NULL
- `value`: TEXT NOT NULL
- PRIMARY KEY (`user_id`, `key`)

### 4. `projects`
Projects and goals.
- `id`: INTEGER PRIMARY KEY AUTOINCREMENT
- `user_id`: INTEGER FK → `users(id)` ON DELETE CASCADE
- `name`: TEXT NOT NULL
- `description`: TEXT
- `status`: TEXT ('active', 'completed', 'archived', 'on_hold')
- `color`: TEXT ('#3b82f6')
- `progress_percentage`: INTEGER (0-100)
- `target_date`: DATETIME

### 5. `milestones`
Project milestones.
- `id`: INTEGER PRIMARY KEY AUTOINCREMENT
- `project_id`: INTEGER FK → `projects(id)` ON DELETE CASCADE
- `title`: TEXT NOT NULL
- `due_date`: DATETIME
- `is_completed`: BOOLEAN

### 6. `tasks`
Task tracking item.
- `id`: INTEGER PRIMARY KEY AUTOINCREMENT
- `user_id`: INTEGER FK → `users(id)` ON DELETE CASCADE
- `project_id`: INTEGER FK → `projects(id)` ON DELETE SET NULL
- `title`: TEXT NOT NULL
- `description`: TEXT
- `status`: TEXT ('todo', 'in_progress', 'done')
- `priority`: TEXT ('low', 'medium', 'high', 'critical')
- `energy_level`: TEXT ('low', 'medium', 'high')
- `due_date`: DATETIME
- `estimated_minutes`, `actual_minutes`: INTEGER

### 7. `task_labels`
Task tag labels.
- `task_id`: INTEGER FK → `tasks(id)` ON DELETE CASCADE
- `label`: TEXT NOT NULL
- PRIMARY KEY (`task_id`, `label`)

### 8. `task_dependencies`
Flat task dependency mapping.
- `task_id`: INTEGER FK → `tasks(id)` ON DELETE CASCADE
- `depends_on_id`: INTEGER FK → `tasks(id)` ON DELETE CASCADE
- PRIMARY KEY (`task_id`, `depends_on_id`)

### 9. `notes`
Knowledge Base documents.
- `id`: INTEGER PRIMARY KEY AUTOINCREMENT
- `user_id`: INTEGER FK → `users(id)` ON DELETE CASCADE
- `project_id`: INTEGER FK → `projects(id)` ON DELETE SET NULL
- `title`: TEXT NOT NULL
- `content`: TEXT NOT NULL
- `is_pinned`: BOOLEAN

### 10. `note_tags` & `note_links`
Note tags and bidirectional note connections.

### 11. `journal_entries` & `journal_items`
Daily entries with mood/energy/sleep metrics + bullet points.

---

## Global Search Index (FTS5)

`search_index` is a SQLite FTS5 virtual table synchronized automatically by database triggers:

```sql
CREATE VIRTUAL TABLE search_index USING fts5(
    entity_type,          -- 'task' | 'note' | 'project' | 'journal'
    entity_id UNINDEXED,
    user_id UNINDEXED,
    title,
    content,
    tags
);
```

### Automatic Sync Triggers
- `trg_tasks_fts_insert` / `update` / `delete`
- `trg_notes_fts_insert` / `update` / `delete`
- `trg_projects_fts_insert` / `update` / `delete`
- `trg_journal_fts_insert` / `update` / `delete`
