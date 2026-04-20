# contribute

A skill that helps you contribute local Claude Code configuration to the
skill registry — either publishing new packages or upstreaming modifications
to existing ones.

## Usage

```
/contribute
/contribute ~/.claude/skills/my-skill/
/contribute go-conventions
```

## What it does

### No argument — scan and suggest

Scans your `~/.claude/` configuration for artifacts that could be shared:
- Local skills and rules not in the registry
- skr-installed packages you've modified locally

### New package flow

Generates `metadata.json`, the artifact file, and `README.md` from your
local file. Shows everything for review, then opens a PR against the
registry.

### Modified package flow

Detects drift in an installed package, diffs against the registry version,
and opens a PR with just the artifact change.

## Requirements

- `gh` CLI installed and authenticated
- `git` available

## System skill

This skill is automatically installed on first `skr install` as part of the
registry core utilities.
