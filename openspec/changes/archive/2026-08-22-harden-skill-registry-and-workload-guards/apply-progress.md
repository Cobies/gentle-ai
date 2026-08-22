# SDD Apply Progress: harden-skill-registry-and-workload-guards

## Status: Complete (13/13 tasks complete)

### Executive Summary
Implemented strict exact-path skill resolution, universal resolver cache invalidation, bounded single Work Unit apply dispatch (<=400 lines), and fail-closed workload risk guards across all 12 SDD orchestrator templates, shared protocols, and skills, verified with comprehensive Go parity test assertions and golden fixtures.

### Implemented Work Units

#### Work Unit 1: Shared Protocols & Skills (Foundation)
- **`internal/assets/skills/_shared/skill-resolver.md`**:
  - Added exact-path injection protocol: orchestrators resolve skills by canonical name and inject absolute paths under `## Skills to load before work`.
  - Added Step 4 Cache Invalidation: any fallback status reported by subagents (`fallback-registry`, `fallback-path`, `none`) forces the delegator to invalidate the cached skill registry immediately and re-read from persistent store.
- **`internal/assets/skills/_shared/sdd-phase-common.md`**:
  - Added exact plain-text guard lines: `Chain strategy: stacked-to-main|feature-branch-chain|size-exception|pending` in Section E alongside `Decision needed before apply`, `Chained PRs recommended`, and `400-line budget risk`.
  - Defined bounded Work Unit execution (<=400 changed lines `additions + deletions`).
  - Added subagent return envelope cache invalidation contract in Section D.
- **`internal/assets/skills/sdd-tasks/SKILL.md`**:
  - Standardized workload forecast plain-text guard output and Suggested Work Units table with focused test commands, runtime harness, and rollback boundaries.
- **`internal/assets/skills/sdd-apply/SKILL.md`**:
  - Enforced single Work Unit batch execution bounds (<=400 changed lines) in Step 2a.
  - Implemented explicit fail-closed stop returning `status: blocked` when unmanaged workload risk is detected.

#### Work Unit 2: Core Orchestrator Templates (Implementation)
Updated all 12 SDD orchestrator prompt templates:
- **`antigravity/sdd-orchestrator.md`**
- **`claude/sdd-orchestrator.md` / `claude/sdd-orchestrator-workflow.md`**
- **`codex/sdd-orchestrator.md`**
- **`cursor/sdd-orchestrator.md`**
- **`gemini/sdd-orchestrator.md`**
- **`generic/sdd-orchestrator.md`**
- **`hermes/sdd-orchestrator.md`**
- **`kimi/sdd-orchestrator.md`**
- **`kiro/sdd-orchestrator.md`**
- **`opencode/sdd-orchestrator.md`**
- **`qwen/sdd-orchestrator.md`**
- **`windsurf/sdd-orchestrator.md`**

Each template now strictly enforces:
1. **Mandatory Tasks Skills Injection**: When delegating `sdd-tasks`, the orchestrator resolves and injects both `work-unit-commits` (registry skill `gentle-ai-work-unit-commits`) and `chained-pr` (registry skill `gentle-ai-chained-pr`) under `## Skills to load before work`.
2. **Bounded Single Work Unit Apply Batch**: The orchestrator delegates `sdd-apply` strictly in single Work Unit batches bounded to <=400 changed lines (`additions + deletions`). Monolithic multi-work-unit apply dispatches are prohibited.
3. **Fail-Closed Workload Risk Guard**: When task summary reports workload risk (`Chained PRs recommended: Yes`, `400-line budget risk: High`, changed lines >400, or `Decision needed before apply: Yes`), the orchestrator blocks fail-closed unless an approved delivery strategy or maintainer `size:exception` is resolved.
4. **Universal Cache Invalidation**: When any subagent or delegated phase returns a `skill_resolution` of `fallback-registry`, `fallback-path`, or `none`, the orchestrator immediately invalidates the cached registry and re-reads from persistent store before subsequent delegations.

#### Work Unit 3: Parity Testing & Golden Fixture Verification (Verification)
- **`internal/assets/assets_test.go`**:
  - `TestOrchestratorsMandateTaskSkillsInjection`: Asserts all 12 orchestrator templates mandate injecting `work-unit-commits` and `chained-pr` under `## Skills to load before work` when delegating `sdd-tasks`.
  - `TestOrchestratorsEnforceBoundedWorkUnitApplyAndFailClosedGuard`: Asserts single Work Unit bounded apply dispatch (<=400 lines) and fail-closed stop on unmanaged risk.
  - `TestUniversalSkillCacheInvalidationOnFallbackResolution`: Asserts universal cache invalidation on non-injected skill resolution across all 12 orchestrators, `skill-resolver.md`, and `sdd-phase-common.md`.
  - `TestSDDTaskAndApplySkills`: Asserts `sdd-tasks/SKILL.md` and `sdd-apply/SKILL.md` contracts.
- **Golden Fixture Verification**:
  - Regenerated and verified golden fixtures across all agent runtimes (`testdata/golden/*`).
  - Ran `go test ./internal/assets/...` (PASS) and `go test ./internal/components/...` (PASS).

### Verification Evidence
```bash
go test -v ./internal/assets -run "TestOrchestrators|TestSkillResolver|TestSDDTaskAndApplySkills"
# PASS (all 12 templates + shared skills)

go test ./internal/assets/...
# ok github.com/gentleman-programming/gentle-ai/v2/internal/assets 13.698s

go test ./internal/components/...
# ok github.com/gentleman-programming/gentle-ai/v2/internal/components 0.564s
# ok github.com/gentleman-programming/gentle-ai/v2/internal/components/agentguidance (cached)
# ok github.com/gentleman-programming/gentle-ai/v2/internal/components/communitytool 93.747s
# ok all subpackages
```
