# Development Log

This document records the detailed development history and key decisions made during the project.

## 2025-08-02 (Build Fix: Corrected `cmd/list.go` Import Statement)

- **Resolved Build Error**:
  - Identified and corrected a syntax error in `cmd/list.go` where the `os` package import was missing quotes (`import (os)` instead of `import ("os")`).
  - This syntax error was the root cause of the persistent "missing import path" build errors, which had previously masked other issues and led to extensive troubleshooting.
  - The fix was applied manually by the user.
  - Verified successful build after the correction.

## 2025-08-02 (Security Enhancements: Configuration File Permissions)

- **Enhanced Security for Configuration File**:
  - Modified `internal/config/config.go` to set more restrictive file permissions for `~/.config/llm-cli/config.json` and its parent directory.
  - Changed `os.WriteFile` permission from `0644` to `0600` for `config.json`.
  - Changed `os.MkdirAll` permission from `0755` to `0700` for the `llm-cli` configuration directory.
  - This ensures that only the file owner can read and write the configuration file, enhancing the security of sensitive information like API keys.

## 2025-08-02 (Release v0.0.5 and Documentation/Build System Enhancements)

- **Release v0.0.5 Preparation**:
  - Updated `CHANGELOG.ja.md` and `CHANGELOG.md` to reflect the `v0.0.5` release.
  - Added application version (`ver.0.0.5`) to `cmd/root.go`.
  - Improved `Short` and `Long` descriptions in `cmd/root.go` for better clarity.
  - Executed `make all` to build and test the application.
  - Created Git tag `v0.0.5` (corrected from `ver.0.0.5`).

- **Code Comment Cleanup**:
  - Standardized all code comments to English only across `cmd/`, `internal/config/`, and `internal/llm/` directories.
  - Added copyright headers to test files (`cmd/edit_test.go`, `cmd/profile_commands_test.go`, `internal/config/config_test.go`).

- **Makefile Enhancements**:
  - Added `install` and `uninstall` targets for easier binary and completion script management.
  - Implemented `PREFIX` variable for flexible installation paths (system-wide or user-local).
  - Ensured macOS builds always produce universal binaries (`amd64` and `arm64`).

- **Documentation Overhaul**:
  - Updated `README.md`, `README.ja.md`, `BUILD.md`, `BUILD.ja.md` to reflect new installation procedures and macOS universal binary builds.
  - Created `CONTRIBUTING.md` by integrating development guidelines from `GEMINI.md` and provider development guides.
  - Removed redundant `DEVELOPING_PROVIDERS.en.md` and `DEVELOPING_PROVIDERS.ja.md`.
  - Created `SECURITY.md` with vulnerability reporting guidelines and security principles.
  - Updated `README.md` and `README.ja.md` to link to `SECURITY.md`.
  - Updated `README.md` and `README.ja.md` to include API key configuration for OpenAI-compatible APIs.
  - Updated `README.md` and `README.ja.md` to include `--profile` option for `prompt` command.
  - Updated `README.md` and `README.ja.md` to include `show` subcommand for `profile` command.
  - Standardized English documentation file names (e.g., `README.en.md` to `README.md`).
  - Updated language policy in `.gemini/GEMINI.md` to reflect the new documentation naming convention.

## 2025-08-02 (Vertex AI SDK Migration and System Prompt Workaround)

- **Initial Attempt with Incorrect SDK**: Began migrating Vertex AI provider to a new SDK based on external developer's provided code (`google.golang.org/genai/vertexai`). This led to persistent build errors due to incorrect package paths and API usage.
- **Misunderstanding of `google.golang.org/genai`**: Repeated attempts to fix build errors by forcing `google.golang.org/genai` to `v0.4.0` via `go.mod` `replace` directives, based on external advice, proved ineffective and highlighted a fundamental misunderstanding of the SDK's structure.
- **Clarification of `genai.NewClient` Signature**: Through direct consultation of `https://pkg.go.dev/google.golang.org/genai`, it was clarified that `genai.NewClient` expects a `*genai.ClientConfig` object, not variable arguments of `ClientOption`s directly.
- **Authentication Implementation**: Correctly implemented service account authentication by reading the JSON key file, parsing it, and using `auth.New2LOTokenProvider` to create a `TokenProvider` which is then set in `genai.ClientConfig.Credentials`.
- **SDK Instability and Future Outlook**: It has become apparent that the Go client library for the GenAI SDK (`google.golang.org/genai`) is still in an early and somewhat unstable state. Key observations include:
    - Lack of a `Close()` method on the `genai.Client` object, which is unusual for client libraries managing network connections.
    - Inconsistent API behavior and documentation discrepancies encountered during the migration process.
    - The need for workarounds (e.g., for system prompts) due to missing direct API support.
    We will continue to monitor the updates to the `google.golang.org/genai` SDK and adapt our implementation as the library matures and stabilizes. Future enhancements will prioritize aligning with official best practices as they evolve.
- **System Prompt Workaround**: Since Vertex AI's GenAI SDK does not directly support system prompts, a workaround was implemented. The system prompt content is now sent as the first message in the chat history for both `Chat` and `ChatStream` functions.
- **Streaming Iterator Fix**: Corrected the usage of `chat.SendMessageStream` to iterate over its results using a `for ... range` loop, resolving `iter.Next undefined` errors.
- **Successful Build and Verification**: After numerous iterations and careful adherence to the official SDK documentation, the application now builds successfully and the Vertex AI provider functions as expected with the new SDK.

## 2025-08-01 (Code Audit and Refactoring)

