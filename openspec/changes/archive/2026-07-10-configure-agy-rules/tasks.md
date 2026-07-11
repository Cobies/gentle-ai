# Tasks: Configure Antigravity Workspace Rules

Decision needed before apply: No
Chained PRs recommended: No
Chain strategy: stacked-to-main
400-line budget risk: Medium

## Review Workload Forecast

| Field | Value |
|-------|-------|
| Estimated changed lines | 250-350 |
| 400-line budget risk | Medium |
| Chained PRs recommended | No |
| Suggested split | Single PR |
| Delivery strategy | ask-on-risk |
| Chain strategy | stacked-to-main |

## Phase 1: Foundation & Adapter Changes

- [x] 1.1 Add `DetectLowModel(homeDir string) bool` to `internal/agents/antigravity/adapter.go` checking env vars and `settings.json`.
- [x] 1.2 Add `GetWorkspaceRules(cwd string) (string, error)` to `internal/agents/antigravity/adapter.go` to load `.agents/rules/sdd-workflow.md`.
- [x] 1.3 Add unit tests for `DetectLowModel` in `internal/agents/antigravity/adapter_test.go` covering env vars and settings.json mocks.
- [x] 1.4 Add unit tests for `GetWorkspaceRules` in `internal/agents/antigravity/adapter_test.go` covering existing/missing rules files.
  - Test command: `go test -v ./internal/agents/antigravity -run TestAdapter_`

## Phase 2: Core/Prompt Sync Changes

- [x] 2.1 Update `internal/components/persona/inject.go` in the `StrategyAppendToFile` case (used by Antigravity) to load workspace rules.
- [x] 2.2 Prepend the low-tier model warning markdown block to workspace rules if `DetectLowModel` is true.
- [x] 2.3 Inject the compiled workspace rules into `~/.gemini/GEMINI.md` within `<!-- gentle-ai:workspace-rules -->` and `<!-- /gentle-ai:workspace-rules -->` markers.
- [x] 2.4 Add unit tests in `internal/components/persona/inject_test.go` to verify rule and warning injection for Antigravity.
  - Test command: `go test -v ./internal/components/persona -run TestInject_`

## Phase 3: Testing & Verification

- [x] 3.1 Verify sync integration by running `gentle-ai sync` and inspecting the generated `~/.gemini/GEMINI.md`.
  - Verification: Check markers and content under real and mock workspace structures.
