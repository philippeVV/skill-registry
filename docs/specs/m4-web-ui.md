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

**Astro with Preact islands — three interactive components, rest is static HTML:**

The original design called for vanilla JS islands. During implementation,
Preact was chosen because the search/filter component needed reactive state
management (search text, active tag filters, loading states, fetched data).
Preact adds ~3KB gzipped and provides hooks-based state that would have
required significantly more boilerplate in vanilla JS.

1. **`PackageBrowser.tsx`** (home page — `client:load`)
   - Fetches `marketplace.json` from GitHub raw content on mount
   - Filters package grid by search text (name, description, tags)
     and clickable tag pills (multi-select)
   - Responsive grid layout (1/2/3 columns)
   - Skeleton loading states while fetching

2. **`PackageDetail.tsx`** (package detail page — `client:load`)
   - Displays full package metadata and install command
   - Fetches and renders README.md as HTML via `marked` library
   - Fetches stats from `PUBLIC_STATS_URL` if configured
   - Shows `—` placeholder when no stats backend

3. **`Leaderboard.tsx`** (leaderboard page — `client:load`)
   - Two tables: top packages (by invocations for skills, installs for
     others) and top authors (aggregate score with package count)
   - Fetches live stats from `PUBLIC_STATS_URL`
   - Skeleton loading states

All islands use Astro's `client:load` directive for immediate hydration.
Navigation, layout, and page shells are static HTML generated at build time.

---

## Build Config

```
web/
  astro.config.mjs        # site: philippevv.github.io, base: /skill-registry/
  tsconfig.json
  package.json
  src/
    pages/
      index.astro          # home page with PackageBrowser island
      packages/
        [name].astro       # generated from marketplace.json, PackageDetail island
      leaderboard.astro    # Leaderboard island
    components/
      PackageBrowser.tsx   # Preact island — search, filter, package grid
      PackageDetail.tsx    # Preact island — package metadata, README, stats
      Leaderboard.tsx      # Preact island — top packages and authors tables
    lib/
      registry.ts          # fetch helpers for marketplace.json, README, stats
      types.ts             # TypeScript interfaces (Package, PackageStats, etc.)
      constants.ts         # type badge styling (Tailwind classes per type)
    layouts/
      Base.astro           # HTML shell, nav header, footer
    styles/
      global.css           # Tailwind imports
  public/
    favicon.svg
    favicon.ico
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
| Deployment (OSS) | GitHub Pages (base: `/skill-registry/`) |
| Deployment (v2) | S3 + CloudFront |
| Pages | `/`, `/packages/<name>`, `/leaderboard` |
| Search/filter | Preact island (`PackageBrowser.tsx`), fetches marketplace.json client-side |
| Stats when no backend | `—` placeholder always shown, never hidden |
| JS architecture | Astro + Preact islands — three interactive components, rest static |
| Island framework | Preact (not vanilla JS as originally spec'd — reactive state needed) |
| Markdown rendering | `marked` library in `PackageDetail.tsx` |
| Styling | Tailwind CSS v4 via `@tailwindcss/vite` plugin |
| Stats config | `PUBLIC_STATS_URL` env var at build time |
| Data fetching | GitHub raw content URLs for marketplace.json and READMEs |
