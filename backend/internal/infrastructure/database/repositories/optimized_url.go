package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"url-shortener/internal/core/domain"
	"url-shortener/internal/infrastructure/database"
)

// OptimizedURLRepository extends URLRepository with performance optimizations
type OptimizedURLRepository struct {
	*urlRepository
	optimizer   *database.OptimizedQuery
	transaction *database.TransactionOptimizer
}

func NewOptimizedURLRepository(db *gorm.DB) *OptimizedURLRepository {
	return &OptimizedURLRepository{
		urlRepository: &urlRepository{db: db},
		optimizer:     database.NewOptimizedQuery(db),
		transaction:   database.NewTransactionOptimizer(db),
	}
}

// GetByShortCodeOptimized uses optimized query with selective column loading
func (r *OptimizedURLRepository) GetByShortCodeOptimized(ctx context.Context, shortCode string) (*domain.ShortURL, error) {
	var url domain.ShortURL
	
	// Only load essential columns for performance
	err := r.db.WithContext(ctx).
		Select("id", "short_code", "original_url", "user_id", "is_active", "expires_at", "password_hash").
		Where("short_code = ? AND is_active = ?", shortCode, true).
		First(&url).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrShortURLNotFound
		}
		return nil, fmt.Errorf("failed to get short URL by code: %w", err)
	}
	
	return &url, nil
}

// GetUserURLsWithPagination provides optimized pagination for large datasets
func (r *OptimizedURLRepository) GetUserURLsWithPagination(ctx context.Context, userID uint, offset, limit int, sortBy string) ([]*domain.ShortURL, int64, error) {
	var urls []*domain.ShortURL
	var total int64
	
	// Get total count efficiently
	baseQuery := r.db.WithContext(ctx).Model(&domain.ShortURL{}).Where("user_id = ?", userID)
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count user URLs: %w", err)
	}
	
	// Build optimized query with sorting
	query := r.db.WithContext(ctx).
		Select("id", "short_code", "original_url", "title", "created_at", "updated_at", "click_count", "is_active").
		Where("user_id = ?", userID).
		Offset(offset).
		Limit(limit)
	
	// Apply sorting
	switch sortBy {
	case "created_desc":
		query = query.Order("created_at DESC")
	case "created_asc":
		query = query.Order("created_at ASC")
	case "clicks_desc":
		query = query.Order("click_count DESC")
	case "clicks_asc":
		query = query.Order("click_count ASC")
	case "title":
		query = query.Order("title ASC")
	default:
		query = query.Order("created_at DESC")
	}
	
	if err := query.Find(&urls).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get user URLs: %w", err)
	}
	
	return urls, total, nil
}

// BulkUpdateClickCounts efficiently updates click counts for multiple URLs
func (r *OptimizedURLRepository) BulkUpdateClickCounts(ctx context.Context, updates map[string]int) error {
	if len(updates) == 0 {
		return nil
	}
	
	return r.transaction.OptimizedTransaction(ctx, func(tx *gorm.DB) error {
		for shortCode, increment := range updates {
			err := tx.Model(&domain.ShortURL{}).
				Where("short_code = ?", shortCode).
				Update("click_count", gorm.Expr("click_count + ?", increment)).Error
			if err != nil {
				return fmt.Errorf("failed to update click count for %s: %w", shortCode, err)
			}
		}
		return nil
	})
}

// GetPopularURLsOptimized retrieves popular URLs with minimal data loading
func (r *OptimizedURLRepository) GetPopularURLsOptimized(ctx context.Context, limit int, timeWindow time.Duration) ([]*domain.ShortURL, error) {
	var urls []*domain.ShortURL
	
	since := time.Now().Add(-timeWindow)
	
	err := r.db.WithContext(ctx).
		Select("short_urls.id", "short_urls.short_code", "short_urls.original_url", "short_urls.title", "COUNT(clicks.id) as recent_clicks").
		Table("short_urls").
		Joins("LEFT JOIN clicks ON short_urls.id = clicks.short_url_id AND clicks.clicked_at > ?", since).
		Where("short_urls.is_active = ? AND short_urls.is_public = ?", true, true).
		Group("short_urls.id").
		Order("recent_clicks DESC").
		Limit(limit).
		Find(&urls).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to get popular URLs: %w", err)
	}
	
	return urls, nil
}

// SearchURLsOptimized provides optimized full-text search
func (r *OptimizedURLRepository) SearchURLsOptimized(ctx context.Context, userID uint, query string, limit int) ([]*domain.ShortURL, error) {
	var urls []*domain.ShortURL
	
	// Use optimized search with indexed columns
	searchQuery := "%" + strings.ToLower(query) + "%"
	
	err := r.db.WithContext(ctx).
		Select("id", "short_code", "original_url", "title", "description", "created_at").
		Where("user_id = ? AND (LOWER(title) LIKE ? OR LOWER(original_url) LIKE ? OR LOWER(short_code) LIKE ?)", 
			userID, searchQuery, searchQuery, searchQuery).
		Order("created_at DESC").
		Limit(limit).
		Find(&urls).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to search URLs: %w", err)
	}
	
	return urls, nil
}

