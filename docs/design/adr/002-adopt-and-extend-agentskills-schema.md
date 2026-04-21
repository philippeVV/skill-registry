# ADR-002: Adopt and Extend the agentskills.io Schema

**Date:** 2026-04-18
**Status:** Accepted

## Context

The agentskills.io open standard (`SKILL.md` + `.claude-plugin/marketplace.json`)
is now co-signed by Anthropic, OpenAI, Cursor, Copilot, Goose, Amp, and
Letta. Claude Code natively supports `/plugin marketplace add <url>` and
`/plugin install <name>` against any compliant marketplace manifest. Rolling
our own schema would break this native compatibility and require a custom
client for every install.

## Decision

We adopt the agentskills.io / Anthropic marketplace schema as the base and
extend it with registry-specific fields. We never remove or rename upstream
fields.

**Contributor-authored fields** (in `metadata.json`):
- `name`, `type`, `description`, `tags`, `author`, `license` — required
- `notes` — optional free text for caveats or recommendations
- `upstream` — optional object (`url`, `path`, `ref`) for external tracking
- `install_target` — optional override for type-based install placement

**CI-generated fields** (added to `marketplace.json` by `ci/build_index.py`,
never written by contributors):
- `version` — resolved from the latest git tag for the package
- `path` — directory path within the repo (`packages/<name>`)
- `files` — list of files in the package directory
- `artifact_hash` — SHA256 of the artifact file(s)

The originally planned `x-registry.*` namespaced fields (reviewer,
eval_score, security_scan_status, etc.) were not implemented for v1.
Trust signal fields like install and invocation counts are served from
a separate stats endpoint and merged client-side — they are not baked
into the index (see `docs/stats-contract.md`). Extension fields for
review status and security scanning are deferred to v2.

The registry index (`marketplace.json`) is a valid agentskills.io manifest.
Users can point Claude Code at it natively with no custom client required.

## Consequences

**Positive:**
- `/plugin marketplace add <registry-url>` works out of the box in Claude Code
- Future Claude Code schema updates apply automatically to our base layer
- Compatibility with any tool that speaks the agentskills.io standard
- No maintenance burden on the core schema

**Negative:**
- Bound to upstream spec decisions we don't control
- The spec is currently under-specified in places; we may hit gaps
- Extension fields won't be indexed or understood by third-party tooling

**Neutral:**
- We maintain a CI schema validator that pins to a specific schema version
  and flags breaking upstream changes before they reach production
