package openai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/magifd2/llm-cli/internal/config"
	"github.com/magifd2/llm-cli/internal/llm"
)

// Provider implements the llm.Provider interface for OpenAI-compatible APIs.
// This includes OpenAI's own API and local LLM servers like LM Studio that mimic OpenAI's API.
type Provider struct {
	Profile config.Profile // The configuration profile for this OpenAI-compatible instance.
}

// openAIRequest represents the JSON structure for requests to the OpenAI Chat Completions API.
type openAIRequest struct {
	Model    string         `json:"model"`          // The name of the model to use.
	Messages []llm.Message `json:"messages"`       // A list of messages in the conversation history.
	Stream   bool           `json:"stream,omitempty"` // Whether to stream the response. Omitted if false.
}

// openAIResponse represents the JSON structure for a non-streaming response from the OpenAI API.
type openAIResponse struct {
	Choices []struct {
		Message llm.Message `json:"message"` // The assistant's message.
	} `json:"choices"` // A list of chat completion choices.
}

// Chat sends a chat request to the OpenAI-compatible API and returns a single, complete response.
func (p *Provider) Chat(systemPrompt, userPrompt string) (string, error) {
	// Determine the API endpoint. Use a default if not specified in the profile.
	endpoint := p.Profile.Endpoint
	if endpoint == "" {
		endpoint = "https://api.openai.com/v1/chat/completions"
	}

	// Build the messages array, including the system prompt if provided.
	messages := []llm.Message{}
	if systemPrompt != "" {
		messages = append(messages, llm.Message{Role: "system", Content: systemPrompt})
	}
	messages = append(messages, llm.Message{Role: "user", Content: userPrompt})

	// Construct the request body for a non-streaming chat.
	reqBody := openAIRequest{
		Model:    p.Profile.Model,
		Messages: messages,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshalling request body: %w", err)
	}

	// Create an HTTP request with a background context.
	req, err := http.NewRequestWithContext(context.Background(), "POST", endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	// Determine the API key to use (from file or direct in profile).
	apiKey := p.Profile.APIKey
	if p.Profile.CredentialsFile != "" {
		fileKey, err := loadOpenAIAPIKeyFromFile(p.Profile.CredentialsFile)
		if err != nil {
			// Chat関数は(string, error)を返すため、エラー時はstringも返す
			return "", fmt.Errorf("failed to load OpenAI API key from file %s: %w", p.Profile.CredentialsFile, err)
		}
		apiKey = fileKey
	}

	// Set necessary headers, including Content-Type and Authorization (if API key is provided).
	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer " + apiKey)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request to openai-compatible api: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-OK HTTP status codes and return an error with the response body.
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("openai-compatible api request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Decode the JSON response.
	var openAIResp openAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return "", fmt.Errorf("error decoding openai response: %w", err)
	}

	// Validate if any choices were returned.
	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from openai-compatible api")
	}

	return openAIResp.Choices[0].Message.Content, nil
}

// openAIStreamResponse represents a chunk of the streaming response from the OpenAI API.
type openAIStreamResponse struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"` // The content delta for the current chunk.
		} `json:"delta"` // The change in content.
	} `json:"choices"` // A list of chat completion choices (usually one for streaming).
}

