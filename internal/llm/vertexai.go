package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/magifd2/llm-cli/internal/config"
	"cloud.google.com/go/auth"
	"google.golang.org/genai"
)

// VertexAIProvider implements the Provider interface for Google Cloud Vertex AI.
// This provider uses the new `google.golang.org/genai` SDK and specifies the Vertex AI backend
// to resolve deprecation warnings and ensure compatibility.
type VertexAIProvider struct {
	Profile config.Profile // The configuration profile for this Vertex AI instance.
}

// newVertexAIClient initializes a Vertex AI client based on the provided profile information.
// It handles project ID, location, and optional service account key authentication.
func (p *VertexAIProvider) newVertexAIClient(ctx context.Context) (*genai.Client, error) {
	projectID := p.Profile.ProjectID
	location := p.Profile.Location
	credentialsFile := p.Profile.CredentialsFile

	// Ensure project_id and location are set in the profile.
	if projectID == "" || location == "" {
		return nil, fmt.Errorf("project_id and location must be set in the profile for vertexai provider")
	}

	// Configure the client with project, location, and backend.
	clientConfig := &genai.ClientConfig{
		Project:  projectID,
		Location: location,
		Backend:  genai.BackendVertexAI,
	}

	// If a credentials file is provided, load and use it for authentication.
	if credentialsFile != "" {
		// Expand the path to the credentials file (e.g., handling ~ for home directory).
		expandedPath, err := expandPath(credentialsFile)
		if err != nil {
			return nil, fmt.Errorf("expanding path for credentials_file: %w", err)
		}
		// Read the credentials file content.
		credsJSON, err := os.ReadFile(expandedPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read credentials file: %w", err)
		}
		// Check if the credentials JSON is empty.
		if len(credsJSON) == 0 {
			return nil, fmt.Errorf("credentials file is empty: %s", expandedPath)
		}

		// Parse the service account JSON to extract necessary information.
		var sa struct {
			ClientEmail string `json:"client_email"`
			PrivateKey  string `json:"private_key"`
			TokenURI    string `json:"token_uri"`
			ProjectID   string `json:"project_id"`
		}
		if err := json.Unmarshal(credsJSON, &sa); err != nil {
			return nil, fmt.Errorf("invalid service account JSON: %w", err)
		}

		// Create a TokenProvider using the service account credentials.
		tp, err := auth.New2LOTokenProvider(&auth.Options2LO{
			Email:      sa.ClientEmail,
			PrivateKey: []byte(sa.PrivateKey),
			TokenURL:   sa.TokenURI,
			Scopes:     []string{"https://www.googleapis.com/auth/cloud-platform"}, // Cloud Platform scope is required.
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create token provider: %w", err)
		}

		// Set the credentials in the client configuration.
		clientConfig.Credentials = auth.NewCredentials(&auth.CredentialsOptions{TokenProvider: tp})
	}

	// Create and return the new genai client.
	client, err := genai.NewClient(ctx, clientConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating new genai client: %w", err)
	}

	return client, nil
}

// expandPath expands a path that starts with ~ to the user's home directory.
func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(homeDir, path[1:]), nil
	}
	return path, nil
}

// Chat sends a chat request to the Vertex AI API and returns a single, complete response.
// System prompts are sent as the first message in the conversation.
func (p *VertexAIProvider) Chat(systemPrompt, userPrompt string) (string, error) {
	ctx := context.Background()
	client, err := p.newVertexAIClient(ctx)
	if err != nil {
		return "", err
	}

	// Create a new chat session.
	chat, err := client.Chats.Create(ctx, p.Profile.Model, nil, nil)
	if err != nil {
		return "", fmt.Errorf("error creating chat: %w", err)
	}

	// If a system prompt is provided, send it as the first message.
	if systemPrompt != "" {
		_, err := chat.SendMessage(ctx, genai.Part{Text: systemPrompt})
		if err != nil {
			return "", fmt.Errorf("error sending system prompt to vertexai: %w", err)
		}
	}

	// Send the user prompt and get the response.
	resp, err := chat.SendMessage(ctx, genai.Part{Text: userPrompt})
	if err != nil {
		return "", fmt.Errorf("error sending message to vertexai: %w", err)
	}

	return extractTextFromResponse(resp), nil
}

// ChatStream sends a streaming chat request to the Vertex AI API and streams response chunks to a channel.
// System prompts are sent as the first message in the conversation.
func (p *VertexAIProvider) ChatStream(ctx context.Context, systemPrompt, userPrompt string, responseChan chan<- string) error {
	client, err := p.newVertexAIClient(ctx)
	if err != nil {
		return err
	}

	// Create a new chat session.
	chat, err := client.Chats.Create(ctx, p.Profile.Model, nil, nil)
	if err != nil {
		return fmt.Errorf("error creating chat: %w", err)
	}

	// If a system prompt is provided, send it as the first message.
	if systemPrompt != "" {
		_, err := chat.SendMessage(ctx, genai.Part{Text: systemPrompt})
		if err != nil {
			return fmt.Errorf("error sending system prompt to vertexai: %w", err)
		}
	}

	// Stream the user prompt and process each response chunk.
	for resp, err := range chat.SendMessageStream(ctx, genai.Part{Text: userPrompt}) {
		if err != nil {
			return fmt.Errorf("error reading stream from vertexai: %w", err)
		}

		chunk := extractTextFromResponse(resp)
		if chunk != "" {
			select {
			case responseChan <- chunk:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	return nil
}

// extractTextFromResponse extracts and concatenates text content from a Vertex AI GenerateContentResponse.
func extractTextFromResponse(resp *genai.GenerateContentResponse) string {
	var sb strings.Builder
	if resp == nil || len(resp.Candidates) == 0 {
		return ""
	}

	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				sb.WriteString(part.Text)
			}
		}
	}
	return sb.String()
}
