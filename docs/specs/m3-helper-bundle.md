# M3 Spec — Helper Bundle

## Goal

The registry ships a first-party bundle of utility skills that make it
self-referential. The bundle is itself published as packages in the registry,
tagged `registry-core`, installable in one command.

Exit criterion: A user can install the bundle with
`skr install --tag registry-core` and successfully run `/conflict-check`
and `/contribute`.

---

## Task List

### 1. `suggest-packages` Skill ✓

**Status:** Complete (delivered pre-M3)

**Trigger:** User invokes `/suggest-packages`

**Behavior:**
1. Shell out to `skr search` to get the full package catalog (uses `skr`'s
   cache — no direct index fetch needed)
2. Read current project context: `CLAUDE.md`, root `README.md` if present,
   and one level of directory structure
3. Match project context against package descriptions and tags
4. Present ranked recommendations with install commands
5. For each recommendation, optionally run `skr info <name>` to surface
   full detail before the user decides

---

### 2. `conflict-check` Skill

**Trigger:** User invokes `/conflict-check` — user-invoked only, no
automatic invocation at install time.

**Behavior:**
1. Inventory all loaded context sources:
   - skr-managed packages (from `~/.config/skr/skr.lock`)
   - User rules (`~/.claude/rules/*.md`, `.claude/rules/*.md`)
   - User skills (`~/.claude/skills/*/SKILL.md`, `.claude/skills/*/SKILL.md`)
   - Knowledge (`~/.claude/knowledge/*/KNOWLEDGE.md`)
   - CLAUDE.md files (project root, `~/.claude/CLAUDE.md`)
2. Cross-reference lockfile to distinguish skr-managed vs user-created
3. Read full artifact content for pairwise comparison
4. Analyse for two failure modes:
   - **Conflict**: two sources give contradictory instructions
   - **Overlap**: two sources cover the same ground (context budget waste)
5. Report findings with specific quotes from each source
6. Present both resolution options neutrally — user decides

**Notes:**
- Scans ALL loaded context, not just skr packages
- Full artifact content analysis, not description-only
- Non-blocking — presents findings, user decides what to do
- Particularly useful to run after `skr install --tag <tag>` bulk installs

---

### 3. `contribute` Skill (merged publish-skill + contribute)

**Trigger:** User invokes `/contribute [path-or-name]`

**Two flows, detected automatically:**

**Flow A — New package:**
1. Scan `~/.claude/skills/`, `~/.claude/rules/`, `.claude/skills/`,
   `.claude/rules/`, cross-reference lockfile
2. Present non-skr artifacts as contribution candidates
3. User picks (or passed via argument, skip scan)
4. Infer type from source location, name from filename/directory
5. Draft description and suggest tags from artifact content
6. Pull author from `git config user.name`, default license to MIT
7. Generate `metadata.json` and `README.md`
8. Show all files to the user for review before proceeding

**Flow B — Modified package:**
1. Confirm modified status via lockfile hash comparison
2. Read installed artifact from disk
3. Only replace the artifact file (SKILL.md/RULE.md/KNOWLEDGE.md) — keep
   registry metadata.json and README.md
4. Draft PR description from diff, ask user for additional context

**Both flows end with:**
1. Clone registry repo to temp directory
2. Create branch from main
3. Commit package files under `packages/<name>/`
4. Open PR via `gh` with description
5. Clean up temp directory

**Notes:**
- Merges the originally separate `publish-skill` and `contribute` skills
- Uses `gh` CLI for PR creation — requires `gh` to be installed and
  authenticated
- Always shows files for review before committing

---

### 4. `eval` Skill — Deferred to v2

See ADR-034. The eval skill's workflow has unresolved design questions
around output capture and session orchestration. Design decisions captured
in ADR-033 and ADR-034 for future implementation.

---

### 5. Publishing All Skills as Packages

All skills are published to the registry under `packages/` with
`registry-core` tag. They must pass `ci/validate.py` and serve as
reference implementations for contributors.

Each includes:
- `SKILL.md` — the skill prompt
- `metadata.json` — with `tags: ["registry-core"]`
- `README.md` — usage instructions

---

## Key Decisions

| Decision | Choice |
|---|---|
| `suggest-packages` index access | Shell out to `skr search` / `skr info` |
| `conflict-check` scope | All loaded context, not just skr packages |
| `conflict-check` trigger | User-invoked only — no automatic install-time check |
| `conflict-check` discovery | Lockfile for skr packages + glob for user config |
| `conflict-check` resolution | Present both options neutrally, user decides |
| `publish-skill` + `contribute` | Merged into single `contribute` skill |
| `contribute` type inference | From source location (rules/ → rule, etc.) |
| `contribute` modified packages | Replace artifact file only, keep metadata/README |
| `contribute` repo access | Clone to temp dir, clean up when done |
| `contribute` PR | Via `gh` CLI, draft description with user input |
| `eval` | Deferred to v2 (ADR-034) |
| Bundle install | `skr install --tag registry-core` (tag install added in M1) |
