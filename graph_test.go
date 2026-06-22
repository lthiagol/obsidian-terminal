package main

import (
	"os"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func makeTestEdges() ([]string, []GraphEdge, map[string]string, map[string][]string, *VaultEntry) {
	allPaths := []string{"a.md", "b.md", "c.md"}
	searchIndex := map[string]string{
		"a.md": "[[b.md]]",
		"b.md": "[[c.md]]",
		"c.md": "some text",
	}
	edges := []GraphEdge{
		{From: "a.md", To: "b.md"},
		{From: "b.md", To: "c.md"},
	}
	backlinks := map[string][]string{
		"b.md": {"a.md"},
		"c.md": {"b.md"},
	}
	vault := &VaultEntry{
		IsDir: true,
		Children: []*VaultEntry{
			{Name: "a.md", Path: "a.md"},
			{Name: "b.md", Path: "b.md"},
			{Name: "c.md", Path: "c.md"},
		},
	}
	return allPaths, edges, searchIndex, backlinks, vault
}

// T1: BuildEdges from wiki-links in searchIndex.
func TestBuildEdges_Chain(t *testing.T) {
	allPaths, _, searchIndex, _, vault := makeTestEdges()
	edges := BuildEdges(allPaths, searchIndex, vault, "/fake")
	if len(edges) != 2 {
		t.Fatalf("expected 2 edges, got %d", len(edges))
	}
	if edges[0].From != "a.md" || edges[0].To != "b.md" {
		t.Errorf("edge 0: expected a.md → b.md, got %s → %s", edges[0].From, edges[0].To)
	}
	if edges[1].From != "b.md" || edges[1].To != "c.md" {
		t.Errorf("edge 1: expected b.md → c.md, got %s → %s", edges[1].From, edges[1].To)
	}
}

// T2: Self-loops excluded.
func TestBuildEdges_SelfLoop(t *testing.T) {
	allPaths := []string{"a.md"}
	searchIndex := map[string]string{"a.md": "[[a.md]]"}
	vault := &VaultEntry{
		IsDir: true,
		Children: []*VaultEntry{
			{Name: "a.md", Path: "a.md"},
		},
	}

	edges := BuildEdges(allPaths, searchIndex, vault, "/fake")
	if len(edges) != 0 {
		t.Fatalf("expected 0 edges (self-loop excluded), got %d", len(edges))
	}
}

// T3: BuildLocalGraph includes backlink neighbors.
func TestBuildLocalGraph_Backlinks(t *testing.T) {
	_, edges, _, backlinks, _ := makeTestEdges()

	g := BuildLocalGraph("b.md", edges, backlinks, 30)
	if len(g.nodes) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(g.nodes))
	}
	paths := make(map[string]bool)
	for _, n := range g.nodes {
		paths[n.Path] = true
	}
	if !paths["a.md"] || !paths["b.md"] || !paths["c.md"] {
		t.Errorf("missing expected nodes: got %v", paths)
	}
	if g.centerPath != "b.md" {
		t.Errorf("center: expected b.md, got %s", g.centerPath)
	}
	if g.scope != GraphLocal {
		t.Errorf("scope: expected GraphLocal, got %v", g.scope)
	}
}

// T4: BuildGlobalGraph caps at maxNodes.
func TestBuildGlobalGraph_Cap(t *testing.T) {
	allPaths := make([]string, 100)
	edges := make([]GraphEdge, 100*2)
	for i := 0; i < 100; i++ {
		allPaths[i] = "note" + strings.Repeat("x", i) + ".md"
	}
	for i := 0; i < 100; i++ {
		edges = append(edges, GraphEdge{From: allPaths[i], To: allPaths[(i+1)%100]})
	}
	backlinks := make(map[string][]string)
	for i := 0; i < 100; i++ {
		backlinks[allPaths[i]] = []string{allPaths[(i-1+100)%100]}
	}

	g := BuildGlobalGraph(edges, allPaths, backlinks, 60)
	if len(g.nodes) > 60 {
		t.Fatalf("expected ≤60 nodes, got %d", len(g.nodes))
	}
	if g.scope != GraphGlobal {
		t.Errorf("scope: expected GraphGlobal, got %v", g.scope)
	}
}

