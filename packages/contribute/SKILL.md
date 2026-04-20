---
name: contribute
description: >
  Contribute local skills, rules, or knowledge to the skill registry.
  Publish new packages or upstream modifications to existing ones.
user-invocable: true
argument-hint: "[path-or-name]"
---

# Contribute

You help the user contribute their local Claude Code configuration to the
skill registry. You handle two flows:

- **New package** — a local skill, rule, or knowledge file that is not in
  the registry yet
- **Modified package** — an skr-installed package the user has changed
  locally

## Step 1 — Determine what to contribute

**If the user provided an argument:** Use it directly. It may be a file
path, a directory path, or the name of an installed package. Skip to the
appropriate flow.

**If no argument:** Scan for contribution candidates:

1. **Non-skr artifacts** — Glob these paths:
   - `~/.claude/skills/*/SKILL.md`
   - `~/.claude/rules/*.md`
   - `.claude/skills/*/SKILL.md`
   - `.claude/rules/*.md`

   Read `~/.config/skr/skr.lock` to get the list of skr-managed packages
   and their install locations. Any artifact found by globbing that is NOT
   at an skr-managed location is a candidate for a new package.

2. **Modified skr packages** — For each entry in the lockfile, compare the
   `hash` (installed content hash) against what is on disk now. If they
   differ, the package has been modified locally and is a candidate for
   upstreaming.

Present the candidates as a numbered list:
```
Artifacts you could contribute:
  1. my-cool-skill (local skill, not in registry)
  2. my-rule (local rule, not in registry)
  3. go-conventions [modified] (local changes vs registry version)
```

Let the user pick.

## Flow A — New package

### Infer metadata

- **type**: Infer from source location — `rules/` → rule, `skills/` →
  skill, `knowledge/` → knowledge
- **name**: Infer from filename or directory name (e.g.,
  `~/.claude/rules/my-rule.md` → `my-rule`,
  `~/.claude/skills/my-skill/` → `my-skill`)
- **description**: Read the artifact content and draft a one-line
  description (max 1024 chars)
- **tags**: Suggest relevant tags based on the artifact content
- **author**: Run `git config user.name` to get the author name
- **license**: Default to MIT

### Generate package files

1. **metadata.json** — All required fields: name, type, description, tags,
   author, license. Must conform to the registry schema:
   - name: kebab-case, pattern `^[a-z0-9][a-z0-9-]*[a-z0-9]$`
   - type: one of `skill`, `rule`, `knowledge`
   - tags: non-empty array of non-empty strings
2. **The artifact file** — Copy the user's file as SKILL.md, RULE.md, or
   KNOWLEDGE.md matching the type
3. **README.md** — Brief description, usage examples, and what the
   package does

### Review

Present all generated files to the user for review. Wait for approval or
edits before proceeding.

## Flow B — Modified package

1. Read the currently installed artifact from disk (location from lockfile)
2. Note: you will compare against the registry version after cloning the
   repo in the next step

## Step 2 — Open a pull request

Both flows converge here:

1. **Clone the registry** to a temporary directory:
   ```
   git clone --depth 1 https://github.com/philippeVV/skill-registry.git /tmp/skr-contribute-<random>
   ```

2. **Create a branch** from main:
   ```
   git checkout -b contribute/<package-name>
   ```

3. **Place the files:**
   - For **new packages**: create `packages/<name>/` with metadata.json,
     the artifact file, and README.md
   - For **modified packages**: replace ONLY the artifact file
     (SKILL.md/RULE.md/KNOWLEDGE.md) in `packages/<name>/`. Keep the
     existing metadata.json and README.md from the registry

4. **For modified packages**: Now that you have the registry version, diff
   it against the user's local version. Draft a PR description summarizing
   what changed. Ask the user: "Anything to add about why you made this
   change?"

5. **For new packages**: Draft a PR description explaining what the package
   does and why it is useful.

6. **Commit and push:**
   ```
   git add packages/<name>/
   git commit -m "feat: add <name> package" (or "feat: update <name> artifact")
   git push origin contribute/<package-name>
   ```

7. **Open a PR** via `gh`:
   ```
   gh pr create --repo philippeVV/skill-registry --title "<title>" --body "<description>"
   ```

8. **Clean up** the temporary directory.

9. **Report the PR URL** to the user.

## Requirements

- `gh` CLI must be installed and authenticated (`gh auth status`)
- `git` must be available
- The user must have permission to push to a fork or the registry repo

## Guidelines

- Always show generated files for review before committing
- Never modify metadata.json or README.md for existing packages in Flow B
- If `gh` is not available, tell the user to install it and authenticate
- If the user's artifact has no frontmatter (for skills), note that they
  may want to add name/description frontmatter before contributing
