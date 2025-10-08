# TinyBrain Memory Storage MCP Server
# Security-focused LLM memory storage application

.PHONY: build install test clean run help

# Variables
BINARY_NAME=tinybrain
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION=$(shell go version | awk '{print $$3}')

# Build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GoVersion=$(GO_VERSION)"

# Default target
all: build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/server
	@echo "Build complete: bin/$(BINARY_NAME)"

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p bin
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-amd64 ./cmd/server
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-amd64 ./cmd/server
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-arm64 ./cmd/server
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-windows-amd64.exe ./cmd/server
	@echo "Multi-platform build complete"

# Install the binary to GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	@go install $(LDFLAGS) ./cmd/server
	@echo "Installation complete"

# Run tests
test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Test coverage report generated: coverage.html"

# Run tests with verbose output
test-verbose:
	@echo "Running tests with verbose output..."
	@go test -v -race ./...

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

# Run the server
run:
	@echo "Running $(BINARY_NAME) server..."
	@go run ./cmd/server

# Run with development database
run-dev:
	@echo "Running $(BINARY_NAME) with development database..."
	@TINYBRAIN_DB_PATH=./dev.db go run ./cmd/server

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@rm -f dev.db*
	@echo "Clean complete"

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Code formatted"

# Run linter
lint:
	@echo "Running linter..."
	@go vet ./...
	@echo "Linting complete"

# Run security checks
security:
	@echo "Running security checks..."
	@go list -json -deps ./... | nancy sleuth
	@echo "Security check complete"

# Generate documentation
docs:
	@echo "Generating documentation..."
	@go doc -all ./... > docs/api.md
	@echo "Documentation generated: docs/api.md"

# Create release package
release: build-all
	@echo "Creating release package..."
	@mkdir -p release
	@tar -czf release/$(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz -C bin $(BINARY_NAME)-linux-amd64
	@tar -czf release/$(BINARY_NAME)-$(VERSION)-darwin-amd64.tar.gz -C bin $(BINARY_NAME)-darwin-amd64
	@tar -czf release/$(BINARY_NAME)-$(VERSION)-darwin-arm64.tar.gz -C bin $(BINARY_NAME)-darwin-arm64
	@zip -j release/$(BINARY_NAME)-$(VERSION)-windows-amd64.zip bin/$(BINARY_NAME)-windows-amd64.exe
	@echo "Release packages created in release/"

# Development setup
dev-setup:
	@echo "Setting up development environment..."
	@go mod download
	@go mod tidy
	@echo "Development setup complete"

# Database operations
db-init:
	@echo "Initializing database..."
	@mkdir -p ~/.tinybrain
	@echo "Database directory created: ~/.tinybrain"

db-reset:
	@echo "Resetting database..."
	@rm -f ~/.tinybrain/memory.db*
	@echo "Database reset complete"

# Docker operations
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(BINARY_NAME):$(VERSION) .
	@echo "Docker image built: $(BINARY_NAME):$(VERSION)"

docker-run:
	@echo "Running Docker container..."
	@docker run --rm -it $(BINARY_NAME):$(VERSION)

# Help
help:
	@echo "TinyBrain Memory Storage MCP Server"
	@echo ""
	@echo "Available targets:"
	@echo "  build        - Build the binary"
	@echo "  build-all    - Build for multiple platforms"
	@echo "  install      - Install to GOPATH/bin"
	@echo "  test         - Run tests with coverage"
	@echo "  test-verbose - Run tests with verbose output"
	@echo "  bench        - Run benchmarks"
	@echo "  run          - Run the server"
	@echo "  run-dev      - Run with development database"
	@echo "  clean        - Clean build artifacts"
	@echo "  fmt          - Format code"
	@echo "  lint         - Run linter"
	@echo "  security     - Run security checks"
	@echo "  docs         - Generate documentation"
	@echo "  release      - Create release packages"
	@echo "  dev-setup    - Setup development environment"
	@echo "  db-init      - Initialize database directory"
	@echo "  db-reset     - Reset database"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  help         - Show this help"
	@echo ""
	@echo "Environment variables:"
	@echo "  TINYBRAIN_DB_PATH - Path to SQLite database (default: ~/.tinybrain/memory.db)"
