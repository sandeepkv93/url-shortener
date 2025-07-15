package services

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"url-shortener/internal/core/domain"
	"url-shortener/internal/core/ports"
)

type reportingService struct {
	analyticsService      ports.AnalyticsService
	biService            ports.BusinessIntelligenceService
	funnelService        ports.FunnelService
	urlRepo              ports.URLRepository
	clickRepo            ports.ClickRepository
	userRepo             ports.UserRepository
	notificationService  ports.NotificationService
	configRepo           ports.ConfigService
	dataExportService    *dataExportService
}

func NewReportingService(
	analyticsService ports.AnalyticsService,
	biService ports.BusinessIntelligenceService,
	funnelService ports.FunnelService,
	urlRepo ports.URLRepository,
	clickRepo ports.ClickRepository,
	userRepo ports.UserRepository,
	notificationService ports.NotificationService,
	configRepo ports.ConfigService,
) ports.ReportingService {
	// Initialize the enhanced data export service
	dataExportService := NewDataExportService(
		analyticsService,
		urlRepo,
		clickRepo,
		userRepo,
		notificationService,
	)

	return &reportingService{
		analyticsService:    analyticsService,
		biService:          biService,
		funnelService:      funnelService,
		urlRepo:            urlRepo,
		clickRepo:          clickRepo,
		userRepo:           userRepo,
		notificationService: notificationService,
		configRepo:         configRepo,
		dataExportService:  dataExportService,
	}
}

// Scheduled Reports

