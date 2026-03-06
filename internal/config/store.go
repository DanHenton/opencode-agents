package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type ConfigStore struct {
	FilePath string
	Data     map[string]interface{}
}

func NewConfigStore(path string) *ConfigStore {
	return &ConfigStore{
		FilePath: path,
		Data:     make(map[string]interface{}),
	}
}

func (s *ConfigStore) Load() error {
	content, err := os.ReadFile(s.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			s.Data = make(map[string]interface{})
			return nil
		}
		return fmt.Errorf("failed to read %s: %w", s.FilePath, err)
	}

	if len(content) == 0 {
		s.Data = make(map[string]interface{})
		return nil
	}

	if err := json.Unmarshal(content, &s.Data); err != nil {
		return fmt.Errorf("failed to parse JSON in %s: %w", s.FilePath, err)
	}

	return nil
}

func (s *ConfigStore) Save() error {
	content, err := json.MarshalIndent(s.Data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return os.WriteFile(s.FilePath, content, 0644)
}

func (s *ConfigStore) GetEnabledAgents() map[string]bool {
	enabled := make(map[string]bool)
	agentMap, ok := s.Data["agent"].(map[string]interface{})
	if !ok {
		return enabled
	}

	for name, v := range agentMap {
		if agentData, ok := v.(map[string]interface{}); ok {
			// In opencode.json, missing "disable" generally implies enabled
			disableVal, exists := agentData["disable"]
			if !exists {
				enabled[name] = true
			} else if disabled, ok := disableVal.(bool); ok && !disabled {
				enabled[name] = true
			}
		}
	}
	return enabled
}

func (s *ConfigStore) UpdateAgent(name, prompt string, disabled bool, metadata map[string]interface{}) {
	if s.Data["agent"] == nil {
		s.Data["agent"] = make(map[string]interface{})
	}

	agentMap, ok := s.Data["agent"].(map[string]interface{})
	if !ok {
		// If "agent" is not a map, override it safely
		agentMap = make(map[string]interface{})
		s.Data["agent"] = agentMap
	}

	var agentData map[string]interface{}
	if existingData, exists := agentMap[name]; exists {
		if dataMap, ok := existingData.(map[string]interface{}); ok {
			agentData = dataMap
		} else {
			agentData = make(map[string]interface{})
		}
	} else {
		agentData = make(map[string]interface{})
	}

	agentData["prompt"] = prompt
	agentData["disable"] = disabled

	// Inject metadata safely
	for k, v := range metadata {
		// Prevent overwriting internal keys
		if k == "prompt" || k == "disable" || k == "name" {
			continue
		}
		agentData[k] = v
	}

	agentMap[name] = agentData
}
