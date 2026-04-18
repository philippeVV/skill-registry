# ADR-009: OpenTelemetry for Install and Invocation Stats

**Date:** 2026-04-18
**Status:** Accepted

## Context

Install count and invocation count are live signals that cannot live in the
static CI-generated index. We need a stats layer that works for both the OSS
version (a personal project / POC) and the enterprise version (self-hosted
at work with existing observability infrastructure).

## Decision

`skr` is an OpenTelemetry-instrumented client. It emits events on install and
(where detectable) on invocation. The operator configures an OTEL collector;
the CLI is agnostic to what backend receives the events.

**Events emitted by `skr`:**
- `skr.package.install` — fired on `skr install <name>`, attributes:
  `package.name`, `package.version`, `package.type`, `registry.url`
- `skr.package.uninstall` — fired on `skr uninstall <name>`

**Invocation events:**
- For `skill` type: Claude Code may emit `claude_code.skill_activated` via its
  own OTEL pipeline (to be verified against actual Claude Code docs). If it
  does, invocation count comes from there, not from `skr`.
- For other types: no invocation signal exists. Install count is the only
  metric for non-skill packages.

**Collector configuration:**
- Set via `OTEL_EXPORTER_OTLP_ENDPOINT` env var or `~/.config/skr/config.toml`
- If no collector is configured, `skr` runs silently with no telemetry
- Telemetry is opt-out, not opt-in, for enterprise deployments (operator
  configures the endpoint centrally via managed settings)

**OSS / personal use:**
The OSS version is a personal project and POC — not a public deployment.
Running a collector is optional. Stats are a non-requirement for OSS; the
architecture simply supports it for free because `skr` is already instrumented.

**Enterprise:**
Events flow into the existing OTEL stack (Datadog, Grafana, Prometheus, etc.).
No custom stats infrastructure required. The web UI queries the OTEL backend
for aggregated counts.

## Consequences

**Positive:**
- Zero custom stats infrastructure — reuses whatever the org already runs
- Standard OTEL means any backend works (Datadog, Grafana, Splunk, etc.)
- `skr` instrumentation is straightforward — standard OTEL SDK
- OSS version carries no operational burden for stats

**Negative:**
- Web UI needs to query the OTEL backend — coupling to whatever backend the
  org chose
- Invocation count for skills depends on Claude Code's own OTEL pipeline,
  which needs to be verified
- Non-skill types have no invocation signal, only install count

**Neutral:**
- OTEL endpoint configuration follows the standard `OTEL_*` env var convention
  so it composes naturally with existing instrumented services
