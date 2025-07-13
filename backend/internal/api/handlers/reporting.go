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

type ReportingHandler struct {
	reportingService ports.ReportingService
}

func NewReportingHandler(reportingService ports.ReportingService) *ReportingHandler {
	return &ReportingHandler{
		reportingService: reportingService,
	}
}

// Scheduled Reports

// CreateScheduledReport handles creating a new scheduled report
func (h *ReportingHandler) CreateScheduledReport(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var req domain.CreateScheduledReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	report, err := h.reportingService.CreateScheduledReport(r.Context(), userID, req)
	if err != nil {
		switch err {
		case domain.ErrInvalidRequest:
			h.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSONResponse(w, report, http.StatusCreated)
}

// GetScheduledReport handles getting a specific scheduled report
func (h *ReportingHandler) GetScheduledReport(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	reportIDStr := chi.URLParam(r, "id")
	reportID, err := strconv.ParseUint(reportIDStr, 10, 32)
	if err != nil {
		h.writeErrorResponse(w, "Invalid report ID", http.StatusBadRequest)
		return
	}

	report, err := h.reportingService.GetScheduledReport(r.Context(), uint(reportID), userID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			h.writeErrorResponse(w, "Report not found", http.StatusNotFound)
		case domain.ErrUnauthorized:
			h.writeErrorResponse(w, "Access denied", http.StatusForbidden)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSONResponse(w, report, http.StatusOK)
}

// GetUserScheduledReports handles getting all scheduled reports for a user
func (h *ReportingHandler) GetUserScheduledReports(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	reports, err := h.reportingService.GetUserScheduledReports(r.Context(), userID)
	if err != nil {
		h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"reports": reports,
		"total":   len(reports),
	}

	h.writeJSONResponse(w, response, http.StatusOK)
}

// UpdateScheduledReport handles updating a scheduled report
func (h *ReportingHandler) UpdateScheduledReport(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	reportIDStr := chi.URLParam(r, "id")
	reportID, err := strconv.ParseUint(reportIDStr, 10, 32)
	if err != nil {
		h.writeErrorResponse(w, "Invalid report ID", http.StatusBadRequest)
		return
	}

	var req domain.CreateScheduledReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	report, err := h.reportingService.UpdateScheduledReport(r.Context(), uint(reportID), userID, req)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			h.writeErrorResponse(w, "Report not found", http.StatusNotFound)
		case domain.ErrUnauthorized:
			h.writeErrorResponse(w, "Access denied", http.StatusForbidden)
		case domain.ErrInvalidRequest:
			h.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSONResponse(w, report, http.StatusOK)
}

// DeleteScheduledReport handles deleting a scheduled report
func (h *ReportingHandler) DeleteScheduledReport(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	reportIDStr := chi.URLParam(r, "id")
	reportID, err := strconv.ParseUint(reportIDStr, 10, 32)
	if err != nil {
		h.writeErrorResponse(w, "Invalid report ID", http.StatusBadRequest)
		return
	}

	err = h.reportingService.DeleteScheduledReport(r.Context(), uint(reportID), userID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			h.writeErrorResponse(w, "Report not found", http.StatusNotFound)
		case domain.ErrUnauthorized:
			h.writeErrorResponse(w, "Access denied", http.StatusForbidden)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Report Execution

// ExecuteReport handles manually executing a scheduled report
func (h *ReportingHandler) ExecuteReport(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	reportIDStr := chi.URLParam(r, "id")
	reportID, err := strconv.ParseUint(reportIDStr, 10, 32)
	if err != nil {
		h.writeErrorResponse(w, "Invalid report ID", http.StatusBadRequest)
		return
	}

	// Verify user owns the report
	_, err = h.reportingService.GetScheduledReport(r.Context(), uint(reportID), userID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			h.writeErrorResponse(w, "Report not found", http.StatusNotFound)
		case domain.ErrUnauthorized:
			h.writeErrorResponse(w, "Access denied", http.StatusForbidden)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Execute the report
	err = h.reportingService.ExecuteReport(r.Context(), uint(reportID))
	if err != nil {
		h.writeErrorResponse(w, "Failed to execute report", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":   "Report execution started",
		"report_id": reportID,
		"status":    "processing",
	}

	h.writeJSONResponse(w, response, http.StatusAccepted)
}

// GetReportHistory handles getting execution history for a report
func (h *ReportingHandler) GetReportHistory(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	reportIDStr := chi.URLParam(r, "id")
	reportID, err := strconv.ParseUint(reportIDStr, 10, 32)
	if err != nil {
		h.writeErrorResponse(w, "Invalid report ID", http.StatusBadRequest)
		return
	}

	history, err := h.reportingService.GetReportHistory(r.Context(), uint(reportID), userID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			h.writeErrorResponse(w, "Report not found", http.StatusNotFound)
		case domain.ErrUnauthorized:
			h.writeErrorResponse(w, "Access denied", http.StatusForbidden)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	response := map[string]interface{}{
		"report_id":  reportID,
		"executions": history,
		"total":      len(history),
	}

	h.writeJSONResponse(w, response, http.StatusOK)
}

// Data Export

// CreateDataExport handles creating a new data export
func (h *ReportingHandler) CreateDataExport(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var req domain.DataExportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	export, err := h.reportingService.CreateDataExport(r.Context(), userID, req)
	if err != nil {
		switch err {
		case domain.ErrInvalidRequest:
			h.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSONResponse(w, export, http.StatusCreated)
}

// GetDataExport handles getting a specific data export
func (h *ReportingHandler) GetDataExport(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	exportIDStr := chi.URLParam(r, "id")
	exportID, err := strconv.ParseUint(exportIDStr, 10, 32)
	if err != nil {
		h.writeErrorResponse(w, "Invalid export ID", http.StatusBadRequest)
		return
	}

	export, err := h.reportingService.GetDataExport(r.Context(), uint(exportID), userID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			h.writeErrorResponse(w, "Export not found", http.StatusNotFound)
		case domain.ErrUnauthorized:
			h.writeErrorResponse(w, "Access denied", http.StatusForbidden)
		default:
			h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSONResponse(w, export, http.StatusOK)
}

// GetUserDataExports handles getting all data exports for a user
func (h *ReportingHandler) GetUserDataExports(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	exports, err := h.reportingService.GetUserDataExports(r.Context(), userID)
	if err != nil {
		h.writeErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"exports": exports,
		"total":   len(exports),
	}

	h.writeJSONResponse(w, response, http.StatusOK)
}

// DownloadExport handles downloading a completed export
func (h *ReportingHandler) DownloadExport(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	exportIDStr := chi.URLParam(r, "id")
	exportID, err := strconv.ParseUint(exportIDStr, 10, 32)
	if err != nil {
		h.writeErrorResponse(w, "Invalid export ID", http.StatusBadRequest)
		return
	}

	data, filename, err := h.reportingService.DownloadExport(r.Context(), uint(exportID), userID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			h.writeErrorResponse(w, "Export not found", http.StatusNotFound)
		case domain.ErrUnauthorized:
			h.writeErrorResponse(w, "Access denied", http.StatusForbidden)
		default:
			h.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	// Set appropriate headers for file download
	w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// Report Generation

// GenerateAnalyticsReport handles generating an analytics report on-demand
func (h *ReportingHandler) GenerateAnalyticsReport(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var config domain.ReportConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	data, err := h.reportingService.GenerateAnalyticsReport(r.Context(), userID, config)
	if err != nil {
		h.writeErrorResponse(w, "Failed to generate report", http.StatusInternalServerError)
		return
	}

	// Determine content type based on config
	contentType := "application/json"
	if config.Grouping == "csv" {
		contentType = "text/csv"
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// GenerateDashboardReport handles generating a dashboard report
func (h *ReportingHandler) GenerateDashboardReport(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	dashboardIDStr := chi.URLParam(r, "dashboardId")
	dashboardID, err := strconv.ParseUint(dashboardIDStr, 10, 32)
	if err != nil {
		h.writeErrorResponse(w, "Invalid dashboard ID", http.StatusBadRequest)
		return
	}

	format := r.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}

	data, err := h.reportingService.GenerateDashboardReport(r.Context(), uint(dashboardID), userID, format)
	if err != nil {
		h.writeErrorResponse(w, "Failed to generate dashboard report", http.StatusInternalServerError)
		return
	}

	contentType := "application/json"
	if format == "pdf" {
		contentType = "application/pdf"
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// GenerateFunnelReport handles generating a funnel report
func (h *ReportingHandler) GenerateFunnelReport(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	funnelIDStr := chi.URLParam(r, "funnelId")
	funnelID, err := strconv.ParseUint(funnelIDStr, 10, 32)
	if err != nil {
		h.writeErrorResponse(w, "Invalid funnel ID", http.StatusBadRequest)
		return
	}

	format := r.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}

	data, err := h.reportingService.GenerateFunnelReport(r.Context(), uint(funnelID), userID, format)
	if err != nil {
		h.writeErrorResponse(w, "Failed to generate funnel report", http.StatusInternalServerError)
		return
	}

	contentType := "application/json"
	if format == "pdf" {
		contentType = "application/pdf"
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// Advanced Reporting

// GetReportTemplates handles getting available report templates
func (h *ReportingHandler) GetReportTemplates(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	templates := []map[string]interface{}{
		{
			"id":          "daily_summary",
			"name":        "Daily Summary",
			"description": "Daily overview of key metrics",
			"type":        "analytics",
			"format":      "pdf",
			"schedule":    "0 8 * * *",
			"config": map[string]interface{}{
				"metrics":    []string{"clicks", "conversions", "top_urls"},
				"charts":     true,
				"comparison": true,
			},
		},
		{
			"id":          "weekly_performance",
			"name":        "Weekly Performance",
			"description": "Comprehensive weekly performance analysis",
			"type":        "analytics",
			"format":      "excel",
			"schedule":    "0 9 * * MON",
			"config": map[string]interface{}{
				"metrics":    []string{"clicks", "conversions", "geographic", "devices"},
				"charts":     true,
				"comparison": true,
				"raw_data":   true,
			},
		},
		{
			"id":          "monthly_funnel",
			"name":        "Monthly Funnel Analysis",
			"description": "Detailed monthly funnel performance report",
			"type":        "funnel",
			"format":      "pdf",
			"schedule":    "0 10 1 * *",
			"config": map[string]interface{}{
				"include_optimizations": true,
				"comparison":           true,
				"cohort_analysis":      true,
			},
		},
		{
			"id":          "bi_dashboard",
			"name":        "Business Intelligence Dashboard",
			"description": "Comprehensive BI report with competitive analysis",
			"type":        "dashboard",
			"format":      "pdf",
			"schedule":    "0 10 * * MON",
			"config": map[string]interface{}{
				"include_predictive": true,
				"include_competitive": true,
				"include_recommendations": true,
			},
		},
	}

	response := map[string]interface{}{
		"templates": templates,
		"total":     len(templates),
	}

	h.writeJSONResponse(w, response, http.StatusOK)
}

// CreateReportFromTemplate handles creating a scheduled report from a template
func (h *ReportingHandler) CreateReportFromTemplate(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	templateID := chi.URLParam(r, "templateId")
	
	var customization struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Recipients  []string `json:"recipients"`
		Schedule    string   `json:"schedule,omitempty"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&customization); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get template configuration (mock implementation)
	template := map[string]interface{}{
		"type":   "analytics",
		"format": "pdf",
		"schedule": "0 9 * * MON",
		"config": domain.ReportConfig{
			DateRange:  "7d",
			Metrics:    []string{"clicks", "conversions"},
			Comparison: true,
			Charts:     true,
		},
	}

	schedule := customization.Schedule
	if schedule == "" {
		schedule = template["schedule"].(string)
	}

	// Create scheduled report request
	req := domain.CreateScheduledReportRequest{
		Name:        customization.Name,
		Description: customization.Description,
		ReportType:  template["type"].(string),
		Schedule:    schedule,
		Recipients:  customization.Recipients,
		Format:      template["format"].(string),
		Config:      template["config"].(domain.ReportConfig),
	}

	report, err := h.reportingService.CreateScheduledReport(r.Context(), userID, req)
	if err != nil {
		h.writeErrorResponse(w, "Failed to create report from template", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"report":       report,
		"template_id":  templateID,
		"message":      "Report created successfully from template",
	}

	h.writeJSONResponse(w, response, http.StatusCreated)
}

// GetReportInsights handles getting insights about reporting usage
func (h *ReportingHandler) GetReportInsights(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Mock insights data
	insights := map[string]interface{}{
		"total_reports":     3,
		"active_reports":    2,
		"reports_sent":      15,
		"last_30_days":      8,
		"success_rate":      94.7,
		"avg_generation_time": 12.5,
		"most_popular_format": "pdf",
		"most_popular_schedule": "weekly",
		"recent_activity": []map[string]interface{}{
			{
				"date":   "2024-07-13",
				"action": "report_sent",
				"report": "Weekly Performance",
				"status": "success",
			},
			{
				"date":   "2024-07-12",
				"action": "report_created",
				"report": "Daily Summary",
				"status": "success",
			},
		},
		"recommendations": []string{
			"Consider enabling weekly funnel analysis reports",
			"Your PDF reports are most engaged with - consider using this format more",
			"Adding competitive analysis to monthly reports could provide valuable insights",
		},
	}

	h.writeJSONResponse(w, insights, http.StatusOK)
}

// Helper methods

func (h *ReportingHandler) writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to encode response"}`))
	}
}

func (h *ReportingHandler) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := map[string]string{"error": message}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.Write([]byte(`{"error": "Internal server error"}`))
	}
}