package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"url-shortener/internal/core/domain"
	"url-shortener/internal/core/ports"
)

// SecurityAuditService handles security event logging and monitoring
type SecurityAuditService interface {
	LogSecurityEvent(ctx context.Context, event *domain.SecurityEvent) error
	LogAuthEvent(ctx context.Context, event *domain.AuthEvent) error
	LogRateLimitEvent(ctx context.Context, event *domain.RateLimitEvent) error
	LogSuspiciousActivity(ctx context.Context, event *domain.SuspiciousActivityEvent) error
	GetSecurityMetrics(ctx context.Context, timeRange time.Duration) (*domain.SecurityMetrics, error)
	GenerateSecurityReport(ctx context.Context, timeRange time.Duration) (*domain.SecurityReport, error)
	CheckSecurityAlerts(ctx context.Context) ([]*domain.SecurityAlert, error)
}

type securityAuditService struct {
	cache  ports.CacheService
	logger *logrus.Logger
}

// NewSecurityAuditService creates a new security audit service
func NewSecurityAuditService(cache ports.CacheService, logger *logrus.Logger) SecurityAuditService {
	return &securityAuditService{
		cache:  cache,
		logger: logger,
	}
}

// LogSecurityEvent logs a general security event
func (s *securityAuditService) LogSecurityEvent(ctx context.Context, event *domain.SecurityEvent) error {
	// Log to structured logger
	logEntry := s.logger.WithFields(logrus.Fields{
		"event_type":    event.Type,
		"severity":      event.Severity,
		"source_ip":     event.SourceIP,
		"user_id":       event.UserID,
		"endpoint":      event.Endpoint,
		"user_agent":    event.UserAgent,
		"timestamp":     event.Timestamp,
		"description":   event.Description,
		"metadata":      event.Metadata,
	})
	
	switch event.Severity {
	case domain.SecuritySeverityLow:
		logEntry.Info("Security event")
	case domain.SecuritySeverityMedium:
		logEntry.Warn("Security event")
	case domain.SecuritySeverityHigh:
		logEntry.Error("Security event")
	case domain.SecuritySeverityCritical:
		logEntry.Error("CRITICAL Security event")
	}
	
	// Store in cache for metrics and alerting
	key := fmt.Sprintf("security_event:%s:%d", event.Type, event.Timestamp.Unix())
	eventData, _ := json.Marshal(event)
	
	if err := s.cache.Set(ctx, key, string(eventData), 24*time.Hour); err != nil {
		s.logger.WithError(err).Error("Failed to store security event in cache")
	}
	
	// Update metrics counters
	s.updateSecurityMetrics(ctx, event)
	
	// Check for alert conditions
	s.checkAndTriggerAlerts(ctx, event)
	
	return nil
}

// LogAuthEvent logs authentication-related events
func (s *securityAuditService) LogAuthEvent(ctx context.Context, event *domain.AuthEvent) error {
	logEntry := s.logger.WithFields(logrus.Fields{
		"event_type":    "auth",
		"action":        event.Action,
		"user_id":       event.UserID,
		"email":         event.Email,
		"source_ip":     event.SourceIP,
		"user_agent":    event.UserAgent,
		"success":       event.Success,
		"failure_reason": event.FailureReason,
		"timestamp":     event.Timestamp,
	})
	
	if event.Success {
		logEntry.Info("Authentication event")
	} else {
		logEntry.Warn("Authentication failure")
	}
	
	// Store for analysis
	key := fmt.Sprintf("auth_event:%s:%d", event.Action, event.Timestamp.Unix())
	eventData, _ := json.Marshal(event)
	
	if err := s.cache.Set(ctx, key, string(eventData), 24*time.Hour); err != nil {
		s.logger.WithError(err).Error("Failed to store auth event in cache")
	}
	
	// Track failed login attempts
	if !event.Success && event.Action == domain.AuthActionLogin {
		s.trackFailedLoginAttempts(ctx, event.SourceIP, event.Email)
	}
	
	return nil
}

// LogRateLimitEvent logs rate limiting events
func (s *securityAuditService) LogRateLimitEvent(ctx context.Context, event *domain.RateLimitEvent) error {
	s.logger.WithFields(logrus.Fields{
		"event_type":    "rate_limit",
		"limit_type":    event.LimitType,
		"source_ip":     event.SourceIP,
		"user_id":       event.UserID,
		"endpoint":      event.Endpoint,
		"method":        event.Method,
		"current_count": event.CurrentCount,
		"limit":         event.Limit,
		"window":        event.Window,
		"timestamp":     event.Timestamp,
	}).Warn("Rate limit exceeded")
	
	// Store for analysis
	key := fmt.Sprintf("rate_limit_event:%s:%d", event.LimitType, event.Timestamp.Unix())
	eventData, _ := json.Marshal(event)
	
	if err := s.cache.Set(ctx, key, string(eventData), time.Hour); err != nil {
		s.logger.WithError(err).Error("Failed to store rate limit event in cache")
	}
	
	// Track rate limit violations per IP
	s.trackRateLimitViolations(ctx, event.SourceIP)
	
	return nil
}