func (s *reportingService) CreateScheduledReport(ctx context.Context, userID uint, req domain.CreateScheduledReportRequest) (*domain.ScheduledReport, error) {
	// Validate request
	if err := s.validateScheduledReportRequest(req); err != nil {
		return nil, err
	}

	// Parse cron expression and calculate next run
	nextRun, err := s.calculateNextRun(req.Schedule)
	if err != nil {
		return nil, fmt.Errorf("invalid schedule format: %w", err)
	}

	report := &domain.ScheduledReport{
		ID:          uint(time.Now().Unix()), // Mock ID
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		ReportType:  req.ReportType,
		Schedule:    req.Schedule,
		Recipients:  req.Recipients,
		Format:      req.Format,
		Config:      req.Config,
		IsActive:    true,
		NextRun:     &nextRun,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// In a real implementation, you would save to repository
	return report, nil
}

func (s *reportingService) GetScheduledReport(ctx context.Context, reportID uint, userID uint) (*domain.ScheduledReport, error) {
	// Note: In a real implementation, you would fetch from repository
	nextRun := time.Now().Add(24 * time.Hour)
	lastRun := time.Now().Add(-24 * time.Hour)
	
	report := &domain.ScheduledReport{
		ID:          reportID,
		UserID:      userID,
		Name:        "Weekly Analytics Report",
		Description: "Comprehensive weekly analytics overview",
		ReportType:  "analytics",
		Schedule:    "0 9 * * MON", // Every Monday at 9 AM
		Recipients:  []string{"user@example.com", "manager@example.com"},
		Format:      "pdf",
		Config: domain.ReportConfig{
			DateRange:  "7d",
			Metrics:    []string{"clicks", "conversions", "geographic"},
			Grouping:   "daily",
			Comparison: true,
			Charts:     true,
			RawData:    false,
		},
		IsActive:  true,
		LastRun:   &lastRun,
		NextRun:   &nextRun,
		CreatedAt: time.Now().AddDate(0, -1, 0),
		UpdatedAt: time.Now(),
	}

	// Verify ownership
	if report.UserID != userID {
		return nil, domain.ErrUnauthorized
	}

	return report, nil
}

func (s *reportingService) GetUserScheduledReports(ctx context.Context, userID uint) ([]*domain.ScheduledReport, error) {
	// Note: In a real implementation, you would fetch from repository
	now := time.Now()
	
	reports := []*domain.ScheduledReport{
		{
			ID:          1,
			UserID:      userID,
			Name:        "Daily Performance Summary",
			Description: "Quick daily overview of key metrics",
			ReportType:  "dashboard",
			Schedule:    "0 8 * * *", // Daily at 8 AM
			Recipients:  []string{"user@example.com"},
			Format:      "excel",
			IsActive:    true,
			LastRun:     func() *time.Time { t := now.Add(-24 * time.Hour); return &t }(),
			NextRun:     func() *time.Time { t := now.Add(8 * time.Hour); return &t }(),
			CreatedAt:   now.AddDate(0, -2, 0),
			UpdatedAt:   now.AddDate(0, 0, -1),
		},
		{
			ID:          2,
			UserID:      userID,
			Name:        "Weekly Funnel Analysis",
			Description: "Detailed funnel performance and optimization insights",
			ReportType:  "funnel",
			Schedule:    "0 9 * * MON", // Weekly on Monday at 9 AM
			Recipients:  []string{"user@example.com", "analytics@example.com"},
			Format:      "pdf",
			IsActive:    true,
			LastRun:     func() *time.Time { t := now.Add(-7 * 24 * time.Hour); return &t }(),
			NextRun:     func() *time.Time { t := now.Add(24 * time.Hour); return &t }(),
			CreatedAt:   now.AddDate(0, -1, 0),
			UpdatedAt:   now.AddDate(0, 0, -3),
		},
		{
			ID:          3,
			UserID:      userID,
			Name:        "Monthly Business Intelligence",
			Description: "Comprehensive monthly BI report with competitive analysis",
			ReportType:  "competitive",
			Schedule:    "0 10 1 * *", // Monthly on 1st at 10 AM
			Recipients:  []string{"user@example.com", "management@example.com"},
			Format:      "pdf",
			IsActive:    false,
			LastRun:     func() *time.Time { t := now.AddDate(0, -1, 0); return &t }(),
			NextRun:     func() *time.Time { t := now.AddDate(0, 1, 0); return &t }(),
			CreatedAt:   now.AddDate(0, -3, 0),
			UpdatedAt:   now.AddDate(0, -1, 0),
		},
	}

	return reports, nil
}

func (s *reportingService) UpdateScheduledReport(ctx context.Context, reportID uint, userID uint, req domain.CreateScheduledReportRequest) (*domain.ScheduledReport, error) {
	// Get existing report
	report, err := s.GetScheduledReport(ctx, reportID, userID)
	if err != nil {
		return nil, err
	}

	// Validate request
	if err := s.validateScheduledReportRequest(req); err != nil {
		return nil, err
	}

	// Calculate next run if schedule changed
	if req.Schedule != report.Schedule {
		nextRun, err := s.calculateNextRun(req.Schedule)
		if err != nil {
			return nil, fmt.Errorf("invalid schedule format: %w", err)
		}
		report.NextRun = &nextRun
	}

	// Update fields
	report.Name = req.Name
	report.Description = req.Description
	report.ReportType = req.ReportType
	report.Schedule = req.Schedule
	report.Recipients = req.Recipients
	report.Format = req.Format
	report.Config = req.Config
	report.UpdatedAt = time.Now()

	return report, nil
}

func (s *reportingService) DeleteScheduledReport(ctx context.Context, reportID uint, userID uint) error {
	// Verify ownership
	report, err := s.GetScheduledReport(ctx, reportID, userID)
	if err != nil {
		return err
	}
	if report.UserID != userID {
		return domain.ErrUnauthorized
	}

	// In a real implementation, you would delete from repository
	return nil
}

// Report Execution

func (s *reportingService) ExecuteReport(ctx context.Context, reportID uint) error {
	// Note: In a real implementation, you would fetch the report from repository
	// For now, we'll simulate report execution
	
	// Get report details (mock)
	report := &domain.ScheduledReport{
		ID:         reportID,
		UserID:     1, // Mock user ID
		Name:       "Weekly Analytics Report",
		ReportType: "analytics",
		Format:     "pdf",
		Recipients: []string{"user@example.com"},
		Config: domain.ReportConfig{
			DateRange:  "7d",
			Metrics:    []string{"clicks", "conversions"},
			Comparison: true,
			Charts:     true,
		},
	}

	// Generate report
	reportData, err := s.generateReport(ctx, report)
	if err != nil {
		return fmt.Errorf("failed to generate report: %w", err)
	}

	// Send report via email
	if s.notificationService != nil {
		err = s.sendReportEmail(ctx, report, reportData)
		if err != nil {
			return fmt.Errorf("failed to send report email: %w", err)
		}
	}

	// Update last run time
	now := time.Now()
	report.LastRun = &now
	
	// Calculate next run
	nextRun, err := s.calculateNextRun(report.Schedule)
	if err == nil {
		report.NextRun = &nextRun
	}

	return nil
}

func (s *reportingService) GetReportHistory(ctx context.Context, reportID uint, userID uint) ([]domain.ReportExecution, error) {
	// Verify ownership
	_, err := s.GetScheduledReport(ctx, reportID, userID)
	if err != nil {
		return nil, err
	}

	// Note: In a real implementation, you would fetch from repository
	now := time.Now()
	
	executions := []domain.ReportExecution{
		{
			ID:          1,
			ReportID:    reportID,
			Status:      "completed",
			StartedAt:   now.Add(-7 * 24 * time.Hour),
			CompletedAt: func() *time.Time { t := now.Add(-7*24*time.Hour + 30*time.Second); return &t }(),
			Duration:    func() *int64 { d := int64(30000); return &d }(), // 30 seconds
			FilePath:    "/reports/weekly_analytics_20240701.pdf",
			FileSize:    2048576, // 2MB
			CreatedAt:   now.Add(-7 * 24 * time.Hour),
			UpdatedAt:   now.Add(-7 * 24 * time.Hour),
		},
		{
			ID:          2,
			ReportID:    reportID,
			Status:      "completed",
			StartedAt:   now.Add(-14 * 24 * time.Hour),
			CompletedAt: func() *time.Time { t := now.Add(-14*24*time.Hour + 45*time.Second); return &t }(),
			Duration:    func() *int64 { d := int64(45000); return &d }(), // 45 seconds
			FilePath:    "/reports/weekly_analytics_20240624.pdf",
			FileSize:    1987432,
			CreatedAt:   now.Add(-14 * 24 * time.Hour),
			UpdatedAt:   now.Add(-14 * 24 * time.Hour),
		},
		{
			ID:          3,
			ReportID:    reportID,
			Status:      "failed",
			StartedAt:   now.Add(-21 * 24 * time.Hour),
			CompletedAt: nil,
			Duration:    func() *int64 { d := int64(5000); return &d }(), // 5 seconds before failure
			FilePath:    "",
			FileSize:    0,
			ErrorMessage: "Database connection timeout",
			CreatedAt:   now.Add(-21 * 24 * time.Hour),
			UpdatedAt:   now.Add(-21 * 24 * time.Hour),
		},
	}

	return executions, nil
}

// Data Export

func (s *reportingService) CreateDataExport(ctx context.Context, userID uint, req domain.DataExportRequest) (*domain.DataExport, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	export := &domain.DataExport{
		ID:          uint(time.Now().Unix()), // Mock ID
		UserID:      userID,
		ExportType:  req.ExportType,
		Format:      req.Format,
		Status:      "pending",
		Config:      req.Config,
		ExpiresAt:   time.Now().Add(7 * 24 * time.Hour), // Expire in 7 days
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// In a real implementation, you would save to repository and queue for processing
	// For now, simulate immediate processing
	go s.processDataExport(ctx, export)

	return export, nil
}

func (s *reportingService) GetDataExport(ctx context.Context, exportID uint, userID uint) (*domain.DataExport, error) {
	// Note: In a real implementation, you would fetch from repository
	export := &domain.DataExport{
		ID:          exportID,
		UserID:      userID,
		ExportType:  "analytics",
		Format:      "csv",
		Status:      "completed",
		FilePath:    "/exports/analytics_data_20240701.csv",
		FileSize:    5242880, // 5MB
		RecordCount: 50000,
		Config: domain.ExportConfig{
			DateRange: domain.DateRange{
				StartDate: "2024-06-01",
				EndDate:   "2024-06-30",
			},
			Columns:     []string{"date", "clicks", "unique_clicks", "country", "device"},
			Compression: true,
		},
		ExpiresAt: time.Now().Add(5 * 24 * time.Hour),
		CreatedAt: time.Now().Add(-2 * time.Hour),
		UpdatedAt: time.Now().Add(-1 * time.Hour),
	}

	// Verify ownership
	if export.UserID != userID {
		return nil, domain.ErrUnauthorized
	}

	return export, nil
}

func (s *reportingService) GetUserDataExports(ctx context.Context, userID uint) ([]*domain.DataExport, error) {
	// Note: In a real implementation, you would fetch from repository
	now := time.Now()
	
	exports := []*domain.DataExport{
		{
			ID:          1,
			UserID:      userID,
			ExportType:  "analytics",
			Format:      "csv",
			Status:      "completed",
			FilePath:    "/exports/analytics_data_20240701.csv",
			FileSize:    5242880,
			RecordCount: 50000,
			ExpiresAt:   now.Add(5 * 24 * time.Hour),
			CreatedAt:   now.Add(-2 * time.Hour),
			UpdatedAt:   now.Add(-1 * time.Hour),
		},
		{
			ID:          2,
			UserID:      userID,
			ExportType:  "urls",
			Format:      "excel",
			Status:      "processing",
			FilePath:    "",
			FileSize:    0,
			RecordCount: 0,
			ExpiresAt:   now.Add(7 * 24 * time.Hour),
			CreatedAt:   now.Add(-30 * time.Minute),
			UpdatedAt:   now.Add(-5 * time.Minute),
		},
		{
			ID:          3,
			UserID:      userID,
			ExportType:  "clicks",
			Format:      "json",
			Status:      "failed",
			FilePath:    "",
			FileSize:    0,
			RecordCount: 0,
			ExpiresAt:   now.Add(6 * 24 * time.Hour),
			CreatedAt:   now.Add(-1 * 24 * time.Hour),
			UpdatedAt:   now.Add(-23 * time.Hour),
		},
	}

	return exports, nil
}

func (s *reportingService) DownloadExport(ctx context.Context, exportID uint, userID uint) ([]byte, string, error) {
	export, err := s.GetDataExport(ctx, exportID, userID)
	if err != nil {
		return nil, "", err
	}

	if export.Status != "completed" {
		return nil, "", fmt.Errorf("export is not ready for download, status: %s", export.Status)
	}

	if time.Now().After(export.ExpiresAt) {
		return nil, "", fmt.Errorf("export has expired")
	}

	// In a real implementation, you would read the file from storage
	// For now, generate sample data
	data, filename := s.generateSampleExportData(export)

	return data, filename, nil
}

// Report Generation

func (s *reportingService) GenerateAnalyticsReport(ctx context.Context, userID uint, config domain.ReportConfig) ([]byte, error) {
	// Get analytics data
	stats, err := s.analyticsService.GetDashboardStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get analytics data: %w", err)
	}

	// Generate report based on format
	switch config.Grouping {
	case "csv":
		return s.generateCSVReport(stats, config)
	case "json":
		return s.generateJSONReport(stats, config)
	default:
		return s.generateTextReport(stats, config)
	}
}

func (s *reportingService) GenerateDashboardReport(ctx context.Context, dashboardID uint, userID uint, format string) ([]byte, error) {
	if s.biService == nil {
		return nil, fmt.Errorf("business intelligence service not available")
	}

	dashboard, err := s.biService.GetDashboard(ctx, dashboardID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboard: %w", err)
	}

	// Generate dashboard report
	report := map[string]interface{}{
		"dashboard_name": dashboard.Name,
		"generated_at":   time.Now(),
		"widget_count":   len(dashboard.Widgets),
		"widgets":        dashboard.Widgets,
	}

	switch format {
	case "json":
		return json.MarshalIndent(report, "", "  ")
	default:
		return json.MarshalIndent(report, "", "  ")
	}
}

func (s *reportingService) GenerateFunnelReport(ctx context.Context, funnelID uint, userID uint, format string) ([]byte, error) {
	if s.funnelService == nil {
		return nil, fmt.Errorf("funnel service not available")
	}

	analytics, err := s.funnelService.GetFunnelAnalytics(ctx, funnelID, userID, "30d")
	if err != nil {
		return nil, fmt.Errorf("failed to get funnel analytics: %w", err)
	}

	switch format {
	case "json":
		return json.MarshalIndent(analytics, "", "  ")
	default:
		return json.MarshalIndent(analytics, "", "  ")
	}
}

// Helper methods

func (s *reportingService) validateScheduledReportRequest(req domain.CreateScheduledReportRequest) error {
	if req.Name == "" {
		return fmt.Errorf("report name is required")
	}
	if req.Schedule == "" {
		return fmt.Errorf("schedule is required")
	}
	if len(req.Recipients) == 0 {
		return fmt.Errorf("at least one recipient is required")
	}
	
	validFormats := []string{"pdf", "excel", "csv", "json"}
	isValidFormat := false
	for _, format := range validFormats {
		if req.Format == format {
			isValidFormat = true
			break
		}
	}
	if !isValidFormat {
		return fmt.Errorf("invalid format: %s", req.Format)
	}

	return nil
}

func (s *reportingService) calculateNextRun(cronExpression string) (time.Time, error) {
	// Simple cron parser for common expressions
	// In a real implementation, you would use a proper cron library
	
	now := time.Now()
	
	// Parse basic patterns
	switch cronExpression {
	case "0 8 * * *": // Daily at 8 AM
		next := time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, now.Location())
		if next.Before(now) {
			next = next.Add(24 * time.Hour)
		}
		return next, nil
	case "0 9 * * MON": // Weekly on Monday at 9 AM
		next := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location())
		for next.Weekday() != time.Monday || next.Before(now) {
			next = next.Add(24 * time.Hour)
		}
		return next, nil
	case "0 10 1 * *": // Monthly on 1st at 10 AM
		next := time.Date(now.Year(), now.Month(), 1, 10, 0, 0, 0, now.Location())
		if next.Before(now) {
			next = next.AddDate(0, 1, 0)
		}
		return next, nil
	default:
		// Default to daily at 9 AM
		next := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location())
		if next.Before(now) {
			next = next.Add(24 * time.Hour)
		}
		return next, nil
	}
}

