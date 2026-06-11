# M1 — Config + Vault Scanner

**Status:** ✅ done

## Goal

Load config from `~/.config/obsidian-terminal/config.yaml` using YAML, scan vault
directory into an in-memory file tree, load individual notes with frontmatter via
yaml.v3, and build the initial search index.

## Files to create

- `config.go` / `config_test.go`
- `vault.go` / `vault_test.go`

## Steps

### 1. `config.go`
- `Config` struct:
  ```go
  type Config struct {
      VaultPath   string   `yaml:"vault_path"`
      Theme       string   `yaml:"theme"`        // "dark" | "light"
      DefaultKeys string   `yaml:"default_keys"` // "vim" | "arrows"
      SkipDirs    []string `yaml:"skip_dirs"`    // extra dirs to skip
  }
  ```
- `LoadConfig(path string) (*Config, error)` — reads YAML file, returns config
- Returns error (not nil) when file is missing — we require explicit vault path
- `DefaultConfig() *Config` — returns config with defaults:
  - Theme: "dark"
  - DefaultKeys: "vim"
  - SkipDirs: [".obsidian", ".git", ".trash", "node_modules", "archive"]
- CLI `--vault` flag overrides `Config.VaultPath`
- Config file location: `~/.config/obsidian-terminal/config.yaml`

### 2. `vault.go`
- `VaultEntry` struct:
  ```go
  type VaultEntry struct {
      Name      string
      Path      string        // relative to vault root
      IsDir     bool
      IsSymlink bool
      Children  []*VaultEntry // nil for files
  }
  ```
- `VaultNote` struct:
  ```go
  type VaultNote struct {
      Path    string
      Title   string   // from frontmatter "title", falls back to filename
      Tags    []string // from frontmatter "tags"
      Aliases []string // from frontmatter "aliases"
      Body    string   // markdown without frontmatter
      RawBody string   // full file content (for content search index)
  }
  ```
- `ScanVault(root string, skipDirs []string) (*VaultEntry, error)`:
  - Recursive walk using `filepath.WalkDir`
  - Skip: dot-prefixed files/dirs, dirs in skipDirs list
  - Only include `.md` and `.markdown` files in tree
  - Detect symlinks: `info.Mode()&os.ModeSymlink != 0` → set `IsSymlink = true`
  - Sort: folders first, then files, both alphabetically (case-insensitive)
  - Build `searchIndex map[string]string` concurrently: filename → plain text for search
- `LoadNote(vaultRoot, relativePath string) (*VaultNote, error)`:
  - Read file
  - Split frontmatter: find `---` at start, find closing `---`
  - Parse frontmatter with `yaml.v3` into struct with `title`, `tags`, `aliases`
  - Extract body (everything after closing `---`)
  - Title fallback: filename without extension, with first letter capitalized
  - Handle notes without frontmatter: Body = full content, Title = filename
  - For symlinks: resolve real path via `filepath.EvalSymlinks` for reading, but keep symlink path

### 3. CLI wiring in `main.go`
```go
func main() {
    vaultFlag := flag.String("vault", "", "path to Obsidian vault")
    configPath := flag.String("config", "", "path to config file")
    flag.Parse()

    cfg, err := LoadConfig(configPathOrDefault(*configPath))
    if err != nil {
        // Config file missing is OK — user may use --vault
        cfg = DefaultConfig()
    }
    if *vaultFlag != "" {
        cfg.VaultPath = *vaultFlag
    }
    if cfg.VaultPath == "" {
        fmt.Fprintln(os.Stderr, "Error: vault path required. Use --vault flag or set vault_path in config.")
        os.Exit(1)
    }
    // ... start bubbletea with cfg
}
```

## Test Spec (8 tests)

| # | Test | File | Description |
|---|------|------|-------------|
| 1 | `TestLoadConfig_Valid` | config_test.go | Parses vault_path, theme, default_keys, skip_dirs from valid YAML |
| 2 | `TestLoadConfig_MissingFile` | config_test.go | Returns error when config file is absent |
| 3 | `TestLoadConfig_CLIOverride` | config_test.go | `--vault` flag wins over config file vault_path |
| 4 | `TestScanVault_Structure` | vault_test.go | Returns correct nested tree from test-vault fixture (all dirs/files present) |
| 5 | `TestScanVault_SkipsExcluded` | vault_test.go | Excludes .obsidian, .git, .trash, node_modules, archive, dot-prefixed |
| 6 | `TestScanVault_SortsFoldersFirst` | vault_test.go | Folders before files, both alphabetically, case-insensitive |
| 7 | `TestLoadNote_Frontmatter` | vault_test.go | Extracts title, tags, aliases from YAML frontmatter |
| 8 | `TestLoadNote_PlainMarkdown` | vault_test.go | Handles notes without frontmatter (Body = full content, Title = filename) |

## Completion Criteria

- [ ] Config loaded from YAML, respects `--vault` override, defaults applied
- [ ] Vault scanner builds correct tree with folders/files sorted
- [ ] Hidden files, skip_dirs, non-.md files all excluded
- [ ] Symlinks detected and flagged (IsSymlink = true)
- [ ] Frontmatter parsed via yaml.v3: title, tags, aliases extracted
- [ ] Search index map built during scan
- [ ] All 8 tests pass
- [ ] `go vet ./...` exits 0
