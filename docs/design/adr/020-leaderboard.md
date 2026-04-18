# ADR-020: Leaderboard

**Date:** 2026-04-18
**Status:** Accepted

## Context

The leaderboard is a gamification feature to drive contribution engagement —
not a critical trust signal. It needs a home and a ranking mechanism.

## Decision

The leaderboard lives on the web UI only. It is computed at page load from
the OTEL stats endpoint — no separate infrastructure, no index entry.

Two views:
- **Top packages** — ranked by invocation count (skill type) or install count
  (all other types)
- **Top authors** — aggregate score across all packages by the same author

The leaderboard is decorative and motivational. It is not a primary trust
signal and carries no weight in install decisions. Its job is to make
contribution feel rewarding and give people something to compete over.

## Consequences

**Positive:**
- Zero additional infrastructure — reuses existing OTEL stats
- Easy to add, easy to ignore if it doesn't drive behavior
- Drives contribution engagement at low cost

**Negative:**
- Leaderboard is only as good as the stats coverage — if OTEL isn't
  configured, it shows nothing

**Neutral:**
- Can be dropped or restyled without touching registry mechanics