- **Conducted a full code audit** to identify potential bugs and vulnerabilities.
- **Identified and fixed a potential OS command injection vulnerability** in the `profile edit` command by validating the editor path with `exec.LookPath`.
- **Refactored configuration path handling** by creating `config.GetConfigPath()` to centralize the logic and prevent inconsistencies between `cmd/edit.go` and `internal/config/config.go`.
- **Improved error handling** for stdin reading and provided more user-friendly error messages for invalid keys in `profile set`.
- **Updated placeholder descriptions** in Cobra commands for better clarity.
- **Identified a significant lack of test coverage** as a major technical debt. Decided to prioritize test implementation in the next development phase.

## 2025-08-01 (LLM Provider Test Strategy Re-evaluation)

- **Attempted to implement comprehensive unit tests for `internal/llm` package**:
  - Initial strategy involved using `httptest.Server` for mocking and `io.Pipe` for streaming control.
  - Encountered severe deadlocking issues during `ChatStream` cancellation tests, primarily due to the blocking nature of `io.Reader` operations and the inability of `context.Context` to directly interrupt them.
  - Multiple attempts were made to resolve the deadlocks by refactoring `ollama.go` (e.g., `bufio.Scanner` to `bufio.NewReader`, goroutine-based `ReadBytes` with channels) and refining test mocks (e.g., `http.CloseNotifier`).
  - Despite extensive debugging and re-evaluation, a stable, non-deadlocking test pattern for streaming cancellation could not be established without fundamentally altering the core `internal/llm` implementation.
- **Decision**: Due to the high complexity and significant development overhead, the implementation of dedicated unit tests for the `internal/llm` package is **temporarily frozen**.
- **Policy Update**: A new policy has been added to `DEVELOPMENT_PLAN.md` and `GEMINI.md` strictly prohibiting modifications to existing, functionally verified code within `internal/llm` unless critical for security or essential functionality, to ensure operational stability.
- **Future Outlook**: Future enhancements or refactoring in this area will require a thoroughly re-evaluated and approved testing strategy, potentially involving a "Block-2" provider model with a redesigned I/O layer.

## 2025-08-01 (Initial Work)

### Amazon Bedrock Nova Model Implementation
...

- **Initial Attempt & Debugging**: Began implementing Amazon Bedrock support. Encountered numerous `ValidationException` errors due to incorrect request/response formats for Nova models (e.g., `required key [messages] not found`, `extraneous key [maxTokenCount] is not permitted`, `extraneous key [temperature] is not permitted`, `extraneous key [topP] is not permitted`, `expected type: JSONArray, found: String`).
- **API Guide Consultation**: Utilized the provided Nova API guide (https://docs.aws.amazon.com/nova/latest/userguide/using-invoke-api.html) to understand the correct JSON structures for requests and streaming responses.
- **Refactoring `internal/llm/bedrock.go`**: Completely re-implemented the Bedrock provider to strictly adhere to Nova Messages API specification.
  - Redefined request/response structures (`novaMessageContent`, `novaMessage`, `novaSystemPrompt`, `novaMessagesAPIRequest`, `novaCombinedAPIResponse`, `novaMessagesAPIStreamChunk`).
  - Adjusted prompt and inference parameter handling (e.g., nesting parameters under `inferenceConfig`).
  - Corrected streaming response parsing to extract text from `contentBlockDelta.delta.text`.
- **Streaming Debugging**: Initially, streaming produced no output. Debugging revealed that `novaMessagesAPIStreamChunk`'s `Type` and `Delta.Type` fields were empty, indicating a mismatch with the actual JSON. Removed these fields from the struct and simplified the streaming logic to directly process `ContentBlockDelta.Delta.Text`.

### Command Renaming and Flag Refactoring

- **`ask` to `prompt`**: Renamed the `ask` command to `prompt` for better intuitiveness.
- **Flag Renaming**: Renamed `--prompt` to `--user-prompt` and introduced shorthand flags (`-p`, `-P`, `-f`, `-F`) for user and system prompts/files.
- **Prompt Validation**: Centralized prompt validation in `cmd/prompt.go` (formerly `cmd/ask.go`), ensuring that a prompt is always provided via flags, positional arguments, or stdin.

### Documentation Multilingual Support

- **File Renaming**: Renamed `README.md`, `BUILD.md`, `CHANGELOG.md` to `README.ja.md`, `BUILD.ja.md`, `CHANGELOG.ja.md` respectively.
- **English Versions**: Created `README.en.md`, `BUILD.en.md`, `CHANGELOG.en.md` by translating the Japanese content.
- **Language Selector**: Created a new root `README.md` to serve as a language selection page.

### Makefile Improvements

- **`make all` Enhancement**: Modified `Makefile` so that `make all` now executes both `make build` (for current OS/ARCH) and `make cross-compile` (for all platforms), generating all binaries and archives.
- **Build Output Paths**: Ensured `make build` outputs binaries to `bin/<OS>-<ARCH>/llm-cli` for better organization.

### Tooling and Development Rules

- **`golangci-lint` Introduction**: Integrated `golangci-lint` for static code analysis. Added `lint` target to `Makefile`.
- **`GEMINI.md` Guidelines**: Established `GEMINI.md` as a project guideline document, including rules for:
  - Code Style and Formatting (including meaningful comments).
  - Testing Principles.
  - Rollback Strategy (commit before major changes).
  - Commit Message Conventions.
  - Bug Fixes (root cause analysis).
  - Code Quality and Security (robustness, safety, maintainability, security-first).
  - Documentation Principles (multilingual, update on feature change).
  - File and Directory Operations (absolute paths).