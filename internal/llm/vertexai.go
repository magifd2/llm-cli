package llm

import (
	"context"
	"encoding/json" // 追加: JSONパース用
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/magifd2/llm-cli/internal/config"
	"cloud.google.com/go/auth" // auth.Credentials を使用するため追加
	"google.golang.org/genai"
)

// VertexAIProvider implements the Provider interface for Google Cloud Vertex AI.
// NOTE: このプロバイダーは、新しい `google.golang.org/genai` SDK を使用し、
// Vertex AI バックエンドを指定することで、非推奨の警告を解決します。
type VertexAIProvider struct {
	Profile config.Profile
}

// newVertexAIClient は、プロファイル情報に基づいてVertex AIクライアントを初期化します。
func (p *VertexAIProvider) newVertexAIClient(ctx context.Context) (*genai.Client, error) {
	projectID := p.Profile.ProjectID
	location := p.Profile.Location
	credentialsFile := p.Profile.CredentialsFile

	if projectID == "" || location == "" {
		return nil, fmt.Errorf("project_id and location must be set in the profile for vertexai provider")
	}

	clientConfig := &genai.ClientConfig{
		Project:  projectID,
		Location: location,
		Backend:  genai.BackendVertexAI,
	}

	if credentialsFile != "" {
		expandedPath, err := expandPath(credentialsFile)
		if err != nil {
			return nil, fmt.Errorf("expanding path for credentials_file: %w", err)
		}
		// クレデンシャルファイルから認証情報をロード
		credsJSON, err := os.ReadFile(expandedPath) // ファイルを読み込む
		if err != nil {
			return nil, fmt.Errorf("failed to read credentials file: %w", err)
		}
		if len(credsJSON) == 0 { // credsJSON が空の場合のチェックを追加
			return nil, fmt.Errorf("credentials file is empty: %s", expandedPath)
		}

		// サービスアカウントJSONをパースして必要な情報を抽出
		var sa struct {
			ClientEmail string `json:"client_email"`
			PrivateKey  string `json:"private_key"`
			TokenURI    string `json:"token_uri"`
			ProjectID   string `json:"project_id"` // JSONからproject_idを取得
		}
		if err := json.Unmarshal(credsJSON, &sa); err != nil {
			return nil, fmt.Errorf("invalid service account JSON: %w", err)
		}

		// TokenProvider を作成
		tp, err := auth.New2LOTokenProvider(&auth.Options2LO{
			Email:      sa.ClientEmail,
			PrivateKey: []byte(sa.PrivateKey),
			TokenURL:   sa.TokenURI,
			Scopes:     []string{"https://www.googleapis.com/auth/cloud-platform"}, // Scope は必須
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create token provider: %w", err)
		}

		clientConfig.Credentials = auth.NewCredentials(&auth.CredentialsOptions{TokenProvider: tp}) // auth.NewCredentials を使用
	}

	// Vertex AI用のクライアントを作成します。
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

// Chat sends a chat request to the Vertex AI API.
func (p *VertexAIProvider) Chat(systemPrompt, userPrompt string) (string, error) {
	ctx := context.Background()
	client, err := p.newVertexAIClient(ctx)
	if err != nil {
		return "", err
	}

	// Chatオブジェクトを生成
	chat, err := client.Chats.Create(ctx, p.Profile.Model, nil, nil) // モデルIDを直接渡す
	if err != nil {
		return "", fmt.Errorf("error creating chat: %w", err)
	}

	// システムプロンプトがある場合、最初に送信
	if systemPrompt != "" {
		_, err := chat.SendMessage(ctx, genai.Part{Text: systemPrompt})
		if err != nil {
			return "", fmt.Errorf("error sending system prompt to vertexai: %w", err)
		}
	}

	resp, err := chat.SendMessage(ctx, genai.Part{Text: userPrompt})
	if err != nil {
		return "", fmt.Errorf("error sending message to vertexai: %w", err)
	}

	return extractTextFromResponse(resp), nil
}

// ChatStream sends a streaming chat request to the Vertex AI API.
func (p *VertexAIProvider) ChatStream(ctx context.Context, systemPrompt, userPrompt string, responseChan chan<- string) error {
	client, err := p.newVertexAIClient(ctx)
	if err != nil {
		return err
	}

	// Chatオブジェクトを生成
	chat, err := client.Chats.Create(ctx, p.Profile.Model, nil, nil) // モデルIDを直接渡す
	if err != nil {
		return fmt.Errorf("error creating chat: %w", err)
	}

	// システムプロンプトがある場合、最初に送信
	if systemPrompt != "" {
		_, err := chat.SendMessage(ctx, genai.Part{Text: systemPrompt})
		if err != nil {
			return fmt.Errorf("error sending system prompt to vertexai: %w", err)
		}
	}

	for resp, err := range chat.SendMessageStream(ctx, genai.Part{Text: userPrompt}) { // for ... range 構文を使用
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
				sb.WriteString(part.Text)
			}
		}
	}
	return sb.String()
}
