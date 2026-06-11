# obsidian-terminal

A terminal-based TUI for browsing and reading [Obsidian](https://obsidian.md) vaults, built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Features

- **File tree navigation** — browse your vault's directory structure with vim-style keybindings
- **Markdown rendering** — view notes with syntax-highlighted headings, bold, italic, inline code, and callouts
- **Wiki-link navigation** — cycle and follow `[[wikilinks]]` between notes
- **Fuzzy file search** — quickly find notes by name (`/`)
- **Full-text search** — search across all note contents (`s`)
- **Auto-rescan** — automatically detects external vault changes every few seconds
- **Configurable** — YAML config for vault path, theme, skip directories, and keybindings
- **Vim and Emacs keys** — configurable navigation styles

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
| `Enter` | Open note / toggle folder |
| `/` | Fuzzy file name search |
| `s` | Full-text content search |
| `Tab` | Cycle wiki-links (in viewer) |
| `Enter` | Follow selected link |
| `Ctrl+R` | Force rescan vault |
| `?` | Help screen |
| `q` | Quit |

## Planned Features

This TUI is read-only by design — no editing, no writing to the vault.

- [ ] Image preview (sixel/kitty protocol)
- [ ] Backlinks panel
- [ ] Tag browsing and filtering
- [ ] Multiple vault profiles
- [ ] Custom CSS-like themes
- [ ] Pinned notes (working set)
- [ ] Outline / table of contents
- [ ] Daily notes + recent notes
- [ ] Checkbox rendering (`- [ ]` / `- [x]`)
- [ ] Frontmatter metadata display
- [ ] Markdown table rendering
- [ ] Command palette
- [ ] Mouse support for tree and viewer
- [ ] Export to PDF/HTML

## Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) | v1.3.10 | TUI framework (Elm Architecture) |
| [charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss) | v1.1.0 | Terminal styling and layout |

## License

MIT
