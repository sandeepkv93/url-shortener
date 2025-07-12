package services

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// MetricsService provides comprehensive metrics collection and reporting
type MetricsService struct {
	counters   map[string]*Counter
	gauges     map[string]*Gauge
	histograms map[string]*Histogram
	timers     map[string]*Timer
	mu         sync.RWMutex
	logger     *LoggingService
}

// Counter represents a monotonically increasing counter
type Counter struct {
	name   string
	help   string
	labels map[string]string
	value  int64
	mu     sync.RWMutex
}

// Gauge represents a value that can go up and down
type Gauge struct {
	name   string
	help   string
	labels map[string]string
	value  float64
	mu     sync.RWMutex
}

// Histogram tracks distributions of values
type Histogram struct {
	name    string
	help    string
	labels  map[string]string
	buckets []float64
	counts  []int64
	sum     float64
	count   int64
	mu      sync.RWMutex
}

// Timer tracks timing information
type Timer struct {
	name      string
	help      string
	labels    map[string]string
	durations []time.Duration
	mu        sync.RWMutex
}

// MetricValue represents a metric reading
type MetricValue struct {
	Name      string            `json:"name"`
	Type      string            `json:"type"`
	Help      string            `json:"help"`
	Labels    map[string]string `json:"labels"`
	Value     interface{}       `json:"value"`
	Timestamp time.Time         `json:"timestamp"`
}

// NewMetricsService creates a new metrics service
func NewMetricsService(logger *LoggingService) *MetricsService {
	ms := &MetricsService{
		counters:   make(map[string]*Counter),
		gauges:     make(map[string]*Gauge),
		histograms: make(map[string]*Histogram),
		timers:     make(map[string]*Timer),
		logger:     logger,
	}
	
	// Initialize standard application metrics
	ms.initializeStandardMetrics()
	
	return ms
}

// NewCounter creates or retrieves a counter
func (ms *MetricsService) NewCounter(name, help string, labels map[string]string) *Counter {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	
	key := ms.buildKey(name, labels)
	if counter, exists := ms.counters[key]; exists {
		return counter
	}
	
	counter := &Counter{
		name:   name,
		help:   help,
		labels: labels,
	}
	ms.counters[key] = counter
	
	ms.logger.Debug("Created new counter", "name", name, "labels", labels)
	return counter
}

// NewGauge creates or retrieves a gauge
func (ms *MetricsService) NewGauge(name, help string, labels map[string]string) *Gauge {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	
	key := ms.buildKey(name, labels)
	if gauge, exists := ms.gauges[key]; exists {
		return gauge
	}
	
	gauge := &Gauge{
		name:   name,
		help:   help,
		labels: labels,
	}
	ms.gauges[key] = gauge
	
	ms.logger.Debug("Created new gauge", "name", name, "labels", labels)
	return gauge
}

// NewHistogram creates or retrieves a histogram
func (ms *MetricsService) NewHistogram(name, help string, buckets []float64, labels map[string]string) *Histogram {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	
	key := ms.buildKey(name, labels)
	if histogram, exists := ms.histograms[key]; exists {
		return histogram
	}
	
	histogram := &Histogram{
		name:    name,
		help:    help,
		labels:  labels,
		buckets: buckets,
		counts:  make([]int64, len(buckets)),
	}
	ms.histograms[key] = histogram
	
	ms.logger.Debug("Created new histogram", "name", name, "labels", labels, "buckets", buckets)
	return histogram
}

// NewTimer creates or retrieves a timer
func (ms *MetricsService) NewTimer(name, help string, labels map[string]string) *Timer {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	
	key := ms.buildKey(name, labels)
	if timer, exists := ms.timers[key]; exists {
		return timer
	}
	
	timer := &Timer{
		name:   name,
		help:   help,
		labels: labels,
	}
	ms.timers[key] = timer
	
	ms.logger.Debug("Created new timer", "name", name, "labels", labels)
	return timer
}

// GetAllMetrics returns all current metric values
func (ms *MetricsService) GetAllMetrics() []MetricValue {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	
	var metrics []MetricValue
	timestamp := time.Now()
	
	// Collect counters
	for _, counter := range ms.counters {
		metrics = append(metrics, MetricValue{
			Name:      counter.name,
			Type:      "counter",
			Help:      counter.help,
			Labels:    counter.labels,
			Value:     counter.Get(),
			Timestamp: timestamp,
		})
	}
	
	// Collect gauges
	for _, gauge := range ms.gauges {
		metrics = append(metrics, MetricValue{
			Name:      gauge.name,
			Type:      "gauge",
			Help:      gauge.help,
			Labels:    gauge.labels,
			Value:     gauge.Get(),
			Timestamp: timestamp,
		})
	}
	
	// Collect histograms
	for _, histogram := range ms.histograms {
		summary := histogram.Summary()
		metrics = append(metrics, MetricValue{
			Name:      histogram.name,
			Type:      "histogram",
			Help:      histogram.help,
			Labels:    histogram.labels,
			Value:     summary,
			Timestamp: timestamp,
		})
	}
	
	// Collect timers
	for _, timer := range ms.timers {
		summary := timer.Summary()
		metrics = append(metrics, MetricValue{
			Name:      timer.name,
			Type:      "timer",
			Help:      timer.help,
			Labels:    timer.labels,
			Value:     summary,
			Timestamp: timestamp,
		})
	}
	
	return metrics
}

