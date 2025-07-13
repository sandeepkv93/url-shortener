package config

import (
	"runtime"
	"runtime/debug"
	"time"
)

// ProductionOptimization contains all production-specific optimization settings
type ProductionOptimization struct {
	// Server optimizations
	Server ServerOptimization `json:"server" yaml:"server"`

	// Database optimizations
	Database DatabaseOptimization `json:"database" yaml:"database"`

	// Cache optimizations
	Cache CacheOptimization `json:"cache" yaml:"cache"`

	// Runtime optimizations
	Runtime RuntimeOptimization `json:"runtime" yaml:"runtime"`

	// Memory optimizations
	Memory MemoryOptimization `json:"memory" yaml:"memory"`

	// Logging optimizations
	Logging LoggingOptimization `json:"logging" yaml:"logging"`

	// Monitoring configurations
	Monitoring MonitoringOptimization `json:"monitoring" yaml:"monitoring"`

	// Resource limits
	Resources ResourceOptimization `json:"resources" yaml:"resources"`
}

// ServerOptimization contains HTTP server optimization settings
type ServerOptimization struct {
	// Connection settings
	MaxConnections        int           `json:"max_connections" yaml:"max_connections"`
	MaxIdleConnections    int           `json:"max_idle_connections" yaml:"max_idle_connections"`
	ConnectionTimeout     time.Duration `json:"connection_timeout" yaml:"connection_timeout"`
	IdleTimeout           time.Duration `json:"idle_timeout" yaml:"idle_timeout"`
	ReadTimeout           time.Duration `json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout          time.Duration `json:"write_timeout" yaml:"write_timeout"`
	ReadHeaderTimeout     time.Duration `json:"read_header_timeout" yaml:"read_header_timeout"`
	KeepAlive             bool          `json:"keep_alive" yaml:"keep_alive"`
	KeepAliveTimeout      time.Duration `json:"keep_alive_timeout" yaml:"keep_alive_timeout"`
	
	// Buffer sizes
	ReadBufferSize        int           `json:"read_buffer_size" yaml:"read_buffer_size"`
	WriteBufferSize       int           `json:"write_buffer_size" yaml:"write_buffer_size"`
	MaxHeaderBytes        int           `json:"max_header_bytes" yaml:"max_header_bytes"`
	
	// Compression
	EnableCompression     bool          `json:"enable_compression" yaml:"enable_compression"`
	CompressionLevel      int           `json:"compression_level" yaml:"compression_level"`
	CompressionThreshold  int           `json:"compression_threshold" yaml:"compression_threshold"`
	
	// Request limits
	MaxRequestSize        int64         `json:"max_request_size" yaml:"max_request_size"`
	MaxMultipartMemory    int64         `json:"max_multipart_memory" yaml:"max_multipart_memory"`
}

// DatabaseOptimization contains database connection and query optimization settings
type DatabaseOptimization struct {
	// Connection pool settings
	MaxOpenConnections    int           `json:"max_open_connections" yaml:"max_open_connections"`
	MaxIdleConnections    int           `json:"max_idle_connections" yaml:"max_idle_connections"`
	ConnectionMaxLifetime time.Duration `json:"connection_max_lifetime" yaml:"connection_max_lifetime"`
	ConnectionMaxIdleTime time.Duration `json:"connection_max_idle_time" yaml:"connection_max_idle_time"`
	
	// Query optimization
	PreparedStatements    bool          `json:"prepared_statements" yaml:"prepared_statements"`
	QueryTimeout          time.Duration `json:"query_timeout" yaml:"query_timeout"`
	SlowQueryThreshold    time.Duration `json:"slow_query_threshold" yaml:"slow_query_threshold"`
	
	// Batch settings
	BatchSize             int           `json:"batch_size" yaml:"batch_size"`
	BatchTimeout          time.Duration `json:"batch_timeout" yaml:"batch_timeout"`
	
	// Health check
	PingInterval          time.Duration `json:"ping_interval" yaml:"ping_interval"`
	PingTimeout           time.Duration `json:"ping_timeout" yaml:"ping_timeout"`
}

// CacheOptimization contains cache-specific optimization settings
type CacheOptimization struct {
	// Redis connection settings
	PoolSize              int           `json:"pool_size" yaml:"pool_size"`
	MinIdleConnections    int           `json:"min_idle_connections" yaml:"min_idle_connections"`
	MaxRetries            int           `json:"max_retries" yaml:"max_retries"`
	RetryDelay            time.Duration `json:"retry_delay" yaml:"retry_delay"`
	DialTimeout           time.Duration `json:"dial_timeout" yaml:"dial_timeout"`
	ReadTimeout           time.Duration `json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout          time.Duration `json:"write_timeout" yaml:"write_timeout"`
	PoolTimeout           time.Duration `json:"pool_timeout" yaml:"pool_timeout"`
	IdleTimeout           time.Duration `json:"idle_timeout" yaml:"idle_timeout"`
	IdleCheckFrequency    time.Duration `json:"idle_check_frequency" yaml:"idle_check_frequency"`
	
	// Cache behavior
	DefaultTTL            time.Duration `json:"default_ttl" yaml:"default_ttl"`
	EnableCompression     bool          `json:"enable_compression" yaml:"enable_compression"`
	CompressionThreshold  int           `json:"compression_threshold" yaml:"compression_threshold"`
	
	// Cache warming
	EnableCacheWarming    bool          `json:"enable_cache_warming" yaml:"enable_cache_warming"`
	WarmupBatchSize       int           `json:"warmup_batch_size" yaml:"warmup_batch_size"`
	WarmupConcurrency     int           `json:"warmup_concurrency" yaml:"warmup_concurrency"`
}

