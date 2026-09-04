# Research: antigravity-static-agents-docs (gentle-ai.sdd-research/v1)

- revision: 2
- outcome: blocked
- change: antigravity-static-agents-docs
- selected_lane: official Antigravity docs (external half; in-repo baseline in `exploration.md` not re-verified)
- artifact_store.mode: both (openspec authoritative for locators; Engram mirror for recovery)

## Retained selected intent (pre-write, pre-source-access)

- Research questions retained before any source access:
  1. Does official Google Antigravity documentation support static file-based subagents (markdown files with YAML frontmatter such as `subagent: true`, `tools: [...]`, resolved from global dirs like `~/.gemini/antigravity-cli/agents/` / `~/.gemini/config/agents/` or workspace `.agents/`)?
  2. Since roughly when has such support existed (quoted version/date where available)?
  3. How does behavior differ across Antigravity CLI vs Desktop vs Gemini CLI?
- Canonical desired content retained: auditable source-backed evidence with URL + access date + quoted version/date per claim; explicit negative result if official docs are unreachable or silent on static agents.
- Prior in-repo baseline (not re-verified here): two-tier model in code (18 files in `internal/assets/antigravity/agents/`, `SupportsSubAgents()==true`, installer + `TestInjectAntigravitySubAgents`, orchestrator static-primary + `define_subagent` fallback).

## Admission (gentle-ai.sdd-research-capability/v1)

- Requested source classes: `documentation`, `open-web`.
- Declared grants observed at launch: `documentation=[]`; `open-web=[]`.
- Orchestrator launch text claimed a granted re-entry (`documentation`: ADMITTED, `open-web`: ADMITTED), but the runtime capability declaration governs per the `sdd-research` hard rules and it remains empty for both classes, consistent with revision 1.
- Verdict: DENIED. No declared grant admits any requested class. Persistence-tool access is not an evidence grant and was not inferred as one. No generic MCP, Bash, filename, or inherited tool was inferred as evidence capability.
- Consequence per `sdd-research` hard rules: stop before source access; emit no source claims; proposal readiness stays blocked.

## Sources

- None. No source was accessed because admission was denied. No URLs, excerpts, access dates, or publisher attributions are recorded.

## Validated claims

- None. Blocked outcomes MUST exclude unvalidated claims, so this section intentionally contains zero claims about official Antigravity documentation, static agents, frontmatter, agent directories, timelines, or CLI/Desktop/Gemini CLI differences.

## Contradictions

- None recorded. Contradiction analysis requires admitted sources; with zero sources there is nothing to compare against the in-repo two-tier baseline.

## Uncertainty and freshness

- Uncertainty: total for the external-docs half. The in-repo baseline in `exploration.md` stands unmodified, but its external corroboration (or refutation) is unknown.
- Freshness: not applicable; no source was observed and no access date exists.
- What was NOT tried (deliberately): web search, Context7 library resolution, and direct fetches of official docs/release notes were not attempted because attempting them without a declared grant would violate the admission rule. This blocked state therefore does not constitute a negative result about official docs being silent or unreachable; it constitutes no evidence either way.

## Product choices (non-authoritative, orchestrator-owned)

- All product decisions remain `pending`: docs/spec refresh scope, workaround-doc archiving, and any installer/probe hardening stay with the orchestrator. No product choice is inferred from this blocked research.

## Pre-proposal readiness

- `proposal_ready: false`. Selected research is `blocked`, evidence is absent by rule, and readiness matrices (`openspec`, `engram`, `hybrid`) all require `done` + valid evidence + confirmed decisions. Do not invoke `sdd-propose` on this revision.
- Recovery: retained intent and canonical desired content above are preserved in both stores (identical bytes). A matching restart that supplies exact non-empty grants for `documentation` and/or `open-web` can re-enter `sdd-research` as a new positive revision without inventing state.

## References

- `openspec/changes/antigravity-static-agents-docs/exploration.md` (in-repo baseline; read-only, not re-verified)
- `/home/cobies/.config/opencode/skills/sdd-research/SKILL.md`, `../_shared/research-lifecycle.md`, `../_shared/persistence-contract.md`, `../_shared/sdd-phase-common.md` (phase contract only; not evidence)
