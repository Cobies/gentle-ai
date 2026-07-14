# Proposal: Robust Sync Pipeline

## Intent
Make the `sync` command resilient by allowing it to gather warnings/errors per component and continue executing remaining components, rather than aborting the entire sync pipeline immediately on a single component or file verification failure.

## Scope

### In Scope
- Modify `pipeline.Runner` and `pipeline.Orchestrator` to support per-step rollback on failure during execution under the `ContinueOnError` policy.
- Configure `RunSyncWithSelection` to run the sync pipeline with `ContinueOnError`.
- Change `runPostSyncVerification` to return verification failures as warnings instead of failing the command.
- Display a clean summary of successes, warnings, and skipped items at the end of the sync run.
- Keep exit code as 0 for non-critical component sync/verification failures; return non-zero only for global system failures.

### Out of Scope
- Changing the install or upgrade pipelines to use non-aborting policies.
- Retrying failed components automatically.

## Capabilities

### New Capabilities
- `robust-sync`: Continues sync pipeline execution on component failures, rolling back only the failed components, and reporting verification failures as warnings.

### Modified Capabilities
None

## Approach
1. **Pipeline Execution Update**:
   - Update `pipeline.Runner` so that when `FailurePolicy == ContinueOnError`, it catches any step failure, runs that step's rollback handler immediately to clean up its partial work, marks the step as failed, and continues.
   - Update `pipeline.Orchestrator` to not run a global rollback if `ContinueOnError` was used and there were step failures.
2. **Sync Command Integration**:
   - In `RunSyncWithSelection`, pass `WithFailurePolicy(pipeline.ContinueOnError)` to `NewOrchestrator`.
   - Update output logic to aggregate all step results and show a clear summary.
3. **Verification Update**:
   - Update `RunSyncWithSelection` post-sync verification step: instead of throwing an error when `!result.Verify.Ready`, report the verification details as warnings.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/pipeline/runner.go` | Modified | Support per-step rollback during `ContinueOnError` |
| `internal/pipeline/orchestrator.go` | Modified | Avoid global rollback for `ContinueOnError` |
| `internal/cli/sync.go` | Modified | Configure pipeline to continue on error, convert verification failure to warnings, and format output summary |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Succeeded components mask errors | Low | Clear CLI summary highlighting warnings/errors |
| Partial updates cause inconsistency | Low | Transactional rollback per failed component cleans up partial state |

## Rollback Plan
Revert changes in `internal/pipeline` and `internal/cli/sync.go` to return to default `StopOnError` policy and strict post-sync verification.

## Dependencies
None

## Success Criteria
- [ ] Sync command runs to completion and exits with code 0 even if a component (e.g., Engram) fails.
- [ ] Failed components are rolled back individually; other components retain updates.
- [ ] Post-sync verification failures are reported as warnings, not command aborts.
- [ ] CLI/TUI displays a summary table of successes, warnings, and skipped items.
