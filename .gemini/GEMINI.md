# Gemini Project Guidelines

This document outlines the development and build rules specific to this project, intended for the Gemini agent.

For a detailed history of development and key decisions, please refer to the [Development Log](DEVELOPMENT_LOG.md).
For future development plans and roadmap, please refer to the [Development Plan](DEVELOPMENT_PLAN.md).

## Development Rules

### Code Style and Formatting

- Adhere to standard Go formatting (`gofmt`).
- Follow idiomatic Go practices.
- Keep functions concise and focused on a single responsibility.

### Concurrency Best Practices
- When implementing asynchronous processes or communication between goroutines, prioritize standard Go concurrency patterns like `sync.WaitGroup` and buffered channels.
- Avoid overly complex `select` controls for state management to prevent race conditions and ensure robust error handling.

### Testing Principles

- Write unit tests for new features and bug fixes.
- For critical bug fixes, especially those related to core logic like API interaction or concurrency, add a regression test to prevent recurrence.
- Ensure tests cover critical paths and edge cases.
- Use `make test` to run tests.

### Linting
- Use `golangci-lint` for static code analysis.
- Ensure all code passes lint checks before committing.

### Commit Message Conventions

- Use the Conventional Commits specification (e.g., `feat:`, `fix:`, `refactor:`, `docs:`).
- For multi-line commit messages, write the message in a temporary file (e.g., `.git/COMMIT_MSG`) and use `git commit -F <file>` to avoid shell interpretation errors. This is the standard procedure.
- Explain *why* a change was made, not just *what* was changed.

### Bug Fixes
- Identify and address the root cause of bugs, avoiding temporary or superficial fixes.

### Documentation Principles

- **Language Policy**: All documentation will be written in Japanese first (as the primary source of truth) and then translated into English.
  - The English version should include a note indicating it is a translation and that the Japanese version takes precedence in case of discrepancies.
- **Scope**: Maintain both user-facing documents (e.g., `README`) and developer-facing documents (e.g., `DEVELOPING_PROVIDERS.md`).
- **Maintenance**: When a feature is changed or added, ensure all relevant documentation is updated accordingly.

### File and Directory Operations
- When deleting or modifying files/directories, always use absolute paths instead of relative paths to prevent unintended operations.

## Build Rules

### Makefile Usage

- Always use `make` commands for building, testing, and cleaning the project.

### Build Commands and Outputs

- `make build`: Builds a binary for the current OS and architecture. Output: `bin/<OS>-<ARCH>/llm-cli`.
- `make cross-compile`: Builds binaries for multiple OS/architectures and creates compressed archives. Output: `bin/llm-cli-<platform>.tar.gz` or `.zip`.
- `make all`: Executes both `make build` and `make cross-compile`, generating all binaries and archives.
- `make test`: Runs all project tests.
- `make clean`: Removes build artifacts and caches.

### Dependency Management

- Use Go Modules for dependency management.
- Run `go mod tidy` after adding or removing dependencies to ensure `go.mod` and `go.sum` are up-to-date.