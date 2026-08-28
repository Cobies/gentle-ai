# Workspace Rules & Execution Efficiency

## 1. Zero-Polling & Reactive Wakeup (CRITICAL: Token & Context Preservation)
- **NEVER poll or loop `manage_task(Action="status")` or `schedule` timers** waiting for a background task (`run_command`) or subagent (`invoke_subagent`) to complete.
- The Antigravity runtime uses **Reactive Wakeup**. When a command or subagent runs asynchronously in the background, you **MUST STOP calling tools immediately** and end your turn.
- The runtime will automatically wake you up with a high-priority message when the task/subagent finishes.
- Any manual busy-waiting or status-polling is an anti-pattern that bloats context history and wastes hundreds of thousands of tokens.

## 2. Exploration & Context Slicing Hierarchy
- **Architecture / Structural queries**: Prefer `codegraph_explore` MCP tool before touching the filesystem.
- **Symbol / Text search**: Use `grep_search` or `find_by_name`. **NEVER** cascade `list_dir` down multiple directory levels.
- **File Reading**: Always specify `StartLine` and `EndLine` in `view_file` to read targeted slices. Never dump whole files into context unless the file is under 100 lines.
- **Subagent Delegation**: Pass only distilled, minimal requirements to subagents. Subagents must return concise summaries and artifact references, never dumping raw logs or entire files back into the parent thread.

## 3. SDD Workflow & Implementation Routing
- **Organic Routing**: Match the route to the task scope. Small/mechanical fixes (1 single file, already understood, no design questions) stay Direct Inline; delegated exploration/writing is used for broader tasks; SDD is used for large architectural features or when explicitly requested.
- **Multi-File & Non-Trivial Delegation Guard**: NEVER perform direct inline writes in the parent orchestrator thread for changes touching 2+ non-trivial files or exceeding 100 lines. Always delegate to a dedicated writer subagent or invoke SDD.
- **SDD Task & Work-Unit Slicing (`sdd-apply`)**: In SDD execution, NEVER implement an entire multi-task feature in a single monolithic `sdd-apply` run. Decompose implementation into atomic work units (<400 changed lines per unit) and dispatch separate `sdd-apply` invocations per work unit or task batch, paired with focused tests.
- **Scope Creep Circuit Breaker**: If an inline edit begins on a single file but expands to dependent files (services, models, templates, tests), STOP inline modifications immediately, report the expansion to the user, and switch to delegated subagents.
- **Absolute Stop (No Auto-Apply)**: NEVER modify code or write files without explicit user confirmation first. Always present the proposed plan and pause for user confirmation before executing changes.
- **Interactive Gates**: When running multi-phase SDD, pause after each major phase (`sdd-explore`, `sdd-propose`, `sdd-spec`, `sdd-tasks`, `sdd-verify`), present a concise summary, and wait for confirmation.
- **TDD Requirement**: When writing new features or fixing bugs under SDD, ensure test coverage and verification in `sdd-verify`.


