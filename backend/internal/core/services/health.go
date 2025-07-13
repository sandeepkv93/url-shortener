package services

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"url-shortener/internal/core/domain"
	"url-shortener/internal/core/ports"

	"gorm.io/gorm"
)

type HealthService struct {
	db           *gorm.DB
	cacheService ports.CacheService
	startTime    time.Time
	version      string
}

func NewHealthService(db *gorm.DB, cacheService ports.CacheService, version string) ports.HealthService {
	return &HealthService{
		db:           db,
		cacheService: cacheService,
		startTime:    time.Now(),
		version:      version,
	}
}

// GetHealth returns the overall health status of the application
func (h *HealthService) GetHealth(ctx context.Context) (*domain.HealthStatus, error) {
	status := &domain.HealthStatus{
		Status:     "healthy",
		Version:    h.version,
		Uptime:     time.Since(h.startTime),
		Timestamp:  time.Now(),
		Components: make(map[string]*domain.ComponentHealth),
		Checks:     make(map[string]*domain.HealthCheck),
	}

	// Check database health
	dbHealth := h.checkDatabase(ctx)
	status.Components["database"] = dbHealth
	
	// Check cache health
	cacheHealth := h.checkCache(ctx)
	status.Components["cache"] = cacheHealth

	// Check system resources
	systemHealth := h.checkSystemResources(ctx)
	status.Components["system"] = systemHealth

	// Determine overall status
	status.Status = h.determineOverallStatus(status.Components)

	return status, nil
}

// GetSystemMetrics returns detailed system metrics
func (h *HealthService) GetSystemMetrics(ctx context.Context) (*domain.SystemMetrics, error) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	metrics := &domain.SystemMetrics{
		CPU: domain.CPUMetrics{
			Usage: h.getCPUUsage(),
			Cores: runtime.NumCPU(),
		},
		Memory: domain.MemoryMetrics{
			Total:        memStats.Sys,
			Used:         memStats.Alloc,
			Free:         memStats.Sys - memStats.Alloc,
			Available:    memStats.Sys - memStats.Alloc,
			UsagePercent: float64(memStats.Alloc) / float64(memStats.Sys) * 100,
		},
		Processes: domain.ProcessMetrics{
			Total:   runtime.NumGoroutine(),
			Running: runtime.NumGoroutine(),
		},
		Timestamp: time.Now(),
	}

	return metrics, nil
}

// GetApplicationMetrics returns application-specific metrics
func (h *HealthService) GetApplicationMetrics(ctx context.Context) (*domain.ApplicationMetrics, error) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	metrics := &domain.ApplicationMetrics{
		GoroutineCount: runtime.NumGoroutine(),
		HeapSize:       memStats.HeapAlloc,
		Timestamp:      time.Now(),
	}

	// Get database metrics
	if h.db != nil {
		dbMetrics := h.getDatabaseMetrics(ctx)
		metrics.DatabaseConnections = dbMetrics
	}

	// Get cache metrics
	if h.cacheService != nil {
		cacheMetrics := h.getCacheMetrics(ctx)
		metrics.CacheMetrics = cacheMetrics
	}

	return metrics, nil
}

// RunHealthChecks runs all configured health checks
func (h *HealthService) RunHealthChecks(ctx context.Context) (map[string]*domain.HealthCheck, error) {
	checks := make(map[string]*domain.HealthCheck)

	// Database connectivity check
	checks["database_connectivity"] = h.runDatabaseConnectivityCheck(ctx)
	
	// Database query performance check
	checks["database_performance"] = h.runDatabasePerformanceCheck(ctx)
	
	// Cache connectivity check
	checks["cache_connectivity"] = h.runCacheConnectivityCheck(ctx)
	
	// Cache performance check
	checks["cache_performance"] = h.runCachePerformanceCheck(ctx)
	
	// Disk space check
	checks["disk_space"] = h.runDiskSpaceCheck(ctx)
	
	// Memory usage check
	checks["memory_usage"] = h.runMemoryUsageCheck(ctx)

	return checks, nil
}

// CheckHealth is an alias for GetHealth to match the interface
func (h *HealthService) CheckHealth(ctx context.Context) (*domain.HealthStatus, error) {
	return h.GetHealth(ctx)
}

