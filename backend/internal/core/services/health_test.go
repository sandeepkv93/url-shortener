package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"url-shortener/internal/core/domain"
)

func TestHealthService_CheckHealth_AllHealthy(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockURLRepo := new(MockURLRepository)
	mockClickRepo := new(MockClickRepository)
	mockCacheRepo := new(MockCacheService)
	mockConfigRepo := new(MockConfigService)

	// Mock successful database check
	mockUserRepo.On("List", mock.Anything, 0, 1).Return([]*domain.User{}, int64(0), nil)

	// Mock successful cache check
	mockCacheRepo.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Return(nil)
	mockCacheRepo.On("Get", mock.Anything, mock.AnythingOfType("string")).Return("test_value", nil)
	mockCacheRepo.On("Delete", mock.Anything, mock.AnythingOfType("string")).Return(nil)

	service := NewHealthService(mockUserRepo, mockURLRepo, mockClickRepo, mockCacheRepo, mockConfigRepo)

	ctx := context.Background()
	status, err := service.CheckHealth(ctx)

	require.NoError(t, err)
	assert.Equal(t, "healthy", status.Status)
	assert.Equal(t, "1.0.0", status.Version)
	assert.NotEmpty(t, status.Timestamp)
	assert.NotEmpty(t, status.Checks)

	// Check database health
	dbHealth, exists := status.Checks["database"]
	assert.True(t, exists)
	assert.Equal(t, "healthy", dbHealth.Status)
	assert.Equal(t, "Database connection successful", dbHealth.Message)

	// Check cache health
	cacheHealth, exists := status.Checks["cache"]
	assert.True(t, exists)
	assert.Equal(t, "healthy", cacheHealth.Status)
	assert.Equal(t, "Cache operations successful", cacheHealth.Message)

	// Check external services
	geoHealth, exists := status.Checks["geolocation"]
	assert.True(t, exists)
	assert.Equal(t, "healthy", geoHealth.Status)

	emailHealth, exists := status.Checks["email"]
	assert.True(t, exists)
	assert.Equal(t, "healthy", emailHealth.Status)

	qrHealth, exists := status.Checks["qr_code"]
	assert.True(t, exists)
	assert.Equal(t, "healthy", qrHealth.Status)

	mockUserRepo.AssertExpectations(t)
	mockCacheRepo.AssertExpectations(t)
}

func TestHealthService_CheckHealth_DatabaseUnhealthy(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockURLRepo := new(MockURLRepository)
	mockClickRepo := new(MockClickRepository)
	mockCacheRepo := new(MockCacheService)
	mockConfigRepo := new(MockConfigService)

	// Mock failed database check
	mockUserRepo.On("List", mock.Anything, 0, 1).Return(nil, int64(0), errors.New("database connection failed"))

	// Mock successful cache check
	mockCacheRepo.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Return(nil)
	mockCacheRepo.On("Get", mock.Anything, mock.AnythingOfType("string")).Return("test_value", nil)
	mockCacheRepo.On("Delete", mock.Anything, mock.AnythingOfType("string")).Return(nil)

	service := NewHealthService(mockUserRepo, mockURLRepo, mockClickRepo, mockCacheRepo, mockConfigRepo)

	ctx := context.Background()
	status, err := service.CheckHealth(ctx)

	require.NoError(t, err)
	assert.Equal(t, "unhealthy", status.Status)

	// Check database health
	dbHealth, exists := status.Checks["database"]
	assert.True(t, exists)
	assert.Equal(t, "unhealthy", dbHealth.Status)
	assert.Contains(t, dbHealth.Message, "Database connection failed")

	mockUserRepo.AssertExpectations(t)
	mockCacheRepo.AssertExpectations(t)
}

