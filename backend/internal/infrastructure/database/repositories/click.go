package repositories

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"url-shortener/internal/core/domain"
	"url-shortener/internal/core/ports"
)

type clickRepository struct {
	db *gorm.DB
}

func NewClickRepository(db *gorm.DB) ports.ClickRepository {
	return &clickRepository{
		db: db,
	}
}

func (r *clickRepository) Create(ctx context.Context, click *domain.Click) error {
	if err := r.db.WithContext(ctx).Create(click).Error; err != nil {
		return fmt.Errorf("failed to create click: %w", err)
	}
	return nil
}

func (r *clickRepository) GetByID(ctx context.Context, id uint) (*domain.Click, error) {
	var click domain.Click
	if err := r.db.WithContext(ctx).
		Preload("ShortURL").
		First(&click, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("click not found")
		}
		return nil, fmt.Errorf("failed to get click by id: %w", err)
	}
	return &click, nil
}

func (r *clickRepository) GetByShortURLID(ctx context.Context, shortURLID uint, offset, limit int) ([]*domain.Click, int64, error) {
	var clicks []*domain.Click
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).
		Model(&domain.Click{}).
		Where("short_url_id = ?", shortURLID).
		Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count clicks: %w", err)
	}

	// Get clicks with pagination
	if err := r.db.WithContext(ctx).
		Where("short_url_id = ?", shortURLID).
		Offset(offset).
		Limit(limit).
		Order("clicked_at DESC").
		Find(&clicks).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list clicks: %w", err)
	}

	return clicks, total, nil
}

func (r *clickRepository) GetClickStats(ctx context.Context, shortURLID uint, period string) (*domain.ClickStats, error) {
	stats := &domain.ClickStats{
		ClicksByDate: make(map[string]int64),
		ClicksByTime: make(map[int]int64),
	}

	// Get total clicks
	if err := r.db.WithContext(ctx).
		Model(&domain.Click{}).
		Where("short_url_id = ?", shortURLID).
		Count(&stats.TotalClicks).Error; err != nil {
		return nil, fmt.Errorf("failed to count total clicks: %w", err)
	}

	// Get unique clicks (distinct IP addresses)
	if err := r.db.WithContext(ctx).
		Model(&domain.Click{}).
		Select("COUNT(DISTINCT ip_address)").
		Where("short_url_id = ?", shortURLID).
		Scan(&stats.UniqueClicks).Error; err != nil {
		return nil, fmt.Errorf("failed to count unique clicks: %w", err)
	}

	// Get clicks by date (last 30 days)
	var dateStats []struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}
	if err := r.db.WithContext(ctx).
		Model(&domain.Click{}).
		Select("DATE(clicked_at) as date, COUNT(*) as count").
		Where("short_url_id = ? AND clicked_at >= ?", shortURLID, time.Now().AddDate(0, 0, -30)).
		Group("DATE(clicked_at)").
		Order("date").
		Scan(&dateStats).Error; err != nil {
		return nil, fmt.Errorf("failed to get clicks by date: %w", err)
	}

	for _, stat := range dateStats {
		stats.ClicksByDate[stat.Date] = stat.Count
	}

	// Get clicks by hour
	var timeStats []struct {
		Hour  int   `json:"hour"`
		Count int64 `json:"count"`
	}
	if err := r.db.WithContext(ctx).
		Model(&domain.Click{}).
		Select("EXTRACT(HOUR FROM clicked_at) as hour, COUNT(*) as count").
		Where("short_url_id = ?", shortURLID).
		Group("EXTRACT(HOUR FROM clicked_at)").
		Order("hour").
		Scan(&timeStats).Error; err != nil {
		return nil, fmt.Errorf("failed to get clicks by time: %w", err)
	}

	for _, stat := range timeStats {
		stats.ClicksByTime[stat.Hour] = stat.Count
	}

	// Get top countries
	topCountries, err := r.GetTopCountries(ctx, shortURLID, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get top countries: %w", err)
	}
	stats.TopCountries = topCountries

	// Get top devices
	topDevices, err := r.GetTopDevices(ctx, shortURLID, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get top devices: %w", err)
	}
	stats.TopDevices = topDevices

	// Get top browsers
	topBrowsers, err := r.GetTopBrowsers(ctx, shortURLID, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get top browsers: %w", err)
	}
	stats.TopBrowsers = topBrowsers

	// Get top referers
	topReferers, err := r.GetTopReferers(ctx, shortURLID, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get top referers: %w", err)
	}
	stats.TopReferers = topReferers

	// Get recent clicks
	recentClicks, err := r.GetRecentClicks(ctx, shortURLID, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent clicks: %w", err)
	}
	stats.RecentClicks = recentClicks

	return stats, nil
}