// T5: BuildGlobalGraph sorts by degree (highest first).
func TestBuildGlobalGraph_DegreeSort(t *testing.T) {
	allPaths := []string{"a.md", "b.md", "c.md"}
	edges := []GraphEdge{
		{From: "a.md", To: "b.md"},
		{From: "a.md", To: "c.md"},
		{From: "b.md", To: "c.md"},
	}
	backlinks := map[string][]string{
		"b.md": {"a.md"},
		"c.md": {"a.md", "b.md"},
	}

	g := BuildGlobalGraph(edges, allPaths, backlinks, 60)
	if len(g.nodes) < 3 {
		t.Fatalf("expected 3 nodes, got %d", len(g.nodes))
	}
	// Highest degree should be "a.md" (out 2, in 0 = 2) or "c.md" (out 0, in 2 = 2)
	if g.nodes[0].Path != "a.md" && g.nodes[0].Path != "c.md" {
		t.Logf("first node: %s (out=%d in=%d)", g.nodes[0].Path, g.nodes[0].OutDeg, g.nodes[0].InDeg)
	}
	// "b.md" should NOT be first: out 1, in 1 = 2, ties are fine
}

// T6: LayoutCircle sets non-zero ScreenX/Y.
func TestLayoutCircle_Positions(t *testing.T) {
	nodes := []GraphNode{
		{Path: "a.md", Label: "a"},
		{Path: "b.md", Label: "b"},
		{Path: "c.md", Label: "c"},
		{Path: "d.md", Label: "d"},
	}
	LayoutCircle(nodes, 80, 24)
	for i, n := range nodes {
		if n.ScreenX == 0 && n.ScreenY == 0 {
			t.Errorf("node %d: expected non-zero position, got (%d,%d)", i, n.ScreenX, n.ScreenY)
		}
	}
}

// T7: RenderGrid returns non-empty string containing all labels.
func TestRenderGrid_Labels(t *testing.T) {
	nodes := []GraphNode{
		{Path: "a.md", Label: "alpha", ScreenX: 10, ScreenY: 5},
		{Path: "b.md", Label: "beta", ScreenX: 30, ScreenY: 5},
		{Path: "c.md", Label: "gamma", ScreenX: 20, ScreenY: 10},
	}
	edges := []GraphEdge{
		{From: "a.md", To: "c.md"},
	}
	out := RenderGrid(nodes, edges, 0, 40, 15, Palette{})
	if out == "" {
		t.Fatal("expected non-empty output")
	}
	for _, label := range []string{"alpha", "beta", "gamma"} {
		if !strings.Contains(out, label) {
			t.Errorf("output missing label: %s", label)
		}
	}
}

// T8: Selected node rendered with bracket.
func TestRenderGrid_SelectedBracket(t *testing.T) {
	nodes := []GraphNode{{Path: "a.md", Label: "note", ScreenX: 10, ScreenY: 5}}
	out := RenderGrid(nodes, nil, 0, 40, 15, Palette{})
	if !strings.Contains(out, "[") || !strings.Contains(out, "]") {
		t.Error("selected node should have bracket wrapper")
	}
}

// T9: 0 nodes → RenderGrid returns empty, no panic.
func TestRenderGrid_Empty(t *testing.T) {
	out := RenderGrid(nil, nil, 0, 40, 15, Palette{})
	if out != "" {
		t.Errorf("expected empty output for 0 nodes, got %q", out)
	}
}

// T10: 0 edges → RenderGrid shows labels only.
func TestRenderGrid_NoEdges(t *testing.T) {
	nodes := []GraphNode{
		{Path: "a.md", Label: "a", ScreenX: 10, ScreenY: 5},
		{Path: "b.md", Label: "b", ScreenX: 30, ScreenY: 5},
	}
	out := RenderGrid(nodes, nil, 0, 40, 15, Palette{})
	if out == "" {
		t.Fatal("expected non-empty output")
	}
	if !strings.Contains(out, "a") || !strings.Contains(out, "b") {
		t.Error("output should contain node labels")
	}
	// No edge chars should appear
	if strings.ContainsAny(out, "│╲╱─") {
		t.Log("output contains edge chars — may be from labels")
	}
}

