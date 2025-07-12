package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server      ServerConfig
	Database    DatabaseConfig
	Redis       RedisConfig
	JWT         JWTConfig
	CORS        CORSConfig
	Rate        RateLimitConfig
	External    ExternalConfig
	App         AppConfig
	Security    SecurityConfig
	Logging     LoggingConfig
	Cache       CacheConfig
	Performance PerformanceConfig
	Monitoring  MonitoringConfig
}

type ServerConfig struct {
	Port string
	Host string
	Env  string
}

type DatabaseConfig struct {
	URL            string
	Host           string
	Port           string
	Name           string
	User           string
	Password       string
	MaxConnections int
	MaxIdle        int
}

type RedisConfig struct {
	URL      string
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret        string
	Expiry        time.Duration
	RefreshExpiry time.Duration
}

type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

type RateLimitConfig struct {
	Requests int
	Window   time.Duration
	Enabled  bool
}

type ExternalConfig struct {
	GeolocationAPIKey string
	GeolocationAPIURL string
}

type AppConfig struct {
	BaseURL              string
	FrontendURL          string
	ShortCodeLength      int
	DefaultExpiryDays    int
	MaxCustomAliasLength int
}

type SecurityConfig struct {
	BcryptCost              int
	MaxRequestSize          string
	EnableHTTPS             bool
	EnableSecurityHeaders   bool
	EnableInputValidation   bool
	EnableXSSProtection     bool
	EnableCSRFProtection    bool
	EnableSanitization      bool
	EnableSecureCookies     bool
	EnableIPRateLimit       bool
	IPRateLimitPerMinute    int
	BlockSuspiciousURLs     bool
	AllowedContentTypes     []string
	BlockedDomains          []string
	CSPDirectives           map[string]string
}

type LoggingConfig struct {
	Level  string
	Format string
}

type CacheConfig struct {
	TTL    time.Duration
	URLTTL time.Duration
}

type PerformanceConfig struct {
	EnableCompression         bool
	CompressionLevel          int
	EnablePerformanceMonitoring bool
	EnableCacheHeaders        bool
	MaxRequestSize            int64
	SlowQueryThreshold        time.Duration
	EnableQueryMonitoring     bool
	EnableBatchOperations     bool
	BatchSize                 int
	MaxConcurrentQueries      int
	ConnectionPoolOptimization bool
}

type MonitoringConfig struct {
	EnableMetrics       bool
	EnableHealthChecks  bool
	EnablePerformanceAlerts bool
	MetricsInterval     time.Duration
	HealthCheckInterval time.Duration
	AlertThresholds     AlertThresholds
}

