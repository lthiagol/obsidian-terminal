package search

import (
	"fmt"
	"sort"
	"strings"
	"unicode"

	"github.com/charmbracelet/lipgloss"
)

// Mode specifies the type of search (name or content).
type Mode int

const (
	maxSearchResults       = 50
	maxContentResults      = 100
	contentResultContextLen = 80
)

const (
	Name    Mode = iota // Name is file name fuzzy search.
	Content             // Content is full-text content search.
)

// Result represents a single search match.
type Result struct {
	Path    string  // Path is the file path of the match.
	Score   float64 // Score is the relevance score (higher = better match).
	Context string  // Context is the matching line excerpt (content search only).
	LineNum int     // LineNum is the line number of the match (content search only).
}

// State holds the current state of a search session.
type State struct {
	mode        Mode
	query       string
	results     []Result
	selected    int
	allPaths    []string
	allLower    []string
	allRunes    [][]rune
	allLowerRunes [][]rune
	searchIndex map[string]string
}

// Style holds colors for search result rendering.
type Style struct {
	Accent        lipgloss.Color // Accent is the primary accent color.
	TextSecondary lipgloss.Color // TextSecondary is the secondary text color.
	TextMuted     lipgloss.Color // TextMuted is the muted text color.
	SelectionText lipgloss.Color // SelectionText is the text color for selected items.
}

// NewState creates a new search state with the given mode and data.
func NewState(mode Mode, paths []string, index map[string]string) State {
	lower := make([]string, len(paths))
	runes := make([][]rune, len(paths))
	lowerRunes := make([][]rune, len(paths))
	for i, p := range paths {
		lower[i] = strings.ToLower(p)
		runes[i] = []rune(p)
		lowerRunes[i] = []rune(lower[i])
	}
	s := State{
		mode:          mode,
		allPaths:      paths,
		allLower:      lower,
		allRunes:      runes,
		allLowerRunes: lowerRunes,
		searchIndex:   index,
		selected:      0,
	}
	if mode == Name {
		s.results = FuzzySearch("", paths, lower, runes, lowerRunes)
	}
	return s
}

// SetQuery updates the search query and re-runs the search.
func (s *State) SetQuery(query string) {
	s.query = query
	s.selected = 0

	switch s.mode {
	case Name:
		s.results = FuzzySearch(query, s.allPaths, s.allLower, s.allRunes, s.allLowerRunes)
	case Content:
		if query == "" {
			s.results = []Result{}
		} else {
			s.results = ContentSearch(query, s.searchIndex)
		}
	}
}

// MoveUp moves the selection cursor up one result.
func (s *State) MoveUp() {
	if s.selected > 0 {
		s.selected--
	}
}

// MoveDown moves the selection cursor down one result.
func (s *State) MoveDown() {
	if s.selected < len(s.results)-1 {
		s.selected++
	}
}

// SetSelected sets the selection cursor to a specific index, clamped to valid range.
func (s *State) SetSelected(i int) {
	if i < 0 {
		i = 0
	}
	if i >= len(s.results) {
		i = len(s.results) - 1
	}
	s.selected = i
}

// SelectedIndex returns the current selection index.
func (s State) SelectedIndex() int {
	return s.selected
}

// ResultCount returns the total number of search results.
func (s State) ResultCount() int {
	return len(s.results)
}

// Query returns the current search query string.
func (s State) Query() string {
	return s.query
}

// SelectedResult returns a pointer to the currently selected result, or nil if out of range.
func (s State) SelectedResult() *Result {
	if s.selected >= 0 && s.selected < len(s.results) {
		return &s.results[s.selected]
	}
	return nil
}

// FuzzyScore computes a match score between query and target.
// queryRunes and queryLowerRunes must be pre-computed from the same query.
// targetOrigRunes and targetLowerRunes must be pre-computed from the same target.
// targetLower is a string copy of targetLowerRunes for boundary checks.
func FuzzyScore(queryRunes, queryLowerRunes, targetOrigRunes, targetLowerRunes []rune, targetLower string) float64 {
	if len(queryLowerRunes) == 0 || len(targetLowerRunes) == 0 {
		return 0
	}

	if string(queryLowerRunes) == targetLower {
		return 100 + float64(len(queryLowerRunes))
	}

	qi := 0
	ti := 0
	consecutive := 0
	gapCount := 0
	var score float64

	exactCaseCount := 0

	for qi < len(queryLowerRunes) && ti < len(targetLowerRunes) {
		q := queryLowerRunes[qi]

		found := false
		for ti < len(targetLowerRunes) {
			t := targetLowerRunes[ti]
			if t == q {
				found = true
				break
			}
			ti++
			gapCount++
		}

		if !found {
			return 0
		}

		if ti < len(targetOrigRunes) && qi < len(queryRunes) {
			if targetOrigRunes[ti] == queryRunes[qi] {
				exactCaseCount++
			}
		}

		if qi == 0 && ti == 0 {
			score += 12
		}

		if ti > 0 && isBoundary(rune(targetLower[ti-1])) {
			score += 10
		}

		if qi > 0 && consecutive > 0 {
			consecutive++
			score += 8
		} else {
			consecutive = 1
		}

		qi++
		ti++
	}

	score += float64(exactCaseCount) * 2
	score -= float64(gapCount)

	score += float64(len(queryLowerRunes))

	return score
}

func pathRune(slice [][]rune, i int) []rune {
	if i < len(slice) {
		return slice[i]
	}
	return nil
}

