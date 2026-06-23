.PHONY: all build run test lint clean pre-check

all: lint build

build:
	@echo "Building..."
	@go build -o discord-ping ./cmd/discord-ping/main.go

run:
	@echo "Running..."
	@export $$(cat .env | xargs) && go run cmd/discord-ping/main.go

test:
	@echo "Testing..."
	@go test ./... -v

lint:
	@echo "Vetting code..."
	@go vet ./...

pre-check:
	@echo "Running pre-check suite..."
	@go fmt ./...
	@make lint
	@make test

clean:
	@echo "Cleaning..."
	@rm -f discord-ping
	@go clean
