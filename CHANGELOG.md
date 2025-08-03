# CHANGELOG

## v0.0.9 - 2025-08-03

### ✨ Features
*   **OpenAI API Key File Support**: Added support for loading OpenAI API keys from a JSON file specified by `credentials-file`.
    *   The JSON file should contain the API key under the `openai_api_key` field.
    *   `credentials-file` takes precedence over `api_key` directly set in the profile.
*   **Enhanced Profile Check Command**: The `llm-cli profile check` command now verifies the existence of credential files specified in profiles.
    *   Displays warnings if a specified credential file does not exist or cannot be resolved.

### 🐛 Bug Fixes
*   **Correct Profile Check Logic**: Fixed an issue in `llm-cli profile check` where it would unnecessarily prompt to update 'limits' settings even when they already matched standard default values.
    *   The command now only prompts for updates if the 'limits' are at their zero value or meaningfully different from the standard defaults.
*   **Build Fix**: Resolved unused import in `cmd/root.go` that caused `govulncheck` errors during build.

### 📝 Documentation
*   **Updated READMEs**: Clarified `~` (tilde) path notation and added usage examples for OpenAI API key file support in `README.md` and `README.ja.md`.

### ♻️ Refactor
*   **Error Handling**: Centralized error handling in `main.go` and `cmd/root.go` for consistent error processing and exit codes.

## v0.0.8 - 2025-08-03

### ✨ Features
*   **Profile Check Command**: Added `llm-cli profile check` command to verify and migrate configuration profiles.
    *   Checks all profiles for consistency, especially for newly introduced settings like 'limits'.
    *   Prompts the user to update profiles with default or unconfigured 'limits' settings to standard default values.
    *   Includes a `--confirm` or `-y` flag for non-interactive operation.
    *   Creates a timestamped backup of the `config.json` file in `~/.config/llm-cli/backups/` before saving any changes, ensuring data safety.
    *   Enhanced the `profile show` command to display 'limits' information.

### 🐛 Bug Fixes
*   **DoS Protection Enhancements**: Addressed remaining issues related to DoS protection and UTF-8 safety.
    *   Modified `readAndProcessStream` to stop reading input once the `MaxPromptSizeBytes` limit is reached, even in "warn" mode, preventing large files from being fully loaded into memory.
    *   Updated `truncateStringByBytes` to be UTF-8 aware, ensuring that string truncation for size limiting does not corrupt multi-byte characters.
*   **Configuration Loading Consistency**: Ensured `Limits` struct is always initialized with default values when loading configuration, even if not explicitly present in the config file. This guarantees consistent behavior across all profiles.


## v0.0.7 - 2025-08-03

### ✨ Features
*   **DoS Protection via Size Limits**: Implemented configurable input and output size limits to prevent accidental resource exhaustion.
    *   Added a new `limits` object to profiles in `config.json`.
    *   `prompt` command now checks these limits, with configurable behavior (`stop` or `warn`) on excess.
    *   `add` command supports setting limits for new profiles via flags (e.g., `--limits-max-prompt-size-bytes`).
    *   `set` command can now configure limits using dot notation (e.g., `limits.on_input_exceeded`).
    *   `list` command now displays the configured limits for each profile, enhancing user visibility.
*   **Improve UX with Spinner**: Added a spinner to the `prompt` command when not using `--stream` mode to provide visual feedback during long-running operations. The spinner is only displayed in interactive terminals.

### 🐛 Bug Fixes
*   **CLI Behavior**: Suppressed the automatic display of usage instructions on runtime errors (e.g., API failures) to provide cleaner error output.

### ♻️ Refactor
*   **Improve Testability of Profile Commands**: Refactored `profile` subcommands (`add`, `set`, `use`, `remove`) to return errors instead of calling `os.Exit`, making them testable. The test suite was updated to execute commands directly and validate their behavior, improving test reliability.

## v0.0.6 - 2025-08-02

### ✨ Features
*   **Vulnerability Check Integration**: Added `govulncheck` to Makefile and integrated it into the build process to automatically scan for known vulnerabilities.
*   **macOS Ad-hoc Signing**: Implemented ad-hoc signing for macOS universal binaries in the Makefile to allow execution on machines other than the build machine.

### 🐛 Bug Fixes
*   **Build Fix**: Corrected a syntax error in `cmd/list.go` where the `os` package import was missing quotes, resolving persistent "missing import path" build errors.

### 🔒 Security
*   **Enhanced Configuration File Permissions**: Implemented more restrictive file permissions for `~/.config/llm-cli/config.json` (changed from `0644` to `0600`) and its parent directory (changed from `0755` to `0700`) to enhance the security of sensitive information like API keys.

## v0.0.5 - 2025-08-02

### ✨ Features
*   **Google Cloud Vertex AI Provider Support**: Added support for interacting with Google Cloud Vertex AI.
*   **Enhanced `profile add` Command**: The `profile add` command now allows specifying parameters such as provider, model, endpoint, API key, AWS credentials, GCP Project ID, location, and credentials file path in a single command.

### ♻️ Refactor
*   **Vertex AI SDK Migration**: Migrated to the latest `google.golang.org/genai` SDK, including fixes for service account authentication, correct `Client` object usage, and proper streaming iterator handling.
*   **Runtime Expansion of Credential File Paths**: Changed the expansion of `credentials_file` paths to occur at runtime instead of at configuration time, providing greater flexibility in dynamic home directory environments.

### 📝 Documentation
*   **Updated Development Log**: Added detailed history of Vertex AI SDK migration and current SDK status to `DEVELOPMENT_LOG.md`.
*   **Updated Related Documentation**: `README.ja.md` and `README.en.md` have been updated to reflect the addition of the Vertex AI provider and the enhanced `profile add` command usage, and system prompt handling approach.
*   **Revised Provider Development Guide**: Removed specific provider implementation details from `DEVELOPING_PROVIDERS.ja.md` and `DEVELOPING_PROVIDERS.en.md`.
*   **Updated Changelogs**: `CHANGELOG.ja.md` and `CHANGELOG.en.md` have been updated.

### ♻️ Refactor
*   **Code Audit and Quality Improvements**: Performed a full code audit and fixed potential bugs and vulnerabilities. Hardened against command injection in `profile edit`, centralized config path management, and improved error messages to enhance robustness and maintainability.

## v0.0.4 - 2025-08-01

### 🐛 Bug Fixes
*   **API Error Handling**: Fixed an issue where API errors during streaming mode were not detected, causing the program to exit silently. Resolved a race condition in asynchronous handling to make error reporting more robust.

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

*   **Amazon Bedrock Nova Model Support**: Added support for interacting with Nova models suchs as `amazon.nova-lite-v1:0`.

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