type AlertThresholds struct {
	SlowQueryThreshold      time.Duration
	HighMemoryThreshold     int64
	HighCPUThreshold        float64
	ErrorRateThreshold      float64
	ResponseTimeThreshold   time.Duration
	CacheHitRateThreshold   float64
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		// It's okay if .env file doesn't exist in production
	}

	config := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Host: getEnv("HOST", "localhost"),
			Env:  getEnv("GO_ENV", "development"),
		},
		Database: DatabaseConfig{
			URL:            getEnv("DATABASE_URL", "postgres://username:password@localhost:5432/urlshortener?sslmode=disable"),
			Host:           getEnv("DATABASE_HOST", "localhost"),
			Port:           getEnv("DATABASE_PORT", "5432"),
			Name:           getEnv("DATABASE_NAME", "urlshortener"),
			User:           getEnv("DATABASE_USER", "username"),
			Password:       getEnv("DATABASE_PASSWORD", "password"),
			MaxConnections: getEnvInt("DATABASE_MAX_CONNECTIONS", 25),
			MaxIdle:        getEnvInt("DATABASE_MAX_IDLE_CONNECTIONS", 5),
		},
		Redis: RedisConfig{
			URL:      getEnv("REDIS_URL", "redis://localhost:6379"),
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:        getEnv("JWT_SECRET", "your-jwt-secret-key-here-change-in-production"),
			Expiry:        getEnvDuration("JWT_EXPIRY", "24h"),
			RefreshExpiry: getEnvDuration("JWT_REFRESH_EXPIRY", "7d"),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnvStringSlice("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000", "http://localhost:5173"}),
			AllowedMethods: getEnvStringSlice("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			AllowedHeaders: getEnvStringSlice("CORS_ALLOWED_HEADERS", []string{"Content-Type", "Authorization"}),
		},
		Rate: RateLimitConfig{
			Requests: getEnvInt("RATE_LIMIT_REQUESTS", 100),
			Window:   getEnvDuration("RATE_LIMIT_WINDOW", "1h"),
			Enabled:  getEnvBool("RATE_LIMIT_ENABLED", true),
		},
		External: ExternalConfig{
			GeolocationAPIKey: getEnv("GEOLOCATION_API_KEY", ""),
			GeolocationAPIURL: getEnv("GEOLOCATION_API_URL", "https://api.ipgeolocation.io/ipgeo"),
		},
		App: AppConfig{
			BaseURL:              getEnv("BASE_URL", "http://localhost:8080"),
			FrontendURL:          getEnv("FRONTEND_URL", "http://localhost:3000"),
			ShortCodeLength:      getEnvInt("SHORT_CODE_LENGTH", 8),
			DefaultExpiryDays:    getEnvInt("DEFAULT_EXPIRY_DAYS", 365),
			MaxCustomAliasLength: getEnvInt("MAX_CUSTOM_ALIAS_LENGTH", 50),
		},
		Security: SecurityConfig{
			BcryptCost:              getEnvInt("BCRYPT_COST", 12),
			MaxRequestSize:          getEnv("MAX_REQUEST_SIZE", "10MB"),
			EnableHTTPS:             getEnvBool("ENABLE_HTTPS", false),
			EnableSecurityHeaders:   getEnvBool("ENABLE_SECURITY_HEADERS", true),
			EnableInputValidation:   getEnvBool("ENABLE_INPUT_VALIDATION", true),
			EnableXSSProtection:     getEnvBool("ENABLE_XSS_PROTECTION", true),
			EnableCSRFProtection:    getEnvBool("ENABLE_CSRF_PROTECTION", true),
			EnableSanitization:      getEnvBool("ENABLE_SANITIZATION", true),
			EnableSecureCookies:     getEnvBool("ENABLE_SECURE_COOKIES", false),
			EnableIPRateLimit:       getEnvBool("ENABLE_IP_RATE_LIMIT", true),
			IPRateLimitPerMinute:    getEnvInt("IP_RATE_LIMIT_PER_MINUTE", 100),
			BlockSuspiciousURLs:     getEnvBool("BLOCK_SUSPICIOUS_URLS", true),
			AllowedContentTypes:     getEnvStringSlice("ALLOWED_CONTENT_TYPES", []string{"application/json", "application/x-www-form-urlencoded", "multipart/form-data"}),
			BlockedDomains:          getEnvStringSlice("BLOCKED_DOMAINS", []string{}),
			CSPDirectives: map[string]string{
				"default-src": getEnv("CSP_DEFAULT_SRC", "'self'"),
				"script-src":  getEnv("CSP_SCRIPT_SRC", "'self' 'unsafe-inline' 'unsafe-eval'"),
				"style-src":   getEnv("CSP_STYLE_SRC", "'self' 'unsafe-inline'"),
				"img-src":     getEnv("CSP_IMG_SRC", "'self' data: https:"),
			},
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
		Cache: CacheConfig{
			TTL:    getEnvDuration("CACHE_TTL", "1h"),
			URLTTL: getEnvDuration("URL_CACHE_TTL", "24h"),
		},
		Performance: PerformanceConfig{
			EnableCompression:          getEnvBool("ENABLE_COMPRESSION", true),
			CompressionLevel:           getEnvInt("COMPRESSION_LEVEL", 6),
			EnablePerformanceMonitoring: getEnvBool("ENABLE_PERFORMANCE_MONITORING", true),
			EnableCacheHeaders:         getEnvBool("ENABLE_CACHE_HEADERS", true),
			MaxRequestSize:             getEnvInt64("MAX_REQUEST_SIZE", 10<<20), // 10MB
			SlowQueryThreshold:         getEnvDuration("SLOW_QUERY_THRESHOLD", "100ms"),
			EnableQueryMonitoring:      getEnvBool("ENABLE_QUERY_MONITORING", true),
			EnableBatchOperations:      getEnvBool("ENABLE_BATCH_OPERATIONS", true),
			BatchSize:                  getEnvInt("BATCH_SIZE", 100),
			MaxConcurrentQueries:       getEnvInt("MAX_CONCURRENT_QUERIES", 50),
			ConnectionPoolOptimization: getEnvBool("CONNECTION_POOL_OPTIMIZATION", true),
		},
		Monitoring: MonitoringConfig{
			EnableMetrics:       getEnvBool("ENABLE_METRICS", true),
			EnableHealthChecks:  getEnvBool("ENABLE_HEALTH_CHECKS", true),
			EnablePerformanceAlerts: getEnvBool("ENABLE_PERFORMANCE_ALERTS", true),
			MetricsInterval:     getEnvDuration("METRICS_INTERVAL", "30s"),
			HealthCheckInterval: getEnvDuration("HEALTH_CHECK_INTERVAL", "10s"),
			AlertThresholds: AlertThresholds{
				SlowQueryThreshold:      getEnvDuration("ALERT_SLOW_QUERY_THRESHOLD", "500ms"),
				HighMemoryThreshold:     getEnvInt64("ALERT_HIGH_MEMORY_THRESHOLD", 1<<30), // 1GB
				HighCPUThreshold:        getEnvFloat64("ALERT_HIGH_CPU_THRESHOLD", 80.0),
				ErrorRateThreshold:      getEnvFloat64("ALERT_ERROR_RATE_THRESHOLD", 5.0),
				ResponseTimeThreshold:   getEnvDuration("ALERT_RESPONSE_TIME_THRESHOLD", "1s"),
				CacheHitRateThreshold:   getEnvFloat64("ALERT_CACHE_HIT_RATE_THRESHOLD", 70.0),
			},
		},
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	if duration, err := time.ParseDuration(defaultValue); err == nil {
		return duration
	}
	return time.Hour
}

func getEnvStringSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvFloat64(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

func (c *Config) IsDevelopment() bool {
	return c.Server.Env == "development"
}

func (c *Config) IsProduction() bool {
	return c.Server.Env == "production"
}

func (c *Config) GetServerAddress() string {
	return c.Server.Host + ":" + c.Server.Port
}

func (c *Config) GetAllowedOrigins() []string {
	return c.CORS.AllowedOrigins
}

// Performance configuration helper methods
func (c *Config) IsCompressionEnabled() bool {
	return c.Performance.EnableCompression
}

func (c *Config) IsPerformanceMonitoringEnabled() bool {
	return c.Performance.EnablePerformanceMonitoring
}

func (c *Config) IsCacheHeadersEnabled() bool {
	return c.Performance.EnableCacheHeaders
}

func (c *Config) IsQueryMonitoringEnabled() bool {
	return c.Performance.EnableQueryMonitoring
}

func (c *Config) IsBatchOperationsEnabled() bool {
	return c.Performance.EnableBatchOperations
}

func (c *Config) IsConnectionPoolOptimizationEnabled() bool {
	return c.Performance.ConnectionPoolOptimization
}

func (c *Config) GetCompressionLevel() int {
	return c.Performance.CompressionLevel
}

func (c *Config) GetMaxRequestSize() int64 {
	return c.Performance.MaxRequestSize
}

func (c *Config) GetSlowQueryThreshold() time.Duration {
	return c.Performance.SlowQueryThreshold
}

func (c *Config) GetBatchSize() int {
	return c.Performance.BatchSize
}

func (c *Config) GetMaxConcurrentQueries() int {
	return c.Performance.MaxConcurrentQueries
}

// Monitoring configuration helper methods
func (c *Config) IsMetricsEnabled() bool {
	return c.Monitoring.EnableMetrics
}

func (c *Config) IsHealthChecksEnabled() bool {
	return c.Monitoring.EnableHealthChecks
}

func (c *Config) IsPerformanceAlertsEnabled() bool {
	return c.Monitoring.EnablePerformanceAlerts
}

func (c *Config) GetMetricsInterval() time.Duration {
	return c.Monitoring.MetricsInterval
}

func (c *Config) GetHealthCheckInterval() time.Duration {
	return c.Monitoring.HealthCheckInterval
}

func (c *Config) GetAlertThresholds() AlertThresholds {
	return c.Monitoring.AlertThresholds
}