package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"url-shortener/internal/config"
)

func TestConfigService_ServerConfig(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host:    "localhost",
			Port:    "8080",
			BaseURL: "http://localhost:8080",
		},
	}

	service := NewConfigService(cfg)

	assert.Equal(t, "localhost", service.GetServerHost())
	assert.Equal(t, "8080", service.GetServerPort())
	assert.Equal(t, "http://localhost:8080", service.GetBaseURL())
}

func TestConfigService_DatabaseConfig(t *testing.T) {
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			URL:             "postgres://user:pass@localhost/db",
			MaxOpenConns:    10,
			MaxIdleConns:    5,
			ConnMaxLifetime: 3600,
		},
	}

	service := NewConfigService(cfg)

	assert.Equal(t, "postgres://user:pass@localhost/db", service.GetDatabaseURL())
	assert.Equal(t, 10, service.GetDatabaseMaxOpenConns())
	assert.Equal(t, 5, service.GetDatabaseMaxIdleConns())
	assert.Equal(t, 3600, service.GetDatabaseConnMaxLifetime())
}

func TestConfigService_RedisConfig(t *testing.T) {
	cfg := &config.Config{
		Redis: config.RedisConfig{
			URL:        "redis://localhost:6379",
			MaxRetries: 3,
			RetryDelay: 1000,
			PoolSize:   10,
			Timeout:    5000,
		},
	}

	service := NewConfigService(cfg)

	assert.Equal(t, "redis://localhost:6379", service.GetRedisURL())
	assert.Equal(t, 3, service.GetRedisMaxRetries())
	assert.Equal(t, 1000, service.GetRedisRetryDelay())
	assert.Equal(t, 10, service.GetRedisPoolSize())
	assert.Equal(t, 5000, service.GetRedisTimeout())
}

func TestConfigService_JWTConfig(t *testing.T) {
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:           "secret-key",
			AccessTokenTTL:   3600,
			RefreshTokenTTL:  86400,
		},
	}

	service := NewConfigService(cfg)

	assert.Equal(t, "secret-key", service.GetJWTSecret())
	assert.Equal(t, 3600, service.GetJWTAccessTokenTTL())
	assert.Equal(t, 86400, service.GetJWTRefreshTokenTTL())
}

func TestConfigService_Environment(t *testing.T) {
	tests := []struct {
		name        string
		environment string
		isProd      bool
		isDev       bool
	}{
		{
			name:        "production environment",
			environment: "production",
			isProd:      true,
			isDev:       false,
		},
		{
			name:        "development environment",
			environment: "development",
			isProd:      false,
			isDev:       true,
		},
		{
			name:        "staging environment",
			environment: "staging",
			isProd:      false,
			isDev:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Environment: tt.environment,
			}

			service := NewConfigService(cfg)

			assert.Equal(t, tt.environment, service.GetEnvironment())
			assert.Equal(t, tt.isProd, service.IsProduction())
			assert.Equal(t, tt.isDev, service.IsDevelopment())
		})
	}
}

func TestConfigService_LoggingConfig(t *testing.T) {
	cfg := &config.Config{
		Logging: config.LoggingConfig{
			Level: "info",
		},
	}

	service := NewConfigService(cfg)

	assert.Equal(t, "info", service.GetLogLevel())
}

func TestConfigService_CORSConfig(t *testing.T) {
	cfg := &config.Config{
		CORS: config.CORSConfig{
			AllowedOrigins: []string{"http://localhost:3000", "https://example.com"},
		},
	}

	service := NewConfigService(cfg)

	origins := service.GetCORSOrigins()
	assert.Len(t, origins, 2)
	assert.Contains(t, origins, "http://localhost:3000")
	assert.Contains(t, origins, "https://example.com")
}

func TestConfigService_RateLimitConfig(t *testing.T) {
	cfg := &config.Config{
		RateLimit: config.RateLimitConfig{
			RequestsPerSecond: 10,
			BurstSize:         20,
		},
	}

	service := NewConfigService(cfg)

	assert.Equal(t, 10, service.GetRateLimitRequests())
	assert.Equal(t, 20, service.GetRateLimitBurst())
}

func TestConfigService_ExternalAPIConfig(t *testing.T) {
	cfg := &config.Config{
		ExternalAPI: config.ExternalAPIConfig{
			Timeout:    5000,
			MaxRetries: 3,
		},
	}

	service := NewConfigService(cfg)

	assert.Equal(t, 5000, service.GetExternalAPITimeout())
	assert.Equal(t, 3, service.GetExternalAPIRetries())
}

