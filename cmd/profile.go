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

// profileCmd represents the base command for managing configuration profiles.
// It serves as a container for subcommands like add, list, use, set, remove, and edit.
var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage configuration profiles",
	Long:  `The profile command and its subcommands help you manage different configurations for various LLM providers.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Show help if no subcommand is given
		_ = cmd.Help()
	},
}

// showCmd represents the 'profile show' command.
// This command displays the detailed configuration of a specified profile.
var showCmd = &cobra.Command{
	Use:   "show [profile_name]",
	Short: "Show details of a specific profile",
	Long:  `Shows the detailed configuration for a specified profile. If no profile name is given, it shows the current active profile.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		profileName := cfg.CurrentProfile
		if len(args) > 0 {
			profileName = args[0]
		}

		profile, ok := cfg.Profiles[profileName]
		if !ok {
			fmt.Fprintf(os.Stderr, "Error: Profile '%s' not found.\n", profileName)
			os.Exit(1)
		}

		showProfile(profile, profileName)
	},
}

// showProfile prints the details of a given profile to the console.
func showProfile(profile config.Profile, name string) {
	fmt.Printf("Profile: %s\n", name)
	fmt.Printf("  Provider: %s\n", profile.Provider)
	fmt.Printf("  Model: %s\n", profile.Model)
	if profile.Endpoint != "" {
		fmt.Printf("  Endpoint: %s\n", profile.Endpoint)
	}
	if profile.APIKey != "" {
		fmt.Printf("  APIKey: %s\n", profile.APIKey)
	}
	if profile.AWSRegion != "" {
		fmt.Printf("  AWSRegion: %s\n", profile.AWSRegion)
	}
	if profile.AWSAccessKeyID != "" {
		fmt.Printf("  AWSAccessKeyID: %s\n", profile.AWSAccessKeyID)
	}
	if profile.AWSSecretAccessKey != "" {
		fmt.Printf("  AWSSecretAccessKey: %s\n", profile.AWSSecretAccessKey)
	}
	if profile.ProjectID != "" {
		fmt.Printf("  ProjectID: %s\n", profile.ProjectID)
	}
	if profile.Location != "" {
		fmt.Printf("  Location: %s\n", profile.Location)
	}
	if profile.CredentialsFile != "" {
		fmt.Printf("  CredentialsFile: %s\n", profile.CredentialsFile)
	}
}

// init function registers the profileCmd with the rootCmd and adds the showCmd as a subcommand.
func init() {
	rootCmd.AddCommand(profileCmd)
	profileCmd.AddCommand(showCmd)
}
