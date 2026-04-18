# ADR-025: Go Module Structure

**Date:** 2026-04-18
**Status:** Accepted

## Context

The repo contains Go source code for `skr` and potentially future components
(MCP server in M7). We need to decide whether to use a single root Go module
or standalone modules per component.

## Decision

Each Go component is a standalone module with its own `go.mod`. For M1,
that means `skr/go.mod` only.

```
skr/
  go.mod    # module github.com/org/skill-registry/skr
  go.sum
  main.go
  ...
```

Future Go components (e.g. `mcp-server/`) get their own module when built.
No root-level `go.mod`.

## Consequences

**Positive:**
- `go install github.com/org/skill-registry/skr@latest` works cleanly
- Each component has independent dependencies and versioning
- No shared module coupling between unrelated components

**Negative:**
- Shared code between components requires a separate internal module or
  duplication (unlikely to be needed in the near term)

**Neutral:**
- Standard pattern for multi-component Go monorepos
