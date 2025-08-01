package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/magifd2/llm-cli/internal/config"
)

// OllamaProvider implements the Provider interface for interacting with Ollama LLMs.
type OllamaProvider struct {
	Profile config.Profile // The configuration profile for this Ollama instance.
}

// ollamaRequest represents the JSON structure for requests to the Ollama chat API.
type ollamaRequest struct {
	Model    string    `json:"model"`    // The name of the model to use.
	Messages []message `json:"messages"` // A list of messages in the conversation history.
	Stream   bool      `json:"stream"`   // Whether to stream the response.
}

// message represents a single message in the chat history, with a role and content.
type message struct {
	Role    string `json:"role"`    // The role of the message sender (e.g., "user", "system", "assistant").
	Content string `json:"content"` // The content of the message.
}

// ollamaResponse represents the JSON structure for responses from the Ollama chat API.
// For streaming, each chunk will contain a message.
type ollamaResponse struct {
	Message message `json:"message"` // The message content from the LLM.
}

// Chat sends a chat request to the Ollama API and returns a single, complete response.
func (p *OllamaProvider) Chat(systemPrompt, userPrompt string) (string, error) {
	// Determine the API endpoint. Use a default if not specified in the profile.
	endpoint := p.Profile.Endpoint
	if endpoint == "" {
		endpoint = "http://localhost:11434/api/chat"
	}

	// Build the messages array, including the system prompt if provided.
	messages := []message{}
	if systemPrompt != "" {
		messages = append(messages, message{Role: "system", Content: systemPrompt})
	}
	messages = append(messages, message{Role: "user", Content: userPrompt})

	// Construct the request body for a non-streaming chat.
	reqBody := ollamaRequest{
		Model:    p.Profile.Model,
		Messages: messages,
		Stream:   false, // Explicitly set to false for non-streaming chat.
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshalling request body: %w", err)
	}

	// Send the HTTP POST request to the Ollama API.
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("error making request to ollama: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-OK HTTP status codes and return an error with the response body.
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Decode the JSON response.
	var ollamaResp ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("error decoding ollama response: %w", err)
	}

	return ollamaResp.Message.Content, nil
}

// ChatStream sends a streaming chat request to the Ollama API and sends response chunks to a channel.
func (p *OllamaProvider) ChatStream(ctx context.Context, systemPrompt, userPrompt string, responseChan chan<- string) error {
	defer close(responseChan)

	// Determine the API endpoint. Use a default if not specified in the profile.
	endpoint := p.Profile.Endpoint
	if endpoint == "" {
		endpoint = "http://localhost:11434/api/chat"
	}

	// Build the messages array, including the system prompt if provided.
	messages := []message{}
	if systemPrompt != "" {
		messages = append(messages, message{Role: "system", Content: systemPrompt})
	}
	messages = append(messages, message{Role: "user", Content: userPrompt})

	// Construct the request body for a streaming chat.
	reqBody := ollamaRequest{
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
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request to ollama: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-OK HTTP status codes.
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ollama API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Read and process the streaming response line by line.
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var streamResp ollamaResponse
		// Unmarshal each line as a JSON response chunk.
		if err := json.Unmarshal([]byte(line), &streamResp); err != nil {
			return fmt.Errorf("error decoding ollama stream response: %w", err)
		}

		// Send the message content to the response channel or stop if context is cancelled.
		select {
		case responseChan <- streamResp.Message.Content:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	// Check for any errors during scanning.
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading stream: %w", err)
	}

	return nil
}
