# Delta for sdd-apply-workload-guard

## ADDED Requirements

### Requirement: Step 2a Fail-Closed Workload Execution Guard

Before reading implementation files or writing code, `sdd-apply` MUST inspect the workload forecast in `tasks.md`. If the forecast indicates `400-line budget risk: High` or `Chained PRs recommended: Yes` without an approved strategy (`auto-chain`, chained/stacked PR slice, or explicit `size:exception`), `sdd-apply` MUST fail closed and return `status: blocked`.

#### Scenario: Unapproved oversized work blocked at Step 2a
- GIVEN `tasks.md` with High budget risk and no approved delivery strategy
- WHEN `sdd-apply` reaches Step 2a
- THEN `sdd-apply` stops immediately before file modification
- AND returns `status: blocked` with a workload decision requirement

### Requirement: Single Work Unit Batch Enforcement

`sdd-apply` MUST enforce execution boundaries to a single assigned Work Unit <=400 lines unless maintainer `size:exception` is explicitly recorded.

#### Scenario: Single work unit batch accepted
- GIVEN an assigned work unit slice bounded to <=400 lines
- WHEN `sdd-apply` validates batch parameters
- THEN implementation proceeds for only that slice
- AND completed tasks are recorded in `apply-progress`
