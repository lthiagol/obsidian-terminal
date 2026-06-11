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
go install github.com/atr0t0s/obsidian-terminal@latest
```

Or build from source:

```bash
git clone https://github.com/atr0t0s/obsidian-terminal.git
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

- [ ] Image preview (sixel/kitty protocol)
- [ ] Backlinks panel
- [ ] Tag browsing and filtering
- [ ] Note editing
- [ ] Daily notes / quick capture
- [ ] Multiple vault profiles
- [ ] Custom CSS-like themes
- [ ] Graph view (ASCII-based)
- [ ] Embedded search results (block embeds)
- [ ] Mouse support for tree and viewer
- [ ] Drag-and-drop reordering
- [ ] Export to PDF/HTML

## License

MIT
