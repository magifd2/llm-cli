# llm-cli

`llm-cli` is a CLI tool for interacting directly with local (Ollama, LM Studio) or remote LLMs (such as OpenAI in the future) from the command line.

## Features

*   **Multi-Provider Support**: Compatible with Ollama and LM Studio (OpenAI-compatible API).
*   **Profile Management**: Save multiple LLM configurations (endpoints, models, etc.) as profiles and switch between them easily.
*   **Flexible Input**: Pass prompts via command-line arguments, files, or standard input (pipes).
*   **Streaming Display**: Display responses from the LLM in real-time.
*   **Single Binary with Go**: Operates as a single executable file (excluding configuration files), making it easy to distribute.

## Usage

### Sending Prompts (Required)

```bash
# Simple prompt
llm-cli ask --prompt "What is the capital of Japan?" # --prompt or --prompt-file is required

# With system prompt
llm-cli ask --prompt "Introduce yourself" --system-prompt "You are a cat. Speak with 'nyan' at the end of your sentences."

# Streaming display
llm-cli ask --prompt "Count from 1 to 100" --stream

# Read prompt from a file (or pipe from standard input)
llm-cli ask --prompt-file ./my_prompt.txt

# Pass via pipe
echo "Summarize this text" | llm-cli ask
```

### Managing Profiles

```bash
# List profiles
llm-cli profile list

# Add a new profile (created by copying the default profile)
llm-cli profile add my-new-profile

# Switch the active profile
llm-cli profile use my-new-profile

# Change settings for the current profile
llm-cli profile set model "new-model-name"
llm-cli profile set endpoint "http://my-endpoint/v1"

# Delete a profile
llm-cli profile remove my-new-profile

# Edit the configuration file directly
llm-cli profile edit
```

### Amazon Bedrock Configuration

To use Amazon Bedrock, you need AWS credentials and region settings.
Credentials can be set directly in the profile or by utilizing the AWS SDK's default credential provider chain (environment variables, IAM roles, etc.).

**Example Bedrock Profile:**

```bash
# Add a new Bedrock profile
llm-cli profile add bedrock-claude

# Set provider to bedrock (Nova models use the Messages API)
llm-cli profile set provider bedrock

# Set model ID (e.g., Amazon Nova Lite v1)
llm-cli profile set model amazon.nova-lite-v1:0

# Set AWS region (e.g., ap-northeast-1)
llm-cli profile set aws_region ap-northeast-1

# To set access key ID and secret access key directly (deprecated: environment variables or IAM roles are recommended)
llm-cli profile set aws_access_key_id YOUR_AWS_ACCESS_KEY_ID
llm-cli profile set aws_secret_access_key YOUR_AWS_SECRET_ACCESS_KEY

# After configuration, switch to this profile
llm-cli profile use bedrock-claude # or llm-cli profile use bedrock
```

**Credential Precedence:**

1.  `aws_access_key_id` and `aws_secret_access_key` directly configured in the `llm-cli` profile.
2.  AWS SDK's default credential provider chain (environment variables `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_REGION`, IAM roles, etc.).

#### Required IAM Policies

To invoke Amazon Bedrock models, your AWS credentials must have the appropriate IAM policies attached. The minimum required actions are `bedrock:InvokeModel` and `bedrock:InvokeModelWithResponseStream`.

**Example of Minimum IAM Policy:**

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
            "Resource": "arn:aws:bedrock:ap-northeast-1::foundation-model/amazon.nova*"
        }
    ]
}
```

**Note**: Replace `<your-aws-region>` and `<your-model-id>` with the actual region and model ID you intend to use. As a security best practice, it is strongly recommended to restrict the `Resource` to the most specific model possible. If you need to use multiple models, you can use wildcards like `"Resource": "arn:aws:bedrock:<your-aws-region>::/foundation-model/*"`, but be aware that this broadens access.

## Configuration

Settings are saved in `~/.config/llm-cli/config.json`. While they can be managed with the `profile` commands, you can also edit them directly using `profile edit`.

**Security Note**: Sensitive information such as API keys is stored in plain text in the configuration file. You are responsible for managing access to this file.

## Acknowledgements

This project was developed with Google's AI assistant "Gemini" as a coding partner.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.