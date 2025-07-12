package domain

import (
	"fmt"
	"time"
)

// SecuritySeverity represents the severity level of a security event
type SecuritySeverity string

const (
	SecuritySeverityLow      SecuritySeverity = "low"
	SecuritySeverityMedium   SecuritySeverity = "medium"
	SecuritySeverityHigh     SecuritySeverity = "high"
	SecuritySeverityCritical SecuritySeverity = "critical"
)

// SecurityEvent represents a general security event
type SecurityEvent struct {
	Type        string                 `json:"type"`
	Severity    SecuritySeverity       `json:"severity"`
	SourceIP    string                 `json:"source_ip"`
	UserID      string                 `json:"user_id,omitempty"`
	Endpoint    string                 `json:"endpoint"`
	Method      string                 `json:"method"`
	UserAgent   string                 `json:"user_agent"`
	Timestamp   time.Time              `json:"timestamp"`
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// AuthAction represents different authentication actions
type AuthAction string

const (
	AuthActionLogin           AuthAction = "login"
	AuthActionRegister        AuthAction = "register"
	AuthActionLogout          AuthAction = "logout"
	AuthActionPasswordReset   AuthAction = "password_reset"
	AuthActionPasswordChange  AuthAction = "password_change"
	AuthActionTokenRefresh    AuthAction = "token_refresh"
	AuthActionAccountLockout  AuthAction = "account_lockout"
)

// AuthEvent represents authentication-related events
type AuthEvent struct {
	Action        AuthAction `json:"action"`
	UserID        string     `json:"user_id,omitempty"`
	Email         string     `json:"email"`
	SourceIP      string     `json:"source_ip"`
	UserAgent     string     `json:"user_agent"`
	Success       bool       `json:"success"`
	FailureReason string     `json:"failure_reason,omitempty"`
	Timestamp     time.Time  `json:"timestamp"`
	SessionID     string     `json:"session_id,omitempty"`
}

// RateLimitEvent represents rate limiting events
type RateLimitEvent struct {
	LimitType    string        `json:"limit_type"`
	SourceIP     string        `json:"source_ip"`
	UserID       string        `json:"user_id,omitempty"`
	Endpoint     string        `json:"endpoint"`
	Method       string        `json:"method"`
	CurrentCount int           `json:"current_count"`
	Limit        int           `json:"limit"`
	Window       time.Duration `json:"window"`
	Timestamp    time.Time     `json:"timestamp"`
}

// SuspiciousActivityEvent represents suspicious activity
type SuspiciousActivityEvent struct {
	ActivityType     string    `json:"activity_type"`
	SourceIP         string    `json:"source_ip"`
	UserID           string    `json:"user_id,omitempty"`
	Endpoint         string    `json:"endpoint"`
	Indicators       []string  `json:"indicators"`
	RiskScore        int       `json:"risk_score"`
	AutomaticAction  string    `json:"automatic_action,omitempty"`
	Timestamp        time.Time `json:"timestamp"`
	Description      string    `json:"description"`
}

// SecurityMetrics contains security-related metrics
type SecurityMetrics struct {
	TimeRange            time.Duration            `json:"time_range"`
	CollectedAt          time.Time                `json:"collected_at"`
	TotalSecurityEvents  int                      `json:"total_security_events"`
	EventCounts          map[string]int           `json:"event_counts"`
	ThreatLevels         map[string]int           `json:"threat_levels"`
	RateLimitViolations  int                      `json:"rate_limit_violations"`
	BlockedRequests      int                      `json:"blocked_requests"`
	TopAttackers         []AttackerInfo           `json:"top_attackers"`
	TopEndpoints         []EndpointInfo           `json:"top_endpoints"`
}

// AttackerInfo contains information about potential attackers
type AttackerInfo struct {
	IP            string    `json:"ip"`
	Country       string    `json:"country,omitempty"`
	EventCount    int       `json:"event_count"`
	ThreatScore   int       `json:"threat_score"`
	LastActivity  time.Time `json:"last_activity"`
	AttackTypes   []string  `json:"attack_types"`
}

// EndpointInfo contains information about targeted endpoints
type EndpointInfo struct {
	Endpoint     string    `json:"endpoint"`
	Method       string    `json:"method"`
	AttackCount  int       `json:"attack_count"`
	LastAttack   time.Time `json:"last_attack"`
	AttackTypes  []string  `json:"attack_types"`
}

// SecurityReport contains a comprehensive security analysis report
type SecurityReport struct {
	GeneratedAt      time.Time        `json:"generated_at"`
	TimeRange        time.Duration    `json:"time_range"`
	Metrics          *SecurityMetrics `json:"metrics"`
	Summary          string           `json:"summary"`
	CriticalFindings []string         `json:"critical_findings"`
	Recommendations  []string         `json:"recommendations"`
}

// AlertType represents different types of security alerts
type AlertType string

const (
	AlertTypeSecurityBreach     AlertType = "security_breach"
	AlertTypeRateLimit          AlertType = "rate_limit"
	AlertTypeSuspiciousActivity AlertType = "suspicious_activity"
	AlertTypeAuthentication     AlertType = "authentication"
	AlertTypeSystemHealth       AlertType = "system_health"
)

// SecurityAlert already exists in notification.go - using that definition

// AlertAction represents an action taken in response to an alert
type AlertAction struct {
	Type        string    `json:"type"`
	Description string    `json:"description"`
	PerformedAt time.Time `json:"performed_at"`
	PerformedBy string    `json:"performed_by"`
	Result      string    `json:"result"`
}

// SecurityPolicy represents security policy configuration
type SecurityPolicy struct {
	ID                   string            `json:"id"`
	Name                 string            `json:"name"`
	Description          string            `json:"description"`
	Enabled              bool              `json:"enabled"`
	Rules                []SecurityRule    `json:"rules"`
	Actions              []PolicyAction    `json:"actions"`
	CreatedAt            time.Time         `json:"created_at"`
	UpdatedAt            time.Time         `json:"updated_at"`
	CreatedBy            string            `json:"created_by"`
}

// SecurityRule represents a security rule within a policy
type SecurityRule struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Condition   string                 `json:"condition"`
	Parameters  map[string]interface{} `json:"parameters"`
	Enabled     bool                   `json:"enabled"`
	Description string                 `json:"description"`
}

