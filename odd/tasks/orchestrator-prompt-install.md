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

- [ ] T1 — Wire orchestrator prompt installation for all non-Pi runtimes + OpenCode/Kilo agent prompt, upgrade replacement of legacy marker blocks, installation tests. Route: delegated (writer + preparation triggers: 12 runtimes, 4+ files).
- [ ] T4 — Gentle Shell parity (authorized by user 2026-09-26): port the runtime-agnostic sections of Gentle Shell's orchestrator prompt (gentle-pi `assets/orchestrator.md`, `orchestrator-delegation.md`, `orchestrator-memory.md`, `orchestrator-skills.md` @89b8de3b5) into the shared ODD sections so every non-Pi runtime receives them: Identity Contract / Core Role / Mental Model / Safety, Work Routing Ladder (inline direct → simple delegation), Canonical Lightweight Workflows, Delivery strategy, Key Learnings closing block, Allowed edit surfaces requirement for writer delegations, Memory organic feature continuity, Skill Registry Protocol / Intent-Driven Skill Discovery (only where not already covered). EXCLUDED: Judgment Day activation/correction batch (stays in the `judgment-day` skill) and Pi-only runtime bindings (ODD phase signaling tool, Pi subagent model routing, Pi delegation/trigger bindings, Pi background policy). Adapt Pi tool names to runtime-neutral wording. Add a parity test listing the required sections for every non-Pi runtime's rendered prompt. Route: delegated (writer, after T1).
- [ ] T2 — Sandbox e2e: 16 runtimes fresh + upgrade from v3.7.0 + idempotent; section diff vs v3.7.0 must lose only SDD sections. Route: inline (bounded action).
- [ ] T3 — Full validation (gofmtcheck, vet, isolated `go test ./...` vs base), commit(s), native review, PR (delivery pending user decision).

## Acceptance criteria

1. Fresh install of each non-Pi runtime: installed prompt contains Agent Teams Orchestrator/coordinator, Lossless Blocking Prompts, Provider Defect Handoff, Delegation Rules, Delegated Verification Gate, Native Checking Contract, Language Domain Contract, Native Compact Review Orchestration, and Implementation Routing/ODD protocol — once each.
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

## Next step

Await T1; parent verifies, commits, runs review assess (agent pi), then T2.
