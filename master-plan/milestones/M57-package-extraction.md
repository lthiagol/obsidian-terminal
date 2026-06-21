# M57 — Package Structure Extraction

**Status:** ⏸ deferred (design phase — reactivation criteria below)  
**Phase:** 99 — Future (Low Priority)  
**Priority:** 🔵 Low  
**Decision:** D-6 in [PHASE-12-EXECUTION-PLAN.md](../PHASE-12-EXECUTION-PLAN.md)

## Goal

Move domain logic from `package main` into `internal/*` packages for clearer boundaries, testability, and reduced compile unit size. Keep the Bubble Tea `Model` in `main`.

## Problem statement

~13k LOC in `package main`. Vault, config, session, yamlmini are reusable domain code trapped in the binary package. Hard to reuse vault/tree/config without importing main. Large compile unit; harder for agents to navigate.

## Reactivation criteria

Execute M57 when **any** of:
- Second contributor needs to import vault logic
- M52 complete and model.go still hard to navigate (post-M59 check)
- Preparing public library API (unlikely for this CLI)
- Test coverage requires unit-testing domain logic in isolation (without Bubble Tea model)

**As of 2026-06-21:** Criteria not met. M59 (handlers split) is the higher-impact maintainability win at lower risk. Defer M57 until the above triggers fire.

## Out of scope

- Moving the Bubble Tea `Model` struct out of `main`
- Moving handler files (`handlers_*.go`, `*_handler.go`) — these are TUI-specific
- Moving render files (`render_layout.go`, `viewer.go`, `tree.go`, etc.) — these are TUI-specific
- Moving `theme.go` (lipgloss styles are TUI-specific)
- Moving `keys.go`, `mouse.go`, `help.go`, `statusbar.go`, `toast.go` (TUI-specific)
- Moving `internal/markdown`, `internal/search`, `internal/ansiext` (already extracted)
- Adding new functionality — pure relocation + API decoupling
- New dependencies

## Dependencies

| Relation | Milestone / artifact |
|----------|----------------------|
| **Blocked by** | M59 (handlers.go split — cleaner import graph after handlers are in focused files) |
| **Blocks** | nothing (optional refactor) |
| **Parallel-safe with** | nothing — large mechanical refactor; do alone |

---

## Design (approved for execution — detailed 2026-06-21)

### Proposed package map

| Package | Sources (current `main` files) | Exports | Deps on other new packages |
|---------|-------------------------------|---------|---------------------------|
| `internal/yamlmini` | `yamlmini.go` | `ScanYAML`, `ParseNestedMap`, `ParseFlatMap` | none |
| `internal/config` | `config.go` + `parseHexColor` + `paletteFromCustom` (from `theme.go`) | `Config`, `Profile`, `CustomTheme`, `LoadConfig`, `ValidateConfig`, `DefaultConfig` | `internal/yamlmini` |
| `internal/vault` | `vault.go` + `wikilink.go` | `VaultEntry`, `VaultNote`, `VaultIndexes`, `ScanVault`, `LoadNote`, `ResolveWikiLink`, `ExtractSection` | `internal/yamlmini` (for frontmatter parsing) |
| `internal/session` | `session.go` | `SessionState`, `SaveSession`, `RestoreSession` | `internal/vault` (uses `VaultEntry`) |

**Stays in `main`:** `Model`, `NewModel`, `Init`, `Update`, `View`, all handlers, all render files, `tree.go`, `viewer.go`, `viewport.go`, `theme.go` (lipgloss), `keys.go`, `mouse.go`, `help.go`, `statusbar.go`, `toast.go`, `command_palette.go`, `profile_picker.go`, `preview.go`, `textinput.go`.

### Circular dependency resolution: `config.go` → `parseHexColor` → `theme.go`

**Problem:** `config.go` calls `parseHexColor` in `ValidateConfig` (custom theme hex validation). If `config.go` moves to `internal/config`, it cannot import `main`'s `parseHexColor`.

**Resolution:** Move `parseHexColor` + `paletteFromCustom` into `internal/config` alongside the `CustomTheme` type. These are pure validation/conversion functions — they belong with the config types, not with lipgloss style building in `theme.go`.

**Impact on `theme.go`:** `theme.go` currently calls `paletteFromCustom` to build the palette from custom theme overrides. After the move, `theme.go` imports `internal/config` and calls `config.PaletteFromCustom(ct, base)`. This is fine — `main` can import `internal/config`.

