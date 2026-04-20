# Invocation Tracking

How skill invocation counts are derived from Claude Code's OTEL pipeline.

## `claude_code.skill_activated` Event

Claude Code emits a `claude_code.skill_activated` OTEL **log event** each time
a skill is invoked during a session.

### Enabling Telemetry

```bash
export CLAUDE_CODE_ENABLE_TELEMETRY=1
export OTEL_LOGS_EXPORTER=otlp          # or "console" for local debugging
export OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318
```

### Event Attributes

| Attribute | Description |
|---|---|
| `skill.name` | Name of the invoked skill |
| `skill.source` | Origin: `"bundled"`, `"userSettings"`, `"projectSettings"`, `"plugin"` |
| `plugin.name` | (Optional) Name of owning plugin |
| `marketplace.name` | (Optional) Marketplace where plugin was installed |
| `session.id` | Unique session identifier |
| `prompt.id` | UUID linking events within a single user prompt |
| `event.timestamp` | ISO 8601 timestamp |
| `event.sequence` | Monotonic counter for ordering within a session |

### Deriving Invocation Counts

To compute invocation counts for a registry package, query OTEL log events
where:

```
event.name == "claude_code.skill_activated"
  AND skill.name == "<package-name>"
```

Group by `skill.name` to get per-skill counts. Filter by time range as needed.

For registry-specific counts, filter on `skill.source` — skills installed via
`skr` appear as `"userSettings"` or `"projectSettings"` depending on scope.

### Example: Sentry Query

When OTEL events are forwarded to Sentry, invocation counts can be queried via
Sentry's Discover feature or the API:

```
event_type:transaction
transaction:"claude_code.skill_activated"
tags[skill.name]:<package-name>
```

The exact query syntax depends on how the OTEL collector maps log events to
Sentry events.

## Relationship to `skr` Events

`skr` emits its own OTEL **span events** for install/uninstall/update
operations (`skr.package.install`, `skr.package.uninstall`,
`skr.package.update`). These are separate from Claude Code's invocation events.

| Signal | Source | Event |
|---|---|---|
| Install count | `skr` CLI | `skr.package.install` |
| Uninstall count | `skr` CLI | `skr.package.uninstall` |
| Invocation count | Claude Code | `claude_code.skill_activated` |

Both streams feed into the same OTEL backend when configured, enabling
`skr info` to display unified trust signals in a future milestone.
