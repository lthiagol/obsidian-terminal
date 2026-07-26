package main

import (
	"os"
	"time"
)

func vaultStateFrom(scanErrorCount int) VaultState {
	if scanErrorCount > 0 {
		return VaultStatePartial
	}
	return VaultStateOK
}

func (m *Model) checkVaultChanges() {
	if time.Since(m.lastRescan) < 2*time.Second {
		return
	}

	info, err := os.Stat(m.config.VaultPath)
	if err != nil {
		if m.vaultState != VaultStateBroken {
			m.vaultState = VaultStateBroken
			m.addToast("Vault inaccessible: "+err.Error(), ToastError)
		}
		return
	}

	if m.vaultState == VaultStateBroken {
		m.addToast("Vault is accessible again — rescanning", ToastInfo)
		m.rescanVault()
		return
	}

	if !info.ModTime().After(m.lastRootModTime) {
		return
	}

	m.rescanVault()
}

func (m *Model) rescanVault() {
	m.lastRescan = time.Now()

	info, err := os.Stat(m.config.VaultPath)
	if err != nil {
		m.vaultState = VaultStateBroken
		m.addToast("Cannot rescan vault: "+err.Error(), ToastError)
		return
	}
	m.lastRootModTime = info.ModTime()

	tree, indexes, scanErrors, err := ScanVault(m.config.VaultPath, m.config.SkipDirs)
	if err != nil {
		m.vaultState = VaultStateBroken
		m.addToast("Vault scan failed: "+err.Error(), ToastError)
		return
	}
	m.scanErrors = scanErrors
	m.vaultState = vaultStateFrom(len(scanErrors))

	oldActivePath := ""
	if m.activeNote != nil {
		oldActivePath = m.activeNote.Path
	}

	m.vault = tree
	m.searchIndex = indexes.Search
	m.backlinkIndex = indexes.Backlinks
	m.tagIndex = indexes.Tags
	m.allPaths = allPaths(tree)
	m.fileTree = NewFileTree(tree, m.palette)
	m.validatePins()

	if oldActivePath != "" {
		note, err := LoadNote(m.config.VaultPath, oldActivePath)
		if err != nil {
			m.addToast("Note was deleted: "+oldActivePath, ToastWarning)
			m.mode = ModeBrowse
			m.activeNote = nil
		} else {
			m.loadNote(note.Path, navReload)
		}
	}
}

func countFiles(entry *VaultEntry) int {
	if entry == nil {
		return 0
	}
	count := 0
	for _, child := range entry.Children {
		if child.IsDir {
			count += countFiles(child)
		} else {
			count++
		}
	}
	return count
}
