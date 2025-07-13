package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"url-shortener/internal/api/handlers"
	"url-shortener/internal/core/domain"
	"url-shortener/tests/mocks"
)

type WebhookHandlerTestSuite struct {
	suite.Suite
	webhookService *mocks.MockWebhookService
	logger         *mocks.MockLogger
	handler        *handlers.WebhookHandler
}

func TestWebhookHandlerSuite(t *testing.T) {
	suite.Run(t, new(WebhookHandlerTestSuite))
}

func (suite *WebhookHandlerTestSuite) SetupTest() {
	suite.webhookService = &mocks.MockWebhookService{}
	suite.logger = &mocks.MockLogger{}
	suite.handler = handlers.NewWebhookHandler(suite.webhookService, suite.logger)
}

func (suite *WebhookHandlerTestSuite) TestCreateWebhook_Success() {
	req := domain.WebhookCreateRequest{
		Name:   "Test Webhook",
		URL:    "https://example.com/webhook",
		Events: []domain.WebhookEvent{domain.WebhookEventURLCreated},
		Secret: "test-secret",
	}
	
	expectedWebhook := &domain.Webhook{
		ID:     123,
		UserID: 456,
		Name:   req.Name,
		URL:    req.URL,
		Events: req.Events,
		Status: domain.WebhookStatusActive,
	}
	
	suite.webhookService.On("CreateWebhook", mock.Anything, uint64(456), req).Return(expectedWebhook, nil)
	
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/webhooks", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	
	// Add user ID to context
	ctx := context.WithValue(httpReq.Context(), "user_id", uint(456))
	httpReq = httpReq.WithContext(ctx)
	
	w := httptest.NewRecorder()
	suite.handler.CreateWebhook(w, httpReq)
	
	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	
	var response domain.Webhook
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedWebhook.ID, response.ID)
	assert.Equal(suite.T(), expectedWebhook.Name, response.Name)
	
	suite.webhookService.AssertExpectations(suite.T())
}

func (suite *WebhookHandlerTestSuite) TestCreateWebhook_InvalidJSON() {
	httpReq := httptest.NewRequest("POST", "/webhooks", bytes.NewBufferString("invalid json"))
	httpReq.Header.Set("Content-Type", "application/json")
	
	// Add user ID to context
	ctx := context.WithValue(httpReq.Context(), "user_id", uint(456))
	httpReq = httpReq.WithContext(ctx)
	
	w := httptest.NewRecorder()
	suite.handler.CreateWebhook(w, httpReq)
	
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Contains(suite.T(), response["error"], "Invalid request body")
}

func (suite *WebhookHandlerTestSuite) TestCreateWebhook_ServiceError() {
	req := domain.WebhookCreateRequest{
		Name:   "Test Webhook",
		URL:    "invalid-url",
		Events: []domain.WebhookEvent{domain.WebhookEventURLCreated},
	}
	
	validationErr := &domain.ValidationError{
		Field:   "url",
		Message: "Invalid URL format",
	}
	
	suite.webhookService.On("CreateWebhook", mock.Anything, uint64(456), req).Return(nil, validationErr)
	
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/webhooks", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	
	// Add user ID to context
	ctx := context.WithValue(httpReq.Context(), "user_id", uint(456))
	httpReq = httpReq.WithContext(ctx)
	
	w := httptest.NewRecorder()
	suite.handler.CreateWebhook(w, httpReq)
	
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), validationErr.Error(), response["error"])
	assert.Equal(suite.T(), validationErr.Field, response["field"])
	
	suite.webhookService.AssertExpectations(suite.T())
}

