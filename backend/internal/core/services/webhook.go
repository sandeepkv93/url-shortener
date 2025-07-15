package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"url-shortener/internal/core/domain"
	"url-shortener/internal/core/ports"
)

type webhookService struct {
	webhookRepo         ports.WebhookRepository
	deliveryRepo        ports.WebhookDeliveryRepository
	deliveryService     ports.WebhookDeliveryService
	logger              ports.Logger
}

func NewWebhookService(
	webhookRepo ports.WebhookRepository,
	deliveryRepo ports.WebhookDeliveryRepository,
	deliveryService ports.WebhookDeliveryService,
	logger ports.Logger,
) ports.WebhookService {
	return &webhookService{
		webhookRepo:     webhookRepo,
		deliveryRepo:    deliveryRepo,
		deliveryService: deliveryService,
		logger:          logger,
	}
}

func (s *webhookService) CreateWebhook(ctx context.Context, userID uint64, req domain.WebhookCreateRequest) (*domain.Webhook, error) {
	// Validate the request
	if err := domain.ValidateWebhookEvents(req.Events); err != nil {
		return nil, err
	}
	
	// Validate the URL
	if err := s.deliveryService.ValidateWebhookURL(ctx, req.URL); err != nil {
		return nil, &domain.ValidationError{
			Field:   "url",
			Message: "Invalid webhook URL: " + err.Error(),
		}
	}
	
	// Generate secret if not provided
	secret := req.Secret
	if secret == "" {
		var err error
		secret, err = s.generateSecret()
		if err != nil {
			return nil, fmt.Errorf("failed to generate webhook secret: %w", err)
		}
	}
	
	// Set default values
	maxRetries := 3
	if req.MaxRetries != nil {
		maxRetries = *req.MaxRetries
	}
	
	timeoutSeconds := 30
	if req.TimeoutSeconds != nil {
		timeoutSeconds = *req.TimeoutSeconds
	}
	
	webhook := &domain.Webhook{
		UserID:          userID,
		Name:            req.Name,
		URL:             req.URL,
		Events:          req.Events,
		Secret:          secret,
		Status:          domain.WebhookStatusActive,
		MaxRetries:      maxRetries,
		TimeoutSeconds:  timeoutSeconds,
		RetryBackoffMS:  1000, // 1 second default backoff
	}
	
	if err := s.webhookRepo.Create(ctx, webhook); err != nil {
		return nil, fmt.Errorf("failed to create webhook: %w", err)
	}
	
	s.logger.InfoContext(ctx, "Webhook created successfully", map[string]interface{}{
		"webhook_id": webhook.ID,
		"user_id":    userID,
		"url":        webhook.URL,
		"events":     webhook.Events,
	})
	
	return webhook, nil
}

func (s *webhookService) GetWebhook(ctx context.Context, webhookID uint64, userID uint64) (*domain.Webhook, error) {
	webhook, err := s.webhookRepo.GetByID(ctx, webhookID)
	if err != nil {
		return nil, err
	}
	
	if webhook.UserID != userID {
		return nil, domain.NewNotFoundError("webhook")
	}
	
	return webhook, nil
}

func (s *webhookService) GetUserWebhooks(ctx context.Context, userID uint64, offset, limit int) ([]*domain.Webhook, int64, error) {
	return s.webhookRepo.GetByUserID(ctx, userID, offset, limit)
}

