# Skill Registry

OSS marketplace for Claude Code config artifacts (skills, rules, knowledge).
Developers publish packages via PR, CI validates and builds the index,
users install with `skr`.

## Structure

- `packages/` — registry content (one dir per package)
- `skr/` — Go CLI source (`github.com/philippeVV/skill-registry/skr`)
- `ci/` — Python scripts for validation and index building
- `docs/` — vision, roadmap, ADRs, specs
- `web/` — Astro frontend (M4, placeholder)
- `marketplace.json` — generated index (do not hand-edit)

## Dev Setup

- **Go 1.24** for the `skr` CLI
- **Python 3** for CI scripts

## Common Commands

```bash
# Validate a package
python ci/validate.py packages/<name>

# Validate all packages
python ci/validate.py

# Rebuild the index
python ci/build_index.py

# Run CLI tests
cd skr && go test ./...

# Build the CLI
cd skr && go build -o skr .
```

## Package Types

| Type | Artifact | Install target |
|------|----------|---------------|
| skill | SKILL.md (directory) | ~/.claude/skills/<name>/ |
| rule | RULE.md | ~/.claude/rules/<name>.md |
| knowledge | KNOWLEDGE.md | ~/.claude/knowledge/<name>/ |

See `docs/package-types.md` for full details.

## Contributing Packages

Each package lives at `packages/<name>/` with:
- `metadata.json` — name, type, description, tags, author, license
- Artifact file matching the type (SKILL.md, RULE.md, or KNOWLEDGE.md)
- `README.md`

Run `python ci/validate.py packages/<name>` before submitting a PR.
