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
| `Ctrl+D` | Open daily note | Navigate to today's note |
| `Ctrl+O` | Recent notes / go back | In View mode: go back in history; otherwise: toggle recent notes |
| `Ctrl+K` | Command palette | Open command palette overlay |
| `Ctrl+←` | Shrink tree panel | Decrease tree width by 5 |
| `Ctrl+→` | Grow tree panel | Increase tree width by 5 |
| `Ctrl+\` | Reset tree width | Restore default tree width (1/4 of screen) |
| `q` / `Q` | Quit | Graceful exit |
| `?` | Toggle help | Show/hide help panel |
| `Esc` | Cancel / go back | Context-dependent |

### Browse Mode

| Key | Action | Notes |
|-----|--------|-------|
| `j` / `↓` | Move down | File tree navigation |
| `k` / `↑` | Move up | File tree navigation |
| `h` / `←` | Collapse directory | Tree navigation |
| `l` / `→` | Expand directory | Tree navigation |
| `g` | Jump to top | Move to first item |
| `G` | Jump to bottom | Move to last item |
| `PgUp` | Page up | Scroll tree up |
| `PgDn` | Page down | Scroll tree down |
| `Enter` | Open note / toggle folder | Primary action |
| `/` | Fuzzy file search | Enter search mode |
| `s` | Full-text content search | Enter find mode |
| `T` | Browse tags | Enter tag browser mode |
| `p` | Pin/unpin current note | Toggle pin on selected entry |
| `P` | Switch profile | Open profile picker (if profiles configured) |
| `Ctrl+[` | Cycle to previous pinned note | Navigate pinned notes |
| `Ctrl+]` | Cycle to next pinned note | Navigate pinned notes |
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
| `b` | Toggle backlinks panel | Show/hide backlinks overlay |
| `t` | Toggle outline/TOC | Show/hide table of contents overlay |
| `p` | Pin/unpin current note | Toggle pin on active note |
| `Ctrl+[` | Cycle to previous pinned note | Navigate pinned notes |
| `Ctrl+]` | Cycle to next pinned note | Navigate pinned notes |
| `/` | Activate in-note search | Search within current note |
| `n` / `N` | Next / previous in-note match | Cycle search results (when in-note search active) |
| `s` | Full-text content search | Enter find mode across all notes |
| `[` | Go back in history | Navigate to previous note |
| `]` | Go forward in history | Navigate to next note |
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

### Backlink Mode (overlay in View Mode)

| Key | Action | Notes |
|-----|--------|-------|
| `j` / `↓` | Move down | Backlink list navigation |
| `k` / `↑` | Move up | Backlink list navigation |
| `Enter` | Open selected backlink | Navigate to linked note |
| `Esc` / `b` | Close backlinks | Return to view mode |

### Tag Browser Mode

| Key | Action | Notes |
|-----|--------|-------|
| `j` / `↓` | Move down | Tag list navigation |
| `k` / `↑` | Move up | Tag list navigation |
| `Enter` | Filter by selected tag | Apply tag filter to file tree |
| `Esc` | Cancel | Return to browse mode |

### Help Mode

| Key | Action | Notes |
|-----|--------|-------|
| `j` / `↓` | Scroll down | Help navigation |
| `k` / `↑` | Scroll up | Help navigation |
| `Esc` | Close help | Return to previous mode |

### Command Palette Overlay

| Key | Action | Notes |
|-----|--------|-------|
| `j` / `↓` | Move down in results | Command navigation |
| `k` / `↑` | Move up in results | Command navigation |
| `Enter` | Execute selected command | |
| `Esc` | Close palette | Cancel |
| `Backspace` | Delete character | Query editing |
| *(any rune)* | Append to query | Query editing |

### Recent Notes Overlay

| Key | Action | Notes |
|-----|--------|-------|
| `j` / `↓` | Move down in list | Recent notes navigation |
| `k` / `↑` | Move up in list | Recent notes navigation |
| `Enter` | Open selected note | |
| `Esc` | Close list | Cancel |

### Outline Overlay

| Key | Action | Notes |
|-----|--------|-------|
| `j` / `↓` | Move down | Outline navigation |
| `k` / `↑` | Move up | Outline navigation |
| `Enter` | Jump to heading | Scroll viewer to selected heading |
| `Esc` / `t` | Close outline | Return to view mode |

### Profile Picker Mode

| Key | Action | Notes |
|-----|--------|-------|
| `j` / `↓` | Move down | Profile list navigation |
| `k` / `↑` | Move up | Profile list navigation |
| `Enter` | Switch to selected profile | |
| `Esc` | Cancel / quit | Returns to browse or quit if no vault loaded |

### Broken Vault State

| Key | Action | Notes |
|-----|--------|-------|
| `r` | Retry rescan | Attempt to reload vault |
| `q` / `Q` | Quit | Exit application |
| `?` | Toggle help | Show/hide help panel |

## Mouse Support

| Action | Target | Behavior |
|--------|--------|----------|
| Left click | Tree | Move cursor to clicked item |
| Left double-click | Tree | Open item (same as Enter) |
| Wheel up/down | Tree | Scroll 3 lines |
| Wheel up/down | Viewer | Scroll (viewport handles) |
| Left click | Search results | Select result |
| Left double-click | Search results | Open result |
| Wheel up/down | Help | Scroll help text |

## Reserved Keys

**Do not allocate these without updating this document:**

- All navigation keys: `j/k/h/l`, arrows, `g/G`
- Core actions: `Enter`, `Esc`, `Tab`, `/`, `s`, `?`, `q`
- System keys: `Ctrl+C`, `Ctrl+R`, `Ctrl+D`, `Ctrl+O`, `Ctrl+K`
- Tree resize: `Ctrl+←`, `Ctrl+→`, `Ctrl+\`
- Already in use: `b`, `t`, `T`, `p`, `P`, `n`, `N`, `[`, `]`, `r`
- Pin cycling: `Ctrl+[`, `Ctrl+]`

## Keybinding Conflict Resolution

When adding a new keybinding:

1. **Check this document** - Is the key already allocated?
2. **Consider mode scope** - Can the key be mode-specific?
3. **Prefer mnemonics** - `b` for backlinks, `t` for tags/TOC
4. **Use Shift variants** - If lowercase is taken, try uppercase
5. **Use Ctrl combinations** - For global actions
6. **Update this document** - Add to "Reserved" table before implementing
7. **Update help text** - Add to `help.go` renderHelp() function

## Testing Keybindings

When implementing new keybindings:

1. Add to `KeyMap` struct in `keys.go`
2. Add to `DefaultKeys()` function
3. Add handler case in appropriate `handle*Key()` function
4. Add to help text in `help.go`
5. Update this document
6. Write tests in `keys_test.go` or mode-specific test file
7. Test for conflicts with existing bindings
