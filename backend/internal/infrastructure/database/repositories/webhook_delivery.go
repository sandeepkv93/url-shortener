package repositories

import (
	"context"
	"database/sql"
	"time"

	"gorm.io/gorm"
	"url-shortener/internal/core/domain"
	"url-shortener/internal/core/ports"
)

type webhookDeliveryRepository struct {
	db *gorm.DB
}

func NewWebhookDeliveryRepository(db *gorm.DB) ports.WebhookDeliveryRepository {
	return &webhookDeliveryRepository{db: db}
}

func (r *webhookDeliveryRepository) Create(ctx context.Context, delivery *domain.WebhookDelivery) error {
	delivery.CreatedAt = time.Now()
	delivery.UpdatedAt = time.Now()
	
	return r.db.WithContext(ctx).Create(delivery).Error
}

func (r *webhookDeliveryRepository) GetByID(ctx context.Context, id uint64) (*domain.WebhookDelivery, error) {
	var delivery domain.WebhookDelivery
	err := r.db.WithContext(ctx).
		Preload("Webhook").
		First(&delivery, id).Error
	
	if err != nil {
		return nil, err
	}
	
	return &delivery, nil
}

func (r *webhookDeliveryRepository) Update(ctx context.Context, delivery *domain.WebhookDelivery) error {
	delivery.UpdatedAt = time.Now()
	
	return r.db.WithContext(ctx).
		Select("*").
		Where("id = ?", delivery.ID).
		Updates(delivery).Error
}

func (r *webhookDeliveryRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).
		Delete(&domain.WebhookDelivery{}, "id = ?", id).Error
}

