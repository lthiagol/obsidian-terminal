# M17 — Performance: Hot Path Optimizations

**Status:** ⏳ pending

## Goal

Optimize the hottest render/search paths to reduce per-frame and per-keystroke allocations.

## Files to modify

- `tree.go` — pre-compute lipgloss styles, cache depth prefixes
- `internal/search/search.go` — pre-compute lowercase paths in NewState
- `help.go` — cache static help text

## Steps

### 1. Pre-compute tree styles and prefix cache

In `tree.go`, define package-level styles so `View()` doesn't allocate 3-5 lipgloss styles per item per frame:

```go
var (
    treeFileStyle     = lipgloss.NewStyle().Foreground(TextSecondary)
    treeDirStyle      = lipgloss.NewStyle().Foreground(AccentSecondary)
    treeSelFileStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#000000")).Background(Accent).Bold(true)
    treeSelDirStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#000000")).Background(Accent).Bold(true)
)
```

Also add a `depthPrefix` cache on `FileTree` to avoid `strings.Repeat("  ", depth)` per item.

### 2. Pre-compute lowercase paths for FuzzySearch

In `NewState`, when mode is `Name`, store lowercase versions of all paths once:

```go
func NewState(mode Mode, paths []string, index map[string]string) State {
    s := State{mode: mode, allPaths: paths, searchIndex: index}
    if mode == Name {
        s.lowerPaths = make([]string, len(paths))
        for i, p := range paths {
            s.lowerPaths[i] = strings.ToLower(p)
        }
        s.results = FuzzySearch("", paths)
    }
    return s
}
```

Update `FuzzySearch` and `FuzzyScore` to accept/use pre-lowered paths.

### 3. Cache help text

In `help.go`, build the help text once into a package-level var:

```go
var helpLines []string

func initHelp() {
    // build groups + lines once
}
```

Call `initHelp()` once at startup. `renderHelp()` just slices `helpLines[start:end]`.

## Completion Criteria

- [ ] Pre-computed tree styles + prefix cache in tree.go
- [ ] Pre-computed lowercase paths in search State
- [ ] Help text built once, not every frame
- [ ] All 98 tests pass
- [ ] `make build && make vet` exit 0