func TestHealthService_CheckHealth_CacheUnhealthy(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockURLRepo := new(MockURLRepository)
	mockClickRepo := new(MockClickRepository)
	mockCacheRepo := new(MockCacheService)
	mockConfigRepo := new(MockConfigService)

	// Mock successful database check
	mockUserRepo.On("List", mock.Anything, 0, 1).Return([]*domain.User{}, int64(0), nil)

	// Mock failed cache check
	mockCacheRepo.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("time.Duration")).Return(errors.New("cache write failed"))

	service := NewHealthService(mockUserRepo, mockURLRepo, mockClickRepo, mockCacheRepo, mockConfigRepo)

	ctx := context.Background()
	status, err := service.CheckHealth(ctx)

	require.NoError(t, err)
	assert.Equal(t, "degraded", status.Status)

	// Check cache health
	cacheHealth, exists := status.Checks["cache"]
	assert.True(t, exists)
	assert.Equal(t, "unhealthy", cacheHealth.Status)
	assert.Contains(t, cacheHealth.Message, "Cache write failed")

	mockUserRepo.AssertExpectations(t)
	mockCacheRepo.AssertExpectations(t)
}

func TestHealthService_CheckDatabaseHealth_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockURLRepo := new(MockURLRepository)
	mockClickRepo := new(MockClickRepository)
	mockCacheRepo := new(MockCacheService)
	mockConfigRepo := new(MockConfigService)

	mockUserRepo.On("List", mock.Anything, 0, 1).Return([]*domain.User{}, int64(0), nil)

	service := NewHealthService(mockUserRepo, mockURLRepo, mockClickRepo, mockCacheRepo, mockConfigRepo)

	ctx := context.Background()
	health, err := service.CheckDatabaseHealth(ctx)

	require.NoError(t, err)
	assert.Equal(t, "healthy", health.Status)
	assert.Equal(t, "Database connection successful", health.Message)
	assert.NotEmpty(t, health.LastChecked)
	assert.Greater(t, health.Duration, time.Duration(0))
	assert.Contains(t, health.Metadata, "response_time_ms")

	mockUserRepo.AssertExpectations(t)
}

func TestHealthService_CheckDatabaseHealth_Failure(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockURLRepo := new(MockURLRepository)
	mockClickRepo := new(MockClickRepository)
	mockCacheRepo := new(MockCacheService)
	mockConfigRepo := new(MockConfigService)

	mockUserRepo.On("List", mock.Anything, 0, 1).Return(nil, int64(0), errors.New("connection timeout"))

	service := NewHealthService(mockUserRepo, mockURLRepo, mockClickRepo, mockCacheRepo, mockConfigRepo)

	ctx := context.Background()
	health, err := service.CheckDatabaseHealth(ctx)

	assert.Error(t, err)
	assert.Equal(t, "unhealthy", health.Status)
	assert.Contains(t, health.Message, "Database connection failed")
	assert.NotEmpty(t, health.LastChecked)
	assert.Greater(t, health.Duration, time.Duration(0))

	mockUserRepo.AssertExpectations(t)
}

func TestHealthService_CheckCacheHealth_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockURLRepo := new(MockURLRepository)
	mockClickRepo := new(MockClickRepository)
	mockCacheRepo := new(MockCacheService)
	mockConfigRepo := new(MockConfigService)

	mockCacheRepo.On("Set", mock.Anything, mock.AnythingOfType("string"), "test_value", mock.AnythingOfType("time.Duration")).Return(nil)
	mockCacheRepo.On("Get", mock.Anything, mock.AnythingOfType("string")).Return("test_value", nil)
	mockCacheRepo.On("Delete", mock.Anything, mock.AnythingOfType("string")).Return(nil)

	service := NewHealthService(mockUserRepo, mockURLRepo, mockClickRepo, mockCacheRepo, mockConfigRepo)

	ctx := context.Background()
	health, err := service.CheckCacheHealth(ctx)

	require.NoError(t, err)
	assert.Equal(t, "healthy", health.Status)
	assert.Equal(t, "Cache operations successful", health.Message)
	assert.NotEmpty(t, health.LastChecked)
	assert.Greater(t, health.Duration, time.Duration(0))
	assert.Contains(t, health.Metadata, "response_time_ms")

	mockCacheRepo.AssertExpectations(t)
}

