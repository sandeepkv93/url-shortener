package integration

import (
	"net/http"
	"strings"
	"time"
)

// MiddlewareIntegrationTestSuite tests middleware integration and request flow
type MiddlewareIntegrationTestSuite struct {
	IntegrationTestSuite
}

// TestAuthMiddlewareIntegration tests authentication middleware integration
func (s *MiddlewareIntegrationTestSuite) TestAuthMiddlewareIntegration() {
	// Test endpoints that require authentication
	protectedEndpoints := []struct {
		method string
		path   string
	}{
		{"GET", "/api/v1/auth/profile"},
		{"PUT", "/api/v1/auth/profile"},
		{"POST", "/api/v1/auth/change-password"},
		{"POST", "/api/v1/auth/logout"},
		{"GET", "/api/v1/urls/"},
		{"POST", "/api/v1/urls/"},
		{"GET", "/api/v1/analytics/dashboard"},
		{"GET", "/api/v1/analytics/global"},
	}

	for _, endpoint := range protectedEndpoints {
		// Test without authorization header
		resp := s.makeRequest(endpoint.method, endpoint.path, nil, nil)
		s.assertErrorResponse(resp, http.StatusUnauthorized, "authorization header required")

		// Test with invalid authorization header format
		invalidAuthResp := s.makeRequest(endpoint.method, endpoint.path, nil, map[string]string{
			"Authorization": "InvalidFormat token",
		})
		s.assertErrorResponse(invalidAuthResp, http.StatusUnauthorized, "invalid authorization header format")

		// Test with invalid token
		invalidTokenResp := s.makeRequest(endpoint.method, endpoint.path, nil, map[string]string{
			"Authorization": "Bearer invalid.jwt.token",
		})
		s.assertErrorResponse(invalidTokenResp, http.StatusUnauthorized, "invalid token")

		// Test with valid token (should proceed to actual handler)
		validResp := s.makeAuthenticatedRequest(endpoint.method, endpoint.path, nil)
		// Should not return 401, might return other status codes based on handler logic
		s.NotEqual(http.StatusUnauthorized, validResp.StatusCode)
		validResp.Body.Close()
	}
}

// TestCORSMiddlewareIntegration tests CORS middleware integration
func (s *MiddlewareIntegrationTestSuite) TestCORSMiddlewareIntegration() {
	// Test preflight request
	preflightHeaders := map[string]string{
		"Origin":                         "https://example.com",
		"Access-Control-Request-Method":  "POST",
		"Access-Control-Request-Headers": "Content-Type,Authorization",
	}

	preflightResp := s.makeRequest("OPTIONS", "/api/v1/auth/login", nil, preflightHeaders)
	s.Equal(http.StatusOK, preflightResp.StatusCode)

	// Check CORS headers
	s.NotEmpty(preflightResp.Header.Get("Access-Control-Allow-Origin"))
	s.NotEmpty(preflightResp.Header.Get("Access-Control-Allow-Methods"))
	s.NotEmpty(preflightResp.Header.Get("Access-Control-Allow-Headers"))

	preflightResp.Body.Close()

	// Test actual request with CORS headers
	corsHeaders := map[string]string{
		"Origin": "https://example.com",
	}

	loginReq := map[string]string{
		"email":    s.testUser.Email,
		"password": "password123",
	}

	corsResp := s.makeRequest("POST", "/api/v1/auth/login", loginReq, corsHeaders)
	s.Equal(http.StatusOK, corsResp.StatusCode)

	// Verify CORS headers are present in response
	s.NotEmpty(corsResp.Header.Get("Access-Control-Allow-Origin"))
	corsResp.Body.Close()
}

// TestLoggingMiddlewareIntegration tests logging middleware integration
func (s *MiddlewareIntegrationTestSuite) TestLoggingMiddlewareIntegration() {
	// This test verifies that requests are logged properly
	// In a real scenario, you might capture logs and verify their content
	
	// Make various requests to test logging
	requests := []struct {
		method string
		path   string
		body   interface{}
		auth   bool
	}{
		{"GET", "/health", nil, false},
		{"POST", "/api/v1/auth/login", map[string]string{"email": s.testUser.Email, "password": "password123"}, false},
		{"GET", "/api/v1/auth/profile", nil, true},
		{"GET", "/nonexistent", nil, false},
	}

	for _, req := range requests {
		var resp *http.Response
		if req.auth {
			resp = s.makeAuthenticatedRequest(req.method, req.path, req.body)
		} else {
			resp = s.makeRequest(req.method, req.path, req.body, nil)
		}
		
		// Verify request was processed (logging should be transparent)
		s.NotNil(resp)
		resp.Body.Close()
	}
}

