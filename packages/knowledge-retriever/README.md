# knowledge-retriever

An always-on skill that retrieves relevant domain knowledge from the local
knowledge base managed by `skr`.

## How it works

Knowledge packages installed via `skr install` are stored at
`~/.claude/knowledge/<name>/` and indexed in `~/.claude/knowledge/index.json`.
This skill reads the index, evaluates relevance to the current task, and
loads the appropriate knowledge files on demand.

This implements a retrieval-augmented pattern: knowledge consumes zero context
budget at startup and is loaded only when relevant.

## System skill

This skill is automatically installed on first `skr install` as part of the
registry core utilities. It has `user-invocable: false` — Claude triggers
it automatically when domain knowledge would help.
