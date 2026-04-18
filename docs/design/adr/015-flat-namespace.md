# ADR-015: Flat Package Namespace, Tags for Team Attribution

**Date:** 2026-04-18
**Status:** Accepted

## Context

Enterprise registries often use team namespaces and multi-tier promotion
workflows (team → global) to control package visibility and ownership.
We need to decide whether to implement this upfront or keep the namespace
flat and revisit based on real usage.

## Decision

The registry uses a flat namespace. All packages live under `packages/<name>/`
regardless of which team authored them. There is no team namespace, no
promotion workflow, and no scope field.

Teams use tags to attribute and filter packages:

```json
"tags": ["team-platform", "code-review", "quality"]
```

The web UI and `skr search` support tag filtering, so teams can effectively
browse their own packages without formal namespacing.

Two-tier promotion, team-scoped registries, and RBAC are not on the roadmap.
They will be revisited if real usage creates a genuine need for them.

## Consequences

**Positive:**
- Lowest possible barrier to contribution — no team registration, no scope
  approval, no promotion step
- Tags are already a first-class concept — team attribution costs nothing
- Simpler CI, simpler index, simpler CLI

**Negative:**
- No formal ownership model — package maintainership is implicit
- Name collisions possible if two teams want the same package name (first
  merged wins)
- No way to make a package team-private without a separate registry instance

**Neutral:**
- If team namespaces become necessary, the migration path is
  `packages/<name>/` → `packages/<team>/<name>/` with redirects in the index
