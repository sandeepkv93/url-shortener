package handlers

import (
	"encoding/json"
	"net/http"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"url-shortener/internal/config"
	"url-shortener/internal/core/services"
	"url-shortener/internal/infrastructure/performance"

	"github.com/go-chi/chi/v5"
)

// PerformanceHandler handles performance monitoring and optimization endpoints
type PerformanceHandler struct {
	config             *config.Config
	logger             *services.LoggingService
	optimizationService *performance.OptimizationService
}

// NewPerformanceHandler creates a new performance handler
func NewPerformanceHandler(cfg *config.Config, logger *services.LoggingService, optimizationService *performance.OptimizationService) *PerformanceHandler {
	return &PerformanceHandler{
		config:              cfg,
		logger:              logger,
		optimizationService: optimizationService,
	}
}

// RegisterRoutes registers performance monitoring routes
func (ph *PerformanceHandler) RegisterRoutes(r chi.Router) {
	r.Route("/performance", func(r chi.Router) {
		r.Get("/report", ph.GetPerformanceReport)
		r.Get("/metrics", ph.GetMetrics)
		r.Get("/runtime", ph.GetRuntimeInfo)
		r.Get("/memory", ph.GetMemoryInfo)
		r.Get("/gc", ph.GetGCInfo)
		r.Post("/gc/trigger", ph.TriggerGC)
		r.Get("/optimizations", ph.GetOptimizations)
		r.Post("/optimizations/apply", ph.ApplyOptimizations)
		r.Get("/health", ph.GetHealthStatus)
		r.Get("/dashboard", ph.GetDashboard)
	})
}

