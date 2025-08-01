# Development Plan

This document outlines the future development roadmap and planned feature enhancements for the `llm-cli` project.

## Short-Term Goals

### 1. Amazon Bedrock - Support for Additional Models

- **Objective**: Extend `llm-cli` to support other foundational models available through Amazon Bedrock, beyond the currently implemented Nova models.
- **Target Models**: Prioritize widely used models such as Anthropic Claude (e.g., `anthropic.claude-v2`, `anthropic.claude-3-sonnet-20240229-v1:0`).
- **Implementation Strategy**: Following the successful refactoring for Nova models, implement new provider files (e.g., `bedrock_claude.go`) for each model family, adhering to their specific API request/response formats and streaming protocols.

### 2. Google Cloud Platform (GCP) Vertex AI Generative AI - Re-attempt Integration

- **Objective**: Integrate `llm-cli` with GCP Vertex AI's Generative AI services.
- **Lessons Learned from Previous Attempt**: Analyze the reasons for the previous failure (e.g., authentication issues, API format mismatches, SDK usage). Focus on robust error handling and clear debugging outputs.
- **Target Models**: Initially focus on text-based models (e.g., `gemini-pro`).
- **Implementation Strategy**: 
  - Research and understand the latest Vertex AI API documentation for generative models.
  - Implement a dedicated provider file (e.g., `gcp_vertex_ai.go`) within `internal/llm/`.
  - Ensure proper authentication mechanisms (e.g., service accounts, ADC) are supported.
  - Implement both buffered and streaming chat functionalities.

### 3. Enhanced LLM Provider Testing Strategy (テスト戦略の強化)

- **Objective**: Implement robust and reliable automated tests for LLM providers (`internal/llm` package) to ensure correctness of request/response handling, especially for streaming APIs. This aims to prevent regressions and improve development efficiency.
- **Motivation**: Previous attempts to test streaming providers encountered significant challenges, including deadlocks and complex interactions between `context.Context` and blocking I/O operations. This highlighted the need for a more controlled testing environment.
- **Strategy**: Instead of modifying the core provider code for testability, we will leverage `httptest.Server` to create local mock HTTP servers. For streaming tests, `io.Pipe` will be used to precisely control data flow and simulate connection closures, allowing for reliable testing of context cancellation without modifying the production code.
- **Implementation Approach (Block-1 vs Block-2 Providers):**
  - **Block-1 Providers**: Existing provider implementations (`ollama.go`, `openai.go`, `bedrock_nova.go`) will remain as they are, representing the current stable version.
  - **Block-2 Providers**: New provider implementations (e.g., `ollama_block2.go`) will be created. These will be identical in functionality to Block-1 but will be designed to utilize the new, testable I/O module (if necessary) or simply serve as the target for the new testing methodology.
  - **Test-Specific Provider**: A dedicated `test_provider.go` will be created to validate the testing framework itself, ensuring `httptest.Server` and `io.Pipe` interactions are correctly handled before applying the pattern to real providers.
  - **Build-Time Switching**: The `cmd/prompt.go` logic will be updated to allow switching between Block-1 and Block-2 providers at build time (e.g., using Go build tags), enabling testing of the new implementations without affecting the default build.

## Mid-to-Long-Term Goals

- (To be defined based on user feedback, market trends, and project priorities.)