// TestRateLimitingMiddlewareIntegration tests rate limiting middleware integration
func (s *MiddlewareIntegrationTestSuite) TestRateLimitingMiddlewareIntegration() {
	if s.cacheService == nil {
		s.T().Skip("Redis not available, skipping rate limiting tests")
		return
	}

	// Test rate limiting on login endpoint
	loginReq := map[string]string{
		"email":    "ratelimit-test@example.com",
		"password": "password123",
	}

	// Make rapid requests to trigger rate limiting
	var rateLimitTriggered bool
	var successfulRequests int

	for i := 0; i < 50; i++ {
		resp := s.makeRequest("POST", "/api/v1/auth/login", loginReq, nil)
		
		if resp.StatusCode == http.StatusTooManyRequests {
			rateLimitTriggered = true
			
			// Verify rate limit headers are present
			s.NotEmpty(resp.Header.Get("X-RateLimit-Limit"))
			s.NotEmpty(resp.Header.Get("X-RateLimit-Remaining"))
			s.NotEmpty(resp.Header.Get("Retry-After"))
			
			resp.Body.Close()
			break
		} else if resp.StatusCode == http.StatusUnauthorized {
			// Expected for invalid credentials
			successfulRequests++
		}
		
		resp.Body.Close()
		time.Sleep(10 * time.Millisecond)
	}

	// Log results for debugging
	s.T().Logf("Rate limiting triggered: %v, Successful requests before limit: %d", rateLimitTriggered, successfulRequests)
}

// TestMiddlewareChainOrder tests the order of middleware execution
func (s *MiddlewareIntegrationTestSuite) TestMiddlewareChainOrder() {
	// Test that CORS middleware runs before auth middleware
	// by sending a preflight request to a protected endpoint
	
	preflightHeaders := map[string]string{
		"Origin":                         "https://example.com",
		"Access-Control-Request-Method":  "GET",
		"Access-Control-Request-Headers": "Authorization",
	}

	preflightResp := s.makeRequest("OPTIONS", "/api/v1/auth/profile", nil, preflightHeaders)
	
	// Should return 200 (CORS handled) not 401 (auth middleware)
	s.Equal(http.StatusOK, preflightResp.StatusCode)
	s.NotEmpty(preflightResp.Header.Get("Access-Control-Allow-Origin"))
	preflightResp.Body.Close()

	// Test that auth middleware runs before business logic
	// by sending an unauthenticated request to a protected endpoint
	
	unauthedResp := s.makeRequest("GET", "/api/v1/auth/profile", nil, nil)
	
	// Should return 401 (auth middleware) not the business logic response
	s.Equal(http.StatusUnauthorized, unauthedResp.StatusCode)
	unauthedResp.Body.Close()
}

// TestContentTypeMiddlewareIntegration tests content type handling
func (s *MiddlewareIntegrationTestSuite) TestContentTypeMiddlewareIntegration() {
	// Test JSON content type handling
	jsonReq := map[string]string{
		"email":    s.testUser.Email,
		"password": "password123",
	}

	jsonHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	jsonResp := s.makeRequest("POST", "/api/v1/auth/login", jsonReq, jsonHeaders)
	s.Equal(http.StatusOK, jsonResp.StatusCode)
	jsonResp.Body.Close()

	// Test unsupported content type
	formHeaders := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	formResp := s.makeRequest("POST", "/api/v1/auth/login", "email=test&password=test", formHeaders)
	// Should handle gracefully (might return 400 or process as best effort)
	s.NotEqual(http.StatusOK, formResp.StatusCode)
	formResp.Body.Close()

	// Test missing content type (should default to JSON)
	noContentTypeResp := s.makeRequest("POST", "/api/v1/auth/login", jsonReq, nil)
	s.Equal(http.StatusOK, noContentTypeResp.StatusCode)
	noContentTypeResp.Body.Close()
}

