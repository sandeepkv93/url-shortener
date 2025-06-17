package ports

import (
	"context"
	"time"
)

type CacheService interface {
	// Basic cache operations
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, key string) (bool, error)
	TTL(ctx context.Context, key string) (time.Duration, error)

	// Counter operations for analytics
	Incr(ctx context.Context, key string) (int64, error)
	IncrBy(ctx context.Context, key string, value int64) (int64, error)

	// Set operations for unique tracking
	SAdd(ctx context.Context, key string, members ...interface{}) error
	SIsMember(ctx context.Context, key string, member interface{}) (bool, error)
	SCard(ctx context.Context, key string) (int64, error)

	// Hash operations for complex data
	HSet(ctx context.Context, key string, values ...interface{}) error
	HGet(ctx context.Context, key, field string) (string, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HDel(ctx context.Context, key string, fields ...string) error

	// URL-specific cache operations
	CacheURL(ctx context.Context, shortCode, originalURL string, userID uint, expiration time.Duration) error
	GetCachedURL(ctx context.Context, shortCode string) (string, uint, error)
	InvalidateURL(ctx context.Context, shortCode string) error

	// Rate limiting operations
	IsRateLimited(ctx context.Context, key string, limit int64, window time.Duration) (bool, error)
	IncrementRateLimit(ctx context.Context, key string, window time.Duration) (int64, error)

	// Session management
	SetSession(ctx context.Context, token string, userID uint, expiration time.Duration) error
	GetSession(ctx context.Context, token string) (uint, error)
	InvalidateSession(ctx context.Context, token string) error

	// Analytics caching
	CacheClickCount(ctx context.Context, shortCode string, count int64) error
	GetClickCount(ctx context.Context, shortCode string) (int64, error)
	IncrementClickCount(ctx context.Context, shortCode string) (int64, error)
	CacheUniqueClick(ctx context.Context, shortCode, ipAddress string) (bool, error)
	GetUniqueClickCount(ctx context.Context, shortCode string) (int64, error)

	// Health and monitoring
	Ping(ctx context.Context) error
	FlushDB(ctx context.Context) error
	Info(ctx context.Context) (string, error)
	Close() error
}