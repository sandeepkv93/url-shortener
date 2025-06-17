package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"url-shortener/internal/core/ports"
)

type CacheServiceImpl struct {
	redis *RedisClient
}

func NewCacheService(redis *RedisClient) ports.CacheService {
	return &CacheServiceImpl{
		redis: redis,
	}
}

// Basic cache operations
func (c *CacheServiceImpl) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.redis.Set(ctx, key, value, expiration)
}

func (c *CacheServiceImpl) Get(ctx context.Context, key string) (string, error) {
	return c.redis.Get(ctx, key)
}

func (c *CacheServiceImpl) Del(ctx context.Context, keys ...string) error {
	return c.redis.Del(ctx, keys...)
}

func (c *CacheServiceImpl) Exists(ctx context.Context, key string) (bool, error) {
	return c.redis.Exists(ctx, key)
}

func (c *CacheServiceImpl) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.redis.TTL(ctx, key)
}

// Counter operations
func (c *CacheServiceImpl) Incr(ctx context.Context, key string) (int64, error) {
	return c.redis.Incr(ctx, key)
}

func (c *CacheServiceImpl) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return c.redis.IncrBy(ctx, key, value)
}

// Set operations
func (c *CacheServiceImpl) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return c.redis.SAdd(ctx, key, members...)
}

func (c *CacheServiceImpl) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return c.redis.SIsMember(ctx, key, member)
}

func (c *CacheServiceImpl) SCard(ctx context.Context, key string) (int64, error) {
	return c.redis.SCard(ctx, key)
}

// Hash operations
func (c *CacheServiceImpl) HSet(ctx context.Context, key string, values ...interface{}) error {
	return c.redis.HSet(ctx, key, values...)
}

func (c *CacheServiceImpl) HGet(ctx context.Context, key, field string) (string, error) {
	return c.redis.HGet(ctx, key, field)
}

func (c *CacheServiceImpl) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return c.redis.HGetAll(ctx, key)
}

func (c *CacheServiceImpl) HDel(ctx context.Context, key string, fields ...string) error {
	return c.redis.HDel(ctx, key, fields...)
}

// URL-specific cache operations
func (c *CacheServiceImpl) CacheURL(ctx context.Context, shortCode, originalURL string, userID uint, expiration time.Duration) error {
	urlData := map[string]interface{}{
		"url":     originalURL,
		"user_id": userID,
		"cached_at": time.Now().Unix(),
	}
	
	jsonData, err := json.Marshal(urlData)
	if err != nil {
		return fmt.Errorf("failed to marshal URL data: %w", err)
	}
	
	key := fmt.Sprintf("url:%s", shortCode)
	return c.redis.Set(ctx, key, string(jsonData), expiration)
}

func (c *CacheServiceImpl) GetCachedURL(ctx context.Context, shortCode string) (string, uint, error) {
	key := fmt.Sprintf("url:%s", shortCode)
	data, err := c.redis.Get(ctx, key)
	if err != nil {
		return "", 0, err
	}
	
	var urlData map[string]interface{}
	if err := json.Unmarshal([]byte(data), &urlData); err != nil {
		return "", 0, fmt.Errorf("failed to unmarshal URL data: %w", err)
	}
	
	url, ok := urlData["url"].(string)
	if !ok {
		return "", 0, fmt.Errorf("invalid URL data format")
	}
	
	userIDFloat, ok := urlData["user_id"].(float64)
	if !ok {
		return "", 0, fmt.Errorf("invalid user ID data format")
	}
	
	return url, uint(userIDFloat), nil
}

func (c *CacheServiceImpl) InvalidateURL(ctx context.Context, shortCode string) error {
	key := fmt.Sprintf("url:%s", shortCode)
	return c.redis.Del(ctx, key)
}

// Rate limiting operations
func (c *CacheServiceImpl) IsRateLimited(ctx context.Context, key string, limit int64, window time.Duration) (bool, error) {
	current, err := c.redis.Get(ctx, key)
	if err != nil {
		// Key doesn't exist, not rate limited
		return false, nil
	}
	
	count, err := strconv.ParseInt(current, 10, 64)
	if err != nil {
		return false, fmt.Errorf("failed to parse rate limit count: %w", err)
	}
	
	return count >= limit, nil
}

