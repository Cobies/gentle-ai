# OpenCode agent parity with Gentle Shell (fixes #4471, #4684, #4758)

Locator: `odd/tasks/opencode-agent-parity.md` · Engram mirror: `odd/opencode-agent-parity/tasks`
Branch: `fix/4471-opencode-agent-parity` (worktree `gentle-ai-worktrees/4471-opencode-agent-parity`, base `origin/main` 6c7f162f4)

## Objective

A fresh or upgraded `gentle-ai install/sync --agent opencode` yields OpenCode subagents functionally at parity with Gentle Shell's global agents, and no managed OpenCode agent entry carries `__managed_by`. Every runtime regains the non-SDD assets the SDD retirement dropped by accident (T4).

## Problem

- Up to v3.7.0, `internal/assets/opencode/sdd-overlay-*.json` defined 23 OpenCode agents, each with `"__managed_by": "gentle-ai/sdd"`. OpenCode forwards unknown agent fields as provider request options; strict providers (opencode-go, Fireworks, NVIDIA NIM) reject every call (#4471, dups #4684, #4758).
- Commit e219644b2 (retire SDD) deleted that overlay. It was the sole owner of ALL OpenCode agents, not only SDD ones, so JD (`jd-judge-a/b`, `jd-fix-agent`) and the four review lenses were dropped as collateral. Only `gentle-orchestrator` (permission patch), `review-refuter`, `review-validator` (run.go `installOpenCodeReviewProviderRoles`) and `gentleman` remain.
- Observed: fresh install from main in isolated HOME → `opencode.json` agents = gentle-orchestrator, gentleman, review-refuter, review-validator.
- Writes use `filemerge.MergeJSONObjects` (deep merge, never deletes) → upgraded users keep stale v3.7.0 entries including `__managed_by`; no migration, no test.
- `sync.go` model assignment, TUI model picker, and `opencode-review-transport.ts` still expect JD/lens agents.

## Why

User requirement: only SDD was meant to be retired; OpenCode agents must be functionally the same as Gentle Shell's (names may differ).

## Scope (authorized)

Parity set (canonical source: gentle-pi `assets/agents/*.md` at 89b8de3b5): `gentle-ai-explore`, `gentle-ai-verify`, `gentle-ai-worker`, `jd-judge-a`, `jd-judge-b`, `jd-fix-agent`, `review-risk`, `review-readability`, `review-reliability`, `review-resilience`. Plus retained `review-refuter`, `review-validator`, `gentle-orchestrator`, `gentleman`.

- Names: identical to Gentle Shell (decision: simplest parity; no collision with OpenCode built-in `explore`/`general`).
- Format: JSON agent entries in `opencode.json` (same owner style as `installOpenCodeReviewProviderRoles`), prompt bodies embedded under `internal/assets/opencode/agents/`. No `__managed_by` or any non-AgentConfig key.
- Routing: OpenCode `orchestrator.md` delegates to `gentle-ai-explore` / `gentle-ai-worker` / `gentle-ai-verify` like Gentle Shell.
- Migration on install/sync: `__managed_by == "gentle-ai/sdd"` is ownership proof. Owned entries for current roles are replaced cleanly; owned retired entries (`sdd-*`, v3.7.0 `general`/`explore` overrides) are removed; the field is stripped from any remaining entry. Unmarked user entries untouched. `managedConfigPriority` keeps working via shape heuristic / new ownership signal.

## Constraints

- Preserve unrelated user config; idempotent re-runs.
- Model assignment (`sync.go`) and TUI model picker cover the parity set.
- ~400 authored lines per task is advisory only; embedded prompt bodies are copied assets.

## Tasks

- [x] T1 — Install parity agents for OpenCode (embedded prompts, permissions from Gentle Shell tool sets), orchestrator routing, model assignment/picker allowlist, parity inventory test. Route: delegated (writer trigger: 2+ non-trivial files).
- [ ] T2 — Upgrade migration: strip `__managed_by`, remove owned retired entries, replace owned current-role entries; regression tests from a v3.7.0 config fixture. Route: delegated (same writer, sequential).
- [ ] T4 — Restore non-SDD assets dropped by the SDD retirement (authorized by user 2026-09-25): install `skills/_shared/*` still referenced by orchestrators/skills (skill-resolver, engram-convention, persistence-contract, research-lifecycle, review-ledger-contract(-pi), README) for every runtime that installed them in v3.7.0, under a non-SDD owner; restore `skill-creator`/`skill-registry` commands for OpenCode, Kilo, Qwen; restore OpenCode `default_agent: gentle-orchestrator` (with existing opencodedefault ownership file) and `share: disabled`. Tests: per-runtime install inventory asserts every `_shared/*.md` referenced by installed orchestrator/skills exists (no dangling refs); OpenCode default_agent/share asserted. Route: delegated (writer trigger).
- [ ] T3 — Sandbox e2e: fresh install + upgrade from v3.7.0 config in isolated HOME (no Homebrew on PATH), idempotent re-run. Route: inline (bounded action).

## Acceptance criteria

1. Fresh install: all 10 parity agents + refuter/validator present; no agent entry contains `__managed_by`.
2. Upgrade from v3.7.0 config: no `__managed_by` anywhere; no `sdd-*`, owned `general`/`explore` removed; user-defined agents untouched.
3. Second run produces byte-identical config.
4. T4: no installed orchestrator or skill references a `_shared/*` file that is not installed; OpenCode/Kilo/Qwen skill commands present; OpenCode `default_agent` = gentle-orchestrator and `share` = disabled on fresh install; ownership preserved on uninstall.
5. `go test ./...` green (or known env failures documented); `go vet`, `gofmt` clean.

## Checks

- `go test ./internal/cli/... ./internal/opencode/... ./internal/components/... ./internal/tui/... ./internal/assets/...`
- `go test ./...`, `go vet ./...`, `gofmt -l .`
- Sandbox install per T3.

## Delivery

Forecast: ~700 authored lines + ~600 copied prompt lines. Strategy: ask-on-risk (default) — ask before PR. RDD: enabled by the user on 2026-09-25 (clone-local unset, global on) → per work-unit commit `review assess --committed-only` with base = branch point 6c7f162f4.

## Progress

- 2026-09-25: regression reproduced (fresh install, 4 agents only). Feature doc created.

- 2026-09-25: T1 delegated to gentle-ai-worker (task muhcno1f-3-f75u).
- 2026-09-25: cross-runtime sandbox audit (v3.7.0 vs main, 16 runtimes) found further SDD-retirement collateral, authorized by the user the same day as T4: `skills/_shared/*` no longer installed anywhere but still referenced by orchestrators and judgment-day/skill-registry skills; OpenCode/Kilo/Qwen lose `skill-creator`/`skill-registry` commands; OpenCode loses `default_agent: gentle-orchestrator` and `share: disabled`. Engram: `bug/sdd-retirement-collateral-regressions`.

- 2026-09-25: T1 done (writer muhcno1f-3-f75u + parent inline fixes: dropped Pi-only host-relay line from the 4 lens prompts; removed the `codegraph` CLI fallback from gentle-ai-explore because it has bash denied). Commit 72519acad (20 files, +924/-9 incl. ~600 copied prompt lines). Evidence: RED `undefined: installOpenCodeODDParityAgents`; GREEN `TestOpenCodeInstallWritesParityAgentsWithoutManagedByMarker`, `TestOpenCodeInstallParityAgentsAreIdempotent`; assets/opencode/tui/opencodedefault/telemetry packages ok; gofmt/vet clean. Known env failures on base 6c7f162f4 too (parent re-ran on clean main): `TestSyncPersonaOnlyRollbackRestoresOpenCodeSettingsAfterGentlemanCleanup`, `TestSyncRollbackRestoresLegacyOpenCodePluginAndRemovesReplacement` (+ Pi `~/.gentle-shell` leakage tests per writer). Review assess (agent pi, base 6c7f162f4): medium, review_due=true (slice_budget_reached) → START returned consent envelope, awaiting user answer. Lineage review-482425bf66d8e129.

- 2026-09-25: T1 review granted by user → lineage review-482425bf66d8e129 (1 lens: review-reliability) APPROVED and acknowledged (authority burned). Advisory non-blocking finding R3-001 (WARNING, internal/cli/run.go:1150). Reviewed boundary → 72519acad.
- 2026-09-25: user authorized finishing everything through PR + merge while away, testing that everything works first; user pre-authorized granting review consent for remaining candidates. T2 delegated (gentle-ai-worker muhfqpvl-4-q4gc).

## Next step

T2 verify/commit/review → T4 → T3 e2e → PR → merge.
