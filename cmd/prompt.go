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
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/magifd2/llm-cli/internal/config"
	"github.com/magifd2/llm-cli/internal/llm"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

// promptCmd represents the 'prompt' command.
// This command sends a user prompt to the configured LLM provider and prints the response.
var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Send a prompt to the LLM",
	Long:  `Sends a prompt to the configured LLM and prints the response.`, // Corrected: Removed unnecessary escaping of backticks.
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1. Get prompt values from flags.
		userPrompt, _ := cmd.Flags().GetString("user-prompt")
		userPromptFile, _ := cmd.Flags().GetString("user-prompt-file")
		systemPrompt, _ := cmd.Flags().GetString("system-prompt")
		systemPromptFile, _ := cmd.Flags().GetString("system-prompt-file")

		// 2. Load prompts from direct input, file, or stdin.
		userPromptStr := loadPrompt(userPrompt, userPromptFile)
		systemPromptStr := loadPrompt(systemPrompt, systemPromptFile)

		// If no user prompt is provided via flags, check for positional arguments.
		if userPromptStr == "" && len(args) > 0 {
			userPromptStr = args[0] // Take the first positional argument as the prompt.
		}

		// If userPromptStr is still empty, it's an error as a user prompt is mandatory.
		if userPromptStr == "" {
			return fmt.Errorf("no user prompt provided. Please use --user-prompt, --user-prompt-file, provide a positional argument, or pipe input to stdin")
		}

		// 3. Load configuration and determine the active LLM provider.
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("error loading config: %w", err)
		}

		profileName, _ := cmd.Flags().GetString("profile")
		var activeProfile config.Profile
		var ok bool

		// If a specific profile is requested via flag, use it; otherwise, use the current active profile.
		if profileName != "" {
			activeProfile, ok = cfg.Profiles[profileName]
			if !ok {
				return fmt.Errorf("profile '%s' not found", profileName)
			}
		} else {
			activeProfile, ok = cfg.Profiles[cfg.CurrentProfile]
			if !ok {
				return fmt.Errorf("active profile '%s' not found", cfg.CurrentProfile)
			}
		}

		// 4. Apply limits
		if activeProfile.Limits.Enabled {
			onInputExceeded := activeProfile.Limits.OnInputExceeded
			if cmd.Flags().Changed("on-input-exceeded") {
				onInputExceeded, _ = cmd.Flags().GetString("on-input-exceeded")
			}

			promptSize := int64(len(userPromptStr) + len(systemPromptStr))
			if promptSize > activeProfile.Limits.MaxPromptSizeBytes {
				if onInputExceeded == "stop" {
					return fmt.Errorf("input size (%d bytes) exceeds the limit of %d bytes", promptSize, activeProfile.Limits.MaxPromptSizeBytes)
				} else if onInputExceeded == "warn" {
					fmt.Fprintf(os.Stderr, "Warning: Input size (%d bytes) exceeds the limit of %d bytes. Truncating...\n", promptSize, activeProfile.Limits.MaxPromptSizeBytes)
					// Truncate the user prompt
					combinedLen := int64(len(systemPromptStr))
					if combinedLen < activeProfile.Limits.MaxPromptSizeBytes {
						userPromptStr = userPromptStr[:activeProfile.Limits.MaxPromptSizeBytes-combinedLen]
					} else {
						userPromptStr = ""
					}
				}
			}
		}

		var provider llm.Provider
		// Initialize the appropriate LLM provider based on the active profile's provider type.
		switch activeProfile.Provider {
		case "ollama":
			provider = &llm.OllamaProvider{Profile: activeProfile}
		case "openai":
			provider = &llm.OpenAIProvider{Profile: activeProfile}
		case "bedrock":
			// For Bedrock, check the model ID to determine if it's a Nova model.
			// If it's a Nova model, use NovaBedrockProvider; otherwise, use a mock provider for unsupported models.
			if strings.HasPrefix(activeProfile.Model, "amazon.nova") {
				provider = &llm.NovaBedrockProvider{Profile: activeProfile}
			} else {
				fmt.Fprintf(os.Stderr, "Error: Bedrock model '%s' not supported yet. Using mock provider.\n", activeProfile.Model)
				provider = &llm.MockProvider{}
			}
		case "vertexai":
			provider = &llm.VertexAIProvider{Profile: activeProfile}
		default:
			// If the provider is not recognized, default to a mock provider and issue a warning.
			fmt.Fprintf(os.Stderr, "Warning: Provider '%s' not recognized. Using mock provider.\n", activeProfile.Provider)
			provider = &llm.MockProvider{}
		}

		// 5. Get and print the LLM response, either streaming or as a single response.
		stream, _ := cmd.Flags().GetBool("stream")
		if stream {
			var wg sync.WaitGroup
			// Use a buffered channel to prevent the goroutine from blocking if the main goroutine is slow.
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

			var totalResponseSize int64
			// Read from the response channel until it's closed and print tokens.
			for token := range responseChan {
				if activeProfile.Limits.Enabled {
					onOutputExceeded := activeProfile.Limits.OnOutputExceeded
					if cmd.Flags().Changed("on-output-exceeded") {
						onOutputExceeded, _ = cmd.Flags().GetString("on-output-exceeded")
					}

					totalResponseSize += int64(len(token))
					if totalResponseSize > activeProfile.Limits.MaxResponseSizeBytes {
						if onOutputExceeded == "stop" {
							return fmt.Errorf("\nError: Output size exceeded the limit of %d bytes.", activeProfile.Limits.MaxResponseSizeBytes)
						} else if onOutputExceeded == "warn" {
							fmt.Fprintf(os.Stderr, "\nWarning: Output size exceeded the limit of %d bytes. Truncating...\n", activeProfile.Limits.MaxResponseSizeBytes)
							break
						}
					}
				}
				fmt.Print(token)
			}

			// Wait for the goroutine to finish completely and check for any errors.
			wg.Wait()
			close(errChan)

			if err := <-errChan; err != nil {
				return fmt.Errorf("\nError: %w", err)
			}

			fmt.Println() // Print a final newline after the stream ends for clean output.
		} else {
			var response string
			var err error

			// isatty.IsTerminal で、標準出力がターミナルかどうかを判定
			if isatty.IsTerminal(os.Stdout.Fd()) {
				// 【ターミナルの場合：スピナーを表示】
				s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
				s.Suffix = "  Generating response..."
				s.Start()

				// 非同期でレスポンスを取得
				done := make(chan bool)
				go func() {
					response, err = provider.Chat(systemPromptStr, userPromptStr)
					done <- true
				}()
				<-done

				s.Stop()
			} else {
				// 【リダイレクト/パイプの場合：スピナーを表示しない】
				// 単純にレスポンスを待つ
				response, err = provider.Chat(systemPromptStr, userPromptStr)
			}

			if err != nil {
				return fmt.Errorf("Error getting response: %w", err)
			}

			if activeProfile.Limits.Enabled {
				onOutputExceeded := activeProfile.Limits.OnOutputExceeded
				if cmd.Flags().Changed("on-output-exceeded") {
					onOutputExceeded, _ = cmd.Flags().GetString("on-output-exceeded")
				}

				if int64(len(response)) > activeProfile.Limits.MaxResponseSizeBytes {
					if onOutputExceeded == "stop" {
						return fmt.Errorf("Error: Output size (%d bytes) exceeds the limit of %d bytes.", len(response), activeProfile.Limits.MaxResponseSizeBytes)
					} else if onOutputExceeded == "warn" {
						fmt.Fprintf(os.Stderr, "Warning: Output size (%d bytes) exceeds the limit of %d bytes. Truncating...\n", len(response), activeProfile.Limits.MaxResponseSizeBytes)
						response = response[:activeProfile.Limits.MaxResponseSizeBytes]
					}
				}
			}
			fmt.Println(response)
		}
		return nil
	},
}

