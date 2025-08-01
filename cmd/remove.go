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

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove [profile_name]",
	Short: "Remove a profile",
	Long:  `Removes a specified profile from the configuration.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		profileName := args[0]
		if err := removeProfile(profileName); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Profile '%s' removed.\n", profileName)
	},
}

func removeProfile(profileName string) error {
	if profileName == "default" {
		return fmt.Errorf("the 'default' profile cannot be removed")
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	if _, ok := cfg.Profiles[profileName]; !ok {
		return fmt.Errorf("profile '%s' not found", profileName)
	}

	if cfg.CurrentProfile == profileName {
		return fmt.Errorf("cannot remove the currently active profile. Please switch to another profile first")
	}

	delete(cfg.Profiles, profileName)

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}
	return nil
}

func init() {
	profileCmd.AddCommand(removeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
