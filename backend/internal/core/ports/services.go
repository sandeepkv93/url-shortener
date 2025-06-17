package ports

import (
	"context"

	"url-shortener/internal/core/domain"
)

type AuthService interface {
	// Authentication
	Register(ctx context.Context, req domain.RegisterRequest) (*domain.AuthResponse, error)
	Login(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*domain.AuthResponse, error)
	Logout(ctx context.Context, userID uint) error
	
	// Token management
	ValidateToken(ctx context.Context, token string) (*domain.TokenClaims, error)
	
	// User management
	GetProfile(ctx context.Context, userID uint) (*domain.UserResponse, error)
	UpdateProfile(ctx context.Context, userID uint, req domain.UpdateProfileRequest) (*domain.UserResponse, error)
	ChangePassword(ctx context.Context, userID uint, req domain.ChangePasswordRequest) error
}

type URLService interface {
	// URL shortening
	ShortenURL(ctx context.Context, req domain.ShortenURLRequest) (*domain.ShortURL, error)
	GetOriginalURL(ctx context.Context, shortCode string) (*domain.ShortURL, error)
	
	// URL management
	GetUserURLs(ctx context.Context, userID uint, offset, limit int) ([]*domain.ShortURL, int64, error)
	UpdateURL(ctx context.Context, id uint, userID uint, req domain.UpdateURLRequest) (*domain.ShortURL, error)
	DeleteURL(ctx context.Context, id uint, userID uint) error
	
	// URL operations
	RecordClick(ctx context.Context, shortURL *domain.ShortURL, clickData domain.ClickData) error
	ValidatePassword(ctx context.Context, shortCode, password string) (bool, error)
	
	// URL utilities
	GetURLStats(ctx context.Context, id uint, userID uint) (*domain.URLStats, error)
	GetPopularURLs(ctx context.Context, limit int) ([]*domain.ShortURL, error)
	CleanupExpiredURLs(ctx context.Context) error
}

type AnalyticsService interface {
	// Dashboard analytics
	GetDashboardStats(ctx context.Context, userID uint) (*domain.DashboardStats, error)
	GetGlobalStats(ctx context.Context) (*domain.GlobalStats, error)
	
	// URL analytics
	GetURLAnalytics(ctx context.Context, shortURLID uint, userID uint) (*domain.URLAnalytics, error)
	GetTopPerformingURLs(ctx context.Context, userID uint, limit int) ([]*domain.URLPerformance, error)
	
	// Detailed analytics
	GetClickTimeline(ctx context.Context, shortURLID uint, userID uint, period string) (*domain.TimelineStats, error)
	GetGeographicStats(ctx context.Context, shortURLID uint, userID uint) (*domain.GeoStats, error)
	GetDeviceStats(ctx context.Context, shortURLID uint, userID uint) (*domain.DeviceStats, error)
	GetReferrerStats(ctx context.Context, shortURLID uint, userID uint) ([]domain.RefererStat, error)
	
	// Export functionality
	ExportAnalytics(ctx context.Context, userID uint, format string, dateRange domain.DateRange) ([]byte, error)
}

type QRService interface {
	// QR code generation
	GenerateQRCode(ctx context.Context, req domain.QRCodeRequest) (*domain.QRCodeResponse, error)
	GenerateQRCodeForURL(ctx context.Context, shortCode string, options domain.QRCodeOptions) (*domain.QRCodeResponse, error)
	
	// QR code utilities
	GetQRCodeFormats(ctx context.Context) []string
	GetQRCodeSizes(ctx context.Context) []int
	ValidateQRCodeOptions(ctx context.Context, options domain.QRCodeOptions) error
}

type NotificationService interface {
	// Email notifications
	SendWelcomeEmail(ctx context.Context, user *domain.User) error
	SendPasswordResetEmail(ctx context.Context, user *domain.User, resetToken string) error
	SendPasswordChangedNotification(ctx context.Context, user *domain.User) error
	
	// Analytics notifications
	SendAnalyticsDigest(ctx context.Context, user *domain.User, digest *domain.AnalyticsDigest) error
	SendClickAlert(ctx context.Context, user *domain.User, alert *domain.ClickAlert) error
	
	// System notifications
	SendMaintenanceNotification(ctx context.Context, users []*domain.User, message string) error
	SendSecurityAlert(ctx context.Context, user *domain.User, alert *domain.SecurityAlert) error
}

type GeolocationService interface {
	// IP geolocation
	GetLocationFromIP(ctx context.Context, ipAddress string) (*domain.GeoLocation, error)
	
	// Batch geolocation
	GetLocationsBatch(ctx context.Context, ipAddresses []string) (map[string]*domain.GeoLocation, error)
	
	// Location utilities
	ValidateLocation(ctx context.Context, location *domain.GeoLocation) error
	GetCountryCode(ctx context.Context, countryName string) (string, error)
}

type HealthService interface {
	// Health checks
	CheckHealth(ctx context.Context) (*domain.HealthStatus, error)
	CheckDatabaseHealth(ctx context.Context) (*domain.ComponentHealth, error)
	CheckCacheHealth(ctx context.Context) (*domain.ComponentHealth, error)
	CheckExternalServices(ctx context.Context) (map[string]*domain.ComponentHealth, error)
	
	// System metrics
	GetSystemMetrics(ctx context.Context) (*domain.SystemMetrics, error)
	GetApplicationMetrics(ctx context.Context) (*domain.ApplicationMetrics, error)
}

// Additional service interfaces needed by the service implementations

type JWTService interface {
	GenerateAccessToken(userID uint, email string) (string, error)
	GenerateRefreshToken(userID uint) (string, error)
	ValidateAccessToken(token string) (*domain.TokenClaims, error)
	ValidateRefreshToken(token string) (*domain.TokenClaims, error)
}

type ConfigService interface {
	GetBaseURL() string
	GetJWTSecret() string
	GetDatabaseURL() string
	GetRedisURL() string
}

type QRCodeProvider interface {
	GenerateQRCode(url string, options domain.QRGenerationOptions) ([]byte, error)
}