---
name: knowledge-retriever
description: >
  Retrieve relevant domain knowledge from the local knowledge base when the
  current task could benefit from additional context. Consult the knowledge
  index to find and load relevant knowledge files.
user-invocable: false
---

# Knowledge Retriever

You have a local knowledge base at `~/.claude/knowledge/` managed by the
skill registry. It contains domain knowledge packages that are NOT loaded
automatically — you must retrieve them when relevant.

## How it works

1. **Read the knowledge index** at `~/.claude/knowledge/index.json`. This
   file contains a list of available knowledge packages with their name,
   description, and tags.

2. **Evaluate relevance.** Based on the current task, determine which
   knowledge entries (if any) would provide useful context. Consider:
   - Does the description match the domain of the current task?
   - Do the tags overlap with the technologies or concepts involved?
   - Would this knowledge help you make better decisions?

3. **Load relevant knowledge.** For each relevant entry, read the file at
   `~/.claude/knowledge/<name>/KNOWLEDGE.md`. The path field in the index
   is relative to the knowledge directory.

4. **Apply the knowledge.** Use the loaded content as context for your work.
   Do not quote it verbatim — integrate it into your understanding.

## When to trigger

- When working on domain-specific tasks (payments, auth, infrastructure)
- When you encounter unfamiliar concepts or conventions
- When the user asks about how something works in their system
- When making architectural decisions that benefit from domain context

## When NOT to trigger

- For simple, self-contained tasks (formatting, typo fixes)
- When the knowledge index is empty
- When you've already loaded the relevant knowledge in this session
