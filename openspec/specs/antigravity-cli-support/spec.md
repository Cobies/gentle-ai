# Antigravity CLI Support Specification

## Purpose

Defines the required system behavior and file layout for supporting the **Google Antigravity CLI** (`antigravity-cli`) agent in `gentle-ai`, including its settings, MCP configuration, skill injection, and dynamic subagent orchestration.

---

## Requirements

### Requirement: Configuration Directories and Settings Path

The system MUST manage configurations and install the settings JSON file specifically inside the Antigravity CLI home directory.

- Home directory path: `~/.gemini/antigravity-cli/`
- Settings file path: `~/.gemini/antigravity-cli/settings.json`

#### Scenario: Settings are written to the correct CLI directory

- GIVEN the installer runs for the `antigravity-cli` agent
- WHEN the configuration is synchronized
- THEN the settings JSON is merged into `~/.gemini/antigravity-cli/settings.json`

---

### Requirement: MCP Server Configuration Path

The system MUST merge MCP server definitions (including Engram) into the Antigravity CLI MCP configuration file at `~/.gemini/antigravity-cli/mcp_config.json`.

#### Scenario: MCP servers are merged into the CLI mcp_config.json

- GIVEN the MCP injector runs for the `antigravity-cli` agent
- WHEN Engram is enabled
- THEN the Engram MCP configuration is merged into `~/.gemini/antigravity-cli/mcp_config.json`
- AND the Engram server args use plain `engram mcp` so Antigravity CLI receives Engram's standard MCP toolset instead of Pi-specific direct-tool profiles

---

### Requirement: Agent Skills Installation Path

The system MUST install global skills in `~/.gemini/antigravity-cli/skills/`.

#### Scenario: Global skills are copied to the CLI skills directory

- GIVEN the skill injector runs for the `antigravity-cli` agent
- WHEN a skill like `sdd-explore` is injected
- THEN it is written to `~/.gemini/antigravity-cli/skills/sdd-explore/SKILL.md`

---

### Requirement: Dynamic Subagents Support

The agent adapter for `antigravity-cli` MUST report that it does not support static subagent configuration files in Go (meaning `SupportsSubAgents()` returns `false`), so that no static agent files are copied to disk. However, the orchestrator prompt MUST guide the LLM to define and invoke subagents dynamically at runtime.

#### Scenario: Adapter reports false for static subagent files

- GIVEN the `antigravity-cli` agent adapter is queried
- WHEN `SupportsSubAgents()` is called
- THEN it returns `false`

---

### Requirement: Gemini-Compatible Prompt Surface

The system MUST treat `antigravity-cli` as a Gemini-compatible successor surface that intentionally writes global rules to `~/.gemini/GEMINI.md`. When `gemini-cli` and `antigravity-cli` are selected together, the installer MUST warn that the last synced SDD orchestrator owns the shared `gentle-ai:sdd-orchestrator` section.

#### Scenario: Antigravity CLI and Gemini CLI are selected together

- GIVEN both `gemini-cli` and `antigravity-cli` are selected
- WHEN verification checks run
- THEN the user receives a soft warning that Antigravity CLI shares the Gemini-compatible global prompt surface
- AND the warning recommends preferring Antigravity CLI for new installs
