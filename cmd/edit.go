/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit the configuration file",
	Long:  `Opens the configuration file in the default editor ($EDITOR).`,
	Run: func(cmd *cobra.Command, args []string) {
		editor := os.Getenv("EDITOR")
		if editor == "" {
			// Fallback to vim or nano if EDITOR is not set
			if _, err := exec.LookPath("vim"); err == nil {
				editor = "vim"
			} else if _, err := exec.LookPath("nano"); err == nil {
				editor = "nano"
			} else {
				fmt.Fprintf(os.Stderr, "Error: EDITOR environment variable not set, and vim/nano not found.\n")
				os.Exit(1)
			}
		}

		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
			os.Exit(1)
		}
		// Note: This path needs to be in sync with internal/config/config.go
		configPath := filepath.Join(home, ".config", "llm-cli", "config.json")

		// Ensure the directory exists before trying to open the file
		if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating config directory: %v\n", err)
			os.Exit(1)
		}

		execCmd := exec.Command(editor, configPath)
		execCmd.Stdin = os.Stdin
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr

		if err := execCmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error opening editor: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	profileCmd.AddCommand(editCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// editCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// editCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
