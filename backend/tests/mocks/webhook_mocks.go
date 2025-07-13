package mocks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
	"url-shortener/internal/core/domain"
)

// MockWebhookRepository implements ports.WebhookRepository
type MockWebhookRepository struct {
	mock.Mock
}

func (m *MockWebhookRepository) Create(ctx context.Context, webhook *domain.Webhook) error {
	args := m.Called(ctx, webhook)
	return args.Error(0)
}

func (m *MockWebhookRepository) GetByID(ctx context.Context, id uint64) (*domain.Webhook, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Webhook), args.Error(1)
}

func (m *MockWebhookRepository) Update(ctx context.Context, webhook *domain.Webhook) error {
	args := m.Called(ctx, webhook)
	return args.Error(0)
}

func (m *MockWebhookRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockWebhookRepository) GetByUserID(ctx context.Context, userID uint64, offset, limit int) ([]*domain.Webhook, int64, error) {
	args := m.Called(ctx, userID, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Webhook), args.Get(1).(int64), args.Error(2)
}

func (m *MockWebhookRepository) GetActiveWebhooks(ctx context.Context, event domain.WebhookEvent) ([]*domain.Webhook, error) {
	args := m.Called(ctx, event)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Webhook), args.Error(1)
}

func (m *MockWebhookRepository) GetByUserIDAndEvent(ctx context.Context, userID uint64, event domain.WebhookEvent) ([]*domain.Webhook, error) {
	args := m.Called(ctx, userID, event)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Webhook), args.Error(1)
}

func (m *MockWebhookRepository) UpdateStatus(ctx context.Context, id uint64, status domain.WebhookStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockWebhookRepository) UpdateStatistics(ctx context.Context, id uint64, totalDeliveries, successDeliveries, failedDeliveries int64) error {
	args := m.Called(ctx, id, totalDeliveries, successDeliveries, failedDeliveries)
	return args.Error(0)
}

func (m *MockWebhookRepository) UpdateLastDelivery(ctx context.Context, id uint64, success bool) error {
	args := m.Called(ctx, id, success)
	return args.Error(0)
}

func (m *MockWebhookRepository) Find(ctx context.Context, filter *domain.WebhookFilter, offset, limit int) ([]*domain.Webhook, int64, error) {
	args := m.Called(ctx, filter, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Webhook), args.Get(1).(int64), args.Error(2)
}

// MockWebhookDeliveryRepository implements ports.WebhookDeliveryRepository
type MockWebhookDeliveryRepository struct {
	mock.Mock
}

func (m *MockWebhookDeliveryRepository) Create(ctx context.Context, delivery *domain.WebhookDelivery) error {
	args := m.Called(ctx, delivery)
	return args.Error(0)
}

func (m *MockWebhookDeliveryRepository) GetByID(ctx context.Context, id uint64) (*domain.WebhookDelivery, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.WebhookDelivery), args.Error(1)
}

func (m *MockWebhookDeliveryRepository) Update(ctx context.Context, delivery *domain.WebhookDelivery) error {
	args := m.Called(ctx, delivery)
	return args.Error(0)
}

func (m *MockWebhookDeliveryRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockWebhookDeliveryRepository) GetByWebhookID(ctx context.Context, webhookID uint64, offset, limit int) ([]*domain.WebhookDelivery, int64, error) {
	args := m.Called(ctx, webhookID, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.WebhookDelivery), args.Get(1).(int64), args.Error(2)
}

func (m *MockWebhookDeliveryRepository) GetPendingRetries(ctx context.Context, limit int) ([]*domain.WebhookDelivery, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.WebhookDelivery), args.Error(1)
}

func (m *MockWebhookDeliveryRepository) GetFailedDeliveries(ctx context.Context, webhookID uint64, limit int) ([]*domain.WebhookDelivery, error) {
	args := m.Called(ctx, webhookID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.WebhookDelivery), args.Error(1)
}

func (m *MockWebhookDeliveryRepository) UpdateStatus(ctx context.Context, id uint64, status domain.WebhookDeliveryStatus, errorMessage string) error {
	args := m.Called(ctx, id, status, errorMessage)
	return args.Error(0)
}

func (m *MockWebhookDeliveryRepository) ScheduleRetry(ctx context.Context, id uint64, nextRetryAt *time.Time) error {
	args := m.Called(ctx, id, nextRetryAt)
	return args.Error(0)
}

func (m *MockWebhookDeliveryRepository) IncrementAttempt(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockWebhookDeliveryRepository) Find(ctx context.Context, filter *domain.WebhookDeliveryFilter, offset, limit int) ([]*domain.WebhookDelivery, int64, error) {
	args := m.Called(ctx, filter, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.WebhookDelivery), args.Get(1).(int64), args.Error(2)
}

func (m *MockWebhookDeliveryRepository) GetDeliveryStats(ctx context.Context, webhookID uint64) (*domain.WebhookDeliveryStats, error) {
	args := m.Called(ctx, webhookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.WebhookDeliveryStats), args.Error(1)
}

func (m *MockWebhookDeliveryRepository) DeleteOldDeliveries(ctx context.Context, olderThan time.Time) (int64, error) {
	args := m.Called(ctx, olderThan)
	return args.Get(0).(int64), args.Error(1)
}

// MockWebhookService implements ports.WebhookService
type MockWebhookService struct {
	mock.Mock
}

func (m *MockWebhookService) CreateWebhook(ctx context.Context, userID uint64, req domain.WebhookCreateRequest) (*domain.Webhook, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Webhook), args.Error(1)
}

