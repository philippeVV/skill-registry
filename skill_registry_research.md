# The Claude Code registry space is crowded — but the niche you want is still unclaimed

**Verdict up front: Do NOT build from scratch, but DO build.** Adopt Anthropic's `marketplace.json` + `SKILL.md` as your wire format (it's now an open standard), fork one of the two mature community registry stacks (`davila7/claude-code-templates` or `Kamalnrf/claude-plugins`) for the web UI and index, and spend your engineering budget on the three things nobody has shipped well: **LLM-judge eval scoring baked into the registry, private/air-gapped governance, and true agent self-discovery.** The Claude Code config-registry category has ~10 serious community projects and one mature cross-tool competitor (Continue Hub), but none combine eval-backed trust signals, internal/private governance, and agent-native discovery. That triad is your wedge.

The landscape changed drastically between October and December 2025. Anthropic shipped **Agent Skills** (Oct 16), **Plugins** (Oct 9), and the **Plugin Marketplace system** (GA Dec 18), then opened the Skills spec as an industry standard at `agentskills.io` with OpenAI Codex, Cursor, Copilot, goose, Amp, and Letta all signing on. The format wars are over; the distribution primitives are defined. What's contested is the *layer above*: curation, trust, eval, analytics, and enterprise governance. That's where you play.

## What Anthropic already shipped — and what it means for you

Anthropic's `.claude-plugin/marketplace.json` IS a registry format. A "marketplace" in Claude Code is literally a Git repo with a manifest listing plugins; each plugin is a directory bundling skills, subagents, slash commands, hooks, output styles, MCP servers, and LSP servers. Users subscribe with `/plugin marketplace add your-org/claude-plugins` and install with `/plugin install formatter@your-team-tools`. Scopes (user/project/local/managed), caching (`~/.claude/plugins/cache/`), persistent data (`${CLAUDE_PLUGIN_DATA}`), auto-update semantics, namespacing (`plugin-name:skill-name`), and enterprise-controlled `extraKnownMarketplaces` settings all exist. **Reinventing any of this is wasted work** — and worse, any registry that doesn't speak this schema can't be consumed by `/plugin` natively.

The critical gotcha: Anthropic's official schema URL (`https://anthropic.com/claude-code/marketplace.schema.json`) **returns 404**. The de facto schema lives in the docs at `https://code.claude.com/docs/en/plugins-reference` and in community JSON Schemas (`hesreallyhim/claude-code-json-schema`). The `agentskills.io` spec itself is — per Simon Willison — "quite heavily under-specified." This is good news for you: there's room to ship a stricter internal schema variant with CI validation that Anthropic hasn't gotten around to.

What Anthropic does **not** provide, and what your registry therefore legitimately adds: private/org-scoped discovery UI, install analytics, per-skill eval scoring, vulnerability and prompt-injection scanning, audit trails, air-gap support, and non-plugin knowledge-snippet distribution. The official `claude.com/plugins` directory has an "Anthropic Verified" badge but no install counts, no ratings, and no eval framework.

## Direct competitors and adoption candidates

Six projects meaningfully overlap with your vision. The maturity assessment below is honest — some will close gaps you're targeting before you ship if you're slow.

**`davila7/claude-code-templates` (aitmpl.com)** is the most complete public competitor. MIT-licensed, ~1000+ components (600+ agents, 200+ commands, 60+ settings, 39+ hooks, skills, MCPs), `npx claude-code-templates@latest` interactive CLI, **Supabase-backed install analytics at aitmpl.com/download-stats**, real-time dashboard, category browsing, PR-based contributions. Missing from your checklist: LLM-judge eval scoring, leaderboards, true private-mode. Active — v1.27.0 shipped November 2025, weekly releases. https://github.com/davila7/claude-code-templates

**`davepoon/buildwithclaude` (buildwithclaude.com)** is the same idea with a cleaner UX: 55 plugins / 125 skills / 117 subagents / 175 commands / 28 hooks cataloged, `bwc-cli` on npm, PR-merges auto-deploy to the web UI, categories and tags, partial install counts. Also MIT. Arguably closer aesthetically to what you're sketching. https://github.com/davepoon/claude-code-subagents-collection

**`Kamalnrf/claude-plugins` (claude-plugins.dev)** is architecturally the closest to your spec: a **Val.town-backed auto-indexer scrapes all public Claude Code plugin repos on GitHub every 10 minutes** and exposes a JSON REST API with `{stars, downloads, verified, metadata}`, plus `npx claude-plugins install owner/marketplace/plugin`. Open source, lightweight. This is what "lightweight JSON index on S3" looks like as a finished product. https://claude-plugins.dev

