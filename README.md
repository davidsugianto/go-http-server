# go-http-server

A Go HTTP Server implementing Clean Architecture principles, built with [Gin](https://github.com/gin-gonic/gin).

## Features

- **Clean Architecture**: Separation of concerns across handlers, use cases, and repositories.
- **Graceful Shutdown**: Safe shutdown handling using `grace`.
- **Structured Logging**: Pre-configured error, info, and debug log files.
- **Middleware Included**: Request ID injection, CORS, and Panic Recovery.
- **Docker Ready**: Multi-stage Dockerfile for minimal production images.
- **Kubernetes Ready**: Kustomize overlays available in `deployments/kubernetes`.

## Prerequisites

- Go 1.25+
- Docker (optional, for containerized builds)

## Configuration

Configuration is managed via YAML files located in `configs/`. The application determines which config to load based on the `ENV` environment variable.

For example, when `ENV=development`, the server loads `configs/config.development.yaml`.

## Running the Application

### Local Development

1. Ensure dependencies are downloaded:
   ```bash
   go mod download
   ```

2. Run the server using the development environment:
   ```bash
   ENV=development go run cmd/server/*.go
   ```
   By default, the server will start on port `7979`.

### Using Docker

1. Build the Docker image:
   ```bash
   docker build -t go-http-server .
   ```

2. Run the container:
   ```bash
   docker run -p 7979:7979 -e ENV=development go-http-server
   ```

## API Endpoints

### Health Check

- **GET** `/v1/ping`
  - Returns the health status of the application.

## Project Structure

```text
.
├── cmd/
│   └── server/          # Application entrypoints
├── configs/             # Configuration files by environment
├── deployments/         # K8s manifests and deployment scripts
├── internal/
│   ├── handler/         # HTTP handlers and middleware
│   └── pkg/             # Internal packages (config, etc)
├── Dockerfile           # Multi-stage Docker build
└── go.mod               # Go module dependencies
```