**Impact on `model.go`:** `model.go` calls `parseHexColor` in `NewModel` (for custom theme setup). After the move, it calls `config.ParseHexColor(s)`.

### Cross-package type references (verified 2026-06-21)

| Symbol | Used by (non-test files) | Location after extraction |
|--------|-------------------------|--------------------------|
| `VaultEntry` | tree.go, viewer.go, vault_rescan.go, session.go, model.go, mouse.go | `internal/vault` |
| `VaultNote` | handlers_note.go, model.go, viewer.go | `internal/vault` |
| `VaultIndexes` | model.go | `internal/vault` |
| `ScanVault` | model.go | `internal/vault` |
| `LoadNote` | handlers_note.go, model.go, vault_rescan.go | `internal/vault` |
| `ResolveWikiLink` | handlers_note.go (applyNote embed resolver), handlers_view.go (wiki-link Enter) | `internal/vault` (from `wikilink.go`) |
| `ExtractSection` | handlers_note.go (embed resolver) | `internal/vault` (from `vault.go`) |
| `allPaths` | model.go, vault_rescan.go | `internal/vault` (unexported → `AllPaths` exported) |
| `Config`, `Profile`, `CustomTheme` | model.go, handlers_note.go, profile_handler.go, theme.go, config_test.go | `internal/config` |
| `LoadConfig`, `ValidateConfig`, `DefaultConfig` | main.go, model.go | `internal/config` |
| `parseHexColor` | config.go, theme.go, model.go | `internal/config` (→ `ParseHexColor`) |
| `paletteFromCustom` | theme.go | `internal/config` (→ `PaletteFromCustom`) |
| `SessionState` | model.go | `internal/session` |
| `saveSession`, `restoreSession` | model.go (quit flow, startup) | `internal/session` (→ `SaveSession`, `RestoreSession`) |
| `scanYAML`, `parseNestedMap`, `parseFlatMap` | config.go, vault.go | `internal/yamlmini` (→ `ScanYAML`, `ParseNestedMap`, `ParseFlatMap`) |

### Session/Model decoupling (required refactor)

**Problem:** `saveSession(m Model)` reads `m.config.VaultPath`, `m.fileTree.Items()`, `m.fileTree.SelectedEntry()`. `restoreSession(m *Model)` writes to `m.fileTree`. These take the Bubble Tea `Model` directly, which stays in `main`. Can't move to `internal/session` with this API.

**Resolution:** Decouple by introducing a data-only struct:

```go
// In internal/session
type SessionData struct {
    VaultPath   string
    Expanded    []string  // expanded directory paths
    CursorPath  string
}

// SaveSession writes session state to disk.
func SaveSession(data SessionData) error

// RestoreSession reads session state from disk.
// Returns SessionData and nil if no session file exists.
func RestoreSession(vaultPath string) (SessionData, error)
```

**Call site changes in `main`:**

```go
// Before (saveSession in model.go quit flow):
saveSession(m)

// After:
data := session.SessionData{
    VaultPath:  m.config.VaultPath,
    Expanded:   m.collectExpandedDirs(),  // new helper on Model
    CursorPath: m.cursorPath(),           // new helper on Model
}
if err := session.SaveSession(data); err != nil {
    m.addToast("Could not save session: "+err.Error(), ToastError)
}

// Before (restoreSession in NewModel):
restoreSession(m)

// After:
data, err := session.RestoreSession(cfg.VaultPath)
if err == nil {
    m.applySessionData(data)  // new helper on Model: expand dirs, set cursor
}
```

**New helpers on `Model` (stay in `main`):**
- `func (m Model) collectExpandedDirs() []string` — reads `m.fileTree.Items()`
- `func (m Model) cursorPath() string` — reads `m.fileTree.SelectedEntry()`
- `func (m *Model) applySessionData(data session.SessionData)` — expands dirs, sets cursor

### Export capitalization changes

Moving to `internal/*` requires unexported functions to become exported:

