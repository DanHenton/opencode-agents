package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/frontmatter"
)

type AgentConfig struct {
	Name        string
	Description string
	Prompt      string
	Metadata    map[string]interface{}
}

type CommandManager struct {
	SourceDir string
}

func NewCommandManager(dir string) *CommandManager {
	return &CommandManager{
		SourceDir: dir,
	}
}

func (m *CommandManager) LoadAgents() ([]AgentConfig, error) {
	var agents []AgentConfig

	entries, err := os.ReadDir(m.SourceDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("source directory does not exist: %s. Please create it or use --dir flag", m.SourceDir)
		}
		return nil, fmt.Errorf("failed to read source directory %s: %w", m.SourceDir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		filePath := filepath.Join(m.SourceDir, entry.Name())
		f, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
		}

		var metadata map[string]interface{}
		body, err := frontmatter.Parse(f, &metadata)
		f.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to parse frontmatter in %s: %w", filePath, err)
		}

		// Defaults from filename
		name := strings.TrimSuffix(entry.Name(), ".md")
		description := ""

		// Overrides from frontmatter
		if metadata != nil {
			if n, ok := metadata["name"].(string); ok && n != "" {
				name = n
			}
			if d, ok := metadata["description"].(string); ok && d != "" {
				description = d
			}
		} else {
			metadata = make(map[string]interface{})
		}

		agents = append(agents, AgentConfig{
			Name:        name,
			Description: description,
			Prompt:      strings.TrimSpace(string(body)),
			Metadata:    metadata,
		})
	}

	return agents, nil
}
