package middleware

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/microcosm-cc/bluemonday"
	"github.com/sirupsen/logrus"
)

// SecurityConfig holds security middleware configuration
type SecurityConfig struct {
	MaxRequestSize       int64
	AllowedContentTypes  []string
	EnableXSSProtection  bool
	EnableSanitization   bool
	BlockSuspiciousHosts []string
	AllowedURLSchemes    []string
}

// DefaultSecurityConfig returns default security configuration
func DefaultSecurityConfig() *SecurityConfig {
	return &SecurityConfig{
		MaxRequestSize:      10 << 20, // 10MB
		AllowedContentTypes: []string{"application/json", "application/x-www-form-urlencoded", "multipart/form-data"},
		EnableXSSProtection: true,
		EnableSanitization:  true,
		BlockSuspiciousHosts: []string{
			"localhost", "127.0.0.1", "0.0.0.0", "::1",
			"10.", "172.", "192.168.", "169.254.",
		},
		AllowedURLSchemes: []string{"http", "https"},
	}
}

// SecurityHeadersMiddleware adds comprehensive security headers
func SecurityHeadersMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Content Security Policy
			csp := "default-src 'self'; " +
				"script-src 'self' 'unsafe-inline' 'unsafe-eval'; " +
				"style-src 'self' 'unsafe-inline'; " +
				"img-src 'self' data: https:; " +
				"font-src 'self' data:; " +
				"connect-src 'self'; " +
				"frame-ancestors 'none'; " +
				"base-uri 'self'; " +
				"form-action 'self'"
			
			w.Header().Set("Content-Security-Policy", csp)
			
			// Additional security headers
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
			
			// HSTS header (only for HTTPS)
			if r.TLS != nil {
				w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
			}
			
			// Remove server information
			w.Header().Set("Server", "")
			
			next.ServeHTTP(w, r)
		})
	}
}

// InputValidationMiddleware provides comprehensive input validation and sanitization
func InputValidationMiddleware(config *SecurityConfig) func(next http.Handler) http.Handler {
	if config == nil {
		config = DefaultSecurityConfig()
	}
	
	// Initialize HTML sanitizer
	policy := bluemonday.StrictPolicy()
	
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check request size
			if r.ContentLength > config.MaxRequestSize {
				logrus.WithFields(logrus.Fields{
					"content_length": r.ContentLength,
					"max_size":       config.MaxRequestSize,
					"ip":             middleware.GetReqID(r.Context()),
				}).Warn("Request size exceeded limit")
				
				http.Error(w, "Request entity too large", http.StatusRequestEntityTooLarge)
				return
			}
			
			// Validate content type for non-GET requests
			if r.Method != http.MethodGet && r.Method != http.MethodHead {
				contentType := r.Header.Get("Content-Type")
				if contentType != "" {
					// Extract base content type (ignore charset, boundary, etc.)
					baseContentType := strings.Split(contentType, ";")[0]
					baseContentType = strings.TrimSpace(baseContentType)
					
					allowed := false
					for _, allowedType := range config.AllowedContentTypes {
						if baseContentType == allowedType {
							allowed = true
							break
						}
					}
					
					if !allowed {
						logrus.WithFields(logrus.Fields{
							"content_type": contentType,
							"ip":           middleware.GetReqID(r.Context()),
						}).Warn("Invalid content type")
						
						http.Error(w, "Unsupported content type", http.StatusUnsupportedMediaType)
						return
					}
				}
			}
			
			// Sanitize query parameters
			if config.EnableSanitization {
				query := r.URL.Query()
				sanitized := false
				
				for key, values := range query {
					for i, value := range values {
						// HTML sanitization
						sanitizedValue := policy.Sanitize(value)
						
						// Additional XSS protection
						if config.EnableXSSProtection {
							sanitizedValue = html.EscapeString(sanitizedValue)
							sanitizedValue = sanitizeForXSS(sanitizedValue)
						}
						
						if sanitizedValue != value {
							query[key][i] = sanitizedValue
							sanitized = true
						}
					}
				}
				
				if sanitized {
					r.URL.RawQuery = query.Encode()
				}
			}
			
			// Validate and sanitize headers
			validateHeaders(r, config)
			
			next.ServeHTTP(w, r)
		})
	}
}

