# Phase 3 Retrospective — Layout Shell

> **Phase**: 3
> **Started**: 2026-07-26
> **Completed**: 2026-07-26
> **Time Estimated**: 1 day
> **Time Actual**: 1 day

---

## What Was Built

Phase 3 delivered the master layout shell and design system for Atlas:
- Master Base Shell (`web/templates/layout/base.templ` & `base_templ.go`): Integrates head, theme engine script, sidebar, topbar, HTMX main content target container (`#main-content`), and global toast container.
- Sidebar Navigation (`web/templates/layout/sidebar.templ` & `sidebar_templ.go`): Sidebar with active route highlighting, HTMX boost (`hx-boost="true"`), collapse/expand toggles, and user account status footer.
- Topbar Controls (`web/templates/layout/topbar.templ` & `topbar_templ.go`): Page title breadcrumbs, `Ctrl+K` global search trigger button, and Light / Dark / System theme switcher dropdown.
- Toast Notifications (`web/templates/layout/toast.templ` & `toast_templ.go`): Auto-dismissing toast notification component supporting success, error, warning, and info levels.
- Design System CSS (`web/static/css/app.css`): CSS custom properties for slate dark mode and light mode, glassmorphism panels (`.glass-panel`), custom subtle scrollbars, and toast animations.
- Local Vendor Assets (`web/static/js/`): Downloaded self-hosted HTMX (2.0.4), Alpine.js (3.14.8), and Lucide Icons UMD scripts to guarantee 100% offline self-hosted operation.
- Server Route Integration (`cmd/atlas/main.go`): Mounted `/static/*` HTTP file server and updated dashboard route `/` to render inside the master base layout.
- Unit Testing (`tests/unit/layout_test.go`): Verified component rendering for Sidebar, Topbar, and Toast elements.

---

## What Went Well

- Local bundling of HTMX, Alpine.js, and Lucide icons ensures Atlas operates completely self-contained without external CDN network calls.
- Theme listener script in `<head>` prevents FOUC (Flash of Unstyled Content) by applying saved theme preferences from `localStorage` instantly before first paint.

---

## What Was Harder Than Expected

- Ensuring helper functions (`navClass`, `toastStyle`, `toastIcon`) were co-located inside compiled `*_templ.go` files for standalone Go package compilation.

---

## Decisions Made During Build

| Decision | Reason | Impact |
|----------|--------|--------|
| `hx-boost="true"` on Sidebar Links | Gives instant SPA-like page swaps without complex frontend frameworks | Smooth navigation transitions with server-rendered Go |
| Local Vendor Script Bundling | Self-hosted privacy and offline independence | Zero external CDN requests |

---

## Bugs Found and Fixed

| Bug | Root Cause | Fix |
|-----|-----------|-----|
| Missing helper functions in `_templ.go` files | Templ helper functions were omitted during manual template compilation | Added `navClass`, `toastStyle`, and `toastIcon` directly into `sidebar_templ.go` and `toast_templ.go` |

---

## What to Carry Forward

- Use `#main-content` as the default swap target for module pages.
- Reuse `layouts.Toast` for HTMX partial responses to trigger instant user notifications.

---

## Links

- Phase 3 Commit: feat(layout): implement responsive application layout shell, theme engine, and static asset delivery
- Documentation: `progress.md`, `ARCHITECTURE.md`