func (s *reportingService) generateReport(ctx context.Context, report *domain.ScheduledReport) ([]byte, error) {
	switch report.ReportType {
	case "analytics":
		return s.GenerateAnalyticsReport(ctx, report.UserID, report.Config)
	case "dashboard":
		// Get user's default dashboard or first dashboard
		if s.biService != nil {
			dashboards, err := s.biService.GetUserDashboards(ctx, report.UserID)
			if err == nil && len(dashboards) > 0 {
				return s.GenerateDashboardReport(ctx, dashboards[0].ID, report.UserID, report.Format)
			}
		}
		return []byte("No dashboard available"), nil
	case "funnel":
		// Get user's first funnel
		if s.funnelService != nil {
			funnels, err := s.funnelService.GetUserFunnels(ctx, report.UserID)
			if err == nil && len(funnels) > 0 {
				return s.GenerateFunnelReport(ctx, funnels[0].ID, report.UserID, report.Format)
			}
		}
		return []byte("No funnel available"), nil
	default:
		return []byte("Unsupported report type"), nil
	}
}

func (s *reportingService) sendReportEmail(ctx context.Context, report *domain.ScheduledReport, data []byte) error {
	// Use the notification service to send the report via email
	if s.notificationService == nil {
		return fmt.Errorf("notification service not available")
	}
	
	// Send the report to all recipients
	err := s.notificationService.SendScheduledReport(ctx, report.Recipients, report, data)
	if err != nil {
		// Send failure notification to report owner
		user, userErr := s.userRepo.GetByID(ctx, report.UserID)
		if userErr == nil {
			s.notificationService.SendReportFailureAlert(ctx, user, report, err.Error())
		}
		return fmt.Errorf("failed to send report email: %w", err)
	}
	
	// Send success notification to report owner
	user, userErr := s.userRepo.GetByID(ctx, report.UserID)
	if userErr == nil {
		s.notificationService.SendReportGenerationNotification(ctx, user, report, true, "")
	}
	
	return nil
}

