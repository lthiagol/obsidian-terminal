package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/lthiagol/obsidian-terminal/internal/markdown"
	"github.com/lthiagol/obsidian-terminal/internal/search"
)

func createTestTree(count int) *VaultEntry {
	root := &VaultEntry{Name: ".", Path: "", IsDir: true}
	for i := 0; i < count; i++ {
		dir := &VaultEntry{
			Name:  "dir" + strings.Repeat("x", 10),
			Path:  "dir" + strings.Repeat("x", 10) + string(rune('0'+i%10)),
			IsDir: true,
		}
		for j := 0; j < 10; j++ {
			file := &VaultEntry{
				Name:  "note" + strings.Repeat("y", 15) + string(rune('0'+j)),
				Path:  "path/note" + strings.Repeat("y", 15) + string(rune('0'+j)) + ".md",
				IsDir: false,
			}
			dir.Children = append(dir.Children, file)
		}
		root.Children = append(root.Children, dir)
	}
	return root
}

func BenchmarkFileTreeView(b *testing.B) {
	tree := createTestTree(100)
	ft := NewFileTree(tree, newDarkPalette())
	ft.SetSize(30, 50)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ft.View()
	}
}

func BenchmarkFuzzySearch(b *testing.B) {
	paths := make([]string, 10000)
	for i := range paths {
		paths[i] = "vault/path/note" + strings.Repeat("x", 20) + string(rune('0'+i%10)) + ".md"
	}
	lower := make([]string, len(paths))
	for i, p := range paths {
		lower[i] = strings.ToLower(p)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = search.FuzzySearch("test", paths, lower)
	}
}

func BenchmarkMarkdownRender(b *testing.B) {
	content := "# Heading\n\n" + strings.Repeat("This is a paragraph with **bold** and *italic* text. ", 50) + "\n\n"
	content += "## Subheading\n\n- item 1\n- item 2\n- item 3\n\n"
	content += "```go\nfunc main() {\n\tfmt.Println(\"hello\")\n}\n```\n\n"
	content += strings.Repeat("More content here. ", 100)

	lines := markdown.ParseMarkdown(content)
	style := markdown.RendererStyle{
		Accent:          "#a78bfa",
		AccentSecondary: "#fbbf24",
		AccentTertiary:  "#2dd4bf",
		TextSecondary:   "#9ca3af",
		TextDim:         "#4b5563",
		Success:         "#34d399",
		CodeBackground:  "#1f2937",
		Heading1:        "#fbbf24",
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = markdown.RenderMarkdown(lines, 80, style)
	}
}

func BenchmarkHelpRender(b *testing.B) {
	cfg := &Config{VaultPath: "testdata/test-vault", SkipDirs: DefaultConfig().SkipDirs}
	m := NewModel(cfg)
	m.width = 120
	m.height = 40

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = m.renderHelp()
	}
}

func BenchmarkViewportView(b *testing.B) {
	vp := newViewport(80, 24)
	content := strings.Repeat("Line of text in the viewport\n", 500)
	vp.SetContent(content)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = vp.View()
	}
}

func BenchmarkScanVault(b *testing.B) {
	skipDirs := []string{".obsidian", ".git", ".trash", "node_modules", "archive"}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _, _, _ = ScanVault("testdata/test-vault", skipDirs)
	}
}

func BenchmarkScanVault_1k(b *testing.B) {
	dir := b.TempDir()
	generateBenchVault(dir, 100, 10)
	skipDirs := []string{".obsidian", ".git", ".trash", "node_modules", "archive"}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _, _, _ = ScanVault(dir, skipDirs)
	}
}

func BenchmarkScanVault_5k(b *testing.B) {
	dir := b.TempDir()
	generateBenchVault(dir, 500, 10)
	skipDirs := []string{".obsidian", ".git", ".trash", "node_modules", "archive"}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _, _, _ = ScanVault(dir, skipDirs)
	}
}

func generateBenchVault(root string, dirs, filesPerDir int) {
	for i := 0; i < dirs; i++ {
		dirPath := fmt.Sprintf("%s/dir%d", root, i)
		os.MkdirAll(dirPath, 0755)
		for j := 0; j < filesPerDir; j++ {
			content := fmt.Sprintf("---\ntitle: Note %d-%d\ntags:\n  - bench\n---\n\n# Note %d-%d\n\nThis is benchmark content.\n", i, j, i, j)
			os.WriteFile(fmt.Sprintf("%s/note%d.md", dirPath, j), []byte(content), 0644)
		}
	}
}
