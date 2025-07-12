package services

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"sync"
	"time"
)

// ErrorTrackingService provides comprehensive error tracking and alerting
type ErrorTrackingService struct {
	errors         map[string]*ErrorSummary
	alerts         map[string]*Alert
	mu             sync.RWMutex
	logger         *LoggingService
	metricsService *MetricsService
	alertHandlers  []AlertHandler
	config         *ErrorTrackingConfig
}

// ErrorTrackingConfig configures error tracking behavior
type ErrorTrackingConfig struct {
	MaxErrors         int           `json:"max_errors"`
	AlertThreshold    int           `json:"alert_threshold"`
	TimeWindow        time.Duration `json:"time_window"`
	CriticalThreshold int           `json:"critical_threshold"`
	EnableAlerting    bool          `json:"enable_alerting"`
	RetentionPeriod   time.Duration `json:"retention_period"`
}

// ErrorSummary tracks error patterns and frequency
type ErrorSummary struct {
	Type           string                 `json:"type"`
	Message        string                 `json:"message"`
	Count          int64                  `json:"count"`
	FirstSeen      time.Time              `json:"first_seen"`
	LastSeen       time.Time              `json:"last_seen"`
	Component      string                 `json:"component"`
	Severity       ErrorSeverity          `json:"severity"`
	Stack          string                 `json:"stack"`
	Context        map[string]interface{} `json:"context"`
	Occurrences    []ErrorOccurrence      `json:"occurrences"`
	Fingerprint    string                 `json:"fingerprint"`
	UserID         string                 `json:"user_id,omitempty"`
	RequestID      string                 `json:"request_id,omitempty"`
	Tags           map[string]string      `json:"tags"`
	Resolved       bool                   `json:"resolved"`
	ResolvedAt     *time.Time             `json:"resolved_at,omitempty"`
	ResolvedBy     string                 `json:"resolved_by,omitempty"`
}

// ErrorOccurrence represents a single error occurrence
type ErrorOccurrence struct {
	Timestamp time.Time              `json:"timestamp"`
	Context   map[string]interface{} `json:"context"`
	UserID    string                 `json:"user_id,omitempty"`
	RequestID string                 `json:"request_id,omitempty"`
	Stack     string                 `json:"stack"`
}

// ErrorSeverity represents the severity level of an error
type ErrorSeverity string

const (
	SeverityCritical ErrorSeverity = "critical"
	SeverityHigh     ErrorSeverity = "high"
	SeverityMedium   ErrorSeverity = "medium"
	SeverityLow      ErrorSeverity = "low"
	SeverityInfo     ErrorSeverity = "info"
)

// Alert represents an alert triggered by error patterns
type Alert struct {
	ID          string                 `json:"id"`
	Type        AlertType              `json:"type"`
	Severity    ErrorSeverity          `json:"severity"`
	Title       string                 `json:"title"`
	Message     string                 `json:"message"`
	Timestamp   time.Time              `json:"timestamp"`
	Component   string                 `json:"component"`
	ErrorCount  int64                  `json:"error_count"`
	TimeWindow  time.Duration          `json:"time_window"`
	Context     map[string]interface{} `json:"context"`
	Resolved    bool                   `json:"resolved"`
	ResolvedAt  *time.Time             `json:"resolved_at,omitempty"`
	Fingerprint string                 `json:"fingerprint"`
}

// AlertType represents different types of alerts
type AlertType string

const (
	AlertTypeErrorRate      AlertType = "error_rate"
	AlertTypeNewError       AlertType = "new_error"
	AlertTypeCriticalError  AlertType = "critical_error"
	AlertTypeErrorSpike     AlertType = "error_spike"
	AlertTypeSystemHealth   AlertType = "system_health"
	AlertTypePerformance    AlertType = "performance"
)

// AlertHandler defines the interface for handling alerts
type AlertHandler interface {
	HandleAlert(ctx context.Context, alert *Alert) error
}

