# Stats API Contract

The web UI fetches trust signals from an optional stats backend
configured via the `PUBLIC_STATS_URL` environment variable. When the
variable is empty or unset, all stats display as `—` placeholders.

---

## Endpoints

### `GET {STATS_URL}/packages/{name}`

Returns stats for a single package.

**Response (200):**

```json
{
  "installs": 142,
  "invocations": 891
}
```

**Response (404):** Package not found. Web UI treats this the same as
no stats configured (displays `—`).

---

### `GET {STATS_URL}/summary`

Returns stats for all packages. Used by the leaderboard page.

**Response (200):**

```json
[
  { "name": "code-review-checklist", "installs": 142, "invocations": 891 },
  { "name": "suggest-packages", "installs": 87, "invocations": 2034 }
]
```

Packages with zero activity may be omitted from the array.

---

## CORS

The stats backend must allow requests from the GitHub Pages origin:

```
Access-Control-Allow-Origin: https://philippevv.github.io
```

For local development, also allow `http://localhost:4321`.

---

## Error Handling

The web UI treats any non-200 response or network error the same as
"no stats available" — it renders `—` placeholders. There is no retry
logic; the page shows whatever it gets on the first attempt.

---

## Notes

- The OTEL backend (M2/M6) is the expected data source. This contract
  defines what the web UI expects, not how the backend aggregates data.
- Field values are integers (counts), not floats.
- The `eval_score` field from the M4 spec is deferred until the eval
  system is implemented (see ADR-034).