func (s *webhookService) UpdateWebhook(ctx context.Context, webhookID uint64, userID uint64, req domain.WebhookUpdateRequest) (*domain.Webhook, error) {
	webhook, err := s.GetWebhook(ctx, webhookID, userID)
	if err != nil {
		return nil, err
	}
	
	// Update fields if provided
	if req.Name != nil {
		webhook.Name = *req.Name
	}
	
	if req.URL != nil {
		if err := s.deliveryService.ValidateWebhookURL(ctx, *req.URL); err != nil {
			return nil, &domain.ValidationError{
				Field:   "url",
				Message: "Invalid webhook URL: " + err.Error(),
			}
		}
		webhook.URL = *req.URL
	}
	
	if req.Events != nil {
		if err := domain.ValidateWebhookEvents(req.Events); err != nil {
			return nil, err
		}
		webhook.Events = req.Events
	}
	
	if req.Secret != nil {
		webhook.Secret = *req.Secret
	}
	
	if req.Status != nil {
		webhook.Status = *req.Status
	}
	
	if req.MaxRetries != nil {
		webhook.MaxRetries = *req.MaxRetries
	}
	
	if req.TimeoutSeconds != nil {
		webhook.TimeoutSeconds = *req.TimeoutSeconds
	}
	
	if err := s.webhookRepo.Update(ctx, webhook); err != nil {
		return nil, fmt.Errorf("failed to update webhook: %w", err)
	}
	
	s.logger.InfoContext(ctx, "Webhook updated successfully", map[string]interface{}{
		"webhook_id": webhook.ID,
		"user_id":    userID,
	})
	
	return webhook, nil
}

func (s *webhookService) DeleteWebhook(ctx context.Context, webhookID uint64, userID uint64) error {
	webhook, err := s.GetWebhook(ctx, webhookID, userID)
	if err != nil {
		return err
	}
	
	if err := s.webhookRepo.Delete(ctx, webhookID); err != nil {
		return fmt.Errorf("failed to delete webhook: %w", err)
	}
	
	s.logger.InfoContext(ctx, "Webhook deleted successfully", map[string]interface{}{
		"webhook_id": webhook.ID,
		"user_id":    userID,
	})
	
	return nil
}

func (s *webhookService) ActivateWebhook(ctx context.Context, webhookID uint64, userID uint64) error {
	_, err := s.GetWebhook(ctx, webhookID, userID)
	if err != nil {
		return err
	}
	
	if err := s.webhookRepo.UpdateStatus(ctx, webhookID, domain.WebhookStatusActive); err != nil {
		return fmt.Errorf("failed to activate webhook: %w", err)
	}
	
	s.logger.InfoContext(ctx, "Webhook activated", map[string]interface{}{
		"webhook_id": webhookID,
		"user_id":    userID,
	})
	
	return nil
}

func (s *webhookService) DeactivateWebhook(ctx context.Context, webhookID uint64, userID uint64) error {
	_, err := s.GetWebhook(ctx, webhookID, userID)
	if err != nil {
		return err
	}
	
	if err := s.webhookRepo.UpdateStatus(ctx, webhookID, domain.WebhookStatusInactive); err != nil {
		return fmt.Errorf("failed to deactivate webhook: %w", err)
	}
	
	s.logger.InfoContext(ctx, "Webhook deactivated", map[string]interface{}{
		"webhook_id": webhookID,
		"user_id":    userID,
	})
	
	return nil
}

func (s *webhookService) TestWebhook(ctx context.Context, webhookID uint64, userID uint64) (*domain.WebhookDelivery, error) {
	webhook, err := s.GetWebhook(ctx, webhookID, userID)
	if err != nil {
		return nil, err
	}
	
	// Create test payload
	payload := &domain.WebhookPayload{
		ID:        fmt.Sprintf("test_%d_%d", webhookID, time.Now().Unix()),
		Event:     "system.test",
		Data: map[string]interface{}{
			"test":    true,
			"message": "This is a test webhook delivery",
		},
		Timestamp: time.Now(),
		UserID:    userID,
		Version:   "1.0",
	}
	
	delivery, err := s.deliveryService.DeliverWebhook(ctx, webhook, payload)
	if err != nil {
		s.logger.ErrorContext(ctx, "Webhook test failed", map[string]interface{}{
			"webhook_id": webhookID,
			"user_id":    userID,
			"error":      err.Error(),
		})
		return nil, fmt.Errorf("webhook test failed: %w", err)
	}
	
	s.logger.InfoContext(ctx, "Webhook test completed", map[string]interface{}{
		"webhook_id":  webhookID,
		"user_id":     userID,
		"delivery_id": delivery.ID,
		"status":      delivery.Status,
	})
	
	return delivery, nil
}

