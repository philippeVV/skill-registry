# ADR-026: Index Rebuild Script Language

**Date:** 2026-04-18
**Status:** Accepted

## Context

CI needs a script to read all `packages/<name>/metadata.json` files and
generate a valid `marketplace.json` on every merge to main. This is a CI
utility, not a distributed component.

## Decision

The indexer is a Python script at `ci/build_index.py`. No package manager
setup required — standard library only where possible. Runs directly in
GitHub Actions with the default Python environment.

```
ci/
  build_index.py    # reads packages/, writes marketplace.json
  validate.py       # schema validation on PR
```

## Consequences

**Positive:**
- Fast to write, easy to read and modify
- No module setup or build step
- Python's standard library handles JSON and filesystem traversal cleanly

**Negative:**
- Inconsistent language with `skr` — acceptable since it's purely a CI
  utility never installed by users

**Neutral:**
- If the indexer grows complex enough to warrant tests or packaging,
  revisit at that point
