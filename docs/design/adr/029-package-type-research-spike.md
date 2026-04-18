# ADR-029: Package Type Finalization as M1 Prerequisite

**Date:** 2026-04-18
**Status:** Accepted

## Context

The type taxonomy (`skill`, `knowledge`, `fragment`, `config`) was defined
early in the design process before fully understanding how Claude Code
actually loads and uses different kinds of context. Types drive install
placement logic, schema validation, example packages, and CLI behavior.
Building on provisional types risks rework.

## Decision

"Research and finalize package types" is the first task of M1, before any
CLI placement logic, schema code, or example packages are written.

The research spike produces a short document at `docs/package-types.md`
covering:
- What context sources Claude Code actually loads and how (skills, CLAUDE.md,
  settings.json, plugins, hooks — what's real vs. assumed)
- Which of these map cleanly to installable package types
- Install semantics for each finalized type
- Evidence from Claude Code's actual behavior and documentation

All M1 work that depends on types (placement logic, schema validator, example
packages) is blocked until this document is accepted.

Example packages will cover: one per finalized type, the registry helper
bundle skills, and one package with an `upstream` field to validate the
schema even if the Renovate job isn't built until M5.

## Consequences

**Positive:**
- Prevents building on assumptions that don't match Claude Code's reality
- One focused research task unblocks everything else cleanly
- The output document serves as contributor reference going forward

**Negative:**
- Delays the start of implementation by one research task
- Types may still evolve after M1 as real usage reveals gaps

**Neutral:**
- The `skill` type is already well-understood and can be developed in
  parallel with the research spike
