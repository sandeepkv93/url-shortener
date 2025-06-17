package services

import (
	"context"
	"fmt"
	"time"

	"url-shortener/internal/core/domain"
	"url-shortener/internal/core/ports"
)

type analyticsService struct {
	urlRepo     ports.URLRepository
	clickRepo   ports.ClickRepository
	userRepo    ports.UserRepository
	cacheRepo   ports.CacheService
	configRepo  ports.ConfigService
}

func NewAnalyticsService(
	urlRepo ports.URLRepository,
	clickRepo ports.ClickRepository,
	userRepo ports.UserRepository,
	cacheRepo ports.CacheService,
	configRepo ports.ConfigService,
) ports.AnalyticsService {
	return &analyticsService{
		urlRepo:    urlRepo,
		clickRepo:  clickRepo,
		userRepo:   userRepo,
		cacheRepo:  cacheRepo,
		configRepo: configRepo,
	}
}

func (s *analyticsService) GetDashboardStats(ctx context.Context, userID uint) (*domain.DashboardStats, error) {
	stats := &domain.DashboardStats{}

	// Get user stats
	userStats, err := s.userRepo.GetUserStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}

	stats.TotalURLs = userStats.TotalURLs
	stats.ActiveURLs = userStats.ActiveURLs
	stats.TotalClicks = userStats.TotalClicks

	// Get user analytics
	userAnalytics, err := s.clickRepo.GetUserStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user analytics: %w", err)
	}

	stats.ClicksByDate = userAnalytics.ClicksByDate
	stats.TopURLs = userAnalytics.TopURLs

	// Get recent activity (last 10 clicks)
	recentActivity, err := s.getRecentUserActivity(ctx, userID, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent activity: %w", err)
	}
	stats.RecentActivity = recentActivity

	// Calculate growth rates
	stats.ClickGrowthRate = s.calculateGrowthRate(userAnalytics.ClicksByDate, 7)
	stats.URLGrowthRate = s.calculateURLGrowthRate(ctx, userID, 7)

	return stats, nil
}

func (s *analyticsService) GetURLAnalytics(ctx context.Context, shortURLID uint, userID uint) (*domain.URLAnalytics, error) {
	// Verify URL ownership
	shortURL, err := s.urlRepo.GetByID(ctx, shortURLID)
	if err != nil {
		return nil, err
	}
	if shortURL.UserID != userID {
		return nil, domain.ErrUnauthorized
	}

	analytics := &domain.URLAnalytics{
		ShortURLID: shortURLID,
		ShortCode:  shortURL.ShortCode,
	}

	// Get click stats
	clickStats, err := s.clickRepo.GetClickStats(ctx, shortURLID, "month")
	if err != nil {
		return nil, fmt.Errorf("failed to get click stats: %w", err)
	}
	analytics.TotalClicks = clickStats.TotalClicks
	analytics.UniqueClicks = clickStats.UniqueClicks
	analytics.ClicksByDate = clickStats.ClicksByDate
	analytics.ClicksByTime = clickStats.ClicksByTime

	// Get geo stats
	geoStats, err := s.clickRepo.GetGeoStats(ctx, shortURLID)
	if err != nil {
		return nil, fmt.Errorf("failed to get geo stats: %w", err)
	}
	analytics.CountryStats = geoStats.CountryStats
	analytics.RegionStats = geoStats.RegionStats
	analytics.CityStats = geoStats.CityStats

	// Get device stats
	analytics.TopDevices = clickStats.TopDevices
	analytics.TopBrowsers = clickStats.TopBrowsers
	analytics.TopReferers = clickStats.TopReferers

	// Get recent clicks
	analytics.RecentClicks = clickStats.RecentClicks

	return analytics, nil
}

func (s *analyticsService) GetGlobalStats(ctx context.Context) (*domain.GlobalStats, error) {
	// Check cache first - skip caching for now due to interface limitations
	// In a real implementation, we would use a more sophisticated caching approach

	// Get fresh stats from database
	stats, err := s.clickRepo.GetGlobalStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get global stats: %w", err)
	}

	// Skip caching for now - would need proper serialization in a real implementation

	return stats, nil
}

