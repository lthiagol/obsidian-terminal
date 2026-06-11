# M21 — Multiple Vault Profiles

**Status:** ⏳ pending

## Goal

Configure multiple vault profiles in config, switch between them with `--profile` flag and in-app picker (`P` in Browse mode).

## Keybinding

**Key:** `P` (uppercase)
**Mode:** Browse mode only
**Rationale:** Mnemonic for "Profile", uppercase to be distinct

See [KEYBINDINGS.md](../../KEYBINDINGS.md) for complete keybinding reference.

## Implementation Plan

### 1. Config changes (`config.go`)

```go
type Profile struct {
    Path     string   `yaml:"path"`
    Theme    string   `yaml:"theme"`
    SkipDirs []string `yaml:"skip_dirs"`
}
// Add to Config: Profile string; Profiles map[string]Profile
```

### 2. CLI flag (`main.go`)

Add `--profile` flag. Resolution logic (order):
1. Load config
2. Apply `--vault` to cfg.VaultPath (if set)
3. Apply `--profile` to cfg.Profile (if set)
4. If no `--vault` but profile is set: resolve profile → populate cfg.VaultPath/Theme/SkipDirs
5. If neither and profiles exist → enter profile picker mode

### 3. New mode: `ModeProfilePicker`

New file: `profile_picker.go`:
```go
type ProfilePicker struct { profiles []string; cursor int }
func NewProfilePicker(profiles map[string]Profile) ProfilePicker  // sorted names
func (pp *ProfilePicker) MoveUp/Down()
func (pp ProfilePicker) Selected() string
func (pp ProfilePicker) View() string  // centered list
```

### 4. Model changes (`model.go`)

Add `ModeProfilePicker` to Mode. Add `profilePicker ProfilePicker`, `pendingProfileSwitch string` fields.

In `NewModel`: if vault path empty and profiles exist → return model in picker mode (don't error).

### 5. Handlers (`handlers.go`)

`handleProfilePickerKey()`: Esc quits (startup) / returns to prev mode (in-app). Enter selects profile, triggers rescan. Arrow keys move cursor.

In `handleBrowseKey`: add `ProfileSwitch rune` (`P`) that opens picker from browse mode.

In `Update()`: check `pendingProfileSwitch`, apply profile settings, rescan vault.

Note: Consider using a Cmd-based approach instead of `pendingProfileSwitch` field for better Bubble Tea idiomatic style.

### 6. View() — picker rendering

Render picker full-screen centered when in `ModeProfilePicker`.

### Edge cases

- `--vault` + `--profile`: vault flag wins
- Missing profile: error + exit
- No path in profile: error
- In-app switch while viewing note: note discarded, full rescan
- No profiles defined: picker never shown, `P` no-op
- Profile with invalid theme: show warning, use default theme

### Implementation order

1. Add Profile struct + fields to config.go
2. Update main.go with --profile flag + resolution
3. Add ModeProfilePicker + picker fields to Model
4. Modify NewModel for picker mode
5. Create profile_picker.go
6. Add handleProfilePickerKey
7. Add ProfileSwitch keybinding (`P`)
8. Wire dispatch + rendering
9. Handle in-app switching
10. Write tests

## Completion Criteria

- [ ] Profile struct added to config.go
- [ ] `profiles` map in config for multiple vault profiles
- [ ] `--profile` CLI flag implemented
- [ ] Profile picker mode with centered list UI
- [ ] `P` keybinding opens picker in Browse mode
- [ ] Enter selects profile and triggers rescan
- [ ] Esc returns to previous mode
- [ ] Profile switch applies theme and skip_dirs
- [ ] `--vault` flag takes precedence over profile
- [ ] Missing/invalid profile shows error
- [ ] Help text updated
- [ ] KEYBINDINGS.md updated
- [ ] `make test` passes
- [ ] `make vet` exits 0
- [ ] Manual test: profile switching works end-to-end
