package services

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
	"url-shortener/internal/core/domain"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, offset, limit int) ([]*domain.User, int64, error) {
	args := m.Called(ctx, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) GetUserStats(ctx context.Context, userID uint) (*domain.UserStats, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserStats), args.Error(1)
}

func (m *MockUserRepository) Exists(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

// MockURLRepository is a mock implementation of URLRepository
type MockURLRepository struct {
	mock.Mock
}

func (m *MockURLRepository) Create(ctx context.Context, url *domain.ShortURL) error {
	args := m.Called(ctx, url)
	return args.Error(0)
}

func (m *MockURLRepository) GetByShortCode(ctx context.Context, shortCode string) (*domain.ShortURL, error) {
	args := m.Called(ctx, shortCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ShortURL), args.Error(1)
}

func (m *MockURLRepository) GetByID(ctx context.Context, id uint) (*domain.ShortURL, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ShortURL), args.Error(1)
}

func (m *MockURLRepository) Update(ctx context.Context, url *domain.ShortURL) error {
	args := m.Called(ctx, url)
	return args.Error(0)
}

func (m *MockURLRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockURLRepository) GetByUserID(ctx context.Context, userID uint, offset, limit int) ([]*domain.ShortURL, int64, error) {
	args := m.Called(ctx, userID, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.ShortURL), args.Get(1).(int64), args.Error(2)
}

func (m *MockURLRepository) GetTotalURLs(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockURLRepository) GetExpiredURLs(ctx context.Context, limit int) ([]*domain.ShortURL, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ShortURL), args.Error(1)
}

func (m *MockURLRepository) GetPopularURLs(ctx context.Context, limit int) ([]*domain.ShortURL, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ShortURL), args.Error(1)
}

func (m *MockURLRepository) UpdateClickCount(ctx context.Context, id uint, count int64) error {
	args := m.Called(ctx, id, count)
	return args.Error(0)
}

func (m *MockURLRepository) GetURLsByDateRange(ctx context.Context, userID uint, start, end time.Time) ([]*domain.ShortURL, error) {
	args := m.Called(ctx, userID, start, end)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ShortURL), args.Error(1)
}

func (m *MockURLRepository) BulkCreate(ctx context.Context, urls []*domain.ShortURL) error {
	args := m.Called(ctx, urls)
	return args.Error(0)
}

func (m *MockURLRepository) BulkUpdate(ctx context.Context, urls []*domain.ShortURL) error {
	args := m.Called(ctx, urls)
	return args.Error(0)
}

func (m *MockURLRepository) BulkDelete(ctx context.Context, ids []uint) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

func (m *MockURLRepository) ExistsByShortCode(ctx context.Context, shortCode string) (bool, error) {
	args := m.Called(ctx, shortCode)
	return args.Bool(0), args.Error(1)
}

func (m *MockURLRepository) GetActiveByShortCode(ctx context.Context, shortCode string) (*domain.ShortURL, error) {
	args := m.Called(ctx, shortCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ShortURL), args.Error(1)
}

func (m *MockURLRepository) GetTotalURLsByUser(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockURLRepository) IncrementClickCount(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockClickRepository is a mock implementation of ClickRepository
type MockClickRepository struct {
	mock.Mock
}

func (m *MockClickRepository) Create(ctx context.Context, click *domain.Click) error {
	args := m.Called(ctx, click)
	return args.Error(0)
}

func (m *MockClickRepository) GetByID(ctx context.Context, id uint) (*domain.Click, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Click), args.Error(1)
}

func (m *MockClickRepository) GetByShortURLID(ctx context.Context, shortURLID uint, offset, limit int) ([]*domain.Click, int64, error) {
	args := m.Called(ctx, shortURLID, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Click), args.Get(1).(int64), args.Error(2)
}

func (m *MockClickRepository) GetClickStats(ctx context.Context, shortURLID uint, period string) (*domain.ClickStats, error) {
	args := m.Called(ctx, shortURLID, period)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ClickStats), args.Error(1)
}

func (m *MockClickRepository) GetGeoStats(ctx context.Context, shortURLID uint) (*domain.GeoStats, error) {
	args := m.Called(ctx, shortURLID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.GeoStats), args.Error(1)
}

func (m *MockClickRepository) GetGlobalStats(ctx context.Context) (*domain.GlobalStats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.GlobalStats), args.Error(1)
}

