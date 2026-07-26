# M26 — Daily Notes + Recent Notes

**Status:** ✅ done

## Goal

Quick-open today's daily note (`Ctrl+D`) and track recently opened notes (`Ctrl+O` overlay).

## Keybindings

**Keys:**
- `Ctrl+D` — open today's daily note (global)
- `Ctrl+O` — toggle recent notes overlay (global)

See [KEYBINDINGS.md](../../KEYBINDINGS.md) for complete keybinding reference.

## Implementation Plan

### 1. Config changes (`config.go`)

Add fields: `DailyNotesDir string` (e.g. "Journal"), `DailyNotesFormat string` (default "2006-01-02").

### 2. Model fields (`model.go`)

`recentNotes []string` (most recent last, cap 50), `recentVisible bool`, `recentCursor int`.

### 3. New methods (`model.go`)

- `buildDailyNotePath() string` — format today's date with config pattern, prepend dir
- `openDailyNote()` — loads note (or shows empty if nonexistent), adds to recents
- `addRecentNote(path string)` — dedup, append, cap at 50
- `toggleRecents()` — toggle overlay, default cursor to most recent
- `openRecentNote(index int)` — load note, close overlay, add to recents
- `renderRecents() string` — styled list, newest first, cursor highlight

### 4. Handler (`handlers.go`)

`handleRecentsKey()` — Esc/`o` dismiss, j/k navigate, Enter opens.

In `Update()` (global, before mode dispatch):
```go
case tea.KeyCtrlD: m.openDailyNote(); return m, nil
case tea.KeyCtrlO: m.toggleRecents(); return m, nil
```

Check `recentVisible` before mode dispatch → route to `handleRecentsKey`.

### 5. Note-load sites

Call `addRecentNote(path)` at all note-open locations (same as M25 buildOutline sites).

### 6. Help section

Add "Daily & Recent" bindings.

### Edge cases

- Daily note doesn't exist → show empty note with "Daily: date" title
- Daily dir doesn't exist → same behavior (LoadNote fails, empty note shown)
- Recent overlay over search → search mode preserved (overlay uses prevMode pattern)
- 50+ recents → oldest evicted
- Deleted note in recents → openRecentNote catches error, removes entry

### Implementation order

1. Add config fields
2. Add recent fields to Model
3. Implement daily note + recent note methods
4. Add handleRecentsKey
5. Wire Ctrl+D/Ctrl+O globally in Update
6. Add addRecentNote at all note-open sites
7. Add help section
8. Write tests

## Design Decisions

**Recent notes persistence:** For v1, recent notes are stored in memory only (lost on restart). Future enhancement could add file-based persistence if users request it.

## Completion Criteria

- [x] `daily_notes_dir` and `daily_notes_format` config fields added
- [x] `Ctrl+D` opens today's daily note
- [x] Daily note path built from config format (default: "2006-01-02")
- [x] Missing daily note shows empty note with date title
- [x] `Ctrl+O` toggles recent notes overlay
- [x] Recent notes list shows last 50 opened notes (newest first)
- [x] Enter on recent note opens it
- [x] Recent notes updated at all note-open sites
- [x] Deleted notes handled gracefully (removed from recents)
- [x] Help text updated
- [x] KEYBINDINGS.md updated
- [x] `make test` passes
- [x] `make vet` exits 0
- [x] Manual test: daily note and recent notes work
