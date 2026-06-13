# obsidian-terminal

A terminal-based TUI for browsing and reading [Obsidian](https://obsidian.md) vaults, built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Features

- **File tree navigation** — browse your vault's directory structure with vim-style keybindings
- **Markdown rendering** — view notes with syntax-highlighted headings, bold, italic, inline code, callouts, tables, checkboxes, and frontmatter
- **Wiki-link navigation** — cycle and follow `[[wikilinks]]` between notes
- **Fuzzy file search** — quickly find notes by name (`/`)
- **Full-text search** — search across all note contents (`s`)
- **Backlinks** — see which notes link to the current note (`b`)
- **Tag browser** — browse and filter notes by tags (`T`)
- **Outline / table of contents** — jump to headings within a note (`t`)
- **Pinned notes** — keep frequently used notes in a working set (`p`)
- **Daily notes** — quickly open today's daily note (`Ctrl+D`)
- **Recent notes** — revisit recently opened notes (`Ctrl+O`)
- **Command palette** — discover and run commands (`Ctrl+K`)
- **Multiple vault profiles** — switch between vaults (`P`)
- **Custom themes** — 7 built-in palettes + per-color overrides
- **Mouse support** — click tree items, scroll viewer, drag split
- **Resizable panels** — adjust tree/viewer width (`Ctrl+←`/`→`)
- **Auto-rescan** — automatically detects external vault changes every few seconds
- **Graceful degradation** — broken vault shows error screen with retry
- **Config validation** — invalid config values are auto-fixed with helpful warnings
- **Configurable** — YAML config for vault path, theme, skip directories, and keybindings

## Installation

### Homebrew (macOS / Linux)

```bash
brew install lthiagol/tap/obsidian-terminal
```

### Go install

```bash
go install github.com/lthiagol/obsidian-terminal@latest
```

### Build from source

```bash
git clone https://github.com/lthiagol/obsidian-terminal.git
cd obsidian-terminal
make build
```

Available `make` targets:

| Command | Description |
|---------|-------------|
| `make build` | Compile the binary |
| `make run` | Run directly with `go run` |
| `make test` | Run all tests |
| `make test-race` | Run tests with race detector |
| `make vet` | Run `go vet` |
| `make lint` | Run golangci-lint |
| `make fmt` | Format code |
| `make clean` | Remove built binary |
| `make install` | Install to `$GOPATH/bin` |
| `make bench` | Run benchmarks (5s default) |
| `make bench-short` | Run benchmarks (1s) |

## Usage

```bash
# Open a vault directly
obsidian-terminal --vault /path/to/your/vault

# With custom config
obsidian-terminal --config ~/.config/obsidian-terminal/config.yaml
```

A YAML config is optional. If omitted, defaults apply. The config file is loaded from:

- `$XDG_CONFIG_HOME/obsidian-terminal/config.yaml` if set
- `~/.config/obsidian-terminal/config.yaml` otherwise

Create the config file:

```yaml
vault_path: "/Users/you/notes"
theme: "dark"
default_keys: "vim"
skip_dirs:
  - .obsidian
  - .git
  - .trash
  - node_modules
  - archive
```

## Keybindings

See [KEYBINDINGS.md](KEYBINDINGS.md) for the complete keybinding reference, including planned features and conflict resolution rules.

### Quick Reference

| Key | Action |
|-----|--------|
| `j` / `↓` | Move down |
| `k` / `↑` | Move up |
| `h` / `←` | Collapse / back |
| `l` / `→` | Expand / forward |
| `g` | Jump to top |
| `G` | Jump to bottom |
| `Enter` | Open note / toggle folder |
| `Esc` | Go back / cancel |
| `/` | Fuzzy file name search |
| `s` | Full-text content search |
| `t` | Outline / table of contents |
| `b` | Backlinks panel |
| `T` | Tag browser |
| `p` | Pin / unpin note |
| `P` | Switch profile |
| `Tab` | Cycle wiki-links (in viewer) |
| `Ctrl+D` | Open daily note |
| `Ctrl+O` | Recent notes / go back |
| `Ctrl+R` | Force rescan vault |
| `Ctrl+K` | Command palette |
| `Ctrl+←` / `Ctrl+→` | Resize tree panel |
| `Ctrl+\` | Reset tree panel width |
| `Ctrl+[` / `Ctrl+]` | Cycle pinned notes |
| `[` / `]` | History back / forward |
| `?` | Help screen |
| `q` | Quit |

## Planned Features

This TUI is read-only by design — no editing, no writing to the vault.

- [ ] Image preview (sixel/kitty protocol)
- [ ] Export to PDF/HTML

## Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) | v1.3.10 | TUI framework (Elm Architecture) |
| [charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss) | v1.1.0 | Terminal styling and layout |

## License

MIT
