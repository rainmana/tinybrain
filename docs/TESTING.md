# TinyBrain testing strategy

TinyBrain uses layered testing so user-facing MCP behavior is verified beyond isolated unit tests.

## Test layers

- Unit tests: package-level tests for models, repositories, services, and database behavior.
- In-process MCP user-flow tests: JSON-RPC-shaped calls against the registered MCP server to verify documented aliases and core tool behavior without process startup noise.
- Executable MCP user-flow tests: black-box tests that build the TinyBrain server binary, launch it as a subprocess, communicate over stdio JSON-RPC, and validate the workflow through public MCP tools only.
- CI: GitHub Actions runs the full test suite, vetting, and builds on Linux and Windows.

## Executable user-flow coverage

`TestMCPExecutableAuthorizedAssessmentFlow` verifies the behavior an LLM client depends on:

- MCP initialization and `notifications/initialized` handling.
- Tool discovery for core memory/session/search/relationship/context/task/export tools.
- Health checks against a temporary SQLite database.
- Session creation with metadata.
- Memory create, read, update, search, tag/category filtering, and delete.
- Relationship creation and traversal.
- Context snapshots and context summaries.
- Task progress create, update, and list.
- Session export.
- Database statistics.
- Process restart with the same database to prove persistence across server lifetimes.

This is intentionally more than a smoke test: the test fails if the executable cannot complete the workflow through the same stdio transport a Codex or Claude-style MCP client uses.

## Commands

Run the full suite:

```bash
go test -v ./...
```

Run only the executable MCP user-flow suite:

```bash
go test -v ./cmd/server -run TestMCPExecutableAuthorizedAssessmentFlow -count=1
```

Run vet and build checks:

```bash
go vet ./...
go build -o bin/tinybrain ./cmd/server
```
