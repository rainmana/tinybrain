# Getting Started

## Installation

### From Source
```bash
go install github.com/rainmana/tinybrain/cmd/server@latest
```

### Docker
```bash
docker pull rainmana/tinybrain:latest
docker run -p 8080:8080 rainmana/tinybrain
```

### Pre-built Binaries
Download from [Releases](https://github.com/rainmana/tinybrain/releases)

## Basic Usage

### 1. Start the Server
```bash
tinybrain-server --config config.json
```

### 2. Create a Session
```bash
curl -X POST http://localhost:8080/sessions \
  -H "Content-Type: application/json" \
  -d '{"name": "Security Assessment", "task_type": "penetration_test", "intelligence_type": "osint"}'
```

### 3. Store Intelligence
```bash
curl -X POST http://localhost:8080/memory \
  -H "Content-Type: application/json" \
  -d '{"session_id": "session-id", "title": "OSINT Finding", "content": "Social media analysis reveals...", "category": "intelligence", "intelligence_type": "osint", "threat_level": "medium"}'
```

## Configuration

See [Configuration](Configuration) for detailed setup options.

## Next Steps

- [Core Features](Core-Features)
- [Intelligence & Reconnaissance](Intelligence-&-Reconnaissance)
- [API Reference](API-Reference)
