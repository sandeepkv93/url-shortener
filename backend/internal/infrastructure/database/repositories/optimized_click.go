package repositories

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"url-shortener/internal/core/domain"
	"url-shortener/internal/infrastructure/database"
)

// OptimizedClickRepository extends ClickRepository with performance optimizations
type OptimizedClickRepository struct {
	*clickRepository
	optimizer   *database.OptimizedQuery
	transaction *database.TransactionOptimizer
}

func NewOptimizedClickRepository(db *gorm.DB) *OptimizedClickRepository {
	return &OptimizedClickRepository{
		clickRepository: &clickRepository{db: db},
		optimizer:       database.NewOptimizedQuery(db),
		transaction:     database.NewTransactionOptimizer(db),
	}
}

// RecordClickOptimized efficiently records clicks with batch processing
func (r *OptimizedClickRepository) RecordClickOptimized(ctx context.Context, click *domain.Click) error {
	// Use a more efficient insert that doesn't return the full record
	return r.db.WithContext(ctx).
		Select("short_url_id", "ip_address", "user_agent", "referer", "country", "city", "clicked_at").
		Create(click).Error
}

// BatchRecordClicks efficiently records multiple clicks
func (r *OptimizedClickRepository) BatchRecordClicks(ctx context.Context, clicks []*domain.Click) error {
	if len(clicks) == 0 {
		return nil
	}
	
	return r.optimizer.BatchInsert(ctx, clicks, 500) // Process in larger batches for clicks
}

// GetClickStatsOptimized retrieves click statistics with optimized queries
func (r *OptimizedClickRepository) GetClickStatsOptimized(ctx context.Context, shortURLID uint, startDate, endDate time.Time) (*domain.ClickStats, error) {
	stats := &domain.ClickStats{
		ShortURLID: shortURLID,
		StartDate:  &startDate,
		EndDate:    &endDate,
	}
	
	// Use read-only transaction for better performance
	return stats, r.transaction.ReadOnlyTransaction(ctx, func(tx *gorm.DB) error {
		// Get total clicks in date range
		err := tx.Model(&domain.Click{}).
			Where("short_url_id = ? AND clicked_at BETWEEN ? AND ?", shortURLID, startDate, endDate).
			Count(&stats.TotalClicks).Error
		if err != nil {
			return fmt.Errorf("failed to count total clicks: %w", err)
		}
		
		// Get unique visitors (distinct IP addresses)
		err = tx.Model(&domain.Click{}).
			Select("COUNT(DISTINCT ip_address)").
			Where("short_url_id = ? AND clicked_at BETWEEN ? AND ?", shortURLID, startDate, endDate).
			Scan(&stats.UniqueVisitors).Error
		if err != nil {
			return fmt.Errorf("failed to count unique visitors: %w", err)
		}
		
		// Get clicks by hour (for timeline)
		var hourlyData []struct {
			Hour   int   `json:"hour"`
			Clicks int64 `json:"clicks"`
		}
		
		err = tx.Model(&domain.Click{}).
			Select("EXTRACT(hour from clicked_at) as hour, COUNT(*) as clicks").
			Where("short_url_id = ? AND clicked_at BETWEEN ? AND ?", shortURLID, startDate, endDate).
			Group("EXTRACT(hour from clicked_at)").
			Order("hour").
			Scan(&hourlyData).Error
		if err != nil {
			return fmt.Errorf("failed to get hourly clicks: %w", err)
		}
		
		// Convert to map for easier access
		stats.ClicksByHour = make(map[int]int64)
		for _, data := range hourlyData {
			stats.ClicksByHour[data.Hour] = data.Clicks
		}
		
		return nil
	})
}

// GetGeoStatsOptimized retrieves geographic statistics efficiently
func (r *OptimizedClickRepository) GetGeoStatsOptimized(ctx context.Context, shortURLID uint, startDate, endDate time.Time, limit int) ([]*domain.GeoStats, error) {
	var geoStats []*domain.GeoStats
	
	err := r.db.WithContext(ctx).
		Model(&domain.Click{}).
		Select("country, city, COUNT(*) as click_count").
		Where("short_url_id = ? AND clicked_at BETWEEN ? AND ? AND country IS NOT NULL", shortURLID, startDate, endDate).
		Group("country, city").
		Order("click_count DESC").
		Limit(limit).
		Find(&geoStats).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to get geo stats: %w", err)
	}
	
	return geoStats, nil
}

