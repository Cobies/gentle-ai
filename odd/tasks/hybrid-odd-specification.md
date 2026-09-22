# Feature: Hybrid ODD Unified Specification Pattern

## 1. Diagnosis & Technical Proposal
- **Problem**:
  `odd/tasks/<feature>.md` was previously described as a loose checklist with basic metadata (objective, problem, why, scope, constraints). As a result, technical proposals, architectural trade-offs, and explicit contracts were either omitted or debated at length in the live chat session, causing rapid context window bloat and forcing workers to implement against vague prompts instead of concrete technical specifications.
- **Chosen Approach**:
  Formalize the "Hybrid ODD Unified Specification Pattern" across canonical guidance (`gentle-ai`) and the orchestrator harness (`gentle-pi`). The feature document `odd/tasks/<feature-name>.md` is standardized into four mandatory sections:
  1. Diagnosis & Technical Proposal (root cause, approach, trade-offs)
  2. Technical Specification & Contracts (interfaces, constraints, edge cases, authorized edit surfaces)
  3. Tasks & Evidence (actionable checklist with stable IDs, route declarations, and commit SHAs)
  4. Acceptance Criteria & Verification (deterministic verification commands, checks, and observed outcome)
  This offloads token weight from the live chat session into the physical file and provides workers with self-contained technical contracts.
- **Trade-offs**:
  Slightly more structured feature file upfront, but avoids OpenSpec multi-file ceremony while preventing context degradation and vague worker delegations.

## 2. Technical Specification & Contracts
- **Canonical Authority (`gentle-ai`)**:
  - `internal/components/agentguidance/routing.go`: Update `RenderRouting()` to specify the four structured sections in ODD step 5 and the feature document guidance bullet.
  - `internal/components/agentguidance/routing_test.go`: Assert the presence of the 4 sections, context offloading, and contract requirements.
  - `internal/components/sdd/odd_integration_test.go`: Ensure integration checks pass.
- **Harness Mirror (`gentle-pi`)**:
  - `assets/orchestrator-memory.md`: Detail the 4-part structure and token offloading benefit.
  - `assets/orchestrator-delegation.md`: Reference the 4 sections in ODD tracking.
  - `extensions/gentle-ai.ts`: Mirror the updated guidance in the orchestrator system prompt.
  - `docs/readme-reference.md`: Document the 4-part feature specification pattern.
  - `fixtures/odd-routing-canonical.md`: Update via `scripts/mirror-odd-routing.mjs`.
  - `tests/odd-routing-contract.test.ts` & `tests/odd-routing-canonical-ratchet.test.ts`: Verify ratchet alignment and zero drift.
- **Authorized Edit Surfaces**:
  - `odd/tasks/hybrid-odd-specification.md`
  - `internal/components/agentguidance/routing.go`
  - `internal/components/agentguidance/routing_test.go`
  - `internal/components/sdd/odd_integration_test.go`
  - `/mnt/c/Users/Cobies/Desktop/Proyectos/COLAB/gentle-pi/assets/orchestrator-memory.md`
  - `/mnt/c/Users/Cobies/Desktop/Proyectos/COLAB/gentle-pi/assets/orchestrator-delegation.md`
  - `/mnt/c/Users/Cobies/Desktop/Proyectos/COLAB/gentle-pi/extensions/gentle-ai.ts`
  - `/mnt/c/Users/Cobies/Desktop/Proyectos/COLAB/gentle-pi/docs/readme-reference.md`
  - `/mnt/c/Users/Cobies/Desktop/Proyectos/COLAB/gentle-pi/fixtures/odd-routing-canonical.md`
  - `/mnt/c/Users/Cobies/Desktop/Proyectos/COLAB/gentle-pi/tests/odd-routing-contract.test.ts`
  - `/mnt/c/Users/Cobies/Desktop/Proyectos/COLAB/gentle-pi/tests/odd-routing-canonical-ratchet.test.ts`

## 3. Tasks & Evidence
- [x] 1. Create feature branch `feat/hybrid-odd-specification` in `gentle-ai` and `gentle-pi`.
- [x] 2. `gentle-ai`: Update `routing_test.go` with failing tests for the 4-section hybrid pattern (TDD RED).
- [x] 3. `gentle-ai`: Update `routing.go` to satisfy tests and format the 4-section pattern (TDD GREEN).
- [x] 4. `gentle-ai`: Run tests, create work-unit commit on `feat/hybrid-odd-specification` (`72025280`).
- [x] 5. `gentle-pi`: Re-mirror canonical guidance to `fixtures/odd-routing-canonical.md` and update orchestrator guidance assets/extension.
- [x] 6. `gentle-pi`: Update contract tests, run `tests/odd-routing-contract.test.ts` & `tests/odd-routing-canonical-ratchet.test.ts`, and create work-unit commit on `feat/hybrid-odd-specification` (`1e647dce`).
- [x] 7. Push feature branches to GitHub for both repositories.

## 4. Acceptance Criteria & Verification
- `go test ./internal/components/agentguidance/...` passes in `gentle-ai` (all 16 supported agents pass).
- `node --test tests/odd-routing-contract.test.ts` (12/12 pass) and `node --test tests/odd-routing-canonical-ratchet.test.ts` (7/7 pass) pass in `gentle-pi`.
- Both repositories have work-unit commits pushed to GitHub (`origin/feat/hybrid-odd-specification`).
