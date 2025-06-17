package ports

import (
	"context"

	"url-shortener/internal/core/domain"
)

type UserRepository interface {
	// User management
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uint) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uint) error
	
	// User queries
	Exists(ctx context.Context, email string) (bool, error)
	List(ctx context.Context, offset, limit int) ([]*domain.User, int64, error)
	
	// User statistics
	GetUserStats(ctx context.Context, userID uint) (*domain.UserStats, error)
}

type URLRepository interface {
	// URL management
	Create(ctx context.Context, url *domain.ShortURL) error
	GetByID(ctx context.Context, id uint) (*domain.ShortURL, error)
	GetByShortCode(ctx context.Context, shortCode string) (*domain.ShortURL, error)
	Update(ctx context.Context, url *domain.ShortURL) error
	Delete(ctx context.Context, id uint) error
	
	// URL queries
	ExistsByShortCode(ctx context.Context, shortCode string) (bool, error)
	GetByUserID(ctx context.Context, userID uint, offset, limit int) ([]*domain.ShortURL, int64, error)
	GetActiveByShortCode(ctx context.Context, shortCode string) (*domain.ShortURL, error)
	
	// URL operations
	IncrementClickCount(ctx context.Context, id uint) error
	GetExpiredURLs(ctx context.Context, limit int) ([]*domain.ShortURL, error)
	
	// URL statistics
	GetTotalURLs(ctx context.Context) (int64, error)
	GetTotalURLsByUser(ctx context.Context, userID uint) (int64, error)
	GetPopularURLs(ctx context.Context, limit int) ([]*domain.ShortURL, error)
}

type ClickRepository interface {
	// Click management
	Create(ctx context.Context, click *domain.Click) error
	GetByID(ctx context.Context, id uint) (*domain.Click, error)
	
	// Click queries
	GetByShortURLID(ctx context.Context, shortURLID uint, offset, limit int) ([]*domain.Click, int64, error)
	GetClickStats(ctx context.Context, shortURLID uint, period string) (*domain.ClickStats, error)
	GetGeoStats(ctx context.Context, shortURLID uint) (*domain.GeoStats, error)
	GetTimelineStats(ctx context.Context, shortURLID uint, period string) (*domain.TimelineStats, error)
	
	// Analytics queries
	GetTotalClicks(ctx context.Context, shortURLID uint) (int64, error)
	GetUniqueClicks(ctx context.Context, shortURLID uint) (int64, error)
	GetClicksByDateRange(ctx context.Context, shortURLID uint, startDate, endDate string) ([]*domain.Click, error)
	GetTopCountries(ctx context.Context, shortURLID uint, limit int) ([]domain.CountryStat, error)
	GetTopDevices(ctx context.Context, shortURLID uint, limit int) ([]domain.DeviceStat, error)
	GetTopBrowsers(ctx context.Context, shortURLID uint, limit int) ([]domain.BrowserStat, error)
	GetTopReferers(ctx context.Context, shortURLID uint, limit int) ([]domain.RefererStat, error)
	GetRecentClicks(ctx context.Context, shortURLID uint, limit int) ([]domain.RecentClickStat, error)
	
	// Global analytics
	GetGlobalStats(ctx context.Context) (*domain.GlobalStats, error)
	GetUserStats(ctx context.Context, userID uint) (*domain.UserAnalytics, error)
}