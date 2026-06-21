package main

func (m *Model) switchToProfile(profileName string) {
	profile, ok := m.config.Profiles[profileName]
	if !ok {
		m.addToast("Profile not found: "+profileName, ToastError)
		return
	}

	// Apply profile settings
	if profile.Path != "" {
		m.config.VaultPath = profile.Path
	}
	if profile.Theme != "" {
		m.config.Theme = profile.Theme
		m.setTheme(profile.Theme)
	}
	if len(profile.SkipDirs) > 0 {
		m.config.SkipDirs = profile.SkipDirs
	}

	// Rescan vault with new settings
	m.rescanVault()
	m.mode = ModeBrowse
	m.addToast("Switched to profile: "+profileName, ToastInfo)
}

func (m *Model) setTheme(themeName string) {
	palette, err := lookupPalette(themeName)
	if err != nil {
		return
	}
	m.palette = palette
	m.viewer.renderStyle = markdownStyleFrom(palette, m.config.LineSpacing)
	m.searchStyle = searchStyleFrom(palette)
	m.fileTree.SetPalette(palette)
}
