package main

import (
	"strings"
	"testing"

	"github.com/lthiagol/obsidian-terminal/internal/markdown"
)

func TestTable_BasicTable(t *testing.T) {
	content := "| Name | Age |\n|------|-----|\n| John | 30  |\n| Jane | 25  |"
	lines := markdown.ParseMarkdown(content)

	tableCount := 0
	for _, line := range lines {
		if line.BlockType == markdown.BlockTable {
			tableCount++
			if len(line.TableCells) != 2 {
				t.Errorf("expected 2 cells, got %d", len(line.TableCells))
			}
		}
	}

	if tableCount != 3 {
		t.Fatalf("expected 3 table lines (header + 2 data), got %d", tableCount)
	}

	rendered := markdown.RenderMarkdown(lines, 80, defaultRendererStyle())
	if !strings.Contains(rendered, "┌") || !strings.Contains(rendered, "┐") {
		t.Error("rendered table should have box-drawing characters")
	}
	if !strings.Contains(rendered, "Name") || !strings.Contains(rendered, "John") {
		t.Error("rendered table should contain cell text")
	}
}

func TestTable_Alignment(t *testing.T) {
	content := "| Left | Center | Right |\n|:-----|:------:|------:|\n| a    |   b    |    c |"
	lines := markdown.ParseMarkdown(content)

	for _, line := range lines {
		if line.BlockType != markdown.BlockTable {
			continue
		}
		if len(line.TableCells) == 3 {
			if line.TableCells[0] != "Left" && line.TableCells[0] != "a" {
				continue
			}
			if line.TableCells[0] == "Left" {
				if line.TableCells[0] != "Left" {
					t.Errorf("cell 0 = %q, want Left", line.TableCells[0])
				}
				if line.TableCells[1] != "Center" {
					t.Errorf("cell 1 = %q, want Center", line.TableCells[1])
				}
				if line.TableCells[2] != "Right" {
					t.Errorf("cell 2 = %q, want Right", line.TableCells[2])
				}
			}
		}
	}
}

func TestTable_EscapedPipe(t *testing.T) {
	content := "| Col1 |\n|------|\n| a \\| b |"
	lines := markdown.ParseMarkdown(content)

	found := false
	for _, line := range lines {
		if line.BlockType == markdown.BlockTable && len(line.TableCells) > 0 && line.TableCells[0] == "a | b" {
			found = true
		}
	}
	if !found {
		t.Error("escaped pipe should be preserved as literal pipe in cell")
	}
}

func TestTable_SingleColumn(t *testing.T) {
	content := "| Only |\n|------|\n| one  |"
	lines := markdown.ParseMarkdown(content)

	for _, line := range lines {
		if line.BlockType == markdown.BlockTable {
			if len(line.TableCells) != 1 {
				t.Fatalf("expected 1 cell, got %d", len(line.TableCells))
			}
		}
	}

	rendered := markdown.RenderMarkdown(lines, 80, defaultRendererStyle())
	if !strings.Contains(rendered, "Only") {
		t.Error("single column table should contain cell text")
	}
}

func TestTable_NotTable(t *testing.T) {
	content := "| This is just | a paragraph with pipes"
	lines := markdown.ParseMarkdown(content)

	for _, line := range lines {
		if line.BlockType == markdown.BlockTable {
			t.Error("single pipe line without separator should not be a table")
		}
	}
}

func TestTable_ManyRows(t *testing.T) {
	var content string
	content = "| A | B |\n|---|---|\n"
	for i := 0; i < 10; i++ {
		content += "| x | y |\n"
	}
	content = strings.TrimSuffix(content, "\n")
	lines := markdown.ParseMarkdown(content)

	tableCount := 0
	for _, line := range lines {
		if line.BlockType == markdown.BlockTable {
			tableCount++
		}
	}
	if tableCount != 11 {
		t.Errorf("expected 11 table lines (header + 10 data), got %d", tableCount)
	}
}
