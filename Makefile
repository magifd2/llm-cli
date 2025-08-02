# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

# Project details
BINARY_NAME=llm-cli
OUTPUT_DIR=bin
# LDFLAGS for smaller binaries (-s strips symbol table, -w strips DWARF debug info)
LDFLAGS=-ldflags="-s -w"

# Installation paths
PREFIX?=/usr/local
BIN_DIR=$(PREFIX)/bin
COMPLETION_DIR=$(PREFIX)/share/zsh/site-functions # Zsh specific, adjust for others

.PHONY: all build clean test cross-compile install uninstall build-mac-universal build-linux build-windows package-all vulncheck

all: build cross-compile

# Build for the current OS/Arch
build: vulncheck
	@echo "Building for $(shell go env GOOS)/$(shell go env GOARCH)..."
	@mkdir -p $(OUTPUT_DIR)/$(shell go env GOOS)-$(shell go env GOARCH)
	@$(GOBUILD) $(LDFLAGS) -o $(OUTPUT_DIR)/$(shell go env GOOS)-$(shell go env GOARCH)/$(BINARY_NAME) .

# Run tests
test:
	@echo "Running tests..."
	@$(GOTEST) -v ./...

# Run linters
lint:
	@echo "Running linters..."
	@$(GOCMD) run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run ./...

# Run vulnerability check
vulncheck:
	@echo "Running vulnerability check..."
	@$(GOCMD) run golang.org/x/vuln/cmd/govulncheck@latest ./...

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	@$(GOCLEAN)
	@rm -f $(BINARY_NAME)
	@rm -rf $(OUTPUT_DIR)
	@rm -f extract_release_notes.go release_notes.txt

# Install the binary and completion scripts
# Usage: make install (installs to /usr/local/bin)
#        make install PREFIX=~ (installs to ~/bin)
install: build
	@echo "Installing $(BINARY_NAME) to $(BIN_DIR)..."
	@mkdir -p $(BIN_DIR)
	@if [ "$(shell go env GOOS)" = "darwin" ]; then \
		cp $(OUTPUT_DIR)/darwin-universal/$(BINARY_NAME) $(BIN_DIR)/; \
	else \
		cp $(OUTPUT_DIR)/$(shell go env GOOS)-$(shell go env GOARCH)/$(BINARY_NAME) $(BIN_DIR)/; \
	fi
	@echo "Generating and installing Zsh completion script to $(COMPLETION_DIR)..."
	@mkdir -p $(COMPLETION_DIR)
	@$(BIN_DIR)/$(BINARY_NAME) completion zsh > $(COMPLETION_DIR)/_$(BINARY_NAME)
	@echo "Installation complete. Remember to run 'compinit' in Zsh or restart your shell."

# Uninstall the binary and completion scripts
# Note: This does NOT remove configuration files (e.g., ~/.config/llm-cli/config.json).
# Usage: make uninstall (uninstalls from /usr/local/bin)
#        make uninstall PREFIX=~ (uninstalls from ~/bin)
uninstall:
	@echo "Uninstalling $(BINARY_NAME) from $(BIN_DIR)..."
	@rm -f $(BIN_DIR)/$(BINARY_NAME)
	@echo "Removing Zsh completion script from $(COMPLETION_DIR)..."
	@rm -f $(COMPLETION_DIR)/_$(BINARY_NAME)
	@echo "Uninstallation complete. Configuration files are kept."

# Cross-compile for all target platforms
cross-compile: build-mac-universal build-linux build-windows package-all
	@echo "Cross-compilation and packaging finished. Release assets are in the $(OUTPUT_DIR)/ directory."

# Build for Linux (amd64)
build-linux:
	@echo "Building for Linux (amd64)..."
	@mkdir -p $(OUTPUT_DIR)/linux-amd64
	@GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(OUTPUT_DIR)/linux-amd64/$(BINARY_NAME) .

# Build for Windows (amd64)
build-windows:
	@echo "Building for Windows (amd64)..."
	@mkdir -p $(OUTPUT_DIR)/windows-amd64
	@GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(OUTPUT_DIR)/windows-amd64/$(BINARY_NAME).exe .

# Build macOS Universal Binary
build-mac-universal:
	@echo "Building for macOS (Universal)..."
	@mkdir -p $(OUTPUT_DIR)/darwin-universal
	# Build for amd64
	@GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-amd64 .
	# Build for arm64
	@GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-arm64 .
	# Combine with lipo
	@lipo -create -output $(OUTPUT_DIR)/darwin-universal/$(BINARY_NAME) $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-amd64 $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-arm64
	# Ad-hoc sign the universal binary
	@codesign -s - $(OUTPUT_DIR)/darwin-universal/$(BINARY_NAME)
	# Clean up intermediate files
	@rm $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-amd64 $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-arm64
	@echo "Created Universal binary at $(OUTPUT_DIR)/darwin-universal/$(BINARY_NAME)"

# Package all binaries into archives
package-all: package-darwin package-linux package-windows

# Package macOS binary
package-darwin:
	@echo "Packaging macOS binary..."
	@cd $(OUTPUT_DIR)/darwin-universal && tar -czvf ../$(BINARY_NAME)-darwin-universal.tar.gz $(BINARY_NAME)

# Package Linux binary
package-linux:
	@echo "Packaging Linux binary..."
	@cd $(OUTPUT_DIR)/linux-amd64 && tar -czvf ../$(BINARY_NAME)-linux-amd64.tar.gz $(BINARY_NAME)

# Package Windows binary
package-windows:
	@echo "Packaging Windows binary..."
	@cd $(OUTPUT_DIR)/windows-amd64 && zip -r ../$(BINARY_NAME)-windows-amd64.zip $(BINARY_NAME).exe