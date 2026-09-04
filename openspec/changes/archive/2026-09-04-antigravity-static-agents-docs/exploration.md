## Exploration: antigravity-static-agents-docs

### Current State
Antigravity already runs a two-tier delegation model: 18 static subagent files ship embedded under `internal/assets/antigravity/agents/` (10 SDD phases + 5 review lenses + jd-judge-a/b + jd-fix-agent), each with `subagent: true` frontmatter, and the installer copies them to `SubAgentsDir` (`<GlobalConfigDir>/agents/`, covered by `TestInjectAntigravitySubAgents`). `Adapter.SupportsSubAgents()` returns `true` via `FileSubAgents: true` (`internal/agents/antigravity/adapter.go:147-157`; manifest `internal/agents/capabilitymanifest/manifest.go:291-293`). The orchestrator asset (`internal/assets/antigravity/sdd-orchestrator.md:7,66-67`) mandates `invoke_subagent` against the pre-registered static set as the primary route with `define_subagent` only as resilient fallback, and the hardening message (`antigravity_sdd_agents.go:59`) already states "static subagent invocation as primary with dynamic subagent creation as resilient fallback". Docs and spec lag behind: `docs/agents.md:19,44,62-64` still lists Antigravity as "Solo-agent + Mission Control" with "custom sub-agent creation not yet available"; `docs/antigravity-sdd-workaround.md` still prescribes single-threaded inline execution with filesystem save-state; `README.md:111` still says "Full (dynamic subagents)"; `openspec/specs/antigravity-support/spec.md:30-46` still requires "dynamic subagents rather than static subagent files" (`define_subagent` then `invoke_subagent`). One stale code comment also lags: `antigravity_sdd_agents.go:17-21` still claims "Antigravity has no static sub-agent registry". Git chronology confirms the inversion: static definitions landed via `398b0f51 feat(antigravity): add static subagent definitions and sync support` plus follow-ups (`5db7014c`, `b9b95d1d`, `57decb2f`, `16382faa`), after the July-2026 dynamic-only era.

### Affected Areas
- `docs/agents.md` — agent matrix row, delegation-models table, and "Antigravity + Mission Control" section claim solo-agent/dynamic-unavailable
- `docs/antigravity-sdd-workaround.md` — entire single-threaded workaround premise (inline phases, `.sdd/*.md` save-state, Engram-restricted-to-global-decisions) contradicts static delegation
- `README.md` — agent matrix row says "Full (dynamic subagents)"
- `openspec/specs/antigravity-support/spec.md` — "Antigravity uses dynamic subagents" requirement + both scenarios mandate dynamic-only
- `internal/components/sdd/antigravity_sdd_agents.go` — header comment (lines 17-21) claims no static registry; hardening message itself is already correct
- `internal/assets/antigravity/sdd-orchestrator.md` — already two-tier; only ~4 stale lines (`170`, `76`, `98`, `501` area wording: "runs in dynamic subagent context", "does not install static subagent files") need prose touch-up, no behavior change

### Approaches
1. **Minimal docs + spec refresh to two-tier** — rewrite the five surfaces above to describe static-primary + dynamic-fallback; fix the stale header comment; touch up the 4 stale orchestrator lines
   - Pros: smallest diff, zero behavior change, fully within 400-line review budget
   - Cons: leaves the legacy workaround doc discoverable (redirect needed at minimum)
   - Effort: Low

2. **Refresh + archive legacy workaround doc** — option A plus convert `docs/antigravity-sdd-workaround.md` into a dated historical note pointing at `docs/agents.md`
   - Pros: removes the most misleading surface (it prescribes inline execution, the exact behavior the orchestrator now forbids); prevents new users following dead guidance
   - Cons: slightly larger diff; needs a redirect pointer so old links don't rot
   - Effort: Low

3. **Refresh + installer/probe hardening** — options A/B plus runtime verification that static agents actually resolve at install/sync time (probe, warn when registry misses force fallback)
   - Pros: closes the observability gap (today fallback is silent); makes the two-tier claim verifiable
   - Cons: behavior-adjacent code changes, new tests, exceeds a docs-fix change; should be its own change
   - Effort: Medium

### Recommendation
Option B. The core fix is prose-only (A), but the workaround doc is actively harmful — it instructs inline execution that the orchestrator's fail-closed contract (`sdd-orchestrator.md:58`) now refuses — so archiving it in the same change is justified and still low-effort. Defer option C to a follow-up change with its own spec.

### Risks
- Shared `~/.gemini/GEMINI.md` surface with Gemini CLI: rewording must not imply Gemini CLI inherits Antigravity static agents (collision warning in `docs/agents.md:181` stays valid)
- Desktop vs CLI plugin asymmetry (`antigravity_sdd_agents.go:34-38`): docs must not claim the hardening plugin covers the desktop variant
- Engram memory `#4470/#4520` (July 2026 dynamic-only) vs `#5315` (Aug 2026 static era): proposal should cite `#5315` and the `398b0f51` commit as superseding evidence
- No CRITICAL risks; no source behavior changes proposed

### Ready for Proposal
Yes — orchestrator should tell the user: exploration confirms the two-tier model is real and tested (18 files, `SupportsSubAgents()==true`, `TestInjectAntigravitySubAgents`), the fix is docs+spec+comment prose only, recommend scope B (refresh + archive workaround), and ask for confirmation to proceed to `propose`.