// RuntimeOptimization contains Go runtime optimization settings
type RuntimeOptimization struct {
	// Garbage collection
	GCTargetPercentage    int           `json:"gc_target_percentage" yaml:"gc_target_percentage"`
	MaxGCPauseTarget      time.Duration `json:"max_gc_pause_target" yaml:"max_gc_pause_target"`
	EnableGCTrace         bool          `json:"enable_gc_trace" yaml:"enable_gc_trace"`
	
	// CPU and concurrency
	MaxProcs              int           `json:"max_procs" yaml:"max_procs"`
	WorkerPoolSize        int           `json:"worker_pool_size" yaml:"worker_pool_size"`
	QueueSize             int           `json:"queue_size" yaml:"queue_size"`
	
	// Memory ballast
	MemoryBallastSize     int64         `json:"memory_ballast_size" yaml:"memory_ballast_size"`
	EnableMemoryBallast   bool          `json:"enable_memory_ballast" yaml:"enable_memory_ballast"`
}

// MemoryOptimization contains memory management optimization settings
type MemoryOptimization struct {
	// Memory limits
	MaxMemoryUsage        int64         `json:"max_memory_usage" yaml:"max_memory_usage"`
	MemoryWarningThreshold float64      `json:"memory_warning_threshold" yaml:"memory_warning_threshold"`
	MemoryPanicThreshold  float64       `json:"memory_panic_threshold" yaml:"memory_panic_threshold"`
	
	// Buffer pools
	EnableBufferPools     bool          `json:"enable_buffer_pools" yaml:"enable_buffer_pools"`
	BufferPoolSize        int           `json:"buffer_pool_size" yaml:"buffer_pool_size"`
	MaxBufferSize         int           `json:"max_buffer_size" yaml:"max_buffer_size"`
	
	// Object pools
	EnableObjectPools     bool          `json:"enable_object_pools" yaml:"enable_object_pools"`
	ObjectPoolSize        int           `json:"object_pool_size" yaml:"object_pool_size"`
	
	// Memory monitoring
	MemoryCheckInterval   time.Duration `json:"memory_check_interval" yaml:"memory_check_interval"`
	EnableMemoryProfiling bool          `json:"enable_memory_profiling" yaml:"enable_memory_profiling"`
}

