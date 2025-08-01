package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_DefaultConfig(t *testing.T) {
	// Create a temporary directory for the test to ensure a clean environment
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	// When no config file exists, Load should return the default config
	cfg, err := Load()
	require.NoError(t, err)

	// Assert that the default values are correct
	assert.Equal(t, "default", cfg.CurrentProfile)
	require.Contains(t, cfg.Profiles, "default")
	assert.Equal(t, "ollama", cfg.Profiles["default"].Provider)
	assert.Equal(t, "llama3", cfg.Profiles["default"].Model)
}

func TestSaveAndLoad_Cycle(t *testing.T) {
	// Setup a temporary home directory for the test
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	// 1. Create a custom config to save
	originalCfg := &Config{
		CurrentProfile: "test_profile",
		Profiles: map[string]Profile{
			"default": {
				Provider: "ollama",
				Model:    "llama3",
			},
			"test_profile": {
				Provider: "openai",
				Endpoint: "http://localhost:1234/v1",
				Model:    "test-model",
				APIKey:   "test-key",
			},
		},
	}

	// 2. Save the config
	err := originalCfg.Save()
	require.NoError(t, err)

	// Verify the file was actually created
	configPath, err := GetConfigPath()
	require.NoError(t, err)
	assert.FileExists(t, configPath)

	// 3. Load the config back from the file
	loadedCfg, err := Load()
	require.NoError(t, err)

	// 4. Assert that the loaded config is identical to the original
	assert.Equal(t, originalCfg, loadedCfg)
}

func TestGetConfigPath(t *testing.T) {
	// Setup a temporary home directory
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	// Call the function to test
	path, err := GetConfigPath()
	require.NoError(t, err)

	// Construct the expected path
	expectedPath := filepath.Join(tempDir, ".config", "llm-cli", "config.json")

	// Assert that the path is what we expect
	assert.Equal(t, expectedPath, path)
}
