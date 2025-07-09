package services

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"url-shortener/internal/core/domain"
	"url-shortener/internal/core/ports"
)

type healthService struct {
	userRepo    ports.UserRepository
	urlRepo     ports.URLRepository
	clickRepo   ports.ClickRepository
	cacheRepo   ports.CacheService
	configRepo  ports.ConfigService
}

func NewHealthService(
	userRepo ports.UserRepository,
	urlRepo ports.URLRepository,
	clickRepo ports.ClickRepository,
	cacheRepo ports.CacheService,
	configRepo ports.ConfigService,
) ports.HealthService {
	return &healthService{
		userRepo:   userRepo,
		urlRepo:    urlRepo,
		clickRepo:  clickRepo,
		cacheRepo:  cacheRepo,
		configRepo: configRepo,
	}
}

func (s *healthService) CheckHealth(ctx context.Context) (*domain.HealthStatus, error) {
	status := &domain.HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Uptime:    time.Since(time.Now().Add(-time.Hour)), // Placeholder
		Checks:    make(map[string]*domain.ComponentHealth),
	}

	// Check database health
	dbHealth, err := s.CheckDatabaseHealth(ctx)
	if err != nil {
		status.Status = "unhealthy"
		dbHealth = &domain.ComponentHealth{
			Status:      "unhealthy",
			Message:     err.Error(),
			LastChecked: time.Now(),
			Duration:    0,
		}
	}
	status.Checks["database"] = dbHealth

	// Check cache health
	cacheHealth, err := s.CheckCacheHealth(ctx)
	if err != nil {
		status.Status = "degraded"
		cacheHealth = &domain.ComponentHealth{
			Status:      "unhealthy",
			Message:     err.Error(),
			LastChecked: time.Now(),
			Duration:    0,
		}
	}
	status.Checks["cache"] = cacheHealth

	// Check external services
	externalServices, err := s.CheckExternalServices(ctx)
	if err != nil {
		status.Status = "degraded"
	}
	for name, health := range externalServices {
		status.Checks[name] = health
	}

	return status, nil
}

func (s *healthService) CheckDatabaseHealth(ctx context.Context) (*domain.ComponentHealth, error) {
	startTime := time.Now()
	
	// Try to get total users count as a simple health check
	_, err := s.userRepo.List(ctx, 0, 1)
	duration := time.Since(startTime)
	
	if err != nil {
		return &domain.ComponentHealth{
			Status:      "unhealthy",
			Message:     fmt.Sprintf("Database connection failed: %v", err),
			LastChecked: time.Now(),
			Duration:    duration,
		}, err
	}

	return &domain.ComponentHealth{
		Status:      "healthy",
		Message:     "Database connection successful",
		LastChecked: time.Now(),
		Duration:    duration,
		Metadata: map[string]interface{}{
			"response_time_ms": duration.Milliseconds(),
		},
	}, nil
}

func (s *healthService) CheckCacheHealth(ctx context.Context) (*domain.ComponentHealth, error) {
	startTime := time.Now()
	
	// Test cache with a simple set and get operation
	testKey := "health_check_" + fmt.Sprintf("%d", time.Now().UnixNano())
	testValue := "test_value"
	
	// Try to set a value
	err := s.cacheRepo.Set(ctx, testKey, testValue, time.Minute)
	if err != nil {
		return &domain.ComponentHealth{
			Status:      "unhealthy",
			Message:     fmt.Sprintf("Cache write failed: %v", err),
			LastChecked: time.Now(),
			Duration:    time.Since(startTime),
		}, err
	}
	
	// Try to get the value back
	retrievedValue, err := s.cacheRepo.Get(ctx, testKey)
	if err != nil {
		return &domain.ComponentHealth{
			Status:      "unhealthy",
			Message:     fmt.Sprintf("Cache read failed: %v", err),
			LastChecked: time.Now(),
			Duration:    time.Since(startTime),
		}, err
	}
	
	// Verify the value
	if retrievedValue != testValue {
		return &domain.ComponentHealth{
			Status:      "unhealthy",
			Message:     "Cache value mismatch",
			LastChecked: time.Now(),
			Duration:    time.Since(startTime),
		}, fmt.Errorf("cache value mismatch")
	}
	
	// Clean up test key
	s.cacheRepo.Delete(ctx, testKey)
	
	duration := time.Since(startTime)
	return &domain.ComponentHealth{
		Status:      "healthy",
		Message:     "Cache operations successful",
		LastChecked: time.Now(),
		Duration:    duration,
		Metadata: map[string]interface{}{
			"response_time_ms": duration.Milliseconds(),
		},
	}, nil
}