// URLValidationMiddleware validates URLs in requests
func URLValidationMiddleware(config *SecurityConfig) func(next http.Handler) http.Handler {
	if config == nil {
		config = DefaultSecurityConfig()
	}
	
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only validate POST/PUT requests that might contain URLs
			if r.Method == http.MethodPost || r.Method == http.MethodPut {
				contentType := r.Header.Get("Content-Type")
				
				// For JSON requests, parse and validate URLs
				if strings.Contains(contentType, "application/json") {
					if err := validateJSONURLs(r, config); err != nil {
						logrus.WithFields(logrus.Fields{
							"error": err.Error(),
							"ip":    middleware.GetReqID(r.Context()),
						}).Warn("Invalid URL in request")
						
						http.Error(w, "Invalid URL provided", http.StatusBadRequest)
						return
					}
				}
			}
			
			next.ServeHTTP(w, r)
		})
	}
}

// RequestSanitizationMiddleware sanitizes request body content
func RequestSanitizationMiddleware() func(next http.Handler) http.Handler {
	policy := bluemonday.StrictPolicy()
	
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only process JSON requests
			contentType := r.Header.Get("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				next.ServeHTTP(w, r)
				return
			}
			
			// Parse and sanitize JSON body
			if r.Body != nil && r.ContentLength > 0 {
				var body map[string]interface{}
				decoder := json.NewDecoder(r.Body)
				
				if err := decoder.Decode(&body); err == nil {
					// Sanitize string fields recursively
					sanitizeMapRecursive(body, policy)
					
					// Re-encode the sanitized body
					sanitizedBody, err := json.Marshal(body)
					if err == nil {
						r.Body = http.NoBody
						r.ContentLength = int64(len(sanitizedBody))
						r.Header.Set("Content-Length", fmt.Sprintf("%d", len(sanitizedBody)))
						
						// Create new request with sanitized body
						r.Body = &requestBody{data: sanitizedBody}
					}
				}
			}
			
			next.ServeHTTP(w, r)
		})
	}
}

// Helper functions

// sanitizeForXSS performs additional XSS sanitization
func sanitizeForXSS(input string) string {
	// Remove common XSS patterns
	xssPatterns := []string{
		`(?i)<script[^>]*>.*?</script>`,
		`(?i)<iframe[^>]*>.*?</iframe>`,
		`(?i)<object[^>]*>.*?</object>`,
		`(?i)<embed[^>]*>.*?</embed>`,
		`(?i)<link[^>]*>`,
		`(?i)<meta[^>]*>`,
		`(?i)javascript:`,
		`(?i)vbscript:`,
		`(?i)onload=`,
		`(?i)onerror=`,
		`(?i)onclick=`,
		`(?i)onmouseover=`,
	}
	
	result := input
	for _, pattern := range xssPatterns {
		re := regexp.MustCompile(pattern)
		result = re.ReplaceAllString(result, "")
	}
	
	return result
}

// validateHeaders checks for suspicious headers
func validateHeaders(r *http.Request, config *SecurityConfig) {
	suspiciousHeaders := []string{
		"X-Forwarded-Host",
		"X-Original-URL",
		"X-Rewrite-URL",
	}
	
	for _, header := range suspiciousHeaders {
		if value := r.Header.Get(header); value != "" {
			logrus.WithFields(logrus.Fields{
				"header": header,
				"value":  value,
				"ip":     middleware.GetReqID(r.Context()),
			}).Warn("Suspicious header detected")
		}
	}
}

// validateJSONURLs validates URLs in JSON request body
func validateJSONURLs(r *http.Request, config *SecurityConfig) error {
	var body map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	
	if err := decoder.Decode(&body); err != nil {
		return nil // If we can't parse it, let the handler deal with it
	}
	
	// Look for URL fields
	urlFields := []string{"url", "original_url", "destination_url", "redirect_url"}
	
	for _, field := range urlFields {
		if urlValue, exists := body[field]; exists {
			if urlStr, ok := urlValue.(string); ok {
				if err := validateURL(urlStr, config); err != nil {
					return err
				}
			}
		}
	}
	
	return nil
}

