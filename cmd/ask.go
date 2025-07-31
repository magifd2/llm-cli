/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/magifd2/llm-cli/internal/config"
	"github.com/magifd2/llm-cli/internal/llm"
	"github.com/spf13/cobra"
)

// askCmd represents the ask command
var askCmd = &cobra.Command{
	Use:   "ask",
	Short: "Send a prompt to the LLM",
	Long:  `Sends a prompt to the configured LLM and prints the response.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Get prompt values
		prompt, _ := cmd.Flags().GetString("prompt")
		promptFile, _ := cmd.Flags().GetString("prompt-file")
		systemPrompt, _ := cmd.Flags().GetString("system-prompt")
		systemPromptFile, _ := cmd.Flags().GetString("system-prompt-file")

		// 2. Load prompts
		userPromptStr := loadPrompt(prompt, promptFile)
		systemPromptStr := loadPrompt(systemPrompt, systemPromptFile)

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
            provider = &llm.BedrockProvider{Profile: activeProfile}
        default:
            // For now, default to mock provider if not ollama
            fmt.Fprintf(os.Stderr, "Warning: Provider '%s' not recognized. Using mock provider.\n", activeProfile.Provider)
            provider = &llm.MockProvider{}
        }

        // 4. Get and print response
        stream, _ := cmd.Flags().GetBool("stream")
        if stream {
            responseChan := make(chan string)
            go func() {
                err := provider.ChatStream(cmd.Context(), systemPromptStr, userPromptStr, responseChan)
                if err != nil {
                    fmt.Fprintf(os.Stderr, "\nError during streaming: %v\n", err)
                }
            }()

            for token := range responseChan {
                fmt.Print(token)
            }
            fmt.Println()
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
	rootCmd.AddCommand(askCmd)

	askCmd.Flags().String("prompt", "", "Prompt to send to the LLM")
	askCmd.Flags().String("prompt-file", "", "Path to a file containing the prompt. Use '-' for stdin.")
	askCmd.Flags().String("system-prompt", "", "System prompt to send to the LLM")
	askCmd.Flags().String("system-prompt-file", "", "Path to a file containing the system prompt.")
	askCmd.Flags().Bool("stream", false, "Enable streaming response")
}
