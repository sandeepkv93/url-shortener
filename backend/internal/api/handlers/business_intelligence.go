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

type BusinessIntelligenceHandler struct {
	biService      ports.BusinessIntelligenceService
	funnelService  ports.FunnelService
	reportService  ports.ReportingService
}

func NewBusinessIntelligenceHandler(
	biService ports.BusinessIntelligenceService,
	funnelService ports.FunnelService,
	reportService ports.ReportingService,
) *BusinessIntelligenceHandler {
	return &BusinessIntelligenceHandler{
		biService:     biService,
		funnelService: funnelService,
		reportService: reportService,
	}
}

// Dashboard Management

// CreateDashboard handles creating a new BI dashboard
func (h *BusinessIntelligenceHandler) CreateDashboard(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var req domain.CreateDashboardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	dashboard, err := h.biService.CreateDashboard(r.Context(), userID, req)
	if err != nil {
		switch err {
		case domain.ErrInvalidRequest:
			h.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSONResponse(w, dashboard, http.StatusCreated)
}

// GetDashboard handles getting a specific dashboard
func (h *BusinessIntelligenceHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	dashboardIDStr := chi.URLParam(r, "id")
	dashboardID, err := strconv.ParseUint(dashboardIDStr, 10, 32)
	if err != nil {
		h.writeErrorResponse(w, "Invalid dashboard ID", http.StatusBadRequest)
		return
	}

	dashboard, err := h.biService.GetDashboard(r.Context(), uint(dashboardID), userID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			h.writeErrorResponse(w, "Dashboard not found", http.StatusNotFound)
		case domain.ErrUnauthorized:
			h.writeErrorResponse(w, "Access denied", http.StatusForbidden)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSONResponse(w, dashboard, http.StatusOK)
}

// GetUserDashboards handles getting all dashboards for a user
func (h *BusinessIntelligenceHandler) GetUserDashboards(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	dashboards, err := h.biService.GetUserDashboards(r.Context(), userID)
	if err != nil {
		h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := domain.DashboardListResponse{
		Dashboards: dashboards,
		Total:      int64(len(dashboards)),
		Page:       1,
		PageSize:   len(dashboards),
		TotalPages: 1,
	}

	h.writeJSONResponse(w, response, http.StatusOK)
}

// UpdateDashboard handles updating a dashboard
func (h *BusinessIntelligenceHandler) UpdateDashboard(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	dashboardIDStr := chi.URLParam(r, "id")
	dashboardID, err := strconv.ParseUint(dashboardIDStr, 10, 32)
	if err != nil {
		h.writeErrorResponse(w, "Invalid dashboard ID", http.StatusBadRequest)
		return
	}

	var req domain.UpdateDashboardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	dashboard, err := h.biService.UpdateDashboard(r.Context(), uint(dashboardID), userID, req)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			h.writeErrorResponse(w, "Dashboard not found", http.StatusNotFound)
		case domain.ErrUnauthorized:
			h.writeErrorResponse(w, "Access denied", http.StatusForbidden)
		case domain.ErrInvalidRequest:
			h.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSONResponse(w, dashboard, http.StatusOK)
}

// DeleteDashboard handles deleting a dashboard
func (h *BusinessIntelligenceHandler) DeleteDashboard(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	dashboardIDStr := chi.URLParam(r, "id")
	dashboardID, err := strconv.ParseUint(dashboardIDStr, 10, 32)
	if err != nil {
		h.writeErrorResponse(w, "Invalid dashboard ID", http.StatusBadRequest)
		return
	}

	err = h.biService.DeleteDashboard(r.Context(), uint(dashboardID), userID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			h.writeErrorResponse(w, "Dashboard not found", http.StatusNotFound)
		case domain.ErrUnauthorized:
			h.writeErrorResponse(w, "Access denied", http.StatusForbidden)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Widget Management

// CreateWidget handles creating a new widget
func (h *BusinessIntelligenceHandler) CreateWidget(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var req domain.CreateWidgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	widget, err := h.biService.CreateWidget(r.Context(), userID, req)
	if err != nil {
		switch err {
		case domain.ErrInvalidRequest:
			h.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSONResponse(w, widget, http.StatusCreated)
}

// GetWidget handles getting a specific widget
func (h *BusinessIntelligenceHandler) GetWidget(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	widgetIDStr := chi.URLParam(r, "id")
	widgetID, err := strconv.ParseUint(widgetIDStr, 10, 32)
	if err != nil {
		h.writeErrorResponse(w, "Invalid widget ID", http.StatusBadRequest)
		return
	}

	widget, err := h.biService.GetWidget(r.Context(), uint(widgetID), userID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			h.writeErrorResponse(w, "Widget not found", http.StatusNotFound)
		case domain.ErrUnauthorized:
			h.writeErrorResponse(w, "Access denied", http.StatusForbidden)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSONResponse(w, widget, http.StatusOK)
}

// GetWidgetData handles getting data for a specific widget
func (h *BusinessIntelligenceHandler) GetWidgetData(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	widgetIDStr := chi.URLParam(r, "id")
	widgetID, err := strconv.ParseUint(widgetIDStr, 10, 32)
	if err != nil {
		h.writeErrorResponse(w, "Invalid widget ID", http.StatusBadRequest)
		return
	}

	data, err := h.biService.GetWidgetData(r.Context(), uint(widgetID), userID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			h.writeErrorResponse(w, "Widget not found", http.StatusNotFound)
		case domain.ErrUnauthorized:
			h.writeErrorResponse(w, "Access denied", http.StatusForbidden)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSONResponse(w, data, http.StatusOK)
}

// UpdateWidget handles updating a widget
func (h *BusinessIntelligenceHandler) UpdateWidget(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	widgetIDStr := chi.URLParam(r, "id")
	widgetID, err := strconv.ParseUint(widgetIDStr, 10, 32)
	if err != nil {
		h.writeErrorResponse(w, "Invalid widget ID", http.StatusBadRequest)
		return
	}

	var req domain.UpdateWidgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	widget, err := h.biService.UpdateWidget(r.Context(), uint(widgetID), userID, req)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			h.writeErrorResponse(w, "Widget not found", http.StatusNotFound)
		case domain.ErrUnauthorized:
			h.writeErrorResponse(w, "Access denied", http.StatusForbidden)
		case domain.ErrInvalidRequest:
			h.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSONResponse(w, widget, http.StatusOK)
}

// DeleteWidget handles deleting a widget
func (h *BusinessIntelligenceHandler) DeleteWidget(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	widgetIDStr := chi.URLParam(r, "id")
	widgetID, err := strconv.ParseUint(widgetIDStr, 10, 32)
	if err != nil {
		h.writeErrorResponse(w, "Invalid widget ID", http.StatusBadRequest)
		return
	}

	err = h.biService.DeleteWidget(r.Context(), uint(widgetID), userID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			h.writeErrorResponse(w, "Widget not found", http.StatusNotFound)
		case domain.ErrUnauthorized:
			h.writeErrorResponse(w, "Access denied", http.StatusForbidden)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Advanced Analytics

// GetAdvancedAnalytics handles getting comprehensive analytics
func (h *BusinessIntelligenceHandler) GetAdvancedAnalytics(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	period := r.URL.Query().Get("period")
	if period == "" {
		period = "30d"
	}

	analytics, err := h.biService.GetAdvancedAnalytics(r.Context(), userID, period)
	if err != nil {
		h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, analytics, http.StatusOK)
}

// GetPerformanceMetrics handles getting performance metrics
func (h *BusinessIntelligenceHandler) GetPerformanceMetrics(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	period := r.URL.Query().Get("period")
	if period == "" {
		period = "30d"
	}

	metrics, err := h.biService.GetPerformanceMetrics(r.Context(), userID, period)
	if err != nil {
		h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, metrics, http.StatusOK)
}

// GetAudienceInsights handles getting audience insights
func (h *BusinessIntelligenceHandler) GetAudienceInsights(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	period := r.URL.Query().Get("period")
	if period == "" {
		period = "30d"
	}

	insights, err := h.biService.GetAudienceInsights(r.Context(), userID, period)
	if err != nil {
		h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, insights, http.StatusOK)
}

// GetContentAnalytics handles getting content analytics
func (h *BusinessIntelligenceHandler) GetContentAnalytics(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	period := r.URL.Query().Get("period")
	if period == "" {
		period = "30d"
	}

	analytics, err := h.biService.GetContentAnalytics(r.Context(), userID, period)
	if err != nil {
		h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, analytics, http.StatusOK)
}

// Competitive Analysis

// GetCompetitiveAnalysis handles getting competitive analysis
func (h *BusinessIntelligenceHandler) GetCompetitiveAnalysis(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	analysis, err := h.biService.GetCompetitiveAnalysis(r.Context(), userID)
	if err != nil {
		h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, analysis, http.StatusOK)
}

// GetMarketPosition handles getting market position analysis
func (h *BusinessIntelligenceHandler) GetMarketPosition(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	position, err := h.biService.GetMarketPosition(r.Context(), userID)
	if err != nil {
		h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, position, http.StatusOK)
}

// GetBenchmarkData handles getting benchmark data
func (h *BusinessIntelligenceHandler) GetBenchmarkData(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	metric := r.URL.Query().Get("metric")
	if metric == "" {
		metric = "all"
	}

	benchmarks, err := h.biService.GetBenchmarkData(r.Context(), userID, metric)
	if err != nil {
		h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, benchmarks, http.StatusOK)
}

// Predictive Insights

// GetPredictiveInsights handles getting predictive insights
func (h *BusinessIntelligenceHandler) GetPredictiveInsights(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	insights, err := h.biService.GetPredictiveInsights(r.Context(), userID)
	if err != nil {
		h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, insights, http.StatusOK)
}

// GetForecastData handles getting forecast data
func (h *BusinessIntelligenceHandler) GetForecastData(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	metric := r.URL.Query().Get("metric")
	if metric == "" {
		metric = "clicks"
	}

	period := r.URL.Query().Get("period")
	if period == "" {
		period = "30d"
	}

	forecast, err := h.biService.GetForecastData(r.Context(), userID, metric, period)
	if err != nil {
		h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, forecast, http.StatusOK)
}

// DetectAnomalies handles detecting anomalies
func (h *BusinessIntelligenceHandler) DetectAnomalies(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	metric := r.URL.Query().Get("metric")
	if metric == "" {
		metric = "clicks"
	}

	anomalies, err := h.biService.DetectAnomalies(r.Context(), userID, metric)
	if err != nil {
		h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"anomalies": anomalies,
		"metric":    metric,
		"count":     len(anomalies),
	}

	h.writeJSONResponse(w, response, http.StatusOK)
}

// GetTrendPrediction handles getting trend predictions
func (h *BusinessIntelligenceHandler) GetTrendPrediction(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	metric := r.URL.Query().Get("metric")
	if metric == "" {
		metric = "clicks"
	}

	prediction, err := h.biService.GetTrendPrediction(r.Context(), userID, metric)
	if err != nil {
		h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, prediction, http.StatusOK)
}

// Recommendations

// GetRecommendations handles getting all recommendations
func (h *BusinessIntelligenceHandler) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	recommendations, err := h.biService.GetRecommendations(r.Context(), userID)
	if err != nil {
		h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, recommendations, http.StatusOK)
}

// GetOptimizationSuggestions handles getting optimization suggestions
func (h *BusinessIntelligenceHandler) GetOptimizationSuggestions(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	suggestions, err := h.biService.GetOptimizationSuggestions(r.Context(), userID)
	if err != nil {
		h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"suggestions": suggestions,
		"count":       len(suggestions),
	}

	h.writeJSONResponse(w, response, http.StatusOK)
}

// GetContentRecommendations handles getting content recommendations
func (h *BusinessIntelligenceHandler) GetContentRecommendations(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	recommendations, err := h.biService.GetContentRecommendations(r.Context(), userID)
	if err != nil {
		h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"recommendations": recommendations,
		"count":           len(recommendations),
	}

	h.writeJSONResponse(w, response, http.StatusOK)
}

// GetAudienceRecommendations handles getting audience recommendations
func (h *BusinessIntelligenceHandler) GetAudienceRecommendations(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	recommendations, err := h.biService.GetAudienceRecommendations(r.Context(), userID)
	if err != nil {
		h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"recommendations": recommendations,
		"count":           len(recommendations),
	}

	h.writeJSONResponse(w, response, http.StatusOK)
}

// Helper methods

func (h *BusinessIntelligenceHandler) writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to encode response"}`))
	}
}

func (h *BusinessIntelligenceHandler) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := map[string]string{"error": message}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.Write([]byte(`{"error": "Internal server error"}`))
	}
}