package services

import (
	"url-shortener/internal/config"
	"url-shortener/internal/core/ports"
)

type configService struct {
	config *config.Config
}

func NewConfigService(cfg *config.Config) ports.ConfigService {
	return &configService{
		config: cfg,
	}
}

func (s *configService) GetBaseURL() string {
	return s.config.Server.BaseURL
}

func (s *configService) GetJWTSecret() string {
	return s.config.JWT.Secret
}

func (s *configService) GetDatabaseURL() string {
	return s.config.Database.URL
}

func (s *configService) GetRedisURL() string {
	return s.config.Redis.URL
}

func (s *configService) GetServerPort() string {
	return s.config.Server.Port
}

func (s *configService) GetServerHost() string {
	return s.config.Server.Host
}

func (s *configService) IsProduction() bool {
	return s.config.Environment == "production"
}

func (s *configService) IsDevelopment() bool {
	return s.config.Environment == "development"
}

func (s *configService) GetEnvironment() string {
	return s.config.Environment
}

func (s *configService) GetLogLevel() string {
	return s.config.Logging.Level
}

func (s *configService) GetCORSOrigins() []string {
	return s.config.CORS.AllowedOrigins
}

func (s *configService) GetRateLimitRequests() int {
	return s.config.RateLimit.RequestsPerSecond
}

func (s *configService) GetRateLimitBurst() int {
	return s.config.RateLimit.BurstSize
}

func (s *configService) GetJWTAccessTokenTTL() int {
	return s.config.JWT.AccessTokenTTL
}

func (s *configService) GetJWTRefreshTokenTTL() int {
	return s.config.JWT.RefreshTokenTTL
}

func (s *configService) GetDatabaseMaxOpenConns() int {
	return s.config.Database.MaxOpenConns
}

func (s *configService) GetDatabaseMaxIdleConns() int {
	return s.config.Database.MaxIdleConns
}

func (s *configService) GetDatabaseConnMaxLifetime() int {
	return s.config.Database.ConnMaxLifetime
}

func (s *configService) GetRedisMaxRetries() int {
	return s.config.Redis.MaxRetries
}

func (s *configService) GetRedisRetryDelay() int {
	return s.config.Redis.RetryDelay
}

func (s *configService) GetRedisPoolSize() int {
	return s.config.Redis.PoolSize
}

func (s *configService) GetRedisTimeout() int {
	return s.config.Redis.Timeout
}

func (s *configService) GetExternalAPITimeout() int {
	return s.config.ExternalAPI.Timeout
}

func (s *configService) GetExternalAPIRetries() int {
	return s.config.ExternalAPI.MaxRetries
}

func (s *configService) GetFileUploadMaxSize() int64 {
	return s.config.FileUpload.MaxSize
}

func (s *configService) GetFileUploadAllowedTypes() []string {
	return s.config.FileUpload.AllowedTypes
}

func (s *configService) GetEmailSMTPHost() string {
	return s.config.Email.SMTPHost
}

func (s *configService) GetEmailSMTPPort() int {
	return s.config.Email.SMTPPort
}

func (s *configService) GetEmailSMTPUser() string {
	return s.config.Email.SMTPUser
}

func (s *configService) GetEmailSMTPPassword() string {
	return s.config.Email.SMTPPassword
}

func (s *configService) GetEmailFromAddress() string {
	return s.config.Email.FromAddress
}

func (s *configService) GetEmailFromName() string {
	return s.config.Email.FromName
}

func (s *configService) GetMonitoringEnabled() bool {
	return s.config.Monitoring.Enabled
}

func (s *configService) GetMonitoringPort() string {
	return s.config.Monitoring.Port
}

func (s *configService) GetTracingEnabled() bool {
	return s.config.Tracing.Enabled
}

func (s *configService) GetTracingEndpoint() string {
	return s.config.Tracing.Endpoint
}

