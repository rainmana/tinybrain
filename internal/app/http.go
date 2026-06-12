package app

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// toolHandlerFunc matches the signature of the MCP tool handlers so the REST
// API reuses them directly, guaranteeing identical behavior on both surfaces.
type toolHandlerFunc func(ctx context.Context, params map[string]interface{}) (interface{}, error)

// StartHTTP serves the REST API documented at
// https://rainmana.github.io/tinybrain/api-reference/ plus a small dashboard.
// Endpoints are registered both at the root (POST /sessions) and under /api/
// (POST /api/sessions).
func (t *TinyBrainServer) StartHTTP(addr string) error {
	mux := http.NewServeMux()

	register := func(pattern string, handler http.HandlerFunc) {
		method, path, _ := strings.Cut(pattern, " ")
		mux.HandleFunc(pattern, handler)
		mux.HandleFunc(method+" /api"+path, handler)
	}

	register("POST /auth", t.httpAuth)
	register("GET /health", t.httpHealth)
	register("POST /sessions", t.httpToolJSON(t.handleCreateSession, http.StatusCreated))
	register("GET /sessions", t.httpToolQuery(t.handleListSessions, nil))
	register("GET /sessions/{id}", t.httpToolPath(t.handleGetSession, "session_id"))
	register("POST /memory", t.httpToolJSON(t.handleStoreMemory, http.StatusCreated))
	register("GET /memory/search", t.httpToolQuery(t.handleSearchMemories, []string{"limit", "offset", "min_priority", "min_confidence"}))
	register("GET /memory/{id}", t.httpToolPath(t.handleGetMemory, "memory_id"))

	mux.HandleFunc("GET /api/", t.httpAPIIndex)
	mux.HandleFunc("GET /api", t.httpAPIIndex)
	mux.HandleFunc("GET /_/", t.httpDashboard)
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/_/", http.StatusFound)
	})

	server := &http.Server{
		Addr:              addr,
		Handler:           t.authMiddleware(mux),
		ReadHeaderTimeout: 10 * time.Second,
	}
	return server.ListenAndServe()
}

// authMiddleware enforces bearer-token auth on data endpoints when
// TINYBRAIN_API_TOKEN is configured. The dashboard, API index, health check,
// and /auth itself stay open.
func (t *TinyBrainServer) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := os.Getenv("TINYBRAIN_API_TOKEN")
		path := strings.TrimPrefix(r.URL.Path, "/api")
		open := path == "/" || path == "" || path == "/auth" || path == "/health" || strings.HasPrefix(r.URL.Path, "/_/")
		if token != "" && !open {
			auth := r.Header.Get("Authorization")
			provided, _ := strings.CutPrefix(auth, "Bearer ")
			if subtle.ConstantTimeCompare([]byte(provided), []byte(token)) != 1 {
				writeJSON(w, http.StatusUnauthorized, map[string]interface{}{
					"error": "unauthorized: missing or invalid bearer token",
				})
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// httpAuth validates a token posted as {"token": "..."} against
// TINYBRAIN_API_TOKEN. With no token configured, the API is open.
func (t *TinyBrainServer) httpAuth(w http.ResponseWriter, r *http.Request) {
	configured := os.Getenv("TINYBRAIN_API_TOKEN")
	if configured == "" {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"authenticated": true,
			"note":          "no API token configured; set TINYBRAIN_API_TOKEN to require authentication",
		})
		return
	}

	var body struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]interface{}{"error": "invalid JSON body"})
		return
	}
	if subtle.ConstantTimeCompare([]byte(body.Token), []byte(configured)) != 1 {
		writeJSON(w, http.StatusUnauthorized, map[string]interface{}{"authenticated": false})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"authenticated": true,
		"token_type":    "Bearer",
		"usage":         "send 'Authorization: Bearer <token>' on subsequent requests",
	})
}

func (t *TinyBrainServer) httpHealth(w http.ResponseWriter, r *http.Request) {
	if err := t.db.HealthCheck(); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]interface{}{"status": "unhealthy", "error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"status": "healthy", "version": Version})
}

