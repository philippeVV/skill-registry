# ADR-023: CI Workflows

**Date:** 2026-04-18
**Status:** Accepted

## Context

The publish pipeline (ADR-005, ADR-006, ADR-017) and frontend deployment
need to run on appropriate triggers without doing unnecessary work on every
commit.

## Decision

Two GitHub Actions workflows, with path scoping to avoid redundant runs:

**`pr.yml` — runs on every PR, no path scoping**
- Schema validation on changed packages
- Conflict detection against existing index
- LLM review gate (v2 only)
- Must pass before merge is allowed

**`publish.yml` — runs on merge to main, path-scoped:**

| Changed path | Jobs triggered |
|---|---|
| `packages/**` | Rebuild index → bump version tag → deploy frontend |
| `web/**` | Deploy frontend only |
| `skr/**` | Build and release Go binary |
| `ci/**` | No deploy |
| `docs/**` | No deploy |

Changes to `packages/` always trigger a frontend deploy because the index
changed and the static site must be rebuilt from the new index.

## Consequences

**Positive:**
- No unnecessary frontend deploys when only the CLI or docs change
- No unnecessary binary releases when only packages change
- Clear separation of concerns per workflow

**Negative:**
- Path scoping adds complexity to workflow files
- A change touching both `packages/` and `skr/` triggers both job sets —
  acceptable and correct behaviour

**Neutral:**
- `pr.yml` has no path scoping — every PR gets validated regardless of
  what changed, including docs and CI script changes
