# Apply Progress: Update Experience Overhaul (Slices 2 & 3)

**Change**: update-experience
**Mode**: Strict TDD

## Progress Overview

### Slice 2 (Update-check cooldown)
Complete. Introduces a 6h cooldown window for update checks, persisting the last successful check timestamp in the `.gentle-ai/state.json` file. Repeated launches check this cached timestamp to avoid rate-limiting or unnecessary remote requests to the GitHub API.

### Slice 3 (Channel-Honoring Upgrade)
Complete. Consolidates the stable and beta engram download paths into a single `DownloadLatestBinary` call that accepts the channel. Stable upgrades now bypass the GitHub Releases API and use the pinned `versions.EngramCore` version ("1.3.0") directly. To prevent a Go import cycle (since `internal/cli` imports `internal/components/engram`), the channel parameter in the `engram` package is typed as `string`, and the callers in `internal/update/upgrade` and `internal/cli` cast their typed `InstallChannel` to `string`.

### TDD Cycle Evidence

| Task | Test File | Layer | Safety Net | RED | GREEN | TRIANGULATE | REFACTOR |
|------|-----------|-------|------------|-----|-------|-------------|----------|
| 2.1 | `internal/state/state_test.go` | Unit | N/A (new field) | ✅ Written | ✅ Passed | ✅ 4 cases | ➖ None needed |
| 2.2 | `internal/update/cooldown_test.go` | Unit | N/A (new file) | ✅ Written | ✅ Passed | ✅ 4 cases | ➖ None needed |
| 2.3 | `internal/state/state_test.go` | Unit | N/A (new field) | ✅ Written | ✅ Passed | ✅ 4 cases | ➖ None needed |
| 2.4 | `internal/update/cooldown_test.go` | Unit | N/A (new file) | ✅ Written | ✅ Passed | ✅ 2 cases | ➖ None needed |
| 2.5 | `internal/update/cooldown_test.go` | Unit | N/A (new file) | ✅ Written | ✅ Passed | ✅ 2 cases | ➖ None needed |
| 2.6 | `internal/update/cooldown_test.go` | Unit | N/A (new file) | ✅ Written | ✅ Passed | ✅ 2 cases | ➖ None needed |
| 2.7 | `internal/update/cooldown_test.go` | Unit | N/A (new file) | ✅ Written | ✅ Passed | ✅ 8 cases | ✅ Clock injected |
| 3.1 | `internal/update/upgrade/strategy_test.go` | Unit | ✅ Passed | ✅ Written | ✅ Passed | ✅ 4 cases | ➖ None needed |
| 3.2 | `internal/components/engram/download_test.go` | Unit | ✅ Passed | ✅ Written | ✅ Passed | ✅ 2 cases | ➖ None needed |
| 3.3 | `internal/components/engram/download_test.go` | Unit | ✅ Passed | ✅ Written | ✅ Passed | ✅ 2 cases | ➖ None needed |
| 3.4 | `internal/update/upgrade/strategy_test.go` | Unit | ✅ Passed | ✅ Written | ✅ Passed | ✅ 4 cases | ➖ None needed |
| 3.5 | `internal/update/upgrade/strategy_test.go` | Unit | ✅ Passed | ✅ Written | ✅ Passed | ✅ 4 cases | ➖ None needed |
| 3.6 | `internal/cli/channel_test.go` | Unit | ✅ Passed | ✅ Pre-existing | ✅ Passed | ✅ 2 cases | ➖ None needed |

### Test Summary
- **Total tests written**: 15 (12 in Slice 2 + 3 in Slice 3)
- **Total tests passing**: 15 (verified offline via Go standard runner)
- **Layers used**: Unit (15)
- **Approval tests** (refactoring): None — no legacy behavior refactoring
- **Pure functions created**: 1 (`checkSucceeded` in `internal/update/cooldown.go`)

## Files Changed

| File | Action | What Was Done |
|------|--------|---------------|
| `internal/state/state.go` | Modified | Added `LastUpdateCheck *time.Time` field to `InstallState` and updated `MergeAgents`. |
| `internal/state/state_test.go` | Modified | Added round-trip, omit-empty, backward-compatibility, and merge-preservation tests for `LastUpdateCheck`. |
| `internal/update/cooldown.go` | Created | Added `CheckAllWithCooldown` with clock injection (`nowFn`) and `checkAllFn` stubs to manage update cooldown. |
| `internal/update/cooldown_test.go` | Created | Added comprehensive suite of unit tests for `CheckAllWithCooldown` verifying TTL skips, stale refreshes, error retention, clock injection, and invalid directory fallbacks. |
| `internal/app/selfupdate.go` | Modified | Integrated `CheckAllWithCooldown` in the startup `selfUpdate` flow with `UpdateCheckTTL` (6 hours). |
| `internal/tui/model.go` | Modified | Integrated `CheckAllWithCooldown` in Bubbletea TUI model `Init()` utilizing `tuiNowFn` for cooldown check. |
| `internal/versions/versions.go` | Modified | Added `EngramCore = "1.3.0"` constant. |
| `internal/components/engram/download.go` | Modified | Updated `DownloadLatestBinary` signature to take `channel string` and use `versions.EngramCore` for stable channel downloads. |
| `internal/components/engram/download_test.go` | Modified | Updated tests to pass `string` parameters and added `TestDownloadLatestBinary_StableChannelUsesPinnedVersionDirectly` to verify bypassing API queries. |
| `internal/update/upgrade/strategy.go` | Modified | Updated `engramDownloadFn` signature to take `cli.InstallChannel`, removed `engramBetaInstallFn`, and updated `engramBinaryUpgrade` to call `engramDownloadFn` directly casting channel to string. |
| `internal/update/upgrade/strategy_test.go` | Modified | Added `TestEngramBinaryUpgrade_ChannelHonoring` and updated all mock overrides of `engramDownloadFn`. |
| `internal/cli/run.go` | Modified | Updated `engramDownloadFn` call to pass `string(ChannelStable)`. |

## Deviations from Design
- To prevent a Go compilation import cycle (since `internal/cli` imports `internal/components/engram` to use `DownloadLatestBinary`, and `components/engram` would import `internal/cli` to use `cli.InstallChannel`), the channel parameter in `DownloadLatestBinary` is typed as a standard `string`. Callers cast their `cli.InstallChannel` values to `string` where needed.

## Issues Found
- The test `TestConfigPathsForBackup_GGAExtrasAreIncluded` in `internal/update/upgrade/executor_test.go` has a pre-existing path resolution failure on Windows environments, which is unrelated to the changes introduced in this slice.
- Windows-specific path slash issues exist in `TestEngramGoInstallFromMain_UsesGoEnvForBinDir` and `TestEngramGoInstallFromMain_BypassesPublicGoProxy` (pre-existing test failures on Windows).
