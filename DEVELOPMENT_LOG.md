# Development Log

This document records the detailed development history and key decisions made during the project.

## 2025-08-03 (Fix: Correct Profile Check Logic for Limits Settings)

- **Problem**: The `llm-cli profile check` command would unnecessarily prompt to update 'limits' settings even when they already matched standard default values.
- **Analysis**: The logic for identifying 'default or unconfigured' limits was too broad, triggering prompts for profiles that had implicitly correct settings.
- **Solution**: Refined the comparison logic to only prompt for updates if the 'limits' are at their zero value or meaningfully different from the standard defaults.
- **Outcome**: Improved user experience by reducing unnecessary prompts during profile checks.

## 2025-08-03 (Feature: OpenAI API Key File Support)

- **Objective**: To enhance security by allowing OpenAI API keys to be loaded from external JSON files, aligning with best practices for sensitive credential management.
- **Key Changes & Decisions**:
    - Reused the existing `CredentialsFile` field in the `Profile` struct for OpenAI API key file paths, ensuring consistency with other providers (Bedrock, Vertex AI).
    - Defined the JSON file format to use `"openai_api_key"` as the key for the API token, providing clear identification.
    - Implemented logic in `internal/llm/openai.go` to prioritize the API key from `CredentialsFile` if specified, falling back to `APIKey` directly in the profile otherwise.
    - Ensured `config.ResolvePath` is used for secure path resolution of the credentials file.
- **Implementation Details**:
    - `internal/llm/openai.go`: Added `openAIAPIKey` struct and `loadOpenAIAPIKeyFromFile` function. Modified `Chat` and `ChatStream` methods to load API key from `CredentialsFile`.
    - CLI commands (`add`, `set`): No specific changes needed as they already handle `credentials-file`.
    - Documentation (`README.md`, `README.ja.md`): Updated to reflect new usage.
- **Outcome**: Improved security posture for OpenAI API key management.

## 2025-08-03 (Enhancement: Profile Check for Credential File Existence)

- **Objective**: To improve user experience by proactively identifying missing credential files linked to profiles.
- **Key Changes & Decisions**:
    - Extended the `llm-cli profile check` command to verify the existence of files specified in `profile.CredentialsFile`.
    - Implemented checks for both path resolution errors and file non-existence.
    - Displayed informative warning messages to the user for any issues found.
- **Implementation Details**: 
    - `cmd/profile.go`: Modified `checkCmd` to include credential file existence checks using `config.ResolvePath` and `os.Stat`.
- **Outcome**: Users can now more easily diagnose and correct issues related to missing credential files.

## 2025-08-03 (Fix: Build Errors after Error Handling Refactor)

- **Problem**: Introduced build errors (`too many return values`, `undefined: os`) after refactoring `main.go` and `cmd/root.go` for centralized error handling.
- **Analysis**: 
    - `too many return values`: Caused by `Chat` function (which returns `(string, error)`) attempting to return only `error` from `loadOpenAIAPIKeyFromFile`'s error path.
    - `undefined: os`: Caused by missing `os` import in `internal/llm/openai.go` after `os.Exit(1)` was removed from `cmd/root.go`.
- **Solution**: 
    - Corrected `Chat` function's error return to `return "", fmt.Errorf(...)`.
    - Added `os` import to `internal/llm/openai.go`.
- **Outcome**: Resolved compilation errors, ensuring the project builds successfully.

## 2025-08-03 (Feature: Enhanced Bedrock Credentials Handling and Profile Display)

