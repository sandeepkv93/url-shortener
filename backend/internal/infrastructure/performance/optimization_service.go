package performance

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"runtime/debug"
	"sync"
	"time"

	"url-shortener/internal/config"
	"url-shortener/internal/core/services"
)

// OptimizationService manages production optimizations and performance monitoring
type OptimizationService struct {
	config        *config.Config
	logger        *services.LoggingService
	
	// Performance monitoring
	requestMetrics   *RequestMetrics
	resourceMetrics  *ResourceMetrics
	optimizations    *config.ProductionOptimization
	
	// Worker pools and buffers
	requestPool      *sync.Pool
	responsePool     *sync.Pool
	bufferPools      map[int]*sync.Pool
	
	// Monitoring
	ctx              context.Context
	cancel           context.CancelFunc
	monitoringTicker *time.Ticker
	
	mu               sync.RWMutex
}

// RequestMetrics tracks request-level performance metrics
type RequestMetrics struct {
	TotalRequests     int64         `json:"total_requests"`
	ActiveRequests    int64         `json:"active_requests"`
	AverageLatency    time.Duration `json:"average_latency"`
	P50Latency        time.Duration `json:"p50_latency"`
	P95Latency        time.Duration `json:"p95_latency"`
	P99Latency        time.Duration `json:"p99_latency"`
	ErrorRate         float64       `json:"error_rate"`
	ThroughputRPS     float64       `json:"throughput_rps"`
	LastUpdated       time.Time     `json:"last_updated"`
	
	latencyHistogram  []time.Duration
	mu                sync.RWMutex
}

// ResourceMetrics tracks system resource usage
type ResourceMetrics struct {
	CPUUsage          float64       `json:"cpu_usage"`
	MemoryUsage       int64         `json:"memory_usage"`
	MemoryPercentage  float64       `json:"memory_percentage"`
	GoroutineCount    int           `json:"goroutine_count"`
	GCPauseTotal      time.Duration `json:"gc_pause_total"`
	GCCycles          uint32        `json:"gc_cycles"`
	HeapSize          uint64        `json:"heap_size"`
	StackSize         uint64        `json:"stack_size"`
	OpenConnections   int           `json:"open_connections"`
	LastUpdated       time.Time     `json:"last_updated"`
}

// PerformanceReport contains comprehensive performance analysis
type PerformanceReport struct {
	Timestamp       time.Time       `json:"timestamp"`
	RequestMetrics  *RequestMetrics `json:"request_metrics"`
	ResourceMetrics *ResourceMetrics `json:"resource_metrics"`
	Optimizations   map[string]interface{} `json:"optimizations_applied"`
	Recommendations []string        `json:"recommendations"`
	HealthStatus    string          `json:"health_status"`
}

// NewOptimizationService creates a new optimization service
func NewOptimizationService(cfg *config.Config, logger *services.LoggingService) *OptimizationService {
	ctx, cancel := context.WithCancel(context.Background())
	
	service := &OptimizationService{
		config:        cfg,
		logger:        logger,
		ctx:           ctx,
		cancel:        cancel,
		optimizations: cfg.GetProductionOptimization(),
		requestMetrics: &RequestMetrics{
			latencyHistogram: make([]time.Duration, 0, 10000),
		},
		resourceMetrics: &ResourceMetrics{},
		bufferPools:     make(map[int]*sync.Pool),
	}
	
	// Initialize object pools
	service.initializeObjectPools()
	
	// Apply optimizations
	if err := service.applyOptimizations(); err != nil {
		logger.Error("Failed to apply optimizations", "error", err)
	}
	
	// Start monitoring
	service.startMonitoring()
	
	return service
}

// initializeObjectPools sets up object pools for better memory management
func (os *OptimizationService) initializeObjectPools() {
	memOpt := os.optimizations.Memory
	
	if memOpt.EnableObjectPools {
		// Request pool
		os.requestPool = &sync.Pool{
			New: func() interface{} {
				return &http.Request{}
			},
		}
		
		// Response pool
		os.responsePool = &sync.Pool{
			New: func() interface{} {
				return &http.Response{}
			},
		}
		
		// Buffer pools for different sizes
		bufferSizes := []int{1024, 4096, 8192, 32768, 65536}
		for _, size := range bufferSizes {
			os.bufferPools[size] = &sync.Pool{
				New: func() interface{} {
					return make([]byte, size)
				},
			}
		}
		
		os.logger.Info("Object pools initialized",
			"request_pool", true,
			"response_pool", true,
			"buffer_pools", len(bufferSizes),
		)
	}
}

