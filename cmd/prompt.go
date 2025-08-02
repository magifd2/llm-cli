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
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/briandowns/spinner"
	"github.com/magifd2/llm-cli/internal/config"
	"github.com/magifd2/llm-cli/internal/llm"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

// promptCmd represents the 'prompt' command.
var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Send a prompt to the LLM",
	Long:  `Sends a prompt to the configured LLM and prints the response.`, // Corrected: Removed unnecessary escaping of backticks
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1. Load configuration and determine the active profile.
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("error loading config: %w", err)
		}

		profileName, _ := cmd.Flags().GetString("profile")
		var activeProfile config.Profile
		var ok bool

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

		// 2. Determine limit settings from profile and flags.
		limits := activeProfile.Limits
		onInputExceeded := limits.OnInputExceeded
		if cmd.Flags().Changed("on-input-exceeded") {
			onInputExceeded, _ = cmd.Flags().GetString("on-input-exceeded")
		}
		onOutputExceeded := limits.OnOutputExceeded
		if cmd.Flags().Changed("on-output-exceeded") {
			onOutputExceeded, _ = cmd.Flags().GetString("on-output-exceeded")
		}

		// 3. Load and validate prompts.
		userPrompt, _ := cmd.Flags().GetString("user-prompt")
		userPromptFile, _ := cmd.Flags().GetString("user-prompt-file")
		systemPrompt, _ := cmd.Flags().GetString("system-prompt")
		systemPromptFile, _ := cmd.Flags().GetString("system-prompt-file")

		userPromptStr, err := loadPrompt(userPrompt, userPromptFile, limits, onInputExceeded)
		if err != nil {
			return err
		}
		systemPromptStr, err := loadPrompt(systemPrompt, systemPromptFile, limits, onInputExceeded)
		if err != nil {
			return err
		}

		if userPromptStr == "" && len(args) > 0 {
			userPromptStr = args[0]
		}

		if userPromptStr == "" {
			return fmt.Errorf("no user prompt provided")
		}

		// 4. Initialize provider.
		var provider llm.Provider
		switch activeProfile.Provider {
		case "ollama":
			provider = &llm.OllamaProvider{Profile: activeProfile}
		case "openai":
			provider = &llm.OpenAIProvider{Profile: activeProfile}
		case "bedrock":
			if strings.HasPrefix(activeProfile.Model, "amazon.nova") {
				provider = &llm.NovaBedrockProvider{Profile: activeProfile}
			} else {
				fmt.Fprintf(os.Stderr, "Error: Bedrock model '%s' not supported yet. Using mock provider.\n", activeProfile.Model)
				provider = &llm.MockProvider{}
			}
		case "vertexai":
			provider = &llm.VertexAIProvider{Profile: activeProfile}
		default:
			fmt.Fprintf(os.Stderr, "Warning: Provider '%s' not recognized. Using mock provider.\n", activeProfile.Provider)
			provider = &llm.MockProvider{}
		}

		// 5. Execute and get response.
		stream, _ := cmd.Flags().GetBool("stream")
		if stream {
			return handleStreamResponse(cmd, provider, systemPromptStr, userPromptStr, activeProfile, onOutputExceeded)
		} else {
			return handleSingleResponse(provider, systemPromptStr, userPromptStr, activeProfile, onOutputExceeded)
		}
	},
}

func handleSingleResponse(provider llm.Provider, systemPrompt, userPrompt string, profile config.Profile, onOutputExceeded string) error {
	var response string
	var err error

	var s *spinner.Spinner
	if isatty.IsTerminal(os.Stdout.Fd()) {
		s = spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = "  Generating response..."
		s.Start()
	}

	response, err = provider.Chat(systemPrompt, userPrompt)

	if s != nil {
		s.Stop()
	}

	if err != nil {
		return fmt.Errorf("error getting response: %w", err)
	}

	// Sanitize and check output size limit.
	if profile.Limits.Enabled {
		response = sanitizeUTF8(response, "output")
		if int64(len(response)) > profile.Limits.MaxResponseSizeBytes {
			if onOutputExceeded == "stop" {
				return fmt.Errorf("output size (%d bytes) exceeds the limit of %d bytes", len(response), profile.Limits.MaxResponseSizeBytes)
			} else if onOutputExceeded == "warn" {
				fmt.Fprintf(os.Stderr, "Warning: Output size (%d bytes) exceeds the limit of %d bytes. Truncating...\n", len(response), profile.Limits.MaxResponseSizeBytes)
				response = truncateStringByBytes(response, profile.Limits.MaxResponseSizeBytes)
			}
		}
	}

	fmt.Println(response)
	return nil
}

