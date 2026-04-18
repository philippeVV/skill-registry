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

**Registry-managed fields (written by CI or stats service, never by contributors):**

```json
{
  "x-registry.reviewer": "alice",
  "x-registry.security_scan_status": "passed",
  "x-registry.eval_score": 8.4,
  "x-registry.eval_runs": 12
}
```

**Live stats (not in metadata.json — served separately):**

Install count and invocation count are live signals that change with every
use. They are NOT baked into `metadata.json` by CI. Instead, they are served
from a separate stats endpoint and merged client-side by `skr` and the web UI.

The index (`marketplace.json`) contains only static metadata. Stats are a
separate concern with their own architecture (see future ADR on stats layer).

## Consequences

**Positive:**
- Contributor surface is minimal — 7-8 fields, all obvious
- Trust signal fields cannot be gamed by hand-editing metadata
- Live stats are always current, not stale from the last CI run
- Dropping `min_claude_version` removes a high-maintenance field with low value

**Negative:**
- `skr` and the web UI must make two fetches (index + stats) and merge them
- Stats architecture needs its own design (not settled here)

**Neutral:**
- `notes` gives authors an escape hatch for anything not covered by structured
  fields, without bloating the schema
