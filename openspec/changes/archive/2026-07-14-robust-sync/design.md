# Design: Robust Sync Pipeline

## Technical Approach

Make the `sync` command resilient by allowing it to gather warnings/errors per component and continue executing remaining components, rather than aborting the entire sync pipeline immediately on a single component or file verification failure.

- Modify `pipeline.Runner` to support per-step rollback on failure during execution under the `ContinueOnError` policy.
- Configure the Orchestrator with `ContinueOnError` and disable the global rollback when this policy is active.
- Change the `sync` command (`RunSyncWithSelection`) to run the sync pipeline with `ContinueOnError` and report verification failures as warnings instead of failing the command.
- Display a clean summary of successes, warnings, and skipped items at the end of the sync run.

## Architecture Decisions

| Option | Tradeoff | Decision |
|--------|----------|----------|
| Immediate per-step rollback under `ContinueOnError` | Cleans up partial work immediately but requires immediate detection. | Run step's `Rollback()` immediately if it implements `RollbackStep` when a step fails under `ContinueOnError`. |
| Skip global rollback under `ContinueOnError` | Avoids reverting successful steps, which aligns with robust sync but leaves succeeded steps intact. | If policy is `ContinueOnError`, bypass the orchestrator's global rollback. |
| Zero exit code on step/verification failures | May mask issues if the output summary is not read, but ensures automation does not break on non-critical component sync warnings. | Return `nil` error (exit 0) from the sync command for component execution and verification failures. |

## Data Flow

```
   Sync Command (cli.RunSync)
         │
         ▼
   Initialize Orchestrator (with ContinueOnError)
         │
         ▼
   Execute Stages:
     1. Prepare (Failures abort execution)
     2. Apply (Step Failure -> Run local Step Rollback -> Continue)
         │
         ▼
   Post-Sync Verification (Failures converted to Warnings)
         │
         ▼
   Render Summary Table & Exit (Exit Code 0)
```

## File Changes

| File | Action | Description |
|------|--------|-------------|
| [internal/pipeline/runner.go](file:///mnt/c/Users/Cobies/Desktop/Proyectos/COLAB/gentle-ai/internal/pipeline/runner.go) | Modify | Update `Runner.Run` to check for `ContinueOnError` and immediately execute the failed step's rollback handler if it implements `RollbackStep`. |
| [internal/pipeline/orchestrator.go](file:///mnt/c/Users/Cobies/Desktop/Proyectos/COLAB/gentle-ai/internal/pipeline/orchestrator.go) | Modify | Update `Orchestrator.Execute` to skip global rollback if `o.runner.FailurePolicy == ContinueOnError`. |
| [internal/cli/sync.go](file:///mnt/c/Users/Cobies/Desktop/Proyectos/COLAB/gentle-ai/internal/cli/sync.go) | Modify | Pass `WithFailurePolicy(pipeline.ContinueOnError)` to `NewOrchestrator`, do not abort the sync command on apply stage failures or verification failures, and display a summary table highlighting successes, warnings, and failures. |

## Interfaces / Contracts

No new interfaces or API contracts are introduced. The existing `pipeline.FailurePolicy` and `pipeline.RollbackStep` are utilized.

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit | Step failure under `ContinueOnError` triggers individual rollback and continues. | Add tests in `runner_test.go` and `orchestrator_test.go`. |
| Integration | Robust sync command execution under failures. | Add integration tests in `sync_test.go` simulating step failures and verification failures. |

## Threat Matrix

N/A — no routing, shell, subprocess, VCS/PR automation, executable-file classification, or process-integration boundary.

## Migration / Rollout

No migration required.
