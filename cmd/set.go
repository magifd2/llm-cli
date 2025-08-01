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
	"fmt"
	"strconv"
	"strings"

	"github.com/magifd2/llm-cli/internal/config"
	"github.com/spf13/cobra"
)

// setCmd represents the 'profile set' command.
// This command allows users to set a specific configuration value for the currently active profile.
var setCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a value in the current profile",
	Long:  `Set a configuration value for the currently active profile.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("Error loading config: %w", err)
		}

		if err := setProfileValue(args[0], args[1]); err != nil {
			return fmt.Errorf("Error: %w", err)
		}
		// Success message moved here
		fmt.Printf("Set %s = %s in profile %s\n", args[0], args[1], cfg.CurrentProfile)
		return nil
	},
}

// setProfileValue updates a specific key-value pair in the currently active profile.
// It handles different configuration keys and returns an error for unknown keys.
func setProfileValue(key, value string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	profile, ok := cfg.Profiles[cfg.CurrentProfile]
	if !ok {
		return fmt.Errorf("active profile '%s' not found", cfg.CurrentProfile)
	}

	switch key {
	case "model":
		profile.Model = value
	case "provider":
		profile.Provider = value
	case "endpoint":
		profile.Endpoint = value
	case "api_key":
		profile.APIKey = value
	case "aws_region":
		profile.AWSRegion = value
	case "aws_access_key_id":
		profile.AWSAccessKeyID = value
	case "aws_secret_access_key":
		profile.AWSSecretAccessKey = value
	case "project_id":
		profile.ProjectID = value
	case "location":
		profile.Location = value
	case "credentials_file":
		profile.CredentialsFile = value
	case "limits.enabled":
		enabled, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid boolean value for limits.enabled: %s", value)
		}
		profile.Limits.Enabled = enabled
	case "limits.on_input_exceeded":
		if value != "stop" && value != "warn" {
			return fmt.Errorf("invalid value for limits.on_input_exceeded: must be 'stop' or 'warn'")
		}
		profile.Limits.OnInputExceeded = value
	case "limits.on_output_exceeded":
		if value != "stop" && value != "warn" {
			return fmt.Errorf("invalid value for limits.on_output_exceeded: must be 'stop' or 'warn'")
		}
		profile.Limits.OnOutputExceeded = value
	case "limits.max_prompt_size_bytes":
		size, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid integer value for limits.max_prompt_size_bytes: %s", value)
		}
		profile.Limits.MaxPromptSizeBytes = size
	case "limits.max_response_size_bytes":
		size, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid integer value for limits.max_response_size_bytes: %s", value)
		}
		profile.Limits.MaxResponseSizeBytes = size
	default:
		availableKeys := []string{
			"model", "provider", "endpoint", "api_key", "aws_region", "aws_access_key_id", "aws_secret_access_key", "project_id", "location", "credentials_file",
			"limits.enabled", "limits.on_input_exceeded", "limits.on_output_exceeded", "limits.max_prompt_size_bytes", "limits.max_response_size_bytes",
		}
		return fmt.Errorf("unknown configuration key '%s'.\nAvailable keys: %s", key, strings.Join(availableKeys, ", "))
	}

	cfg.Profiles[cfg.CurrentProfile] = profile
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}
	return nil
}

// init function registers the setCmd with the profileCmd.
func init() {
	profileCmd.AddCommand(setCmd)
}