func (s *configService) GetCacheDefaultTTL() int {
	return s.config.Cache.DefaultTTL
}

func (s *configService) GetCacheURLTTL() int {
	return s.config.Cache.URLTTL
}

func (s *configService) GetCacheSessionTTL() int {
	return s.config.Cache.SessionTTL
}

func (s *configService) GetCacheAnalyticsTTL() int {
	return s.config.Cache.AnalyticsTTL
}

func (s *configService) GetSecurityEnabled() bool {
	return s.config.Security.Enabled
}

func (s *configService) GetSecurityHSTSEnabled() bool {
	return s.config.Security.HSTSEnabled
}

func (s *configService) GetSecurityCSPEnabled() bool {
	return s.config.Security.CSPEnabled
}

func (s *configService) GetSecurityXSSProtectionEnabled() bool {
	return s.config.Security.XSSProtectionEnabled
}

func (s *configService) GetSecurityContentTypeNosniffEnabled() bool {
	return s.config.Security.ContentTypeNosniffEnabled
}

func (s *configService) GetSecurityFrameOptionsEnabled() bool {
	return s.config.Security.FrameOptionsEnabled
}

func (s *configService) GetBackupEnabled() bool {
	return s.config.Backup.Enabled
}

func (s *configService) GetBackupInterval() int {
	return s.config.Backup.Interval
}

func (s *configService) GetBackupRetention() int {
	return s.config.Backup.Retention
}

func (s *configService) GetBackupStoragePath() string {
	return s.config.Backup.StoragePath
}

func (s *configService) GetNotificationEnabled() bool {
	return s.config.Notification.Enabled
}

func (s *configService) GetNotificationWebhookURL() string {
	return s.config.Notification.WebhookURL
}

func (s *configService) GetFeatureQRCodeEnabled() bool {
	return s.config.Features.QRCodeEnabled
}

func (s *configService) GetFeatureAnalyticsEnabled() bool {
	return s.config.Features.AnalyticsEnabled
}

func (s *configService) GetFeatureCustomDomainsEnabled() bool {
	return s.config.Features.CustomDomainsEnabled
}

func (s *configService) GetFeaturePasswordProtectionEnabled() bool {
	return s.config.Features.PasswordProtectionEnabled
}

func (s *configService) GetFeatureExpirationEnabled() bool {
	return s.config.Features.ExpirationEnabled
}

func (s *configService) GetFeatureBulkOperationsEnabled() bool {
	return s.config.Features.BulkOperationsEnabled
}

func (s *configService) GetFeatureAPIEnabled() bool {
	return s.config.Features.APIEnabled
}

func (s *configService) GetFeatureWebhooksEnabled() bool {
	return s.config.Features.WebhooksEnabled
}

func (s *configService) GetFeatureEmailNotificationsEnabled() bool {
	return s.config.Features.EmailNotificationsEnabled
}

func (s *configService) GetFeatureGeoLocationEnabled() bool {
	return s.config.Features.GeoLocationEnabled
}

func (s *configService) GetFeatureDeviceTrackingEnabled() bool {
	return s.config.Features.DeviceTrackingEnabled
}

func (s *configService) GetFeatureUserRegistrationEnabled() bool {
	return s.config.Features.UserRegistrationEnabled
}

func (s *configService) GetFeatureGuestModeEnabled() bool {
	return s.config.Features.GuestModeEnabled
}

func (s *configService) GetFeatureRateLimitingEnabled() bool {
	return s.config.Features.RateLimitingEnabled
}

func (s *configService) GetFeatureSpamDetectionEnabled() bool {
	return s.config.Features.SpamDetectionEnabled
}

func (s *configService) GetFeatureAuditLogsEnabled() bool {
	return s.config.Features.AuditLogsEnabled
}

func (s *configService) GetRawConfig() *config.Config {
	return s.config
}

func (s *configService) Validate() error {
	return s.config.Validate()
}