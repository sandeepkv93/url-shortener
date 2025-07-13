package ports

import (
	"context"
	"time"

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
	
	// Report notifications
	SendScheduledReport(ctx context.Context, recipients []string, report *domain.ScheduledReport, data []byte) error
	SendReportGenerationNotification(ctx context.Context, user *domain.User, report *domain.ScheduledReport, success bool, errorMsg string) error
	SendDataExportNotification(ctx context.Context, user *domain.User, export *domain.DataExport) error
	SendReportFailureAlert(ctx context.Context, user *domain.User, report *domain.ScheduledReport, error string) error
	
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
	GetHealth(ctx context.Context) (*domain.HealthStatus, error)
	CheckHealth(ctx context.Context) (*domain.HealthStatus, error)
	CheckDatabaseHealth(ctx context.Context) (*domain.ComponentHealth, error)
	CheckCacheHealth(ctx context.Context) (*domain.ComponentHealth, error)
	CheckExternalServices(ctx context.Context) (map[string]*domain.ComponentHealth, error)
	RunHealthChecks(ctx context.Context) (map[string]*domain.HealthCheck, error)
	IsHealthy(ctx context.Context) bool
	
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
	GetTokenTTL(tokenType string) time.Duration
	ExtractTokenFromHeader(authHeader string) (string, error)
	GetUserIDFromToken(tokenString string) (uint, error)
	IsTokenExpired(tokenString string) bool
	GetTokenClaims(tokenString string) (*domain.TokenClaims, error)
	RevokeToken(tokenString string) error
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

// Business Intelligence Services

type BusinessIntelligenceService interface {
	// Dashboard management
	CreateDashboard(ctx context.Context, userID uint, req domain.CreateDashboardRequest) (*domain.Dashboard, error)
	GetDashboard(ctx context.Context, dashboardID uint, userID uint) (*domain.Dashboard, error)
	GetUserDashboards(ctx context.Context, userID uint) ([]domain.DashboardResponse, error)
	UpdateDashboard(ctx context.Context, dashboardID uint, userID uint, req domain.UpdateDashboardRequest) (*domain.Dashboard, error)
	DeleteDashboard(ctx context.Context, dashboardID uint, userID uint) error
	
	// Widget management
	CreateWidget(ctx context.Context, userID uint, req domain.CreateWidgetRequest) (*domain.DashboardWidget, error)
	GetWidget(ctx context.Context, widgetID uint, userID uint) (*domain.DashboardWidget, error)
	GetWidgetData(ctx context.Context, widgetID uint, userID uint) (*domain.WidgetDataResponse, error)
	UpdateWidget(ctx context.Context, widgetID uint, userID uint, req domain.UpdateWidgetRequest) (*domain.DashboardWidget, error)
	DeleteWidget(ctx context.Context, widgetID uint, userID uint) error
	
	// Advanced analytics
	GetAdvancedAnalytics(ctx context.Context, userID uint, period string) (*domain.AdvancedAnalytics, error)
	GetPerformanceMetrics(ctx context.Context, userID uint, period string) (*domain.PerformanceMetrics, error)
	GetAudienceInsights(ctx context.Context, userID uint, period string) (*domain.AudienceInsights, error)
	GetContentAnalytics(ctx context.Context, userID uint, period string) (*domain.ContentAnalytics, error)
	
	// Competitive analysis
	GetCompetitiveAnalysis(ctx context.Context, userID uint) (*domain.CompetitiveAnalysis, error)
	GetMarketPosition(ctx context.Context, userID uint) (*domain.MarketPosition, error)
	GetBenchmarkData(ctx context.Context, userID uint, metric string) (*domain.BenchmarkData, error)
	
	// Predictive insights
	GetPredictiveInsights(ctx context.Context, userID uint) (*domain.PredictiveInsights, error)
	GetForecastData(ctx context.Context, userID uint, metric string, period string) (*domain.ForecastData, error)
	DetectAnomalies(ctx context.Context, userID uint, metric string) ([]domain.Anomaly, error)
	GetTrendPrediction(ctx context.Context, userID uint, metric string) (*domain.TrendPrediction, error)
	
	// Recommendations
	GetRecommendations(ctx context.Context, userID uint) (*domain.RecommendationEngine, error)
	GetOptimizationSuggestions(ctx context.Context, userID uint) ([]domain.OptimizationSuggestion, error)
	GetContentRecommendations(ctx context.Context, userID uint) ([]domain.ContentRec, error)
	GetAudienceRecommendations(ctx context.Context, userID uint) ([]domain.AudienceRec, error)
}

type FunnelService interface {
	// Funnel management
	CreateFunnel(ctx context.Context, userID uint, req domain.CreateFunnelRequest) (*domain.ConversionFunnel, error)
	GetFunnel(ctx context.Context, funnelID uint, userID uint) (*domain.ConversionFunnel, error)
	GetUserFunnels(ctx context.Context, userID uint) ([]*domain.ConversionFunnel, error)
	UpdateFunnel(ctx context.Context, funnelID uint, userID uint, req domain.CreateFunnelRequest) (*domain.ConversionFunnel, error)
	DeleteFunnel(ctx context.Context, funnelID uint, userID uint) error
	
	// Funnel analytics
	GetFunnelAnalytics(ctx context.Context, funnelID uint, userID uint, period string) (*domain.FunnelAnalyticsResponse, error)
	GetFunnelStepAnalytics(ctx context.Context, funnelID uint, stepID uint, userID uint) (*domain.FunnelStepAnalytics, error)
	GetConversionTrend(ctx context.Context, funnelID uint, userID uint, period string) (map[string]int64, error)
	
	// Funnel optimization
	GetFunnelOptimizations(ctx context.Context, funnelID uint, userID uint) ([]domain.OptimizationSuggestion, error)
	AnalyzeFunnelDropOffs(ctx context.Context, funnelID uint, userID uint) ([]domain.FunnelStepAnalytics, error)
}

type ReportingService interface {
	// Scheduled reports
	CreateScheduledReport(ctx context.Context, userID uint, req domain.CreateScheduledReportRequest) (*domain.ScheduledReport, error)
	GetScheduledReport(ctx context.Context, reportID uint, userID uint) (*domain.ScheduledReport, error)
	GetUserScheduledReports(ctx context.Context, userID uint) ([]*domain.ScheduledReport, error)
	UpdateScheduledReport(ctx context.Context, reportID uint, userID uint, req domain.CreateScheduledReportRequest) (*domain.ScheduledReport, error)
	DeleteScheduledReport(ctx context.Context, reportID uint, userID uint) error
	
	// Report execution
	ExecuteReport(ctx context.Context, reportID uint) error
	GetReportHistory(ctx context.Context, reportID uint, userID uint) ([]domain.ReportExecution, error)
	
	// Data export
	CreateDataExport(ctx context.Context, userID uint, req domain.DataExportRequest) (*domain.DataExport, error)
	GetDataExport(ctx context.Context, exportID uint, userID uint) (*domain.DataExport, error)
	GetUserDataExports(ctx context.Context, userID uint) ([]*domain.DataExport, error)
	DownloadExport(ctx context.Context, exportID uint, userID uint) ([]byte, string, error)
	
	// Report generation
	GenerateAnalyticsReport(ctx context.Context, userID uint, config domain.ReportConfig) ([]byte, error)
	GenerateDashboardReport(ctx context.Context, dashboardID uint, userID uint, format string) ([]byte, error)
	GenerateFunnelReport(ctx context.Context, funnelID uint, userID uint, format string) ([]byte, error)
}

// Additional helper services for BI

type PredictiveAnalyticsService interface {
	// Machine learning predictions
	PredictClickTrend(ctx context.Context, userID uint, daysAhead int) (map[string]float64, error)
	PredictUserBehavior(ctx context.Context, userID uint) (*domain.BehaviorPatterns, error)
	PredictContentPerformance(ctx context.Context, contentID string) (*domain.ContentPerformance, error)
	
	// Anomaly detection
	DetectClickAnomalies(ctx context.Context, userID uint, sensitivity float64) ([]domain.Anomaly, error)
	DetectTrafficSpikes(ctx context.Context, userID uint) ([]domain.Anomaly, error)
	DetectBehaviorAnomalies(ctx context.Context, userID uint) ([]domain.Anomaly, error)
	
	// Risk assessment
	AssessPerformanceRisk(ctx context.Context, userID uint) (*domain.RiskAssessment, error)
	AssessSecurityRisk(ctx context.Context, userID uint) (*domain.RiskAssessment, error)
	
	// Optimization recommendations
	RecommendOptimizations(ctx context.Context, userID uint) ([]domain.OptimizationSuggestion, error)
	RecommendContentStrategy(ctx context.Context, userID uint) ([]domain.ContentRec, error)
	RecommendAudienceTargeting(ctx context.Context, userID uint) ([]domain.AudienceRec, error)
}

type CompetitiveIntelligenceService interface {
	// Market analysis
	AnalyzeMarketPosition(ctx context.Context, userID uint) (*domain.MarketPosition, error)
	GetCompetitorData(ctx context.Context, competitorID string) (*domain.CompetitorData, error)
	GetMarketTrends(ctx context.Context, industry string) (*domain.MarketTrends, error)
	
	// Benchmarking
	GetIndustryBenchmarks(ctx context.Context, industry string) (*domain.BenchmarkData, error)
	ComparePerformance(ctx context.Context, userID uint, competitorID string) (map[string]float64, error)
	GetPerformanceGaps(ctx context.Context, userID uint) ([]domain.OpportunityGap, error)
	
	// Opportunity identification
	IdentifyMarketOpportunities(ctx context.Context, userID uint) ([]domain.OpportunityGap, error)
	AnalyzeCompetitorWeaknesses(ctx context.Context, competitorID string) ([]string, error)
	GetEmergingTrends(ctx context.Context, industry string) ([]domain.Trend, error)
}

type SchedulerService interface {
	// Report scheduling
	ScheduleReport(ctx context.Context, report *domain.ScheduledReport) error
	UnscheduleReport(ctx context.Context, reportID uint) error
	GetDueReports(ctx context.Context) ([]*domain.ScheduledReport, error)
	
	// Scheduler management
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	IsRunning() bool
	
	// Manual execution
	ExecuteReportNow(ctx context.Context, reportID uint) error
	GetJobStatus(ctx context.Context, jobID string) (*domain.JobStatus, error)
}