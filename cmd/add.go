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

	"github.com/magifd2/llm-cli/internal/config"
	"github.com/spf13/cobra"
)

// addCmd represents the 'profile add' command.
// This command allows users to add a new profile for LLM configurations.
var addCmd = &cobra.Command{
	Use:   "add [profile_name]",
	Short: "Add a new profile",
	Long:  `Adds a new profile. If no specific parameters are provided, it copies settings from the default profile. Otherwise, it creates a new profile with the specified parameters.`, 
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profileName := args[0]

		cfg, err := config.Load(cfgFile)
		if err != nil {
			return fmt.Errorf("Error loading config: %w", err)
		}

		if _, ok := cfg.Profiles[profileName]; ok {
			return fmt.Errorf("Error: Profile '%s' already exists", profileName)
		}

		newProfile := config.Profile{}
		// If no specific flags are provided, the new profile copies settings from the default profile.
		// Otherwise, it creates a new profile using the values provided by the flags.
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
				return fmt.Errorf("Error: Default profile not found. Cannot create new profile without parameters.")
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

		// Populate limits with flag values, or use defaults
		limitsEnabled, _ := cmd.Flags().GetBool("limits-enabled")
		onInputExceeded, _ := cmd.Flags().GetString("limits-on-input-exceeded")
		onOutputExceeded, _ := cmd.Flags().GetString("limits-on-output-exceeded")
		maxPromptSizeBytes, _ := cmd.Flags().GetInt64("limits-max-prompt-size-bytes")
		maxResponseSizeBytes, _ := cmd.Flags().GetInt64("limits-max-response-size-bytes")

		newProfile.Limits = config.Limits{
			Enabled:              limitsEnabled,
			OnInputExceeded:      onInputExceeded,
			OnOutputExceeded:     onOutputExceeded,
			MaxPromptSizeBytes:   maxPromptSizeBytes,
			MaxResponseSizeBytes: maxResponseSizeBytes,
		}

		cfg.Profiles[profileName] = newProfile

		if err := cfg.Save(cfgFile); err != nil {
			return fmt.Errorf("Error saving config: %w", err)
		}

		fmt.Printf("Profile '%s' added.\n", profileName)
		fmt.Printf("To switch to the new profile, run: llm-cli profile use %s\n", profileName)
		return nil
	},
}

// init function registers the addCmd with the profileCmd and defines its flags.
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

	// Flags for limits
	addCmd.Flags().Bool("limits-enabled", true, "Enable limits for the profile")
	addCmd.Flags().String("limits-on-input-exceeded", "stop", "Action on input size limit exceeded (stop or warn)")
	addCmd.Flags().String("limits-on-output-exceeded", "stop", "Action on output size limit exceeded (stop or warn)")
	addCmd.Flags().Int64("limits-max-prompt-size-bytes", 10485760, "Max prompt size in bytes (10MB)")
	addCmd.Flags().Int64("limits-max-response-size-bytes", 20971520, "Max response size in bytes (20MB)")
}