func (suite *WebhookHandlerTestSuite) TestGetUserWebhooks_Success() {
	expectedWebhooks := []*domain.Webhook{
		{
			ID:     123,
			UserID: 456,
			Name:   "Webhook 1",
			URL:    "https://example.com/webhook1",
			Status: domain.WebhookStatusActive,
		},
		{
			ID:     124,
			UserID: 456,
			Name:   "Webhook 2",
			URL:    "https://example.com/webhook2",
			Status: domain.WebhookStatusActive,
		},
	}
	
	suite.webhookService.On("GetUserWebhooks", mock.Anything, uint64(456), 0, 10).Return(expectedWebhooks, int64(2), nil)
	
	httpReq := httptest.NewRequest("GET", "/webhooks?limit=10&offset=0", nil)
	
	// Add user ID to context
	ctx := context.WithValue(httpReq.Context(), "user_id", uint(456))
	httpReq = httpReq.WithContext(ctx)
	
	w := httptest.NewRecorder()
	suite.handler.GetUserWebhooks(w, httpReq)
	
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	webhooks := response["webhooks"].([]interface{})
	assert.Equal(suite.T(), 2, len(webhooks))
	assert.Equal(suite.T(), float64(2), response["total"])
	
	suite.webhookService.AssertExpectations(suite.T())
}

func (suite *WebhookHandlerTestSuite) TestGetWebhook_Success() {
	expectedWebhook := &domain.Webhook{
		ID:     123,
		UserID: 456,
		Name:   "Test Webhook",
		URL:    "https://example.com/webhook",
		Status: domain.WebhookStatusActive,
	}
	
	suite.webhookService.On("GetWebhook", mock.Anything, uint64(123), uint64(456)).Return(expectedWebhook, nil)
	
	httpReq := httptest.NewRequest("GET", "/webhooks/123", nil)
	
	// Add user ID to context
	ctx := context.WithValue(httpReq.Context(), "user_id", uint(456))
	httpReq = httpReq.WithContext(ctx)
	
	// Add URL parameter to context using chi
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "123")
	httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), chi.RouteCtxKey, rctx))
	
	w := httptest.NewRecorder()
	suite.handler.GetWebhook(w, httpReq)
	
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response domain.Webhook
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedWebhook.ID, response.ID)
	assert.Equal(suite.T(), expectedWebhook.Name, response.Name)
	
	suite.webhookService.AssertExpectations(suite.T())
}

func (suite *WebhookHandlerTestSuite) TestGetWebhook_InvalidID() {
	httpReq := httptest.NewRequest("GET", "/webhooks/invalid", nil)
	
	// Add user ID to context
	ctx := context.WithValue(httpReq.Context(), "user_id", uint(456))
	httpReq = httpReq.WithContext(ctx)
	
	// Add invalid URL parameter to context using chi
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "invalid")
	httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), chi.RouteCtxKey, rctx))
	
	w := httptest.NewRecorder()
	suite.handler.GetWebhook(w, httpReq)
	
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Contains(suite.T(), response["error"], "Invalid webhook ID")
}

func (suite *WebhookHandlerTestSuite) TestGetWebhook_NotFound() {
	notFoundErr := &domain.NotFoundError{
		Resource: "webhook",
		ID:       "123",
	}
	
	suite.webhookService.On("GetWebhook", mock.Anything, uint64(123), uint64(456)).Return(nil, notFoundErr)
	
	httpReq := httptest.NewRequest("GET", "/webhooks/123", nil)
	
	// Add user ID to context
	ctx := context.WithValue(httpReq.Context(), "user_id", uint(456))
	httpReq = httpReq.WithContext(ctx)
	
	// Add URL parameter to context using chi
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "123")
	httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), chi.RouteCtxKey, rctx))
	
	w := httptest.NewRecorder()
	suite.handler.GetWebhook(w, httpReq)
	
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
	
	suite.webhookService.AssertExpectations(suite.T())
}

