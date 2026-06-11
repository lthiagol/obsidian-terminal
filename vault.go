package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

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
	Title   string   `yaml:"title"`
	Tags    []string `yaml:"tags"`
	Aliases []string `yaml:"aliases"`
}

// ScanVault walks the vault directory and builds the file tree and search index.
func ScanVault(root string, skipDirs []string) (*VaultEntry, map[string]string, []string, error) {
	skipSet := make(map[string]bool)
	for _, d := range skipDirs {
		skipSet[d] = true
	}

	var entries []*VaultEntry
	searchIndex := make(map[string]string)
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

		// Check if any path component should be skipped
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

			// Build search index
			data, readErr := os.ReadFile(path)
			if readErr != nil {
				scanErrors = append(scanErrors, fmt.Sprintf("read: %s: %v", relPath, readErr))
			} else {
				searchIndex[relPath] = stripFrontmatter(string(data))
			}
		}

		return nil
	})

	if err != nil {
		return nil, nil, scanErrors, fmt.Errorf("walking vault: %w", err)
	}

	tree := buildTree(entries)
	return tree, searchIndex, scanErrors, nil
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
	if err := yaml.Unmarshal([]byte(yamlBlock), &fm); err != nil {
		return frontmatterData{Title: fm.Title}, content
	}
	return fm, content[yamlEnd+5:]
}

func stripFrontmatter(content string) string {
	_, yamlEnd, ok := findFrontmatterBounds(content)
	if !ok {
		return content
	}
	return content[yamlEnd+5:]
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