func (r *clickRepository) GetGeoStats(ctx context.Context, shortURLID uint) (*domain.GeoStats, error) {
	stats := &domain.GeoStats{
		CountryStats: make(map[string]int64),
		RegionStats:  make(map[string]int64),
		CityStats:    make(map[string]int64),
	}

	// Get country stats
	var countryStats []struct {
		Country string `json:"country"`
		Count   int64  `json:"count"`
	}
	if err := r.db.WithContext(ctx).
		Model(&domain.Click{}).
		Select("country, COUNT(*) as count").
		Where("short_url_id = ? AND country != ''", shortURLID).
		Group("country").
		Order("count DESC").
		Scan(&countryStats).Error; err != nil {
		return nil, fmt.Errorf("failed to get country stats: %w", err)
	}

	for _, stat := range countryStats {
		stats.CountryStats[stat.Country] = stat.Count
	}

	// Get region stats
	var regionStats []struct {
		Region string `json:"region"`
		Count  int64  `json:"count"`
	}
	if err := r.db.WithContext(ctx).
		Model(&domain.Click{}).
		Select("region, COUNT(*) as count").
		Where("short_url_id = ? AND region != ''", shortURLID).
		Group("region").
		Order("count DESC").
		Limit(20).
		Scan(&regionStats).Error; err != nil {
		return nil, fmt.Errorf("failed to get region stats: %w", err)
	}

	for _, stat := range regionStats {
		stats.RegionStats[stat.Region] = stat.Count
	}

	// Get city stats
	var cityStats []struct {
		City  string `json:"city"`
		Count int64  `json:"count"`
	}
	if err := r.db.WithContext(ctx).
		Model(&domain.Click{}).
		Select("city, COUNT(*) as count").
		Where("short_url_id = ? AND city != ''", shortURLID).
		Group("city").
		Order("count DESC").
		Limit(20).
		Scan(&cityStats).Error; err != nil {
		return nil, fmt.Errorf("failed to get city stats: %w", err)
	}

	for _, stat := range cityStats {
		stats.CityStats[stat.City] = stat.Count
	}

	return stats, nil
}

func (r *clickRepository) GetTimelineStats(ctx context.Context, shortURLID uint, period string) (*domain.TimelineStats, error) {
	stats := &domain.TimelineStats{
		Period: period,
		Data:   make(map[string]int64),
	}

	var dateFormat string
	var startDate time.Time

	switch period {
	case "day":
		dateFormat = "DATE_FORMAT(clicked_at, '%Y-%m-%d %H:00:00')"
		startDate = time.Now().AddDate(0, 0, -1)
	case "week":
		dateFormat = "DATE(clicked_at)"
		startDate = time.Now().AddDate(0, 0, -7)
	case "month":
		dateFormat = "DATE(clicked_at)"
		startDate = time.Now().AddDate(0, -1, 0)
	case "year":
		dateFormat = "DATE_FORMAT(clicked_at, '%Y-%m')"
		startDate = time.Now().AddDate(-1, 0, 0)
	default:
		dateFormat = "DATE(clicked_at)"
		startDate = time.Now().AddDate(0, 0, -30)
	}

	var timelineData []struct {
		Period string `json:"period"`
		Count  int64  `json:"count"`
	}

	query := fmt.Sprintf(`
		SELECT %s as period, COUNT(*) as count 
		FROM clicks 
		WHERE short_url_id = ? AND clicked_at >= ? 
		GROUP BY %s 
		ORDER BY period`,
		dateFormat, dateFormat)

	if err := r.db.WithContext(ctx).
		Raw(query, shortURLID, startDate).
		Scan(&timelineData).Error; err != nil {
		return nil, fmt.Errorf("failed to get timeline stats: %w", err)
	}

	for _, data := range timelineData {
		stats.Data[data.Period] = data.Count
	}

	return stats, nil
}