// TestSecurityHeadersMiddleware tests security headers middleware
func (s *MiddlewareIntegrationTestSuite) TestSecurityHeadersMiddleware() {
	resp := s.makeRequest("GET", "/health", nil, nil)
	
	// Check for common security headers
	headers := resp.Header
	
	// These headers might be set by security middleware
	securityHeaders := []string{
		"X-Content-Type-Options",
		"X-Frame-Options", 
		"X-XSS-Protection",
		"Strict-Transport-Security",
		"Content-Security-Policy",
	}

	for _, header := range securityHeaders {
		value := headers.Get(header)
		if value != "" {
			s.T().Logf("Security header %s: %s", header, value)
		}
	}

	resp.Body.Close()
}

// TestErrorHandlingMiddleware tests error handling middleware
func (s *MiddlewareIntegrationTestSuite) TestErrorHandlingMiddleware() {
	// Test various error scenarios to ensure consistent error responses

	// 1. 404 Not Found
	notFoundResp := s.makeRequest("GET", "/api/v1/nonexistent", nil, nil)
	s.Equal(http.StatusNotFound, notFoundResp.StatusCode)
	
	var notFoundError map[string]interface{}
	s.parseJSONResponse(notFoundResp, &notFoundError)
	s.Contains(notFoundError, "error")

	// 2. 405 Method Not Allowed
	methodNotAllowedResp := s.makeRequest("PATCH", "/api/v1/auth/login", nil, nil)
	s.Equal(http.StatusMethodNotAllowed, methodNotAllowedResp.StatusCode)
	methodNotAllowedResp.Body.Close()

	// 3. 400 Bad Request (malformed JSON)
	malformedJSONResp := s.makeRequest("POST", "/api/v1/auth/login", "invalid json", map[string]string{
		"Content-Type": "application/json",
	})
	s.Equal(http.StatusBadRequest, malformedJSONResp.StatusCode)
	malformedJSONResp.Body.Close()

	// 4. 401 Unauthorized
	unauthorizedResp := s.makeRequest("GET", "/api/v1/auth/profile", nil, nil)
	s.Equal(http.StatusUnauthorized, unauthorizedResp.StatusCode)
	
	var unauthorizedError map[string]interface{}
	s.parseJSONResponse(unauthorizedResp, &unauthorizedError)
	s.Contains(unauthorizedError, "error")
}

// TestMiddlewareWithDifferentHTTPMethods tests middleware with various HTTP methods
func (s *MiddlewareIntegrationTestSuite) TestMiddlewareWithDifferentHTTPMethods() {
	methods := []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"}
	
	for _, method := range methods {
		// Test public endpoint
		var resp *http.Response
		
		if method == "HEAD" {
			resp = s.makeRequest(method, "/health", nil, nil)
		} else if method == "OPTIONS" {
			// Test CORS preflight
			resp = s.makeRequest(method, "/api/v1/auth/login", nil, map[string]string{
				"Origin": "https://example.com",
				"Access-Control-Request-Method": "POST",
			})
		} else if method == "GET" {
			resp = s.makeRequest(method, "/health", nil, nil)
		} else {
			// For other methods, test against a protected endpoint to verify auth middleware
			resp = s.makeRequest(method, "/api/v1/auth/profile", nil, nil)
		}
		
		// Verify response is received (middleware chain executed)
		s.NotNil(resp)
		
		// For protected endpoints without auth, should return 401
		if method != "GET" && method != "HEAD" && method != "OPTIONS" {
			s.Equal(http.StatusUnauthorized, resp.StatusCode)
		}
		
		resp.Body.Close()
	}
}

