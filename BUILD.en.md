# Build Instructions

This document describes how to build `llm-cli` from source code.

## Prerequisites

*   [Go](https://go.dev/doc/install) (version 1.21 or later recommended)
*   [Git](https://git-scm.com/)
*   `make` command
    *   Standard on macOS/Linux.
    *   For Windows, please install [Make for Windows](http://gnuwin32.sourceforge.net/packages/make.htm) or similar.

## Build

This project recommends building using `Makefile`.

### 1. Clone the repository

```bash
git clone https://github.com/magifd2/llm-cli.git
cd llm-cli
```

### 2. Build Commands

The following `make` commands are available. The built binaries will be generated in the `bin/` directory.

*   **`make build`**
    *   Builds only one binary for the currently used OS and architecture. The built binary will be placed in `bin/<OS>-<ARCH>/llm-cli`. This is convenient for testing during development.

*   **`make cross-compile`**
    *   For distribution, builds binaries for multiple OS and architectures at once and creates compressed archives. The artifacts will be generated in the `bin/` directory.
        *   `bin/llm-cli-darwin-universal.tar.gz` (macOS Universal Binary)
        *   `bin/llm-cli-linux-amd64.tar.gz` (Linux amd64)
        *   `bin/llm-cli-windows-amd64.zip` (Windows amd64)

*   **`make all`**
    *   Executes both `make build` and `make cross-compile`. Creates a binary for the current OS and architecture, as well as all cross-compiled binaries and archives.
*   **`make test`**
    *   Runs the project's tests.

*   **`make clean`**
    *   Deletes the `bin/` directory and build cache.