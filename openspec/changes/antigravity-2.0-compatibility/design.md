# Design: Antigravity 2.0 and CLI Compatibility

## Technical Approach
Add a dedicated `AgentAntigravityCLI` agent to `gentle-ai` with configuration paths mapping to `~/.gemini/antigravity-cli/`. Introduce a new asset `internal/assets/antigravitycli/sdd-orchestrator.md` specifying dynamic subagent execution using `define_subagent` and `invoke_subagent`.

## Architecture Decisions

| Option | Tradeoffs | Decision |
|--------|-----------|----------|
| **Dedicated Agent ID** | + Clean separation of paths and TUI installation screens.<br>- Adds a new type/constant. | **Go with dedicated ID (`antigravity-cli`)** |
| **Dynamic Subagents via Skills** | + No file duplication on disk. Runs natively on CLI session.<br>- Requires LLM to read skill file dynamically. | **Instruct LLM to read and define subagents dynamically** |

## Data Flow
```
[User Command] ──→ [gentle-ai Installer] ──→ Writes config to ~/.gemini/antigravity-cli/settings.json
                                         ──→ Merges MCP to ~/.gemini/antigravity/mcp_config.json
                                         ──→ Copies skills to ~/.gemini/antigravity-cli/skills/

[Antigravity CLI Session] ──→ Reads ~/.gemini/antigravity-cli/skills/{phase}/SKILL.md
                         ──→ Calls define_subagent(prompt=SKILL)
                         ──→ Calls invoke_subagent(prompt=task)
```

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `internal/model/types.go` | Modify | Define `AgentAntigravityCLI` constant. |
| `internal/catalog/agents.go` | Modify | Add `AgentAntigravityCLI` to agent catalog. |
| `internal/agents/factory.go` | Modify | Register new agent and switch cases. |
| `internal/agents/antigravitycli/adapter.go` | Create | Create adapter for `antigravity-cli`. |
| `internal/components/sdd/inject.go` | Modify | Map `AgentAntigravityCLI` to its orchestrator asset. |
| `internal/assets/antigravitycli/sdd-orchestrator.md` | Create | Dynamic delegation instructions for Antigravity CLI. |
| `internal/cli/install_test.go` | Modify | Update install CLI tests to cover new agent. |
| `internal/cli/run.go` / `validate.go` | Modify | Add compatibility logic. |
| `internal/tui/model.go` | Modify | Map selected agent strings in TUI. |

## Interfaces / Contracts

```go
package model

const (
	AgentAntigravityCLI AgentID = "antigravity-cli"
)
```

For `antigravitycli.Adapter`:
- `GlobalConfigDir` returns `~/.gemini/antigravity-cli`
- `MCPConfigPath` returns `~/.gemini/antigravity/mcp_config.json`
- `SupportsSubAgents()` returns `false` (dynamic definition only)
- `SupportsSkills()` returns `true`
- `MCPStrategy()` returns `model.StrategyMCPConfigFile`

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit | Adapter methods, factory registration, injection mapping. | Standard Go unit tests (`adapter_test.go`, `factory_test.go`). |
| Unit | Integration of MCP configurations and settings overlays. | Validate settings JSON structure merge results. |

## Migration / Rollout
No migration required.
