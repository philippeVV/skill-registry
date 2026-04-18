# ADR-030: Local Index Cache with TTL

**Date:** 2026-04-18
**Status:** Accepted

## Context

`skr search`, `skr info`, and `skr install` all need the `marketplace.json`
index. Fetching it on every command adds latency and requires network access.
We need a caching strategy that keeps commands fast without serving stale data.

## Decision

The index is cached locally at `~/.config/skr/cache/marketplace.json` with
a 1-hour TTL.

**Cache behavior:**
- On any read command (`search`, `info`, `install`), `skr` checks the cache
  age
- If cache is fresh (< 1 hour old): use cache, no network call
- If cache is stale or missing: fetch from registry URL, write to cache,
  proceed
- `skr update` forces a cache refresh regardless of TTL
- All commands work offline if a valid cache exists, regardless of TTL

**Cache location:** `~/.config/skr/cache/marketplace.json`

The TTL of 1 hour is a default. It may be made configurable in `config.toml`
if real usage shows it needs tuning.

## Consequences

**Positive:**
- Read commands feel instant after first fetch
- Works offline with cached data
- Reduces load on the registry source

**Negative:**
- Users may see packages up to 1 hour out of date without knowing
- Cache must be invalidated manually after publishing a new package
  (run `skr update`)

**Neutral:**
- `skr update` serves double duty: refresh cache and update installed packages
