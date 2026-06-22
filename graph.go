package main

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
)

// GraphScope determines whether the graph shows all notes or a local neighborhood.
type GraphScope int

const (
	GraphGlobal GraphScope = iota
	GraphLocal
)

// GraphNode represents a single note in the graph visualization.
type GraphNode struct {
	Path    string
	Label   string
	OutDeg  int
	InDeg   int
	ScreenX int
	ScreenY int
}

// GraphEdge represents a directed link from one note to another.
type GraphEdge struct {
	From string
	To   string
}

// GraphModel holds the state for the graph view mode.
type GraphModel struct {
	scope      GraphScope
	nodes      []GraphNode
	edges      []GraphEdge
	selected   int
	centerPath string
}

// BuildEdges constructs directed edges from wiki-links found in note contents.
func BuildEdges(allPaths []string, searchIndex map[string]string, vault *VaultEntry, vaultPath string) []GraphEdge {
	var edges []GraphEdge
	for _, path := range allPaths {
		content, ok := searchIndex[path]
		if !ok {
			continue
		}
		for _, target := range extractWikiLinkTargets(content) {
			resolved := ResolveWikiLink(target, vault, vaultPath)
			if resolved != "" && resolved != path {
				edges = append(edges, GraphEdge{From: path, To: resolved})
			}
		}
	}
	return edges
}

// BuildLocalGraph builds a graph centered on centerPath with its out-neighbors and in-neighbors.
func BuildLocalGraph(center string, edges []GraphEdge, backlinks map[string][]string, maxNodes int) GraphModel {
	labelMap := labelMap(edges, backlinks, center)
	nodeSet := make(map[string]bool)
	nodeSet[center] = true

	outNeighbors := make(map[string]bool)
	for _, e := range edges {
		if e.From == center {
			outNeighbors[e.To] = true
		}
	}
	inNeighbors := make(map[string]bool)
	for _, p := range backlinks[center] {
		inNeighbors[p] = true
	}

	for p := range outNeighbors {
		if len(nodeSet) < maxNodes {
			nodeSet[p] = true
		}
	}
	for p := range inNeighbors {
		if len(nodeSet) < maxNodes {
			nodeSet[p] = true
		}
	}

	var nodes []GraphNode
	for p := range nodeSet {
		nodes = append(nodes, GraphNode{Path: p, Label: labelMap[p]})
	}

	var localEdges []GraphEdge
	for _, e := range edges {
		if nodeSet[e.From] && nodeSet[e.To] {
			localEdges = append(localEdges, e)
		}
	}

	return GraphModel{
		scope:      GraphLocal,
		nodes:      nodes,
		edges:      localEdges,
		selected:   0,
		centerPath: center,
	}
}

// BuildGlobalGraph builds a graph of the top-connected notes.
func BuildGlobalGraph(edges []GraphEdge, allPaths []string, backlinks map[string][]string, maxNodes int) GraphModel {
	deg := degreeMap(edges, backlinks, allPaths)
	labelMap := make(map[string]string)
	for _, p := range allPaths {
		labelMap[p] = basename(p)
	}

	type pathDeg struct {
		path  string
		deg   int
	}
	var sorted []pathDeg
	for p, d := range deg {
		sorted = append(sorted, pathDeg{p, d})
	}
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].deg > sorted[j].deg })

	if len(sorted) > maxNodes {
		sorted = sorted[:maxNodes]
	}

	nodeSet := make(map[string]bool)
	var nodes []GraphNode
	for _, pd := range sorted {
		nodeSet[pd.path] = true
		nodes = append(nodes, GraphNode{Path: pd.path, Label: labelMap[pd.path], OutDeg: degOut(pd.path, edges), InDeg: degIn(pd.path, backlinks)})
	}

	var globalEdges []GraphEdge
	for _, e := range edges {
		if nodeSet[e.From] && nodeSet[e.To] {
			globalEdges = append(globalEdges, e)
		}
	}

	return GraphModel{
		scope:    GraphGlobal,
		nodes:    nodes,
		edges:    globalEdges,
		selected: 0,
	}
}

// labelMap builds a map of path → basename for all nodes in edges + backlinks + an explicit center.
func labelMap(edges []GraphEdge, backlinks map[string][]string, center string) map[string]string {
	m := make(map[string]string)
	m[center] = basename(center)
	for _, e := range edges {
		m[e.From] = basename(e.From)
		m[e.To] = basename(e.To)
	}
	for from, tos := range backlinks {
		m[from] = basename(from)
		for _, to := range tos {
			m[to] = basename(to)
		}
	}
	return m
}