| File | Old name (unexported) | New name (exported) | Package |
|------|----------------------|---------------------|---------|
| `yamlmini.go` | `scanYAML` | `ScanYAML` | `internal/yamlmini` |
| `yamlmini.go` | `parseNestedMap` | `ParseNestedMap` | `internal/yamlmini` |
| `yamlmini.go` | `parseFlatMap` | `ParseFlatMap` | `internal/yamlmini` |
| `vault.go` | `allPaths` | `AllPaths` | `internal/vault` |
| `vault.go` | `extractWikiLinkTargets` | `ExtractWikiLinkTargets` | `internal/vault` |
| `vault.go` | `extractSection` | `ExtractSection` | `internal/vault` |
| `vault.go` | `normalizeWikiLinkTarget` | `NormalizeWikiLinkTarget` | `internal/vault` |
| `vault.go` | `parseFrontmatter` | `ParseFrontmatter` | `internal/vault` |
| `vault.go` | `stripFrontmatter` | `StripFrontmatter` | `internal/vault` |
| `vault.go` | `findFrontmatterBounds` | stays unexported | `internal/vault` (internal helper) |
| `session.go` | `saveSession` | `SaveSession` | `internal/session` |
| `session.go` | `restoreSession` | `RestoreSession` | `internal/session` |
| `session.go` | `stateFilePath` | stays unexported | `internal/session` (internal helper) |
| `theme.go` | `parseHexColor` | `ParseHexColor` | `internal/config` |
| `theme.go` | `paletteFromCustom` | `PaletteFromCustom` | `internal/config` |

### Key decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Safe extraction order: yamlmini → config → vault → session | Yes | yamlmini has zero deps; config depends on yamlmini; vault depends on yamlmini; session depends on vault |
| Move `parseHexColor` + `paletteFromCustom` to `internal/config` | Yes | They're pure validation/conversion; belong with config types, not lipgloss style building |
| Decouple session from `Model` via `SessionData` struct | Yes | Can't move to internal package if it takes `Model` (which stays in main); data-only struct is clean boundary |
| Export previously-unexported functions | Yes | Required for cross-package calls; `internal/*` packages only export capitalized symbols |
| Keep `wikilink.go` with `vault.go` in `internal/vault` | Yes | `ResolveWikiLink` takes `*VaultEntry` — tightly coupled to vault types |
| Keep `findFrontmatterBounds` unexported in `internal/vault` | Yes | Only used by `ParseFrontmatter` in the same package — no need to export |
| Keep `stateFilePath` unexported in `internal/session` | Yes | Only used by `SaveSession`/`RestoreSession` in the same package |

---

## Work packages

> **Rule:** One WP = one commit. Run `make test && make vet` after each WP. Strict order: WP1 → WP2 → WP3 → WP4. Do not parallelize.

### WP1 — Extract `internal/yamlmini` (3h)

**Functions to move:** `scanYAML`, `parseNestedMap`, `parseFlatMap`, `findKeyColon`, `stripInlineComment`, `stripQuotes`, `parseInlineArray`, `splitArrayItems`, `yamlPair` type.

**Steps:**
1. Create directory: `mkdir -p internal/yamlmini`
2. Create `internal/yamlmini/yamlmini.go` with `package yamlmini` header
3. Move all functions from `yamlmini.go` (root) to `internal/yamlmini/yamlmini.go`
4. Capitalize exported functions: `scanYAML` → `ScanYAML`, `parseNestedMap` → `ParseNestedMap`, `parseFlatMap` → `ParseFlatMap`
5. Keep `findKeyColon`, `stripInlineComment`, `stripQuotes`, `parseInlineArray`, `splitArrayItems`, `yamlPair` unexported (internal helpers)
6. Delete root `yamlmini.go`
7. Move `yamlmini_test.go` → `internal/yamlmini/yamlmini_test.go`; change `package main` → `package yamlmini`; update call sites to use capitalized names
8. Update callers in root package:
   - `config.go`: `scanYAML(...)` → `yamlmini.ScanYAML(...)`, `parseNestedMap(...)` → `yamlmini.ParseNestedMap(...)`, `parseFlatMap(...)` → `yamlmini.ParseFlatMap(...)`. Add import `"github.com/lthiagol/obsidian-terminal/internal/yamlmini"`
   - `vault.go`: `scanYAML(...)` → `yamlmini.ScanYAML(...)` in `parseFrontmatter`. Add import.
9. Run `goimports -w internal/yamlmini/yamlmini.go internal/yamlmini/yamlmini_test.go config.go vault.go`
10. Verify `make build` succeeds

**Import budget:**
- `internal/yamlmini/yamlmini.go`: none (pure stdlib: `strings`, `strconv` — already imported)
- `config.go`: add `"github.com/lthiagol/obsidian-terminal/internal/yamlmini"`
- `vault.go`: add `"github.com/lthiagol/obsidian-terminal/internal/yamlmini"`

