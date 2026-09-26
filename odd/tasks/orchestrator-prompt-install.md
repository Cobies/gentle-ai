# Install the ODD/RDD orchestrator prompt for every runtime

Locator: `odd/tasks/orchestrator-prompt-install.md` · Engram mirror: `odd/orchestrator-prompt-install/tasks`
Branch: `fix/orchestrator-prompt-install` (worktree `gentle-ai-worktrees/orchestrator-prompt-install`, base `origin/main` c014fbe33)

## Objective

Every supported runtime (all agents except Pi, whose prompt is owned by the Gentle Shell package) gets the full ODD + RDD orchestrator instructions installed again, exactly as v3.7.0 installed its orchestrator section, minus SDD-only content.

## Problem

- SDD retirement (e219644b2) renamed `internal/assets/<runtime>/sdd-orchestrator.md` → `orchestrator.md` (12 runtimes: antigravity, claude, codex, cursor, gemini, generic, hermes, kimi, kiro, opencode, qwen, windsurf) keeping the ODD/RDD content, but the only reader was the retired SDD component (`internal/components/sdd/orchestrator.go`). No production code reads these assets on main.
- Only `agentguidance.RenderRouting` ("Implementation Routing" / ODD protocol) is installed. Sandbox fresh install v3.7.0 vs main: Claude `CLAUDE.md` 73KB→29KB, Codex `AGENTS.md` 100KB→33KB.
- Lost non-SDD sections: Agent Teams Orchestrator (coordinator role), Lossless Blocking Prompts, Provider Defect Handoff, Delegation Rules, Delegated Verification Gate, Native Checking Contract, Language Domain Contract, Native Compact Review Orchestration (RDD consent/capture/terminal procedure), Model Assignments, Optional Research.
- OpenCode/Kilo `gentle-orchestrator` agent prompt likewise only carries RenderRouting.
- `TestPrimaryODDOnlyOrchestrator` pins asset text but no test asserts installation → CI blind spot. Not yet released (v3.7.0 is fine); must land before the next release.

## Why

User requirement (2026-09-26): all agents must have ODD and RDD.

## Scope (authorized)

- Install each runtime's `orchestrator.md` (expanded with `_shared/odd-orchestrator-sections.md` placeholders as v3.7.0 did) into the same system-prompt target v3.7.0 used, under a non-SDD owner (agentguidance), idempotent, with rollback/backup/uninstall coverage.
- Upgrade: stale v3.7.0 `sdd-orchestrator` marker blocks are replaced, never duplicated.
- OpenCode/Kilo: the `gentle-orchestrator` agent prompt carries the orchestrator content.
- Pi untouched.
- Tests: per-runtime installed prompt contains the key ODD/RDD sections; no SDD-only sections; idempotent; sandbox diff vs v3.7.0 loses only SDD sections.

## Constraints

- No duplicate ODD protocol blocks (RenderRouting vs orchestrator content) — mirror v3.7.0 layout.
- Preserve user content outside managed markers; preserve file modes (`filemerge.ExistingFileMode`).
- ~400 authored lines per task is advisory only.

## Tasks

- [x] T1 — Wire orchestrator prompt installation for all non-Pi runtimes + OpenCode/Kilo agent prompt, upgrade replacement of legacy marker blocks, installation tests. Route: delegated (writer + preparation triggers: 12 runtimes, 4+ files).
- [x] T4 — Gentle Shell parity (authorized by user 2026-09-26): port the runtime-agnostic sections of Gentle Shell's orchestrator prompt (gentle-pi `assets/orchestrator.md`, `orchestrator-delegation.md`, `orchestrator-memory.md`, `orchestrator-skills.md` @89b8de3b5) into the shared ODD sections so every non-Pi runtime receives them: Identity Contract / Core Role / Mental Model / Safety, Work Routing Ladder (inline direct → simple delegation), Canonical Lightweight Workflows, Delivery strategy, Key Learnings closing block, Allowed edit surfaces requirement for writer delegations, Memory organic feature continuity, Skill Registry Protocol / Intent-Driven Skill Discovery (only where not already covered). EXCLUDED: Judgment Day activation/correction batch (stays in the `judgment-day` skill) and Pi-only runtime bindings (ODD phase signaling tool, Pi subagent model routing, Pi delegation/trigger bindings, Pi background policy). Adapt Pi tool names to runtime-neutral wording. Add a parity test listing the required sections for every non-Pi runtime's rendered prompt. Route: delegated (writer, after T1).
- [ ] T5 — RDD gating (user decision 2026-09-26: "RDD solo para OpenCode, Codex, Claude Code y Pi, nada más"). RDD-capable runtimes: claude-code, codex, opencode (+ pi, owned by Gentle Shell, untouched). Every other runtime (gemini-cli, cursor, kilocode, kiro-ide, qwen-code, windsurf, antigravity, hermes, kimi, vscode-copilot, openclaw, trae-ide) gets ODD only: no RDD content in the installed orchestrator or routing blocks (no Native Compact Review Orchestration, no Provider Defect Handoff consent-envelope handling, no "Receipt-driven development is user-owned", no RDD branches in the Delegated Verification Gate / ODD protocol / delivery text — use the non-RDD path wording). Single capability source of truth (e.g. an RDD-capable runtime set next to the review-transport manifest), tests asserting presence for the 3 runtimes and absence for all others. ALSO RDD agents (user clarification: each runtime has its own agents; RDD agents go only to these runtimes): `review-risk`, `review-readability`, `review-reliability`, `review-resilience`, `review-refuter`, `review-validator` must be installed only for claude-code, codex (if applicable), opencode (Pi separate). Today `reviewassets.NativeAgentManifest` installs them to cursor, kiro-ide, kimi (md+yaml) and PR #4998 added them to kilocode (`openCodeFamilyManagedRoles`/parity agents in internal/cli/run.go) → remove for those runtimes; Judgment Day agents (jd-*) are NOT RDD and stay. Upgrade: remove previously installed review agents that gentle-ai owns (native agent ownership record for cursor/kiro/kimi; Kilo entries matching the managed shape), never user-owned ones; drop the gentle-orchestrator task permissions for removed agents. Tests for install/upgrade removal and for presence on claude/opencode. Route: delegated (writer, after T4).
- [ ] T2 — Sandbox e2e: 16 runtimes fresh + upgrade from v3.7.0 + idempotent; section diff vs v3.7.0 must lose only SDD sections. Route: inline (bounded action).
- [ ] T3 — Full validation (gofmtcheck, vet, isolated `go test ./...` vs base), commit(s), native review, PR (delivery pending user decision).

