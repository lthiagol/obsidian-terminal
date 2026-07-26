package main

import (
	"strings"
	"testing"

	"github.com/lthiagol/obsidian-terminal/internal/markdown"
)

func TestEmbed_ParseNotation(t *testing.T) {
	content := "![[target]]"
	lines := markdown.ParseMarkdown(content)

	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}

	if lines[0].BlockType != markdown.BlockEmbed {
		t.Errorf("expected BlockEmbed, got %v", lines[0].BlockType)
	}
	if lines[0].EmbedTarget != "target" {
		t.Errorf("EmbedTarget = %q, want target", lines[0].EmbedTarget)
	}
	if lines[0].EmbedHeading != "" {
		t.Errorf("EmbedHeading should be empty, got %q", lines[0].EmbedHeading)
	}
}

func TestEmbed_ParseWithHeading(t *testing.T) {
	content := "![[note#heading]]"
	lines := markdown.ParseMarkdown(content)

	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}

	if lines[0].EmbedTarget != "note" {
		t.Errorf("EmbedTarget = %q, want note", lines[0].EmbedTarget)
	}
	if lines[0].EmbedHeading != "heading" {
		t.Errorf("EmbedHeading = %q, want heading", lines[0].EmbedHeading)
	}
}

func TestEmbed_ResolveEmbeds(t *testing.T) {
	content := "![[target]]"
	lines := markdown.ParseMarkdown(content)

	resolver := func(target, heading string) (string, error) {
		return "# Resolved: " + target, nil
	}

	resolved := markdown.ResolveEmbeds(lines, resolver)

	foundStart := false
	foundEnd := false
	foundContent := false
	for _, line := range resolved {
		if line.BlockType == markdown.BlockEmbedStart {
			foundStart = true
		}
		if line.BlockType == markdown.BlockEmbedEnd {
			foundEnd = true
		}
		if line.BlockType == markdown.BlockHeading {
			foundContent = true
		}
	}

	if !foundStart {
		t.Error("should have BlockEmbedStart")
	}
	if !foundEnd {
		t.Error("should have BlockEmbedEnd")
	}
	if !foundContent {
		t.Error("should have content from resolver")
	}
}

func TestEmbed_NotFound(t *testing.T) {
	content := "![[nonexistent]]"
	lines := markdown.ParseMarkdown(content)

	resolver := func(target, heading string) (string, error) {
		return "", nil
	}

	resolved := markdown.ResolveEmbeds(lines, resolver)

	found := false
	for _, line := range resolved {
		if line.BlockType == markdown.BlockParagraph &&
			len(line.Segments) > 0 &&
			strings.Contains(line.Segments[0].Text, "not found") {
			found = true
		}
	}
	if !found {
		t.Error("should show 'embed not found' message")
	}
}

func TestEmbed_CircularDetection(t *testing.T) {
	content := "![[self]]"
	lines := markdown.ParseMarkdown(content)

	resolver := func(target, heading string) (string, error) {
		return "![[self]]", nil // embeds itself
	}

	resolved := markdown.ResolveEmbeds(lines, resolver)

	found := false
	for _, line := range resolved {
		if line.BlockType == markdown.BlockEmbedStart &&
			strings.Contains(line.EmbedTarget, "circular") {
			found = true
		}
	}
	if !found {
		t.Error("should detect circular embeds")
	}
}

func TestEmbed_DepthLimit(t *testing.T) {
	content := "![[a]]"
	lines := markdown.ParseMarkdown(content)

	resolver := func(target, heading string) (string, error) {
		return "![[" + target + "2]]", nil // each level creates another embed
	}

	resolved := markdown.ResolveEmbeds(lines, resolver)

	embedCount := 0
	for _, line := range resolved {
		if line.BlockType == markdown.BlockEmbedStart {
			embedCount++
		}
	}

	if embedCount > 3 {
		t.Errorf("should limit embed depth: got %d embed blocks", embedCount)
	}
}

func TestExtractSection(t *testing.T) {
	content := "# Section 1\n\nContent 1\n\n## Subsection 1.1\n\nSubcontent\n\n# Section 2\n\nContent 2"

	section := extractSection(content, "Section 1")
	if !strings.Contains(section, "Content 1") {
		t.Error("section should contain Content 1")
	}
	if strings.Contains(section, "Content 2") {
		t.Error("section should not contain Content 2 (different heading level)")
	}
	if !strings.Contains(section, "Subsection 1.1") {
		t.Error("section should contain subsections")
	}
}

func TestExtractSection_NotFound(t *testing.T) {
	content := "# Section 1\n\nContent 1"
	section := extractSection(content, "Nonexistent")
	if section != content {
		t.Error("should return full content when heading not found")
	}
}