func (r *webhookDeliveryRepository) GetByWebhookID(ctx context.Context, webhookID uint64, offset, limit int) ([]*domain.WebhookDelivery, int64, error) {
	var deliveries []*domain.WebhookDelivery
	var total int64
	
	// Get total count
	err := r.db.WithContext(ctx).
		Model(&domain.WebhookDelivery{}).
		Where("webhook_id = ?", webhookID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	
	// Get paginated results
	err = r.db.WithContext(ctx).
		Where("webhook_id = ?", webhookID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&deliveries).Error
	
	return deliveries, total, err
}

func (r *webhookDeliveryRepository) GetPendingRetries(ctx context.Context, limit int) ([]*domain.WebhookDelivery, error) {
	var deliveries []*domain.WebhookDelivery
	
	err := r.db.WithContext(ctx).
		Where("status = ? AND (next_retry_at IS NULL OR next_retry_at <= ?)", 
			domain.WebhookDeliveryStatusFailed, time.Now()).
		Where("attempt_count < ?", 10). // Hard limit to prevent infinite retries
		Order("created_at ASC").
		Limit(limit).
		Find(&deliveries).Error
	
	return deliveries, err
}

func (r *webhookDeliveryRepository) GetFailedDeliveries(ctx context.Context, webhookID uint64, limit int) ([]*domain.WebhookDelivery, error) {
	var deliveries []*domain.WebhookDelivery
	
	err := r.db.WithContext(ctx).
		Where("webhook_id = ? AND status IN (?)", 
			webhookID, []domain.WebhookDeliveryStatus{
				domain.WebhookDeliveryStatusFailed,
				domain.WebhookDeliveryStatusAbandoned,
			}).
		Order("created_at DESC").
		Limit(limit).
		Find(&deliveries).Error
	
	return deliveries, err
}

func (r *webhookDeliveryRepository) UpdateStatus(ctx context.Context, id uint64, status domain.WebhookDeliveryStatus, errorMessage string) error {
	updates := map[string]interface{}{
		"status":        status,
		"error_message": errorMessage,
		"updated_at":    time.Now(),
	}
	
	return r.db.WithContext(ctx).
		Model(&domain.WebhookDelivery{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *webhookDeliveryRepository) ScheduleRetry(ctx context.Context, id uint64, nextRetryAt *time.Time) error {
	updates := map[string]interface{}{
		"status":        domain.WebhookDeliveryStatusRetrying,
		"next_retry_at": nextRetryAt,
		"updated_at":    time.Now(),
	}
	
	return r.db.WithContext(ctx).
		Model(&domain.WebhookDelivery{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *webhookDeliveryRepository) IncrementAttempt(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).
		Model(&domain.WebhookDelivery{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"attempt_count": gorm.Expr("attempt_count + 1"),
			"updated_at":    time.Now(),
		}).Error
}

func (r *webhookDeliveryRepository) Find(ctx context.Context, filter *domain.WebhookDeliveryFilter, offset, limit int) ([]*domain.WebhookDelivery, int64, error) {
	var deliveries []*domain.WebhookDelivery
	var total int64
	
	query := r.db.WithContext(ctx).Model(&domain.WebhookDelivery{})
	
	// Apply filters
	if filter.WebhookID != nil {
		query = query.Where("webhook_id = ?", *filter.WebhookID)
	}
	
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	
	if filter.EventType != nil {
		query = query.Where("event_type = ?", *filter.EventType)
	}
	
	if filter.From != nil {
		query = query.Where("created_at >= ?", *filter.From)
	}
	
	if filter.To != nil {
		query = query.Where("created_at <= ?", *filter.To)
	}
	
	// Get total count
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	
	// Get paginated results
	err = query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&deliveries).Error
	
	return deliveries, total, err
}

func (r *webhookDeliveryRepository) GetDeliveryStats(ctx context.Context, webhookID uint64) (*domain.WebhookDeliveryStats, error) {
	var stats domain.WebhookDeliveryStats
	
	// Get basic counts
	var totalCount, successCount, failedCount int64
	
	err := r.db.WithContext(ctx).
		Model(&domain.WebhookDelivery{}).
		Where("webhook_id = ?", webhookID).
		Count(&totalCount).Error
	if err != nil {
		return nil, err
	}
	
	err = r.db.WithContext(ctx).
		Model(&domain.WebhookDelivery{}).
		Where("webhook_id = ? AND status = ?", webhookID, domain.WebhookDeliveryStatusSuccess).
		Count(&successCount).Error
	if err != nil {
		return nil, err
	}
	
	err = r.db.WithContext(ctx).
		Model(&domain.WebhookDelivery{}).
		Where("webhook_id = ? AND status IN (?)", webhookID, []domain.WebhookDeliveryStatus{
			domain.WebhookDeliveryStatusFailed,
			domain.WebhookDeliveryStatusAbandoned,
		}).
		Count(&failedCount).Error
	if err != nil {
		return nil, err
	}
	
	stats.TotalDeliveries = totalCount
	stats.SuccessDeliveries = successCount
	stats.FailedDeliveries = failedCount
	
	// Calculate success rate
	if totalCount > 0 {
		stats.SuccessRate = float64(successCount) / float64(totalCount) * 100.0
	}
	
	// Get average response time
	var avgDuration sql.NullFloat64
	err = r.db.WithContext(ctx).
		Model(&domain.WebhookDelivery{}).
		Where("webhook_id = ? AND status = ? AND duration > 0", webhookID, domain.WebhookDeliveryStatusSuccess).
		Select("AVG(duration)").
		Scan(&avgDuration).Error
	if err != nil {
		return nil, err
	}
	
	if avgDuration.Valid {
		stats.AverageResponseTime = int64(avgDuration.Float64)
	}
	
	// Get timestamp information
	var lastDelivery, lastSuccess, lastFailure *time.Time
	
	// Last delivery
	err = r.db.WithContext(ctx).
		Model(&domain.WebhookDelivery{}).
		Where("webhook_id = ?", webhookID).
		Order("created_at DESC").
		Limit(1).
		Select("created_at").
		Scan(&lastDelivery).Error
	if err == nil {
		stats.LastDeliveryAt = lastDelivery
	}
	
	// Last success
	err = r.db.WithContext(ctx).
		Model(&domain.WebhookDelivery{}).
		Where("webhook_id = ? AND status = ?", webhookID, domain.WebhookDeliveryStatusSuccess).
		Order("created_at DESC").
		Limit(1).
		Select("created_at").
		Scan(&lastSuccess).Error
	if err == nil {
		stats.LastSuccessAt = lastSuccess
	}
	
	// Last failure
	err = r.db.WithContext(ctx).
		Model(&domain.WebhookDelivery{}).
		Where("webhook_id = ? AND status IN (?)", webhookID, []domain.WebhookDeliveryStatus{
			domain.WebhookDeliveryStatusFailed,
			domain.WebhookDeliveryStatusAbandoned,
		}).
		Order("created_at DESC").
		Limit(1).
		Select("created_at").
		Scan(&lastFailure).Error
	if err == nil {
		stats.LastFailureAt = lastFailure
	}
	
	return &stats, nil
}

func (r *webhookDeliveryRepository) DeleteOldDeliveries(ctx context.Context, olderThan time.Time) (int64, error) {
	result := r.db.WithContext(ctx).
		Where("created_at < ?", olderThan).
		Delete(&domain.WebhookDelivery{})
	
	return result.RowsAffected, result.Error
}