func (m *MockClickRepository) GetTimelineStats(ctx context.Context, shortURLID uint, period string) (*domain.TimelineStats, error) {
	args := m.Called(ctx, shortURLID, period)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TimelineStats), args.Error(1)
}

func (m *MockClickRepository) GetTotalClicks(ctx context.Context, shortURLID uint) (int64, error) {
	args := m.Called(ctx, shortURLID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockClickRepository) GetUniqueClicks(ctx context.Context, shortURLID uint) (int64, error) {
	args := m.Called(ctx, shortURLID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockClickRepository) GetTopCountries(ctx context.Context, shortURLID uint, limit int) ([]domain.CountryStat, error) {
	args := m.Called(ctx, shortURLID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.CountryStat), args.Error(1)
}

func (m *MockClickRepository) GetTopDevices(ctx context.Context, shortURLID uint, limit int) ([]domain.DeviceStat, error) {
	args := m.Called(ctx, shortURLID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.DeviceStat), args.Error(1)
}

func (m *MockClickRepository) GetTopBrowsers(ctx context.Context, shortURLID uint, limit int) ([]domain.BrowserStat, error) {
	args := m.Called(ctx, shortURLID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.BrowserStat), args.Error(1)
}

func (m *MockClickRepository) GetTopReferers(ctx context.Context, shortURLID uint, limit int) ([]domain.RefererStat, error) {
	args := m.Called(ctx, shortURLID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.RefererStat), args.Error(1)
}

func (m *MockClickRepository) GetUserStats(ctx context.Context, userID uint) (*domain.UserAnalytics, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserAnalytics), args.Error(1)
}

func (m *MockClickRepository) GetClicksByDateRange(ctx context.Context, shortURLID uint, startDate, endDate string) ([]*domain.Click, error) {
	args := m.Called(ctx, shortURLID, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Click), args.Error(1)
}

func (m *MockClickRepository) GetClicksByCountry(ctx context.Context, urlID uint) (map[string]int64, error) {
	args := m.Called(ctx, urlID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]int64), args.Error(1)
}

func (m *MockClickRepository) GetClicksByDevice(ctx context.Context, urlID uint) (map[string]int64, error) {
	args := m.Called(ctx, urlID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]int64), args.Error(1)
}

func (m *MockClickRepository) GetClicksByReferrer(ctx context.Context, urlID uint) (map[string]int64, error) {
	args := m.Called(ctx, urlID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]int64), args.Error(1)
}

func (m *MockClickRepository) GetUserClickStats(ctx context.Context, userID uint) (*domain.UserClickStats, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserClickStats), args.Error(1)
}

func (m *MockClickRepository) GetRecentClicks(ctx context.Context, shortURLID uint, limit int) ([]domain.RecentClickStat, error) {
	args := m.Called(ctx, shortURLID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.RecentClickStat), args.Error(1)
}

func (m *MockClickRepository) GetHourlyStats(ctx context.Context, urlID uint, date time.Time) ([]domain.HourlyStats, error) {
	args := m.Called(ctx, urlID, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.HourlyStats), args.Error(1)
}

func (m *MockClickRepository) GetDailyStats(ctx context.Context, urlID uint, start, end time.Time) ([]domain.DailyStats, error) {
	args := m.Called(ctx, urlID, start, end)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.DailyStats), args.Error(1)
}

func (m *MockClickRepository) DeleteByURLID(ctx context.Context, urlID uint) error {
	args := m.Called(ctx, urlID)
	return args.Error(0)
}

// MockCacheService is a mock implementation of CacheService
type MockCacheService struct {
	mock.Mock
}

func (m *MockCacheService) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockCacheService) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	args := m.Called(ctx, key, value, ttl)
	return args.Error(0)
}

func (m *MockCacheService) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockCacheService) GetCounter(ctx context.Context, key string) (int64, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCacheService) IncrementCounter(ctx context.Context, key string, delta int64, ttl time.Duration) error {
	args := m.Called(ctx, key, delta, ttl)
	return args.Error(0)
}

func (m *MockCacheService) GetMultiple(ctx context.Context, keys []string) (map[string]string, error) {
	args := m.Called(ctx, keys)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *MockCacheService) SetMultiple(ctx context.Context, values map[string]string, ttl time.Duration) error {
	args := m.Called(ctx, values, ttl)
	return args.Error(0)
}