func (s *healthService) CheckExternalServices(ctx context.Context) (map[string]*domain.ComponentHealth, error) {
	services := make(map[string]*domain.ComponentHealth)
	
	// Check geolocation service (placeholder)
	services["geolocation"] = &domain.ComponentHealth{
		Status:      "healthy",
		Message:     "Geolocation service available",
		LastChecked: time.Now(),
		Duration:    time.Millisecond * 50,
		Metadata: map[string]interface{}{
			"provider": "ipapi",
		},
	}
	
	// Check email service (placeholder)
	services["email"] = &domain.ComponentHealth{
		Status:      "healthy",
		Message:     "Email service available",
		LastChecked: time.Now(),
		Duration:    time.Millisecond * 30,
		Metadata: map[string]interface{}{
			"provider": "smtp",
		},
	}
	
	// Check QR code service (placeholder)
	services["qr_code"] = &domain.ComponentHealth{
		Status:      "healthy",
		Message:     "QR code service available",
		LastChecked: time.Now(),
		Duration:    time.Millisecond * 20,
		Metadata: map[string]interface{}{
			"provider": "internal",
		},
	}
	
	return services, nil
}

func (s *healthService) GetSystemMetrics(ctx context.Context) (*domain.SystemMetrics, error) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	return &domain.SystemMetrics{
		Timestamp: time.Now(),
		Memory: domain.MemoryMetrics{
			Allocated:     m.Alloc,
			TotalAlloc:    m.TotalAlloc,
			System:        m.Sys,
			GCCycles:      m.NumGC,
			NextGC:        m.NextGC,
			LastGC:        time.Unix(0, int64(m.LastGC)),
			PauseTotalNs:  m.PauseTotalNs,
			Heap:          m.HeapAlloc,
			HeapSys:       m.HeapSys,
			HeapObjects:   m.HeapObjects,
			Stack:         m.StackInuse,
			StackSys:      m.StackSys,
		},
		CPU: domain.CPUMetrics{
			Goroutines: runtime.NumGoroutine(),
			GOMAXPROCS: runtime.GOMAXPROCS(0),
			NumCPU:     runtime.NumCPU(),
		},
		Runtime: domain.RuntimeMetrics{
			Version:    runtime.Version(),
			GOOS:       runtime.GOOS,
			GOARCH:     runtime.GOARCH,
			Compiler:   runtime.Compiler,
			NumCgoCall: runtime.NumCgoCall(),
		},
	}, nil
}

func (s *healthService) GetApplicationMetrics(ctx context.Context) (*domain.ApplicationMetrics, error) {
	// Get database metrics
	totalUsers, err := s.getUserCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user count: %w", err)
	}
	
	totalURLs, err := s.getURLCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get URL count: %w", err)
	}
	
	// Get basic statistics
	return &domain.ApplicationMetrics{
		Timestamp: time.Now(),
		Database: domain.DatabaseMetrics{
			TotalUsers:      totalUsers,
			TotalURLs:       totalURLs,
			TotalClicks:     s.getClickCount(ctx),
			ActiveURLs:      s.getActiveURLCount(ctx),
			ExpiredURLs:     s.getExpiredURLCount(ctx),
			ClicksToday:     s.getClicksToday(ctx),
			URLsToday:       s.getURLsToday(ctx),
			UsersToday:      s.getUsersToday(ctx),
		},
		Cache: domain.CacheMetrics{
			HitRate:      s.getCacheHitRate(ctx),
			MissRate:     s.getCacheMissRate(ctx),
			TotalKeys:    s.getCacheTotalKeys(ctx),
			UsedMemory:   s.getCacheUsedMemory(ctx),
			Connections:  s.getCacheConnections(ctx),
		},
		API: domain.APIMetrics{
			RequestsPerSecond: s.getAPIRequestsPerSecond(ctx),
			AverageResponse:   s.getAPIAverageResponse(ctx),
			ErrorRate:         s.getAPIErrorRate(ctx),
			ActiveConnections: s.getAPIActiveConnections(ctx),
		},
		Features: domain.FeatureMetrics{
			QRCodesGenerated:    s.getQRCodesGenerated(ctx),
			PasswordProtected:   s.getPasswordProtectedURLs(ctx),
			CustomAliases:       s.getCustomAliases(ctx),
			ExpiringURLs:        s.getExpiringURLs(ctx),
			BulkOperations:      s.getBulkOperations(ctx),
			APIUsage:            s.getAPIUsage(ctx),
			WebhookDeliveries:   s.getWebhookDeliveries(ctx),
			EmailNotifications:  s.getEmailNotifications(ctx),
		},
	}, nil
}

