# ADR-017: CI-Owned Semver via GitHub Tags

**Date:** 2026-04-18
**Status:** Accepted

## Context

Package versions need to be tracked in the index so `skr` can pin, update,
and detect drift. Having contributors manually write a version in `metadata.json`
is an anti-pattern: it causes merge conflicts, forces mechanical decisions on
authors, and pollutes PRs with version bumps that carry no information.

## Decision

Versions are owned by CI, not contributors. `metadata.json` contains no
`version` field.

**Mechanism:**
- Contributors commit changes to `packages/<name>/` with no version bump
- CI uses conventional commits to determine the semver increment:
  - `fix:` → patch
  - `feat:` → minor
  - `feat!:` or `BREAKING CHANGE:` → major
- On merge to main, CI creates a GitHub tag: `<name>@<semver>`
  (e.g. `code-review-expert@1.2.0`)
- The index rebuild (ADR-006) reads the latest tag per package and writes
  the resolved version into `marketplace.json`

**`skr` behavior:**
- `skr install <name>` installs the version currently in the index (latest)
- `skr install <name>@1.1.0` installs a specific tagged version
- `skr update` resolves the latest tag per installed package from the index
- The lockfile records the installed version (derived from the tag at install
  time, not from `metadata.json`)

**Initial version:**
First publish of a package tags it as `1.0.0`. CI detects a new package
(no prior tag exists) and bootstraps at `1.0.0` regardless of commit type.

## Consequences

**Positive:**
- Contributors focus on content, not versioning mechanics
- No version-related merge conflicts
- Git tags are the authoritative version record — auditable, immutable
- Conventional commits enforce meaningful commit messages as a side effect

**Negative:**
- Requires conventional commit discipline from contributors
- CI needs write access to create tags
- First-time contributors need to learn conventional commit format

**Neutral:**
- Per-package tags (`name@version`) rather than repo-level tags keep
  version histories independent across packages in the monorepo
