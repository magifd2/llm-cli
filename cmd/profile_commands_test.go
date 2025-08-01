package cmd

import (
	"testing"

	"github.com/magifd2/llm-cli/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestEnvironment creates a temporary directory and a dummy config file.
func setupTestEnvironment(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	// Create an initial config to work with
	cfg := &config.Config{
		CurrentProfile: "default",
		Profiles: map[string]config.Profile{
			"default": {
				Provider: "ollama",
				Model:    "llama3",
			},
			"existing_profile": {
				Provider: "openai",
				Model:    "gpt-4",
			},
		},
	}
	err := cfg.Save()
	require.NoError(t, err)
}

func TestAddCommand(t *testing.T) {
	setupTestEnvironment(t)

	// Test adding a new profile
	err := addProfile("new_profile")
	assert.NoError(t, err)

	// Verify the profile was added
	cfg, err := config.Load()
	require.NoError(t, err)
	assert.Contains(t, cfg.Profiles, "new_profile")
	assert.Equal(t, "ollama", cfg.Profiles["new_profile"].Provider) // Should copy from default

	// Test adding a profile that already exists (should fail)
	err = addProfile("default")
	assert.Error(t, err)
}

func TestUseCommand(t *testing.T) {
	setupTestEnvironment(t)

	// Test switching to an existing profile
	err := useProfile("existing_profile")
	assert.NoError(t, err)

	// Verify the current profile was changed
	cfg, err := config.Load()
	require.NoError(t, err)
	assert.Equal(t, "existing_profile", cfg.CurrentProfile)

	// Test switching to a non-existent profile (should fail)
	err = useProfile("non_existent_profile")
	assert.Error(t, err)
}

func TestSetCommand(t *testing.T) {
	setupTestEnvironment(t)

	// Test setting a new model for the default profile
	err := setProfileValue("model", "test-model-123")
	assert.NoError(t, err)

	// Verify the model was set
	cfg, err := config.Load()
	require.NoError(t, err)
	assert.Equal(t, "test-model-123", cfg.Profiles["default"].Model)

	// Test setting an invalid key (should fail)
	err = setProfileValue("invalid_key", "test_value")
	assert.Error(t, err)
}

func TestRemoveCommand(t *testing.T) {
	setupTestEnvironment(t)

	// Test removing an existing, non-active profile
	err := removeProfile("existing_profile")
	assert.NoError(t, err)

	// Verify the profile was removed
	cfg, err := config.Load()
	require.NoError(t, err)
	assert.NotContains(t, cfg.Profiles, "existing_profile")

	// Test removing the default profile (should fail)
	err = removeProfile("default")
	assert.Error(t, err)

	// Test removing the active profile (should fail)
	// (The active profile is 'default' in this setup)
	// To test this, we need to ensure 'default' is the active profile
	// and then try to remove it.
	cfg, err = config.Load()
	require.NoError(t, err)
	cfg.CurrentProfile = "default"
	require.NoError(t, cfg.Save())

	err = removeProfile("default")
	assert.Error(t, err)
}
