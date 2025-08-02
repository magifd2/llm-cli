# llm-cli

`llm-cli` is a command-line interface tool for interacting directly with local and remote Large Language Models (LLMs). It provides a unified way to send prompts and manage configurations for various providers like Ollama, LM Studio, and Amazon Bedrock.

## Features

*   **Multi-Provider Support**: Works seamlessly with Ollama, LM Studio (and other OpenAI-compatible APIs), and Amazon Bedrock.
*   **Profile Management**: Save multiple LLM configurations (endpoints, models, API keys) as profiles and easily switch between them.
*   **Flexible Input**: Pass prompts via command-line arguments, files, or standard input (pipes).
*   **Streaming Display**: Display responses from the LLM in real-time using the `--stream` flag.
*   **Single Binary**: Operates as a single executable file (excluding configuration files), making it easy to distribute and use.

## Installation

1.  **Download the Binary**: Go to the [Releases page](https://github.com/magifd2/llm-cli/releases) of the project repository.
2.  Download the appropriate binary for your operating system and architecture.
3.  **Place it in your PATH**: Move the downloaded executable to a directory included in your system's `PATH` (e.g., `/usr/local/bin` on macOS/Linux or a custom directory on Windows).
4.  **Make it Executable**: On macOS and Linux, you may need to grant execute permissions:
    ```bash
    chmod +x /path/to/your/llm-cli
    ```

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

# Switch to the newly created profile
llm-cli profile use lmstudio
```

You can now send prompts to your LM Studio model.

#### 3. Amazon Bedrock

To use Amazon Bedrock, you need valid AWS credentials and a specified region.

**Credential Precedence:**
1.  Credentials set directly in the `llm-cli` profile (`aws_access_key_id`, `aws_secret_access_key`).
2.  Standard AWS SDK credential chain (e.g., environment variables, shared credentials file, IAM roles).

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

# Switch to the Bedrock profile
llm-cli profile use bedrock-nova
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

**Note:** For `credentials-file`, you can specify the path to your service account key JSON file using `~` (tilde) or an absolute path. The `~` will be expanded to your home directory at runtime.

**Required IAM Roles:**
Your service account needs permissions to invoke Vertex AI models.
*   `Vertex AI User` role

## Command Reference

### `llm-cli prompt`

Sends a prompt to the currently active LLM.

| Flag                 | Shorthand | Description                                                 |
| -------------------- | --------- | ----------------------------------------------------------- |
| `--user-prompt`      | `-p`      | The main prompt text to send to the model.                  |
| `--user-prompt-file` | `-f`      | Path to a file containing the user prompt. Use `-` for stdin. |
| `--system-prompt`    | `-P`      | An optional system-level instruction for the model.         |
| `--system-prompt-file`| `-F`      | Path to a file containing the system prompt.                |
| `--stream`           |           | Whether to display the response as a real-time stream.      |

*If no prompt flag is provided, the first positional argument is used as the prompt. If that is also missing, input is read from stdin.*

### `llm-cli profile`

Manages configuration profiles.

| Subcommand | Description                                                        |
| ---------- | ------------------------------------------------------------------ |
| `list`     | Shows all available profiles and indicates the active one.         |
| `use`      | Switches the active profile. `llm-cli profile use <profile-name>`  |
| `add`      | Creates a new profile. If no parameters are specified, it copies settings from the default profile. `llm-cli profile add <new-name> [--provider <provider>] [--model <model>] [...]` |
| `set`      | Modifies a key in the current profile. `llm-cli profile set <key> <value>` |
| `remove`   | Deletes a profile. `llm-cli profile remove <profile-name>`         |
| `edit`     | Opens the `config.json` file in your default text editor.          |

## Contributing & Development

Contributions, such as adding new features or fixing bugs, are welcome.

If you are interested in adding a new LLM provider, please see the [Provider Development Guide](./DEVELOPING_PROVIDERS.en.md).

## Acknowledgements

This project was developed with Google's AI assistant "Gemini" as a coding partner.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
