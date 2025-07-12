package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"url-shortener/internal/core/ports"

	"github.com/go-chi/chi/v5"
)

type HealthHandler struct {
	healthService ports.HealthService
}

func NewHealthHandler(healthService ports.HealthService) *HealthHandler {
	return &HealthHandler{
		healthService: healthService,
	}
}

// GetHealth returns the overall health status of the application
func (h *HealthHandler) GetHealth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	health, err := h.healthService.GetHealth(ctx)
	if err != nil {
		http.Error(w, "Failed to get health status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	
	// Set appropriate status code based on health
	statusCode := http.StatusOK
	switch health.Status {
	case "unhealthy":
		statusCode = http.StatusServiceUnavailable
	case "degraded":
		statusCode = http.StatusPartialContent
	}
	
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(health)
}

// GetHealthChecks returns detailed health check results
func (h *HealthHandler) GetHealthChecks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	checks, err := h.healthService.RunHealthChecks(ctx)
	if err != nil {
		http.Error(w, "Failed to run health checks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	
	// Determine status code based on check results
	statusCode := http.StatusOK
	for _, check := range checks {
		if check.Critical && check.Status == "fail" {
			statusCode = http.StatusServiceUnavailable
			break
		} else if check.Status == "warn" || check.Status == "fail" {
			statusCode = http.StatusPartialContent
		}
	}
	
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"checks": checks,
		"timestamp": checks["database_connectivity"].LastRun,
	})
}

// GetLiveness returns a simple liveness check
func (h *HealthHandler) GetLiveness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "alive",
		"timestamp": r.Header.Get("Date"),
	})
}

// GetReadiness returns a readiness check that verifies critical dependencies
func (h *HealthHandler) GetReadiness(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Check if service is ready (database and cache are accessible)
	ready := h.healthService.IsHealthy(ctx)
	
	w.Header().Set("Content-Type", "application/json")
	
	if ready {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "ready",
			"timestamp": r.Header.Get("Date"),
		})
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "not ready",
			"timestamp": r.Header.Get("Date"),
		})
	}
}

// GetSystemMetrics returns system-level metrics
func (h *HealthHandler) GetSystemMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	metrics, err := h.healthService.GetSystemMetrics(ctx)
	if err != nil {
		http.Error(w, "Failed to get system metrics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(metrics)
}

// GetApplicationMetrics returns application-specific metrics
func (h *HealthHandler) GetApplicationMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	metrics, err := h.healthService.GetApplicationMetrics(ctx)
	if err != nil {
		http.Error(w, "Failed to get application metrics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(metrics)
}

// GetDatabaseHealth returns database health status
func (h *HealthHandler) GetDatabaseHealth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	health, err := h.healthService.CheckDatabaseHealth(ctx)
	if err != nil {
		http.Error(w, "Failed to check database health", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	
	statusCode := http.StatusOK
	if health.Status == "down" {
		statusCode = http.StatusServiceUnavailable
	} else if health.Status == "degraded" {
		statusCode = http.StatusPartialContent
	}
	
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(health)
}

// GetCacheHealth returns cache health status
func (h *HealthHandler) GetCacheHealth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	health, err := h.healthService.CheckCacheHealth(ctx)
	if err != nil {
		http.Error(w, "Failed to check cache health", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	
	statusCode := http.StatusOK
	if health.Status == "down" {
		statusCode = http.StatusServiceUnavailable
	} else if health.Status == "degraded" {
		statusCode = http.StatusPartialContent
	}
	
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(health)
}

// GetExternalServicesHealth returns external services health status
func (h *HealthHandler) GetExternalServicesHealth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	services, err := h.healthService.CheckExternalServices(ctx)
	if err != nil {
		http.Error(w, "Failed to check external services health", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	
	// Determine overall status
	statusCode := http.StatusOK
	overallStatus := "healthy"
	
	for _, service := range services {
		if service.Status == "down" {
			statusCode = http.StatusServiceUnavailable
			overallStatus = "unhealthy"
			break
		} else if service.Status == "degraded" {
			statusCode = http.StatusPartialContent
			overallStatus = "degraded"
		}
	}
	
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": overallStatus,
		"services": services,
		"count": len(services),
	})
}

// GetHealthVersion returns version and build information
func (h *HealthHandler) GetHealthVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	json.NewEncoder(w).Encode(map[string]interface{}{
		"service": "url-shortener",
		"version": "1.0.0",
		"build_time": "2024-01-15T10:00:00Z", // This would be set at build time
		"commit": "latest", // This would be set from git commit
		"go_version": "1.21",
	})
}

// SetupHealthRoutes configures health check routes
func (h *HealthHandler) SetupHealthRoutes() http.Handler {
	r := chi.NewRouter()
	
	// Main health endpoints
	r.Get("/", h.GetHealth)
	r.Get("/health", h.GetHealth)
	r.Get("/healthz", h.GetHealth)
	
	// Kubernetes-style health checks
	r.Get("/livez", h.GetLiveness)
	r.Get("/readyz", h.GetReadiness)
	
	// Detailed health information
	r.Get("/checks", h.GetHealthChecks)
	r.Get("/metrics/system", h.GetSystemMetrics)
	r.Get("/metrics/application", h.GetApplicationMetrics)
	
	// Component-specific health checks
	r.Get("/database", h.GetDatabaseHealth)
	r.Get("/cache", h.GetCacheHealth)
	r.Get("/external", h.GetExternalServicesHealth)
	
	// Version and build information
	r.Get("/version", h.GetHealthVersion)
	r.Get("/info", h.GetHealthVersion)
	
	return r
}

// GetHealthSummary returns a summarized health status for quick checks
func (h *HealthHandler) GetHealthSummary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	isHealthy := h.healthService.IsHealthy(ctx)
	
	w.Header().Set("Content-Type", "application/json")
	
	if isHealthy {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "healthy",
			"ok": true,
		})
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "unhealthy",
			"ok": false,
		})
	}
}

// HealthMiddleware provides health check middleware
func (h *HealthHandler) HealthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add health-related headers
		w.Header().Set("X-Health-Check-Available", "/health")
		w.Header().Set("X-Service-Name", "url-shortener")
		w.Header().Set("X-Service-Version", "1.0.0")
		
		next.ServeHTTP(w, r)
	})
}

// ParseHealthParams parses common health check parameters
func parseHealthParams(r *http.Request) (detailed bool, format string) {
	detailed = false
	format = "json"

	if detailedParam := r.URL.Query().Get("detailed"); detailedParam != "" {
		if parsed, err := strconv.ParseBool(detailedParam); err == nil {
			detailed = parsed
		}
	}

	if formatParam := r.URL.Query().Get("format"); formatParam != "" {
		format = formatParam
	}

	return detailed, format
}