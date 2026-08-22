# Skill Resolver — Universal Protocol

Any agent that **delegates work to sub-agents** MUST use this protocol to resolve relevant skills and pass them safely.

## Why This Exists

Sub-agents start with no project skill context. The registry gives delegators a cheap index of available skills without rewriting or summarizing those skills.

## When to Apply

Before every sub-agent launch that involves reading, writing, reviewing, testing, documenting, or creating project artifacts. Skip only for purely mechanical commands.

## The Protocol

### Step 1: Obtain the Skill Registry

The registry is an **index** of skill names, triggers, scopes, and exact `SKILL.md` paths. It is not a compact-rules bundle.

Resolution order:
1. Use the session cache if present.
2. `mem_search(query: "skill-registry", project: "{project}")` → `mem_get_observation(id)` for full content.
3. Fallback: read `.atl/skill-registry.md` from the project root.
4. No registry found → proceed without project skills and warn the user to run `gentle-ai skill-registry refresh`.

### Step 2: Match Relevant Skills

Match on two dimensions:

| Context | Match against |
| --- | --- |
| Code/files | Registry trigger/description mentions the language, framework, tool, or path context |
| Task/action | Registry trigger/description mentions actions like PR, review, docs, tests, Jira, comments, release |

Prefer the smallest useful set. If more than five skills match, keep the five most relevant and prioritize code context over task context.

### Step 3: Pass Skill Paths

Inject exact paths, not summaries:

```markdown
## Skills to load before work

Read these exact files before reading, writing, reviewing, testing, or creating artifacts:

- /absolute/path/to/skills/go-testing/SKILL.md
- /absolute/path/to/skills/typescript/SKILL.md
```

Delegating agents MUST resolve skills from the skill registry by canonical name and inject absolute filesystem paths to `SKILL.md` files under `## Skills to load before work`. Delegators MUST NOT inject inline summaries or hardcoded fallback paths.

The sub-agent MUST read those files before task-specific work. `SKILL.md` is the runtime contract and source of truth.

### Step 4: Report Resolution and Cache Invalidation

Sub-agents MUST report `skill_resolution` in their return envelope:

- `paths-injected` — received exact skill paths from the delegator and loaded them.
- `fallback-registry` — no paths received, self-loaded paths from the registry.
- `fallback-path` — loaded an explicit fallback path outside the registry.
- `none` — no skills loaded.

When any sub-agent returns a `skill_resolution` value other than `paths-injected` (`fallback-registry`, `fallback-path`, or `none`), the delegator MUST invalidate its cached skill registry and re-read the registry from persistent store before the next delegation.


## Compaction Safety

- The registry persists in Engram and `.atl/skill-registry.md`.
- Delegators can recover selected paths after compaction by re-reading the registry.
- Sub-agents receive exact files to read, so skill meaning is not degraded by generated summaries.

## Integration Points

- **ATL Orchestrator**: resolves paths for all SDD and non-SDD delegations.
- **judgment-day**: resolves paths before Judge A, Judge B, and Fix Agent.
- **pr-review and future delegators**: use this protocol when launching sub-agents.
