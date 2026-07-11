# Proposal: Workspace Rule Configuration for agy CLI

## Intent

To establish a mechanism in the `agy` (Google Antigravity) CLI for loading and enforcing workspace-specific rules configured in `.agents/rules/sdd-workflow.md`. This ensures that when executing tasks using a low model (e.g., Gemini 3.5 Flash Low), the agent strongly adheres to workflow boundaries—specifically, that the orchestrator must define and invoke dynamic subagents for each phase rather than executing tasks inline—without impacting other editors or agents in the workspace.

## Scope

### In Scope
- Load workspace-specific rules from `.agents/rules/sdd-workflow.md` exclusively during `antigravity` agent prompt synchronization and execution.
- Append these rules into `~/.gemini/GEMINI.md` within a dedicated marker block (`<!-- gentle-ai:workspace-rules -->`) to isolate them to `antigravity`.
- Enforce explicit rules instructing the orchestrator that it **MUST** use `define_subagent` and `invoke_subagent` to run dynamic subagents (e.g., `sdd-explore`, `sdd-propose`, `sdd-spec`, `sdd-design`, `sdd-tasks`, `sdd-apply`, `sdd-verify`) for each individual phase, forbidding inline execution when a low model is active.
- Ensure the rule loading mechanism is completely isolated to the `antigravity` adapter and does not affect Cursor, Windsurf, Claude Code, or Pi.

### Out of Scope
- Modifying prompt injection mechanisms or rule formats for other non-Antigravity agents.
- General modifications to the prompt structures of other editors.

## Capabilities

### New Capabilities
None

### Modified Capabilities
- `antigravity-support`: Extend this capability to detect `.agents/rules/sdd-workflow.md` in the active workspace and inject it into the compiled system prompt surface (`~/.gemini/GEMINI.md`), wrapping it in dynamic subagent enforcement instructions.

## Approach

1. **Rule Discovery & Loading**:
   - In the `antigravity` adapter (or during `ComponentPersona` sync for `antigravity`), check for the existence of `{workspaceDir}/.agents/rules/sdd-workflow.md`.
   - Read the contents of the file if present.

2. **System Prompt Injection**:
   - Update `internal/components/persona/inject.go` or the `antigravity` prompt building logic to read the workspace rules.
   - Inject the rules into `~/.gemini/GEMINI.md` inside a marker:
     ```markdown
     <!-- gentle-ai:workspace-rules -->
     ...
     <!-- /gentle-ai:workspace-rules -->
     ```
   - Prepend model-hardening instructions to these workspace rules when the active model is classified as a "low model" (e.g., Gemini 3.5 Flash Low), instructing the model that it **MUST NOT** execute phases inline and **MUST** delegate using dynamic subagents.

3. **Low-Model Guard Rails**:
   - Define a detection helper for low-tier models (checking configuration/environment variables or active model IDs) to inject additional enforcement warnings directly into the prompt stream.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/agents/antigravity/adapter.go` | Modified | Add support/helpers for detecting workspace rules. |
| `internal/components/persona/inject.go` | Modified | Integrate workspace rules loading and target prompt markdown injection. |
| `openspec/specs/antigravity-support/spec.md` | Modified | Add specs for workspace rule injection and subagent enforcement. |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Rule duplication or bloating of `GEMINI.md` | Low | Use strict marker-based block replacement (`<!-- gentle-ai:workspace-rules -->`) to ensure idempotency. |
| Low model ignores injected rules due to context window noise | Medium | Inject the dynamic subagent enforcement rules at the very end of the prompt stream (or within high-priority zones) when a low model is active. |

## Rollback Plan

- Delete/revert modifications in `internal/agents/antigravity/adapter.go` and `internal/components/persona/inject.go`.
- Run `gentle-ai sync` with the restored CLI to clean up and remove the injected rule markers from `~/.gemini/GEMINI.md`.

## Dependencies

- None

## Success Criteria

- [ ] `gentle-ai sync` detects `.agents/rules/sdd-workflow.md` and correctly injects it into `~/.gemini/GEMINI.md`.
- [ ] Injected content contains explicit instructions requiring the orchestrator to define and invoke dynamic subagents for all SDD phases.
- [ ] No other agent prompt surfaces (e.g., `.cursor/rules/`, `.codeium/windsurf/memories/`, `.claude/`) are modified or affected.
- [ ] All unit and integration tests run and pass successfully (`go test ./...`).
