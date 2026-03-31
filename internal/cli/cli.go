package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/danhenton/opencode-agents/internal/config"
	"github.com/danhenton/opencode-agents/internal/manager"
)

type CLI struct {
	TargetGlobal bool
	TargetLocal  bool
	SyncAll      bool
	SourceDir    string

	Manager *manager.CommandManager
	Store   *config.ConfigStore
}

func NewCLI(global, local, syncAll bool, dir string) *CLI {
	if dir == "" {
		dir = os.Getenv("OPENCODE_AGENTS_DIR")
		if dir == "" {
			home, err := os.UserHomeDir()
			if err == nil {
				dir = filepath.Join(home, ".config", "opencode-agents", "agents")
			} else {
				dir = "agents" // fallback
			}
		}
	}

	return &CLI{
		TargetGlobal: global,
		TargetLocal:  local,
		SyncAll:      syncAll,
		SourceDir:    dir,
		Manager:      manager.NewCommandManager(dir),
	}
}

func (c *CLI) Run() error {
	// Ensure configuration exists before proceeding
	if err := c.EnsureConfig(); err != nil {
		return err
	}

	var targetPath string
	home, err := os.UserHomeDir()
	globalPath := ""
	if err == nil {
		globalPath = filepath.Join(home, ".config", "opencode", "opencode.json")
	}

	if c.TargetGlobal {
		targetPath = globalPath
	} else if c.TargetLocal {
		targetPath = "./opencode.json"
	} else {
		var targetChoice string
		err := huh.NewSelect[string]().
			Title("Which configuration do you want to update?").
			Options(
				huh.NewOption(fmt.Sprintf("Global (%s)", globalPath), "global"),
				huh.NewOption("Local (./opencode.json)", "local"),
			).
			Value(&targetChoice).
			Run()

		if err != nil {
			return fmt.Errorf("prompt cancelled: %w", err)
		}

		if targetChoice == "global" {
			targetPath = globalPath
		} else {
			targetPath = "./opencode.json"
		}
	}

	c.Store = config.NewConfigStore(targetPath)
	if err := c.Store.Load(); err != nil {
		return fmt.Errorf("failed to load opencode.json: %w", err)
	}

	agents, err := c.Manager.LoadAgents()
	if err != nil {
		return fmt.Errorf("failed to load agents: %w", err)
	}

	if len(agents) == 0 {
		fmt.Printf("No markdown agents found in %s.\n", c.SourceDir)
		return nil
	}

	enabledInStore := c.Store.GetEnabledAgents()
	var selectedAgents []string
	var options []huh.Option[string]

	for _, a := range agents {
		desc := a.Name
		if a.Description != "" {
			desc = fmt.Sprintf("%s - %s", a.Name, a.Description)
		}
		opt := huh.NewOption(desc, a.Name)

		// Pre-select if currently enabled in opencode.json
		if enabledInStore[a.Name] {
			selectedAgents = append(selectedAgents, a.Name)
			opt = opt.Selected(true)
		}

		options = append(options, opt)
	}

	if c.SyncAll {
		// Just sync all to selectedAgents
		selectedAgents = nil
		for _, a := range agents {
			selectedAgents = append(selectedAgents, a.Name)
		}
	} else {
		err := huh.NewMultiSelect[string]().
			Title(fmt.Sprintf("Select agents to sync into %s", targetPath)).
			Options(options...).
			Value(&selectedAgents).
			Run()

		if err != nil {
			return fmt.Errorf("prompt cancelled: %w", err)
		}
	}

	// Create a quick lookup for selected agents
	selectedMap := make(map[string]bool)
	for _, name := range selectedAgents {
		selectedMap[name] = true
	}

	// Apply mutations
	enabledCount := 0
	disabledCount := 0

	for _, a := range agents {
		disabled := !selectedMap[a.Name]
		c.Store.UpdateAgent(a.Name, a.Prompt, disabled, a.Metadata)

		if disabled {
			disabledCount++
		} else {
			enabledCount++
		}
	}

	// Save back to JSON
	if err := c.Store.Save(); err != nil {
		return fmt.Errorf("failed to save opencode.json: %w", err)
	}

	fmt.Printf("\nSuccess! Updated %s\n", targetPath)
	fmt.Printf("- %d agents enabled\n", enabledCount)
	fmt.Printf("- %d agents disabled\n", disabledCount)
	return nil
}