**Verification:**
- [ ] `test -f internal/yamlmini/yamlmini.go && ! test -f yamlmini.go` passes
- [ ] `internal/yamlmini/yamlmini_test.go` passes (19 tests)
- [ ] `config_test.go` passes (17 tests — config parsing uses yamlmini)
- [ ] `vault_test.go` passes (14 tests — frontmatter parsing uses yamlmini)
- [ ] `make test && make vet` pass

---

### WP2 — Extract `internal/config` (4h)

**Types/functions to move from `config.go`:** `Profile`, `CustomTheme`, `Config`, `DefaultConfig`, `LoadConfig`, `parseConfigYAML`, `configPathOrDefault`, `ValidLineSpacing`, `ValidateConfig`, `stringInSlice`, `isValidDateFormat`, `validateCustomThemeColors`.

**Functions to move from `theme.go`:** `parseHexColor`, `paletteFromCustom`.

**Steps:**
1. Create directory: `mkdir -p internal/config`
2. Create `internal/config/config.go` with `package config` header
3. Move all types and functions from `config.go` (root) to `internal/config/config.go`
4. Move `parseHexColor` and `paletteFromCustom` from `theme.go` to `internal/config/config.go`
5. Capitalize: `parseHexColor` → `ParseHexColor`, `paletteFromCustom` → `PaletteFromCustom`, `parseConfigYAML` → stays unexported (internal helper), `configPathOrDefault` → `ConfigPathOrDefault`, `stringInSlice` → stays unexported, `isValidDateFormat` → stays unexported, `validateCustomThemeColors` → stays unexported
6. Add import to `internal/config/config.go`: `"github.com/lthiagol/obsidian-terminal/internal/yamlmini"` (for `ScanYAML`, `ParseNestedMap`, `ParseFlatMap`)
7. `internal/config/config.go` also needs `"github.com/charmbracelet/lipgloss"` for `ParseHexColor` return type `lipgloss.Color` and `PaletteFromCustom` return type `Palette` — **wait**, `Palette` is defined in `theme.go` (main package). This is a problem.
   - **Resolution:** Move `Palette` type + `themeData` + `lookupPalette` + `newDarkPalette` etc. to `internal/config` as well? No — that's too much. Instead:
   - `ParseHexColor` returns `(lipgloss.Color, error)` — only depends on lipgloss, fine.
   - `PaletteFromCustom` takes `CustomTheme` and `base Palette` — but `Palette` is in `theme.go` (main). **Problem.**
   - **Better resolution:** Move `ParseHexColor` to `internal/config` (returns `lipgloss.Color`, no `Palette` dependency). Keep `paletteFromCustom` in `theme.go` (main) — it builds a `Palette` which is a TUI concept. `paletteFromCustom` calls `config.ParseHexColor` instead of the local `parseHexColor`.
   - Update `config.go`'s `validateCustomThemeColors` to call `config.ParseHexColor` instead of `parseHexColor`.
8. Delete root `config.go`
9. Move `config_test.go` → `internal/config/config_test.go`; change `package main` → `package config`; update call sites to use capitalized names
10. Update callers in root package:
    - `main.go`: `LoadConfig(...)` → `config.LoadConfig(...)`, `DefaultConfig()` → `config.DefaultConfig()`. Add import.
    - `model.go`: `DefaultConfig()` → `config.DefaultConfig()`, `ValidateConfig(...)` → `config.ValidateConfig(...)`, `parseHexColor(...)` → `config.ParseHexColor(...)`. Add import.
    - `profile_handler.go`: uses `Config`, `Profile` types → `config.Config`, `config.Profile`. Add import.
    - `theme.go`: `paletteFromCustom(...)` calls `parseHexColor` → `config.ParseHexColor(...)`. Add import. Also needs `config.CustomTheme` type reference.
    - `config_test.go` moved — root tests that reference config types need updating
11. Run `goimports -w` on all changed files
12. Verify `make build` succeeds

**Import budget:**
- `internal/config/config.go`: `os`, `path/filepath`, `strings`, `time`, `github.com/charmbracelet/lipgloss`, `github.com/lthiagol/obsidian-terminal/internal/yamlmini`
- `main.go`: add `"github.com/lthiagol/obsidian-terminal/internal/config"`
- `model.go`: add `"github.com/lthiagol/obsidian-terminal/internal/config"`
- `profile_handler.go`: add `"github.com/lthiagol/obsidian-terminal/internal/config"`
- `theme.go`: add `"github.com/lthiagol/obsidian-terminal/internal/config"`

