/*
Copyright © 2025 magifd2

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/magifd2/llm-cli/internal/config"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestEnvironment creates a temporary directory and a dummy config file.
func setupTestEnvironment(t *testing.T) string {
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
	err := cfg.Save(cfgFile)
	require.NoError(t, err)
	return tempDir
}

// executeCommand is a helper function to execute cobra commands and capture their output.
func executeCommand(root *cobra.Command, args ...string) (string, string, error) {
	var outBuf, errBuf bytes.Buffer
	root.SetOut(&outBuf)
	root.SetErr(&errBuf)
	root.SetArgs(args)

	err := root.Execute()

	return outBuf.String(), errBuf.String(), err
}

func TestAddCommand(t *testing.T) {
	_ = setupTestEnvironment(t)

	// Test adding a new profile
	_, _, err := executeCommand(rootCmd, "profile", "add", "new_profile")
	assert.NoError(t, err)

	// Verify the profile was added
	cfg, err := config.Load(cfgFile)
	require.NoError(t, err)
	assert.Contains(t, cfg.Profiles, "new_profile")
	assert.Equal(t, "ollama", cfg.Profiles["new_profile"].Provider) // Should copy from default

	// Test adding a profile that already exists (should fail)
	_, _, err = executeCommand(rootCmd, "profile", "add", "default")
	assert.Error(t, err)
}

func TestUseCommand(t *testing.T) {
	_ = setupTestEnvironment(t)

	// Test switching to an existing profile
	_, _, err := executeCommand(rootCmd, "profile", "use", "existing_profile")
	assert.NoError(t, err)

	// Verify the current profile was changed
	cfg, err := config.Load(cfgFile)
	require.NoError(t, err)
	assert.Equal(t, "existing_profile", cfg.CurrentProfile)

	// Test switching to a non-existent profile (should fail)
	_, _, err = executeCommand(rootCmd, "profile", "use", "non_existent_profile")
	assert.Error(t, err)
}

func TestSetCommand(t *testing.T) {
	_ = setupTestEnvironment(t)

	// Test setting a new model for the default profile
	_, _, err := executeCommand(rootCmd, "profile", "set", "model", "test-model-123")
	assert.NoError(t, err)

	// Verify the model was set
	cfg, err := config.Load(cfgFile)
	require.NoError(t, err)
	assert.Equal(t, "test-model-123", cfg.Profiles["default"].Model)

	// Test setting an invalid key (should fail)
	_, _, err = executeCommand(rootCmd, "profile", "set", "invalid_key", "test_value")
	assert.Error(t, err)
}

func TestRemoveCommand(t *testing.T) {
	_ = setupTestEnvironment(t)

	// Test removing an existing, non-active profile
	_, _, err := executeCommand(rootCmd, "profile", "remove", "existing_profile")
	assert.NoError(t, err)

	// Verify the profile was removed
	cfg, err := config.Load(cfgFile)
	require.NoError(t, err)
	assert.NotContains(t, cfg.Profiles, "existing_profile")

	// Test removing the default profile (should fail)
	_, _, err = executeCommand(rootCmd, "profile", "remove", "default")
	assert.Error(t, err)

	// Test removing the active profile (should fail)
	// (The active profile is 'default' in this setup)
	// To test this, we need to ensure 'default' is the active profile
	// and then try to remove it.
	cfg, err = config.Load(cfgFile)
	require.NoError(t, err)
	cfg.CurrentProfile = "default"
	require.NoError(t, cfg.Save(cfgFile))

	_, _, err = executeCommand(rootCmd, "profile", "remove", "default")
	assert.Error(t, err)
}

// TestMain is required to reset the command state between tests.
func TestMain(m *testing.M) {
	// This is a bit of a hack to allow cobra's state to be reset between tests.
	// It's not ideal, but it's a common pattern for testing cobra apps.
	var (
		originalOut = os.Stdout
		originalErr = os.Stderr
	)

	// Run the tests
	code := m.Run()

	// Restore original stdout and stderr
	profileCmd.SetOut(originalOut)
	profileCmd.SetErr(originalErr)

	os.Exit(code)
}