// EmailAlertHandler sends alerts via email
type EmailAlertHandler struct {
	notificationService *NotificationService
}

// SlackAlertHandler sends alerts to Slack
type SlackAlertHandler struct {
	webhookURL string
	logger     *LoggingService
}

// WebhookAlertHandler sends alerts to a generic webhook
type WebhookAlertHandler struct {
	webhookURL string
	logger     *LoggingService
}

// NewErrorTrackingService creates a new error tracking service
func NewErrorTrackingService(logger *LoggingService, metricsService *MetricsService, config *ErrorTrackingConfig) *ErrorTrackingService {
	if config == nil {
		config = &ErrorTrackingConfig{
			MaxErrors:         1000,
			AlertThreshold:    10,
			TimeWindow:        time.Hour,
			CriticalThreshold: 5,
			EnableAlerting:    true,
			RetentionPeriod:   24 * time.Hour,
		}
	}

	service := &ErrorTrackingService{
		errors:         make(map[string]*ErrorSummary),
		alerts:         make(map[string]*Alert),
		logger:         logger,
		metricsService: metricsService,
		config:         config,
	}

	// Start background cleanup goroutine
	go service.cleanup()

	return service
}

// TrackError records an error occurrence
func (ets *ErrorTrackingService) TrackError(ctx context.Context, err error, component string, severity ErrorSeverity, context map[string]interface{}) {
	if err == nil {
		return
	}

	// Extract context information
	userID := extractUserID(ctx)
	requestID := extractRequestID(ctx)

	// Generate fingerprint for error grouping
	fingerprint := ets.generateFingerprint(err, component)

	// Get stack trace
	stack := ets.getStackTrace()

	// Create or update error summary
	ets.mu.Lock()
	defer ets.mu.Unlock()

	summary, exists := ets.errors[fingerprint]
	if !exists {
		summary = &ErrorSummary{
			Type:        fmt.Sprintf("%T", err),
			Message:     err.Error(),
			Count:       0,
			FirstSeen:   time.Now(),
			Component:   component,
			Severity:    severity,
			Stack:       stack,
			Context:     context,
			Fingerprint: fingerprint,
			UserID:      userID,
			RequestID:   requestID,
			Tags:        make(map[string]string),
			Occurrences: make([]ErrorOccurrence, 0),
		}
		ets.errors[fingerprint] = summary

		// Log new error discovery
		ets.logger.Error("New error detected",
			"fingerprint", fingerprint,
			"type", summary.Type,
			"message", summary.Message,
			"component", component,
			"severity", severity,
		)

		// Check if this should trigger a new error alert
		if ets.config.EnableAlerting && severity == SeverityCritical {
			ets.triggerAlert(AlertTypeNewError, summary)
		}
	}

	// Update summary
	summary.Count++
	summary.LastSeen = time.Now()

	// Add occurrence
	occurrence := ErrorOccurrence{
		Timestamp: time.Now(),
		Context:   context,
		UserID:    userID,
		RequestID: requestID,
		Stack:     stack,
	}
	summary.Occurrences = append(summary.Occurrences, occurrence)

	// Limit occurrences to prevent memory bloat
	if len(summary.Occurrences) > 100 {
		summary.Occurrences = summary.Occurrences[1:]
	}

	// Record metrics
	if ets.metricsService != nil {
		counter := ets.metricsService.NewCounter(
			"errors_total",
			"Total number of errors",
			map[string]string{
				"component": component,
				"severity":  string(severity),
				"type":      summary.Type,
			},
		)
		counter.Inc()
	}

	// Check for alert conditions
	if ets.config.EnableAlerting {
		ets.checkAlertConditions(summary)
	}

	// Log the error
	ets.logError(summary, occurrence)
}

// GetErrorSummaries returns all error summaries
func (ets *ErrorTrackingService) GetErrorSummaries() []*ErrorSummary {
	ets.mu.RLock()
	defer ets.mu.RUnlock()

	summaries := make([]*ErrorSummary, 0, len(ets.errors))
	for _, summary := range ets.errors {
		summaries = append(summaries, summary)
	}

	return summaries
}