func (c *CacheServiceImpl) IncrementRateLimit(ctx context.Context, key string, window time.Duration) (int64, error) {
	count, err := c.redis.Incr(ctx, key)
	if err != nil {
		return 0, err
	}
	
	// Set expiration only on first increment
	if count == 1 {
		if err := c.redis.Expire(ctx, key, window); err != nil {
			return count, fmt.Errorf("failed to set expiration: %w", err)
		}
	}
	
	return count, nil
}

// Session management
func (c *CacheServiceImpl) SetSession(ctx context.Context, token string, userID uint, expiration time.Duration) error {
	key := fmt.Sprintf("session:%s", token)
	return c.redis.Set(ctx, key, userID, expiration)
}

func (c *CacheServiceImpl) GetSession(ctx context.Context, token string) (uint, error) {
	key := fmt.Sprintf("session:%s", token)
	val, err := c.redis.Get(ctx, key)
	if err != nil {
		return 0, err
	}
	
	userID, err := strconv.ParseUint(val, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to parse user ID: %w", err)
	}
	
	return uint(userID), nil
}

func (c *CacheServiceImpl) InvalidateSession(ctx context.Context, token string) error {
	key := fmt.Sprintf("session:%s", token)
	return c.redis.Del(ctx, key)
}

// Analytics caching
func (c *CacheServiceImpl) CacheClickCount(ctx context.Context, shortCode string, count int64) error {
	key := fmt.Sprintf("clicks:%s", shortCode)
	return c.redis.Set(ctx, key, count, 24*time.Hour) // Cache for 24 hours
}

func (c *CacheServiceImpl) GetClickCount(ctx context.Context, shortCode string) (int64, error) {
	key := fmt.Sprintf("clicks:%s", shortCode)
	val, err := c.redis.Get(ctx, key)
	if err != nil {
		return 0, err
	}
	
	count, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse click count: %w", err)
	}
	
	return count, nil
}

func (c *CacheServiceImpl) IncrementClickCount(ctx context.Context, shortCode string) (int64, error) {
	key := fmt.Sprintf("clicks:%s", shortCode)
	count, err := c.redis.Incr(ctx, key)
	if err != nil {
		return 0, err
	}
	
	// Set expiration for cache invalidation
	if count == 1 {
		if err := c.redis.Expire(ctx, key, 24*time.Hour); err != nil {
			return count, fmt.Errorf("failed to set expiration: %w", err)
		}
	}
	
	return count, nil
}

func (c *CacheServiceImpl) CacheUniqueClick(ctx context.Context, shortCode, ipAddress string) (bool, error) {
	key := fmt.Sprintf("unique_clicks:%s", shortCode)
	
	// Check if IP already exists
	exists, err := c.redis.SIsMember(ctx, key, ipAddress)
	if err != nil {
		return false, err
	}
	
	if exists {
		return false, nil // Not a unique click
	}
	
	// Add IP to set
	if err := c.redis.SAdd(ctx, key, ipAddress); err != nil {
		return false, err
	}
	
	// Set expiration for cleanup
	if err := c.redis.Expire(ctx, key, 24*time.Hour); err != nil {
		return true, fmt.Errorf("failed to set expiration: %w", err)
	}
	
	return true, nil // Unique click
}

func (c *CacheServiceImpl) GetUniqueClickCount(ctx context.Context, shortCode string) (int64, error) {
	key := fmt.Sprintf("unique_clicks:%s", shortCode)
	return c.redis.SCard(ctx, key)
}

// Health and monitoring
func (c *CacheServiceImpl) Ping(ctx context.Context) error {
	return c.redis.Ping(ctx)
}

func (c *CacheServiceImpl) FlushDB(ctx context.Context) error {
	return c.redis.FlushDB(ctx)
}

func (c *CacheServiceImpl) Info(ctx context.Context) (string, error) {
	return c.redis.Info(ctx)
}

func (c *CacheServiceImpl) Close() error {
	return c.redis.Close()
}