func (m *MockCacheService) DeleteMultiple(ctx context.Context, keys []string) error {
	args := m.Called(ctx, keys)
	return args.Error(0)
}

func (m *MockCacheService) Exists(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

func (m *MockCacheService) Clear(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockCacheService) SetSession(ctx context.Context, token string, userID uint, ttl time.Duration) error {
	args := m.Called(ctx, token, userID, ttl)
	return args.Error(0)
}

func (m *MockCacheService) GetSession(ctx context.Context, token string) (uint, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockCacheService) Del(ctx context.Context, keys ...string) error {
	args := m.Called(ctx, keys)
	return args.Error(0)
}

func (m *MockCacheService) TTL(ctx context.Context, key string) (time.Duration, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(time.Duration), args.Error(1)
}

func (m *MockCacheService) Incr(ctx context.Context, key string) (int64, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCacheService) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	args := m.Called(ctx, key, value)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCacheService) SAdd(ctx context.Context, key string, members ...interface{}) error {
	args := m.Called(ctx, key, members)
	return args.Error(0)
}

func (m *MockCacheService) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	args := m.Called(ctx, key, member)
	return args.Bool(0), args.Error(1)
}

func (m *MockCacheService) SCard(ctx context.Context, key string) (int64, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCacheService) HSet(ctx context.Context, key string, values ...interface{}) error {
	args := m.Called(ctx, key, values)
	return args.Error(0)
}

func (m *MockCacheService) HGet(ctx context.Context, key, field string) (string, error) {
	args := m.Called(ctx, key, field)
	return args.String(0), args.Error(1)
}

func (m *MockCacheService) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *MockCacheService) HDel(ctx context.Context, key string, fields ...string) error {
	args := m.Called(ctx, key, fields)
	return args.Error(0)
}

func (m *MockCacheService) CacheURL(ctx context.Context, shortCode, originalURL string, userID uint, expiration time.Duration) error {
	args := m.Called(ctx, shortCode, originalURL, userID, expiration)
	return args.Error(0)
}

func (m *MockCacheService) GetCachedURL(ctx context.Context, shortCode string) (string, uint, error) {
	args := m.Called(ctx, shortCode)
	return args.String(0), args.Get(1).(uint), args.Error(2)
}

func (m *MockCacheService) InvalidateURL(ctx context.Context, shortCode string) error {
	args := m.Called(ctx, shortCode)
	return args.Error(0)
}

func (m *MockCacheService) IsRateLimited(ctx context.Context, key string, limit int64, window time.Duration) (bool, error) {
	args := m.Called(ctx, key, limit, window)
	return args.Bool(0), args.Error(1)
}

func (m *MockCacheService) IncrementRateLimit(ctx context.Context, key string, window time.Duration) (int64, error) {
	args := m.Called(ctx, key, window)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCacheService) InvalidateSession(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockCacheService) CacheClickCount(ctx context.Context, shortCode string, count int64) error {
	args := m.Called(ctx, shortCode, count)
	return args.Error(0)
}

func (m *MockCacheService) GetClickCount(ctx context.Context, shortCode string) (int64, error) {
	args := m.Called(ctx, shortCode)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCacheService) IncrementClickCount(ctx context.Context, shortCode string) (int64, error) {
	args := m.Called(ctx, shortCode)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCacheService) CacheUniqueClick(ctx context.Context, shortCode, ipAddress string) (bool, error) {
	args := m.Called(ctx, shortCode, ipAddress)
	return args.Bool(0), args.Error(1)
}

func (m *MockCacheService) GetUniqueClickCount(ctx context.Context, shortCode string) (int64, error) {
	args := m.Called(ctx, shortCode)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockCacheService) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockCacheService) FlushDB(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockCacheService) Info(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *MockCacheService) Close() error {
	args := m.Called()
	return args.Error(0)
}

// MockConfigService is a mock implementation of ConfigService
type MockConfigService struct {
	mock.Mock
}