// LoggingOptimization contains logging optimization settings
type LoggingOptimization struct {
	// Async logging
	EnableAsyncLogging    bool          `json:"enable_async_logging" yaml:"enable_async_logging"`
	LogBufferSize         int           `json:"log_buffer_size" yaml:"log_buffer_size"`
	LogFlushInterval      time.Duration `json:"log_flush_interval" yaml:"log_flush_interval"`
	
	// Log rotation
	MaxLogFileSize        int64         `json:"max_log_file_size" yaml:"max_log_file_size"`
	MaxLogFiles           int           `json:"max_log_files" yaml:"max_log_files"`
	LogRetentionDays      int           `json:"log_retention_days" yaml:"log_retention_days"`
	
	// Performance logging
	EnablePerformanceLogs bool          `json:"enable_performance_logs" yaml:"enable_performance_logs"`
	SlowRequestThreshold  time.Duration `json:"slow_request_threshold" yaml:"slow_request_threshold"`
	
	// Sampling
	LogSamplingRate       float64       `json:"log_sampling_rate" yaml:"log_sampling_rate"`
	ErrorLogSamplingRate  float64       `json:"error_log_sampling_rate" yaml:"error_log_sampling_rate"`
}

// MonitoringOptimization contains monitoring and metrics optimization settings
type MonitoringOptimization struct {
	// Metrics collection
	MetricsInterval       time.Duration `json:"metrics_interval" yaml:"metrics_interval"`
	EnableSystemMetrics   bool          `json:"enable_system_metrics" yaml:"enable_system_metrics"`
	EnableRuntimeMetrics  bool          `json:"enable_runtime_metrics" yaml:"enable_runtime_metrics"`
	EnableCustomMetrics   bool          `json:"enable_custom_metrics" yaml:"enable_custom_metrics"`
	
	// Health checks
	HealthCheckInterval   time.Duration `json:"health_check_interval" yaml:"health_check_interval"`
	HealthCheckTimeout    time.Duration `json:"health_check_timeout" yaml:"health_check_timeout"`
	
	// Profiling
	EnableProfiling       bool          `json:"enable_profiling" yaml:"enable_profiling"`
	ProfilingInterval     time.Duration `json:"profiling_interval" yaml:"profiling_interval"`
	EnableCPUProfiling    bool          `json:"enable_cpu_profiling" yaml:"enable_cpu_profiling"`
	EnableMemoryProfiling bool          `json:"enable_memory_profiling" yaml:"enable_memory_profiling"`
	
	// Tracing
	EnableTracing         bool          `json:"enable_tracing" yaml:"enable_tracing"`
	TracingSampleRate     float64       `json:"tracing_sample_rate" yaml:"tracing_sample_rate"`
	TraceExporter         string        `json:"trace_exporter" yaml:"trace_exporter"`
}

// ResourceOptimization contains resource limit optimization settings
type ResourceOptimization struct {
	// CPU limits
	CPUQuota              float64       `json:"cpu_quota" yaml:"cpu_quota"`
	CPUShares             int           `json:"cpu_shares" yaml:"cpu_shares"`
	CPUThrottleThreshold  float64       `json:"cpu_throttle_threshold" yaml:"cpu_throttle_threshold"`
	
	// Network limits
	NetworkBandwidthLimit int64         `json:"network_bandwidth_limit" yaml:"network_bandwidth_limit"`
	MaxConcurrentRequests int           `json:"max_concurrent_requests" yaml:"max_concurrent_requests"`
	
	// File descriptor limits
	MaxFileDescriptors    int           `json:"max_file_descriptors" yaml:"max_file_descriptors"`
	
	// Disk I/O limits
	DiskIOThreshold       int64         `json:"disk_io_threshold" yaml:"disk_io_threshold"`
	MaxDiskUsage          int64         `json:"max_disk_usage" yaml:"max_disk_usage"`
}

// GetProductionOptimization returns optimized configuration for production
func (c *Config) GetProductionOptimization() *ProductionOptimization {
	if c.Production.Optimization == nil {
		c.Production.Optimization = c.getDefaultProductionOptimization()
	}
	return c.Production.Optimization
}

