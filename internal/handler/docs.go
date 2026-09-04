package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// OpenAPIHandler serves the OpenAPI/Swagger specification and an interactive
// Swagger UI at /api/docs. Non-essential — failures are non-fatal.
type OpenAPIHandler struct {
	specPath string
	spec     []byte
}

// NewOpenAPIHandler returns a handler that locates the OpenAPI YAML in the
// repo's docs/ directory (best-effort via several candidate paths).
func NewOpenAPIHandler() *OpenAPIHandler {
	h := &OpenAPIHandler{}
	h.locate()
	return h
}

func (h *OpenAPIHandler) locate() {
	candidates := []string{
		"docs/openapi.yaml",
		"../docs/openapi.yaml",
		"../../docs/openapi.yaml",
		filepath.Join(filepath.Dir(os.Args[0]), "docs", "openapi.yaml"),
		filepath.Join(filepath.Dir(os.Args[0]), "..", "docs", "openapi.yaml"),
		filepath.Join(filepath.Dir(os.Args[0]), "..", "..", "docs", "openapi.yaml"),
	}
	if exe, err := os.Executable(); err == nil {
		dir := filepath.Dir(exe)
		for i := 0; i < 4; i++ {
			candidates = append(candidates, filepath.Join(dir, "..", strings.Repeat("../", i), "docs", "openapi.yaml"))
		}
	}
	for _, c := range candidates {
		if b, err := os.ReadFile(c); err == nil {
			h.specPath = c
			h.spec = b
			return
		}
	}
}

// ServeHTTP routes /api/docs and /api/docs/ (Swagger UI).
func (h *OpenAPIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}
	path := strings.TrimSuffix(r.URL.Path, "/")
	if strings.HasSuffix(path, "/ui") {
		h.serveUI(w)
		return
	}
	h.serveSpec(w, r)
}

func (h *OpenAPIHandler) serveSpec(w http.ResponseWriter, r *http.Request) {
	if h.spec == nil {
		WriteError(w, http.StatusServiceUnavailable, "SPEC_UNAVAILABLE", "OpenAPI spec not found")
		return
	}
	w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(h.spec)
}

func (h *OpenAPIHandler) serveUI(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>BATIQA API — Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
  <style>html{box-sizing:border-box}*,*:before,*:after{box-sizing:inherit}body{margin:0}</style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = function(){
      window.ui = SwaggerUIBundle({
        url: '/api/docs',
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIBundle.SwaggerUIStandalonePreset
        ],
        layout: 'BaseLayout'
      });
    };
  </script>
</body>
</html>`
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(html))
}
