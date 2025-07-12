package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

// LogLevel represents different log levels
type LogLevel string

const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
	FatalLevel LogLevel = "fatal"
	PanicLevel LogLevel = "panic"
)

// LoggingService provides structured logging capabilities
type LoggingService struct {
	logger     *logrus.Logger
	serviceName string
	version    string
	environment string
}

// LogEntry represents a structured log entry
type LogEntry struct {
	Level       LogLevel               `json:"level"`
	Message     string                 `json:"message"`
	Timestamp   time.Time              `json:"timestamp"`
	Service     string                 `json:"service"`
	Version     string                 `json:"version"`
	Environment string                 `json:"environment"`
	RequestID   string                 `json:"request_id,omitempty"`
	UserID      string                 `json:"user_id,omitempty"`
	Component   string                 `json:"component,omitempty"`
	Operation   string                 `json:"operation,omitempty"`
	Duration    *time.Duration         `json:"duration,omitempty"`
	Error       *ErrorDetails          `json:"error,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Stack       string                 `json:"stack,omitempty"`
}

// ErrorDetails provides structured error information
type ErrorDetails struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
	Stack   string `json:"stack,omitempty"`
}

// NewLoggingService creates a new logging service
func NewLoggingService(serviceName, version, environment string) *LoggingService {
	logger := logrus.New()
	
	// Configure logger based on environment
	if environment == "production" {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		})
		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
		logger.SetLevel(logrus.DebugLevel)
	}
	
	logger.SetOutput(os.Stdout)
	
	return &LoggingService{
		logger:      logger,
		serviceName: serviceName,
		version:     version,
		environment: environment,
	}
}

// WithContext creates a logger with context information
func (ls *LoggingService) WithContext(ctx context.Context) *ContextLogger {
	return &ContextLogger{
		service: ls,
		ctx:     ctx,
	}
}

// WithComponent creates a logger for a specific component
func (ls *LoggingService) WithComponent(component string) *ComponentLogger {
	return &ComponentLogger{
		service:   ls,
		component: component,
	}
}

// Debug logs a debug message
func (ls *LoggingService) Debug(message string, fields ...interface{}) {
	ls.log(DebugLevel, message, fields...)
}

// Info logs an info message
func (ls *LoggingService) Info(message string, fields ...interface{}) {
	ls.log(InfoLevel, message, fields...)
}

// Warn logs a warning message
func (ls *LoggingService) Warn(message string, fields ...interface{}) {
	ls.log(WarnLevel, message, fields...)
}

// Error logs an error message
func (ls *LoggingService) Error(message string, fields ...interface{}) {
	ls.log(ErrorLevel, message, fields...)
}

// Fatal logs a fatal message and exits
func (ls *LoggingService) Fatal(message string, fields ...interface{}) {
	ls.log(FatalLevel, message, fields...)
}

// Panic logs a panic message and panics
func (ls *LoggingService) Panic(message string, fields ...interface{}) {
	ls.log(PanicLevel, message, fields...)
}

// LogOperation logs an operation with duration
func (ls *LoggingService) LogOperation(operation string, duration time.Duration, success bool, fields ...interface{}) {
	level := InfoLevel
	if !success {
		level = ErrorLevel
	}
	
	allFields := append(fields, "operation", operation, "duration_ms", duration.Milliseconds(), "success", success)
	ls.log(level, fmt.Sprintf("Operation %s completed", operation), allFields...)
}

// LogHTTPRequest logs an HTTP request
func (ls *LoggingService) LogHTTPRequest(method, path, userAgent, remoteAddr string, requestID string, fields ...interface{}) {
	allFields := append(fields, 
		"method", method,
		"path", path,
		"user_agent", userAgent,
		"remote_addr", remoteAddr,
		"request_id", requestID,
		"type", "http_request",
	)
	ls.log(InfoLevel, "HTTP request received", allFields...)
}

// LogHTTPResponse logs an HTTP response
func (ls *LoggingService) LogHTTPResponse(method, path string, statusCode int, duration time.Duration, responseSize int, requestID string, fields ...interface{}) {
	level := InfoLevel
	if statusCode >= 500 {
		level = ErrorLevel
	} else if statusCode >= 400 {
		level = WarnLevel
	}
	
	allFields := append(fields,
		"method", method,
		"path", path,
		"status_code", statusCode,
		"duration_ms", duration.Milliseconds(),
		"response_size", responseSize,
		"request_id", requestID,
		"type", "http_response",
	)
	
	message := fmt.Sprintf("HTTP %s %s %d", method, path, statusCode)
	ls.log(level, message, allFields...)
}

// LogDatabaseOperation logs a database operation
func (ls *LoggingService) LogDatabaseOperation(operation, table string, duration time.Duration, rowsAffected int64, err error, fields ...interface{}) {
	level := InfoLevel
	if err != nil {
		level = ErrorLevel
	}
	
	allFields := append(fields,
		"operation", operation,
		"table", table,
		"duration_ms", duration.Milliseconds(),
		"rows_affected", rowsAffected,
		"type", "database_operation",
	)
	
	if err != nil {
		allFields = append(allFields, "error", err.Error())
	}
	
	message := fmt.Sprintf("Database %s on %s", operation, table)
	ls.log(level, message, allFields...)
}

// LogCacheOperation logs a cache operation
func (ls *LoggingService) LogCacheOperation(operation, key string, duration time.Duration, hit bool, err error, fields ...interface{}) {
	level := InfoLevel
	if err != nil {
		level = ErrorLevel
	}
	
	allFields := append(fields,
		"operation", operation,
		"key", key,
		"duration_ms", duration.Milliseconds(),
		"cache_hit", hit,
		"type", "cache_operation",
	)
	
	if err != nil {
		allFields = append(allFields, "error", err.Error())
	}
	
	message := fmt.Sprintf("Cache %s for key %s", operation, key)
	ls.log(level, message, allFields...)
}

// LogSecurityEvent logs a security-related event
func (ls *LoggingService) LogSecurityEvent(eventType, description string, severity string, userID, requestID string, fields ...interface{}) {
	level := WarnLevel
	switch severity {
	case "critical":
		level = ErrorLevel
	case "high":
		level = ErrorLevel
	case "medium":
		level = WarnLevel
	case "low":
		level = InfoLevel
	}
	
	allFields := append(fields,
		"event_type", eventType,
		"severity", severity,
		"user_id", userID,
		"request_id", requestID,
		"type", "security_event",
	)
	
	ls.log(level, fmt.Sprintf("Security event: %s", description), allFields...)
}

// LogBusinessEvent logs a business logic event
func (ls *LoggingService) LogBusinessEvent(eventType, description string, userID string, fields ...interface{}) {
	allFields := append(fields,
		"event_type", eventType,
		"user_id", userID,
		"type", "business_event",
	)
	
	ls.log(InfoLevel, fmt.Sprintf("Business event: %s", description), allFields...)
}

// Internal logging method
func (ls *LoggingService) log(level LogLevel, message string, fields ...interface{}) {
	entry := ls.logger.WithFields(ls.buildFields(fields...))
	
	switch level {
	case DebugLevel:
		entry.Debug(message)
	case InfoLevel:
		entry.Info(message)
	case WarnLevel:
		entry.Warn(message)
	case ErrorLevel:
		entry.Error(message)
	case FatalLevel:
		entry.Fatal(message)
	case PanicLevel:
		entry.Panic(message)
	}
}

// buildFields converts field pairs to logrus.Fields
func (ls *LoggingService) buildFields(fields ...interface{}) logrus.Fields {
	logFields := logrus.Fields{
		"service":     ls.serviceName,
		"version":     ls.version,
		"environment": ls.environment,
		"timestamp":   time.Now(),
	}
	
	// Add provided fields
	for i := 0; i < len(fields)-1; i += 2 {
		if key, ok := fields[i].(string); ok {
			logFields[key] = fields[i+1]
		}
	}
	
	return logFields
}

// ContextLogger provides logging with context
type ContextLogger struct {
	service *LoggingService
	ctx     context.Context
}

// Debug logs a debug message with context
func (cl *ContextLogger) Debug(message string, fields ...interface{}) {
	cl.logWithContext(DebugLevel, message, fields...)
}

// Info logs an info message with context
func (cl *ContextLogger) Info(message string, fields ...interface{}) {
	cl.logWithContext(InfoLevel, message, fields...)
}

// Warn logs a warning message with context
func (cl *ContextLogger) Warn(message string, fields ...interface{}) {
	cl.logWithContext(WarnLevel, message, fields...)
}

// Error logs an error message with context
func (cl *ContextLogger) Error(message string, fields ...interface{}) {
	cl.logWithContext(ErrorLevel, message, fields...)
}

func (cl *ContextLogger) logWithContext(level LogLevel, message string, fields ...interface{}) {
	// Extract context information
	contextFields := cl.extractContextFields()
	allFields := append(fields, contextFields...)
	cl.service.log(level, message, allFields...)
}

func (cl *ContextLogger) extractContextFields() []interface{} {
	var fields []interface{}
	
	// Extract request ID from context
	if requestID := cl.ctx.Value("request_id"); requestID != nil {
		fields = append(fields, "request_id", requestID)
	}
	
	// Extract user ID from context
	if userID := cl.ctx.Value("user_id"); userID != nil {
		fields = append(fields, "user_id", userID)
	}
	
	// Extract trace ID from context
	if traceID := cl.ctx.Value("trace_id"); traceID != nil {
		fields = append(fields, "trace_id", traceID)
	}
	
	return fields
}

// ComponentLogger provides logging for a specific component
type ComponentLogger struct {
	service   *LoggingService
	component string
}

// Debug logs a debug message for the component
func (cl *ComponentLogger) Debug(message string, fields ...interface{}) {
	cl.logWithComponent(DebugLevel, message, fields...)
}

// Info logs an info message for the component
func (cl *ComponentLogger) Info(message string, fields ...interface{}) {
	cl.logWithComponent(InfoLevel, message, fields...)
}

// Warn logs a warning message for the component
func (cl *ComponentLogger) Warn(message string, fields ...interface{}) {
	cl.logWithComponent(WarnLevel, message, fields...)
}

// Error logs an error message for the component
func (cl *ComponentLogger) Error(message string, fields ...interface{}) {
	cl.logWithComponent(ErrorLevel, message, fields...)
}

func (cl *ComponentLogger) logWithComponent(level LogLevel, message string, fields ...interface{}) {
	allFields := append(fields, "component", cl.component)
	cl.service.log(level, message, allFields...)
}

// PerformanceTracker tracks operation performance
type PerformanceTracker struct {
	logger    *LoggingService
	operation string
	startTime time.Time
	metadata  map[string]interface{}
}

// NewPerformanceTracker creates a new performance tracker
func (ls *LoggingService) NewPerformanceTracker(operation string) *PerformanceTracker {
	return &PerformanceTracker{
		logger:    ls,
		operation: operation,
		startTime: time.Now(),
		metadata:  make(map[string]interface{}),
	}
}

// AddMetadata adds metadata to the performance tracker
func (pt *PerformanceTracker) AddMetadata(key string, value interface{}) {
	pt.metadata[key] = value
}

// Finish completes the performance tracking and logs the result
func (pt *PerformanceTracker) Finish(success bool, err error) {
	duration := time.Since(pt.startTime)
	
	fields := []interface{}{
		"operation", pt.operation,
		"duration_ms", duration.Milliseconds(),
		"success", success,
	}
	
	// Add metadata
	for key, value := range pt.metadata {
		fields = append(fields, key, value)
	}
	
	// Add error information if present
	if err != nil {
		fields = append(fields, "error", err.Error())
	}
	
	if success {
		pt.logger.Info(fmt.Sprintf("Operation %s completed successfully", pt.operation), fields...)
	} else {
		pt.logger.Error(fmt.Sprintf("Operation %s failed", pt.operation), fields...)
	}
}

// LoggerInterface provides a standard interface for logging
type LoggerInterface interface {
	Debug(message string, fields ...interface{})
	Info(message string, fields ...interface{})
	Warn(message string, fields ...interface{})
	Error(message string, fields ...interface{})
}

// Ensure LoggingService implements LoggerInterface
var _ LoggerInterface = (*LoggingService)(nil)

// Helper functions for structured logging

// LogWithError logs with structured error information
func LogWithError(logger LoggerInterface, level LogLevel, message string, err error, fields ...interface{}) {
	if err != nil {
		// Get stack trace
		buf := make([]byte, 1024)
		stackSize := runtime.Stack(buf, false)
		
		errorFields := append(fields,
			"error", err.Error(),
			"error_type", fmt.Sprintf("%T", err),
			"stack", string(buf[:stackSize]),
		)
		
		switch level {
		case ErrorLevel:
			logger.Error(message, errorFields...)
		case WarnLevel:
			logger.Warn(message, errorFields...)
		default:
			logger.Info(message, errorFields...)
		}
	} else {
		switch level {
		case ErrorLevel:
			logger.Error(message, fields...)
		case WarnLevel:
			logger.Warn(message, fields...)
		default:
			logger.Info(message, fields...)
		}
	}
}

// LogJSON logs a JSON-serializable object
func LogJSON(logger LoggerInterface, level LogLevel, message string, obj interface{}, fields ...interface{}) {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		logger.Error("Failed to marshal object to JSON", "error", err.Error(), "object_type", fmt.Sprintf("%T", obj))
		return
	}
	
	allFields := append(fields, "json_data", string(jsonData))
	
	switch level {
	case ErrorLevel:
		logger.Error(message, allFields...)
	case WarnLevel:
		logger.Warn(message, allFields...)
	default:
		logger.Info(message, allFields...)
	}
}