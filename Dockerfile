# Multi-stage build for TinyBrain Memory Storage MCP Server
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Pure-Go SQLite driver (modernc.org/sqlite), so CGO is not required
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-s -w' -o tinybrain ./cmd/tinybrain

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates wget

RUN addgroup -g 1001 -S tinybrain && \
    adduser -u 1001 -S tinybrain -G tinybrain

WORKDIR /app

COPY --from=builder /app/tinybrain .

RUN mkdir -p /app/data && \
    chown -R tinybrain:tinybrain /app

USER tinybrain

# REST API and dashboard
EXPOSE 8090

ENV TINYBRAIN_DB_PATH=/app/data/memory.db

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget -q -O /dev/null http://127.0.0.1:8090/health || exit 1

# Bind to all interfaces so the published port is reachable from the host
ENTRYPOINT ["./tinybrain", "serve", "--host", "0.0.0.0", "--port", "8090"]
