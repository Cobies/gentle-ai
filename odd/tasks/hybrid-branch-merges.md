# ODD Task Specification: Merge Main into Hybrid Branches (gentle-ai & gentle-pi) - Cycle 2 (v4 & Petal Rail)

## 1. Diagnosis & Technical Proposal

### Problem
`upstream/main` advanced in both repositories:
- In `gentle-ai`: 22 new commits (`bb0b3f5e..9bce022a`), featuring the major `/v4` module path migration (`github.com/gentleman-programming/gentle-ai/v4`), Kimi Code paths, Antigravity Engram plugin updates, and Windows test portability.
- In `gentle-pi`: 19 new commits (`9879a1f9..5cc591f7`), featuring petal rail quiet tool calls, rounded card quiet tool frame, user message contrast styling, symlink materialization, and unref timeout fixtures.
Both repositories maintain their active `feat/hybrid-odd-specification` branches with custom architectural capabilities (Hybrid ODD unified spec, Pure Thinker guardrails, Dynamic Subagent Specialization, Cross-Repository Consent Mandate, Environment HUD Card).

### Technical Proposal & Integration Strategy
1. **In `gentle-pi`:**
   - Saved WIP for Environment HUD Card via `git stash`.
   - Updated local `main` to `upstream/main` (`5cc591f7`).
   - Merged `upstream/main` into `feat/hybrid-odd-specification`.
   - Restored WIP via `git stash pop` cleanly.
   - Verified unit tests for HUD and sidebar layout (51 passing) plus shell suites (164 passing).

2. **In `gentle-ai`:**
   - Updated local `main` to `upstream/main` (`9bce022a`).
   - Merged `upstream/main` into `feat/hybrid-odd-specification`.
   - Resolved merge conflict in `internal/components/agentguidance/routing_test.go` by migrating imports to `/v4` (`github.com/gentleman-programming/gentle-ai/v4/internal/opencode`) while retaining `init()` and the Hybrid ODD unified specification test case.
   - Verified `go test ./internal/components/agentguidance/...` and related packages.

---

## 2. Technical Specification & Contracts

### Allowed Edit Surfaces
- `odd/tasks/hybrid-branch-merges.md`
- `/mnt/c/Users/Cobies/Desktop/Proyectos/COLAB/gentle-ai/internal/components/agentguidance/routing_test.go`
- `/mnt/c/Users/Cobies/Desktop/Proyectos/COLAB/gentle-pi/extensions/gentle-shell.ts`
- `/mnt/c/Users/Cobies/Desktop/Proyectos/COLAB/gentle-pi/odd/tasks/environment-hud-card.md`

### Integration Commits
- `gentle-pi`: Commit `9d56a3e4f00fb749fc534b8ec7158e57cd7c1661` (`chore(merge): integrate upstream/main into feat/hybrid-odd-specification`)
- `gentle-ai`: Commit `9ca4053449b39c9d5d6edd9e7fbf4c8cd45a8f22` (`chore(merge): integrate upstream/main into feat/hybrid-odd-specification adapting /v4 module imports`)

---

## 3. Tasks & Evidence

- [x] Task 1: Clean working tree in gentle-pi, merge `upstream/main` (5cc591f7) into `feat/hybrid-odd-specification`, update local `main`, and run test suite.
  - Evidence: Commit `9d56a3e4`, stash popped with 0 conflicts, 51/51 tests pass in `tests/shell-hud.test.ts` & `tests/shell-sidebar-layout.test.ts`, 164/164 tests pass across shell test suites.
- [x] Task 2: Merge `upstream/main` (9bce022a) into `feat/hybrid-odd-specification` in gentle-ai, resolve `routing_test.go` conflict adapting the `/v4` import, update local `main`, and run `go test ./internal/components/agentguidance/...`.
  - Evidence: Commit `9ca40534`, `routing_test.go` migrated to `/v4`, `go test ./internal/components/agentguidance/...` PASS.
- [x] Task 3: Final cross-repository verification and status check.
  - Evidence: Verified clean branch state and test green in both gentle-ai and gentle-pi.

---

## 4. Acceptance Criteria & Verification

1. [x] Both repositories have their `feat/hybrid-odd-specification` branches updated with their respective `upstream/main` heads.
2. [x] All custom architectural features (Hybrid ODD, DSS, HUD card, guards) remain intact and functioning.
3. [x] Relevant test suites pass in both repositories.
4. [x] Local `main` in both repositories mirrors `upstream/main` 1:1.
