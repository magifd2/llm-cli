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
	"path/filepath"
	"time"

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

// checkCmd represents the 'profile check' command.
// It checks configuration profiles for consistency and offers to migrate them.
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check and migrate configuration profiles",
	Long: `Checks all configuration profiles for consistency, especially for newly introduced settings like 'limits'.
If a profile's settings are found to be at their default zero values (indicating they might be from an older version or not explicitly set),

the command will prompt to update them to the current standard default values.`, 
	RunE: func(cmd *cobra.Command, args []string) error {
		confirm, _ := cmd.Flags().GetBool("confirm")

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("error loading config: %w", err)
		}

		modified := false
		for name, profile := range cfg.Profiles {
			// Check for missing credentials file
			if profile.CredentialsFile != "" {
				resolvedPath, err := config.ResolvePath(profile.CredentialsFile)
				if err != nil {
					fmt.Printf("Profile '%s': Error resolving credentials file path '%s': %v\n", name, profile.CredentialsFile, err)
				} else {
					_, err := os.Stat(resolvedPath)
					if os.IsNotExist(err) {
						fmt.Printf("Profile '%s': Credentials file '%s' (resolved to '%s') does not exist.\n", name, profile.CredentialsFile, resolvedPath)
					} else if err != nil {
						fmt.Printf("Profile '%s': Error checking credentials file '%s' (resolved to '%s'): %v\n", name, profile.CredentialsFile, resolvedPath, err)
					}
				}
			}

			// Define the standard default limits for comparison
			standardDefaultLimits := config.Limits{
				Enabled:              true,
				OnInputExceeded:      "stop",
				OnOutputExceeded:     "stop",
				MaxPromptSizeBytes:   10485760, // 10MB
				MaxResponseSizeBytes: 20971520, // 20MB
			}

			// Check if current profile's limits are different from standard defaults
			// or if they are the zero value (indicating they were never explicitly set)
			if profile.Limits != standardDefaultLimits {
				// If limits are the zero value, or if they are different from standard defaults,
				// and they are not already explicitly set to something else, prompt for update.
				// This condition ensures we don't prompt if user has intentionally set custom limits.
				if profile.Limits == (config.Limits{}) || (profile.Limits.Enabled == standardDefaultLimits.Enabled &&
					 profile.Limits.OnInputExceeded == standardDefaultLimits.OnInputExceeded &&
					 profile.Limits.OnOutputExceeded == standardDefaultLimits.OnOutputExceeded &&
					 profile.Limits.MaxPromptSizeBytes == standardDefaultLimits.MaxPromptSizeBytes &&
					 profile.Limits.MaxResponseSizeBytes == standardDefaultLimits.MaxResponseSizeBytes) {
					
					fmt.Printf("Profile '%s' has default or unconfigured 'limits' settings.\n", name)
					if !confirm {
						fmt.Printf("Do you want to update them to standard default values? (y/N): ")
						var response string
						if _, err := fmt.Scanln(&response); err != nil {
							// Handle EOF as a 'No' answer
							if err.Error() == "unexpected newline" || err.Error() == "EOF" {
								fmt.Println("Skipping profile due to no input.")
								continue
							}
							return fmt.Errorf("failed to read response: %w", err)
						}
						if ! (response == "y" || response == "Y") {
							fmt.Printf("Skipping profile '%s'.\n", name)
							continue
						}
					}
					profile.Limits = standardDefaultLimits
					cfg.Profiles[name] = profile
					modified = true
					fmt.Printf("Profile '%s' 'limits' updated to standard defaults.\n", name)
				} else {
					fmt.Printf("Profile '%s' 'limits' settings are configured.\n", name)
				}
			} else {
				fmt.Printf("Profile '%s' 'limits' settings are up-to-date.\n", name)
			}
		}

		if modified {
			fmt.Println("Configuration changes detected.")
			if !confirm {
				fmt.Printf("Do you want to save the changes? (y/N): ")
				var response string
				if _, err := fmt.Scanln(&response); err != nil {
					// Handle EOF as a 'No' answer
					if err.Error() == "unexpected newline" || err.Error() == "EOF" {
						fmt.Println("Changes not saved due to no input.")
						return nil
					}
					return fmt.Errorf("failed to read response: %w", err)
				}
				if ! (response == "y" || response == "Y") {
					fmt.Println("Changes not saved.")
					return nil
				}
			}

			// Backup before saving
			if err := backupConfigFile(); err != nil {
				return fmt.Errorf("failed to backup config file: %w", err)
			}
			fmt.Println("Configuration file backed up.")

			if err := cfg.Save(); err != nil {
				return fmt.Errorf("error saving config: %w", err)
			}
			fmt.Println("Configuration saved successfully.")
		} else {
			fmt.Println("All profiles are up-to-date. No changes needed.")
		}

		return nil
	},
}

// backupConfigFile creates a timestamped backup of the config.json file.
func backupConfigFile() error {
	configPath, err := config.GetConfigPath()
	if err != nil {
		return fmt.Errorf("could not get config path: %w", err)
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil // No config file to backup
	}

	backupDir := filepath.Join(filepath.Dir(configPath), "backups")
	if err := os.MkdirAll(backupDir, 0700); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	backupFileName := fmt.Sprintf("config_%s.json.bak", timestamp)
	backupPath := filepath.Join(backupDir, backupFileName)

	input, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file for backup: %w", err)
	}

	err = os.WriteFile(backupPath, input, 0600)
	if err != nil {
		return fmt.Errorf("failed to write backup file: %w", err)
	}

	return nil
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
		resolvedPath, err := config.ResolvePath(profile.CredentialsFile)
		if err != nil {
			fmt.Printf("  CredentialsFile: %s (Error resolving path: %v)\n", profile.CredentialsFile, err)
		} else {
			fmt.Printf("  CredentialsFile: %s (Resolved: %s)\n", profile.CredentialsFile, resolvedPath)
		}
	}
	// Display Limits if enabled or if any limit is non-zero/non-empty
	if profile.Limits.Enabled ||
		profile.Limits.OnInputExceeded != "" ||
		profile.Limits.OnOutputExceeded != "" ||
		profile.Limits.MaxPromptSizeBytes != 0 ||
		profile.Limits.MaxResponseSizeBytes != 0 {
		fmt.Printf("  Limits:\n")
		fmt.Printf("    Enabled: %t\n", profile.Limits.Enabled)
		fmt.Printf("    OnInputExceeded: %s\n", profile.Limits.OnInputExceeded)
		fmt.Printf("    OnOutputExceeded: %s\n", profile.Limits.OnOutputExceeded)
		fmt.Printf("    MaxPromptSizeBytes: %d\n", profile.Limits.MaxPromptSizeBytes)
		fmt.Printf("    MaxResponseSizeBytes: %d\n", profile.Limits.MaxResponseSizeBytes)
	}
}

// init function registers the profileCmd with the rootCmd and adds the showCmd and checkCmd as subcommands.
func init() {
	rootCmd.AddCommand(profileCmd)
	profileCmd.AddCommand(showCmd)
	profileCmd.AddCommand(checkCmd) // Register the new checkCmd

	// Add flags for checkCmd
	checkCmd.Flags().BoolP("confirm", "y", false, "Confirm all prompts automatically (non-interactive)")
}