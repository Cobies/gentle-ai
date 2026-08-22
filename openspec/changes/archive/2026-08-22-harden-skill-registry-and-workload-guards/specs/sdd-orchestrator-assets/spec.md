# Delta for sdd-orchestrator-assets

## ADDED Requirements

### Requirement: Mandatory Tasks Skill Injection

All 12 SDD orchestrator templates MUST resolve and inject both `work-unit-commits` and `chained-pr` skills under `## Skills to load before work` whenever delegating `sdd-tasks`.

#### Scenario: Orchestrator delegates sdd-tasks with required skills
- GIVEN an orchestrator prepares to delegate `sdd-tasks`
- WHEN the launch prompt is constructed
- THEN both `work-unit-commits` and `chained-pr` paths are injected under `## Skills to load before work`
- AND neither skill is omitted

### Requirement: Bounded Work Unit Apply Dispatch

The orchestrator MUST delegate `sdd-apply` strictly in single Work Unit batches bounded to <=400 changed lines (`additions + deletions`). The orchestrator MUST NOT dispatch multi-unit batches exceeding 400 lines in a single delegation.

#### Scenario: Orchestrator dispatches single work unit batch
- GIVEN task planning identified multiple work units exceeding 400 lines total
- WHEN `sdd-apply` is launched
- THEN the delegation prompt targets exactly one autonomous work unit <=400 lines
- AND subsequent units are deferred to future apply dispatches

### Requirement: Fail-Closed Unmanaged Workload Risk Blocking

When `sdd-tasks` reports `400-line budget risk: High`, `Chained PRs recommended: Yes`, or `Decision needed before apply: Yes`, the orchestrator MUST block fail-closed before delegating `sdd-apply` unless a canonical delivery strategy or `size:exception` has been explicitly resolved.

#### Scenario: Unmanaged workload risk halts apply launch
- GIVEN `sdd-tasks` forecasts High budget risk with no resolved delivery strategy
- WHEN the orchestrator evaluates the transition to `sdd-apply`
- THEN execution halts immediately with a blocking decision request
- AND no source files or apply agents are dispatched

### Requirement: Skill Registry Cache Invalidation Loop

When any sub-agent returns a `skill_resolution` value other than `paths-injected`, the orchestrator MUST invalidate its cached skill registry and re-read the registry from persistent store before the next delegation.

#### Scenario: Non-injected resolution triggers cache invalidation
- GIVEN a delegated sub-agent returns `skill_resolution: fallback-registry`
- WHEN the orchestrator receives the return envelope
- THEN the session skill cache is invalidated immediately
- AND the next delegation re-reads the registry from Engram or `.atl/skill-registry.md`
