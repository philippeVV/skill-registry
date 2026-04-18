# M4 Spec — Web UI

## Goal

A static Astro site that serves as the discovery and excitement surface for
the registry. Three pages, minimal JS, deployed automatically on every index
rebuild.

Exit criterion: Someone can browse the registry in a browser, filter by tag,
and find something worth installing.

---

## Deployment

**OSS:** GitHub Pages, deployed from `gh-pages` branch via `publish.yml`
when `packages/**` or `web/**` changes.

**v2:** Same static output deployed to S3 + CloudFront. Only the
deploy step in `publish.yml` changes — the build is identical.

---

## Pages

### `/` — Home

- Hero section: name, one-line description of the registry, install command
  for `skr`
- Search bar (island — client-side JS)
- Tag filter pills (island — same JS as search)
- Package grid: cards showing name, type, description, tags, and stats
  placeholder

**Package card:**
```
┌─────────────────────────────────┐
│ code-review-expert    [skill]   │
│ Expert-level code review        │
│ #code-review #quality           │
│ ↓ 142 installs  ⚡ 891 invokes  │
│ (or placeholder — if no stats)  │
└─────────────────────────────────┘
```

Stats show placeholder (`—`) when no `stats_url` is configured. The field
is always present so the layout is consistent.

### `/packages/<name>` — Package Detail

- Name, type, version, author, license
- Tags
- Install command: `skr install <name>`
- Trust signals (stats island):
  - Install count: `142` or `—`
  - Invocation count: `891` or `—`
  - Eval score: `8.4 (12 runs)` or `—`
- Notes (if present)
- Full README rendered as markdown

Stats island fetches from `stats_url` at page load. Shows `—` placeholder
while loading and when no backend is configured.

### `/leaderboard` — Leaderboard

- Top packages by invocation count (skills) or install count (other types)
- Top authors by aggregate score across their packages
- Stats island fetches from `stats_url` — entire page shows placeholder
  state when no backend configured (ranked list with `—` values, not hidden)

---

## Architecture

**Astro islands — two interactive components, rest is static HTML:**

1. **Search + filter island** (home page only)
   - Loads `marketplace.json` in the browser on mount
   - Filters package grid by search text (name, description, tags)
     and active tag pills
   - Vanilla JS — no search library needed at this scale

2. **Stats island** (package detail + leaderboard)
   - Fetches from `stats_url` at page load
   - Renders counts when available, `—` placeholder otherwise
   - Configured at Astro build time via environment variable:
     `PUBLIC_STATS_URL` (empty = no backend)

Everything else — package cards, README rendering, navigation, leaderboard
skeleton — is pure static HTML generated at build time from `marketplace.json`.

---

## Build Config

```
web/
  astro.config.mjs
  src/
    pages/
      index.astro
      packages/
        [name].astro     # generated from marketplace.json
      leaderboard.astro
    components/
      PackageCard.astro
      SearchFilter.tsx   # island
      StatsBlock.tsx     # island
    layouts/
      Base.astro
  public/
```

**Environment variables:**
- `PUBLIC_STATS_URL` — OTEL stats endpoint for client-side fetching.
  Empty by default (OSS). Set in CI for v2 deploys.

---

## Styling

Tailwind CSS. Visual aesthetic inspired by claude-code-templates — clean,
developer-focused, dark-mode friendly. No custom design system for M4.

---

## CI Integration

`publish.yml` deploy step (path-scoped to `packages/**` and `web/**`):
1. Run `astro build` with env vars injected
2. Push static output to `gh-pages` branch (OSS) or S3 (v2)

The index (`marketplace.json`) is always rebuilt before the frontend build
so the site reflects the current state of packages.

---

## Key Decisions

| Decision | Choice |
|---|---|
| Deployment (OSS) | GitHub Pages |
| Deployment (v2) | S3 + CloudFront |
| Pages | `/`, `/packages/<name>`, `/leaderboard` |
| Search/filter | Client-side island, loads marketplace.json in browser |
| Stats when no backend | `—` placeholder always shown, never hidden |
| JS architecture | Astro islands — search + stats only, rest static |
| Styling | Tailwind CSS |
| Stats config | `PUBLIC_STATS_URL` env var at build time |
