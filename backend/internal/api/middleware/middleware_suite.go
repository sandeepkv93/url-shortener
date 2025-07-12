package middleware

import (
	"net/http"
	
	"github.com/go-chi/chi/v5"
	"url-shortener/internal/core/ports"
)

// SecurityMiddlewareSuite provides a comprehensive security middleware stack
type SecurityMiddlewareSuite struct {
	config             *SecurityConfig
	cacheService       ports.CacheService
	enableHTTPS        bool
	enableSecureCookies bool
}

// NewSecurityMiddlewareSuite creates a new security middleware suite
func NewSecurityMiddlewareSuite(cacheService ports.CacheService, enableHTTPS bool) *SecurityMiddlewareSuite {
	return &SecurityMiddlewareSuite{
		config:             DefaultSecurityConfig(),
		cacheService:       cacheService,
		enableHTTPS:        enableHTTPS,
		enableSecureCookies: enableHTTPS,
	}
}

// Apply applies all security middleware to a router
func (s *SecurityMiddlewareSuite) Apply(r chi.Router) {
	// Security headers (should be applied early)
	r.Use(SecurityHeadersMiddleware())
	
	// HTTPS enforcement (if enabled)
	if s.enableHTTPS {
		r.Use(HTTPSEnforcementMiddleware(true))
	}
	
	// Secure cookie settings
	r.Use(SecureCookieMiddleware(s.enableSecureCookies))
	
	// Input validation and sanitization
	r.Use(InputValidationMiddleware(s.config))
	
	// URL validation for URL-related endpoints
	r.Use(URLValidationMiddleware(s.config))
	
	// Request sanitization
	r.Use(RequestSanitizationMiddleware())
	
	// Additional IP-based rate limiting for security
	if s.cacheService != nil {
		r.Use(RateLimitByIPMiddleware(100)) // 100 requests per minute per IP
	}
}

// ApplyToEndpoint applies security middleware to specific endpoints
func (s *SecurityMiddlewareSuite) ApplyToEndpoint(handler http.HandlerFunc) http.HandlerFunc {
	// Create a chain of middleware
	return func(w http.ResponseWriter, r *http.Request) {
		// Apply security headers
		SecurityHeadersMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Apply input validation
			InputValidationMiddleware(s.config)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Apply the actual handler
				handler(w, r)
			})).ServeHTTP(w, r)
		})).ServeHTTP(w, r)
	}
}

// WithConfig allows customizing the security configuration
func (s *SecurityMiddlewareSuite) WithConfig(config *SecurityConfig) *SecurityMiddlewareSuite {
	s.config = config
	return s
}

// WithHTTPS enables or disables HTTPS enforcement
func (s *SecurityMiddlewareSuite) WithHTTPS(enable bool) *SecurityMiddlewareSuite {
	s.enableHTTPS = enable
	s.enableSecureCookies = enable
	return s
}

// WithSecureCookies enables or disables secure cookie settings
func (s *SecurityMiddlewareSuite) WithSecureCookies(enable bool) *SecurityMiddlewareSuite {
	s.enableSecureCookies = enable
	return s
}