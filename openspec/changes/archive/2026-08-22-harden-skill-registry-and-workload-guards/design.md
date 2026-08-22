# Design: Harden Skill Registry and Workload Guards

## Technical Approach

Enforce reviewer workload safety (<=400 lines budget) and deterministic skill resolution across all 12 SDD orchestrators. The approach tightens orchestrator templates and phase contracts at four integration boundaries:
1. **Task Planning Injection**: Orchestrators mandate resolving and injecting `work-unit-commits` and `chained-pr` under `## Skills to load before work` when delegating `sdd-tasks`.
2. **Single Work Unit Apply Delegation**: Orchestrators and `sdd-apply` restrict execution batches to single autonomous Work Units (<=400 lines). Monolithic task execution is forbidden.
3. **Fail-Closed Workload Risk Guard**: When workload forecast flags >400 lines or High risk, the orchestrator and `sdd-apply` halt fail-closed before any file edits unless an approved delivery strategy (`auto-chain`, chained/stacked PR slice) or maintainer `size:exception` is recorded.
4. **Resolver Invalidation Loop**: Subagents report `skill_resolution` mode (`paths-injected`, `fallback-registry`, `fallback-path`, `none`). Any non-injected return purges orchestrator session cache and forces registry reload.

## Architecture Decisions

| Decision | Options Considered | Tradeoffs | Selected Decision |
|---|---|---|---|
| **Task Skills Injection** | A: Ad-hoc subagent self-discovery<br>B: Hardcoded skill file paths<br>C: Mandatory registry name resolution | A causes token waste and missed rules.<br>B breaks across install layouts.<br>C ensures deterministic path loading. | **C: Mandatory registry resolution** of `work-unit-commits` and `chained-pr` by canonical name injected as exact paths. |
| **Apply Batch Boundary** | A: Full tasks.md delegation<br>B: Phase-based batching<br>C: Single Work Unit batch (<=400 lines) | A causes oversized PRs and cognitive fatigue.<br>B still permits large multi-file diffs.<br>C guarantees bounded review slices. | **C: Single Work Unit batch (<=400 lines)** with focused test, runtime harness, and rollback boundaries. |
| **Workload Risk Enforcement** | A: Warning-only advisory<br>B: Conversational prompt fallback<br>C: Fail-closed gate at Step 2a & Orchestrator | A allows oversized diffs to slip to PR.<br>B introduces ambiguous execution state.<br>C guarantees safe delivery before edits. | **C: Fail-closed gate** halting orchestrator dispatch and `sdd-apply` Step 2a on unmanaged risk. |
| **Skill Cache Invalidation** | A: Static session cache<br>B: Time-based TTL<br>C: Feedback loop on non-injected status | A risks stale/missing skill paths.<br>B causes unpredictable reload overhead.<br>C reacts immediately on actual fallback. | **C: Event-driven cache invalidation** whenever subagent returns non-`paths-injected` resolution. |

## Data Flow

