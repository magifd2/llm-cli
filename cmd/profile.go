/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
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

func init() {
	rootCmd.AddCommand(profileCmd)
}
