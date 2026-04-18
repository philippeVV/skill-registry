# ADR-014: `skr` CLI — Commands, Language, and Lockfile

**Date:** 2026-04-18
**Status:** Accepted

## Context

`skr` is the primary interface for humans and agents to interact with the
registry. We need to settle the full command surface, implementation language,
and lockfile format before implementation begins.

## Decision

### Language

`skr` is written in Go. Single binary, no runtime dependency, easy to
distribute internally.

### Commands

```
skr install <name>[@version]   # fetch package, place at correct scope,
                                # record in lockfile
skr install --tag <tag>        # install all packages matching a tag
skr uninstall <name>           # remove package files, clean lockfile entry
skr update [<name>]            # update one or all packages to latest version
skr search <query>             # search index by name, tag, description
skr list                       # show installed packages with version and
                                # drift status
skr info <name>                # package detail: description, tags, trust
                                # signals, eval score
skr --registry <url> <cmd>     # override registry source for any command
```

`skr publish` and `skr contribute` are intentionally absent. Package
submission and upstreaming local changes are handled by the `publish-skill`
and `contribute` skills in the registry helper bundle — the coding agent
is better suited to reasoning about diffs and writing PR descriptions than
CLI code.

### Lockfile

The lockfile lives at `~/.config/skr/skr.lock` (user-level installs).
It tracks every installed package with enough information to detect drift
and reproduce the install state exactly.

```json
{
  "registry": "https://github.com/org/skill-registry",
  "packages": {
    "code-review-expert": {
      "version": "1.2.0",
      "type": "skill",
      "location": "~/.claude/skills/code-review-expert.md",
      "hash": "sha256:abc123...",
      "installed_at": "2026-04-18T10:00:00Z",
      "registry_hash": "sha256:abc123..."
    }
  }
}
```

- `hash` — SHA256 of the currently installed artifact file
- `registry_hash` — SHA256 of the artifact at install time from the registry
- If `hash != registry_hash`, the file was manually modified locally
- `skr list` shows a `[modified]` flag next to drifted packages
- the `contribute` helper skill reads the installed file directly and opens
  a PR — no `skr diff` command needed

### Install-time conflict check

When `skr install` runs and the `conflict-check` helper skill is already
installed, `skr` invokes it against the incoming package before completing.
Advisory only — user can proceed past a warning.

## Consequences

**Positive:**
- Go binary is easy to distribute, no runtime to manage
- Lockfile enables reproducible installs and drift detection
- Hash-based drift detection is exact — no false positives
- Contribution workflow delegated to the coding agent via helper skills —
  better reasoning, better PR descriptions, no CLI complexity

**Negative:**
- Lockfile at user level means no shared lockfile for team installs
  (project-level install support is a future concern)

**Neutral:**
- `registry_hash` vs `hash` comparison is the core drift detection mechanism
  and must be maintained correctly across updates and reinstalls