// GetErrorSummary returns a specific error summary by fingerprint
func (ets *ErrorTrackingService) GetErrorSummary(fingerprint string) (*ErrorSummary, bool) {
	ets.mu.RLock()
	defer ets.mu.RUnlock()

	summary, exists := ets.errors[fingerprint]
	return summary, exists
}

// ResolveError marks an error as resolved
func (ets *ErrorTrackingService) ResolveError(fingerprint, resolvedBy string) error {
	ets.mu.Lock()
	defer ets.mu.Unlock()

	summary, exists := ets.errors[fingerprint]
	if !exists {
		return fmt.Errorf("error with fingerprint %s not found", fingerprint)
	}

	now := time.Now()
	summary.Resolved = true
	summary.ResolvedAt = &now
	summary.ResolvedBy = resolvedBy

	ets.logger.Info("Error resolved",
		"fingerprint", fingerprint,
		"resolved_by", resolvedBy,
		"error_count", summary.Count,
	)

	return nil
}

// GetAlerts returns all alerts
func (ets *ErrorTrackingService) GetAlerts() []*Alert {
	ets.mu.RLock()
	defer ets.mu.RUnlock()

	alerts := make([]*Alert, 0, len(ets.alerts))
	for _, alert := range ets.alerts {
		alerts = append(alerts, alert)
	}

	return alerts
}

// AddAlertHandler adds an alert handler
func (ets *ErrorTrackingService) AddAlertHandler(handler AlertHandler) {
	ets.alertHandlers = append(ets.alertHandlers, handler)
}

// GetErrorStats returns error statistics
func (ets *ErrorTrackingService) GetErrorStats() map[string]interface{} {
	ets.mu.RLock()
	defer ets.mu.RUnlock()

	stats := map[string]interface{}{
		"total_errors":     len(ets.errors),
		"total_alerts":     len(ets.alerts),
		"critical_errors":  0,
		"high_errors":      0,
		"medium_errors":    0,
		"low_errors":       0,
		"resolved_errors":  0,
		"recent_errors":    0,
	}

	now := time.Now()
	recentThreshold := now.Add(-time.Hour)

	for _, summary := range ets.errors {
		switch summary.Severity {
		case SeverityCritical:
			stats["critical_errors"] = stats["critical_errors"].(int) + 1
		case SeverityHigh:
			stats["high_errors"] = stats["high_errors"].(int) + 1
		case SeverityMedium:
			stats["medium_errors"] = stats["medium_errors"].(int) + 1
		case SeverityLow:
			stats["low_errors"] = stats["low_errors"].(int) + 1
		}

		if summary.Resolved {
			stats["resolved_errors"] = stats["resolved_errors"].(int) + 1
		}

		if summary.LastSeen.After(recentThreshold) {
			stats["recent_errors"] = stats["recent_errors"].(int) + 1
		}
	}

	return stats
}

// Private methods

func (ets *ErrorTrackingService) generateFingerprint(err error, component string) string {
	// Create a fingerprint based on error type, message, and component
	// This groups similar errors together
	fingerprint := fmt.Sprintf("%T:%s:%s", err, err.Error(), component)
	return fmt.Sprintf("%x", fingerprint)
}

func (ets *ErrorTrackingService) getStackTrace() string {
	buf := make([]byte, 4096)
	stackSize := runtime.Stack(buf, false)
	return string(buf[:stackSize])
}

