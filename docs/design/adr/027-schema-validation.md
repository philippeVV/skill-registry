# ADR-027: Schema Validation

**Date:** 2026-04-18
**Status:** Accepted

## Context

Every PR that adds or modifies a package must be validated before human
review. Validation needs to catch both structural errors (wrong field types,
missing required fields) and semantic errors (artifact file missing,
name doesn't match directory).

## Decision

Validation runs in `ci/validate.py` as part of `pr.yml`. Two layers in one
script:

**Structural — JSON Schema:**
- `metadata.json` is validated against a JSON Schema file at
  `ci/schema/metadata.schema.json`
- Uses the `jsonschema` Python library
- Catches missing required fields, wrong types, invalid formats

**Semantic — custom logic:**
- Package directory name matches `metadata.json` `name` field
- Artifact file exists and matches the declared `type`
  (e.g. `type: skill` → `SKILL.md` must be present)
- Tags are non-empty strings
- `upstream.url` is a valid URL if `upstream` is present

Both layers run in a single CI step. Errors are reported with clear messages
pointing to the specific field or file that failed.

## Consequences

**Positive:**
- JSON Schema is declarative and easy to extend as the schema evolves
- Semantic rules catch errors that JSON Schema cannot express
- Single script, single CI step — simple for contributors to understand
  and reproduce locally

**Negative:**
- `jsonschema` is a dependency — needs to be installed in the CI environment
- Semantic rules must be kept in sync with schema evolution manually

**Neutral:**
- Contributors can run `python ci/validate.py packages/<name>` locally
  before submitting a PR