func TestHealthService_CheckCacheHealth_SetFailure(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockURLRepo := new(MockURLRepository)
	mockClickRepo := new(MockClickRepository)
	mockCacheRepo := new(MockCacheService)
	mockConfigRepo := new(MockConfigService)

	mockCacheRepo.On("Set", mock.Anything, mock.AnythingOfType("string"), "test_value", mock.AnythingOfType("time.Duration")).Return(errors.New("cache write failed"))

	service := NewHealthService(mockUserRepo, mockURLRepo, mockClickRepo, mockCacheRepo, mockConfigRepo)

	ctx := context.Background()
	health, err := service.CheckCacheHealth(ctx)

	assert.Error(t, err)
	assert.Equal(t, "unhealthy", health.Status)
	assert.Contains(t, health.Message, "Cache write failed")

	mockCacheRepo.AssertExpectations(t)
}

func TestHealthService_CheckCacheHealth_GetFailure(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockURLRepo := new(MockURLRepository)
	mockClickRepo := new(MockClickRepository)
	mockCacheRepo := new(MockCacheService)
	mockConfigRepo := new(MockConfigService)

	mockCacheRepo.On("Set", mock.Anything, mock.AnythingOfType("string"), "test_value", mock.AnythingOfType("time.Duration")).Return(nil)
	mockCacheRepo.On("Get", mock.Anything, mock.AnythingOfType("string")).Return("", errors.New("cache read failed"))

	service := NewHealthService(mockUserRepo, mockURLRepo, mockClickRepo, mockCacheRepo, mockConfigRepo)

	ctx := context.Background()
	health, err := service.CheckCacheHealth(ctx)

	assert.Error(t, err)
	assert.Equal(t, "unhealthy", health.Status)
	assert.Contains(t, health.Message, "Cache read failed")

	mockCacheRepo.AssertExpectations(t)
}

func TestHealthService_CheckCacheHealth_ValueMismatch(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockURLRepo := new(MockURLRepository)
	mockClickRepo := new(MockClickRepository)
	mockCacheRepo := new(MockCacheService)
	mockConfigRepo := new(MockConfigService)

	mockCacheRepo.On("Set", mock.Anything, mock.AnythingOfType("string"), "test_value", mock.AnythingOfType("time.Duration")).Return(nil)
	mockCacheRepo.On("Get", mock.Anything, mock.AnythingOfType("string")).Return("wrong_value", nil)

	service := NewHealthService(mockUserRepo, mockURLRepo, mockClickRepo, mockCacheRepo, mockConfigRepo)

	ctx := context.Background()
	health, err := service.CheckCacheHealth(ctx)

	assert.Error(t, err)
	assert.Equal(t, "unhealthy", health.Status)
	assert.Equal(t, "Cache value mismatch", health.Message)

	mockCacheRepo.AssertExpectations(t)
}

func TestHealthService_GetSystemMetrics(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockURLRepo := new(MockURLRepository)
	mockClickRepo := new(MockClickRepository)
	mockCacheRepo := new(MockCacheService)
	mockConfigRepo := new(MockConfigService)

	service := NewHealthService(mockUserRepo, mockURLRepo, mockClickRepo, mockCacheRepo, mockConfigRepo)

	ctx := context.Background()
	metrics, err := service.GetSystemMetrics(ctx)

	require.NoError(t, err)
	assert.NotEmpty(t, metrics.Timestamp)
	assert.Greater(t, metrics.Memory.Allocated, uint64(0))
	assert.Greater(t, metrics.Memory.System, uint64(0))
	assert.Greater(t, metrics.CPU.Goroutines, 0)
	assert.Greater(t, metrics.CPU.NumCPU, 0)
	assert.NotEmpty(t, metrics.Runtime.Version)
	assert.NotEmpty(t, metrics.Runtime.GOOS)
	assert.NotEmpty(t, metrics.Runtime.GOARCH)
}