// PolicyAction represents an action to take when a security rule is triggered
type PolicyAction struct {
	Type        string                 `json:"type"`
	Parameters  map[string]interface{} `json:"parameters"`
	Enabled     bool                   `json:"enabled"`
	Description string                 `json:"description"`
}

// SecurityIncident represents a security incident
type SecurityIncident struct {
	ID               string            `json:"id"`
	Title            string            `json:"title"`
	Description      string            `json:"description"`
	Severity         SecuritySeverity  `json:"severity"`
	Status           IncidentStatus    `json:"status"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	ResolvedAt       *time.Time        `json:"resolved_at,omitempty"`
	AssignedTo       string            `json:"assigned_to,omitempty"`
	Reporter         string            `json:"reporter"`
	AffectedSystems  []string          `json:"affected_systems"`
	RelatedEvents    []string          `json:"related_events"`
	Actions          []IncidentAction  `json:"actions"`
	Tags             []string          `json:"tags"`
	Metadata         map[string]string `json:"metadata,omitempty"`
}

// IncidentStatus represents the status of a security incident
type IncidentStatus string

const (
	IncidentStatusNew        IncidentStatus = "new"
	IncidentStatusInProgress IncidentStatus = "in_progress"
	IncidentStatusResolved   IncidentStatus = "resolved"
	IncidentStatusClosed     IncidentStatus = "closed"
	IncidentStatusReopened   IncidentStatus = "reopened"
)

// IncidentAction represents an action taken during incident response
type IncidentAction struct {
	Type        string    `json:"type"`
	Description string    `json:"description"`
	PerformedAt time.Time `json:"performed_at"`
	PerformedBy string    `json:"performed_by"`
	Result      string    `json:"result"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// SecurityDashboard represents security dashboard data
type SecurityDashboard struct {
	GeneratedAt          time.Time         `json:"generated_at"`
	TimeRange            time.Duration     `json:"time_range"`
	OverallSecurityScore int               `json:"overall_security_score"`
	ThreatLevel          SecuritySeverity  `json:"threat_level"`
	ActiveAlerts         int               `json:"active_alerts"`
	RecentIncidents      int               `json:"recent_incidents"`
	BlockedAttacks       int               `json:"blocked_attacks"`
	SecurityMetrics      *SecurityMetrics  `json:"security_metrics"`
	TopThreats           []ThreatInfo      `json:"top_threats"`
	SystemHealth         *ComponentHealth  `json:"system_health"` // Use ComponentHealth instead of SystemHealth
}

// ThreatInfo represents information about a specific threat
type ThreatInfo struct {
	Type        string    `json:"type"`
	Count       int       `json:"count"`
	Severity    SecuritySeverity `json:"severity"`
	LastSeen    time.Time `json:"last_seen"`
	Description string    `json:"description"`
	Mitigation  string    `json:"mitigation"`
}

// SystemHealth and ComponentHealth already exist in health.go - using those definitions

// Validation methods

// Validate validates a SecurityEvent
func (e *SecurityEvent) Validate() error {
	if e.Type == "" {
		return fmt.Errorf("security event type is required")
	}
	if e.SourceIP == "" {
		return fmt.Errorf("source IP is required")
	}
	if e.Endpoint == "" {
		return fmt.Errorf("endpoint is required")
	}
	return nil
}

// Validate validates an AuthEvent
func (e *AuthEvent) Validate() error {
	if e.Action == "" {
		return fmt.Errorf("auth action is required")
	}
	if e.Email == "" && e.UserID == "" {
		return fmt.Errorf("either email or user ID is required")
	}
	if e.SourceIP == "" {
		return fmt.Errorf("source IP is required")
	}
	return nil
}

// Validate validates a SecurityAlert - using the existing SecurityAlert from notification.go
func (a *SecurityAlert) Validate() error {
	if a.Type == "" {
		return fmt.Errorf("alert type is required")
	}
	// Note: Using the SecurityAlert from notification.go which might have different fields
	return nil
}

// Helper methods

// IsHighSeverity returns true if the event is high or critical severity
func (e *SecurityEvent) IsHighSeverity() bool {
	return e.Severity == SecuritySeverityHigh || e.Severity == SecuritySeverityCritical
}

// IsCritical returns true if the event is critical severity
func (e *SecurityEvent) IsCritical() bool {
	return e.Severity == SecuritySeverityCritical
}

// IsAuthFailure returns true if this is a failed authentication event
func (e *AuthEvent) IsAuthFailure() bool {
	return !e.Success
}

// GetSeverityLevel returns a numeric representation of severity for comparison
func (s SecuritySeverity) GetSeverityLevel() int {
	switch s {
	case SecuritySeverityLow:
		return 1
	case SecuritySeverityMedium:
		return 2
	case SecuritySeverityHigh:
		return 3
	case SecuritySeverityCritical:
		return 4
	default:
		return 0
	}
}

// Note: SecurityAlert methods removed since we're using the SecurityAlert from notification.go
// which has different fields