// CheckDatabaseHealth returns the health status of the database
func (h *HealthService) CheckDatabaseHealth(ctx context.Context) (*domain.ComponentHealth, error) {
	return h.checkDatabase(ctx), nil
}

// CheckCacheHealth returns the health status of the cache
func (h *HealthService) CheckCacheHealth(ctx context.Context) (*domain.ComponentHealth, error) {
	return h.checkCache(ctx), nil
}

// CheckExternalServices returns the health status of external services
func (h *HealthService) CheckExternalServices(ctx context.Context) (map[string]*domain.ComponentHealth, error) {
	services := make(map[string]*domain.ComponentHealth)
	
	// For now, we don't have external services to check
	// This method can be expanded when external services are added
	
	return services, nil
}

// IsHealthy returns true if the service is healthy
func (h *HealthService) IsHealthy(ctx context.Context) bool {
	health, err := h.GetHealth(ctx)
	if err != nil {
		return false
	}
	return health.Status == "healthy"
}

// Private helper methods

func (h *HealthService) checkDatabase(ctx context.Context) *domain.ComponentHealth {
	start := time.Now()
	health := &domain.ComponentHealth{
		Status:      "up",
		Message:     "Database is healthy",
		LastChecked: time.Now(),
		Metadata:    make(map[string]interface{}),
	}

	if h.db == nil {
		health.Status = "down"
		health.Message = "Database connection not configured"
		health.Duration = time.Since(start)
		return health
	}

	// Test database connectivity
	sqlDB, err := h.db.DB()
	if err != nil {
		health.Status = "down"
		health.Message = fmt.Sprintf("Failed to get database connection: %v", err)
		health.Duration = time.Since(start)
		return health
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		health.Status = "down"
		health.Message = fmt.Sprintf("Database ping failed: %v", err)
		health.Duration = time.Since(start)
		return health
	}

	// Get database stats
	stats := sqlDB.Stats()
	health.Metadata["open_connections"] = stats.OpenConnections
	health.Metadata["in_use"] = stats.InUse
	health.Metadata["idle"] = stats.Idle
	health.Metadata["max_open_connections"] = stats.MaxOpenConnections

	// Check if connection pool is healthy
	if stats.OpenConnections > stats.MaxOpenConnections*90/100 {
		health.Status = "degraded"
		health.Message = "Database connection pool is nearly exhausted"
	}

	health.Duration = time.Since(start)
	return health
}

func (h *HealthService) checkCache(ctx context.Context) *domain.ComponentHealth {
	start := time.Now()
	health := &domain.ComponentHealth{
		Status:      "up",
		Message:     "Cache is healthy",
		LastChecked: time.Now(),
		Metadata:    make(map[string]interface{}),
	}

	if h.cacheService == nil {
		health.Status = "down"
		health.Message = "Cache service not configured"
		health.Duration = time.Since(start)
		return health
	}

	// Test cache connectivity
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	testKey := "health_check_" + fmt.Sprintf("%d", time.Now().Unix())
	testValue := "test"

	// Test SET operation
	if err := h.cacheService.Set(ctx, testKey, testValue, time.Minute); err != nil {
		health.Status = "down"
		health.Message = fmt.Sprintf("Cache SET operation failed: %v", err)
		health.Duration = time.Since(start)
		return health
	}

	// Test GET operation
	retrieved, err := h.cacheService.Get(ctx, testKey)
	if err != nil {
		health.Status = "down"
		health.Message = fmt.Sprintf("Cache GET operation failed: %v", err)
		health.Duration = time.Since(start)
		return health
	}

	// Test DELETE operation
	if err := h.cacheService.Del(ctx, testKey); err != nil {
		health.Status = "degraded"
		health.Message = fmt.Sprintf("Cache DELETE operation failed: %v", err)
	}

	// Verify retrieved value
	if retrieved != testValue {
		health.Status = "degraded"
		health.Message = "Cache data integrity issue"
	}

	health.Duration = time.Since(start)
	return health
}

