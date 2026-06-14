package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	vaultFlag := flag.String("vault", "", "path to Obsidian vault")
	configPathFlag := flag.String("config", "", "path to config file")
	profileFlag := flag.String("profile", "", "vault profile to use")
	flag.Parse()

	cfgPath := configPathOrDefault(*configPathFlag)

	cfg, err := LoadConfig(cfgPath)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}
		cfg = DefaultConfig()
	}

	// Apply --vault flag (takes precedence)
	if *vaultFlag != "" {
		cfg.VaultPath = *vaultFlag
	}

	// Apply --profile flag
	if *profileFlag != "" {
		if cfg.Profiles == nil {
			fmt.Fprintf(os.Stderr, "Error: no profiles defined in config\n")
			os.Exit(1)
		}
		profile, ok := cfg.Profiles[*profileFlag]
		if !ok {
			fmt.Fprintf(os.Stderr, "Error: profile %q not found\n", *profileFlag)
			os.Exit(1)
		}
		// Apply profile settings (vault flag takes precedence)
		if *vaultFlag == "" && profile.Path != "" {
			cfg.VaultPath = profile.Path
		}
		if profile.Theme != "" {
			cfg.Theme = profile.Theme
		}
		if len(profile.SkipDirs) > 0 {
			cfg.SkipDirs = profile.SkipDirs
		}
	}

	// If no vault path but profiles exist, enter picker mode
	if cfg.VaultPath == "" && len(cfg.Profiles) > 0 {
		m := NewModel(cfg)
		if m.err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", m.err)
			os.Exit(1)
		}
		p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if cfg.VaultPath == "" {
		fmt.Fprintln(os.Stderr, "Error: vault path required. Use --vault flag or set vault_path in config.")
		fmt.Fprintf(os.Stderr, "Config file location: %s\n", cfgPath)
		os.Exit(1)
	}

	m := NewModel(cfg)
	if m.err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", m.err)
		os.Exit(1)
	}

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