**Verification:**
- [ ] `test -f internal/config/config.go && ! test -f config.go` passes
- [ ] `internal/config/config_test.go` passes (17 tests)
- [ ] `custom_theme_test.go` passes (10 tests — uses `CustomTheme`, `ParseHexColor` via `paletteFromCustom`)
- [ ] `profiles_test.go` passes (7 tests — uses `Config`, `Profile`)
- [ ] `make test && make vet` pass

---

### WP3 — Extract `internal/vault` (5h)

**This is the largest WP — 9 non-test files reference vault types.**

**Types/functions to move from `vault.go`:** `VaultIndexes`, `VaultEntry`, `VaultNote`, `frontmatterData`, `ScanVault`, `buildTree`, `sortVaultEntries`, `LoadNote`, `findFrontmatterBounds`, `parseFrontmatter`, `stripFrontmatter`, `allPaths`, `collectPaths`, `extractTagsFromFrontmatter`, `extractWikiLinkTargets`, `normalizeWikiLinkTarget`, `extractSection`, `wikiLinkTargetRe`.

**Types/functions to move from `wikilink.go`:** `ResolveWikiLink`, `findAlias`, `findBasename`.

**Steps:**
1. Create directory: `mkdir -p internal/vault`
2. Create `internal/vault/vault.go` with `package vault` header
3. Move all types and functions from `vault.go` (root) to `internal/vault/vault.go`
4. Create `internal/vault/wikilink.go` with `package vault` header
5. Move `ResolveWikiLink`, `findAlias`, `findBasename` from `wikilink.go` (root) to `internal/vault/wikilink.go`
6. Capitalize: `allPaths` → `AllPaths`, `extractWikiLinkTargets` → `ExtractWikiLinkTargets`, `extractSection` → `ExtractSection`, `normalizeWikiLinkTarget` → `NormalizeWikiLinkTarget`, `parseFrontmatter` → `ParseFrontmatter`, `stripFrontmatter` → `StripFrontmatter`
7. Keep unexported: `buildTree`, `sortVaultEntries`, `findFrontmatterBounds`, `collectPaths`, `extractTagsFromFrontmatter`, `frontmatterData`, `wikiLinkTargetRe`, `findAlias`, `findBasename`
8. Add import to `internal/vault/vault.go`: `"github.com/lthiagol/obsidian-terminal/internal/yamlmini"` (for `ScanYAML` in `ParseFrontmatter`)
9. Delete root `vault.go` and `wikilink.go`
10. Move `vault_test.go` → `internal/vault/vault_test.go`; change `package main` → `package vault`; update call sites
11. Move `wikilink_test.go` → `internal/vault/wikilink_test.go`; change `package main` → `package vault`
12. Update callers in root package (9 non-test files — discovery command below):
    - `model.go`: `ScanVault(...)` → `vault.ScanVault(...)`, `LoadNote(...)` → `vault.LoadNote(...)`, `allPaths(...)` → `vault.AllPaths(...)`, `VaultEntry` → `vault.VaultEntry`, `VaultNote` → `vault.VaultNote`, `VaultIndexes` → `vault.VaultIndexes`. Add import.
    - `handlers_note.go`: `LoadNote(...)` → `vault.LoadNote(...)`, `ResolveWikiLink(...)` → `vault.ResolveWikiLink(...)`, `extractSection(...)` → `vault.ExtractSection(...)`, `VaultNote` → `vault.VaultNote`. Add import.
    - `handlers_view.go`: `ResolveWikiLink(...)` → `vault.ResolveWikiLink(...)`. Add import.
    - `vault_rescan.go`: `VaultEntry` → `vault.VaultEntry`, `ScanVault(...)` → `vault.ScanVault(...)`. Add import.
    - `tree.go`: `VaultEntry` → `vault.VaultEntry`. Add import.
    - `viewer.go`: `VaultNote` → `vault.VaultNote` (if referenced). Add import.
    - `session.go`: `VaultEntry` → `vault.VaultEntry` (if referenced — but session will be extracted in WP4, so this may be temporary). Add import.
    - `mouse.go`: `VaultEntry` → `vault.VaultEntry` (if referenced). Add import.
13. Run `goimports -w` on all changed files
14. Verify `make build` succeeds — this is the critical step that catches all missed references

**Discovery command for executing agent:**
```bash
# Before starting WP3, run this to find all files that reference vault types:
rg -l '\bVaultEntry\b|\bVaultNote\b|\bVaultIndexes\b|\bScanVault\b|\bLoadNote\b|\bResolveWikiLink\b|\ballPaths\b|\bextractSection\b' --glob '*.go' | grep -v '_test.go' | grep -v 'vault.go' | grep -v 'wikilink.go'
```