func handleStreamResponse(cmd *cobra.Command, provider llm.Provider, systemPrompt, userPrompt string, profile config.Profile, onOutputExceeded string) error {
	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	responseChan := make(chan string)

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(responseChan)
		err := provider.ChatStream(cmd.Context(), systemPrompt, userPrompt, responseChan)
		if err != nil {
			errChan <- err
		}
	}()

	var totalResponseSize int64
	var truncated bool
	for token := range responseChan {
		sanitizedToken := sanitizeUTF8(token, "output")

		if profile.Limits.Enabled && !truncated {
			if totalResponseSize+int64(len(sanitizedToken)) > profile.Limits.MaxResponseSizeBytes {
				if onOutputExceeded == "stop" {
					return fmt.Errorf("\nError: Output size exceeded the limit of %d bytes", profile.Limits.MaxResponseSizeBytes)
				} else if onOutputExceeded == "warn" {
					remainingBytes := profile.Limits.MaxResponseSizeBytes - totalResponseSize
					fmt.Print(truncateStringByBytes(sanitizedToken, remainingBytes))
					fmt.Fprintf(os.Stderr, "\nWarning: Output size exceeded the limit of %d bytes. Truncating...\n", profile.Limits.MaxResponseSizeBytes)
					truncated = true
					break
				}
			}
		}
		totalResponseSize += int64(len(sanitizedToken))
		fmt.Print(sanitizedToken)
	}

	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		return fmt.Errorf("\nError: %w", err)
	}

	if !truncated {
		fmt.Println()
	}
	return nil
}

func loadPrompt(directValue, filePath string, limits config.Limits, onExceeded string) (string, error) {
	if directValue != "" {
		return handlePromptData([]byte(directValue), "argument", limits, onExceeded)
	}
	if filePath != "" {
		if filePath == "-" {
			return readAndProcessStream(os.Stdin, "stdin", limits, onExceeded)
		}
		file, err := os.Open(filePath)
		if err != nil {
			return "", fmt.Errorf("error opening prompt file: %w", err)
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			return "", fmt.Errorf("error getting file stats: %w", err)
		}

		if limits.Enabled && stat.Size() > limits.MaxPromptSizeBytes {
			if onExceeded == "stop" {
				return "", fmt.Errorf("input file size (%d bytes) exceeds the limit of %d bytes", stat.Size(), limits.MaxPromptSizeBytes)
			} else if onExceeded == "warn" {
				fmt.Fprintf(os.Stderr, "Warning: Input file size (%d bytes) exceeds the limit of %d bytes. Reading up to the limit...\n", stat.Size(), limits.MaxPromptSizeBytes)
			}
		}
		return readAndProcessStream(file, fmt.Sprintf("file '%s'", filePath), limits, onExceeded)
	}

	stat, err := os.Stdin.Stat()
	if err != nil || (stat.Mode()&os.ModeCharDevice) != 0 {
		return "", nil // No pipe, no problem
	}

	return readAndProcessStream(os.Stdin, "stdin", limits, onExceeded)
}

func readAndProcessStream(r io.Reader, source string, limits config.Limits, onExceeded string) (string, error) {
	reader := bufio.NewReader(r)
	var buf bytes.Buffer
	var totalBytes int64

	for {
		chunk := make([]byte, 4096)
		n, err := reader.Read(chunk)
		if n > 0 {
			totalBytes += int64(n)
			if limits.Enabled && totalBytes > limits.MaxPromptSizeBytes {
				if onExceeded == "stop" {
					return "", fmt.Errorf("input from %s exceeds size limit of %d bytes", source, limits.MaxPromptSizeBytes)
				}
				// For warn, we just stop reading and will truncate later
				buf.Write(chunk[:n])
				break
			}
			buf.Write(chunk[:n])
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("error reading from %s: %w", source, err)
		}
	}

	return handlePromptData(buf.Bytes(), source, limits, onExceeded)
}

func handlePromptData(data []byte, source string, limits config.Limits, onExceeded string) (string, error) {
	// 1. Sanitize
	sanitizedStr := sanitizeUTF8(string(data), source)

	// 2. Check size and truncate if needed
	if limits.Enabled && int64(len(sanitizedStr)) > limits.MaxPromptSizeBytes {
		if onExceeded == "warn" {
			// This case is for stdin that was read up to the buffer limit
			fmt.Fprintf(os.Stderr, "Warning: Input from %s exceeds the limit of %d bytes. Truncating...\n", source, limits.MaxPromptSizeBytes)
			return truncateStringByBytes(sanitizedStr, limits.MaxPromptSizeBytes), nil
		}
		// Stop case should have been handled earlier for files, but as a fallback for stdin
		return "", fmt.Errorf("input from %s exceeds size limit of %d bytes", source, limits.MaxPromptSizeBytes)
	}

	return sanitizedStr, nil
}

func sanitizeUTF8(s, source string) string {
	if !utf8.ValidString(s) {
		fmt.Fprintf(os.Stderr, "Warning: Invalid UTF-8 sequence detected in %s. Non-UTF-8 characters will be replaced.\n", source)
		return strings.ToValidUTF8(s, "\uFFFD")
	}
	return s
}

func truncateStringByBytes(s string, maxBytes int64) string {
	if int64(len(s)) <= maxBytes {
		return s
	}
	// Find the last rune start position that is within the limit
	endIndex := 0
	for i := range s {
		if int64(i) > maxBytes {
			return s[:endIndex]
		}
		endIndex = i
	}
	return s[:endIndex]
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
