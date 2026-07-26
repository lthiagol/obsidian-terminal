# M7 — File Watcher + Polish

**Status:** ✅ done

## Goal

Add auto-refresh when vault files change (poll-based with modtime fast-path), wiki-link
resolution, and toast notifications. Wire up the complete user experience.

## Files to modify

- `model.go` — watcher tick, toast system
- `vault.go` — wiki-link resolution

## Steps

### 1. File watcher (poll-based)
- `tickMsg` — internal message type for periodic ticks
- `Init()` returns `tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg { return tickMsg{} })`
- `lastRootModTime time.Time` in Model
- On each tick:
  1. `os.Stat(vaultPath)` → get root modtime
  2. If unchanged since last check → skip (fast path)
  3. If changed → rescan vault via `ScanVault()`, update model's tree
  4. If currently viewed note was deleted (path not in new tree):
     - Show toast: `"Note was deleted: <path>"`
     - Set mode to browse, clear activeNote
  5. If currently viewed note was modified (modtime changed):
     - Reload note body, re-render viewer
  6. Update `lastRootModTime`
- Debounce: don't trigger within 500ms of last rescan (use `lastRescan time.Time`)

### 2. Wiki-link resolution
- `ResolveWikiLink(target string, vault *VaultEntry) string`:
  1. Strip `#heading-name` and `#^block-id` suffixes
  2. Strip `|display` text (already handled by ExtractWikiLinks)
  3. Exact relative path match (append `.md` if no extension):
     - `projects/api-design` → `projects/api-design.md`
  4. Walk all files; match by basename (filename without path):
     - `api-design` → `projects/api-design.md`
  5. Walk all files; case-insensitive basename match:
     - `API-Design` → `projects/api-design.md`
  6. Walk all notes' `aliases` frontmatter; case-insensitive alias match
  7. Return empty string if no match found
- On link follow (Enter with selectedLink >= 0):
  - Call ResolveWikiLink
  - If found: load target note, set activeNote, reset selectedLink
  - If not found: show toast `"Link not found: <target>"`

### 3. Toast notification system
- `Toast` struct:
  ```go
  type Toast struct {
      Message string
      Type    ToastType  // info, success, warning, error
      TTL     time.Duration
      Created time.Time
  }
  ```
- `AddToast(model *Model, message string, t ToastType)` — adds toast with 3s TTL
- On each tick: remove expired toasts
- Render toasts stacked in bottom-right corner (above status bar, within viewer panel)
- Toast styling:
  - Info: blue border + `ℹ` icon
  - Success: green border + `✔` icon
  - Warning: amber border + `⚠` icon
  - Error: red border + `✖` icon

### 4. Polish details
- Wiki-link highlighting: when selectedLink >= 0, status bar shows `"→ <target>"` in teal
- Viewer scroll progress: optional `"Line X/Y (Z%)"` in status bar
- Manual refresh: `Ctrl+R` triggers immediate vault rescan (overrides debounce)
- Smooth resize: recalculate `treeWidth` on WindowSizeMsg; clamp to reasonable min/max

## Test Spec (9 tests)

| # | Test | File | Description |
|---|------|------|-------------|
| 1 | `TestWatcher_DetectsNewFile` | model_test.go | Adding .md file triggers rescan; tree reflects new file |
| 2 | `TestWatcher_DetectsModify` | model_test.go | Modifying file content triggers note reload in viewer |
| 3 | `TestWatcher_DetectsDelete` | model_test.go | Deleting .md file removes it from tree |
| 4 | `TestWatcher_IgnoresNonMd` | model_test.go | Non-.md file changes do not trigger rescan |
| 5 | `TestWatcher_NoteDeletedWhileViewing` | model_test.go | Viewing a note that gets deleted → mode=browse + toast |
| 6 | `TestWikiLinkResolution_ExactPath` | vault_test.go | `projects/api-design` resolves to `projects/api-design.md` |
| 7 | `TestWikiLinkResolution_Basename` | vault_test.go | `api-design` resolves to `projects/api-design.md` |
| 8 | `TestWikiLinkResolution_CaseInsensitive` | vault_test.go | `API-Design` resolves to `projects/api-design.md` |
| 9 | `TestWikiLinkResolution_AliasMatch` | vault_test.go | Link target matching an alias resolves to correct note |

## Completion Criteria

- [ ] File watcher detects new, modified, deleted `.md` files
- [ ] Modtime fast-path avoids unnecessary rescans
- [ ] Deleted note returns to browse with toast
- [ ] Modified note content auto-refreshes in viewer
- [ ] Wiki-link resolution: exact path, basename, case-insensitive, aliases
- [ ] Wiki-link fragments (#heading, #^block, |display) stripped
- [ ] Toast notifications display and auto-dismiss (3s TTL)
- [ ] `Ctrl+R` manual refresh
- [ ] All 9 tests pass
- [ ] `go vet ./...` exits 0
