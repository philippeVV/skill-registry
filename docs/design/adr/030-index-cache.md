# ADR-030: Local Index Cache with TTL

**Date:** 2026-04-18
**Status:** Accepted

## Context

`skr search`, `skr info`, and `skr install` all need the `marketplace.json`
index. Fetching it on every command adds latency and requires network access.
We need a caching strategy that keeps commands fast without serving stale data.

## Decision

`skr` maintains a **shallow clone of the entire registry repo** at
`~/.config/skr/cache/repo/` with a 1-hour TTL. The index JSON is also
cached separately at `~/.config/skr/cache/marketplace.json` for fast reads.

**Why clone the whole repo instead of just fetching marketplace.json:**
`skr install` needs access to the actual package files (SKILL.md, RULE.md,
etc.), not just the index metadata. Cloning the repo once gives access to
all package artifacts locally. This avoids per-package HTTP fetches and
means `skr install` after `skr search` requires no additional network calls
if the cache is fresh. The clone is shallow (`--depth 1`) to minimize disk
usage.

**Cache behavior:**
- On any read command (`search`, `info`, `install`), `skr` checks the
  repo freshness via a `.skr-fetched` marker file
- If cache is fresh (< 1 hour since last fetch): use local clone, no
  network call
- If cache is stale or missing: `git fetch --depth 1` + `git reset --hard
  origin/main` (or full clone if no cache exists). A copy of
  `marketplace.json` is also written to the separate cache location.
- If fetch fails (network down), falls back to the stale repo clone
- `skr update` forces a repo refresh regardless of TTL
- All commands work offline if a valid clone exists, regardless of TTL

**Cache locations:**
- `~/.config/skr/cache/repo/` — shallow git clone of the registry
- `~/.config/skr/cache/marketplace.json` — extracted index for fast reads
  when only index data is needed

The TTL of 1 hour is a default. It may be made configurable in `config.json`
if real usage shows it needs tuning.

## Consequences

**Positive:**
- Read commands feel instant after first fetch
- `skr install` after `skr search` is fully local — no additional fetch
- Works offline with cached clone and index
- Reduces load on the registry source

**Negative:**
- Users may see packages up to 1 hour out of date without knowing
- Cache must be invalidated manually after publishing a new package
  (run `skr update`)
- Shallow clone uses more disk than just caching `marketplace.json`,
  though `--depth 1` keeps this minimal (~MBs for a typical registry)
- Requires `git` on the user's PATH (acceptable for the developer audience)

**Neutral:**
- `skr update` serves double duty: refresh cache and update installed packages
- If `git pull` fails (e.g. force-push upstream), `skr` nukes and re-clones
  automatically — no manual intervention needed
