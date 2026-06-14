package ansiext

import (
	"strings"
	"testing"
)

func TestUndercurl_NonEmpty(t *testing.T) {
	out := Undercurl("hello")
	if !strings.Contains(out, "\033[4:3m") {
		t.Error("undercurl should contain 4:3 SGR")
	}
	if !strings.Contains(out, "\033[4:0m") {
		t.Error("undercurl should close with 4:0 SGR")
	}
	if !strings.Contains(out, "hello") {
		t.Error("undercurl should contain the text")
	}
}

func TestUndercurl_Empty(t *testing.T) {
	out := Undercurl("")
	if out != "" {
		t.Errorf("empty undercurl should be empty, got %q", out)
	}
}

func TestOverline_NonEmpty(t *testing.T) {
	out := Overline("world")
	if !strings.Contains(out, "\033[53m") {
		t.Error("overline should contain 53m SGR")
	}
	if !strings.Contains(out, "\033[55m") {
		t.Error("overline should close with 55m SGR")
	}
	if !strings.Contains(out, "world") {
		t.Error("overline should contain the text")
	}
}

func TestOverline_Empty(t *testing.T) {
	out := Overline("")
	if out != "" {
		t.Errorf("empty overline should be empty, got %q", out)
	}
}
