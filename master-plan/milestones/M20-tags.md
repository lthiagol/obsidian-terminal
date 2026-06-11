# M20 — Tag Browsing & Filtering

**Status:** ⏳ pending

## Goal

Add a tag browser mode to explore all tags in the vault and filter notes by tag. Tags are already parsed from frontmatter.

## Steps

### 1. Build tag index

During `ScanVault`, build a `map[string][]string` mapping each tag to the file paths that use it.

### 2. Add tag browser mode

New mode `ModeTags` triggered by `t` key. Shows all tags with file counts. Enter on a tag filters the tree to only show files with that tag.

### 3. Add tag-based filtering

When a tag is selected, the file tree shows only matching files. Press `Esc` to clear the filter.

## Completion Criteria

- [ ] Tag browser shows all tags with file counts
- [ ] Enter on a tag filters the tree
- [ ] Esc clears tag filter
- [ ] `make test && make vet` pass
