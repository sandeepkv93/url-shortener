package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"url-shortener/internal/core/ports"
)

type CacheServiceTestSuite struct {
	suite.Suite
	cache ports.CacheService
	redis *RedisClient
	ctx   context.Context
}

func (suite *CacheServiceTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	
	// Try connecting to Redis with and without password
	var client *RedisClient
	var err error
	
	// First try with docker-compose password
	client, err = NewRedisClient("localhost:6379", "redis123", 15)
	if err != nil {
		// Try without password (local Redis)
		client, err = NewRedisClient("localhost:6379", "", 15)
		if err != nil {
			suite.T().Skip("Redis not available for testing")
			return
		}
	}
	
	suite.redis = client
	suite.cache = NewCacheService(client)
	
	// Clear test database
	err = suite.redis.FlushDB(suite.ctx)
	suite.Require().NoError(err)
}

func (suite *CacheServiceTestSuite) TearDownSuite() {
	if suite.redis != nil {
		// Clear test database
		_ = suite.redis.FlushDB(suite.ctx)
		_ = suite.redis.Close()
	}
}

func (suite *CacheServiceTestSuite) SetupTest() {
	if suite.redis != nil {
		// Clear database before each test
		_ = suite.redis.FlushDB(suite.ctx)
	}
}

func (suite *CacheServiceTestSuite) TestURLCaching() {
	if suite.cache == nil {
		suite.T().Skip("Redis not available")
		return
	}
	
	shortCode := "test123"
	originalURL := "https://example.com"
	userID := uint(1)
	
	// Test CacheURL
	err := suite.cache.CacheURL(suite.ctx, shortCode, originalURL, userID, time.Hour)
	suite.NoError(err)
	
	// Test GetCachedURL
	cachedURL, cachedUserID, err := suite.cache.GetCachedURL(suite.ctx, shortCode)
	suite.NoError(err)
	suite.Equal(originalURL, cachedURL)
	suite.Equal(userID, cachedUserID)
	
	// Test InvalidateURL
	err = suite.cache.InvalidateURL(suite.ctx, shortCode)
	suite.NoError(err)
	
	_, _, err = suite.cache.GetCachedURL(suite.ctx, shortCode)
	suite.Error(err)
}

func (suite *CacheServiceTestSuite) TestRateLimiting() {
	if suite.cache == nil {
		suite.T().Skip("Redis not available")
		return
	}
	
	key := "rate_limit:test_ip"
	limit := int64(5)
	window := time.Minute
	
	// First check should not be rate limited
	isLimited, err := suite.cache.IsRateLimited(suite.ctx, key, limit, window)
	suite.NoError(err)
	suite.False(isLimited)
	
	// Increment counter multiple times
	for i := 0; i < 5; i++ {
		count, err := suite.cache.IncrementRateLimit(suite.ctx, key, window)
		suite.NoError(err)
		suite.Equal(int64(i+1), count)
	}
	
	// Now should be rate limited
	isLimited, err = suite.cache.IsRateLimited(suite.ctx, key, limit, window)
	suite.NoError(err)
	suite.True(isLimited)
	
	// One more increment should still work but exceed limit
	count, err := suite.cache.IncrementRateLimit(suite.ctx, key, window)
	suite.NoError(err)
	suite.Equal(int64(6), count)
	
	isLimited, err = suite.cache.IsRateLimited(suite.ctx, key, limit, window)
	suite.NoError(err)
	suite.True(isLimited)
}

func (suite *CacheServiceTestSuite) TestSessionManagement() {
	if suite.cache == nil {
		suite.T().Skip("Redis not available")
		return
	}
	
	token := "test_token_123"
	userID := uint(42)
	
	// Test SetSession
	err := suite.cache.SetSession(suite.ctx, token, userID, time.Hour)
	suite.NoError(err)
	
	// Test GetSession
	retrievedUserID, err := suite.cache.GetSession(suite.ctx, token)
	suite.NoError(err)
	suite.Equal(userID, retrievedUserID)
	
	// Test InvalidateSession
	err = suite.cache.InvalidateSession(suite.ctx, token)
	suite.NoError(err)
	
	_, err = suite.cache.GetSession(suite.ctx, token)
	suite.Error(err)
}