// validateURL performs comprehensive URL validation
func validateURL(urlStr string, config *SecurityConfig) error {
	if urlStr == "" {
		return fmt.Errorf("empty URL")
	}
	
	// Parse URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}
	
	// Check scheme
	schemeAllowed := false
	for _, scheme := range config.AllowedURLSchemes {
		if parsedURL.Scheme == scheme {
			schemeAllowed = true
			break
		}
	}
	
	if !schemeAllowed {
		return fmt.Errorf("unsupported URL scheme: %s", parsedURL.Scheme)
	}
	
	// Check for suspicious hosts
	hostname := strings.ToLower(parsedURL.Hostname())
	for _, suspicious := range config.BlockSuspiciousHosts {
		if strings.Contains(hostname, suspicious) {
			return fmt.Errorf("suspicious hostname detected: %s", hostname)
		}
	}
	
	// Check for localhost/private IP patterns
	if isPrivateOrLocalhost(hostname) {
		return fmt.Errorf("localhost/private IP not allowed: %s", hostname)
	}
	
	return nil
}

// isPrivateOrLocalhost checks if hostname is localhost or private IP
func isPrivateOrLocalhost(hostname string) bool {
	if hostname == "localhost" || hostname == "127.0.0.1" || hostname == "::1" {
		return true
	}
	
	// Check for private IP ranges
	privatePatterns := []string{
		"^10\\.",
		"^172\\.(1[6-9]|2[0-9]|3[01])\\.",
		"^192\\.168\\.",
		"^169\\.254\\.", // Link-local
	}
	
	for _, pattern := range privatePatterns {
		matched, _ := regexp.MatchString(pattern, hostname)
		if matched {
			return true
		}
	}
	
	return false
}

// sanitizeMapRecursive recursively sanitizes string values in a map
func sanitizeMapRecursive(data map[string]interface{}, policy *bluemonday.Policy) {
	for key, value := range data {
		switch v := value.(type) {
		case string:
			data[key] = policy.Sanitize(v)
		case map[string]interface{}:
			sanitizeMapRecursive(v, policy)
		case []interface{}:
			for i, item := range v {
				if str, ok := item.(string); ok {
					v[i] = policy.Sanitize(str)
				} else if nested, ok := item.(map[string]interface{}); ok {
					sanitizeMapRecursive(nested, policy)
				}
			}
		}
	}
}

// requestBody implements io.ReadCloser for modified request bodies
type requestBody struct {
	data []byte
	pos  int
}

func (rb *requestBody) Read(p []byte) (int, error) {
	if rb.pos >= len(rb.data) {
		return 0, fmt.Errorf("EOF")
	}
	
	n := copy(p, rb.data[rb.pos:])
	rb.pos += n
	return n, nil
}

func (rb *requestBody) Close() error {
	return nil
}

// HTTPSEnforcementMiddleware redirects HTTP to HTTPS in production
func HTTPSEnforcementMiddleware(enforceHTTPS bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if enforceHTTPS && r.TLS == nil {
				// Check for forwarded protocol headers (for reverse proxies)
				proto := r.Header.Get("X-Forwarded-Proto")
				if proto != "https" {
					// Redirect to HTTPS
					target := url.URL{
						Scheme:   "https",
						Host:     r.Host,
						Path:     r.URL.Path,
						RawQuery: r.URL.RawQuery,
					}
					
					logrus.WithFields(logrus.Fields{
						"original_url": r.URL.String(),
						"redirect_url": target.String(),
					}).Info("Redirecting HTTP to HTTPS")
					
					http.Redirect(w, r, target.String(), http.StatusMovedPermanently)
					return
				}
			}
			
			next.ServeHTTP(w, r)
		})
	}
}

