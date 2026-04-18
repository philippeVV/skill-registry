# M2 Spec — Trust Signals

## Goal

The registry starts accumulating evidence. `skr` emits OTEL events on
install and update. `skr info` shows rich package detail including live
stats when available. `skr update` keeps installed packages current.
`skr list` shows drift status via hash comparison.

Exit criterion: Install and invocation counts are verifiable via OTEL
console output. `skr info` shows full package detail usable by both humans
and Claude (via the `suggest-packages` skill).

---

## Task List

### 1. OTEL Instrumentation in `skr`

Add the OTEL Go SDK to `skr/go.mod`.

**Backend selection via config:**
- `otel_endpoint` empty (default) → console exporter: events printed to
  stdout. Sufficient for POC verification.
- `otel_endpoint` set → OTLP exporter: events sent to configured endpoint.

**Events to emit:**

| Event | Trigger | Attributes |
|---|---|---|
| `skr.package.install` | `skr install` completes | `package.name`, `package.version`, `package.type`, `registry.url` |
| `skr.package.uninstall` | `skr uninstall` completes | `package.name`, `package.version` |
| `skr.package.update` | `skr update` installs a new version | `package.name`, `package.version_from`, `package.version_to` |

Telemetry fires after the operation succeeds — never on failure.

---

### 2. Verify Claude Code Invocation Events (Research Task)

`claude_code.skill_activated` is the expected OTEL event from Claude Code's
own pipeline when a skill is invoked. This must be verified manually:

1. Enable OTEL on Claude Code: `CLAUDE_CODE_ENABLE_TELEMETRY=1` with a
   console or OTLP exporter configured
2. Install a skill and invoke it in a session
3. Confirm `claude_code.skill_activated` fires with attributes including
   `skill.name` and `skill.source`
4. Document findings in `docs/invocation-tracking.md`

If the event exists: invocation count is derivable from the OTEL pipeline
without any changes to `skr`. Document the attribute names and how to
query them.

If the event does not exist: document what is available and what the
fallback approach is (e.g. session-level tracking via `Stop` hook).

---

### 3. `skr update [<name>]`

Updates one or all installed packages to the latest version in the index.

**Behavior:**
- Force-refresh the index cache regardless of TTL
- For each package to update (all if no name given):
  - Compare installed version (from lockfile) against latest in index
  - If newer version available: reinstall at new version, update lockfile,
    emit `skr.package.update` event
  - If already latest: print `<name> is up to date`
- Print summary of what was updated

---

### 4. `skr info <name>` — Rich Output

`skr info` is used by both humans and Claude (via `suggest-packages` skill
as a bridge before the MCP server exists in M7). Output must be
comprehensive enough for Claude to make informed recommendations.

**Output format — structured plain text:**

```
Name:        code-review-expert
Version:     1.2.0
Type:        skill
Author:      team-platform
License:     MIT
Tags:        code-review, quality, engineering

Description:
  Adds expert-level code review behavior to Claude, focusing on
  correctness, maintainability, and security.

Notes:
  Works best on tasks with clear review scope.

Install target: skills (default)

Trust signals:
  Install count:    142
  Invocation count: 891
  Eval score:       8.4 (12 runs)

  Stats unavailable — set otel_endpoint in config to see live counts.
  (shown when no endpoint configured)

README:
  [full README.md content]
```

Stats lines are omitted entirely when no endpoint is configured and no
cached stats exist — replaced by the single "Stats unavailable" note.
README content is always included.

---

### 5. `skr list` — Drift Detection

Already planned in M1. M2 adds the drift flag powered by hash comparison.

**Output:**
```
code-review-expert   1.2.0   skill      ✓
api-conventions      1.0.1   knowledge  [modified]
suggest-packages     1.1.0   skill      ✓
```

`[modified]` means `hash != registry_hash` in the lockfile — the installed
artifact has been changed locally since install.

---

## Key Decisions

| Decision | Choice |
|---|---|
| OTEL backend for POC | Console exporter (stdout) |
| OTEL backend for enterprise | OTLP endpoint via config |
| `skr diff` | Dropped — `skr list [modified]` is sufficient |
| `skr info` format | Structured plain text, no `--json` flag |
| `skr info` audience | Humans and Claude (bridge until MCP server) |
| Invocation count | Research task — verify `claude_code.skill_activated` |
