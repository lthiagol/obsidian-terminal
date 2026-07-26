# M0 — Environment Setup + Test Infrastructure

**Status:** ✅ done

## Goal

Set up Go toolchain, initialize the project module, install dependencies, and create the test fixture vault with comprehensive markdown samples.

## Steps

### 1. Install Go
```bash
brew install go
go version  # verify >= 1.24
```

### 2. Initialize Go module
```bash
cd /Users/thiago/Code/Pessoal/obsidian-terminal-interface
go mod init github.com/atr0t0s/obsidian-terminal
```

### 3. Add dependencies (4 total)
```bash
go get github.com/charmbracelet/bubbletea
go get github.com/charmbracelet/bubbles
go get github.com/charmbracelet/lipgloss
go get gopkg.in/yaml.v3
```

### 4. Project structure (flat layout)
```
obsidian-terminal-interface/
├── main.go
├── config.go / config_test.go
├── vault.go / vault_test.go
├── model.go / model_test.go / model_e2e_test.go
├── tree.go / tree_test.go
├── markdown.go / markdown_test.go
├── viewer.go / viewer_test.go
├── search.go / search_test.go
├── keys.go / keys_test.go
├── statusbar.go
├── theme.go
├── go.mod / go.sum
└── testdata/
    └── test-vault/      # (see below)
```

### 5. Create test fixture vault

```
testdata/test-vault/
├── index.md                    # Frontmatter: title, tags [array], aliases
│                                # Simple body: heading + paragraph + list
├── readme.md                   # Inline formatting: **bold**, *italic*, `code`,
│                                # ~~strikethrough~~, ==highlight==
├── projects/
│   ├── api-design.md           # Fenced code blocks (go, yaml, bash) + h1/h2/h3
│   ├── database.md             # [[wiki-links]] to api-design and infrastructure
│   ├── infrastructure.md       # Target note for wiki-links
│   └── deep/
│       └── nested/
│           └── buried.md       # Deeply nested (3+ levels) for tree nav tests
├── notes/
│   ├── meeting.md              # Long note (40+ lines) with lists, blockquotes
│   ├── callouts.md             # All callout types: note, tip, warning, danger,
│   │                            # info, todo, question, success, bug, example
│   ├── frontmatter-test.md     # Complex YAML: tags as array, aliases, multiline
│   └── no-frontmatter.md       # Plain markdown, no --- markers at all
├── .obsidian/
│   └── app.json                # Minimal config: {}
├── .hidden-dir/
│   └── secret.md               # Should be SKIPPED by scanner (dot-prefixed dir)
├── .gitignore                  # Should be SKIPPED (dot-prefixed file)
└── readme-symlink.md -> readme.md  # Symlink (shown as-is in tree)
```

### 6. Verify
```bash
go build ./...     # should exit 0 (no source files yet, just infra check)
go test ./...      # should exit 0 (no tests yet)
```

## Test Spec

_None — infrastructure only._

## Completion Criteria

- [ ] `go version` reports >= 1.24
- [ ] `go.mod` and `go.sum` exist with all 4 dependencies
- [ ] `go build ./...` exits 0
- [ ] Test fixture vault created with all 13+ files including symlink
- [ ] `go test ./...` exits 0
