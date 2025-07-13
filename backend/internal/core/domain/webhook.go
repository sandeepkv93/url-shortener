package domain

import (
	"encoding/json"
	"time"
)

// WebhookEvent represents the type of event that triggers a webhook
type WebhookEvent string

const (
	// URL Events
	WebhookEventURLCreated   WebhookEvent = "url.created"
	WebhookEventURLUpdated   WebhookEvent = "url.updated"
	WebhookEventURLDeleted   WebhookEvent = "url.deleted"
	WebhookEventURLClicked   WebhookEvent = "url.clicked"
	WebhookEventURLExpired   WebhookEvent = "url.expired"
	
	// Analytics Events
	WebhookEventAnalyticsThreshold WebhookEvent = "analytics.threshold"
	WebhookEventAnalyticsReport    WebhookEvent = "analytics.report"
	
	// User Events
	WebhookEventUserRegistered WebhookEvent = "user.registered"
	WebhookEventUserUpdated    WebhookEvent = "user.updated"
	
	// System Events
	WebhookEventSystemError   WebhookEvent = "system.error"
	WebhookEventSystemAlert   WebhookEvent = "system.alert"
)

// WebhookStatus represents the current status of a webhook
type WebhookStatus string

const (
	WebhookStatusActive    WebhookStatus = "active"
	WebhookStatusInactive  WebhookStatus = "inactive"
	WebhookStatusFailed    WebhookStatus = "failed"
	WebhookStatusSuspended WebhookStatus = "suspended"
)

// WebhookDeliveryStatus represents the status of a webhook delivery attempt
type WebhookDeliveryStatus string

const (
	WebhookDeliveryStatusPending   WebhookDeliveryStatus = "pending"
	WebhookDeliveryStatusSuccess   WebhookDeliveryStatus = "success"
	WebhookDeliveryStatusFailed    WebhookDeliveryStatus = "failed"
	WebhookDeliveryStatusRetrying  WebhookDeliveryStatus = "retrying"
	WebhookDeliveryStatusAbandoned WebhookDeliveryStatus = "abandoned"
)

// Webhook represents a webhook configuration
type Webhook struct {
	ID          uint64          `json:"id" gorm:"primaryKey"`
	UserID      uint64          `json:"user_id" gorm:"not null;index"`
	Name        string          `json:"name" gorm:"size:255;not null"`
	URL         string          `json:"url" gorm:"size:2048;not null"`
	Events      []WebhookEvent  `json:"events" gorm:"type:text"`
	Secret      string          `json:"-" gorm:"size:255"`
	Status      WebhookStatus   `json:"status" gorm:"default:'active'"`
	
	// Configuration
	MaxRetries      int           `json:"max_retries" gorm:"default:3"`
	TimeoutSeconds  int           `json:"timeout_seconds" gorm:"default:30"`
	RetryBackoffMS  int           `json:"retry_backoff_ms" gorm:"default:1000"`
	
	// Statistics
	TotalDeliveries   int64     `json:"total_deliveries" gorm:"default:0"`
	SuccessDeliveries int64     `json:"success_deliveries" gorm:"default:0"`
	FailedDeliveries  int64     `json:"failed_deliveries" gorm:"default:0"`
	LastDeliveryAt    *time.Time `json:"last_delivery_at"`
	LastSuccessAt     *time.Time `json:"last_success_at"`
	LastFailureAt     *time.Time `json:"last_failure_at"`
	
	// Metadata
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// Relationships
	User       User                `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Deliveries []WebhookDelivery   `json:"deliveries,omitempty" gorm:"foreignKey:WebhookID"`
}

// WebhookDelivery represents a webhook delivery attempt
type WebhookDelivery struct {
	ID          uint64                `json:"id" gorm:"primaryKey"`
	WebhookID   uint64                `json:"webhook_id" gorm:"not null;index"`
	EventType   WebhookEvent          `json:"event_type" gorm:"size:255;not null"`
	Status      WebhookDeliveryStatus `json:"status" gorm:"default:'pending'"`
	
	// Request details
	RequestURL     string            `json:"request_url" gorm:"size:2048"`
	RequestHeaders json.RawMessage   `json:"request_headers" gorm:"type:jsonb"`
	RequestBody    json.RawMessage   `json:"request_body" gorm:"type:jsonb"`
	
	// Response details
	ResponseStatus  int             `json:"response_status"`
	ResponseHeaders json.RawMessage `json:"response_headers" gorm:"type:jsonb"`
	ResponseBody    string          `json:"response_body" gorm:"type:text"`
	
	// Timing
	Duration      int64     `json:"duration_ms"`
	AttemptCount  int       `json:"attempt_count" gorm:"default:1"`
	NextRetryAt   *time.Time `json:"next_retry_at"`
	
	// Error handling
	ErrorMessage  string    `json:"error_message" gorm:"type:text"`
	
	// Metadata
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// Relationships
	Webhook Webhook `json:"webhook,omitempty" gorm:"foreignKey:WebhookID"`
}

// WebhookPayload represents the payload sent to webhook endpoints
type WebhookPayload struct {
	ID        string                 `json:"id"`
	Event     WebhookEvent           `json:"event"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
	UserID    uint64                 `json:"user_id"`
	Version   string                 `json:"version"`
}

// WebhookEventData contains event-specific data structures
type WebhookEventData struct {
	// URL-related events
	URL *ShortURL `json:"url,omitempty"`
	
	// Click-related events
	Click *Click `json:"click,omitempty"`
	
	// Analytics-related events
	Analytics *AnalyticsData `json:"analytics,omitempty"`
	
	// User-related events
	User *User `json:"user,omitempty"`
	
	// Custom data
	Custom map[string]interface{} `json:"custom,omitempty"`
}

