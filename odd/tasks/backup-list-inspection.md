# Tasks: Backup Listing and Disk Space Inspection (#2304)

Feature: `backup-list-inspection`
Issue: #2304 (`feat(cli): add backup listing and disk space inspection (gentle-ai backup list)`)
Branch: `feat/2304-backup-list`

## Tasks

- [x] 1. Implement backup listing and disk usage calculation in `internal/backup/list.go` with unit tests in `internal/backup/list_test.go`.
- [x] 2. Implement `gentle-ai backup list` and `gentle-ai backup clean` CLI command handlers with `--json` and `--keep` flags in `internal/cli/backup.go` and `internal/cli/backup_test.go`.
- [x] 3. Wire `backup` command dispatch in `internal/app/app.go` and help text in `internal/app/help.go`.
- [x] 4. Add backup footprint diagnostic check in `internal/cli/doctor.go` and `internal/doctor/doctor.go` with a 500 MB advisory threshold.
- [x] 5. Run full test suite (`go test ./internal/backup/...`, `go test ./internal/cli/...`, `go test ./internal/app/...`, `go test ./internal/doctor/...`) and verify line budget.
