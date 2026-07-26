# Phase 10 Retrospective — Settings & System Polish

## What Went Well
- Built dedicated Settings page providing toggle between local Ollama and Google Gemini AI providers.
- Integrated theme switcher and account preferences with instant persistence.

## Harder Than Expected
- Avoiding circular import dependencies when enriching template views.

## Decisions Made
- Kept UI template-specific enrichment models inside `web/templates/...` helpers to ensure strict layer isolation.
