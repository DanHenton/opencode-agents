package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestConfigStore_LoadAndSave(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "opencode.json")

	initialJSON := `{
		"server": {"port": 8080},
		"agent": {
			"plan": {
				"disable": false,
				"prompt": "I am a planner."
			}
		}
	}`
	err := os.WriteFile(filePath, []byte(initialJSON), 0644)
	if err != nil {
		t.Fatalf("Failed to write initial json: %v", err)
	}

	store := NewConfigStore(filePath)
	if err := store.Load(); err != nil {
		t.Fatalf("Failed to load store: %v", err)
	}

	// Verify unknown fields are preserved
	if store.Data["server"] == nil {
		t.Errorf("Expected 'server' key to be preserved")
	}

	// Verify agent loading
	agentMap, ok := store.Data["agent"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected 'agent' key to be a map")
	}
	if planAgent, exists := agentMap["plan"]; !exists {
		t.Errorf("Expected 'plan' agent to exist")
	} else {
		planMap := planAgent.(map[string]interface{})
		if planMap["prompt"] != "I am a planner." {
			t.Errorf("Expected prompt to be 'I am a planner.', got %v", planMap["prompt"])
		}
	}

	// Save and verify output
	if err := store.Save(); err != nil {
		t.Fatalf("Failed to save store: %v", err)
	}

	savedContent, _ := os.ReadFile(filePath)
	var savedData map[string]interface{}
	json.Unmarshal(savedContent, &savedData)

	if savedData["server"] == nil {
		t.Errorf("Expected 'server' key in saved JSON")
	}
}

func TestConfigStore_LoadMissingFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "does-not-exist.json")

	store := NewConfigStore(filePath)
	if err := store.Load(); err != nil {
		t.Fatalf("Expected no error when file is missing, got %v", err)
	}

	if store.Data == nil {
		t.Fatalf("Expected Data to be initialized, got nil")
	}
}

func TestConfigStore_GetEnabledAgents(t *testing.T) {
	store := NewConfigStore("dummy.json")
	store.Data = map[string]interface{}{
		"agent": map[string]interface{}{
			"agent-enabled": map[string]interface{}{
				"disable": false,
			},
			"agent-disabled": map[string]interface{}{
				"disable": true,
			},
			"agent-implicit": map[string]interface{}{
				"prompt": "Implicitly enabled because disable is missing",
			},
		},
	}

	enabled := store.GetEnabledAgents()

	if !enabled["agent-enabled"] {
		t.Errorf("Expected agent-enabled to be true")
	}
	if enabled["agent-disabled"] {
		t.Errorf("Expected agent-disabled to be false")
	}
	if !enabled["agent-implicit"] {
		t.Errorf("Expected agent-implicit to be true")
	}
}

func TestConfigStore_UpdateAgent(t *testing.T) {
	store := NewConfigStore("dummy.json")
	store.Data = map[string]interface{}{
		"server": "keep-me",
	}

	metadata := map[string]interface{}{
		"model":       "claude-3-opus",
		"temperature": 0.5,
		"name":        "should-be-ignored", // Explicitly ignored in UpdateAgent
	}

	store.UpdateAgent("test-agent", "You are a test agent.", false, metadata)

	agentMap, ok := store.Data["agent"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected 'agent' map to be created")
	}

	testAgent, exists := agentMap["test-agent"]
	if !exists {
		t.Fatalf("Expected 'test-agent' to be created")
	}

	taMap, ok := testAgent.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected test-agent to be a map")
	}

	if taMap["prompt"] != "You are a test agent." {
		t.Errorf("Expected prompt to match, got %v", taMap["prompt"])
	}
	if taMap["disable"] != false {
		t.Errorf("Expected disable to be false")
	}
	if taMap["model"] != "claude-3-opus" {
		t.Errorf("Expected metadata 'model' to be 'claude-3-opus', got %v", taMap["model"])
	}
	if taMap["name"] == "should-be-ignored" {
		t.Errorf("Expected 'name' in metadata to be ignored, got %v", taMap["name"])
	}

	// Ensure other keys survived
	if store.Data["server"] != "keep-me" {
		t.Errorf("Expected root level keys to survive mutation")
	}
}
