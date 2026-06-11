# M20 — Tag Browsing & Filtering

**Status:** ⏳ pending

## Goal

Add a tag browser mode to explore all vault tags and filter the file tree by selected tag.

## Implementation Plan

### 1. Build tag index in `ScanVault` (`vault.go`)

Change `ScanVault` signature: add `map[string][]string` return (tag → file paths).

During WalkDir, after setting `searchIndex[relPath]`, parse frontmatter and extract tags. Normalize: strip leading `#`, lowercase, skip empty.

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

Add fields: `tagIndex map[string][]string`, `tagList TagList`, `tagCursor int`, `tagFilter string`, `treeUnfiltered []treeItem`

Add methods: `applyTagFilter(tag string)`, `clearTagFilter()`

### 4. Tree filtering (`tree.go`)

Add `ApplyPathFilter(paths map[string]bool)` — rebuilds `items` to only include matching files (directories kept if any descendant matches).

### 5. Keybinding + handlers (`handlers.go`, `keys.go`)

Add `BrowseTags rune` (`t`) to KeyMap.  
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

1. Add tag extraction to ScanVault (vault.go)
2. Create tags.go with TagList
3. Add ModeTags + fields to Model
4. Add ApplyPathFilter to tree.go
5. Add applyTagFilter/clearTagFilter to model.go
6. Add BrowseTags key, enterTagsMode, handleTagsKey
7. Wire dispatch + rendering
8. Add ModeTags color to theme.go
9. Update statusbar hints
10. Write tests
