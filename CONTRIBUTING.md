# Contributing to llm-cli

We welcome contributions to `llm-cli`! Whether it's bug fixes, new features, or documentation improvements, your help is greatly appreciated.

Please take a moment to review this document to understand how to contribute effectively.

## Table of Contents
- [How to Contribute](#how-to-contribute)
- [Development Environment Setup](#development-environment-setup)
- [Code Style and Quality](#code-style-and-quality)
- [Testing](#testing)
- [Commit Message Guidelines](#commit-message-guidelines)
- [Security](#security)
- [Documentation](#documentation)
- [Adding a New LLM Provider](#adding-a-new-llm-provider)

## How to Contribute
1.  **Fork the Repository**: Start by forking the `llm-cli` repository on GitHub.
2.  **Clone Your Fork**: Clone your forked repository to your local machine.
    ```bash
    git clone https://github.com/YOUR_USERNAME/llm-cli.git
    cd llm-cli
    ```
3.  **Create a New Branch**: Create a new branch for your feature or bug fix.
    ```bash
    git checkout -b feature/your-feature-name
    ```
4.  **Make Your Changes**: Implement your changes, adhering to the [Code Style and Quality](#code-style-and-quality) guidelines.
5.  **Test Your Changes**: Ensure your changes pass existing tests and add new tests if necessary. See [Testing](#testing).
6.  **Commit Your Changes**: Write clear and concise commit messages. See [Commit Message Guidelines](#commit-message-guidelines).
7.  **Push Your Branch**: Push your changes to your fork on GitHub.
    ```bash
    git push origin feature/your-feature-name
    ```
8.  **Open a Pull Request**: Open a pull request from your branch to the `main` branch of the original `llm-cli` repository. Provide a clear description of your changes.

## Development Environment Setup

### Prerequisites

*   [Go](https://go.dev/doc/install) (version 1.21 or later recommended)
*   [Git](https://git-scm.com/)
*   `make` command
    *   Standard on macOS/Linux.
    *   For Windows, please install [Make for Windows](http://gnuwin32.sourceforge.net/packages/make.htm) or similar.

### Build Commands

The following `make` commands are available. The built binaries will be generated in the `bin/` directory.

*   **`make build`**
    *   Builds a binary for the currently used OS and architecture. For macOS, this will automatically produce a universal binary (supporting both `amd64` and `arm64` architectures). The built binary will be placed in `bin/<OS>-<ARCH>/llm-cli` (or `bin/darwin-universal/llm-cli` for macOS). This is convenient for testing during development.

*   **`make cross-compile`**
    *   For distribution, builds binaries for multiple OS and architectures at once and creates compressed archives. The artifacts will be generated in the `bin/` directory.
        *   `bin/llm-cli-darwin-universal.tar.gz` (macOS Universal Binary)
        *   `bin/llm-cli-linux-amd64.tar.gz` (Linux amd64)
        *   `bin/llm-cli-windows-amd64.zip` (Windows amd64)

*   **`make all`**
    *   Executes both `make build` and `make cross-compile`. Creates a binary for the current OS and architecture, as well as all cross-compiled binaries and archives.
*   **`make test`**
    *   Runs the project's tests.

*   **`make clean`**
    *   Deletes the `bin/` directory and build cache.

### Installation and Uninstallation

`llm-cli` can be installed and uninstalled using the `Makefile` targets.

#### `make install`

This target builds the `llm-cli` binary and installs it to a specified directory, along with the Zsh shell completion script. The default installation path is `/usr/local/bin`.

*   **Default Installation (System-wide):**
    To install `llm-cli` to `/usr/local/bin` (requires `sudo`):
    ```bash
    sudo make install
    ```

*   **User-local Installation:**
    To install `llm-cli` to `~/bin` (recommended for non-root users, ensure `~/bin` is in your `PATH`):
    ```bash
    make install PREFIX=~
    ```

*   **Custom Directory Installation:**
    To install `llm-cli` to a custom directory (e.g., `/opt/llm-cli`):
    ```bash
    sudo make install PREFIX=/opt/llm-cli
    ```

After installation, for Zsh users, you might need to run `compinit` or restart your shell for the completion script to take effect.

#### `make uninstall`

This target removes the `llm-cli` binary and its associated completion script from the installation directory. It is crucial to use the same `PREFIX` value that was used during installation.

*   **Default Uninstallation:**
    ```bash
    sudo make uninstall
    ```

*   **User-local Uninstallation:**
    ```bash
    make uninstall PREFIX=~
    ```

*   **Custom Directory Uninstallation:
    ```bash
    sudo make uninstall PREFIX=/opt/llm-cli
    ```

**Note:** The uninstallation process does NOT remove your configuration files located at `~/.config/llm-cli/config.json`. These files contain your LLM profiles and are preserved across installations/uninstallations.

## Code Style and Quality

### Code Style and Formatting

- Adhere to standard Go formatting (`gofmt`).
- Follow idiomatic Go practices.
- Keep functions concise and focused on a single responsibility.

### Linting
- Use `golangci-lint` for static code analysis.
- Ensure all code passes lint checks before committing.

## Testing

### Testing Principles

- Write unit tests for new features and bug fixes.
- For critical bug fixes, especially those related to core logic like API interaction or concurrency, add a regression test to prevent recurrence.
- Ensure tests cover critical paths and edge cases.
- Use `make test` to run tests.

## Commit Message Guidelines

### Commit Message Conventions

- Use the Conventional Commits specification (e.g., `feat:`, `fix:`, `refactor:`, `docs:`).
- For multi-line commit messages, write the message in a temporary file (e.g., `.git/COMMIT_MSG`) and use `git commit -F <file>` to avoid shell interpretation errors. This is the standard procedure.
- Explain *why* a change was made, not just *what* was changed.

## Security

### Security First Principle

- Security is the highest priority, overriding all other considerations such as functionality or performance.
- All code, dependencies, and configurations must be reviewed for potential security vulnerabilities before being committed.
- Never trust user input, including environment variables. All external inputs must be validated and sanitized to prevent injection attacks.
- Sensitive information (API keys, credentials) must never be hardcoded or stored in insecure locations.

### Secure Development Lifecycle

- **Threat Modeling at Design Phase:** Before implementing a new feature, consider potential threats. For example, when adding a feature that interacts with the filesystem, evaluate risks like path traversal.
- **Security-Focused Code Reviews:** All code reviews must include a specific check for security vulnerabilities. Do not approve pull requests that have not been reviewed from a security perspective.
- **Safe Testing Practices:** When testing for vulnerabilities, use harmless proof-of-concept payloads. Before running tests that involve external inputs like environment variables, always inspect their contents first.
- **Dependency Scanning:** Regularly scan project dependencies for known vulnerabilities using tools like `govulncheck`.

## Documentation

### Documentation Principles

- **Language Policy**: All documentation will be written in Japanese first (as the primary source of truth) and then translated into English.
  - The English version should include a note indicating it is a translation and that the Japanese version takes precedence in case of discrepancies.
- **Scope**: Maintain both user-facing documents (e.g., `README`) and developer-facing documents (e.g., `DEVELOPING_PROVIDERS.md`).
- **Maintenance**: When a feature is changed or added, ensure all relevant documentation is updated accordingly.

## Adding a New LLM Provider

The core of the provider system is the `Provider` interface, defined in `internal/llm/provider.go`. Any new provider must implement this interface.

```go
package llm

import (
	"context"
)

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

### Step-by-Step Implementation Guide

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

	return "", fmt.Errorf("ChatStream not implemented for MyProvider")
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