**Import budget:**
- `internal/vault/vault.go`: `os`, `path/filepath`, `regexp`, `sort`, `strings`, `github.com/lthiagol/obsidian-terminal/internal/yamlmini`
- `internal/vault/wikilink.go`: `strings`, `path/filepath` (already imported)
- All 9 caller files in root: add `"github.com/lthiagol/obsidian-terminal/internal/vault"`

**Verification:**
- [ ] `test -f internal/vault/vault.go && ! test -f vault.go && ! test -f wikilink.go` passes
- [ ] `internal/vault/vault_test.go` passes (14 tests)
- [ ] `internal/vault/wikilink_test.go` passes (wikilink tests)
- [ ] `tree_test.go` passes (11 tests — uses `VaultEntry`)
- [ ] `viewer_test.go` passes (12 tests — may use `VaultNote`)
- [ ] `model_test.go` passes (23 tests — uses `ScanVault`, `LoadNote`)
- [ ] `embeds_test.go` passes (8 tests — uses `ResolveWikiLink`)
- [ ] `make test && make vet` pass

---

### WP4 — Extract `internal/session` (3h)

**Types/functions to move from `session.go`:** `SessionState`, `SaveSession` (was `saveSession`), `RestoreSession` (was `restoreSession`), `stateFilePath` (stays unexported), `sessionVersion` (stays unexported).

**API refactor required:** Decouple from `Model` — see "Session/Model decoupling" in Design section.

**Steps:**
1. Create directory: `mkdir -p internal/session`
2. Create `internal/session/session.go` with `package session` header
3. Define `SessionData` struct (replaces `SessionState` — or keep `SessionState` name if preferred):
   ```go
   type SessionState struct {
       VaultPath  string   `json:"vault_path"`
       Version    int      `json:"version"`
       Expanded   []string `json:"expanded"`
       CursorPath string   `json:"cursor_path"`
   }
   ```
4. Rewrite `SaveSession` to take `SessionState` (data only, no `Model`):
   ```go
   func SaveSession(s SessionState) error
   ```
5. Rewrite `RestoreSession` to return `(SessionState, error)`:
   ```go
   func RestoreSession(vaultPath string) (SessionState, error)
   ```
6. Move `stateFilePath`, `sessionVersion` to `internal/session/session.go` (keep unexported)
7. Delete root `session.go`
8. Move `session_test.go` → `internal/session/session_test.go`; change `package main` → `package session`; update tests to use new API (construct `SessionState` directly instead of via `Model`)
9. Add new helpers to `Model` in `model.go` (or a new `session_handler.go` in root):
   ```go
   func (m Model) collectExpandedDirs() []string {
       var expanded []string
       for _, item := range m.fileTree.Items() {
           if item.entry.IsDir && item.expanded {
               expanded = append(expanded, item.entry.Path)
           }
       }
       return expanded
   }

   func (m Model) cursorPath() string {
       entry := m.fileTree.SelectedEntry()
       if entry != nil {
           return entry.Path
       }
       return ""
   }

   func (m *Model) applySessionData(data session.SessionState) {
       for _, path := range data.Expanded {
           m.fileTree.ExpandPath(path)  // may need new FileTree method
       }
       if data.CursorPath != "" {
           m.fileTree.SetCursorPath(data.CursorPath)  // may need new FileTree method
       }
   }
   ```
10. Update call sites in `model.go`:
    - Quit flow: `saveSession(m)` → construct `session.SessionState{...}` + `session.SaveSession(data)`
    - Startup (`NewModel`): `restoreSession(m)` → `data, err := session.RestoreSession(cfg.VaultPath); if err == nil { m.applySessionData(data) }`
11. **Check `FileTree` for `ExpandPath` / `SetCursorPath` methods** — if they don't exist, add them to `tree.go` (or use existing methods like `expand`/`toggleExpand` in a loop). Run: `grep -n 'func (t \*FileTree)' tree.go` to see available methods.
12. Run `goimports -w` on all changed files
13. Verify `make build` succeeds

**Import budget:**
- `internal/session/session.go`: `os`, `path/filepath`, `encoding/json`
- `model.go` (or new `session_handler.go`): add `"github.com/lthiagol/obsidian-terminal/internal/session"`, `"github.com/lthiagol/obsidian-terminal/internal/vault"` (for `VaultEntry` in `collectExpandedDirs` — but `FileTree` stores `VaultEntry` internally, so the helper reads `item.entry.Path` which is a string; the `vault` import may not be needed if `FileTree` exposes the path as a string)

