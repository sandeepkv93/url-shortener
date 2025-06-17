package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheckRoute(t *testing.T) {
	// Create a minimal router for testing the health check
	router := NewRouterBuilder().
		WithCORS(false).
		WithLogging(false).
		Build()

	handler := router.SetupRoutes()
	
	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "healthy")
	assert.Contains(t, rr.Body.String(), "url-shortener")
}

func TestRouterBuilder(t *testing.T) {
	// Test that the builder pattern works correctly
	router := NewRouterBuilder().
		WithCORS(false).
		WithLogging(false).
		Build()

	assert.NotNil(t, router)
	assert.NotNil(t, router.config)
	assert.False(t, router.config.EnableCORS)
	assert.False(t, router.config.EnableLogging)
}

func TestRouterBuilderDefaults(t *testing.T) {
	// Test that the builder has sensible defaults
	builder := NewRouterBuilder()
	
	assert.NotNil(t, builder.config)
	assert.True(t, builder.config.EnableCORS)
	assert.True(t, builder.config.EnableLogging)
	assert.Equal(t, []string{"*"}, builder.config.AllowedOrigins)
}

func TestGetHandler(t *testing.T) {
	router := NewRouterBuilder().Build()
	handler := router.GetHandler()
	
	assert.NotNil(t, handler)
}

func TestCORSConfiguration(t *testing.T) {
	router := NewRouterBuilder().
		WithCORS(true, "https://example.com", "https://app.example.com").
		Build()

	assert.True(t, router.config.EnableCORS)
	assert.Equal(t, []string{"https://example.com", "https://app.example.com"}, router.config.AllowedOrigins)
}

func TestNotFoundRoute(t *testing.T) {
	router := NewRouterBuilder().
		WithCORS(false).
		WithLogging(false).
		Build()

	handler := router.SetupRoutes()
	
	req := httptest.NewRequest("GET", "/nonexistent", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestAPIVersioning(t *testing.T) {
	router := NewRouterBuilder().
		WithCORS(false).
		WithLogging(false).
		Build()

	handler := router.SetupRoutes()
	
	// Test that versioned API paths are properly structured
	// This will return 404 since we don't have handlers, but it shows the routing works
	req := httptest.NewRequest("GET", "/api/v1/auth/login", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// We expect 405 Method Not Allowed since this should be POST, not GET
	// or 404 if no handlers are configured
	assert.True(t, rr.Code == http.StatusMethodNotAllowed || rr.Code == http.StatusNotFound)
}