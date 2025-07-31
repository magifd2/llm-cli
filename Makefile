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

.PHONY: all build clean test cross-compile build-mac-universal build-linux build-windows package-all

all: build cross-compile

# Build for the current OS/Arch
build:
	@echo "Building for $(shell go env GOOS)/$(shell go env GOARCH)..."
	@mkdir -p $(OUTPUT_DIR)/$(shell go env GOOS)-$(shell go env GOARCH)
	@$(GOBUILD) $(LDFLAGS) -o $(OUTPUT_DIR)/$(shell go env GOOS)-$(shell go env GOARCH)/$(BINARY_NAME) .

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
	@rm -f extract_release_notes.go release_notes.txt

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