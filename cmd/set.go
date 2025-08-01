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

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a value in the current profile",
	Long:  `Set a configuration value for the currently active profile.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		if err := setProfileValue(args[0], args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		// Success message moved here
		fmt.Printf("Set %s = %s in profile %s\n", args[0], args[1], cfg.CurrentProfile)
	},
}

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
	default:
		return fmt.Errorf("unknown configuration key '%s'.\nAvailable keys: model, provider, endpoint, api_key, aws_region, aws_access_key_id, aws_secret_access_key", key)
	}

	cfg.Profiles[cfg.CurrentProfile] = profile
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}
	return nil
}

func init() {
	profileCmd.AddCommand(setCmd)
}
