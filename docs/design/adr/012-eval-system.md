# ADR-012: Eval System

**Date:** 2026-04-18
**Status:** Accepted

## Context

We want users to be able to measure whether a package actually improves
Claude's output on their real work. The eval needs to run with and without
the package in context, compare results, and give the user a concrete signal.
It must not interfere with real work, especially for tasks involving file
operations.

## Decision

The eval system is a **skill** (`eval`) that is itself a package in the
registry — not a CLI feature. The user invokes it before a task they believe
the installed package would improve.

**Flow:**
1. User installs the `eval` skill from the registry
2. Before a task, user invokes `/eval` and describes what they're about to do
3. The eval skill orchestrates:
   - **Shadow sub-agent** — runs the same task without the package in context
   - **Judge sub-agent** — compares the two outputs and scores which performed better
4. Results are surfaced to the user at the end. User can optionally publish
   the result back to the registry as an evidence point for that package.

**Isolation for file operations:**
When the task involves file operations, the shadow sub-agent works in a git
worktree — an isolated working copy of the current branch. Real work happens
in the main tree. No collision, no need to copy the full filesystem.

**Package-level eval guidelines:**
Each package may include an `eval/guidelines.md` file. The eval skill reads
this before orchestrating the comparison — it tells the eval skill what to
watch for, what constitutes a meaningful improvement for this specific package,
and how to frame the judge's comparison. This makes evals more precise than
a generic side-by-side.

**Publishing results:**
Eval results (score, summary, package version, timestamp) can be submitted
back to the registry as an evidence point. Over time, a package accumulates
results across diverse real tasks from multiple users — a distribution more
meaningful than any synthetic benchmark.

## Consequences

**Positive:**
- Eval is composable — it's a package like any other, installed and invoked
  the same way
- Real work is the benchmark — no synthetic test maintenance
- Git worktree isolation is clean for file operations without full filesystem
  duplication
- Per-package eval guidelines make comparisons more meaningful and directed
- Results aggregate into registry-level evidence over time

**Negative:**
- Running three concurrent agents (main, shadow, judge) has real API cost
- Worktree creation adds latency before the task starts
- Judge quality depends on the judge prompt — needs careful design
- Not all tasks are eval-able (e.g. tasks that require external API calls
  or have non-deterministic side effects)

**Neutral:**
- The `eval` skill eats its own dog food — it is a reference implementation
  of what a sophisticated skill looks like
- `eval/guidelines.md` is optional — packages without it get a generic
  comparison