// LogSuspiciousActivity logs suspicious activity events
func (s *securityAuditService) LogSuspiciousActivity(ctx context.Context, event *domain.SuspiciousActivityEvent) error {
	s.logger.WithFields(logrus.Fields{
		"event_type":        "suspicious_activity",
		"activity_type":     event.ActivityType,
		"source_ip":         event.SourceIP,
		"user_id":           event.UserID,
		"endpoint":          event.Endpoint,
		"indicators":        event.Indicators,
		"risk_score":        event.RiskScore,
		"automatic_action":  event.AutomaticAction,
		"timestamp":         event.Timestamp,
		"description":       event.Description,
	}).Error("Suspicious activity detected")
	
	// Store for analysis
	key := fmt.Sprintf("suspicious_activity:%s:%d", event.ActivityType, event.Timestamp.Unix())
	eventData, _ := json.Marshal(event)
	
	if err := s.cache.Set(ctx, key, string(eventData), 7*24*time.Hour); err != nil {
		s.logger.WithError(err).Error("Failed to store suspicious activity event in cache")
	}
	
	// Update risk score for IP
	s.updateIPRiskScore(ctx, event.SourceIP, event.RiskScore)
	
	return nil
}

// GetSecurityMetrics retrieves security metrics for a given time range
func (s *securityAuditService) GetSecurityMetrics(ctx context.Context, timeRange time.Duration) (*domain.SecurityMetrics, error) {
	metrics := &domain.SecurityMetrics{
		TimeRange:     timeRange,
		CollectedAt:   time.Now(),
		EventCounts:   make(map[string]int),
		ThreatLevels:  make(map[string]int),
		TopAttackers:  []domain.AttackerInfo{},
		TopEndpoints:  []domain.EndpointInfo{},
	}
	
	// Get metrics from cache counters
	now := time.Now()
	startTime := now.Add(-timeRange)
	
	// Count different types of security events
	eventTypes := []string{
		"failed_login", "rate_limit_violation", "suspicious_url",
		"xss_attempt", "sql_injection_attempt", "path_traversal",
	}
	
	for _, eventType := range eventTypes {
		key := fmt.Sprintf("security_metrics:%s", eventType)
		count, err := s.cache.GetCounter(ctx, key)
		if err == nil {
			metrics.EventCounts[eventType] = int(count)
		}
	}
	
	// Get rate limit violations
	rateLimitKey := "security_metrics:rate_limit_total"
	rateLimitCount, err := s.cache.GetCounter(ctx, rateLimitKey)
	if err == nil {
		metrics.RateLimitViolations = int(rateLimitCount)
	}
	
	// Get blocked requests
	blockedKey := "security_metrics:blocked_requests"
	blockedCount, err := s.cache.GetCounter(ctx, blockedKey)
	if err == nil {
		metrics.BlockedRequests = int(blockedCount)
	}
	
	// Calculate threat distribution
	totalEvents := 0
	for _, count := range metrics.EventCounts {
		totalEvents += count
	}
	metrics.TotalSecurityEvents = totalEvents
	
	// Get top attackers (simplified - in production, this would be more sophisticated)
	metrics.TopAttackers = s.getTopAttackers(ctx, 10)
	metrics.TopEndpoints = s.getTopTargetedEndpoints(ctx, 10)
	
	return metrics, nil
}

// GenerateSecurityReport generates a comprehensive security report
func (s *securityAuditService) GenerateSecurityReport(ctx context.Context, timeRange time.Duration) (*domain.SecurityReport, error) {
	metrics, err := s.GetSecurityMetrics(ctx, timeRange)
	if err != nil {
		return nil, err
	}
	
	report := &domain.SecurityReport{
		GeneratedAt:      time.Now(),
		TimeRange:        timeRange,
		Metrics:          metrics,
		Recommendations:  []string{},
		CriticalFindings: []string{},
		Summary:          "",
	}
	
	// Analyze metrics and generate recommendations
	s.analyzeSecurityMetrics(report, metrics)
	
	return report, nil
}

