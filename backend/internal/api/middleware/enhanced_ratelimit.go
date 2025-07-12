package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"url-shortener/internal/core/ports"
)

// EndpointRateLimitConfig defines rate limiting rules for specific endpoints
type EndpointRateLimitConfig struct {
	// Endpoint pattern (supports wildcards)
	Pattern string
	// HTTP methods to apply rate limiting to
	Methods []string
	// Requests per minute for this endpoint
	RequestsPerMinute int
	// Requests per hour for this endpoint
	RequestsPerHour int
	// Requests per day for this endpoint
	RequestsPerDay int
	// Burst allowance (short-term burst of requests)
	BurstLimit int
	// Custom key generator for this endpoint
	KeyGenerator func(*http.Request) string
	// Custom message when rate limit is exceeded
	LimitExceededMessage string
}

// EnhancedRateLimitConfig provides comprehensive rate limiting configuration
type EnhancedRateLimitConfig struct {
	// Global rate limiting (applies to all endpoints)
	GlobalLimits struct {
		RequestsPerMinute int
		RequestsPerHour   int
		RequestsPerDay    int
		BurstLimit        int
	}
	
	// Per-endpoint rate limiting
	EndpointLimits []EndpointRateLimitConfig
	
	// Per-user rate limiting (requires authentication)
	UserLimits struct {
		RequestsPerMinute int
		RequestsPerHour   int
		RequestsPerDay    int
		URLCreationPerDay int
	}
	
	// IP-based rate limiting
	IPLimits struct {
		RequestsPerMinute int
		RequestsPerHour   int
		MaxConcurrent     int
	}
	
	// Whitelist/Blacklist
	WhitelistedIPs []string
	BlacklistedIPs []string
	
	// Advanced features
	EnableAdaptiveRateLimit bool  // Adjust limits based on system load
	EnableGeoRateLimit      bool  // Different limits for different countries
	EnableHeaderInspection  bool  // Rate limit based on specific headers
	
	// Monitoring and alerts
	EnableRateLimitLogging bool
	AlertThresholds        struct {
		HighTrafficPercent    float64 // Alert when traffic exceeds this % of limit
		SuspiciousActivityIPs int     // Alert when this many IPs hit rate limits
	}
}

// EnhancedRateLimitMiddleware provides advanced rate limiting functionality
type EnhancedRateLimitMiddleware struct {
	cache            ports.CacheService
	config           *EnhancedRateLimitConfig
	logger           *logrus.Logger
	concurrentTracker map[string]int // Track concurrent requests per IP
}

// NewEnhancedRateLimitMiddleware creates a new enhanced rate limiting middleware
func NewEnhancedRateLimitMiddleware(cache ports.CacheService, config *EnhancedRateLimitConfig, logger *logrus.Logger) *EnhancedRateLimitMiddleware {
	if config == nil {
		config = DefaultEnhancedRateLimitConfig()
	}
	
	return &EnhancedRateLimitMiddleware{
		cache:             cache,
		config:            config,
		logger:            logger,
		concurrentTracker: make(map[string]int),
	}
}

// DefaultEnhancedRateLimitConfig returns default enhanced rate limiting configuration
func DefaultEnhancedRateLimitConfig() *EnhancedRateLimitConfig {
	return &EnhancedRateLimitConfig{
		GlobalLimits: struct {
			RequestsPerMinute int
			RequestsPerHour   int
			RequestsPerDay    int
			BurstLimit        int
		}{
			RequestsPerMinute: 60,
			RequestsPerHour:   3600,
			RequestsPerDay:    86400,
			BurstLimit:        10,
		},
		EndpointLimits: []EndpointRateLimitConfig{
			{
				Pattern:           "/api/v1/auth/login",
				Methods:           []string{"POST"},
				RequestsPerMinute: 5,
				RequestsPerHour:   30,
				BurstLimit:        2,
				LimitExceededMessage: "Too many login attempts. Please try again later.",
			},
			{
				Pattern:           "/api/v1/auth/register",
				Methods:           []string{"POST"},
				RequestsPerMinute: 3,
				RequestsPerHour:   10,
				BurstLimit:        1,
				LimitExceededMessage: "Too many registration attempts. Please try again later.",
			},
			{
				Pattern:           "/api/v1/urls",
				Methods:           []string{"POST"},
				RequestsPerMinute: 20,
				RequestsPerHour:   200,
				RequestsPerDay:    1000,
				BurstLimit:        5,
				LimitExceededMessage: "URL creation rate limit exceeded. Please slow down.",
			},
			{
				Pattern:           "/api/v1/auth/forgot-password",
				Methods:           []string{"POST"},
				RequestsPerMinute: 2,
				RequestsPerHour:   5,
				BurstLimit:        1,
				LimitExceededMessage: "Too many password reset requests. Please try again later.",
			},
		},
		UserLimits: struct {
			RequestsPerMinute int
			RequestsPerHour   int
			RequestsPerDay    int
			URLCreationPerDay int
		}{
			RequestsPerMinute: 120,
			RequestsPerHour:   7200,
			RequestsPerDay:    172800,
			URLCreationPerDay: 500,
		},
		IPLimits: struct {
			RequestsPerMinute int
			RequestsPerHour   int
			MaxConcurrent     int
		}{
			RequestsPerMinute: 100,
			RequestsPerHour:   6000,
			MaxConcurrent:     10,
		},
		WhitelistedIPs:          []string{},
		BlacklistedIPs:          []string{},
		EnableAdaptiveRateLimit: true,
		EnableGeoRateLimit:      false,
		EnableHeaderInspection:  true,
		EnableRateLimitLogging:  true,
		AlertThresholds: struct {
			HighTrafficPercent    float64
			SuspiciousActivityIPs int
		}{
			HighTrafficPercent:    80.0,
			SuspiciousActivityIPs: 100,
		},
	}
}