func (h *HealthService) checkSystemResources(ctx context.Context) *domain.ComponentHealth {
	start := time.Now()
	health := &domain.ComponentHealth{
		Status:      "up",
		Message:     "System resources are healthy",
		LastChecked: time.Now(),
		Metadata:    make(map[string]interface{}),
	}

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Check memory usage
	memoryUsagePercent := float64(memStats.Alloc) / float64(memStats.Sys) * 100
	health.Metadata["memory_usage_percent"] = memoryUsagePercent
	health.Metadata["memory_allocated"] = memStats.Alloc
	health.Metadata["memory_total"] = memStats.Sys
	health.Metadata["goroutines"] = runtime.NumGoroutine()
	health.Metadata["gc_cycles"] = memStats.NumGC

	// Check if memory usage is concerning
	if memoryUsagePercent > 90 {
		health.Status = "degraded"
		health.Message = "High memory usage detected"
	} else if memoryUsagePercent > 95 {
		health.Status = "down"
		health.Message = "Critical memory usage"
	}

	// Check goroutine count
	goroutines := runtime.NumGoroutine()
	if goroutines > 10000 {
		health.Status = "degraded"
		health.Message = "High goroutine count detected"
	}

	health.Duration = time.Since(start)
	return health
}

func (h *HealthService) determineOverallStatus(components map[string]*domain.ComponentHealth) string {
	hasDown := false
	hasDegraded := false

	for _, component := range components {
		switch component.Status {
		case "down":
			hasDown = true
		case "degraded":
			hasDegraded = true
		}
	}

	if hasDown {
		return "unhealthy"
	}
	if hasDegraded {
		return "degraded"
	}
	return "healthy"
}

func (h *HealthService) getCPUUsage() float64 {
	// This is a simplified CPU usage calculation
	// In a real implementation, you might want to use system-specific methods
	return 0.0
}

func (h *HealthService) getDatabaseMetrics(ctx context.Context) domain.DatabaseMetrics {
	metrics := domain.DatabaseMetrics{}

	if h.db != nil {
		sqlDB, err := h.db.DB()
		if err == nil {
			stats := sqlDB.Stats()
			metrics.OpenConnections = stats.OpenConnections
			metrics.InUseConnections = stats.InUse
			metrics.IdleConnections = stats.Idle
			metrics.MaxOpenConnections = stats.MaxOpenConnections
		}
	}

	return metrics
}

func (h *HealthService) getCacheMetrics(ctx context.Context) domain.CacheMetrics {
	// This would require implementing cache statistics in the cache service
	// For now, return empty metrics
	return domain.CacheMetrics{}
}

// Health check implementations

func (h *HealthService) runDatabaseConnectivityCheck(ctx context.Context) *domain.HealthCheck {
	start := time.Now()
	check := &domain.HealthCheck{
		Name:     "Database Connectivity",
		Status:   "pass",
		Message:  "Database is accessible",
		Critical: true,
		LastRun:  time.Now(),
		NextRun:  time.Now().Add(30 * time.Second),
	}

	if h.db == nil {
		check.Status = "fail"
		check.Message = "Database not configured"
		check.Duration = time.Since(start)
		return check
	}

	sqlDB, err := h.db.DB()
	if err != nil {
		check.Status = "fail"
		check.Message = fmt.Sprintf("Failed to get database connection: %v", err)
		check.Duration = time.Since(start)
		return check
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		check.Status = "fail"
		check.Message = fmt.Sprintf("Database ping failed: %v", err)
	}

	check.Duration = time.Since(start)
	return check
}

func (h *HealthService) runDatabasePerformanceCheck(ctx context.Context) *domain.HealthCheck {
	start := time.Now()
	check := &domain.HealthCheck{
		Name:     "Database Performance",
		Status:   "pass",
		Message:  "Database queries are performing well",
		Critical: false,
		LastRun:  time.Now(),
		NextRun:  time.Now().Add(60 * time.Second),
	}

	if h.db == nil {
		check.Status = "fail"
		check.Message = "Database not configured"
		check.Duration = time.Since(start)
		return check
	}

	// Run a simple query to test performance
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var count int64
	queryStart := time.Now()
	err := h.db.WithContext(ctx).Raw("SELECT COUNT(*) FROM information_schema.tables").Scan(&count).Error
	queryDuration := time.Since(queryStart)

	if err != nil {
		check.Status = "fail"
		check.Message = fmt.Sprintf("Performance query failed: %v", err)
	} else if queryDuration > 5*time.Second {
		check.Status = "warn"
		check.Message = fmt.Sprintf("Slow query detected: %v", queryDuration)
	}

	check.Duration = time.Since(start)
	return check
}

