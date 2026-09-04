# Design: Antigravity Static Agents Docs Refresh

## Technical Approach

Prose-only refresh across six surfaces to state the two-tier model already in runtime: `invoke_subagent` against the 18-file pre-registered static set as primary, `define_subagent` as resilient fallback. Maps to proposal Option B and delta spec `antigravity-support` (static-primary requirement + consistency, historical-note, and scope-guard requirements). No behavior, installer, or probe change.

```
Reader ──→ docs/agents.md + README:111 (matrix/delegation wording)
         ──→ openspec/specs/antigravity-support/spec.md (main promotion)
         ──→ docs/antigravity-sdd-workaround.md (dated note + redirect)
         ──→ antigravity_sdd_agents.go:17-21 + sdd-orchestrator.md ~4 lines (comment touch-up)
```

## Architecture Decisions

| Option | Tradeoff | Decision |
|---|---|---|
| Edit wording in place per surface vs. restructure docs | Restructure risks link rot, exceeds budget | Edit in place; keep headings, anchors, file paths |
| Delete workaround doc vs. dated note + redirect | Delete breaks old links | Keep file; dated historical banner + pointer to `docs/agents.md`, normative inline-execution rules struck |
| Touch orchestrator prompt lines vs. freeze prompt asset | Prompt edits risk behavior drift | Touch ~4 stale noun phrases only (L76, L98, L170 + fail-closed noun); keep delegation logic byte-identical |
| Broaden to installer/probe hardening vs. prose-only | Verifiable fallback is valuable but behavior-adjacent | Deferred to follow-up change per proposal out-of-scope |

## File Changes

| File | Action | Old → New strategy |
|---|---|---|
| `docs/agents.md:19` | Modify | `Solo-agent + Mission Control` → `Full (static subagents) + Mission Control` |
| `docs/agents.md:40,43` | Modify | Move `Antigravity` from `Solo-agent` row to `Full (sub-agents)` row |
| `docs/agents.md:62-64` | Modify | `custom sub-agent creation not yet available, phases run inline` → `invoke_subagent` vs static set primary, `define_subagent` fallback, inline forbidden |
| `docs/agents.md:177-183` | Modify | Add static path (`agents/` under `~/.gemini/antigravity-cli/`); keep `GEMINI.md` collision warning verbatim |
| `README.md:111` | Modify | `Full (dynamic subagents)` + `Dynamic runtime subagents via define_subagent/invoke_subagent` → `Full (static subagents + dynamic fallback)` + `Pre-registered agents via invoke_subagent, define_subagent fallback` |
| `openspec/specs/antigravity-support/spec.md:30-46` | Modify | Promote per delta spec: static-primary requirement, both scenarios `invoke_subagent` primary + `define_subagent` fallback |
| `docs/antigravity-sdd-workaround.md` | Modify | Dated banner (superseded, date, pointer to `docs/agents.md`); demote inline-execution rules to historical, non-normative |
| `internal/components/sdd/antigravity_sdd_agents.go:17-21` | Modify | `no static registry, define/invoke dynamic-only` → two-tier; keep hardening message + CLI-only plugin wording intact |
| `internal/assets/antigravity/sdd-orchestrator.md` (~4 lines: L76, L98, L170, fail-closed noun) | Modify | `dynamic subagent context` → invoked-subagent wording; `does not install static subagent files` → installed static set; keep fail-closed semantics |

## Interfaces / Contracts

No new interfaces. Wording contract for every touched surface: "static-primary via `invoke_subagent`, `define_subagent` fallback; never inline." Guards: `GEMINI.md` collision warning verbatim-valid; hardening stays CLI-scoped (desktop excluded).

## Testing Strategy

| Layer | What | Approach |
|---|---|---|
| Unit/regression | No behavior change | `go test ./internal/components/sdd/... ./internal/agents/antigravity/...` must pass unchanged |
| Consistency | Zero stale claims | `rg -n "Solo-agent.*Antigravity|Antigravity.*Solo-agent|dynamic-only|no static (sub-agent )?registry|custom sub-agent creation is not yet available|does not install static subagent files|Full \(dynamic subagents\)"` on touched surfaces returns zero hits |
| Links | Old links resolve | Open workaround doc via old path; banner + redirect present |

## Threat Matrix

N/A — no routing, shell, subprocess, VCS/PR automation, executable-file classification, or process-integration boundary. Prose-only; research rev2 external gap stays an assumption, not a design input.

## Migration / Rollout

No migration. Single-commit revert restores all surfaces including the workaround file from git.

## Open Questions

None blocking. Assumption: in-repo evidence (`exploration.md`, commit `398b0f51`, memory #5315) suffices; no external corroboration.
