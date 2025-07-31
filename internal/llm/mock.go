package llm

import (
	"context"
	"fmt"
)

// MockProvider is a dummy provider for testing.
type MockProvider struct{}

// Chat implements the Provider interface for MockProvider.
func (p *MockProvider) Chat(systemPrompt, userPrompt string) (string, error) {
	response := fmt.Sprintf("\n--- Mock Response ---\nSystem Prompt: %s\nUser Prompt: %s\n---------------------\n", systemPrompt, userPrompt)
	return response, nil
}

// ChatStream implements the Provider interface for MockProvider.
func (p *MockProvider) ChatStream(ctx context.Context, systemPrompt, userPrompt string, responseChan chan<- string) error {
	defer close(responseChan)
	response, _ := p.Chat(systemPrompt, userPrompt)
	select {
	case responseChan <- response:
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

