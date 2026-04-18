# M6 Spec — Enterprise / Stats Foundation

## Goal

Define the stats contract that `skr info` and the web UI depend on.
Implement the daily heartbeat event in `skr` so the telemetry pipeline is
proven end-to-end even before a backend exists.

The stats aggregation backend (receive, compile, serve) is v2 and out of
scope for this milestone.

Exit criterion: `skr` emits `skr.packages.heartbeat` daily with the list
of installed packages, verifiable via console exporter output.

---

## Stats Contract

Any backend that implements the stats URL must serve this JSON shape:

```json
{
  "packages": {
    "code-review-expert": {
      "install_count": 142,
      "invocation_count": 891,
      "eval_score": 8.4,
      "eval_runs": 12
    }
  }
}
```

`skr info` and the web UI `StatsBlock` island fetch this endpoint and
display the values. Fields are omitted (shown as `—`) when absent.

This contract is the interface. What serves it is an implementation detail
— Prometheus, a custom aggregator, a static JSON file, anything.

---

## Task List

### 1. Daily Heartbeat Event in `skr`

On the first `skr` invocation of each day, `skr` emits:

```
event: skr.packages.heartbeat
attributes:
  registry:  "https://github.com/org/skill-registry"
  packages:  ["code-review-expert@1.2.0", "suggest-packages@1.0.0"]
```

**Implementation:**
- Track last heartbeat timestamp in `~/.config/skr/state.json`
- On any `skr` command: check if 24h have elapsed since last heartbeat
- If yes: read installed packages from lockfile, emit event, update timestamp
- Emission uses the same OTEL pipeline as install/uninstall events — console
  exporter in dev, OTLP in enterprise

**Verification:** Run any `skr` command after clearing the timestamp,
confirm `skr.packages.heartbeat` appears in console output with correct
package list.

---

### 2. Stats Contract Documentation

Document `docs/stats-contract.md`:
- The JSON schema above
- How `skr info` and the web UI consume it (`PUBLIC_STATS_URL` / `otel_endpoint`)
- What a backend needs to implement to serve it
- Notes on what `install_count` means (unique heartbeats containing the
  package) vs `invocation_count` (Claude Code `skill_activated` events)

---

## What's Deferred (v2)

- Stats aggregation backend — receives OTEL events, compiles per-package
  counts, serves the stats contract JSON
- LLM review gate in the publish pipeline
- Prompt-injection scanning in CI
- Enterprise-specific auth, VPN access controls, branch protection rules
  (these are ops configuration, not code)

---

## Key Decisions

| Decision | Choice |
|---|---|
| Stats contract | Simple JSON per package — implementation-agnostic |
| Heartbeat cadence | Daily — first `skr` invocation of each day |
| Heartbeat content | Registry URL + list of `name@version` strings |
| Heartbeat state | `~/.config/skr/state.json` tracks last emission timestamp |
| Backend | v2 — not built in this milestone |
| LLM review gate | v2 |
| Security scanning | v2 |
