# M20 — Tag Browsing & Filtering

**Status:** ✅ done

## Goal

Add a tag browser mode to explore all vault tags and filter the file tree by selected tag.

## Keybinding

**Key:** `T` (uppercase)
**Mode:** Browse mode only
**Rationale:** Mnemonic for "Tags", uppercase to avoid conflict with `t` (outline) in View mode

See [KEYBINDINGS.md](../../KEYBINDINGS.md) for complete keybinding reference.

## Prerequisites

- M18.5 (Vault Index System) must be completed first

## Implementation Plan

### 1. Use tag index from VaultIndexes

The tag index is already built in M18.5. Access via `m.tagIndex` (populated from `indexes.Tags`).

No changes needed to `ScanVault` - just use the existing index.

### 2. New file: `tags.go`

```go
type TagEntry struct { Name string; Count int; Files []string }
type TagList struct { entries []TagEntry; cursor int; width, height int }
func NewTagList(tagIndex map[string][]string) TagList  // sorted by count desc
func (tl *TagList) MoveUp/Down()
func (tl TagList) SelectedTag() string
func (tl TagList) SelectedFiles() []string
func (tl TagList) View() string  // styled list "#tagname (count)"
```

### 3. Model changes (`model.go`)

Add `ModeTags` to Mode enum + String() case.

Add fields: `tagList TagList`, `tagFilter string`, `treeUnfiltered []treeItem`

Note: `tagIndex` is already in Model from M18.5.

Add methods: `applyTagFilter(tag string)`, `clearTagFilter()`

### 4. Tree filtering (`tree.go`)

Add `ApplyPathFilter(paths map[string]bool)` — rebuilds `items` to only include matching files (directories kept if any descendant matches).

### 5. Keybinding + handlers (`handlers.go`, `keys.go`)

Add `BrowseTags rune` (`T`) to KeyMap.  
In `handleBrowseKey`: `case MatchRune(msg, m.keys.BrowseTags): m.enterTagsMode()`

New handler: `handleTagsKey()` — Esc back, j/k move cursor, Enter applies tag filter and returns to browse mode with toast.

In Update(): `case ModeTags: return m.handleTagsKey(msg)`

### 6. Palette + statusbar

Add `ModeTags lipgloss.Color` to Palette → each palette constructor.

Statusbar: ModeTags shows tag count + active filter.

### Edge cases

- Tags with `#` prefix (some users write `#tag`) → normalize by stripping `#`
- Multiple tags per file → file appears under each tag
- Tag filter + rescan → re-apply filter after rescan
- Nested tags (`status/done`) → stored as-is, exact match filtering

### Implementation order

1. Create tags.go with TagList
2. Add ModeTags + fields to Model
3. Add ApplyPathFilter to tree.go
4. Add applyTagFilter/clearTagFilter to model.go
5. Add BrowseTags key (`T`), enterTagsMode, handleTagsKey
6. Wire dispatch + rendering
7. Add ModeTags color to theme.go
8. Update statusbar hints
9. Write tests

## Completion Criteria

- [x] `tags.go` created with TagList and TagEntry types
- [x] Tag browser shows all tags sorted by count
- [x] `T` keybinding opens tag browser in Browse mode
- [x] Enter on tag filters file tree to show only matching files
- [x] Directories kept if any descendant matches filter
- [x] Tag normalization: `#Tag`, `TAG` all become `tag`
- [x] Multiple tags per file handled correctly
- [x] Filter reapplied after rescan
- [x] Help text updated
- [x] KEYBINDINGS.md updated
- [x] `make test` passes
- [x] `make vet` exits 0
- [x] Manual test: tag browser and filtering work correctly
