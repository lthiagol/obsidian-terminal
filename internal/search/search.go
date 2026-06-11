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
	Name Mode = iota
	Content
)

// Result represents a single search match.
type Result struct {
	Path    string
	Score   float64
	Context string
	LineNum int
}

// State holds the current state of a search session.
type State struct {
	mode        Mode
	query       string
	results     []Result
	selected    int
	allPaths    []string
	searchIndex map[string]string
}

// Style holds colors for search result rendering.
type Style struct {
	Accent        lipgloss.Color
	TextSecondary lipgloss.Color
	TextMuted     lipgloss.Color
}

// NewState creates a new search state with the given mode and data.
func NewState(mode Mode, paths []string, index map[string]string) State {
	s := State{
		mode:        mode,
		allPaths:    paths,
		searchIndex: index,
		selected:    0,
	}
	if mode == Name {
		s.results = FuzzySearch("", paths)
	}
	return s
}

func (s *State) SetQuery(query string) {
	s.query = query
	s.selected = 0

	switch s.mode {
	case Name:
		s.results = FuzzySearch(query, s.allPaths)
	case Content:
		if query == "" {
			s.results = []Result{}
		} else {
			s.results = ContentSearch(query, s.searchIndex)
		}
	}
}

func (s *State) MoveUp() {
	if s.selected > 0 {
		s.selected--
	}
}

func (s *State) MoveDown() {
	if s.selected < len(s.results)-1 {
		s.selected++
	}
}

func (s State) ResultCount() int {
	return len(s.results)
}

func (s State) Query() string {
	return s.query
}

func (s State) SelectedResult() *Result {
	if s.selected >= 0 && s.selected < len(s.results) {
		return &s.results[s.selected]
	}
	return nil
}

// FuzzyScore computes a match score between query and target.
func FuzzyScore(query, target string) float64 {
	queryLower := strings.ToLower(query)
	targetLower := strings.ToLower(target)

	if queryLower == "" {
		return 0
	}

	if queryLower == targetLower {
		return 100 + float64(len(queryLower))
	}

	queryRunes := []rune(query)
	queryLowerRunes := []rune(queryLower)
	targetLowerRunes := []rune(targetLower)
	targetOrigRunes := []rune(target)

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

func isBoundary(r rune) bool {
	return r == '/' || r == '-' || r == '_' || r == ' ' || r == '.'
}

// FuzzySearch performs fuzzy matching on file paths.
func FuzzySearch(query string, paths []string) []Result {
	if query == "" {
		results := make([]Result, len(paths))
		for i, path := range paths {
			results[i] = Result{Path: path, Score: 0}
		}
		sort.Slice(results, func(i, j int) bool {
			return strings.ToLower(results[i].Path) < strings.ToLower(results[j].Path)
		})
		if len(results) > 50 {
			results = results[:50]
		}
		return results
	}

	var results []Result
	for _, path := range paths {
		score := FuzzyScore(query, path)
		if score > 0 {
			results = append(results, Result{Path: path, Score: score})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	if len(results) > 50 {
		results = results[:50]
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
		contentLower := strings.ToLower(content)
		if strings.Contains(contentLower, queryLower) {
			lines := strings.Split(content, "\n")
			for lineNum, line := range lines {
				if strings.Contains(strings.ToLower(line), queryLower) {
					trimmed := strings.TrimSpace(line)
					if len(trimmed) > 80 {
						trimmed = trimmed[:80] + "..."
					}
					results = append(results, Result{
						Path:    path,
						Context: trimmed,
						LineNum: lineNum + 1,
					})
					if len(results) >= 100 {
						break
					}
				}
			}
		}
		if len(results) >= 100 {
			break
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
			line = lipgloss.NewStyle().Background(style.Accent).Foreground(lipgloss.Color("#000000")).Bold(true).Render(line)
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
		return lipgloss.NewStyle().Background(style.Accent).Foreground(lipgloss.Color("#000000")).Bold(true).Render(line)
	}
	return lipgloss.NewStyle().Foreground(style.TextSecondary).Render(line)
}