func (r *clickRepository) GetTotalClicks(ctx context.Context, shortURLID uint) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&domain.Click{}).
		Where("short_url_id = ?", shortURLID).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count total clicks: %w", err)
	}
	return count, nil
}

func (r *clickRepository) GetUniqueClicks(ctx context.Context, shortURLID uint) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&domain.Click{}).
		Select("COUNT(DISTINCT ip_address)").
		Where("short_url_id = ?", shortURLID).
		Scan(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count unique clicks: %w", err)
	}
	return count, nil
}

func (r *clickRepository) GetClicksByDateRange(ctx context.Context, shortURLID uint, startDate, endDate string) ([]*domain.Click, error) {
	var clicks []*domain.Click
	if err := r.db.WithContext(ctx).
		Where("short_url_id = ? AND DATE(clicked_at) BETWEEN ? AND ?", shortURLID, startDate, endDate).
		Order("clicked_at DESC").
		Find(&clicks).Error; err != nil {
		return nil, fmt.Errorf("failed to get clicks by date range: %w", err)
	}
	return clicks, nil
}

func (r *clickRepository) GetTopCountries(ctx context.Context, shortURLID uint, limit int) ([]domain.CountryStat, error) {
	var stats []domain.CountryStat
	if err := r.db.WithContext(ctx).
		Model(&domain.Click{}).
		Select("country, COUNT(*) as count").
		Where("short_url_id = ? AND country != ''", shortURLID).
		Group("country").
		Order("count DESC").
		Limit(limit).
		Scan(&stats).Error; err != nil {
		return nil, fmt.Errorf("failed to get top countries: %w", err)
	}
	return stats, nil
}

func (r *clickRepository) GetTopDevices(ctx context.Context, shortURLID uint, limit int) ([]domain.DeviceStat, error) {
	var stats []domain.DeviceStat
	if err := r.db.WithContext(ctx).
		Model(&domain.Click{}).
		Select("device, COUNT(*) as count").
		Where("short_url_id = ? AND device != ''", shortURLID).
		Group("device").
		Order("count DESC").
		Limit(limit).
		Scan(&stats).Error; err != nil {
		return nil, fmt.Errorf("failed to get top devices: %w", err)
	}
	return stats, nil
}

func (r *clickRepository) GetTopBrowsers(ctx context.Context, shortURLID uint, limit int) ([]domain.BrowserStat, error) {
	var stats []domain.BrowserStat
	if err := r.db.WithContext(ctx).
		Model(&domain.Click{}).
		Select("browser, COUNT(*) as count").
		Where("short_url_id = ? AND browser != ''", shortURLID).
		Group("browser").
		Order("count DESC").
		Limit(limit).
		Scan(&stats).Error; err != nil {
		return nil, fmt.Errorf("failed to get top browsers: %w", err)
	}
	return stats, nil
}

func (r *clickRepository) GetTopReferers(ctx context.Context, shortURLID uint, limit int) ([]domain.RefererStat, error) {
	var stats []domain.RefererStat
	if err := r.db.WithContext(ctx).
		Model(&domain.Click{}).
		Select("referer, COUNT(*) as count").
		Where("short_url_id = ? AND referer != ''", shortURLID).
		Group("referer").
		Order("count DESC").
		Limit(limit).
		Scan(&stats).Error; err != nil {
		return nil, fmt.Errorf("failed to get top referers: %w", err)
	}
	return stats, nil
}

func (r *clickRepository) GetRecentClicks(ctx context.Context, shortURLID uint, limit int) ([]domain.RecentClickStat, error) {
	var stats []domain.RecentClickStat
	if err := r.db.WithContext(ctx).
		Model(&domain.Click{}).
		Select("country, city, device, browser, referer, clicked_at").
		Where("short_url_id = ?", shortURLID).
		Order("clicked_at DESC").
		Limit(limit).
		Scan(&stats).Error; err != nil {
		return nil, fmt.Errorf("failed to get recent clicks: %w", err)
	}
	return stats, nil
}