// GetTopCountriesOptimized gets top countries by click count
func (r *OptimizedClickRepository) GetTopCountriesOptimized(ctx context.Context, shortURLID uint, startDate, endDate time.Time, limit int) ([]*domain.CountryStats, error) {
	var countryStats []*domain.CountryStats
	
	err := r.db.WithContext(ctx).
		Model(&domain.Click{}).
		Select("country, COUNT(*) as click_count, COUNT(DISTINCT ip_address) as unique_visitors").
		Where("short_url_id = ? AND clicked_at BETWEEN ? AND ? AND country IS NOT NULL", shortURLID, startDate, endDate).
		Group("country").
		Order("click_count DESC").
		Limit(limit).
		Find(&countryStats).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to get country stats: %w", err)
	}
	
	return countryStats, nil
}

// GetDeviceStatsOptimized retrieves device statistics efficiently
func (r *OptimizedClickRepository) GetDeviceStatsOptimized(ctx context.Context, shortURLID uint, startDate, endDate time.Time) (*domain.DeviceStats, error) {
	stats := &domain.DeviceStats{
		ShortURLID: shortURLID,
	}
	
	// Get device type statistics
	var deviceData []struct {
		DeviceType string `json:"device_type"`
		Count      int64  `json:"count"`
	}
	
	err := r.db.WithContext(ctx).
		Model(&domain.Click{}).
		Select("CASE WHEN user_agent LIKE '%Mobile%' THEN 'Mobile' WHEN user_agent LIKE '%Tablet%' THEN 'Tablet' ELSE 'Desktop' END as device_type, COUNT(*) as count").
		Where("short_url_id = ? AND clicked_at BETWEEN ? AND ?", shortURLID, startDate, endDate).
		Group("device_type").
		Scan(&deviceData).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to get device stats: %w", err)
	}
	
	// Convert to structured format
	stats.DeviceTypes = make(map[string]int64)
	for _, data := range deviceData {
		stats.DeviceTypes[data.DeviceType] = data.Count
	}
	
	// Get browser statistics
	var browserData []struct {
		Browser string `json:"browser"`
		Count   int64  `json:"count"`
	}
	
	err = r.db.WithContext(ctx).
		Model(&domain.Click{}).
		Select("CASE WHEN user_agent LIKE '%Chrome%' THEN 'Chrome' WHEN user_agent LIKE '%Firefox%' THEN 'Firefox' WHEN user_agent LIKE '%Safari%' THEN 'Safari' WHEN user_agent LIKE '%Edge%' THEN 'Edge' ELSE 'Other' END as browser, COUNT(*) as count").
		Where("short_url_id = ? AND clicked_at BETWEEN ? AND ?", shortURLID, startDate, endDate).
		Group("browser").
		Order("count DESC").
		Scan(&browserData).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to get browser stats: %w", err)
	}
	
	stats.Browsers = make(map[string]int64)
	for _, data := range browserData {
		stats.Browsers[data.Browser] = data.Count
	}
	
	return stats, nil
}

// GetReferrerStatsOptimized retrieves referrer statistics efficiently
func (r *OptimizedClickRepository) GetReferrerStatsOptimized(ctx context.Context, shortURLID uint, startDate, endDate time.Time, limit int) ([]*domain.ReferrerStats, error) {
	var referrerStats []*domain.ReferrerStats
	
	err := r.db.WithContext(ctx).
		Model(&domain.Click{}).
		Select("COALESCE(referer, 'Direct') as referrer, COUNT(*) as click_count").
		Where("short_url_id = ? AND clicked_at BETWEEN ? AND ?", shortURLID, startDate, endDate).
		Group("referrer").
		Order("click_count DESC").
		Limit(limit).
		Find(&referrerStats).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to get referrer stats: %w", err)
	}
	
	return referrerStats, nil
}

