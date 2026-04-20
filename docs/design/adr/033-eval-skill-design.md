# ADR-033: Eval Skill Design

**Date:** 2026-04-18
**Status:** Deferred to v2 (see ADR-034)

## Context

The eval skill A/B tests a package against real user work to produce a
concrete signal of whether the package helps. The design must balance
isolation quality against what Claude Code actually supports for sub-agent
context control.

## Decision

**Three participants — not three sub-agents:**

1. **Main session** — the user's actual Claude Code session, with the
   package loaded. This IS the "with package" execution. The user performs
   their real task here as normal.

2. **Shadow sub-agent** — spawned by the eval skill to perform the same
   task without the package in context. This is the "without package"
   execution.

3. **Judge sub-agent** — receives both outputs and scores which performed
   better, with reasoning.

**Shadow sub-agent isolation — intent:**
The shadow sub-agent should run without awareness of the package being
evaluated. The desired approach (subject to what Claude Code supports):
- Spawn with a clean or restricted system prompt that excludes the package
- Potentially deny read access to the package file on disk
- At minimum: explicitly instruct the sub-agent to ignore the package
  content (weaker but practical fallback)

**Package eval guidelines:**
Each package may include `eval/guidelines.md` describing what a meaningful
improvement looks like for that specific package. The eval skill reads this
before spawning the shadow and judge sub-agents to focus the comparison.

**File operation isolation:**
When the task involves file operations, the shadow sub-agent works in a
git worktree to avoid interfering with the user's real work.

## Open Questions (resolve during M3 implementation)

- Can Claude Code sub-agents be spawned with a restricted context that
  excludes specific loaded skills?
- Can file-level read access be denied for a sub-agent's scope?
- How does the main session output get captured and passed to the judge —
  does the eval skill observe the conversation, or does the user summarize?
- What is the judge's output format — score only, or score + reasoning?

## Consequences

**Positive:**
- Main session is the real "with package" signal — no synthetic reproduction
- Real work is the benchmark — no test maintenance
- Judge sub-agent keeps evaluation reasoning in Claude, not in code

**Negative:**
- Shadow sub-agent isolation is uncertain — weaker isolation means noisier
  signal
- Parallel agent execution has real API cost
- Implementation depends on Claude Code sub-agent capabilities that need
  to be verified during M3

**Neutral:**
- The eval skill is the most experimental part of the system — expect
  iteration on the design after first implementation
