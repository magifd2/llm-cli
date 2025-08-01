/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/magifd2/llm-cli/internal/config"
	"github.com/magifd2/llm-cli/internal/llm"
	"github.com/spf13/cobra"
)

// promptCmd represents the prompt command
var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Send a prompt to the LLM",
	Long:  `Sends a prompt to the configured LLM and prints the response.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Get prompt values
		userPrompt, _ := cmd.Flags().GetString("user-prompt")
		userPromptFile, _ := cmd.Flags().GetString("user-prompt-file")
		systemPrompt, _ := cmd.Flags().GetString("system-prompt")
		systemPromptFile, _ := cmd.Flags().GetString("system-prompt-file")

		// 2. Load prompts
		userPromptStr := loadPrompt(userPrompt, userPromptFile)
		systemPromptStr := loadPrompt(systemPrompt, systemPromptFile)

		// If no user prompt is provided via flags or stdin, check for positional arguments
		if userPromptStr == "" && len(args) > 0 {
			userPromptStr = args[0] // Take the first positional argument as the prompt
		}

		// If userPromptStr is still empty, it's an error.
		if userPromptStr == "" {
			fmt.Fprintf(os.Stderr, "Error: No user prompt provided. Please use --user-prompt, --user-prompt-file, provide a positional argument, or pipe input to stdin.\n")
			os.Exit(1)
		}

		// 3. Get LLM provider
        cfg, err := config.Load()
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
            os.Exit(1)
        }

        activeProfile, ok := cfg.Profiles[cfg.CurrentProfile]
        if !ok {
            fmt.Fprintf(os.Stderr, "Error: Active profile '%s' not found.\n", cfg.CurrentProfile)
            os.Exit(1)
        }

        var provider llm.Provider
        switch activeProfile.Provider {
        case "ollama":
            provider = &llm.OllamaProvider{Profile: activeProfile}
        case "openai":
            provider = &llm.OpenAIProvider{Profile: activeProfile}
        case "bedrock":
            // Check the model ID to determine which Bedrock provider to use
            if strings.HasPrefix(activeProfile.Model, "amazon.nova") {
                provider = &llm.NovaBedrockProvider{Profile: activeProfile}
            } else {
                // Fallback for other Bedrock models not yet implemented
                fmt.Fprintf(os.Stderr, "Error: Bedrock model '%s' not supported yet. Using mock provider.\n", activeProfile.Model)
                provider = &llm.MockProvider{}
            }
        default:
            // For now, default to mock provider if not ollama
            fmt.Fprintf(os.Stderr, "Warning: Provider '%s' not recognized. Using mock provider.\n", activeProfile.Provider)
            provider = &llm.MockProvider{}
        }

        // 4. Get and print response
        stream, _ := cmd.Flags().GetBool("stream")
        if stream {
            var wg sync.WaitGroup
            // Use a buffered channel to prevent the goroutine from blocking.
            errChan := make(chan error, 1)
            responseChan := make(chan string)

            wg.Add(1)
            go func() {
                defer wg.Done()
                defer close(responseChan)
                err := provider.ChatStream(cmd.Context(), systemPromptStr, userPromptStr, responseChan)
                if err != nil {
                    errChan <- err
                }
            }()

            // Read from the response channel until it's closed
            for token := range responseChan {
                fmt.Print(token)
            }

            // Wait for the goroutine to finish completely
            wg.Wait()
            close(errChan)

            // After the goroutine is done, check for any error it might have sent.
            if err := <-errChan; err != nil {
                fmt.Fprintf(os.Stderr, "\nError: %v\n", err)
                os.Exit(1)
            }

            fmt.Println() // Print a final newline after the stream ends
        } else {
            response, err := provider.Chat(systemPromptStr, userPromptStr)
            if err != nil {
                fmt.Fprintf(os.Stderr, "Error getting response: %v\n", err)
                os.Exit(1)
            }
            fmt.Println(response)
        }
    },
}

// loadPrompt determines the prompt string based on direct input, file path, or stdin.
func loadPrompt(directValue, filePath string) string {
	if directValue != "" {
		return directValue
	}
	if filePath != "" {
		var content []byte
		var err error
		if filePath == "-" {
			content, err = io.ReadAll(os.Stdin)
		} else {
			content, err = os.ReadFile(filePath)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading prompt file: %v\n", err)
			os.Exit(1)
		}
		return string(content)
	}
	// Check if stdin is being piped
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		content, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
			os.Exit(1)
		}
		return string(content)
	}

	return ""
}

func init() {
	rootCmd.AddCommand(promptCmd)

	promptCmd.Flags().StringP("user-prompt", "p", "", "User prompt to send to the LLM")
	promptCmd.Flags().StringP("user-prompt-file", "f", "", "Path to a file containing the user prompt. Use '-' for stdin.")
	promptCmd.Flags().StringP("system-prompt", "P", "", "System prompt to send to the LLM")
	promptCmd.Flags().StringP("system-prompt-file", "F", "", "Path to a file containing the system prompt.")
	promptCmd.Flags().Bool("stream", false, "Enable streaming response")
}