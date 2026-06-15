# Apply Progress: Update Experience Overhaul (Slices 2, 3, 4 & 5)

**Change**: update-experience
**Mode**: Strict TDD

## Progress Overview

### Slice 2 (Update-check cooldown)
Complete. Introduces a 6h cooldown window for update checks, persisting the last successful check timestamp in the `.gentle-ai/state.json` file. Repeated launches check this cached timestamp to avoid rate-limiting or unnecessary remote requests to the GitHub API.

### Slice 3 (Channel-Honoring Upgrade)
Complete. Consolidates the stable and beta engram download paths into a single `DownloadLatestBinary` call that accepts the channel. Stable upgrades dynamically resolve the latest core engram version from the GitHub Releases API (excluding gentle-engram/pi tags) to share the same source of truth as the update-check path. Beta upgrades are installed from the main branch. To prevent a Go import cycle (since `internal/cli` imports `internal/components/engram`), the channel parameter in the `engram` package is typed as `string`, and the callers in `internal/update/upgrade` and `internal/cli` cast their typed `InstallChannel` to `string`.

### Slice 4 (Upgrade+Sync Deferred via `pending_sync`)
Complete. Implemented deferred synchronization using a `pending_sync` state flag. When `gentle-ai` undergoes a self-upgrade, the config sync phase is deferred to the next launch. The binary writes `PendingSync = true` to `state.json` and exits gracefully. Upon restart, `gentle-ai` detects the flag, runs config synchronization, and clears the flag on success. This prevents Windows binary lock issues and ensures a consistent restart path across all platforms.

### Slice 5 (CLI Prompt Default + Apply-Then-Close)
Complete. Implemented unconditional CLI prompt default and converged exit behavior. The `GENTLE_AI_CONFIRM_UPDATE` environment variable is removed and ignored, making the update prompt unconditional when updates are available. The prompt is default-Yes `[Y/n]` and accepts Enter (empty input) or explicit `y`/`yes` as acceptance. In non-interactive environments (non-TTY stdin), the prompt auto-declines to prevent hang/blockage. The `--yes` flag (via `GENTLE_AI_YES=1`) bypasses the interactive prompt and automatically accepts.

## TDD Cycle Evidence

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
| 3.4 | `internal/update/upgrade/strategy.go` | Unit | ✅ Passed | ✅ Written | ✅ Passed | ✅ 4 cases | ➖ None needed |
| 3.5 | `internal/update/upgrade/strategy.go` | Unit | ✅ Passed | ✅ Written | ✅ Passed | ✅ 4 cases | ➖ None needed |
| 3.6 | `internal/cli/channel_test.go` | Unit | ✅ Passed | ✅ Pre-existing | ✅ Passed | ✅ 2 cases | ➖ None needed |
| 4.1 | `internal/state/state_test.go` | Unit | ✅ Passed | ✅ Written | ✅ Passed | ✅ 4 cases | ➖ None needed |
| 4.2 | `internal/app/selfupdate_test.go` | Unit | ✅ Passed | ✅ Written | ✅ Passed | ✅ 3 cases | ➖ None needed |
| 4.3 | `internal/app/app_test.go` | Unit | ✅ Passed | ✅ Written | ✅ Passed | ✅ 4 cases | ➖ None needed |
| 4.4 | `internal/state/state_test.go` | Unit | ✅ Passed | ✅ Written | ✅ Passed | ✅ 4 cases | ➖ None needed |
| 4.5 | `internal/app/selfupdate_test.go` | Unit | ✅ Passed | ✅ Written | ✅ Passed | ✅ 3 cases | ➖ None needed |
| 4.6 | `internal/app/selfupdate_test.go` | Unit | ✅ Passed | ✅ Written | ✅ Passed | ✅ 3 cases | ➖ None needed |
| 4.7 | `internal/app/app_test.go` | Unit | ✅ Passed | ✅ Written | ✅ Passed | ✅ 4 cases | ➖ None needed |
| 4.8 | `internal/tui/model_test.go` | Unit | ✅ Passed | ✅ Written | ✅ Passed | ✅ 2 cases | ➖ None needed |
| 4.9 | `internal/app/app_test.go` | Unit | ✅ Passed | ✅ Written | ✅ Passed | ✅ 4 cases | ✅ Comment added |
| 5.1 | `internal/app/selfupdate_test.go` | Unit | ✅ Passed | ✅ Written | ✅ Passed | ✅ 3 cases | ➖ None needed |
| 5.2 | `internal/app/selfupdate_test.go` | Unit | ✅ Passed | ✅ Written | ✅ Passed | ✅ 1 case  | ➖ None needed |
| 5.3 | `internal/app/selfupdate_test.go` | Unit | ✅ Passed | ✅ Written | ✅ Passed | ✅ 2 cases | ➖ None needed |
| 5.4 | `internal/app/selfupdate_test.go` | Unit | ✅ Passed | ✅ Written | ✅ Passed | ✅ 3 cases | ➖ None needed |
| 5.5 | `internal/app/selfupdate_test.go` | Unit | ✅ Passed | ✅ Written | ✅ Passed | ✅ 2 cases | ➖ None needed |
| 5.6 | `internal/app/selfupdate_test.go` | Unit | ✅ Passed | ✅ Written | ✅ Passed | ✅ 3 cases | ➖ None needed |
| 5.7 | `internal/app/selfupdate_test.go` | Unit | ✅ Passed | ✅ Written | ✅ Passed | ✅ 1 case  | ✅ Comment added |

