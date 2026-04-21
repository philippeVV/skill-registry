# ADR-010: Frontend Stack — Astro Static Site

**Date:** 2026-04-18
**Status:** Accepted

## Context

The web UI is a vitrine — a discovery and excitement surface, not a product
feature. Its job is to let people browse packages, filter by tag, and get
interested in what the registry offers. It is expected to be lightly used
but worth investing in over time. Heavy framework overhead is not justified
at this stage.

## Decision

The frontend is a static Astro site with Tailwind CSS and **Preact islands**
for interactive components. At build time, Astro reads `marketplace.json`
(fetched from GitHub raw content) and generates static package listing and
detail pages. No server required — deployed to GitHub Pages (OSS) or
S3/CloudFront (v2).

**Island architecture:** Three Preact components handle all client-side
interactivity:
- `PackageBrowser.tsx` — home page search, tag filtering, and package grid
  (fetches `marketplace.json` client-side for live filtering)
- `PackageDetail.tsx` — package detail page with stats and rendered README
  (uses `marked` for Markdown rendering)
- `Leaderboard.tsx` — top packages and top authors tables with live stats

Everything else is static HTML generated at build time. The Preact islands
use Astro's `client:load` directive for immediate hydration.

The visual aesthetic is inspired by claude-code-templates but the stack is
not forked from it. Astro gives component-based structure without the runtime
JS overhead of a full React/Next.js app.

The stack is not considered critical and may be replaced as needs evolve.
Choosing Astro is a pragmatic default, not a strong commitment.

**Why Preact over vanilla JS:** The M4 spec originally called for vanilla JS
islands. During implementation, Preact was chosen because the search/filter
island needed reactive state management (search text, active tag filters,
loading states, fetched data). Preact adds ~3KB gzipped and provides
hooks-based state management that would have required significantly more
boilerplate in vanilla JS. The same reasoning applied to the stats and
leaderboard islands.

## Consequences

**Positive:**
- Zero server infrastructure for the frontend
- Fast static pages, good for a discovery surface
- Easy to evolve or replace — it's just a vitrine
- CI rebuild on index change keeps it always current

**Negative:**
- Live stats (install/invocation counts) require a client-side fetch at
  page load since they are not baked into the static build
- Astro is unfamiliar to the team — low risk given the UI's limited scope

**Neutral:**
- Stack can be replaced without touching registry mechanics