// GetTimelineStatsOptimized retrieves timeline statistics with configurable granularity
func (r *OptimizedClickRepository) GetTimelineStatsOptimized(ctx context.Context, shortURLID uint, startDate, endDate time.Time, granularity string) ([]*domain.TimelineStats, error) {
	var timelineStats []*domain.TimelineStats
	
	var selectClause string
	var groupClause string
	
	switch granularity {
	case "hour":
		selectClause = "DATE_TRUNC('hour', clicked_at) as time_bucket, COUNT(*) as click_count"
		groupClause = "DATE_TRUNC('hour', clicked_at)"
	case "day":
		selectClause = "DATE_TRUNC('day', clicked_at) as time_bucket, COUNT(*) as click_count"
		groupClause = "DATE_TRUNC('day', clicked_at)"
	case "week":
		selectClause = "DATE_TRUNC('week', clicked_at) as time_bucket, COUNT(*) as click_count"
		groupClause = "DATE_TRUNC('week', clicked_at)"
	case "month":
		selectClause = "DATE_TRUNC('month', clicked_at) as time_bucket, COUNT(*) as click_count"
		groupClause = "DATE_TRUNC('month', clicked_at)"
	default:
		selectClause = "DATE_TRUNC('day', clicked_at) as time_bucket, COUNT(*) as click_count"
		groupClause = "DATE_TRUNC('day', clicked_at)"
	}
	
	err := r.db.WithContext(ctx).
		Model(&domain.Click{}).
		Select(selectClause).
		Where("short_url_id = ? AND clicked_at BETWEEN ? AND ?", shortURLID, startDate, endDate).
		Group(groupClause).
		Order("time_bucket").
		Scan(&timelineStats).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to get timeline stats: %w", err)
	}
	
	return timelineStats, nil
}

// GetClickHeatmapOptimized generates heatmap data for click patterns
func (r *OptimizedClickRepository) GetClickHeatmapOptimized(ctx context.Context, shortURLID uint, startDate, endDate time.Time) (*domain.ClickHeatmap, error) {
	heatmap := &domain.ClickHeatmap{
		ShortURLID: shortURLID,
		Data:       make(map[string]map[int]int64), // [day_of_week][hour] = click_count
	}
	
	var heatmapData []struct {
		DayOfWeek int   `json:"day_of_week"`
		Hour      int   `json:"hour"`
		Clicks    int64 `json:"clicks"`
	}
	
	err := r.db.WithContext(ctx).
		Model(&domain.Click{}).
		Select("EXTRACT(dow from clicked_at) as day_of_week, EXTRACT(hour from clicked_at) as hour, COUNT(*) as clicks").
		Where("short_url_id = ? AND clicked_at BETWEEN ? AND ?", shortURLID, startDate, endDate).
		Group("EXTRACT(dow from clicked_at), EXTRACT(hour from clicked_at)").
		Scan(&heatmapData).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to get heatmap data: %w", err)
	}
	
	// Convert to structured format
	dayNames := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	
	for _, data := range heatmapData {
		dayName := dayNames[data.DayOfWeek]
		if heatmap.Data[dayName] == nil {
			heatmap.Data[dayName] = make(map[int]int64)
		}
		heatmap.Data[dayName][data.Hour] = data.Clicks
	}
	
	return heatmap, nil
}

