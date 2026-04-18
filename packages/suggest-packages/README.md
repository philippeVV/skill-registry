# suggest-packages

A skill that analyzes your current project and recommends relevant packages
from the skill registry.

## Usage

```
/suggest-packages
/suggest-packages backend
/suggest-packages testing
```

## What it does

1. Reads your project structure, languages, and frameworks
2. Checks what's already installed via `skr list`
3. Searches the registry for relevant packages
4. Presents recommendations with reasoning

## System skill

This skill is automatically installed on first `skr install` as part of the
registry core utilities.
