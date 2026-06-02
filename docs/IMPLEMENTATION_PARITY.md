# TinyBrain implementation parity

This note tracks the current implementation against the published GitHub Pages documentation at https://rainmana.github.io/tinybrain/. It should be updated whenever the MCP surface, storage model, or documentation changes.

## Current shipped core

TinyBrain is currently a SQLite-backed MCP stdio server for authorized security assessment memory workflows. The verified core is:

- MCP JSON-RPC over stdio with `initialize`, `tools/list`, and `tools/call`.
- Session management for security review, penetration testing, exploit development, vulnerability analysis, threat modeling, incident response, intelligence analysis, reverse engineering, malware analysis, and general workflows.
- Memory create/read/update/delete with priority, confidence, source, tags, security-focused categories, and access tracking.
- Search by semantic/FTS-style text, exact text, tags, category, priority, confidence, and session.
- Relationship creation and traversal, including documented aliases such as `references`.
- Context snapshots and context summaries for LLM continuity.
- Task progress create/update/list/get for multi-stage assessment work.
- Batch memory operations, export/import, cleanup, diagnostics, local embeddings, similarity, duplicate checks, and notifications.
- Security repository/query scaffolding for NVD, MITRE ATT&CK, and OWASP data, with the core memory path not depending on network downloads.

## Documentation aliases supported

The server accepts both the local tool names and the names shown in the published MCP examples:

- `mcp_tinybrain-mcp-server_<tool>` is accepted as an alias for `<tool>`.
- `create_memory` is accepted as an alias for `store_memory`.
- `search_memory` is accepted as an alias for `search_memories`.
- `get_related_entries` is accepted as an alias for `get_related_memories`.

MCP arguments now accept LLM-friendly JSON arrays or JSON array strings for fields such as `tags` and `categories`.

## Intentionally deferred

These items appear in the published GitHub Pages documentation but are not part of the current Go core:

- PocketBase admin dashboard, REST API, GraphQL API, SDKs, API keys, JWT auth, and HTTP rate limiting.
- Live NVD synchronization as a core startup/runtime dependency.
- Cloud integrations, SIEM integrations, Kubernetes deployment, and multi-user access control.
- Full reverse-engineering, malware sandbox, debugger, decompiler, and external tool orchestration.

For NVD/CVE-scale data, prefer manual file loading first. If network/API ingestion is added later, keep it behind a plugin-style integration so the core MCP memory system remains deterministic and usable offline.

## Python fork reference

The Python fork at https://github.com/rainmana/tinybrain-python is a useful reference for MCP ergonomics, deterministic local similarity, and a future web UI shape. It also documents a later "librarian agent" direction. Do not port the known scoring bug where entries tend to cluster around a rating of about `7`; scoring changes should have regression tests with varied expected scores.

## Verification

Current verification includes:

- Repository unit tests for sessions, memory CRUD, search, relationships, context snapshots, task progress, CVE mapping, risk correlation, compliance mapping, and security repositories.
- MCP user-flow tests in `cmd/server` that exercise documented prefixed tool names, aliases, intelligence sessions, memory CRUD, search, relationships, context snapshots, and task progress through JSON-RPC-shaped calls.
- Pure-Go SQLite storage via `modernc.org/sqlite`, so tests and local builds work without CGO or a C compiler.
- Dockerfile aligned to Go 1.24 and pure-Go SQLite.

Run the suite with:

```bash
go test ./...
go build -o tinybrain ./cmd/server
go run ./cmd/server --health-check
```
