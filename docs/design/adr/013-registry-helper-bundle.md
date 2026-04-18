# ADR-013: Registry Helper Bundle

**Date:** 2026-04-18
**Status:** Accepted

## Context

Several registry features — eval, conflict detection, self-discovery — require
reasoning that is better handled by Claude than by CLI code. Rather than
implementing complex logic in `skr`, these features are themselves packages
in the registry: a first-party helper bundle that ships with the registry and
is installed as part of onboarding.

## Decision

The registry publishes a first-party **helper bundle** — a curated set of
utility skills tagged `registry-core`. Users install the bundle once:

```
skr install --tag registry-core
```

**Included skills:**

- **`eval`** — A/B tests a package against the user's real work using a shadow
  sub-agent and a judge sub-agent. See ADR-012.

- **`conflict-check`** — Analyses the full loaded context from a fresh
  session (before the first user message), reading complete artifact content
  — not just descriptions — to identify duplicated information or
  contradictory instructions. User-invoked only: `/conflict-check` at
  session start or on demand to audit everything currently loaded in context.

- **`suggest-packages`** — Fetches `marketplace.json`, reads the current repo
  context, and recommends relevant packages to install. The v1 self-discovery
  mechanism (see ADR-011).

- **`publish-skill`** — Helps the agent draft a new package from an observed
  pattern and submit it as a pull request to the registry. Enforces the
  package structure and metadata schema before submission.

- **`contribute`** — Handles upstreaming local modifications to an installed
  package. Runs `skr diff <name>` to get the diff, understands the changes,
  writes a PR description, and opens a PR against the registry repo via `gh`.
  Preferred over a CLI command because the agent can reason about the diff,
  write a meaningful PR description, and handle edge cases.

**`skr` and skills — division of responsibility:**
`skr` handles mechanics: fetching artifacts, placing files at the right scope,
tracking installs via OTEL, enforcing security policy. Claude skills handle
reasoning: conflict detection, eval comparison, discovery recommendations,
package authoring. The CLI delegates to skills for anything requiring judgment.

## Consequences

**Positive:**
- Reasoning stays in Claude, not in fragile CLI heuristics
- The helper bundle is itself a showcase of what the registry can do —
  it eats its own dog food
- Adding new helper skills is just publishing a package, not a CLI release
- `publish-skill` closes the creation loop: usage → pattern → new package
- No dependency on `skr` shelling out to Claude — skills are user-invoked only

**Negative:**
- Helper skills consume context tokens on every invocation
- Quality of conflict detection depends on Claude's reasoning, not
  deterministic logic

**Neutral:**
- The helper bundle is pinned to a version like any other package
- `registry-core` tag makes the bundle discoverable and installable as a unit