func (m *MockWebhookService) GetWebhook(ctx context.Context, webhookID uint64, userID uint64) (*domain.Webhook, error) {
	args := m.Called(ctx, webhookID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Webhook), args.Error(1)
}

func (m *MockWebhookService) GetUserWebhooks(ctx context.Context, userID uint64, offset, limit int) ([]*domain.Webhook, int64, error) {
	args := m.Called(ctx, userID, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Webhook), args.Get(1).(int64), args.Error(2)
}

func (m *MockWebhookService) UpdateWebhook(ctx context.Context, webhookID uint64, userID uint64, req domain.WebhookUpdateRequest) (*domain.Webhook, error) {
	args := m.Called(ctx, webhookID, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Webhook), args.Error(1)
}

func (m *MockWebhookService) DeleteWebhook(ctx context.Context, webhookID uint64, userID uint64) error {
	args := m.Called(ctx, webhookID, userID)
	return args.Error(0)
}

func (m *MockWebhookService) ActivateWebhook(ctx context.Context, webhookID uint64, userID uint64) error {
	args := m.Called(ctx, webhookID, userID)
	return args.Error(0)
}

func (m *MockWebhookService) DeactivateWebhook(ctx context.Context, webhookID uint64, userID uint64) error {
	args := m.Called(ctx, webhookID, userID)
	return args.Error(0)
}

func (m *MockWebhookService) TestWebhook(ctx context.Context, webhookID uint64, userID uint64) (*domain.WebhookDelivery, error) {
	args := m.Called(ctx, webhookID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.WebhookDelivery), args.Error(1)
}

func (m *MockWebhookService) TriggerEvent(ctx context.Context, event domain.WebhookEvent, data interface{}, userID uint64) error {
	args := m.Called(ctx, event, data, userID)
	return args.Error(0)
}

func (m *MockWebhookService) TriggerEventForUser(ctx context.Context, event domain.WebhookEvent, data interface{}, userID uint64) error {
	args := m.Called(ctx, event, data, userID)
	return args.Error(0)
}

func (m *MockWebhookService) TriggerEventGlobally(ctx context.Context, event domain.WebhookEvent, data interface{}) error {
	args := m.Called(ctx, event, data)
	return args.Error(0)
}

func (m *MockWebhookService) GetWebhookDeliveries(ctx context.Context, webhookID uint64, userID uint64, offset, limit int) ([]*domain.WebhookDelivery, int64, error) {
	args := m.Called(ctx, webhookID, userID, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.WebhookDelivery), args.Get(1).(int64), args.Error(2)
}

func (m *MockWebhookService) GetDelivery(ctx context.Context, deliveryID uint64, userID uint64) (*domain.WebhookDelivery, error) {
	args := m.Called(ctx, deliveryID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.WebhookDelivery), args.Error(1)
}

func (m *MockWebhookService) RetryDelivery(ctx context.Context, deliveryID uint64, userID uint64) error {
	args := m.Called(ctx, deliveryID, userID)
	return args.Error(0)
}

func (m *MockWebhookService) GetWebhookStats(ctx context.Context, webhookID uint64, userID uint64) (*domain.WebhookDeliveryStats, error) {
	args := m.Called(ctx, webhookID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.WebhookDeliveryStats), args.Error(1)
}

func (m *MockWebhookService) GetFailedDeliveries(ctx context.Context, webhookID uint64, userID uint64, limit int) ([]*domain.WebhookDelivery, error) {
	args := m.Called(ctx, webhookID, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.WebhookDelivery), args.Error(1)
}

func (m *MockWebhookService) ProcessPendingRetries(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockWebhookService) CleanupOldDeliveries(ctx context.Context, olderThanDays int) (int64, error) {
	args := m.Called(ctx, olderThanDays)
	return args.Get(0).(int64), args.Error(1)
}

// MockWebhookDeliveryService implements ports.WebhookDeliveryService
type MockWebhookDeliveryService struct {
	mock.Mock
}

func (m *MockWebhookDeliveryService) DeliverWebhook(ctx context.Context, webhook *domain.Webhook, payload *domain.WebhookPayload) (*domain.WebhookDelivery, error) {
	args := m.Called(ctx, webhook, payload)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.WebhookDelivery), args.Error(1)
}

func (m *MockWebhookDeliveryService) RetryFailedDelivery(ctx context.Context, delivery *domain.WebhookDelivery) (*domain.WebhookDelivery, error) {
	args := m.Called(ctx, delivery)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.WebhookDelivery), args.Error(1)
}

func (m *MockWebhookDeliveryService) ScheduleRetry(ctx context.Context, delivery *domain.WebhookDelivery) error {
	args := m.Called(ctx, delivery)
	return args.Error(0)
}

func (m *MockWebhookDeliveryService) ProcessPendingDeliveries(ctx context.Context, limit int) (int, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).(int), args.Error(1)
}

func (m *MockWebhookDeliveryService) ProcessFailedDeliveries(ctx context.Context, limit int) (int, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).(int), args.Error(1)
}

func (m *MockWebhookDeliveryService) ValidateWebhookURL(ctx context.Context, url string) error {
	args := m.Called(ctx, url)
	return args.Error(0)
}

func (m *MockWebhookDeliveryService) GenerateSignature(payload []byte, secret string) string {
	args := m.Called(payload, secret)
	return args.String(0)
}

func (m *MockWebhookDeliveryService) VerifySignature(payload []byte, signature, secret string) bool {
	args := m.Called(payload, signature, secret)
	return args.Bool(0)
}