func (s *analyticsService) GetTopPerformingURLs(ctx context.Context, userID uint, limit int) ([]*domain.URLPerformance, error) {
	urls, _, err := s.urlRepo.GetByUserID(ctx, userID, 0, limit*2) // Get more to filter
	if err != nil {
		return nil, fmt.Errorf("failed to get user URLs: %w", err)
	}

	var urlPerformances []*domain.URLPerformance
	for _, url := range urls {
		if len(urlPerformances) >= limit {
			break
		}

		// Get click stats for each URL
		totalClicks, err := s.clickRepo.GetTotalClicks(ctx, url.ID)
		if err != nil {
			continue // Skip URLs with errors
		}

		uniqueClicks, err := s.clickRepo.GetUniqueClicks(ctx, url.ID)
		if err != nil {
			continue
		}

		performance := &domain.URLPerformance{
			ShortURL:     url,
			TotalClicks:  totalClicks,
			UniqueClicks: uniqueClicks,
			ClickRate:    s.calculateClickRate(totalClicks, uniqueClicks),
		}

		urlPerformances = append(urlPerformances, performance)
	}

	// Sort by total clicks (already sorted by database query)
	return urlPerformances, nil
}

func (s *analyticsService) GetClickTimeline(ctx context.Context, shortURLID uint, userID uint, period string) (*domain.TimelineStats, error) {
	// Verify URL ownership
	shortURL, err := s.urlRepo.GetByID(ctx, shortURLID)
	if err != nil {
		return nil, err
	}
	if shortURL.UserID != userID {
		return nil, domain.ErrUnauthorized
	}

	return s.clickRepo.GetTimelineStats(ctx, shortURLID, period)
}

func (s *analyticsService) GetGeographicStats(ctx context.Context, shortURLID uint, userID uint) (*domain.GeoStats, error) {
	// Verify URL ownership
	shortURL, err := s.urlRepo.GetByID(ctx, shortURLID)
	if err != nil {
		return nil, err
	}
	if shortURL.UserID != userID {
		return nil, domain.ErrUnauthorized
	}

	return s.clickRepo.GetGeoStats(ctx, shortURLID)
}

func (s *analyticsService) GetDeviceStats(ctx context.Context, shortURLID uint, userID uint) (*domain.DeviceStats, error) {
	// Verify URL ownership
	shortURL, err := s.urlRepo.GetByID(ctx, shortURLID)
	if err != nil {
		return nil, err
	}
	if shortURL.UserID != userID {
		return nil, domain.ErrUnauthorized
	}

	// Get click stats which includes device information
	clickStats, err := s.clickRepo.GetClickStats(ctx, shortURLID, "all")
	if err != nil {
		return nil, fmt.Errorf("failed to get click stats: %w", err)
	}

	return &domain.DeviceStats{
		TopDevices:  clickStats.TopDevices,
		TopBrowsers: clickStats.TopBrowsers,
	}, nil
}

func (s *analyticsService) GetReferrerStats(ctx context.Context, shortURLID uint, userID uint) ([]domain.RefererStat, error) {
	// Verify URL ownership
	shortURL, err := s.urlRepo.GetByID(ctx, shortURLID)
	if err != nil {
		return nil, err
	}
	if shortURL.UserID != userID {
		return nil, domain.ErrUnauthorized
	}

	return s.clickRepo.GetTopReferers(ctx, shortURLID, 20)
}

