# Tasks: Robust Sync Pipeline

Decision needed before apply: No
Chained PRs recommended: No
Chain strategy: pending
400-line budget risk: Low

## Review Workload Forecast

| Field | Value |
|-------|-------|
| Estimated changed lines | 150-250 |
| 400-line budget risk | Low |
| Chained PRs recommended | No |
| Suggested split | Single PR |
| Delivery strategy | ask-on-risk |
| Chain strategy | pending |

### Suggested Work Units

| Unit | Goal | Likely PR | Focused test command | Runtime harness | Rollback boundary |
|------|------|-----------|----------------------|-----------------|-------------------|
| 1 | Robust pipeline runner & sync integration | PR 1 | `go test ./internal/pipeline/... ./internal/cli/...` | `gentle-ai sync --dry-run` | `internal/pipeline/`, `internal/cli/sync.go` |

## Phase 1: Pipeline Engine Enhancements (TDD)

- [x] 1.1 RED: Add tests in [orchestrator_test.go](file:///mnt/c/Users/Cobies/Desktop/Proyectos/COLAB/gentle-ai/internal/pipeline/orchestrator_test.go) asserting step-level rollback and no global rollback under `ContinueOnError`.
- [x] 1.2 GREEN: Update `Runner.Run` in [runner.go](file:///mnt/c/Users/Cobies/Desktop/Proyectos/COLAB/gentle-ai/internal/pipeline/runner.go) to invoke `Rollback()` on failed steps immediately if `FailurePolicy == ContinueOnError`.
- [x] 1.3 GREEN: Update `Orchestrator.Execute` in [orchestrator.go](file:///mnt/c/Users/Cobies/Desktop/Proyectos/COLAB/gentle-ai/internal/pipeline/orchestrator.go) to skip global rollback under `ContinueOnError`.
- [x] 1.4 REFACTOR: Clean up pipeline conditional branches and verify all pipeline unit tests pass green.

## Phase 2: CLI Sync Command Integration (TDD)

- [x] 2.1 RED: Add integration tests in [sync_test.go](file:///mnt/c/Users/Cobies/Desktop/Proyectos/COLAB/gentle-ai/internal/cli/sync_test.go) simulating component/verification failures and asserting exit code 0.
- [x] 2.2 GREEN: In [sync.go](file:///mnt/c/Users/Cobies/Desktop/Proyectos/COLAB/gentle-ai/internal/cli/sync.go), pass `ContinueOnError` to orchestrator and handle step failures gracefully without command abort.
- [x] 2.3 GREEN: In [sync.go](file:///mnt/c/Users/Cobies/Desktop/Proyectos/COLAB/gentle-ai/internal/cli/sync.go), convert verification failure to warnings rather than returning an error.
- [x] 2.4 GREEN: In [sync.go](file:///mnt/c/Users/Cobies/Desktop/Proyectos/COLAB/gentle-ai/internal/cli/sync.go), render summary table highlighting Success, Failed, Warning, and Skipped states.
- [x] 2.5 REFACTOR: Clean up CLI output rendering code and run full test suite to verify green.
