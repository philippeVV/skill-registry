# conflict-check

A skill that scans all loaded Claude Code context for conflicts and
redundant overlap between installed packages and user configuration.

## Usage

```
/conflict-check
```

## What it does

1. Inventories all context sources: skr packages, user rules, user skills,
   knowledge files, and CLAUDE.md files
2. Reads the full content of every artifact
3. Compares pairwise for contradictions (conflicts) and redundant coverage
   (overlaps)
4. Reports findings with quoted evidence and resolution options

## When to use

- After installing new packages with `skr install`
- After a bulk install with `skr install --tag <tag>`
- At session start if your configuration has grown large
- When you notice Claude giving inconsistent guidance

## System skill

This skill is automatically installed on first `skr install` as part of the
registry core utilities.
