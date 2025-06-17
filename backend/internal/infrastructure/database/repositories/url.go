package repositories

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"url-shortener/internal/core/domain"
	"url-shortener/internal/core/ports"
)

type urlRepository struct {
	db *gorm.DB
}

func NewURLRepository(db *gorm.DB) ports.URLRepository {
	return &urlRepository{
		db: db,
	}
}

func (r *urlRepository) Create(ctx context.Context, url *domain.ShortURL) error {
	if err := r.db.WithContext(ctx).Create(url).Error; err != nil {
		if isDuplicateKeyError(err) {
			return domain.ErrShortCodeExists
		}
		return fmt.Errorf("failed to create short URL: %w", err)
	}
	return nil
}

func (r *urlRepository) GetByID(ctx context.Context, id uint) (*domain.ShortURL, error) {
	var url domain.ShortURL
	if err := r.db.WithContext(ctx).
		Preload("User").
		First(&url, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrShortURLNotFound
		}
		return nil, fmt.Errorf("failed to get short URL by id: %w", err)
	}
	return &url, nil
}

func (r *urlRepository) GetByShortCode(ctx context.Context, shortCode string) (*domain.ShortURL, error) {
	var url domain.ShortURL
	if err := r.db.WithContext(ctx).
		Preload("User").
		Where("short_code = ?", shortCode).
		First(&url).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrShortURLNotFound
		}
		return nil, fmt.Errorf("failed to get short URL by code: %w", err)
	}
	return &url, nil
}

func (r *urlRepository) Update(ctx context.Context, url *domain.ShortURL) error {
	if err := r.db.WithContext(ctx).Save(url).Error; err != nil {
		if isDuplicateKeyError(err) {
			return domain.ErrShortCodeExists
		}
		return fmt.Errorf("failed to update short URL: %w", err)
	}
	return nil
}

func (r *urlRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&domain.ShortURL{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete short URL: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrShortURLNotFound
	}
	return nil
}

func (r *urlRepository) ExistsByShortCode(ctx context.Context, shortCode string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&domain.ShortURL{}).
		Where("short_code = ?", shortCode).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check short code existence: %w", err)
	}
	return count > 0, nil
}

func (r *urlRepository) GetByUserID(ctx context.Context, userID uint, offset, limit int) ([]*domain.ShortURL, int64, error) {
	var urls []*domain.ShortURL
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).
		Model(&domain.ShortURL{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count user URLs: %w", err)
	}

	// Get URLs with pagination
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&urls).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list user URLs: %w", err)
	}

	return urls, total, nil
}

func (r *urlRepository) GetActiveByShortCode(ctx context.Context, shortCode string) (*domain.ShortURL, error) {
	var url domain.ShortURL
	now := time.Now()
	if err := r.db.WithContext(ctx).
		Where("short_code = ? AND is_active = ? AND (expires_at IS NULL OR expires_at > ?)", shortCode, true, now).
		First(&url).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrShortURLNotFound
		}
		return nil, fmt.Errorf("failed to get active short URL: %w", err)
	}
	return &url, nil
}

func (r *urlRepository) IncrementClickCount(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).
		Model(&domain.ShortURL{}).
		Where("id = ?", id).
		Update("click_count", gorm.Expr("click_count + ?", 1)).Error; err != nil {
		return fmt.Errorf("failed to increment click count: %w", err)
	}
	return nil
}

func (r *urlRepository) GetExpiredURLs(ctx context.Context, limit int) ([]*domain.ShortURL, error) {
	var urls []*domain.ShortURL
	now := time.Now()
	if err := r.db.WithContext(ctx).
		Where("expires_at IS NOT NULL AND expires_at <= ?", now).
		Limit(limit).
		Find(&urls).Error; err != nil {
		return nil, fmt.Errorf("failed to get expired URLs: %w", err)
	}
	return urls, nil
}

func (r *urlRepository) GetTotalURLs(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&domain.ShortURL{}).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count total URLs: %w", err)
	}
	return count, nil
}

func (r *urlRepository) GetTotalURLsByUser(ctx context.Context, userID uint) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&domain.ShortURL{}).
		Where("user_id = ?", userID).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count user URLs: %w", err)
	}
	return count, nil
}

func (r *urlRepository) GetPopularURLs(ctx context.Context, limit int) ([]*domain.ShortURL, error) {
	var urls []*domain.ShortURL
	if err := r.db.WithContext(ctx).
		Preload("User").
		Where("is_active = ?", true).
		Order("click_count DESC").
		Limit(limit).
		Find(&urls).Error; err != nil {
		return nil, fmt.Errorf("failed to get popular URLs: %w", err)
	}
	return urls, nil
}