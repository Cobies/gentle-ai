# Tasks: Antigravity Dynamic Runtime Hardening

## Review Workload Forecast

| Field | Value |
|-------|-------|
| Estimated changed lines | 40–80 |
| 400-line budget risk | Low |
| Chained PRs recommended | No |
| Suggested split | Single PR |
| Delivery strategy | single-pr |
| Chain strategy | pending |

Decision needed before apply: No
Chained PRs recommended: No
Chain strategy: pending
400-line budget risk: Low

### Suggested Work Units

| Unit | Goal | Likely PR | Focused test command | Runtime harness | Rollback boundary |
|------|------|-----------|----------------------|-----------------|-------------------|
| 1 | Hardening message and test verification | PR 1 | `go test ./internal/components/sdd` | N/A | [antigravity_sdd_agents.go](file:///mnt/c/Users/Cobies/Desktop/Proyectos/COLAB/gentle-ai/internal/components/sdd/antigravity_sdd_agents.go) |

## Phase 1: Infrastructure & Test Setup

- [x] 1.1 RED: Update `TestAntigravitySddAgentsHardeningContractPhrases` in [antigravity_sdd_agents_test.go](file:///mnt/c/Users/Cobies/Desktop/Proyectos/COLAB/gentle-ai/internal/components/sdd/antigravity_sdd_agents_test.go) to assert that the hardening message contains the phrases: `"Strict TDD"`, `"Red phase"`, and `"Red-Green-Refactor"`.
- [x] 1.2 RED: Add a test case in [antigravity_sdd_agents_test.go](file:///mnt/c/Users/Cobies/Desktop/Proyectos/COLAB/gentle-ai/internal/components/sdd/antigravity_sdd_agents_test.go) verifying all dynamic 4R and Judgment Day roles are present and correctly mapped in the canonical role-scope validation table.

## Phase 2: Core Implementation

- [x] 2.1 GREEN: Add the Strict TDD enforcement rules text to the `antigravitySddAgentsHardeningMessage` constant in [antigravity_sdd_agents.go](file:///mnt/c/Users/Cobies/Desktop/Proyectos/COLAB/gentle-ai/internal/components/sdd/antigravity_sdd_agents.go).
- [x] 2.2 GREEN: Add the TDD assertion phrases (`"Strict TDD"`, `"Red phase"`, `"Red-Green-Refactor"`) to `antigravitySddAgentsHardeningContractPhrases` in [antigravity_sdd_agents.go](file:///mnt/c/Users/Cobies/Desktop/Proyectos/COLAB/gentle-ai/internal/components/sdd/antigravity_sdd_agents.go).
- [x] 2.3 REFACTOR: Refactor `antigravitySddAgentsHardeningMessage` formatting to ensure no single quotes are present and lines are cleanly structured.

## Phase 3: Verification

- [x] 3.1 Verify formatting: run `gofmt -s -w` on `internal/components/sdd/`.
- [x] 3.2 Run local component tests: `go test -v ./internal/components/sdd` to ensure all tests pass.
- [x] 3.3 Run full workspace verification: `go test ./...` and `go vet ./...` to verify no regressions across other agents.
