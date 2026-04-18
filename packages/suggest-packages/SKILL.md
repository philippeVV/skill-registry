---
name: suggest-packages
description: >
  Analyze the current project and recommend relevant packages from the skill
  registry. Use when starting work on a new project, onboarding to a codebase,
  or when the user asks what packages might help.
user-invocable: true
argument-hint: "[focus area]"
---

# Suggest Packages

You have access to the skill registry via the `skr` CLI. Your job is to
analyze the current project and recommend packages that would be useful.

## How to discover packages

Use the `skr` CLI to explore the registry:

- `skr search <query>` — search by name, description, or tags
- `skr info <name>` — get full details on a package
- `skr list` — see what's already installed

## Process

1. **Understand the project context.** Read the project structure, look at
   file types, frameworks, languages, and existing CLAUDE.md or settings.
   Check what's already installed with `skr list`.

2. **Search the registry.** Run multiple `skr search` queries based on what
   you observe:
   - Search by language/framework (e.g., `skr search go`, `skr search react`)
   - Search by activity (e.g., `skr search review`, `skr search testing`)
   - Search by domain if the user mentioned one

3. **Evaluate relevance.** For each candidate, run `skr info <name>` to read
   the full description. Consider:
   - Does it address something this project actually needs?
   - Is it already covered by an installed package?
   - Would it conflict with existing rules or skills?

4. **Present recommendations.** For each recommended package:
   - Name and type (skill, rule, or knowledge)
   - Why it's relevant to this specific project
   - The install command: `skr install <name>`

5. **Let the user decide.** Present your recommendations and wait for the
   user to choose which to install. Do not install packages without explicit
   approval.

## Tips

- Prefer fewer, highly relevant recommendations over a long list
- If the user specified a focus area via the argument, narrow your search
- Knowledge packages are especially valuable for domain-specific projects
- Rules are most useful when they match the project's language and conventions
- Check `skr search registry-core` for the registry's own utility skills