// CheckSecurityAlerts checks for security conditions that require alerts
func (s *securityAuditService) CheckSecurityAlerts(ctx context.Context) ([]*domain.SecurityAlert, error) {
	alerts := []*domain.SecurityAlert{}
	
	// Check for high rate of failed login attempts
	failedLoginKey := "security_metrics:failed_login"
	failedLogins, err := s.cache.GetCounter(ctx, failedLoginKey)
	if err == nil && failedLogins > 100 { // More than 100 failed logins in the last period
		alerts = append(alerts, &domain.SecurityAlert{
			ID:          fmt.Sprintf("failed_login_%d", time.Now().Unix()),
			Type:        domain.AlertTypeSecurityBreach,
			Severity:    domain.SecuritySeverityHigh,
			Title:       "High Rate of Failed Login Attempts",
			Description: fmt.Sprintf("Detected %d failed login attempts in recent period", failedLogins),
			CreatedAt:   time.Now(),
			Resolved:    false,
		})
	}
	
	// Check for high rate limit violations
	rateLimitKey := "security_metrics:rate_limit_violations_per_hour"
	rateLimitViolations, err := s.cache.GetCounter(ctx, rateLimitKey)
	if err == nil && rateLimitViolations > 1000 {
		alerts = append(alerts, &domain.SecurityAlert{
			ID:          fmt.Sprintf("rate_limit_%d", time.Now().Unix()),
			Type:        domain.AlertTypeRateLimit,
			Severity:    domain.SecuritySeverityMedium,
			Title:       "High Rate Limit Violations",
			Description: fmt.Sprintf("Detected %d rate limit violations in the last hour", rateLimitViolations),
			CreatedAt:   time.Now(),
			Resolved:    false,
		})
	}
	
	// Check for suspicious activity patterns
	suspiciousKey := "security_metrics:suspicious_activity_score"
	suspiciousScore, err := s.cache.GetCounter(ctx, suspiciousKey)
	if err == nil && suspiciousScore > 500 {
		alerts = append(alerts, &domain.SecurityAlert{
			ID:          fmt.Sprintf("suspicious_%d", time.Now().Unix()),
			Type:        domain.AlertTypeSuspiciousActivity,
			Severity:    domain.SecuritySeverityCritical,
			Title:       "High Suspicious Activity Score",
			Description: fmt.Sprintf("Cumulative suspicious activity score: %d", suspiciousScore),
			CreatedAt:   time.Now(),
			Resolved:    false,
		})
	}
	
	return alerts, nil
}

// Helper methods

func (s *securityAuditService) updateSecurityMetrics(ctx context.Context, event *domain.SecurityEvent) {
	// Update general counters
	key := fmt.Sprintf("security_metrics:%s", event.Type)
	s.cache.IncrementCounter(key, 1, 24*time.Hour)
	
	// Update severity counters
	severityKey := fmt.Sprintf("security_metrics:severity:%s", event.Severity)
	s.cache.IncrementCounter(severityKey, 1, 24*time.Hour)
	
	// Update IP-based counters
	ipKey := fmt.Sprintf("security_metrics:ip:%s", event.SourceIP)
	s.cache.IncrementCounter(ipKey, 1, 24*time.Hour)
}

func (s *securityAuditService) trackFailedLoginAttempts(ctx context.Context, sourceIP, email string) {
	// Track by IP
	ipKey := fmt.Sprintf("failed_login:ip:%s", sourceIP)
	count, _ := s.cache.IncrementCounter(ipKey, 1, time.Hour)
	
	// Track by email
	emailKey := fmt.Sprintf("failed_login:email:%s", email)
	s.cache.IncrementCounter(emailKey, 1, time.Hour)
	
	// If too many failed attempts, consider blocking
	if count > 10 {
		s.logger.WithFields(logrus.Fields{
			"source_ip": sourceIP,
			"count":     count,
		}).Error("IP exceeded failed login threshold")
		
		// Could trigger automatic IP blocking here
		blockKey := fmt.Sprintf("blocked_ip:%s", sourceIP)
		s.cache.Set(ctx, blockKey, "failed_login_attempts", time.Hour)
	}
}

func (s *securityAuditService) trackRateLimitViolations(ctx context.Context, sourceIP string) {
	key := fmt.Sprintf("rate_limit_violations:ip:%s", sourceIP)
	count, _ := s.cache.IncrementCounter(key, 1, time.Hour)
	
	// Update global counter
	globalKey := "security_metrics:rate_limit_violations_per_hour"
	s.cache.IncrementCounter(globalKey, 1, time.Hour)
	
	if count > 50 { // More than 50 violations per hour from single IP
		s.logger.WithFields(logrus.Fields{
			"source_ip": sourceIP,
			"count":     count,
		}).Error("IP exceeded rate limit violation threshold")
	}
}

