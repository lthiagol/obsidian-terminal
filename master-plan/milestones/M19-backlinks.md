# M19 — Backlinks Panel

**Status:** ⏳ pending

## Goal

Show notes linking TO the current note. Build a reverse wiki-link index during `ScanVault`, display as panel below viewer, Enter to navigate.

## Implementation Plan

### 1. Build reverse index in `ScanVault` (`vault.go`)

Add `var wikiLinkRawRe = regexp.MustCompile(`\[\[([^\]|#]+)`)`  
New function `extractWikiLinkTargetsFromRaw(content string) []string` — regex extracts all `[[target` references, normalizes (lowercase, append .md), deduplicates.

Change `ScanVault` signature: add `map[string][]string` return (normalized target → source paths). After walk loop, for each file's body, extract targets and populate index.

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

Add fields: `backlinkIndex map[string][]string`, `backlinkPanel BacklinkPanel`, `backlinkMode bool`

In `NewModel` and `rescanVault`: capture new return value from `ScanVault`.

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

1. Add `extractWikiLinkTargetsFromRaw()` + regex to vault.go
2. Change ScanVault signature, build backlink index
3. Update NewModel/rescanVault in model.go
4. Create `backlinks.go`
5. Add fields to Model, populate on note open
6. Add backlink focus mode in handleViewKey
7. Split viewer area in View()
8. Add BacklinkToggle to KeyMap
9. Update statusbar + help
10. Write tests