func TestConfigService_FileUploadConfig(t *testing.T) {
	cfg := &config.Config{
		FileUpload: config.FileUploadConfig{
			MaxSize:      10485760, // 10MB
			AllowedTypes: []string{"image/jpeg", "image/png"},
		},
	}

	service := NewConfigService(cfg)

	assert.Equal(t, int64(10485760), service.GetFileUploadMaxSize())
	
	allowedTypes := service.GetFileUploadAllowedTypes()
	assert.Len(t, allowedTypes, 2)
	assert.Contains(t, allowedTypes, "image/jpeg")
	assert.Contains(t, allowedTypes, "image/png")
}

func TestConfigService_EmailConfig(t *testing.T) {
	cfg := &config.Config{
		Email: config.EmailConfig{
			SMTPHost:     "smtp.gmail.com",
			SMTPPort:     587,
			SMTPUser:     "user@gmail.com",
			SMTPPassword: "password",
			FromAddress:  "noreply@example.com",
			FromName:     "URL Shortener",
		},
	}

	service := NewConfigService(cfg)

	assert.Equal(t, "smtp.gmail.com", service.GetEmailSMTPHost())
	assert.Equal(t, 587, service.GetEmailSMTPPort())
	assert.Equal(t, "user@gmail.com", service.GetEmailSMTPUser())
	assert.Equal(t, "password", service.GetEmailSMTPPassword())
	assert.Equal(t, "noreply@example.com", service.GetEmailFromAddress())
	assert.Equal(t, "URL Shortener", service.GetEmailFromName())
}

func TestConfigService_MonitoringConfig(t *testing.T) {
	cfg := &config.Config{
		Monitoring: config.MonitoringConfig{
			Enabled: true,
			Port:    "9090",
		},
	}

	service := NewConfigService(cfg)

	assert.True(t, service.GetMonitoringEnabled())
	assert.Equal(t, "9090", service.GetMonitoringPort())
}

func TestConfigService_TracingConfig(t *testing.T) {
	cfg := &config.Config{
		Tracing: config.TracingConfig{
			Enabled:  true,
			Endpoint: "http://jaeger:14268/api/traces",
		},
	}

	service := NewConfigService(cfg)

	assert.True(t, service.GetTracingEnabled())
	assert.Equal(t, "http://jaeger:14268/api/traces", service.GetTracingEndpoint())
}

func TestConfigService_CacheConfig(t *testing.T) {
	cfg := &config.Config{
		Cache: config.CacheConfig{
			DefaultTTL:   3600,
			URLTTL:       7200,
			SessionTTL:   1800,
			AnalyticsTTL: 300,
		},
	}

	service := NewConfigService(cfg)

	assert.Equal(t, 3600, service.GetCacheDefaultTTL())
	assert.Equal(t, 7200, service.GetCacheURLTTL())
	assert.Equal(t, 1800, service.GetCacheSessionTTL())
	assert.Equal(t, 300, service.GetCacheAnalyticsTTL())
}

func TestConfigService_SecurityConfig(t *testing.T) {
	cfg := &config.Config{
		Security: config.SecurityConfig{
			Enabled:                     true,
			HSTSEnabled:                 true,
			CSPEnabled:                  true,
			XSSProtectionEnabled:        true,
			ContentTypeNosniffEnabled:   true,
			FrameOptionsEnabled:         true,
		},
	}

	service := NewConfigService(cfg)

	assert.True(t, service.GetSecurityEnabled())
	assert.True(t, service.GetSecurityHSTSEnabled())
	assert.True(t, service.GetSecurityCSPEnabled())
	assert.True(t, service.GetSecurityXSSProtectionEnabled())
	assert.True(t, service.GetSecurityContentTypeNosniffEnabled())
	assert.True(t, service.GetSecurityFrameOptionsEnabled())
}

func TestConfigService_BackupConfig(t *testing.T) {
	cfg := &config.Config{
		Backup: config.BackupConfig{
			Enabled:     true,
			Interval:    3600,
			Retention:   168,
			StoragePath: "/backups",
		},
	}

	service := NewConfigService(cfg)

	assert.True(t, service.GetBackupEnabled())
	assert.Equal(t, 3600, service.GetBackupInterval())
	assert.Equal(t, 168, service.GetBackupRetention())
	assert.Equal(t, "/backups", service.GetBackupStoragePath())
}

