package bedrock

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

// NovaProvider implements the llm.Provider interface for Amazon Bedrock's Anthropic Claude 3 (Nova) models.
// It handles authentication and communication with the Bedrock Runtime service.
type NovaProvider struct {
	Profile appconfig.Profile // The configuration profile for this Bedrock instance.
}

// novaMessageContent defines the structure for content within a message for Nova models.
// Currently, it supports text content.
type novaMessageContent struct {
	Text string `json:"text"` // The text content of the message.
}

// novaMessage defines the structure for a single message in the conversation for Nova models.
// It includes the role of the sender and an array of content blocks.
type novaMessage struct {
	Role    string               `json:"role"`    // The role of the message sender (e.g., "user", "assistant").
	Content []novaMessageContent `json:"content"` // An array of content blocks, typically containing text.
}

// novaSystemPrompt defines the structure for a system prompt for Nova models.
// System prompts are used to set the behavior or context for the model.
type novaSystemPrompt struct {
	Text string `json:"text"` // The text content of the system prompt.
}

// inferenceConfig defines the structure for inference parameters for Nova models.
// These parameters control the model's generation behavior.
type inferenceConfig struct {
	MaxTokens     int      `json:"maxTokens,omitempty"`     // The maximum number of tokens to generate in the response.
	Temperature   float64  `json:"temperature,omitempty"`   // Controls the randomness of the output. Higher values mean more random.
	TopP          float64  `json:"topP,omitempty"`          // Controls diversity via nucleus sampling.
	TopK          int      `json:"topK,omitempty"`          // Controls diversity by limiting the number of highest probability tokens.
	StopSequences []string `json:"stopSequences,omitempty"` // A list of sequences that will cause the model to stop generating.
}

// novaMessagesAPIRequest represents the request body for Nova models using the Messages API.
// This is the primary request format for Claude 3 models on Bedrock.
type novaMessagesAPIRequest struct {
	SchemaVersion   string           `json:"schemaVersion"`     // The schema version for the API request (e.g., "messages-v1").
	Messages        []novaMessage    `json:"messages"`          // The conversation history, including user and assistant messages.
	System          []novaSystemPrompt `json:"system,omitempty"`    // Optional system prompts to guide the model's behavior.
	InferenceConfig inferenceConfig  `json:"inferenceConfig,omitempty"` // Optional inference parameters.
}

// novaCombinedAPIResponse represents the full response structure for Nova Messages API.
// It includes the generated message, stop reason, and token usage information.
type novaCombinedAPIResponse struct {
	Output struct {
		Message struct {
			Content []struct {
				Text string `json:"text"` // The text content of the assistant's response.
			} `json:"content"` // Content blocks of the message.
			Role string `json:"role"` // The role of the message sender (e.g., "assistant").
		} `json:"message"` // The generated message from the model.
	} `json:"output"` // The main output block.
	StopReason string `json:"stopReason"` // The reason the model stopped generating (e.g., "end_turn", "max_tokens").
	Usage      struct {
		InputTokens            int `json:"inputTokens"`            // Number of input tokens.
		OutputTokens           int `json:"outputTokens"`           // Number of output tokens.
		TotalTokens            int `json:"totalTokens"`            // Total number of tokens.
		CacheReadInputTokenCount  int `json:"cacheReadInputTokenCount"`  // Number of input tokens read from cache.
		CacheWriteInputTokenCount int `json:"cacheWriteInputTokenCount"` // Number of input tokens written to cache.
	} `json:"usage"` // Token usage statistics.
}

// novaMessagesAPIStreamChunk represents a single chunk of a streaming response from a Nova Messages API model.
// It typically contains a delta of content.
type novaMessagesAPIStreamChunk struct {
	ContentBlockDelta struct {
		Delta struct {
			Text string `json:"text"` // The incremental text content.
		} `json:"delta"` // The delta of content.
	} `json:"contentBlockDelta"` // The content block delta event.
}

// bedrockErrorResponse defines the structure for an error response from Bedrock.
// This is used to parse and report errors returned by the Bedrock API.
// e.g., {"message": "...", "type": "ValidationException"}
type bedrockErrorResponse struct {
	Message string `json:"message"` // The error message.
	Type    string `json:"type"`    // The type of error (e.g., "ValidationException").
}

// newBedrockClient creates a new Bedrock Runtime client.
// It configures the client with the specified AWS region and optional static credentials.
func newBedrockClient(ctx context.Context, profile appconfig.Profile) (*bedrockruntime.Client, error) {
	var opts []func(*config.LoadOptions) error
	// Set the AWS region from the profile.
	opts = append(opts, config.WithRegion(profile.AWSRegion))

	// If AWS credentials file is provided, load credentials from it.
	if profile.CredentialsFile != "" { // Changed from profile.AWSCredentialsFile
		creds, err := loadAWSCredentialsFromFile(profile.CredentialsFile) // Changed from profile.AWSCredentialsFile
		if err != nil {
			return nil, fmt.Errorf("failed to load AWS credentials from file %s: %w", profile.CredentialsFile, err) // Changed from profile.AWSCredentialsFile
		}
		opts = append(opts, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(creds.AWSAccessKeyID, creds.AWSSecretAccessKey, "")))
	} else if profile.AWSAccessKeyID != "" && profile.AWSSecretAccessKey != "" {
		// If AWS access key ID and secret access key are provided directly, use static credentials.
		opts = append(opts, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(profile.AWSAccessKeyID, profile.AWSSecretAccessKey, "")))
	}

	// Load the default AWS configuration with the specified options.
	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create and return a new Bedrock Runtime client from the loaded configuration.
	return bedrockruntime.NewFromConfig(cfg), nil
}

