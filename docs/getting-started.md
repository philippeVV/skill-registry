# Getting Started

## Install the CLI

```bash
go install github.com/philippeVV/skill-registry/skr@latest
```

Requires Go 1.24+.

## Basic Usage

```bash
# Search for packages
skr search code-review

# View package details
skr info suggest-packages

# Install a package
skr install suggest-packages

# Install all registry helper skills
skr install --tag registry-core

# List installed packages
skr list

# Check for updates
skr update

# Compare local changes against upstream
skr diff suggest-packages

# Send local improvements back upstream
skr contribute suggest-packages
```

## How It Works

1. **Authors** submit packages via pull request to `packages/<name>/`
2. **CI** validates the package (schema + semantic rules) on every PR
3. **On merge**, CI tags a version and rebuilds the index (`marketplace.json`)
4. **Users** install with `skr install <name>` — artifacts go to the correct
   Claude Code location based on package type

### Install Locations

| Type | Install target |
|------|---------------|
| skill | `~/.claude/skills/<name>/` |
| rule | `~/.claude/rules/<name>.md` |
| knowledge | `~/.claude/knowledge/<name>/` |

## Contributing a Package

Each package lives at `packages/<name>/` with three files:

- **`metadata.json`** — name, type, description, tags, author, license
- **Artifact file** — `SKILL.md`, `RULE.md`, or `KNOWLEDGE.md` (matches the type)
- **`README.md`** — documentation

Validate locally before submitting a PR:

```bash
python ci/validate.py packages/<your-package>
```

Or use the `contribute` skill to automate the process:

```
/contribute
```

## Development

```bash
# Validate a single package
python ci/validate.py packages/<name>

# Validate all packages
python ci/validate.py

# Rebuild the index
python ci/build_index.py

# Run CLI tests
cd skr && go test ./...

# Build the CLI locally
cd skr && go build -o skr .
```

## Further Reading

- [docs/VISION.md](VISION.md) — full project vision
- [docs/ROADMAP.md](ROADMAP.md) — milestone plan
- [docs/package-types.md](package-types.md) — detailed type documentation
