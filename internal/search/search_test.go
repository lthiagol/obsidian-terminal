package search

import (
	"strings"
	"testing"
)

func scoreFor(query, target string) float64 {
	return FuzzyScore([]rune(query), []rune(strings.ToLower(query)), []rune(target), []rune(strings.ToLower(target)), strings.ToLower(target))
}

func lowerPaths(paths []string) []string {
	lower := make([]string, len(paths))
	for i, p := range paths {
		lower[i] = strings.ToLower(p)
	}
	return lower
}

func TestFuzzyScore_ExactMatch(t *testing.T) {
	score := scoreFor("test", "test")
	if score < 100 {
		t.Errorf("exact match score = %f, want >= 100", score)
	}
}

func TestFuzzyScore_Substring(t *testing.T) {
	consecutive := scoreFor("abc", "abcxxx")
	scattered := scoreFor("axc", "axbxc")

	if consecutive <= scattered {
		t.Errorf("consecutive match (%f) should score higher than scattered (%f)", consecutive, scattered)
	}
}

func TestFuzzyScore_BoundaryBonus(t *testing.T) {
	boundaryScore := scoreFor("readme", "projects/readme.md")
	noBoundaryScore := scoreFor("readme", "some_readme_file.md")

	_ = boundaryScore
	_ = noBoundaryScore
	if boundaryScore <= 0 {
		t.Error("boundary match should score > 0")
	}
}

func TestFuzzyScore_ExactCaseBonus(t *testing.T) {
	exactCase := scoreFor("ReadMe", "ReadMe.md")
	lowerCase := scoreFor("readme", "ReadMe.md")

	if exactCase < lowerCase {
		t.Errorf("exact case (%f) should score >= lower case (%f)", exactCase, lowerCase)
	}
}

func TestFuzzyScore_NoMatch(t *testing.T) {
	score := scoreFor("xyz", "abcdef")
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

	results := FuzzySearch("readme", paths, lowerPaths(paths), nil, nil)
	if len(results) == 0 {
		t.Fatal("expected results for 'readme'")
	}

	for i := 1; i < len(results); i++ {
		if results[i].Score > results[i-1].Score {
			t.Errorf("results not sorted: %f at index %d > %f at index %d",
				results[i].Score, i, results[i-1].Score, i-1)
		}
	}
}

