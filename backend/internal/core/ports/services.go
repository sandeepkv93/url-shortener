package ports

import (
	"context"

	"url-shortener/internal/core/domain"
)

type AuthService interface {
	// Authentication
	Register(ctx context.Context, req *domain.CreateUserRequest) (*domain.UserResponse, error)
	Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenResponse, error)
	Logout(ctx context.Context, token string) error
	
	// Password management
	ResetPassword(ctx context.Context, email string) error
	ConfirmPasswordReset(ctx context.Context, token, newPassword string) error
	ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error
	
	// Token management
	ValidateToken(ctx context.Context, token string) (*domain.TokenClaims, error)
	GenerateTokens(ctx context.Context, user *domain.User) (*domain.TokenResponse, error)
	RevokeToken(ctx context.Context, token string) error
	
	// User management
	GetUserProfile(ctx context.Context, userID uint) (*domain.UserResponse, error)
	UpdateUserProfile(ctx context.Context, userID uint, req *domain.UpdateUserRequest) (*domain.UserResponse, error)
	DeleteAccount(ctx context.Context, userID uint) error
}

type URLService interface {
	// URL shortening
	ShortenURL(ctx context.Context, req *domain.CreateShortURLRequest, userID *uint) (*domain.ShortURLResponse, error)
	GetOriginalURL(ctx context.Context, shortCode string) (*domain.ShortURLResponse, error)
	
	// URL management
	GetURL(ctx context.Context, id uint, userID *uint) (*domain.ShortURLResponse, error)
	UpdateURL(ctx context.Context, id uint, req *domain.UpdateShortURLRequest, userID uint) (*domain.ShortURLResponse, error)
	DeleteURL(ctx context.Context, id uint, userID uint) error
	
	// URL queries
	ListUserURLs(ctx context.Context, userID uint, page, pageSize int, filter *domain.URLFilter) (*domain.ShortURLListResponse, error)
	SearchURLs(ctx context.Context, userID uint, query string, page, pageSize int) (*domain.ShortURLListResponse, error)
	
	// URL operations
	RedirectURL(ctx context.Context, shortCode, ipAddress, userAgent, referer string) (string, error)
	ToggleURLStatus(ctx context.Context, id uint, userID uint) (*domain.ShortURLResponse, error)
	
	// URL utilities
	GenerateShortCode(ctx context.Context, length int) (string, error)
	ValidateCustomAlias(ctx context.Context, alias string) error
	CheckURLAvailability(ctx context.Context, shortCode string) (bool, error)
	
	// Bulk operations
	BulkCreateURLs(ctx context.Context, urls []domain.CreateShortURLRequest, userID uint) ([]*domain.ShortURLResponse, error)
	BulkUpdateURLs(ctx context.Context, updates []domain.BulkUpdateRequest, userID uint) error
	BulkDeleteURLs(ctx context.Context, ids []uint, userID uint) error
}

type AnalyticsService interface {
	// Click recording
	RecordClick(ctx context.Context, req *domain.RecordClickRequest) error
	
	// Analytics queries
	GetURLAnalytics(ctx context.Context, shortURLID uint, userID uint, period string) (*domain.ClickStats, error)
	GetGeoAnalytics(ctx context.Context, shortURLID uint, userID uint) (*domain.GeoStats, error)
	GetTimelineAnalytics(ctx context.Context, shortURLID uint, userID uint, period string) (*domain.TimelineStats, error)
	
	// Dashboard analytics
	GetUserDashboard(ctx context.Context, userID uint) (*domain.UserDashboard, error)
	GetGlobalDashboard(ctx context.Context) (*domain.GlobalDashboard, error)
	
	// Report generation
	GenerateAnalyticsReport(ctx context.Context, shortURLID uint, userID uint, format string) (*domain.AnalyticsReport, error)
	ExportAnalytics(ctx context.Context, shortURLID uint, userID uint, format string) ([]byte, error)
	
	// Real-time analytics
	GetRealTimeStats(ctx context.Context, shortURLID uint, userID uint) (*domain.RealTimeStats, error)
	SubscribeToRealTimeUpdates(ctx context.Context, shortURLID uint, userID uint) (<-chan *domain.RealTimeUpdate, error)
}

type QRService interface {
	// QR code generation
	GenerateQRCode(ctx context.Context, shortURL string, options *domain.QROptions) (*domain.QRCodeResponse, error)
	GenerateQRCodeBytes(ctx context.Context, shortURL string, format string, size int) ([]byte, error)
	
	// QR code customization
	GenerateCustomQRCode(ctx context.Context, req *domain.CustomQRRequest) (*domain.QRCodeResponse, error)
	
	// QR code management
	GetQRCodeHistory(ctx context.Context, userID uint, page, pageSize int) (*domain.QRCodeListResponse, error)
	
	// Bulk QR generation
	BulkGenerateQRCodes(ctx context.Context, urls []string, options *domain.QROptions) ([]*domain.QRCodeResponse, error)
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