package handlers

import (
	"net/http"
	"path/filepath"
	"strings"
)

type DocsHandler struct {
	staticPath string
}

func NewDocsHandler(staticPath string) *DocsHandler {
	return &DocsHandler{
		staticPath: staticPath,
	}
}

// ServeDocumentation serves the API documentation files
func (h *DocsHandler) ServeDocumentation(w http.ResponseWriter, r *http.Request) {
	// Security: prevent directory traversal
	path := filepath.Clean(r.URL.Path)
	if strings.Contains(path, "..") {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	// Remove /docs prefix
	path = strings.TrimPrefix(path, "/docs")
	if path == "" || path == "/" {
		path = "/index.html"
	}

	// Construct file path
	filePath := filepath.Join(h.staticPath, path)

	// Set appropriate content type
	switch filepath.Ext(path) {
	case ".html":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	case ".css":
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	case ".js":
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	case ".yaml", ".yml":
		w.Header().Set("Content-Type", "application/x-yaml; charset=utf-8")
	case ".json":
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	// Add security headers
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")

	// Serve the file
	http.ServeFile(w, r, filePath)
}

// GetAPISpec returns the OpenAPI specification
func (h *DocsHandler) GetAPISpec(w http.ResponseWriter, r *http.Request) {
	specPath := filepath.Join(h.staticPath, "openapi.yaml")
	
	w.Header().Set("Content-Type", "application/x-yaml; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	
	http.ServeFile(w, r, specPath)
}

// GetSwaggerUI serves the Swagger UI
func (h *DocsHandler) GetSwaggerUI(w http.ResponseWriter, r *http.Request) {
	uiPath := filepath.Join(h.staticPath, "swagger-ui.html")
	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "SAMEORIGIN")
	
	http.ServeFile(w, r, uiPath)
}

// HealthCheck for docs service
func (h *DocsHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{
		"status": "healthy",
		"service": "docs",
		"message": "Documentation service is running"
	}`))
}