# Proposal: Antigravity Static Agents Docs Refresh

## Intent

Antigravity already runs static-primary + dynamic-fallback (18 embedded files, `SupportsSubAgents()==true`, installer + test, orchestrator mandates `invoke_subagent` first), but docs/spec still claim dynamic-only/solo-agent. Users follow dead inline-execution guidance the orchestrator now refuses.

## Scope

### In Scope
- Refresh `docs/agents.md` matrix row, delegation table, Mission Control section to two-tier
- Refresh `README.md:111` row to static-primary + dynamic-fallback
- Update `openspec/specs/antigravity-support/spec.md:30-46` requirement + scenarios
- Fix stale header `internal/components/sdd/antigravity_sdd_agents.go:17-21` + ~4 stale orchestrator lines
- Convert `docs/antigravity-sdd-workaround.md` into dated historical note pointing at `docs/agents.md`

### Out of Scope
- Zero runtime behavior change; no installer/probe hardening (separate future change)
- No external Antigravity docs corroboration (deselected; recorded as risk)

## Capabilities

### New Capabilities
- None — prose-only refresh, no new spec-level behavior.

### Modified Capabilities
- `antigravity-support`: requirement changes from dynamic-only (`define_subagent` then `invoke_subagent`) to static-primary (`invoke_subagent` against pre-registered set) with dynamic creation as resilient fallback.

## Approach

Prose-only edits on the six surfaces above, citing commit `398b0f51` and memory #5315 (supersedes #4470/#4520). Keep `GEMINI.md` collision warning valid; do not claim the hardening plugin covers desktop.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `docs/agents.md` | Modified | Matrix, delegation table, Mission Control section |
| `README.md:111` | Modified | Matrix row wording |
| `openspec/specs/antigravity-support/spec.md:30-46` | Modified | Requirement + 2 scenarios |
| `docs/antigravity-sdd-workaround.md` | Modified | Historical note + redirect |
| `internal/components/sdd/antigravity_sdd_agents.go:17-21` | Modified | Stale header comment only |
| `internal/assets/antigravity/sdd-orchestrator.md` | Modified | ~4 stale lines; no behavior |
| `internal/agents/antigravity/` | None | Reference only; no code change |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| No external corroboration (research rev2 blocked; runtime `documentation=[]`, `open-web=[]`) | High | Scope claims to in-repo evidence; cite exploration.md |
| Implying Gemini CLI inherits static agents via shared `GEMINI.md` | Low | Keep collision warning verbatim |
| Implying hardening plugin covers desktop variant | Low | Preserve CLI-only wording |

## Rollback Plan

Revert the single docs/spec/comment commit; restore workaround file from git. No migration, no data, no runtime state.

## Dependencies

- None. Evidence: `exploration.md`, commit `398b0f51` + follow-ups, Engram #5315. `research.md` rev2 (blocked) cited as known gap only.

## Success Criteria

- [ ] No surface claims Antigravity is solo-agent, dynamic-only, or lacks custom sub-agent creation
- [ ] Every touched surface states static-primary + `define_subagent` fallback consistently
- [ ] Workaround doc is dated history with pointer; old links resolve
- [ ] `go test ./internal/components/sdd/... ./internal/agents/antigravity/...` passes (no behavior change)