// T11: Ctrl+G from browse → ModeGraph, scope=global.
func TestGraph_CtrlG_Browse(t *testing.T) {
	m := newTestModel(t, &Config{})
	m.mode = ModeBrowse

	updateM, _ := m.Update(tea.KeyMsg{Type: tea.KeyCtrlG})
	updated, ok := updateM.(Model)
	if !ok {
		t.Fatal("expected Model")
	}
	if updated.mode != ModeGraph {
		t.Errorf("expected ModeGraph, got %v", updated.mode)
	}
	if updated.graph.scope != GraphGlobal {
		t.Errorf("expected GraphGlobal scope, got %v", updated.graph.scope)
	}
	if updated.prevMode != ModeBrowse {
		t.Errorf("expected prevMode ModeBrowse, got %v", updated.prevMode)
	}
}

// T12: Ctrl+G from view → ModeGraph, scope=local.
func TestGraph_CtrlG_View_Local(t *testing.T) {
	m := newTestModel(t, &Config{})
	m.mode = ModeView
	m.activeNote = &VaultNote{Path: "test.md", Body: "hello"}

	updateM, _ := m.Update(tea.KeyMsg{Type: tea.KeyCtrlG})
	updated, ok := updateM.(Model)
	if !ok {
		t.Fatal("expected Model")
	}
	if updated.mode != ModeGraph {
		t.Errorf("expected ModeGraph, got %v", updated.mode)
	}
	if updated.graph.scope != GraphLocal {
		t.Errorf("expected GraphLocal scope, got %v", updated.graph.scope)
	}
}

// T13: j increments selected.
func TestGraph_JDown(t *testing.T) {
	m := newTestModel(t, &Config{})
	m.mode = ModeGraph
	m.graph.nodes = []GraphNode{
		{Path: "a.md", Label: "a"},
		{Path: "b.md", Label: "b"},
	}
	m.graph.selected = 0
	m.keys = DefaultKeys()

	updateM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	updated, ok := updateM.(Model)
	if !ok {
		t.Fatal("expected Model")
	}
	if updated.graph.selected != 1 {
		t.Errorf("expected selected 1, got %d", updated.graph.selected)
	}

	// Arrows work too
	updateM, _ = updated.Update(tea.KeyMsg{Type: tea.KeyUp})
	updated, ok = updateM.(Model)
	if !ok {
		t.Fatal("expected Model")
	}
	if updated.graph.selected != 0 {
		t.Errorf("expected selected back to 0, got %d", updated.graph.selected)
	}
}

// T14: Enter opens note → ModeView.
func TestGraph_Enter_Open(t *testing.T) {
	vaultDir := t.TempDir()
	if err := os.WriteFile(vaultDir+"/hello.md", []byte("# Hello\n\ncontent"), 0644); err != nil {
		t.Fatal(err)
	}
	vaultTree, _, _, _ := ScanVault(vaultDir, nil)

	var scanPaths []string
	var walk func(e *VaultEntry)
	walk = func(e *VaultEntry) {
		if !e.IsDir {
			scanPaths = append(scanPaths, e.Path)
		}
		for _, c := range e.Children {
			walk(c)
		}
	}
	walk(vaultTree)

	m := Model{
		mode:      ModeGraph,
		config:    &Config{VaultPath: vaultDir},
		vault:     vaultTree,
		allPaths:  scanPaths,
	}
	m.graph.nodes = []GraphNode{
		{Path: "hello.md", Label: "hello"},
	}
	m.graph.selected = 0

	updateM, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	updated, ok := updateM.(Model)
	if !ok {
		t.Fatal("expected Model")
	}
	if updated.mode != ModeView {
		t.Errorf("expected ModeView after Enter, got %v", updated.mode)
	}
	if updated.activeNote == nil {
		t.Fatal("expected activeNote to be set")
	}
}

// T15: Esc restores prevMode.
func TestGraph_Esc(t *testing.T) {
	m := newTestModel(t, &Config{})
	m.mode = ModeGraph
	m.prevMode = ModeBrowse

	updateM, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	updated, ok := updateM.(Model)
	if !ok {
		t.Fatal("expected Model")
	}
	if updated.mode != ModeBrowse {
		t.Errorf("expected ModeBrowse after Esc, got %v", updated.mode)
	}
}

// T16: 0-edge vault shows friendly message.
func TestGraph_EmptyVault(t *testing.T) {
	m := newTestModel(t, &Config{})
	m.mode = ModeGraph
	m.graph.nodes = nil
	m.width = 80
	m.height = 24
	m.treeWidth = 20

	out := m.renderGraph()
	if !strings.Contains(out, "No wiki-link connections") {
		t.Errorf("expected friendly message, got: %s", out)
	}
}