// applyOptimizations applies all production optimizations
func (os *OptimizationService) applyOptimizations() error {
	opt := os.optimizations
	
	// Apply runtime optimizations
	if err := os.applyRuntimeOptimizations(opt.Runtime); err != nil {
		return fmt.Errorf("failed to apply runtime optimizations: %w", err)
	}
	
	// Apply memory optimizations
	if err := os.applyMemoryOptimizations(opt.Memory); err != nil {
		return fmt.Errorf("failed to apply memory optimizations: %w", err)
	}
	
	os.logger.Info("Production optimizations applied successfully")
	return nil
}

// applyRuntimeOptimizations applies Go runtime optimizations
func (os *OptimizationService) applyRuntimeOptimizations(opt config.RuntimeOptimization) error {
	// Set GOMAXPROCS
	if opt.MaxProcs > 0 {
		oldProcs := runtime.GOMAXPROCS(opt.MaxProcs)
		os.logger.Info("GOMAXPROCS updated",
			"old_value", oldProcs,
			"new_value", opt.MaxProcs,
		)
	}
	
	// Set GC target percentage
	if opt.GCTargetPercentage > 0 {
		oldPercent := debug.SetGCPercent(opt.GCTargetPercentage)
		os.logger.Info("GC target percentage updated",
			"old_value", oldPercent,
			"new_value", opt.GCTargetPercentage,
		)
	}
	
	// Create memory ballast if enabled
	if opt.EnableMemoryBallast && opt.MemoryBallastSize > 0 {
		ballast := make([]byte, opt.MemoryBallastSize)
		runtime.KeepAlive(ballast)
		os.logger.Info("Memory ballast created",
			"size_mb", opt.MemoryBallastSize/(1024*1024),
		)
	}
	
	return nil
}

// applyMemoryOptimizations applies memory management optimizations
func (os *OptimizationService) applyMemoryOptimizations(opt config.MemoryOptimization) error {
	// Start memory monitoring if enabled
	if opt.MemoryCheckInterval > 0 {
		go os.monitorMemoryUsage(opt)
		os.logger.Info("Memory monitoring started",
			"check_interval", opt.MemoryCheckInterval,
			"warning_threshold", opt.MemoryWarningThreshold,
			"panic_threshold", opt.MemoryPanicThreshold,
		)
	}
	
	return nil
}

// monitorMemoryUsage monitors memory usage and takes action if needed
func (os *OptimizationService) monitorMemoryUsage(opt config.MemoryOptimization) {
	ticker := time.NewTicker(opt.MemoryCheckInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-os.ctx.Done():
			return
		case <-ticker.C:
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			
			currentUsage := m.Alloc
			usagePercent := float64(currentUsage) / float64(opt.MaxMemoryUsage)
			
			// Update metrics
			os.mu.Lock()
			os.resourceMetrics.MemoryUsage = int64(currentUsage)
			os.resourceMetrics.MemoryPercentage = usagePercent
			os.resourceMetrics.HeapSize = m.HeapAlloc
			os.resourceMetrics.StackSize = m.StackInuse
			os.resourceMetrics.GoroutineCount = runtime.NumGoroutine()
			os.resourceMetrics.LastUpdated = time.Now()
			os.mu.Unlock()
			
			if usagePercent > opt.MemoryPanicThreshold {
				// Critical memory usage - force GC and consider panic
				runtime.GC()
				runtime.ReadMemStats(&m)
				newUsage := float64(m.Alloc) / float64(opt.MaxMemoryUsage)
				
				os.logger.Error("Critical memory usage detected",
					"usage_percent", usagePercent*100,
					"usage_after_gc", newUsage*100,
					"threshold", opt.MemoryPanicThreshold*100,
				)
				
				if newUsage > opt.MemoryPanicThreshold {
					// Memory still critical after GC - this is a panic situation
					panic(fmt.Sprintf("Critical memory usage: %.2f%% (threshold: %.2f%%)",
						newUsage*100, opt.MemoryPanicThreshold*100))
				}
			} else if usagePercent > opt.MemoryWarningThreshold {
				// Warning threshold - force GC
				runtime.GC()
				os.logger.Warn("High memory usage detected, forced GC",
					"usage_percent", usagePercent*100,
					"threshold", opt.MemoryWarningThreshold*100,
				)
			}
		}
	}
}