// ChatStream sends a streaming chat request to the OpenAI-compatible API and sends response chunks to a channel.
func (p *Provider) ChatStream(ctx context.Context, systemPrompt, userPrompt string, responseChan chan<- string) error {
	// Note: The caller is responsible for closing the responseChan.

	// Determine the API endpoint. Use a default if not specified in the profile.
	endpoint := p.Profile.Endpoint
	if endpoint == "" {
		endpoint = "https://api.openai.com/v1/chat/completions"
	}

	// Build the messages array, including the system prompt if provided.
	messages := []llm.Message{}
	if systemPrompt != "" {
		messages = append(messages, llm.Message{Role: "system", Content: systemPrompt})
	}
	messages = append(messages, llm.Message{Role: "user", Content: userPrompt})

	// Construct the request body for a streaming chat.
	reqBody := openAIRequest{
		Model:    p.Profile.Model,
		Messages: messages,
		Stream:   true,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("error marshalling request body: %w", err)
	}

	// Create an HTTP request with context for cancellation.
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// Determine the API key to use (from file or direct in profile).
	apiKey := p.Profile.APIKey
	if p.Profile.CredentialsFile != "" {
		fileKey, err := loadOpenAIAPIKeyFromFile(p.Profile.CredentialsFile)
		if err != nil {
			// ChatStream関数はerrorのみを返すため、エラー時はerrorのみ返す
			return fmt.Errorf("failed to load OpenAI API key from file %s: %w", p.Profile.CredentialsFile, err)
		}
		apiKey = fileKey
	}

	// Set necessary headers, including Content-Type and Authorization (if API key is provided).
	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer " + apiKey)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request to openai-compatible api: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-OK HTTP status codes.
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("openai-compatible api request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Read and process the streaming response line by line.
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		// Skip lines that are not data or are the DONE signal.
		if !strings.HasPrefix(line, "data: ") || line == "data: [DONE]" {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")

        // Check for stream errors within the data payload.
        if strings.Contains(data, "\"error\":") {
            return fmt.Errorf("streaming error: %s", data)
        }

        var streamResp openAIStreamResponse
		// Unmarshal each line as a JSON response chunk.
		if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
			// Ignore JSON parsing errors for now, as some streams might contain metadata or empty deltas.
			continue
		}

		// If content delta is present, send it to the response channel.
		if len(streamResp.Choices) > 0 {
			select {
			case responseChan <- streamResp.Choices[0].Delta.Content:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	// Check for any errors during scanning.
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading stream: %w", err)
	}

	return nil
}

// openAIAPIKey represents the structure of the OpenAI API key JSON file.
type openAIAPIKey struct {
	OpenAIAPIKey string `json:"openai_api_key"`
}

// loadOpenAIAPIKeyFromFile loads OpenAI API key from a specified JSON file.
func loadOpenAIAPIKeyFromFile(filePath string) (string, error) {
	resolvedPath, err := config.ResolvePath(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve credentials file path %s: %w", filePath, err)
	}

	data, err := os.ReadFile(resolvedPath)
	if err != nil {
		return "", fmt.Errorf("failed to read credentials file %s: %w", resolvedPath, err)
	}

	var key openAIAPIKey
	if err := json.Unmarshal(data, &key); err != nil {
		return "", fmt.Errorf("failed to unmarshal credentials from file %s: %w", resolvedPath, err)
	}

	if key.OpenAIAPIKey == "" {
		return "", fmt.Errorf("openai_api_key is missing in credentials file %s", resolvedPath)
	}

	return key.OpenAIAPIKey, nil
}

// NewProvider is a factory function that returns a new OpenAI provider.
func NewProvider(p config.Profile) (llm.Provider, error) {
	return &Provider{Profile: p}, nil
}

// ValidateConfig checks if the OpenAI provider's configuration is valid.
// It requires a model and, for the default OpenAI endpoint, an API key or credentials file.
func (p *Provider) ValidateConfig() error {
	// Model is always required.
	if p.Profile.Model == "" {
		return fmt.Errorf("OpenAI provider requires a 'model' to be specified in the profile")
	}

	// Determine the effective endpoint. If empty, it will use the default OpenAI endpoint.
	endpoint := p.Profile.Endpoint
	if endpoint == "" {
		endpoint = "https://api.openai.com/v1/chat/completions"
	}

	// If using the default OpenAI endpoint, an API key or credentials file is mandatory.
	if endpoint == "https://api.openai.com/v1/chat/completions" {
		if p.Profile.APIKey == "" && p.Profile.CredentialsFile == "" {
			return fmt.Errorf("OpenAI provider (using default endpoint) requires either 'api-key' or 'credentials-file' to be set in the profile")
		}
	} else {
		// If a custom endpoint is used (e.g., LM Studio), the endpoint itself is the primary requirement.
		// We don't strictly require API key/credentials file for custom endpoints, as they might not need it.
		// However, the endpoint must not be empty if it's not the default.
		if p.Profile.Endpoint == "" { // Check original Profile.Endpoint, not the resolved one
			return fmt.Errorf("OpenAI provider (using custom endpoint) requires 'endpoint' to be specified in the profile")
		}
	}

	// If a credentials file is provided, attempt to resolve its path and check existence.
	if p.Profile.CredentialsFile != "" {
		resolvedPath, err := config.ResolvePath(p.Profile.CredentialsFile)
		if err != nil {
			return fmt.Errorf("failed to resolve credentials file path %s: %w", p.Profile.CredentialsFile, err)
		}
		// Check if the file exists and is readable.
		if _, err := os.Stat(resolvedPath); os.IsNotExist(err) {
			return fmt.Errorf("credentials file not found at %s", resolvedPath)
		}
	}

	return nil
}
