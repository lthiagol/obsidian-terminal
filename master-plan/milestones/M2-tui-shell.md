# M2 — Basic TUI Shell

**Status:** ✅ done

## Goal

Set up the Bubble Tea `Model`, mode state machine, key dispatcher (vim + arrows),
screen layout (25/75 split with lipgloss), and theme styles.

## Files to create

- `model.go` / `model_test.go`
- `keys.go` / `keys_test.go`
- `theme.go`

## Steps

### 1. `model.go`
- `Model` struct:
  ```go
  type Model struct {
      mode        Mode       // "browse" | "view" | "search" | "find" | "help"
      prevMode    Mode       // mode to return to on Esc
      vault       *VaultEntry
      activeNote  *VaultNote // nil in browse/search mode
      searchIndex map[string]string // filename → plain text (built in M1)

      // Panels
      tree     FileTree
      viewer   MarkdownViewer
      search   SearchState

      // Layout
      width, height int
      treeWidth     int    // max(25, 25% of width)
      helpOffset    int    // scroll offset for help panel

      // State
      config      *Config
      ready       bool   // true after first WindowSizeMsg
      quitting    bool
  }
  ```
- `NewModel(cfg *Config) Model` — initializes all fields, scans vault
- `Init() tea.Cmd` — returns `tea.SetWindowTitle("obsidian-terminal")` + `tea.EnterAltScreen()`
- `Update(msg tea.Msg) (tea.Model, tea.Cmd)`:
  - `tea.KeyMsg` → dispatch to mode-specific handler
  - `tea.WindowSizeMsg` → recalculate `treeWidth`, set ready=true
  - `tickMsg` → file watcher poll (M7)
- `View()` — compose via lipgloss:
  ```
  leftPanel = tree.View()           (treeWidth wide)
  rightPanel = viewer.View() or search.View() or help.View()
  main = lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, separator, rightPanel)
  full = lipgloss.JoinVertical(lipgloss.Top, main, statusBar)
  ```
- Mode handlers (stubs in M2, implemented in later milestones):
  - `handleBrowseKey(key tea.KeyMsg)` — delegates to tree component
  - `handleViewKey(key tea.KeyMsg)` — delegates to viewer component
  - `handleSearchKey(key tea.KeyMsg)` — delegates to search component
  - `handleFindKey(key tea.KeyMsg)` — delegates to search component
  - `handleHelpKey(key tea.KeyMsg)` — j/k scroll, Esc back

### 2. `keys.go`
- `KeyMap` struct:
  ```go
  type KeyMap struct {
      Up        []tea.KeyType
      Down      []tea.KeyType
      Left      []tea.KeyType
      Right     []tea.KeyType
      Top       []tea.KeyType
      Bottom    []tea.KeyType
      Enter     tea.KeyType
      Esc       tea.KeyType
      Quit      []tea.KeyType
      Search    tea.KeyType
      ContentSearch tea.KeyType
      Help      tea.KeyType
      Tab       tea.KeyType
      PageUp    tea.KeyType
      PageDown  tea.KeyType
      CtrlC     tea.KeyType
  }
  ```
- `MatchKey(msg tea.KeyMsg, keys []tea.KeyType) bool` — matches any key in the list
- `DefaultKeys() KeyMap` — returns the standard mapping (vim + arrows both supported via multi-key arrays)

### 3. `theme.go`
- Colors (hex → lipgloss.Color):
  - `Accent` = `"#a78bfa"` — violet (primary)
  - `AccentSecondary` = `"#fbbf24"` — amber
  - `AccentTertiary` = `"#2dd4bf"` — teal (links)
  - `TextSecondary` = `"#9ca3af"` — gray
  - `TextMuted` = `"#6b7280"` — darker gray
  - `TextDim` = `"#4b5563"` — very dim
  - `Success` = `"#34d399"` — emerald
  - `Warning` = `"#fbbf24"` — amber
  - `Error` = `"#f87171"` — red
  - `Info` = `"#60a5fa"` — blue
- Unicode icons:
  - `IconFolderOpen` = `"▾ "`
  - `IconFolderClosed` = `"▸ "`
  - `IconFile` = `"◇ "`
  - `IconVertical` = `"│"`
  - `IconDiamond` = `"◆"`
- Lipgloss styles:
  - `TreeStyle` — left panel, violet border left
  - `ViewerStyle` — right panel, default
  - `StatusStyle` — bottom bar, dark bg (`#1f2937`)
  - `HelpStyle` — help overlay
  - `SearchStyle` — search overlay
  - `ModeColors` — map[Mode]lipgloss.Color for status bar badges

### 4. Wire up `main.go`
```go
func main() {
    // ... config loading (from M1)
    m := NewModel(cfg)
    p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
    if _, err := p.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```

## Test Spec (7 tests)

| # | Test | File | Description |
|---|------|------|-------------|
| 1 | `TestModel_InitialMode_IsBrowse` | model_test.go | NewModel starts in "browse" mode |
| 2 | `TestKeyDispatch_Browse_JK` | model_test.go | j/k/↑/↓ change tree selection index |
| 3 | `TestKeyDispatch_View_Esc` | model_test.go | Esc in view returns to browse |
| 4 | `TestKeyDispatch_Search_OpenClose` | model_test.go | `/` opens search; Esc closes, restores prevMode |
| 5 | `TestKeyDispatch_Help_Toggle` | model_test.go | `?` opens help; Esc returns to prevMode |
| 6 | `TestModeTransitions` | model_test.go | browse→view, view→browse, browse→search, search→browse, browse→help, help→browse, view→help, help→view |
| 7 | `TestKeyDispatch_BothVimAndArrows` | keys_test.go | j and ↓ both match Down; k and ↑ both match Up; h/← match Left; l/→ match Right |

## Completion Criteria

- [x] Model initializes, scans vault, starts in browse mode
- [x] 5 modes: browse, view, search, find, help
- [x] Key dispatcher with both vim and arrow key support (via KeyMap arrays)
- [x] Layout: TreeStyle (25%) | ViewerStyle (75%) via lipgloss
- [x] Theme colors, icons, and styles defined
- [x] App exits on `q` / `Ctrl+C`
- [x] Window resize recalculates treeWidth
- [x] All 7 tests pass
- [x] `go vet ./...` exits 0

## Verification Evidence

- `go build ./...` exits 0
- `go test ./...` — 18/18 tests pass (7 new M2 tests + 11 from M1)
- `go vet ./...` exits 0
- Files created: `model.go`, `model_test.go`, `keys.go`, `keys_test.go`, `theme.go`
- `main.go` wired up with `tea.NewProgram`
