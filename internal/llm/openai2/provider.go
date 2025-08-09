package openai2

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

// Provider implements the llm.Provider interface for enhanced OpenAI-compatible APIs.
// It adds support for dynamically selecting models by querying the `/v1/models` endpoint,
// allowing for features like an "auto" mode.
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

// openAIStreamResponse represents a chunk of the streaming response from the OpenAI API.
type openAIStreamResponse struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"` // The content delta for the current chunk.
		} `json:"delta"` // The change in content.
	} `json:"choices"` // A list of chat completion choices (usually one for streaming).
}

// openAIAPIKey represents the structure of the OpenAI API key JSON file.
type openAIAPIKey struct {
	OpenAIAPIKey string `json:"openai_api_key"`
}

// openAIModelsResponse defines the structure for the response from the /v1/models endpoint.
type openAIModelsResponse struct {
	Data []struct {
		ID string `json:"id"`
	} `json:"data"`
}

// getAPIKey retrieves the API key from the profile or a credentials file.
func (p *Provider) getAPIKey() (string, error) {
	if p.Profile.CredentialsFile != "" {
		return loadOpenAIAPIKeyFromFile(p.Profile.CredentialsFile)
	}
	return p.Profile.APIKey, nil
}

// getAvailableModels fetches the list of available models from the /v1/models endpoint.
func (p *Provider) getAvailableModels() ([]string, error) {
	// Trim specific suffixes to get the base endpoint URL.
	baseEndpoint := strings.TrimSuffix(p.Profile.Endpoint, "/v1/chat/completions")
	baseEndpoint = strings.TrimSuffix(baseEndpoint, "/v1")
	modelsEndpoint := baseEndpoint + "/v1/models"

	req, err := http.NewRequestWithContext(context.Background(), "GET", modelsEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request for models: %w", err)
	}

	apiKey, err := p.getAPIKey()
	if err != nil {
		return nil, err
	}

	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request to /v1/models: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// It's okay if this fails (e.g., endpoint not found), so we don't return a hard error.
		// We'll just fall back to the user-specified model.
		return nil, fmt.Errorf("models endpoint returned status %d", resp.StatusCode)
	}

	var modelsResp openAIModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&modelsResp); err != nil {
		return nil, fmt.Errorf("error decoding models response: %w", err)
	}

	var modelIDs []string
	for _, model := range modelsResp.Data {
		modelIDs = append(modelIDs, model.ID)
	}

	return modelIDs, nil
}

// resolveModel determines the final model name to use based on a prioritized list from the user's profile setting.
// It supports combinations like "model1,auto,model2".
func (p *Provider) resolveModel() (string, error) {
	userModelSetting := p.Profile.Model
	priorityList := strings.Split(userModelSetting, ",")

	// First, try to get the list of available models from the endpoint.
	availableModels, err := p.getAvailableModels()

	// Create a map for quick lookup of available models.
	availableModelsMap := make(map[string]bool)
	if err == nil {
		for _, m := range availableModels {
			availableModelsMap[m] = true
		}
	}

	// Iterate through the user's prioritized list.
	for _, candidate := range priorityList {
		trimmedCandidate := strings.TrimSpace(candidate)
		if trimmedCandidate == "" {
			continue
		}

		if trimmedCandidate == "auto" {
			// If 'auto' is specified, try to use the first model from the available list.
			if err == nil && len(availableModels) > 0 {
				return availableModels[0], nil
			}
		} else {
			// If a specific model is specified, check if it's in the available list.
			if err == nil {
				if _, ok := availableModelsMap[trimmedCandidate]; ok {
					return trimmedCandidate, nil
				}
			}
		}
	}

	// If no model could be resolved from the priority list after checking the endpoint:
	// As a final fallback, if the model list couldn't be fetched (e.g., endpoint doesn't support /v1/models),
	// let's try the first non-"auto" model from the user's list.
	if err != nil {
		for _, candidate := range priorityList {
			trimmedCandidate := strings.TrimSpace(candidate)
			if trimmedCandidate != "auto" && trimmedCandidate != "" {
				return trimmedCandidate, nil
			}
		}
	}

	return "", fmt.Errorf("could not resolve a valid model from the priority list: [%s]", userModelSetting)
}

// Chat sends a chat request to the OpenAI-compatible API and returns a single, complete response.
func (p *Provider) Chat(systemPrompt, userPrompt string) (string, error) {
	model, err := p.resolveModel()
	if err != nil {
		return "", err
	}

	endpoint := p.Profile.Endpoint
	if endpoint == "" {
		endpoint = "https://api.openai.com/v1/chat/completions"
	}

	messages := []llm.Message{}
	if systemPrompt != "" {
		messages = append(messages, llm.Message{Role: "system", Content: systemPrompt})
	}
	messages = append(messages, llm.Message{Role: "user", Content: userPrompt})

	reqBody := openAIRequest{
		Model:    model,
		Messages: messages,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshalling request body: %w", err)
	}

	req, err := http.NewRequestWithContext(context.Background(), "POST", endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	apiKey, err := p.getAPIKey()
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request to openai-compatible api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("openai-compatible api request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var openAIResp openAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return "", fmt.Errorf("error decoding openai response: %w", err)
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from openai-compatible api")
	}

	return openAIResp.Choices[0].Message.Content, nil
}

// ChatStream sends a streaming chat request to the OpenAI-compatible API and sends response chunks to a channel.
func (p *Provider) ChatStream(ctx context.Context, systemPrompt, userPrompt string, responseChan chan<- string) error {
	model, err := p.resolveModel()
	if err != nil {
		return err
	}

	endpoint := p.Profile.Endpoint
	if endpoint == "" {
		endpoint = "https://api.openai.com/v1/chat/completions"
	}

	messages := []llm.Message{}
	if systemPrompt != "" {
		messages = append(messages, llm.Message{Role: "system", Content: systemPrompt})
	}
	messages = append(messages, llm.Message{Role: "user", Content: userPrompt})

	reqBody := openAIRequest{
		Model:    model,
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

	apiKey, err := p.getAPIKey()
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request to openai-compatible api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("openai-compatible api request failed with status %d: %s", resp.StatusCode, string(body))
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") || line == "data: [DONE]" {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")

		if strings.Contains(data, `"error":`) {
			return fmt.Errorf("streaming error: %s", data)
		}

		var streamResp openAIStreamResponse
		if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
			continue
		}

		if len(streamResp.Choices) > 0 {
			select {
			case responseChan <- streamResp.Choices[0].Delta.Content:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading stream: %w", err)
	}

	return nil
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

// NewProvider is a factory function that returns a new OpenAI2 provider.
func NewProvider(p config.Profile) (llm.Provider, error) {
	return &Provider{Profile: p}, nil
}
