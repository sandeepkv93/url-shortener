package services

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"url-shortener/internal/core/domain"
	"url-shortener/internal/core/ports"
)

type webhookDeliveryService struct {
	deliveryRepo ports.WebhookDeliveryRepository
	webhookRepo  ports.WebhookRepository
	httpClient   *http.Client
	logger       ports.Logger
}

func NewWebhookDeliveryService(
	deliveryRepo ports.WebhookDeliveryRepository,
	webhookRepo ports.WebhookRepository,
	logger ports.Logger,
) ports.WebhookDeliveryService {
	return &webhookDeliveryService{
		deliveryRepo: deliveryRepo,
		webhookRepo:  webhookRepo,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:       10,
				IdleConnTimeout:    30 * time.Second,
				DisableCompression: false,
			},
		},
		logger: logger,
	}
}

func (s *webhookDeliveryService) DeliverWebhook(ctx context.Context, webhook *domain.Webhook, payload *domain.WebhookPayload) (*domain.WebhookDelivery, error) {
	startTime := time.Now()
	
	// Create delivery record
	delivery := &domain.WebhookDelivery{
		WebhookID:   webhook.ID,
		EventType:   payload.Event,
		Status:      domain.WebhookDeliveryStatusPending,
		RequestURL:  webhook.URL,
		AttemptCount: 1,
	}
	
	// Serialize payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		delivery.Status = domain.WebhookDeliveryStatusFailed
		delivery.ErrorMessage = fmt.Sprintf("Failed to serialize payload: %v", err)
		s.deliveryRepo.Create(ctx, delivery)
		return delivery, fmt.Errorf("failed to serialize payload: %w", err)
	}
	
	delivery.RequestBody = json.RawMessage(payloadBytes)
	
	// Prepare HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", webhook.URL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		delivery.Status = domain.WebhookDeliveryStatusFailed
		delivery.ErrorMessage = fmt.Sprintf("Failed to create HTTP request: %v", err)
		s.deliveryRepo.Create(ctx, delivery)
		return delivery, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "URLShortener-Webhook/1.0")
	req.Header.Set("X-Webhook-Event", string(payload.Event))
	req.Header.Set("X-Webhook-ID", payload.ID)
	req.Header.Set("X-Webhook-Timestamp", strconv.FormatInt(payload.Timestamp.Unix(), 10))
	
	// Add signature if secret is configured
	if webhook.Secret != "" {
		signature := s.GenerateSignature(payloadBytes, webhook.Secret)
		req.Header.Set("X-Webhook-Signature", signature)
	}
	
	// Store request headers
	requestHeaders, _ := json.Marshal(req.Header)
	delivery.RequestHeaders = json.RawMessage(requestHeaders)
	
	// Set custom timeout if configured
	client := s.httpClient
	if webhook.TimeoutSeconds > 0 {
		client = &http.Client{
			Timeout: time.Duration(webhook.TimeoutSeconds) * time.Second,
			Transport: s.httpClient.Transport,
		}
	}
	
	// Make the request
	resp, err := client.Do(req)
	duration := time.Since(startTime).Milliseconds()
	delivery.Duration = duration
	
	if err != nil {
		delivery.Status = domain.WebhookDeliveryStatusFailed
		delivery.ErrorMessage = fmt.Sprintf("HTTP request failed: %v", err)
		s.deliveryRepo.Create(ctx, delivery)
		
		// Schedule retry if within retry limits
		if delivery.AttemptCount < webhook.MaxRetries {
			s.ScheduleRetry(ctx, delivery)
		}
		
		return delivery, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()
	
	// Read response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		delivery.ResponseStatus = resp.StatusCode
		delivery.ErrorMessage = fmt.Sprintf("Failed to read response body: %v", err)
		delivery.Status = domain.WebhookDeliveryStatusFailed
	} else {
		delivery.ResponseStatus = resp.StatusCode
		delivery.ResponseBody = string(responseBody)
		
		// Store response headers
		responseHeaders, _ := json.Marshal(resp.Header)
		delivery.ResponseHeaders = json.RawMessage(responseHeaders)
		
		// Determine success based on status code
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			delivery.Status = domain.WebhookDeliveryStatusSuccess
			
			// Update webhook statistics
			s.updateWebhookStatistics(ctx, webhook.ID, true)
		} else {
			delivery.Status = domain.WebhookDeliveryStatusFailed
			delivery.ErrorMessage = fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(responseBody))
			
			// Update webhook statistics
			s.updateWebhookStatistics(ctx, webhook.ID, false)
			
			// Schedule retry if within retry limits and status suggests retry could work
			if delivery.AttemptCount < webhook.MaxRetries && s.shouldRetry(resp.StatusCode) {
				s.ScheduleRetry(ctx, delivery)
			}
		}
	}
	
	// Save delivery record
	if err := s.deliveryRepo.Create(ctx, delivery); err != nil {
		s.logger.ErrorContext(ctx, "Failed to save webhook delivery", map[string]interface{}{
			"webhook_id": webhook.ID,
			"error":      err.Error(),
		})
	}
	
	// Update webhook last delivery timestamp
	s.webhookRepo.UpdateLastDelivery(ctx, webhook.ID, delivery.Status == domain.WebhookDeliveryStatusSuccess)
	
	s.logger.InfoContext(ctx, "Webhook delivery completed", map[string]interface{}{
		"webhook_id":      webhook.ID,
		"delivery_id":     delivery.ID,
		"status":          delivery.Status,
		"response_status": delivery.ResponseStatus,
		"duration_ms":     delivery.Duration,
		"attempt":         delivery.AttemptCount,
	})
	
	return delivery, nil
}

