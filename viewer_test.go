package main

import (
	"strings"
	"testing"
)

func TestViewer_RendersMarkdown(t *testing.T) {
	v := NewViewer(defaultMarkdownStyle())
	v.SetContent("# Test\n\nHello **world**.", 80)

	view := v.View()
	if !strings.Contains(view, "Test") {
		t.Error("viewer should render 'Test' heading")
	}
	if !strings.Contains(view, "world") {
		t.Error("viewer should render 'world' text")
	}
}

func TestViewer_ScrollDown(t *testing.T) {
	v := NewViewer(defaultMarkdownStyle())
	var content strings.Builder
	for i := 0; i < 50; i++ {
		content.WriteString("Line\n")
	}
	v.SetContent(content.String(), 80)

	v.ScrollDown(5)
	if v.viewport.YOffset == 0 {
		t.Error("YOffset should be > 0 after scrolling down")
	}
}

func TestViewer_ScrollUp(t *testing.T) {
	v := NewViewer(defaultMarkdownStyle())
	var content strings.Builder
	for i := 0; i < 50; i++ {
		content.WriteString("Line\n")
	}
	v.SetContent(content.String(), 80)

	v.ScrollDown(10)
	offsetAfterDown := v.viewport.YOffset
	v.ScrollUp(5)
	if v.viewport.YOffset >= offsetAfterDown {
		t.Error("YOffset should be less after scrolling up")
	}
}

func TestViewer_ScrollToTop(t *testing.T) {
	v := NewViewer(defaultMarkdownStyle())
	var content strings.Builder
	for i := 0; i < 50; i++ {
		content.WriteString("Line\n")
	}
	v.SetContent(content.String(), 80)

	v.ScrollDown(20)
	v.ScrollTop()
	if v.viewport.YOffset != 0 {
		t.Errorf("YOffset should be 0 after ScrollTop, got %d", v.viewport.YOffset)
	}
}

func TestViewer_ScrollToBottom(t *testing.T) {
	v := NewViewer(defaultMarkdownStyle())
	var content strings.Builder
	for i := 0; i < 50; i++ {
		content.WriteString("Line\n")
	}
	v.SetContent(content.String(), 80)

	v.ScrollBottom()
	maxY := v.viewport.TotalLineCount() - v.viewport.Height
	if v.viewport.YOffset < maxY {
		t.Errorf("should be near bottom: YOffset=%d, max=%d", v.viewport.YOffset, maxY)
	}

	// Can still scroll up
	v.ScrollUp(1)
	if v.viewport.YOffset >= maxY {
		t.Error("should scroll up from bottom")
	}
}

func TestViewer_WikiLinkExtraction(t *testing.T) {
	v := NewViewer(defaultMarkdownStyle())
	v.SetContent("See [[projects/database]] and [[notes/meeting|Meeting Notes]].", 80)

	if v.LinkCount() != 2 {
		t.Fatalf("expected 2 wiki links, got %d", v.LinkCount())
	}
	if v.links[0].Target != "projects/database" {
		t.Errorf("link 0 target = %q", v.links[0].Target)
	}
	if v.links[1].Target != "notes/meeting" {
		t.Errorf("link 1 target = %q", v.links[1].Target)
	}
	if v.links[1].Display != "Meeting Notes" {
		t.Errorf("link 1 display = %q", v.links[1].Display)
	}
}

func TestViewer_CycleLinks(t *testing.T) {
	v := NewViewer(defaultMarkdownStyle())
	v.SetContent("[[a]] [[b]] [[c]]", 80)

	if v.SelectedLinkIndex() != -1 {
		t.Error("initial selected link should be -1")
	}

	v.CycleLink()
	if v.SelectedLinkIndex() != 0 {
		t.Errorf("after 1st Tab: selected = %d, want 0", v.SelectedLinkIndex())
	}
	if v.SelectedLinkPath() != "a" {
		t.Errorf("selected path = %q, want 'a'", v.SelectedLinkPath())
	}

	v.CycleLink()
	if v.SelectedLinkIndex() != 1 {
		t.Errorf("after 2nd Tab: selected = %d, want 1", v.SelectedLinkIndex())
	}

	v.CycleLink()
	if v.SelectedLinkIndex() != 2 {
		t.Errorf("after 3rd Tab: selected = %d, want 2", v.SelectedLinkIndex())
	}

	v.CycleLink()
	if v.SelectedLinkIndex() != 0 {
		t.Errorf("after 4th Tab (wrap): selected = %d, want 0", v.SelectedLinkIndex())
	}
}

func TestViewer_NoLinks(t *testing.T) {
	v := NewViewer(defaultMarkdownStyle())
	v.SetContent("Plain text with no links.", 80)

	if v.LinkCount() != 0 {
		t.Error("should have no links")
	}
	v.CycleLink() // Should not panic
	if v.SelectedLinkIndex() != -1 {
		t.Error("selected link should remain -1")
	}
}

func TestSetSize_NegativeDimensions(t *testing.T) {
	v := NewViewer(defaultMarkdownStyle())
	v.SetContent("# Test", 80)
	v.SetSize(-10, -5)
	if v.viewport.Width < 5 || v.viewport.Height < 3 {
		t.Error("SetSize should clamp negative dimensions")
	}
}

func TestViewer_SetSize(t *testing.T) {
	v := NewViewer(defaultMarkdownStyle())
	v.SetContent("# Test\n\nHello.", 80)
	v.SetSize(40, 15)

	if v.viewport.Width != 38 {
		t.Errorf("width = %d, want 38", v.viewport.Width)
	}
	if v.viewport.Height != 13 {
		t.Errorf("height = %d, want 13", v.viewport.Height)
	}
}
