/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/magifd2/llm-cli/internal/config"
	"github.com/spf13/cobra"
)

// profileCmd represents the profile command
var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage configuration profiles",
	Long:  `The profile command and its subcommands help you manage different configurations for various LLM providers.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Show help if no subcommand is given
		_ = cmd.Help()
	},
}

// showCmd represents the show command
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

func init() {
	rootCmd.AddCommand(profileCmd)
	profileCmd.AddCommand(showCmd)
}
