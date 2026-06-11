# M26 — Daily Notes + Recent Notes

**Status:** ⏳ pending

## Goal

Quick-open daily notes by date pattern and track recently opened notes for fast re-access.

## Steps

### 1. Daily notes

Add config field for the daily note directory (e.g., `daily/`) and date format (e.g., `2006-01-02`). On `Ctrl+D`, open today's daily note. If it doesn't exist, show a toast — read-only, no file creation.

### 2. Recent notes

Track the last 10 opened notes in memory with timestamps. New mode `ModeRecent` triggered by `Ctrl+O`. Shows a list of recent notes with last-accessed time. Enter to open.

### 3. UI

- `Ctrl+D` — open today's daily note (if it exists)
- `Ctrl+O` — open recent notes list
- Recent notes list shows: `title  —  opened 2m ago`
- Both stats shown in status bar hints

## Completion Criteria

- [ ] `Ctrl+D` opens today's daily note
- [ ] Daily note directory and date format configurable
- [ ] `Ctrl+O` shows recent notes list
- [ ] Recent notes updated on every note open
- [ ] `make test && make vet` pass
