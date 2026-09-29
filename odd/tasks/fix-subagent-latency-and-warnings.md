# ODD Task Specification: Fix Subagent Latency, Host Warnings, and Startup Isolation

## 1. Diagnosis & Technical Proposal

### Problem
1. **Host-provided extension warnings**:
   Pi v0.99.1 emitted warnings on startup:
   - `gentle-pi/package.json`: `@earendil-works/pi-ai` and `@earendil-works/pi-tui` were in `dependencies` instead of `peerDependencies` with `*`.
   - `gentle-engram/package.json`: `typebox` was in `dependencies` instead of `peerDependencies` with `*`.
2. **Subagent startup failure (`Cannot find module startup-banner.ts`) and WSL2 I/O latency**:
   - Headless subagents running as child processes unnecessarily loaded interactive TUI modules (`startup-banner.ts`, `gentle-shell.ts`, `gentle-todo.ts`, `pi-pretty.ts`).
   - WSL2 DrvFs / 9P bridge locks during concurrent operations caused module resolution errors.
3. **Subagent reasoning latency**:
   - `~/.pi/agent/subagents.json` had all 24 subagent profiles pinned to `default_effort: "high"`, causing 800–2000 hidden reasoning tokens even on simple file read / grep actions.

### Implementation Summary
1. **Package Manifests**:
   - In `gentle-pi/package.json`: Moved `@earendil-works/pi-ai` and `@earendil-works/pi-tui` to `peerDependencies` with `"*"`, leaving only `@heyhuynhgiabuu/pi-pretty` in `dependencies`.
   - In `gentle-engram/package.json`: Moved `typebox` to `peerDependencies: { "typebox": "*" }`.
2. **Subagent Child Isolation**:
   - Added early child process guards (`if (process.env.GENTLE_PI_AGENTS_CHILD === "1") return;`) to `extensions/startup-banner.ts`, `extensions/gentle-shell.ts`, `extensions/gentle-todo.ts`, and `extensions/pi-pretty.ts`.
   - In subagent child processes, interactive UI elements, banner animation hooks, and TUI status handlers are skipped immediately.
3. **Subagent Thinking Effort Optimization**:
   - Updated `~/.pi/agent/subagents.json` to `default_effort: "medium"`.
   - Set `"effort": "low"` for fast read/scout/verify subagents (`gentle-ai-explore`, `gentle-ai-verify`, `sdd-explore`, `sdd-status`, `sdd-init`), and `"effort": "high"` for deep reasoning subagents (`sdd-design`, `jd-judge-a`, `jd-judge-b`).

---

## 2. Technical Specification & Contracts

### Allowed Edit Surfaces
- `/mnt/c/Users/Cobies/Desktop/Proyectos/COLAB/gentle-pi/package.json`
- `/home/cobies/.pi/agent/npm/node_modules/gentle-engram/package.json`
- `/mnt/c/Users/Cobies/Desktop/Proyectos/COLAB/gentle-pi/extensions/startup-banner.ts`
- `/mnt/c/Users/Cobies/Desktop/Proyectos/COLAB/gentle-pi/extensions/gentle-shell.ts`
- `/mnt/c/Users/Cobies/Desktop/Proyectos/COLAB/gentle-pi/extensions/gentle-todo.ts`
- `/mnt/c/Users/Cobies/Desktop/Proyectos/COLAB/gentle-pi/extensions/pi-pretty.ts`
- `/mnt/c/Users/Cobies/Desktop/Proyectos/COLAB/gentle-pi/tests/startup-banner.test.ts`
- `/home/cobies/.pi/agent/subagents.json`
- `odd/tasks/fix-subagent-latency-and-warnings.md`

---

## 3. Tasks & Evidence

- [x] Task 1: Fix package.json manifests in gentle-pi and gentle-engram to eliminate host-provided extension warnings.
  - Evidence: Verified clean diffs, valid JSON syntax, `@earendil-works/pi-ai`, `@earendil-works/pi-tui`, and `typebox` now in `peerDependencies: { ...: "*" }`.
- [x] Task 2: Implement subagent child process isolation (`GENTLE_PI_AGENTS_CHILD === "1"` guard in interactive extensions) and verify tests.
  - Evidence: `extensions/startup-banner.ts`, `gentle-shell.ts`, `gentle-todo.ts`, `pi-pretty.ts` updated. Unit test in `tests/startup-banner.test.ts` confirms 0 hooks/commands registered in child processes. 84/84 tests pass in `tests/startup-banner.test.ts` and `tests/agents-runner.test.ts`.
- [x] Task 3: Optimize thinking effort profiles in `~/.pi/agent/subagents.json` for lightweight scout/verification agents.
  - Evidence: `~/.pi/agent/subagents.json` updated with `default_effort: "medium"` and `"effort": "low"` on scout/verify profiles. Verified instant subagent execution on subsequent dispatches.
- [x] Task 4: Recompile gentle-ai v4 binary to `/home/cobies/go/bin/gentle-ai` and clean `/home/cobies/.pi/agent/mcp.json`.
  - Evidence: `rm -f /home/cobies/.pi/agent/mcp.json` executed successfully. `go install ./cmd/gentle-ai` completed; binary updated at `/home/cobies/go/bin/gentle-ai` (version `4.0.0-20260928233356-9ca4053449b3+dirty`).

---

## 4. Acceptance Criteria & Verification

1. [x] Pi starts without host-provided extension warnings.
2. [x] Subagents spawn cleanly without loading interactive TUI banners or crashing with `Cannot find module startup-banner.ts`.
3. [x] All 84 runner and startup-banner tests pass without regressions.
4. [x] Subagent execution latency is reduced significantly.
5. [x] gentle-ai v4 binary installed and /home/cobies/.pi/agent/mcp.json cleaned.