// getDefaultProductionOptimization returns default production optimization settings
func (c *Config) getDefaultProductionOptimization() *ProductionOptimization {
	numCPU := runtime.NumCPU()
	
	return &ProductionOptimization{
		Server: ServerOptimization{
			MaxConnections:        10000,
			MaxIdleConnections:    1000,
			ConnectionTimeout:     30 * time.Second,
			IdleTimeout:           120 * time.Second,
			ReadTimeout:           15 * time.Second,
			WriteTimeout:          15 * time.Second,
			ReadHeaderTimeout:     10 * time.Second,
			KeepAlive:             true,
			KeepAliveTimeout:      15 * time.Second,
			ReadBufferSize:        32 * 1024,  // 32KB
			WriteBufferSize:       32 * 1024,  // 32KB
			MaxHeaderBytes:        64 * 1024,  // 64KB
			EnableCompression:     true,
			CompressionLevel:      6,
			CompressionThreshold:  1024,       // 1KB
			MaxRequestSize:        32 * 1024 * 1024, // 32MB
			MaxMultipartMemory:    32 * 1024 * 1024, // 32MB
		},
		Database: DatabaseOptimization{
			MaxOpenConnections:    numCPU * 4,
			MaxIdleConnections:    numCPU * 2,
			ConnectionMaxLifetime: 1 * time.Hour,
			ConnectionMaxIdleTime: 15 * time.Minute,
			PreparedStatements:    true,
			QueryTimeout:          30 * time.Second,
			SlowQueryThreshold:    1 * time.Second,
			BatchSize:             1000,
			BatchTimeout:          5 * time.Second,
			PingInterval:          1 * time.Minute,
			PingTimeout:           5 * time.Second,
		},
		Cache: CacheOptimization{
			PoolSize:              numCPU * 10,
			MinIdleConnections:    numCPU * 2,
			MaxRetries:            3,
			RetryDelay:            100 * time.Millisecond,
			DialTimeout:           5 * time.Second,
			ReadTimeout:           3 * time.Second,
			WriteTimeout:          3 * time.Second,
			PoolTimeout:           4 * time.Second,
			IdleTimeout:           5 * time.Minute,
			IdleCheckFrequency:    1 * time.Minute,
			DefaultTTL:            1 * time.Hour,
			EnableCompression:     true,
			CompressionThreshold:  1024, // 1KB
			EnableCacheWarming:    true,
			WarmupBatchSize:       100,
			WarmupConcurrency:     numCPU,
		},
		Runtime: RuntimeOptimization{
			GCTargetPercentage:    100,
			MaxGCPauseTarget:      10 * time.Millisecond,
			EnableGCTrace:         false,
			MaxProcs:              numCPU,
			WorkerPoolSize:        numCPU * 2,
			QueueSize:             10000,
			MemoryBallastSize:     100 * 1024 * 1024, // 100MB
			EnableMemoryBallast:   true,
		},
		Memory: MemoryOptimization{
			MaxMemoryUsage:        2 * 1024 * 1024 * 1024, // 2GB
			MemoryWarningThreshold: 0.8,
			MemoryPanicThreshold:  0.95,
			EnableBufferPools:     true,
			BufferPoolSize:        1000,
			MaxBufferSize:         64 * 1024, // 64KB
			EnableObjectPools:     true,
			ObjectPoolSize:        500,
			MemoryCheckInterval:   30 * time.Second,
			EnableMemoryProfiling: false,
		},
		Logging: LoggingOptimization{
			EnableAsyncLogging:    true,
			LogBufferSize:         10000,
			LogFlushInterval:      5 * time.Second,
			MaxLogFileSize:        100 * 1024 * 1024, // 100MB
			MaxLogFiles:           10,
			LogRetentionDays:      30,
			EnablePerformanceLogs: true,
			SlowRequestThreshold:  1 * time.Second,
			LogSamplingRate:       1.0,
			ErrorLogSamplingRate:  1.0,
		},
		Monitoring: MonitoringOptimization{
			MetricsInterval:       30 * time.Second,
			EnableSystemMetrics:   true,
			EnableRuntimeMetrics:  true,
			EnableCustomMetrics:   true,
			HealthCheckInterval:   30 * time.Second,
			HealthCheckTimeout:    5 * time.Second,
			EnableProfiling:       false,
			ProfilingInterval:     10 * time.Minute,
			EnableCPUProfiling:    false,
			EnableMemoryProfiling: false,
			EnableTracing:         true,
			TracingSampleRate:     0.1,
			TraceExporter:         "jaeger",
		},
		Resources: ResourceOptimization{
			CPUQuota:              float64(numCPU),
			CPUShares:             1024,
			CPUThrottleThreshold:  0.9,
			NetworkBandwidthLimit: 1000 * 1024 * 1024, // 1Gbps
			MaxConcurrentRequests: 10000,
			MaxFileDescriptors:    65536,
			DiskIOThreshold:       100 * 1024 * 1024, // 100MB/s
			MaxDiskUsage:          10 * 1024 * 1024 * 1024, // 10GB
		},
	}
}