**`iflytek/skillhub`** (Apache-2.0, beta) is the enterprise-self-hostable reference implementation you probably didn't know existed. Java 21 + React, **semantic versioning with beta/stable/latest tags, team namespaces with two-tier promotion workflow (team-admin reviews → platform-admin gates global promotion), pluggable storage (local FS or S3/MinIO), RBAC, API tokens with prefix-hashed storage, audit logs, Docker Compose + K8s manifests, Prometheus+Grafana monitoring**, and a ClawHub protocol compatibility layer so existing CLI clients can re-target via env var. This is the closest thing to your enterprise vision shipped as open source. Worth cloning architectural patterns directly even if you don't use the code. https://github.com/iflytek/skillhub

**Continue Hub (hub.continue.dev)** is the most mature cross-tool competitor — composable "blocks" (rules, prompts, models, context, docs, MCP, data), `uses: owner/slug@version` composition syntax like npm deps, three publish paths (web UI, GitHub Action, template repo sync), first-class public/private/org visibility flags, semver. **Closed-source hub despite open-source IDE extension** — a meaningful strategic weakness you should exploit by being open. No install counts visible. https://hub.continue.dev

**`wshobson/agents`** (~32k stars) and **`VoltAgent/awesome-claude-code-subagents`** (~8k stars, MIT) are content-heavy collections worth mining for seed content. VoltAgent notably ships a meta-agent called `agent-installer` that queries the GitHub API to browse and install agents — the only clear example of **agent self-discovery** shipped publicly.

A few close-but-narrower projects: **`aig787/ccpm`** (Rust, Cargo-style `ccpm.toml` + `ccpm.lock` for reproducible installs from any Git repo), **`sjnims/cc-plugin-eval`** (4-stage LLM-judge eval framework for plugins — this is basically the missing eval piece you want, as a library), **`LiteLLM`** (exposes `/claude-code/marketplace.json` from a gateway — the canonical enterprise-proxy pattern for private registries).

Outside the Claude-specific bubble, three adjacent projects have patterns worth borrowing: **Smithery** runs a registry-as-MCP-server so agents query `search_servers` / `get_server_details` at runtime; **Vercel's skills.sh** populates its leaderboard entirely from `npx skills add` telemetry (zero-submission-friction publishing) and runs Snyk security scans server-side before install; and **PRPM (prpm.dev)** has a canonical-format-plus-server-side-fanout model that converts one manifest into Cursor/Claude/Copilot/Continue/Windsurf/Kiro/Aider formats.

One strong warning: **`hesreallyhim/awesome-claude-code` is CC BY-NC-ND 4.0** — no commercial use, no derivatives. You can reference it; you cannot fork its data for an internal tool at a company. Most other projects mentioned are MIT or Apache-2.0.

## Design ideas worth stealing, grouped

### Publishing patterns

The highest-leverage idea is **auto-detect from repo structure**, shipped by `cursor.directory` as their "Open Plugins" spec: paste a GitHub URL, the registry scans conventional paths (`skills/*/SKILL.md`, `.mcp.json`, `agents/*.md`, `hooks/hooks.json`), and auto-extracts metadata with zero manual manifest. Your PR-based model is fine for governance, but combine it with auto-detect so reviewers just approve rather than hand-authoring metadata. The **`strict: false` mode** in Anthropic's marketplace.json (the marketplace entry itself IS the manifest, skipping per-plugin `plugin.json`) is the lightweight variant of this.

**Goose's `GOOSE_RECIPE_GITHUB_REPO` env-var pattern** is the simplest private-registry implementation on earth: one env var points at a private GitHub repo and that repo becomes the registry. You should support this as a fallback/bootstrap mode even if your "real" registry is S3-backed. **Multi-path publishing** (web form + GitHub Action + PR to template repo) as Continue does captures every user archetype with one shared backend.

For enterprise governance, **iflytek/skillhub's two-tier promotion workflow (team namespace → global, with separate admin approval at each stage)** is the governance model you actually need. Tie it to signed commits and you have SOC-2-ready audit trails out of the box.

### Discovery patterns

The table-stakes signals are search, tags, install counts, and categories. Everyone is converging here. What's differentiated:

