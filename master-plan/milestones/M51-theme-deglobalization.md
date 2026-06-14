# M51 — Complete Theme De-globalization

**Status:** ✅ done (finishes partial **M37**)  
**Finding:** F-2 in [ARCHITECTURE-REVIEW-2026-06-13.md](../ARCHITECTURE-REVIEW-2026-06-13.md)  
**Execution plan:** [PHASE-12-EXECUTION-PLAN.md](../PHASE-12-EXECUTION-PLAN.md) §6

## Goal

Single source of truth for theme state: **`Model.palette`**. No runtime mutation of package-level color/style globals.

## Problem statement

M37 delivered data-driven `themeData`, `m.palette`, and fixed `switchToProfile` — but widgets still read `Accent`, `TreeStyle`, etc., and `activatePalette()` overwrites globals on every theme switch.

## Out of scope

- New themes or custom_theme format changes
- Markdown renderer color logic (already uses `RendererStyle` from palette)
- Lipgloss version upgrade

## Dependencies

- **After:** M50, M55 (CI catches regressions)
- **Blocks:** M52 WP5–WP6 (render extractions should use palette threading)
- **Not parallel with:** M52

---

## Work packages

### WP1 — Global read audit (1h)

**Steps:**
1. Run: `rg '\b(Accent|AccentSecondary|TreeStyle|ViewerStyle|StatusStyle|HelpStyle|SearchStyle|TextMuted|TextDim|ModeColors)\b' --glob '*.go' -l`
2. Paste file list into milestone completion notes
3. Classify each file: **must fix** vs **test-only** vs **theme.go definition**

**Verification:**
- [ ] Checklist table complete (expect ~15 files)

---

### WP2 — Pass palette into tree + status + help (3h)

**Preferred pattern** — avoid more globals:

```go
// Option A: methods take Palette
func (t FileTree) View(p Palette) string

// Option B: FileTree stores palette (updated on theme switch)
func (t *FileTree) SetPalette(p Palette)
```

**Steps:**
1. Choose Option A or B (document in WP1 notes — **recommend Option B** for fewer signature changes)
2. Update `FileTree`, `renderStatusBar`, `renderHelp` to use `m.palette`
3. Update call sites in `View()` / `Update()`

**Verification:**
- [ ] `tree_test.go`, status bar tests pass
- [ ] Switch theme in `TestStatePreservation_ThemeSwitch` still passes

---

### WP3 — Handlers, toast, overlays (2h)

**Files:** `toast.go`, `command_palette.go`, `tags.go`, `backlinks.go`, `profile_picker.go`, `handlers.go` (`renderInNoteSearch`), `model.go` render helpers

**Steps:**
1. Replace global color reads with `m.palette` or passed `Palette`
2. Replace `ModeColors[mode]` with palette mode color fields

**Verification:**
- [ ] `command_palette_test.go`, `custom_theme_test.go` pass

---

### WP4 — Remove runtime global mutation (1h)

**Steps:**
1. Change `setTheme` to only assign `m.palette`, `m.viewer.renderStyle`, `m.searchStyle`, and call widget `SetPalette` if used
2. Remove `activatePalette(palette)` call from `setTheme` — or reduce to test-only helper
3. `NewModel`: build palette via `lookupPalette`; do not rely on globals being pre-set

**Verification:**
- [ ] Grep: `activatePalette` only in tests or deprecated comment
- [ ] Profile switch + custom theme tests pass

---

### WP5 — Shrink globals + document (1h)

**Steps:**
1. Keep in `theme.go`: `themeData`, `lookupPalette`, `rebuildDerivedStyles`, `IconFolderOpen` constants (non-mutable)
2. Remove or mark deprecated: mutable `var Accent`, `TreeStyle`, …
3. Update M37 milestone → ✅ only after this WP
4. M53 will update AGENTS styling section

**Verification:**
- [ ] `make test && make vet` pass
- [ ] Two consecutive theme switches in test show different `m.palette` without relying on global side effect

---

## Acceptance criteria

- [ ] All 5 WPs complete
- [ ] Zero runtime reads of mutable theme globals in non-test code (grep verification)
- [ ] M37 can be marked ✅ in STATUS (supersedes partial)
- [ ] `make test && make vet` pass

## Rollback / risk

| Risk | Mitigation |
|------|------------|
| Missed global read → wrong colors after switch | WP1 audit + theme integration test |
| Large diff | One WP per commit; CI from M55 |

## Handoff notes

Do not extract files (M52) in same PR as WP2–WP4. Thread palette first, split files second.

## Estimated total

8–10 hours (1–2 days)

## Priority

🟡 High (Track B)
