package llm

import (
	"context"
	"fmt"

	"cloud.google.com/go/vertexai/genai"
	"github.com/magifd2/llm-cli/internal/config"
)

// VertexAIProvider implements the Provider interface for Google Cloud Vertex AI.
type VertexAIProvider struct {
	Profile config.Profile
}

// newVertexAIClient creates a new Vertex AI client.
func newVertexAIClient(ctx context.Context, projectID, location string) (*genai.Client, error) {
	return genai.NewClient(ctx, projectID, location)
}

// Chat sends a chat request to the Vertex AI Gemini API.
func (p *VertexAIProvider) Chat(systemPrompt, userPrompt string) (string, error) {
	ctx := context.Background()
	client, err := newVertexAIClient(ctx, p.Profile.ProjectID, p.Profile.Location)
	if err != nil {
		return "", fmt.Errorf("error creating vertex ai client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel(p.Profile.Model)
	cs := model.StartChat()

	var parts []genai.Part
	if systemPrompt != "" {
		model.SystemInstruction = &genai.Content{Parts: []genai.Part{genai.Text(systemPrompt)}}
	}
	parts = append(parts, genai.Text(userPrompt))

	resp, err := cs.SendMessage(ctx, parts...)
	if err != nil {
		return "", fmt.Errorf("error sending message to vertex ai: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from vertex ai")
	}

	return fmt.Sprint(resp.Candidates[0].Content.Parts[0]), nil
}

// ChatStream sends a streaming chat request to the Vertex AI Gemini API.
func (p *VertexAIProvider) ChatStream(ctx context.Context, systemPrompt, userPrompt string, responseChan chan<- string) error {
	defer close(responseChan)

	client, err := newVertexAIClient(ctx, p.Profile.ProjectID, p.Profile.Location)
	if err != nil {
		return fmt.Errorf("error creating vertex ai client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel(p.Profile.Model)
	cs := model.StartChat()

	var parts []genai.Part
	if systemPrompt != "" {
		model.SystemInstruction = &genai.Content{Parts: []genai.Part{genai.Text(systemPrompt)}}
	}
	parts = append(parts, genai.Text(userPrompt))

	stream := cs.SendMessageStream(ctx, parts...)

	for {
		resp, err := stream.Next()
		if err != nil {
			if err.Error() == "iterator done" { // Correctly check for end of stream
				break
			}
			return fmt.Errorf("error receiving stream from vertex ai: %w", err)
		}

		if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
			token := fmt.Sprint(resp.Candidates[0].Content.Parts[0])
			select {
			case responseChan <- token:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	return nil
}
