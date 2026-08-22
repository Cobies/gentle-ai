# Proposal: Harden Skill Registry and Workload Guards

## Intent

Enforce reviewer workload safety and skill resolution determinism across the SDD lifecycle by mandating `work-unit-commits` and `chained-pr` injection during `sdd-tasks` delegation, bounding `sdd-apply` to single Work Unit batches (<=400 lines) with fail-closed unmanaged risk guards, and hardening skill resolution cache invalidation across all 12 orchestrator templates.

## Scope

### In Scope
- Mandate passing `work-unit-commits` and `chained-pr` when delegating `sdd-tasks` across all 12 orchestrator templates.
- Enforce bounding `sdd-apply` to single Work Unit batches (<=400 lines) and fail-closed blocking when workload risk is unmanaged.
- Harden universal skill resolver protocol (`internal/assets/skills/_shared/skill-resolver.md`) for exact path injection and cache invalidation on non-injected resolutions.
- Update parity test assertions in `internal/assets/assets_test.go` and regenerate golden fixtures.

### Out of Scope
- Changes to git commit creation primitives or low-level VCS hooks.
- Altering the 400-line budget calculation formula or excluding non-golden files.
- Introducing new delivery strategies outside the 4 canonical values (`ask-on-risk`, `auto-chain`, `single-pr`, `exception-ok`).

## Capabilities

### New Capabilities
None.

### Modified Capabilities
- `sdd-orchestrator-assets`: Enforce mandatory `work-unit-commits` and `chained-pr` skill injection in `sdd-tasks`, single Work Unit batch constraints (<=400 lines) and fail-closed unmanaged risk blocking in `sdd-apply`, and tightened skill resolver cache protocols.

## Approach

Update all 12 SDD orchestrator templates and the shared `skill-resolver.md` to guarantee that every `sdd-tasks` launch injects `work-unit-commits` and `chained-pr`, and every `sdd-apply` execution is bounded to <=400 lines per Work Unit batch with fail-closed rejection on unmanaged risk. Refresh static parity checks in `assets_test.go` and regenerate golden fixtures across all agent runtimes.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/assets/templates/sdd-orchestrator/*.md` | Modified | All 12 orchestrators: inject tasks skills and bound apply work units |
| `internal/assets/skills/_shared/skill-resolver.md` | Modified | Exact path injection & cache invalidation loop protocol |
| `internal/assets/skills/sdd-tasks/SKILL.md` | Modified | Workload forecast and work-unit commit structuring |
| `internal/assets/skills/sdd-apply/SKILL.md` | Modified | Single Work Unit batch enforcement and fail-closed blocking |
| `internal/assets/assets_test.go` | Modified | Parity assertions for work-unit skills, batch bounds, and resolver |
| `internal/assets/testdata/golden/*` | Modified | Regenerated fixtures reflecting hardened guidance |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Orchestrator template divergence across 12 targets | Medium | Comprehensive table-driven parity tests in `assets_test.go` covering all agent families |
| Oversized tasks slip past unmanaged risk guards | Low | Hard fail-closed stop in `sdd-apply` unless explicit exception or chained slice |

## Rollback Plan

Revert template and skill markdown changes together with test updates and regenerated golden fixtures.

## Dependencies

- Existing canonical 4-value delivery strategy domain (`ask-on-risk`, `auto-chain`, `single-pr`, `exception-ok`).

## Success Criteria

- [ ] All 12 orchestrator templates instruct injecting `work-unit-commits` and `chained-pr` for `sdd-tasks`.
- [ ] All 12 orchestrator templates and `sdd-apply` enforce single Work Unit batches (<=400 lines) with fail-closed stop on unmanaged risk.
- [ ] Universal skill resolver mandates cache invalidation when sub-agent reports non-injected resolution.
- [ ] Parity test suite and golden fixtures pass without regression.
