package domain

import "time"

// URLStats represents comprehensive URL statistics
type URLStats struct {
	ShortCode      string     `json:"short_code"`
	URL            *ShortURL  `json:"url"`
	TotalClicks    int64      `json:"total_clicks"`
	ClicksInRange  int64      `json:"clicks_in_range"`
	UniqueVisitors int64      `json:"unique_visitors"`
	StartDate      time.Time  `json:"start_date"`
	EndDate        time.Time  `json:"end_date"`
}

// URLStats specific to performance monitoring (different from click.go ClickStats)

// DeviceStats represents device and browser statistics
type DeviceStats struct {
	ShortURLID  uint                      `json:"short_url_id"`
	DeviceTypes map[string]int64          `json:"device_types"`
	Browsers    map[string]int64          `json:"browsers"`
	TopDevices  []map[string]interface{}  `json:"top_devices"`
	TopBrowsers []map[string]interface{}  `json:"top_browsers"`
}

// ReferrerStats represents referrer statistics
type ReferrerStats struct {
	Referrer   string `json:"referrer" gorm:"column:referrer"`
	ClickCount int64  `json:"click_count" gorm:"column:click_count"`
}

// Removed TimelineStats - exists in click.go

// ClickHeatmap represents click patterns by day and hour
type ClickHeatmap struct {
	ShortURLID uint                      `json:"short_url_id"`
	Data       map[string]map[int]int64  `json:"data"` // [day_of_week][hour] = click_count
}

// Removed GlobalStats - exists in click.go

// URLPerformance represents URL performance metrics
type URLPerformance struct {
	ShortURL       string `json:"short_url" gorm:"column:short_code"`
	OriginalURL    string `json:"original_url" gorm:"column:original_url"`
	Title          string `json:"title" gorm:"column:title"`
	TotalClicks    int64  `json:"total_clicks" gorm:"column:click_count"`
	UniqueClicks   int64  `json:"unique_clicks" gorm:"column:unique_visitors"`
	ClickRate      float64 `json:"click_rate"`
}

// RealtimeStats represents real-time application statistics
type RealtimeStats struct {
	Since          time.Time                `json:"since"`
	RecentClicks   int64                    `json:"recent_clicks"`
	ActiveSessions int64                    `json:"active_sessions"`
	TopURLs        []RealtimeURLPerformance `json:"top_urls"`
}

// RealtimeURLPerformance represents real-time URL performance
type RealtimeURLPerformance struct {
	ShortCode  string `json:"short_code" gorm:"column:short_code"`
	ClickCount int64  `json:"click_count" gorm:"column:click_count"`
}

// PerformanceMetrics represents application performance metrics
type PerformanceMetrics struct {
	Timestamp           time.Time     `json:"timestamp"`
	AverageResponseTime time.Duration `json:"average_response_time"`
	RequestsPerSecond   float64       `json:"requests_per_second"`
	ErrorRate           float64       `json:"error_rate"`
	CacheHitRate        float64       `json:"cache_hit_rate"`
	DatabaseConnections int           `json:"database_connections"`
	MemoryUsage         int64         `json:"memory_usage_bytes"`
	CPUUsage            float64       `json:"cpu_usage_percent"`
}

// Removed duplicate types - they exist in health.go

// OptimizationRecommendation represents a performance optimization suggestion
type OptimizationRecommendation struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"` // "index", "query", "cache", "configuration"
	Priority    string    `json:"priority"` // "high", "medium", "low"
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Impact      string    `json:"impact"`
	Effort      string    `json:"effort"`
	Query       string    `json:"query,omitempty"`
	Suggestion  string    `json:"suggestion"`
	CreatedAt   time.Time `json:"created_at"`
}

// PerformanceAlert represents a performance alert
type PerformanceAlert struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Type        string    `json:"type"` // "slow_query", "high_memory", "high_cpu", "error_rate"
	Severity    string    `json:"severity"` // "critical", "warning", "info"
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Threshold   float64   `json:"threshold"`
	ActualValue float64   `json:"actual_value"`
	Query       string    `json:"query,omitempty"`
	Resolved    bool      `json:"resolved"`
	CreatedAt   time.Time `json:"created_at"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
}

// QueryExecutionPlan represents a database query execution plan
type QueryExecutionPlan struct {
	Query          string    `json:"query"`
	Plan           []string  `json:"plan"`
	EstimatedCost  float64   `json:"estimated_cost"`
	ActualCost     float64   `json:"actual_cost"`
	ExecutionTime  time.Duration `json:"execution_time"`
	RowsProcessed  int64     `json:"rows_processed"`
	IndexesUsed    []string  `json:"indexes_used"`
	Suggestions    []string  `json:"suggestions"`
	CreatedAt      time.Time `json:"created_at"`
}

// BulkOperationResult represents the result of a bulk operation
type BulkOperationResult struct {
	TotalRecords     int64         `json:"total_records"`
	ProcessedRecords int64         `json:"processed_records"`
	FailedRecords    int64         `json:"failed_records"`
	Errors           []string      `json:"errors,omitempty"`
	Duration         time.Duration `json:"duration"`
	StartTime        time.Time     `json:"start_time"`
	EndTime          time.Time     `json:"end_time"`
}

// CacheOperation represents cache operation statistics
type CacheOperation struct {
	Operation string        `json:"operation"` // "get", "set", "delete", "exists"
	Key       string        `json:"key"`
	Hit       bool          `json:"hit"`
	Duration  time.Duration `json:"duration"`
	Size      int64         `json:"size,omitempty"`
	TTL       time.Duration `json:"ttl,omitempty"`
	Timestamp time.Time     `json:"timestamp"`
}

// Methods for performance domain objects

// IsHealthy method moved to health.go

// SystemHealth methods moved to health.go

// IsSlowQuery checks if a query execution time exceeds the threshold
func (qep *QueryExecutionPlan) IsSlowQuery(threshold time.Duration) bool {
	return qep.ExecutionTime > threshold
}

// Success checks if the bulk operation was successful
func (bor *BulkOperationResult) Success() bool {
	return bor.FailedRecords == 0
}

// SuccessRate calculates the success rate of the bulk operation
func (bor *BulkOperationResult) SuccessRate() float64 {
	if bor.TotalRecords == 0 {
		return 0
	}
	return float64(bor.ProcessedRecords) / float64(bor.TotalRecords) * 100
}

// HitRate calculates the cache hit rate for a series of operations
func CalculateCacheHitRate(operations []CacheOperation) float64 {
	if len(operations) == 0 {
		return 0
	}
	
	hits := 0
	for _, op := range operations {
		if op.Hit && op.Operation == "get" {
			hits++
		}
	}
	
	getOperations := 0
	for _, op := range operations {
		if op.Operation == "get" {
			getOperations++
		}
	}
	
	if getOperations == 0 {
		return 0
	}
	
	return float64(hits) / float64(getOperations) * 100
}

// AverageExecutionTime calculates the average execution time for query plans
func AverageExecutionTime(plans []QueryExecutionPlan) time.Duration {
	if len(plans) == 0 {
		return 0
	}
	
	var total time.Duration
	for _, plan := range plans {
		total += plan.ExecutionTime
	}
	
	return total / time.Duration(len(plans))
}