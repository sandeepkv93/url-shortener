package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type RedisTestSuite struct {
	suite.Suite
	redis *RedisClient
	ctx   context.Context
}

func (suite *RedisTestSuite) SetupSuite() {
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
	
	// Clear test database
	err = suite.redis.FlushDB(suite.ctx)
	suite.Require().NoError(err)
}

func (suite *RedisTestSuite) TearDownSuite() {
	if suite.redis != nil {
		// Clear test database
		_ = suite.redis.FlushDB(suite.ctx)
		_ = suite.redis.Close()
	}
}

func (suite *RedisTestSuite) SetupTest() {
	if suite.redis != nil {
		// Clear database before each test
		_ = suite.redis.FlushDB(suite.ctx)
	}
}

func (suite *RedisTestSuite) TestConnection() {
	if suite.redis == nil {
		suite.T().Skip("Redis not available")
		return
	}
	
	err := suite.redis.Ping(suite.ctx)
	suite.NoError(err)
}

func (suite *RedisTestSuite) TestBasicOperations() {
	if suite.redis == nil {
		suite.T().Skip("Redis not available")
		return
	}
	
	// Test Set and Get
	err := suite.redis.Set(suite.ctx, "test:key", "test_value", time.Minute)
	suite.NoError(err)
	
	value, err := suite.redis.Get(suite.ctx, "test:key")
	suite.NoError(err)
	suite.Equal("test_value", value)
	
	// Test Exists
	exists, err := suite.redis.Exists(suite.ctx, "test:key")
	suite.NoError(err)
	suite.True(exists)
	
	// Test Del
	err = suite.redis.Del(suite.ctx, "test:key")
	suite.NoError(err)
	
	exists, err = suite.redis.Exists(suite.ctx, "test:key")
	suite.NoError(err)
	suite.False(exists)
}

func (suite *RedisTestSuite) TestCounterOperations() {
	if suite.redis == nil {
		suite.T().Skip("Redis not available")
		return
	}
	
	// Test Incr
	count, err := suite.redis.Incr(suite.ctx, "test:counter")
	suite.NoError(err)
	suite.Equal(int64(1), count)
	
	count, err = suite.redis.Incr(suite.ctx, "test:counter")
	suite.NoError(err)
	suite.Equal(int64(2), count)
	
	// Test IncrBy
	count, err = suite.redis.IncrBy(suite.ctx, "test:counter", 5)
	suite.NoError(err)
	suite.Equal(int64(7), count)
}

func (suite *RedisTestSuite) TestSetOperations() {
	if suite.redis == nil {
		suite.T().Skip("Redis not available")
		return
	}
	
	// Test SAdd
	err := suite.redis.SAdd(suite.ctx, "test:set", "member1", "member2", "member3")
	suite.NoError(err)
	
	// Test SIsMember
	isMember, err := suite.redis.SIsMember(suite.ctx, "test:set", "member1")
	suite.NoError(err)
	suite.True(isMember)
	
	isMember, err = suite.redis.SIsMember(suite.ctx, "test:set", "member4")
	suite.NoError(err)
	suite.False(isMember)
	
	// Test SCard
	count, err := suite.redis.SCard(suite.ctx, "test:set")
	suite.NoError(err)
	suite.Equal(int64(3), count)
}

func (suite *RedisTestSuite) TestHashOperations() {
	if suite.redis == nil {
		suite.T().Skip("Redis not available")
		return
	}
	
	// Test HSet
	err := suite.redis.HSet(suite.ctx, "test:hash", "field1", "value1", "field2", "value2")
	suite.NoError(err)
	
	// Test HGet
	value, err := suite.redis.HGet(suite.ctx, "test:hash", "field1")
	suite.NoError(err)
	suite.Equal("value1", value)
	
	// Test HGetAll
	all, err := suite.redis.HGetAll(suite.ctx, "test:hash")
	suite.NoError(err)
	suite.Equal(map[string]string{
		"field1": "value1",
		"field2": "value2",
	}, all)
	
	// Test HDel
	err = suite.redis.HDel(suite.ctx, "test:hash", "field1")
	suite.NoError(err)
	
	_, err = suite.redis.HGet(suite.ctx, "test:hash", "field1")
	suite.Error(err)
}

func (suite *RedisTestSuite) TestExpiration() {
	if suite.redis == nil {
		suite.T().Skip("Redis not available")
		return
	}
	
	// Set key with short expiration
	err := suite.redis.Set(suite.ctx, "test:expire", "value", 500*time.Millisecond)
	suite.NoError(err)
	
	// Verify key exists immediately
	exists, err := suite.redis.Exists(suite.ctx, "test:expire")
	suite.NoError(err)
	suite.True(exists, "Key should exist immediately after setting")
	
	// Wait for expiration
	time.Sleep(600 * time.Millisecond)
	
	// Verify key no longer exists
	exists, err = suite.redis.Exists(suite.ctx, "test:expire")
	suite.NoError(err)
	suite.False(exists, "Key should not exist after expiration")
}

func TestRedisTestSuite(t *testing.T) {
	suite.Run(t, new(RedisTestSuite))
}