// GetPrometheusFormat returns metrics in Prometheus exposition format
func (ms *MetricsService) GetPrometheusFormat() string {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	
	var output string
	
	// Export counters
	for _, counter := range ms.counters {
		output += fmt.Sprintf("# HELP %s %s\n", counter.name, counter.help)
		output += fmt.Sprintf("# TYPE %s counter\n", counter.name)
		labelStr := ms.formatLabels(counter.labels)
		output += fmt.Sprintf("%s%s %d\n", counter.name, labelStr, counter.Get())
	}
	
	// Export gauges
	for _, gauge := range ms.gauges {
		output += fmt.Sprintf("# HELP %s %s\n", gauge.name, gauge.help)
		output += fmt.Sprintf("# TYPE %s gauge\n", gauge.name)
		labelStr := ms.formatLabels(gauge.labels)
		output += fmt.Sprintf("%s%s %f\n", gauge.name, labelStr, gauge.Get())
	}
	
	// Export histograms
	for _, histogram := range ms.histograms {
		output += fmt.Sprintf("# HELP %s %s\n", histogram.name, histogram.help)
		output += fmt.Sprintf("# TYPE %s histogram\n", histogram.name)
		
		summary := histogram.Summary()
		labelStr := ms.formatLabels(histogram.labels)
		
		// Export buckets
		for i, bucket := range histogram.buckets {
			bucketLabels := make(map[string]string)
			for k, v := range histogram.labels {
				bucketLabels[k] = v
			}
			bucketLabels["le"] = fmt.Sprintf("%f", bucket)
			bucketLabelStr := ms.formatLabels(bucketLabels)
			output += fmt.Sprintf("%s_bucket%s %d\n", histogram.name, bucketLabelStr, histogram.counts[i])
		}
		
		// Export sum and count
		output += fmt.Sprintf("%s_sum%s %f\n", histogram.name, labelStr, summary["sum"])
		output += fmt.Sprintf("%s_count%s %d\n", histogram.name, labelStr, summary["count"])
	}
	
	return output
}

// RecordHTTPRequest records HTTP request metrics
func (ms *MetricsService) RecordHTTPRequest(method, path string, statusCode int, duration time.Duration, responseSize int) {
	// HTTP request counter
	httpRequestsCounter := ms.NewCounter(
		"http_requests_total",
		"Total number of HTTP requests",
		map[string]string{
			"method": method,
			"path":   path,
			"status": fmt.Sprintf("%d", statusCode),
		},
	)
	httpRequestsCounter.Inc()
	
	// HTTP request duration histogram
	durationBuckets := []float64{0.1, 0.5, 1.0, 2.5, 5.0, 10.0}
	httpDurationHistogram := ms.NewHistogram(
		"http_request_duration_seconds",
		"HTTP request duration in seconds",
		durationBuckets,
		map[string]string{
			"method": method,
			"path":   path,
		},
	)
	httpDurationHistogram.Observe(duration.Seconds())
	
	// HTTP response size histogram
	sizeBuckets := []float64{100, 1000, 10000, 100000, 1000000}
	httpSizeHistogram := ms.NewHistogram(
		"http_response_size_bytes",
		"HTTP response size in bytes",
		sizeBuckets,
		map[string]string{
			"method": method,
			"path":   path,
		},
	)
	httpSizeHistogram.Observe(float64(responseSize))
}

// RecordDatabaseOperation records database operation metrics
func (ms *MetricsService) RecordDatabaseOperation(operation, table string, duration time.Duration, success bool) {
	// Database operation counter
	status := "success"
	if !success {
		status = "error"
	}
	
	dbOperationsCounter := ms.NewCounter(
		"database_operations_total",
		"Total number of database operations",
		map[string]string{
			"operation": operation,
			"table":     table,
			"status":    status,
		},
	)
	dbOperationsCounter.Inc()
	
	// Database operation duration histogram
	durationBuckets := []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0}
	dbDurationHistogram := ms.NewHistogram(
		"database_operation_duration_seconds",
		"Database operation duration in seconds",
		durationBuckets,
		map[string]string{
			"operation": operation,
			"table":     table,
		},
	)
	dbDurationHistogram.Observe(duration.Seconds())
}

