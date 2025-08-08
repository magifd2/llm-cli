# Development Plan

This document outlines the future development roadmap and planned feature enhancements for the `llm-cli` project.

## Core Architectural Principle: Provider Modularity

As of version `v0.0.12`, the project has undergone a significant refactoring to enforce strict modularity for LLM providers. All providers now reside in their own self-contained packages under `internal/llm/`. For example, the Ollama provider is located in `internal/llm/ollama`, and the Vertex AI provider is in `internal/llm/vertexai`.

This architectural change was made to:
- **Enhance Maintainability**: Isolate provider-specific logic to prevent changes in one provider from affecting others.
- **Improve Testability**: Create a clear path toward implementing robust, independent unit tests for each provider, resolving previous challenges with testing the monolithic `internal/llm` package.
- **Ensure Stability**: Guarantee that individual providers can be added, modified, or removed safely without unintended side effects.

All future development must adhere to this principle. New providers must be created in their own dedicated packages.

## Short-Term Goals

### 1. Enhanced LLM Provider Testing Strategy

- **Objective**: With the new modular architecture in place, the next critical step is to implement a robust and reliable automated testing strategy for each LLM provider package.
- **Motivation**: The previous monolithic structure of the `internal/llm` package made unit testing exceptionally difficult, leading to a freeze on test implementation. The new package-per-provider structure removes this roadblock.
- **Strategy**: Leverage `httptest.Server` to create local mock HTTP servers for each provider test suite. This will allow for controlled testing of request/response handling, including streaming and error conditions, without relying on live external APIs.

## Mid-to-Long-Term Goals

### 1. Amazon Bedrock - Support for Additional Models

- **Objective**: Extend `llm-cli` to support other foundational models available through Amazon Bedrock.
- **Target Models**: Prioritize widely used models such as Anthropic Claude (e.g., `anthropic.claude-v2`, `anthropic.claude-3-sonnet-20240229-v1:0`).
- **Implementation Strategy**: Following the modular architecture, new model families will be implemented within the `internal/llm/bedrock` package, potentially as separate files (e.g., `claude.go`) or sub-packages if their API logic differs significantly from the existing Nova implementation.

### 2. Input Validation and Security Hardening (Ongoing)

- **Objective**: Continue to strengthen the application against potential security vulnerabilities.
- **Status**: Core functionalities for input size limiting and UTF-8 validation are complete and stable.
- **Next Steps**: Regularly review and incorporate security best practices. Future provider implementations must also include appropriate validation and security considerations.

---

### Historical/Completed Goals

- **Google Cloud Platform (GCP) Vertex AI Integration**: Completed as of `v0.0.5`. The provider logic is now located in `internal/llm/vertexai` and `internal/llm/vertexai2`.
- **LLM Provider Unit Testing and Code Stability Policy (Superseded)**: The policy of freezing modifications to the `internal/llm` package has been superseded by the new modular architecture, which is designed to enable safe modifications and testing.
