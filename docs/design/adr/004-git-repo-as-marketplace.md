# ADR-004: Git Repo as the Marketplace

**Date:** 2026-04-18
**Status:** Accepted

## Context

In Claude Code's native plugin system, a marketplace is simply a Git
repository containing a `marketplace.json` manifest at a well-known path.
There is no server, no API, no special infrastructure. We need to decide
where the index and package artifacts live for both the OSS and v2
versions.

## Decision

The registry is a Git repository. The OSS version of this project IS the
marketplace — `marketplace.json` lives at the repo root, package artifacts
live under `packages/<name>/`. Users subscribe with:

```
/plugin marketplace add https://github.com/org/skill-registry
```

or via the CLI:

```
skr --registry https://github.com/org/skill-registry install <name>
```

The OSS repo serves three roles simultaneously:
1. Reference implementation of the registry software
2. Working marketplace with example packages
3. The default registry source for `skr`

For the v2 version, the same format is hosted on a private Git repo
or an S3-backed URL. The `skr` CLI accepts a `--registry <url>` flag (or a
config file entry) that overrides the default. No code changes required to
point at a private registry — only configuration.

## Consequences

**Positive:**
- Zero infrastructure for the OSS version — Git IS the distribution mechanism
- Native `/plugin marketplace add` works out of the box with no extra tooling
- PR-based publishing workflow maps naturally to Git
- Versioning, history, and rollback come free from Git
- Private v2 registry is purely a configuration change

**Negative:**
- Large binary assets in the repo could grow Git history over time — mitigated
  by keeping package artifacts as text (SKILL.md, markdown, JSON)
- Public OSS repo means all example packages are public — intentional for OSS,
  v2 uses private repo

**Neutral:**
- `skr` defaults to the OSS repo as its registry source; v2 teams
  override via `~/.config/skr/config.toml` or the `SKR_REGISTRY` env var
