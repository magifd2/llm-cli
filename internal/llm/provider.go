package llm

import "context"

// Message represents a single message in the chat history, with a role and content.
type Message struct {
	Role    string `json:"role"`    // The role of the message sender (e.g., "user", "system", "assistant").
	Content string `json:"content"` // The content of the message.
}

// Provider defines the interface for interacting with a Large Language Model (LLM).
// It specifies methods for both single-response chat and streaming chat interactions.
type Provider interface {
	// Chat sends a single user prompt and an optional system prompt to the LLM and returns a single response.
	Chat(systemPrompt, userPrompt string) (string, error)
	// ChatStream sends a user prompt and an optional system prompt to the LLM and streams the response.
	// The context allows for cancellation of the streaming operation.
	// Response tokens are sent to the provided response channel.
	ChatStream(ctx context.Context, systemPrompt, userPrompt string, responseChan chan<- string) error
}