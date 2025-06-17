package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RateLimitMiddlewareTestSuite struct {
	suite.Suite
	mockCache *MockCacheService
}

func TestRateLimitMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(RateLimitMiddlewareTestSuite))
}

func (suite *RateLimitMiddlewareTestSuite) SetupTest() {
	suite.mockCache = &MockCacheService{}
}

func (suite *RateLimitMiddlewareTestSuite) TestRateLimitMiddleware_AllowedRequest() {
	config := &RateLimitConfig{
		RequestsPerWindow: 10,
		WindowDuration:    time.Minute,
	}
	middleware := NewRateLimitMiddleware(suite.mockCache, config)

	// Mock cache responses
	suite.mockCache.On("IsRateLimited", mock.Anything, mock.AnythingOfType("string"), int64(10), time.Minute).Return(false, nil)
	suite.mockCache.On("IncrementRateLimit", mock.Anything, mock.AnythingOfType("string"), time.Minute).Return(int64(1), nil)

	handler := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	assert.Equal(suite.T(), "10", rr.Header().Get("X-RateLimit-Limit"))
	assert.Equal(suite.T(), "9", rr.Header().Get("X-RateLimit-Remaining"))
	assert.NotEmpty(suite.T(), rr.Header().Get("X-RateLimit-Reset"))
	
	suite.mockCache.AssertExpectations(suite.T())
}

func (suite *RateLimitMiddlewareTestSuite) TestRateLimitMiddleware_RateLimited() {
	config := &RateLimitConfig{
		RequestsPerWindow: 10,
		WindowDuration:    time.Minute,
	}
	middleware := NewRateLimitMiddleware(suite.mockCache, config)

	// Mock cache responses - rate limited
	suite.mockCache.On("IsRateLimited", mock.Anything, mock.AnythingOfType("string"), int64(10), time.Minute).Return(true, nil)

	handler := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.T().Error("Handler should not be called when rate limited")
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusTooManyRequests, rr.Code)
	assert.Contains(suite.T(), rr.Body.String(), "Rate limit exceeded")
	assert.Equal(suite.T(), "60", rr.Header().Get("Retry-After"))
	
	suite.mockCache.AssertExpectations(suite.T())
}

func (suite *RateLimitMiddlewareTestSuite) TestRateLimitMiddleware_CacheError() {
	config := &RateLimitConfig{
		RequestsPerWindow: 10,
		WindowDuration:    time.Minute,
	}
	middleware := NewRateLimitMiddleware(suite.mockCache, config)

	// Mock cache error - should fail open
	suite.mockCache.On("IsRateLimited", mock.Anything, mock.AnythingOfType("string"), int64(10), time.Minute).Return(false, assert.AnError)

	handler := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Should allow request through when cache fails
	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	
	suite.mockCache.AssertExpectations(suite.T())
}

func (suite *RateLimitMiddlewareTestSuite) TestRateLimitMiddleware_IncrementError() {
	config := &RateLimitConfig{
		RequestsPerWindow: 10,
		WindowDuration:    time.Minute,
	}
	middleware := NewRateLimitMiddleware(suite.mockCache, config)

	// Mock cache responses
	suite.mockCache.On("IsRateLimited", mock.Anything, mock.AnythingOfType("string"), int64(10), time.Minute).Return(false, nil)
	suite.mockCache.On("IncrementRateLimit", mock.Anything, mock.AnythingOfType("string"), time.Minute).Return(int64(0), assert.AnError)

	handler := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Should allow request through when increment fails
	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	
	suite.mockCache.AssertExpectations(suite.T())
}

func (suite *RateLimitMiddlewareTestSuite) TestRateLimitMiddleware_SkipPaths() {
	config := &RateLimitConfig{
		RequestsPerWindow: 10,
		WindowDuration:    time.Minute,
		SkipPaths:         []string{"/health", "/metrics"},
	}
	middleware := NewRateLimitMiddleware(suite.mockCache, config)

	handler := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Test health endpoint
	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(suite.T(), http.StatusOK, rr.Code)

	// Test metrics endpoint
	req = httptest.NewRequest("GET", "/metrics", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(suite.T(), http.StatusOK, rr.Code)

	// No cache methods should have been called
	suite.mockCache.AssertExpectations(suite.T())
}

func (suite *RateLimitMiddlewareTestSuite) TestRateLimitMiddleware_CustomKeyGenerator() {
	customKey := "custom-key"
	config := &RateLimitConfig{
		RequestsPerWindow: 10,
		WindowDuration:    time.Minute,
		KeyGenerator: func(r *http.Request) string {
			return customKey
		},
	}
	middleware := NewRateLimitMiddleware(suite.mockCache, config)

	suite.mockCache.On("IsRateLimited", mock.Anything, customKey, int64(10), time.Minute).Return(false, nil)
	suite.mockCache.On("IncrementRateLimit", mock.Anything, customKey, time.Minute).Return(int64(1), nil)

	handler := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusOK, rr.Code)
	suite.mockCache.AssertExpectations(suite.T())
}

