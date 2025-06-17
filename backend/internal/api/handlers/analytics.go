package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"url-shortener/internal/api/middleware"
	"url-shortener/internal/core/domain"
	"url-shortener/internal/core/ports"
)

type AnalyticsHandler struct {
	analyticsService ports.AnalyticsService
}

func NewAnalyticsHandler(analyticsService ports.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

// GetDashboard handles getting dashboard statistics for authenticated user
func (h *AnalyticsHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	stats, err := h.analyticsService.GetDashboardStats(r.Context(), userID)
	if err != nil {
		h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, stats, http.StatusOK)
}

// GetURLAnalytics handles getting analytics for a specific URL
func (h *AnalyticsHandler) GetURLAnalytics(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Parse URL ID
	urlIDStr := chi.URLParam(r, "id")
	urlID, err := strconv.ParseUint(urlIDStr, 10, 32)
	if err != nil {
		h.writeErrorResponse(w, "Invalid URL ID", http.StatusBadRequest)
		return
	}

	analytics, err := h.analyticsService.GetURLAnalytics(r.Context(), uint(urlID), userID)
	if err != nil {
		switch err {
		case domain.ErrUnauthorized:
			h.writeErrorResponse(w, "Access denied", http.StatusForbidden)
		case domain.ErrURLNotFound:
			h.writeErrorResponse(w, "URL not found", http.StatusNotFound)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSONResponse(w, analytics, http.StatusOK)
}

// GetTopPerformingURLs handles getting top performing URLs for user
func (h *AnalyticsHandler) GetTopPerformingURLs(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Parse limit parameter
	limit := 10
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 50 {
			limit = parsedLimit
		}
	}

	urls, err := h.analyticsService.GetTopPerformingURLs(r.Context(), userID, limit)
	if err != nil {
		h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"urls":  urls,
		"limit": limit,
	}

	h.writeJSONResponse(w, response, http.StatusOK)
}

// GetClickTimeline handles getting click timeline for a specific URL
func (h *AnalyticsHandler) GetClickTimeline(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Parse URL ID
	urlIDStr := chi.URLParam(r, "id")
	urlID, err := strconv.ParseUint(urlIDStr, 10, 32)
	if err != nil {
		h.writeErrorResponse(w, "Invalid URL ID", http.StatusBadRequest)
		return
	}

	// Parse period parameter
	period := r.URL.Query().Get("period")
	if period == "" {
		period = "week" // default period
	}

	// Validate period
	validPeriods := []string{"day", "week", "month", "year"}
	isValidPeriod := false
	for _, validPeriod := range validPeriods {
		if period == validPeriod {
			isValidPeriod = true
			break
		}
	}
	if !isValidPeriod {
		h.writeErrorResponse(w, "Invalid period. Valid values: day, week, month, year", http.StatusBadRequest)
		return
	}

	timeline, err := h.analyticsService.GetClickTimeline(r.Context(), uint(urlID), userID, period)
	if err != nil {
		switch err {
		case domain.ErrUnauthorized:
			h.writeErrorResponse(w, "Access denied", http.StatusForbidden)
		case domain.ErrURLNotFound:
			h.writeErrorResponse(w, "URL not found", http.StatusNotFound)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSONResponse(w, timeline, http.StatusOK)
}

// GetGeographicStats handles getting geographic statistics for a specific URL
func (h *AnalyticsHandler) GetGeographicStats(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Parse URL ID
	urlIDStr := chi.URLParam(r, "id")
	urlID, err := strconv.ParseUint(urlIDStr, 10, 32)
	if err != nil {
		h.writeErrorResponse(w, "Invalid URL ID", http.StatusBadRequest)
		return
	}

	geoStats, err := h.analyticsService.GetGeographicStats(r.Context(), uint(urlID), userID)
	if err != nil {
		switch err {
		case domain.ErrUnauthorized:
			h.writeErrorResponse(w, "Access denied", http.StatusForbidden)
		case domain.ErrURLNotFound:
			h.writeErrorResponse(w, "URL not found", http.StatusNotFound)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSONResponse(w, geoStats, http.StatusOK)
}

// GetDeviceStats handles getting device statistics for a specific URL
func (h *AnalyticsHandler) GetDeviceStats(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Parse URL ID
	urlIDStr := chi.URLParam(r, "id")
	urlID, err := strconv.ParseUint(urlIDStr, 10, 32)
	if err != nil {
		h.writeErrorResponse(w, "Invalid URL ID", http.StatusBadRequest)
		return
	}

	deviceStats, err := h.analyticsService.GetDeviceStats(r.Context(), uint(urlID), userID)
	if err != nil {
		switch err {
		case domain.ErrUnauthorized:
			h.writeErrorResponse(w, "Access denied", http.StatusForbidden)
		case domain.ErrURLNotFound:
			h.writeErrorResponse(w, "URL not found", http.StatusNotFound)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSONResponse(w, deviceStats, http.StatusOK)
}

// GetReferrerStats handles getting referrer statistics for a specific URL
func (h *AnalyticsHandler) GetReferrerStats(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Parse URL ID
	urlIDStr := chi.URLParam(r, "id")
	urlID, err := strconv.ParseUint(urlIDStr, 10, 32)
	if err != nil {
		h.writeErrorResponse(w, "Invalid URL ID", http.StatusBadRequest)
		return
	}

	referrerStats, err := h.analyticsService.GetReferrerStats(r.Context(), uint(urlID), userID)
	if err != nil {
		switch err {
		case domain.ErrUnauthorized:
			h.writeErrorResponse(w, "Access denied", http.StatusForbidden)
		case domain.ErrURLNotFound:
			h.writeErrorResponse(w, "URL not found", http.StatusNotFound)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	response := map[string]interface{}{
		"referrers": referrerStats,
	}

	h.writeJSONResponse(w, response, http.StatusOK)
}

// GetGlobalStats handles getting global platform statistics (admin only)
func (h *AnalyticsHandler) GetGlobalStats(w http.ResponseWriter, r *http.Request) {
	// For now, allow any authenticated user to see global stats
	// In production, you'd want to check for admin privileges
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	stats, err := h.analyticsService.GetGlobalStats(r.Context())
	if err != nil {
		h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, stats, http.StatusOK)
}

// ExportAnalytics handles exporting analytics data
func (h *AnalyticsHandler) ExportAnalytics(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Parse format parameter
	format := r.URL.Query().Get("format")
	if format == "" {
		format = "json" // default format
	}

	// Validate format
	validFormats := []string{"json", "csv"}
	isValidFormat := false
	for _, validFormat := range validFormats {
		if format == validFormat {
			isValidFormat = true
			break
		}
	}
	if !isValidFormat {
		h.writeErrorResponse(w, "Invalid format. Valid values: json, csv", http.StatusBadRequest)
		return
	}

	// Parse date range parameters
	dateRange, err := h.parseDateRange(r)
	if err != nil {
		h.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Export analytics
	data, err := h.analyticsService.ExportAnalytics(r.Context(), userID, format, dateRange)
	if err != nil {
		h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set appropriate content type and headers
	switch format {
	case "csv":
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename=analytics.csv")
	case "json":
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition", "attachment; filename=analytics.json")
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// Helper methods

func (h *AnalyticsHandler) parseDateRange(r *http.Request) (domain.DateRange, error) {
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	// Default to last 30 days if not specified
	if startDateStr == "" {
		startDate := time.Now().AddDate(0, 0, -30)
		startDateStr = startDate.Format("2006-01-02")
	}
	if endDateStr == "" {
		endDateStr = time.Now().Format("2006-01-02")
	}

	// Validate date format
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return domain.DateRange{}, err
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return domain.DateRange{}, err
	}

	// Validate date range
	if startDate.After(endDate) {
		return domain.DateRange{}, domain.ErrInvalidRequest
	}

	return domain.DateRange{
		StartDate: startDateStr,
		EndDate:   endDateStr,
	}, nil
}

func (h *AnalyticsHandler) writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to encode response"}`))
	}
}

func (h *AnalyticsHandler) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := map[string]string{"error": message}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.Write([]byte(`{"error": "Internal server error"}`))
	}
}