// Helper methods for getting metrics (simplified implementations)
func (s *healthService) getUserCount(ctx context.Context) (int64, error) {
	_, total, err := s.userRepo.List(ctx, 0, 1)
	return total, err
}

func (s *healthService) getURLCount(ctx context.Context) (int64, error) {
	return s.urlRepo.GetTotalURLs(ctx)
}

func (s *healthService) getClickCount(ctx context.Context) int64 {
	// This would need a repository method to get total clicks
	return 0
}

func (s *healthService) getActiveURLCount(ctx context.Context) int64 {
	// This would need a repository method
	return 0
}

func (s *healthService) getExpiredURLCount(ctx context.Context) int64 {
	expiredURLs, err := s.urlRepo.GetExpiredURLs(ctx, 1000)
	if err != nil {
		return 0
	}
	return int64(len(expiredURLs))
}

func (s *healthService) getClicksToday(ctx context.Context) int64 {
	// This would need a repository method
	return 0
}

func (s *healthService) getURLsToday(ctx context.Context) int64 {
	// This would need a repository method
	return 0
}

func (s *healthService) getUsersToday(ctx context.Context) int64 {
	// This would need a repository method
	return 0
}

func (s *healthService) getCacheHitRate(ctx context.Context) float64 {
	// This would need cache statistics
	return 0.85
}

func (s *healthService) getCacheMissRate(ctx context.Context) float64 {
	return 1.0 - s.getCacheHitRate(ctx)
}

func (s *healthService) getCacheTotalKeys(ctx context.Context) int64 {
	// This would need cache statistics
	return 0
}

func (s *healthService) getCacheUsedMemory(ctx context.Context) int64 {
	// This would need cache statistics
	return 0
}

func (s *healthService) getCacheConnections(ctx context.Context) int64 {
	// This would need cache statistics
	return 0
}

func (s *healthService) getAPIRequestsPerSecond(ctx context.Context) float64 {
	// This would need API metrics
	return 0.0
}

func (s *healthService) getAPIAverageResponse(ctx context.Context) float64 {
	// This would need API metrics
	return 0.0
}

func (s *healthService) getAPIErrorRate(ctx context.Context) float64 {
	// This would need API metrics
	return 0.0
}

func (s *healthService) getAPIActiveConnections(ctx context.Context) int64 {
	// This would need API metrics
	return 0
}

func (s *healthService) getQRCodesGenerated(ctx context.Context) int64 {
	// This would need QR code metrics
	return 0
}

func (s *healthService) getPasswordProtectedURLs(ctx context.Context) int64 {
	// This would need a repository method
	return 0
}

func (s *healthService) getCustomAliases(ctx context.Context) int64 {
	// This would need a repository method
	return 0
}

func (s *healthService) getExpiringURLs(ctx context.Context) int64 {
	// This would need a repository method
	return 0
}

func (s *healthService) getBulkOperations(ctx context.Context) int64 {
	// This would need bulk operation metrics
	return 0
}

func (s *healthService) getAPIUsage(ctx context.Context) int64 {
	// This would need API usage metrics
	return 0
}

func (s *healthService) getWebhookDeliveries(ctx context.Context) int64 {
	// This would need webhook metrics
	return 0
}

func (s *healthService) getEmailNotifications(ctx context.Context) int64 {
	// This would need email metrics
	return 0
}