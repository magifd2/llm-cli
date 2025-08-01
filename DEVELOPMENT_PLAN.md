# Development Plan

This document outlines the future development roadmap and planned feature enhancements for the `llm-cli` project.

## Short-Term Goals

### 1. Google Cloud Platform (GCP) Vertex AI Generative AI - Integration (Completed)

- **Objective**: Integrated `llm-cli` with GCP Vertex AI's Generative AI services.
- **Status**: Completed as of `v0.0.5`. The integration involved addressing authentication issues, API format mismatches, and proper SDK usage.
- **Target Models**: Supports text-based models (e.g., `gemini-pro`).
- **Implementation Details**: 
  - Implemented a dedicated provider file (`vertexai.go`) within `internal/llm/`.
  - Ensured proper authentication mechanisms (e.g., service accounts) are supported.
  - Implemented both buffered and streaming chat functionalities.
  - Ensure proper authentication mechanisms (e.g., service accounts, ADC) are supported.
  - Implement both buffered and streaming chat functionalities.

### 2. Enhanced LLM Provider Testing Strategy (テスト戦略の強化)

- **Objective**: Implement robust and reliable automated tests for LLM providers (`internal/llm` package) to ensure correctness of request/response handling, especially for streaming APIs. This aims to prevent regressions and improve development efficiency.
- **Motivation**: Previous attempts to test streaming providers encountered significant challenges, including deadlocks and complex interactions between `context.Context` and blocking I/O operations. This highlighted the need for a more controlled testing environment.
- **Strategy**: Instead of modifying the core provider code for testability, we will leverage `httptest.Server` to create local mock HTTP servers. For streaming tests, `io.Pipe` will be used to precisely control data flow and simulate connection closures, allowing for reliable testing of context cancellation without modifying the production code.
- **Implementation Approach (Block-1 vs Block-2 Providers):**
  - **Block-1 Providers**: Existing provider implementations (`ollama.go`, `openai.go`, `bedrock_nova.go`) will remain as they are, representing the current stable version.
  - **Block-2 Providers**: New provider implementations (e.g., `ollama_block2.go`) will be created. These will be identical in functionality to Block-1 but will be designed to utilize the new, testable I/O module (if necessary) or simply serve as the target for the new testing methodology.
  - **Test-Specific Provider**: A dedicated `test_provider.go` will be created to validate the testing framework itself, ensuring `httptest.Server` and `io.Pipe` interactions are correctly handled before applying the pattern to real providers.
  - **Build-Time Switching**: The `cmd/prompt.go` logic will be updated to allow switching between Block-1 and Block-2 providers at build time (e.g., using Go build tags), enabling testing of the new implementations without affecting the default build.

### 3. LLM Provider Unit Testing and Code Stability Policy (LLMプロバイダーのユニットテストとコード安定性ポリシー)

- **Current Status**: Due to the inherent complexity of testing streaming API interactions and the blocking nature of network I/O in Go, implementing comprehensive unit tests for the `internal/llm` package has proven to be exceptionally challenging. Previous attempts to establish a robust testing framework for these components have resulted in significant development overhead and unresolved issues like deadlocks.
- **Policy**: For the foreseeable future, the implementation of dedicated unit tests for the `internal/llm` package (LLM providers) is **frozen**. This decision is made to prioritize overall project stability and development efficiency.
- **Code Stability Mandate**: Any modifications to the existing, functionally verified code within the `internal/llm` package are **strictly prohibited** unless absolutely critical for security or essential functionality. This measure is put in place to prevent regressions and maintain the current operational stability of the LLM interaction core. Future enhancements or refactoring in this area will require a thoroughly re-evaluated and approved testing strategy.

## Mid-to-Long-Term Goals

### 1. Amazon Bedrock - Support for Additional Models

- **Objective**: Extend `llm-cli` to support other foundational models available through Amazon Bedrock, beyond the currently implemented Nova models.
- **Target Models**: Prioritize widely used models such as Anthropic Claude (e.g., `anthropic.claude-v2`, `anthropic.claude-3-sonnet-20240229-v1:0`).
- **Implementation Strategy**: Following the successful refactoring for Nova models, implement new provider files (e.g., `bedrock_claude.go`) for each model family, adhering to their specific API request/response formats and streaming protocols.

### 2. Input Validation and Security Hardening (Completed / Ongoing)

- **Objective**: Strengthen the application against potential security vulnerabilities related to user input.
- **Status**: Significant progress has been made in this area.
- **Key Initiatives**:
  - **Input Size Limitation**: **[COMPLETED]** Implemented a configurable size limit for input prompts, especially from standard input, to mitigate denial-of-service (DoS) risks. Further refined to prevent excessive memory consumption even in "warn" mode.
  - **UTF-8 Validation**: **[COMPLETED]** Enforced UTF-8 validation and ensured UTF-8 aware string truncation for all incoming prompt data to prevent issues arising from malformed character encodings.
- **Next Steps**: Continue to monitor for new security best practices and potential vulnerabilities.

- (To be defined based on user feedback, market trends, and project priorities.)