func (suite *WebhookHandlerTestSuite) TestUpdateWebhook_Success() {
	newName := "Updated Webhook"
	req := domain.WebhookUpdateRequest{
		Name: &newName,
	}
	
	expectedWebhook := &domain.Webhook{
		ID:     123,
		UserID: 456,
		Name:   newName,
		URL:    "https://example.com/webhook",
		Status: domain.WebhookStatusActive,
	}
	
	suite.webhookService.On("UpdateWebhook", mock.Anything, uint64(123), uint64(456), req).Return(expectedWebhook, nil)
	
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("PUT", "/webhooks/123", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	
	// Add user ID to context
	ctx := context.WithValue(httpReq.Context(), "user_id", uint(456))
	httpReq = httpReq.WithContext(ctx)
	
	// Add URL parameter to context using chi
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "123")
	httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), chi.RouteCtxKey, rctx))
	
	w := httptest.NewRecorder()
	suite.handler.UpdateWebhook(w, httpReq)
	
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response domain.Webhook
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedWebhook.Name, response.Name)
	
	suite.webhookService.AssertExpectations(suite.T())
}

func (suite *WebhookHandlerTestSuite) TestDeleteWebhook_Success() {
	suite.webhookService.On("DeleteWebhook", mock.Anything, uint64(123), uint64(456)).Return(nil)
	
	httpReq := httptest.NewRequest("DELETE", "/webhooks/123", nil)
	
	// Add user ID to context
	ctx := context.WithValue(httpReq.Context(), "user_id", uint(456))
	httpReq = httpReq.WithContext(ctx)
	
	// Add URL parameter to context using chi
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "123")
	httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), chi.RouteCtxKey, rctx))
	
	w := httptest.NewRecorder()
	suite.handler.DeleteWebhook(w, httpReq)
	
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Webhook deleted successfully", response["message"])
	
	suite.webhookService.AssertExpectations(suite.T())
}

func (suite *WebhookHandlerTestSuite) TestActivateWebhook_Success() {
	suite.webhookService.On("ActivateWebhook", mock.Anything, uint64(123), uint64(456)).Return(nil)
	
	httpReq := httptest.NewRequest("POST", "/webhooks/123/activate", nil)
	
	// Add user ID to context
	ctx := context.WithValue(httpReq.Context(), "user_id", uint(456))
	httpReq = httpReq.WithContext(ctx)
	
	// Add URL parameter to context using chi
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "123")
	httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), chi.RouteCtxKey, rctx))
	
	w := httptest.NewRecorder()
	suite.handler.ActivateWebhook(w, httpReq)
	
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Webhook activated successfully", response["message"])
	
	suite.webhookService.AssertExpectations(suite.T())
}

func (suite *WebhookHandlerTestSuite) TestTestWebhook_Success() {
	expectedDelivery := &domain.WebhookDelivery{
		ID:            789,
		WebhookID:     123,
		Status:        domain.WebhookDeliveryStatusSuccess,
		ResponseStatus: http.StatusOK,
	}
	
	suite.webhookService.On("TestWebhook", mock.Anything, uint64(123), uint64(456)).Return(expectedDelivery, nil)
	
	httpReq := httptest.NewRequest("POST", "/webhooks/123/test", nil)
	
	// Add user ID to context
	ctx := context.WithValue(httpReq.Context(), "user_id", uint(456))
	httpReq = httpReq.WithContext(ctx)
	
	// Add URL parameter to context using chi
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "123")
	httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), chi.RouteCtxKey, rctx))
	
	w := httptest.NewRecorder()
	suite.handler.TestWebhook(w, httpReq)
	
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response domain.WebhookDelivery
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedDelivery.ID, response.ID)
	assert.Equal(suite.T(), expectedDelivery.Status, response.Status)
	
	suite.webhookService.AssertExpectations(suite.T())
}

