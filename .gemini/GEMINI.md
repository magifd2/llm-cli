# Gemini Project Guidelines

This document outlines the development and build rules specific to this project, intended for the Gemini agent.

For a detailed history of development and key decisions, please refer to the [Development Log](DEVELOPMENT_LOG.md).
For future development plans and roadmap, please refer to the [Development Plan](DEVELOPMENT_PLAN.md).

## Development Rules

### Security First Principle (セキュリティ第一原則)

- **Security is the highest priority, overriding all other considerations such as functionality or performance. (セキュリティは、機能性やパフォーマンスといった他のいかなる考慮事項よりも優先される、絶対的な最優先事項である。)**
- **All code, dependencies, and configurations must be reviewed for potential security vulnerabilities before being committed. (全てのコード、依存関係、設定は、コミットされる前に、潜在的なセキュリティ脆弱性についてレビューされなければならない。)**
- **Never trust user input, including environment variables. All external inputs must be validated and sanitized to prevent injection attacks. (環境変数を含む、いかなるユーザー入力も信頼してはならない。全ての外部入力は、インジェクション攻撃を防ぐために検証され、無害化されなければならない。)**
- **Sensitive information (API keys, credentials) must never be hardcoded or stored in insecure locations. (APIキーや認証情報などの機密情報は、決してハードコードされたり、安全でない場所に保存されたりしてはならない。)**

### Secure Development Lifecycle (セキュア開発ライフサイクル)

- **Threat Modeling at Design Phase (設計段階での脅威モデリング):** Before implementing a new feature, consider potential threats. For example, when adding a feature that interacts with the filesystem, evaluate risks like path traversal.
- **Security-Focused Code Reviews (セキュリティを重視したコードレビュー):** All code reviews must include a specific check for security vulnerabilities. Do not approve pull requests that have not been reviewed from a security perspective.
- **Safe Testing Practices (安全なテストの実施):** When testing for vulnerabilities, use harmless proof-of-concept payloads. Before running tests that involve external inputs like environment variables, always inspect their contents first.
- **Dependency Scanning (依存関係のスキャン):** Regularly scan project dependencies for known vulnerabilities using tools like `govulncheck`.

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