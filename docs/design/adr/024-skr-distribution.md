# ADR-024: `skr` Distribution

**Date:** 2026-04-18
**Status:** Accepted

## Context

Developers need a way to install `skr` before they can use the registry.
The OSS version is a developer POC, not a consumer product, so distribution
complexity is not justified at this stage.

## Decision

`skr` is installed from source via `go install`:

```
go install github.com/org/skill-registry/skr@latest
```

No pre-built binaries, no release pipeline, no Homebrew formula for M1.
This is sufficient for the developer POC audience.

Binary releases and package manager distribution are deferred to when
there are actual external users who need them.

## Consequences

**Positive:**
- Zero distribution infrastructure to maintain
- Works immediately for any developer with Go installed
- No CI jobs needed for binary builds in M1

**Negative:**
- Requires Go toolchain — not suitable for non-developer users
- No pinned binary installs — `@latest` always pulls current main

**Neutral:**
- Enterprise distribution (managed install, pinned versions) is a future
  concern addressed when moving to M6
