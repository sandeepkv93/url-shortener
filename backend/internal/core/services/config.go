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
	return s.config.App.BaseURL
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
	return s.config.Server.Env == "production"
}

func (s *configService) IsDevelopment() bool {
	return s.config.Server.Env == "development"
}

func (s *configService) GetEnvironment() string {
	return s.config.Server.Env
}

func (s *configService) GetLogLevel() string {
	return s.config.Logging.Level
}

func (s *configService) GetCORSOrigins() []string {
	return s.config.CORS.AllowedOrigins
}

func (s *configService) GetRateLimitRequests() int {
	return s.config.Rate.Requests
}

func (s *configService) GetRateLimitBurst() int {
	return s.config.Rate.Requests // Use same field as burst for now
}

func (s *configService) GetJWTAccessTokenTTL() int {
	return int(s.config.JWT.Expiry.Seconds())
}

func (s *configService) GetJWTRefreshTokenTTL() int {
	return int(s.config.JWT.RefreshExpiry.Seconds())
}

func (s *configService) GetDatabaseMaxOpenConns() int {
	return s.config.Database.MaxConnections
}

func (s *configService) GetDatabaseMaxIdleConns() int {
	return s.config.Database.MaxIdle
}

func (s *configService) GetDatabaseConnMaxLifetime() int {
	return 3600 // Default 1 hour in seconds, field doesn't exist in config
}

func (s *configService) GetRedisMaxRetries() int {
	return 3 // Default value
}

func (s *configService) GetRedisRetryDelay() int {
	return 100 // Default value in milliseconds
}

func (s *configService) GetRedisPoolSize() int {
	return 10 // Default value
}

func (s *configService) GetRedisTimeout() int {
	return 5000 // Default value in milliseconds
}

func (s *configService) GetExternalAPITimeout() int {
	return 30 // Default value in seconds
}

func (s *configService) GetExternalAPIRetries() int {
	return 3 // Default value
}

func (s *configService) GetFileUploadMaxSize() int64 {
	return 10 * 1024 * 1024 // Default 10MB
}

func (s *configService) GetFileUploadAllowedTypes() []string {
	return []string{"image/jpeg", "image/png", "image/gif"} // Default allowed types
}

func (s *configService) GetEmailSMTPHost() string {
	return "localhost" // Default value
}

func (s *configService) GetEmailSMTPPort() int {
	return 587 // Default SMTP port
}

func (s *configService) GetEmailSMTPUser() string {
	return "" // Default empty
}

func (s *configService) GetEmailSMTPPassword() string {
	return "" // Default empty
}

func (s *configService) GetEmailFromAddress() string {
	return "noreply@localhost" // Default value
}

func (s *configService) GetEmailFromName() string {
	return "URL Shortener" // Default value
}

func (s *configService) GetMonitoringEnabled() bool {
	return false // Default disabled
}

func (s *configService) GetMonitoringPort() string {
	return "9090" // Default port
}

func (s *configService) GetTracingEnabled() bool {
	return false // Default disabled
}

func (s *configService) GetTracingEndpoint() string {
	return "http://localhost:14268/api/traces" // Default jaeger endpoint
}

func (s *configService) GetCacheDefaultTTL() int {
	return int(s.config.Cache.TTL.Seconds())
}

func (s *configService) GetCacheURLTTL() int {
	return int(s.config.Cache.URLTTL.Seconds())
}

func (s *configService) GetCacheSessionTTL() int {
	return 86400 // Default 24 hours in seconds
}

func (s *configService) GetCacheAnalyticsTTL() int {
	return 3600 // Default 1 hour in seconds
}

func (s *configService) GetSecurityEnabled() bool {
	return s.config.Security.EnableHTTPS
}

func (s *configService) GetSecurityHSTSEnabled() bool {
	return s.config.Security.EnableHTTPS // Use EnableHTTPS as proxy
}

func (s *configService) GetSecurityCSPEnabled() bool {
	return true // Default enabled
}

func (s *configService) GetSecurityXSSProtectionEnabled() bool {
	return true // Default enabled
}

func (s *configService) GetSecurityContentTypeNosniffEnabled() bool {
	return true // Default enabled
}

func (s *configService) GetSecurityFrameOptionsEnabled() bool {
	return true // Default enabled
}

func (s *configService) GetBackupEnabled() bool {
	return false // Default disabled
}

func (s *configService) GetBackupInterval() int {
	return 86400 // Default 24 hours in seconds
}

func (s *configService) GetBackupRetention() int {
	return 30 // Default 30 days
}

func (s *configService) GetBackupStoragePath() string {
	return "./backups" // Default path
}

func (s *configService) GetNotificationEnabled() bool {
	return false // Default disabled
}

func (s *configService) GetNotificationWebhookURL() string {
	return "" // Default empty
}

func (s *configService) GetFeatureQRCodeEnabled() bool {
	return true // Default enabled
}

func (s *configService) GetFeatureAnalyticsEnabled() bool {
	return true // Default enabled
}

func (s *configService) GetFeatureCustomDomainsEnabled() bool {
	return false // Default disabled
}

func (s *configService) GetFeaturePasswordProtectionEnabled() bool {
	return true // Default enabled
}

func (s *configService) GetFeatureExpirationEnabled() bool {
	return true // Default enabled
}

func (s *configService) GetFeatureBulkOperationsEnabled() bool {
	return false // Default disabled
}

func (s *configService) GetFeatureAPIEnabled() bool {
	return true // Default enabled
}

func (s *configService) GetFeatureWebhooksEnabled() bool {
	return false // Default disabled
}

func (s *configService) GetFeatureEmailNotificationsEnabled() bool {
	return false // Default disabled
}

func (s *configService) GetFeatureGeoLocationEnabled() bool {
	return true // Default enabled
}

func (s *configService) GetFeatureDeviceTrackingEnabled() bool {
	return true // Default enabled
}

func (s *configService) GetFeatureUserRegistrationEnabled() bool {
	return true // Default enabled
}

func (s *configService) GetFeatureGuestModeEnabled() bool {
	return false // Default disabled
}

func (s *configService) GetFeatureRateLimitingEnabled() bool {
	return s.config.Rate.Enabled
}

func (s *configService) GetFeatureSpamDetectionEnabled() bool {
	return false // Default disabled
}

func (s *configService) GetFeatureAuditLogsEnabled() bool {
	return false // Default disabled
}

func (s *configService) GetRawConfig() *config.Config {
	return s.config
}

func (s *configService) Validate() error {
	// Validation not implemented in config struct yet
	return nil
}