package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CORSMiddlewareTestSuite struct {
	suite.Suite
}

func TestCORSMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(CORSMiddlewareTestSuite))
}

func (suite *CORSMiddlewareTestSuite) TestDefaultCORSConfig() {
	config := DefaultCORSConfig()
	
	assert.Equal(suite.T(), []string{"*"}, config.AllowedOrigins)
	assert.Contains(suite.T(), config.AllowedMethods, "GET")
	assert.Contains(suite.T(), config.AllowedMethods, "POST")
	assert.Contains(suite.T(), config.AllowedMethods, "OPTIONS")
	assert.Contains(suite.T(), config.AllowedHeaders, "Authorization")
	assert.True(suite.T(), config.AllowCredentials)
	assert.Equal(suite.T(), 86400, config.MaxAge)
}

func (suite *CORSMiddlewareTestSuite) TestProductionCORSConfig() {
	allowedOrigins := []string{"https://example.com", "https://app.example.com"}
	config := ProductionCORSConfig(allowedOrigins)
	
	assert.Equal(suite.T(), allowedOrigins, config.AllowedOrigins)
	assert.True(suite.T(), config.AllowCredentials)
}

func (suite *CORSMiddlewareTestSuite) TestCORSMiddleware_WithWildcardOrigin() {
	middleware := NewCORSMiddleware(DefaultCORSConfig())
	
	handler := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	assert.Equal(suite.T(), "*", rr.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(suite.T(), "true", rr.Header().Get("Access-Control-Allow-Credentials"))
}

func (suite *CORSMiddlewareTestSuite) TestCORSMiddleware_WithSpecificOrigin() {
	config := &CORSConfig{
		AllowedOrigins:   []string{"https://example.com", "https://app.example.com"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}
	
	middleware := NewCORSMiddleware(config)
	
	handler := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	assert.Equal(suite.T(), "https://example.com", rr.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(suite.T(), "true", rr.Header().Get("Access-Control-Allow-Credentials"))
	assert.Equal(suite.T(), "GET, POST", rr.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(suite.T(), "Content-Type, Authorization", rr.Header().Get("Access-Control-Allow-Headers"))
}

func (suite *CORSMiddlewareTestSuite) TestCORSMiddleware_UnallowedOrigin() {
	config := &CORSConfig{
		AllowedOrigins:   []string{"https://example.com"},
		AllowCredentials: true,
	}
	
	middleware := NewCORSMiddleware(config)
	
	handler := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://malicious.com")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	assert.Equal(suite.T(), "", rr.Header().Get("Access-Control-Allow-Origin"))
}

func (suite *CORSMiddlewareTestSuite) TestCORSMiddleware_WildcardSubdomain() {
	config := &CORSConfig{
		AllowedOrigins: []string{"*.example.com"},
	}
	
	middleware := NewCORSMiddleware(config)
	
	handler := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Test subdomain
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://app.example.com")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(suite.T(), "https://app.example.com", rr.Header().Get("Access-Control-Allow-Origin"))

	// Test root domain
	req = httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(suite.T(), "https://example.com", rr.Header().Get("Access-Control-Allow-Origin"))

	// Test non-matching domain
	req = httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://other.com")
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(suite.T(), "", rr.Header().Get("Access-Control-Allow-Origin"))
}

func (suite *CORSMiddlewareTestSuite) TestCORSMiddleware_PreflightRequest() {
	config := &CORSConfig{
		AllowedOrigins: []string{"https://example.com"},
		AllowedMethods: []string{"GET", "POST", "PUT"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		MaxAge:         3600,
	}
	
	middleware := NewCORSMiddleware(config)
	
	handler := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This should not be called for OPTIONS requests
		suite.T().Error("Handler should not be called for OPTIONS request")
	}))

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusNoContent, rr.Code)
	assert.Equal(suite.T(), "https://example.com", rr.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(suite.T(), "GET, POST, PUT", rr.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(suite.T(), "Content-Type, Authorization", rr.Header().Get("Access-Control-Allow-Headers"))
}

func (suite *CORSMiddlewareTestSuite) TestCORSMiddleware_ExposedHeaders() {
	config := &CORSConfig{
		AllowedOrigins: []string{"*"},
		ExposedHeaders: []string{"X-Total-Count", "X-Page-Count"},
	}
	
	middleware := NewCORSMiddleware(config)
	
	handler := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(suite.T(), "X-Total-Count, X-Page-Count", rr.Header().Get("Access-Control-Expose-Headers"))
}

func (suite *CORSMiddlewareTestSuite) TestCORSMiddleware_NoCredentials() {
	config := &CORSConfig{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: false,
	}
	
	middleware := NewCORSMiddleware(config)
	
	handler := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(suite.T(), "", rr.Header().Get("Access-Control-Allow-Credentials"))
}

func (suite *CORSMiddlewareTestSuite) TestCORSConvenienceFunction() {
	// Test with no origins (should use default)
	handler := CORS()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(suite.T(), "*", rr.Header().Get("Access-Control-Allow-Origin"))

	// Test with specific origins
	handler = CORS("https://example.com", "https://app.example.com")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req = httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(suite.T(), "https://example.com", rr.Header().Get("Access-Control-Allow-Origin"))
}

func (suite *CORSMiddlewareTestSuite) TestIsOriginAllowed() {
	config := &CORSConfig{
		AllowedOrigins: []string{"https://example.com", "*.app.com", "*"},
	}
	middleware := NewCORSMiddleware(config)

	// Test exact match
	assert.True(suite.T(), middleware.isOriginAllowed("https://example.com"))
	
	// Test wildcard
	assert.True(suite.T(), middleware.isOriginAllowed("https://anything.com"))
	
	// Test wildcard subdomain
	assert.True(suite.T(), middleware.isOriginAllowed("https://test.app.com"))
	assert.True(suite.T(), middleware.isOriginAllowed("https://app.com"))
	
	// Test non-matching
	config.AllowedOrigins = []string{"https://example.com"}
	middleware = NewCORSMiddleware(config)
	assert.False(suite.T(), middleware.isOriginAllowed("https://other.com"))
}