// Handler returns the HTTP middleware handler
func (m *EnhancedRateLimitMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if IP is blacklisted
		clientIP := m.getClientIP(r)
		if m.isBlacklisted(clientIP) {
			m.logRateLimit("blacklisted_ip", clientIP, r.URL.Path, r.Method)
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}
		
		// Check if IP is whitelisted (skip rate limiting)
		if m.isWhitelisted(clientIP) {
			next.ServeHTTP(w, r)
			return
		}
		
		// Check concurrent requests limit
		if m.config.IPLimits.MaxConcurrent > 0 {
			if m.checkConcurrentLimit(clientIP) {
				m.logRateLimit("concurrent_limit", clientIP, r.URL.Path, r.Method)
				http.Error(w, "Too many concurrent requests", http.StatusTooManyRequests)
				return
			}
		}
		
		// Track concurrent request
		m.incrementConcurrent(clientIP)
		defer m.decrementConcurrent(clientIP)
		
		// Check endpoint-specific rate limits
		if exceeded, message := m.checkEndpointRateLimit(r); exceeded {
			m.logRateLimit("endpoint_limit", clientIP, r.URL.Path, r.Method)
			http.Error(w, message, http.StatusTooManyRequests)
			return
		}
		
		// Check user-specific rate limits (if authenticated)
		if userID := m.getUserID(r); userID != "" {
			if exceeded, message := m.checkUserRateLimit(r, userID); exceeded {
				m.logRateLimit("user_limit", clientIP, r.URL.Path, r.Method)
				http.Error(w, message, http.StatusTooManyRequests)
				return
			}
		}
		
		// Check IP-based rate limits
		if exceeded, message := m.checkIPRateLimit(r, clientIP); exceeded {
			m.logRateLimit("ip_limit", clientIP, r.URL.Path, r.Method)
			http.Error(w, message, http.StatusTooManyRequests)
			return
		}
		
		// Check global rate limits
		if exceeded, message := m.checkGlobalRateLimit(r, clientIP); exceeded {
			m.logRateLimit("global_limit", clientIP, r.URL.Path, r.Method)
			http.Error(w, message, http.StatusTooManyRequests)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// checkEndpointRateLimit checks rate limits for specific endpoints
func (m *EnhancedRateLimitMiddleware) checkEndpointRateLimit(r *http.Request) (bool, string) {
	path := r.URL.Path
	method := r.Method
	
	for _, endpointConfig := range m.config.EndpointLimits {
		// Check if path matches pattern
		if m.matchesPattern(path, endpointConfig.Pattern) {
			// Check if method matches
			if len(endpointConfig.Methods) > 0 && !m.containsMethod(endpointConfig.Methods, method) {
				continue
			}
			
			clientIP := m.getClientIP(r)
			keyGenerator := endpointConfig.KeyGenerator
			if keyGenerator == nil {
				keyGenerator = func(r *http.Request) string {
					return fmt.Sprintf("endpoint:%s:%s:%s", endpointConfig.Pattern, clientIP, method)
				}
			}
			
			key := keyGenerator(r)
			
			// Check minute limit
			if endpointConfig.RequestsPerMinute > 0 {
				if exceeded := m.checkRateLimit(key+":minute", endpointConfig.RequestsPerMinute, time.Minute); exceeded {
					message := endpointConfig.LimitExceededMessage
					if message == "" {
						message = "Rate limit exceeded for this endpoint"
					}
					return true, message
				}
			}
			
			// Check hour limit
			if endpointConfig.RequestsPerHour > 0 {
				if exceeded := m.checkRateLimit(key+":hour", endpointConfig.RequestsPerHour, time.Hour); exceeded {
					message := endpointConfig.LimitExceededMessage
					if message == "" {
						message = "Hourly rate limit exceeded for this endpoint"
					}
					return true, message
				}
			}
			
			// Check day limit
			if endpointConfig.RequestsPerDay > 0 {
				if exceeded := m.checkRateLimit(key+":day", endpointConfig.RequestsPerDay, 24*time.Hour); exceeded {
					message := endpointConfig.LimitExceededMessage
					if message == "" {
						message = "Daily rate limit exceeded for this endpoint"
					}
					return true, message
				}
			}
			
			// Endpoint-specific limit found and checked
			break
		}
	}
	
	return false, ""
}

// checkUserRateLimit checks rate limits for authenticated users
func (m *EnhancedRateLimitMiddleware) checkUserRateLimit(r *http.Request, userID string) (bool, string) {
	baseKey := fmt.Sprintf("user:%s", userID)
	
	// Check minute limit
	if m.config.UserLimits.RequestsPerMinute > 0 {
		if exceeded := m.checkRateLimit(baseKey+":minute", m.config.UserLimits.RequestsPerMinute, time.Minute); exceeded {
			return true, "User rate limit exceeded (per minute)"
		}
	}
	
	// Check hour limit
	if m.config.UserLimits.RequestsPerHour > 0 {
		if exceeded := m.checkRateLimit(baseKey+":hour", m.config.UserLimits.RequestsPerHour, time.Hour); exceeded {
			return true, "User rate limit exceeded (per hour)"
		}
	}
	
	// Check day limit
	if m.config.UserLimits.RequestsPerDay > 0 {
		if exceeded := m.checkRateLimit(baseKey+":day", m.config.UserLimits.RequestsPerDay, 24*time.Hour); exceeded {
			return true, "User rate limit exceeded (per day)"
		}
	}
	
	// Check URL creation limit for POST /urls endpoint
	if r.Method == "POST" && strings.Contains(r.URL.Path, "/urls") {
		if m.config.UserLimits.URLCreationPerDay > 0 {
			if exceeded := m.checkRateLimit(baseKey+":urls:day", m.config.UserLimits.URLCreationPerDay, 24*time.Hour); exceeded {
				return true, "Daily URL creation limit exceeded"
			}
		}
	}
	
	return false, ""
}

// checkIPRateLimit checks rate limits based on IP address
func (m *EnhancedRateLimitMiddleware) checkIPRateLimit(r *http.Request, clientIP string) (bool, string) {
	baseKey := fmt.Sprintf("ip:%s", clientIP)
	
	// Check minute limit
	if m.config.IPLimits.RequestsPerMinute > 0 {
		if exceeded := m.checkRateLimit(baseKey+":minute", m.config.IPLimits.RequestsPerMinute, time.Minute); exceeded {
			return true, "IP rate limit exceeded (per minute)"
		}
	}
	
	// Check hour limit
	if m.config.IPLimits.RequestsPerHour > 0 {
		if exceeded := m.checkRateLimit(baseKey+":hour", m.config.IPLimits.RequestsPerHour, time.Hour); exceeded {
			return true, "IP rate limit exceeded (per hour)"
		}
	}
	
	return false, ""
}

// checkGlobalRateLimit checks global rate limits
func (m *EnhancedRateLimitMiddleware) checkGlobalRateLimit(r *http.Request, clientIP string) (bool, string) {
	baseKey := fmt.Sprintf("global:%s", clientIP)
	
	// Check minute limit
	if m.config.GlobalLimits.RequestsPerMinute > 0 {
		if exceeded := m.checkRateLimit(baseKey+":minute", m.config.GlobalLimits.RequestsPerMinute, time.Minute); exceeded {
			return true, "Global rate limit exceeded (per minute)"
		}
	}
	
	// Check hour limit
	if m.config.GlobalLimits.RequestsPerHour > 0 {
		if exceeded := m.checkRateLimit(baseKey+":hour", m.config.GlobalLimits.RequestsPerHour, time.Hour); exceeded {
			return true, "Global rate limit exceeded (per hour)"
		}
	}
	
	// Check day limit
	if m.config.GlobalLimits.RequestsPerDay > 0 {
		if exceeded := m.checkRateLimit(baseKey+":day", m.config.GlobalLimits.RequestsPerDay, 24*time.Hour); exceeded {
			return true, "Global rate limit exceeded (per day)"
		}
	}
	
	return false, ""
}

// checkRateLimit checks if a specific rate limit has been exceeded
func (m *EnhancedRateLimitMiddleware) checkRateLimit(key string, limit int, window time.Duration) bool {
	// Increment counter
	count, err := m.cache.IncrBy(context.Background(), key, 1)
	if err != nil {
		// If cache is unavailable, allow the request (fail open)
		if m.logger != nil {
			m.logger.WithError(err).Warn("Rate limit cache error")
		}
		return false
	}
	
	// Set expiration on first increment
	if count == 1 {
		// This is a simplified approach - in production, you'd want atomic operations
		go m.cache.Set(context.Background(), key, count, window)
	}
	
	return count > int64(limit)
}

// checkConcurrentLimit checks if concurrent request limit is exceeded
func (m *EnhancedRateLimitMiddleware) checkConcurrentLimit(clientIP string) bool {
	count, exists := m.concurrentTracker[clientIP]
	if !exists {
		count = 0
	}
	return count >= m.config.IPLimits.MaxConcurrent
}

// incrementConcurrent increments concurrent request counter
func (m *EnhancedRateLimitMiddleware) incrementConcurrent(clientIP string) {
	m.concurrentTracker[clientIP]++
}

// decrementConcurrent decrements concurrent request counter
func (m *EnhancedRateLimitMiddleware) decrementConcurrent(clientIP string) {
	if count, exists := m.concurrentTracker[clientIP]; exists {
		if count <= 1 {
			delete(m.concurrentTracker, clientIP)
		} else {
			m.concurrentTracker[clientIP]--
		}
	}
}

// Helper functions

func (m *EnhancedRateLimitMiddleware) getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// Take the first IP in the chain
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}
	
	// Check X-Real-IP header
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return strings.TrimSpace(xri)
	}
	
	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	if strings.Contains(ip, ":") {
		ip = strings.Split(ip, ":")[0]
	}
	
	return ip
}