```text
[Orchestrator] ──(Inject work-unit-commits & chained-pr)──► [sdd-tasks]
      │                                                          │
      │ ◄──(Review Workload Forecast + Plain-Text Guard Lines)───┘
      ▼
[Workload Risk Gate] ──[High Risk & Unmanaged]──► [Fail-Closed STOP]
      │
      ├──[Auto-chain / Slice / size:exception]
      ▼
[sdd-apply] (Work Unit <=400 lines) ──(Step 2a Guard)──► [Code & Test Edits]
      │
      ▼
[Return Envelope: skill_resolution] ──[non-paths-injected]──► [Invalidate Cache]
```

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `internal/assets/antigravity/sdd-orchestrator.md` | Modify | Mandate task skills, single work-unit apply batch, fail-closed guard, resolver reload |
| `internal/assets/claude/sdd-orchestrator-workflow.md` | Modify | Mandate task skills, single work-unit apply batch, fail-closed guard, resolver reload |
| `internal/assets/codex/sdd-orchestrator.md` | Modify | Mandate task skills, single work-unit apply batch, fail-closed guard, resolver reload |
| `internal/assets/cursor/sdd-orchestrator.md` | Modify | Mandate task skills, single work-unit apply batch, fail-closed guard, resolver reload |
| `internal/assets/gemini/sdd-orchestrator.md` | Modify | Mandate task skills, single work-unit apply batch, fail-closed guard, resolver reload |
| `internal/assets/generic/sdd-orchestrator.md` | Modify | Mandate task skills, single work-unit apply batch, fail-closed guard, resolver reload |
| `internal/assets/hermes/sdd-orchestrator.md` | Modify | Mandate task skills, single work-unit apply batch, fail-closed guard, resolver reload |
| `internal/assets/kimi/sdd-orchestrator.md` | Modify | Mandate task skills, single work-unit apply batch, fail-closed guard, resolver reload |
| `internal/assets/kiro/sdd-orchestrator.md` | Modify | Mandate task skills, single work-unit apply batch, fail-closed guard, resolver reload |
| `internal/assets/opencode/sdd-orchestrator.md` | Modify | Mandate task skills, single work-unit apply batch, fail-closed guard, resolver reload |
| `internal/assets/qwen/sdd-orchestrator.md` | Modify | Mandate task skills, single work-unit apply batch, fail-closed guard, resolver reload |
| `internal/assets/windsurf/sdd-orchestrator.md` | Modify | Mandate task skills, single work-unit apply batch, fail-closed guard, resolver reload |
| `internal/assets/skills/_shared/skill-resolver.md` | Modify | Universal exact-path injection and delegator cache invalidation protocol |
| `internal/assets/skills/_shared/sdd-phase-common.md` | Modify | Workload guard line matching and skill resolution status contract |
| `internal/assets/skills/sdd-tasks/SKILL.md` | Modify | Workload forecast plain-text guard contract and work-unit structuring |
| `internal/assets/skills/sdd-apply/SKILL.md` | Modify | Step 2a fail-closed stop on unmanaged risk and single work unit execution |
| `internal/assets/assets_test.go` | Modify | Parity tests for task skills, batch bounds, fail-closed guards, cache invalidation |
| `testdata/golden/*` | Modify | Regenerated golden fixtures for all affected templates and skills |

## Interfaces / Contracts

### 1. Task Delegation Prompt Contract
```markdown
## Skills to load before work
Read these exact files before reading, writing, reviewing, testing, or creating artifacts:
- /path/to/skills/work-unit-commits/SKILL.md
- /path/to/skills/chained-pr/SKILL.md
```

### 2. Workload Forecast Guard Lines (`tasks.md`)
```text
Decision needed before apply: Yes|No
Chained PRs recommended: Yes|No
Chain strategy: stacked-to-main|feature-branch-chain|size-exception|pending
400-line budget risk: Low|Medium|High
```

### 3. Work Unit Evidence Table Contract (`apply-progress`)
| Evidence | Required value |
|---|---|
| Focused test command and exact result | Smallest command proving this unit; command, exit/result, relevant counts |
| Runtime harness command/scenario and exact result | Real integration/runtime path; explicit `N/A` with reason if no runtime boundary |
| Rollback boundary | Exact files/behavior removable without reverting unrelated work |

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit (Assets Parity) | All 12 orchestrators mandate `work-unit-commits` and `chained-pr` in `sdd-tasks` | Table-driven string assertions in `assets_test.go` |
| Unit (Assets Parity) | All 12 orchestrators enforce single Work Unit apply (<=400 lines) and fail-closed stop | Table-driven string assertions in `assets_test.go` |
| Unit (Assets Parity) | Universal cache invalidation on non-injected resolution in orchestrators & resolver | Table-driven assertions in `assets_test.go` |
| Integration (Goldens) | End-to-end template generation consistency across all platforms | Golden test suite in `internal/components/golden_test.go` |

## Threat Matrix

| Boundary | Minimum Adversarial Cases | Applicability | Design Response | Planned RED Tests |
|---|---|---|---|---|
| Documentation-like paths | `requirements.txt`, executable Markdown | N/A | No executable classification changes | None |
| Git repository selection | `git -C`, relative/absolute paths | N/A | No repository routing alterations | None |
| Commit state | Staged, empty index, multi-unit conflation | **Applicable** | Work units mandate atomic commits <=400 lines with rollback boundaries | `assets_test.go` asserts work-unit commit guidance and <=400 line bounds across all templates |
| Push state | Tracking branch, first push | N/A | No direct push mechanisms modified | None |
| PR commands | Target branches, chain strategies (`stacked-to-main`, `feature-branch-chain`) | **Applicable** | PR chain strategy and fail-closed workload stop prevent oversized or mistargeted PRs | `assets_test.go` asserts fail-closed blocking and chain strategy prompt contracts |

## Migration / Rollout

No migration or database schema change required. Templates and skills will be deployed via standard `gentle-ai update` or asset re-injection.

## Open Questions

- None.
