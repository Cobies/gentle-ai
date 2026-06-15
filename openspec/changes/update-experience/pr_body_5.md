## 🔗 Linked Issue

Closes #872

---

## 🏷️ PR Type

- [ ] `type:bug` — Bug fix (non-breaking change that fixes an issue)
- [x] `type:feature` — New feature (non-breaking change that adds functionality)
- [ ] `type:docs` — Documentation only
- [ ] `type:refactor` — Code refactoring (no functional changes)
- [ ] `type:chore` — Build, CI, or tooling changes
- [ ] `type:breaking-change` — Breaking change (fix or feature that changes existing behavior)

---

## 📝 Summary

This PR completes Slice 5 (CLI Prompt Default + Apply-Then-Close) of the Update Experience Overhaul:
1. Prompts for update unconditionally on TTY/interactive environments.
2. Prompts default to Yes `[Y/n]` (Enter accepts).
3. Automatically declines on non-TTY (CI/scripts) to prevent hanging.
4. Bypasses prompt via `GENTLE_AI_YES=1` (or `--yes` flag).
5. Resolves Windows compatibility bugs across the test harness (APPDATA isolation, path separator normalization, WSL bash checking, and command mocking).

---

## ⚖️ Size Exception Rationale

This PR exceeds the 400 changed lines limit because:
1. **Stacked PR Chain:** It is stacked on top of PR #894 and PR #895. The actual changes introduced solely by Slice 5 are ~160 lines, but GitHub compares against `main` and counts the cumulative changes of all three slices.
2. **Windows Compatibility Fixes:** To ensure a green CI on all systems, we resolved 7 pre-existing Windows-specific test harness bugs (isolated APPDATA/USERPROFILE caches, normalized paths, mocked executables/cmds) which added multiple platform-specific test configurations.

---


## 📂 Changes

| File / Area | What Changed |
|-------------|-------------|
| `internal/cli/run_engram_download_test.go` | Appended `.exe` to binary names on Windows. |
| `internal/components/engram/download_test.go` | Fixed path separator comparison and mocked `go` using `go.cmd` batch script on Windows. |
| `internal/components/gga/config_test.go` | Isolated `APPDATA` directory inside temporary directories for tests. |
| `internal/tui/model_test.go` | Mocked `USERPROFILE` environment variable alongside `HOME` for proper cache isolation. |
| `internal/update/check_test.go` | Used `mockCmd` instead of direct `exec.Command` for `echo` commands. |
| `internal/update/install_script_test.go` | Skipped test when bash/WSL is broken or missing. |
| `internal/update/upgrade/executor_test.go` | Used dynamic `gga` path resolvers and isolated `APPDATA` in tests. |
| `openspec/changes/update-experience/apply-progress.md` | Documented Slice 5 progress. |
| `openspec/changes/update-experience/tasks.md` | Marked tasks 5.1 to 5.7 as complete. |
| `openspec/changes/update-experience/verify-report.md` | Updated verify report for Slices 2 to 5. |

---

## 🧪 Test Plan

**Unit Tests**
```bash
go test ./...
```

- [x] Unit tests pass (`go test ./...`)
- [ ] E2E tests pass (`cd e2e && ./docker-test.sh`)
- [x] Manually tested locally

---

## 🤖 Automated Checks

The following checks run automatically on this PR:

| Check | Status | Description |
|-------|--------|-------------|
| Check PR Cognitive Load | ⏳ | PR should stay within 400 changed lines (`additions + deletions`) or use `size:exception` |
| Check Issue Reference | ⏳ | PR body must contain `Closes/Fixes/Resolves #N` |
| Check Issue Has `status:approved` | ⏳ | Linked issue must have been approved before work began |
| Check PR Has `type:*` Label | ⏳ | Exactly one `type:*` label must be applied |
| Unit Tests | ⏳ | `go test ./...` must pass |
| E2E Tests | ⏳ | `cd e2e && ./docker-test.sh` must pass |

---

## ✅ Contributor Checklist

- [x] PR is linked to an issue with `status:approved`
- [x] PR stays within 400 changed lines, or I have requested/obtained maintainer-applied `size:exception` with rationale documented
- [x] I have added the appropriate `type:*` label to this PR
- [x] Unit tests pass (`go test ./...`)
- [ ] E2E tests pass (`cd e2e && ./docker-test.sh`)
- [x] I have updated documentation if necessary
- [x] My commits follow [Conventional Commits](https://www.conventionalcommits.org/) format
- [x] My commits do not include `Co-Authored-By` trailers