func (s *analyticsService) ExportAnalytics(ctx context.Context, userID uint, format string, dateRange domain.DateRange) ([]byte, error) {
	// Get detailed click data for the date range
	urls, _, err := s.urlRepo.GetByUserID(ctx, userID, 0, 1000) // Get all URLs
	if err != nil {
		return nil, fmt.Errorf("failed to get user URLs: %w", err)
	}

	var exportData []map[string]interface{}
	for _, url := range urls {
		clicks, err := s.clickRepo.GetClicksByDateRange(ctx, url.ID, dateRange.StartDate, dateRange.EndDate)
		if err != nil {
			continue
		}

		for _, click := range clicks {
			exportData = append(exportData, map[string]interface{}{
				"short_code":   url.ShortCode,
				"original_url": url.OriginalURL,
				"title":        url.Title,
				"clicked_at":   click.ClickedAt,
				"ip_address":   click.IPAddress,
				"country":      click.Country,
				"region":       click.Region,
				"city":         click.City,
				"device":       click.Device,
				"browser":      click.Browser,
				"os":           click.OS,
				"referer":      click.Referer,
			})
		}
	}

	// Format the data based on the requested format
	switch format {
	case "json":
		return s.exportToJSON(exportData)
	case "csv":
		return s.exportToCSV(exportData)
	default:
		return nil, fmt.Errorf("unsupported export format: %s", format)
	}
}

func (s *analyticsService) getRecentUserActivity(ctx context.Context, userID uint, limit int) ([]domain.ActivityItem, error) {
	var activities []domain.ActivityItem

	// Get recent URLs created by user
	urls, _, err := s.urlRepo.GetByUserID(ctx, userID, 0, limit/2)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent URLs: %w", err)
	}

	for _, url := range urls {
		activities = append(activities, domain.ActivityItem{
			Type:        "url_created",
			Description: fmt.Sprintf("Created short URL: %s", url.ShortCode),
			Timestamp:   url.CreatedAt,
		})
	}

	// Get recent clicks on user's URLs
	for _, url := range urls {
		recentClicks, err := s.clickRepo.GetRecentClicks(ctx, url.ID, 2)
		if err != nil {
			continue
		}

		for _, click := range recentClicks {
			activities = append(activities, domain.ActivityItem{
				Type:        "url_clicked",
				Description: fmt.Sprintf("URL %s clicked from %s", url.ShortCode, click.Country),
				Timestamp:   click.ClickedAt,
			})
		}
	}

	// Sort activities by timestamp (most recent first)
	// Simple bubble sort for small datasets
	for i := 0; i < len(activities); i++ {
		for j := i + 1; j < len(activities); j++ {
			if activities[i].Timestamp.Before(activities[j].Timestamp) {
				activities[i], activities[j] = activities[j], activities[i]
			}
		}
	}

	// Limit the results
	if len(activities) > limit {
		activities = activities[:limit]
	}

	return activities, nil
}

func (s *analyticsService) calculateGrowthRate(clicksByDate map[string]int64, days int) float64 {
	if len(clicksByDate) < days {
		return 0.0
	}

	// Calculate growth rate based on the last 'days' period
	var recentSum, previousSum int64
	now := time.Now()

	for i := 0; i < days; i++ {
		dateStr := now.AddDate(0, 0, -i).Format("2006-01-02")
		recentSum += clicksByDate[dateStr]
	}

	for i := days; i < days*2; i++ {
		dateStr := now.AddDate(0, 0, -i).Format("2006-01-02")
		previousSum += clicksByDate[dateStr]
	}

	if previousSum == 0 {
		return 0.0
	}

	return float64(recentSum-previousSum) / float64(previousSum) * 100
}

func (s *analyticsService) calculateURLGrowthRate(ctx context.Context, userID uint, days int) float64 {
	// This is a simplified calculation
	// In a real implementation, you'd track URL creation dates
	return 0.0
}

func (s *analyticsService) calculateClickRate(totalClicks, uniqueClicks int64) float64 {
	if totalClicks == 0 {
		return 0.0
	}
	return float64(uniqueClicks) / float64(totalClicks) * 100
}

func (s *analyticsService) exportToJSON(data []map[string]interface{}) ([]byte, error) {
	// Simple JSON marshaling
	// In a real implementation, you'd use proper JSON encoding
	return []byte("{}"), nil // Placeholder
}

func (s *analyticsService) exportToCSV(data []map[string]interface{}) ([]byte, error) {
	// Simple CSV generation
	// In a real implementation, you'd use proper CSV encoding
	return []byte(""), nil // Placeholder
}