func TestHealthService_GetApplicationMetrics(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockURLRepo := new(MockURLRepository)
	mockClickRepo := new(MockClickRepository)
	mockCacheRepo := new(MockCacheService)
	mockConfigRepo := new(MockConfigService)

	// Mock user count
	mockUserRepo.On("List", mock.Anything, 0, 1).Return([]*domain.User{}, int64(100), nil)

	// Mock URL count
	mockURLRepo.On("GetTotalURLs", mock.Anything).Return(int64(500), nil)

	// Mock expired URLs
	mockURLRepo.On("GetExpiredURLs", mock.Anything, 1000).Return([]*domain.ShortURL{}, nil)

	service := NewHealthService(mockUserRepo, mockURLRepo, mockClickRepo, mockCacheRepo, mockConfigRepo)

	ctx := context.Background()
	metrics, err := service.GetApplicationMetrics(ctx)

	require.NoError(t, err)
	assert.NotEmpty(t, metrics.Timestamp)
	assert.Equal(t, int64(100), metrics.Database.TotalUsers)
	assert.Equal(t, int64(500), metrics.Database.TotalURLs)
	assert.Equal(t, int64(0), metrics.Database.ExpiredURLs)

	// Test cache metrics (these return default values)
	assert.Equal(t, 0.85, metrics.Cache.HitRate)
	assert.Equal(t, 0.15, metrics.Cache.MissRate)

	mockUserRepo.AssertExpectations(t)
	mockURLRepo.AssertExpectations(t)
}

func TestHealthService_GetApplicationMetrics_UserCountError(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockURLRepo := new(MockURLRepository)
	mockClickRepo := new(MockClickRepository)
	mockCacheRepo := new(MockCacheService)
	mockConfigRepo := new(MockConfigService)

	// Mock user count error
	mockUserRepo.On("List", mock.Anything, 0, 1).Return(nil, int64(0), errors.New("database error"))

	service := NewHealthService(mockUserRepo, mockURLRepo, mockClickRepo, mockCacheRepo, mockConfigRepo)

	ctx := context.Background()
	_, err := service.GetApplicationMetrics(ctx)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get user count")

	mockUserRepo.AssertExpectations(t)
}

func TestHealthService_GetApplicationMetrics_URLCountError(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockURLRepo := new(MockURLRepository)
	mockClickRepo := new(MockClickRepository)
	mockCacheRepo := new(MockCacheService)
	mockConfigRepo := new(MockConfigService)

	// Mock user count success
	mockUserRepo.On("List", mock.Anything, 0, 1).Return([]*domain.User{}, int64(100), nil)

	// Mock URL count error
	mockURLRepo.On("GetTotalURLs", mock.Anything).Return(int64(0), errors.New("database error"))

	service := NewHealthService(mockUserRepo, mockURLRepo, mockClickRepo, mockCacheRepo, mockConfigRepo)

	ctx := context.Background()
	_, err := service.GetApplicationMetrics(ctx)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get URL count")

	mockUserRepo.AssertExpectations(t)
	mockURLRepo.AssertExpectations(t)
}

func TestHealthService_CheckExternalServices(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockURLRepo := new(MockURLRepository)
	mockClickRepo := new(MockClickRepository)
	mockCacheRepo := new(MockCacheService)
	mockConfigRepo := new(MockConfigService)

	service := NewHealthService(mockUserRepo, mockURLRepo, mockClickRepo, mockCacheRepo, mockConfigRepo)

	ctx := context.Background()
	services, err := service.CheckExternalServices(ctx)

	require.NoError(t, err)
	assert.Len(t, services, 3)

	// Check geolocation service
	geoService, exists := services["geolocation"]
	assert.True(t, exists)
	assert.Equal(t, "healthy", geoService.Status)
	assert.Equal(t, "Geolocation service available", geoService.Message)
	assert.Equal(t, "ipapi", geoService.Metadata["provider"])

	// Check email service
	emailService, exists := services["email"]
	assert.True(t, exists)
	assert.Equal(t, "healthy", emailService.Status)
	assert.Equal(t, "Email service available", emailService.Message)
	assert.Equal(t, "smtp", emailService.Metadata["provider"])

	// Check QR code service
	qrService, exists := services["qr_code"]
	assert.True(t, exists)
	assert.Equal(t, "healthy", qrService.Status)
	assert.Equal(t, "QR code service available", qrService.Message)
	assert.Equal(t, "internal", qrService.Metadata["provider"])
}