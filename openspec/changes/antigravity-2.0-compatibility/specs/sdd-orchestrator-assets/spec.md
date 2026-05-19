# Delta for sdd-orchestrator-assets

## ADDED Requirements

### Requirement: Antigravity CLI Dynamic Delegation Semantics

The `antigravity-cli` SDD orchestrator asset MUST instruct the agent to use dynamic subagent definition (`define_subagent`) and invocation (`invoke_subagent`).

The instructions MUST:
- Direct the agent to read skill files dynamically from `~/.gemini/antigravity-cli/skills/{phase}/SKILL.md` (or workspace `.agents/skills/{phase}/SKILL.md`).
- Define the subagent dynamically using the `define_subagent` tool with the skill content as the system prompt.
- Invoke the subagent using the `invoke_subagent` tool.
- Advise the agent of the strict nesting depth limit of 10 levels.

#### Scenario: Orchestrator asset for antigravity-cli includes dynamic subagent delegation

- GIVEN the `antigravity-cli` orchestrator asset is reviewed
- WHEN delegation instructions are checked
- THEN it includes guidelines on calling `define_subagent` and `invoke_subagent`
- AND it specifies reading skills from the global CLI skills directory or workspace `.agents/skills`
- AND it enforces the subagent nesting depth limit of 10 levels
