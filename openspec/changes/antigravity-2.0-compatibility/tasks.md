# Tasks: Antigravity 2.0 and CLI Compatibility

## Review Workload Forecast

| Field | Value |
|-------|-------|
| Estimated changed lines | 300-350 lines |
| 400-line budget risk | Low |
| Chained PRs recommended | No |
| Suggested split | Single PR |
| Delivery strategy | ask-on-risk |
| Chain strategy | stacked-to-main |

Decision needed before apply: No
Chained PRs recommended: No
Chain strategy: stacked-to-main
400-line budget risk: Low

### Suggested Work Units

| Unit | Goal | Likely PR | Notes |
|------|------|-----------|-------|
| 1 | Register agent and implement adapter | PR 1 | Base branch; tests/docs included |

## Phase 1: Foundation

- [ ] 1.1 Add `AgentAntigravityCLI` constant to `internal/model/types.go`
- [ ] 1.2 Catalog `antigravity-cli` in `internal/catalog/agents.go`
- [ ] 1.3 Register `antigravity-cli` in `internal/agents/factory.go`

## Phase 2: Core Implementation

- [ ] 2.1 Implement `internal/agents/antigravitycli/adapter.go` with settings in `~/.gemini/antigravity-cli/`
- [ ] 2.2 Add `internal/assets/antigravitycli/sdd-orchestrator.md` with dynamic subagent delegation prompts
- [ ] 2.3 Map `antigravity-cli` to `antigravitycli/sdd-orchestrator.md` in `internal/components/sdd/inject.go`
- [ ] 2.4 Update `internal/cli/run.go`, `internal/cli/validate.go`, and `internal/tui/model.go` for `antigravity-cli` mapping

## Phase 3: Testing & Verification

- [ ] 3.1 Update `internal/cli/install_test.go` and factory tests to cover `antigravity-cli`
- [ ] 3.2 Run tests via `go test ./...` and verify clean execution
