# Developer Guide: Adding a New Provider

> **Note:** This is a translation of the Japanese version (`DEVELOPING_PROVIDERS.ja.md`). If there are any discrepancies, the Japanese version takes precedence.

This guide explains how to add support for a new LLM provider to `llm-cli`.

## The `Provider` Interface

The core of the provider system is the `Provider` interface, defined in `internal/llm/provider.go`. Any new provider must implement this interface.

```go
package llm

import (
	"context"
)

// Provider defines the interface for interacting with different LLMs.
type Provider interface {
	// Chat sends a standard, non-streaming request to the LLM.
	Chat(systemPrompt, userPrompt string) (string, error)

	// ChatStream sends a streaming request to the LLM.
	ChatStream(ctx context.Context, systemPrompt, userPrompt string, responseChan chan<- string) error
}
```

### Method Details

#### `Chat(systemPrompt, userPrompt string) (string, error)`

*   This method handles a simple request-response cycle.
*   It should send the `systemPrompt` (if provided) and the `userPrompt` to the LLM's API.
*   It must block until the full response is received.
*   It should return the complete response text as a `string`.
*   If any error occurs (network, API error, etc.), it should return an `error`.

#### `ChatStream(ctx context.Context, systemPrompt, userPrompt string, responseChan chan<- string) error`

*   This method handles real-time, streaming responses.
*   It sends the prompts to the LLM's streaming API endpoint.
*   As response chunks (tokens) are received, they should be sent to the `responseChan` as `string`s.
*   **Crucial Convention**: The `ChatStream` implementation must **NEVER** close the `responseChan`. The channel's lifecycle is managed by the caller in `cmd/prompt.go`. Your implementation should simply send data to it.
*   If an error occurs at any point (before or during the stream), the function should stop processing and return an `error`.
*   The `context.Context` should be respected to handle cancellation requests from the user (e.g., Ctrl+C).

---

## Step-by-Step Implementation Guide

Here is how to create and integrate a new provider.

### Step 1: Create the Provider File

Create a new file in the `internal/llm/` directory. For example, `internal/llm/my_provider.go`.

### Step 2: Implement the Interface

In your new file, define a struct for your provider and implement the two required methods. You can use the following template as a starting point:

```go
package llm

import (
	"context"
	"fmt"

	appconfig "github.com/magifd2/llm-cli/internal/config"
)

// MyProvider implements the Provider interface for our new service.
type MyProvider struct {
	Profile appconfig.Profile
}

// Chat handles non-streaming requests for MyProvider.
func (p *MyProvider) Chat(systemPrompt, userPrompt string) (string, error) {
	// TODO: Implement the logic to call your provider's API.
	// 1. Construct the request body using the prompts.
	// 2. Send the HTTP request to the API endpoint (p.Profile.Endpoint).
	// 3. Handle the API response, checking for errors.
	// 4. Parse the response body to extract the message content.
	// 5. Return the content and a nil error.

	return "", fmt.Errorf("Chat not implemented for MyProvider")
}

// ChatStream handles streaming requests for MyProvider.
func (p *MyProvider) ChatStream(ctx context.Context, systemPrompt, userPrompt string, responseChan chan<- string) error {
	// TODO: Implement the logic for streaming.
	// 1. Construct the request for a streaming response.
	// 2. Send the HTTP request.
	// 3. Check for API errors before starting the stream.
	// 4. Read the response body line-by-line or chunk-by-chunk.
	// 5. For each chunk, parse it and send the text content to responseChan.
	// 6. Respect the context for cancellation (e.g., in your read loop).
	// 7. If an error occurs, return it immediately.

	return fmt.Errorf("ChatStream not implemented for MyProvider")
}

```

### Step 3: Activate the Provider

Finally, make the CLI aware of your new provider. Open `cmd/prompt.go` and find the `switch` statement inside the `Run` function. Add a new `case` for your provider.

```go
// cmd/prompt.go

// ...
        var provider llm.Provider
        switch activeProfile.Provider {
        case "ollama":
            provider = &llm.OllamaProvider{Profile: activeProfile}
        case "openai":
            provider = &llm.OpenAIProvider{Profile: activeProfile}
        case "bedrock":
            // ... (Bedrock logic)

        // Add your new provider here
        case "my_provider": // This string must match the 'provider' value in the config
            provider = &llm.MyProvider{Profile: activeProfile}

        default:
            fmt.Fprintf(os.Stderr, "Warning: Provider '%s' not recognized...\n", activeProfile.Provider)
            provider = &llm.MockProvider{}
        }
// ...
```

After these steps, a user can set `provider: my_provider` in their profile, and `llm-cli` will use your new implementation.
