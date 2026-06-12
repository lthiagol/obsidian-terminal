package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/lthiagol/obsidian-terminal/internal/markdown"
)

// VaultIndexes holds all indexes built during vault scanning.
type VaultIndexes struct {
	Search    map[string]string
	Backlinks map[string][]string
	Tags      map[string][]string
}

// VaultEntry represents a file or directory in the vault tree.
type VaultEntry struct {
	Name      string
	Path      string
	IsDir     bool
	IsSymlink bool
	Children  []*VaultEntry
}

// VaultNote represents a parsed markdown note with frontmatter.
type VaultNote struct {
	Path    string
	Title   string
	Tags    []string
	Aliases []string
	Body    string
	RawBody string
}

type frontmatterData struct {
	Title   string
	Tags    []string
	Aliases []string
}

// ScanVault walks the vault directory and builds the file tree and all indexes.
func ScanVault(root string, skipDirs []string) (*VaultEntry, *VaultIndexes, []string, error) {
	skipSet := make(map[string]bool)
	for _, d := range skipDirs {
		skipSet[d] = true
	}

	var entries []*VaultEntry
	indexes := &VaultIndexes{
		Search:    make(map[string]string),
		Backlinks: make(map[string][]string),
		Tags:      make(map[string][]string),
	}
	var scanErrors []string

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			scanErrors = append(scanErrors, fmt.Sprintf("walk: %s: %v", path, err))
			if path == root {
				return err
			}
			return nil
		}

		relPath, _ := filepath.Rel(root, path)
		if relPath == "." {
			return nil
		}

		parts := strings.Split(relPath, string(filepath.Separator))

		for _, part := range parts {
			if strings.HasPrefix(part, ".") || skipSet[part] {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		isSymlink := d.Type()&os.ModeSymlink != 0

		ext := strings.ToLower(filepath.Ext(d.Name()))
		isMd := ext == ".md" || ext == ".markdown"

		if d.IsDir() {
			entries = append(entries, &VaultEntry{
				Name:      d.Name(),
				Path:      relPath,
				IsDir:     true,
				IsSymlink: isSymlink,
			})
		} else if isMd {
			entries = append(entries, &VaultEntry{
				Name:      d.Name(),
				Path:      relPath,
				IsDir:     false,
				IsSymlink: isSymlink,
			})

			data, readErr := os.ReadFile(path)
			if readErr != nil {
				scanErrors = append(scanErrors, fmt.Sprintf("read: %s: %v", relPath, readErr))
			} else {
				content := string(data)

				indexes.Search[relPath] = stripFrontmatter(content)

				for _, tag := range extractTagsFromFrontmatter(content) {
					indexes.Tags[tag] = append(indexes.Tags[tag], relPath)
				}

				for _, target := range extractWikiLinkTargets(content) {
					normalized := normalizeWikiLinkTarget(target)
					indexes.Backlinks[normalized] = append(indexes.Backlinks[normalized], relPath)
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, nil, scanErrors, fmt.Errorf("walking vault: %w", err)
	}

	tree := buildTree(entries)
	return tree, indexes, scanErrors, nil
}

func buildTree(entries []*VaultEntry) *VaultEntry {
	root := &VaultEntry{
		Name:     ".",
		Path:     "",
		IsDir:    true,
		Children: nil,
	}

	for _, entry := range entries {
		parts := strings.Split(entry.Path, string(filepath.Separator))
		current := root

		for i, part := range parts {
			if i == len(parts)-1 {
				current.Children = append(current.Children, entry)
			} else {
				found := false
				for _, child := range current.Children {
					if child.Name == part && child.IsDir {
						current = child
						found = true
						break
					}
				}
				if !found {
					dirPath := strings.Join(parts[:i+1], string(filepath.Separator))
					dir := &VaultEntry{
						Name:     part,
						Path:     dirPath,
						IsDir:    true,
						Children: nil,
					}
					current.Children = append(current.Children, dir)
					current = dir
				}
			}
		}
	}

	sortVaultEntries(root.Children)
	return root
}

func sortVaultEntries(entries []*VaultEntry) {
	sort.SliceStable(entries, func(i, j int) bool {
		if entries[i].IsDir != entries[j].IsDir {
			return entries[i].IsDir
		}
		return strings.ToLower(entries[i].Name) < strings.ToLower(entries[j].Name)
	})

	for _, e := range entries {
		if e.IsDir && e.Children != nil {
			sortVaultEntries(e.Children)
		}
	}
}

// LoadNote reads a markdown file and parses its frontmatter and body.
func LoadNote(vaultRoot, relativePath string) (*VaultNote, error) {
	fullPath := filepath.Join(vaultRoot, relativePath)

	info, err := os.Lstat(fullPath)
	if err != nil {
		return nil, fmt.Errorf("stat: %w", err)
	}

	// Resolve symlinks for reading content
	readPath := fullPath
	if info.Mode()&os.ModeSymlink != 0 {
		resolved, err := filepath.EvalSymlinks(fullPath)
		if err != nil {
			return nil, fmt.Errorf("resolving symlink: %w", err)
		}
		readPath = resolved
	}

	data, err := os.ReadFile(readPath)
	if err != nil {
		return nil, fmt.Errorf("reading: %w", err)
	}

	content := string(data)
	fm, body := parseFrontmatter(content)

	title := fm.Title
	if title == "" {
		base := filepath.Base(relativePath)
		ext := filepath.Ext(base)
		name := strings.TrimSuffix(base, ext)
		if len(name) > 0 {
			title = strings.ToUpper(name[:1]) + name[1:]
		} else {
			title = "Untitled"
		}
	}

	return &VaultNote{
		Path:    relativePath,
		Title:   title,
		Tags:    fm.Tags,
		Aliases: fm.Aliases,
		Body:    body,
		RawBody: content,
	}, nil
}

func findFrontmatterBounds(content string) (yamlStart, yamlEnd int, ok bool) {
	if !strings.HasPrefix(content, "---\n") && !strings.HasPrefix(content, "---\r\n") {
		return 0, 0, false
	}
	rest := content[3:]
	idx := strings.Index(rest, "\n---\n")
	if idx == -1 {
		idx = strings.Index(rest, "\n---\r\n")
	}
	if idx == -1 && strings.HasSuffix(rest, "\n---") {
		return 3, len(content) - 4, true
	}
	if idx == -1 {
		return 0, 0, false
	}
	return 3, 3 + idx, true
}

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

	bodyStart := yamlEnd + 5
	if bodyStart > len(content) {
		bodyStart = len(content)
	}

	return fm, content[bodyStart:]
}

func stripFrontmatter(content string) string {
	_, yamlEnd, ok := findFrontmatterBounds(content)
	if !ok {
		return content
	}
	bodyStart := yamlEnd + 5
	if bodyStart > len(content) {
		return ""
	}
	return content[bodyStart:]
}

func allPaths(vault *VaultEntry) []string {
	var paths []string
	collectPaths(vault, "", &paths)
	return paths
}

func collectPaths(entry *VaultEntry, prefix string, paths *[]string) {
	if !entry.IsDir {
		*paths = append(*paths, entry.Path)
		return
	}
	for _, child := range entry.Children {
		collectPaths(child, "", paths)
	}
}

var wikiLinkTargetRe = regexp.MustCompile(`\[\[([^\]|#]+)`)

func extractTagsFromFrontmatter(content string) []string {
	fm, _ := parseFrontmatter(content)
	seen := make(map[string]bool)
	var tags []string
	for _, tag := range fm.Tags {
		tag = strings.TrimPrefix(tag, "#")
		tag = strings.ToLower(tag)
		if tag != "" && !seen[tag] {
			seen[tag] = true
			tags = append(tags, tag)
		}
	}
	return tags
}

func extractWikiLinkTargets(content string) []string {
	matches := wikiLinkTargetRe.FindAllStringSubmatch(content, -1)
	seen := make(map[string]bool)
	var targets []string
	for _, match := range matches {
		target := match[1]
		if !seen[target] {
			seen[target] = true
			targets = append(targets, target)
		}
	}
	return targets
}

func normalizeWikiLinkTarget(target string) string {
	target = strings.ToLower(target)
	if !strings.HasSuffix(target, ".md") {
		target += ".md"
	}
	return target
}

func extractSection(content, heading string) string {
	lines := strings.Split(content, "\n")

	headingLevel := 0
	startIdx := -1

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if markdown.IsHeading(trimmed) {
			level := markdown.HeadingLevel(trimmed)
			text := strings.TrimSpace(trimmed[level:])
			if strings.ToLower(text) == strings.ToLower(heading) {
				headingLevel = level
				startIdx = i
				break
			}
		}
	}

	if startIdx < 0 {
		return content
	}

	var result []string
	for i := startIdx; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if i > startIdx && markdown.IsHeading(trimmed) {
			level := markdown.HeadingLevel(trimmed)
			if level <= headingLevel {
				break
			}
		}
		result = append(result, lines[i])
	}

	return strings.Join(result, "\n")
}

// isMarkdownHeading and countHeadingLevel were removed.
// Use markdown.IsHeading and markdown.HeadingLevel instead.
