# Workspace Rules & Execution Efficiency

## 1. Zero-Polling & Reactive Wakeup (CRITICAL: Token & Context Preservation)
- **NEVER poll or loop `manage_task(Action="status")` or `schedule` timers** waiting for a background task (`run_command`) or subagent (`invoke_subagent`) to complete.
- The Antigravity runtime uses **Reactive Wakeup**. When a command or subagent runs asynchronously in the background, you **MUST STOP calling tools immediately** and end your turn.
- The runtime will automatically wake you up with a high-priority message when the task/subagent finishes.
- Any manual busy-waiting or status-polling is an anti-pattern that bloats context history and wastes hundreds of thousands of tokens.

## 2. Exploration & Context Slicing Hierarchy
- **Mandatory CodeGraph First**: For all architecture, dependency mapping, symbol references, callers/callees, and structural questions, `codegraph_explore` (via MCP) MUST be used BEFORE falling back to filesystem tools (`grep_search`, `find_by_name`, `list_dir`). Broad grepping or full sweeps in the parent thread are prohibited.
- **Exploration Circuit Breaker**: The orchestrator in the parent chat is capped at a HARD LIMIT of at most 2 targeted file reads (`view_file`) or 1 search. If the answer is not found within 2 reads, the orchestrator is STRICTLY FORBIDDEN from continuing reading or grepping inline in the parent thread. It MUST either invoke `codegraph_explore` or delegate exploration to a dedicated subagent (`research` or `sdd-explore`).
- **Architectural & Flow Exploration Mandate (`sdd-explore`)**: Any investigation of system flows, cross-module integration, architectural questions, or codebase exploration touching 4+ files **MUST delegate to `sdd-explore`** (or `research` for purely conceptual lookups). Reading files inline beyond the 2-read circuit breaker or performing inline architectural investigations in the parent thread is STRICTLY PROHIBITED.
- **Symbol / Text search**: Use `grep_search` or `find_by_name`. **NEVER** cascade `list_dir` down multiple directory levels.
- **File Reading**: Always specify `StartLine` and `EndLine` in `view_file` to read targeted slices. Never dump whole files into context unless the file is under 100 lines.
- **Subagent Delegation**: Pass only distilled, minimal requirements to subagents. Subagents must return concise summaries and artifact references, never dumping raw logs or entire files back into the parent thread.

## 3. SDD Workflow & Implementation Routing
- **Pure Orchestrator Mandate (Zero Context Bloat & Token Preservation)**:
  - The parent orchestrator is strictly a strategist, thinker, and coordinator. It MUST NOT perform multi-file edits (`replace_file_content`), file creations (`write_to_file`), or broad test/build executions (`run_command`) inline in the parent thread.
  - All heavy lifting (multi-file reads, edits, bash commands) MUST be delegated to subagents so that detailed diffs and execution logs stay in isolated child contexts, keeping parent context under 20k tokens.
- **Unified SDD Routing**: Match the route to the task scope under a single SDD framework:
  - **Direct Inline**: STRICTLY for 1 single mechanical file with zero design or architectural ambiguity (e.g. typo fix, single constant update). Max 2 targeted reads inline.
  - **Lean SDD (Fast-Path)**: The DEFAULT automatic route for bounded changes touching 2–3 files within a single domain. Skips proposal, spec, and design ceremonies. Runs `sdd-explore` (targeted mapping with mandatory Engram persistence) → `sdd-apply` (strict TDD) → `sdd-verify`. Does NOT require explicit user prompt to activate.
  - **Full SDD**: For cross-layer / multi-domain tasks (e.g. backend + frontend, DB + API), core contracts/architecture changes, or substantial ambiguity. Runs the full phased lifecycle (`sdd-explore` → `sdd-propose` → `sdd-spec` → `sdd-design` → `sdd-tasks` → `sdd-apply` → `sdd-verify` → `sdd-archive`).
- **Dynamic Subagents as Universal Delegation Fallback (`define_subagent`)**:
  - While pre-registered static subagents (`sdd-*`, `research`) are the primary route, static subagents may not cover every ad-hoc or generic scenario (e.g. general multi-file refactoring, isolated script runs, custom test runners, or tasks requiring specific MCP tools).
  - If a static subagent is absent, unsuited, or requires direct MCP access (such as CodeGraph or Engram), the orchestrator MUST dynamically define a specialized subagent via `define_subagent` and invoke it via `invoke_subagent`.
  - Tool scoping rules for dynamic subagents:
    - `enable_subagent_tools: false` for all delegated workers (prevents uncontrolled recursion).
    - `enable_write_tools: false` for read-only exploration, mapper, review, and audit workers.
    - `enable_write_tools: true` for code modification, build, and test execution workers.
    - `enable_mcp_tools: true` whenever the worker needs Engram, CodeGraph, Obscura, or ClickUp access.
- **Mandatory Engram Exploration Persistence**: In ALL SDD routes (both Lean and Full), `sdd-explore` MUST be executed first to map dependencies, data flows, and symbol references. It MUST persist its findings, observations, and context into Engram (`mem_save` under topic key `sdd/{change-name}/explore`) before proceeding to code modification (`sdd-apply`). Never skip exploration or Engram indexing.
- **Multi-File & Non-Trivial Delegation Guard**: NEVER perform direct inline writes in the parent orchestrator thread for changes touching 2+ non-trivial files or exceeding 100 lines. Always delegate via SDD (`sdd-apply`) or a dynamically defined subagent worker.
- **Scope Detection (Lean vs Full)**: The orchestrator automatically detects whether a change is Lean or Full based on scope: changes confined to a single component/domain without contract changes qualify for Lean SDD; any cross-module, contract-breaking, or multi-domain change mandates Full SDD with user proposal confirmation.
- **SDD Task & Work-Unit Slicing (`sdd-apply`)**: In SDD execution, NEVER implement an entire multi-task feature in a single monolithic `sdd-apply` run. Decompose implementation into atomic work units (<400 changed lines per unit) and dispatch separate `sdd-apply` invocations per work unit or task batch, paired with focused tests.
- **Scope Creep Circuit Breaker**: If an inline edit begins on a single file but expands to dependent files (services, models, templates, tests), STOP inline modifications immediately, report the expansion to the user, and switch to Lean/Full SDD or dynamic subagent delegation.
- **Absolute Stop (No Auto-Apply)**: NEVER modify code or write files without explicit user confirmation first. Always present the proposed plan and pause for user confirmation before executing changes.
- **Interactive Gates**: When running multi-phase SDD, pause after each major phase (`sdd-explore`, `sdd-propose`, `sdd-spec`, `sdd-tasks`, `sdd-verify`), present a concise summary, and wait for confirmation.
- **Parallel Execution & Worktree Isolation**: When launching multiple concurrent subagents or parallel work units, always set `Workspace: "share"` in `invoke_subagent` to isolate execution in Git worktrees and prevent file lock contention.
- **TDD Requirement**: When writing new features or fixing bugs under SDD, ensure test coverage and verification in `sdd-verify`.



