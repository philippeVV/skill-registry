# ADR-021: Progressive Disclosure as a Per-Type Recommended Pattern

**Date:** 2026-04-18
**Status:** Accepted

## Context

Packages loaded at session start consume context tokens. A user with many
installed packages — especially always-on types like `rule` — can exhaust
useful context budget before any real work begins. The agentskills.io
standard addresses this with a three-tier loading model. We need to decide
whether to enforce this at the registry level or leave it to package authors.

## Decision

Progressive disclosure is a **recommended authoring pattern**, not a registry
requirement. It is per-package-type guidance, not a schema enforcement.

**Recommended tiers by type:**

`skill` — follows the agentskills.io three-tier model:
- Tier 1 (always loaded): SKILL.md frontmatter — name, description, trigger
  hints (~50-100 tokens)
- Tier 2 (loaded on activation): SKILL.md body — full instructions
- Tier 3 (loaded on demand): `references/` files — examples, extended docs

`rule` — always loaded at session start, so authors are responsible for
keeping content lean. The registry contributor guide recommends:
- Keep under 200 lines
- Lead with the most load-bearing instructions
- Use optional YAML frontmatter `paths` field for path scoping

`knowledge` — loaded on demand via the `knowledge-retriever` skill, so
context cost is only paid when relevant. Authors should still keep the
primary `KNOWLEDGE.md` focused and lead with load-bearing content.

The registry does not enforce token limits or tier structure at publish time.
This may be revisited if context bloat becomes an observed problem.

## Consequences

**Positive:**
- Packages that follow the pattern are good citizens in any context window
- Guidance lives in the contributor docs, not in CI — no false positives
- Consistent with the agentskills.io standard for skill types

**Negative:**
- Without enforcement, poorly-authored packages will bloat context
- No tooling to warn authors about token footprint at PR time (v2 concern)

**Neutral:**
- The conflict-check helper skill (ADR-013) partially addresses this by
  flagging overlap — duplicate context is a symptom of the same problem
