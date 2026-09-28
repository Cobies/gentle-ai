# ODD Specification: Migrate Pi MCP Adapter Config to `mcp-adapter.json` (v3.1.0+)

## 1. Diagnosis & Technical Proposal

### Root Cause
`pi-mcp-adapter` v3.0.0+ introduced a breaking change: it renamed its configuration file from `<agentDir>/mcp.json` (and `<workspace>/.pi/mcp.json`) to `<agentDir>/mcp-adapter.json` (and `<workspace>/.pi/mcp-adapter.json`), reserving `mcp.json` for Pi's upcoming native MCP support. When `mcp.json` is present and `mcp-adapter.json` is absent, `pi-mcp-adapter` v3 ignores the configured servers (`codegraph`, `context7`, `engram`, etc.) and emits a runtime migration warning on every session launch:
`Warning: pi-mcp-adapter no longer reads ~/.pi/agent/mcp.json. Move it with: mv "~/.pi/agent/mcp.json" "~/.pi/agent/mcp-adapter.json"`.

Additionally, `gentle-ai` (`internal/agents/pi/adapter.go`) still pins `piMCPAdapterVersion = "2.6.0"` / `^2.6.0` and targets `mcp.json` via `piEngramMCPConfigFile = "mcp.json"`, causing version skew in `<agentDir>/npm/package.json` during sync/install and recreating the deprecated `mcp.json` file whenever `gentle-ai` provisions MCP servers.

### Technical Proposal
1. **Update `internal/agents/pi/adapter.go`**:
   - Update `piMCPAdapterVersion` to `"3.1.0"` and `piMCPAdapterVersionRange` to `"^3.1.0"`.
   - Change `piEngramMCPConfigFile` from `"mcp.json"` to `"mcp-adapter.json"`, keeping a `piLegacyMCPConfigFile = "mcp.json"` constant for automatic migration.
   - Add a legacy migration helper (`MigrateLegacyPiMCPConfig` or inside `ProvisionEngramMCP` / `reconcilePiMCP` / `injectMCPConfigFile`) so that if `mcp.json` exists and `mcp-adapter.json` does not, existing server definitions are preserved/migrated into `mcp-adapter.json` and the legacy `mcp.json` is removed (or merged cleanly if both exist).
   - Update `EffectiveCodeGraphMCPPath` to inspect `.pi/mcp-adapter.json` (with fallback recognition for `.pi/mcp.json`).
2. **Update Tests and Documentation**:
   - Update `internal/agents/pi/adapter_test.go`, `internal/components/communitytool/pi_codegraph_test.go`, `internal/cli/sync_test.go`, `internal/cli/run_integration_test.go`, `internal/components/uninstall/service_test.go`, `internal/update/upgrade/executor_test.go`, and `internal/cli/doctor_test.go` to expect `mcp-adapter.json` and `^3.1.0`.
   - Update `docs/pi.md` and `docs/quickstart.md` where `~/.pi/agent/mcp.json` is documented.

---

## 2. Technical Specification & Contracts

### Constants (`internal/agents/pi/adapter.go`)
```go
const (
	piMCPAdapterPackage         = "npm:pi-mcp-adapter"
	piMCPAdapterPackageSpec     = "npm:pi-mcp-adapter"
	piGentleEngramPackageSource = "npm:gentle-engram"
	piMCPAdapterDependency      = "pi-mcp-adapter"
	piMCPAdapterVersion         = "3.1.0"
	piMCPAdapterVersionRange    = "^3.1.0"
	piAppendSystemFile          = "APPEND_SYSTEM.md"
	piEngramMCPConfigFile       = "mcp-adapter.json"
	piLegacyMCPConfigFile       = "mcp.json"
	piSettingsFile              = "settings.json"
	piNPMDirectory              = "npm"
	piNPMPackageFile            = "package.json"
)
```

### Legacy Config Migration Contract
When `ProvisionEngramMCP(homeDir)` or MCP config writers run for Pi:
- Check `<agentDir>/mcp.json`. If it exists and `<agentDir>/mcp-adapter.json` does not exist, rename/migrate `<agentDir>/mcp.json` to `<agentDir>/mcp-adapter.json` before applying overlays so user/engram servers are preserved and `pi-mcp-adapter` v3 stops emitting warnings.

---

## 3. Tasks & Evidence

- [x] **Task 1**: Update `internal/agents/pi/adapter.go`, `internal/components/communitytool/pi_codegraph.go`, and docs (`docs/pi.md`, `docs/quickstart.md`) to use `mcp-adapter.json` and `pi-mcp-adapter` `^3.1.0`, with automatic migration from legacy `mcp.json`.
  - **Evidence**: Updated `piMCPAdapterVersion` to `"3.1.0"`, `piMCPAdapterVersionRange` to `"^3.1.0"`, `piEngramMCPConfigFile` to `"mcp-adapter.json"`, added `MigrateLegacyPiMCPConfig`, and updated workspace fallback in `EffectiveCodeGraphMCPPath`.
- [x] **Task 2**: Update unit and integration tests (`internal/agents/pi/adapter_test.go`, `internal/components/communitytool/pi_codegraph_test.go`, `internal/cli/*_test.go`, `internal/components/uninstall/service_test.go`, `internal/update/upgrade/executor_test.go`) and verify all Go tests pass.
  - **Evidence**: Added `TestProvisionEngramMCPMigratesLegacyMCPConfig` and `TestEffectiveCodeGraphMCPPathWorkspaceFallback` in `internal/agents/pi/adapter_test.go`; verified `go test -count=1 ./internal/agents/pi/... ./internal/components/communitytool/... ./internal/components/engram/... ./internal/components/uninstall/... ./internal/update/upgrade/...` and `./internal/cli -run "Test.*Pi|Test.*CodeGraph|Test.*MCP"` pass.

---

## 4. Acceptance Criteria & Verification

1. `Adapter.MCPConfigPath(homeDir, ...)` and `CodeGraphPaths(homeDir).MCPConfig` resolve to `<agentDir>/mcp-adapter.json`.
2. `ProvisionEngramMCP` writes `"pi-mcp-adapter": "^3.1.0"` into `<agentDir>/npm/package.json` and migrates any legacy `<agentDir>/mcp.json` to `<agentDir>/mcp-adapter.json`.
3. `go test ./internal/agents/pi/... ./internal/components/communitytool/... ./internal/components/uninstall/... ./internal/update/upgrade/...` pass with zero failures.
