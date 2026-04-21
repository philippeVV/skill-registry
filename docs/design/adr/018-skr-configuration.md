# ADR-018: `skr` Configuration

**Date:** 2026-04-18
**Status:** Accepted

## Context

`skr` needs to know where the registry index lives and where to send OTEL
events. These values differ between personal OSS use and v2 deployments.
We need a configuration mechanism that works for both without requiring a
`skr config` command.

## Decision

Configuration lives in `~/.config/skr/config.json`. Environment variables
override file values. Defaults are baked into the binary.

**Default config (baked into binary):**
```json
{
  "registry": "https://github.com/philippeVV/skill-registry",
  "claude_config_dir": "~/.claude",
  "otel_endpoint": ""
}
```

- `registry` — URL of the Git repository serving as the marketplace
- `claude_config_dir` — root directory for Claude Code config (where
  skills/, rules/, knowledge/ live). Tilde expansion is handled by `skr`.
- `otel_endpoint` — OTLP HTTP endpoint for telemetry. Empty = disabled.

**Example v2 config:**
```json
{
  "registry": "https://github.com/internal-org/skill-registry",
  "otel_endpoint": "otel.internal.yourorg.com"
}
```

**Environment variable overrides:**
- `SKR_REGISTRY` — overrides `registry`
- `SKR_CLAUDE_CONFIG_DIR` — overrides `claude_config_dir`
- `SKR_OTEL_ENDPOINT` — overrides `otel_endpoint`

**Why `SKR_*` env vars instead of standard `OTEL_*`:** The original design
used `OTEL_EXPORTER_OTLP_ENDPOINT`. During implementation this was changed
to `SKR_OTEL_ENDPOINT` to avoid conflicting with other OTEL-instrumented
tools on the same machine. A developer running `skr` alongside other
instrumented services would otherwise have their OTEL endpoint shared
across all tools unintentionally.

**Why JSON instead of TOML:** The original design specified TOML for human
editability. JSON was chosen during implementation because Go's standard
library handles JSON natively with no external dependency. Given that most
users will never hand-edit this file (defaults work out of the box, env
vars handle overrides), the readability trade-off was acceptable.

For v2 onboarding, the org distributes `~/.config/skr/config.json`
via their standard dotfile or onboarding tooling (same pattern as `.gitconfig`
or `~/.npmrc`). No `skr config` command is provided in v1 — editing the
JSON directly is sufficient.

## Consequences

**Positive:**
- Zero friction for personal use — defaults work out of the box
- v2 configuration is a single file, distributable via existing
  onboarding tooling
- `SKR_*` env vars avoid conflicts with other OTEL-instrumented tools
- No `skr config` command to maintain

**Negative:**
- Users must hand-edit JSON to change settings — acceptable for v1
- Config file location is Unix-centric (`~/.config/`) — Windows support
  would need a separate path convention
- JSON is less human-friendly than TOML for hand-editing, but most users
  won't need to — defaults and env vars cover the common cases

**Neutral:**
- Config options will grow as features are added; JSON handles this cleanly
