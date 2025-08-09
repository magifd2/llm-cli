package mock

import (
	"context"
	"fmt"

	"github.com/magifd2/llm-cli/internal/config"
	"github.com/magifd2/llm-cli/internal/llm"
)

// Provider is a dummy implementation of the LLM Provider interface.
// It's used for testing and as a fallback when a real provider is not configured or recognized.
type Provider struct{}

// NewProvider is a factory function that returns a new mock provider.
func NewProvider(p config.Profile) (llm.Provider, error) {
	return &Provider{}, nil
}

// Chat provides a mock response for a single chat interaction.
// It returns a formatted string containing the system and user prompts.
func (p *Provider) Chat(systemPrompt, userPrompt string) (string, error) {
	response := fmt.Sprintf("\n--- Mock Response ---\nSystem Prompt: %s\nUser Prompt: %s\n---------------------\n", systemPrompt, userPrompt)
	return response, nil
}

// ChatStream provides a mock streaming response.
// It sends the full mock response as a single chunk to the response channel.
// The context is checked for cancellation.
func (p *Provider) ChatStream(ctx context.Context, systemPrompt, userPrompt string, responseChan chan<- string) error {
	defer close(responseChan)
	response, _ := p.Chat(systemPrompt, userPrompt)
	select {
	case responseChan <- response:
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}