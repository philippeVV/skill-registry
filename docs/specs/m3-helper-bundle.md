# M3 Spec — Helper Bundle

## Goal

The registry ships a first-party bundle of utility skills that make it
self-referential. The bundle is itself published as packages in the registry,
tagged `registry-core`, installable in one command.

Exit criterion: A user can install the bundle with
`skr install --tag registry-core` and run `/eval` against a real task.

---

## Task List

### 1. `suggest-packages` Skill

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

**Notes:**
- `skr` handles caching — the skill doesn't need to manage it
- Recommendations should include the install command so the user can
  copy-paste: `skr install <name>`
- This is the v1 agent self-discovery mechanism (bridge until MCP in M7)

---

### 2. `conflict-check` Skill

**Trigger:** User invokes `/conflict-check` — user-invoked only, no
automatic invocation at install time.

**Behavior:**
1. Run at session start (before significant work) or on demand
2. Read the full content of all currently loaded packages from their
   install locations
3. Analyse for two failure modes:
   - **Conflict**: two packages give contradictory instructions
   - **Overlap**: two packages cover the same ground (context budget waste)
4. Report findings with specific package names and the conflicting content
5. Suggest which package to uninstall or which to keep

**Notes:**
- Full artifact content analysis, not description-only
- Non-blocking — presents findings, user decides what to do
- Particularly useful to run after `skr install --tag <tag>` bulk installs

---

### 3. `eval` Skill

**Trigger:** User invokes `/eval <package-name>` before performing a task
they believe the package improves.

**Behavior:**
1. Read `eval/guidelines.md` from the package if present
2. Ask the user to describe the task they are about to perform
3. Spawn a **shadow sub-agent** to perform the same task without the
   package in context:
   - Preferred: spawn with restricted context excluding the package
   - Fallback: explicitly instruct the sub-agent to ignore the package content
   - For file operations: shadow sub-agent works in a git worktree
4. Main session proceeds with the real task as normal (with package loaded)
5. Spawn a **judge sub-agent** that receives both outputs and scores which
   performed better, with reasoning
6. Present results to the user: score, reasoning, which performed better
7. Optionally publish result back to the registry as an evidence point

**Open questions to resolve during implementation:**
- Can sub-agents be spawned with restricted context excluding specific skills?
- Can file read access be denied for a sub-agent's scope?
- How does the main session output get captured for the judge?

See ADR-033 for full design rationale.

---

### 4. `publish-skill` Skill

**Trigger:** User invokes `/publish-skill`

**Behavior:**
1. Ask the user to describe the pattern or workflow they want to package
2. Determine the appropriate package type based on the description
3. Draft the artifact in the correct format for the type
4. Generate `metadata.json` with all required fields — prompt user for
   any missing values (name, tags, author, license)
5. Generate `README.md` with description and usage examples
6. Show all files to the user for review before proceeding
7. On user approval: create a branch, commit the files under
   `packages/<name>/`, open a PR against the registry repo via `gh`

**Notes:**
- Enforces package structure per `docs/package-types.md`
- Validates `metadata.json` locally before committing (mirrors `ci/validate.py`
  rules so contributors get early feedback)
- PR description explains what the package does and why it's useful

---

### 5. `contribute` Skill

**Trigger:** User invokes `/contribute <package-name>`

**Behavior:**
1. Check `skr list` output — confirm the package shows `[modified]`
2. Read the currently installed artifact from its install location
3. Fetch the original artifact from the registry (via `skr info`)
4. Understand what changed and why (prompt user if intent is unclear)
5. Create a branch, update `packages/<name>/` with the local version,
   open a PR via `gh` with a meaningful description of the improvement

**Notes:**
- Relies on `skr list` drift detection to identify modified packages
- Uses `gh` CLI for PR creation — requires `gh` to be installed and
  authenticated

---

### 6. Publishing All Five as Packages

All five skills are published to the registry under `packages/` with
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
| `conflict-check` trigger | User-invoked only — no automatic install-time check |
| `eval` participant count | Main session + shadow sub-agent + judge sub-agent |
| Shadow agent isolation | TBD — depends on Claude Code sub-agent capabilities |
| `publish-skill` PR step | Skill opens PR via `gh` after user review |
| `contribute` drift detection | Reads `skr list [modified]` output |
| Bundle install | `skr install --tag registry-core` (tag install added in M1) |