func TestConfigService_NotificationConfig(t *testing.T) {
	cfg := &config.Config{
		Notification: config.NotificationConfig{
			Enabled:    true,
			WebhookURL: "https://hooks.slack.com/services/...",
		},
	}

	service := NewConfigService(cfg)

	assert.True(t, service.GetNotificationEnabled())
	assert.Equal(t, "https://hooks.slack.com/services/...", service.GetNotificationWebhookURL())
}

func TestConfigService_FeatureFlags(t *testing.T) {
	cfg := &config.Config{
		Features: config.FeaturesConfig{
			QRCodeEnabled:             true,
			AnalyticsEnabled:          true,
			CustomDomainsEnabled:      false,
			PasswordProtectionEnabled: true,
			ExpirationEnabled:         true,
			BulkOperationsEnabled:     false,
			APIEnabled:                true,
			WebhooksEnabled:           false,
			EmailNotificationsEnabled: true,
			GeoLocationEnabled:        true,
			DeviceTrackingEnabled:     true,
			UserRegistrationEnabled:   true,
			GuestModeEnabled:          false,
			RateLimitingEnabled:       true,
			SpamDetectionEnabled:      true,
			AuditLogsEnabled:          false,
		},
	}

	service := NewConfigService(cfg)

	// Test enabled features
	assert.True(t, service.GetFeatureQRCodeEnabled())
	assert.True(t, service.GetFeatureAnalyticsEnabled())
	assert.True(t, service.GetFeaturePasswordProtectionEnabled())
	assert.True(t, service.GetFeatureExpirationEnabled())
	assert.True(t, service.GetFeatureAPIEnabled())
	assert.True(t, service.GetFeatureEmailNotificationsEnabled())
	assert.True(t, service.GetFeatureGeoLocationEnabled())
	assert.True(t, service.GetFeatureDeviceTrackingEnabled())
	assert.True(t, service.GetFeatureUserRegistrationEnabled())
	assert.True(t, service.GetFeatureRateLimitingEnabled())
	assert.True(t, service.GetFeatureSpamDetectionEnabled())

	// Test disabled features
	assert.False(t, service.GetFeatureCustomDomainsEnabled())
	assert.False(t, service.GetFeatureBulkOperationsEnabled())
	assert.False(t, service.GetFeatureWebhooksEnabled())
	assert.False(t, service.GetFeatureGuestModeEnabled())
	assert.False(t, service.GetFeatureAuditLogsEnabled())
}

func TestConfigService_GetRawConfig(t *testing.T) {
	cfg := &config.Config{
		Environment: "test",
		Server: config.ServerConfig{
			Host: "localhost",
			Port: "8080",
		},
	}

	service := NewConfigService(cfg)

	rawConfig := service.GetRawConfig()
	assert.Equal(t, cfg, rawConfig)
	assert.Equal(t, "test", rawConfig.Environment)
	assert.Equal(t, "localhost", rawConfig.Server.Host)
}

func TestConfigService_Validate(t *testing.T) {
	// Create a minimal valid config
	cfg := &config.Config{
		Environment: "test",
		Server: config.ServerConfig{
			Host: "localhost",
			Port: "8080",
		},
	}

	service := NewConfigService(cfg)

	// Note: This test depends on the actual implementation of cfg.Validate()
	// which should be implemented in the config package
	err := service.Validate()
	// For now, we'll just ensure the method exists and doesn't panic
	// In a real implementation, you'd test various validation scenarios
	assert.NotPanics(t, func() { service.Validate() })
	_ = err // We're not testing the actual validation logic here
}

func TestConfigService_NilConfig(t *testing.T) {
	// Test that the service handles nil config gracefully
	service := NewConfigService(nil)
	
	// These should not panic, but may return zero values
	assert.NotPanics(t, func() { service.GetServerHost() })
	assert.NotPanics(t, func() { service.GetServerPort() })
	assert.NotPanics(t, func() { service.GetBaseURL() })
}

func TestConfigService_DefaultValues(t *testing.T) {
	// Test with empty config to ensure default/zero values are handled
	cfg := &config.Config{}
	service := NewConfigService(cfg)

	// These should return zero values without panicking
	assert.Equal(t, "", service.GetServerHost())
	assert.Equal(t, "", service.GetServerPort())
	assert.Equal(t, "", service.GetBaseURL())
	assert.Equal(t, "", service.GetEnvironment())
	assert.False(t, service.IsProduction())
	assert.False(t, service.IsDevelopment())
	assert.Equal(t, 0, service.GetRateLimitRequests())
	assert.Equal(t, 0, service.GetRateLimitBurst())
}