# M19 — Backlinks Panel

**Status:** ✅ done

## Goal

Show notes linking TO the current note. Use the backlink index from VaultIndexes (built in M18.5), display as panel below viewer, Enter to navigate.

## Keybinding

**Key:** `b`
**Mode:** View mode only
**Rationale:** Mnemonic for "backlinks"

See [KEYBINDINGS.md](../../KEYBINDINGS.md) for complete keybinding reference.

## Prerequisites

- M18.5 (Vault Index System) must be completed first

## Implementation Plan

### 1. Use backlink index from VaultIndexes

The backlink index is already built in M18.5. Access via `m.backlinkIndex` (populated from `indexes.Backlinks`).

No changes needed to `ScanVault` - just use the existing index.

### 2. New file: `backlinks.go`

```go
type BacklinkPanel struct { links []string; cursor int; width int }
func NewBacklinkPanel(notePath string, backlinkIndex map[string][]string) BacklinkPanel
func (bp *BacklinkPanel) MoveUp/Down()
func (bp BacklinkPanel) SelectedPath() string
func (bp BacklinkPanel) Count() int
func (bp BacklinkPanel) View() string  // styled list with cursor highlight
```

### 3. Model changes (`model.go`)

Add fields: `backlinkPanel BacklinkPanel`, `backlinkMode bool`

Note: `backlinkIndex` is already in Model from M18.5.

In `handleBrowseKey` + `handleViewKey` note-load paths: populate `m.backlinkPanel = NewBacklinkPanel(note.Path, m.backlinkIndex)`.

### 4. handleViewKey additions (`handlers.go`)

Add `BacklinkToggle` rune (`b`) to `KeyMap`. When pressed: `m.backlinkMode = !m.backlinkMode`.

When `backlinkMode` true: arrow keys move backlink cursor, Enter navigates to source note, Esc returns focus.

### 5. View() split layout (`model.go`)

When viewing note with backlinks: divide right panel — top 70% viewer, separator border, bottom 30% backlinks.

### 6. Status bar + help

View mode hint: `"b backlinks | ..."`  
Help: add Backlinks section

### Edge cases

- Self-referencing links (note links to itself) → appears in backlinks (legitimate)
- Dead links → LoadNote fails, toast shown
- Rescan → rebuild backlinkPanel for current note
- No backlinks → "No backlinks" text

### Implementation order

1. Create `backlinks.go` with BacklinkPanel
2. Add backlinkPanel and backlinkMode fields to Model
3. Populate backlinkPanel on note open
4. Add backlink focus mode in handleViewKey
5. Split viewer area in View()
6. Add BacklinkToggle to KeyMap
7. Update statusbar + help
8. Write tests

## Completion Criteria

- [x] `backlinks.go` created with BacklinkPanel type
- [x] BacklinkPanel displays list of notes linking to current note
- [x] `b` keybinding toggles backlink focus in View mode
- [x] Enter on backlink navigates to source note
- [x] Split layout: viewer top 70%, backlinks bottom 30%
- [x] "No backlinks" shown when none exist
- [x] Backlinks rebuild on rescan
- [x] Help text updated
- [x] KEYBINDINGS.md updated
- [x] `make test` passes
- [x] `make vet` exits 0
- [x] Manual test: backlinks work for notes with incoming links
