# Educational Game Database Makefile

# Variables
BINARY_NAME=educational-game-db
MAIN_PATH=cmd/cli/main.go
BUILD_DIR=bin
PORT=8081

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME)
	rm -f *.db
	@echo "Clean complete"

# Run the web server
.PHONY: web
web: build
	@echo "Starting web server on port $(PORT)..."
	./$(BUILD_DIR)/$(BINARY_NAME) web --port $(PORT)

# Run interactive mode
.PHONY: interactive
interactive: build
	@echo "Starting interactive mode..."
	./$(BUILD_DIR)/$(BINARY_NAME) interactive

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test ./...

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

# Lint code
.PHONY: lint
lint:
	@echo "Linting code..."
	golangci-lint run

# Development setup
.PHONY: dev-setup
dev-setup: deps
	@echo "Setting up development environment..."
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Create sample data
.PHONY: sample-data
sample-data: build
	@echo "Creating sample data..."
	echo -e "alice\nalice@school.edu\npassword123\nAlice\nJohnson\n3\nElementary School" | ./$(BUILD_DIR)/$(BINARY_NAME) create
	echo -e "bob\nbob@school.edu\npassword123\nBob\nSmith\n4\nElementary School" | ./$(BUILD_DIR)/$(BINARY_NAME) create
	echo -e "charlie\ncharlie@school.edu\npassword123\nCharlie\nBrown\n5\nElementary School" | ./$(BUILD_DIR)/$(BINARY_NAME) create
	@echo "Sample data created"

# Show help
.PHONY: help
help:
	@echo "Educational Game Database - Available Commands:"
	@echo ""
	@echo "  make build         - Build the application"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make web           - Start web server (port $(PORT))"
	@echo "  make interactive   - Start interactive mode"
	@echo "  make deps          - Install dependencies"
	@echo "  make test          - Run tests"
	@echo "  make fmt           - Format code"
	@echo "  make lint          - Lint code"
	@echo "  make dev-setup     - Setup development environment"
	@echo "  make sample-data   - Create sample accounts"
	@echo "  make help          - Show this help"
	@echo ""
	@echo "Web Interface: http://localhost:$(PORT)"
	@echo "Admin Panel:   http://localhost:$(PORT)/admin"

# Docker support
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t educational-game-db .

.PHONY: docker-run
docker-run:
	@echo "Running Docker container..."
	docker run -p $(PORT):$(PORT) educational-game-db
