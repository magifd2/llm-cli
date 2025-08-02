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
	Long:  `Adds a new profile. If no specific parameters are provided, it copies settings from the default profile. Otherwise, it creates a new profile with the specified parameters.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		profileName := args[0]

		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		if _, ok := cfg.Profiles[profileName]; ok {
			fmt.Fprintf(os.Stderr, "Error: Profile '%s' already exists\n", profileName)
			os.Exit(1)
		}

		newProfile := config.Profile{}
		// If no flags are provided, copy from default profile
		if !cmd.Flags().Changed("provider") &&
			!cmd.Flags().Changed("model") &&
			!cmd.Flags().Changed("endpoint") &&
			!cmd.Flags().Changed("api-key") &&
			!cmd.Flags().Changed("aws-region") &&
			!cmd.Flags().Changed("aws-access-key-id") &&
			!cmd.Flags().Changed("aws-secret-access-key") &&
			!cmd.Flags().Changed("project-id") &&
			!cmd.Flags().Changed("location") &&
			!cmd.Flags().Changed("credentials-file") {

			defaultProfile, ok := cfg.Profiles["default"]
			if !ok {
				fmt.Fprintf(os.Stderr, "Error: Default profile not found. Cannot create new profile without parameters.\n")
				os.Exit(1)
			}
			newProfile = defaultProfile
		} else {
			// Populate newProfile with flag values
			provider, _ := cmd.Flags().GetString("provider")
			model, _ := cmd.Flags().GetString("model")
			endpoint, _ := cmd.Flags().GetString("endpoint")
			apiKey, _ := cmd.Flags().GetString("api-key")
			awsRegion, _ := cmd.Flags().GetString("aws-region")
			awsAccessKeyID, _ := cmd.Flags().GetString("aws-access-key-id")
			awsSecretAccessKey, _ := cmd.Flags().GetString("aws-secret-access-key")
			projectID, _ := cmd.Flags().GetString("project-id")
			location, _ := cmd.Flags().GetString("location")
			credentialsFile, _ := cmd.Flags().GetString("credentials-file")

			newProfile.Provider = provider
			newProfile.Model = model
			newProfile.Endpoint = endpoint
			newProfile.APIKey = apiKey
			newProfile.AWSRegion = awsRegion
			newProfile.AWSAccessKeyID = awsAccessKeyID
			newProfile.AWSSecretAccessKey = awsSecretAccessKey
			newProfile.ProjectID = projectID
			newProfile.Location = location
			newProfile.CredentialsFile = credentialsFile
		}

		cfg.Profiles[profileName] = newProfile

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

	addCmd.Flags().String("provider", "", "LLM provider (e.g., ollama, openai, bedrock, vertexai)")
	addCmd.Flags().String("model", "", "Model name (e.g., llama3, gpt-4, gemini-1.5-pro-001)")
	addCmd.Flags().String("endpoint", "", "API endpoint URL")
	addCmd.Flags().String("api-key", "", "API key for the provider")
	addCmd.Flags().String("aws-region", "", "AWS region for Bedrock")
	addCmd.Flags().String("aws-access-key-id", "", "AWS Access Key ID for Bedrock")
	addCmd.Flags().String("aws-secret-access-key", "", "AWS Secret Access Key for Bedrock")
	addCmd.Flags().String("project-id", "", "GCP Project ID for Vertex AI")
	addCmd.Flags().String("location", "", "GCP Location for Vertex AI")
	addCmd.Flags().String("credentials-file", "", "Path to GCP credentials file for Vertex AI")
}
