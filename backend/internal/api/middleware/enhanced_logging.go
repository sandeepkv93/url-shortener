package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"url-shortener/internal/core/services"
)

// EnhancedLoggingMiddleware provides comprehensive request/response logging with structured logging
type EnhancedLoggingMiddleware struct {
	logger         *services.LoggingService
	metricsService *services.MetricsService
	errorTracking  *services.ErrorTrackingService
	config         *EnhancedLoggingConfig
}

// EnhancedLoggingConfig configures the enhanced logging middleware
type EnhancedLoggingConfig struct {
	SkipPaths           []string `json:"skip_paths"`
	SkipSuccessLogs     bool     `json:"skip_success_logs"`
	LogRequestBody      bool     `json:"log_request_body"`
	LogResponseBody     bool     `json:"log_response_body"`
	MaxBodyLogSize      int      `json:"max_body_log_size"`
	TrackPerformance    bool     `json:"track_performance"`
	TrackErrors         bool     `json:"track_errors"`
	SlowRequestThreshold time.Duration `json:"slow_request_threshold"`
}

// NewEnhancedLoggingMiddleware creates a new enhanced logging middleware
func NewEnhancedLoggingMiddleware(
	logger *services.LoggingService,
	metricsService *services.MetricsService,
	errorTracking *services.ErrorTrackingService,
	config *EnhancedLoggingConfig,
) *EnhancedLoggingMiddleware {
	if config == nil {
		config = &EnhancedLoggingConfig{
			SkipPaths: []string{
				"/health",
				"/metrics",
				"/favicon.ico",
			},
			SkipSuccessLogs:      false,
			LogRequestBody:       false,
			LogResponseBody:      false,
			MaxBodyLogSize:       1024,
			TrackPerformance:     true,
			TrackErrors:          true,
			SlowRequestThreshold: 2 * time.Second,
		}
	}

	return &EnhancedLoggingMiddleware{
		logger:         logger,
		metricsService: metricsService,
		errorTracking:  errorTracking,
		config:         config,
	}
}

// Handler provides the middleware handler
func (elm *EnhancedLoggingMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip logging for certain paths
		if elm.shouldSkipPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Generate request ID if not present
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = elm.generateRequestID()
		}

		// Add request ID to context and response header
		ctx := context.WithValue(r.Context(), "request_id", requestID)
		w.Header().Set("X-Request-ID", requestID)

		// Create enhanced response writer
		erw := &enhancedResponseWriter{
			ResponseWriter: w,
			statusCode:     0,
			size:          0,
			body:          make([]byte, 0),
		}

		start := time.Now()

		// Log request start
		elm.logRequestStart(r, requestID)

		// Start performance tracking
		var performanceTracker *services.PerformanceTracker
		if elm.config.TrackPerformance && elm.logger != nil {
			performanceTracker = elm.logger.NewPerformanceTracker(
				fmt.Sprintf("http_%s_%s", r.Method, r.URL.Path),
			)
			performanceTracker.AddMetadata("method", r.Method)
			performanceTracker.AddMetadata("path", r.URL.Path)
			performanceTracker.AddMetadata("user_agent", r.Header.Get("User-Agent"))
			performanceTracker.AddMetadata("remote_addr", r.RemoteAddr)
		}

		// Process request with panic recovery
		var requestError error
		func() {
			defer func() {
				if rec := recover(); rec != nil {
					requestError = fmt.Errorf("panic in request handler: %v", rec)
					erw.statusCode = http.StatusInternalServerError
					
					// Track the panic as an error
					if elm.config.TrackErrors && elm.errorTracking != nil {
						elm.errorTracking.TrackError(
							ctx,
							requestError,
							"http_handler",
							services.SeverityCritical,
							map[string]interface{}{
								"method":      r.Method,
								"path":        r.URL.Path,
								"user_agent":  r.Header.Get("User-Agent"),
								"remote_addr": r.RemoteAddr,
							},
						)
					}
				}
			}()

			next.ServeHTTP(erw, r.WithContext(ctx))
		}()

		duration := time.Since(start)

		// Determine final status code
		statusCode := erw.statusCode
		if statusCode == 0 {
			statusCode = http.StatusOK
		}

		// Log response
		elm.logRequestComplete(r, erw, requestID, duration, requestError)

		// Record metrics
		if elm.metricsService != nil {
			elm.metricsService.RecordHTTPRequest(
				r.Method,
				r.URL.Path,
				statusCode,
				duration,
				erw.size,
			)
		}

		// Complete performance tracking
		if performanceTracker != nil {
			success := statusCode < 400 && requestError == nil
			performanceTracker.AddMetadata("status_code", statusCode)
			performanceTracker.AddMetadata("response_size", erw.size)
			performanceTracker.Finish(success, requestError)
		}

		// Track errors for 4xx and 5xx responses
		if elm.config.TrackErrors && elm.errorTracking != nil && statusCode >= 400 {
			severity := services.SeverityMedium
			if statusCode >= 500 {
				severity = services.SeverityHigh
			}

			err := fmt.Errorf("HTTP %d: %s %s", statusCode, r.Method, r.URL.Path)
			elm.errorTracking.TrackError(
				ctx,
				err,
				"http_response",
				severity,
				map[string]interface{}{
					"method":       r.Method,
					"path":         r.URL.Path,
					"status_code":  statusCode,
					"duration_ms":  duration.Milliseconds(),
					"response_size": erw.size,
					"user_agent":   r.Header.Get("User-Agent"),
					"remote_addr":  r.RemoteAddr,
				},
			)
		}

		// Log slow requests
		if duration > elm.config.SlowRequestThreshold {
			elm.logSlowRequest(r, duration, statusCode, requestID)
		}
	})
}

