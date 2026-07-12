# Design: Dynamic Subagent Hardening & Strict TDD for Antigravity

## 1. Architectural Decisions & Rationale

### A. Dynamic Subagent Scoping
- **Decision**: Update `antigravitySddAgentsRoleScopes` in `internal/components/sdd/antigravity_sdd_agents.go` to explicitly include all new dynamic 4R and Judgment Day roles:
  - `review-risk`, `review-readability`, `review-reliability`, `review-resilience`, `review-refuter`, `jd-judge-a`, `jd-judge-b` (all read-only, no write tool privileges).
  - `jd-fix-agent` (write tools enabled, restricted to confirmed ledger edits only).
- **Rationale**: Keeps the role-to-scope table canonical, ensuring tests and static checkers can validate that subagents are restricted from escalating permissions.

### B. Prompt-Based Strict TDD Enforcement
- **Decision**: Update the constant `antigravitySddAgentsHardeningMessage` in `antigravity_sdd_agents.go` to include clear rules for Strict TDD:
  - When `strict_tdd: true` is active, `sdd-apply` is prohibited from editing production files without first writing or modifying test files and running the test runner to observe test failure (Red phase).
  - `sdd-verify` must run tests to verify behavior and is prohibited from editing source code.
  - Any attempt to bypass the TDD Red-Green-Refactor sequence must fail closed.
- **Rationale**: Since Antigravity lacks static overlays, embedding these guidelines in the injected PreInvocation contract prompt is the safest way to enforce TDD phases dynamically.

## 2. Component Details

### `internal/components/sdd/antigravity_sdd_agents.go`
- Update `antigravitySddAgentsRoleScopes` array to match the new roles exactly.
- Update `antigravitySddAgentsHardeningMessage` text.
- Add "Strict TDD" phrases to `antigravitySddAgentsHardeningContractPhrases` slice to ensure the tests check for the presence of the TDD rules in the contract.
