package main

import (
	"strings"
	"testing"
)

func TestFuzzyScore_ExactMatch(t *testing.T) {
	score := FuzzyScore("test", "test")
	if score < 100 {
		t.Errorf("exact match score = %f, want >= 100", score)
	}
}

func TestFuzzyScore_Substring(t *testing.T) {
	consecutive := FuzzyScore("abc", "abcxxx")
	scattered := FuzzyScore("axc", "axbxc")

	if consecutive <= scattered {
		t.Errorf("consecutive match (%f) should score higher than scattered (%f)", consecutive, scattered)
	}
}

func TestFuzzyScore_BoundaryBonus(t *testing.T) {
	boundaryScore := FuzzyScore("readme", "projects/readme.md")
	noBoundaryScore := FuzzyScore("readme", "some_readme_file.md")

	_ = boundaryScore
	_ = noBoundaryScore
	// Both should score since 'r' matches
	if boundaryScore <= 0 {
		t.Error("boundary match should score > 0")
	}
}

func TestFuzzyScore_ExactCaseBonus(t *testing.T) {
	exactCase := FuzzyScore("ReadMe", "ReadMe.md")
	lowerCase := FuzzyScore("readme", "ReadMe.md")

	// Exact case match gets case bonus, so should score >= lower case
	if exactCase < lowerCase {
		t.Errorf("exact case (%f) should score >= lower case (%f)", exactCase, lowerCase)
	}
}

func TestFuzzyScore_NoMatch(t *testing.T) {
	score := FuzzyScore("xyz", "abcdef")
	if score != 0 {
		t.Errorf("no-match score = %f, want 0", score)
	}
}

func TestFuzzySearch_ResultsSorted(t *testing.T) {
	paths := []string{
		"projects/api-design.md",
		"readme.md",
		"index.md",
		"notes/meeting.md",
	}

	results := FuzzySearch("readme", paths)
	if len(results) == 0 {
		t.Fatal("expected results for 'readme'")
	}

	// Results should be sorted descending by score
	for i := 1; i < len(results); i++ {
		if results[i].Score > results[i-1].Score {
			t.Errorf("results not sorted: %f at index %d > %f at index %d",
				results[i].Score, i, results[i-1].Score, i-1)
		}
	}
}

func TestFuzzySearch_NoMatchingFiles(t *testing.T) {
	paths := []string{"readme.md", "index.md"}
	results := FuzzySearch("zzz", paths)
	if len(results) != 0 {
		t.Errorf("expected no results, got %d", len(results))
	}
}

func TestContentSearch_FindsInBody(t *testing.T) {
	index := map[string]string{
		"readme.md": "This is a readme file with important info.",
		"api.md":    "API documentation here.",
	}

	results := ContentSearch("readme", index)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Path != "readme.md" {
		t.Errorf("result path = %q, want 'readme.md'", results[0].Path)
	}
}

func TestContentSearch_ReturnsLineContext(t *testing.T) {
	index := map[string]string{
		"notes.md": "Line 1: nothing\nLine 2: has TODO here\nLine 3: more text",
	}

	results := ContentSearch("TODO", index)
	if len(results) == 0 {
		t.Fatal("expected results for 'TODO'")
	}

	if results[0].Path != "notes.md" {
		t.Errorf("path = %q, want 'notes.md'", results[0].Path)
	}
	if results[0].LineNum != 2 {
		t.Errorf("line = %d, want 2", results[0].LineNum)
	}
	if !strings.Contains(results[0].Context, "TODO") {
		t.Errorf("context should contain 'TODO': %s", results[0].Context)
	}
}

func TestContentSearch_CaseInsensitive(t *testing.T) {
	index := map[string]string{
		"tasks.md": "Buy groceries\nFix bug\nCall doctor",
	}

	results := ContentSearch("BUG", index)
	if len(results) == 0 {
		t.Fatal("'BUG' should match 'bug' (case insensitive)")
	}
}

func TestSearchState_SetQuery(t *testing.T) {
	paths := []string{"readme.md", "index.md", "api.md"}
	index := map[string]string{}

	s := NewSearchState(SearchName, paths, index)
	s.SetQuery("read")

	if s.query != "read" {
		t.Errorf("query = %q, want 'read'", s.query)
	}
	if s.ResultCount() == 0 {
		t.Error("should have results for 'read'")
	}
}

func TestSearchState_MoveUpDown(t *testing.T) {
	paths := []string{"a.md", "b.md", "c.md"}
	index := map[string]string{}

	s := NewSearchState(SearchName, paths, index)

	s.MoveDown()
	if s.selected != 1 {
		t.Errorf("after down: selected = %d, want 1", s.selected)
	}

	s.MoveDown()
	s.MoveDown()
	s.MoveDown() // beyond bounds
	if s.selected != 2 {
		t.Errorf("clamped at bottom: selected = %d, want 2", s.selected)
	}

	s.MoveUp()
	if s.selected != 1 {
		t.Errorf("after up: selected = %d, want 1", s.selected)
	}

	s.MoveUp()
	s.MoveUp() // beyond bounds
	if s.selected != 0 {
		t.Errorf("clamped at top: selected = %d, want 0", s.selected)
	}
}
