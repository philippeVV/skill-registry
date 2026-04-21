# ADR-028: Package Install Placement

**Date:** 2026-04-18
**Status:** Accepted

## Context

`skr install` needs to know where to place each package artifact on the
user's machine. Placement varies by package type, but some packages may
need non-default locations that the type alone can't express.

## Decision

**Default placement by type:**

| Type | Default target | Structure |
|---|---|---|
| `skill` | `~/.claude/skills/<name>/` | Directory containing `SKILL.md` and optional subdirs (scripts/, references/, etc.). `metadata.json` and `README.md` are excluded during copy — they are registry-only files. |
| `rule` | `~/.claude/rules/<name>.md` | Single file. `RULE.md` from the package is copied and renamed. |
| `knowledge` | `~/.claude/knowledge/<name>/` | Directory containing `KNOWLEDGE.md`. On install, `skr` also updates `~/.claude/knowledge/index.json` — a local index mapping installed knowledge packages to their paths, descriptions, and tags. The `knowledge-retriever` skill reads this index to find relevant knowledge on demand. |

**Why skills install as directories, not single files:** Claude Code loads
skills from `~/.claude/skills/<name>/SKILL.md` and supports progressive
disclosure — the SKILL.md frontmatter is loaded at startup, the body on
invocation, and `references/` files on demand. A directory structure is
required to support this three-tier model and any supplementary files
(scripts, templates, examples).

**Why knowledge uses a local index:** Knowledge packages are loaded on
demand, not at startup. The `index.json` file allows the `knowledge-retriever`
skill to scan available knowledge without reading every file. Each entry
records the package name, description, tags, and relative path. `skr`
manages this index automatically on install and uninstall.

**Claude config root:**
Defaults to `~/.claude/`. User can override in `~/.config/skr/config.json`:

```json
{
  "claude_config_dir": "/custom/path/.claude"
}
```

Or via environment variable: `SKR_CLAUDE_CONFIG_DIR=/custom/path/.claude`

**Per-package override — `install_target` in `metadata.json`:**
An optional field that overrides the type-based default when the author
knows the placement better:

```json
"install_target": "skills" | "rules" | "knowledge"
```

`skr` resolves placement in this order:
1. `install_target` from `metadata.json` (if present)
2. Type-based default (skill→skills, rule→rules, knowledge→knowledge)
3. `claude_config_dir` from user config (for the root path)

## Consequences

**Positive:**
- Sensible defaults mean most packages need no `install_target`
- Package authors can declare intent explicitly for non-obvious types
- User override of Claude config root handles non-standard setups
- Directory-based placement for skills and knowledge supports progressive
  disclosure and supplementary files out of the box
- Knowledge index enables efficient on-demand retrieval without loading
  all packages at startup

**Negative:**
- `install_target` adds a field to the schema that most packages won't use
- Knowledge index requires `skr` to manage a separate file that could
  become inconsistent if manually edited

**Neutral:**
- Skills and knowledge use directory-based uninstall (`os.RemoveAll`);
  rules use single-file removal — both are clean and complete