// SecureCookieMiddleware configures secure cookie settings
func SecureCookieMiddleware(secureCookies bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Wrap response writer to modify cookies
			wrapped := &secureResponseWriter{
				ResponseWriter: w,
				secureCookies:  secureCookies,
				isHTTPS:       r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https",
			}
			
			next.ServeHTTP(wrapped, r)
		})
	}
}

// secureResponseWriter wraps http.ResponseWriter to modify cookies
type secureResponseWriter struct {
	http.ResponseWriter
	secureCookies bool
	isHTTPS       bool
}

func (w *secureResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *secureResponseWriter) Write(data []byte) (int, error) {
	// Modify Set-Cookie headers before writing
	w.modifyCookies()
	return w.ResponseWriter.Write(data)
}

func (w *secureResponseWriter) WriteHeader(statusCode int) {
	// Modify Set-Cookie headers before writing header
	w.modifyCookies()
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *secureResponseWriter) modifyCookies() {
	if !w.secureCookies {
		return
	}
	
	cookies := w.Header()["Set-Cookie"]
	for i, cookie := range cookies {
		// Parse existing cookie
		parts := strings.Split(cookie, ";")
		if len(parts) == 0 {
			continue
		}
		
		// Add security attributes
		var newParts []string
		newParts = append(newParts, parts[0]) // Keep the name=value part
		
		// Add security flags
		hasHttpOnly := false
		hasSecure := false
		hasSameSite := false
		
		// Check existing attributes
		for j := 1; j < len(parts); j++ {
			part := strings.TrimSpace(parts[j])
			lower := strings.ToLower(part)
			
			if strings.HasPrefix(lower, "httponly") {
				hasHttpOnly = true
			} else if strings.HasPrefix(lower, "secure") {
				hasSecure = true
			} else if strings.HasPrefix(lower, "samesite") {
				hasSameSite = true
			}
			
			newParts = append(newParts, part)
		}
		
		// Add missing security attributes
		if !hasHttpOnly {
			newParts = append(newParts, "HttpOnly")
		}
		
		if !hasSecure && w.isHTTPS {
			newParts = append(newParts, "Secure")
		}
		
		if !hasSameSite {
			newParts = append(newParts, "SameSite=Strict")
		}
		
		// Reconstruct cookie
		cookies[i] = strings.Join(newParts, "; ")
	}
	
	// Update header
	if len(cookies) > 0 {
		w.Header()["Set-Cookie"] = cookies
	}
}

// RateLimitByIPMiddleware provides additional IP-based rate limiting for security
func RateLimitByIPMiddleware(requestsPerMinute int) func(next http.Handler) http.Handler {
	type rateLimitEntry struct {
		count     int
		resetTime time.Time
	}
	
	ipLimits := make(map[string]*rateLimitEntry)
	
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			if forwardedIP := r.Header.Get("X-Forwarded-For"); forwardedIP != "" {
				ip = strings.Split(forwardedIP, ",")[0]
			}
			ip = strings.TrimSpace(ip)
			
			now := time.Now()
			
			// Clean up old entries
			for k, v := range ipLimits {
				if now.After(v.resetTime) {
					delete(ipLimits, k)
				}
			}
			
			// Check current IP
			entry, exists := ipLimits[ip]
			if !exists {
				entry = &rateLimitEntry{
					count:     1,
					resetTime: now.Add(time.Minute),
				}
				ipLimits[ip] = entry
			} else {
				if now.After(entry.resetTime) {
					entry.count = 1
					entry.resetTime = now.Add(time.Minute)
				} else {
					entry.count++
				}
			}
			
			if entry.count > requestsPerMinute {
				logrus.WithFields(logrus.Fields{
					"ip":    ip,
					"count": entry.count,
					"limit": requestsPerMinute,
				}).Warn("IP rate limit exceeded")
				
				w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", requestsPerMinute))
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", entry.resetTime.Unix()))
				
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			
			// Add rate limit headers
			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", requestsPerMinute))
			w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", requestsPerMinute-entry.count))
			w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", entry.resetTime.Unix()))
			
			next.ServeHTTP(w, r)
		})
	}
}