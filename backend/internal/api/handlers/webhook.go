package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"url-shortener/internal/core/domain"
	"url-shortener/internal/core/ports"
)

type WebhookHandler struct {
	webhookService ports.WebhookService
	logger         ports.Logger
}

func NewWebhookHandler(webhookService ports.WebhookService, logger ports.Logger) *WebhookHandler {
	return &WebhookHandler{
		webhookService: webhookService,
		logger:         logger,
	}
}

// CreateWebhook handles POST /webhooks
func (h *WebhookHandler) CreateWebhook(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	
	var req domain.WebhookCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, &domain.ValidationError{
			Field:   "body",
			Message: "Invalid request body",
		}, http.StatusBadRequest)
		return
	}
	
	webhook, err := h.webhookService.CreateWebhook(r.Context(), uint64(userID), req)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.sendJSON(w, webhook, http.StatusCreated)
}

// GetUserWebhooks handles GET /webhooks
func (h *WebhookHandler) GetUserWebhooks(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	
	offset, limit := getPaginationParams(r)
	
	webhooks, total, err := h.webhookService.GetUserWebhooks(r.Context(), uint64(userID), offset, limit)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	response := map[string]interface{}{
		"webhooks": webhooks,
		"total":    total,
		"offset":   offset,
		"limit":    limit,
	}
	
	h.sendJSON(w, response, http.StatusOK)
}

// GetWebhook handles GET /webhooks/{id}
func (h *WebhookHandler) GetWebhook(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	webhookID := h.getWebhookIDFromURL(w, r)
	if webhookID == 0 {
		return
	}
	
	webhook, err := h.webhookService.GetWebhook(r.Context(), webhookID, uint64(userID))
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.sendJSON(w, webhook, http.StatusOK)
}

// UpdateWebhook handles PUT /webhooks/{id}
func (h *WebhookHandler) UpdateWebhook(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	webhookID := h.getWebhookIDFromURL(w, r)
	if webhookID == 0 {
		return
	}
	
	var req domain.WebhookUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, &domain.ValidationError{
			Field:   "body",
			Message: "Invalid request body",
		}, http.StatusBadRequest)
		return
	}
	
	webhook, err := h.webhookService.UpdateWebhook(r.Context(), webhookID, uint64(userID), req)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.sendJSON(w, webhook, http.StatusOK)
}

// DeleteWebhook handles DELETE /webhooks/{id}
func (h *WebhookHandler) DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	webhookID := h.getWebhookIDFromURL(w, r)
	if webhookID == 0 {
		return
	}
	
	err := h.webhookService.DeleteWebhook(r.Context(), webhookID, uint64(userID))
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.sendJSON(w, map[string]string{"message": "Webhook deleted successfully"}, http.StatusOK)
}

// ActivateWebhook handles POST /webhooks/{id}/activate
func (h *WebhookHandler) ActivateWebhook(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	webhookID := h.getWebhookIDFromURL(w, r)
	if webhookID == 0 {
		return
	}
	
	err := h.webhookService.ActivateWebhook(r.Context(), webhookID, uint64(userID))
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.sendJSON(w, map[string]string{"message": "Webhook activated successfully"}, http.StatusOK)
}

// DeactivateWebhook handles POST /webhooks/{id}/deactivate
func (h *WebhookHandler) DeactivateWebhook(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	webhookID := h.getWebhookIDFromURL(w, r)
	if webhookID == 0 {
		return
	}
	
	err := h.webhookService.DeactivateWebhook(r.Context(), webhookID, uint64(userID))
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.sendJSON(w, map[string]string{"message": "Webhook deactivated successfully"}, http.StatusOK)
}

// TestWebhook handles POST /webhooks/{id}/test
func (h *WebhookHandler) TestWebhook(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	webhookID := h.getWebhookIDFromURL(w, r)
	if webhookID == 0 {
		return
	}
	
	delivery, err := h.webhookService.TestWebhook(r.Context(), webhookID, uint64(userID))
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.sendJSON(w, delivery, http.StatusOK)
}

// GetWebhookDeliveries handles GET /webhooks/{id}/deliveries
func (h *WebhookHandler) GetWebhookDeliveries(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	webhookID := h.getWebhookIDFromURL(w, r)
	if webhookID == 0 {
		return
	}
	
	offset, limit := getPaginationParams(r)
	
	deliveries, total, err := h.webhookService.GetWebhookDeliveries(r.Context(), webhookID, uint64(userID), offset, limit)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	response := map[string]interface{}{
		"deliveries": deliveries,
		"total":      total,
		"offset":     offset,
		"limit":      limit,
	}
	
	h.sendJSON(w, response, http.StatusOK)
}

