# ADR-032: OTEL Development Approach

**Date:** 2026-04-18
**Status:** Accepted

## Context

M2 adds OTEL instrumentation to `skr`. For the OSS POC, we don't need a
real OTEL backend — we just need to verify that the right events are being
emitted with the right attributes.

## Decision

`skr` uses the standard OTEL Go SDK. Backend selection is purely config:

- `otel_endpoint` empty (default) → console exporter: events printed to
  stdout. Used during development to verify emission.
- `otel_endpoint` set → OTLP exporter: events sent to the configured
  endpoint. Used in enterprise deployments.

No Docker Compose, no local Prometheus stack required for M2. Verification
is done by reading stdout.

**Events emitted by `skr`:**
- `skr.package.install` — attributes: `package.name`, `package.version`,
  `package.type`, `registry.url`
- `skr.package.uninstall` — attributes: `package.name`, `package.version`
- `skr.package.update` — attributes: `package.name`, `package.version_from`,
  `package.version_to`

**Invocation count via Claude Code:**
`claude_code.skill_activated` is the expected event from Claude Code's own
OTEL pipeline. Verification is a manual test: enable OTEL on Claude Code,
invoke a skill, confirm the event fires with expected attributes. This is
a research/verification task in M2, not a build task.

## Consequences

**Positive:**
- Zero backend infrastructure for POC verification
- Same `skr` binary works for dev (console) and enterprise (OTLP)
- Console exporter output is human-readable and easy to assert in tests

**Negative:**
- Console output is not queryable — no aggregation or dashboards in dev
- Invocation count verification depends on Claude Code's OTEL behavior
  which must be confirmed manually

**Neutral:**
- OTEL SDK adds a dependency to `skr/go.mod` but is well-maintained and
  standard
