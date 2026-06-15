# =============================================================================
# Build stage
# =============================================================================
FROM golang:1.25-alpine AS builder
RUN apk add --no-cache git ca-certificates tzdata
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /app/go-http-server ./cmd/server

# =============================================================================
# Production stage
# =============================================================================
FROM alpine:3.19 AS production
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
RUN mkdir -p log/go-http-server \
    && adduser -D -g '' -u 1000 appuser \
    && chown -R appuser:appuser log/
COPY --from=builder --chown=appuser:appuser /app/configs/ ./configs/
COPY --from=builder --chown=appuser:appuser /app/go-http-server .
USER appuser
EXPOSE 7979
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:7979/v1/ping || exit 1
ENTRYPOINT ["./go-http-server"]