func (s *webhookDeliveryService) RetryFailedDelivery(ctx context.Context, delivery *domain.WebhookDelivery) (*domain.WebhookDelivery, error) {
	if !delivery.CanRetry() {
		return nil, fmt.Errorf("delivery cannot be retried")
	}
	
	// Get the webhook
	webhook, err := s.webhookRepo.GetByID(ctx, delivery.WebhookID)
	if err != nil {
		return nil, fmt.Errorf("failed to get webhook: %w", err)
	}
	
	// Check if delivery should be abandoned
	if delivery.ShouldAbandon(webhook.MaxRetries) {
		delivery.Status = domain.WebhookDeliveryStatusAbandoned
		s.deliveryRepo.Update(ctx, delivery)
		return delivery, nil
	}
	
	// Increment attempt count
	s.deliveryRepo.IncrementAttempt(ctx, delivery.ID)
	delivery.AttemptCount++
	
	// Reconstruct payload from stored request body
	var payload domain.WebhookPayload
	if err := json.Unmarshal(delivery.RequestBody, &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal stored payload: %w", err)
	}
	
	// Reattempt delivery
	newDelivery, err := s.DeliverWebhook(ctx, webhook, &payload)
	if err != nil {
		return delivery, err
	}
	
	return newDelivery, nil
}

func (s *webhookDeliveryService) ScheduleRetry(ctx context.Context, delivery *domain.WebhookDelivery) error {
	// Calculate exponential backoff
	backoffMs := int64(1000) // Start with 1 second
	for i := 1; i < delivery.AttemptCount; i++ {
		backoffMs *= 2 // Double the backoff for each retry
	}
	
	// Cap at 5 minutes
	if backoffMs > 300000 {
		backoffMs = 300000
	}
	
	nextRetryAt := time.Now().Add(time.Duration(backoffMs) * time.Millisecond)
	
	return s.deliveryRepo.ScheduleRetry(ctx, delivery.ID, &nextRetryAt)
}