func (m *MockConfigService) GetBaseURL() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfigService) GetJWTSecret() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfigService) GetDatabaseURL() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfigService) GetRedisURL() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfigService) GetServerPort() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfigService) GetServerHost() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfigService) IsProduction() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockConfigService) IsDevelopment() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockConfigService) GetEnvironment() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfigService) GetLogLevel() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfigService) GetCORSOrigins() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *MockConfigService) GetRateLimitRequests() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockConfigService) GetRateLimitBurst() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockConfigService) GetJWTAccessTokenTTL() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockConfigService) GetJWTRefreshTokenTTL() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockConfigService) GetDatabaseMaxOpenConns() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockConfigService) GetDatabaseMaxIdleConns() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockConfigService) GetDatabaseConnMaxLifetime() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockConfigService) GetRedisMaxRetries() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockConfigService) GetRedisRetryDelay() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockConfigService) GetRedisPoolSize() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockConfigService) GetRedisTimeout() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockConfigService) GetExternalAPITimeout() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockConfigService) GetExternalAPIRetries() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockConfigService) GetFileUploadMaxSize() int64 {
	args := m.Called()
	return args.Get(0).(int64)
}

func (m *MockConfigService) GetFileUploadAllowedTypes() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *MockConfigService) GetEmailSMTPHost() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfigService) GetEmailSMTPPort() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockConfigService) GetEmailSMTPUser() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfigService) GetEmailSMTPPassword() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfigService) GetEmailFromAddress() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfigService) GetEmailFromName() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfigService) GetMonitoringEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockConfigService) GetMonitoringPort() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfigService) GetTracingEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockConfigService) GetTracingEndpoint() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfigService) GetCacheDefaultTTL() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockConfigService) GetCacheURLTTL() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockConfigService) GetCacheSessionTTL() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockConfigService) GetCacheAnalyticsTTL() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockConfigService) GetSecurityEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockConfigService) GetSecurityHSTSEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockConfigService) GetSecurityCSPEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockConfigService) GetSecurityXSSProtectionEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockConfigService) GetSecurityContentTypeNosniffEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockConfigService) GetSecurityFrameOptionsEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockConfigService) GetBackupEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockConfigService) GetBackupInterval() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockConfigService) GetBackupRetention() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockConfigService) GetBackupStoragePath() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfigService) GetNotificationEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockConfigService) GetNotificationWebhookURL() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockConfigService) GetFeatureQRCodeEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockConfigService) GetFeatureAnalyticsEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockConfigService) GetFeatureCustomDomainsEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockConfigService) GetFeaturePasswordProtectionEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockConfigService) GetFeatureExpirationEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockConfigService) GetFeatureBulkOperationsEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockConfigService) GetFeatureAPIEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockConfigService) GetFeatureWebhooksEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockConfigService) GetFeatureEmailNotificationsEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockConfigService) GetFeatureGeoLocationEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockConfigService) GetFeatureDeviceTrackingEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockConfigService) GetFeatureUserRegistrationEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockConfigService) GetFeatureGuestModeEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockConfigService) GetFeatureRateLimitingEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockConfigService) GetFeatureSpamDetectionEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockConfigService) GetFeatureAuditLogsEnabled() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockConfigService) GetRawConfig() interface{} {
	args := m.Called()
	return args.Get(0)
}

func (m *MockConfigService) Validate() error {
	args := m.Called()
	return args.Error(0)
}

// MockJWTService is a mock implementation of JWTService
type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateAccessToken(userID uint, email string) (string, error) {
	args := m.Called(userID, email)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) GenerateRefreshToken(userID uint) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateAccessToken(token string) (*domain.TokenClaims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TokenClaims), args.Error(1)
}

func (m *MockJWTService) ValidateRefreshToken(token string) (*domain.TokenClaims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TokenClaims), args.Error(1)
}

func (m *MockJWTService) GetTokenTTL(tokenType string) time.Duration {
	args := m.Called(tokenType)
	return args.Get(0).(time.Duration)
}

func (m *MockJWTService) ExtractTokenFromHeader(authHeader string) (string, error) {
	args := m.Called(authHeader)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) GetUserIDFromToken(tokenString string) (uint, error) {
	args := m.Called(tokenString)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockJWTService) IsTokenExpired(tokenString string) bool {
	args := m.Called(tokenString)
	return args.Bool(0)
}

func (m *MockJWTService) GetTokenClaims(tokenString string) (*domain.TokenClaims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TokenClaims), args.Error(1)
}

func (m *MockJWTService) RevokeToken(tokenString string) error {
	args := m.Called(tokenString)
	return args.Error(0)
}