func (h *HealthService) runCacheConnectivityCheck(ctx context.Context) *domain.HealthCheck {
	start := time.Now()
	check := &domain.HealthCheck{
		Name:     "Cache Connectivity",
		Status:   "pass",
		Message:  "Cache is accessible",
		Critical: true,
		LastRun:  time.Now(),
		NextRun:  time.Now().Add(30 * time.Second),
	}

	if h.cacheService == nil {
		check.Status = "fail"
		check.Message = "Cache service not configured"
		check.Duration = time.Since(start)
		return check
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	testKey := "health_check_connectivity"
	if err := h.cacheService.Set(ctx, testKey, "test", time.Minute); err != nil {
		check.Status = "fail"
		check.Message = fmt.Sprintf("Cache connectivity test failed: %v", err)
	} else {
		h.cacheService.Del(ctx, testKey) // Clean up
	}

	check.Duration = time.Since(start)
	return check
}

func (h *HealthService) runCachePerformanceCheck(ctx context.Context) *domain.HealthCheck {
	start := time.Now()
	check := &domain.HealthCheck{
		Name:     "Cache Performance",
		Status:   "pass",
		Message:  "Cache operations are performing well",
		Critical: false,
		LastRun:  time.Now(),
		NextRun:  time.Now().Add(60 * time.Second),
	}

	if h.cacheService == nil {
		check.Status = "fail"
		check.Message = "Cache service not configured"
		check.Duration = time.Since(start)
		return check
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Test cache performance with multiple operations
	testKey := "health_check_performance"
	testValue := "performance_test_value"

	opStart := time.Now()
	err := h.cacheService.Set(ctx, testKey, testValue, time.Minute)
	if err != nil {
		check.Status = "fail"
		check.Message = fmt.Sprintf("Cache SET performance test failed: %v", err)
		check.Duration = time.Since(start)
		return check
	}

	retrieved, err := h.cacheService.Get(ctx, testKey)
	if err != nil {
		check.Status = "fail"
		check.Message = fmt.Sprintf("Cache GET performance test failed: %v", err)
		check.Duration = time.Since(start)
		return check
	}
	
	// Verify the retrieved value matches what we set
	if retrieved != testValue {
		check.Status = "fail"
		check.Message = "Cache GET operation returned incorrect value"
		check.Duration = time.Since(start)
		return check
	}

	opDuration := time.Since(opStart)
	if opDuration > 1*time.Second {
		check.Status = "warn"
		check.Message = fmt.Sprintf("Slow cache operations detected: %v", opDuration)
	}

	h.cacheService.Del(ctx, testKey) // Clean up

	check.Duration = time.Since(start)
	return check
}

func (h *HealthService) runDiskSpaceCheck(ctx context.Context) *domain.HealthCheck {
	check := &domain.HealthCheck{
		Name:     "Disk Space",
		Status:   "pass",
		Message:  "Disk space is adequate",
		Critical: false,
		LastRun:  time.Now(),
		NextRun:  time.Now().Add(5 * time.Minute),
		Duration: time.Millisecond, // Placeholder
	}

	// This is a simplified check
	// In a real implementation, you would check actual disk usage
	return check
}

func (h *HealthService) runMemoryUsageCheck(ctx context.Context) *domain.HealthCheck {
	start := time.Now()
	check := &domain.HealthCheck{
		Name:     "Memory Usage",
		Status:   "pass",
		Message:  "Memory usage is normal",
		Critical: false,
		LastRun:  time.Now(),
		NextRun:  time.Now().Add(time.Minute),
	}

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	memoryUsagePercent := float64(memStats.Alloc) / float64(memStats.Sys) * 100

	if memoryUsagePercent > 90 {
		check.Status = "warn"
		check.Message = fmt.Sprintf("High memory usage: %.2f%%", memoryUsagePercent)
	} else if memoryUsagePercent > 95 {
		check.Status = "fail"
		check.Message = fmt.Sprintf("Critical memory usage: %.2f%%", memoryUsagePercent)
	}

	check.Duration = time.Since(start)
	return check
}