func (s *webhookDeliveryService) ProcessPendingDeliveries(ctx context.Context, limit int) (int, error) {
	deliveries, err := s.deliveryRepo.GetPendingRetries(ctx, limit)
	if err != nil {
		return 0, fmt.Errorf("failed to get pending deliveries: %w", err)
	}
	
	processedCount := 0
	for _, delivery := range deliveries {
		if _, err := s.RetryFailedDelivery(ctx, delivery); err != nil {
			s.logger.ErrorContext(ctx, "Failed to process pending delivery", map[string]interface{}{
				"delivery_id": delivery.ID,
				"webhook_id":  delivery.WebhookID,
				"error":       err.Error(),
			})
			continue
		}
		processedCount++
	}
	
	return processedCount, nil
}

func (s *webhookDeliveryService) ProcessFailedDeliveries(ctx context.Context, limit int) (int, error) {
	// This could be expanded to handle specific retry logic for failed deliveries
	return s.ProcessPendingDeliveries(ctx, limit)
}

func (s *webhookDeliveryService) ValidateWebhookURL(ctx context.Context, webhookURL string) error {
	// Parse URL
	u, err := url.Parse(webhookURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}
	
	// Check scheme
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("invalid URL scheme: %s (must be http or https)", u.Scheme)
	}
	
	// Check host
	if u.Host == "" {
		return fmt.Errorf("URL must have a host")
	}
	
	// Block localhost and private IPs in production
	host := strings.Split(u.Host, ":")[0]
	if s.isPrivateOrLocalhost(host) {
		return fmt.Errorf("webhook URLs cannot point to localhost or private IP addresses")
	}
	
	// Optionally, make a test request to validate the endpoint
	// This is commented out as it might be too aggressive for validation
	/*
	testReq, _ := http.NewRequestWithContext(ctx, "HEAD", webhookURL, nil)
	resp, err := s.httpClient.Do(testReq)
	if err != nil {
		return fmt.Errorf("webhook URL is not reachable: %w", err)
	}
	resp.Body.Close()
	
	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook URL returned error status: %d", resp.StatusCode)
	}
	*/
	
	return nil
}

func (s *webhookDeliveryService) GenerateSignature(payload []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(payload)
	return "sha256=" + hex.EncodeToString(h.Sum(nil))
}

func (s *webhookDeliveryService) VerifySignature(payload []byte, signature, secret string) bool {
	expectedSignature := s.GenerateSignature(payload, secret)
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// Helper methods

func (s *webhookDeliveryService) shouldRetry(statusCode int) bool {
	// Retry on 5xx errors and specific 4xx errors
	switch statusCode {
	case 408, 409, 429: // Timeout, Conflict, Too Many Requests
		return true
	case 500, 502, 503, 504: // Server errors
		return true
	default:
		return false
	}
}

func (s *webhookDeliveryService) isPrivateOrLocalhost(host string) bool {
	// Simple check for localhost and common private IP ranges
	privateHosts := []string{
		"localhost",
		"127.0.0.1",
		"::1",
	}
	
	for _, private := range privateHosts {
		if strings.EqualFold(host, private) {
			return true
		}
	}
	
	// Check for private IP ranges (simplified)
	if strings.HasPrefix(host, "10.") ||
		strings.HasPrefix(host, "192.168.") ||
		strings.HasPrefix(host, "172.") {
		return true
	}
	
	return false
}

func (s *webhookDeliveryService) updateWebhookStatistics(ctx context.Context, webhookID uint64, success bool) {
	// This could be optimized with atomic updates or caching
	webhook, err := s.webhookRepo.GetByID(ctx, webhookID)
	if err != nil {
		return
	}
	
	webhook.TotalDeliveries++
	if success {
		webhook.SuccessDeliveries++
	} else {
		webhook.FailedDeliveries++
	}
	
	s.webhookRepo.UpdateStatistics(ctx, webhookID, webhook.TotalDeliveries, webhook.SuccessDeliveries, webhook.FailedDeliveries)
}