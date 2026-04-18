# go-conventions

A rule package that installs Go coding conventions as always-on context.
Activates only when working on `.go` files (path-scoped).

## What it covers

- Naming conventions (MixedCaps, acronyms, interfaces)
- Error handling patterns (wrapping, checking, errors.Is/As)
- Package structure (internal/, minimal main.go)
- Testing conventions (table-driven, t.Helper, t.TempDir)
- Dependency management

## Install

```bash
skr install go-conventions
```

Installs to `~/.claude/rules/go-conventions.md`.
