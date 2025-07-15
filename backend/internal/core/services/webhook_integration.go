package services

import (
	"context"
	"time"

	"url-shortener/internal/core/domain"
	"url-shortener/internal/core/ports"
)

// WebhookIntegrationService provides convenient methods to trigger webhooks from other services
type WebhookIntegrationService struct {
	webhookService ports.WebhookService
	logger         ports.Logger
}

func NewWebhookIntegrationService(
	webhookService ports.WebhookService,
	logger ports.Logger,
) *WebhookIntegrationService {
	return &WebhookIntegrationService{
		webhookService: webhookService,
		logger:         logger,
	}
}

// URL Events

func (s *WebhookIntegrationService) TriggerURLCreated(ctx context.Context, url *domain.ShortURL) {
	s.triggerEventAsync(ctx, domain.WebhookEventURLCreated, map[string]interface{}{
		"url": url,
	}, uint64(url.UserID))
}

func (s *WebhookIntegrationService) TriggerURLUpdated(ctx context.Context, url *domain.ShortURL, changes map[string]interface{}) {
	s.triggerEventAsync(ctx, domain.WebhookEventURLUpdated, map[string]interface{}{
		"url":     url,
		"changes": changes,
	}, uint64(url.UserID))
}

func (s *WebhookIntegrationService) TriggerURLDeleted(ctx context.Context, url *domain.ShortURL) {
	s.triggerEventAsync(ctx, domain.WebhookEventURLDeleted, map[string]interface{}{
		"url": map[string]interface{}{
			"id":           url.ID,
			"short_code":   url.ShortCode,
			"original_url": url.OriginalURL,
			"title":        url.Title,
			"deleted_at":   time.Now(),
		},
	}, uint64(url.UserID))
}

func (s *WebhookIntegrationService) TriggerURLClicked(ctx context.Context, url *domain.ShortURL, click *domain.Click) {
	s.triggerEventAsync(ctx, domain.WebhookEventURLClicked, map[string]interface{}{
		"url":   url,
		"click": click,
		"stats": map[string]interface{}{
			"total_clicks": url.ClickCount + 1,
		},
	}, uint64(url.UserID))
}

func (s *WebhookIntegrationService) TriggerURLExpired(ctx context.Context, url *domain.ShortURL) {
	s.triggerEventAsync(ctx, domain.WebhookEventURLExpired, map[string]interface{}{
		"url": url,
		"expired_at": time.Now(),
	}, uint64(url.UserID))
}

// Analytics Events

func (s *WebhookIntegrationService) TriggerAnalyticsThreshold(ctx context.Context, userID uint64, urlCode string, threshold *domain.AnalyticsThreshold) {
	s.triggerEventAsync(ctx, domain.WebhookEventAnalyticsThreshold, map[string]interface{}{
		"url_code":        urlCode,
		"threshold":       threshold,
		"current_value":   threshold.CurrentValue,
		"threshold_value": threshold.ThresholdValue,
		"threshold_type":  threshold.ThresholdType,
		"triggered_at":    time.Now(),
	}, userID)
}

func (s *WebhookIntegrationService) TriggerAnalyticsReport(ctx context.Context, userID uint64, report *domain.AnalyticsReport) {
	s.triggerEventAsync(ctx, domain.WebhookEventAnalyticsReport, map[string]interface{}{
		"report":       report,
		"generated_at": time.Now(),
	}, userID)
}

// User Events

func (s *WebhookIntegrationService) TriggerUserRegistered(ctx context.Context, user *domain.User) {
	// Remove sensitive information before sending
	safeUser := map[string]interface{}{
		"id":         user.ID,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	}
	
	s.triggerEventAsync(ctx, domain.WebhookEventUserRegistered, map[string]interface{}{
		"user": safeUser,
	}, uint64(user.ID))
}

func (s *WebhookIntegrationService) TriggerUserUpdated(ctx context.Context, user *domain.User, changes map[string]interface{}) {
	// Remove sensitive information before sending
	safeUser := map[string]interface{}{
		"id":         user.ID,
		"email":      user.Email,
		"updated_at": user.UpdatedAt,
	}
	
	s.triggerEventAsync(ctx, domain.WebhookEventUserUpdated, map[string]interface{}{
		"user":    safeUser,
		"changes": changes,
	}, uint64(user.ID))
}

// System Events

func (s *WebhookIntegrationService) TriggerSystemError(ctx context.Context, error *domain.SystemError) {
	s.triggerGlobalEventAsync(ctx, domain.WebhookEventSystemError, map[string]interface{}{
		"error":      error,
		"time":       time.Now(),
		"severity":   error.Severity,
		"component":  error.Component,
		"message":    error.Message,
	})
}

func (s *WebhookIntegrationService) TriggerSystemAlert(ctx context.Context, alert *domain.SystemAlert) {
	s.triggerGlobalEventAsync(ctx, domain.WebhookEventSystemAlert, map[string]interface{}{
		"alert":     alert,
		"time":      time.Now(),
		"type":      alert.Type,
		"message":   alert.Message,
		"severity":  alert.Severity,
	})
}

