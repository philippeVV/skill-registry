---
name: conflict-check
description: >
  Scan all loaded Claude Code context for conflicts and redundant overlap.
  Use after installing new packages, at session start, or whenever you
  suspect your configuration has grown inconsistent.
user-invocable: true
disable-model-invocation: true
---

# Conflict Check

You audit **all** context sources loaded in this Claude Code session —
not just packages installed via `skr`, but also hand-written rules, skills,
knowledge, and CLAUDE.md files. Your goal is to find contradictions and
wasted context budget.

## Step 1 — Inventory all context sources

Gather every loaded artifact:

**skr-managed packages:** Read `~/.config/skr/skr.lock`. For each entry,
note the `name`, `type`, and `location`. These are packages installed via
`skr`.

**User rules:** Glob `~/.claude/rules/*.md` and `.claude/rules/*.md`.

**User skills:** Glob `~/.claude/skills/*/SKILL.md` and
`.claude/skills/*/SKILL.md`.

**Knowledge:** Glob `~/.claude/knowledge/*/KNOWLEDGE.md`.

**CLAUDE.md files:** Read `CLAUDE.md` in the project root and
`~/.claude/CLAUDE.md` if they exist.

Cross-reference the lockfile locations against the globbed paths. Anything
in the lockfile is "installed via skr"; everything else is "user-created".

## Step 2 — Read full content

Read the full content of every artifact found in Step 1. You need the
actual text to detect conflicts — descriptions and names are not enough.

## Step 3 — Analyse pairwise

Compare every pair of sources for two failure modes:

### Conflict (high priority)

Two sources give **contradictory** instructions. Examples:
- One says "use snake_case for function names", another says "use camelCase"
- One says "always add error handling", another says "don't add error
  handling unless at system boundaries"
- One says "prefer external packages", another says "prefer stdlib"

### Overlap (lower priority)

Two sources cover the **same ground** without contradicting each other.
This wastes context budget — the same instruction loaded twice consumes
tokens without benefit. Examples:
- A rule and a CLAUDE.md section both describe the same coding conventions
- Two skills with overlapping functionality

## Step 4 — Report findings

For each finding, present:

1. **Category**: Conflict or Overlap
2. **Sources**: The two artifacts involved — name, type, and whether
   skr-managed or user-created
3. **Evidence**: Quote the specific lines from each source that conflict
   or overlap
4. **Resolution options**: Present both options neutrally:
   - If one source is an skr package: "Uninstall with `skr uninstall <name>`"
   - If one source is user config: "Remove the relevant section from
     `<path>`"
   - If both are skr packages: suggest which to uninstall based on which
     is more specific or comprehensive
   - If overlap: note which source is more detailed and suggest keeping
     that one

If no conflicts or overlaps are found, say so.

## Guidelines

- Be thorough but concise — quote specific lines, not entire files
- Do not auto-fix anything. Present findings and let the user decide
- Group findings by category (conflicts first, then overlaps)
- If an artifact is empty or unreadable, note it and move on