// GetDelivery handles GET /webhooks/deliveries/{id}
func (h *WebhookHandler) GetDelivery(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	deliveryID := h.getDeliveryIDFromURL(w, r)
	if deliveryID == 0 {
		return
	}
	
	delivery, err := h.webhookService.GetDelivery(r.Context(), deliveryID, uint64(userID))
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.sendJSON(w, delivery, http.StatusOK)
}

// RetryDelivery handles POST /webhooks/deliveries/{id}/retry
func (h *WebhookHandler) RetryDelivery(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	deliveryID := h.getDeliveryIDFromURL(w, r)
	if deliveryID == 0 {
		return
	}
	
	err := h.webhookService.RetryDelivery(r.Context(), deliveryID, uint64(userID))
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.sendJSON(w, map[string]string{"message": "Delivery retry scheduled successfully"}, http.StatusOK)
}

// GetWebhookStats handles GET /webhooks/{id}/stats
func (h *WebhookHandler) GetWebhookStats(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	webhookID := h.getWebhookIDFromURL(w, r)
	if webhookID == 0 {
		return
	}
	
	stats, err := h.webhookService.GetWebhookStats(r.Context(), webhookID, uint64(userID))
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.sendJSON(w, stats, http.StatusOK)
}

// GetFailedDeliveries handles GET /webhooks/{id}/failed-deliveries
func (h *WebhookHandler) GetFailedDeliveries(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	webhookID := h.getWebhookIDFromURL(w, r)
	if webhookID == 0 {
		return
	}
	
	limit := 50 // Default limit for failed deliveries
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}
	
	deliveries, err := h.webhookService.GetFailedDeliveries(r.Context(), webhookID, uint64(userID), limit)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	
	h.sendJSON(w, map[string]interface{}{
		"failed_deliveries": deliveries,
		"limit":             limit,
	}, http.StatusOK)
}

// GetWebhookEvents handles GET /webhooks/events
func (h *WebhookHandler) GetWebhookEvents(w http.ResponseWriter, r *http.Request) {
	events := map[string]interface{}{
		"url_events": []string{
			string(domain.WebhookEventURLCreated),
			string(domain.WebhookEventURLUpdated),
			string(domain.WebhookEventURLDeleted),
			string(domain.WebhookEventURLClicked),
			string(domain.WebhookEventURLExpired),
		},
		"analytics_events": []string{
			string(domain.WebhookEventAnalyticsThreshold),
			string(domain.WebhookEventAnalyticsReport),
		},
		"user_events": []string{
			string(domain.WebhookEventUserRegistered),
			string(domain.WebhookEventUserUpdated),
		},
		"system_events": []string{
			string(domain.WebhookEventSystemError),
			string(domain.WebhookEventSystemAlert),
		},
	}
	
	h.sendJSON(w, events, http.StatusOK)
}

// Helper methods

func (h *WebhookHandler) getWebhookIDFromURL(w http.ResponseWriter, r *http.Request) uint64 {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		h.sendError(w, &domain.ValidationError{
			Field:   "id",
			Message: "Webhook ID is required",
		}, http.StatusBadRequest)
		return 0
	}
	
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.sendError(w, &domain.ValidationError{
			Field:   "id",
			Message: "Invalid webhook ID",
		}, http.StatusBadRequest)
		return 0
	}
	
	return id
}

func (h *WebhookHandler) getDeliveryIDFromURL(w http.ResponseWriter, r *http.Request) uint64 {
	idStr := chi.URLParam(r, "delivery_id")
	if idStr == "" {
		idStr = chi.URLParam(r, "id")
	}
	
	if idStr == "" {
		h.sendError(w, &domain.ValidationError{
			Field:   "id",
			Message: "Delivery ID is required",
		}, http.StatusBadRequest)
		return 0
	}
	
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.sendError(w, &domain.ValidationError{
			Field:   "id",
			Message: "Invalid delivery ID",
		}, http.StatusBadRequest)
		return 0
	}
	
	return id
}

func (h *WebhookHandler) handleServiceError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case *domain.ValidationError:
		h.sendError(w, e, http.StatusBadRequest)
	case *domain.NotFoundError:
		h.sendError(w, e, http.StatusNotFound)
	case *domain.ConflictError:
		h.sendError(w, e, http.StatusConflict)
	default:
		h.logger.Error("Webhook service error", map[string]interface{}{
			"error": err.Error(),
		})
		h.sendError(w, &domain.InternalError{
			Message: "Internal server error",
		}, http.StatusInternalServerError)
	}
}

func (h *WebhookHandler) sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode JSON response", map[string]interface{}{
			"error": err.Error(),
		})
	}
}

func (h *WebhookHandler) sendError(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	errorResponse := map[string]interface{}{
		"error": err.Error(),
	}
	
	if validationErr, ok := err.(*domain.ValidationError); ok {
		errorResponse["field"] = validationErr.Field
	}
	
	json.NewEncoder(w).Encode(errorResponse)
}