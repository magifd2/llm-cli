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

.PHONY: all build clean test cross-compile build-mac-universal build-linux build-windows

all: build

# Build for the current OS/Arch
build:
	@echo "Building for $(shell go env GOOS)/$(shell go env GOARCH)..."
	@$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) .

# Run tests
test:
	@echo "Running tests..."
	@$(GOTEST) -v ./...

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	@$(GOCLEAN)
	@rm -f $(BINARY_NAME)
	@rm -rf $(OUTPUT_DIR)

# Cross-compile for all target platforms
cross-compile: build-mac-universal build-linux build-windows
	@echo "Cross-compilation finished. Binaries are in the $(OUTPUT_DIR)/ directory."

# Build for Linux (amd64)
build-linux:
	@echo "Building for Linux (amd64)..."
	@mkdir -p $(OUTPUT_DIR)
	@GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-linux-amd64 .

# Build for Windows (amd64)
build-windows:
	@echo "Building for Windows (amd64)..."
	@mkdir -p $(OUTPUT_DIR)
	@GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-windows-amd64.exe .

# Build macOS Universal Binary
build-mac-universal:
	@echo "Building for macOS (Universal)..."
	@mkdir -p $(OUTPUT_DIR)
	# Build for amd64
	@GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-amd64 .
	# Build for arm64
	@GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-arm64 .
	# Combine with lipo
	@lipo -create -output $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-universal $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-amd64 $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-arm64
	# Clean up intermediate files
	@rm $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-amd64 $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-arm64
	@echo "Created Universal binary at $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-universal"
