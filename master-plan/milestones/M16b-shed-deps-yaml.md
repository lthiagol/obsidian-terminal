# M16b — Replace YAML Dependency

**Status:** ⏳ pending

## Goal

Replace `gopkg.in/yaml.v3` with a minimal custom YAML parser. This eliminates the last external dependency beyond the Bubble Tea ecosystem.

## Motivation

The project uses YAML in only two places:
1. Config file parsing (`config.go`)
2. Frontmatter parsing (`vault.go`)

Both use simple YAML features: key-value pairs, arrays, quoted strings. A custom parser can handle these cases without the full YAML spec complexity.

## Scope

**Supported YAML features:**
- Key-value pairs: `key: value`
- Quoted strings: `key: "value"` or `key: 'value'`
- Inline arrays: `key: [a, b, c]`
- Block arrays: `key:\n  - item1\n  - item2`
- Comments: `# comment`
- Blank lines
- Windows line endings: `\r\n`

**NOT supported (out of scope):**
- Nested objects
- Multi-line strings (|, >)
- Anchors and aliases
- Type coercion (all values are strings)
- Complex keys

## Implementation Plan

### 1. Create `yamlmini.go`

Shared mini YAML parser:

```go
type yamlKeyValue struct {
    Key   string
    Value string
    Items []string  // for arrays
}

func scanYAML(data []byte, fn func(key, value string, items []string))
func stripQuotes(s string) string
func parseInlineArray(s string) []string
```

**scanYAML algorithm:**
1. Split input into lines
2. For each line:
   - Skip blank lines and comments
   - If line matches `key: value`, extract and call fn
   - If line matches `key: [a, b]`, parse inline array
   - If line matches `key:`, start collecting block array items
3. Block array: collect subsequent `- item` lines until non-array line

### 2. Update `config.go`

Replace yaml.Unmarshal with custom parser:

```go
func parseConfigYAML(data []byte, cfg *Config) error {
    scanYAML(data, func(key, value string, items []string) {
        switch key {
        case "vault_path":
            cfg.VaultPath = value
        case "theme":
            cfg.Theme = value
        case "default_keys":
            cfg.DefaultKeys = value
        case "skip_dirs":
            if len(items) > 0 {
                cfg.SkipDirs = items
            } else if value != "" {
                cfg.SkipDirs = []string{value}
            }
        }
    })
    return nil
}
```

Remove yaml struct tags from Config struct.

### 3. Update `vault.go`

Replace yaml.Unmarshal in parseFrontmatter:

```go
func parseFrontmatter(content string) (frontmatterData, string) {
    var fm frontmatterData
    yamlStart, yamlEnd, ok := findFrontmatterBounds(content)
    if !ok {
        return fm, content
    }
    yamlBlock := content[yamlStart:yamlEnd]
    
    scanYAML([]byte(yamlBlock), func(key, value string, items []string) {
        switch key {
        case "title":
            fm.Title = value
        case "tags":
            if len(items) > 0 {
                fm.Tags = items
            } else if value != "" {
                fm.Tags = []string{value}
            }
        case "aliases":
            if len(items) > 0 {
                fm.Aliases = items
            } else if value != "" {
                fm.Aliases = []string{value}
            }
        }
    })
    
    return fm, content[yamlEnd+5:]
}
```

Remove yaml struct tags from frontmatterData struct.

### 4. Handle edge cases

**Edge cases to handle:**
- Quoted strings with colons: `title: "Note: A Story"`
- Empty values: `tags:` (no value, no array)
- Mixed arrays: `tags: [fiction, "sci-fi"]`
- Windows line endings: `\r\n`
- Comments after values: `theme: dark # comment`
- Colons in values: `url: https://example.com`

**Edge cases to reject gracefully:**
- Invalid YAML syntax: skip line, continue parsing
- Nested structures: ignore, treat as scalar
- Unknown keys: ignore silently

### 5. Remove yaml.v3 dependency

Run `go mod tidy` to remove:
- `gopkg.in/yaml.v3`

## Testing Strategy

- Unit tests for scanYAML with various inputs
- Unit tests for parseConfigYAML
- Unit tests for parseFrontmatter
- Integration tests: load real config files
- Integration tests: parse real frontmatter from test vault
- Edge case tests: quoted strings, arrays, comments, empty values

## Fallback Plan

If the custom parser fails on real-world YAML:
1. Add logging for parse failures
2. Fall back to default values for failed fields
3. Show warning toast to user
4. Consider keeping yaml.v3 as optional dependency

## Completion Criteria

- [ ] Custom YAML parser in `yamlmini.go`
- [ ] `config.go` uses custom parser, no yaml.v3 import
- [ ] `vault.go` uses custom parser, no yaml.v3 import
- [ ] All existing tests pass
- [ ] New tests cover edge cases (quoted strings, arrays, comments)
- [ ] `make test` passes
- [ ] `make vet` exits 0
- [ ] `go mod tidy` removes yaml.v3 dependency
- [ ] Manual test: load config with all field types
- [ ] Manual test: parse frontmatter with tags and aliases
