# Proposal: Antigravity 2.0 and CLI Compatibility

## Intent
Add native support for the newly released **Google Antigravity CLI** (`antigravity-cli`) agent in `gentle-ai`. This allows CLI users to run SDD workflows using dynamic subagent delegation (via `define_subagent` and `invoke_subagent` tools) and Engram persistent memory.

## Scope

### In Scope
- Define `AgentAntigravityCLI` in `internal/model/types.go`.
- Implement `internal/agents/antigravitycli/adapter.go` for the CLI settings (`~/.gemini/antigravity-cli/`) and MCP servers.
- Map the new agent to a dedicated dynamic SDD orchestrator asset.
- Add `internal/assets/antigravitycli/sdd-orchestrator.md` containing dynamic subagent definition and invocation instructions.
- Ensure all test suites pass.

### Out of Scope
- Modifying legacy `AgentAntigravity` (Desktop IDE) behavior.
- Implementing automatic installation of the CLI executable (`agy`) itself.

## Capabilities

### New Capabilities
- `antigravity-cli-support`: Full support for Google Antigravity CLI with dynamic subagents and Engram.

### Modified Capabilities
- `sdd-orchestrator-assets`: Add `antigravitycli/sdd-orchestrator.md` orchestrator asset.

## Approach
Create a new agent adapter for the CLI (`antigravity-cli`). It will return `SupportsSubAgents() = false` in Go (as it doesn't use static files in a home folder), but its specific prompt `sdd-orchestrator.md` will instruct the LLM to dynamically call `define_subagent` and `invoke_subagent` in caliente by reading the skills installed under `~/.gemini/antigravity-cli/skills/`.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/model/types.go` | Modified | Add `AgentAntigravityCLI` constant. |
| `internal/agents/antigravitycli/adapter.go` | New | CLI agent adapter implementation. |
| `internal/agents/factory.go` | Modified | Register the new CLI agent. |
| `internal/components/sdd/inject.go` | Modified | Map new agent to its orchestrator. |
| `internal/assets/antigravitycli/sdd-orchestrator.md` | New | Dynamic subagent orchestrator prompt. |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Dynamic delegation failure | Low | Provide clear errors/fallbacks in orchestrator prompt. |
| Configuration path conflict | Low | Explicitly isolate paths using `antigravity-cli` subdirectory. |

## Rollback Plan
Revert changes using `git checkout` and delete the new `antigravitycli/` directories.

## Success Criteria
- [ ] `go test ./...` passes.
- [ ] Orchestrator template successfully generated.
