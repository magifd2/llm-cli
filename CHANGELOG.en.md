# CHANGELOG

## v0.0.3 - 2025-08-01

### 🚨 Breaking Changes
*   **Command Name and Flag Changes**: The `ask` command has been renamed to `prompt`. The `--prompt` flag has been renamed to `--user-prompt`, and `--prompt-file` has been renamed to `--user-prompt-file`. Existing scripts and workflows will need to be updated.

### ✨ Features
*   **Command Name and Flag Refactoring**: Renamed the `ask` command to `prompt` and the `--prompt` flag to `--user-prompt`. Added shorthand flags for `--user-prompt` (`-p`), `--system-prompt` (`-P`), `--user-prompt-file` (`-f`), and `--system-prompt-file` (`-F`).

### 📝 Documentation
*   **Added Development Log and Plan**: `DEVELOPMENT_LOG.md` and `DEVELOPMENT_PLAN.md` were added to record project history and future plans.
*   **Updated Development Rules**: `GEMINI.md` was updated with new development rules (e.g., absolute paths for file operations, root cause bug fixes, documentation updates, pre-commit commits, code quality and security, comment principles).

## v0.0.2 - 2025-08-01

### ✨ Features

*   **Amazon Bedrock Nova Model Support**: Added support for interacting with Nova models such as `amazon.nova-lite-v1:0`.

### ♻️ Refactor

*   **Bedrock Provider Refactoring**: The internal implementation of the Bedrock provider has been redesigned to strictly conform to the Nova model's Messages API specification.
    *   Updated request/response structs (`novaMessageContent`, `novaMessage`, `novaSystemPrompt`, `novaMessagesAPIRequest`, `novaCombinedAPIResponse`, `novaMessagesAPIStreamChunk`).
    *   Adjusted prompt and inference parameter handling to align with the Nova API.
    *   Fixed streaming response parsing logic.
*   **Improved Provider Selection Logic**: Modified `cmd/ask.go` to dynamically select `NovaBedrockProvider` based on the model ID prefix (`amazon.nova`) when the `bedrock` provider is chosen.

### 🐛 Bug Fixes

*   **Prompt Validation Fix**: Centralized prompt validation in `cmd/ask.go` and removed redundant validation from `internal/llm/bedrock_nova.go`. This makes a prompt from `--prompt`, `--prompt-file`, or standard input mandatory.

### 📝 Documentation

*   **Documentation Updates**: Updated `README.md` and `BUILD.md` to reflect the new Bedrock setup procedures and build process changes.
*   **Makefile Improvements**: Corrected the `make all` command to also perform cross-compilation and organized the build output directory.

## v0.0.1 - 2025-07-31

### ✨ Features

*   **LLM Interaction Functionality**:
    *   Supports interaction with Ollama and LM Studio (OpenAI-compatible API).
    *   Supports prompt input from command-line arguments, files, and standard input.
    *   Supports streaming response display (`--stream` flag).
*   **Profile Management Functionality**:
    *   Manages multiple LLM configurations as profiles (`profile list`, `profile use`, `profile add`, `profile set`, `profile remove`, `profile edit`).
*   **Build System**:
    *   Initialized Go modules and introduced the Cobra CLI framework.
    *   Supports build, test, and cross-compilation (macOS Universal, Linux, Windows) via `Makefile`.
    *   Organizes build artifacts into platform-specific directories.

### 📝 Documentation

*   Created `README.md` and added feature descriptions.
*   Created `BUILD.md` and detailed build procedures.
*   Explicitly stated Gemini development support and MIT License in `README.md`.
*   Reflected improvements suggested in code reviews in `BUILD.md` and `README.md`.