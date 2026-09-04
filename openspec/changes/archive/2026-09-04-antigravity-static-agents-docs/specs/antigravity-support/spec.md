# Delta for antigravity-support

## MODIFIED Requirements

### Requirement: Antigravity uses static-primary subagents with dynamic fallback

The Antigravity orchestrator MUST delegate SDD phases via `invoke_subagent` against the pre-registered static subagent set as the primary route and SHALL use runtime `define_subagent` creation only as a resilient fallback. When running on a low-tier model, the system MUST enforce subagent delegation and MUST NOT execute SDD phases (such as explore, propose, spec, design, tasks, apply, verify) inline.

(Previously: dynamic-only delegation via `define_subagent` then `invoke_subagent`, with no static subagent files.)

#### Scenario: SDD orchestration runs in Antigravity

- GIVEN the Antigravity SDD orchestrator is installed with the static subagent set registered
- WHEN an SDD phase requires a subagent
- THEN the prompt instructs Antigravity to call `invoke_subagent` against the pre-registered set
- AND names `define_subagent` only as fallback when a static agent is unavailable.

#### Scenario: Low-model subagent enforcement

- GIVEN a low-tier model is active in the `antigravity` agent CLI
- WHEN the orchestrator compiles system instructions
- THEN the prompt MUST require delegating each phase (`sdd-explore`, `sdd-propose`, `sdd-spec`, `sdd-design`, `sdd-tasks`, `sdd-apply`, `sdd-verify`) via `invoke_subagent`, with `define_subagent` as fallback
- AND the prompt MUST explicitly forbid inline phase execution.

## ADDED Requirements

### Requirement: Documentation states the two-tier model consistently

Every refreshed surface MUST describe Antigravity as static-primary with `define_subagent` fallback and MUST NOT claim solo-agent operation, dynamic-only delegation, or unavailable custom sub-agent creation.

#### Scenario: Matrix rows and Mission Control section agree

- GIVEN a reader opens `docs/agents.md`, `README.md`, or the Mission Control section
- WHEN they read the Antigravity delegation description
- THEN each surface states static-primary delegation with dynamic fallback
- AND none claims solo-agent or dynamic-only behavior.

#### Scenario: Code comments match runtime behavior

- GIVEN a reader opens `internal/components/sdd/antigravity_sdd_agents.go` or `internal/assets/antigravity/sdd-orchestrator.md`
- WHEN they read the header comment and delegation lines
- THEN the text describes `invoke_subagent` against the static set as primary with `define_subagent` fallback
- AND no line claims Antigravity has no static registry or runs dynamic-only.

### Requirement: Legacy workaround doc is a dated historical note

`docs/antigravity-sdd-workaround.md` MUST be a dated historical note that points readers to `docs/agents.md`, MUST NOT prescribe inline single-threaded execution, and old links to it MUST still resolve.

#### Scenario: Reader lands on the workaround doc

- GIVEN a reader opens `docs/antigravity-sdd-workaround.md` (including via an old link)
- WHEN they read the page
- THEN they see a dated notice marking it historical with a pointer to `docs/agents.md`
- AND no instruction to execute SDD phases inline remains normative.

### Requirement: Prose-only scope guards

The change MUST NOT alter runtime behavior, MUST keep the shared `GEMINI.md` collision warning valid, and MUST NOT claim the hardening plugin covers the desktop variant. Installer and probe hardening are out of scope.

#### Scenario: Guard rails hold after refresh

- GIVEN the refreshed docs, spec, and comments
- WHEN a reader checks behavior claims, the collision warning, and hardening coverage
- THEN no new runtime behavior is described, the `GEMINI.md` warning is intact, and hardening stays CLI-scoped.

## RENAMED Requirements

### Requirement: Antigravity uses dynamic subagents → Antigravity uses static-primary subagents with dynamic fallback

(Reason: the old name states the superseded dynamic-only model.)
(Migration: update spec references, tests, and docs to the new requirement name.)

## Assumptions

- No external corroboration: research rev2 blocked (`documentation=[]`, `open-web=[]`); claims rest on in-repo evidence (`exploration.md`, commit `398b0f51`, memory #5315 superseding #4470/#4520).
