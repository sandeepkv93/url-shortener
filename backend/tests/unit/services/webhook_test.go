package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"url-shortener/internal/core/domain"
	"url-shortener/internal/core/services"
	"url-shortener/tests/mocks"
)

type WebhookServiceTestSuite struct {
	suite.Suite
	webhookRepo     *mocks.MockWebhookRepository
	deliveryRepo    *mocks.MockWebhookDeliveryRepository
	deliveryService *mocks.MockWebhookDeliveryService
	logger          *mocks.MockLogger
	service         *services.WebhookService
}

func TestWebhookServiceSuite(t *testing.T) {
	suite.Run(t, new(WebhookServiceTestSuite))
}

func (suite *WebhookServiceTestSuite) SetupTest() {
	suite.webhookRepo = &mocks.MockWebhookRepository{}
	suite.deliveryRepo = &mocks.MockWebhookDeliveryRepository{}
	suite.deliveryService = &mocks.MockWebhookDeliveryService{}
	suite.logger = &mocks.MockLogger{}
	
	suite.service = services.NewWebhookService(
		suite.webhookRepo,
		suite.deliveryRepo,
		suite.deliveryService,
		suite.logger,
	)
}

func (suite *WebhookServiceTestSuite) TestCreateWebhook_Success() {
	ctx := context.Background()
	userID := uint64(123)
	
	req := domain.WebhookCreateRequest{
		Name:   "Test Webhook",
		URL:    "https://example.com/webhook",
		Events: []domain.WebhookEvent{domain.WebhookEventURLCreated},
		Secret: "test-secret",
	}
	
	suite.deliveryService.On("ValidateWebhookURL", ctx, req.URL).Return(nil)
	suite.webhookRepo.On("Create", ctx, mock.AnythingOfType("*domain.Webhook")).Return(nil)
	suite.logger.On("InfoContext", ctx, "Webhook created successfully", mock.Anything).Return()
	
	webhook, err := suite.service.CreateWebhook(ctx, userID, req)
	
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), webhook)
	assert.Equal(suite.T(), req.Name, webhook.Name)
	assert.Equal(suite.T(), req.URL, webhook.URL)
	assert.Equal(suite.T(), req.Events, webhook.Events)
	assert.Equal(suite.T(), userID, webhook.UserID)
	assert.Equal(suite.T(), domain.WebhookStatusActive, webhook.Status)
	
	suite.webhookRepo.AssertExpectations(suite.T())
	suite.deliveryService.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

func (suite *WebhookServiceTestSuite) TestCreateWebhook_InvalidEvents() {
	ctx := context.Background()
	userID := uint64(123)
	
	req := domain.WebhookCreateRequest{
		Name:   "Test Webhook",
		URL:    "https://example.com/webhook",
		Events: []domain.WebhookEvent{"invalid.event"},
	}
	
	webhook, err := suite.service.CreateWebhook(ctx, userID, req)
	
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), webhook)
	assert.IsType(suite.T(), &domain.ValidationError{}, err)
}

func (suite *WebhookServiceTestSuite) TestCreateWebhook_InvalidURL() {
	ctx := context.Background()
	userID := uint64(123)
	
	req := domain.WebhookCreateRequest{
		Name:   "Test Webhook",
		URL:    "https://example.com/webhook",
		Events: []domain.WebhookEvent{domain.WebhookEventURLCreated},
	}
	
	suite.deliveryService.On("ValidateWebhookURL", ctx, req.URL).Return(assert.AnError)
	
	webhook, err := suite.service.CreateWebhook(ctx, userID, req)
	
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), webhook)
	assert.IsType(suite.T(), &domain.ValidationError{}, err)
	
	suite.deliveryService.AssertExpectations(suite.T())
}

func (suite *WebhookServiceTestSuite) TestGetWebhook_Success() {
	ctx := context.Background()
	webhookID := uint64(456)
	userID := uint64(123)
	
	expectedWebhook := &domain.Webhook{
		ID:     webhookID,
		UserID: userID,
		Name:   "Test Webhook",
		URL:    "https://example.com/webhook",
		Status: domain.WebhookStatusActive,
	}
	
	suite.webhookRepo.On("GetByID", ctx, webhookID).Return(expectedWebhook, nil)
	
	webhook, err := suite.service.GetWebhook(ctx, webhookID, userID)
	
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedWebhook, webhook)
	
	suite.webhookRepo.AssertExpectations(suite.T())
}

func (suite *WebhookServiceTestSuite) TestGetWebhook_NotOwner() {
	ctx := context.Background()
	webhookID := uint64(456)
	userID := uint64(123)
	
	webhook := &domain.Webhook{
		ID:     webhookID,
		UserID: uint64(999), // Different user
		Name:   "Test Webhook",
	}
	
	suite.webhookRepo.On("GetByID", ctx, webhookID).Return(webhook, nil)
	
	result, err := suite.service.GetWebhook(ctx, webhookID, userID)
	
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.IsType(suite.T(), &domain.NotFoundError{}, err)
	
	suite.webhookRepo.AssertExpectations(suite.T())
}