func (s *reportingService) processDataExport(ctx context.Context, export *domain.DataExport) {
	// Update status to processing
	export.Status = "processing"
	export.UpdatedAt = time.Now()
	
	var data []byte
	var err error
	var filename string
	
	// Generate export using enhanced data export service
	switch export.ExportType {
	case "analytics":
		filename = fmt.Sprintf("analytics_export_%s.%s", time.Now().Format("20060102"), export.Format)
		switch export.Format {
		case "csv":
			data, err = s.dataExportService.ExportAnalyticsToCSV(ctx, export.UserID, export.Config)
		case "excel":
			data, err = s.dataExportService.ExportAnalyticsToExcel(ctx, export.UserID, export.Config)
		case "pdf":
			data, err = s.dataExportService.ExportAnalyticsToPDF(ctx, export.UserID, export.Config)
		case "json":
			data, err = s.dataExportService.ExportAnalyticsToJSON(ctx, export.UserID, export.Config)
		}
	case "urls":
		filename = fmt.Sprintf("urls_export_%s.%s", time.Now().Format("20060102"), export.Format)
		switch export.Format {
		case "csv":
			data, err = s.dataExportService.ExportURLsToCSV(ctx, export.UserID, export.Config)
		default:
			// For other formats, use the original method
			data, filename = s.generateSampleExportData(export)
		}
	default:
		// Use original method for other export types
		data, filename = s.generateSampleExportData(export)
	}
	
	if err != nil {
		// Export failed
		export.Status = "failed"
		export.UpdatedAt = time.Now()
		
		// Send failure notification
		if s.notificationService != nil {
			user, userErr := s.userRepo.GetByID(ctx, export.UserID)
			if userErr == nil {
				s.notificationService.SendDataExportNotification(ctx, user, export)
			}
		}
		return
	}
	
	// Calculate actual record count based on export type and data
	recordCount := s.calculateRecordCount(export, data)
	
	// Update export status to completed
	export.Status = "completed"
	export.FilePath = "/exports/" + filename
	export.FileSize = int64(len(data))
	export.RecordCount = recordCount
	export.UpdatedAt = time.Now()
	
	// Send success notification to user that export is ready
	if s.notificationService != nil {
		user, err := s.userRepo.GetByID(ctx, export.UserID)
		if err == nil {
			s.notificationService.SendDataExportNotification(ctx, user, export)
		}
	}
}

