package domain

import (
	"time"
)

type AnalyticsDigest struct {
	UserID      uint      `json:"user_id"`
	Period      string    `json:"period"`     // daily, weekly, monthly
	TotalClicks int64     `json:"total_clicks"`
	TotalURLs   int64     `json:"total_urls"`
	TopURLs     []TopURLStat `json:"top_urls"`
	Summary     string    `json:"summary"`
	GeneratedAt time.Time `json:"generated_at"`
}

type ClickAlert struct {
	UserID      uint      `json:"user_id"`
	ShortURLID  uint      `json:"short_url_id"`
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	Threshold   int64     `json:"threshold"`
	CurrentCount int64    `json:"current_count"`
	AlertType   string    `json:"alert_type"` // milestone, threshold_exceeded, spike_detected
	TriggeredAt time.Time `json:"triggered_at"`
}

type SecurityAlert struct {
	UserID      uint      `json:"user_id"`
	AlertType   string    `json:"alert_type"` // suspicious_login, password_changed, account_locked
	Description string    `json:"description"`
	IPAddress   string    `json:"ip_address"`
	UserAgent   string    `json:"user_agent"`
	Location    string    `json:"location"`
	Severity    string    `json:"severity"` // low, medium, high, critical
	TriggeredAt time.Time `json:"triggered_at"`
	Action      string    `json:"action,omitempty"` // account_locked, password_reset_required
}

type NotificationPreferences struct {
	UserID              uint `json:"user_id"`
	EmailDigests        bool `json:"email_digests"`
	ClickAlerts         bool `json:"click_alerts"`
	SecurityAlerts      bool `json:"security_alerts"`
	MaintenanceNotices  bool `json:"maintenance_notices"`
	MarketingEmails     bool `json:"marketing_emails"`
	DigestFrequency     string `json:"digest_frequency"` // daily, weekly, monthly
	ClickAlertThreshold int64  `json:"click_alert_threshold"`
}

type EmailTemplate struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Subject     string    `json:"subject"`
	HTMLContent string    `json:"html_content"`
	TextContent string    `json:"text_content"`
	Variables   []string  `json:"variables"`
	Category    string    `json:"category"` // welcome, reset, alert, digest
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type NotificationLog struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	Type        string    `json:"type"`
	Channel     string    `json:"channel"` // email, sms, push
	Status      string    `json:"status"`  // sent, failed, pending
	Subject     string    `json:"subject"`
	Content     string    `json:"content"`
	Error       string    `json:"error,omitempty"`
	SentAt      *time.Time `json:"sent_at"`
	CreatedAt   time.Time `json:"created_at"`
}