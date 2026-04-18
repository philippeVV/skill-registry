# Package Types

Defines the package types supported by the skill registry, how Claude Code
loads each context source, and the install semantics for `skr`.

---

## Type Summary

| Type | Artifact | Install target | Load behavior |
|------|----------|---------------|---------------|
| `skill` | `SKILL.md` (directory) | `~/.claude/skills/<name>/` | Progressive: description at startup, body on invocation |
| `rule` | `RULE.md` | `~/.claude/rules/<name>.md` | Always-on: loaded at session start |
| `knowledge` | `KNOWLEDGE.md` | `~/.claude/knowledge/<name>/` | On-demand: retrieved via `knowledge-retriever` skill |

---

## skill

### How Claude Code loads skills

Skills follow the [Agent Skills open standard](https://agentskills.io/specification).
Claude Code discovers skills from multiple scopes:

- **Personal:** `~/.claude/skills/<name>/SKILL.md`
- **Project:** `.claude/skills/<name>/SKILL.md`
- **Plugin:** `<plugin>/skills/<name>/SKILL.md` (namespaced as `plugin:skill`)
- **Enterprise:** system-level managed skills

**Progressive disclosure model:**

1. **Tier 1 — Metadata** (always loaded, ~100 tokens per skill): `name` and
   `description` from YAML frontmatter. Claude reads these at session start
   to decide relevance. The budget is 1% of context window (~8,000 chars
   default).
2. **Tier 2 — Instructions** (loaded on activation, <5,000 tokens recommended):
   the Markdown body of SKILL.md. Loaded when Claude decides to invoke the
   skill or the user invokes it via `/skill-name`.
3. **Tier 3 — Resources** (loaded on demand): supporting files in `scripts/`,
   `references/`, `assets/`, `examples/`, `template.md`. Loaded only when
   the skill body references them.

**Auto-invocation:** Claude reads skill descriptions and decides whether to
invoke based on task relevance. Set `disable-model-invocation: true` in
frontmatter to make a skill manual-only. Set `user-invocable: false` to
hide from the `/` menu while keeping Claude's auto-invocation.

**Compaction behavior:** When context fills, skills are re-attached after
summary with first 5,000 tokens. Multiple skills share a 25,000-token
budget; oldest dropped if space runs out.

### SKILL.md frontmatter fields

**Required (agentskills.io standard):**
- `name` — kebab-case, 1-64 chars
- `description` — 1-1024 chars, used for discovery matching

**Optional (Claude Code extensions):**
- `disable-model-invocation` — `true` to prevent auto-invocation
- `user-invocable` — `false` to hide from `/` menu
- `argument-hint` — autocomplete hint (e.g., `[filename]`)
- `allowed-tools` — space-separated pre-approved tools
- `model` — override default model
- `effort` — override effort level
- `context` — `fork` to run in isolated subagent
- `agent` — subagent type when `context: fork`
- `hooks` — lifecycle hooks scoped to this skill
- `paths` — glob patterns limiting activation
- `when_to_use` — additional trigger context

**Optional (agentskills.io standard):**
- `license`, `compatibility`, `metadata`, `allowed-tools`

### Artifact structure

```
<name>/
  SKILL.md              # required
  scripts/              # optional: executable scripts
  references/           # optional: detailed reference docs
  assets/               # optional: static resources
  examples/             # optional: example outputs
  template.md           # optional: templates for Claude to fill
```

### Install semantics

`skr install <name>` copies the entire package directory to
`~/.claude/skills/<name>/`. The directory structure is preserved.
`skr uninstall <name>` removes the directory.

### Dynamic context

SKILL.md supports string substitution (`$ARGUMENTS`, `${CLAUDE_SKILL_DIR}`)
and dynamic context injection (`` !`shell-command` `` or ` ```! ` code blocks).

---

## rule

### How Claude Code loads rules

Rules live in `.claude/rules/` and are loaded at session start as part of
Claude's instruction context. They function like modular CLAUDE.md fragments
with optional path scoping.

**Scopes:**
- **User:** `~/.claude/rules/<name>.md` — applies to all projects
- **Project:** `.claude/rules/<name>.md` — shared via source control

**Path scoping:** Rules support optional YAML frontmatter with a `paths`
field. When present, the rule only activates when Claude is working on files
matching the glob patterns:

```yaml
---
paths:
  - "src/api/**/*.ts"
  - "src/**/*.{ts,tsx}"
---

# Your rule content
```

Without frontmatter, the rule applies globally within its scope.

**Loading behavior:** All rules in `~/.claude/rules/` and `.claude/rules/`
are loaded at session start. Path-scoped rules are filtered at activation
time. Rules are always-on context — they consume context budget regardless
of whether the current task is relevant.

**Size guideline:** Keep rules under 200 lines. For longer content, use the
knowledge type instead (on-demand retrieval, no startup context cost).

### Artifact structure

A single file: `RULE.md`. No directory structure needed.

### Install semantics

`skr install <name>` copies `RULE.md` to `~/.claude/rules/<name>.md`.
`skr uninstall <name>` removes the file.

---

## knowledge

### How Claude Code loads knowledge (it doesn't — we manage it)

Knowledge is the registry's own concept. Claude Code has no native
"knowledge" loading mechanism. Instead, the registry implements a
retrieval-augmented pattern:

1. Knowledge files are stored at `~/.claude/knowledge/<name>/KNOWLEDGE.md`
2. A local index at `~/.claude/knowledge/index.json` catalogs all installed
   knowledge with name, description, and tags
3. An always-on skill (`knowledge-retriever`) reads the index and loads
   relevant knowledge files on demand based on task context

This is the RLM (Retrieval Language Model) pattern: knowledge is never
loaded at startup, consuming zero context budget until needed. The retriever
skill's description (~100 tokens) is the only startup cost.

### Knowledge index

`~/.claude/knowledge/index.json` — maintained by `skr`, never hand-edited:

```json
{
  "packages": [
    {
      "name": "payment-system",
      "description": "How our payment processing pipeline works, including Stripe integration, retry logic, and idempotency keys",
      "tags": ["payments", "stripe", "backend"],
      "path": "payment-system/KNOWLEDGE.md"
    }
  ]
}
```

The retriever skill reads this index (~100 tokens per entry), evaluates
relevance to the current task, and reads the full knowledge file only when
it determines a match.

### The knowledge-retriever skill

An always-on skill (`user-invocable: false`, auto-invocation enabled) that
ships as a system package with the registry. It is auto-installed on the
first `skr install` of any package type. Its job:

1. Read `~/.claude/knowledge/index.json`
2. Evaluate which entries are relevant to the current task
3. Read the relevant `KNOWLEDGE.md` files
4. Incorporate the knowledge into its response

The skill has no manual trigger — Claude invokes it automatically when
domain knowledge would help.

### Artifact structure

```
<name>/
  KNOWLEDGE.md          # required: the knowledge content
```

### Install semantics

`skr install <name>`:
1. Copy `KNOWLEDGE.md` to `~/.claude/knowledge/<name>/KNOWLEDGE.md`
2. Add an entry to `~/.claude/knowledge/index.json` with the package's
   name, description, and tags (read from metadata.json)

`skr uninstall <name>`:
1. Remove `~/.claude/knowledge/<name>/` directory
2. Remove the entry from `~/.claude/knowledge/index.json`

---

## Deferred types

### hook (M7/v2)

Hooks are event-driven automation in Claude Code (27 lifecycle events).
They live as JSON config in settings.json. Packaging hooks requires
merging into and removing from a shared config file — deferred due to
complexity of the merge/unmerge logic.

### plugin (M3+)

Plugins are full distribution bundles (skills + hooks + commands + agents +
MCP servers + output styles). They have their own manifest, caching, and
scope system. Supporting plugin distribution requires understanding and
aligning with Anthropic's plugin install pipeline — deferred to a later
milestone.

---

## Evidence sources

- [Extend Claude with skills](https://code.claude.com/docs/en/skills)
- [Agent Skills specification](https://agentskills.io/specification)
- [How Claude remembers your project](https://code.claude.com/docs/en/memory)
- [Claude Code settings](https://code.claude.com/docs/en/settings)
- [Hooks in Claude Code](https://code.claude.com/docs/en/hooks)
- [Plugins reference](https://code.claude.com/docs/en/plugins-reference)
