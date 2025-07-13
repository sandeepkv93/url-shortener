package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"url-shortener/internal/api/middleware"
	"url-shortener/internal/core/domain"
	"url-shortener/internal/core/ports"
)

type FunnelHandler struct {
	funnelService ports.FunnelService
}

func NewFunnelHandler(funnelService ports.FunnelService) *FunnelHandler {
	return &FunnelHandler{
		funnelService: funnelService,
	}
}

// Funnel Management

// CreateFunnel handles creating a new conversion funnel
func (h *FunnelHandler) CreateFunnel(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var req domain.CreateFunnelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	funnel, err := h.funnelService.CreateFunnel(r.Context(), userID, req)
	if err != nil {
		switch err {
		case domain.ErrInvalidRequest:
			h.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSONResponse(w, funnel, http.StatusCreated)
}

// GetFunnel handles getting a specific funnel
func (h *FunnelHandler) GetFunnel(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	funnelIDStr := chi.URLParam(r, "id")
	funnelID, err := strconv.ParseUint(funnelIDStr, 10, 32)
	if err != nil {
		h.writeErrorResponse(w, "Invalid funnel ID", http.StatusBadRequest)
		return
	}

	funnel, err := h.funnelService.GetFunnel(r.Context(), uint(funnelID), userID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			h.writeErrorResponse(w, "Funnel not found", http.StatusNotFound)
		case domain.ErrUnauthorized:
			h.writeErrorResponse(w, "Access denied", http.StatusForbidden)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSONResponse(w, funnel, http.StatusOK)
}

// GetUserFunnels handles getting all funnels for a user
func (h *FunnelHandler) GetUserFunnels(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	funnels, err := h.funnelService.GetUserFunnels(r.Context(), userID)
	if err != nil {
		h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"funnels": funnels,
		"total":   len(funnels),
	}

	h.writeJSONResponse(w, response, http.StatusOK)
}

// UpdateFunnel handles updating a funnel
func (h *FunnelHandler) UpdateFunnel(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	funnelIDStr := chi.URLParam(r, "id")
	funnelID, err := strconv.ParseUint(funnelIDStr, 10, 32)
	if err != nil {
		h.writeErrorResponse(w, "Invalid funnel ID", http.StatusBadRequest)
		return
	}

	var req domain.CreateFunnelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	funnel, err := h.funnelService.UpdateFunnel(r.Context(), uint(funnelID), userID, req)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			h.writeErrorResponse(w, "Funnel not found", http.StatusNotFound)
		case domain.ErrUnauthorized:
			h.writeErrorResponse(w, "Access denied", http.StatusForbidden)
		case domain.ErrInvalidRequest:
			h.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSONResponse(w, funnel, http.StatusOK)
}

// DeleteFunnel handles deleting a funnel
func (h *FunnelHandler) DeleteFunnel(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	funnelIDStr := chi.URLParam(r, "id")
	funnelID, err := strconv.ParseUint(funnelIDStr, 10, 32)
	if err != nil {
		h.writeErrorResponse(w, "Invalid funnel ID", http.StatusBadRequest)
		return
	}

	err = h.funnelService.DeleteFunnel(r.Context(), uint(funnelID), userID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			h.writeErrorResponse(w, "Funnel not found", http.StatusNotFound)
		case domain.ErrUnauthorized:
			h.writeErrorResponse(w, "Access denied", http.StatusForbidden)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Funnel Analytics

// GetFunnelAnalytics handles getting analytics for a specific funnel
func (h *FunnelHandler) GetFunnelAnalytics(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	funnelIDStr := chi.URLParam(r, "id")
	funnelID, err := strconv.ParseUint(funnelIDStr, 10, 32)
	if err != nil {
		h.writeErrorResponse(w, "Invalid funnel ID", http.StatusBadRequest)
		return
	}

	period := r.URL.Query().Get("period")
	if period == "" {
		period = "30d"
	}

	// Validate period
	validPeriods := []string{"7d", "30d", "90d", "1y"}
	isValidPeriod := false
	for _, validPeriod := range validPeriods {
		if period == validPeriod {
			isValidPeriod = true
			break
		}
	}
	if !isValidPeriod {
		h.writeErrorResponse(w, "Invalid period. Valid values: 7d, 30d, 90d, 1y", http.StatusBadRequest)
		return
	}

	analytics, err := h.funnelService.GetFunnelAnalytics(r.Context(), uint(funnelID), userID, period)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			h.writeErrorResponse(w, "Funnel not found", http.StatusNotFound)
		case domain.ErrUnauthorized:
			h.writeErrorResponse(w, "Access denied", http.StatusForbidden)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSONResponse(w, analytics, http.StatusOK)
}

// GetFunnelStepAnalytics handles getting analytics for a specific funnel step
func (h *FunnelHandler) GetFunnelStepAnalytics(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	funnelIDStr := chi.URLParam(r, "id")
	funnelID, err := strconv.ParseUint(funnelIDStr, 10, 32)
	if err != nil {
		h.writeErrorResponse(w, "Invalid funnel ID", http.StatusBadRequest)
		return
	}

	stepIDStr := chi.URLParam(r, "stepId")
	stepID, err := strconv.ParseUint(stepIDStr, 10, 32)
	if err != nil {
		h.writeErrorResponse(w, "Invalid step ID", http.StatusBadRequest)
		return
	}

	analytics, err := h.funnelService.GetFunnelStepAnalytics(r.Context(), uint(funnelID), uint(stepID), userID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			h.writeErrorResponse(w, "Funnel or step not found", http.StatusNotFound)
		case domain.ErrUnauthorized:
			h.writeErrorResponse(w, "Access denied", http.StatusForbidden)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSONResponse(w, analytics, http.StatusOK)
}

// GetConversionTrend handles getting conversion trend data for a funnel
func (h *FunnelHandler) GetConversionTrend(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	funnelIDStr := chi.URLParam(r, "id")
	funnelID, err := strconv.ParseUint(funnelIDStr, 10, 32)
	if err != nil {
		h.writeErrorResponse(w, "Invalid funnel ID", http.StatusBadRequest)
		return
	}

	period := r.URL.Query().Get("period")
	if period == "" {
		period = "30d"
	}

	trend, err := h.funnelService.GetConversionTrend(r.Context(), uint(funnelID), userID, period)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			h.writeErrorResponse(w, "Funnel not found", http.StatusNotFound)
		case domain.ErrUnauthorized:
			h.writeErrorResponse(w, "Access denied", http.StatusForbidden)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	response := map[string]interface{}{
		"funnel_id":        funnelID,
		"period":           period,
		"conversion_trend": trend,
	}

	h.writeJSONResponse(w, response, http.StatusOK)
}

// Funnel Optimization

// GetFunnelOptimizations handles getting optimization suggestions for a funnel
func (h *FunnelHandler) GetFunnelOptimizations(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	funnelIDStr := chi.URLParam(r, "id")
	funnelID, err := strconv.ParseUint(funnelIDStr, 10, 32)
	if err != nil {
		h.writeErrorResponse(w, "Invalid funnel ID", http.StatusBadRequest)
		return
	}

	optimizations, err := h.funnelService.GetFunnelOptimizations(r.Context(), uint(funnelID), userID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			h.writeErrorResponse(w, "Funnel not found", http.StatusNotFound)
		case domain.ErrUnauthorized:
			h.writeErrorResponse(w, "Access denied", http.StatusForbidden)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	response := map[string]interface{}{
		"funnel_id":      funnelID,
		"optimizations":  optimizations,
		"count":          len(optimizations),
	}

	h.writeJSONResponse(w, response, http.StatusOK)
}

// AnalyzeFunnelDropOffs handles analyzing drop-offs in a funnel
func (h *FunnelHandler) AnalyzeFunnelDropOffs(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	funnelIDStr := chi.URLParam(r, "id")
	funnelID, err := strconv.ParseUint(funnelIDStr, 10, 32)
	if err != nil {
		h.writeErrorResponse(w, "Invalid funnel ID", http.StatusBadRequest)
		return
	}

	dropOffs, err := h.funnelService.AnalyzeFunnelDropOffs(r.Context(), uint(funnelID), userID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			h.writeErrorResponse(w, "Funnel not found", http.StatusNotFound)
		case domain.ErrUnauthorized:
			h.writeErrorResponse(w, "Access denied", http.StatusForbidden)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	response := map[string]interface{}{
		"funnel_id":        funnelID,
		"drop_off_analysis": dropOffs,
		"total_steps":      len(dropOffs),
		"highest_drop_off": func() *domain.FunnelStepAnalytics {
			if len(dropOffs) > 0 {
				return &dropOffs[0]
			}
			return nil
		}(),
	}

	h.writeJSONResponse(w, response, http.StatusOK)
}

// Advanced Funnel Analytics

// GetFunnelComparisonReport handles comparing multiple funnels
func (h *FunnelHandler) GetFunnelComparisonReport(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Parse funnel IDs from query parameters
	funnelIDs := r.URL.Query()["funnel_ids"]
	if len(funnelIDs) == 0 {
		h.writeErrorResponse(w, "At least one funnel ID is required", http.StatusBadRequest)
		return
	}

	period := r.URL.Query().Get("period")
	if period == "" {
		period = "30d"
	}

	// Get analytics for each funnel
	comparisons := make([]map[string]interface{}, 0, len(funnelIDs))
	
	for _, funnelIDStr := range funnelIDs {
		funnelID, err := strconv.ParseUint(funnelIDStr, 10, 32)
		if err != nil {
			continue // Skip invalid IDs
		}

		analytics, err := h.funnelService.GetFunnelAnalytics(r.Context(), uint(funnelID), userID, period)
		if err != nil {
			continue // Skip funnels that can't be accessed
		}

		funnel, err := h.funnelService.GetFunnel(r.Context(), uint(funnelID), userID)
		if err != nil {
			continue
		}

		comparison := map[string]interface{}{
			"funnel_id":        funnelID,
			"funnel_name":      funnel.Name,
			"conversion_rate":  analytics.ConversionRate,
			"drop_off_rate":    analytics.DropOffRate,
			"total_entries":    analytics.TotalEntries,
			"conversions":      analytics.Conversions,
			"step_count":       len(analytics.StepAnalytics),
		}

		comparisons = append(comparisons, comparison)
	}

	response := map[string]interface{}{
		"period":      period,
		"comparisons": comparisons,
		"count":       len(comparisons),
		"summary": map[string]interface{}{
			"best_performing": func() map[string]interface{} {
				if len(comparisons) == 0 {
					return nil
				}
				best := comparisons[0]
				for _, comp := range comparisons[1:] {
					if comp["conversion_rate"].(float64) > best["conversion_rate"].(float64) {
						best = comp
					}
				}
				return best
			}(),
			"worst_performing": func() map[string]interface{} {
				if len(comparisons) == 0 {
					return nil
				}
				worst := comparisons[0]
				for _, comp := range comparisons[1:] {
					if comp["conversion_rate"].(float64) < worst["conversion_rate"].(float64) {
						worst = comp
					}
				}
				return worst
			}(),
		},
	}

	h.writeJSONResponse(w, response, http.StatusOK)
}

// GetFunnelABTestResults handles getting A/B test results for funnel variations
func (h *FunnelHandler) GetFunnelABTestResults(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	funnelIDStr := chi.URLParam(r, "id")
	funnelID, err := strconv.ParseUint(funnelIDStr, 10, 32)
	if err != nil {
		h.writeErrorResponse(w, "Invalid funnel ID", http.StatusBadRequest)
		return
	}

	// Mock A/B test results
	results := map[string]interface{}{
		"test_id":    "test_" + funnelIDStr,
		"funnel_id":  funnelID,
		"status":     "completed",
		"start_date": "2024-06-01",
		"end_date":   "2024-06-30",
		"variations": []map[string]interface{}{
			{
				"variation":       "A",
				"name":           "Original",
				"traffic_split":  50.0,
				"entries":        5000,
				"conversions":    425,
				"conversion_rate": 8.5,
				"confidence":     95.2,
			},
			{
				"variation":       "B",
				"name":           "Optimized",
				"traffic_split":  50.0,
				"entries":        5000,
				"conversions":    562,
				"conversion_rate": 11.24,
				"confidence":     98.7,
			},
		},
		"winner": map[string]interface{}{
			"variation":     "B",
			"improvement":   32.2,
			"significance":  99.1,
			"recommendation": "Implement variation B for improved conversion rates",
		},
		"statistical_summary": map[string]interface{}{
			"p_value":      0.009,
			"chi_square":   8.235,
			"sample_size":  10000,
			"power":        99.1,
		},
	}

	h.writeJSONResponse(w, results, http.StatusOK)
}

// GetFunnelCohortAnalysis handles getting cohort analysis for a funnel
func (h *FunnelHandler) GetFunnelCohortAnalysis(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	funnelIDStr := chi.URLParam(r, "id")
	funnelID, err := strconv.ParseUint(funnelIDStr, 10, 32)
	if err != nil {
		h.writeErrorResponse(w, "Invalid funnel ID", http.StatusBadRequest)
		return
	}

	period := r.URL.Query().Get("period")
	if period == "" {
		period = "weekly"
	}

	// Mock cohort analysis
	cohorts := []map[string]interface{}{
		{
			"cohort":           "2024-W01",
			"initial_users":    1000,
			"week_0":          100.0,
			"week_1":          85.2,
			"week_2":          72.8,
			"week_3":          65.5,
			"week_4":          58.9,
		},
		{
			"cohort":           "2024-W02",
			"initial_users":    1200,
			"week_0":          100.0,
			"week_1":          87.5,
			"week_2":          75.2,
			"week_3":          68.1,
			"week_4":          61.4,
		},
		{
			"cohort":           "2024-W03",
			"initial_users":    950,
			"week_0":          100.0,
			"week_1":          89.1,
			"week_2":          78.5,
			"week_3":          71.2,
			"week_4":          64.8,
		},
	}

	response := map[string]interface{}{
		"funnel_id": funnelID,
		"period":    period,
		"cohorts":   cohorts,
		"summary": map[string]interface{}{
			"average_retention_week_1": 87.3,
			"average_retention_week_4": 61.7,
			"best_cohort":             "2024-W03",
			"worst_cohort":            "2024-W01",
		},
	}

	h.writeJSONResponse(w, response, http.StatusOK)
}

// Helper methods

func (h *FunnelHandler) writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to encode response"}`))
	}
}

func (h *FunnelHandler) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := map[string]string{"error": message}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.Write([]byte(`{"error": "Internal server error"}`))
	}
}