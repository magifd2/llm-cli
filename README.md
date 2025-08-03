# llm-cli

`llm-cli` is a command-line interface tool for interacting directly with local and remote Large Language Models (LLMs). It provides a unified way to send prompts and manage configurations for various providers like Ollama, LM Studio, and Amazon Bedrock.

## Features

*   **Multi-Provider Support**: Works seamlessly with Ollama, LM Studio (and other OpenAI-compatible APIs), Amazon Bedrock, and Google Cloud Vertex AI.
*   **Profile Management**: Save multiple LLM configurations (endpoints, models, API keys) as profiles and easily switch between them.
*   **Flexible Input**: Pass prompts via command-line arguments, files, or standard input (pipes).
*   **Streaming Display**: Display responses from the LLM in real-time using the `--stream` flag.
*   **Single Binary**: Operates as a single executable file (excluding configuration files), making it easy to distribute and use.

## Installation

`llm-cli` can be easily installed using the provided `Makefile`.

### Using `make install`

This method builds the `llm-cli` binary and installs it to a specified directory, along with the Zsh shell completion script.

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
    To install `llm-cli` to a custom directory (e.g., `/opt/llm-cli/bin`):
    ```bash
    sudo make install PREFIX=/opt/llm-cli
    ```

After installation, for Zsh users, you might need to run `compinit` or restart your shell for the completion script to take effect.

### Uninstallation

To uninstall `llm-cli` and its completion script, use `make uninstall` with the same `PREFIX` used during installation.

*   **Default Uninstallation:**
    ```bash
    sudo make uninstall
    ```

*   **User-local Uninstallation:**
    ```bash
    make uninstall PREFIX=~
    ```

*   **Custom Directory Uninstallation:**
    ```bash
    sudo make uninstall PREFIX=/opt/llm-cli
    ```

**Note:** The uninstallation process does NOT remove your configuration files located at `~/.config/llm-cli/config.json`.

**Note on Path Notation:** The `~` (tilde) character is a common shorthand for the user's home directory. `llm-cli` correctly expands this to the appropriate home directory path on all supported operating systems (Linux, macOS, Windows), ensuring cross-platform compatibility.

## Quick Start

Once installed and configured, you can immediately start interacting with your LLM.

```bash
# Send a simple prompt to the default LLM
llm-cli prompt "What is the distance between Earth and the Moon?"

# Get a streaming response
llm-cli prompt "Tell me a short story about a robot who discovers music." --stream
```

## Configuration

`llm-cli` manages all its settings in a single configuration file located at `~/.config/llm-cli/config.json`. While you can edit this file directly with `llm-cli profile edit`, it is recommended to use the `profile` subcommands.

### Provider-Specific Setup

#### 1. Ollama

If you are running Ollama on its default address (`http://localhost:11434`), `llm-cli` works out of the box. The default profile is pre-configured for this setup.

To use a specific model you have pulled with Ollama:
```bash
# Switch to the default profile (if not already active)
llm-cli profile use default

# Set the model you want to use
llm-cli profile set model "llama3" 
```

#### 2. LM Studio (or other OpenAI-compatible APIs)

To use LM Studio, you first need to start its local server.

1.  **Start the Server**: In LM Studio, go to the "Local Server" tab (the `<->` icon).
2.  **Load a Model**: Select a model to load and wait for it to be ready.
3.  **Start Server**: Click the "Start Server" button. Note the server URL displayed at the top (e.g., `http://localhost:1234/v1`).

Now, configure `llm-cli` to use this server:

```bash
# Add a new profile for LM Studio
llm-cli profile add lmstudio

# Set the provider to "openai"
llm-cli profile set provider openai

# Set the endpoint to the URL from LM Studio
llm-cli profile set endpoint "http://localhost:1234/v1"

# The model name is often arbitrary for local servers, but it must be set.
# You can typically use the model identifier from LM Studio.
llm-cli profile set model "gemma-2-9b-it"

# (Optional) Set API key directly in the profile
# llm-cli profile set api_key "YOUR_API_KEY"

# (Optional) Use a credentials file for OpenAI API key (e.g., ~/.openai/api_key.json)
# llm-cli profile set credentials-file "~/path/to/your/openai-api-key.json"

# Switch to the newly created profile
llm-cli profile use lmstudio
```

**Note:** For `credentials-file`, the JSON file should contain the API key under the `openai_api_key` field, like this example:

```json
{
  "openai_api_key": "sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
}
```

You can now send prompts to your LM Studio model.

#### 3. Amazon Bedrock

To use Amazon Bedrock, you need valid AWS credentials and a specified region.

**Credential Precedence:**
1.  Credentials loaded from a specified `credentials-file` in the `llm-cli` profile.
2.  Credentials set directly in the `llm-cli` profile (`aws_access_key_id`, `aws_secret_access_key`).
3.  Standard AWS SDK credential chain (e.g., environment variables, shared credentials file, IAM roles).

