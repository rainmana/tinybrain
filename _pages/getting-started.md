---
layout: default
title: Getting Started
permalink: /getting-started/
---

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
  -d '{
    "name": "Security Assessment",
    "task_type": "penetration_test",
    "intelligence_type": "osint"
  }'
```

### 3. Store Intelligence
```bash
curl -X POST http://localhost:8080/memory \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "session-id",
    "title": "OSINT Finding",
    "content": "Social media analysis reveals...",
    "category": "intelligence",
    "intelligence_type": "osint",
    "threat_level": "medium"
  }'
```

## Configuration

### Basic Configuration
Create a `config.json` file:

```json
{
  "database": {
    "path": "./tinybrain.db",
    "max_connections": 10
  },
  "server": {
    "host": "localhost",
    "port": 8080
  },
  "security": {
    "classification_levels": ["unclassified", "confidential", "secret"],
    "threat_levels": ["low", "medium", "high", "critical"]
  }
}
```

### Environment Variables
```bash
export TINYBRAIN_DB_PATH="./tinybrain.db"
export TINYBRAIN_HOST="localhost"
export TINYBRAIN_PORT="8080"
```

## Next Steps

- [Core Features](core-features/) - Learn about memory management and sessions
- [Intelligence & Reconnaissance](intelligence/) - Explore intelligence gathering capabilities
- [API Reference](api-reference/) - Complete API documentation
