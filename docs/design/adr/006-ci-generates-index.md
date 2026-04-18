# ADR-006: CI Generates the marketplace.json Index

**Date:** 2026-04-18
**Status:** Accepted

## Context

The `marketplace.json` index must stay consistent with the packages actually
in the repo. We need to decide who owns keeping the index up to date —
contributors, CI, or a runtime generator.

## Decision

CI regenerates `marketplace.json` automatically on every merge to main.
Contributors author only their package files under `packages/<name>/`.
They never touch the index directly.

On merge, CI:
1. Reads all package manifests from `packages/`
2. Rebuilds `marketplace.json` from scratch
3. Commits the updated index back to main

The index in the repo is always the authoritative, current state.

## Consequences

**Positive:**
- Index drift is impossible — it is always derived from the packages present
- Contributors have no index-related burden or merge conflict surface
- A single source of truth: the package files
- Index rebuild is also a validation step — malformed packages fail the build

**Negative:**
- CI must have write access to main to commit the index back
- Adds a commit to main on every merge (minor noise in git log)

**Neutral:**
- The index rebuild script is a core piece of registry infrastructure and
  must be maintained alongside the schema
- For the enterprise S3 version, CI pushes the generated index to S3 instead
  of committing it back to the repo
