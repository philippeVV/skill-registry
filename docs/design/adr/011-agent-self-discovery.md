# ADR-011: Agent Self-Discovery

**Date:** 2026-04-18
**Status:** Accepted

## Context

The registry needs a way for Claude Code agents to discover and install
packages at runtime — not just humans browsing the web UI. Two mechanisms
are available: a bundled skill that reads the index and recommends installs,
and a registry-as-MCP-server that exposes search and install as tools.

## Decision

Self-discovery is handled entirely by the `suggest-packages` helper skill
(M3). The skill shells out to `skr search` and `skr info`, reads current
repo context, and recommends packages. The agent runs `skr install <name>`
via bash to install.

Claude Code agents have native bash access, so `skr` via bash already
provides full runtime discovery capability. A registry MCP server was
considered but dropped — it would add infrastructure and complexity without
adding capability beyond what `skr` already provides.

MCP will be revisited if a concrete use case emerges that `skr` cannot serve
(e.g. non-bash contexts, cross-tool consumers).

## Consequences

**Positive:**
- No MCP infrastructure to build or host
- `skr` is the single interface for both humans and agents
- `suggest-packages` skill ships in M3 with no additional dependencies

**Negative:**
- Agents must have bash access — pure API contexts cannot use the registry
- Discovery is recommendation-only; agent proposes, user approves install

**Neutral:**
- The `suggest-packages` skill is itself a package in the registry —
  it eats its own dog food and serves as a reference implementation