// RecordCacheOperation records cache operation metrics
func (ms *MetricsService) RecordCacheOperation(operation, key string, duration time.Duration, hit bool, success bool) {
	// Cache operation counter
	status := "success"
	if !success {
		status = "error"
	}
	
	cacheOperationsCounter := ms.NewCounter(
		"cache_operations_total",
		"Total number of cache operations",
		map[string]string{
			"operation": operation,
			"status":    status,
		},
	)
	cacheOperationsCounter.Inc()
	
	// Cache hit/miss counter
	if operation == "get" {
		hitStatus := "miss"
		if hit {
			hitStatus = "hit"
		}
		
		cacheHitsCounter := ms.NewCounter(
			"cache_hits_total",
			"Total number of cache hits and misses",
			map[string]string{
				"status": hitStatus,
			},
		)
		cacheHitsCounter.Inc()
	}
	
	// Cache operation duration histogram
	durationBuckets := []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.025, 0.05, 0.1}
	cacheDurationHistogram := ms.NewHistogram(
		"cache_operation_duration_seconds",
		"Cache operation duration in seconds",
		durationBuckets,
		map[string]string{
			"operation": operation,
		},
	)
	cacheDurationHistogram.Observe(duration.Seconds())
}

// UpdateSystemMetrics updates system-level metrics
func (ms *MetricsService) UpdateSystemMetrics(ctx context.Context) {
	// This should be called periodically to update system metrics
	
	// CPU usage (simplified)
	cpuGauge := ms.NewGauge(
		"system_cpu_usage_percent",
		"Current CPU usage percentage",
		nil,
	)
	// In a real implementation, you would get actual CPU usage
	cpuGauge.Set(0.0) // Placeholder
	
	// Memory usage
	memoryGauge := ms.NewGauge(
		"system_memory_usage_bytes",
		"Current memory usage in bytes",
		nil,
	)
	// In a real implementation, you would get actual memory usage
	memoryGauge.Set(0.0) // Placeholder
	
	// Goroutine count
	goroutineGauge := ms.NewGauge(
		"system_goroutines_count",
		"Current number of goroutines",
		nil,
	)
	// This would be set with actual runtime.NumGoroutine()
	goroutineGauge.Set(0.0) // Placeholder
}

// Helper methods

func (ms *MetricsService) buildKey(name string, labels map[string]string) string {
	key := name
	for k, v := range labels {
		key += fmt.Sprintf("_%s_%s", k, v)
	}
	return key
}

func (ms *MetricsService) formatLabels(labels map[string]string) string {
	if len(labels) == 0 {
		return ""
	}
	
	var labelPairs []string
	for k, v := range labels {
		labelPairs = append(labelPairs, fmt.Sprintf(`%s="%s"`, k, v))
	}
	
	return fmt.Sprintf("{%s}", fmt.Sprintf("%v", labelPairs))
}

func (ms *MetricsService) initializeStandardMetrics() {
	// Initialize standard application metrics
	
	// Application uptime
	uptimeGauge := ms.NewGauge(
		"application_uptime_seconds",
		"Application uptime in seconds",
		nil,
	)
	uptimeGauge.Set(0)
	
	// Application info
	infoGauge := ms.NewGauge(
		"application_info",
		"Application information",
		map[string]string{
			"version": "1.0.0", // This should come from build
			"service": "url-shortener",
		},
	)
	infoGauge.Set(1)
}

// Counter methods

func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *Counter) Add(value int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value += value
}

func (c *Counter) Get() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.value
}

// Gauge methods

func (g *Gauge) Set(value float64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.value = value
}

func (g *Gauge) Inc() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.value++
}

func (g *Gauge) Dec() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.value--
}

func (g *Gauge) Add(value float64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.value += value
}

func (g *Gauge) Get() float64 {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.value
}

// Histogram methods

func (h *Histogram) Observe(value float64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	h.sum += value
	h.count++
	
	for i, bucket := range h.buckets {
		if value <= bucket {
			h.counts[i]++
		}
	}
}

func (h *Histogram) Summary() map[string]interface{} {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	return map[string]interface{}{
		"count": h.count,
		"sum":   h.sum,
		"avg":   h.sum / float64(h.count),
	}
}

// Timer methods

func (t *Timer) Start() *TimingContext {
	return &TimingContext{
		timer:     t,
		startTime: time.Now(),
	}
}

func (t *Timer) Record(duration time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.durations = append(t.durations, duration)
}

func (t *Timer) Summary() map[string]interface{} {
	t.mu.RLock()
	defer t.mu.RUnlock()
	
	if len(t.durations) == 0 {
		return map[string]interface{}{
			"count": 0,
			"avg":   0,
			"min":   0,
			"max":   0,
		}
	}
	
	var total time.Duration
	min := t.durations[0]
	max := t.durations[0]
	
	for _, d := range t.durations {
		total += d
		if d < min {
			min = d
		}
		if d > max {
			max = d
		}
	}
	
	return map[string]interface{}{
		"count": len(t.durations),
		"avg":   total / time.Duration(len(t.durations)),
		"min":   min,
		"max":   max,
	}
}

// TimingContext for measuring durations
type TimingContext struct {
	timer     *Timer
	startTime time.Time
}

func (tc *TimingContext) Stop() {
	duration := time.Since(tc.startTime)
	tc.timer.Record(duration)
}