// loadPrompt determines the prompt string based on direct input, file path, or stdin.
// It prioritizes direct value, then file content (supporting '-' for stdin), and finally checks for piped stdin.
func loadPrompt(directValue, filePath string) string {
	if directValue != "" {
		return directValue
	}
	if filePath != "" {
		var content []byte
		var err error
		// If filePath is "-", read from stdin.
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
	// Check if stdin is being piped. This is a non-critical check.
	stat, err := os.Stdin.Stat()
	if err != nil {
		// If Stat fails, it usually means stdin is not available or not a character device.
		// We can ignore this error and return an empty string, as it's not a critical failure.
		return ""
	}

	// If stdin is not a character device, it means input is being piped.
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

// init function registers the promptCmd with the rootCmd and defines its flags.
func init() {
	rootCmd.AddCommand(promptCmd)

	promptCmd.Flags().StringP("user-prompt", "p", "", "User prompt to send to the LLM")
	promptCmd.Flags().StringP("user-prompt-file", "f", "", "Path to a file containing the user prompt. Use '-' for stdin.")
	promptCmd.Flags().StringP("system-prompt", "P", "", "System prompt to send to the LLM")
	promptCmd.Flags().StringP("system-prompt-file", "F", "", "Path to a file containing the system prompt.")
	promptCmd.Flags().Bool("stream", false, "Enable streaming response")
	promptCmd.Flags().String("profile", "", "Use a specific profile for this command (overrides current active profile)")

	// Flags for limits
	promptCmd.Flags().String("on-input-exceeded", "", "Action on input size limit exceeded (stop or warn)")
	promptCmd.Flags().String("on-output-exceeded", "", "Action on output size limit exceeded (stop or warn)")
}
