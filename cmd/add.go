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

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [profile_name]",
	Short: "Add a new profile",
	Long:  `Adds a new profile by copying the settings from the default profile.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		profileName := args[0]

		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		if _, ok := cfg.Profiles[profileName]; ok {
			fmt.Fprintf(os.Stderr, "Error: Profile '%s' already exists.\n", profileName)
			os.Exit(1)
		}

		// Copy settings from the default profile
		defaultProfile, ok := cfg.Profiles["default"]
		if !ok {
			// This should not happen if Load() works correctly
			fmt.Fprintf(os.Stderr, "Error: Default profile not found.\n")
			os.Exit(1)
		}

		cfg.Profiles[profileName] = defaultProfile

		if err := cfg.Save(); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Profile '%s' added.\n", profileName)
		fmt.Printf("To switch to the new profile, run: llm-cli profile use %s\n", profileName)
	},
}

func init() {
	profileCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
