# Keybindings Reference

This document contains the complete keybinding reference for obsidian-terminal. All milestones that add new keybindings must update this file to avoid conflicts.

## Keybinding Allocation Rules

1. **Single-letter keys** (`a-z`, `A-Z`): Check this document before allocating
2. **Ctrl keys** (`Ctrl+X`): Prefer unused combinations, document in Global section
3. **Mode-specific keys**: Same key can have different actions in different modes
4. **Navigation keys** (`j/k/h/l`, arrows): Reserved, never allocate for features
5. **Special keys** (`Enter`, `Esc`, `Tab`): Reserved for core navigation

## Current Keybindings

### Global (All Modes)

| Key | Action | Notes |
|-----|--------|-------|
| `Ctrl+C` | Quit | Immediate exit |
| `Ctrl+R` | Force rescan vault | Rebuild file tree and indexes |
| `q` | Quit | Graceful exit |
| `Q` | Quit | Alternative quit |

### Browse Mode

| Key | Action | Notes |
|-----|--------|-------|
| `j` / `↓` | Move down | File tree navigation |
| `k` / `↑` | Move up | File tree navigation |
| `h` / `←` | Collapse directory | Tree navigation |
| `l` / `→` | Expand directory | Tree navigation |
| `g` | Jump to top | Move to first item |
| `G` | Jump to bottom | Move to last item |
| `Enter` | Open note / toggle folder | Primary action |
| `/` | Fuzzy file search | Enter search mode |
| `s` | Full-text content search | Enter find mode |
| `?` | Toggle help | Show/hide help panel |

### View Mode

| Key | Action | Notes |
|-----|--------|-------|
| `j` / `↓` | Scroll down | Viewer navigation |
| `k` / `↑` | Scroll up | Viewer navigation |
| `g` | Jump to top | Scroll to beginning |
| `G` | Jump to bottom | Scroll to end |
| `PgUp` | Page up | Half-page scroll |
| `PgDn` | Page down | Half-page scroll |
| `Tab` | Cycle wiki-links | Highlight next link |
| `Enter` | Follow selected link | Navigate to linked note |
| `h` / `Esc` | Back to browse | Exit viewer |
| `/` | Fuzzy file search | Enter search mode |
| `s` | Full-text content search | Enter find mode |
| `?` | Toggle help | Show/hide help panel |

### Search/Find Mode

| Key | Action | Notes |
|-----|--------|-------|
| `j` / `↓` | Move down in results | Result navigation |
| `k` / `↑` | Move up in results | Result navigation |
| `Enter` | Open selected result | Load note and enter view mode |
| `Esc` | Cancel search | Return to previous mode |
| `Backspace` | Delete character | Query editing |
| *(any rune)* | Append to query | Query editing |

### Help Mode

| Key | Action | Notes |
|-----|--------|-------|
| `j` / `↓` | Scroll down | Help navigation |
| `k` / `↑` | Scroll up | Help navigation |
| `Esc` | Close help | Return to previous mode |

## Planned Keybindings

### Reserved for Future Features

| Key | Feature | Mode | Milestone | Status |
|-----|---------|------|-----------|--------|
| `b` | Toggle backlinks panel | View | M19 | ⏳ pending |
| `t` | Toggle outline/TOC | View | M25 | ⏳ pending |
| `T` | Open tag browser | Browse | M20 | ⏳ pending |
| `p` | Pin/unpin current note | Browse, View | M24 | ⏳ pending |
| `P` | Open profile picker | Browse | M21 | ⏳ pending |
| `o` | Open outline (alternative) | View | M25 | 🔄 backup |
| `Ctrl+D` | Open daily note | Global | M26 | ⏳ pending |
| `Ctrl+O` | Open recent notes | Global | M26 | ⏳ pending |
| `Ctrl+K` | Open command palette | Global | M29 | ⏳ pending |
| `Ctrl+[` | Cycle to previous pin | Global | M24 | ⏳ pending |
| `Ctrl+]` | Cycle to next pin | Global | M24 | ⏳ pending |

### Mouse Support (M18)

| Action | Target | Behavior |
|--------|--------|----------|
| Left click | Tree | Move cursor to clicked item |
| Left double-click | Tree | Open item (same as Enter) |
| Wheel up/down | Tree | Scroll 3 lines |
| Wheel up/down | Viewer | Scroll (viewport handles) |
| Left click | Search results | Select result |
| Left double-click | Search results | Open result |
| Wheel up/down | Help | Scroll help text |

## Keybinding Conflict Resolution

When adding a new keybinding:

1. **Check this document** - Is the key already allocated?
2. **Consider mode scope** - Can the key be mode-specific?
3. **Prefer mnemonics** - `b` for backlinks, `t` for tags/TOC
4. **Use Shift variants** - If lowercase is taken, try uppercase
5. **Use Ctrl combinations** - For global actions
6. **Update this document** - Add to "Reserved" table before implementing
7. **Update help text** - Add to `help.go` renderHelp() function

## Reserved Keys

**Do not allocate these without updating this document:**

- All navigation keys: `j/k/h/l`, arrows, `g/G`
- Core actions: `Enter`, `Esc`, `Tab`, `/`, `s`, `?`, `q`
- System keys: `Ctrl+C`, `Ctrl+R`
- Already planned: See "Reserved for Future Features" table

## Testing Keybindings

When implementing new keybindings:

1. Add to `KeyMap` struct in `keys.go`
2. Add to `DefaultKeys()` function
3. Add handler case in appropriate `handle*Key()` function
4. Add to help text in `help.go`
5. Update this document
6. Write tests in `keys_test.go` or mode-specific test file
7. Test for conflicts with existing bindings
