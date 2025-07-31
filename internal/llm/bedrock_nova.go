package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	appconfig "github.com/magifd2/llm-cli/internal/config"
)

// NovaBedrockProvider implements the Provider interface for Amazon Bedrock.
type NovaBedrockProvider struct {
	Profile appconfig.Profile
}

// novaMessageContent defines the structure for content within a message for Nova.
type novaMessageContent struct {
	Text string `json:"text"`
}

// novaMessage defines the structure for a single message in the conversation for Nova.
type novaMessage struct {
	Role    string               `json:"role"`
	Content []novaMessageContent `json:"content"`
}

// novaSystemPrompt defines the structure for a system prompt for Nova.
type novaSystemPrompt struct {
	Text string `json:"text"`
}

// inferenceConfig defines the structure for inference parameters.
type inferenceConfig struct {
	MaxTokens     int      `json:"maxTokens,omitempty"`
	Temperature   float64  `json:"temperature,omitempty"`
	TopP          float64  `json:"topP,omitempty"`
	TopK          int      `json:"topK,omitempty"`
	StopSequences []string `json:"stopSequences,omitempty"`
}

// novaMessagesAPIRequest represents the request body for Nova models using the Messages API.
type novaMessagesAPIRequest struct {
	SchemaVersion   string           `json:"schemaVersion"`
	Messages        []novaMessage    `json:"messages"`
	System          []novaSystemPrompt `json:"system,omitempty"`
	InferenceConfig inferenceConfig  `json:"inferenceConfig,omitempty"`
}

// novaCombinedAPIResponse represents the response structure for Nova Messages API.
type novaCombinedAPIResponse struct {
	Output struct {
		Message struct {
			Content []struct {
				Text string `json:"text"`
			} `json:"content"`
			Role string `json:"role"`
		} `json:"message"`
	} `json:"output"`
	StopReason string `json:"stopReason"`
	Usage      struct {
		InputTokens            int `json:"inputTokens"`
		OutputTokens           int `json:"outputTokens"`
		TotalTokens            int `json:"totalTokens"`
		CacheReadInputTokenCount  int `json:"cacheReadInputTokenCount"`
		CacheWriteInputTokenCount int `json:"cacheWriteInputTokenCount"`
	} `json:"usage"`
}

// novaMessagesAPIStreamChunk represents a chunk of a streaming response from a Nova Messages API model.
type novaMessagesAPIStreamChunk struct {
	ContentBlockDelta struct {
		Delta struct {
			Text string `json:"text"`
		} `json:"delta"`
	} `json:"contentBlockDelta"`
}

// newBedrockClient creates a new Bedrock Runtime client.
func newBedrockClient(ctx context.Context, profile appconfig.Profile) (*bedrockruntime.Client, error) {
	var opts []func(*config.LoadOptions) error
	opts = append(opts, config.WithRegion(profile.AWSRegion))

	if profile.AWSAccessKeyID != "" && profile.AWSSecretAccessKey != "" {
		opts = append(opts, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(profile.AWSAccessKeyID, profile.AWSSecretAccessKey, "")))
	}

	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return bedrockruntime.NewFromConfig(cfg), nil
}

// Chat sends a chat request to the Amazon Bedrock API using the Messages API format.
func (p *NovaBedrockProvider) Chat(systemPromptText, userPrompt string) (string, error) {

	ctx := context.Background()
	client, err := newBedrockClient(ctx, p.Profile)
	if err != nil {
		return "", err
	}

	// Construct user message
	messages := []novaMessage{
		{
			Role: "user",
			Content: []novaMessageContent{
				{Text: userPrompt},
			},
		},
	}

	// Construct system prompt as a slice of structs, only if it's not empty.
	var systemContent []novaSystemPrompt
	if systemPromptText != "" {
		systemContent = append(systemContent, novaSystemPrompt{Text: systemPromptText})
	}

	// Initialize InferenceConfig directly in the struct literal.
	reqBody := novaMessagesAPIRequest{
		SchemaVersion: "messages-v1",
		Messages:      messages,
		System:        systemContent,
		InferenceConfig: inferenceConfig{
			MaxTokens:   500,
			Temperature: 0.7,
			TopP:        0.9,
			TopK:        20,
		},
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

	responseBodyBytes := output.Body

	// Unmarshal into the new response struct
	var novaResp novaCombinedAPIResponse
	if err := json.Unmarshal(responseBodyBytes, &novaResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	// Extract text from the new response structure
	if len(novaResp.Output.Message.Content) > 0 {
		return novaResp.Output.Message.Content[0].Text, nil
	}

	return "", fmt.Errorf("no content found in response")
}

// ChatStream sends a streaming chat request to the Amazon Bedrock API using the Messages API format.
func (p *NovaBedrockProvider) ChatStream(ctx context.Context, systemPromptText, userPrompt string, responseChan chan<- string) error {
	defer close(responseChan)

	client, err := newBedrockClient(ctx, p.Profile)
	if err != nil {
		return fmt.Errorf("error creating bedrock client: %w", err)
	}

	// Construct user message
	messages := []novaMessage{
		{
			Role: "user",
			Content: []novaMessageContent{
				{Text: userPrompt},
			},
		},
	}

	// Construct system prompt as a slice of structs, only if it's not empty.
	var systemContent []novaSystemPrompt
	if systemPromptText != "" {
		systemContent = append(systemContent, novaSystemPrompt{Text: systemPromptText})
	}

	// Initialize InferenceConfig directly in the struct literal.
	reqBody := novaMessagesAPIRequest{
		SchemaVersion: "messages-v1",
		Messages:      messages,
		System:        systemContent,
		InferenceConfig: inferenceConfig{
			MaxTokens:   500,
			Temperature: 0.7,
			TopP:        0.9,
			TopK:        20,
		},
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

	stream := output.GetStream()
	for event := range stream.Events() {
		select {
		case <-ctx.Done():
			stream.Close()
			return ctx.Err()
		default:
		}

		switch v := event.(type) {
		case *types.ResponseStreamMemberChunk:
			var chunk novaMessagesAPIStreamChunk
			if err := json.Unmarshal(v.Value.Bytes, &chunk); err != nil {
				fmt.Fprintf(os.Stderr, "Error unmarshaling stream chunk: %v\n", err)
				continue
			}
			responseChan <- chunk.ContentBlockDelta.Delta.Text
		default:
			// Handle other event types if necessary.
			return fmt.Errorf("unhandled stream event type: %T", v) // Added for debugging unhandled types
		}
	}

	return nil
}