### Test Summary
- **Total tests written**: 39 (12 in Slice 2, 3 in Slice 3, 14 in Slice 4, 10 in Slice 5)
- **Total tests passing**: 39 (verified offline via Go standard runner)
- **Layers used**: Unit (39)
- **Approval tests** (refactoring): None — no legacy behavior refactoring
- **Pure functions created**: 1 (`checkSucceeded` in `internal/update/cooldown.go`)

## Files Changed

| File | Action | What Was Done |
|------|--------|---------------|
| `internal/state/state.go` | Modified | Added `LastUpdateCheck *time.Time` field to `InstallState` and updated `MergeAgents`. Added `PendingSync bool` field to `InstallState` and updated `MergeAgents` to carry it. |
| `internal/state/state_test.go` | Modified | Added tests for `LastUpdateCheck` and `PendingSync` round-tripping, omit-empty, backward-compatibility, and merge-preservation. |
| `internal/update/cooldown.go` | Created | Added `CheckAllWithCooldown` with clock injection (`nowFn`) and `checkAllFn` stubs to manage update cooldown. |
| `internal/update/cooldown_test.go` | Created | Added comprehensive suite of unit tests for `CheckAllWithCooldown` verifying TTL skips, stale refreshes, error retention, clock injection, and invalid directory fallbacks. |
| `internal/app/selfupdate.go` | Modified | Integrated `CheckAllWithCooldown` in the startup `selfUpdate` flow. Updated to set `PendingSync = true` in state on successful self-upgrade. Removed `GENTLE_AI_CONFIRM_UPDATE` check to prompt unconditionally. Implemented default-Yes prompt `[Y/n]` with Enter-to-accept, non-TTY auto-decline, and `GENTLE_AI_YES=1` bypass. Converged Unix/Windows restart flow. |
| `internal/app/selfupdate_test.go` | Modified | Added tests for `PendingSync` state writing, and comprehensive suite of 10 unit tests verifying unconditional prompt, default-Yes, TTY checks, GENTLE_AI_YES bypass, and env deletion. |
| `internal/app/app.go` | Modified | Integrated deferred sync checks on startup; checks `PendingSync` from state, runs sync automatically, and clears the flag on success. |
| `internal/app/app_test.go` | Modified | Added tests for `PendingSync` startup runner verifying successful sync/clear, failure persistence, writing warnings to stdout, and no-op when false. |
| `internal/tui/model.go` | Modified | Integrated `CheckAllWithCooldown` in Bubbletea TUI model `Init()`. Added logic to set `PendingSync = true` in state when Upgrade+Sync detects a self-upgrade event in TUI. |
| `internal/tui/model_test.go` | Modified | Added TUI model tests verifying `PendingSync` flag writing when upgrading in TUI. |
| `internal/versions/versions.go` | Modified | Removed unused `EngramCore` constant. |
| `internal/components/engram/download.go` | Modified | Updated `DownloadLatestBinary` signature to take `channel string` and restore dynamic version resolution for the stable channel. |
| `internal/components/engram/download_test.go` | Modified | Updated tests to pass `string` parameters and removed `TestDownloadLatestBinary_StableChannelUsesPinnedVersionDirectly` as it asserted incorrect hard-pinned behavior. |
| `internal/update/upgrade/strategy.go` | Modified | Updated signature of `engramDownloadFn` to take `cli.InstallChannel`, removed `engramBetaInstallFn`, and updated `engramBinaryUpgrade` to call `engramDownloadFn` directly. |
| `internal/update/upgrade/strategy_test.go` | Modified | Added `TestEngramBinaryUpgrade_ChannelHonoring` and updated mocks. |
| `internal/cli/run.go` | Modified | Updated `engramDownloadFn` call to pass `string(ChannelStable)`. |

## Deviations from Design
- To prevent a Go compilation import cycle (since `internal/cli` imports `internal/components/engram` to use `DownloadLatestBinary`, and `components/engram` would import `internal/cli` to use `cli.InstallChannel`), the channel parameter in `DownloadLatestBinary` is typed as a standard `string`. Callers cast their `cli.InstallChannel` values to `string` where needed.

## Issues Found
- **None**: All pre-existing Windows-specific path resolution, executable extensions, shell script mocks, and env isolation issues (such as `TestConfigPathsForBackup_GGAExtrasAreIncluded`, `TestEngramGoInstallFromMain_UsesGoEnvForBinDir`, and `TestEngramGoInstallFromMain_BypassesPublicGoProxy`) have been successfully resolved and are now passing cleanly on Windows.

