# M21 — Multiple Vault Profiles

**Status:** ⏳ pending

## Goal

Allow switching between multiple vault profiles without restarting the app. Each profile has its own vault path, theme, and settings.

## Steps

### 1. Add `profiles` config section

```yaml
profiles:
  personal:
    vault_path: ~/notes/personal
    theme: catppuccin-mocha
  work:
    vault_path: ~/notes/work
    theme: dracula
```

### 2. Add profile picker on startup

If multiple profiles are configured, show a picker to select one. `--profile` flag to bypass the picker.

### 3. Add in-app profile switching

`Ctrl+P` to open profile switcher. Rescans the vault when switching.

## Completion Criteria

- [ ] Multiple profiles configurable in YAML
- [ ] Profile picker on startup (or skip with --profile)
- [ ] In-app profile switching with Ctrl+P
- [ ] `make test && make vet` pass