func (ets *ErrorTrackingService) checkAlertConditions(summary *ErrorSummary) {
	now := time.Now()
	windowStart := now.Add(-ets.config.TimeWindow)

	// Count recent occurrences
	recentCount := int64(0)
	for _, occurrence := range summary.Occurrences {
		if occurrence.Timestamp.After(windowStart) {
			recentCount++
		}
	}

	// Check for error rate threshold
	if recentCount >= int64(ets.config.AlertThreshold) {
		ets.triggerAlert(AlertTypeErrorRate, summary)
	}

	// Check for critical error threshold
	if summary.Severity == SeverityCritical && recentCount >= int64(ets.config.CriticalThreshold) {
		ets.triggerAlert(AlertTypeCriticalError, summary)
	}

	// Check for error spike (more than 10x normal rate)
	if recentCount > summary.Count/10 && recentCount > 5 {
		ets.triggerAlert(AlertTypeErrorSpike, summary)
	}
}

func (ets *ErrorTrackingService) triggerAlert(alertType AlertType, summary *ErrorSummary) {
	alertID := fmt.Sprintf("%s_%s_%d", alertType, summary.Fingerprint, time.Now().Unix())

	// Check if similar alert already exists
	for _, existingAlert := range ets.alerts {
		if existingAlert.Fingerprint == summary.Fingerprint && 
		   existingAlert.Type == alertType && 
		   !existingAlert.Resolved {
			return // Don't create duplicate alerts
		}
	}

	alert := &Alert{
		ID:          alertID,
		Type:        alertType,
		Severity:    summary.Severity,
		Title:       ets.generateAlertTitle(alertType, summary),
		Message:     ets.generateAlertMessage(alertType, summary),
		Timestamp:   time.Now(),
		Component:   summary.Component,
		ErrorCount:  summary.Count,
		TimeWindow:  ets.config.TimeWindow,
		Context:     summary.Context,
		Fingerprint: summary.Fingerprint,
	}

	ets.alerts[alertID] = alert

	// Send alert through handlers
	ctx := context.Background()
	for _, handler := range ets.alertHandlers {
		go func(h AlertHandler) {
			if err := h.HandleAlert(ctx, alert); err != nil {
				ets.logger.Error("Failed to send alert", "error", err.Error(), "alert_id", alertID)
			}
		}(handler)
	}

	ets.logger.Warn("Alert triggered",
		"alert_id", alertID,
		"alert_type", alertType,
		"component", summary.Component,
		"error_count", summary.Count,
	)
}

func (ets *ErrorTrackingService) generateAlertTitle(alertType AlertType, summary *ErrorSummary) string {
	switch alertType {
	case AlertTypeNewError:
		return fmt.Sprintf("New Error Detected: %s", summary.Type)
	case AlertTypeErrorRate:
		return fmt.Sprintf("High Error Rate: %s", summary.Component)
	case AlertTypeCriticalError:
		return fmt.Sprintf("Critical Error: %s", summary.Message)
	case AlertTypeErrorSpike:
		return fmt.Sprintf("Error Spike Detected: %s", summary.Component)
	default:
		return fmt.Sprintf("Error Alert: %s", summary.Type)
	}
}

func (ets *ErrorTrackingService) generateAlertMessage(alertType AlertType, summary *ErrorSummary) string {
	switch alertType {
	case AlertTypeNewError:
		return fmt.Sprintf("A new error has been detected in component %s: %s", summary.Component, summary.Message)
	case AlertTypeErrorRate:
		return fmt.Sprintf("Component %s has exceeded the error rate threshold with %d errors in the last %v", 
			summary.Component, summary.Count, ets.config.TimeWindow)
	case AlertTypeCriticalError:
		return fmt.Sprintf("Critical error detected in %s: %s (count: %d)", 
			summary.Component, summary.Message, summary.Count)
	case AlertTypeErrorSpike:
		return fmt.Sprintf("Error spike detected in %s with %d recent errors", summary.Component, summary.Count)
	default:
		return fmt.Sprintf("Error in %s: %s", summary.Component, summary.Message)
	}
}

func (ets *ErrorTrackingService) logError(summary *ErrorSummary, occurrence ErrorOccurrence) {
	ets.logger.LogSecurityEvent(
		"error_tracked",
		fmt.Sprintf("Error tracked: %s", summary.Message),
		string(summary.Severity),
		occurrence.UserID,
		occurrence.RequestID,
		"component", summary.Component,
		"error_type", summary.Type,
		"count", summary.Count,
		"fingerprint", summary.Fingerprint,
	)
}

