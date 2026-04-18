# ADR-005: Three-Layer Publish Gate Pipeline

**Date:** 2026-04-18
**Status:** Accepted

## Context

The PR is the publish gate — nothing enters the registry without a reviewed
merge. We need to define what checks run before a package is accepted, and
in what order, to balance thoroughness against friction.

## Decision

Every package PR passes through three layers in sequence before merge is
allowed:

**Layer 1 — Automated gates (blocking, runs on every PR)**
- Schema validation: package conforms to extended marketplace.json format
- Prompt-injection scanning: static analysis on package content for
  instruction hijacking patterns
- Conflict detection: new package does not conflict with existing packages
- These must all pass before layers 2 and 3 are triggered

**Layer 2 — LLM review (blocking, runs after Layer 1 passes)**
- Quality and completeness review of package content and documentation
- Redundancy check: does this meaningfully differ from existing packages?
- Security intent analysis: flags suspicious patterns that static analysis
  missed
- Produces a structured review comment on the PR for the human reviewer

**Layer 3 — Human review (blocking, final approval)**
- Human reviewer reads LLM review output alongside the diff
- Focused on intent, organizational fit, and judgment calls
- Merge = publish, no additional step required

**OSS vs. v2:**
The OSS version runs Layer 1 and Layer 3 only. Layer 2 (LLM review) is
a v2 feature, reflecting the higher trust requirements and larger
package volume of internal registries. The OSS pipeline is intentionally
lighter to reduce friction for open contributions.

## Consequences

**Positive:**
- Human review time is spent on judgment, not mechanical checks
- LLM review catches quality issues that static analysis misses
- Prompt-injection scanning addresses a real and documented threat for this
  artifact class
- Clear separation between OSS and v2 value

**Negative:**
- LLM review adds latency and API cost to every v2 PR
- False positives in automated gates will frustrate contributors
- Maintaining three gate layers adds pipeline complexity

**Neutral:**
- CI runs gates 1 and 2; gate 3 is a GitHub branch protection rule
  requiring at least one approved review before merge