// Chat sends a chat request to the Amazon Bedrock API using the Messages API format.
// It returns a single, complete response from the model.
func (p *NovaProvider) Chat(systemPromptText, userPrompt string) (string, error) {

	ctx := context.Background()
	// Create a new Bedrock client.
	client, err := newBedrockClient(ctx, p.Profile)
	if err != nil {
		return "", err
	}

	// Construct the user message for the API request.
	messages := []novaMessage{
		{
			Role: "user",
			Content: []novaMessageContent{
				{Text: userPrompt},
			},
		},
	}

	// Construct the system prompt as a slice of structs, only if it's not empty.
	var systemContent []novaSystemPrompt
	if systemPromptText != "" {
		systemContent = append(systemContent, novaSystemPrompt{Text: systemPromptText})
	}

	// Build the request body for the InvokeModel API call.
	// InferenceConfig is initialized directly with default values.
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

	// Invoke the Bedrock model.
	output, err := client.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(p.Profile.Model),
		ContentType: aws.String("application/json"),
		Body:        jsonBody,
	})
	if err != nil {
		return "", fmt.Errorf("failed to invoke model: %w", err)
	}

	responseBodyBytes := output.Body

	// Attempt to unmarshal into the success response structure.
	var novaResp novaCombinedAPIResponse
	if err := json.Unmarshal(responseBodyBytes, &novaResp); err == nil {
		// Success case: Extract text from the response.
		if len(novaResp.Output.Message.Content) > 0 {
			return novaResp.Output.Message.Content[0].Text, nil
		}
		return "", fmt.Errorf("no content found in response")
	}

	// If unmarshaling into the success structure fails, try to unmarshal into the error structure.
	var errorResp bedrockErrorResponse
	if err := json.Unmarshal(responseBodyBytes, &errorResp); err == nil {
		return "", fmt.Errorf("model error (%s): %s", errorResp.Type, errorResp.Message)
	}

	// If both fail, return a generic error with the raw response body.
	return "", fmt.Errorf("failed to unmarshal response body: %s", string(responseBodyBytes))
}

// ChatStream sends a streaming chat request to the Amazon Bedrock API using the Messages API format.
// It streams response chunks to the provided channel.
func (p *NovaProvider) ChatStream(ctx context.Context, systemPromptText, userPrompt string, responseChan chan<- string) error {
	// Note: The caller is responsible for closing the responseChan.

	// Create a new Bedrock client.
	client, err := newBedrockClient(ctx, p.Profile)
	if err != nil {
		return fmt.Errorf("error creating bedrock client: %w", err)
	}

	// Construct the user message for the API request.
	messages := []novaMessage{
		{
			Role: "user",
			Content: []novaMessageContent{
				{Text: userPrompt},
			},
		},
	}

	// Construct the system prompt as a slice of structs, only if it's not empty.
	var systemContent []novaSystemPrompt
	if systemPromptText != "" {
		systemContent = append(systemContent, novaSystemPrompt{Text: systemPromptText})
	}

	// Build the request body for the InvokeModelWithResponseStream API call.
	// InferenceConfig is initialized directly with default values.
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

	// Invoke the Bedrock model with streaming response.
	output, err := client.InvokeModelWithResponseStream(ctx, &bedrockruntime.InvokeModelWithResponseStreamInput{
		ModelId:     aws.String(p.Profile.Model),
		ContentType: aws.String("application/json"),
		Body:        jsonBody,
	})
	if err != nil {
		return fmt.Errorf("failed to invoke model with stream: %w", err)
	}

	// Process the streaming events.
	stream := output.GetStream()
	for event := range stream.Events() {
		select {
		case <-ctx.Done():
			// If context is cancelled, close the stream and return the context error.
			stream.Close()
			return ctx.Err()
		default:
		}

		switch v := event.(type) {
		case *types.ResponseStreamMemberChunk:
			var chunk novaMessagesAPIStreamChunk
			// Unmarshal the chunk bytes into the streaming response structure.
			if err := json.Unmarshal(v.Value.Bytes, &chunk); err != nil {
				// Log the error but continue processing, as some chunks might be malformed or unexpected.
				fmt.Fprintf(os.Stderr, "Error unmarshaling stream chunk: %v\n", err)
				continue
			}
			// Send the text content to the response channel.
			responseChan <- chunk.ContentBlockDelta.Delta.Text
		// Handle other event types if necessary.
			// For example, *types.ResponseStreamMemberContentBlockStart, *types.ResponseStreamMemberContentBlockStop, etc.
			// fmt.Fprintf(os.Stderr, "unhandled stream event type: %T\n", v)
		}
	}

	// After the loop, check for any errors that occurred during streaming.
	if err := stream.Err(); err != nil {
		return fmt.Errorf("streaming error: %w", err)
	}

	return nil
}

// awsCredentials represents the structure of the AWS credentials JSON file.
type awsCredentials struct {
	AWSAccessKeyID     string `json:"aws_access_key_id"`
	AWSSecretAccessKey string `json:"aws_secret_access_key"`
}

// loadAWSCredentialsFromFile loads AWS credentials from a specified JSON file.
func loadAWSCredentialsFromFile(filePath string) (*awsCredentials, error) {
	resolvedPath, err := appconfig.ResolvePath(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve credentials file path %s: %w", filePath, err)
	}

	data, err := os.ReadFile(resolvedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials file %s: %w", resolvedPath, err)
	}

	var creds awsCredentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, fmt.Errorf("failed to unmarshal credentials from file %s: %w", resolvedPath, err)
	}

	if creds.AWSAccessKeyID == "" || creds.AWSSecretAccessKey == "" {
		return nil, fmt.Errorf("aws_access_key_id or aws_secret_access_key is missing in credentials file %s", resolvedPath)
	}

	return &creds, nil
}
