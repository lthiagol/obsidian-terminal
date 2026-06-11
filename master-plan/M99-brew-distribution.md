# M99 — Homebrew Distribution

**Status:** ⏳ pending

## Goal

Package obsidian-terminal as a Homebrew formula so users can install it with
`brew install lthiagol/tap/obsidian-terminal`.

## Steps

### 1. Create the tap repository

- Create a public repo: `github.com/lthiagol/homebrew-tap`
- Homebrew requires the repo name to start with `homebrew-`

### 2. Create the formula

- Write `Formula/obsidian-terminal.rb` in the tap repo
- Formula should:
  - Build from source using `go build`
  - Reference the `github.com/lthiagol/obsidian-terminal` repo
  - Include a version tag (e.g. `v0.1.0`)
  - Set a description and homepage

### 3. Tag a release

- Tag the first-version commit as a release (e.g. `v0.1.0`)
- Push the tag to GitHub
- Create a GitHub Release with pre-built binaries for macOS (amd64 + arm64)

### 4. Test the formula

```bash
brew install --build-from-source ./Formula/obsidian-terminal.rb
brew test obsidian-terminal
brew audit --strict obsidian-terminal
```

### 5. Document in README

- README already references `brew install lthiagol/tap/obsidian-terminal`
- Add a Homebrew badge once the tap is live

## Completion Criteria

- [ ] `homebrew-tap` repository created
- [ ] `obsidian-terminal.rb` formula written and committed
- [ ] Release tag created and pushed
- [ ] `brew install lthiagol/tap/obsidian-terminal` works end-to-end
- [ ] README badge added