func (suite *CacheServiceTestSuite) TestClickAnalytics() {
	if suite.cache == nil {
		suite.T().Skip("Redis not available")
		return
	}
	
	shortCode := "test456"
	
	// Test initial click count
	count, err := suite.cache.IncrementClickCount(suite.ctx, shortCode)
	suite.NoError(err)
	suite.Equal(int64(1), count)
	
	// Test multiple increments
	for i := 2; i <= 5; i++ {
		count, err := suite.cache.IncrementClickCount(suite.ctx, shortCode)
		suite.NoError(err)
		suite.Equal(int64(i), count)
	}
	
	// Test GetClickCount
	retrievedCount, err := suite.cache.GetClickCount(suite.ctx, shortCode)
	suite.NoError(err)
	suite.Equal(int64(5), retrievedCount)
	
	// Test CacheClickCount (overwrite)
	err = suite.cache.CacheClickCount(suite.ctx, shortCode, 100)
	suite.NoError(err)
	
	retrievedCount, err = suite.cache.GetClickCount(suite.ctx, shortCode)
	suite.NoError(err)
	suite.Equal(int64(100), retrievedCount)
}

func (suite *CacheServiceTestSuite) TestUniqueClickTracking() {
	if suite.cache == nil {
		suite.T().Skip("Redis not available")
		return
	}
	
	shortCode := "test789"
	ip1 := "192.168.1.1"
	ip2 := "192.168.1.2"
	
	// Test first unique click
	isUnique, err := suite.cache.CacheUniqueClick(suite.ctx, shortCode, ip1)
	suite.NoError(err)
	suite.True(isUnique)
	
	// Test second unique click from different IP
	isUnique, err = suite.cache.CacheUniqueClick(suite.ctx, shortCode, ip2)
	suite.NoError(err)
	suite.True(isUnique)
	
	// Test duplicate click from same IP
	isUnique, err = suite.cache.CacheUniqueClick(suite.ctx, shortCode, ip1)
	suite.NoError(err)
	suite.False(isUnique)
	
	// Test unique click count
	uniqueCount, err := suite.cache.GetUniqueClickCount(suite.ctx, shortCode)
	suite.NoError(err)
	suite.Equal(int64(2), uniqueCount)
}

func (suite *CacheServiceTestSuite) TestHealthOperations() {
	if suite.cache == nil {
		suite.T().Skip("Redis not available")
		return
	}
	
	// Test Ping
	err := suite.cache.Ping(suite.ctx)
	suite.NoError(err)
	
	// Test Info
	info, err := suite.cache.Info(suite.ctx)
	suite.NoError(err)
	suite.NotEmpty(info)
	suite.Contains(info, "redis_version")
}

func (suite *CacheServiceTestSuite) TestBasicCacheOperations() {
	if suite.cache == nil {
		suite.T().Skip("Redis not available")
		return
	}
	
	// Test Set and Get
	err := suite.cache.Set(suite.ctx, "test:basic", "test_value", time.Minute)
	suite.NoError(err)
	
	value, err := suite.cache.Get(suite.ctx, "test:basic")
	suite.NoError(err)
	suite.Equal("test_value", value)
	
	// Test Exists
	exists, err := suite.cache.Exists(suite.ctx, "test:basic")
	suite.NoError(err)
	suite.True(exists)
	
	// Test TTL
	ttl, err := suite.cache.TTL(suite.ctx, "test:basic")
	suite.NoError(err)
	suite.True(ttl > 0)
	suite.True(ttl <= time.Minute)
	
	// Test Del
	err = suite.cache.Del(suite.ctx, "test:basic")
	suite.NoError(err)
	
	exists, err = suite.cache.Exists(suite.ctx, "test:basic")
	suite.NoError(err)
	suite.False(exists)
}

func TestCacheServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CacheServiceTestSuite))
}