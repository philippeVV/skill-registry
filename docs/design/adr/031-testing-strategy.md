# ADR-031: Testing Strategy for `skr`

**Date:** 2026-04-18
**Status:** Accepted

## Context

The `skr` CLI has two categories of logic: index operations (parsing,
search, cache, version resolution) and filesystem operations (install,
uninstall, placement, lockfile). Both need reliable test coverage.

## Decision

**Unit tests** — for pure functions:
- Index parsing and search logic
- Version resolution from tags
- Lockfile read/write and drift detection
- Config file parsing

**E2E tests** — for filesystem interaction:
- A fake Claude config directory structure is created in a temp dir before
  each test and torn down after
- Tests exercise the full install/uninstall flow against this fake structure
  to verify correct file placement, CLAUDE.md fence block handling, and
  lockfile state
- `claude_config_dir` is pointed at the temp dir via config override
- Assertions check actual file contents and structure, not mocks

No mocking of the filesystem. No live registry calls in tests — the index
is loaded from a fixture file.

Standard Go `testing` package. No external test framework.

## Consequences

**Positive:**
- E2E tests catch real interaction bugs with Claude config structure
- Temp dir approach is clean, isolated, and fast
- No mock drift — tests break when real behavior breaks

**Negative:**
- E2E tests are slower than pure unit tests
- Fake config structure must be kept in sync with real Claude config
  layout as Claude Code evolves

**Neutral:**
- Fixture `marketplace.json` files serve as both test data and
  documentation of expected index shape
