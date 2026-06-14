package main

import (
	"os"
	"path/filepath"
	"strings"
)

// ResolveWikiLink resolves a [[wiki-link]] target to a vault file path.
func ResolveWikiLink(target string, vault *VaultEntry, vaultRoot string) string {
	if target == "" {
		return ""
	}

	target = strings.SplitN(target, "#", 2)[0]

	if !strings.HasSuffix(target, ".md") {
		exact := findExactPath(vault, "", target+".md")
		if exact != "" {
			return exact
		}
		exact = findExactPath(vault, "", target+".markdown")
		if exact != "" {
			return exact
		}
	} else {
		exact := findExactPath(vault, "", target)
		if exact != "" {
			return exact
		}
	}

	basename := strings.ToLower(target)
	if !strings.HasSuffix(basename, ".md") {
		basename += ".md"
	}
	result := findBasename(vault, "", basename)
	if result != "" {
		return result
	}

	result = findAlias(vault, target, vaultRoot)
	if result != "" {
		return result
	}

	return ""
}

func findAlias(vault *VaultEntry, alias string, vaultRoot string) string {
	aliasLower := strings.ToLower(alias)
	for _, child := range vault.Children {
		if child.IsDir {
			found := findAlias(child, alias, vaultRoot)
			if found != "" {
				return found
			}
			continue
		}
		aliasEntries, err := extractAliasesFromFile(vaultRoot, child.Path)
		if err != nil {
			continue
		}
		for _, a := range aliasEntries {
			if strings.ToLower(a) == aliasLower {
				return child.Path
			}
		}
	}
	return ""
}

func extractAliasesFromFile(vaultRoot, relativePath string) ([]string, error) {
	fullPath := relativePath
	if vaultRoot != "" {
		fullPath = filepath.Join(vaultRoot, relativePath)
	}
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}

	content := string(data)
	fm, body := parseFrontmatter(content)
	if fm.Aliases == nil && body == content {
		return nil, nil
	}
	return fm.Aliases, nil
}

func findExactPath(vault *VaultEntry, prefix, target string) string {
	for _, child := range vault.Children {
		childPath := child.Path
		if prefix != "" {
			childPath = prefix + "/" + child.Name
		}
		if childPath == target && !child.IsDir {
			return childPath
		}
		if child.IsDir {
			found := findExactPath(child, childPath, target)
			if found != "" {
				return found
			}
		}
	}
	return ""
}

func findBasename(vault *VaultEntry, prefix, target string) string {
	targetLower := strings.ToLower(target)
	for _, child := range vault.Children {
		childPath := child.Path
		if prefix != "" {
			childPath = prefix + "/" + child.Name
		}
		nameLower := strings.ToLower(child.Name)
		if nameLower == targetLower && !child.IsDir {
			return childPath
		}
		if child.IsDir {
			found := findBasename(child, childPath, target)
			if found != "" {
				return found
			}
		}
	}
	return ""
}