func (suite *WebhookHandlerTestSuite) TestGetWebhookDeliveries_Success() {
	expectedDeliveries := []*domain.WebhookDelivery{
		{
			ID:            123,
			WebhookID:     456,
			Status:        domain.WebhookDeliveryStatusSuccess,
			ResponseStatus: http.StatusOK,
		},
		{
			ID:            124,
			WebhookID:     456,
			Status:        domain.WebhookDeliveryStatusFailed,
			ResponseStatus: http.StatusInternalServerError,
		},
	}
	
	suite.webhookService.On("GetWebhookDeliveries", mock.Anything, uint64(456), uint64(123), 0, 10).Return(expectedDeliveries, int64(2), nil)
	
	httpReq := httptest.NewRequest("GET", "/webhooks/456/deliveries?limit=10&offset=0", nil)
	
	// Add user ID to context
	ctx := context.WithValue(httpReq.Context(), "user_id", uint(123))
	httpReq = httpReq.WithContext(ctx)
	
	// Add URL parameter to context using chi
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "456")
	httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), chi.RouteCtxKey, rctx))
	
	w := httptest.NewRecorder()
	suite.handler.GetWebhookDeliveries(w, httpReq)
	
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	deliveries := response["deliveries"].([]interface{})
	assert.Equal(suite.T(), 2, len(deliveries))
	assert.Equal(suite.T(), float64(2), response["total"])
	
	suite.webhookService.AssertExpectations(suite.T())
}

func (suite *WebhookHandlerTestSuite) TestGetWebhookEvents_Success() {
	httpReq := httptest.NewRequest("GET", "/webhooks/events", nil)
	
	w := httptest.NewRecorder()
	suite.handler.GetWebhookEvents(w, httpReq)
	
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	urlEvents := response["url_events"].([]interface{})
	analyticsEvents := response["analytics_events"].([]interface{})
	userEvents := response["user_events"].([]interface{})
	systemEvents := response["system_events"].([]interface{})
	
	assert.NotEmpty(suite.T(), urlEvents)
	assert.NotEmpty(suite.T(), analyticsEvents)
	assert.NotEmpty(suite.T(), userEvents)
	assert.NotEmpty(suite.T(), systemEvents)
	
	// Verify specific events are included
	assert.Contains(suite.T(), urlEvents, string(domain.WebhookEventURLCreated))
	assert.Contains(suite.T(), analyticsEvents, string(domain.WebhookEventAnalyticsThreshold))
	assert.Contains(suite.T(), userEvents, string(domain.WebhookEventUserRegistered))
	assert.Contains(suite.T(), systemEvents, string(domain.WebhookEventSystemError))
}

func (suite *WebhookHandlerTestSuite) TestGetWebhookStats_Success() {
	expectedStats := &domain.WebhookDeliveryStats{
		TotalDeliveries:   100,
		SuccessDeliveries: 95,
		FailedDeliveries:  5,
		SuccessRate:       95.0,
		AverageResponseTime: 250,
	}
	
	suite.webhookService.On("GetWebhookStats", mock.Anything, uint64(123), uint64(456)).Return(expectedStats, nil)
	
	httpReq := httptest.NewRequest("GET", "/webhooks/123/stats", nil)
	
	// Add user ID to context
	ctx := context.WithValue(httpReq.Context(), "user_id", uint(456))
	httpReq = httpReq.WithContext(ctx)
	
	// Add URL parameter to context using chi
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "123")
	httpReq = httpReq.WithContext(context.WithValue(httpReq.Context(), chi.RouteCtxKey, rctx))
	
	w := httptest.NewRecorder()
	suite.handler.GetWebhookStats(w, httpReq)
	
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response domain.WebhookDeliveryStats
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedStats.TotalDeliveries, response.TotalDeliveries)
	assert.Equal(suite.T(), expectedStats.SuccessRate, response.SuccessRate)
	
	suite.webhookService.AssertExpectations(suite.T())
}

// Helper function to get user ID from context (this would normally be in a middleware)
func GetUserIDFromContext(ctx context.Context) uint {
	if userID, ok := ctx.Value("user_id").(uint); ok {
		return userID
	}
	return 0
}