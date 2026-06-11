# M29 — Command Palette

**Status:** ⏳ pending

## Goal

Add a fuzzy command palette (`Ctrl+K`) to quickly access any action without remembering keybindings.

## Steps

### 1. Define command registry

```go
type Command struct {
    Name        string
    Description string
    Action      func(*Model)
}
```

Register all available commands: open search, toggle help, switch theme, pin note, go to outline, etc.

### 2. Add command palette UI

- `Ctrl+K` — open command palette (fuzzy search over commands)
- Type to filter commands
- Enter to execute selected command
- Esc to close

Reuse the existing fuzzy search UI pattern but search over commands instead of files.

### 3. Register key commands

About 15-20 commands initially — every action with a keybinding should also be discoverable via the palette. Show the keybinding next to the command name (e.g., `Toggle Help  (?)`).

## Completion Criteria

- [ ] `Ctrl+K` opens command palette
- [ ] Fuzzy search over commands
- [ ] Enter executes selected command
- [ ] Keybindings shown next to command names
- [ ] `make test && make vet` pass