func (s *reportingService) generateSampleExportData(export *domain.DataExport) ([]byte, string) {
	switch export.Format {
	case "csv":
		return s.generateSampleCSV(export), fmt.Sprintf("%s_data_%s.csv", export.ExportType, time.Now().Format("20060102"))
	case "json":
		return s.generateSampleJSON(export), fmt.Sprintf("%s_data_%s.json", export.ExportType, time.Now().Format("20060102"))
	default:
		return []byte("Sample export data"), fmt.Sprintf("%s_data_%s.txt", export.ExportType, time.Now().Format("20060102"))
	}
}

func (s *reportingService) generateSampleCSV(export *domain.DataExport) []byte {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	
	// Write header
	header := []string{"Date", "Clicks", "Unique Clicks", "Country", "Device"}
	writer.Write(header)
	
	// Write sample data
	for i := 0; i < 10; i++ {
		record := []string{
			time.Now().AddDate(0, 0, -i).Format("2006-01-02"),
			strconv.Itoa(100 + i*10),
			strconv.Itoa(80 + i*8),
			"US",
			"desktop",
		}
		writer.Write(record)
	}
	
	writer.Flush()
	return buf.Bytes()
}

func (s *reportingService) generateSampleJSON(export *domain.DataExport) []byte {
	data := map[string]interface{}{
		"export_type": export.ExportType,
		"generated_at": time.Now(),
		"records": []map[string]interface{}{
			{
				"date":          time.Now().Format("2006-01-02"),
				"clicks":        150,
				"unique_clicks": 120,
				"country":       "US",
				"device":        "desktop",
			},
		},
	}
	
	jsonData, _ := json.MarshalIndent(data, "", "  ")
	return jsonData
}

