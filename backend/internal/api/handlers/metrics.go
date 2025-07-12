package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"url-shortener/internal/core/services"

	"github.com/go-chi/chi/v5"
)

// MetricsHandler handles metrics endpoints
type MetricsHandler struct {
	metricsService *services.MetricsService
	logger         *services.LoggingService
	errorTracking  *services.ErrorTrackingService
}

// NewMetricsHandler creates a new metrics handler
func NewMetricsHandler(
	metricsService *services.MetricsService,
	logger *services.LoggingService,
	errorTracking *services.ErrorTrackingService,
) *MetricsHandler {
	return &MetricsHandler{
		metricsService: metricsService,
		logger:         logger,
		errorTracking:  errorTracking,
	}
}

// GetMetrics returns all metrics in JSON format
func (mh *MetricsHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	if mh.metricsService == nil {
		http.Error(w, "Metrics service not available", http.StatusServiceUnavailable)
		return
	}

	metrics := mh.metricsService.GetAllMetrics()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"metrics": metrics,
		"count":   len(metrics),
	}); err != nil {
		mh.logger.Error("Failed to encode metrics response", "error", err.Error())
		http.Error(w, "Failed to encode metrics", http.StatusInternalServerError)
	}
}

// GetPrometheusMetrics returns metrics in Prometheus exposition format
func (mh *MetricsHandler) GetPrometheusMetrics(w http.ResponseWriter, r *http.Request) {
	if mh.metricsService == nil {
		http.Error(w, "Metrics service not available", http.StatusServiceUnavailable)
		return
	}

	prometheusFormat := mh.metricsService.GetPrometheusFormat()

	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	
	if _, err := w.Write([]byte(prometheusFormat)); err != nil {
		mh.logger.Error("Failed to write Prometheus metrics", "error", err.Error())
	}
}

// GetErrorSummaries returns error tracking summaries
func (mh *MetricsHandler) GetErrorSummaries(w http.ResponseWriter, r *http.Request) {
	if mh.errorTracking == nil {
		http.Error(w, "Error tracking service not available", http.StatusServiceUnavailable)
		return
	}

	summaries := mh.errorTracking.GetErrorSummaries()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"error_summaries": summaries,
		"count":          len(summaries),
	}); err != nil {
		mh.logger.Error("Failed to encode error summaries response", "error", err.Error())
		http.Error(w, "Failed to encode error summaries", http.StatusInternalServerError)
	}
}

