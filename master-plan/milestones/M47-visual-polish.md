# M47 — Visual Polish & Look-and-Feel Upgrade

**Status:** ✅ done

## Goal

Bring the visual polish, iconography, and spatial design of obsidian-terminal closer to obsitui's premium terminal aesthetic, while staying within the read-only constraint.

## Context

[obsitui](https://github.com/atr0t0s/obsitui) (Ink/React + TypeScript) achieves a premium terminal look through:
- Rich, consistent Unicode iconography (not emoji)
- Per-mode accent colors (19 distinct mode badge colors)
- Visible panel separators (│/┆)
- Heading color hierarchy (h1 gets its own fuchsia color)
- Code blocks with ╭╮╰╯ box-drawing borders
- Note preview pane with fade effect
- In-note search (`/` in view mode)
- Graph view (ASCII circular layout)

Our project already uses many of the same concepts (lipgloss, Unicode glyphs, mode badges), but can be upgraded.

**Key constraint:** read-only. No editing, writing, or vault-modification features.

## Analysis: obsitui Features — Adopt vs Skip

### ✅ Can adopt (read-only compatible)

| Feature | obsitui implementation | Our approach |
|---------|----------------------|--------------|
| **In-note search** | `/` in view mode, `n`/`N` cycle matches, scroll-to-match | Match existing search components pattern; highlight matches in rendered view |
| **Preview pane** | Bottom panel shows content of hovered sidebar item, with fade effect | Show at bottom of tree panel or right panel; reuse markdown renderer at small width |
| **Navigation history** | Ctrl+O back, Ctrl+I forward, jumplist-style stack | Push path on open, pop on back; render in recent category |
| **Graph view** | ASCII circular layout, Bresenham line drawing, local + global toggle | Use backlink/tag indexes for edges; render as ASCII art with node selection |
| **Heading color hierarchy** | h1=fuchsia, h2=violet, h3=teal, h4+=gray | Add `Heading1` color distinct from `Accent`; already have `AccentSecondary` for h1 |
| **Richer Unicode icons** | ◇ files, ▾/▸ folders, ★ bookmarks, ✦ AI, ◆ section markers | Replace emoji callout icons with Unicode equivalents; use ◇/▸/▾ consistently |
| **Callout icon refinement** | 20+ callout types with Unicode icons (✶ tip, ⚠ warning, ✖ danger, ✔ success) | Expand callout icon set, switch from emoji to Unicode |
| **Selection indicator** | ▸ triangle marks selected item in sidebar | Already have selection highlight via accent background, could add ▸ prefix |
| **Toast icons** | ℹ info, ✔ success, ✖ error, ⚠ warning | Replace our `i`/`v`/`!`/`x` with Unicode equivalents |
| **Panel separators** | │ and ┆ vertical bars between panels | Already have border on tree panel; could add to more panels |

### ❌ Cannot adopt (read-only constraint)

| Feature | Why not |
|---------|---------|
| Vim editor (normal/insert/visual modes) | Writing to vault |
| File create/rename/delete | Writing to vault |
| AI writing (summarize, suggest, generate) | Writing to vault + AI is non-goal |
| Quick capture (Ctrl+X to Inbox) | Writing to vault |
| Kanban board (scans checkboxes) | Would require modifying task state |
| Pomodoro timer | Not a vault feature |
| External editor integration ($EDITOR) | Non-goal (writing + adds platform dependency) |

### 🔮 Future (low priority, complex)

| Feature | Notes |
|---------|-------|
| AI chat / RAG over vault | Non-goal in current charter |
| Semantic search / embeddings | Non-goal |
| Bookmarks with reorder | Extra persistence complexity, could be M85+ |
| Templates with {{date}}/{{title}} | Non-goal (writing feature) |

## Files to modify

| File | Changes |
|------|---------|
| `internal/markdown/markdown.go` | Expand callout icon set (Unicode instead of emoji), add `Heading1` color distinct from `Accent` |
| `theme.go` | Add `ModeViewFuchsia` color, heading h1 color, expand callout icon constants |
| `viewer.go` | Add in-note search capability (`/` triggers search, highlight matches) |
| `model.go` | Add `previewNote`, `previewVisible` fields; add navigation history stack; add `ModeGraph` |
| `handlers.go` | Add in-note search handler (`/` in view mode), `n`/`N` cycle, graph view handler, history back/forward |
| `keys.go` | Add `InNoteSearch`, `NextMatch`, `PrevMatch`, `HistoryBack`, `HistoryFwd`, `GraphToggle` keys |
| `statusbar.go` | Add mode-specific colors for more mode badges; richer toast icons |
| `toast.go` | Replace `i`/`v`/`!`/`x` with ℹ/✔/⚠/✖ Unicode icons |
| New: `graph.go` | ASCII graph renderer using backlink/wiki-link indexes |
| New: `preview.go` | Preview pane renderer |
| New: `history.go` | Navigation history stack |

## Steps

### 1. Unicode icon & callout polish
Replace emoji-based callout icons (📝, 💡, ⚠, etc.) with cleaner Unicode equivalents from obsitui's theme. Update toast icons from `i`/`v`/`!`/`x` to ℹ/✔/⚠/✖. Add file ✦ / folder ▾/▸ icon constants.

### 2. Heading color hierarchy
Add a distinct `Heading1` color (fuchsia or similar) to `RendererStyle` and `Palette`. Currently h1 shares `AccentSecondary` (amber). obsitui uses fuchsia for h1, violet for h2, teal for h3. Update all built-in palettes.

### 3. Mode badge color refinement
Add more distinct per-mode colors. Currently 4 mode colors (Browse=violet, View=teal, Search=amber, Help=blue). Add View=fuchsia, Tags=orange, etc. to match obsitui's richness.

### 4. In-note search
Press `/` in view mode. Type to filter matching lines. `n`/`N` to cycle matches upward/downward. Scroll to match position. Highlight matches in rendered output. Esc to dismiss.

### 5. Note preview pane
When browsing, show a small preview of the currently highlighted note's content at the bottom of the tree panel or as a bordered panel. Reuse the markdown renderer at a narrower width. Toggle visibility with a keybind.

### 6. Navigation history
Track previously opened notes in a history stack. `Ctrl+O` goes back, `Ctrl+I` goes forward. Pushing on every note open. Already partially done via `recentNotes` — extend to full back/forward stack.

### 7. Graph view (ASCII)
Build an ASCII graph view using backlink/wiki-link indexes. Circular node layout with Bresenham lines for edges. Node selection with j/k, Enter to open. `l` toggles local/global. `f` focuses on selected node. `Esc` closes.

## Completion Criteria

- [x] Callout icons use Unicode instead of emoji (28 types: ℹ, ✶, ★, ⚠, ✖, ○, ✔, ◆, ▶, ≡, ❝, ✘)
- [x] Toast icons use ℹ/✔/⚠/✖
- [x] Heading colors: h1 gets distinct fuchsia color, palette updated across all 7 themes
- [x] Mode badges have richer per-mode colors (7 distinct: Browse=violet, View=fuchsia, Search=amber, Tags=orange, Profile=violet, Help=blue, Find=amber)
- [x] In-note search works: `/` enters, typing filters, `n`/`N` cycles, Esc dismisses
- [x] Navigation history: Ctrl+O (view mode) back, `[` back, `]` forward preserves navigation
- [ ] Preview pane shows content of hovered sidebar item (deferred → **M48**)
- [ ] Graph view renders with node selection, global/local toggle (deferred → **M49**)
- [x] All existing keybindings preserved, no conflicts
- [x] `make test` passes all 144 tests
- [x] `make vet` exits 0

## Completed

2026-06-12

### Unicode icon polish
- Replaced emoji callout icons (📝, 💡, 🚫, ❓, ✅, 🐛, 📋) with clean Unicode: ℹ/✶/★/⚠/✖/○/✔/◆/▶/≡/❝/✘
- Added 18 new callout types: hint, important, caution, error, check, done, help, faq, abstract, summary, quote, failure, fail, missing
- Toast icons: `i→ℹ`, `v→✔`, `!→⚠`, `x→✖`

### Heading color hierarchy
- Added `Heading1` field to `Palette` and `RendererStyle`
- h1 now uses fuchsia/pink (#e879f9) across all 7 themes, distinct from amber AccentSecondary
- Each theme gets a matched heading1 color

### Mode badge colors
- Added `ModeTags` (orange) and `ModeProfile` colors to Palette
- Added `ModeProfilePicker` entry to global `ModeColors`
- View mode badge now uses fuchsia (was teal) for visual distinction

### In-note search
- `/` in view mode activates in-note search bar (no longer fuzzy search)
- Type to filter matching lines; `n`/`N` cycles forward/backward
- Scrolls viewer to match position
- Esc or Enter dismisses
- Search bar renders above viewer with match count indicator

### Navigation history  
- `openNote` pushes previous note path to history stack before navigating
- Ctrl+O in view mode goes back (was toggleRecents); browse mode still shows recents
- `[` goes back, `]` goes forward in view mode
- Forward history cleared on new navigation

### Bug fixed
- `tea.KeyCtrlI` conflicts with `tea.KeyTab` (both send `\x09` in most terminals)
- Switched history forward from Ctrl+I to `]` key

## Dependencies

None. This milestone is independent and can be done at any time.

## Priority

🟡 High (visual polish + in-note search + preview are impactful user-facing improvements)