// GetGlobalStatsOptimized retrieves global statistics efficiently
func (r *OptimizedClickRepository) GetGlobalStatsOptimized(ctx context.Context, startDate, endDate time.Time) (*domain.GlobalStats, error) {
	stats := &domain.GlobalStats{}
	
	return stats, r.transaction.ReadOnlyTransaction(ctx, func(tx *gorm.DB) error {
		// Total clicks in date range
		err := tx.Model(&domain.Click{}).
			Where("clicked_at BETWEEN ? AND ?", startDate, endDate).
			Count(&stats.TotalClicks).Error
		if err != nil {
			return fmt.Errorf("failed to count total clicks: %w", err)
		}
		
		// Unique visitors
		err = tx.Model(&domain.Click{}).
			Select("COUNT(DISTINCT ip_address)").
			Where("clicked_at BETWEEN ? AND ?", startDate, endDate).
			Scan(&stats.UniqueVisitors).Error
		if err != nil {
			return fmt.Errorf("failed to count unique visitors: %w", err)
		}
		
		// Active URLs count
		err = tx.Model(&domain.ShortURL{}).
			Where("is_active = ?", true).
			Count(&stats.ActiveURLs).Error
		if err != nil {
			return fmt.Errorf("failed to count active URLs: %w", err)
		}
		
		// Total URLs count
		err = tx.Model(&domain.ShortURL{}).
			Count(&stats.TotalURLs).Error
		if err != nil {
			return fmt.Errorf("failed to count total URLs: %w", err)
		}
		
		return nil
	})
}

// CleanupOldClicksOptimized removes old click records efficiently
func (r *OptimizedClickRepository) CleanupOldClicksOptimized(ctx context.Context, olderThan time.Time, batchSize int) (int64, error) {
	result := r.db.WithContext(ctx).
		Where("clicked_at < ?", olderThan).
		Limit(batchSize).
		Delete(&domain.Click{})
	
	if result.Error != nil {
		return 0, fmt.Errorf("failed to cleanup old clicks: %w", result.Error)
	}
	
	return result.RowsAffected, nil
}

// GetTopPerformingURLsOptimized retrieves top performing URLs efficiently
func (r *OptimizedClickRepository) GetTopPerformingURLsOptimized(ctx context.Context, startDate, endDate time.Time, limit int) ([]*domain.URLPerformance, error) {
	var performance []*domain.URLPerformance
	
	err := r.db.WithContext(ctx).
		Table("clicks").
		Select("short_urls.short_code, short_urls.original_url, short_urls.title, COUNT(clicks.id) as click_count, COUNT(DISTINCT clicks.ip_address) as unique_visitors").
		Joins("JOIN short_urls ON clicks.short_url_id = short_urls.id").
		Where("clicks.clicked_at BETWEEN ? AND ? AND short_urls.is_active = ?", startDate, endDate, true).
		Group("short_urls.id, short_urls.short_code, short_urls.original_url, short_urls.title").
		Order("click_count DESC").
		Limit(limit).
		Find(&performance).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to get top performing URLs: %w", err)
	}
	
	return performance, nil
}

// GetRealtimeStatsOptimized retrieves real-time statistics for the last few minutes
func (r *OptimizedClickRepository) GetRealtimeStatsOptimized(ctx context.Context, minutes int) (*domain.RealtimeStats, error) {
	stats := &domain.RealtimeStats{
		Since: time.Now().Add(-time.Duration(minutes) * time.Minute),
	}
	
	return stats, r.transaction.ReadOnlyTransaction(ctx, func(tx *gorm.DB) error {
		// Recent clicks count
		err := tx.Model(&domain.Click{}).
			Where("clicked_at >= ?", stats.Since).
			Count(&stats.RecentClicks).Error
		if err != nil {
			return fmt.Errorf("failed to count recent clicks: %w", err)
		}
		
		// Active sessions (unique IPs in the last period)
		err = tx.Model(&domain.Click{}).
			Select("COUNT(DISTINCT ip_address)").
			Where("clicked_at >= ?", stats.Since).
			Scan(&stats.ActiveSessions).Error
		if err != nil {
			return fmt.Errorf("failed to count active sessions: %w", err)
		}
		
		// Top URLs in real-time
		err = tx.Table("clicks").
			Select("short_urls.short_code, COUNT(clicks.id) as click_count").
			Joins("JOIN short_urls ON clicks.short_url_id = short_urls.id").
			Where("clicks.clicked_at >= ?", stats.Since).
			Group("short_urls.short_code").
			Order("click_count DESC").
			Limit(10).
			Scan(&stats.TopURLs).Error
		if err != nil {
			return fmt.Errorf("failed to get top URLs: %w", err)
		}
		
		return nil
	})
}