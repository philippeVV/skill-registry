# Skill Registry

A functional OSS registry for Claude Code config artifacts. Publish, discover,
and install skills, rules, and knowledge packages.

## Quickstart

```bash
# Install the CLI
go install github.com/philippeVV/skill-registry/skr@latest

# Search for packages
skr search code-review

# Install a package
skr install suggest-packages

# List installed packages
skr list
```

## How It Works

1. **Authors** submit packages via pull request to `packages/<name>/`
2. **CI** validates the package (schema + semantic rules) on every PR
3. **On merge**, CI tags a version and rebuilds the index (`marketplace.json`)
4. **Users** install with `skr install <name>` — artifacts are placed in
   the correct Claude Code location based on package type

## Package Types

| Type | What it is | Where it goes |
|------|-----------|---------------|
| **skill** | Slash commands and behaviors for Claude Code | `~/.claude/skills/<name>/` |
| **rule** | Always-on instructions (path-scoped optional) | `~/.claude/rules/<name>.md` |
| **knowledge** | Domain knowledge retrieved on demand | `~/.claude/knowledge/<name>/` |

## Contributing

Each package needs:
- `metadata.json` with name, type, description, tags, author, license
- An artifact file matching the type (`SKILL.md`, `RULE.md`, or `KNOWLEDGE.md`)
- A `README.md`

Validate locally before submitting:

```bash
python ci/validate.py packages/<your-package>
```

See `docs/VISION.md` for the full project vision and `docs/ROADMAP.md` for
the milestone plan.
