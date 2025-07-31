package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"github.com/magifd2/llm-cli/internal/config"
)

// BedrockProvider implements the Provider interface for Amazon Bedrock.
type BedrockProvider struct {
	Profile config.Profile
}

// claudeRequest represents the request body for Anthropic Claude models.
type claudeRequest struct {
	Prompt            string `json:"prompt"`
	MaxTokensToSample int    `json:"max_tokens_to_sample"`
	Temperature       float64 `json:"temperature,omitempty"`
	TopP              float64 `json:"top_p,omitempty"`
	TopK              int    `json:"top_k,omitempty"`
	StopSequences     []string `json:"stop_sequences,omitempty"`
}

// claudeResponse represents the response body from Anthropic Claude models.
type claudeResponse struct {
	Completion string `json:"completion"`
}

// claudeStreamResponseChunk represents a chunk from the streaming response.
type claudeStreamResponseChunk struct {
	Bytes string `json:"bytes"`
}

// newBedrockClient creates a new Bedrock Runtime client.
func newBedrockClient(ctx context.Context, profile config.Profile) (*bedrockruntime.Client, error) {
	var opts []func(*config.LoadOptions) error
	opts = append(opts, config.WithRegion(profile.AWSRegion))

	if profile.AWSAccessKeyID != "" && profile.AWSSecretAccessKey != "" {
		opts = append(opts, config.WithCredentialsProvider(aws.NewStaticCredentialsProvider(profile.AWSAccessKeyID, profile.AWSSecretAccessKey, "")))
	}

	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return bedrockruntime.NewFromConfig(cfg), nil
}

// Chat sends a chat request to the Amazon Bedrock API.
func (p *BedrockProvider) Chat(systemPrompt, userPrompt string) (string, error) {
	ctx := context.Background()
	client, err := newBedrockClient(ctx, p.Profile)
	if err != nil {
		return "", err
	}

	// Claude requires a specific prompt format.
	fullPrompt := fmt.Sprintf("\n\nHuman: %s\n\nAssistant:", userPrompt)
	if systemPrompt != "" {
		fullPrompt = fmt.Sprintf("\n\nHuman: %s\n%s\n\nAssistant:", systemPrompt, userPrompt)
	}

	reqBody := claudeRequest{
		Prompt:            fullPrompt,
		MaxTokensToSample: 2000,
		Temperature:       0.7,
		TopP:              1.0,
		TopK:              250,
		StopSequences:     []string{"\n\nHuman:"},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	output, err := client.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(p.Profile.Model),
		ContentType: aws.String("application/json"),
		Body:        jsonBody,
	})
	if err != nil {
		return "", fmt.Errorf("failed to invoke model: %w", err)
	}

	var claudeResp claudeResponse
	if err := json.Unmarshal(output.Body, &claudeResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return claudeResp.Completion, nil
}

// ChatStream sends a streaming chat request to the Amazon Bedrock API.
func (p *BedrockProvider) ChatStream(ctx context.Context, systemPrompt, userPrompt string, responseChan chan<- string) error {
	defer close(responseChan)

	client, err := newBedrockClient(ctx, p.Profile)
	if err != nil {
		return fmt.Errorf("error creating bedrock client: %w", err)
	}

	// Claude requires a specific prompt format.
	fullPrompt := fmt.Sprintf("\n\nHuman: %s\n\nAssistant:", userPrompt)
	if systemPrompt != "" {
		fullPrompt = fmt.Sprintf("\n\nHuman: %s\n%s\n\nAssistant:", systemPrompt, userPrompt)
	}

	reqBody := claudeRequest{
		Prompt:            fullPrompt,
		MaxTokensToSample: 2000,
		Temperature:       0.7,
		TopP:              1.0,
		TopK:              250,
		StopSequences:     []string{"\n\nHuman:"},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	output, err := client.InvokeModelWithResponseStream(ctx, &bedrockruntime.InvokeModelWithResponseStreamInput{
		ModelId:     aws.String(p.Profile.Model),
		ContentType: aws.String("application/json"),
		Body:        jsonBody,
	})
	if err != nil {
		return fmt.Errorf("failed to invoke model with stream: %w", err)
	}

	for event := range output.GetStream().Events() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		switch v := event.(type) {
		case *types.ResponseStreamMemberChunk:
			var chunk claudeStreamResponseChunk
			if err := json.Unmarshal(v.Bytes, &chunk); err != nil {
				return fmt.Errorf("failed to unmarshal stream chunk: %w", err)
			}
			responseChan <- chunk.Bytes // Claude's streaming response is just the completion text in 'bytes'
		case *types.ResponseStreamMemberError:
			return fmt.Errorf("bedrock stream error: %s - %s", *v.Error.ErrorCode, *v.Error.Message)
		default:
			// Ignore other event types
		}
	}

	return nil
}
