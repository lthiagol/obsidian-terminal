package main

import "strings"

type yamlPair struct {
	Key   string
	Value string
	Items []string
}

func scanYAML(data []byte, fn func(key, value string, items []string)) {
	text := string(data)
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	lines := strings.Split(text, "\n")
	i := 0
	for i < len(lines) {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			i++
			continue
		}

		colonIdx := findKeyColon(trimmed)
		if colonIdx < 0 {
			i++
			continue
		}

		key := strings.TrimSpace(trimmed[:colonIdx])
		rest := strings.TrimSpace(trimmed[colonIdx+1:])

		rest = stripInlineComment(rest)

		if rest == "" {
			var items []string
			i++
			for i < len(lines) {
				itemLine := lines[i]
				itemTrimmed := strings.TrimSpace(itemLine)
				if itemTrimmed == "" || strings.HasPrefix(itemTrimmed, "#") {
					i++
					continue
				}
				if !strings.HasPrefix(itemTrimmed, "- ") && itemTrimmed != "-" {
					break
				}
				if len(itemLine) > 0 && (itemLine[0] != ' ' && itemLine[0] != '\t') {
					break
				}
				val := strings.TrimSpace(itemTrimmed[1:])
				val = stripInlineComment(val)
				val = stripQuotes(val)
				items = append(items, val)
				i++
			}
			fn(key, "", items)
			continue
		}

		if strings.HasPrefix(rest, "[") {
			if !strings.HasSuffix(rest, "]") {
				i++
				continue
			}
			items := parseInlineArray(rest)
			fn(key, "", items)
			i++
			continue
		}

		value := stripQuotes(rest)
		fn(key, value, nil)
		i++
	}
}

func findKeyColon(line string) int {
	inSingle := false
	inDouble := false
	for i, c := range line {
		switch c {
		case '\'':
			if !inDouble {
				inSingle = !inSingle
			}
		case '"':
			if !inSingle {
				inDouble = !inDouble
			}
		case ':':
			if !inSingle && !inDouble {
				if i+1 >= len(line) || line[i+1] == ' ' || line[i+1] == '\t' {
					return i
				}
			}
		}
	}
	return -1
}

func stripInlineComment(s string) string {
	inSingle := false
	inDouble := false
	depth := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '\'':
			if !inDouble {
				inSingle = !inSingle
			}
		case '"':
			if !inSingle {
				inDouble = !inDouble
			}
		case '[':
			if !inSingle && !inDouble {
				depth++
			}
		case ']':
			if !inSingle && !inDouble {
				depth--
			}
		case '#':
			if !inSingle && !inDouble && depth == 0 && i > 0 && s[i-1] == ' ' {
				return strings.TrimRight(s[:i-1], " ")
			}
		}
	}
	return s
}

func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func parseInlineArray(s string) []string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "[") {
		s = s[1:]
	}
	if strings.HasSuffix(s, "]") {
		s = s[:len(s)-1]
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}

	parts := splitArrayItems(s)
	var items []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		p = stripQuotes(p)
		if p != "" {
			items = append(items, p)
		}
	}
	return items
}

func splitArrayItems(s string) []string {
	var parts []string
	var current strings.Builder
	inSingle := false
	inDouble := false
	depth := 0

	for i := 0; i < len(s); i++ {
		c := s[i]
		switch c {
		case '\'':
			if !inDouble {
				inSingle = !inSingle
				current.WriteByte(c)
			}
		case '"':
			if !inSingle {
				inDouble = !inDouble
				current.WriteByte(c)
			}
		case '[':
			if !inSingle && !inDouble {
				depth++
				current.WriteByte(c)
			}
		case ']':
			if !inSingle && !inDouble {
				depth--
				current.WriteByte(c)
			}
		case ',':
			if !inSingle && !inDouble && depth == 0 {
				parts = append(parts, current.String())
				current.Reset()
				continue
			}
			current.WriteByte(c)
		default:
			current.WriteByte(c)
		}
	}
	if current.Len() > 0 {
		parts = append(parts, current.String())
	}
	return parts
}

// parseNestedMap parses a nested YAML map structure like:
// profiles:
//   work:
//     path: /path/to/work
//     theme: dark
// Returns a map of profile name -> map of key -> value
func parseNestedMap(data []byte, rootKey string) map[string]map[string]string {
	text := string(data)
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	lines := strings.Split(text, "\n")
	result := make(map[string]map[string]string)

	// Find the root key
	rootIdx := -1
	rootIndent := 0
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		indent := len(line) - len(strings.TrimLeft(line, " \t"))
		colonIdx := findKeyColon(trimmed)
		if colonIdx >= 0 {
			key := strings.TrimSpace(trimmed[:colonIdx])
			if key == rootKey {
				rootIdx = i
				rootIndent = indent
				break
			}
		}
	}

	if rootIdx < 0 {
		return result
	}

	// Parse nested structure
	currentProfile := ""
	for i := rootIdx + 1; i < len(lines); i++ {
		line := lines[i]
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		indent := len(line) - len(strings.TrimLeft(line, " \t"))
		if indent <= rootIndent {
			break
		}

		trimmed := strings.TrimSpace(line)
		colonIdx := findKeyColon(trimmed)
		if colonIdx < 0 {
			continue
		}

		key := strings.TrimSpace(trimmed[:colonIdx])
		value := strings.TrimSpace(trimmed[colonIdx+1:])
		value = stripInlineComment(value)
		value = stripQuotes(value)

		// Check if this is a profile name (second level) or a property (third level)
		profileIndent := rootIndent + 2
		if indent == profileIndent && value == "" {
			// This is a profile name
			currentProfile = key
			if _, exists := result[currentProfile]; !exists {
				result[currentProfile] = make(map[string]string)
			}
		} else if indent > profileIndent && currentProfile != "" && value != "" {
			// This is a property of the current profile
			result[currentProfile][key] = value
		}
	}

	return result
}
