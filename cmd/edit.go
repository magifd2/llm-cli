/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/magifd2/llm-cli/internal/config"
	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit the configuration file",
	Long:  `Opens the configuration file in the default editor ($EDITOR).`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runEditCommand(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// runEditCommand contains the core logic for the edit command.
func runEditCommand() error {
	editorEnv := os.Getenv("EDITOR")
	if editorEnv == "" {
		editorEnv = "vim"
	}

	// Find the absolute path of the editor executable to prevent command injection.
	editorPath, err := exec.LookPath(editorEnv)
	if err != nil {
		// If the primary editor is not found, try nano as a fallback.
		editorPath, err = exec.LookPath("nano")
		if err != nil {
			return fmt.Errorf("EDITOR environment variable not set, and vim/nano not found in PATH")
		}
	}

	configPath, err := config.GetConfigPath()
	if err != nil {
		return fmt.Errorf("getting config path: %w", err)
	}

	// Ensure the directory exists before trying to open the file
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	execCmd := exec.Command(editorPath, configPath)
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	if err := execCmd.Run(); err != nil {
		return fmt.Errorf("opening editor: %w", err)
	}
	return nil
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
