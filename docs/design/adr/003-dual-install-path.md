# ADR-003: Dual Install Path — Native Plugin for Compatibility, Own CLI as Primary

**Date:** 2026-04-18
**Status:** Accepted

## Context

The registry index is a valid agentskills.io `marketplace.json` manifest,
meaning Claude Code's native `/plugin marketplace add <url>` + `/plugin install
<name>` works out of the box. However, native plugin install gives us no
control over the install process: no conflict detection, no type-specific
install semantics, no install count tracking, no security policy enforcement,
and no version pinning beyond what the standard supports.

## Decision

We support two install paths with distinct roles:

**Native `/plugin` path (compatibility layer)**
- The registry index is a fully valid `marketplace.json`
- Users can do `/plugin marketplace add <registry-url>` and `/plugin install
  <name>` with no additional tooling
- This path is a fallback and onboarding convenience, not the recommended flow
- No registry-specific features are available on this path

**Registry CLI (`skr`) — primary interface**
- `skr install <name>[@version]` is the recommended install command
- The CLI clones the registry repo (shallow, depth=1) and reads the
  `marketplace.json` index and package files from the local clone
- Handles type-specific install semantics transparently:
  - `skill` → copy directory to `~/.claude/skills/<name>/`
  - `rule` → copy `RULE.md` to `~/.claude/rules/<name>.md`
  - `knowledge` → copy `KNOWLEDGE.md` to `~/.claude/knowledge/<name>/`
    and update `~/.claude/knowledge/index.json`
- Emits OTEL install events for telemetry (when configured)
- Supports version pinning and lockfile generation with hash-based
  drift detection
- On first install, auto-bootstraps `registry-core` tagged system skills

Both paths consume the same index. There is no schema divergence.

## Consequences

**Positive:**
- Zero-friction entry point via native `/plugin` — no CLI required to start
- Full control over install semantics and policy enforcement via CLI
- CLI can evolve independently of Claude Code's plugin system
- Lockfile support enables reproducible installs for teams

**Negative:**
- Two paths to maintain and document
- Users who use native `/plugin` bypass conflict detection and telemetry
- Install counts will undercount if users prefer the native path

**Neutral:**
- CLI name `skr` (skill registry) — short, typeable, distinct from `skill`
  which is overloaded
