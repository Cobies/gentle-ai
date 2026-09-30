# ODD Task Specification: Dense Environment HUD & Zero-Scroll Sidebar Layout

## 1. Diagnosis & Technical Proposal

### Problem
In `gentle-pi`, when running in fullscreen mode ($\ge 140$ columns), the sidebar rail stacks `hud`, `footer` (Status), `agents`, and `todo`.
- The current `renderHudCard` consumes 16–19 vertical lines due to uppercase multi-line subheaders (`[ PROJECT TARGET ACTIVE ]`, `[ MODEL CONTEXT PROTOCOL ]`, `[ EXECUTION TELEMETRY ]`), blank separator lines, and vertical single-line MCP server listings.
- The `footer` (Status card, `renderShellSidebarBar`) directly duplicates the target and changes telemetry (CWD, Branch, Profile, and Git Diff/Changes) already rendered in the HUD card.
- In standard terminal heights (35–45 lines), this consumes 32–36 vertical lines before the `todo` card is even rendered. When ODD tasks are populated, the `ScrollView` overflows, triggering an unsightly vertical scrollbar that pushes tasks below the viewport fold.

### Technical Proposal & Architecture
1. **Ultra-Compact Dense HUD (`lib/shell-hud.ts`)**:
   - Redesign `renderHudCard` into a 5-line high-density card:
     - **Line 1 (Top border)**: `╭─ ✿ ENVIRONMENT HUD ─────────────────────────────╮`
     - **Line 2 (Target & Git)**: `${shortCwd} · ${branch} [${profile}] · ${diffSummary}`
     - **Line 3 (MCP Ecosystem)**: `MCP (${readyCount}/${total} · ${tools} tools) ${serverGlyphs}` (horizontal status badges with `●`/`○`/`✖`).
     - **Line 4 (Execution Telemetry)**: `${cost} (${latency}ms) · Ctx ${gauge} ${tokens} (${pct}%)` (combining cost, turn latency, and gauge into a single row).
     - **Line 5 (Bottom border)**: `╰─────────────────────────────────────────────────╯`
   - Omit the blank lines and redundant uppercase subheadings.
   - Retain backward-compatible fallbacks and bound truncation to prevent horizontal wrapping.

2. **Deduplicate Footer / Status in Sidebar (`lib/shell-bar.ts`)**:
   - In `renderShellSidebarBar`, when `hud` is active, avoid duplicating the `Project` and `Changes` groups, or present only unique integration/RDD status lines.

3. **Update & Extend Tests (`tests/shell-hud.test.ts`, `tests/shell-sidebar-layout.test.ts`)**:
   - Verify the 5-line bounded card height.
   - Verify that all telemetry data points (target, diff, MCP counts/glyphs, cost, latency, gauge) remain intact and tested.
   - Verify clean integration into `renderSidebarRail` without triggering overflow scrollbars.

---

## 2. Technical Specification & Contracts

### Allowed Edit Surfaces
- `/mnt/c/Users/Cobies/Desktop/Proyectos/COLAB/gentle-pi/lib/shell-hud.ts`
- `/mnt/c/Users/Cobies/Desktop/Proyectos/COLAB/gentle-pi/lib/shell-bar.ts`
- `/mnt/c/Users/Cobies/Desktop/Proyectos/COLAB/gentle-pi/tests/shell-hud.test.ts`
- `/mnt/c/Users/Cobies/Desktop/Proyectos/COLAB/gentle-pi/tests/shell-sidebar-layout.test.ts`
- `/mnt/c/Users/Cobies/Desktop/Proyectos/COLAB/gentle-pi/odd/tasks/dense-environment-hud.md`

---

## 3. Tasks & Evidence

- [ ] Task 1: Refactor `lib/shell-hud.ts` to implement the 5-line Dense HUD layout and horizontal MCP glyphs.
  - Evidence: Pending
- [ ] Task 2: Deduplicate redundant target/changes output in sidebar footer (`lib/shell-bar.ts`) when HUD is rendered.
  - Evidence: Pending
- [ ] Task 3: Update unit test suites (`tests/shell-hud.test.ts`, `tests/shell-sidebar-layout.test.ts`) and verify zero regressions.
  - Evidence: Pending

---

## 4. Acceptance Criteria & Verification

1. `renderHudCard` outputs exactly 5 lines (frame + 3 content lines) when rendered with normal telemetry.
2. All critical telemetry metrics (CWD, branch, profile, diff stats, MCP ready count & tool count, session cost, turn latency, context gauge) are visible and readable.
3. Total rail height prior to TODO is reduced by $\ge 15$ lines, eliminating vertical scrollbar when standard ODD task sets (4–8 tasks) are active.
4. All unit tests in `tests/shell-hud.test.ts` and `tests/shell-sidebar-layout.test.ts` pass cleanly.
