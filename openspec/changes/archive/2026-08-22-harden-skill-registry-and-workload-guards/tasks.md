# Tasks: Harden Skill Registry and Workload Guards

## Review Workload Forecast

| Field | Value |
|-------|-------|
| Estimated changed lines | ~350 lines across templates, skills, and tests |
| 400-line budget risk | Low |
| Chained PRs recommended | No |
| Suggested split | Single PR (or 3 Work Unit slices) |
| Delivery strategy | ask-on-risk |
| Chain strategy | pending |

Decision needed before apply: No
Chained PRs recommended: No
Chain strategy: pending
400-line budget risk: Low

### Suggested Work Units

| Unit | Goal | Likely PR | Focused test command | Runtime harness | Rollback boundary |
|------|------|-----------|----------------------|-----------------|-------------------|
| 1 | Shared Protocols & Skills (`skill-resolver.md`, `sdd-phase-common.md`, `sdd-tasks/SKILL.md`, `sdd-apply/SKILL.md`) | PR 1 | `go test ./internal/assets -run TestSDDTaskAndApplySkills` | N/A: static asset contracts | `internal/assets/skills/_shared/*`, `internal/assets/skills/sdd-*` |
| 2 | Core Orchestrator Templates (`antigravity`, `claude`, `opencode`, `codex`, `gemini`, `cursor`, `hermes`, `kimi`, `kiro`, `qwen`, `windsurf`, `generic`) | PR 2 | `go test ./internal/assets -run TestSDDOrchestrator` | N/A: markdown prompt templates | `internal/assets/{antigravity,claude,codex,cursor,gemini,generic,hermes,kimi,kiro,opencode,qwen,windsurf}/*.md` |
| 3 | Parity Tests & Golden Fixtures (`assets_test.go` parity suite, golden fixtures) | PR 3 | `go test -v ./internal/assets -run "TestOrchestrators\|TestSkillResolver"` | `go test ./internal/components` | `internal/assets/assets_test.go`, `testdata/golden/*` |

## Phase 1: Shared Protocols & Skills (Foundation)

- [x] 1.1 Update `internal/assets/skills/_shared/skill-resolver.md` with exact-path injection protocol and delegator cache invalidation contract.
- [x] 1.2 Update `internal/assets/skills/_shared/sdd-phase-common.md` with workload guard line matching and skill resolution status contract.
- [x] 1.3 Update `internal/assets/skills/sdd-tasks/SKILL.md` with workload forecast plain-text guard contract and work-unit structuring.
- [x] 1.4 Update `internal/assets/skills/sdd-apply/SKILL.md` with Step 2a fail-closed stop on unmanaged risk and single work unit execution bounds.

## Phase 2: Core Orchestrator Templates (Implementation)

- [x] 2.1 Update Antigravity, Claude, and Codex orchestrator templates (`antigravity/sdd-orchestrator.md`, `claude/sdd-orchestrator-workflow.md`, `codex/sdd-orchestrator.md`) with mandatory tasks skills injection, single work-unit apply batch, and resolver cache invalidation.
- [x] 2.2 Update Cursor, Gemini, and Generic orchestrator templates (`cursor/sdd-orchestrator.md`, `gemini/sdd-orchestrator.md`, `generic/sdd-orchestrator.md`) with mandatory tasks skills injection, single work-unit apply batch, and resolver cache invalidation.
- [x] 2.3 Update Hermes, Kimi, and Kiro orchestrator templates (`hermes/sdd-orchestrator.md`, `kimi/sdd-orchestrator.md`, `kiro/sdd-orchestrator.md`) with mandatory tasks skills injection, single work-unit apply batch, and resolver cache invalidation.
- [x] 2.4 Update OpenCode, Qwen, and Windsurf orchestrator templates (`opencode/sdd-orchestrator.md`, `qwen/sdd-orchestrator.md`, `windsurf/sdd-orchestrator.md`) with mandatory tasks skills injection, single work-unit apply batch, and resolver cache invalidation.

## Phase 3: Parity Testing & Golden Fixture Verification (Verification)

- [x] 3.1 RED Test: Add assertions in `internal/assets/assets_test.go` verifying all 12 orchestrators mandate `work-unit-commits` and `chained-pr` injection in `sdd-tasks`.
- [x] 3.2 RED Test: Add assertions in `internal/assets/assets_test.go` verifying single Work Unit apply dispatch (<=400 lines) and fail-closed stop on unmanaged risk.
- [x] 3.3 RED Test: Add assertions in `internal/assets/assets_test.go` verifying universal cache invalidation on non-injected skill resolution across all templates.
- [x] 3.4 GREEN: Run parity test suite to verify all template assertions pass (`go test ./internal/assets -run "TestOrchestrators|TestSkillResolver"`).
- [x] 3.5 Regenerate and verify golden fixtures across all agent runtimes (`go test ./internal/components`).