func (suite *WebhookServiceTestSuite) TestUpdateWebhook_Success() {
	ctx := context.Background()
	webhookID := uint64(456)
	userID := uint64(123)
	
	existingWebhook := &domain.Webhook{
		ID:     webhookID,
		UserID: userID,
		Name:   "Old Name",
		URL:    "https://old.example.com/webhook",
		Status: domain.WebhookStatusActive,
	}
	
	newName := "New Name"
	newURL := "https://new.example.com/webhook"
	req := domain.WebhookUpdateRequest{
		Name: &newName,
		URL:  &newURL,
	}
	
	suite.webhookRepo.On("GetByID", ctx, webhookID).Return(existingWebhook, nil)
	suite.deliveryService.On("ValidateWebhookURL", ctx, newURL).Return(nil)
	suite.webhookRepo.On("Update", ctx, mock.AnythingOfType("*domain.Webhook")).Return(nil)
	suite.logger.On("InfoContext", ctx, "Webhook updated successfully", mock.Anything).Return()
	
	webhook, err := suite.service.UpdateWebhook(ctx, webhookID, userID, req)
	
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), webhook)
	assert.Equal(suite.T(), newName, webhook.Name)
	assert.Equal(suite.T(), newURL, webhook.URL)
	
	suite.webhookRepo.AssertExpectations(suite.T())
	suite.deliveryService.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

