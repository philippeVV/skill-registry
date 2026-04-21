# ADR-001: Package as the Top-Level Publishable Unit

**Date:** 2026-04-18
**Status:** Accepted

## Context

The registry stores multiple artifact types: Claude Code skills (slash
commands), knowledge bites, CLAUDE.md fragments, and prompt configs. We needed
a single umbrella term for the thing you publish, discover, and install —
usable consistently across the CLI, index, web UI, and documentation. The word
"skill" is overloaded by Anthropic's own terminology and too narrow to cover
all artifact types.

## Decision

The top-level publishable unit is called a **package**. Every package has a
`type` field that subdivides it:

- `skill` — slash command or system prompt behavior for Claude Code
  (installed as a directory under `~/.claude/skills/<name>/`)
- `rule` — always-on coding convention or constraint
  (installed as `~/.claude/rules/<name>.md`)
- `knowledge` — domain knowledge snippet loaded on demand
  (installed as a directory under `~/.claude/knowledge/<name>/`)

These three types were finalized in a research spike (ADR-029) that
mapped Claude Code's actual context-loading behavior to installable types.
Earlier drafts included `fragment` (CLAUDE.md block) and `config`
(settings.json artifact) — both were dropped because Claude Code's
context model didn't support clean install/uninstall semantics for them.
`hook` and `plugin` types were identified as future candidates but
deferred due to complexity (see `docs/package-types.md`).

The CLI, index, and web UI all use "package" as the primary noun. Type is a
filter, not a separate concept.

## Consequences

**Positive:**
- Single noun simplifies CLI design (`registry install <name>`, not separate
  commands per type)
- Mirrors established ecosystems (npm, pip) — mental model transfers
- Extensible: new types can be added without changing top-level concepts

**Negative:**
- "Package" is generic and doesn't immediately signal Claude-specificity
- Slight mismatch with Anthropic's own language ("skill", "plugin") that users
  may already know

**Neutral:**
- The `type` field becomes a required schema field for all packages
