# ADR-034: Defer Eval Skill to v2

**Date:** 2026-04-20
**Status:** Accepted
**Supersedes:** ADR-033 (Eval Skill Design) — status changed to Deferred

## Context

During M3 planning, the eval skill's workflow proved to have unresolved
design questions around output capture and session orchestration. The core
A/B testing concept is sound, but the implementation depends on Claude Code
sub-agent capabilities that need experimentation. Shipping the other M3
skills (conflict-check, contribute) should not be blocked by eval.

## Decision

Defer the eval skill to v2. Capture the design work done so far for future
implementation.

### Design decisions reached

**Three-participant architecture:**
1. **Main session** — the user's real Claude Code session with the package
   loaded. This IS the "with package" execution.
2. **"Without" sub-agent** — spawned with explicit instructions to ignore
   the package being evaluated. Performs the same task in a worktree.
3. **Judge sub-agent** — controls the "without" sub-agent, reads the main
   session history from `~/.claude/projects/<slug>/<session>.jsonl`, and
   compares both outputs.

**Isolation approach:** Explicit instruction (tell the sub-agent to ignore
the package). Claude Code does not currently support restricted context
exclusion or file-path denial for sub-agents. Explicit instruction is the
pragmatic v1 — weaker isolation but zero side effects.

**Judge as orchestrator:** The judge sub-agent controls the "without"
sub-agent rather than having the eval skill manage all three participants.
This keeps test outputs isolated from the main session.

**Per-package eval guidelines:** Each package may include
`eval/guidelines.md` describing what a meaningful improvement looks like.

### Open questions (to resolve during v2)

- **Timing:** When does the judge run? The user performs their real task in
  the main session — how do they signal completion? Options explored:
  explicit second invocation (`/eval-judge`), or the skill watches for a
  signal.
- **Output capture:** The judge reads session JSONL, but the format and
  reliability of extracting "what the agent did" from the transcript needs
  experimentation.
- **Cost:** Running parallel agents has real API cost. Need to understand
  whether users find the signal valuable enough to justify it.

## Consequences

**Positive:**
- M3 ships without being blocked by experimental eval design
- Design decisions are preserved for v2 implementation
- Other helper skills (conflict-check, contribute) provide immediate value

**Negative:**
- The original M3 exit criterion ("run `/eval` against a real task") is
  reduced. New criterion: install bundle and run `/conflict-check` and
  `/contribute`.
