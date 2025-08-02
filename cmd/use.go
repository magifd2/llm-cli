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
	os "os"

	"github.com/magifd2/llm-cli/internal/config"
	"github.com/spf13/cobra"
)

// useCmd represents the 'profile use' command.
// This command sets the specified profile as the currently active profile.
var useCmd = &cobra.Command{
	Use:   "use [profile_name]",
	Short: "Set the active profile",
	Long:  `Set the active profile for llm-cli.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		profileName := args[0]
		if err := useProfile(profileName); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Switched to profile: %s\n", profileName)
	},
}

// useProfile contains the core logic for switching the active profile.
// It loads the configuration, validates the profile name, and saves the updated configuration.
func useProfile(profileName string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// Check if the specified profile exists.
	if _, ok := cfg.Profiles[profileName]; !ok {
		return fmt.Errorf("profile '%s' not found", profileName)
	}

	cfg.CurrentProfile = profileName
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}
	return nil
}

// init function registers the useCmd with the profileCmd.
func init() {
	profileCmd.AddCommand(useCmd)
}
