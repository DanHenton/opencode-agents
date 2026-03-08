package cli

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
)

//go:embed templates/AGENTS.md
var agentsTemplate embed.FS

// EnsureConfig checks if the agents directory and base configuration exist.
// If not, it prompts the user to initialize them.
func (c *CLI) EnsureConfig() error {
	// c.SourceDir is the agents directory (e.g. .../opencode-agents/agents)
	// The base config file AGENTS.md should be in the parent directory (.../opencode-agents/AGENTS.md)
	agentsDir := c.SourceDir
	configDir := filepath.Dir(agentsDir)
	agentsFile := filepath.Join(configDir, "AGENTS.md")

	// Check if agents directory exists
	dirExists := false
	if info, err := os.Stat(agentsDir); err == nil && info.IsDir() {
		dirExists = true
	}

	// Check if AGENTS.md exists
	fileExists := false
	if _, err := os.Stat(agentsFile); err == nil {
		fileExists = true
	}

	if dirExists && fileExists {
		return nil
	}

	// Prompt user
	var initConfig bool
	err := huh.NewConfirm().
		Title("Configuration missing").
		Description(fmt.Sprintf("Initialize config at %s?", configDir)).
		Affirmative("Yes").
		Negative("No").
		Value(&initConfig).
		Run()

	if err != nil {
		return fmt.Errorf("prompt cancelled: %w", err)
	}

	if !initConfig {
		return fmt.Errorf("configuration required to proceed")
	}

	// Create directories (including parents)
	if err := os.MkdirAll(agentsDir, 0755); err != nil {
		return fmt.Errorf("failed to create agents directory: %w", err)
	}
	if !dirExists {
		fmt.Printf("Created directory: %s\n", agentsDir)
	}

	// Write AGENTS.md if missing
	if !fileExists {
		content, err := agentsTemplate.ReadFile("templates/AGENTS.md")
		if err != nil {
			return fmt.Errorf("failed to read embedded template: %w", err)
		}
		if err := os.WriteFile(agentsFile, content, 0644); err != nil {
			return fmt.Errorf("failed to write AGENTS.md: %w", err)
		}
		fmt.Printf("Created file: %s\n", agentsFile)
	}

	return nil
}
