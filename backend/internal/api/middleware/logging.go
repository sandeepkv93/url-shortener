package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"time"
)

type LoggingConfig struct {
	Logger          Logger
	SkipPaths       []string
	SkipSuccessLogs bool
}

type Logger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
}

// Default logger implementation using standard log package
type defaultLogger struct{}

func (l *defaultLogger) Info(msg string, fields ...interface{}) {
	log.Printf("[INFO] %s %v", msg, fields)
}

func (l *defaultLogger) Error(msg string, fields ...interface{}) {
	log.Printf("[ERROR] %s %v", msg, fields)
}

func (l *defaultLogger) Warn(msg string, fields ...interface{}) {
	log.Printf("[WARN] %s %v", msg, fields)
}

func (l *defaultLogger) Debug(msg string, fields ...interface{}) {
	log.Printf("[DEBUG] %s %v", msg, fields)
}

type LoggingMiddleware struct {
	config *LoggingConfig
}

func NewLoggingMiddleware(config *LoggingConfig) *LoggingMiddleware {
	if config == nil {
		config = &LoggingConfig{
			Logger: &defaultLogger{},
		}
	}
	if config.Logger == nil {
		config.Logger = &defaultLogger{}
	}
	return &LoggingMiddleware{
		config: config,
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}
	size, err := rw.ResponseWriter.Write(data)
	rw.size += size
	return size, err
}

func (m *LoggingMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip logging for certain paths
		if m.shouldSkipPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Generate request ID
		requestID := m.generateRequestID()
		
		// Add request ID to context and response header
		ctx := context.WithValue(r.Context(), "request_id", requestID)
		w.Header().Set("X-Request-ID", requestID)
		
		// Create response writer wrapper
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     0,
			size:          0,
		}

		start := time.Now()
		
		// Log request
		m.logRequest(r, requestID)
		
		// Process request
		next.ServeHTTP(rw, r.WithContext(ctx))
		
		duration := time.Since(start)
		
		// Log response
		m.logResponse(r, rw, requestID, duration)
	})
}

func (m *LoggingMiddleware) logRequest(r *http.Request, requestID string) {
	userAgent := r.Header.Get("User-Agent")
	if userAgent == "" {
		userAgent = "unknown"
	}

	m.config.Logger.Info("HTTP Request",
		"request_id", requestID,
		"method", r.Method,
		"path", r.URL.Path,
		"query", r.URL.RawQuery,
		"remote_addr", r.RemoteAddr,
		"user_agent", userAgent,
		"content_length", r.ContentLength,
	)
}

func (m *LoggingMiddleware) logResponse(r *http.Request, rw *responseWriter, requestID string, duration time.Duration) {
	statusCode := rw.statusCode
	if statusCode == 0 {
		statusCode = http.StatusOK
	}

	// Skip success logs if configured
	if m.config.SkipSuccessLogs && statusCode < 400 {
		return
	}

	logLevel := "Info"
	if statusCode >= 500 {
		logLevel = "Error"
	} else if statusCode >= 400 {
		logLevel = "Warn"
	}

	fields := []interface{}{
		"request_id", requestID,
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", statusCode,
		"response_size", rw.size,
		"duration_ms", duration.Milliseconds(),
		"duration", duration.String(),
	}

	switch logLevel {
	case "Error":
		m.config.Logger.Error("HTTP Response", fields...)
	case "Warn":
		m.config.Logger.Warn("HTTP Response", fields...)
	default:
		m.config.Logger.Info("HTTP Response", fields...)
	}
}

func (m *LoggingMiddleware) shouldSkipPath(path string) bool {
	for _, skipPath := range m.config.SkipPaths {
		if path == skipPath {
			return true
		}
	}
	return false
}

func (m *LoggingMiddleware) generateRequestID() string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based ID
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}

// Helper function to get request ID from context
func GetRequestIDFromContext(ctx context.Context) string {
	if requestID, ok := ctx.Value("request_id").(string); ok {
		return requestID
	}
	return ""
}

// Convenience function for easy setup
func Logging(logger Logger) func(http.Handler) http.Handler {
	config := &LoggingConfig{
		Logger: logger,
		SkipPaths: []string{
			"/health",
			"/metrics",
			"/favicon.ico",
		},
	}
	
	middleware := NewLoggingMiddleware(config)
	return middleware.Handler
}

// RequestLogging is a simpler version that just logs requests
func RequestLogging() func(http.Handler) http.Handler {
	return Logging(&defaultLogger{})
}