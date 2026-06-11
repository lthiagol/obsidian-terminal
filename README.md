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

```bash
go install github.com/lthiagol/obsidian-terminal@latest
```

Or build from source:

```bash
git clone https://github.com/lthiagol/obsidian-terminal.git
cd obsidian-terminal
go build -o obsidian-terminal
```

## Usage

```bash
# Open a vault directly
obsidian-terminal --vault /path/to/your/vault

# With custom config
obsidian-terminal --config ~/.config/obsidian-terminal/config.yaml
```

A YAML config is optional. If omitted, defaults apply. Create `~/.config/obsidian-terminal/config.yaml`:

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
- [ ] Graph view (ASCII-based)
- [ ] Embedded search results (block embeds)
- [ ] Mouse support for tree and viewer
- [ ] Export to PDF/HTML

## Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) | v1.3.10 | TUI framework (Elm Architecture) |
| [charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss) | v1.1.0 | Terminal styling and layout |
| [charmbracelet/bubbles](https://github.com/charmbracelet/bubbles) | v1.0.0 | TUI components (viewport) |
| [charmbracelet/x/ansi](https://github.com/charmbracelet/x) | v0.11.6 | ANSI escape sequence handling |
| [charmbracelet/x/cellbuf](https://github.com/charmbracelet/x) | v0.0.15 | Terminal cell buffer |
| [charmbracelet/x/term](https://github.com/charmbracelet/x) | v0.2.2 | Terminal capabilities |
| [charmbracelet/colorprofile](https://github.com/charmbracelet/colorprofile) | v0.4.1 | Terminal color profile detection |
| [go-yaml/yaml](https://github.com/go-yaml/yaml) | v3.0.1 | YAML config parsing |
| [muesli/termenv](https://github.com/muesli/termenv) | v0.16.0 | Terminal output helpers |
| [muesli/ansi](https://github.com/muesli/ansi) | — | ANSI sequence utilities |
| [muesli/cancelreader](https://github.com/muesli/cancelreader) | v0.2.2 | Cancellable io.Reader |
| [mattn/go-runewidth](https://github.com/mattn/go-runewidth) | v0.0.19 | East Asian character width |
| [mattn/go-isatty](https://github.com/mattn/go-isatty) | v0.0.20 | TTY detection |
| [mattn/go-localereader](https://github.com/mattn/go-localereader) | v0.0.1 | Locale-aware reader |
| [rivo/uniseg](https://github.com/rivo/uniseg) | v0.4.7 | Unicode segmentation |
| [lucasb-eyer/go-colorful](https://github.com/lucasb-eyer/go-colorful) | v1.3.0 | Color manipulation |
| [clipperhouse/uax29](https://github.com/clipperhouse/uax29) | v2.5.0 | Unicode word segmentation |
| [clipperhouse/displaywidth](https://github.com/clipperhouse/displaywidth) | v0.9.0 | Display width helpers |
| [clipperhouse/stringish](https://github.com/clipperhouse/stringish) | v0.1.1 | String utilities |
| [erikgeiser/coninput](https://github.com/erikgeiser/coninput) | — | Console input handling |
| [xo/terminfo](https://github.com/xo/terminfo) | — | Terminfo database |
| [aymanbagabas/go-osc52](https://github.com/aymanbagabas/go-osc52) | v2.0.1 | OSC52 clipboard |

## License

MIT