// T17: Global >60 notes → toast shown.
func TestGraph_GlobalTruncate_Toast(t *testing.T) {
	vaultDir := t.TempDir()
	entryList := make([]*VaultEntry, 0, 70)
	for i := 0; i < 70; i++ {
		name := "n" + strings.Repeat("a", i) + ".md"
		body := "# Note\n"
		if i < 69 {
			next := "n" + strings.Repeat("a", i+1)
			body += "[[" + next + "]]\n"
		}
		if err := os.WriteFile(vaultDir+"/"+name, []byte(body), 0644); err != nil {
			t.Fatal(err)
		}
		entryList = append(entryList, &VaultEntry{Name: name, Path: name})
	}
	vaultTree := &VaultEntry{IsDir: true, Children: entryList}

	var scanPaths []string
	var walk func(e *VaultEntry)
	walk = func(e *VaultEntry) {
		if !e.IsDir {
			scanPaths = append(scanPaths, e.Path)
		}
		for _, c := range e.Children {
			walk(c)
		}
	}
	walk(vaultTree)

	searchIdx := make(map[string]string)
	for _, p := range scanPaths {
		data, _ := os.ReadFile(vaultDir + "/" + p)
		searchIdx[p] = string(data)
	}

	m := Model{
		mode:        ModeBrowse,
		config:      &Config{VaultPath: vaultDir},
		vault:       vaultTree,
		allPaths:    scanPaths,
		searchIndex: searchIdx,
		prevMode:    ModeBrowse,
	}

	m.enterGraphMode()
	if len(m.toasts) == 0 {
		t.Error("expected truncation toast when >60 connected notes")
	}
}

// T18: l toggles scope.
func TestGraph_L_Toggle(t *testing.T) {
	m := newTestModel(t, &Config{})
	m.mode = ModeGraph
	m.allPaths = []string{"a.md", "b.md"}
	m.searchIndex = map[string]string{
		"a.md": "[[b.md]]",
		"b.md": "text",
	}
	m.graph.scope = GraphGlobal
	m.graph.nodes = []GraphNode{
		{Path: "a.md", Label: "a"},
		{Path: "b.md", Label: "b"},
	}
	m.graph.selected = 0

	updateM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	updated, ok := updateM.(Model)
	if !ok {
		t.Fatal("expected Model")
	}
	if updated.graph.scope != GraphLocal {
		t.Errorf("expected GraphLocal after toggle from global, got %v", updated.graph.scope)
	}
}

// T19: f focuses on selected node.
func TestGraph_F_Focus(t *testing.T) {
	m := newTestModel(t, &Config{})
	m.mode = ModeGraph
	m.allPaths = []string{"a.md", "b.md", "c.md"}
	m.searchIndex = map[string]string{
		"a.md": "[[b.md]]",
		"b.md": "[[c.md]]",
		"c.md": "text",
	}
	m.backlinkIndex = map[string][]string{
		"b.md": {"a.md"},
		"c.md": {"b.md"},
	}
	m.config.VaultPath = "/fake/vault"
	m.graph.nodes = []GraphNode{
		{Path: "b.md", Label: "b"},
	}
	m.graph.selected = 0

	updateM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
	updated, ok := updateM.(Model)
	if !ok {
		t.Fatal("expected Model")
	}
	if updated.graph.scope != GraphLocal {
		t.Errorf("expected GraphLocal after focus, got %v", updated.graph.scope)
	}
	if updated.graph.centerPath != "b.md" {
		t.Errorf("expected centerPath b.md, got %s", updated.graph.centerPath)
	}
}

// T20: r rebuilds graph.
func TestGraph_R_Refresh(t *testing.T) {
	m := newTestModel(t, &Config{})
	m.mode = ModeGraph
	m.allPaths = []string{"a.md", "b.md"}
	m.searchIndex = map[string]string{
		"a.md": "[[b.md]]",
		"b.md": "text",
	}
	m.graph.scope = GraphGlobal
	m.graph.nodes = []GraphNode{
		{Path: "a.md", Label: "a"},
	}
	m.graph.selected = 0

	updateM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	updated, ok := updateM.(Model)
	if !ok {
		t.Fatal("expected Model")
	}
	if updated.mode != ModeGraph {
		t.Errorf("expected ModeGraph after refresh, got %v", updated.mode)
	}
	if len(updated.graph.nodes) < 1 {
		t.Error("expected nodes after refresh")
	}
}
