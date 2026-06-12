package app

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/rainmana/tinybrain/internal/mcp"
)

// Version is the binary version, overridable at build time via
// -ldflags "-X github.com/rainmana/tinybrain/internal/app.Version=v1.2.3"
var Version = "dev"

const usage = `TinyBrain 🧠 - Security-focused memory storage MCP server

Usage:
  tinybrain                MCP server on stdio (for MCP client configs)
  tinybrain serve          MCP server on stdio + REST API/dashboard over HTTP
  tinybrain version        Print version
  tinybrain help           Show this help

Flags for serve:
  --dir <path>     Data directory (default: ~/.tinybrain)
  --host <host>    HTTP listen host (default: 127.0.0.1, env: TINYBRAIN_HOST)
  --port <port>    HTTP listen port (default: 8090, env: TINYBRAIN_PORT)
  --http=<bool>    Enable the HTTP REST API (default: true, env: TINYBRAIN_HTTP)

Environment variables:
  TINYBRAIN_DB_PATH    SQLite database path (default: <dir>/memory.db)
  TINYBRAIN_HOST       HTTP listen host
  TINYBRAIN_PORT       HTTP listen port
  TINYBRAIN_HTTP       Enable HTTP API in serve mode (true/false)
  TINYBRAIN_API_TOKEN  If set, REST API requires this bearer token

Endpoints (serve mode):
  http://<host>:<port>/_/           Dashboard
  http://<host>:<port>/api/         REST API index
  MCP protocol                      stdio (JSON-RPC 2.0)

Documentation: https://rainmana.github.io/tinybrain/
`

// Main is the shared entry point for the tinybrain and server binaries.
func Main() {
	args := os.Args[1:]

	command := ""
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		command = args[0]
		args = args[1:]
	}

	switch command {
	case "version", "--version":
		fmt.Printf("tinybrain %s\n", Version)
	case "help":
		fmt.Print(usage)
	case "", "serve":
		runServe(command == "serve", args)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n%s", command, usage)
		os.Exit(1)
	}
}

// runServe starts the MCP stdio server, optionally with the HTTP REST API
// (enabled by default in `serve` mode, disabled in bare stdio mode).
func runServe(serveMode bool, args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	fs.Usage = func() { fmt.Fprint(os.Stderr, usage) }

	dir := fs.String("dir", "", "Data directory (default: ~/.tinybrain)")
	host := fs.String("host", envOr("TINYBRAIN_HOST", "127.0.0.1"), "HTTP listen host")
	port := fs.String("port", envOr("TINYBRAIN_PORT", "8090"), "HTTP listen port")
	httpEnabled := fs.Bool("http", serveMode && envOr("TINYBRAIN_HTTP", "true") != "false", "Enable HTTP REST API")
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}

	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    false,
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
		Prefix:          "TinyBrain 🧠 ",
		Level:           log.InfoLevel,
	})

	dbPath, err := resolveDBPath(*dir)
	if err != nil {
		logger.Fatal("Failed to resolve database path", "error", err)
	}

	server, err := NewTinyBrainServer(dbPath, logger)
	if err != nil {
		logger.Fatal("Failed to initialize TinyBrain", "error", err)
	}
	defer server.Close()

	mcpServer := mcp.NewServer("TinyBrain Memory Storage", Version,
		"Security-focused LLM memory storage MCP server", logger)
	server.registerTools(mcpServer)

	if *httpEnabled {
		addr := fmt.Sprintf("%s:%s", *host, *port)
		go func() {
			if err := server.StartHTTP(addr); err != nil {
				// HTTP failure (e.g. port in use) must not kill the MCP
				// server that an AI client is attached to over stdio
				logger.Error("HTTP server stopped", "error", err)
			}
		}()
		logger.Info("HTTP API enabled", "dashboard", fmt.Sprintf("http://%s/_/", addr), "api", fmt.Sprintf("http://%s/api/", addr))
	}

	logger.Info("Starting TinyBrain MCP Server", "version", Version, "db_path", dbPath)

	if err := mcpServer.ServeStdio(); err != nil {
		logger.Fatal("MCP server error", "error", err)
	}

	// In serve mode, keep the HTTP API alive after stdin closes (e.g. when
	// run standalone in a terminal or with stdin redirected from /dev/null)
	if serveMode && *httpEnabled {
		logger.Info("stdin closed; HTTP API continues. Press Ctrl+C to stop.")
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		logger.Info("Shutting down")
	}
}

// resolveDBPath determines the SQLite path: explicit --dir wins, then
// TINYBRAIN_DB_PATH, then ~/.tinybrain/memory.db.
func resolveDBPath(dir string) (string, error) {
	if dir != "" {
		expanded, err := expandHome(dir)
		if err != nil {
			return "", err
		}
		return filepath.Join(expanded, "memory.db"), nil
	}

	if dbPath := os.Getenv("TINYBRAIN_DB_PATH"); dbPath != "" {
		return expandHome(dbPath)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(homeDir, ".tinybrain", "memory.db"), nil
}

// expandHome expands a leading ~ since MCP client configs pass paths like
// "~/.tinybrain" without shell expansion.
func expandHome(path string) (string, error) {
	if path == "~" || strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %w", err)
		}
		return filepath.Join(homeDir, strings.TrimPrefix(strings.TrimPrefix(path, "~"), "/")), nil
	}
	return path, nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