func isBoundary(r rune) bool {
	return r == '/' || r == '-' || r == '_' || r == ' ' || r == '.'
}

// FuzzySearch performs fuzzy matching on file paths.
// pathsLower must contain strings.ToLower for each path.
// pathsRunes and pathsLowerRunes are pre-computed []rune slices (optional — pass nil to compute on the fly).
func FuzzySearch(query string, paths, pathsLower []string, pathsRunes, pathsLowerRunes [][]rune) []Result {
	if query == "" {
		n := maxSearchResults
		if len(paths) < n {
			n = len(paths)
		}
		results := make([]Result, n)
		for i := 0; i < n; i++ {
			results[i] = Result{Path: paths[i], Score: 0}
		}
		if n > 1 {
			sort.Slice(results, func(i, j int) bool {
				return pathsLower[i] < pathsLower[j]
			})
		}
		return results
	}

	queryLower := strings.ToLower(query)
	queryRunes := []rune(query)
	queryLowerRunes := []rune(queryLower)
	results := make([]Result, 0, maxSearchResults)
	for i, path := range paths {
		targetOrigRunes := pathRune(pathsRunes, i)
		targetLowerRunes := pathRune(pathsLowerRunes, i)
		if targetOrigRunes == nil {
			targetOrigRunes = []rune(path)
		}
		if targetLowerRunes == nil {
			targetLowerRunes = []rune(pathsLower[i])
		}
		score := FuzzyScore(queryRunes, queryLowerRunes, targetOrigRunes, targetLowerRunes, pathsLower[i])
		if score > 0 {
			results = append(results, Result{Path: path, Score: score})
			if len(results) >= maxSearchResults {
				break
			}
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	if len(results) > maxSearchResults {
		results = results[:maxSearchResults]
	}
	return results
}

// HighlightMatches highlights query characters in a path string.
func HighlightMatches(query, path string) string {
	if query == "" {
		return path
	}

	queryLower := strings.ToLower(query)
	pathLower := strings.ToLower(path)

	var result strings.Builder
	qi := 0
	for i, r := range path {
		if qi < len(queryLower) && i < len(pathLower) && unicode.ToLower(r) == rune(queryLower[qi]) {
			result.WriteString(lipgloss.NewStyle().Bold(true).Underline(true).Render(string(r)))
			qi++
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// ContentSearch performs full-text search across file contents.
func ContentSearch(query string, index map[string]string) []Result {
	queryLower := strings.ToLower(query)
	var results []Result

	for path, content := range index {
		lowerContent := strings.ToLower(content)
		if !strings.Contains(lowerContent, queryLower) {
			continue
		}

		remaining := content
		lineNum := 0
		for remaining != "" {
			var line string
			if idx := strings.Index(remaining, "\n"); idx >= 0 {
				line = remaining[:idx]
				remaining = remaining[idx+1:]
			} else {
				line = remaining
				remaining = ""
			}
			lineNum++

			if strings.Contains(strings.ToLower(line), queryLower) {
				trimmed := strings.TrimSpace(line)
				if len(trimmed) > contentResultContextLen {
					trimmed = trimmed[:contentResultContextLen] + "..."
				}
				results = append(results, Result{
					Path:    path,
					Context: trimmed,
					LineNum: lineNum,
				})
				if len(results) >= maxContentResults {
					return results
				}
			}
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return strings.ToLower(results[i].Path) < strings.ToLower(results[j].Path)
	})

	return results
}

// RenderResults renders search results for display in the TUI.
func RenderResults(state State, width int, style Style) string {
	if state.query == "" && state.mode == Name {
		return renderFileList(state, width, style)
	}
	if state.query == "" && state.mode == Content {
		return lipgloss.NewStyle().Foreground(style.TextMuted).Render("Type to search note contents...")
	}
	if len(state.results) == 0 {
		return lipgloss.NewStyle().Foreground(style.TextMuted).Render("No results")
	}

	var sb strings.Builder
	for i, r := range state.results {
		if i > 0 {
			sb.WriteString("\n")
		}

		line := formatResult(r, state.mode, i == state.selected, width, style)
		sb.WriteString(line)
	}
	return sb.String()
}

func renderFileList(state State, width int, style Style) string {
	var sb strings.Builder
	for i, r := range state.results {
		if i > 0 {
			sb.WriteString("\n")
		}

		var line string
		if state.mode == Name {
			highlighted := HighlightMatches(state.query, r.Path)
			line = fmt.Sprintf("  %s", highlighted)
		} else {
			line = fmt.Sprintf("  %s", r.Path)
		}

		if i == state.selected {
			line = styleSelected(line, style)
		} else {
			line = lipgloss.NewStyle().Foreground(style.TextSecondary).Render(line)
		}

		sb.WriteString(line)
	}
	return sb.String()
}

func formatResult(r Result, mode Mode, selected bool, width int, style Style) string {
	var line string
	if mode == Name {
		line = fmt.Sprintf("  %s  (%.0f)", r.Path, r.Score)
	} else {
		line = fmt.Sprintf("  %s:%d  %s", r.Path, r.LineNum, r.Context)
	}

	if selected {
		return styleSelected(line, style)
	}
	return lipgloss.NewStyle().Foreground(style.TextSecondary).Render(line)
}

func styleSelected(line string, style Style) string {
	return lipgloss.NewStyle().Background(style.Accent).Foreground(style.SelectionText).Bold(true).Render(line)
}