func (suite *RateLimitMiddlewareTestSuite) TestRateLimitMiddleware_CustomRateLimitExceeded() {
	customResponse := "Custom rate limit message"
	config := &RateLimitConfig{
		RequestsPerWindow: 10,
		WindowDuration:    time.Minute,
		OnRateLimitExceeded: func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(customResponse))
		},
	}
	middleware := NewRateLimitMiddleware(suite.mockCache, config)

	suite.mockCache.On("IsRateLimited", mock.Anything, mock.AnythingOfType("string"), int64(10), time.Minute).Return(true, nil)

	handler := middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		suite.T().Error("Handler should not be called")
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusTooManyRequests, rr.Code)
	assert.Equal(suite.T(), customResponse, rr.Body.String())
	
	suite.mockCache.AssertExpectations(suite.T())
}

func (suite *RateLimitMiddlewareTestSuite) TestDefaultKeyGenerator() {
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	
	key := defaultKeyGenerator(req)
	assert.Contains(suite.T(), key, "rate_limit:")
	assert.Contains(suite.T(), key, "192.168.1.1")
}

func (suite *RateLimitMiddlewareTestSuite) TestDefaultKeyGenerator_WithHeaders() {
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "10.0.0.1:12345"
	req.Header.Set("X-Real-IP", "192.168.1.1")
	
	key := defaultKeyGenerator(req)
	assert.Contains(suite.T(), key, "192.168.1.1")

	req.Header.Del("X-Real-IP")
	req.Header.Set("X-Forwarded-For", "203.0.113.1")
	
	key = defaultKeyGenerator(req)
	assert.Contains(suite.T(), key, "203.0.113.1")
}

func (suite *RateLimitMiddlewareTestSuite) TestUserOrIPKeyGenerator() {
	// Test with authenticated user
	req := httptest.NewRequest("GET", "/test", nil)
	ctx := context.WithValue(req.Context(), "user_id", uint(123))
	req = req.WithContext(ctx)
	
	key := UserOrIPKeyGenerator(req)
	assert.Contains(suite.T(), key, "user:123")

	// Test without authenticated user
	req = httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	
	key = UserOrIPKeyGenerator(req)
	assert.Contains(suite.T(), key, "rate_limit:")
	assert.Contains(suite.T(), key, "192.168.1.1")
}

func (suite *RateLimitMiddlewareTestSuite) TestAPIKeyGenerator() {
	keyGen := APIKeyGenerator(100, 20)
	
	// Test with authenticated user
	req := httptest.NewRequest("GET", "/test", nil)
	ctx := context.WithValue(req.Context(), "user_id", uint(123))
	req = req.WithContext(ctx)
	
	key := keyGen(req)
	assert.Contains(suite.T(), key, "user:123:100")

	// Test without authenticated user
	req = httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	
	key = keyGen(req)
	assert.Contains(suite.T(), key, "ip:")
	assert.Contains(suite.T(), key, ":20")
}

func (suite *RateLimitMiddlewareTestSuite) TestGetClientIP() {
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "10.0.0.1:12345"
	
	// Test with X-Real-IP
	req.Header.Set("X-Real-IP", "192.168.1.1")
	ip := getClientIP(req)
	assert.Equal(suite.T(), "192.168.1.1", ip)

	// Test with X-Forwarded-For
	req.Header.Del("X-Real-IP")
	req.Header.Set("X-Forwarded-For", "203.0.113.1")
	ip = getClientIP(req)
	assert.Equal(suite.T(), "203.0.113.1", ip)

	// Test with RemoteAddr
	req.Header.Del("X-Forwarded-For")
	ip = getClientIP(req)
	assert.Equal(suite.T(), "10.0.0.1:12345", ip)
}

func (suite *RateLimitMiddlewareTestSuite) TestPredefinedConfigurations() {
	// Test GlobalRateLimit
	handler := GlobalRateLimit(suite.mockCache)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	assert.NotNil(suite.T(), handler)

	// Test APIRateLimit
	handler = APIRateLimit(suite.mockCache)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	assert.NotNil(suite.T(), handler)

	// Test AuthRateLimit
	handler = AuthRateLimit(suite.mockCache)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	assert.NotNil(suite.T(), handler)

	// Test URLCreationRateLimit
	handler = URLCreationRateLimit(suite.mockCache)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	assert.NotNil(suite.T(), handler)
}