// enhancedResponseWriter wraps http.ResponseWriter to capture response data
type enhancedResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
	body       []byte
}

func (erw *enhancedResponseWriter) WriteHeader(statusCode int) {
	erw.statusCode = statusCode
	erw.ResponseWriter.WriteHeader(statusCode)
}

func (erw *enhancedResponseWriter) Write(data []byte) (int, error) {
	if erw.statusCode == 0 {
		erw.statusCode = http.StatusOK
	}
	size, err := erw.ResponseWriter.Write(data)
	erw.size += size

	// Optionally capture response body for logging
	if len(erw.body) < 1024 { // Limit body capture size
		erw.body = append(erw.body, data...)
	}

	return size, err
}

func (elm *EnhancedLoggingMiddleware) logRequestStart(r *http.Request, requestID string) {
	userAgent := r.Header.Get("User-Agent")
	if userAgent == "" {
		userAgent = "unknown"
	}

	elm.logger.LogHTTPRequest(
		r.Method,
		r.URL.Path,
		userAgent,
		r.RemoteAddr,
		requestID,
		"query", r.URL.RawQuery,
		"content_length", r.ContentLength,
		"protocol", r.Proto,
		"host", r.Host,
		"referer", r.Header.Get("Referer"),
	)
}

func (elm *EnhancedLoggingMiddleware) logRequestComplete(
	r *http.Request,
	erw *enhancedResponseWriter,
	requestID string,
	duration time.Duration,
	requestError error,
) {
	statusCode := erw.statusCode
	if statusCode == 0 {
		statusCode = http.StatusOK
	}

	// Skip success logs if configured
	if elm.config.SkipSuccessLogs && statusCode < 400 && requestError == nil {
		return
	}

	fields := []interface{}{
		"query", r.URL.RawQuery,
		"protocol", r.Proto,
		"host", r.Host,
		"referer", r.Header.Get("Referer"),
		"user_agent", r.Header.Get("User-Agent"),
	}

	// Add error information if present
	if requestError != nil {
		fields = append(fields, "error", requestError.Error())
	}

	elm.logger.LogHTTPResponse(
		r.Method,
		r.URL.Path,
		statusCode,
		duration,
		erw.size,
		requestID,
		fields...,
	)
}

