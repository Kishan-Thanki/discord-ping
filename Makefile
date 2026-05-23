.PHONY: all build run test lint clean pre-check

# Default target
all: lint build

# Build the binary
build:
	@echo "Building..."
	@go build -o ping-bot main.go

# Run the bot
run:
	@echo "Running..."
	@go run main.go

# Run tests
test:
	@echo "Testing..."
	@go test ./... -v

# Run linter/vet
lint:
	@echo "Vetting code..."
	@go vet ./...

# Pre-push security and compilation check
pre-check:
	@echo "Running pre-check suite..."
	@go fmt ./...
	@make lint
	@make test

# Clean build cache and binary
clean:
	@echo "Cleaning..."
	@rm -f ping-bot
	@go clean