**Verification:**
- [ ] `test -f internal/session/session.go && ! test -f session.go` passes
- [ ] `internal/session/session_test.go` passes (session roundtrip, vault mismatch, corrupted file)
- [ ] `session_test.go` (root, if any remaining) passes or is fully moved
- [ ] `model_test.go` passes (session restore in `NewModel`)
- [ ] `model_integration_test.go` passes (session restore integration)
- [ ] `make test && make vet` pass

---

### WP5 — Verify + update docs (1h)

**Steps:**
1. Verify package structure: `go list ./...` shows all packages:
   - `github.com/lthiagol/obsidian-terminal` (main)
   - `github.com/lthiagol/obsidian-terminal/internal/ansiext`
   - `github.com/lthiagol/obsidian-terminal/internal/config`
   - `github.com/lthiagol/obsidian-terminal/internal/markdown`
   - `github.com/lthiagol/obsidian-terminal/internal/search`
   - `github.com/lthiagol/obsidian-terminal/internal/session`
   - `github.com/lthiagol/obsidian-terminal/internal/vault`
   - `github.com/lthiagol/obsidian-terminal/internal/yamlmini`
2. Verify no circular imports: `go build ./...` succeeds
3. Verify no `main` import from `internal` packages (main imports internal, not vice versa):
   ```bash
   rg 'github.com/lthiagol/obsidian-terminal"' internal/  # should return 0 matches
   ```
