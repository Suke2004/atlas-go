# Phase 14 Retrospective — v0.7.0 Analytics & Insights Engine

## What Went Well
- Built a 365-day GitHub-style personal contribution heatmap aggregating activity across all 6 core modules (Tasks, Notes, Journal, Learning, Finance, Documents).
- Integrated AI Weekly Review synthesis endpoint leveraging active AI backend (`Gemini` or `Ollama`).
- Seamless template-layer data structure separation avoiding Go import cycles.

## Harder Than Expected
- Efficiently aggregating daily activity count queries across 6 separate module tables in SQLite.

## Decisions Made
- Consolidated contribution data into an intensity scale (0-4) inside the analytics service to render compact CSS grid elements directly in HTML.
