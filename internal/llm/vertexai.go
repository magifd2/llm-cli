package llm

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cloud.google.com/go/vertexai/genai"
	"github.com/magifd2/llm-cli/internal/config"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// VertexAIProvider implements the Provider interface for Google Cloud Vertex AI.
// 認証は、プロファイルで `credentials_file` が指定されている場合はそのサービスアカウントキーを、
// 指定されていない場合は Application Default Credentials (ADC) を使用します。
type VertexAIProvider struct {
	Profile config.Profile
}

// newVertexAIClient は、プロファイル情報に基づいてVertex AIクライアントを初期化します。
func (p *VertexAIProvider) newVertexAIClient(ctx context.Context) (*genai.Client, error) {
	projectID := p.Profile.ProjectID
	location := p.Profile.Location
	credentialsFile := p.Profile.CredentialsFile // サービスアカウントキーファイルのパス

	if projectID == "" || location == "" {
		return nil, fmt.Errorf("project_id and location must be set in the profile for vertexai provider")
	}

	// credentials_fileが指定されていれば、それを使用してクライアントを初期化
	if credentialsFile != "" {
		expandedPath, err := expandPath(credentialsFile)
		if err != nil {
			return nil, fmt.Errorf("expanding path for credentials_file: %w", err)
		}
		return genai.NewClient(ctx, projectID, location, option.WithCredentialsFile(expandedPath))
	}

	// 指定されていない場合は、ADC (Application Default Credentials) を使用
	return genai.NewClient(ctx, projectID, location)
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

// Chat sends a chat request to the Vertex AI API.
func (p *VertexAIProvider) Chat(systemPrompt, userPrompt string) (string, error) {
	ctx := context.Background()
	client, err := p.newVertexAIClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	model := client.GenerativeModel(p.Profile.Model)
	// Vertex AI (Gemini) では、SystemInstructionsはチャットセッションの最初のターンにのみ有効です。
	// ここでは、リクエストごとに設定することで、対話的でない `Chat` メソッドの動作を模倣します。
	if systemPrompt != "" {
		model.SystemInstruction = &genai.Content{
			Parts: []genai.Part{genai.Text(systemPrompt)},
		}
	}

	resp, err := model.GenerateContent(ctx, genai.Text(userPrompt))
	if err != nil {
		return "", fmt.Errorf("error generating content from vertexai: %w", err)
	}

	return extractTextFromResponse(resp), nil
}

// ChatStream sends a streaming chat request to the Vertex AI API.
func (p *VertexAIProvider) ChatStream(ctx context.Context, systemPrompt, userPrompt string, responseChan chan<- string) error {
	client, err := p.newVertexAIClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	model := client.GenerativeModel(p.Profile.Model)
	if systemPrompt != "" {
		model.SystemInstruction = &genai.Content{
			Parts: []genai.Part{genai.Text(systemPrompt)},
		}
	}

	iter := model.GenerateContentStream(ctx, genai.Text(userPrompt))

	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
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

// extractTextFromResponse は、Vertex AIのレスポンスからテキスト部分を抽出して結合します。
func extractTextFromResponse(resp *genai.GenerateContentResponse) string {
	var sb strings.Builder
	if resp == nil || len(resp.Candidates) == 0 {
		return ""
	}

	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				if txt, ok := part.(genai.Text); ok {
					sb.WriteString(string(txt))
				}
			}
		}
	}
	return sb.String()
}