func (elm *EnhancedLoggingMiddleware) logSlowRequest(r *http.Request, duration time.Duration, statusCode int, requestID string) {
	elm.logger.Warn("Slow request detected",
		"method", r.Method,
		"path", r.URL.Path,
		"duration_ms", duration.Milliseconds(),
		"status_code", statusCode,
		"request_id", requestID,
		"user_agent", r.Header.Get("User-Agent"),
		"remote_addr", r.RemoteAddr,
		"type", "slow_request",
	)
}

func (elm *EnhancedLoggingMiddleware) shouldSkipPath(path string) bool {
	for _, skipPath := range elm.config.SkipPaths {
		if path == skipPath {
			return true
		}
	}
	return false
}

func (elm *EnhancedLoggingMiddleware) generateRequestID() string {
	// Generate a unique request ID
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// MonitoringMiddleware provides comprehensive monitoring capabilities
type MonitoringMiddleware struct {
	logger         *services.LoggingService
	metricsService *services.MetricsService
	errorTracking  *services.ErrorTrackingService
}

// NewMonitoringMiddleware creates a new monitoring middleware
func NewMonitoringMiddleware(
	logger *services.LoggingService,
	metricsService *services.MetricsService,
	errorTracking *services.ErrorTrackingService,
) *MonitoringMiddleware {
	return &MonitoringMiddleware{
		logger:         logger,
		metricsService: metricsService,
		errorTracking:  errorTracking,
	}
}

// Handler provides comprehensive monitoring for HTTP requests
func (mm *MonitoringMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add monitoring headers
		w.Header().Set("X-Monitoring-Enabled", "true")
		w.Header().Set("X-Service-Name", "url-shortener")
		w.Header().Set("X-Service-Version", "1.0.0")

		// Add request start time to context
		ctx := context.WithValue(r.Context(), "request_start_time", time.Now())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ErrorTrackingMiddleware provides automatic error tracking
type ErrorTrackingMiddleware struct {
	errorTracking *services.ErrorTrackingService
	logger        *services.LoggingService
}

// NewErrorTrackingMiddleware creates a new error tracking middleware
func NewErrorTrackingMiddleware(
	errorTracking *services.ErrorTrackingService,
	logger *services.LoggingService,
) *ErrorTrackingMiddleware {
	return &ErrorTrackingMiddleware{
		errorTracking: errorTracking,
		logger:        logger,
	}
}

// Handler provides automatic error tracking for panics and errors
func (etm *ErrorTrackingMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				// Track panic as critical error
				err := fmt.Errorf("panic in HTTP handler: %v", rec)
				etm.errorTracking.TrackError(
					r.Context(),
					err,
					"http_middleware",
					services.SeverityCritical,
					map[string]interface{}{
						"method":      r.Method,
						"path":        r.URL.Path,
						"user_agent":  r.Header.Get("User-Agent"),
						"remote_addr": r.RemoteAddr,
						"panic_value": rec,
					},
				)

				// Log the panic
				etm.logger.Error("Panic in HTTP handler",
					"error", err.Error(),
					"method", r.Method,
					"path", r.URL.Path,
					"panic_value", rec,
				)

				// Return internal server error
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// MetricsExportMiddleware adds metrics export headers
type MetricsExportMiddleware struct {
	metricsService *services.MetricsService
}

// NewMetricsExportMiddleware creates a new metrics export middleware
func NewMetricsExportMiddleware(metricsService *services.MetricsService) *MetricsExportMiddleware {
	return &MetricsExportMiddleware{
		metricsService: metricsService,
	}
}

// Handler adds metrics information to response headers
func (mem *MetricsExportMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add metrics endpoint information
		w.Header().Set("X-Metrics-Endpoint", "/metrics")
		w.Header().Set("X-Metrics-Format", "prometheus")

		next.ServeHTTP(w, r)
	})
}