func (s *securityAuditService) updateIPRiskScore(ctx context.Context, sourceIP string, riskScore int) {
	key := fmt.Sprintf("ip_risk_score:%s", sourceIP)
	currentScore, _ := s.cache.GetCounter(ctx, key)
	newScore := currentScore + int64(riskScore)
	
	s.cache.SetCounter(ctx, key, newScore, 24*time.Hour)
	
	if newScore > 100 { // High risk threshold
		s.logger.WithFields(logrus.Fields{
			"source_ip":  sourceIP,
			"risk_score": newScore,
		}).Error("IP exceeded risk score threshold")
		
		// Could trigger automatic actions here
	}
}

func (s *securityAuditService) checkAndTriggerAlerts(ctx context.Context, event *domain.SecurityEvent) {
	// Check for alert conditions based on event type and severity
	if event.Severity == domain.SecuritySeverityCritical {
		// Immediate alert for critical events
		s.logger.WithFields(logrus.Fields{
			"event_type": event.Type,
			"source_ip":  event.SourceIP,
			"endpoint":   event.Endpoint,
		}).Error("CRITICAL security event - immediate attention required")
	}
	
	// Count events per IP to detect patterns
	ipEventKey := fmt.Sprintf("security_events:ip:%s:count", event.SourceIP)
	count, _ := s.cache.IncrementCounter(ipEventKey, 1, time.Hour)
	
	if count > 20 { // More than 20 security events per hour from single IP
		s.logger.WithFields(logrus.Fields{
			"source_ip": event.SourceIP,
			"count":     count,
		}).Error("IP generating excessive security events")
	}
}

func (s *securityAuditService) getTopAttackers(ctx context.Context, limit int) []domain.AttackerInfo {
	// This is a simplified implementation
	// In production, this would query stored events and aggregate data
	return []domain.AttackerInfo{}
}

func (s *securityAuditService) getTopTargetedEndpoints(ctx context.Context, limit int) []domain.EndpointInfo {
	// This is a simplified implementation
	// In production, this would query stored events and aggregate data
	return []domain.EndpointInfo{}
}

func (s *securityAuditService) analyzeSecurityMetrics(report *domain.SecurityReport, metrics *domain.SecurityMetrics) {
	// Analyze failed login attempts
	if failedLogins, exists := metrics.EventCounts["failed_login"]; exists && failedLogins > 100 {
		report.CriticalFindings = append(report.CriticalFindings, 
			fmt.Sprintf("High number of failed login attempts: %d", failedLogins))
		report.Recommendations = append(report.Recommendations, 
			"Consider implementing additional authentication measures like CAPTCHA or account lockouts")
	}
	
	// Analyze rate limit violations
	if metrics.RateLimitViolations > 500 {
		report.CriticalFindings = append(report.CriticalFindings, 
			fmt.Sprintf("High rate limit violations: %d", metrics.RateLimitViolations))
		report.Recommendations = append(report.Recommendations, 
			"Review and potentially tighten rate limiting rules")
	}
	
	// Generate summary
	if len(report.CriticalFindings) == 0 {
		report.Summary = "No critical security issues detected in the analyzed period."
	} else {
		report.Summary = fmt.Sprintf("Detected %d critical security findings requiring attention.", 
			len(report.CriticalFindings))
	}
}

// CreateSecurityEventFromRequest creates a security event from an HTTP request
func CreateSecurityEventFromRequest(r *http.Request, eventType string, severity domain.SecuritySeverity, description string) *domain.SecurityEvent {
	return &domain.SecurityEvent{
		Type:        eventType,
		Severity:    severity,
		SourceIP:    getClientIP(r),
		UserID:      getUserIDFromContext(r.Context()),
		Endpoint:    r.URL.Path,
		Method:      r.Method,
		UserAgent:   r.UserAgent(),
		Timestamp:   time.Now(),
		Description: description,
		Metadata: map[string]interface{}{
			"query_params": r.URL.RawQuery,
			"headers":      filterSensitiveHeaders(r.Header),
		},
	}
}

func getClientIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		return xff
	}
	return r.RemoteAddr
}

func getUserIDFromContext(ctx context.Context) string {
	if userID := ctx.Value("user_id"); userID != nil {
		return fmt.Sprintf("%v", userID)
	}
	return ""
}

func filterSensitiveHeaders(headers http.Header) map[string]string {
	filtered := make(map[string]string)
	sensitiveHeaders := map[string]bool{
		"authorization": true,
		"cookie":        true,
		"x-api-key":     true,
	}
	
	for key, values := range headers {
		lowerKey := strings.ToLower(key)
		if !sensitiveHeaders[lowerKey] && len(values) > 0 {
			filtered[key] = values[0]
		}
	}
	
	return filtered
}