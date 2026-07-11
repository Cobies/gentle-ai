# Design: Configure Antigravity Workspace Rules

## Technical Approach

We will extend the `antigravity` adapter and `ComponentPersona` sync pipeline to load workspace-specific workflow rules from `.agents/rules/sdd-workflow.md` in the current working directory. The loaded rules, along with dynamic subagent enforcement warnings if a low-tier model is detected, will be injected into `~/.gemini/GEMINI.md` within a dedicated `<!-- gentle-ai:workspace-rules -->` marker block using the existing atomic write and markdown section injection infrastructure.

## Architecture Decisions

| Option | Tradeoff | Decision |
|--------|----------|----------|
| Load workspace rules during `sync` / `install` | Rules are persisted statically in `GEMINI.md` at sync time. Requires a sync run to update. | **Chosen**. Fits the existing prompt compilation pipeline in `inject.go` and keeps runtime overhead minimal. |
| Parse rules dynamically on every CLI execution | Adds runtime I/O overhead to every command execution. | Rejected. |

| Option | Tradeoff | Decision |
|--------|----------|----------|
| Classify active model via environment / config | Simple environment/config checking (`GEMINI_MODEL`, `ANTIGRAVITY_MODEL` envs or `settings.json` model ID key). | **Chosen**. Uses `model.ModelCapability(modelID)` to check for `flash`, `mini`, `haiku`, etc. |
| Hardcode model list in adapter | Less flexible, requires code change for new models. | Rejected. |

## Data Flow

```text
Active Workspace (.agents/rules/sdd-workflow.md) ──┐
                                                    ├──► Load & Combine ──► Inject into ~/.gemini/GEMINI.md
Active Model ID (from settings.json / environment) ─┘
```

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `internal/agents/antigravity/adapter.go` | Modify | Add helper methods to discover `.agents/rules/sdd-workflow.md` in the current workspace, read settings.json for model detection, and check for low-tier models. |
| `internal/components/persona/inject.go` | Modify | Update persona injection for `StrategyAppendToFile` (used by Antigravity) to load discovered workspace rules, compile them with model warning instructions if applicable, and inject them into `GEMINI.md` inside `<!-- gentle-ai:workspace-rules -->` markers. |

## Interfaces / Contracts

### Low-Model Detection in Adapter
We will add helper methods to `internal/agents/antigravity/adapter.go`:

```go
func (a *Adapter) DetectLowModel(homeDir string) bool {
	// 1. Check env vars
	for _, env := range []string{"GEMINI_MODEL", "ANTIGRAVITY_MODEL"} {
		if val := os.Getenv(env); val != "" {
			return model.ModelCapability(val) == "small"
		}
	}
	// 2. Fall back to settings.json
	settingsPath := a.SettingsPath(homeDir)
	if data, err := os.ReadFile(settingsPath); err == nil {
		var cfg struct {
			Model   string `json:"model"`
			ModelID string `json:"modelId"`
		}
		if err := json.Unmarshal(data, &cfg); err == nil {
			m := cfg.Model
			if m == "" {
				m = cfg.ModelID
			}
			if m != "" {
				return model.ModelCapability(m) == "small"
			}
		}
	}
	return false
}
```

### Low-Model Warning Prompt
When `DetectLowModel` is true, the following warning is prepended to the workspace rules block:

```markdown
> [!IMPORTANT]
> ACTIVE MODEL CLASSIFIED AS LOW-TIER.
> You MUST NOT execute SDD phases (sdd-explore, sdd-propose, sdd-spec, sdd-design, sdd-tasks, sdd-apply, sdd-verify) inline within the parent thread.
> You MUST define and invoke dynamic subagents using define_subagent and invoke_subagent for each phase. Inline execution is strictly forbidden.
```

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit | Rule Discovery & Parsing | Verify workspace rules are loaded if the file exists, and handled gracefully if missing. |
| Unit | Low Model Detection | Test `DetectLowModel` with environment variable overrides and mocked settings.json contents. |
| Integration | End-to-end Sync | Perform sync for `antigravity` agent and verify `~/.gemini/GEMINI.md` contains the `workspace-rules` section with markers. |

## Threat Matrix

N/A — no routing, shell, subprocess, VCS/PR automation, executable-file classification, or process-integration boundary.

## Migration / Rollout

No migration required.

## Open Questions

None.
