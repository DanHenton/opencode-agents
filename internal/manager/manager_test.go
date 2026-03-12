package manager

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCommandManager_LoadAgents(t *testing.T) {
	tmpDir := t.TempDir()

	file1Content := `---
name: overridden-name
description: A basic plan agent
model: claude-3-opus
temperature: 0.7
---
I am a plan agent.`

	file2Content := `I am a raw file without frontmatter.`

	// Valid markdown with frontmatter
	err := os.WriteFile(filepath.Join(tmpDir, "plan.md"), []byte(file1Content), 0644)
	if err != nil {
		t.Fatalf("Failed to write plan.md: %v", err)
	}

	// Valid markdown without frontmatter
	err = os.WriteFile(filepath.Join(tmpDir, "raw.md"), []byte(file2Content), 0644)
	if err != nil {
		t.Fatalf("Failed to write raw.md: %v", err)
	}

	// Invalid extension
	err = os.WriteFile(filepath.Join(tmpDir, "ignore.txt"), []byte("Ignore me"), 0644)
	if err != nil {
		t.Fatalf("Failed to write ignore.txt: %v", err)
	}

	manager := NewCommandManager(tmpDir)
	agents, err := manager.LoadAgents()
	if err != nil {
		t.Fatalf("Expected LoadAgents to succeed, got %v", err)
	}

	if len(agents) != 2 {
		t.Fatalf("Expected 2 agents to be loaded, got %d", len(agents))
	}

	// Validate agent 1
	var planAgent AgentConfig
	for _, a := range agents {
		if a.Name == "overridden-name" {
			planAgent = a
		}
	}

	if planAgent.Name != "overridden-name" {
		t.Errorf("Expected name 'overridden-name', got %s", planAgent.Name)
	}
	if planAgent.Description != "A basic plan agent" {
		t.Errorf("Expected description 'A basic plan agent', got %s", planAgent.Description)
	}
	if planAgent.Prompt != "I am a plan agent." {
		t.Errorf("Expected prompt 'I am a plan agent.', got %s", planAgent.Prompt)
	}
	if planAgent.Metadata["model"] != "claude-3-opus" {
		t.Errorf("Expected model 'claude-3-opus', got %v", planAgent.Metadata["model"])
	}

	// Validate agent 2
	var rawAgent AgentConfig
	for _, a := range agents {
		if a.Name == "raw" {
			rawAgent = a
		}
	}

	if rawAgent.Name != "raw" {
		t.Errorf("Expected name 'raw', got %s", rawAgent.Name)
	}
	if rawAgent.Prompt != "I am a raw file without frontmatter." {
		t.Errorf("Expected prompt 'I am a raw file without frontmatter.', got %s", rawAgent.Prompt)
	}
}

func TestCommandManager_LoadAgentsWithPermissions(t *testing.T) {
	tmpDir := t.TempDir()

	fileContent := `---
name: secure-agent
permission:
  edit: deny
  bash: ask
  read:
    internal/secrets: deny
---
Prompt body`

	err := os.WriteFile(filepath.Join(tmpDir, "secure.md"), []byte(fileContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write secure.md: %v", err)
	}

	manager := NewCommandManager(tmpDir)
	agents, err := manager.LoadAgents()
	if err != nil {
		t.Fatalf("Expected LoadAgents to succeed, got %v", err)
	}

	if len(agents) != 1 {
		t.Fatalf("Expected 1 agent, got %d", len(agents))
	}

	agent := agents[0]
	permissions, ok := agent.Metadata["permission"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected permission to be a map[string]interface{}, got %T", agent.Metadata["permission"])
	}

	if permissions["edit"] != "deny" {
		t.Errorf("Expected edit: deny, got %v", permissions["edit"])
	}
	if permissions["bash"] != "ask" {
		t.Errorf("Expected bash: ask, got %v", permissions["bash"])
	}

	readPerms, ok := permissions["read"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected permission.read to be a map[string]interface{}, got %T", permissions["read"])
	}

	if readPerms["internal/secrets"] != "deny" {
		t.Errorf("Expected read.internal/secrets: deny, got %v", readPerms["internal/secrets"])
	}
}

func TestCommandManager_DirNotExist(t *testing.T) {
	manager := NewCommandManager("/path/to/nowhere/that/does/not/exist/ever")
	agents, err := manager.LoadAgents()
	if err == nil {
		t.Fatalf("Expected error when directory doesn't exist, got nil")
	}
	if len(agents) != 0 {
		t.Fatalf("Expected 0 agents on error, got %d", len(agents))
	}
}
