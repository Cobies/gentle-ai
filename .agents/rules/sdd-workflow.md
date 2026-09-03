# Workspace Rules & Execution Efficiency

## 1. Zero-Polling & Reactive Wakeup (CRITICAL: Token & Context Preservation)
- **NEVER poll or loop `manage_task(Action="status")` or `schedule` timers** waiting for a background task (`run_command`) or subagent (`invoke_subagent`) to complete.
- The Antigravity runtime uses **Reactive Wakeup**. When a command or subagent runs asynchronously in the background, you **MUST STOP calling tools immediately** and end your turn.
- The runtime will automatically wake you up with a high-priority message when the task/subagent finishes.
- Any manual busy-waiting or status-polling is an anti-pattern that bloats context history and wastes hundreds of thousands of tokens.

## 2. Exploration & Context Slicing Hierarchy
- **Mandatory CodeGraph First**: For all architecture, dependency mapping, symbol references, callers/callees, and structural questions, `codegraph_explore` (via MCP) MUST be used BEFORE falling back to filesystem tools (`grep_search`, `find_by_name`, `list_dir`). Broad grepping or full sweeps in the parent thread are prohibited.
- **Exploration Circuit Breaker**: The orchestrator in the parent chat is capped at a HARD LIMIT of at most 2 targeted file reads (`view_file`) or 1 search. If the answer is not found within 2 reads, the orchestrator is STRICTLY FORBIDDEN from continuing reading or grepping inline in the parent thread. It MUST either invoke `codegraph_explore` or delegate exploration to a dedicated subagent (`research` or `sdd-explore`).
- **Symbol / Text search**: Use `grep_search` or `find_by_name`. **NEVER** cascade `list_dir` down multiple directory levels.
- **File Reading**: Always specify `StartLine` and `EndLine` in `view_file` to read targeted slices. Never dump whole files into context unless the file is under 100 lines.
- **Subagent Delegation**: Pass only distilled, minimal requirements to subagents. Subagents must return concise summaries and artifact references, never dumping raw logs or entire files back into the parent thread.

## 3. SDD Workflow & Implementation Routing
- **Organic Routing**: Match the route to the task scope: Direct Inline for 1 single mechanical file with no design questions; Delegated Direct (dedicated direct-writer subagent) for 2–3 bounded files with zero architectural ambiguity; Optional SDD for multi-module features, shared core/contracts changes, or architectural ambiguity.
- **Multi-File & Non-Trivial Delegation Guard**: NEVER perform direct inline writes in the parent orchestrator thread for changes touching 2+ non-trivial files or exceeding 100 lines. Always delegate to the dedicated direct-writer subagent or invoke SDD.
- **SDD Task & Work-Unit Slicing (`sdd-apply`)**: In SDD execution, NEVER implement an entire multi-task feature in a single monolithic `sdd-apply` run. Decompose implementation into atomic work units (<400 changed lines per unit) and dispatch separate `sdd-apply` invocations per work unit or task batch, paired with focused tests.
- **Scope Creep Circuit Breaker**: If an inline edit begins on a single file but expands to dependent files (services, models, templates, tests), STOP inline modifications immediately, report the expansion to the user, and switch to delegated subagents.
- **Absolute Stop (No Auto-Apply)**: NEVER modify code or write files without explicit user confirmation first. Always present the proposed plan and pause for user confirmation before executing changes.
- **Interactive Gates**: When running multi-phase SDD, pause after each major phase (`sdd-explore`, `sdd-propose`, `sdd-spec`, `sdd-tasks`, `sdd-verify`), present a concise summary, and wait for confirmation.
- **Parallel Execution & Worktree Isolation**: When launching multiple concurrent subagents or parallel work units, always set `Workspace: "share"` in `invoke_subagent` to isolate execution in Git worktrees and prevent file lock contention.
- **TDD Requirement**: When writing new features or fixing bugs under SDD, ensure test coverage and verification in `sdd-verify`.


