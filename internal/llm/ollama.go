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

// OllamaProvider implements the Provider interface for Ollama.
type OllamaProvider struct {
	Profile config.Profile
}

// ollamaRequest represents the request body for the Ollama API.
type ollamaRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
	Stream   bool      `json:"stream"`
}

// message represents a single message in the chat history.
type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ollamaResponse represents the response body from the Ollama API.
type ollamaResponse struct {
	Message message `json:"message"`
}

// Chat sends a chat request to the Ollama API.
func (p *OllamaProvider) Chat(systemPrompt, userPrompt string) (string, error) {
	endpoint := p.Profile.Endpoint
	if endpoint == "" {
		// Use a default endpoint if not specified in the profile
		endpoint = "http://localhost:11434/api/chat"
	}

	messages := []message{}
	if systemPrompt != "" {
		messages = append(messages, message{Role: "system", Content: systemPrompt})
	}
	messages = append(messages, message{Role: "user", Content: userPrompt})

	reqBody := ollamaRequest{
		Model:    p.Profile.Model,
		Messages: messages,
		Stream:   false, // For now, we don't support streaming
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshalling request body: %w", err)
	}

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("error making request to ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var ollamaResp ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("error decoding ollama response: %w", err)
	}

	return ollamaResp.Message.Content, nil
}

// ChatStream sends a streaming chat request to the Ollama API.
func (p *OllamaProvider) ChatStream(ctx context.Context, systemPrompt, userPrompt string, responseChan chan<- string) error {
	defer close(responseChan)

	endpoint := p.Profile.Endpoint
	if endpoint == "" {
		endpoint = "http://localhost:11434/api/chat"
	}

	messages := []message{}
	if systemPrompt != "" {
		messages = append(messages, message{Role: "system", Content: systemPrompt})
	}
	messages = append(messages, message{Role: "user", Content: userPrompt})

	reqBody := ollamaRequest{
		Model:    p.Profile.Model,
		Messages: messages,
		Stream:   true,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("error marshalling request body: %w", err)
	}

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

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ollama API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var streamResp ollamaResponse
		if err := json.Unmarshal([]byte(line), &streamResp); err != nil {
			return fmt.Errorf("error decoding ollama stream response: %w", err)
		}

		select {
		case responseChan <- streamResp.Message.Content:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading stream: %w", err)
	}

	return nil
}
