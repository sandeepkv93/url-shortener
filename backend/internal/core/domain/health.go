package domain

import (
	"time"
)

type HealthStatus struct {
	Status      string                        `json:"status"` // healthy, degraded, unhealthy
	Version     string                        `json:"version"`
	Uptime      time.Duration                 `json:"uptime"`
	Timestamp   time.Time                     `json:"timestamp"`
	Components  map[string]*ComponentHealth   `json:"components"`
	Checks      map[string]*HealthCheck       `json:"checks"`
}

type ComponentHealth struct {
	Status      string            `json:"status"` // up, down, degraded
	Message     string            `json:"message"`
	LastChecked time.Time         `json:"last_checked"`
	Duration    time.Duration     `json:"duration"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type HealthCheck struct {
	Name        string        `json:"name"`
	Status      string        `json:"status"`
	Message     string        `json:"message"`
	Duration    time.Duration `json:"duration"`
	Critical    bool          `json:"critical"`
	LastRun     time.Time     `json:"last_run"`
	NextRun     time.Time     `json:"next_run"`
}

type SystemMetrics struct {
	CPU         CPUMetrics     `json:"cpu"`
	Memory      MemoryMetrics  `json:"memory"`
	Disk        DiskMetrics    `json:"disk"`
	Network     NetworkMetrics `json:"network"`
	Processes   ProcessMetrics `json:"processes"`
	LoadAverage []float64      `json:"load_average"`
	Timestamp   time.Time      `json:"timestamp"`
}

type CPUMetrics struct {
	Usage     float64 `json:"usage_percent"`
	Cores     int     `json:"cores"`
	Idle      float64 `json:"idle_percent"`
	System    float64 `json:"system_percent"`
	User      float64 `json:"user_percent"`
	IOWait    float64 `json:"iowait_percent"`
}

type MemoryMetrics struct {
	Total       uint64  `json:"total_bytes"`
	Used        uint64  `json:"used_bytes"`
	Free        uint64  `json:"free_bytes"`
	Available   uint64  `json:"available_bytes"`
	UsagePercent float64 `json:"usage_percent"`
	Cached      uint64  `json:"cached_bytes"`
	Buffers     uint64  `json:"buffers_bytes"`
}

type DiskMetrics struct {
	Total       uint64  `json:"total_bytes"`
	Used        uint64  `json:"used_bytes"`
	Free        uint64  `json:"free_bytes"`
	UsagePercent float64 `json:"usage_percent"`
	Inodes      uint64  `json:"inodes"`
	InodesUsed  uint64  `json:"inodes_used"`
	ReadOps     uint64  `json:"read_ops"`
	WriteOps    uint64  `json:"write_ops"`
}

type NetworkMetrics struct {
	BytesReceived uint64 `json:"bytes_received"`
	BytesSent     uint64 `json:"bytes_sent"`
	PacketsReceived uint64 `json:"packets_received"`
	PacketsSent   uint64 `json:"packets_sent"`
	ErrorsReceived uint64 `json:"errors_received"`
	ErrorsSent    uint64 `json:"errors_sent"`
}

type ProcessMetrics struct {
	Total     int `json:"total"`
	Running   int `json:"running"`
	Sleeping  int `json:"sleeping"`
	Stopped   int `json:"stopped"`
	Zombie    int `json:"zombie"`
}

type ApplicationMetrics struct {
	RequestsTotal       int64             `json:"requests_total"`
	RequestsPerSecond   float64           `json:"requests_per_second"`
	ResponseTimes       ResponseTimeMetrics `json:"response_times"`
	ErrorRate           float64           `json:"error_rate"`
	ActiveConnections   int64             `json:"active_connections"`
	DatabaseConnections DatabaseMetrics   `json:"database_connections"`
	CacheMetrics        CacheMetrics      `json:"cache_metrics"`
	GoroutineCount      int               `json:"goroutine_count"`
	HeapSize            uint64            `json:"heap_size_bytes"`
	GCPauses            []time.Duration   `json:"gc_pauses"`
	Timestamp           time.Time         `json:"timestamp"`
}

type ResponseTimeMetrics struct {
	P50 time.Duration `json:"p50"`
	P90 time.Duration `json:"p90"`
	P95 time.Duration `json:"p95"`
	P99 time.Duration `json:"p99"`
	Mean time.Duration `json:"mean"`
	Max  time.Duration `json:"max"`
}

type DatabaseMetrics struct {
	OpenConnections    int           `json:"open_connections"`
	InUseConnections   int           `json:"in_use_connections"`
	IdleConnections    int           `json:"idle_connections"`
	MaxOpenConnections int           `json:"max_open_connections"`
	QueryDuration      time.Duration `json:"avg_query_duration"`
	SlowQueries        int64         `json:"slow_queries"`
}

type CacheMetrics struct {
	HitRate        float64 `json:"hit_rate"`
	MissRate       float64 `json:"miss_rate"`
	KeyCount       int64   `json:"key_count"`
	MemoryUsage    uint64  `json:"memory_usage_bytes"`
	EvictedKeys    int64   `json:"evicted_keys"`
	ExpiredKeys    int64   `json:"expired_keys"`
	CommandsProcessed int64 `json:"commands_processed"`
}