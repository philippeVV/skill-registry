# ADR-018: `skr` Configuration

**Date:** 2026-04-18
**Status:** Accepted

## Context

`skr` needs to know where the registry index lives and where to send OTEL
events. These values differ between personal OSS use and enterprise deployments.
We need a configuration mechanism that works for both without requiring a
`skr config` command.

## Decision

Configuration lives in `~/.config/skr/config.toml`. Environment variables
override file values. Defaults are baked into the binary.

**Default config (baked into binary):**
```toml
registry = "https://github.com/org/skill-registry"
otel_endpoint = ""   # empty = telemetry disabled
```

**Example enterprise config:**
```toml
registry = "https://registry.internal.yourorg.com/marketplace.json"
otel_endpoint = "https://otel.internal.yourorg.com"
```

**Environment variable overrides:**
- `SKR_REGISTRY` — overrides `registry`
- `OTEL_EXPORTER_OTLP_ENDPOINT` — standard OTEL env var, overrides
  `otel_endpoint`

For enterprise onboarding, the org distributes `~/.config/skr/config.toml`
via their standard dotfile or onboarding tooling (same pattern as `.gitconfig`
or `~/.npmrc`). No `skr config` command is provided in v1 — editing the
TOML directly is sufficient.

The exact shape of the config file may evolve as new configurable options
are identified. The TOML format is chosen for human editability.

## Consequences

**Positive:**
- Zero friction for personal use — defaults work out of the box
- Enterprise configuration is a single file, distributable via existing
  onboarding tooling
- Standard OTEL env var means `skr` composes naturally with existing
  instrumented environments
- No `skr config` command to maintain

**Negative:**
- Users must hand-edit TOML to change settings — acceptable for v1
- Config file location is Unix-centric (`~/.config/`) — Windows support
  would need a separate path convention

**Neutral:**
- Config options will grow as features are added; TOML handles this cleanly