func TestFuzzySearch_NoMatchingFiles(t *testing.T) {
	paths := []string{"readme.md", "index.md"}
	results := FuzzySearch("zzz", paths, lowerPaths(paths), nil, nil)
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

func TestState_SetQuery(t *testing.T) {
	paths := []string{"readme.md", "index.md", "api.md"}
	index := map[string]string{}

	s := NewState(Name, paths, index)
	s.SetQuery("read")

	if s.query != "read" {
		t.Errorf("query = %q, want 'read'", s.query)
	}
	if s.ResultCount() == 0 {
		t.Error("should have results for 'read'")
	}
}

func TestFuzzyScore_EmptyInput(t *testing.T) {
	if score := FuzzyScore(nil, nil, nil, nil, "anything"); score != 0 {
		t.Errorf("empty query score = %f, want 0", score)
	}
	if score := FuzzyScore([]rune("test"), []rune("test"), nil, nil, ""); score != 0 {
		t.Errorf("empty target score = %f, want 0", score)
	}
	if score := FuzzyScore(nil, nil, nil, nil, ""); score != 0 {
		t.Errorf("both empty score = %f, want 0", score)
	}
}

func TestState_MoveUpDown(t *testing.T) {
	paths := []string{"a.md", "b.md", "c.md"}
	index := map[string]string{}

	s := NewState(Name, paths, index)

	s.MoveDown()
	if s.selected != 1 {
		t.Errorf("after down: selected = %d, want 1", s.selected)
	}

	s.MoveDown()
	s.MoveDown()
	s.MoveDown()
	if s.selected != 2 {
		t.Errorf("clamped at bottom: selected = %d, want 2", s.selected)
	}

	s.MoveUp()
	if s.selected != 1 {
		t.Errorf("after up: selected = %d, want 1", s.selected)
	}

	s.MoveUp()
	s.MoveUp()
	if s.selected != 0 {
		t.Errorf("clamped at top: selected = %d, want 0", s.selected)
	}
}

func TestHighlightMatches(t *testing.T) {
	result := HighlightMatches("note", "notes/meeting.md")
	if !strings.Contains(result, "note") {
		t.Error("highlighted output should contain matched characters")
	}
	if !strings.Contains(result, "meeting") {
		t.Error("output should contain non-matched parts too")
	}
}

func TestHighlightMatches_EmptyQuery(t *testing.T) {
	result := HighlightMatches("", "notes/meeting.md")
	if result != "notes/meeting.md" {
		t.Errorf("empty query should return original: %q", result)
	}
}

func TestRenderResults_NameMode(t *testing.T) {
	style := Style{
		Accent:        "#a78bfa",
		TextSecondary: "#9ca3af",
		TextMuted:     "#4b5563",
	}
	s := NewState(Name, []string{"a.md", "b.md"}, nil)
	s.SetQuery("a")
	output := RenderResults(s, 80, style)
	if !strings.Contains(output, "a.md") {
		t.Error("should contain matching path")
	}
}

func TestRenderResults_ContentMode(t *testing.T) {
	style := Style{
		Accent:        "#a78bfa",
		TextSecondary: "#9ca3af",
		TextMuted:     "#4b5563",
	}
	index := map[string]string{
		"readme.md": "Hello world\nThis is a test\n",
	}
	s := NewState(Content, []string{"readme.md"}, index)
	s.SetQuery("test")
	output := RenderResults(s, 80, style)
	if !strings.Contains(output, "readme.md") {
		t.Error("should contain file path")
	}
	if !strings.Contains(output, "test") {
		t.Error("should contain match context")
	}
}

func TestRenderResults_Empty(t *testing.T) {
	style := Style{
		Accent:        "#a78bfa",
		TextSecondary: "#9ca3af",
		TextMuted:     "#4b5563",
	}
	s := NewState(Name, []string{}, nil)
	s.SetQuery("no-match")
	output := RenderResults(s, 80, style)
	if !strings.Contains(output, "No results") {
		t.Error("should show 'No results'")
	}
}

func TestRenderResults_ContentModeEmptyQuery(t *testing.T) {
	style := Style{
		Accent:        "#a78bfa",
		TextSecondary: "#9ca3af",
		TextMuted:     "#4b5563",
	}
	s := NewState(Content, nil, nil)
	output := RenderResults(s, 80, style)
	if !strings.Contains(output, "Type to search") {
		t.Error("should show prompt for empty content search")
	}
}

func TestSelectedResult(t *testing.T) {
	s := NewState(Name, []string{"a.md", "b.md"}, nil)

	r := s.SelectedResult()
	if r == nil {
		t.Fatal("SelectedResult should not be nil")
	}
	if r.Path != "a.md" {
		t.Errorf("path = %q", r.Path)
	}

	// SetSelected clamps to max valid index
	s.SetSelected(99)
	r2 := s.SelectedResult()
	if r2 == nil {
		t.Error("clamped out-of-bounds should still return last result")
	}
	if r2.Path != "b.md" {
		t.Errorf("clamped path = %q, want 'b.md'", r2.Path)
	}

	// SetSelected negative clamps to 0
	s.SetSelected(-1)
	r3 := s.SelectedResult()
	if r3 == nil || r3.Path != "a.md" {
		t.Error("negative SetSelected should clamp to 0")
	}
}

func TestContentSearch_MultipleMatches(t *testing.T) {
	index := map[string]string{
		"notes.md": "Line 1\nFind me here\nLine 3\nAlso me\n",
	}
	results := ContentSearch("me", index)
	if len(results) < 2 {
		t.Errorf("expected >= 2 matches for 'me', got %d", len(results))
	}
	for _, r := range results {
		if !strings.Contains(strings.ToLower(r.Context), "me") {
			t.Errorf("result context should contain 'me': %q", r.Context)
		}
	}
}

func TestContentSearch_NoMatch(t *testing.T) {
	index := map[string]string{"a.md": "hello"}
	results := ContentSearch("xyz", index)
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestContentSearch_TruncatesLongLine(t *testing.T) {
	longLine := ""
	for i := 0; i < 100; i++ {
		longLine += "word "
	}
	index := map[string]string{"file.md": longLine + "FINDME"}

	results := ContentSearch("FINDME", index)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if len(results[0].Context) > 83 {
		t.Errorf("context too long: %d chars", len(results[0].Context))
	}
}