// GetPerformanceReport returns a comprehensive performance report
func (ph *PerformanceHandler) GetPerformanceReport(w http.ResponseWriter, r *http.Request) {
	report := ph.optimizationService.GetPerformanceReport()
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(report); err != nil {
		ph.logger.Error("Failed to encode performance report", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetMetrics returns current performance metrics
func (ph *PerformanceHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	metrics := map[string]interface{}{
		"timestamp": time.Now(),
		"runtime": map[string]interface{}{
			"goroutines":     runtime.NumGoroutine(),
			"gc_cycles":      m.NumGC,
			"gomaxprocs":     runtime.GOMAXPROCS(0),
			"num_cpu":        runtime.NumCPU(),
			"version":        runtime.Version(),
		},
		"memory": map[string]interface{}{
			"alloc":           m.Alloc,
			"total_alloc":     m.TotalAlloc,
			"sys":             m.Sys,
			"heap_alloc":      m.HeapAlloc,
			"heap_sys":        m.HeapSys,
			"heap_idle":       m.HeapIdle,
			"heap_inuse":      m.HeapInuse,
			"heap_released":   m.HeapReleased,
			"heap_objects":    m.HeapObjects,
			"stack_inuse":     m.StackInuse,
			"stack_sys":       m.StackSys,
			"gc_target_percent": debug.SetGCPercent(-1),
		},
		"gc": map[string]interface{}{
			"num_gc":           m.NumGC,
			"num_forced_gc":    m.NumForcedGC,
			"gc_cpu_fraction":  m.GCCPUFraction,
			"last_gc":          time.Unix(0, int64(m.LastGC)),
			"pause_total_ns":   m.PauseTotalNs,
			"pause_ns":         m.PauseNs,
		},
	}
	
	// Reset GC percent to original value
	debug.SetGCPercent(100)
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		ph.logger.Error("Failed to encode metrics", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetRuntimeInfo returns Go runtime information
func (ph *PerformanceHandler) GetRuntimeInfo(w http.ResponseWriter, r *http.Request) {
	runtimeInfo := map[string]interface{}{
		"version":       runtime.Version(),
		"goos":          runtime.GOOS,
		"goarch":        runtime.GOARCH,
		"compiler":      runtime.Compiler,
		"num_cpu":       runtime.NumCPU(),
		"gomaxprocs":    runtime.GOMAXPROCS(0),
		"num_goroutine": runtime.NumGoroutine(),
		"num_cgo_call":  runtime.NumCgoCall(),
	}
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(runtimeInfo); err != nil {
		ph.logger.Error("Failed to encode runtime info", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetMemoryInfo returns detailed memory information
func (ph *PerformanceHandler) GetMemoryInfo(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	memoryInfo := map[string]interface{}{
		"general": map[string]interface{}{
			"alloc":       m.Alloc,
			"total_alloc": m.TotalAlloc,
			"sys":         m.Sys,
			"lookups":     m.Lookups,
			"mallocs":     m.Mallocs,
			"frees":       m.Frees,
		},
		"heap": map[string]interface{}{
			"heap_alloc":    m.HeapAlloc,
			"heap_sys":      m.HeapSys,
			"heap_idle":     m.HeapIdle,
			"heap_inuse":    m.HeapInuse,
			"heap_released": m.HeapReleased,
			"heap_objects":  m.HeapObjects,
		},
		"stack": map[string]interface{}{
			"stack_inuse": m.StackInuse,
			"stack_sys":   m.StackSys,
		},
		"off_heap": map[string]interface{}{
			"mspan_inuse": m.MSpanInuse,
			"mspan_sys":   m.MSpanSys,
			"mcache_inuse": m.MCacheInuse,
			"mcache_sys":   m.MCacheSys,
			"buck_hash_sys": m.BuckHashSys,
			"gc_sys":        m.GCSys,
			"other_sys":     m.OtherSys,
		},
		"optimizations": ph.config.GetProductionOptimization().Memory,
	}
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(memoryInfo); err != nil {
		ph.logger.Error("Failed to encode memory info", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetGCInfo returns garbage collection information
func (ph *PerformanceHandler) GetGCInfo(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	gcInfo := map[string]interface{}{
		"num_gc":           m.NumGC,
		"num_forced_gc":    m.NumForcedGC,
		"gc_cpu_fraction":  m.GCCPUFraction,
		"last_gc":          time.Unix(0, int64(m.LastGC)),
		"pause_total_ns":   m.PauseTotalNs,
		"pause_end":        m.PauseEnd,
		"pause_ns":         m.PauseNs[:m.NumGC%256],
		"gc_target_percent": debug.SetGCPercent(-1),
		"next_gc":          m.NextGC,
		"enable_gc":        m.EnableGC,
		"debug_gc":         m.DebugGC,
	}
	
	// Reset GC percent
	debug.SetGCPercent(100)
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(gcInfo); err != nil {
		ph.logger.Error("Failed to encode GC info", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// TriggerGC manually triggers garbage collection
func (ph *PerformanceHandler) TriggerGC(w http.ResponseWriter, r *http.Request) {
	var beforeStats, afterStats runtime.MemStats
	runtime.ReadMemStats(&beforeStats)
	
	startTime := time.Now()
	runtime.GC()
	gcDuration := time.Since(startTime)
	
	runtime.ReadMemStats(&afterStats)
	
	result := map[string]interface{}{
		"timestamp":     time.Now(),
		"gc_duration":   gcDuration,
		"before": map[string]interface{}{
			"alloc":       beforeStats.Alloc,
			"heap_alloc":  beforeStats.HeapAlloc,
			"num_gc":      beforeStats.NumGC,
		},
		"after": map[string]interface{}{
			"alloc":       afterStats.Alloc,
			"heap_alloc":  afterStats.HeapAlloc,
			"num_gc":      afterStats.NumGC,
		},
		"memory_freed": int64(beforeStats.Alloc) - int64(afterStats.Alloc),
	}
	
	ph.logger.Info("Manual GC triggered",
		"duration", gcDuration,
		"memory_freed", result["memory_freed"],
	)
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		ph.logger.Error("Failed to encode GC result", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetOptimizations returns current optimization settings
func (ph *PerformanceHandler) GetOptimizations(w http.ResponseWriter, r *http.Request) {
	optimizations := ph.config.GetProductionOptimization()
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(optimizations); err != nil {
		ph.logger.Error("Failed to encode optimizations", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// ApplyOptimizations applies optimization settings (for runtime adjustments)
func (ph *PerformanceHandler) ApplyOptimizations(w http.ResponseWriter, r *http.Request) {
	var request struct {
		GCPercent   *int    `json:"gc_percent,omitempty"`
		MaxProcs    *int    `json:"max_procs,omitempty"`
		ForceGC     bool    `json:"force_gc,omitempty"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	results := make(map[string]interface{})
	
	// Apply GC percent change
	if request.GCPercent != nil {
		oldPercent := debug.SetGCPercent(*request.GCPercent)
		results["gc_percent"] = map[string]interface{}{
			"old_value": oldPercent,
			"new_value": *request.GCPercent,
		}
		ph.logger.Info("GC percent updated",
			"old_value", oldPercent,
			"new_value", *request.GCPercent,
		)
	}
	
	// Apply GOMAXPROCS change
	if request.MaxProcs != nil {
		oldProcs := runtime.GOMAXPROCS(*request.MaxProcs)
		results["max_procs"] = map[string]interface{}{
			"old_value": oldProcs,
			"new_value": *request.MaxProcs,
		}
		ph.logger.Info("GOMAXPROCS updated",
			"old_value", oldProcs,
			"new_value", *request.MaxProcs,
		)
	}
	
	// Force GC if requested
	if request.ForceGC {
		startTime := time.Now()
		runtime.GC()
		duration := time.Since(startTime)
		results["forced_gc"] = map[string]interface{}{
			"duration": duration,
			"timestamp": time.Now(),
		}
		ph.logger.Info("Forced GC completed", "duration", duration)
	}
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		ph.logger.Error("Failed to encode optimization results", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetHealthStatus returns performance-based health status
func (ph *PerformanceHandler) GetHealthStatus(w http.ResponseWriter, r *http.Request) {
	report := ph.optimizationService.GetPerformanceReport()
	
	healthStatus := map[string]interface{}{
		"status":          report.HealthStatus,
		"timestamp":       report.Timestamp,
		"recommendations": report.Recommendations,
		"metrics": map[string]interface{}{
			"memory_usage":     report.ResourceMetrics.MemoryPercentage,
			"goroutines":       report.ResourceMetrics.GoroutineCount,
			"error_rate":       report.RequestMetrics.ErrorRate,
			"avg_latency":      report.RequestMetrics.AverageLatency,
			"active_requests":  report.RequestMetrics.ActiveRequests,
		},
	}
	
	// Set appropriate HTTP status based on health
	statusCode := http.StatusOK
	switch report.HealthStatus {
	case "CRITICAL":
		statusCode = http.StatusServiceUnavailable
	case "WARNING":
		statusCode = http.StatusAccepted
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(healthStatus); err != nil {
		ph.logger.Error("Failed to encode health status", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GetDashboard returns a performance dashboard (HTML view)
func (ph *PerformanceHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	// In a production environment, you might serve a static HTML file or template
	// For now, we'll return a simple JSON response that can be used by a frontend
	
	report := ph.optimizationService.GetPerformanceReport()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	dashboard := map[string]interface{}{
		"title": "URL Shortener Performance Dashboard",
		"timestamp": time.Now(),
		"performance_report": report,
		"quick_stats": map[string]interface{}{
			"uptime":         time.Since(time.Now().Add(-1 * time.Hour)), // Placeholder
			"total_requests": report.RequestMetrics.TotalRequests,
			"avg_latency":    report.RequestMetrics.AverageLatency,
			"memory_usage":   report.ResourceMetrics.MemoryPercentage,
			"goroutines":     report.ResourceMetrics.GoroutineCount,
			"gc_cycles":      report.ResourceMetrics.GCCycles,
		},
		"charts_data": map[string]interface{}{
			"latency_percentiles": map[string]interface{}{
				"p50": report.RequestMetrics.P50Latency,
				"p95": report.RequestMetrics.P95Latency,
				"p99": report.RequestMetrics.P99Latency,
			},
			"memory_breakdown": map[string]interface{}{
				"heap":  m.HeapAlloc,
				"stack": m.StackInuse,
				"sys":   m.Sys,
			},
		},
		"endpoints": map[string]string{
			"metrics":       "/api/performance/metrics",
			"runtime":       "/api/performance/runtime",
			"memory":        "/api/performance/memory",
			"gc":            "/api/performance/gc",
			"health":        "/api/performance/health",
			"optimizations": "/api/performance/optimizations",
		},
	}
	
	// Check if client wants HTML or JSON
	accept := r.Header.Get("Accept")
	if accept == "text/html" || r.URL.Query().Get("format") == "html" {
		ph.renderHTMLDashboard(w, dashboard)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(dashboard); err != nil {
		ph.logger.Error("Failed to encode dashboard", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// renderHTMLDashboard renders a simple HTML dashboard
func (ph *PerformanceHandler) renderHTMLDashboard(w http.ResponseWriter, data map[string]interface{}) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>Performance Dashboard</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .metric { margin: 10px 0; padding: 10px; background: #f0f0f0; border-radius: 5px; }
        .critical { background: #ffebee; }
        .warning { background: #fff3e0; }
        .healthy { background: #e8f5e8; }
        .recommendations { margin-top: 20px; }
        .recommendation { padding: 5px; margin: 5px 0; background: #fff59d; border-radius: 3px; }
    </style>
    <script>
        function refreshData() {
            location.reload();
        }
        setInterval(refreshData, 30000); // Refresh every 30 seconds
    </script>
</head>
<body>
    <h1>URL Shortener Performance Dashboard</h1>
    <p>Last updated: ` + time.Now().Format(time.RFC3339) + `</p>
    
    <h2>Quick Stats</h2>
    <div class="metrics">
        <div class="metric">Memory Usage: ` + strconv.FormatFloat(data["performance_report"].(*performance.PerformanceReport).ResourceMetrics.MemoryPercentage*100, 'f', 2, 64) + `%</div>
        <div class="metric">Active Goroutines: ` + strconv.Itoa(data["performance_report"].(*performance.PerformanceReport).ResourceMetrics.GoroutineCount) + `</div>
        <div class="metric">Total Requests: ` + strconv.FormatInt(data["performance_report"].(*performance.PerformanceReport).RequestMetrics.TotalRequests, 10) + `</div>
        <div class="metric">Average Latency: ` + data["performance_report"].(*performance.PerformanceReport).RequestMetrics.AverageLatency.String() + `</div>
    </div>
    
    <h2>Health Status</h2>
    <div class="metric ` + strings.ToLower(data["performance_report"].(*performance.PerformanceReport).HealthStatus) + `">
        Status: ` + data["performance_report"].(*performance.PerformanceReport).HealthStatus + `
    </div>
    
    <div class="recommendations">
        <h3>Recommendations</h3>`
    
    for _, rec := range data["performance_report"].(*performance.PerformanceReport).Recommendations {
        html += `<div class="recommendation">` + rec + `</div>`
    }
    
    html += `</div>
    
    <h2>API Endpoints</h2>
    <ul>
        <li><a href="/api/performance/metrics">Metrics</a></li>
        <li><a href="/api/performance/runtime">Runtime Info</a></li>
        <li><a href="/api/performance/memory">Memory Info</a></li>
        <li><a href="/api/performance/gc">GC Info</a></li>
        <li><a href="/api/performance/health">Health Status</a></li>
    </ul>
</body>
</html>`
	
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}