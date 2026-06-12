package main

import (
	"sort"
	"testing"
)

func TestScanYAML_SimpleKeyValue(t *testing.T) {
	var pairs []yamlPair
	scanYAML([]byte("title: Hello World\n"), func(k, v string, items []string) {
		pairs = append(pairs, yamlPair{k, v, items})
	})
	if len(pairs) != 1 || pairs[0].Key != "title" || pairs[0].Value != "Hello World" {
		t.Errorf("got %+v", pairs)
	}
}

func TestScanYAML_ArrayItems(t *testing.T) {
	var pairs []yamlPair
	data := "tags:\n  - foo\n  - bar\n"
	scanYAML([]byte(data), func(k, v string, items []string) {
		pairs = append(pairs, yamlPair{k, v, items})
	})
	if len(pairs) != 1 {
		t.Fatalf("expected 1 pair, got %d", len(pairs))
	}
	if pairs[0].Key != "tags" {
		t.Errorf("key = %q", pairs[0].Key)
	}
	if len(pairs[0].Items) != 2 || pairs[0].Items[0] != "foo" || pairs[0].Items[1] != "bar" {
		t.Errorf("items = %v", pairs[0].Items)
	}
}

func TestScanYAML_QuotedStrings(t *testing.T) {
	var pairs []yamlPair
	data := `title: "Hello World"
author: 'John Doe'
`
	scanYAML([]byte(data), func(k, v string, items []string) {
		pairs = append(pairs, yamlPair{k, v, items})
	})
	if pairs[0].Value != "Hello World" {
		t.Errorf("quoted value = %q", pairs[0].Value)
	}
	if pairs[1].Value != "John Doe" {
		t.Errorf("single-quoted value = %q", pairs[1].Value)
	}
}

func TestScanYAML_InlineArray(t *testing.T) {
	var pairs []yamlPair
	scanYAML([]byte("aliases: [Alias One, Alias Two]\n"), func(k, v string, items []string) {
		pairs = append(pairs, yamlPair{k, v, items})
	})
	if len(pairs) != 1 {
		t.Fatalf("expected 1 pair, got %d", len(pairs))
	}
	if len(pairs[0].Items) != 2 {
		t.Fatalf("expected 2 items, got %v", pairs[0].Items)
	}
	if pairs[0].Items[0] != "Alias One" || pairs[0].Items[1] != "Alias Two" {
		t.Errorf("items = %v", pairs[0].Items)
	}
}

func TestScanYAML_EmptyInput(t *testing.T) {
	called := false
	scanYAML([]byte(""), func(k, v string, items []string) {
		called = true
	})
	if called {
		t.Error("should not call fn for empty input")
	}
}

func TestScanYAML_Comments(t *testing.T) {
	var pairs []yamlPair
	data := "# This is a comment\ntitle: test # inline comment\ntags:\n  - foo # item comment\n"
	scanYAML([]byte(data), func(k, v string, items []string) {
		pairs = append(pairs, yamlPair{k, v, items})
	})
	if len(pairs) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(pairs))
	}
	if pairs[0].Value != "test" {
		t.Errorf("inline comment not stripped: %q", pairs[0].Value)
	}
	if pairs[1].Items[0] != "foo" {
		t.Errorf("item comment not stripped: %v", pairs[1].Items)
	}
}

func TestScanYAML_BlankLines(t *testing.T) {
	var pairs []yamlPair
	data := "\n\ntitle: test\n\n\ntags:\n  - a\n\n  - b\n"
	scanYAML([]byte(data), func(k, v string, items []string) {
		pairs = append(pairs, yamlPair{k, v, items})
	})
	if len(pairs) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(pairs))
	}
	if len(pairs[1].Items) != 2 {
		t.Errorf("items = %v", pairs[1].Items)
	}
}

func TestScanYAML_CRLF(t *testing.T) {
	var pairs []yamlPair
	scanYAML([]byte("title: hi\r\ntags:\r\n  - one\r\n"), func(k, v string, items []string) {
		pairs = append(pairs, yamlPair{k, v, items})
	})
	if len(pairs) != 2 || pairs[0].Key != "title" || pairs[0].Value != "hi" {
		t.Errorf("CRLF: %+v", pairs)
	}
	if len(pairs[1].Items) != 1 || pairs[1].Items[0] != "one" {
		t.Errorf("CRLF items: %v", pairs[1].Items)
	}
}

func TestScanYAML_QuotedValueWithColon(t *testing.T) {
	var pairs []yamlPair
	scanYAML([]byte("title: 'Section: Intro'\n"), func(k, v string, items []string) {
		pairs = append(pairs, yamlPair{k, v, items})
	})
	if len(pairs) != 1 {
		t.Fatalf("expected 1 pair, got %d", len(pairs))
	}
	if pairs[0].Value != "Section: Intro" {
		t.Errorf("value with colon: %q", pairs[0].Value)
	}
}

func TestFindKeyColon_QuotedColon(t *testing.T) {
	idx := findKeyColon("title: \"hello: world\"")
	if idx != 5 {
		t.Errorf("findKeyColon should find first colon outside quotes: got %d", idx)
	}
}

func TestFindKeyColon_NoSpaceAfterColon(t *testing.T) {
	idx := findKeyColon("key:value")
	if idx != -1 {
		t.Errorf("colon without space should not be a key separator: got %d", idx)
	}
}

func TestStripQuotes_Double(t *testing.T) {
	if v := stripQuotes("\"hello\""); v != "hello" {
		t.Errorf("got %q", v)
	}
}

func TestStripQuotes_Single(t *testing.T) {
	if v := stripQuotes("'hello'"); v != "hello" {
		t.Errorf("got %q", v)
	}
}

func TestStripQuotes_NoQuotes(t *testing.T) {
	if v := stripQuotes("hello"); v != "hello" {
		t.Errorf("got %q", v)
	}
}

func TestStripInlineComment_NestedBrackets(t *testing.T) {
	if v := stripInlineComment("[a, b] # comment"); v != "[a, b]" {
		t.Errorf("comment in brackets: %q", v)
	}
}

func TestParseInlineArray_SingleItem(t *testing.T) {
	items := parseInlineArray("[single]")
	if len(items) != 1 || items[0] != "single" {
		t.Errorf("got %v", items)
	}
}

func TestParseInlineArray_Empty(t *testing.T) {
	items := parseInlineArray("[]")
	if len(items) != 0 {
		t.Errorf("empty array: %v", items)
	}
}

func TestSplitArrayItems_QuotedComma(t *testing.T) {
	parts := splitArrayItems(`"hello, world", "foo"`)
	if len(parts) != 2 {
		t.Fatalf("expected 2 parts, got %v", parts)
	}
}

func TestScanYAML_NestedArray(t *testing.T) {
	// Items inside a nested YAML block with indentation
	var pairs []yamlPair
	data := "tags:\n  - foo\n  - bar\n  - baz\ntitle: My Note\n"
	scanYAML([]byte(data), func(k, v string, items []string) {
		pairs = append(pairs, yamlPair{k, v, items})
	})
	if len(pairs) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(pairs))
	}
	sort.Strings(pairs[0].Items)
	if len(pairs[0].Items) != 3 {
		t.Errorf("array items count: %d", len(pairs[0].Items))
	}
	if pairs[1].Value != "My Note" {
		t.Errorf("title value: %q", pairs[1].Value)
	}
}