// startMonitoring starts the performance monitoring goroutine
func (os *OptimizationService) startMonitoring() {
	monitoringOpt := os.optimizations.Monitoring
	
	if monitoringOpt.MetricsInterval > 0 {
		os.monitoringTicker = time.NewTicker(monitoringOpt.MetricsInterval)
		
		go func() {
			for {
				select {
				case <-os.ctx.Done():
					return
				case <-os.monitoringTicker.C:
					os.collectMetrics()
				}
			}
		}()
		
		os.logger.Info("Performance monitoring started",
			"interval", monitoringOpt.MetricsInterval,
		)
	}
}

// collectMetrics collects and updates performance metrics
func (os *OptimizationService) collectMetrics() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	os.mu.Lock()
	defer os.mu.Unlock()
	
	// Update resource metrics
	os.resourceMetrics.MemoryUsage = int64(m.Alloc)
	os.resourceMetrics.HeapSize = m.HeapAlloc
	os.resourceMetrics.StackSize = m.StackInuse
	os.resourceMetrics.GoroutineCount = runtime.NumGoroutine()
	os.resourceMetrics.GCCycles = m.NumGC
	
	// Calculate GC pause time
	if len(m.PauseNs) > 0 {
		totalPause := time.Duration(0)
		for _, pause := range m.PauseNs {
			totalPause += time.Duration(pause)
		}
		os.resourceMetrics.GCPauseTotal = totalPause
	}
	
	os.resourceMetrics.LastUpdated = time.Now()
	
	// Update request metrics
	os.updateRequestMetrics()
}

// updateRequestMetrics updates request-level metrics
func (os *OptimizationService) updateRequestMetrics() {
	// Calculate latency percentiles
	if len(os.requestMetrics.latencyHistogram) > 0 {
		hist := make([]time.Duration, len(os.requestMetrics.latencyHistogram))
		copy(hist, os.requestMetrics.latencyHistogram)
		
		// Sort for percentile calculation
		for i := 0; i < len(hist); i++ {
			for j := i + 1; j < len(hist); j++ {
				if hist[i] > hist[j] {
					hist[i], hist[j] = hist[j], hist[i]
				}
			}
		}
		
		length := len(hist)
		if length > 0 {
			os.requestMetrics.P50Latency = hist[length*50/100]
			os.requestMetrics.P95Latency = hist[length*95/100]
			os.requestMetrics.P99Latency = hist[length*99/100]
		}
		
		// Calculate average latency
		var total time.Duration
		for _, latency := range hist {
			total += latency
		}
		os.requestMetrics.AverageLatency = total / time.Duration(length)
	}
	
	os.requestMetrics.LastUpdated = time.Now()
}

// RecordRequest records a request for performance tracking
func (os *OptimizationService) RecordRequest(duration time.Duration, isError bool) {
	os.mu.Lock()
	defer os.mu.Unlock()
	
	os.requestMetrics.TotalRequests++
	
	// Add to latency histogram (keep last 1000 requests)
	os.requestMetrics.latencyHistogram = append(os.requestMetrics.latencyHistogram, duration)
	if len(os.requestMetrics.latencyHistogram) > 1000 {
		os.requestMetrics.latencyHistogram = os.requestMetrics.latencyHistogram[1:]
	}
	
	// Update error rate
	if isError {
		errorCount := 0
		for i := len(os.requestMetrics.latencyHistogram) - 1; i >= 0 && i >= len(os.requestMetrics.latencyHistogram)-100; i-- {
			// This is a simplified error tracking - in production you'd want proper error counting
		}
		os.requestMetrics.ErrorRate = float64(errorCount) / float64(len(os.requestMetrics.latencyHistogram))
	}
}

// GetBuffer returns a buffer from the pool
func (os *OptimizationService) GetBuffer(size int) []byte {
	memOpt := os.optimizations.Memory
	if !memOpt.EnableBufferPools {
		return make([]byte, size)
	}
	
	// Find the closest buffer pool size
	poolSize := size
	for s := range os.bufferPools {
		if s >= size {
			poolSize = s
			break
		}
	}
	
	if pool, exists := os.bufferPools[poolSize]; exists {
		if buffer := pool.Get(); buffer != nil {
			buf := buffer.([]byte)
			return buf[:size]
		}
	}
	
	return make([]byte, size)
}