## Acceptance criteria

1. Fresh install of each non-Pi runtime: installed prompt contains the ODD orchestrator sections (coordinator, Lossless Blocking Prompts, Delegation Rules, Delegated Verification Gate, Native Checking Contract, Language Domain Contract, Implementation Routing/ODD protocol, ported Gentle Shell sections) once each. RDD sections (Native Compact Review Orchestration, Provider Defect Handoff, RDD user-owned switch, RDD branches) only on claude-code, codex, opencode; absent everywhere else (T5).
2. No SDD-only sections (SDD Workflow, SDD Init Guard, Native SDD Dispatcher Guard, Artifact Store…).
3. Upgrade from v3.7.0: single orchestrator block, no leftover `sdd-orchestrator` block.
4. Second install/sync byte-identical for the orchestrator block.
5. CI green; isolated `go test ./...` no new failures vs base.

## Checks

- Focused package tests (isolated `env -i` HOME).
- `go run ./internal/gofmtcheck`, `go vet ./...`, `go test ./...` vs base.
- Sandbox matrix (`/tmp/e2e.sh`-style, v3.7.0 binary `/tmp/gai-v370-bin` must be rebuilt if missing).

## Delivery

Forecast: ~400–700 authored lines (+ tests). Strategy: ask-on-risk. RDD on (global).

## Progress

- 2026-09-26: Defect confirmed (no reader of orchestrator.md; sandbox size/section diff). Feature doc created.

- 2026-09-26: T1 delegated to gentle-ai-worker (task mui3gc39-b-ufp1).

- 2026-09-26: Parity check vs Gentle Shell: T1 restores the shared ODD/RDD core only; Gentle Shell has extra runtime-agnostic sections never ported → user added T4 (excluding Judgment Day, which stays in its skill).

- 2026-09-26: T1 done (writer mui3gc39-b-ufp1): `agentguidance.RenderOrchestrator` + managed `<!-- gentle-ai:orchestrator -->` section ahead of `agent-routing` (in-place migration of v3.7.0 `sdd-orchestrator`), OpenCode/Kilo gentle-orchestrator prompt, Kimi via `agent-routing.md`, ODD protocol single-sourced in routing, ExistingFileMode. Review contract via injectable source (import cycle reviewassets→agentguidance). Native Compact Review Orchestration only on runtimes advertising review transport (claude, opencode, codex) as in v3.7.0 — open product question. Isolated tests: only known env failures. Parent sandbox spot check (claude/codex/gemini/cursor): every ODD/RDD section exactly once, zero SDD sections. Writer incidents: created+removed temp dir `internal/zz_probe` (nothing left); an unisolated test run tried the real `~/.gentle-shell` APPEND_SYSTEM.md (dangling symlink, nothing written). Commit fb4b59a75 (+961/-17). Review assess (base c014fbe33): high, due → consent envelope relayed to user.

- 2026-09-26: T1 review granted by user → lineage review-51740b4a588cc84d (4 lenses) APPROVED + acknowledged; advisory only (R2/R3 OpenCode fallback drift unchecked; R4 orchestrator render failure blocks routing; section-bounds fail-open) → candidate follow-ups. Boundary → fb4b59a75. User now pre-authorizes all RDD consents. T4 delegated (gentle-ai-worker mui4zekm-c-s34e).

- 2026-09-26: User decided RDD applies only to OpenCode, Codex, Claude Code and Pi → added T5; T4 parity test must not require RDD sections on non-RDD runtimes (adjusted in T5).

- 2026-09-26: T4 done (writer mui4zekm-c-s34e): 3 new shared sections (Orchestrator Identity and Role; Orchestrator Routing and Delivery; Skill Registry Protocol for Claude only) referenced by all 12 orchestrator assets; parity test TestInstalledOrchestratorHasGentleShellParity (RED 15/15 runtimes → GREEN). Already covered, not duplicated: organic feature continuity (RenderRouting), memory lifecycle (Engram protocol), language boundary (Language Domain Contract), skill registry in 11 runtimes. Excluded: Pi-only bindings, RDD, JD dispatch. Commit 1ba7d888f; review assess (base fb4b59a75): medium, under_budget.
- 2026-09-26: ODD "assumption challenge" (routing step 3) is distinct from RDD `review-refuter`; T5 keeps the challenge on all runtimes and drops only the RDD-refuter clause on non-RDD runtimes. T5 delegated (gentle-ai-worker mui56usx-d-6w2a).

## Next step

T5 verify/commit/review → T2 → T3.
