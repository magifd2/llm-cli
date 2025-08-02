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
	Run: func(cmd *cobra.Command, args []string) {
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
			fmt.Fprintf(os.Stderr, "Error: No user prompt provided. Please use --user-prompt, --user-prompt-file, provide a positional argument, or pipe input to stdin.\n")
			os.Exit(1)
		}

		// 3. Load configuration and determine the active LLM provider.
		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		profileName, _ := cmd.Flags().GetString("profile")
		var activeProfile config.Profile
		var ok bool

		// If a specific profile is requested via flag, use it; otherwise, use the current active profile.
		if profileName != "" {
			activeProfile, ok = cfg.Profiles[profileName]
			if !ok {
				fmt.Fprintf(os.Stderr, "Error: Profile '%s' not found.\n", profileName)
				os.Exit(1)
			}
		} else {
			activeProfile, ok = cfg.Profiles[cfg.CurrentProfile]
			if !ok {
				fmt.Fprintf(os.Stderr, "Error: Active profile '%s' not found.\n", cfg.CurrentProfile)
				os.Exit(1)
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

		// 4. Get and print the LLM response, either streaming or as a single response.
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

			// Read from the response channel until it's closed and print tokens.
			for token := range responseChan {
				fmt.Print(token)
			}

			// Wait for the goroutine to finish completely and check for any errors.
			wg.Wait()
			close(errChan)

			if err := <-errChan; err != nil {
				fmt.Fprintf(os.Stderr, "\nError: %v\n", err)
				os.Exit(1)
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
				fmt.Fprintf(os.Stderr, "Error getting response: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(response)
		}
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
}
