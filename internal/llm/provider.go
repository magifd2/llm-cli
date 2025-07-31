package llm

import "context"

// Provider defines the interface for interacting with an LLM.
type Provider interface {
	Chat(systemPrompt, userPrompt string) (string, error)
	ChatStream(ctx context.Context, systemPrompt, userPrompt string, responseChan chan<- string) error
}