// TestMiddlewarePerformance tests middleware performance impact
func (s *MiddlewareIntegrationTestSuite) TestMiddlewarePerformance() {
	numRequests := 100
	
	// Test simple endpoint to measure middleware overhead
	start := time.Now()
	
	for i := 0; i < numRequests; i++ {
		resp := s.makeRequest("GET", "/health", nil, nil)
		resp.Body.Close()
	}
	
	duration := time.Since(start)
	avgLatency := duration / time.Duration(numRequests)
	
	s.T().Logf("Processed %d requests in %v (avg: %v per request)", numRequests, duration, avgLatency)
	
	// Performance assertion - middleware should not add significant overhead
	s.Less(avgLatency, 50*time.Millisecond, "Average request latency with middleware should be reasonable")
}

// TestMiddlewareWithConcurrentRequests tests middleware under concurrent load
func (s *MiddlewareIntegrationTestSuite) TestMiddlewareWithConcurrentRequests() {
	numConcurrent := 20
	done := make(chan bool, numConcurrent)
	
	for i := 0; i < numConcurrent; i++ {
		go func(index int) {
			defer func() { done <- true }()
			
			// Mix of authenticated and unauthenticated requests
			if index%2 == 0 {
				resp := s.makeRequest("GET", "/health", nil, nil)
				s.Equal(http.StatusOK, resp.StatusCode)
				resp.Body.Close()
			} else {
				resp := s.makeAuthenticatedRequest("GET", "/api/v1/auth/profile", nil)
				s.Equal(http.StatusOK, resp.StatusCode)
				resp.Body.Close()
			}
		}(i)
	}
	
	// Wait for all requests to complete
	for i := 0; i < numConcurrent; i++ {
		<-done
	}
}

// TestCustomHeadersHandling tests custom headers handling in middleware
func (s *MiddlewareIntegrationTestSuite) TestCustomHeadersHandling() {
	customHeaders := map[string]string{
		"X-Request-ID":     "test-request-123",
		"X-Client-Version": "1.0.0",
		"X-User-Agent":     "URLShortener-Test/1.0",
		"Accept":           "application/json",
		"Accept-Language":  "en-US",
		"Accept-Encoding":  "gzip, deflate",
	}
	
	resp := s.makeRequest("GET", "/health", nil, customHeaders)
	s.Equal(http.StatusOK, resp.StatusCode)
	
	// Some middleware might echo back certain headers or add response headers
	// based on request headers
	responseHeaders := resp.Header
	s.T().Logf("Response headers: %v", responseHeaders)
	
	resp.Body.Close()
}

// TestMiddlewareWithLargePayloads tests middleware with large request payloads
func (s *MiddlewareIntegrationTestSuite) TestMiddlewareWithLargePayloads() {
	// Create a large description
	largeDescription := strings.Repeat("This is a very long description. ", 1000)
	
	largeURLReq := map[string]interface{}{
		"originalUrl": "https://example.com/large-payload-test",
		"title":       "Large Payload Test",
		"description": largeDescription,
	}
	
	resp := s.makeAuthenticatedRequest("POST", "/api/v1/urls/", largeURLReq)
	
	// Should handle large payloads gracefully
	// Might return 413 if payload is too large, or 201 if accepted
	s.True(resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusRequestEntityTooLarge)
	resp.Body.Close()
}

// TestMiddlewareErrorRecovery tests middleware error recovery
func (s *MiddlewareIntegrationTestSuite) TestMiddlewareErrorRecovery() {
	// Test various scenarios that might cause middleware errors
	
	// 1. Extremely large headers
	largeHeaderValue := strings.Repeat("x", 8192)
	largeHeaderResp := s.makeRequest("GET", "/health", nil, map[string]string{
		"X-Large-Header": largeHeaderValue,
	})
	
	// Should either process successfully or return appropriate error
	s.True(largeHeaderResp.StatusCode == http.StatusOK || 
		   largeHeaderResp.StatusCode == http.StatusRequestHeaderFieldsTooLarge ||
		   largeHeaderResp.StatusCode == http.StatusBadRequest)
	largeHeaderResp.Body.Close()

	// 2. Invalid UTF-8 in headers (if applicable)
	// This test depends on how the server handles invalid UTF-8
	
	// 3. Multiple rapid requests from same IP (rate limiting stress test)
	if s.cacheService != nil {
		for i := 0; i < 10; i++ {
			resp := s.makeRequest("GET", "/health", nil, nil)
			resp.Body.Close()
		}
	}
}