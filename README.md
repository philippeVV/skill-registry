# Skill Registry

**The package manager for Claude Code.**

Skills, rules, and knowledge packages — discovered, installed, and kept
current with a single CLI. No more copying prompt files from Slack or
guessing which config snippets actually work.

## Why

Claude Code is only as good as the context it runs with. Teams are already
sharing skills and prompt configs — but through Slack threads, gists, and
copy-paste. There's no way to find what exists, no way to know if it works,
and no way to keep it up to date.

Skill Registry fixes that. One registry, one CLI, one command to install
what you need.

## What You Get

**One-command install.** `skr install suggest-packages` puts the right
files in the right place. Skills, rules, and knowledge each have their own
install semantics — `skr` handles all of it.

**Search and discovery.** Browse packages by tag, search by name, or let
the `suggest-packages` skill analyze your project and recommend what to
install.

**Drift detection.** The lockfile tracks content hashes. If a package is
modified locally, `skr list` flags it. `skr contribute` sends your improvements back upstream.

**Trust signals.** Install and invocation counts tracked via OpenTelemetry.
See what people actually use, not just what exists.

**Web UI.** Browse the registry in a browser at
[philippevv.github.io/skill-registry](https://philippevv.github.io/skill-registry/).
Filter by tag, read READMEs, check stats — no CLI needed.

**CI-powered quality gate.** Every package enters through a pull request.
Schema validation, semantic checks, and human review before anything gets
published.

**Helper bundle.** First-party skills that use the registry itself:
conflict checking, package suggestions, and contributing back upstream.
Install them all with `skr install --tag registry-core`.

## Package Types

| Type | What it is | Example |
|------|-----------|---------|
| **skill** | Slash commands and behaviors | `/code-review-checklist` |
| **rule** | Always-on instructions (path-scoped) | Go conventions |
| **knowledge** | Domain context retrieved on demand | Domain glossary |

## Getting Started

See [docs/getting-started.md](docs/getting-started.md) for installation,
usage, and contributing instructions.

## Coming Soon

**External package tracking** — Point the registry at upstream sources.
A scheduled CI job watches for changes and opens sync PRs automatically.
Renovate-style freshness for your Claude Code configs.

**A/B eval system** — Test whether a package actually improves Claude's
output on your real tasks. Two sub-agents run in parallel (with and without
the package), a judge scores the difference. Evidence that accumulates
across users and tasks.

**Private registry support** — Same format, same CLI, pointed at your
org's private repo. Internal access controls, OTEL wired to your
observability stack, LLM review gate, and prompt-injection scanning.

**Stats backend** — Live install and invocation counts served from an
OTEL-powered backend. Leaderboards, author rankings, and real trust
signals on the web UI.

## License

MIT
