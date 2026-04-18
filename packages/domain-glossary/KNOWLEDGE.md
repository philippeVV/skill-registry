# Domain Glossary — Example

This is an example knowledge package demonstrating the knowledge type.
Replace this content with your organization's domain glossary.

## Terms

**Package** — A publishable unit in the skill registry. Contains an artifact
file (SKILL.md, RULE.md, or KNOWLEDGE.md), metadata.json, and README.md.

**Skill** — A Claude Code extension that provides behaviors, slash commands,
or automated workflows. Loaded progressively: description at startup, full
body on invocation.

**Rule** — Always-on instructions installed to `~/.claude/rules/`. Loaded at
session start. Can be path-scoped to activate only for matching files.

**Knowledge** — Domain-specific context stored in `~/.claude/knowledge/` and
retrieved on demand via the knowledge-retriever skill. Does not consume
context budget at startup.

**Registry** — The Git repository that serves as both the software (CLI, CI)
and the marketplace (packages, index). `marketplace.json` at the repo root
is the package index.

**Lockfile** — `~/.config/skr/skr.lock` tracks installed packages with
version, location, and content hash for drift detection.

**Drift** — When an installed package's content hash differs from its
registry hash, indicating local modifications.
