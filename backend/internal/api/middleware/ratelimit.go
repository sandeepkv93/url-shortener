package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"url-shortener/internal/core/ports"
)

type RateLimitConfig struct {
	// Requests per window
	RequestsPerWindow int64
	// Window duration
	WindowDuration time.Duration
	// Custom key generator function
	KeyGenerator func(*http.Request) string
	// Skip rate limiting for certain paths
	SkipPaths []string
	// Custom response when rate limit is exceeded
	OnRateLimitExceeded func(http.ResponseWriter, *http.Request)
}

type RateLimitMiddleware struct {
	cache  ports.CacheService
	config *RateLimitConfig
}

func NewRateLimitMiddleware(cache ports.CacheService, config *RateLimitConfig) *RateLimitMiddleware {
	if config == nil {
		config = &RateLimitConfig{
			RequestsPerWindow: 100,
			WindowDuration:    time.Minute,
		}
	}
	
	// Set default key generator if not provided
	if config.KeyGenerator == nil {
		config.KeyGenerator = defaultKeyGenerator
	}
	
	// Set default rate limit exceeded handler if not provided
	if config.OnRateLimitExceeded == nil {
		config.OnRateLimitExceeded = defaultRateLimitExceededHandler
	}
	
	return &RateLimitMiddleware{
		cache:  cache,
		config: config,
	}
}

func (m *RateLimitMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip rate limiting for certain paths
		if m.shouldSkipPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Generate rate limit key
		key := m.config.KeyGenerator(r)
		
		// Check if rate limited
		isLimited, err := m.cache.IsRateLimited(r.Context(), key, m.config.RequestsPerWindow, m.config.WindowDuration)
		if err != nil {
			// If cache is unavailable, log error but allow request through
			// This prevents cache failures from bringing down the service
			next.ServeHTTP(w, r)
			return
		}
		
		if isLimited {
			m.config.OnRateLimitExceeded(w, r)
			return
		}
		
		// Increment rate limit counter
		count, err := m.cache.IncrementRateLimit(r.Context(), key, m.config.WindowDuration)
		if err != nil {
			// If we can't increment, log but continue (fail open)
			next.ServeHTTP(w, r)
			return
		}
		
		// Add rate limit headers
		m.addRateLimitHeaders(w, count)
		
		next.ServeHTTP(w, r)
	})
}

func (m *RateLimitMiddleware) shouldSkipPath(path string) bool {
	for _, skipPath := range m.config.SkipPaths {
		if path == skipPath {
			return true
		}
	}
	return false
}

func (m *RateLimitMiddleware) addRateLimitHeaders(w http.ResponseWriter, currentCount int64) {
	w.Header().Set("X-RateLimit-Limit", strconv.FormatInt(m.config.RequestsPerWindow, 10))
	w.Header().Set("X-RateLimit-Remaining", strconv.FormatInt(m.config.RequestsPerWindow-currentCount, 10))
	w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(m.config.WindowDuration).Unix(), 10))
}

// Default key generator uses IP address
func defaultKeyGenerator(r *http.Request) string {
	// Try to get real IP from headers (in case of proxy/load balancer)
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	
	return fmt.Sprintf("rate_limit:%s", ip)
}

// Default rate limit exceeded handler
func defaultRateLimitExceededHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Retry-After", "60")
	w.WriteHeader(http.StatusTooManyRequests)
	w.Write([]byte(`{"error": "Rate limit exceeded", "retry_after": 60}`))
}

// Key generator that uses user ID if authenticated, otherwise falls back to IP
func UserOrIPKeyGenerator(r *http.Request) string {
	// Check if user is authenticated
	if userID := GetUserIDFromContext(r.Context()); userID != 0 {
		return fmt.Sprintf("rate_limit:user:%d", userID)
	}
	
	// Fall back to IP-based rate limiting
	return defaultKeyGenerator(r)
}

// Key generator for API endpoints that need higher limits for authenticated users
func APIKeyGenerator(authenticatedLimit, unauthenticatedLimit int64) func(*http.Request) string {
	return func(r *http.Request) string {
		if userID := GetUserIDFromContext(r.Context()); userID != 0 {
			return fmt.Sprintf("api_rate_limit:user:%d:%d", userID, authenticatedLimit)
		}
		return fmt.Sprintf("api_rate_limit:ip:%s:%d", getClientIP(r), unauthenticatedLimit)
	}
}

func getClientIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}

// Predefined rate limit configurations
func GlobalRateLimit(cache ports.CacheService) func(http.Handler) http.Handler {
	config := &RateLimitConfig{
		RequestsPerWindow: 1000,
		WindowDuration:    time.Minute,
		KeyGenerator:      defaultKeyGenerator,
		SkipPaths: []string{
			"/health",
			"/metrics",
		},
	}
	
	middleware := NewRateLimitMiddleware(cache, config)
	return middleware.Handler
}

func APIRateLimit(cache ports.CacheService) func(http.Handler) http.Handler {
	config := &RateLimitConfig{
		RequestsPerWindow: 100,
		WindowDuration:    time.Minute,
		KeyGenerator:      UserOrIPKeyGenerator,
	}
	
	middleware := NewRateLimitMiddleware(cache, config)
	return middleware.Handler
}

func AuthRateLimit(cache ports.CacheService) func(http.Handler) http.Handler {
	config := &RateLimitConfig{
		RequestsPerWindow: 5,
		WindowDuration:    time.Minute,
		KeyGenerator:      defaultKeyGenerator,
		OnRateLimitExceeded: func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", "300") // 5 minutes
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error": "Too many authentication attempts", "retry_after": 300}`))
		},
	}
	
	middleware := NewRateLimitMiddleware(cache, config)
	return middleware.Handler
}

func URLCreationRateLimit(cache ports.CacheService) func(http.Handler) http.Handler {
	config := &RateLimitConfig{
		RequestsPerWindow: 50,
		WindowDuration:    time.Hour,
		KeyGenerator:      UserOrIPKeyGenerator,
		OnRateLimitExceeded: func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", "3600") // 1 hour
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error": "URL creation limit exceeded", "retry_after": 3600}`))
		},
	}
	
	middleware := NewRateLimitMiddleware(cache, config)
	return middleware.Handler
}