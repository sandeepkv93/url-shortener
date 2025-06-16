package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	CORS     CORSConfig
	Rate     RateLimitConfig
	External ExternalConfig
	App      AppConfig
	Security SecurityConfig
	Logging  LoggingConfig
	Cache    CacheConfig
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
	BcryptCost     int
	MaxRequestSize string
	EnableHTTPS    bool
}

type LoggingConfig struct {
	Level  string
	Format string
}

type CacheConfig struct {
	TTL    time.Duration
	URLTTL time.Duration
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
			BcryptCost:     getEnvInt("BCRYPT_COST", 12),
			MaxRequestSize: getEnv("MAX_REQUEST_SIZE", "10MB"),
			EnableHTTPS:    getEnvBool("ENABLE_HTTPS", false),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
		Cache: CacheConfig{
			TTL:    getEnvDuration("CACHE_TTL", "1h"),
			URLTTL: getEnvDuration("URL_CACHE_TTL", "24h"),
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

func (c *Config) IsDevelopment() bool {
	return c.Server.Env == "development"
}

func (c *Config) IsProduction() bool {
	return c.Server.Env == "production"
}

func (c *Config) GetServerAddress() string {
	return c.Server.Host + ":" + c.Server.Port
}