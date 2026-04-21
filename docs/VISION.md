# Skill Registry — Vision

## The Problem

Claude Code is only as good as the context it operates in. Skills, knowledge
snippets, CLAUDE.md fragments, and prompt configs are multiplying across teams
— passed around in Slack, copied from repos, pulled from unknown public
sources. There is no single place to find them, no way to know if they work,
and no control over what gets into your agents.

## The Idea

A private, internal registry for all Claude config artifacts. One place to
publish, discover, and install the building blocks that make Claude Code
useful in your organization — with trust signals so people know what actually
works.

The trust boundary is the core value. If it's in the registry, it was
reviewed by someone inside the org before it got there.

## What Lives Here

Every publishable unit is a **package**, subdivided by type:

- **skill** — slash commands and system prompt behaviors for Claude Code
  (directory installed to `~/.claude/skills/<name>/`)
- **rule** — always-on coding conventions and constraints
  (single file installed to `~/.claude/rules/<name>.md`)
- **knowledge** — domain knowledge snippets loaded on demand
  (directory installed to `~/.claude/knowledge/<name>/`)

These types were finalized via a research spike that mapped Claude Code's
actual context-loading behavior to installable types. Earlier drafts
included `fragment` (CLAUDE.md blocks) and `config` (settings.json
artifacts) — both were dropped because Claude Code's context model didn't
support clean install/uninstall for them. `hook` and `plugin` are
identified as future candidates (see `docs/package-types.md`).

MCP servers are out of scope. They involve actions, permissions, and auditing
that belong in a dedicated MCP gateway.

## How It Works

### Publishing

Packages live in Git. Authors submit via pull request. Every PR passes
through a three-layer gate before merge is allowed:

1. **Automated** — schema validation, prompt-injection scanning (v2),
   conflict detection
2. **LLM review** — quality, redundancy, and security intent analysis
   (v2)
3. **Human review** — final approval; merge is publish

Versions are owned by CI, not contributors. Authors never bump a version
field. On merge, CI derives the semver increment from conventional commits
and creates a Git tag (`package-name@1.2.0`). The index is regenerated
automatically from the current package files.

### The Registry is a Git Repo

The registry is a Git repository with a `marketplace.json` index at the
root. Claude Code's native `/plugin marketplace add <repo-url>` works
against it with no additional tooling. The `skr` CLI is the recommended
interface and provides the full feature set.

The OSS version is the repository itself — packages, index, CLI source,
and frontend in one place. For v2, the same format is hosted on a
private repository with internal access controls. Switching is a single
configuration change.

### Discovery

A static Astro frontend built from the index at CI time serves as the
discovery surface — a vitrine for browsing packages by tag, seeing trust
signals, and getting people interested. It is not a heavy product feature.

Tags are the primary organization primitive. Installing everything tagged
for a role or workflow is a one-command operation: `skr install --tag
frontend-dev`. No separate bundle concept needed.

Discovery also surfaces through the **`suggest-packages` helper skill** —
installed as part of the registry helper bundle, it fetches the index,
reads the current repo context, and recommends relevant packages to install.
This is the v1 agent self-discovery mechanism. A registry MCP server
enabling true runtime discovery is planned for v2.

### Installing

`skr` is the primary install interface, written in Go. It handles
type-specific install semantics transparently:

- `skill` → `~/.claude/skills/<name>/` (directory with SKILL.md + optional subdirs)
- `rule` → `~/.claude/rules/<name>.md` (single file, always loaded at startup)
- `knowledge` → `~/.claude/knowledge/<name>/` (directory with KNOWLEDGE.md,
  loaded on demand via `knowledge-retriever` skill)

The lockfile at `~/.config/skr/skr.lock` tracks every installed package
with its version, install location, and a content hash. If a package is
manually modified locally, `skr list` flags it as drifted. `skr diff
<name>` shows what changed. `skr contribute <name>` opens a PR upstream
with local improvements.

Packages default to user-level install. Native `/plugin install` also works
as a fallback for basic installs without the full `skr` feature set.

### External Packages

External packages enter the registry through the same PR process as any
other contribution. A `metadata.json` `upstream` field records the source.
A scheduled CI job tracks upstream changes and opens PRs automatically when
the source changes — a Renovate-style flow that keeps packages current while
preserving the human review gate.

## Trust Signals

The registry accumulates evidence over time, not just presence.

- **Invocation count** — for skill-type packages: how often the skill is
  actually used after install. The primary signal. Tracked via OpenTelemetry.
- **Install count** — for always-on types (rule, knowledge): the primary
  signal. Also tracked via OTEL.
- **Leaderboard** — authors ranked by aggregate invocation and install
  counts. Decorative and motivational, not a critical signal.
- **Evals** — see below.

Stats are live, served from the OTEL backend, and merged client-side by
`skr` and the web UI. They are never baked into the static index.

## Evals (Deferred to v2)

The eval system is designed but deferred to v2 (see ADR-033, ADR-034).

The design: a user installs the `eval` helper skill and invokes it before
a task they expect a package to improve. The eval skill orchestrates two
parallel sub-agents — one with the package loaded, one without — and a
judge sub-agent that compares outputs and scores which performed better.
For tasks involving file operations, the shadow sub-agent runs in a git
worktree to avoid interfering with real work.

Deferred because of unresolved design questions around output capture,
session orchestration, and whether users find the signal valuable enough
to justify the API cost of running parallel agents.

## Conflict and Overlap Detection

The **`conflict-check`** helper skill analyses the full loaded context from
a fresh session — reading complete artifact content, not summaries — to
identify duplicated information or contradictory instructions. It runs on
demand or is invoked by `skr` at install time as an advisory check.

Two failure modes:
- **Conflict** — two packages give contradictory instructions
- **Overlap** — two packages cover the same ground, burning context budget
  twice for no gain

Neither is a hard block. The user sees a warning and decides whether to
proceed.

## Registry Helper Bundle

The registry ships a first-party bundle of utility skills tagged
`registry-core`, installed once with `skr install --tag registry-core`
(or auto-bootstrapped on first `skr install` of any package):

- **`suggest-packages`** — browse the index and recommend installs for the
  current repo
- **`conflict-check`** — audit loaded context for conflicts and overlap
- **`contribute`** — draft new packages or upstream local modifications as
  PRs (merges the originally separate `publish-skill` and `contribute`
  concepts)
- **`knowledge-retriever`** — find and load relevant knowledge packages
  on demand
- **`code-review-checklist`** — structured code review guidance

These skills are themselves packages in the registry — they eat their own
dog food and serve as reference implementations.

## What This Is Not

- Not a public marketplace — this is internal, private, controlled
- Not an MCP gateway — actions and permissions live elsewhere
- Not a heavyweight platform — the OSS version runs with no infrastructure

## Open Source Strategy

The OSS version is the foundation: the registry repository itself, with
example packages, the `skr` CLI, and the Astro frontend. It is the
development and proof-of-concept environment.

The v2 version is the same codebase pointed at a private repository
with internal access controls and an OTEL collector wired in. The interface,
schema, and CLI are identical.

## Primary Target

Claude Code. Everything is designed around how Claude Code loads and uses
context. Other Claude surfaces are secondary.
