---
name: code-review-checklist
description: >
  Guide code reviews with a structured checklist covering correctness, security,
  performance, readability, and testing. Use when reviewing pull requests or
  code changes.
user-invocable: true
argument-hint: "[file or PR]"
---

# Code Review Checklist

When reviewing code, work through this checklist systematically. Report
findings grouped by category with severity (critical, warning, note).

## Checklist

### Correctness
- Does the code do what it claims to do?
- Are edge cases handled (nil, empty, boundary values)?
- Are error paths handled (not silently swallowed)?
- Are return values checked?

### Security
- No hardcoded secrets, tokens, or credentials
- User input is validated before use
- No SQL injection, XSS, or command injection vectors
- Permissions and access controls are correct

### Performance
- No unnecessary allocations in hot paths
- No N+1 query patterns
- No unbounded growth (maps, slices, channels)
- Appropriate use of concurrency primitives

### Readability
- Names are clear and consistent with codebase conventions
- Functions are focused (single responsibility)
- No dead code or commented-out blocks
- Complex logic has comments explaining why (not what)

### Testing
- Are new code paths tested?
- Do tests cover error cases, not just happy paths?
- Are test assertions specific (not just "no error")?

## Output format

For each finding:
```
[severity] category: description
  file:line — specific code reference
  suggestion: what to change
```

Summarize with counts per severity at the end.
