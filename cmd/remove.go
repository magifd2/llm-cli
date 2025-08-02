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

	"github.com/magifd2/llm-cli/internal/config"
	"github.com/spf13/cobra"
)

// removeCmd represents the 'profile remove' command.
// This command removes a specified profile from the configuration.
var removeCmd = &cobra.Command{
	Use:   "remove [profile_name]",
	Short: "Remove a profile",
	Long:  `Removes a specified profile from the configuration.`, 
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profileName := args[0]
		if err := removeProfile(profileName); err != nil {
			return fmt.Errorf("Error: %w", err)
		}
		fmt.Printf("Profile '%s' removed.\n", profileName)
		return nil
	},
}

// removeProfile contains the core logic for removing a profile.
// It prevents removal of the 'default' profile and the currently active profile.
func removeProfile(profileName string) error {
	// Prevent removal of the default profile.
	if profileName == "default" {
		return fmt.Errorf("the 'default' profile cannot be removed")
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// Check if the profile exists.
	if _, ok := cfg.Profiles[profileName]; !ok {
		return fmt.Errorf("profile '%s' not found", profileName)
	}

	// Prevent removal of the currently active profile.
	if cfg.CurrentProfile == profileName {
		return fmt.Errorf("cannot remove the currently active profile. Please switch to another profile first")
	}

	delete(cfg.Profiles, profileName)

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}
	return nil
}

// init function registers the removeCmd with the profileCmd.
func init() {
	profileCmd.AddCommand(removeCmd)
}