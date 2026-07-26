# Phase 9 Retrospective — Global Search Module

## What Went Well
- Ctrl+K Command Palette overlay modal integrates seamlessly across all application layouts.
- SQLite FTS5 index triggers sync changes instantly across Tasks, Notes, Projects, and Documents.

## Harder Than Expected
- Aligning SQLite FTS5 table schemas and triggers across multiple migration files.

## Decisions Made
- Consolidated entity full-text search indexing into a single `search_index` virtual table with entity_type discriminator.
