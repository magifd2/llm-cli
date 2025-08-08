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

### Build Commands

*   **`make build`**: Builds a binary for your current OS/architecture. This is useful for local testing.
*   **`make all`**: Runs the entire release build process, including vulnerability checks, cross-compilation for all target platforms (macOS, Linux, Windows), and packaging the binaries into release archives.
*   **`make test`**: Runs the project's automated tests.
*   **`make clean`**: Deletes build artifacts.

### Installation

To install the binary for local use, run `make install`. This will build the binary and place it in `$(PREFIX)/bin` (default: `/usr/local/bin`).

```bash
# Install to /usr/local/bin (requires sudo)
sudo make install

# Install to your user's bin directory
make install PREFIX=~
```

To enable shell completion, follow the instructions from `llm-cli completion zsh --help` (or `bash`, `fish`, etc.).

## Code Style and Quality

- Adhere to standard Go formatting (`gofmt`).
- Run `make lint` to check for style issues before committing.

## Testing

- Write unit tests for new features and bug fixes.
- Run `make test` to execute the test suite.

## Commit Message Guidelines

- Use the [Conventional Commits](https://www.conventionalcommits.org/) specification (e.g., `feat:`, `fix:`, `refactor:`, `docs:`).
- Explain *why* a change was made, not just *what* was changed.

## Security

- Follow the **Security First Principle**: Security overrides other considerations.
- Scan for vulnerabilities using `make vulncheck`.
- Never hardcode sensitive information. Validate all external inputs.

## Documentation

- **Language Policy**: **English is the primary language** for all code, comments, and documentation. Japanese documentation is provided as a translation. In case of discrepancies, the English version takes precedence.
- When a feature is changed or added, ensure all relevant documentation (`README.md`, `BUILD.md`, etc.) is updated accordingly.

## Adding a New LLM Provider

The project uses a modular, package-per-provider architecture. To add a new provider, you must create a new self-contained package.

### Step 1: Create the Provider Package

Create a new directory under `internal/llm/`. The directory name should be the name of your provider (e.g., `myprovider`).

```bash
mkdir internal/llm/myprovider
```

### Step 2: Implement the `Provider` Interface

Inside your new directory, create a `provider.go` file. In this file, define a struct for your provider and implement the `llm.Provider` interface.

**`internal/llm/myprovider/provider.go`:**
```go
package myprovider

import (
	"context"
	"fmt"

	"github.com/magifd2/llm-cli/internal/config"
	"github.com/magifd2/llm-cli/internal/llm"
)

// Provider implements the llm.Provider interface for MyProvider.
type Provider struct {
	Profile config.Profile
}

// NewProvider is a factory function that returns a new MyProvider.
func NewProvider(p config.Profile) (llm.Provider, error) {
	return &Provider{}, nil
}

// Chat handles non-streaming requests.
func (p *Provider) Chat(systemPrompt, userPrompt string) (string, error) {
	// TODO: Implement the logic to call your provider's API.
	return "", fmt.Errorf("Chat not implemented for MyProvider")
}

// ChatStream handles streaming requests.
func (p *Provider) ChatStream(ctx context.Context, systemPrompt, userPrompt string, responseChan chan<- string) error {
	// TODO: Implement the logic for streaming.
	// Remember: DO NOT close the responseChan. It is managed by the caller.
	return "", fmt.Errorf("ChatStream not implemented for MyProvider")
}
```

### Method Details

#### `NewProvider(p config.Profile) (llm.Provider, error)`

*   This is a factory function that all provider packages must export.
*   It takes a `config.Profile` as input, which contains the configuration settings for the specific provider.
*   It should return an instance of the provider (which implements the `llm.Provider` interface) and an `error` if the provider cannot be created (e.g., due to an invalid profile configuration).
*   This function is used by the central provider registry in `cmd/providers.go` to dynamically create provider instances.

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
*   **Crucial Convention**: The `ChatStream` implementation must **NEVER** close the `responseChan`. The channel's lifecycle is managed by the caller in `cmd/prompt.go`.
*   If an error occurs at any point (before or during the stream), the function should stop processing and return an `error`.
*   The `context.Context` should be respected to handle cancellation requests from the user (e.g., Ctrl+C).

### Step 3: Activate the Provider

Finally, make the CLI aware of your new provider. Open `cmd/providers.go`:

1.  Add your new package to the `import` block.
    ```go
    import (
        // ... other imports
        "github.com/magifd2/llm-cli/internal/llm/myprovider"
    )
    ```

2.  Add your new provider to the `providerRegistry` map.
    ```go
    // cmd/providers.go
    var providerRegistry = map[string]providerFactory{
        // ... other providers
        "myprovider": myprovider.NewProvider,
    }
    ```

After these steps, a user can set `provider: myprovider` in their profile, and `llm-cli` will use your new implementation.