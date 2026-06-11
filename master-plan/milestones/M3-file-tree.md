# M3 — File Tree Navigator

**Status:** ✅ done

## Goal

Implement the left-panel file tree using `bubbles/tree`. Navigate folders/notes,
expand/collapse directories, select and open notes.

## Files to create

- `tree.go` / `tree_test.go`

## Steps

### 1. `tree.go`
- `FileTree` struct:
  ```go
  type FileTree struct {
      model tree.Model
      items []*VaultEntry  // parallel array matching tree nodes
  }
  ```
- `NewFileTree(vault *VaultEntry) FileTree`:
  - Convert VaultEntry tree → bubbles/tree nodes recursively
  - Root node = vault path base name (expanded by default)
  - Dir nodes: set `node.Collapsed = false` initially (root expanded, subdirs collapsed)
  - File nodes: leaf nodes (only `.md`/`.markdown`, already filtered by scanner)
  - Custom `Less` for folder-first, alphabetical sort
  - Store `*VaultEntry` in node `Value` field for path lookup
  - Sets custom item renderer via bubbletea styles:
    - Folder: `▸ folderName` (amber) or `▾ folderName` (amber) when expanded
    - File: `◇ filename.md` (gray)
    - Selected: violet highlight background
    - Active (open): fuchsia highlight
- Methods:
  - `Update(msg tea.Msg) (FileTree, tea.Cmd)` — delegates to bubbles/tree
  - `View() string` — renders tree within styled panel
  - `SelectedEntry() *VaultEntry` — returns the entry at current cursor position
  - `IsDirSelected() bool`
  - `SelectedPath() string`
  - `SetSize(width, height int)`

### 2. Navigation (in model.go `handleBrowseKey`)
- `j/↓` → `tree.MoveDown()`
- `k/↑` → `tree.MoveUp()`
- `g/Home` → `tree.MoveTop()`
- `G/End` → `tree.MoveBottom()`
- `l/→` → expand selected folder; no-op on file
- `h/←` → collapse selected folder; no-op on file (don't navigate to parent in v1)
- `Enter` on directory → toggle expand/collapse
- `Enter` on `.md` file → `LoadNote()`, set `activeNote`, switch mode to "view"
- `PgUp/PgDn` → scroll tree by half visible height

### 3. Edge cases
- **Empty vault** (no `.md` files): show "no notes found" centered in tree panel
- **Single file**: tree shows just that file, root expanded
- **Deeply nested** (10+ levels): tree width accommodates indentation; overflow handled by panel width truncation
- **Very long filenames**: truncated to panel width minus indent
- **Symlinks**: shown as-is with symlink indicator (e.g., `◇ readme-symlink.md → readme.md`)

## Test Spec (6 tests)

| # | Test | File | Description |
|---|------|------|-------------|
| 1 | `TestTree_RootExpanded` | tree_test.go | Root folder starts expanded, immediate children visible |
| 2 | `TestTree_ExpandCollapse` | tree_test.go | Right/→ expands folder; Left/← collapses folder |
| 3 | `TestTree_SelectionClamped` | tree_test.go | Selection stays in bounds: Down at bottom stays; Up at top stays |
| 4 | `TestTree_SymlinkShownInTree` | tree_test.go | testdata symlink (readme-symlink.md) appears as file node |
| 5 | `TestTree_OpenNoteSwitchesToView` | model_test.go | Enter on .md file → mode="view", activeNote set to correct note |
| 6 | `TestTree_EnterOnFolderExpands` | model_test.go | Enter on collapsed directory → expands (children visible) |

## Completion Criteria

- [x] File tree renders all vault `.md`/`.markdown` files with correct nesting
- [x] Folders expand/collapse via `l/→` and `h/←`
- [x] Selection via `j/k/↑/↓` and `g/G`
- [x] `Enter` on `.md` opens note in viewer (mode → "view")
- [x] Symlinks shown in tree (IsSymlink flag tracked)
- [x] Empty vault shows placeholder message
- [x] All 6 tests pass
- [x] `go vet ./...` exits 0

## Verification Evidence

- `go build ./...` exits 0
- `go test ./...` — 24/24 tests pass (6 new M3 tests)
- `go vet ./...` exits 0
- Files created: `tree.go`, `tree_test.go`
- FileTree uses custom interactive component (bubbles/tree not available)
- Integration: navigate, expand/collapse dirs, Enter on file opens viewer
