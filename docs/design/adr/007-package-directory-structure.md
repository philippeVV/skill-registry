# ADR-007: Package Directory Structure

**Date:** 2026-04-18
**Status:** Accepted

## Context

Contributors need a clear, consistent structure for authoring packages. The
structure must accommodate multiple artifact types, remain flexible for types
we haven't defined yet, and keep skill-type packages fully compatible with
Claude Code's native marketplace system.

## Decision

Each package lives under `packages/<name>/` with the following layout:

```
packages/
  <name>/
    <ARTIFACT>        # type-specific artifact file (see below)
    metadata.json     # registry metadata: type, tags, owner, eval scores, etc.
    README.md         # human-readable description rendered by the web UI
    eval/             # optional
      scenarios.yaml  # eval test cases
```

**Artifact filename convention by type:**
- `skill` → `SKILL.md` (agentskills.io standard, native Claude Code compatible)
- `knowledge` → `KNOWLEDGE.md`
- `fragment` → `FRAGMENT.md`
- `config` → `CONFIG.json`
- Future types: follow the `<TYPE>.<ext>` convention, declared in `metadata.json`

`metadata.json` is the authoritative source for type. CI uses it to validate
that the correct artifact file is present and to build the index entry.

**Native Claude Code compatibility:**
`SKILL.md`-based packages are fully agentskills.io compliant and work with
native `/plugin marketplace add` + `/plugin install`. Non-skill types are
registry-specific and handled by `skr` only — Claude Code's native plugin
system is expected to ignore or skip them gracefully in the marketplace index.

**Flexibility:**
Package types beyond the four listed above are not locked in. New types can
be introduced by adding a convention to this ADR. The structure (artifact +
metadata.json + README.md) remains stable regardless of type.

## Consequences

**Positive:**
- Skill packages are drop-in compatible with Claude Code's native system
- Consistent structure across types reduces tooling complexity
- `metadata.json` separation keeps artifact files lean (no token bloat)
- New types can be added without restructuring existing packages

**Negative:**
- Non-skill types don't benefit from native `/plugin` install
- `metadata.json` is a second file to maintain alongside the artifact

**Neutral:**
- CI validates presence of the correct artifact filename for the declared type
- The eval directory is optional; CI picks it up if present
