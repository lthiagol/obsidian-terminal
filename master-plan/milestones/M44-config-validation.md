# M44 — Config Validation

**Status:** ⏳ pending

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

- [ ] All config values are validated
- [ ] Invalid values produce clear error messages with valid options
- [ ] Missing required values (vault_path) produce helpful errors
- [ ] Validation errors don't crash the app (show toast or error screen)
- [ ] `make test` passes all tests (add validation tests)
- [ ] `make vet` exits 0

## Estimated Time

1 day