func (m *EnhancedRateLimitMiddleware) getUserID(r *http.Request) string {
	// Try to extract user ID from context (set by auth middleware)
	if userID := r.Context().Value("user_id"); userID != nil {
		if id, ok := userID.(string); ok {
			return id
		}
		if id, ok := userID.(uint); ok {
			return strconv.FormatUint(uint64(id), 10)
		}
	}
	return ""
}

func (m *EnhancedRateLimitMiddleware) isWhitelisted(ip string) bool {
	for _, whitelistedIP := range m.config.WhitelistedIPs {
		if ip == whitelistedIP {
			return true
		}
	}
	return false
}

func (m *EnhancedRateLimitMiddleware) isBlacklisted(ip string) bool {
	for _, blacklistedIP := range m.config.BlacklistedIPs {
		if ip == blacklistedIP {
			return true
		}
	}
	return false
}

func (m *EnhancedRateLimitMiddleware) matchesPattern(path, pattern string) bool {
	// Simple pattern matching - can be enhanced with regex or chi's route matching
	if pattern == path {
		return true
	}
	
	// Support for Chi-style patterns (simplified)
	if strings.Contains(pattern, "{") {
		// Basic wildcard matching - in production, you'd use proper route matching
		// For now, just check if the base path matches
		patternParts := strings.Split(pattern, "/")
		pathParts := strings.Split(path, "/")
		
		if len(patternParts) != len(pathParts) {
			return false
		}
		
		for i, part := range patternParts {
			if !strings.HasPrefix(part, "{") && part != pathParts[i] {
				return false
			}
		}
		return true
	}
	
	return false
}

func (m *EnhancedRateLimitMiddleware) containsMethod(methods []string, method string) bool {
	for _, m := range methods {
		if m == method {
			return true
		}
	}
	return false
}

func (m *EnhancedRateLimitMiddleware) logRateLimit(limitType, clientIP, path, method string) {
	if m.config.EnableRateLimitLogging && m.logger != nil {
		m.logger.WithFields(logrus.Fields{
			"limit_type": limitType,
			"client_ip":  clientIP,
			"path":       path,
			"method":     method,
			"timestamp":  time.Now(),
		}).Warn("Rate limit exceeded")
	}
}