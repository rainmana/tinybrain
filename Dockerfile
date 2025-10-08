# Multi-stage build for TinyBrain Memory Storage MCP Server
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o tinybrain ./cmd/server

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates sqlite

# Create non-root user
RUN addgroup -g 1001 -S tinybrain && \
    adduser -u 1001 -S tinybrain -G tinybrain

# Create application directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/tinybrain .

# Create data directory
RUN mkdir -p /app/data && \
    chown -R tinybrain:tinybrain /app

# Switch to non-root user
USER tinybrain

# Expose port (if needed for future HTTP transport)
EXPOSE 8080

# Set environment variables
ENV TINYBRAIN_DB_PATH=/app/data/memory.db

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD ./tinybrain --health-check || exit 1

# Run the application
ENTRYPOINT ["./tinybrain"]
