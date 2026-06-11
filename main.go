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

	if *vaultFlag != "" {
		cfg.VaultPath = *vaultFlag
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
