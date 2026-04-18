---
paths:
  - "**/*.go"
---

# Go Conventions

## Naming
- Use MixedCaps, not underscores
- Acronyms are all-caps: `HTTPClient`, `userID`
- Interface names: single-method interfaces end in `-er` (`Reader`, `Writer`)
- Unexported types/functions start with lowercase

## Error Handling
- Always check errors — never use `_` for error returns
- Wrap errors with context: `fmt.Errorf("doing X: %w", err)`
- Return errors, don't panic (except truly unrecoverable cases)
- Use `errors.Is` and `errors.As` for error inspection

## Structure
- One package per directory
- Package name matches directory name (lowercase, no underscores)
- `internal/` for unexported packages
- Keep `main.go` minimal — delegate to a `cmd/` or `internal/` package

## Testing
- Test files: `*_test.go` in the same package
- Table-driven tests for multiple cases
- Use `t.Helper()` in test helpers
- Use `t.TempDir()` for filesystem tests (auto-cleanup)
- Use `testdata/` for fixture files

## Dependencies
- Prefer stdlib over external packages
- Run `go mod tidy` before committing
- Pin with go.sum, don't vendor unless required