func (suite *WebhookServiceTestSuite) TestDeleteWebhook_Success() {
	ctx := context.Background()
	webhookID := uint64(456)
	userID := uint64(123)
	
	webhook := &domain.Webhook{
		ID:     webhookID,
		UserID: userID,
		Name:   "Test Webhook",
	}
	
	suite.webhookRepo.On("GetByID", ctx, webhookID).Return(webhook, nil)
	suite.webhookRepo.On("Delete", ctx, webhookID).Return(nil)
	suite.logger.On("InfoContext", ctx, "Webhook deleted successfully", mock.Anything).Return()
	
	err := suite.service.DeleteWebhook(ctx, webhookID, userID)
	
	assert.NoError(suite.T(), err)
	
	suite.webhookRepo.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

func (suite *WebhookServiceTestSuite) TestTestWebhook_Success() {
	ctx := context.Background()
	webhookID := uint64(456)
	userID := uint64(123)
	
	webhook := &domain.Webhook{
		ID:     webhookID,
		UserID: userID,
		Name:   "Test Webhook",
		URL:    "https://example.com/webhook",
		Secret: "test-secret",
	}
	
	expectedDelivery := &domain.WebhookDelivery{
		ID:        789,
		WebhookID: webhookID,
		Status:    domain.WebhookDeliveryStatusSuccess,
	}
	
	suite.webhookRepo.On("GetByID", ctx, webhookID).Return(webhook, nil)
	suite.deliveryService.On("DeliverWebhook", ctx, webhook, mock.AnythingOfType("*domain.WebhookPayload")).Return(expectedDelivery, nil)
	suite.logger.On("InfoContext", ctx, "Webhook test completed", mock.Anything).Return()
	
	delivery, err := suite.service.TestWebhook(ctx, webhookID, userID)
	
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedDelivery, delivery)
	
	suite.webhookRepo.AssertExpectations(suite.T())
	suite.deliveryService.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

func (suite *WebhookServiceTestSuite) TestTriggerEventForUser_Success() {
	ctx := context.Background()
	userID := uint64(123)
	event := domain.WebhookEventURLCreated
	data := map[string]interface{}{"test": "data"}
	
	webhooks := []*domain.Webhook{
		{
			ID:     456,
			UserID: userID,
			URL:    "https://example.com/webhook1",
			Events: []domain.WebhookEvent{event},
			Status: domain.WebhookStatusActive,
		},
		{
			ID:     789,
			UserID: userID,
			URL:    "https://example.com/webhook2",
			Events: []domain.WebhookEvent{event},
			Status: domain.WebhookStatusActive,
		},
	}
	
	suite.webhookRepo.On("GetByUserIDAndEvent", ctx, userID, event).Return(webhooks, nil)
	suite.deliveryService.On("DeliverWebhook", mock.Anything, webhooks[0], mock.AnythingOfType("*domain.WebhookPayload")).Return(&domain.WebhookDelivery{}, nil)
	suite.deliveryService.On("DeliverWebhook", mock.Anything, webhooks[1], mock.AnythingOfType("*domain.WebhookPayload")).Return(&domain.WebhookDelivery{}, nil)
	suite.logger.On("InfoContext", ctx, "Event triggered for user webhooks", mock.Anything).Return()
	
	err := suite.service.TriggerEventForUser(ctx, event, data, userID)
	
	assert.NoError(suite.T(), err)
	
	// Give goroutines time to complete
	time.Sleep(10 * time.Millisecond)
	
	suite.webhookRepo.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

func (suite *WebhookServiceTestSuite) TestTriggerEventForUser_NoWebhooks() {
	ctx := context.Background()
	userID := uint64(123)
	event := domain.WebhookEventURLCreated
	data := map[string]interface{}{"test": "data"}
	
	suite.webhookRepo.On("GetByUserIDAndEvent", ctx, userID, event).Return([]*domain.Webhook{}, nil)
	
	err := suite.service.TriggerEventForUser(ctx, event, data, userID)
	
	assert.NoError(suite.T(), err)
	
	suite.webhookRepo.AssertExpectations(suite.T())
}

func (suite *WebhookServiceTestSuite) TestActivateWebhook_Success() {
	ctx := context.Background()
	webhookID := uint64(456)
	userID := uint64(123)
	
	webhook := &domain.Webhook{
		ID:     webhookID,
		UserID: userID,
		Status: domain.WebhookStatusInactive,
	}
	
	suite.webhookRepo.On("GetByID", ctx, webhookID).Return(webhook, nil)
	suite.webhookRepo.On("UpdateStatus", ctx, webhookID, domain.WebhookStatusActive).Return(nil)
	suite.logger.On("InfoContext", ctx, "Webhook activated", mock.Anything).Return()
	
	err := suite.service.ActivateWebhook(ctx, webhookID, userID)
	
	assert.NoError(suite.T(), err)
	
	suite.webhookRepo.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

func (suite *WebhookServiceTestSuite) TestDeactivateWebhook_Success() {
	ctx := context.Background()
	webhookID := uint64(456)
	userID := uint64(123)
	
	webhook := &domain.Webhook{
		ID:     webhookID,
		UserID: userID,
		Status: domain.WebhookStatusActive,
	}
	
	suite.webhookRepo.On("GetByID", ctx, webhookID).Return(webhook, nil)
	suite.webhookRepo.On("UpdateStatus", ctx, webhookID, domain.WebhookStatusInactive).Return(nil)
	suite.logger.On("InfoContext", ctx, "Webhook deactivated", mock.Anything).Return()
	
	err := suite.service.DeactivateWebhook(ctx, webhookID, userID)
	
	assert.NoError(suite.T(), err)
	
	suite.webhookRepo.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}

func (suite *WebhookServiceTestSuite) TestGetWebhookStats_Success() {
	ctx := context.Background()
	webhookID := uint64(456)
	userID := uint64(123)
	
	webhook := &domain.Webhook{
		ID:     webhookID,
		UserID: userID,
	}
	
	expectedStats := &domain.WebhookDeliveryStats{
		TotalDeliveries:   100,
		SuccessDeliveries: 95,
		FailedDeliveries:  5,
		SuccessRate:       95.0,
	}
	
	suite.webhookRepo.On("GetByID", ctx, webhookID).Return(webhook, nil)
	suite.deliveryRepo.On("GetDeliveryStats", ctx, webhookID).Return(expectedStats, nil)
	
	stats, err := suite.service.GetWebhookStats(ctx, webhookID, userID)
	
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedStats, stats)
	
	suite.webhookRepo.AssertExpectations(suite.T())
	suite.deliveryRepo.AssertExpectations(suite.T())
}

func (suite *WebhookServiceTestSuite) TestProcessPendingRetries_Success() {
	ctx := context.Background()
	
	pendingDeliveries := []*domain.WebhookDelivery{
		{
			ID:           123,
			WebhookID:    456,
			Status:       domain.WebhookDeliveryStatusFailed,
			AttemptCount: 1,
		},
		{
			ID:           124,
			WebhookID:    457,
			Status:       domain.WebhookDeliveryStatusFailed,
			AttemptCount: 2,
		},
	}
	
	suite.deliveryRepo.On("GetPendingRetries", ctx, 100).Return(pendingDeliveries, nil)
	suite.deliveryService.On("RetryFailedDelivery", ctx, pendingDeliveries[0]).Return(&domain.WebhookDelivery{}, nil)
	suite.deliveryService.On("RetryFailedDelivery", ctx, pendingDeliveries[1]).Return(&domain.WebhookDelivery{}, nil)
	suite.logger.On("InfoContext", ctx, "Processed pending webhook retries", mock.Anything).Return()
	
	err := suite.service.ProcessPendingRetries(ctx)
	
	assert.NoError(suite.T(), err)
	
	suite.deliveryRepo.AssertExpectations(suite.T())
	suite.deliveryService.AssertExpectations(suite.T())
	suite.logger.AssertExpectations(suite.T())
}