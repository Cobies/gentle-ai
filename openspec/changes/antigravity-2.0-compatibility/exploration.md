## Exploration: Antigravity 2.0 and CLI Compatibility

### Current State
Currently, `gentle-ai` only supports the desktop-only "Google Antigravity" IDE (`model.AgentAntigravity`). It configures paths under `~/.gemini/antigravity/` and treats it as a non-installable IDE. It defaults to the solo-agent (inline) SDD orchestrator because it does not support static subagent configuration files.

With Google's release of **Antigravity 2.0** and the **Antigravity CLI**, the platform now supports:
1. A terminal-first CLI (`agy`) with settings in `~/.gemini/antigravity-cli/settings.json` and skills in `~/.gemini/antigravity-cli/skills/`.
2. Dynamic subagent orchestration via native `/agents` command and LLM tools (`define_subagent` and `invoke_subagent`).
3. Persisted memory across sessions using the **Engram MCP server** configured in `~/.gemini/antigravity/mcp_config.json`.

### Affected Areas
- `internal/model/types.go` — Add `AgentAntigravityCLI` (`antigravity-cli`) agent ID.
- `internal/agents/antigravitycli/adapter.go` — [NEW] Create the adapter for the new CLI agent.
- `internal/agents/factory.go` — Register the new `AgentAntigravityCLI` agent.
- `internal/components/sdd/inject.go` — Map `AgentAntigravityCLI` to `antigravitycli/sdd-orchestrator.md`.
- `internal/assets/antigravitycli/sdd-orchestrator.md` — [NEW] Create the orchestrator prompt optimized for dynamic subagent delegation.
- `internal/components/mcp/inject.go` / `context7.go` — Ensure MCP config matches the new paths.
- `internal/components/permissions/inject.go` — Handle permission overlays for the CLI.

### Approaches
1. **Option 1: Add a dedicated Agent ID `AgentAntigravityCLI` ("antigravity-cli")**
   - Pros:
     * Completely isolates CLI settings (`~/.gemini/antigravity-cli/`) from the Desktop IDE (`~/.gemini/antigravity/`).
     * Allows custom auto-install command support for CLI (installable via shell/package manager) while keeping Desktop IDE non-installable.
     * Can customize TUI installer screen and validation specifically for the CLI.
   - Cons:
     * Adds one new Agent ID in `types.go`.
   - Effort: Medium

2. **Option 2: Overload the existing `AgentAntigravity` ("antigravity")**
   - Pros:
     * No new Agent ID needed.
   - Cons:
     * Hard to auto-detect both IDE and CLI accurately in the same adapter.
     * Confuses configuration paths and installation capabilities, as they differ significantly.
   - Effort: High (complexity-wise due to branching paths)

### Recommendation
We recommend **Option 1**. Isolating the CLI agent ensures that path detection, settings merging, and subagent orchestration instructions remain clean, explicit, and easy to maintain without polluting the legacy Desktop IDE adapter.

### Risks
- **Nesting limit:** Antigravity CLI enforces a strict subagent nesting depth limit of 10 levels. The SDD orchestrator must be instructed to limit recursion.
- **Dynamic definition overhead:** The orchestrator must read `SKILL.md` files dynamically from the home directory to call `define_subagent`, which increases token usage per phase transition.

### Ready for Proposal
Yes. We have a clear architectural direction. The orchestrator should proceed to generate the change proposal for user approval.