func (s *webhookService) TriggerEvent(ctx context.Context, event domain.WebhookEvent, data interface{}, userID uint64) error {
	return s.TriggerEventForUser(ctx, event, data, userID)
}

func (s *webhookService) TriggerEventForUser(ctx context.Context, event domain.WebhookEvent, data interface{}, userID uint64) error {
	webhooks, err := s.webhookRepo.GetByUserIDAndEvent(ctx, userID, event)
	if err != nil {
		return fmt.Errorf("failed to get webhooks for event: %w", err)
	}
	
	if len(webhooks) == 0 {
		return nil // No webhooks configured for this event
	}
	
	// Create payload
	payload := &domain.WebhookPayload{
		ID:        fmt.Sprintf("%s_%d_%d", event, userID, time.Now().Unix()),
		Event:     event,
		Data:      s.convertDataToMap(data),
		Timestamp: time.Now(),
		UserID:    userID,
		Version:   "1.0",
	}
	
	// Trigger webhooks
	for _, webhook := range webhooks {
		go func(w *domain.Webhook) {
			if _, err := s.deliveryService.DeliverWebhook(context.Background(), w, payload); err != nil {
				s.logger.ErrorContext(ctx, "Failed to deliver webhook", map[string]interface{}{
					"webhook_id": w.ID,
					"event":      event,
					"user_id":    userID,
					"error":      err.Error(),
				})
			}
		}(webhook)
	}
	
	s.logger.InfoContext(ctx, "Event triggered for user webhooks", map[string]interface{}{
		"event":         event,
		"user_id":       userID,
		"webhook_count": len(webhooks),
	})
	
	return nil
}

func (s *webhookService) TriggerEventGlobally(ctx context.Context, event domain.WebhookEvent, data interface{}) error {
	webhooks, err := s.webhookRepo.GetActiveWebhooks(ctx, event)
	if err != nil {
		return fmt.Errorf("failed to get active webhooks: %w", err)
	}
	
	if len(webhooks) == 0 {
		return nil // No webhooks configured for this event
	}
	
	// Create payload
	payload := &domain.WebhookPayload{
		ID:        fmt.Sprintf("%s_global_%d", event, time.Now().Unix()),
		Event:     event,
		Data:      s.convertDataToMap(data),
		Timestamp: time.Now(),
		Version:   "1.0",
	}
	
	// Trigger webhooks
	for _, webhook := range webhooks {
		payload.UserID = webhook.UserID
		
		go func(w *domain.Webhook, p *domain.WebhookPayload) {
			if _, err := s.deliveryService.DeliverWebhook(context.Background(), w, p); err != nil {
				s.logger.ErrorContext(ctx, "Failed to deliver global webhook", map[string]interface{}{
					"webhook_id": w.ID,
					"event":      event,
					"user_id":    w.UserID,
					"error":      err.Error(),
				})
			}
		}(webhook, payload)
	}
	
	s.logger.InfoContext(ctx, "Global event triggered", map[string]interface{}{
		"event":         event,
		"webhook_count": len(webhooks),
	})
	
	return nil
}

func (s *webhookService) GetWebhookDeliveries(ctx context.Context, webhookID uint64, userID uint64, offset, limit int) ([]*domain.WebhookDelivery, int64, error) {
	// Verify webhook ownership
	_, err := s.GetWebhook(ctx, webhookID, userID)
	if err != nil {
		return nil, 0, err
	}
	
	return s.deliveryRepo.GetByWebhookID(ctx, webhookID, offset, limit)
}

func (s *webhookService) GetDelivery(ctx context.Context, deliveryID uint64, userID uint64) (*domain.WebhookDelivery, error) {
	delivery, err := s.deliveryRepo.GetByID(ctx, deliveryID)
	if err != nil {
		return nil, err
	}
	
	// Verify webhook ownership through the delivery
	webhook, err := s.webhookRepo.GetByID(ctx, delivery.WebhookID)
	if err != nil {
		return nil, err
	}
	
	if webhook.UserID != userID {
		return nil, domain.NewNotFoundError("webhook_delivery")
	}
	
	return delivery, nil
}