// GetURLStatsOptimized retrieves URL statistics efficiently
func (r *OptimizedURLRepository) GetURLStatsOptimized(ctx context.Context, shortCode string, startDate, endDate time.Time) (*domain.URLStats, error) {
	stats := &domain.URLStats{
		ShortCode: shortCode,
	}
	
	// Get basic URL info
	var url domain.ShortURL
	err := r.db.WithContext(ctx).
		Select("id", "short_code", "original_url", "created_at", "click_count").
		Where("short_code = ?", shortCode).
		First(&url).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get URL: %w", err)
	}
	
	stats.URL = &url
	stats.TotalClicks = url.ClickCount
	
	// Get clicks within date range
	var clicksInRange int64
	err = r.db.WithContext(ctx).
		Model(&domain.Click{}).
		Where("short_url_id = ? AND clicked_at BETWEEN ? AND ?", url.ID, startDate, endDate).
		Count(&clicksInRange).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count clicks in range: %w", err)
	}
	
	stats.ClicksInRange = clicksInRange
	
	// Get unique visitors in range
	var uniqueVisitors int64
	err = r.db.WithContext(ctx).
		Model(&domain.Click{}).
		Select("COUNT(DISTINCT ip_address)").
		Where("short_url_id = ? AND clicked_at BETWEEN ? AND ?", url.ID, startDate, endDate).
		Count(&uniqueVisitors).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count unique visitors: %w", err)
	}
	
	stats.UniqueVisitors = uniqueVisitors
	
	return stats, nil
}

// ExistsOptimized checks if a short code exists without loading the full record
func (r *OptimizedURLRepository) ExistsOptimized(ctx context.Context, shortCode string) (bool, error) {
	return r.optimizer.ExistsCheck(ctx, &domain.ShortURL{}, "short_code = ?", shortCode)
}

// BatchCreateOptimized creates multiple URLs efficiently
func (r *OptimizedURLRepository) BatchCreateOptimized(ctx context.Context, urls []*domain.ShortURL) error {
	if len(urls) == 0 {
		return nil
	}
	
	return r.optimizer.BatchInsert(ctx, urls, 100) // Process in batches of 100
}

// CleanupExpiredURLs removes expired URLs efficiently
func (r *OptimizedURLRepository) CleanupExpiredURLs(ctx context.Context, batchSize int) (int64, error) {
	now := time.Now()
	
	result := r.db.WithContext(ctx).
		Where("expires_at IS NOT NULL AND expires_at < ?", now).
		Limit(batchSize).
		Delete(&domain.ShortURL{})
	
	if result.Error != nil {
		return 0, fmt.Errorf("failed to cleanup expired URLs: %w", result.Error)
	}
	
	return result.RowsAffected, nil
}

// GetUserURLCountOptimized efficiently counts user URLs
func (r *OptimizedURLRepository) GetUserURLCountOptimized(ctx context.Context, userID uint) (int64, error) {
	return r.optimizer.CountWithoutLoad(ctx, &domain.ShortURL{}, "user_id = ?", userID)
}

// GetRecentURLsOptimized gets recent URLs with minimal data loading
func (r *OptimizedURLRepository) GetRecentURLsOptimized(ctx context.Context, userID uint, limit int) ([]*domain.ShortURL, error) {
	var urls []*domain.ShortURL
	
	err := r.db.WithContext(ctx).
		Select("id", "short_code", "original_url", "title", "created_at", "click_count").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&urls).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to get recent URLs: %w", err)
	}
	
	return urls, nil
}

// UpdateLastAccessedOptimized efficiently updates last accessed timestamp
func (r *OptimizedURLRepository) UpdateLastAccessedOptimized(ctx context.Context, shortCode string) error {
	return r.db.WithContext(ctx).
		Model(&domain.ShortURL{}).
		Where("short_code = ?", shortCode).
		Update("last_accessed_at", time.Now()).Error
}

// GetURLsByStatusOptimized gets URLs by status with pagination
func (r *OptimizedURLRepository) GetURLsByStatusOptimized(ctx context.Context, userID uint, isActive bool, offset, limit int) ([]*domain.ShortURL, int64, error) {
	var urls []*domain.ShortURL
	var total int64
	
	// Count total
	countQuery := r.db.WithContext(ctx).Model(&domain.ShortURL{}).
		Where("user_id = ? AND is_active = ?", userID, isActive)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count URLs by status: %w", err)
	}
	
	// Get records
	err := r.db.WithContext(ctx).
		Select("id", "short_code", "original_url", "title", "created_at", "click_count", "is_active").
		Where("user_id = ? AND is_active = ?", userID, isActive).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&urls).Error
	
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get URLs by status: %w", err)
	}
	
	return urls, total, nil
}

// ArchiveOldURLsOptimized archives old URLs that haven't been accessed recently
func (r *OptimizedURLRepository) ArchiveOldURLsOptimized(ctx context.Context, olderThan time.Time, batchSize int) (int64, error) {
	result := r.db.WithContext(ctx).
		Model(&domain.ShortURL{}).
		Where("last_accessed_at < ? OR (last_accessed_at IS NULL AND created_at < ?)", olderThan, olderThan).
		Limit(batchSize).
		Update("is_active", false)
	
	if result.Error != nil {
		return 0, fmt.Errorf("failed to archive old URLs: %w", result.Error)
	}
	
	return result.RowsAffected, nil
}