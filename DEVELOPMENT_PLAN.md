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

## Mid-to-Long-Term Goals

- (To be defined based on user feedback, market trends, and project priorities.)
