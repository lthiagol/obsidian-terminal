package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadNote_EmptyFilename(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, ".md"), []byte("# Only extension\n"), 0644)

	note, err := LoadNote(dir, ".md")
	if err != nil {
		t.Fatalf("LoadNote: %v", err)
	}
	if note.Title == "" {
		t.Error("should have fallback title for .md file")
	}
}

func TestParseFrontmatter_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	content := "---\ninvalid: [: yaml\n---\n\n# After frontmatter\n"
	os.WriteFile(filepath.Join(dir, "bad.md"), []byte(content), 0644)

	note, err := LoadNote(dir, "bad.md")
	if err != nil {
		t.Fatalf("LoadNote: %v", err)
	}
	if note.Body == "" || !strings.Contains(note.Body, "After frontmatter") {
		t.Error("body should contain content after frontmatter even with invalid YAML")
	}
}

func testVaultPath(t *testing.T) string {
	t.Helper()
	wd, _ := os.Getwd()
	return filepath.Join(wd, "testdata", "test-vault")
}

func TestScanVault_Structure(t *testing.T) {
	skipDirs := []string{".obsidian", ".git", ".trash", "node_modules", "archive"}
	tree, _, _, err := ScanVault(testVaultPath(t), skipDirs)
	if err != nil {
		t.Fatalf("ScanVault failed: %v", err)
	}

	if tree == nil {
		t.Fatal("expected non-nil tree")
	}

	// Should find projects/ and notes/ dirs
	foundProjects := false
	foundNotes := false
	for _, child := range tree.Children {
		if child.IsDir {
			switch child.Name {
			case "projects":
				foundProjects = true
			case "notes":
				foundNotes = true
			}
		}
	}

	if !foundProjects {
		t.Error("expected projects/ directory in tree")
	}
	if !foundNotes {
		t.Error("expected notes/ directory in tree")
	}
}

func TestScanVault_SkipsExcluded(t *testing.T) {
	skipDirs := []string{".obsidian", ".git", ".trash", "node_modules", "archive"}
	tree, _, _, err := ScanVault(testVaultPath(t), skipDirs)
	if err != nil {
		t.Fatalf("ScanVault failed: %v", err)
	}

	// Walk tree and ensure no excluded paths exist
	var paths []string
	var walk func(e *VaultEntry)
	walk = func(e *VaultEntry) {
		paths = append(paths, e.Path)
		for _, c := range e.Children {
			walk(c)
		}
	}
	walk(tree)

	for _, p := range paths {
		for _, part := range strings.Split(p, string(filepath.Separator)) {
			if strings.HasPrefix(part, ".") && part != "." {
				t.Errorf("found hidden path component in tree: %s", p)
			}
		}
		if strings.HasPrefix(filepath.Base(p), ".gitignore") {
			continue // filtered by walk
		}
	}

	// Verify .hidden-dir is not in tree
	for _, p := range paths {
		if strings.Contains(p, ".hidden") {
			t.Errorf("found hidden directory in tree: %s", p)
		}
		if strings.Contains(p, ".obsidian") {
			t.Errorf("found .obsidian in tree: %s", p)
		}
	}
}

func TestScanVault_SortsFoldersFirst(t *testing.T) {
	skipDirs := []string{".obsidian", ".git", ".trash", "node_modules", "archive"}
	tree, _, _, err := ScanVault(testVaultPath(t), skipDirs)
	if err != nil {
		t.Fatalf("ScanVault failed: %v", err)
	}

	// Check any directory's children are sorted (folders first, then alphabetical)
	var checkSort func(entries []*VaultEntry)
	checkSort = func(entries []*VaultEntry) {
		seenFile := false
		for i, entry := range entries {
			if !entry.IsDir {
				seenFile = true
			} else if seenFile {
				t.Errorf("folder %s found after files at index %d", entry.Name, i)
			}
			if i > 0 {
				prev := entries[i-1]
				if prev.IsDir == entry.IsDir {
					if strings.ToLower(prev.Name) > strings.ToLower(entry.Name) {
						t.Errorf("sort order: %s before %s", prev.Name, entry.Name)
					}
				}
			}
			if entry.Children != nil {
				checkSort(entry.Children)
			}
		}
	}
	checkSort(tree.Children)
}

func TestLoadNote_Frontmatter(t *testing.T) {
	note, err := LoadNote(testVaultPath(t), "index.md")
	if err != nil {
		t.Fatalf("LoadNote failed: %v", err)
	}

	if note.Title != "Welcome to Test Vault" {
		t.Errorf("title = %q, want 'Welcome to Test Vault'", note.Title)
	}

	expectedTags := []string{"test", "vault", "getting-started"}
	if len(note.Tags) != len(expectedTags) {
		t.Errorf("tags count = %d, want %d", len(note.Tags), len(expectedTags))
	}
	for i, tag := range expectedTags {
		if i >= len(note.Tags) || note.Tags[i] != tag {
			t.Errorf("tag[%d] = %q, want %q", i, note.Tags[i], tag)
		}
	}

	if len(note.Aliases) < 2 {
		t.Errorf("aliases count = %d, want at least 2", len(note.Aliases))
	}

	if note.Body == "" {
		t.Error("body should not be empty")
	}
	if strings.Contains(note.Body, "---") {
		t.Error("body should not contain frontmatter markers")
	}
}

