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
extend it with internal-only fields. We never remove or rename upstream fields.
Extensions are namespaced to avoid collisions with future standard additions.

Internal extension fields (added to package metadata):
- `x-registry.reviewer` — who approved the PR
- `x-registry.eval_score` — aggregate score from eval runs
- `x-registry.eval_runs` — number of eval runs contributing to the score
- `x-registry.security_scan_status` — result of prompt-injection scan at publish
- `x-registry.owner_team` — internal team responsible for the package
- `x-registry.install_count` — total installs from the registry
- `x-registry.invocation_count` — total invocations (skill type only)

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
