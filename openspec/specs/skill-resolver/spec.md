# Delta for skill-resolver

## ADDED Requirements

### Requirement: Strict Exact-Path Injection Protocol

Delegating agents MUST resolve skills from the skill registry by canonical name and inject absolute filesystem paths to `SKILL.md` files under `## Skills to load before work`. Delegators MUST NOT inject inline summaries or hardcoded fallback paths.

#### Scenario: Delegator passes resolved skill paths
- GIVEN matched skills for a delegation task
- WHEN the sub-agent prompt is generated
- THEN each matched skill is represented as an exact absolute path to its `SKILL.md`
- AND no summarized or rewritten rules are passed in place of skill paths

### Requirement: Delegator Cache Invalidation Contract

Sub-agents MUST report their `skill_resolution` mode in their return envelope. When `skill_resolution` is `fallback-registry`, `fallback-path`, or `none`, the delegator MUST invalidate its session cache and force a registry reload.

#### Scenario: Fallback resolution triggers registry reload
- GIVEN a sub-agent completes work using `fallback-registry`
- WHEN the delegator processes the result envelope
- THEN the delegator purges its cached skill registry
- AND subsequent delegations reload the registry before matching skills
