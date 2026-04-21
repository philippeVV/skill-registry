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

- **`suggest-packages`** — Fetches the package catalog via `skr search`,
  reads the current repo context, and recommends relevant packages to
  install. The v1 self-discovery mechanism (see ADR-011).

- **`conflict-check`** — Analyses the full loaded context from a fresh
  session (before the first user message), reading complete artifact content
  — not just descriptions — to identify duplicated information or
  contradictory instructions. User-invoked only: `/conflict-check` at
  session start or on demand to audit everything currently loaded in context.

- **`contribute`** — Handles both new package creation and upstreaming
  local modifications. Auto-detects new vs modified packages: for new
  packages, scans for non-skr artifacts and drafts metadata/README; for
  modified packages, reads the installed artifact and opens a PR with the
  changes. Uses `gh` CLI for PR creation. This skill merges the originally
  separate `publish-skill` and `contribute` skills into one.

- **`knowledge-retriever`** — Retrieves relevant domain knowledge from
  installed knowledge packages when the current task could benefit from
  additional context. Reads the knowledge index at
  `~/.claude/knowledge/index.json` to find and load relevant files.

- **`code-review-checklist`** — Guides code reviews with a structured
  checklist covering correctness, security, performance, readability,
  and testing.

**Deferred:**

- **`eval`** — A/B tests a package against the user's real work using a
  shadow sub-agent and a judge sub-agent. Deferred to v2 due to unresolved
  design questions around output capture and session orchestration. See
  ADR-033 and ADR-034.

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
- `contribute` closes the creation loop: usage → pattern → new package
- No dependency on `skr` shelling out to Claude — skills are user-invoked only
- On first `skr install`, all `registry-core` tagged packages are
  auto-bootstrapped as system skills (marked in the lockfile with
  `"system": true`, requiring `--force` to uninstall)

**Negative:**
- Helper skills consume context tokens on every invocation
- Quality of conflict detection depends on Claude's reasoning, not
  deterministic logic

**Neutral:**
- The helper bundle is pinned to a version like any other package
- `registry-core` tag makes the bundle discoverable and installable as a unit
