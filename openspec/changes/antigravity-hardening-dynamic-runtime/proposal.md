# Proposal: Dynamic Subagent Hardening & Strict TDD for Antigravity

## 1. Objectives & Scope
This proposal extends the Antigravity subagent tool hardening mechanism to support:
- Dynamic reviews (4R Lenses) and Judgment Day (Jueces/Arbitration) roles: `review-risk`, `review-readability`, `review-reliability`, `review-resilience`, `review-refuter`, `jd-judge-a`, `jd-judge-b`, and `jd-fix-agent`.
- Strict TDD (Test-Driven Development) constraints at the plugin level, restricting `sdd-apply` and `sdd-verify` behavior when `strict_tdd: true` is active.

## 2. Affected Packages
- `internal/components/sdd/`:
  - [antigravity_sdd_agents.go](file:///mnt/c/Users/Cobies/Desktop/Proyectos/COLAB/gentle-ai/internal/components/sdd/antigravity_sdd_agents.go) - Main plugin definition, role tables, and hardening prompt generation.
  - [antigravity_sdd_agents_test.go](file:///mnt/c/Users/Cobies/Desktop/Proyectos/COLAB/gentle-ai/internal/components/sdd/antigravity_sdd_agents_test.go) - Unit tests.

## 3. Rollback Plan
- Backup the original `antigravity_sdd_agents.go` and `antigravity_sdd_agents_test.go`.
- Restore the files and run `go test ./internal/components/sdd` to verify no regression.
- Since plugin installs are dynamic, running `gentle-ai install` (or the equivalent CLI commands) will overwrite/re-inject the plugin files on disk. An uninstall deletes the plugin directory entirely, removing the hook.
