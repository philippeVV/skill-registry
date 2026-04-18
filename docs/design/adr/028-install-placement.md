# ADR-028: Package Install Placement

**Date:** 2026-04-18
**Status:** Accepted

## Context

`skr install` needs to know where to place each package artifact on the
user's machine. Placement varies by package type, but some packages may
need non-default locations that the type alone can't express.

## Decision

**Default placement by type:**

| Type | Default target |
|---|---|
| `skill` | `~/.claude/skills/<name>.md` |
| `knowledge` | appended to `~/.claude/CLAUDE.md` in a managed fence block |
| `fragment` | appended to `~/.claude/CLAUDE.md` in a managed fence block |
| `config` | merged into `~/.claude/settings.json` |

**Claude config root:**
Defaults to `~/.claude/`. User can override in `~/.config/skr/config.toml`:

```toml
claude_config_dir = "/custom/path/.claude"
```

**Per-package override — `install_target` in `metadata.json`:**
An optional field that overrides the type-based default when the author
knows the placement better:

```json
"install_target": "skills" | "claude_md" | "settings"
```

`skr` resolves placement in this order:
1. `install_target` from `metadata.json` (if present)
2. Type-based default
3. `claude_config_dir` from user config (for the root path)

## Consequences

**Positive:**
- Sensible defaults mean most packages need no `install_target`
- Package authors can declare intent explicitly for non-obvious types
- User override of Claude config root handles non-standard setups

**Negative:**
- `install_target` adds a field to the schema that most packages won't use
- Misuse (author points a skill at `settings`) produces broken installs —
  no type-safety at the target level

**Neutral:**
- The managed fence block for `claude_md` target is what enables clean
  uninstall — `skr uninstall` removes the block by name
