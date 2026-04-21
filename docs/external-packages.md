# External Package Tracking

Packages sourced from external GitHub repositories can be tracked so they
stay current automatically. A daily CI job checks upstream sources and opens
sync PRs when changes are detected.

---

## Importing an External Package

1. Add the package under `packages/<name>/` like any other contribution.
2. Include an `upstream` field in `metadata.json`:

```json
{
  "name": "useful-skill",
  "type": "skill",
  "upstream": {
    "url": "https://github.com/external-org/their-repo",
    "path": "packages/useful-skill",
    "ref": "main"
  }
}
```

- **url** — GitHub repository URL
- **path** — path to the artifact file or directory within the upstream repo
- **ref** — branch, tag, or commit SHA to track

3. The initial artifact content should match upstream at the time of
   submission.
4. Submit a PR. It goes through the normal review pipeline.

---

## How Sync Works

A scheduled GitHub Actions workflow (`renovate.yml`) runs daily at 06:00 UTC:

1. Reads `marketplace.json` to find packages with an `upstream` field.
2. For each tracked package, fetches the upstream content via the GitHub API.
3. Compares upstream file SHAs against local files.
4. If content differs, opens a sync PR (or updates an existing one).

Sync PRs go through the normal `pr.yml` validation pipeline and require
human review before merging. The trust boundary is never bypassed.

### What Gets Synced

All files from the upstream path **except `metadata.json`**, which is
registry-local and contains fields like `upstream`, `tags`, and `author`
that don't exist in the upstream source.

If upstream deletes a file, the sync PR removes the corresponding local
file.

---

## Handling Sync PRs

When a sync PR appears:

- **Review the diff carefully** — the changes come from an external source.
- **Merge** if the upstream changes are acceptable.
- **Close** if the local version has diverged intentionally and you want to
  keep the local version.

If an existing sync PR is still open when upstream changes again, the PR
branch is force-updated with the latest upstream content. GitHub shows a
force-push event in the PR timeline.

---

## Conflict Resolution

If the local artifact has been modified independently from upstream, the
sync PR will contain the full upstream version. The reviewer decides which
wins:

- **Accept upstream:** merge the sync PR.
- **Keep local:** close the sync PR.

There is no automated conflict resolution — it is always a human decision.

---

## Stopping Tracking

To stop tracking an upstream source, remove the `upstream` field from the
package's `metadata.json` and submit a PR. Future sync runs will skip that
package.

---

## Private Upstream Repos

The default `GITHUB_TOKEN` can only access public repositories. To track a
package from a private upstream repo, a Personal Access Token (PAT) with
`repo` scope must be configured as a repository secret.

---

## Repository Settings

Enable **"Automatically delete head branches"** in the repository settings.
This ensures sync branches (`sync/<name>`) are cleaned up after PRs are
merged, allowing the next sync run to create a fresh branch.
