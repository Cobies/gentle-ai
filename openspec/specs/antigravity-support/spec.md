# Antigravity support

Defines the required behavior for supporting Google Antigravity through the public `antigravity` agent ID.

## Requirements

### Requirement: Unified public agent ID

The system MUST expose Antigravity as `antigravity` and MUST NOT expose a separate public `antigravity-cli` agent option.

#### Scenario: Install uses the unified Antigravity agent

- GIVEN the installer is invoked for `antigravity`
- WHEN agent validation runs
- THEN `antigravity` is accepted
- AND the separate `antigravity-cli` option is not listed in the catalog or TUI.

### Requirement: Antigravity writes to the supported config surface

The system MUST write Antigravity settings, MCP config, plugins, and skills under `~/.gemini/antigravity-cli/`.

#### Scenario: Antigravity files are installed

- GIVEN the installer runs for `antigravity`
- WHEN SDD, Engram, or permission components are applied
- THEN settings are initialized at `~/.gemini/antigravity-cli/settings.json`
- AND MCP config is merged at `~/.gemini/antigravity-cli/mcp_config.json`
- AND skills are installed under `~/.gemini/antigravity-cli/skills/`.

### Requirement: Antigravity uses dynamic subagents

The Antigravity orchestrator MUST use runtime dynamic subagent tools rather than static subagent files. When running on a low-tier model, the system MUST enforce dynamic subagent delegation and MUST NOT execute SDD phases (such as explore, propose, spec, design, tasks, apply, verify) inline.

#### Scenario: SDD orchestration runs in Antigravity

- GIVEN the Antigravity SDD orchestrator is installed
- WHEN an SDD phase requires a subagent
- THEN the prompt instructs Antigravity to call `define_subagent`
- AND then call `invoke_subagent`.

#### Scenario: Low-model dynamic subagent enforcement

- GIVEN a low-tier model is active in the `antigravity` agent CLI
- WHEN the orchestrator compiles system instructions
- THEN the prompt MUST include explicit instructions warning the orchestrator to call `define_subagent` and `invoke_subagent` for each phase (`sdd-explore`, `sdd-propose`, `sdd-spec`, `sdd-design`, `sdd-tasks`, `sdd-apply`, `sdd-verify`)
- AND the prompt MUST explicitly forbid inline phase execution.

### Requirement: Antigravity shares the Gemini global prompt surface

The system MUST write global prompt/persona content for Antigravity to `~/.gemini/GEMINI.md`.

#### Scenario: Antigravity and Gemini CLI are selected together

- GIVEN both `gemini-cli` and `antigravity` are selected
- WHEN the installer applies SDD prompt content
- THEN the installer warns that both agents share `~/.gemini/GEMINI.md`.

### Requirement: Workspace rules discovery

The system SHALL check for the existence of workspace-specific workflow rules at `.agents/rules/sdd-workflow.md` when running the sync or execution command in the active workspace. If the file is not present, the system MUST proceed without throwing an error.

#### Scenario: Workspace rules file exists

- GIVEN a workspace with a rules file at `.agents/rules/sdd-workflow.md`
- WHEN the workspace rule discovery runs
- THEN the system loads the contents of `.agents/rules/sdd-workflow.md`

#### Scenario: Workspace rules file is missing

- GIVEN a workspace without a rules file at `.agents/rules/sdd-workflow.md`
- WHEN the workspace rule discovery runs
- THEN the system completes successfully without loading rules

### Requirement: Workspace rules injection

The system MUST inject the discovered workspace rules into `~/.gemini/GEMINI.md` within a dedicated block bounded by `<!-- gentle-ai:workspace-rules -->` and `<!-- /gentle-ai:workspace-rules -->` markers when the `antigravity` agent CLI is active.

#### Scenario: Workspace rules are injected with markers

- GIVEN loaded workspace rules and an active `antigravity` agent CLI
- WHEN the prompt sync process runs
- THEN the rules are written to `~/.gemini/GEMINI.md` inside `<!-- gentle-ai:workspace-rules -->` and `<!-- /gentle-ai:workspace-rules -->`
- AND existing content outside the markers is preserved.
