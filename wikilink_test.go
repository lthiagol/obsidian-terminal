package main

import (
	"testing"
)

func TestFindAlias_FileReadError(t *testing.T) {
	aliasEntries, err := extractAliasesFromFile("/nonexistent/path", "no-such-file.md")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
	if aliasEntries != nil {
		t.Error("expected nil aliases for error")
	}
}

func TestWikiLinkResolution(t *testing.T) {
	skipDirs := DefaultConfig().SkipDirs
	vault, _, _, err := ScanVault(testVaultPath(t), skipDirs)
	if err != nil {
		t.Fatalf("ScanVault: %v", err)
	}

	resolved := ResolveWikiLink("index", vault, testVaultPath(t))
	if resolved != "index.md" {
		t.Errorf("resolved index = %q, want 'index.md'", resolved)
	}

	resolved = ResolveWikiLink("projects/api-design", vault, testVaultPath(t))
	if resolved != "projects/api-design.md" {
		t.Errorf("resolved api-design = %q, want 'projects/api-design.md'", resolved)
	}

	resolved = ResolveWikiLink("database", vault, testVaultPath(t))
	if resolved != "projects/database.md" {
		t.Errorf("resolved database = %q, want 'projects/database.md'", resolved)
	}

	resolved = ResolveWikiLink("nonexistent", vault, testVaultPath(t))
	if resolved != "" {
		t.Errorf("nonexistent should resolve to empty, got %q", resolved)
	}
}
