# Roadmap

## M1 — Working Registry

The OSS repository is a functional marketplace. A developer can publish,
find, and install a package. Nothing more.

- Repo structure scaffolded (`packages/`, `skr/`, `web/`, `ci/`, `docs/`)
- `metadata.json` schema defined and validated in CI
- CI generates `marketplace.json` from `packages/` on every merge
- CI-owned semver via conventional commits and GitHub tags
- `skr install <name>[@version]` — user-level install with type-specific
  placement
- `skr uninstall <name>`
- `skr list` — show installed packages
- `skr search <query>` — search index by name, tag, description
- `skr info <name>` — package detail
- Lockfile at `~/.config/skr/skr.lock` with hash-based drift detection
- `skr --registry <url>` flag for custom registry source
- A handful of real example packages covering each type

**Exit criteria:** `skr install suggest-packages` works against the OSS repo.

---

## M2 — Trust Signals

The registry starts accumulating evidence, not just presence.

- OTEL instrumentation in `skr` — `skr.package.install` and
  `skr.package.uninstall` events
- `skr info` displays install count pulled from OTEL backend
- Invocation count for skill-type packages via Claude Code's OTEL pipeline
  (verify `claude_code.skill_activated` event exists and wire it up)
- `skr update [<name>]` — update one or all installed packages
- `skr diff <name>` — compare local vs upstream
- `skr contribute <name>` — open PR upstream with local changes

**Exit criteria:** Install and invocation counts visible in `skr info`.

---

## M3 — Helper Bundle

The registry ships utility skills that use the registry itself. Self-referential,
and the best demonstration of what the platform can do.

- `suggest-packages` skill — fetches index via `skr search`, reads repo
  context, recommends installs
- `conflict-check` skill — analyses full loaded context for conflicts and
  overlap; user-invoked only (not wired into install)
- `contribute` skill — drafts new packages or upstreams local modifications
  as PRs (merges the originally separate `publish-skill` and `contribute`)
- `knowledge-retriever` skill — finds and loads relevant knowledge packages
  on demand via `~/.claude/knowledge/index.json`
- `code-review-checklist` skill — structured code review guidance
- `eval` skill — deferred to v2 (ADR-034)
- All published under the `registry-core` tag
- `skr install --tag registry-core` installs the bundle in one command
- Auto-bootstrapped as system skills on first `skr install` of any package

**Exit criteria:** A user can install the bundle and run `/conflict-check`
and `/contribute`.

---

## M4 — Web UI

The vitrine. Discovery and excitement surface for the registry.

- Astro static site generated from `marketplace.json` at CI build time
- Package listing with tag filtering and search
- Package detail page — description, type, tags, trust signals, README
- Leaderboard — top packages by invocation/install count, top authors
- Live stats fetched client-side from OTEL backend
- CI deploys on every index rebuild

**Exit criteria:** Someone can browse the registry in a browser and find
something worth installing.

---

## M5 — External Package Tracking

The registry stays current with external sources without manual monitoring.

- `upstream` field in `metadata.json` schema
- Scheduled CI job diffs tracked packages against their upstream source
- Auto-opens PRs when upstream changes; goes through normal review pipeline
- Contributor guide updated with external package import workflow

**Exit criteria:** A tracked external package auto-generates a sync PR when
its upstream changes.

---

## M6 — v2

The OSS registry becomes a private internal registry with organizational
controls.

- Private repo configuration — `SKR_REGISTRY` and `config.json` for
  internal endpoint
- OTEL backend wired to existing org observability stack
- LLM review gate in CI publish pipeline (Layer 2, ADR-005)
- Human review enforced via GitHub branch protection rules
- Prompt-injection scanning in CI (Layer 1 security gate, ADR-016)
- Quality analysis beyond schema conformance in LLM review

**Exit criteria:** Internal registry running behind VPN with full publish
gate and live stats in org's observability stack.

---

## M6 — v2 / Stats Foundation

Define the stats contract. Implement daily heartbeat in `skr`. Stats
aggregation backend and full v2 controls are deferred.

- Stats JSON contract defined (`docs/stats-contract.md`)
- `skr.packages.heartbeat` daily event with installed package list
- Stats backend (receive, aggregate, serve) — v2
- LLM review gate, prompt-injection scanning — v2

**Exit criteria:** `skr` emits `skr.packages.heartbeat` daily, verifiable
via console exporter.

---

> **Note — MCP server (previously M7):** Dropped. `skr` via bash already
> gives Claude Code agents everything needed for runtime discovery.
> The `suggest-packages` helper skill (M3) covers the self-discovery use case.
> MCP will be revisited if a concrete use case emerges that `skr` can't serve.