// WebhookDeliveryStats represents statistics for webhook deliveries
type WebhookDeliveryStats struct {
	TotalDeliveries   int64   `json:"total_deliveries"`
	SuccessDeliveries int64   `json:"success_deliveries"`
	FailedDeliveries  int64   `json:"failed_deliveries"`
	SuccessRate       float64 `json:"success_rate"`
	AverageResponseTime int64 `json:"average_response_time_ms"`
	LastDeliveryAt    *time.Time `json:"last_delivery_at"`
	LastSuccessAt     *time.Time `json:"last_success_at"`
	LastFailureAt     *time.Time `json:"last_failure_at"`
}

// AnalyticsData represents analytics data for webhook events
type AnalyticsData struct {
	URLCode       string `json:"url_code"`
	TotalClicks   int64  `json:"total_clicks"`
	UniqueClicks  int64  `json:"unique_clicks"`
	Period        string `json:"period"`
	Threshold     int64  `json:"threshold,omitempty"`
	ThresholdType string `json:"threshold_type,omitempty"`
}

// WebhookFilter represents filtering criteria for webhooks
type WebhookFilter struct {
	UserID    *uint64         `json:"user_id,omitempty"`
	Status    *WebhookStatus  `json:"status,omitempty"`
	Events    []WebhookEvent  `json:"events,omitempty"`
	Active    *bool           `json:"active,omitempty"`
}

// WebhookDeliveryFilter represents filtering criteria for webhook deliveries
type WebhookDeliveryFilter struct {
	WebhookID *uint64                `json:"webhook_id,omitempty"`
	Status    *WebhookDeliveryStatus `json:"status,omitempty"`
	EventType *WebhookEvent          `json:"event_type,omitempty"`
	From      *time.Time             `json:"from,omitempty"`
	To        *time.Time             `json:"to,omitempty"`
}

// WebhookCreateRequest represents a request to create a webhook
type WebhookCreateRequest struct {
	Name           string         `json:"name" validate:"required,min=1,max=255"`
	URL            string         `json:"url" validate:"required,url,max=2048"`
	Events         []WebhookEvent `json:"events" validate:"required,min=1"`
	Secret         string         `json:"secret,omitempty" validate:"max=255"`
	MaxRetries     *int           `json:"max_retries,omitempty" validate:"omitempty,min=0,max=10"`
	TimeoutSeconds *int           `json:"timeout_seconds,omitempty" validate:"omitempty,min=1,max=300"`
}

// WebhookUpdateRequest represents a request to update a webhook
type WebhookUpdateRequest struct {
	Name           *string        `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	URL            *string        `json:"url,omitempty" validate:"omitempty,url,max=2048"`
	Events         []WebhookEvent `json:"events,omitempty" validate:"omitempty,min=1"`
	Secret         *string        `json:"secret,omitempty" validate:"omitempty,max=255"`
	Status         *WebhookStatus `json:"status,omitempty"`
	MaxRetries     *int           `json:"max_retries,omitempty" validate:"omitempty,min=0,max=10"`
	TimeoutSeconds *int           `json:"timeout_seconds,omitempty" validate:"omitempty,min=1,max=300"`
}

// IsActive returns whether the webhook is currently active
func (w *Webhook) IsActive() bool {
	return w.Status == WebhookStatusActive
}

// ShouldReceiveEvent returns whether the webhook should receive the given event
func (w *Webhook) ShouldReceiveEvent(event WebhookEvent) bool {
	if !w.IsActive() {
		return false
	}
	
	for _, e := range w.Events {
		if e == event {
			return true
		}
	}
	return false
}

// CanRetry returns whether the delivery can be retried
func (wd *WebhookDelivery) CanRetry() bool {
	return wd.Status == WebhookDeliveryStatusFailed && 
		   wd.AttemptCount < 10 && // Hard limit to prevent infinite retries
		   (wd.NextRetryAt == nil || time.Now().After(*wd.NextRetryAt))
}

// ShouldAbandon returns whether the delivery should be abandoned
func (wd *WebhookDelivery) ShouldAbandon(maxRetries int) bool {
	return wd.AttemptCount >= maxRetries
}

// GetSuccessRate returns the success rate of the webhook as a percentage
func (w *Webhook) GetSuccessRate() float64 {
	if w.TotalDeliveries == 0 {
		return 0.0
	}
	return float64(w.SuccessDeliveries) / float64(w.TotalDeliveries) * 100.0
}

// ValidateEvents validates that all events are valid
func ValidateWebhookEvents(events []WebhookEvent) error {
	validEvents := map[WebhookEvent]bool{
		WebhookEventURLCreated:         true,
		WebhookEventURLUpdated:         true,
		WebhookEventURLDeleted:         true,
		WebhookEventURLClicked:         true,
		WebhookEventURLExpired:         true,
		WebhookEventAnalyticsThreshold: true,
		WebhookEventAnalyticsReport:    true,
		WebhookEventUserRegistered:     true,
		WebhookEventUserUpdated:        true,
		WebhookEventSystemError:        true,
		WebhookEventSystemAlert:        true,
	}
	
	for _, event := range events {
		if !validEvents[event] {
			return &ValidationError{
				Field:   "events",
				Message: "Invalid webhook event: " + string(event),
			}
		}
	}
	
	if len(events) == 0 {
		return &ValidationError{
			Field:   "events",
			Message: "At least one event must be specified",
		}
	}
	
	return nil
}