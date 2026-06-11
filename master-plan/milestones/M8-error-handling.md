# M8 — Error Handling + Edge Cases

**Status:** ✅ done

## Goal

Consolidate all error handling, graceful degradation, and edge case hardening
across the entire application.

## Files to modify

- All files — add error handling where missing

## Steps

### 1. Fatal errors (prevent startup)
- **Vault path doesn't exist** → `main.go`: check `os.Stat()`, exit with clear error
- **Vault path is a file, not directory** → same check
- **Vault path is unreadable** → check permissions, exit with error
- **Config file is malformed YAML** → exit with parse error + line number
- **Terminal too small at startup** (< 60×15) → wait for first WindowSizeMsg, warn if too small

### 2. Non-fatal errors (graceful degradation)
- **File read error during scan** → skip file, add to `scanErrors []string`, show count in status bar
- **Malformed frontmatter** → use defaults (title from filename, no tags/aliases), log warning
- **Note deleted between scan and open** → toast "File not found: <path>", stay in browse
- **Custom parser error** → fall back to plain text rendering, log error
- **Markdown code block never closed** → auto-close at EOF, don't hang parser
- **Wiki-link target has invalid chars** → ignore link (don't crash)
- **Search index missing entry** (race condition) → fall back to disk read for that one file

### 3. Runtime edge cases
- **Terminal resize to < 60 columns** → show "Terminal too small" overlay, tree collapses to min width
- **Terminal resize to < 15 rows** → show "Terminal too small — please resize" warning
- **Empty vault** (0 .md files) → show "No notes found in <vault_path>" centered
- **Note with empty body** → show "(empty note)" in viewer
- **Note with only frontmatter** → show "(empty note)" in viewer
- **Note with only whitespace** → show "(empty note)" in viewer
- **Note with binary content** → detect non-UTF8, show "(binary file — cannot display)"
- **Very long single line** (5000+ chars) → truncate or don't wrap (preserve layout)
- **Deeply nested callout/quote** (> 5 levels) → stop indenting at 5, show "..." for deeper
- **Concurrent vault changes during rescan** → atomic swap of vault tree reference (pointer swap)
- **Rapid key presses** → each key processed independently, no queue overflow

### 4. Shutdown
- **`q` key** → set `quitting = true`, return `tea.Quit`
- **`Ctrl+C`** → same as `q`
- **SIGTERM/SIGINT** → bubbletea handles by default (restores terminal)
- **panic() recovery** → wrap `main()` with deferred recovery, print stack trace to stderr

### 5. Status bar warnings
- `scanErrors` count shown when > 0: `"⚠ 3 scan errors"`
- File read errors accumulate and display as count
- Clicking/tabbing to a warning could show details (v2)

### 6. Logging (optional, for debugging)
- `--debug` flag enables file logging to `/tmp/obsidian-terminal.log`
- Log: startup, vault scan, file reads, errors, mode transitions
- Not visible to user — only for `--debug` mode

## Test Spec (8 tests)

| # | Test | File | Description |
|---|------|------|-------------|
| 1 | `TestVaultPath_NotExist` | main_test.go | Error message when vault path doesn't exist |
| 2 | `TestVaultPath_IsFile` | main_test.go | Error when vault path points to a file, not directory |
| 3 | `TestScanVault_ReadError` | vault_test.go | Unreadable files skipped; scanErrors tracked |
| 4 | `TestLoadNote_MalformedFrontmatter` | vault_test.go | Bad YAML → title from filename, no tags |
| 5 | `TestParseMarkdown_UnclosedCodeBlock` | markdown_test.go | Unclosed fence auto-closes at EOF |
| 6 | `TestViewer_EmptyNote` | viewer_test.go | Note with no body shows "(empty note)" |
| 7 | `TestViewer_SingleLongLine` | viewer_test.go | 5000-char line doesn't panic; renders truncated or wrapped |
| 8 | `TestModel_Quit` | model_e2e_test.go | `q` key returns tea.Quit; `Ctrl+C` returns tea.Quit |

## Completion Criteria

- [ ] All fatal startup errors produce clear messages
- [ ] All non-fatal errors degrade gracefully (no crashes)
- [ ] Edge cases handled: empty vault, empty notes, unclosed code blocks, bad frontmatter
- [ ] Clean shutdown via `q`, `Ctrl+C`, SIGTERM
- [ ] Terminal resize handled for all sizes (normal, small, tiny)
- [ ] Concurrent vault changes handled atomically
- [ ] All 8 tests pass
- [ ] `go vet ./...` exits 0