// degreeMap builds a map of path → total-degree (out+in) for all paths.
func degreeMap(edges []GraphEdge, backlinks map[string][]string, allPaths []string) map[string]int {
	deg := make(map[string]int)
	for _, p := range allPaths {
		deg[p] = 0
	}
	for _, e := range edges {
		deg[e.From]++
	}
	for p, tos := range backlinks {
		deg[p] += len(tos)
	}
	return deg
}

func degOut(path string, edges []GraphEdge) int {
	c := 0
	for _, e := range edges {
		if e.From == path {
			c++
		}
	}
	return c
}

func degIn(path string, backlinks map[string][]string) int {
	return len(backlinks[path])
}

// basename strips the directory and .md extension from a path.
func basename(path string) string {
	idx := strings.LastIndexByte(path, '/')
	name := path
	if idx >= 0 {
		name = path[idx+1:]
	}
	if strings.HasSuffix(name, ".md") {
		name = name[:len(name)-3]
	}
	return name
}

func truncateLabel(label string, maxLen int) string {
	if len(label) <= maxLen {
		return label
	}
	return label[:maxLen-1] + "\u2026"
}

// LayoutCircle positions nodes in a circle on a width×height grid.
func LayoutCircle(nodes []GraphNode, width, height int) {
	if len(nodes) == 0 {
		return
	}
	cx := width / 2
	cy := height / 2
	r := min(cx, cy) * 2 / 3
	if r < 1 {
		r = 1
	}
	for i := range nodes {
		angle := 2 * math.Pi * float64(i) / float64(len(nodes))
		nodes[i].ScreenX = cx + int(float64(r)*math.Cos(angle))
		nodes[i].ScreenY = cy + int(float64(r)*math.Sin(angle))
	}
}

