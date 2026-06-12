---
layout: default
title: Claude Code Plugin
permalink: /claude-plugin/
---

# Claude Code Plugin

TinyBrain ships as a [Claude Code](https://claude.com/product/claude-code) plugin that does
two things in a single install:

1. **Registers the TinyBrain MCP server** so Claude can call all of TinyBrain's tools.
2. **Installs skills** that teach Claude *how* to use TinyBrain well — capturing findings as
   it works, building a knowledge graph of an engagement, and recalling it across sessions.

Without the skills, TinyBrain is 50+ tools the model has to figure out on its own. With them,
it knows when to reach for memory, how to categorize and prioritize findings, when to link
relationships, and to retrieve before re-deriving — so the tooling is useful out of the box.

## Prerequisite: install the binary

The plugin's MCP configuration launches `tinybrain` from your `PATH`, so install it first:

```bash
go install github.com/rainmana/tinybrain/v3/cmd/tinybrain@latest
```

Make sure `$GOPATH/bin` (or `$GOBIN`) is on your `PATH`. (Note the `/v3` in the module
path — it's required for v3.x releases.)

## Install the plugin

From inside Claude Code, add the marketplace and install:

```
/plugin marketplace add rainmana/tinybrain
/plugin install tinybrain@tinybrain-marketplace
```

That's it. The `tinybrain` MCP server is registered (pointed at `~/.tinybrain/memory.db` by
default) and the skills activate automatically when you do security work. To use a different
database location, set `TINYBRAIN_DB_PATH` in your environment before launching Claude Code.

## What you get

The plugin bundles four skills:

- **using-tinybrain** — the core skill. The capture → connect → checkpoint → retrieve loop,
  memory hygiene (categories, priority/confidence calibration, relationship types, MITRE
  ATT&CK tagging), and the habit of searching memory before re-investigating. Reference files
  cover the full tool inventory and controlled vocabularies.
- **tinybrain-code-review** — security code review: record findings as categorized memories,
  link evidence → vulnerability → fix, and map to CWE/CVE/OWASP.
- **tinybrain-ctf** — CTF and practice-box workflow: track recon, hypotheses, and working
  payloads; never lose a foothold or repeat a dead end; reconstruct the solve from the graph.
- **tinybrain-threat-intel** — intelligence analysis and threat modeling: indicators, actors,
  campaigns, and TTPs mapped to MITRE ATT&CK, with attribution confidence recorded honestly.

These are for **authorized** security work — code review, CTFs and training targets, threat
modeling, and intelligence analysis.

## Manual MCP setup (without the plugin)

If you'd rather register the MCP server by hand (e.g. for Claude Desktop or another MCP
client), add this to your client's configuration instead:

```json
{
  "mcpServers": {
    "tinybrain": {
      "command": "tinybrain",
      "args": [],
      "env": {
        "TINYBRAIN_DB_PATH": "~/.tinybrain/memory.db"
      }
    }
  }
}
```

Running `tinybrain` with no arguments speaks the MCP protocol over stdio — the right mode for
an MCP client. (Use `tinybrain serve` instead if you also want the REST API and dashboard.)

See [AI Integration](/tinybrain/ai-integration/) for more on using TinyBrain with AI
assistants, and [Getting Started](/tinybrain/getting-started/) for the full installation guide.
