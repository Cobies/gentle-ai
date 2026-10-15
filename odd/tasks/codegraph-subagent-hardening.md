# Technical Specification: CodeGraph Subagent Hardening & Guidance in `gentle-pi`

## Part 1: Diagnosis & Technical Proposal
* **Problem**: 
  1. Subagents (`gentle-ai-explore`) frequently fall back to `read`/`grep` because querying/exploring a repository without a pre-existing `.codegraph/` index triggers an exit code 1 with "not initialized" from CodeGraph CLI. In `extensions/codegraph-tools.ts`, this is caught as a generic `failed` error returning: *"CodeGraph failed to run. Use read, grep, and find for this exploration."*, which causes agents to immediately abandon CodeGraph.
  2. Subagent workers (`gentle-ai-worker.md`) lack the `codegraph` tool in their YAML frontmatter, preventing them from running symbol exploration before bounded edits.
* **Proposed Solution**:
  1. In `extensions/codegraph-tools.ts`: Check if `.codegraph/` directory exists or if error output indicates "not initialized". For `query` and `explore`, return a guided error message informing the agent to run `operation: "init"` first rather than triggering immediate read/grep fallback.
  2. In `assets/agents/gentle-ai-worker.md`: Add `codegraph` to `tools:`.
  3. In `assets/agents/gentle-ai-explore.md`: Explicitly instruct that a missing index is the trigger to run `operation: "init"` rather than an unavailable/failed state.
  4. In `tests/codegraph-tools.test.ts`: Add tests covering uninitialized index handling and verify full test suite passes.

## Part 2: Technical Specification & Contracts
* **Uninitialized status in `codegraph-tools.ts`**:
  - Add `UNINITIALIZED: "uninitialized"` to `CODEGRAPH_STATUS`.
  - When `operation` is `query` or `explore` and `join(cwd, ".codegraph")` does not exist (or execution fails with "not initialized" in stdout/stderr):
    - Return message: `"CodeGraph is not initialized in this workspace. Run codegraph with operation: \"init\" first before querying or exploring."`
    - Do NOT include fallback instructions (`FALLBACK_INSTRUCTIONS`) for uninitialized state so the agent knows to initialize.
* **Agent definitions**:
  - `assets/agents/gentle-ai-worker.md`: Add `codegraph` to `tools:`.
  - `assets/agents/gentle-ai-explore.md`: Update guidance.
* **Allowed edit surfaces in `gentle-pi`**:
  - `extensions/codegraph-tools.ts`
  - `assets/agents/gentle-ai-worker.md`
  - `assets/agents/gentle-ai-explore.md`
  - `tests/codegraph-tools.test.ts`

## Part 3: Tasks & Evidence
- [ ] Task 1: Update `extensions/codegraph-tools.ts` to detect uninitialized workspace on query/explore and return guided initialization instruction.
- [ ] Task 2: Add unit tests in `tests/codegraph-tools.test.ts` asserting uninitialized behavior.
- [ ] Task 3: Update `assets/agents/gentle-ai-worker.md` and `assets/agents/gentle-ai-explore.md`.
- [ ] Task 4: Run test suite in `gentle-pi` and verify all tests pass.

## Part 4: Acceptance Criteria & Verification
- Unit test `CodeGraph tool instructs initialization when queried before init` passes.
- All existing tests in `tests/codegraph-tools.test.ts` pass.
- `gentle-ai-worker.md` includes `codegraph` in tools.