func (ets *ErrorTrackingService) cleanup() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		ets.cleanupOldErrors()
		ets.cleanupOldAlerts()
	}
}

func (ets *ErrorTrackingService) cleanupOldErrors() {
	ets.mu.Lock()
	defer ets.mu.Unlock()

	cutoff := time.Now().Add(-ets.config.RetentionPeriod)
	
	for fingerprint, summary := range ets.errors {
		if summary.LastSeen.Before(cutoff) && summary.Resolved {
			delete(ets.errors, fingerprint)
		}
	}
}

func (ets *ErrorTrackingService) cleanupOldAlerts() {
	ets.mu.Lock()
	defer ets.mu.Unlock()

	cutoff := time.Now().Add(-ets.config.RetentionPeriod)
	
	for alertID, alert := range ets.alerts {
		if alert.Timestamp.Before(cutoff) && alert.Resolved {
			delete(ets.alerts, alertID)
		}
	}
}

// Helper functions

func extractUserID(ctx context.Context) string {
	if userID := ctx.Value("user_id"); userID != nil {
		if id, ok := userID.(string); ok {
			return id
		}
	}
	return ""
}

func extractRequestID(ctx context.Context) string {
	if requestID := ctx.Value("request_id"); requestID != nil {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}

// Alert Handler Implementations

func NewEmailAlertHandler(notificationService *NotificationService) *EmailAlertHandler {
	return &EmailAlertHandler{
		notificationService: notificationService,
	}
}

func (eah *EmailAlertHandler) HandleAlert(ctx context.Context, alert *Alert) error {
	// Format alert as email
	subject := fmt.Sprintf("[ALERT] %s", alert.Title)
	
	body := fmt.Sprintf(`
Alert Details:
- Type: %s
- Severity: %s
- Component: %s
- Message: %s
- Timestamp: %s
- Error Count: %d

Context:
%s
	`, alert.Type, alert.Severity, alert.Component, alert.Message, 
		alert.Timestamp.Format(time.RFC3339), alert.ErrorCount, 
		formatContext(alert.Context))

	// Send via notification service (would need email addresses configured)
	return fmt.Errorf("email notification not implemented yet")
}

func NewSlackAlertHandler(webhookURL string, logger *LoggingService) *SlackAlertHandler {
	return &SlackAlertHandler{
		webhookURL: webhookURL,
		logger:     logger,
	}
}

func (sah *SlackAlertHandler) HandleAlert(ctx context.Context, alert *Alert) error {
	// Format alert as Slack message
	payload := map[string]interface{}{
		"text": fmt.Sprintf("🚨 %s", alert.Title),
		"attachments": []map[string]interface{}{
			{
				"color": sah.getSeverityColor(alert.Severity),
				"fields": []map[string]interface{}{
					{"title": "Component", "value": alert.Component, "short": true},
					{"title": "Severity", "value": string(alert.Severity), "short": true},
					{"title": "Error Count", "value": fmt.Sprintf("%d", alert.ErrorCount), "short": true},
					{"title": "Time", "value": alert.Timestamp.Format(time.RFC3339), "short": true},
				},
				"text": alert.Message,
			},
		},
	}

	// Send to Slack webhook (implementation would use HTTP client)
	sah.logger.Info("Slack alert sent", "alert_id", alert.ID, "webhook", sah.webhookURL)
	return nil
}

func (sah *SlackAlertHandler) getSeverityColor(severity ErrorSeverity) string {
	switch severity {
	case SeverityCritical:
		return "danger"
	case SeverityHigh:
		return "warning"
	case SeverityMedium:
		return "good"
	default:
		return "#439FE0"
	}
}

func formatContext(context map[string]interface{}) string {
	if len(context) == 0 {
		return "No additional context"
	}
	
	data, _ := json.MarshalIndent(context, "", "  ")
	return string(data)
}