# ADR-022: External Package Tracking via Renovate-Style Flow

**Date:** 2026-04-18
**Status:** Accepted

## Context

Teams will find useful packages from external sources (public GitHub repos,
other registries). We need to decide how external packages enter the internal
registry and how they stay current with upstream changes.

Quarantine (a separate holding area for unvetted external packages) was
considered but rejected — the existing pipeline (schema validation + human
review) is sufficient as a trust gate. The gap is staying current with
upstream changes after initial import.

## Decision

**Importing an external package:**
A team submits a PR adding the external package under `packages/<name>/`
like any other contribution. The `metadata.json` includes an optional
`upstream` field pointing to the source:

```json
{
  "upstream": {
    "url": "https://github.com/external-org/their-skills",
    "path": "packages/useful-skill",
    "ref": "main"
  }
}
```

The PR goes through the normal review pipeline (ADR-005). Human review is
responsible for assessing the external source. No special quarantine step.

**Staying current — Renovate-style tracking:**
A scheduled CI job runs periodically against all packages that have an
`upstream` field. For each, it:
1. Fetches the upstream source at the declared `ref`
2. Diffs it against the locally committed artifact
3. If upstream changed, opens a PR in the registry with the diff

The PR title makes the source clear: `chore: sync code-review-expert from
upstream (external-org/their-skills)`. It goes through the normal review
pipeline before merging — a human approves the upstream changes before they
reach the registry. The trust boundary is never bypassed.

If the local version has diverged (via `skr contribute` or direct edits),
the upstream sync PR will show conflicts — the reviewer decides which wins.

## Consequences

**Positive:**
- External packages stay current without manual monitoring
- Trust boundary is fully preserved — every upstream change requires human
  approval before it lands
- `upstream` field makes provenance explicit and auditable
- Same PR pipeline as any other change — no special cases

**Negative:**
- Renovate-style job needs access to external URLs — may require firewall
  exceptions in air-gapped environments
- Upstream sync PRs can pile up if many packages track fast-moving external
  sources
- Conflicts between local edits and upstream require manual resolution

**Neutral:**
- Packages without an `upstream` field are unaffected
- The sync job cadence (daily, weekly) is configurable
