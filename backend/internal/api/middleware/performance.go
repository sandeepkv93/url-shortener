package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

// CompressionMiddleware provides gzip compression for responses
func CompressionMiddleware(level int) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if client accepts gzip encoding
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				next.ServeHTTP(w, r)
				return
			}

			// Skip compression for small responses or specific content types
			if shouldSkipCompression(r) {
				next.ServeHTTP(w, r)
				return
			}

			// Create gzip writer
			gz, err := gzip.NewWriterLevel(w, level)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			defer gz.Close()

			// Set compression headers
			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Set("Vary", "Accept-Encoding")

			// Wrap response writer
			grw := &gzipResponseWriter{
				ResponseWriter: w,
				Writer:         gz,
			}

			next.ServeHTTP(grw, r)
		})
	}
}

// PerformanceMonitoringMiddleware tracks response times and metrics
func PerformanceMonitoringMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap response writer to capture metrics
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// Process request
			next.ServeHTTP(ww, r)

			// Calculate metrics
			duration := time.Since(start)
			status := ww.Status()
			size := ww.BytesWritten()

			// Log performance metrics
			logPerformanceMetrics(r, status, duration, size)

			// Add performance headers for debugging
			if isDevelopmentMode() {
				w.Header().Set("X-Response-Time", duration.String())
				w.Header().Set("X-Response-Size", strconv.Itoa(size))
			}
		})
	}
}

// CacheHeadersMiddleware adds appropriate HTTP caching headers
func CacheHeadersMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set cache headers based on endpoint type
			setCacheHeaders(w, r)
			next.ServeHTTP(w, r)
		})
	}
}

// ETagMiddleware adds ETag support for conditional requests
func ETagMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// For GET requests, check If-None-Match header
			if r.Method == http.MethodGet {
				ifNoneMatch := r.Header.Get("If-None-Match")
				if ifNoneMatch != "" {
					// Generate ETag based on request (simplified)
					etag := generateETag(r)
					if ifNoneMatch == etag {
						w.WriteHeader(http.StatusNotModified)
						return
					}
					w.Header().Set("ETag", etag)
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RequestSizeLimitMiddleware limits request body size
func RequestSizeLimitMiddleware(maxSize int64) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength > maxSize {
				http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
				return
			}

			// Limit reader to prevent memory exhaustion
			r.Body = http.MaxBytesReader(w, r.Body, maxSize)
			next.ServeHTTP(w, r)
		})
	}
}

// Helper types and functions

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (grw *gzipResponseWriter) Write(data []byte) (int, error) {
	return grw.Writer.Write(data)
}

func shouldSkipCompression(r *http.Request) bool {
	// Skip compression for certain content types
	contentType := r.Header.Get("Content-Type")
	skipTypes := []string{
		"image/", "video/", "audio/",
		"application/zip", "application/gzip",
		"application/pdf", "application/octet-stream",
	}

	for _, skipType := range skipTypes {
		if strings.HasPrefix(contentType, skipType) {
			return true
		}
	}

	// Skip for small responses (they don't benefit from compression)
	contentLength := r.Header.Get("Content-Length")
	if contentLength != "" {
		if length, err := strconv.Atoi(contentLength); err == nil && length < 1024 {
			return true
		}
	}

	return false
}

func logPerformanceMetrics(r *http.Request, status int, duration time.Duration, size int) {
	// Create structured log entry
	fields := logrus.Fields{
		"method":        r.Method,
		"path":          r.URL.Path,
		"status":        status,
		"duration_ms":   duration.Milliseconds(),
		"response_size": size,
		"user_agent":    r.Header.Get("User-Agent"),
		"remote_addr":   r.RemoteAddr,
	}

	// Add query parameters for analytics
	if r.URL.RawQuery != "" {
		fields["query"] = r.URL.RawQuery
	}

	// Log based on response status and duration
	logger := logrus.WithFields(fields)
	
	switch {
	case status >= 500:
		logger.Error("Server error response")
	case status >= 400:
		logger.Warn("Client error response")
	case duration > 1*time.Second:
		logger.Warn("Slow response detected")
	case duration > 500*time.Millisecond:
		logger.Info("Response time warning")
	default:
		logger.Debug("Request completed successfully")
	}
}

func setCacheHeaders(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	switch {
	case strings.HasPrefix(path, "/api/v1/auth"):
		// Auth endpoints - no cache
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

	case strings.HasPrefix(path, "/api/v1/analytics"):
		// Analytics - short cache
		w.Header().Set("Cache-Control", "private, max-age=300") // 5 minutes

	case strings.HasPrefix(path, "/api/v1/qr"):
		// QR codes - long cache (they don't change)
		w.Header().Set("Cache-Control", "public, max-age=86400") // 24 hours

	case strings.HasPrefix(path, "/api/v1/urls") && r.Method == http.MethodGet:
		// URL lists - medium cache
		w.Header().Set("Cache-Control", "private, max-age=600") // 10 minutes

	case path == "/health":
		// Health check - no cache
		w.Header().Set("Cache-Control", "no-cache")

	default:
		// Default - short cache
		w.Header().Set("Cache-Control", "private, max-age=60") // 1 minute
	}

	// Add security headers
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
}

func generateETag(r *http.Request) string {
	// Simple ETag generation based on path and query
	// In production, this should be based on actual content
	return `"` + hashString(r.URL.Path+r.URL.RawQuery) + `"`
}

func hashString(s string) string {
	// Simple hash function - in production use crypto/sha256
	h := uint32(2166136261)
	for _, c := range []byte(s) {
		h ^= uint32(c)
		h *= 16777619
	}
	return strconv.FormatUint(uint64(h), 36)
}

func isDevelopmentMode() bool {
	// Check if we're in development mode
	// This should be configured via environment variable
	return true // Simplified for now
}

// Performance monitoring utilities

// ResponseTimeTracker tracks response times for specific endpoints
type ResponseTimeTracker struct {
	endpoint string
	samples  []time.Duration
	maxSize  int
}

func NewResponseTimeTracker(endpoint string, maxSamples int) *ResponseTimeTracker {
	return &ResponseTimeTracker{
		endpoint: endpoint,
		samples:  make([]time.Duration, 0, maxSamples),
		maxSize:  maxSamples,
	}
}

func (rt *ResponseTimeTracker) Record(duration time.Duration) {
	if len(rt.samples) >= rt.maxSize {
		// Remove oldest sample
		rt.samples = rt.samples[1:]
	}
	rt.samples = append(rt.samples, duration)
}

func (rt *ResponseTimeTracker) AverageResponseTime() time.Duration {
	if len(rt.samples) == 0 {
		return 0
	}

	var total time.Duration
	for _, sample := range rt.samples {
		total += sample
	}
	return total / time.Duration(len(rt.samples))
}

func (rt *ResponseTimeTracker) MaxResponseTime() time.Duration {
	if len(rt.samples) == 0 {
		return 0
	}

	max := rt.samples[0]
	for _, sample := range rt.samples[1:] {
		if sample > max {
			max = sample
		}
	}
	return max
}

func (rt *ResponseTimeTracker) MinResponseTime() time.Duration {
	if len(rt.samples) == 0 {
		return 0
	}

	min := rt.samples[0]
	for _, sample := range rt.samples[1:] {
		if sample < min {
			min = sample
		}
	}
	return min
}