func (suite *RateLimitMiddlewareTestSuite) TestDefaultRateLimitExceededHandler() {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	
	defaultRateLimitExceededHandler(rr, req)
	
	assert.Equal(suite.T(), http.StatusTooManyRequests, rr.Code)
	assert.Contains(suite.T(), rr.Body.String(), "Rate limit exceeded")
	assert.Equal(suite.T(), "60", rr.Header().Get("Retry-After"))
}

// Mock cache service for testing
type MockCacheService struct {
	mock.Mock
}

func (m *MockCacheService) IsRateLimited(ctx context.Context, key string, limit int64, window time.Duration) (bool, error) {
	args := m.Called(ctx, key, limit, window)
	return args.Bool(0), args.Error(1)
}

func (m *MockCacheService) IncrementRateLimit(ctx context.Context, key string, window time.Duration) (int64, error) {
	args := m.Called(ctx, key, window)
	return args.Get(0).(int64), args.Error(1)
}

// Implement other required methods (simplified for testing)
func (m *MockCacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return nil
}

func (m *MockCacheService) Get(ctx context.Context, key string) (string, error) {
	return "", nil
}

func (m *MockCacheService) Delete(ctx context.Context, key string) error {
	return nil
}

func (m *MockCacheService) Del(ctx context.Context, keys ...string) error {
	return nil
}

func (m *MockCacheService) Exists(ctx context.Context, key string) (bool, error) {
	return false, nil
}

func (m *MockCacheService) TTL(ctx context.Context, key string) (time.Duration, error) {
	return 0, nil
}

func (m *MockCacheService) Incr(ctx context.Context, key string) (int64, error) {
	return 0, nil
}

func (m *MockCacheService) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return 0, nil
}

func (m *MockCacheService) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return nil
}

func (m *MockCacheService) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return false, nil
}

func (m *MockCacheService) SCard(ctx context.Context, key string) (int64, error) {
	return 0, nil
}

func (m *MockCacheService) HSet(ctx context.Context, key string, values ...interface{}) error {
	return nil
}

func (m *MockCacheService) HGet(ctx context.Context, key, field string) (string, error) {
	return "", nil
}

func (m *MockCacheService) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return nil, nil
}

func (m *MockCacheService) HDel(ctx context.Context, key string, fields ...string) error {
	return nil
}

func (m *MockCacheService) CacheURL(ctx context.Context, shortCode, originalURL string, userID uint, expiration time.Duration) error {
	return nil
}

func (m *MockCacheService) GetCachedURL(ctx context.Context, shortCode string) (string, uint, error) {
	return "", 0, nil
}

func (m *MockCacheService) InvalidateURL(ctx context.Context, shortCode string) error {
	return nil
}

func (m *MockCacheService) SetSession(ctx context.Context, token string, userID uint, expiration time.Duration) error {
	return nil
}

func (m *MockCacheService) GetSession(ctx context.Context, token string) (uint, error) {
	return 0, nil
}

func (m *MockCacheService) InvalidateSession(ctx context.Context, token string) error {
	return nil
}

func (m *MockCacheService) CacheClickCount(ctx context.Context, shortCode string, count int64) error {
	return nil
}

func (m *MockCacheService) GetClickCount(ctx context.Context, shortCode string) (int64, error) {
	return 0, nil
}

func (m *MockCacheService) IncrementClickCount(ctx context.Context, shortCode string) (int64, error) {
	return 0, nil
}

func (m *MockCacheService) CacheUniqueClick(ctx context.Context, shortCode, ipAddress string) (bool, error) {
	return false, nil
}

func (m *MockCacheService) GetUniqueClickCount(ctx context.Context, shortCode string) (int64, error) {
	return 0, nil
}

func (m *MockCacheService) CacheAnalyticsData(ctx context.Context, key string, data interface{}, expiration time.Duration) error {
	return nil
}

func (m *MockCacheService) GetAnalyticsData(ctx context.Context, key string) (interface{}, error) {
	return nil, nil
}

func (m *MockCacheService) Ping(ctx context.Context) error {
	return nil
}

func (m *MockCacheService) FlushDB(ctx context.Context) error {
	return nil
}

func (m *MockCacheService) Info(ctx context.Context) (string, error) {
	return "", nil
}

func (m *MockCacheService) Close() error {
	return nil
}