// ApplyOptimizations applies all production optimizations
func (c *Config) ApplyOptimizations() error {
	if !c.IsProduction() {
		return nil
	}
	
	opt := c.GetProductionOptimization()
	
	// Apply runtime optimizations
	if err := c.applyRuntimeOptimizations(opt.Runtime); err != nil {
		return err
	}
	
	// Apply memory optimizations
	if err := c.applyMemoryOptimizations(opt.Memory); err != nil {
		return err
	}
	
	return nil
}

// applyRuntimeOptimizations applies Go runtime optimizations
func (c *Config) applyRuntimeOptimizations(opt RuntimeOptimization) error {
	// Set GOMAXPROCS
	if opt.MaxProcs > 0 {
		runtime.GOMAXPROCS(opt.MaxProcs)
	}
	
	// Set GC target percentage
	if opt.GCTargetPercentage > 0 {
		debug.SetGCPercent(opt.GCTargetPercentage)
	}
	
	// Create memory ballast if enabled
	if opt.EnableMemoryBallast && opt.MemoryBallastSize > 0 {
		ballast := make([]byte, opt.MemoryBallastSize)
		runtime.KeepAlive(ballast)
	}
	
	return nil
}

// applyMemoryOptimizations applies memory management optimizations
func (c *Config) applyMemoryOptimizations(opt MemoryOptimization) error {
	// Set memory limit if supported (Go 1.19+)
	if opt.MaxMemoryUsage > 0 {
		// Note: This would require Go 1.19+ with GOMEMLIMIT
		// For now, we'll monitor memory usage and implement soft limits
		go c.monitorMemoryUsage(opt)
	}
	
	return nil
}

// monitorMemoryUsage monitors memory usage and takes action if needed
func (c *Config) monitorMemoryUsage(opt MemoryOptimization) {
	ticker := time.NewTicker(opt.MemoryCheckInterval)
	defer ticker.Stop()
	
	for range ticker.C {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		
		currentUsage := m.Alloc
		usagePercent := float64(currentUsage) / float64(opt.MaxMemoryUsage)
		
		if usagePercent > opt.MemoryPanicThreshold {
			// Critical memory usage - force GC and consider panic
			runtime.GC()
			runtime.ReadMemStats(&m)
			newUsage := float64(m.Alloc) / float64(opt.MaxMemoryUsage)
			
			if newUsage > opt.MemoryPanicThreshold {
				// Memory still critical after GC
				panic("Critical memory usage detected")
			}
		} else if usagePercent > opt.MemoryWarningThreshold {
			// Warning threshold - force GC
			runtime.GC()
		}
	}
}

// GetServerOptimization returns server optimization settings
func (c *Config) GetServerOptimization() ServerOptimization {
	return c.GetProductionOptimization().Server
}

// GetDatabaseOptimization returns database optimization settings
func (c *Config) GetDatabaseOptimization() DatabaseOptimization {
	return c.GetProductionOptimization().Database
}

// GetCacheOptimization returns cache optimization settings
func (c *Config) GetCacheOptimization() CacheOptimization {
	return c.GetProductionOptimization().Cache
}

// GetRuntimeOptimization returns runtime optimization settings
func (c *Config) GetRuntimeOptimization() RuntimeOptimization {
	return c.GetProductionOptimization().Runtime
}

// GetMemoryOptimization returns memory optimization settings
func (c *Config) GetMemoryOptimization() MemoryOptimization {
	return c.GetProductionOptimization().Memory
}

// GetLoggingOptimization returns logging optimization settings
func (c *Config) GetLoggingOptimization() LoggingOptimization {
	return c.GetProductionOptimization().Logging
}

// GetMonitoringOptimization returns monitoring optimization settings
func (c *Config) GetMonitoringOptimization() MonitoringOptimization {
	return c.GetProductionOptimization().Monitoring
}

// GetResourceOptimization returns resource optimization settings
func (c *Config) GetResourceOptimization() ResourceOptimization {
	return c.GetProductionOptimization().Resources
}