func (s *reportingService) generateCSVReport(stats *domain.DashboardStats, config domain.ReportConfig) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	
	// Write header
	header := []string{"Metric", "Value"}
	writer.Write(header)
	
	// Write data
	writer.Write([]string{"Total URLs", strconv.FormatInt(stats.TotalURLs, 10)})
	writer.Write([]string{"Active URLs", strconv.FormatInt(stats.ActiveURLs, 10)})
	writer.Write([]string{"Total Clicks", strconv.FormatInt(stats.TotalClicks, 10)})
	writer.Write([]string{"Click Growth Rate", fmt.Sprintf("%.2f%%", stats.ClickGrowthRate)})
	writer.Write([]string{"URL Growth Rate", fmt.Sprintf("%.2f%%", stats.URLGrowthRate)})
	
	writer.Flush()
	return buf.Bytes(), nil
}

func (s *reportingService) generateJSONReport(stats *domain.DashboardStats, config domain.ReportConfig) ([]byte, error) {
	report := map[string]interface{}{
		"generated_at":      time.Now(),
		"period":           config.DateRange,
		"total_urls":       stats.TotalURLs,
		"active_urls":      stats.ActiveURLs,
		"total_clicks":     stats.TotalClicks,
		"click_growth_rate": stats.ClickGrowthRate,
		"url_growth_rate":  stats.URLGrowthRate,
		"clicks_by_date":   stats.ClicksByDate,
		"top_urls":         stats.TopURLs,
		"recent_activity":  stats.RecentActivity,
	}
	
	return json.MarshalIndent(report, "", "  ")
}