// RenderGrid renders the graph as an ASCII string on a width×height grid.
func RenderGrid(nodes []GraphNode, edges []GraphEdge, selected int, width, height int, palette Palette) string {
	if len(nodes) == 0 {
		return ""
	}
	grid := make([][]rune, height)
	for y := range grid {
		grid[y] = make([]rune, width)
		for x := range grid[y] {
			grid[y][x] = ' '
		}
	}

	// Draw edges first (underneath labels)
	for _, e := range edges {
		var from, to *GraphNode
		for i := range nodes {
			if nodes[i].Path == e.From {
				from = &nodes[i]
			}
			if nodes[i].Path == e.To {
				to = &nodes[i]
			}
		}
		if from == nil || to == nil {
			continue
		}
		drawLine(grid, from.ScreenX, from.ScreenY, to.ScreenX, to.ScreenY)
	}

	// Draw node labels
	for i, node := range nodes {
		label := truncateLabel(node.Label, 12)
		x := node.ScreenX - len(label)/2
		y := node.ScreenY
		if selected == i {
			label = "[" + label + "]"
			x = node.ScreenX - len(label)/2
		}
		placeLabel(grid, label, x, y, i == selected, palette)
	}

	var sb strings.Builder
	for y := range grid {
		sb.WriteString(strings.TrimRight(string(grid[y]), " "))
		if y < len(grid)-1 {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func drawLine(grid [][]rune, x1, y1, x2, y2 int) {
	if y1 < 0 || y1 >= len(grid) || y2 < 0 || y2 >= len(grid) {
		return
	}
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	sx := 1
	sy := 1
	if x1 > x2 {
		sx = -1
	}
	if y1 > y2 {
		sy = -1
	}
	err := dx - dy
	steps := 0
	for {
		if x1 >= 0 && x1 < len(grid[0]) && y1 >= 0 && y1 < len(grid) {
			r := grid[y1][x1]
			if r == ' ' {
				if dy > dx {
					grid[y1][x1] = '│'
				} else {
					grid[y1][x1] = '─'
				}
			}
		}
		if x1 == x2 && y1 == y2 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x1 += sx
		}
		if e2 < dx {
			err += dx
			y1 += sy
		}
		steps++
		if steps > 2000 {
			break
		}
	}
}

func placeLabel(grid [][]rune, label string, x, y int, accented bool, palette Palette) {
	if y < 0 || y >= len(grid) {
		return
	}
	for i, r := range label {
		px := x + i
		if px < 0 || px >= len(grid[0]) {
			continue
		}
		grid[y][px] = r
	}
}

// enterGraphMode builds the graph and switches to ModeGraph.
func (m *Model) enterGraphMode() {
	m.prevMode = m.mode
	edges := BuildEdges(m.allPaths, m.searchIndex, m.vault, m.config.VaultPath)
	if m.mode == ModeView && m.activeNote != nil {
		m.graph = BuildLocalGraph(m.activeNote.Path, edges, m.backlinkIndex, 30)
	} else {
		m.graph = BuildGlobalGraph(edges, m.allPaths, m.backlinkIndex, 60)
		if len(m.allPaths) > 60 {
			qualifying := 0
			for _, e := range edges {
				if e.From != "" && e.To != "" {
					qualifying++
				}
			}
			if qualifying > 60 {
				m.addToast("Showing top 60 connected notes", ToastInfo)
			}
		}
	}
	m.mode = ModeGraph
}

func (m *Model) toggleGraphScope() {
	edges := BuildEdges(m.allPaths, m.searchIndex, m.vault, m.config.VaultPath)
	if m.graph.scope == GraphLocal {
		m.graph = BuildGlobalGraph(edges, m.allPaths, m.backlinkIndex, 60)
		m.addToast("Switched to global graph", ToastInfo)
	} else {
		center := m.graph.nodes[m.graph.selected].Path
		m.graph = BuildLocalGraph(center, edges, m.backlinkIndex, 30)
		m.addToast("Switched to local graph", ToastInfo)
	}
}

func (m *Model) focusGraphOn(path string) {
	edges := BuildEdges(m.allPaths, m.searchIndex, m.vault, m.config.VaultPath)
	m.graph = BuildLocalGraph(path, edges, m.backlinkIndex, 30)
	m.addToast("Focused on: "+basename(path), ToastInfo)
}

func (m Model) handleGraphKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case msg.Type == tea.KeyEsc:
		m.mode = m.prevMode
		return m, nil
	case msg.Type == tea.KeyEnter:
		if m.graph.selected < len(m.graph.nodes) {
			m.openNote(m.graph.nodes[m.graph.selected].Path)
		}
		return m, nil
	case m.keys.MatchDown(msg):
		if m.graph.selected < len(m.graph.nodes)-1 {
			m.graph.selected++
		}
		return m, nil
	case m.keys.MatchUp(msg):
		if m.graph.selected > 0 {
			m.graph.selected--
		}
		return m, nil
	case MatchRune(msg, 'l'):
		m.toggleGraphScope()
		return m, nil
	case MatchRune(msg, 'f'):
		if m.graph.selected < len(m.graph.nodes) {
			m.focusGraphOn(m.graph.nodes[m.graph.selected].Path)
		}
		return m, nil
	case MatchRune(msg, 'r'):
		m.enterGraphMode()
		return m, nil
	}
	return m, nil
}

func (m Model) renderGraph() string {
	if len(m.graph.nodes) == 0 {
		msg := lipgloss.NewStyle().Foreground(m.palette.TextDim).Render("No wiki-link connections found in this vault")
		return lipgloss.NewStyle().Width(m.width - m.treeWidth - 2).Align(lipgloss.Center).Render(msg)
	}
	width := m.width - m.treeWidth - 2
	height := m.height - 2
	if width < 20 {
		width = 20
	}
	if height < 5 {
		height = 5
	}

	LayoutCircle(m.graph.nodes, width, height)
	grid := RenderGrid(m.graph.nodes, m.graph.edges, m.graph.selected, width, height, m.palette)

	footer := lipgloss.NewStyle().Foreground(m.palette.TextDim).Render(
		"j/k move · Enter open · l scope · f focus · r refresh · Esc back",
	)

	scopeLabel := "global"
	if m.graph.scope == GraphLocal {
		scopeLabel = "local: " + truncateLabel(m.graph.centerPath, 30)
	}
	status := lipgloss.NewStyle().Foreground(m.palette.Accent).Render(
		fmt.Sprintf("Graph (%s) — %d nodes", scopeLabel, len(m.graph.nodes)),
	)

	return lipgloss.JoinVertical(lipgloss.Left, status, "", grid, "", footer)
}
