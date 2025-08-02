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

- **Language Policy**: The primary documentation will be in English (e.g., `README.md`, `BUILD.md`, `CHANGELOG.md`). Other languages (e.g., Japanese) will be provided as auxiliary documentation with a language suffix (e.g., `README.ja.md`).
- **Scope**: Maintain both user-facing documents (e.g., `README.md`) and developer-facing documents (e.g., `CONTRIBUTING.md`).
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

## Additional Development Principles: Self-Discipline for Robustness and Maintainability
Learning from past failures and to prevent future breakdowns, the following principles must be strictly adhered to in all development and refactoring work.

1. Testability First Principle
* Always prioritize practical "ease of testing" over theoretical "beauty."
* Even if a design pattern (e.g., self-registration in init()) appears clean and extensible, its adoption is prohibited in principle if it significantly complicates unit testing.
* Always ask "How do I test this code?" and do not write code for which you cannot answer. Easily testable code is inherently loosely coupled and easy to understand.

2. Implicit is Dangerous Principle
* Prohibit "implicit" or "magical" implementations such as changes to global state in init() functions or behavior changes simply by importing a package.
* Dependency Injection must always be done explicitly, for example, through function arguments. This makes code behavior traceable and facilitates mocking in tests.

3. Redefinition and Compliance Obligation for "Major Refactoring"
Any change that falls under even one of the following categories is considered a "Major Refactoring" and requires presenting a detailed development plan, including scope of impact and testing plan, and obtaining approval before starting work.
* Changes to Initialization Logic: Changes related to init() functions, global variables, or package initialization order.
* Changes to Core Interfaces: Changes to the signature of central interfaces in the system, such as `Provider`.
* Widespread Impact: If changes require modifications across three or more packages to be completed.
* Introduction of Cross-Cutting Concerns: Adding or changing features that affect multiple components, such as authentication, logging, or caching.

4. Safe Refactoring Protocol
Even if a change does not qualify as a "Major Refactoring," the following steps must be strictly followed when modifying core logic. This is an essential procedure to prevent breakdowns.
1. [Establish Baseline]: Commit a stable state where all tests pass immediately before starting work. Commit message: `refactor: Start refactoring X`
2. [Minimum Unit Change]: Make code changes in the smallest possible units (e.g., extract one function, add one variable).
3. [Immediate Testing]: Immediately run `make test` after making a change.
   * If successful: Commit the change immediately (`feat: Introduce Y` / `refactor: Extract Z`). Then return to step 2.
   * If failed: Discard all changes immediately (`git reset --hard`). Do not modify other code to make the test pass. Recognize that test failures are a danger signal indicating "design is wrong" and fundamentally rethink the approach.
4. [Completion]: Achieve the final goal by accumulating small commits.