func (r *clickRepository) GetGlobalStats(ctx context.Context) (*domain.GlobalStats, error) {
	stats := &domain.GlobalStats{}

	// Get total URLs
	if err := r.db.WithContext(ctx).
		Model(&domain.ShortURL{}).
		Count(&stats.TotalURLs).Error; err != nil {
		return nil, fmt.Errorf("failed to count total URLs: %w", err)
	}

	// Get total clicks
	if err := r.db.WithContext(ctx).
		Model(&domain.Click{}).
		Count(&stats.TotalClicks).Error; err != nil {
		return nil, fmt.Errorf("failed to count total clicks: %w", err)
	}

	// Get total users
	if err := r.db.WithContext(ctx).
		Model(&domain.User{}).
		Count(&stats.TotalUsers).Error; err != nil {
		return nil, fmt.Errorf("failed to count total users: %w", err)
	}

	// Get active URLs
	now := time.Now()
	if err := r.db.WithContext(ctx).
		Model(&domain.ShortURL{}).
		Where("is_active = ? AND (expires_at IS NULL OR expires_at > ?)", true, now).
		Count(&stats.ActiveURLs).Error; err != nil {
		return nil, fmt.Errorf("failed to count active URLs: %w", err)
	}

	// Get today's clicks
	today := time.Now().Format("2006-01-02")
	if err := r.db.WithContext(ctx).
		Model(&domain.Click{}).
		Where("DATE(clicked_at) = ?", today).
		Count(&stats.ClicksToday).Error; err != nil {
		return nil, fmt.Errorf("failed to count today's clicks: %w", err)
	}

	// Get URLs created today
	if err := r.db.WithContext(ctx).
		Model(&domain.ShortURL{}).
		Where("DATE(created_at) = ?", today).
		Count(&stats.URLsCreatedToday).Error; err != nil {
		return nil, fmt.Errorf("failed to count URLs created today: %w", err)
	}

	// Get new users today
	if err := r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("DATE(created_at) = ?", today).
		Count(&stats.NewUsersToday).Error; err != nil {
		return nil, fmt.Errorf("failed to count new users today: %w", err)
	}

	return stats, nil
}

func (r *clickRepository) GetUserStats(ctx context.Context, userID uint) (*domain.UserAnalytics, error) {
	analytics := &domain.UserAnalytics{
		UserID:       userID,
		ClicksByDate: make(map[string]int64),
	}

	// Get total URLs
	if err := r.db.WithContext(ctx).
		Model(&domain.ShortURL{}).
		Where("user_id = ?", userID).
		Count(&analytics.TotalURLs).Error; err != nil {
		return nil, fmt.Errorf("failed to count user URLs: %w", err)
	}

	// Get total clicks
	if err := r.db.WithContext(ctx).
		Table("clicks").
		Joins("JOIN short_urls ON clicks.short_url_id = short_urls.id").
		Where("short_urls.user_id = ?", userID).
		Count(&analytics.TotalClicks).Error; err != nil {
		return nil, fmt.Errorf("failed to count user clicks: %w", err)
	}

	// Get clicks by date (last 30 days)
	var dateStats []struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}
	if err := r.db.WithContext(ctx).
		Table("clicks").
		Select("DATE(clicks.clicked_at) as date, COUNT(*) as count").
		Joins("JOIN short_urls ON clicks.short_url_id = short_urls.id").
		Where("short_urls.user_id = ? AND clicks.clicked_at >= ?", userID, time.Now().AddDate(0, 0, -30)).
		Group("DATE(clicks.clicked_at)").
		Order("date").
		Scan(&dateStats).Error; err != nil {
		return nil, fmt.Errorf("failed to get clicks by date: %w", err)
	}

	for _, stat := range dateStats {
		analytics.ClicksByDate[stat.Date] = stat.Count
	}

	// Get top URLs
	if err := r.db.WithContext(ctx).
		Model(&domain.ShortURL{}).
		Select("short_code, original_url, click_count").
		Where("user_id = ?", userID).
		Order("click_count DESC").
		Limit(10).
		Scan(&analytics.TopURLs).Error; err != nil {
		return nil, fmt.Errorf("failed to get top URLs: %w", err)
	}

	return analytics, nil
}