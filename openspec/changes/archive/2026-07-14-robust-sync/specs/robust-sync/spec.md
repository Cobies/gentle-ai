# Robust Sync Specification

## Purpose

The robust sync pipeline allows the sync command to proceed when individual component sync or post-sync verifications fail. Instead of halting immediately, the pipeline executes all selected components, runs individual rollbacks for failed steps, collects all warnings/errors, and presents a comprehensive CLI/TUI summary at the end.

## Requirements

### Requirement: Continue On Error Policy

The pipeline runner MUST support a `ContinueOnError` failure policy. When this policy is active and a pipeline step fails during execution:
1. The runner MUST run the failed step's rollback handler immediately to clean up its partial work.
2. The runner MUST NOT execute a global rollback of previously succeeded steps.
3. The runner MUST proceed with executing any remaining steps in the pipeline.

#### Scenario: Step fails with ContinueOnError
- GIVEN a pipeline running with `ContinueOnError` policy and three steps: Step A, Step B, and Step C
- WHEN Step B fails during execution
- THEN Step B's rollback handler is executed immediately
- AND Step A's rollback handler is NOT executed
- AND Step C is executed normally

### Requirement: Post-Sync Verification as Warnings

The post-sync verification step MUST NOT abort the sync command if the system state is not ready. Instead:
1. The verification failures MUST be recorded as warnings.
2. The command MUST display these warnings at the end of execution.
3. The command MUST exit with code 0 if only component or post-sync verification failures occurred.

#### Scenario: Verification fails after successful components
- GIVEN a sync command execution where all components sync successfully but post-sync verification fails
- WHEN the sync command completes
- THEN the verification failure is reported as a warning
- AND the exit code of the sync command is 0

### Requirement: CLI/TUI Execution Summary

The CLI/TUI MUST display a summary table at the end of the sync run detailing the outcome of each step.
1. The table MUST list the status (e.g., Success, Failed, Warning, Skipped) for each component.
2. The summary MUST highlight warnings and errors to prevent them from being masked.

#### Scenario: Output summary with failures and warnings
- GIVEN a sync pipeline with Component A succeeding, Component B failing, and post-sync verification warning
- WHEN the sync command completes
- THEN the TUI/CLI displays a summary showing Component A as Success, Component B as Failed, and Verification as Warning