func TestLoadNote_PlainMarkdown(t *testing.T) {
	note, err := LoadNote(testVaultPath(t), "notes/no-frontmatter.md")
	if err != nil {
		t.Fatalf("LoadNote failed: %v", err)
	}

	if note.Title != "No-frontmatter" {
		t.Errorf("title = %q, want 'No-frontmatter' (from filename)", note.Title)
	}

	if len(note.Tags) != 0 {
		t.Errorf("tags should be empty, got %v", note.Tags)
	}

	if note.Body == "" {
		t.Error("body should not be empty")
	}
}

func TestLoadNote_FrontmatterComplex(t *testing.T) {
	note, err := LoadNote(testVaultPath(t), "notes/frontmatter-test.md")
	if err != nil {
		t.Fatalf("LoadNote failed: %v", err)
	}

	if note.Title != "Frontmatter Test" {
		t.Errorf("title = %q", note.Title)
	}

	if len(note.Tags) < 3 {
		t.Errorf("tags count = %d, want at least 3", len(note.Tags))
	}

	if len(note.Aliases) < 2 {
		t.Errorf("aliases count = %d, want at least 2", len(note.Aliases))
	}
}

func TestSearchIndex(t *testing.T) {
	skipDirs := []string{".obsidian", ".git", ".trash", "node_modules", "archive"}
	_, indexes, _, err := ScanVault(testVaultPath(t), skipDirs)
	if err != nil {
		t.Fatalf("ScanVault failed: %v", err)
	}

	if len(indexes.Search) == 0 {
		t.Error("search index should not be empty")
	}

	foundIndex := false
	for path, content := range indexes.Search {
		if strings.Contains(path, "index.md") {
			foundIndex = true
			if !strings.Contains(content, "Welcome") {
				t.Error("index.md content should contain 'Welcome'")
			}
			break
		}
	}
	if !foundIndex {
		t.Error("search index should contain index.md")
	}
}

func TestExtractTagsFromFrontmatter(t *testing.T) {
	content := "---\ntitle: Test\ntags: [Go, PARSER, #yaml]\n---\n\nBody"
	tags := extractTagsFromFrontmatter(content)
	if len(tags) != 3 {
		t.Fatalf("expected 3 tags, got %d: %v", len(tags), tags)
	}
	expected := map[string]bool{"go": true, "parser": true, "yaml": true}
	for _, tag := range tags {
		if !expected[tag] {
			t.Errorf("unexpected tag: %q", tag)
		}
	}
}

func TestExtractWikiLinkTargets(t *testing.T) {
	content := "See [[note-a]] and [[note-b|display]] and [[note-a]] again"
	targets := extractWikiLinkTargets(content)
	if len(targets) != 2 {
		t.Fatalf("expected 2 unique targets, got %d: %v", len(targets), targets)
	}
	if targets[0] != "note-a" {
		t.Errorf("targets[0] = %q, want 'note-a'", targets[0])
	}
	if targets[1] != "note-b" {
		t.Errorf("targets[1] = %q, want 'note-b'", targets[1])
	}
}

func TestNormalizeWikiLinkTarget(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"Note", "note.md"},
		{"note.md", "note.md"},
		{"PATH/Note", "path/note.md"},
		{"Note.MD", "note.md"},
	}
	for _, tt := range tests {
		got := normalizeWikiLinkTarget(tt.input)
		if got != tt.want {
			t.Errorf("normalizeWikiLinkTarget(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestVaultIndexes_Backlinks(t *testing.T) {
	skipDirs := []string{".obsidian", ".git", ".trash", "node_modules", "archive"}
	_, indexes, _, err := ScanVault(testVaultPath(t), skipDirs)
	if err != nil {
		t.Fatalf("ScanVault failed: %v", err)
	}

	if len(indexes.Backlinks) == 0 {
		t.Error("backlink index should not be empty for test vault")
	}
}

func TestVaultIndexes_Tags(t *testing.T) {
	skipDirs := []string{".obsidian", ".git", ".trash", "node_modules", "archive"}
	_, indexes, _, err := ScanVault(testVaultPath(t), skipDirs)
	if err != nil {
		t.Fatalf("ScanVault failed: %v", err)
	}

	if len(indexes.Tags) == 0 {
		t.Error("tag index should not be empty for test vault")
	}

	if _, ok := indexes.Tags["test"]; !ok {
		t.Error("tag index should contain 'test' tag")
	}
}
