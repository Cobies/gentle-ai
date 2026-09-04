# Tasks: Antigravity Static Agents Docs Refresh

## Review Workload Forecast

| Field | Value |
|-------|-------|
| Estimated changed lines | 40–80 |
| 400-line budget risk | Low |
| Chained PRs recommended | No |
| Suggested split | Single PR |
| Delivery strategy | ask-on-risk |
| Chain strategy | pending |

Decision needed before apply: No
Chained PRs recommended: No
Chain strategy: pending
400-line budget risk: Low

### Suggested Work Units

| Unit | Goal | Likely PR | Focused test command | Runtime harness | Rollback boundary |
|------|------|-----------|----------------------|-----------------|-------------------|
| 1 | Prose-only two-tier refresh on all 6 surfaces + verification | Single PR | `go test ./internal/components/sdd/... ./internal/agents/antigravity/...` | N/A — prose-only, no runtime path exercised | Revert single commit; workaround file restores from git |

## Phase 1: Docs Consistency

- [x] 1.1 Edit `docs/agents.md:19` matrix cell to `Full (static subagents) + Mission Control`
- [x] 1.2 Move Antigravity row in `docs/agents.md:40,43` from `Solo-agent` to `Full (sub-agents)`
- [x] 1.3 Rewrite `docs/agents.md:62-64` delegation wording to static-primary `invoke_subagent`, `define_subagent` fallback, inline forbidden
- [x] 1.4 Update `docs/agents.md:177-183` with static path (`agents/` under `~/.gemini/antigravity-cli/`); keep `GEMINI.md` warning verbatim
- [x] 1.5 Edit `README.md:111` row to `Full (static subagents + dynamic fallback)` + pre-registered `invoke_subagent` wording

## Phase 2: Spec Promotion + Historical Note

- [x] 2.1 Promote `openspec/specs/antigravity-support/spec.md:30-46` to static-primary requirement + both scenarios per delta spec
- [x] 2.2 Convert `docs/antigravity-sdd-workaround.md` to dated historical banner pointing to `docs/agents.md`; demote inline rules to non-normative

## Phase 3: Code-Adjacent Comments

- [x] 3.1 Fix header `internal/components/sdd/antigravity_sdd_agents.go:17-21` to two-tier wording; keep hardening + CLI-only text intact
- [x] 3.2 Retouch ~4 stale lines in `internal/assets/antigravity/sdd-orchestrator.md` (L76, L98, L170, fail-closed noun); delegation logic byte-identical

## Phase 4: Verification

- [x] 4.1 Run `go test ./internal/components/sdd/... ./internal/agents/antigravity/...` — must pass unchanged (no behavior change)
- [x] 4.2 Run zero-hit `rg` stale-claim check over touched surfaces (solo-agent, dynamic-only, no-static-registry, not-yet-available, does-not-install-static, dynamic-subagents)
- [x] 4.3 Open `docs/antigravity-sdd-workaround.md` via old path; confirm banner + redirect present and `GEMINI.md` warning intact
