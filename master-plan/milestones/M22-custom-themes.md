# M22 — Custom Themes

**Status:** ⏳ pending

## Goal

Allow users to define custom color overrides in config, building on the palette system from M13.

## Steps

### 1. Add `custom_theme` config section

```yaml
theme: custom
custom_theme:
  accent: "#ff0000"
  background: "#1a1b26"
  success: "#9ece6a"
```

Any field not specified falls back to the dark theme defaults.

### 2. Build palette from custom config

Add `paletteFromCustom(cfg *Config) Palette` that reads `custom_theme` fields and merges with dark defaults.

### 3. Validate custom colors

Reject invalid hex codes. Fall back to dark theme defaults on bad values.

## Completion Criteria

- [ ] Custom colors definable in config
- [ ] Unset fields fall back to dark defaults
- [ ] Invalid hex codes rejected gracefully
- [ ] `make test && make vet` pass
