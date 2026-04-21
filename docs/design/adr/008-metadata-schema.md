# ADR-008: metadata.json Schema

**Date:** 2026-04-18
**Status:** Accepted

## Context

`metadata.json` is the contract between package contributors and the registry.
It must contain everything CI needs to build the index, the web UI needs to
render a package card, and `skr` needs to install correctly. It must also be
minimal enough that contributors can fill it out without friction.

## Decision

**Contributor-authored fields (required unless noted):**

```json
{
  "name": "code-review-expert",
  "type": "skill",
  "description": "One-line description for search and package cards",
  "tags": ["code-review", "quality"],
  "author": "team-platform",
  "license": "MIT",
  "notes": ""
}
```

- `version` is intentionally absent from `metadata.json`. Having the version
  in a committed file is an anti-pattern — it causes conflicts and forces
  contributors to make a mechanical decision. Version is owned by CI via
  GitHub tags (see ADR-017).
- `notes` is optional — free text for caveats, performance characteristics,
  or model tier recommendations (e.g. "works best with a capable model").
  Not structured; not indexed. Rendered as-is in the web UI.
- `min_claude_version` is intentionally omitted. Model version tracking is
  fragile, vendor-specific, and difficult to maintain as model names evolve.
  Authors use `notes` if they need to signal performance characteristics.

**CI-generated fields (written to `marketplace.json` by `ci/build_index.py`,
never by contributors):**

```json
{
  "version": "1.0.0",
  "path": "packages/code-review-expert",
  "files": ["SKILL.md", "metadata.json", "README.md"],
  "artifact_hash": "sha256:abc123..."
}
```

- `version` — resolved from the latest git tag (`<name>@<semver>`)
- `path` — directory path within the repo
- `files` — all files in the package directory
- `artifact_hash` — SHA256 of the type-specific artifact file(s)

The originally planned `x-registry.*` namespaced fields (reviewer,
eval_score, security_scan_status) were not implemented for v1. These are
deferred to v2 along with the LLM review gate and security scanning
pipeline (see ADR-005, ADR-016).

**Live stats (not in metadata.json — served separately):**

Install count and invocation count are live signals that change with every
use. They are NOT baked into `metadata.json` or `marketplace.json` by CI.
Instead, they are served from a separate stats endpoint and merged
client-side by `skr info` and the web UI (see `docs/stats-contract.md`).

## Consequences

**Positive:**
- Contributor surface is minimal — 6 required fields plus optional notes/upstream
- CI-generated fields cannot be gamed by hand-editing metadata
- Live stats are always current, not stale from the last CI run
- Dropping `min_claude_version` removes a high-maintenance field with low value

**Negative:**
- `skr` and the web UI must make two fetches (index + stats) and merge them
- Stats architecture needs its own design (not settled here)

**Neutral:**
- `notes` gives authors an escape hatch for anything not covered by structured
  fields, without bloating the schema
