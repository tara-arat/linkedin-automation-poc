# LinkedIn Automation POC - Makefile

# Variables
BINARY_NAME=linkedin-automation
MAIN_PATH=./cmd/linkedin-automation
BUILD_DIR=./build
GO=go
GOFLAGS=-v

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@$(GO) build $(GOFLAGS) -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BINARY_NAME)"

# Build for multiple platforms
.PHONY: build-all
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=windows GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@GOOS=linux GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	@GOOS=darwin GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	@GOOS=darwin GOARCH=arm64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	@echo "Multi-platform build complete in $(BUILD_DIR)"

# Run the application
.PHONY: run
run:
	@echo "Running $(BINARY_NAME)..."
	@$(GO) run $(MAIN_PATH)

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME)
	@rm -f $(BINARY_NAME).exe
	@rm -rf $(BUILD_DIR)
	@rm -rf data/
	@rm -rf logs/
	@rm -rf sessions/
	@echo "Clean complete"

# Download dependencies
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	@$(GO) mod download
	@echo "Dependencies downloaded"

# Tidy dependencies
.PHONY: tidy
tidy:
	@echo "Tidying dependencies..."
	@$(GO) mod tidy
	@echo "Dependencies tidied"

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@$(GO) fmt ./...
	@echo "Code formatted"

# Run linter
.PHONY: lint
lint:
	@echo "Running linter..."
	@golangci-lint run || echo "golangci-lint not installed, skipping..."

# Install the application
.PHONY: install
install:
	@echo "Installing $(BINARY_NAME)..."
	@$(GO) install $(MAIN_PATH)
	@echo "Installation complete"

# Show help
.PHONY: help
help:
	@echo "LinkedIn Automation POC - Available targets:"
	@echo "  make build      - Build the application"
	@echo "  make build-all  - Build for multiple platforms"
	@echo "  make run        - Run the application"
	@echo "  make clean      - Clean build artifacts and data"
	@echo "  make deps       - Download dependencies"
	@echo "  make tidy       - Tidy dependencies"
	@echo "  make fmt        - Format code"
	@echo "  make lint       - Run linter"
	@echo "  make install    - Install the application"
	@echo "  make help       - Show this help message"