- **Objective**: To improve the security and reusability of Bedrock credentials by allowing them to be loaded from external JSON files, and to enhance user transparency in profile display.
- **Key Changes & Decisions**:
    1.  **Unified `CredentialsFile`**: Initially, a separate `AWSCredentialsFile` field was considered for Bedrock. However, to reduce redundancy and maintain consistency with Vertex AI's existing `CredentialsFile`, it was decided to unify this field. The `Profile` struct now uses a single `CredentialsFile` field for both AWS and GCP credential file paths.
    2.  **External JSON Credential Files**: Bedrock credentials (`aws_access_key_id`, `aws_secret_access_key`) can now be stored in a separate JSON file, improving security by separating sensitive information from the main `config.json`.
    3.  **Path Resolution and Security**: Implemented `config.ResolvePath` to handle `~` (tilde) expansion and convert relative paths to absolute paths for credential files. This addresses potential path injection vulnerabilities and ensures consistent file access.
    4.  **Enhanced `profile show` Transparency**: The `profile show` command was updated to display not only the configured `CredentialsFile` path but also its resolved absolute path. This provides users with clear visibility into which file the application is actually accessing, addressing concerns about path interpretation.
    5.  **CLI Command Adjustments**: `profile add` and `profile set` commands were updated to use the unified `credentials-file` option for setting credential file paths for Bedrock profiles.
- **Implementation Details**:
    - `internal/config/config.go`: `AWSCredentialsFile` removed, `CredentialsFile` comment updated. `ResolvePath` function added for path expansion and resolution.
    - `internal/llm/bedrock_nova.go`: `newBedrockClient` now uses `profile.CredentialsFile` and calls `loadAWSCredentialsFromFile` (which utilizes `config.ResolvePath`) to load AWS credentials from the specified JSON file.
    - `cmd/set.go` & `cmd/add.go`: Removed specific references to `aws_credentials_file` and ensured `credentials-file` is used for Bedrock.
- **Outcome**: The application now offers a more secure and flexible way to manage Bedrock credentials, aligning with best practices for sensitive data handling. The improved transparency in `profile show` enhances user trust and understanding of file operations.

## 2025-08-03 (Enhancements: DoS Protection, Configuration Consistency, and Profile Check Command)

- **Objective**: To address critical issues related to DoS protection and configuration handling, and to introduce a new utility for profile management.
- **Key Issues Addressed**:
    1.  **Incorrect Standard Input Handling for System Prompts**: Previously, system prompts could incorrectly read from standard input, leading to unexpected behavior. This was fixed by refactoring `loadPrompt` into `loadUserPrompt` and `loadSystemPrompt` in `cmd/prompt.go`, ensuring system prompts never consume stdin.
    2.  **Memory Safety Vulnerability**: The application would load entire files into memory before checking size limits, posing a DoS risk. `readAndProcessStream` in `cmd/prompt.go` was modified to stop reading input once `MaxPromptSizeBytes` is reached, even in "warn" mode, preventing excessive memory consumption.
    3.  **Lack of UTF-8 Safety**: String truncation for size limiting was not UTF-8 aware, potentially corrupting multi-byte characters. `truncateStringByBytes` in `cmd/prompt.go` was updated to correctly handle UTF-8 characters during truncation.
    4.  **Configuration Backward Compatibility**: Older configuration files might lack the `Limits` section, leading to inconsistent behavior. `internal/config/config.go` was modified to ensure the `Limits` struct is always initialized with default values when loading configurations, guaranteeing consistent behavior across all profiles.
- **New Feature: `llm-cli profile check` Command**:
    - Introduced a new subcommand under `profile` to verify and migrate configuration profiles.
    - It inspects all profiles and prompts the user to update `limits` settings that are at their default zero values (indicating they might be from an older version or not explicitly set).
    - Includes a `--confirm` (`-y`) flag for non-interactive operation.
    - Before saving any changes, it creates a timestamped backup of the `config.json` file in `~/.config/llm-cli/backups/`, enhancing data safety.
    - The `profile show` command was also enhanced to display `limits` information.
- **Outcome**: The application is now more robust against DoS attacks, provides better backward compatibility for configurations, and offers a new tool for users to manage their profiles effectively. All identified issues from `DEVELOPMENT_LOG.md` related to DoS protection and configuration handling have been addressed.