- **Description-based AI-judged activation** (Cursor's `.mdc` rules, Anthropic's SKILL.md frontmatter): the agent reads short descriptions of everything at session start and decides relevance. Anthropic's `<available_skills>` injection — a system-prompt-level index loaded from frontmatter across user/project/plugin scopes, with `<location>` for scope disambiguation — is the canonical pattern and is what makes progressive disclosure actually work.
- **Three-tier progressive disclosure** (SKILL.md canonical): Tier 1 frontmatter always loaded (~50-100 tokens per skill, 7-13k tokens for a 133-skill catalog), Tier 2 body loaded on activation, Tier 3 `references/` files loaded on-demand. Your knowledge bites and CLAUDE.md fragments should follow the same pattern.
- **Self-referential curation prompts** (github/awesome-copilot's `suggest-awesome-github-copilot-prompts.prompt.md`): a skill whose job is to scan your repo, fetch the registry catalog, and recommend installs. Ship this on day one — it's trivial and it's your best agent-self-discovery primitive.
- **Registry-as-MCP-server** (Smithery): expose your registry as a queryable MCP server so Claude Code can `search_skills` and `install_skill` at runtime. This is the right architecture for true agent self-discovery, beyond the `<available_skills>` injection that handles already-installed skills.
- **VoltAgent's `agent-installer` meta-agent** pattern: a subagent that browses the GitHub API for skills in categories. Simpler than MCP, works today.

### Install UX patterns

Anthropic's **two-step `marketplace add` → `plugin install`** (catalog subscription separate from artifact install) is the right primitive — use it. Map your S3 bucket behind a URL and users `/plugin marketplace add https://claude-registry.internal.yourorg.com/marketplace.json`.

For the CLI, steal from `aig787/ccpm`: **`ccpm.toml` manifest + `ccpm.lock` lockfile** for reproducible installs. Critical for SRE-grade determinism. Plus **`--symlink` install mode** (tech-leads-club/agent-skills) for authors editing a skill live without reinstalling. Plus **content-hash immutability** on the storage layer so a published version can never change underneath you.

**Deep-link URLs** (`claude://` scheme if Anthropic supports it, or deep-link HTTP URLs rendered by your web UI) mean your browse page can have one-click install buttons. **Collections/bundles** (PRPM's `collections/nextjs-pro` installs 20+ packages at once) are how platform teams ship "starter packs" for new projects.

### Trust signals and eval

This is your real differentiation opportunity because almost no one does this well.

**Install count + usage count** with explicit rubric thresholds (Vercel skills.sh codified ">1K installs = trust, <100 = caution" directly in their `find-skills` skill). Author leaderboards are easy wins. **Review-upvote meta-voting** (Smithery) prevents single-review domination in a small-org context where you'll have handful-of-reviews-per-skill problems.

The genuinely novel move: **`scenarios.yaml` bundled in each skill, registry runs promptfoo (or `sjnims/cc-plugin-eval`) on publish, publishes `eval_score`, `pass_rate`, and `last_eval_run` as first-class metadata surfaced in search and install UX**. Your idea of "LLM-judge scores task-with-skill vs task-without-skill" is the right version of this and *no public registry is doing it*. The reference architectures:
- `sjnims/cc-plugin-eval` — 4-stage framework, programmatic activation detection + LLM-judge, multi-sample scoring with confidence scores, conflict detection, cost-saving batches API. Output is `quality_score 0-10`, `trigger_rate`, `avg_quality`. Integrate this as your admission gate.
- `hamelsmu/evals-skills` — skills for eval-driven development of skills (meta).
- MLflow blog on `scorer-registered`/`agent-eval-skill-invoked` — architecture inspiration.
- PRPM's **Playground** (test a package against Claude 3.5 / GPT-4 before install) is the "try before install" UX around the same eval.

Security-adjacent trust signals increasingly matter: **Snyk scans every `npx skills add` on skills.sh server-side**, and **JFrog's Agent Skills Registry (announced April 2026 with NVIDIA) cryptographically signs skills on publish and specifically scans for prompt-injection in skill bodies** — a unique threat for this artifact class. PromptArmor has already demonstrated malicious plugins in Claude Code's auto-discovering marketplace. For an internal tool at a real company, **hook-content static analysis + signed commits on PR-merge** is the bare minimum, and a **prompt-injection scanner** run against skill descriptions is a credible differentiator.

### Schema, versioning, scoping

Adopt `SKILL.md` wholesale as the atomic artifact — it's now the agentskills.io standard. Use Anthropic's `.claude-plugin/plugin.json` as the bundle manifest and `.claude-plugin/marketplace.json` as the registry index. Extend — don't replace — with your internal-only fields (e.g., `reviewer`, `eval_score`, `security_scan_status`, `owner_team`, `slack_channel`). Use **semver + named tags (`beta`/`stable`/`latest`)** from iflytek/skillhub.

For scoping conflicts: Anthropic's precedence (enterprise > personal > project, plus `<location>` tags in the injected index) is good. Add **closest-file-wins from AGENTS.md** semantics when skills cascade in monorepos. For multi-skill overlap, Claude Code's `disable-model-invocation` frontmatter flag is how skills opt out of auto-triggering; make this a first-class registry filter.

One non-obvious pattern from agentic-community's mcp-gateway-registry: a **virtual-gateway layer** that resolves naming conflicts, version-pins to specific backends, and enforces scope-based access control per-tool. Probably overkill for v1 but worth remembering as you scale.

### Interop to steal

Natively consume and export **`AGENTS.md`** (the Linux Foundation-stewarded cross-tool standard, 60k+ repos). Claude Code doesn't natively support AGENTS.md yet — you can bridge it and make your internal registry the place where skills authored for Cursor/Codex/Copilot also work in Claude Code. PRPM's canonical-format-plus-fanout demonstrates feasibility; `dyoshikawa/rulesync` does the format conversion today.

## Gaps your product would uniquely fill

Three gaps are well-defined and defensible.

**Eval-backed trust as first-class metadata.** Every competitor has some install count. None surface eval scores, pass rates, or task-success-delta (with skill vs. without skill) in search and install UX. This is the one place your product can be *qualitatively* better — not "nicer UI" but "this one actually tells you if the skill works." Wire `sjnims/cc-plugin-eval` into your CI, publish scores to the index, rank leaderboards by score-weighted-installs.

**Private, air-gapped, SRE-grade governance for Claude Code specifically.** Continue Hub has private visibility but the hub is closed-source; iflytek/skillhub has governance but isn't Claude-native; LiteLLM has the gateway pattern but not the UX. A product that does "private S3-backed Claude Code registry with signed publishes, audit trail, promotion workflow (team → global), vulnerability + prompt-injection scanning, and speaks `marketplace.json` natively" fills a clean gap. This is your SRE-at-a-Quebec-company demographic's actual need.

**Agent self-discovery that's more than a list.** Anthropic's `<available_skills>` handles already-installed artifacts. Nothing good exists for runtime discovery-then-install. Ship two pieces: (1) expose your registry as an MCP server with `search_skills`/`install_skill` tools, so Claude Code can resolve "I need a skill for X" at runtime; (2) bundle a `suggest-skills` self-referential skill on install, so every onboarded developer's Claude Code instance automatically discovers what's in the registry relevant to their current repo.

A fourth, softer gap worth considering: **non-plugin knowledge artifacts**. Your "knowledge bites" and "CLAUDE.md fragments" don't fit naturally into Anthropic's plugin/skill model (plugins bundle behaviors; CLAUDE.md fragments are inert context). This is genuinely undeserved space. The hack today is wrapping knowledge as a skill with `user-invocable: false`, but that's not how people think about documentation-you-inject. A first-class `knowledge` artifact type with its own install semantics (append-to-CLAUDE.md scoped by glob, rather than copy-to-`.claude/skills/`) is legitimately novel.

## Recommendations on differentiation and build plan

**Stack recommendation.** Use `anthropic/marketplace.json` schema as the wire format (non-negotiable for `/plugin` compat). Fork `davila7/claude-code-templates` or `Kamalnrf/claude-plugins` for the web UI and JSON index — both MIT, both close to your vision. Layer `iflytek/skillhub`'s namespace/promotion/audit model on top. Integrate `sjnims/cc-plugin-eval` as an admission gate. Expose the registry as an MCP server for agent self-discovery. Ship a `suggest-skills` self-referential skill bundled with onboarding. Use LiteLLM's gateway pattern if you already run LiteLLM internally, otherwise a direct S3+CloudFront+VPN setup is simpler.

**Positioning.** "The internal Claude Code registry with evals, governance, and audit." Not "another skill hub." Your three-line pitch to the dev-eng department: skills that actually work (evals), skills that actually got reviewed (PR-merge gate), skills that work offline and stay inside the VPN (S3 + signed artifacts).

**What to skip in v1.** Don't build ratings/reviews (low signal at internal scale; rely on evals and install counts). Don't build a desktop app (CLI + web UI is enough). Don't build cross-tool format conversion (AGENTS.md export is a v2 feature). Don't reinvent the manifest — inherit Anthropic's and extend.

**Risks to watch.** Anthropic could ship an official enterprise registry with eval scoring any quarter; your moat is internal governance tied to your org's systems (IdP, Slack, audit logs) that Anthropic's generic product won't match. `davila7/claude-code-templates` could add eval scoring and leaderboards before you ship; mitigate by being honest about what you're building and scoping tightly to internal use where their public SaaS doesn't apply. The agentskills.io spec is under-specified and could churn; pin to a specific Claude Code version and test upgrades in CI.

## Conclusion

Your instinct to build was right, your concern that the space is crowded is also right, and the two are reconcilable. The *primitives* (Skills, Plugins, Marketplace) are defined and standardized; building those would have been wasted effort six months ago but is now clearly redundant. The *curation/trust/governance layer* above the primitives is contested in public but essentially unbuilt for private/internal Claude Code registries with eval-backed quality signals. You have a narrow but real gap — roughly 6-12 months wide before either Anthropic ships enterprise governance natively or `davila7`/`davepoon` figure out evals. Ship the MVP in weeks not months, inherit the schema wholesale, and put engineering hours into the eval harness and the MCP-served self-discovery interface. Those are the two features that will make developers actually reach for your registry instead of `/plugin marketplace add anthropics/skills`.