**Configuration Steps:**

```bash
# Add a new profile for Bedrock
llm-cli profile add bedrock-nova

# Set the provider to "bedrock"
llm-cli profile set provider bedrock

# Set the model ID for the model you want to use
llm-cli profile set model "amazon.nova-lite-v1:0"

# Set the AWS region where you will invoke the model
llm-cli profile set aws_region "us-east-1"

# (Optional) Set credentials directly if needed
# llm-cli profile set aws_access_key_id "YOUR_KEY_ID"
# llm-cli profile set aws_secret_access_key "YOUR_SECRET_KEY"

# (Optional) Use a credentials file for Bedrock (e.g., ~/.aws/credentials.json)
# llm-cli profile set credentials-file "~/path/to/your/aws-credentials.json"

# Switch to the Bedrock profile
llm-cli profile use bedrock-nova
```

**Note:** For `credentials-file`, you can specify the path to your AWS credentials JSON file using `~` (tilde) or an absolute path. The `~` will be expanded to your home directory at runtime. The JSON file should contain `aws_access_key_id` and `aws_secret_access_key` fields, like this example:

```json
{
  "aws_access_key_id": "YOUR_AWS_ACCESS_KEY_ID",
  "aws_secret_access_key": "YOUR_AWS_SECRET_ACCESS_KEY"
}
```

**Required IAM Policies:**
Your AWS identity must have permissions to invoke Bedrock models.

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "bedrock:InvokeModel",
                "bedrock:InvokeModelWithResponseStream"
            ],
            "Resource": "arn:aws:bedrock:us-east-1::foundation-model/amazon.nova-lite-v1:0"
        }
    ]
}
```
*Note: As a best practice, restrict the `Resource` to the specific models you need.*

#### 4. Google Cloud Vertex AI

To use Google Cloud Vertex AI, you need to set up a GCP project and prepare your credentials.

**Prerequisites:**
1.  Ensure you have a Google Cloud Platform (GCP) project created for using Vertex AI.
2.  Enable the **Vertex AI API** in your target GCP project.
3.  Create a service account key and download it in **JSON** format. Store this key file securely.
    *   Grant the service account the **"Vertex AI User"** role.

**Configuration Steps:**

```bash
# Add a new profile for Vertex AI (one-shot configuration)
llm-cli profile add my-vertex-ai \
  --provider vertexai \
  --model gemini-1.5-pro-001 \
  --project-id "your-gcp-project-id" \
  --location "us-central1" \
  --credentials-file "~/path/to/your/service-account-key.json"

