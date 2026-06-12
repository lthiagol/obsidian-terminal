# M44 — Config Validation

**Status:** ✅ done

## Goal

Validate all configuration values and provide helpful error messages when invalid values are provided.

## Issues

### Missing validation

Config values are not validated. Invalid values are silently ignored:
- `line_spacing: invalid` → silently uses default "compact"
- `theme: nonexistent` → error, but no guidance on valid values
- `vault_path: /nonexistent` → error, but no suggestion to create it
- `daily_notes_format: invalid` → silently uses default
- `skip_dirs: [".git"]` → no validation that paths are valid

### Missing helpful error messages

When config loading fails, the error messages don't guide the user:
- No list of valid theme names
- No list of valid line_spacing values
- No suggestion to check file permissions
- No example of correct config format

## Design

### Validation function

Create a `ValidateConfig(cfg *Config) []error` function that checks:
1. `vault_path` exists and is a directory
2. `theme` is one of the valid theme names (or empty for default)
3. `line_spacing` is one of: "compact", "normal", "relaxed" (or empty for default)
4. `daily_notes_format` is a valid Go time format (try parsing a test date)
5. `skip_dirs` entries are valid path components (no absolute paths, no `..`)
6. `custom_theme` colors are valid hex colors (#RRGGBB or #RGB)
7. `profiles` entries have valid `path` values

### Helpful error messages

When validation fails, provide:
- The invalid value
- The list of valid values (for enum-like fields)
- A suggestion to fix it
- An example of correct usage

Example:
```
Error: invalid theme "nonexistent" in config
Valid themes: dark, catppuccin-latte, catppuccin-frappe, catppuccin-macchiato, catppuccin-mocha, dracula, alucard
Example: theme: dark
```

## Files to modify

| File | Changes |
|------|---------|
| `config.go` | Add `ValidateConfig` function; call it after `LoadConfig` |
| `theme.go` | Export list of valid theme names |
| `main.go` | Display validation errors with helpful messages |

## Completion Criteria

- [x] All config values are validated
- [x] Invalid values produce clear error messages with valid options
- [x] Missing required values (vault_path) produce helpful errors
- [x] Validation errors don't crash the app (show toast or error screen)
- [x] `make test` passes all tests (add validation tests)
- [x] `make vet` exits 0

## Completed

2026-06-12

Added `ValidateConfig(cfg *Config) []string` in config.go that validates and auto-fixes:
- Theme (auto-fixes invalid to "dark", shows valid names)
- Line spacing (auto-fixes invalid to "compact", shows valid values)
- Daily notes format (round-trip validates, auto-fixes to "2006-01-02")
- Skip dirs (warns about invalid entries with path separators)
- Custom theme colors (warns about invalid hex colors per-field)
- Profile paths (warns about empty paths)

`ValidLineSpacing` exported for discoverability. `ValidThemeNames()` already existed.

Wired into `NewModel`: validation runs early, warnings are shown as toasts (non-blocking). Vault path errors in `NewModel` now include actionable suggestions ("directory does not exist, create it first" / "check file permissions").

19 new tests covering valid config, invalid theme, invalid spacing, invalid date format, invalid skip dirs, custom theme colors, empty profiles, and end-to-end model toast display.

## Estimated Time

1 day
