.PHONY: all build test clean proto run install

# Variables
BINARY_NAME=server
BUILD_DIR=bin
SERVER_PATH=./cmd/server/main.go
PROTO_DIR=proto
API_DIR=api/auth

all: test build

# Install dependencies
install:
	go mod download
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Build the server
build:
	@echo "Building server..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(SERVER_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Generate protobuf code
proto:
	@echo "Generating protobuf code..."
	@mkdir -p $(API_DIR)
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/auth.proto
	mv $(PROTO_DIR)/auth.pb.go $(PROTO_DIR)/auth_grpc.pb.go $(API_DIR)/
	@echo "Protobuf code generated"

# Run the server
run: build
	@echo "Starting server..."
	$(BUILD_DIR)/$(BINARY_NAME)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	go clean
	@echo "Clean complete"

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run ./...

# Show help
help:
	@echo "Available targets:"
	@echo "  all            - Run tests and build"
	@echo "  build          - Build the server"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  proto          - Generate protobuf code"
	@echo "  run            - Build and run the server"
	@echo "  clean          - Remove build artifacts"
	@echo "  fmt            - Format code"
	@echo "  lint           - Run linter"
	@echo "  install        - Install dependencies"
	@echo "  help           - Show this help message"