// httpToolJSON adapts an MCP tool handler to a JSON-body POST endpoint.
func (t *TinyBrainServer) httpToolJSON(handler toolHandlerFunc, successStatus int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var params map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]interface{}{"error": "invalid JSON body: " + err.Error()})
			return
		}
		result, err := handler(r.Context(), params)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
			return
		}
		writeJSON(w, successStatus, result)
	}
}

// httpToolQuery adapts an MCP tool handler to a GET endpoint, mapping query
// string values to handler params. Names listed in numeric are converted to
// float64 to mirror JSON number decoding.
func (t *TinyBrainServer) httpToolQuery(handler toolHandlerFunc, numeric []string) http.HandlerFunc {
	isNumeric := make(map[string]bool, len(numeric)+2)
	for _, n := range append([]string{"limit", "offset"}, numeric...) {
		isNumeric[n] = true
	}
	return func(w http.ResponseWriter, r *http.Request) {
		params := make(map[string]interface{})
		for key, values := range r.URL.Query() {
			if len(values) == 0 || values[0] == "" {
				continue
			}
			if isNumeric[key] {
				var f float64
				if _, err := fmt.Sscanf(values[0], "%g", &f); err == nil {
					params[key] = f
					continue
				}
			}
			params[key] = values[0]
		}
		result, err := handler(r.Context(), params)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, result)
	}
}

// httpToolPath adapts an MCP tool handler to a GET endpoint that takes a
// single {id} path parameter.
func (t *TinyBrainServer) httpToolPath(handler toolHandlerFunc, paramName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := handler(r.Context(), map[string]interface{}{paramName: r.PathValue("id")})
		if err != nil {
			status := http.StatusBadRequest
			if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "no rows") {
				status = http.StatusNotFound
			}
			writeJSON(w, status, map[string]interface{}{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, result)
	}
}

func (t *TinyBrainServer) httpAPIIndex(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"name":          "TinyBrain",
		"version":       Version,
		"description":   "Security-focused LLM memory storage",
		"documentation": "https://rainmana.github.io/tinybrain/",
		"endpoints": []string{
			"POST /auth",
			"GET /health",
			"POST /sessions",
			"GET /sessions",
			"GET /sessions/{id}",
			"POST /memory",
			"GET /memory/search?query=...",
			"GET /memory/{id}",
		},
	})
}

func (t *TinyBrainServer) httpDashboard(w http.ResponseWriter, r *http.Request) {
	stats, err := t.db.GetStats()
	if err != nil {
		http.Error(w, "failed to load stats: "+err.Error(), http.StatusInternalServerError)
		return
	}

	statsJSON, _ := json.MarshalIndent(stats, "", "  ")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<title>TinyBrain 🧠 Dashboard</title>
<style>
body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; max-width: 800px; margin: 2rem auto; padding: 0 1rem; background: #0d1117; color: #c9d1d9; }
h1 { color: #58a6ff; }
a { color: #58a6ff; }
pre { background: #161b22; padding: 1rem; border-radius: 8px; overflow-x: auto; }
.card { background: #161b22; border: 1px solid #30363d; border-radius: 8px; padding: 1rem; margin: 1rem 0; }
</style>
</head>
<body>
<h1>TinyBrain 🧠</h1>
<p>Security-focused LLM memory storage &mdash; version %s</p>
<div class="card">
<h2>Database Statistics</h2>
<pre>%s</pre>
</div>
<div class="card">
<h2>REST API</h2>
<p>Index: <a href="/api/">/api/</a> &middot; Health: <a href="/health">/health</a></p>
</div>
<div class="card">
<h2>Documentation</h2>
<p><a href="https://rainmana.github.io/tinybrain/">rainmana.github.io/tinybrain</a></p>
</div>
</body>
</html>`, Version, statsJSON)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		// Headers already sent; nothing more we can do
		_ = err
	}
}