# Switch to the newly created profile
llm-cli profile use my-vertex-ai
```

**Note:** For `credentials-file`, you can specify the path to your service account key JSON file using `~` (tilde) or an absolute path. The `~` will be expanded to your home directory at runtime. This field is now also used for AWS Bedrock credentials files.

**Required IAM Roles:**
Your service account needs permissions to invoke Vertex AI models.
*   `Vertex AI User` role

**System Prompt Handling:**
Vertex AI's GenAI SDK does not directly support system prompts. Therefore, `llm-cli` simulates system prompt behavior by sending the system prompt content as the first message in the chat, followed by the user prompt content.

### Size and Usage Limits (DoS Protection)

To prevent accidental excessive usage or potential misuse that could lead to high costs or system instability, `llm-cli` includes a configurable limiting mechanism. These settings are managed within a `limits` object inside each profile.

By default, these limits are enabled for new profiles.

For users upgrading from older versions, or to ensure all profiles have the latest default limit settings, you can use the `profile check` command:
```bash
llm-cli profile check
```
This command will inspect all your profiles and prompt you to update any `limits` settings that are at their default zero values (indicating they might be from an older version or not explicitly set). It also creates a timestamped backup of your `config.json` before making any changes.

```json
"my-profile": {
    "provider": "openai",
    "model": "gpt-4",
    "limits": {
        "enabled": true,
        "on_input_exceeded": "stop",
        "on_output_exceeded": "stop",
        "max_prompt_size_bytes": 10485760,
        "max_response_size_bytes": 20971520
    }
}
```

*   `enabled`: A boolean (`true` or `false`) to turn limits on or off for the profile.
*   `on_input_exceeded`: Determines the action when the prompt size exceeds the limit.
    *   `"stop"` (default): The command will fail with an error message.
    *   `"warn"`: The command will truncate the prompt, show a warning, and proceed.
*   `on_output_exceeded`: Determines the action when the response size exceeds the limit.
    *   `"stop"` (default): The command will fail (or stop streaming) with an error message.
    *   `"warn"`: The command will truncate the response, show a warning, and exit successfully.
*   `max_prompt_size_bytes`: The maximum allowed size of the combined user and system prompts in bytes. (Default: `10485760` / 10 MB)
*   `max_response_size_bytes`: The maximum allowed size of the response from the LLM in bytes. (Default: `20971520` / 20 MB)

These values can be configured using the `llm-cli profile set` and `llm-cli profile add` commands.

## Command Reference


### `llm-cli prompt`

Sends a prompt to the currently active LLM.

| Flag                      | Shorthand | Description                                                                 |
| ------------------------- | --------- | --------------------------------------------------------------------------- |
| `--user-prompt`           | `-p`      | The main prompt text to send to the model.                                  |
| `--user-prompt-file`      | `-f`      | Path to a file containing the user prompt. Use `-` for stdin.                 |
| `--system-prompt`         | `-P`      | An optional system-level instruction for the model.                         |
| `--system-prompt-file`    | `-F`      | Path to a file containing the system prompt.                                |
| `--stream`                |           | Whether to display the response as a real-time stream.                      |
| `--profile`               |           | Use a specific profile for this command (overrides current active profile). |
| `--on-input-exceeded`     |           | Override profile setting for input limit. (Accepts: `stop`, `warn`)         |
| `--on-output-exceeded`    |           | Override profile setting for output limit. (Accepts: `stop`, `warn`)        |

*If no prompt flag is provided, the first positional argument is used as the prompt. If that is also missing, input is read from stdin.*

### `llm-cli profile`

Manages configuration profiles.

| Subcommand | Description                                                                                             |
| ---------- | ------------------------------------------------------------------------------------------------------- |
| `list`     | Shows all available profiles, their primary settings, and limit configurations.                         |
| `use`      | Switches the active profile. `llm-cli profile use <profile-name>`                                       |
| `add`      | Creates a new profile. If no parameters are specified, it copies settings from the default profile.       |
|            | **Options:**                                                                                            |
|            | `--provider <provider>`: LLM provider (e.g., ollama, openai, bedrock, vertexai)                         |
|            | `--model <model>`: Model name (e.g., llama3, gpt-4, gemini-1.5-pro-001)                                 |
|            | `--endpoint <url>`: API endpoint URL                                                                    |
|            | `--api-key <key>`: API key for the provider                                                             |
|            | `--aws-region <region>`: AWS region for Bedrock                                                         |
|            | `--aws-access-key-id <id>`: AWS Access Key ID for Bedrock                                               |
|            | `--aws-secret-access-key <key>`: AWS Secret Access Key for Bedrock                                      |
|            | `--project-id <id>`: GCP Project ID for Vertex AI                                                       |
|            | `--location <location>`: GCP Location for Vertex AI                                                     |
|            | `--credentials-file <path>`: Path to a credentials file (for GCP service account, AWS Bedrock, or OpenAI API Key).       |
|            | `--limits-enabled <bool>`: Enable or disable limits for this profile. (Default: `true`)                 |
|            | `--limits-on-input-exceeded <action>`: Action for input limit: `stop` or `warn`. (Default: `stop`)       |
|            | `--limits-on-output-exceeded <action>`: Action for output limit: `stop` or `warn`. (Default: `stop`)      |
|            | `--limits-max-prompt-size-bytes <bytes>`: Max prompt size in bytes. (Default: `10485760`)                |
|            | `--limits-max-response-size-bytes <bytes>`: Max response size in bytes. (Default: `20971520`)             |
| `set`      | Modifies a key in the current profile. `llm-cli profile set <key> <value>`. See available keys below.     |
|            | **Available Keys:** `provider`, `model`, `endpoint`, `api_key`, `aws_region`, `aws_access_key_id`, `aws_secret_access_key`, `project_id`, `location`, `credentials_file`, `limits.enabled`, `limits.on_input_exceeded`, `limits.on_output_exceeded`, `limits.max_prompt_size_bytes`, `limits.max_response_size_bytes` |
| `remove`   | Deletes a profile. `llm-cli profile remove <profile-name>`                                              |
| `show`     | Shows all details of a specific profile, including limits. `llm-cli profile show [profile-name]`        |
| `edit`     | Opens the `config.json` file in your default text editor for manual changes.                            |
| `check`    | Checks and migrates configuration profiles, offering to update default settings.                        |

## Contributing & Development

Contributions, such as adding new features or fixing bugs, are welcome.

**Note for macOS Developers:** When building on macOS, `make build` will automatically produce a universal binary (supporting both `amd64` and `arm64` architectures) to ensure broader compatibility.

If you are interested in contributing, please see the [Contributing Guide](./CONTRIBUTING.md).

## Acknowledgements

This project was developed with Google's AI assistant "Gemini" as a coding partner.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.

## Security

For information on how to report security vulnerabilities, please refer to our [Security Policy](./SECURITY.md).