func (s *reportingService) generateTextReport(stats *domain.DashboardStats, config domain.ReportConfig) ([]byte, error) {
	var buf strings.Builder
	
	buf.WriteString("ANALYTICS REPORT\n")
	buf.WriteString("================\n\n")
	buf.WriteString(fmt.Sprintf("Generated: %s\n", time.Now().Format(time.RFC3339)))
	buf.WriteString(fmt.Sprintf("Period: %s\n\n", config.DateRange))
	buf.WriteString("OVERVIEW\n")
	buf.WriteString("--------\n")
	buf.WriteString(fmt.Sprintf("Total URLs: %d\n", stats.TotalURLs))
	buf.WriteString(fmt.Sprintf("Active URLs: %d\n", stats.ActiveURLs))
	buf.WriteString(fmt.Sprintf("Total Clicks: %d\n", stats.TotalClicks))
	buf.WriteString(fmt.Sprintf("Click Growth Rate: %.2f%%\n", stats.ClickGrowthRate))
	buf.WriteString(fmt.Sprintf("URL Growth Rate: %.2f%%\n", stats.URLGrowthRate))
	
	return []byte(buf.String()), nil
}

// calculateRecordCount calculates the number of records in the exported data
func (s *reportingService) calculateRecordCount(export *domain.DataExport, data []byte) int64 {
	if len(data) == 0 {
		return 0
	}
	
	switch export.Format {
	case "csv":
		// Count CSV rows (excluding header)
		lines := strings.Split(string(data), "\n")
		recordCount := len(lines) - 1 // subtract header
		if recordCount < 0 {
			recordCount = 0
		}
		return int64(recordCount)
	case "json":
		// For JSON, try to parse and count array elements
		var jsonData interface{}
		if err := json.Unmarshal(data, &jsonData); err == nil {
			if arr, ok := jsonData.([]interface{}); ok {
				return int64(len(arr))
			}
		}
		return 1 // assume single object
	case "excel":
		// For Excel, we can't easily count without parsing, so estimate based on export type
		switch export.ExportType {
		case "analytics":
			return 10 // sample data rows
		case "urls":
			return 100 // estimated URL count
		case "clicks":
			return 1000 // estimated click count
		default:
			return 1
		}
	default:
		return 1
	}
}