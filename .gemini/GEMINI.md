# Gemini Project Guidelines

This document outlines the development and build rules specific to this project, intended for the Gemini agent.

For a detailed history of development and key decisions, please refer to the [Development Log](DEVELOPMENT_LOG.md).
For future development plans and roadmap, please refer to the [Development Plan](DEVELOPMENT_PLAN.md).

## Development Rules

### Code Style and Formatting

- Adhere to standard Go formatting (`gofmt`).
- Follow idiomatic Go practices.
- Keep functions concise and focused on a single responsibility.

### Testing Principles

- Write unit tests for new features and bug fixes.
- Ensure tests cover critical paths and edge cases.
- Use `make test` to run tests.

### Linting
- Use `golangci-lint` for static code analysis.
- Ensure all code passes lint checks before committing.

### Rollback Strategy
- Before making significant changes, commit the current state to allow for easy rollback if necessary.

### Commit Message Conventions

- Use Conventional Commits specification (e.g., `feat:`, `fix:`, `refactor:`, `docs:`).
- Keep commit messages concise and descriptive.
- Explain *why* a change was made, not just *what* was changed.

### Bug Fixes
- Identify and address the root cause of bugs, avoiding temporary or superficial fixes.

### Code Quality and Security
- Prioritize robustness, security, and maintainability in design and implementation.
- Adopt a "security-first" approach as a fundamental policy in all aspects of development.

### Documentation Principles

- All user-facing documentation (README, BUILD, CHANGELOG) must be maintained in both English (`.en.md`) and Japanese (`.ja.md`).
- The root `README.md` serves as a language selection page.
- Ensure consistency between language versions.
- When a feature is changed or added, ensure all relevant documentation is updated accordingly.

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