// Batch Events

func (s *WebhookIntegrationService) TriggerBulkURLsCreated(ctx context.Context, urls []*domain.ShortURL, userID uint64) {
	s.triggerEventAsync(ctx, domain.WebhookEventURLCreated, map[string]interface{}{
		"urls":  urls,
		"count": len(urls),
		"bulk":  true,
	}, userID)
}

func (s *WebhookIntegrationService) TriggerBulkURLsDeleted(ctx context.Context, urlIDs []uint, userID uint64) {
	s.triggerEventAsync(ctx, domain.WebhookEventURLDeleted, map[string]interface{}{
		"url_ids":    urlIDs,
		"count":      len(urlIDs),
		"bulk":       true,
		"deleted_at": time.Now(),
	}, userID)
}

// Custom Events (for future extensibility)

func (s *WebhookIntegrationService) TriggerCustomEvent(ctx context.Context, eventType string, data interface{}, userID uint64) {
	// This allows for custom events not predefined in the system
	s.triggerEventAsync(ctx, domain.WebhookEvent(eventType), data, userID)
}

func (s *WebhookIntegrationService) TriggerCustomGlobalEvent(ctx context.Context, eventType string, data interface{}) {
	// This allows for custom global events not predefined in the system
	s.triggerGlobalEventAsync(ctx, domain.WebhookEvent(eventType), data)
}

// Helper Methods

func (s *WebhookIntegrationService) triggerEventAsync(ctx context.Context, event domain.WebhookEvent, data interface{}, userID uint64) {
	go func() {
		// Create a new context for the background operation
		bgCtx := context.Background()
		
		if err := s.webhookService.TriggerEventForUser(bgCtx, event, data, userID); err != nil {
			s.logger.ErrorContext(ctx, "Failed to trigger webhook event", map[string]interface{}{
				"event":   event,
				"user_id": userID,
				"error":   err.Error(),
			})
		}
	}()
}

func (s *WebhookIntegrationService) triggerGlobalEventAsync(ctx context.Context, event domain.WebhookEvent, data interface{}) {
	go func() {
		// Create a new context for the background operation
		bgCtx := context.Background()
		
		if err := s.webhookService.TriggerEventGlobally(bgCtx, event, data); err != nil {
			s.logger.ErrorContext(ctx, "Failed to trigger global webhook event", map[string]interface{}{
				"event": event,
				"error": err.Error(),
			})
		}
	}()
}

// Webhook Event Builders

func (s *WebhookIntegrationService) BuildURLEventData(url *domain.ShortURL) map[string]interface{} {
	return map[string]interface{}{
		"id":           url.ID,
		"short_code":   url.ShortCode,
		"original_url": url.OriginalURL,
		"title":        url.Title,
		"description":  url.Description,
		"click_count":  url.ClickCount,
		"is_active":    url.IsActive,
		"expires_at":   url.ExpiresAt,
		"created_at":   url.CreatedAt,
		"updated_at":   url.UpdatedAt,
	}
}

func (s *WebhookIntegrationService) BuildClickEventData(click *domain.Click) map[string]interface{} {
	return map[string]interface{}{
		"id":         click.ID,
		"ip_address": click.IPAddress,
		"user_agent": click.UserAgent,
		"referer":    click.Referer,
		"country":    click.Country,
		"city":       click.City,
		"clicked_at": click.ClickedAt,
	}
}

// Health and Monitoring

func (s *WebhookIntegrationService) GetPendingWebhookCount(ctx context.Context) (int, error) {
	// This would require implementing a method to get pending webhooks count
	// For now, we'll delegate to the webhook service to process pending retries
	return 0, s.webhookService.ProcessPendingRetries(ctx)
}

func (s *WebhookIntegrationService) ProcessPendingWebhooks(ctx context.Context) error {
	return s.webhookService.ProcessPendingRetries(ctx)
}

func (s *WebhookIntegrationService) CleanupOldWebhookDeliveries(ctx context.Context, days int) (int64, error) {
	return s.webhookService.CleanupOldDeliveries(ctx, days)
}

// Additional domain types needed for webhooks (these would be added to domain package)

// These types would be defined in the domain package
type AnalyticsThreshold struct {
	ThresholdType  string      `json:"threshold_type"`
	ThresholdValue interface{} `json:"threshold_value"`
	CurrentValue   interface{} `json:"current_value"`
	URLCode        string      `json:"url_code,omitempty"`
	UserID         uint64      `json:"user_id"`
}

type AnalyticsReport struct {
	ID          uint64    `json:"id"`
	UserID      uint64    `json:"user_id"`
	ReportType  string    `json:"report_type"`
	Period      string    `json:"period"`
	GeneratedAt time.Time `json:"generated_at"`
	Data        map[string]interface{} `json:"data"`
}

type SystemError struct {
	ID        string                 `json:"id"`
	Component string                 `json:"component"`
	Message   string                 `json:"message"`
	Severity  string                 `json:"severity"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

type SystemAlert struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Message   string                 `json:"message"`
	Severity  string                 `json:"severity"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}