func (s *webhookService) RetryDelivery(ctx context.Context, deliveryID uint64, userID uint64) error {
	delivery, err := s.GetDelivery(ctx, deliveryID, userID)
	if err != nil {
		return err
	}
	
	if !delivery.CanRetry() {
		return &domain.ValidationError{
			Field:   "delivery",
			Message: "Delivery cannot be retried",
		}
	}
	
	_, err = s.deliveryService.RetryFailedDelivery(ctx, delivery)
	if err != nil {
		return fmt.Errorf("failed to retry delivery: %w", err)
	}
	
	s.logger.InfoContext(ctx, "Webhook delivery retry requested", map[string]interface{}{
		"delivery_id": deliveryID,
		"webhook_id":  delivery.WebhookID,
		"user_id":     userID,
	})
	
	return nil
}

func (s *webhookService) GetWebhookStats(ctx context.Context, webhookID uint64, userID uint64) (*domain.WebhookDeliveryStats, error) {
	// Verify webhook ownership
	_, err := s.GetWebhook(ctx, webhookID, userID)
	if err != nil {
		return nil, err
	}
	
	return s.deliveryRepo.GetDeliveryStats(ctx, webhookID)
}

func (s *webhookService) GetFailedDeliveries(ctx context.Context, webhookID uint64, userID uint64, limit int) ([]*domain.WebhookDelivery, error) {
	// Verify webhook ownership
	_, err := s.GetWebhook(ctx, webhookID, userID)
	if err != nil {
		return nil, err
	}
	
	return s.deliveryRepo.GetFailedDeliveries(ctx, webhookID, limit)
}

func (s *webhookService) ProcessPendingRetries(ctx context.Context) error {
	deliveries, err := s.deliveryRepo.GetPendingRetries(ctx, 100) // Process up to 100 at a time
	if err != nil {
		return fmt.Errorf("failed to get pending retries: %w", err)
	}
	
	processedCount := 0
	for _, delivery := range deliveries {
		if _, err := s.deliveryService.RetryFailedDelivery(ctx, delivery); err != nil {
			s.logger.ErrorContext(ctx, "Failed to retry delivery", map[string]interface{}{
				"delivery_id": delivery.ID,
				"webhook_id":  delivery.WebhookID,
				"error":       err.Error(),
			})
			continue
		}
		processedCount++
	}
	
	if processedCount > 0 {
		s.logger.InfoContext(ctx, "Processed pending webhook retries", map[string]interface{}{
			"processed_count": processedCount,
			"total_pending":   len(deliveries),
		})
	}
	
	return nil
}

func (s *webhookService) CleanupOldDeliveries(ctx context.Context, olderThanDays int) (int64, error) {
	cutoffDate := time.Now().AddDate(0, 0, -olderThanDays)
	
	deletedCount, err := s.deliveryRepo.DeleteOldDeliveries(ctx, cutoffDate)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup old deliveries: %w", err)
	}
	
	if deletedCount > 0 {
		s.logger.InfoContext(ctx, "Cleaned up old webhook deliveries", map[string]interface{}{
			"deleted_count":     deletedCount,
			"older_than_days":   olderThanDays,
			"cutoff_date":       cutoffDate,
		})
	}
	
	return deletedCount, nil
}

// Helper methods

func (s *webhookService) generateSecret() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *webhookService) convertDataToMap(data interface{}) map[string]interface{} {
	if data == nil {
		return map[string]interface{}{}
	}
	
	// Try to convert to map directly
	if m, ok := data.(map[string]interface{}); ok {
		return m
	}
	
	// Convert to JSON and back to map
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return map[string]interface{}{
			"error": "Failed to serialize data",
			"raw":   fmt.Sprintf("%+v", data),
		}
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return map[string]interface{}{
			"error": "Failed to deserialize data",
			"raw":   string(jsonBytes),
		}
	}
	
	return result
}