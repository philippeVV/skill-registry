# ADR-010: Frontend Stack — Astro Static Site

**Date:** 2026-04-18
**Status:** Accepted

## Context

The web UI is a vitrine — a discovery and excitement surface, not a product
feature. Its job is to let people browse packages, filter by tag, and get
interested in what the registry offers. It is expected to be lightly used
but worth investing in over time. Heavy framework overhead is not justified
at this stage.

## Decision

The frontend is a static Astro site with Tailwind CSS. At CI build time,
Astro fetches `marketplace.json` and generates static package listing and
detail pages. No server required — deployed to S3/CloudFront (v2)
or GitHub Pages (OSS).

The visual aesthetic is inspired by claude-code-templates but the stack is
not forked from it. Astro gives component-based structure without the runtime
JS overhead of a full React/Next.js app.

The stack is not considered critical and may be replaced as needs evolve.
Choosing Astro is a pragmatic default, not a strong commitment.

**LiteLLM note:**
The organization runs LiteLLM as an MCP gateway and AWS model gateway,
primarily for agents. Almost all users access Claude via Claude Code accounts
directly. Potential integration points with the registry (auth, audit,
agent-facing API) are not designed for now. Revisit when LiteLLM usage
patterns are better understood.

## Consequences

**Positive:**
- Zero server infrastructure for the frontend
- Fast static pages, good for a discovery surface
- Easy to evolve or replace — it's just a vitrine
- CI rebuild on index change keeps it always current

**Negative:**
- Live stats (install/invocation counts) require a client-side fetch at
  page load since they are not baked into the static build
- Astro is unfamiliar to the team — low risk given the UI's limited scope

**Neutral:**
- Stack can be replaced without touching registry mechanics
