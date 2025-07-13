package repositories

import (
	"context"
	"time"

	"gorm.io/gorm"
	"url-shortener/internal/core/domain"
	"url-shortener/internal/core/ports"
)

type webhookRepository struct {
	db *gorm.DB
}

func NewWebhookRepository(db *gorm.DB) ports.WebhookRepository {
	return &webhookRepository{db: db}
}

func (r *webhookRepository) Create(ctx context.Context, webhook *domain.Webhook) error {
	webhook.CreatedAt = time.Now()
	webhook.UpdatedAt = time.Now()
	
	return r.db.WithContext(ctx).Create(webhook).Error
}

func (r *webhookRepository) GetByID(ctx context.Context, id uint64) (*domain.Webhook, error) {
	var webhook domain.Webhook
	err := r.db.WithContext(ctx).
		Preload("User").
		First(&webhook, id).Error
	
	if err != nil {
		return nil, err
	}
	
	return &webhook, nil
}

func (r *webhookRepository) Update(ctx context.Context, webhook *domain.Webhook) error {
	webhook.UpdatedAt = time.Now()
	
	return r.db.WithContext(ctx).
		Select("*").
		Where("id = ? AND user_id = ?", webhook.ID, webhook.UserID).
		Updates(webhook).Error
}

func (r *webhookRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).
		Delete(&domain.Webhook{}, "id = ?", id).Error
}

func (r *webhookRepository) GetByUserID(ctx context.Context, userID uint64, offset, limit int) ([]*domain.Webhook, int64, error) {
	var webhooks []*domain.Webhook
	var total int64
	
	// Get total count
	err := r.db.WithContext(ctx).
		Model(&domain.Webhook{}).
		Where("user_id = ?", userID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	
	// Get paginated results
	err = r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&webhooks).Error
	
	return webhooks, total, err
}

func (r *webhookRepository) GetActiveWebhooks(ctx context.Context, event domain.WebhookEvent) ([]*domain.Webhook, error) {
	var webhooks []*domain.Webhook
	
	err := r.db.WithContext(ctx).
		Where("status = ? AND events @> ?", domain.WebhookStatusActive, `["`+string(event)+`"]`).
		Find(&webhooks).Error
	
	return webhooks, err
}

func (r *webhookRepository) GetByUserIDAndEvent(ctx context.Context, userID uint64, event domain.WebhookEvent) ([]*domain.Webhook, error) {
	var webhooks []*domain.Webhook
	
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ? AND events @> ?", 
			userID, domain.WebhookStatusActive, `["`+string(event)+`"]`).
		Find(&webhooks).Error
	
	return webhooks, err
}

func (r *webhookRepository) UpdateStatus(ctx context.Context, id uint64, status domain.WebhookStatus) error {
	return r.db.WithContext(ctx).
		Model(&domain.Webhook{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

func (r *webhookRepository) UpdateStatistics(ctx context.Context, id uint64, totalDeliveries, successDeliveries, failedDeliveries int64) error {
	return r.db.WithContext(ctx).
		Model(&domain.Webhook{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"total_deliveries":   totalDeliveries,
			"success_deliveries": successDeliveries,
			"failed_deliveries":  failedDeliveries,
			"updated_at":         time.Now(),
		}).Error
}

func (r *webhookRepository) UpdateLastDelivery(ctx context.Context, id uint64, success bool) error {
	now := time.Now()
	updates := map[string]interface{}{
		"last_delivery_at": now,
		"updated_at":       now,
	}
	
	if success {
		updates["last_success_at"] = now
	} else {
		updates["last_failure_at"] = now
	}
	
	return r.db.WithContext(ctx).
		Model(&domain.Webhook{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *webhookRepository) Find(ctx context.Context, filter *domain.WebhookFilter, offset, limit int) ([]*domain.Webhook, int64, error) {
	var webhooks []*domain.Webhook
	var total int64
	
	query := r.db.WithContext(ctx).Model(&domain.Webhook{})
	
	// Apply filters
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	
	if filter.Active != nil {
		if *filter.Active {
			query = query.Where("status = ?", domain.WebhookStatusActive)
		} else {
			query = query.Where("status != ?", domain.WebhookStatusActive)
		}
	}
	
	if len(filter.Events) > 0 {
		for _, event := range filter.Events {
			query = query.Where("events @> ?", `["`+string(event)+`"]`)
		}
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
		Find(&webhooks).Error
	
	return webhooks, total, err
}