// PutBuffer returns a buffer to the pool
func (os *OptimizationService) PutBuffer(buffer []byte) {
	memOpt := os.optimizations.Memory
	if !memOpt.EnableBufferPools {
		return
	}
	
	size := cap(buffer)
	if pool, exists := os.bufferPools[size]; exists {
		pool.Put(buffer)
	}
}

// GetPerformanceReport generates a comprehensive performance report
func (os *OptimizationService) GetPerformanceReport() *PerformanceReport {
	os.mu.RLock()
	defer os.mu.RUnlock()
	
	// Copy metrics to avoid race conditions
	requestMetrics := *os.requestMetrics
	resourceMetrics := *os.resourceMetrics
	
	// Generate recommendations
	recommendations := os.generateRecommendations(&requestMetrics, &resourceMetrics)
	
	// Determine health status
	healthStatus := os.determineHealthStatus(&requestMetrics, &resourceMetrics)
	
	return &PerformanceReport{
		Timestamp:       time.Now(),
		RequestMetrics:  &requestMetrics,
		ResourceMetrics: &resourceMetrics,
		Optimizations: map[string]interface{}{
			"object_pools_enabled":    os.optimizations.Memory.EnableObjectPools,
			"buffer_pools_enabled":    os.optimizations.Memory.EnableBufferPools,
			"memory_ballast_enabled":  os.optimizations.Runtime.EnableMemoryBallast,
			"async_logging_enabled":   os.optimizations.Logging.EnableAsyncLogging,
			"compression_enabled":     os.optimizations.Server.EnableCompression,
		},
		Recommendations: recommendations,
		HealthStatus:    healthStatus,
	}
}

// generateRecommendations analyzes metrics and provides optimization recommendations
func (os *OptimizationService) generateRecommendations(req *RequestMetrics, res *ResourceMetrics) []string {
	var recommendations []string
	
	// Memory recommendations
	if res.MemoryPercentage > 0.8 {
		recommendations = append(recommendations, "Consider increasing memory limit or optimizing memory usage")
	}
	
	// Latency recommendations
	if req.P95Latency > 1*time.Second {
		recommendations = append(recommendations, "High P95 latency detected - consider optimizing database queries or enabling caching")
	}
	
	// GC recommendations
	if res.GCPauseTotal > 100*time.Millisecond {
		recommendations = append(recommendations, "High GC pause time - consider tuning GC settings or reducing allocations")
	}
	
	// Goroutine recommendations
	if res.GoroutineCount > 10000 {
		recommendations = append(recommendations, "High goroutine count - check for goroutine leaks")
	}
	
	// Error rate recommendations
	if req.ErrorRate > 0.05 {
		recommendations = append(recommendations, "High error rate detected - investigate error causes")
	}
	
	return recommendations
}

// determineHealthStatus determines overall system health
func (os *OptimizationService) determineHealthStatus(req *RequestMetrics, res *ResourceMetrics) string {
	if res.MemoryPercentage > 0.95 || req.ErrorRate > 0.1 || req.P95Latency > 5*time.Second {
		return "CRITICAL"
	}
	
	if res.MemoryPercentage > 0.8 || req.ErrorRate > 0.05 || req.P95Latency > 2*time.Second {
		return "WARNING"
	}
	
	return "HEALTHY"
}

// OptimizationMiddleware provides HTTP middleware for performance optimization
func (os *OptimizationService) OptimizationMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			// Track active requests
			os.mu.Lock()
			os.requestMetrics.ActiveRequests++
			os.mu.Unlock()
			
			// Process request
			next.ServeHTTP(w, r)
			
			// Record metrics
			duration := time.Since(start)
			isError := false // This would be determined by response status
			
			os.RecordRequest(duration, isError)
			
			// Decrement active requests
			os.mu.Lock()
			os.requestMetrics.ActiveRequests--
			os.mu.Unlock()
		})
	}
}

// Close gracefully shuts down the optimization service
func (os *OptimizationService) Close() error {
	os.cancel()
	
	if os.monitoringTicker != nil {
		os.monitoringTicker.Stop()
	}
	
	os.logger.Info("Optimization service stopped")
	return nil
}