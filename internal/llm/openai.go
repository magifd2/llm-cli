package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/magifd2/llm-cli/internal/config"
)

// OpenAIProvider implements the Provider interface for OpenAI-compatible APIs.
// This includes OpenAI's own API and local LLM servers like LM Studio that mimic OpenAI's API.
type OpenAIProvider struct {
	Profile config.Profile // The configuration profile for this OpenAI-compatible instance.
}

// openAIRequest represents the JSON structure for requests to the OpenAI Chat Completions API.
type openAIRequest struct {
	Model    string    `json:"model"`          // The name of the model to use.
	Messages []message `json:"messages"`       // A list of messages in the conversation history.
	Stream   bool      `json:"stream,omitempty"` // Whether to stream the response. Omitted if false.
}

// openAIResponse represents the JSON structure for a non-streaming response from the OpenAI API.
type openAIResponse struct {
	Choices []struct {
		Message message `json:"message"` // The assistant's message.
	} `json:"choices"` // A list of chat completion choices.
}

// Chat sends a chat request to the OpenAI-compatible API and returns a single, complete response.
func (p *OpenAIProvider) Chat(systemPrompt, userPrompt string) (string, error) {
	// Determine the API endpoint. Use a default if not specified in the profile.
	endpoint := p.Profile.Endpoint
	if endpoint == "" {
		endpoint = "https://api.openai.com/v1/chat/completions"
	}

	// Build the messages array, including the system prompt if provided.
	messages := []message{}
	if systemPrompt != "" {
		messages = append(messages, message{Role: "system", Content: systemPrompt})
	}
	messages = append(messages, message{Role: "user", Content: userPrompt})

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

	// Set necessary headers, including Content-Type and Authorization (if API key is provided).
	req.Header.Set("Content-Type", "application/json")
	if p.Profile.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+p.Profile.APIKey)
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
func (p *OpenAIProvider) ChatStream(ctx context.Context, systemPrompt, userPrompt string, responseChan chan<- string) error {
	// Note: The caller is responsible for closing the responseChan.

	// Determine the API endpoint. Use a default if not specified in the profile.
	endpoint := p.Profile.Endpoint
	if endpoint == "" {
		endpoint = "https://api.openai.com/v1/chat/completions"
	}

	// Build the messages array, including the system prompt if provided.
	messages := []message{}
	if systemPrompt != "" {
		messages = append(messages, message{Role: "system", Content: systemPrompt})
	}
	messages = append(messages, message{Role: "user", Content: userPrompt})

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

	// Set necessary headers, including Content-Type and Authorization (if API key is provided).
	req.Header.Set("Content-Type", "application/json")
	if p.Profile.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+p.Profile.APIKey)
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