// GetErrorSummary returns a specific error summary
func (mh *MetricsHandler) GetErrorSummary(w http.ResponseWriter, r *http.Request) {
	if mh.errorTracking == nil {
		http.Error(w, "Error tracking service not available", http.StatusServiceUnavailable)
		return
	}

	fingerprint := chi.URLParam(r, "fingerprint")
	if fingerprint == "" {
		http.Error(w, "Fingerprint parameter required", http.StatusBadRequest)
		return
	}

	summary, exists := mh.errorTracking.GetErrorSummary(fingerprint)
	if !exists {
		http.Error(w, "Error summary not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	if err := json.NewEncoder(w).Encode(summary); err != nil {
		mh.logger.Error("Failed to encode error summary response", "error", err.Error())
		http.Error(w, "Failed to encode error summary", http.StatusInternalServerError)
	}
}

// ResolveError marks an error as resolved
func (mh *MetricsHandler) ResolveError(w http.ResponseWriter, r *http.Request) {
	if mh.errorTracking == nil {
		http.Error(w, "Error tracking service not available", http.StatusServiceUnavailable)
		return
	}

	fingerprint := chi.URLParam(r, "fingerprint")
	if fingerprint == "" {
		http.Error(w, "Fingerprint parameter required", http.StatusBadRequest)
		return
	}

	var request struct {
		ResolvedBy string `json:"resolved_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.ResolvedBy == "" {
		http.Error(w, "resolved_by field required", http.StatusBadRequest)
		return
	}

	if err := mh.errorTracking.ResolveError(fingerprint, request.ResolvedBy); err != nil {
		mh.logger.Error("Failed to resolve error", "error", err.Error(), "fingerprint", fingerprint)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Error resolved successfully",
		"fingerprint": fingerprint,
		"resolved_by": request.ResolvedBy,
	})
}

// GetAlerts returns all active alerts
func (mh *MetricsHandler) GetAlerts(w http.ResponseWriter, r *http.Request) {
	if mh.errorTracking == nil {
		http.Error(w, "Error tracking service not available", http.StatusServiceUnavailable)
		return
	}

	alerts := mh.errorTracking.GetAlerts()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"alerts": alerts,
		"count":  len(alerts),
	}); err != nil {
		mh.logger.Error("Failed to encode alerts response", "error", err.Error())
		http.Error(w, "Failed to encode alerts", http.StatusInternalServerError)
	}
}

// GetErrorStats returns error tracking statistics
func (mh *MetricsHandler) GetErrorStats(w http.ResponseWriter, r *http.Request) {
	if mh.errorTracking == nil {
		http.Error(w, "Error tracking service not available", http.StatusServiceUnavailable)
		return
	}

	stats := mh.errorTracking.GetErrorStats()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		mh.logger.Error("Failed to encode error stats response", "error", err.Error())
		http.Error(w, "Failed to encode error stats", http.StatusInternalServerError)
	}
}

// TestError creates a test error for testing error tracking
func (mh *MetricsHandler) TestError(w http.ResponseWriter, r *http.Request) {
	if mh.errorTracking == nil {
		http.Error(w, "Error tracking service not available", http.StatusServiceUnavailable)
		return
	}

	// Get severity from query parameter
	severityParam := r.URL.Query().Get("severity")
	severity := services.SeverityMedium
	
	switch severityParam {
	case "critical":
		severity = services.SeverityCritical
	case "high":
		severity = services.SeverityHigh
	case "medium":
		severity = services.SeverityMedium
	case "low":
		severity = services.SeverityLow
	case "info":
		severity = services.SeverityInfo
	}

	// Create test error
	testError := fmt.Errorf("test error for monitoring - severity: %s", severity)
	
	mh.errorTracking.TrackError(
		r.Context(),
		testError,
		"test_component",
		severity,
		map[string]interface{}{
			"test":        true,
			"timestamp":   time.Now().Unix(),
			"user_agent":  r.Header.Get("User-Agent"),
			"remote_addr": r.RemoteAddr,
		},
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Test error created successfully",
		"severity": severity,
		"error": testError.Error(),
	})
}

// GetMonitoringDashboard returns a simple monitoring dashboard
func (mh *MetricsHandler) GetMonitoringDashboard(w http.ResponseWriter, r *http.Request) {
	// Collect all monitoring data
	var metrics []services.MetricValue
	var errorSummaries []*services.ErrorSummary
	var alerts []*services.Alert
	var errorStats map[string]interface{}

	if mh.metricsService != nil {
		metrics = mh.metricsService.GetAllMetrics()
	}

	if mh.errorTracking != nil {
		errorSummaries = mh.errorTracking.GetErrorSummaries()
		alerts = mh.errorTracking.GetAlerts()
		errorStats = mh.errorTracking.GetErrorStats()
	}

	dashboard := map[string]interface{}{
		"service":        "url-shortener",
		"version":        "1.0.0",
		"timestamp":      time.Now(),
		"metrics": map[string]interface{}{
			"count": len(metrics),
			"data":  metrics,
		},
		"errors": map[string]interface{}{
			"summaries_count": len(errorSummaries),
			"alerts_count":    len(alerts),
			"stats":          errorStats,
			"recent_errors":  mh.getRecentErrors(errorSummaries, 10),
		},
		"alerts": map[string]interface{}{
			"count":      len(alerts),
			"active":     mh.getActiveAlerts(alerts),
			"recent":     mh.getRecentAlerts(alerts, 5),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	if err := json.NewEncoder(w).Encode(dashboard); err != nil {
		mh.logger.Error("Failed to encode dashboard response", "error", err.Error())
		http.Error(w, "Failed to encode dashboard", http.StatusInternalServerError)
	}
}

// SetupMetricsRoutes configures metrics routes
func (mh *MetricsHandler) SetupMetricsRoutes() http.Handler {
	r := chi.NewRouter()

	// Metrics endpoints
	r.Get("/", mh.GetMetrics)
	r.Get("/metrics", mh.GetMetrics)
	r.Get("/prometheus", mh.GetPrometheusMetrics)
	
	// Error tracking endpoints
	r.Get("/errors", mh.GetErrorSummaries)
	r.Get("/errors/{fingerprint}", mh.GetErrorSummary)
	r.Post("/errors/{fingerprint}/resolve", mh.ResolveError)
	r.Get("/errors/stats", mh.GetErrorStats)
	
	// Alerts endpoints
	r.Get("/alerts", mh.GetAlerts)
	
	// Testing endpoints (for development)
	r.Post("/test/error", mh.TestError)
	
	// Dashboard endpoint
	r.Get("/dashboard", mh.GetMonitoringDashboard)

	return r
}

// Helper methods

func (mh *MetricsHandler) getRecentErrors(summaries []*services.ErrorSummary, limit int) []*services.ErrorSummary {
	if len(summaries) <= limit {
		return summaries
	}

	// Sort by last seen (most recent first)
	sorted := make([]*services.ErrorSummary, len(summaries))
	copy(sorted, summaries)
	
	// Simple sort by LastSeen (in a real implementation, use sort.Slice)
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i].LastSeen.Before(sorted[j].LastSeen) {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	return sorted[:limit]
}

func (mh *MetricsHandler) getActiveAlerts(alerts []*services.Alert) []*services.Alert {
	var active []*services.Alert
	for _, alert := range alerts {
		if !alert.Resolved {
			active = append(active, alert)
		}
	}
	return active
}

func (mh *MetricsHandler) getRecentAlerts(alerts []*services.Alert, limit int) []*services.Alert {
	active := mh.getActiveAlerts(alerts)
	if len(active) <= limit {
		return active
	}
	return active[:limit]
}