4. Update `ARCHITECTURE.md` / `ARCHITECTURE.md` module map (or defer to M61 if M61 hasn't run yet)
5. Update `STATUS.md`: M57 → ✅ done

**Verification:**
- [ ] `go list ./...` shows 7 packages (was 4)
- [ ] `go build ./...` succeeds (no circular imports)
- [ ] `rg 'github.com/lthiagol/obsidian-terminal"' internal/` returns 0 matches (no internal package imports main)
- [ ] `make test && make vet` pass
- [ ] `STATUS.md` M57 → ✅

---

## Files to modify

| File | Changes |
|------|---------|
| `internal/yamlmini/yamlmini.go` | **New** — all yamlmini functions moved + capitalized |
| `internal/yamlmini/yamlmini_test.go` | **New** — moved from root |
| `internal/config/config.go` | **New** — config types + functions + `ParseHexColor` moved |
| `internal/config/config_test.go` | **New** — moved from root |
| `internal/vault/vault.go` | **New** — vault types + functions moved + capitalized |
| `internal/vault/wikilink.go` | **New** — wikilink functions moved |
| `internal/vault/vault_test.go` | **New** — moved from root |
| `internal/vault/wikilink_test.go` | **New** — moved from root |
| `internal/session/session.go` | **New** — session types + functions, decoupled from Model |
| `internal/session/session_test.go` | **New** — moved from root, updated for new API |
| `yamlmini.go` | **Deleted** |
| `config.go` | **Deleted** |
| `vault.go` | **Deleted** |
| `wikilink.go` | **Deleted** |
| `session.go` | **Deleted** |
| `yamlmini_test.go` | **Deleted** (moved) |
| `config_test.go` | **Deleted** (moved) |
| `vault_test.go` | **Deleted** (moved) |
| `wikilink_test.go` | **Deleted** (moved) |
| `session_test.go` | **Deleted** (moved) |
| `theme.go` | Remove `parseHexColor` + `paletteFromCustom`; add `internal/config` import |
| `main.go` | Add `internal/config` import; update `LoadConfig`/`DefaultConfig` calls |
| `model.go` | Add imports; update all vault/config/session references; add `collectExpandedDirs`/`cursorPath`/`applySessionData` helpers |
| `handlers_note.go` | Add `internal/vault` import; update `LoadNote`/`ResolveWikiLink`/`ExtractSection` calls |
| `handlers_view.go` | Add `internal/vault` import; update `ResolveWikiLink` call |
| `vault_rescan.go` | Add `internal/vault` import; update `ScanVault`/`VaultEntry` references |
| `tree.go` | Add `internal/vault` import; update `VaultEntry` references |
| `viewer.go` | Add `internal/vault` import (if needed) |
| `mouse.go` | Add `internal/vault` import (if needed) |
| `profile_handler.go` | Add `internal/config` import; update `Config`/`Profile` references |
| `STATUS.md` | M57 → ✅ |
| `ARCHITECTURE.md` / `ARCHITECTURE.md` | Module map: add `internal/config`, `internal/vault`, `internal/session`, `internal/yamlmini` |

## Test plan

| ID | Scenario | Type | WP |
|----|----------|------|-----|
| T1 | `internal/yamlmini` tests pass (19 tests) | unit | WP1 |
| T2 | `internal/config` tests pass (17 tests) | unit | WP2 |
| T3 | `internal/vault` tests pass (14 + wikilink tests) | unit | WP3 |
| T4 | `internal/session` tests pass (session roundtrip) | unit | WP4 |
| T5 | `custom_theme_test.go` passes (`ParseHexColor` via `paletteFromCustom`) | regression | WP2 |
| T6 | `profiles_test.go` passes (`Config`/`Profile` types) | regression | WP2 |
| T7 | `tree_test.go` passes (`VaultEntry`) | regression | WP3 |
| T8 | `embeds_test.go` passes (`ResolveWikiLink`) | regression | WP3 |
| T9 | `model_integration_test.go` passes (session restore) | regression | WP4 |
| T10 | Full suite passes | regression | WP5 |
| T11 | No circular imports (`go build ./...`) | build | WP5 |

## Acceptance criteria (when activated)

- [ ] WP1–WP5 complete
- [ ] `go list ./...` shows 7 packages (was 4: main, ansiext, markdown, search + 3 new: config, vault, session, yamlmini = 8 total — wait, 4 existing + 4 new = 8)
- [ ] No circular imports (`go build ./...` succeeds)
- [ ] No `internal/*` package imports `main` (grep verification)
- [ ] Bubble Tea `Model` remains in `main`
- [ ] All existing tests pass (298+)
- [ ] `make test && make vet` pass
- [ ] `STATUS.md` updated: M57 → ✅
- [ ] `ARCHITECTURE.md` / `ARCHITECTURE.md` module map updated (or deferred to M61)

## Rollback / risk

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Circular import (internal package imports main) | medium | WP extraction order (yamlmini → config → vault → session) prevents cycles; `go build ./...` in WP5 catches |
| Missed type reference in root package | high (9 files reference vault) | WP3 discovery command lists all files; `make build` after each WP catches |
| `FileTree` missing `ExpandPath`/`SetCursorPath` methods for session decoupling | medium | WP3 step 11 checks; add methods to `tree.go` if missing |
| `paletteFromCustom` can't move to `internal/config` (needs `Palette` type) | resolved | Keep `paletteFromCustom` in `theme.go`; only move `ParseHexColor` to `internal/config` |
| Large diff in WP3 (vault extraction touches 9 files) | high | WP3 is one commit — review carefully; run full test suite |
| Test files reference root package types | medium | Move test files with source files; update package declaration + imports |

**Rollback:** `git revert` the WP commit. Each WP is independent (strict order means WP2 doesn't start until WP1 is green).

## Handoff notes

**Read first:**
- This milestone file (especially the Cross-package type references table and Session/Model decoupling section)
- Run the WP3 discovery command before starting WP3 to confirm the 9-file list
- **M59 must be done first** — cleaner import graph after handlers are in focused files

**Do not:**
- Move `theme.go`, `tree.go`, `viewer.go`, or any TUI file to `internal/*` — they depend on `Model` and lipgloss
- Move `Palette` type to `internal/config` — it's a TUI concept tied to lipgloss styles
- Parallelize WPs — strict order yamlmini → config → vault → session
- Skip the session API refactor (WP4) — `SaveSession(m Model)` won't compile in `internal/session`

**When stuck:**
- If `go build` fails with "import cycle not allowed": you've created a circular dependency. Check if an `internal/*` package is importing `main` (it shouldn't). The extraction order prevents this — if it happens, you may have moved a function that has an unexpected dependency.
- If a test fails after moving: check if the test file's `package` declaration was updated and if function names are capitalized.
- If `paletteFromCustom` can't compile in `internal/config` because it needs `Palette`: **don't move it** — keep it in `theme.go` and only move `ParseHexColor`.

## Estimated total

3–5 days (3h WP1 + 4h WP2 + 5h WP3 + 3h WP4 + 1h WP5 = 16h ≈ 2-3 days focused, 5 days with review/buffer)

## Priority

⏸ deferred — execute when reactivation criteria met

## Completion log

_Fill when done:_

| Field | Value |
|-------|-------|
| Started | — |
| Completed | — |
| Tests added | 0 (pure relocation) |
| Notes | {paste `go list ./...` output; note any deviations from plan} |
