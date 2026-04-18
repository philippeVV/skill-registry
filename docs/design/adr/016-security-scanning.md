# ADR-016: Security and Quality Scanning

**Date:** 2026-04-18
**Status:** Accepted

## Context

Package artifacts loaded into Claude's context represent a prompt-injection
attack surface. A malicious or poorly-written package could manipulate
agent behavior beyond its stated purpose. We need to decide how much to
invest in scanning at publish time for v1.

## Decision

**v1:** Schema validation and human review are the only gates. The human
reviewer is responsible for catching suspicious content. No automated
security scanning or quality analysis beyond schema conformance.

**v2 and beyond:** Add two-layer scanning to the publish gate:
- *Pattern-based scanner* (fast, deterministic): flags known injection
  patterns — "ignore previous instructions", unusual unicode, hidden
  content, base64-encoded payloads
- *LLM security review* (part of the existing Layer 2 LLM gate, ADR-005):
  reads the full artifact and assesses whether content attempts to
  manipulate Claude beyond its stated purpose

Quality analysis beyond schema conformance (usefulness, completeness,
redundancy) is also deferred to v2 and will be handled by the LLM review
layer.

## Consequences

**Positive:**
- v1 ships without security infrastructure overhead
- Human review provides a baseline trust signal sufficient for early usage
- v2 design is clear and builds on existing gate architecture

**Negative:**
- v1 relies entirely on human reviewers catching malicious content
- Internal trust boundary (ADR-001 context) reduces but does not eliminate
  the risk — insiders can still publish bad packages

**Neutral:**
- The internal-only nature of the registry significantly reduces the
  threat surface compared to a public marketplace
