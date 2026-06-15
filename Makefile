.PHONY: all run build test lint clean docker-up docker-down docker-logs help

# Environment variables
ENV ?= development

# Default target
all: build

# Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?# "} /^[a-zA-Z_-]+:.*?# / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Run the application locally
run: # Run the application locally
	ENV=$(ENV) go run ./cmd/server

# Build the application
build: # Build the application into bin/ directory
	CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/server ./cmd/server

# Run unit tests
test: # Run unit tests with race detection
	go test -v -cover -race ./...

# Run golangci-lint
lint: # Run golangci-lint
	golangci-lint run ./...

# Clean the build output
clean: # Clean the build output
	rm -rf bin/

# Run the application using Docker Compose
docker-up: # Start the application using Docker Compose
	docker-compose up -d app

# Stop the Docker Compose containers and remove images
docker-down: # Stop the Docker Compose containers and remove images
	docker compose down --rmi all

# View the Docker Compose logs
docker-logs: # View the Docker Compose logs
	docker compose logs -f
