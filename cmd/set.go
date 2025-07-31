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
		key := args[0]
		value := args[1]

		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		profile, ok := cfg.Profiles[cfg.CurrentProfile]
		if !ok {
			fmt.Fprintf(os.Stderr, "Error: Active profile '%s' not found.\n", cfg.CurrentProfile)
			os.Exit(1)
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
		case "aws_profile_name":
			profile.AWSProfileName = value
		case "aws_access_key_id":
			profile.AWSAccessKeyID = value
		case "aws_secret_access_key":
			profile.AWSSecretAccessKey = value
		default:
			fmt.Fprintf(os.Stderr, "Error: Unknown configuration key '%s'.\n", key)
			os.Exit(1)
		}

		cfg.Profiles[cfg.CurrentProfile] = profile
		if err := cfg.Save(); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Set %s = %s in profile %s\n", key, value, cfg.CurrentProfile)
	},
}

func init() {
	profileCmd.AddCommand(setCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
