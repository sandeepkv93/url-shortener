package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"url-shortener/internal/core/domain"
	"url-shortener/internal/core/ports"
)

type businessIntelligenceService struct {
	analyticsService      ports.AnalyticsService
	urlRepo              ports.URLRepository
	clickRepo            ports.ClickRepository
	userRepo             ports.UserRepository
	cacheRepo            ports.CacheService
	predictiveService    ports.PredictiveAnalyticsService
	competitiveService   ports.CompetitiveIntelligenceService
}

func NewBusinessIntelligenceService(
	analyticsService ports.AnalyticsService,
	urlRepo ports.URLRepository,
	clickRepo ports.ClickRepository,
	userRepo ports.UserRepository,
	cacheRepo ports.CacheService,
	predictiveService ports.PredictiveAnalyticsService,
	competitiveService ports.CompetitiveIntelligenceService,
) ports.BusinessIntelligenceService {
	return &businessIntelligenceService{
		analyticsService:   analyticsService,
		urlRepo:           urlRepo,
		clickRepo:         clickRepo,
		userRepo:          userRepo,
		cacheRepo:         cacheRepo,
		predictiveService: predictiveService,
		competitiveService: competitiveService,
	}
}

// Dashboard management

func (s *businessIntelligenceService) CreateDashboard(ctx context.Context, userID uint, req domain.CreateDashboardRequest) (*domain.Dashboard, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Create dashboard
	dashboard := &domain.Dashboard{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		IsDefault:   req.IsDefault,
		IsPublic:    req.IsPublic,
		Layout:      req.Layout,
	}

	// Set default layout if not provided
	if dashboard.Layout.GridSize.Columns == 0 {
		dashboard.Layout = s.getDefaultDashboardLayout()
	}

	// If this is set as default, unset other defaults for this user
	if req.IsDefault {
		if err := s.unsetUserDefaultDashboards(ctx, userID); err != nil {
			return nil, fmt.Errorf("failed to unset existing default dashboards: %w", err)
		}
	}

	// Note: In a real implementation, you would use a repository
	// For now, we'll return a mock response
	dashboard.ID = uint(time.Now().Unix()) // Mock ID
	dashboard.CreatedAt = time.Now()
	dashboard.UpdatedAt = time.Now()

	return dashboard, nil
}

func (s *businessIntelligenceService) GetDashboard(ctx context.Context, dashboardID uint, userID uint) (*domain.Dashboard, error) {
	// Note: In a real implementation, you would fetch from repository
	// For now, return a mock dashboard with sample widgets
	dashboard := &domain.Dashboard{
		ID:          dashboardID,
		UserID:      userID,
		Name:        "Advanced Analytics Dashboard",
		Description: "Comprehensive business intelligence dashboard with predictive insights",
		IsDefault:   true,
		IsPublic:    false,
		Layout:      s.getDefaultDashboardLayout(),
		CreatedAt:   time.Now().AddDate(0, -1, 0), // Created a month ago
		UpdatedAt:   time.Now(),
		Widgets:     s.getDefaultWidgets(userID),
	}

	return dashboard, nil
}

func (s *businessIntelligenceService) GetUserDashboards(ctx context.Context, userID uint) ([]domain.DashboardResponse, error) {
	// Note: In a real implementation, you would fetch from repository
	// For now, return mock dashboards
	dashboards := []domain.DashboardResponse{
		{
			Dashboard: &domain.Dashboard{
				ID:          1,
				UserID:      userID,
				Name:        "Main Dashboard",
				Description: "Primary analytics dashboard",
				IsDefault:   true,
				IsPublic:    false,
				Layout:      s.getDefaultDashboardLayout(),
				CreatedAt:   time.Now().AddDate(0, -2, 0),
				UpdatedAt:   time.Now(),
			},
			WidgetCount: 6,
		},
		{
			Dashboard: &domain.Dashboard{
				ID:          2,
				UserID:      userID,
				Name:        "Performance Insights",
				Description: "Advanced performance metrics and predictions",
				IsDefault:   false,
				IsPublic:    false,
				Layout:      s.getDefaultDashboardLayout(),
				CreatedAt:   time.Now().AddDate(0, -1, 0),
				UpdatedAt:   time.Now().AddDate(0, 0, -1),
			},
			WidgetCount: 8,
		},
		{
			Dashboard: &domain.Dashboard{
				ID:          3,
				UserID:      userID,
				Name:        "Competitive Analysis",
				Description: "Market position and competitive intelligence",
				IsDefault:   false,
				IsPublic:    true,
				Layout:      s.getDefaultDashboardLayout(),
				CreatedAt:   time.Now().AddDate(0, 0, -15),
				UpdatedAt:   time.Now().AddDate(0, 0, -2),
			},
			WidgetCount: 4,
		},
	}

	return dashboards, nil
}

func (s *businessIntelligenceService) UpdateDashboard(ctx context.Context, dashboardID uint, userID uint, req domain.UpdateDashboardRequest) (*domain.Dashboard, error) {
	// Get existing dashboard
	dashboard, err := s.GetDashboard(ctx, dashboardID, userID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != nil {
		dashboard.Name = *req.Name
	}
	if req.Description != nil {
		dashboard.Description = *req.Description
	}
	if req.IsDefault != nil {
		dashboard.IsDefault = *req.IsDefault
		// If setting as default, unset other defaults
		if *req.IsDefault {
			if err := s.unsetUserDefaultDashboards(ctx, userID); err != nil {
				return nil, fmt.Errorf("failed to unset existing default dashboards: %w", err)
			}
		}
	}
	if req.IsPublic != nil {
		dashboard.IsPublic = *req.IsPublic
	}
	if req.Layout != nil {
		dashboard.Layout = *req.Layout
	}

	dashboard.UpdatedAt = time.Now()

	return dashboard, nil
}

func (s *businessIntelligenceService) DeleteDashboard(ctx context.Context, dashboardID uint, userID uint) error {
	// Verify ownership
	dashboard, err := s.GetDashboard(ctx, dashboardID, userID)
	if err != nil {
		return err
	}
	if dashboard.UserID != userID {
		return domain.ErrUnauthorized
	}

	// In a real implementation, you would delete from repository
	// For now, just return success
	return nil
}

// Widget management

func (s *businessIntelligenceService) CreateWidget(ctx context.Context, userID uint, req domain.CreateWidgetRequest) (*domain.DashboardWidget, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	widget := &domain.DashboardWidget{
		ID:          uint(time.Now().Unix()), // Mock ID
		UserID:      userID,
		Type:        req.Type,
		Title:       req.Title,
		Description: req.Description,
		Position:    req.Position,
		Size:        req.Size,
		Config:      req.Config,
		IsVisible:   true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return widget, nil
}

func (s *businessIntelligenceService) GetWidget(ctx context.Context, widgetID uint, userID uint) (*domain.DashboardWidget, error) {
	// Note: In a real implementation, you would fetch from repository
	widgets := s.getDefaultWidgets(userID)
	for _, widget := range widgets {
		if widget.ID == widgetID {
			return &widget, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (s *businessIntelligenceService) GetWidgetData(ctx context.Context, widgetID uint, userID uint) (*domain.WidgetDataResponse, error) {
	widget, err := s.GetWidget(ctx, widgetID, userID)
	if err != nil {
		return nil, err
	}

	// Generate mock data based on widget type and configuration
	data, metadata := s.generateWidgetData(ctx, widget, userID)

	return &domain.WidgetDataResponse{
		WidgetID:  widgetID,
		Data:      data,
		Metadata:  metadata,
		UpdatedAt: time.Now(),
	}, nil
}

func (s *businessIntelligenceService) UpdateWidget(ctx context.Context, widgetID uint, userID uint, req domain.UpdateWidgetRequest) (*domain.DashboardWidget, error) {
	widget, err := s.GetWidget(ctx, widgetID, userID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Title != nil {
		widget.Title = *req.Title
	}
	if req.Description != nil {
		widget.Description = *req.Description
	}
	if req.Position != nil {
		widget.Position = *req.Position
	}
	if req.Size != nil {
		widget.Size = *req.Size
	}
	if req.Config != nil {
		widget.Config = *req.Config
	}
	if req.IsVisible != nil {
		widget.IsVisible = *req.IsVisible
	}

	widget.UpdatedAt = time.Now()

	return widget, nil
}

func (s *businessIntelligenceService) DeleteWidget(ctx context.Context, widgetID uint, userID uint) error {
	widget, err := s.GetWidget(ctx, widgetID, userID)
	if err != nil {
		return err
	}
	if widget.UserID != userID {
		return domain.ErrUnauthorized
	}

	// In a real implementation, you would delete from repository
	return nil
}

// Advanced analytics

func (s *businessIntelligenceService) GetAdvancedAnalytics(ctx context.Context, userID uint, period string) (*domain.AdvancedAnalytics, error) {
	// Get basic analytics first
	dashboardStats, err := s.analyticsService.GetDashboardStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboard stats: %w", err)
	}

	// Generate advanced analytics
	analytics := &domain.AdvancedAnalytics{
		UserID:      userID,
		Period:      period,
		GeneratedAt: time.Now(),
	}

	// Performance metrics
	analytics.PerformanceMetrics = s.generatePerformanceMetrics(ctx, userID, dashboardStats)

	// Audience insights
	analytics.AudienceInsights = s.generateAudienceInsights(ctx, userID, period)

	// Content analytics
	analytics.ContentAnalytics = s.generateContentAnalytics(ctx, userID, period)

	// Competitive analysis
	if s.competitiveService != nil {
		competitiveAnalysis, _ := s.competitiveService.AnalyzeMarketPosition(ctx, userID)
		if competitiveAnalysis != nil {
			analytics.CompetitiveAnalysis = domain.CompetitiveAnalysis{
				MarketPosition: *competitiveAnalysis,
			}
		}
	}

	// Predictive insights
	if s.predictiveService != nil {
		predictiveInsights, _ := s.predictiveService.RecommendOptimizations(ctx, userID)
		if predictiveInsights != nil {
			analytics.PredictiveInsights = domain.PredictiveInsights{
				OptimizationSuggestions: predictiveInsights,
			}
		}
	}

	// Recommendation engine
	analytics.RecommendationEngine = s.generateRecommendations(ctx, userID, analytics)

	return analytics, nil
}

func (s *businessIntelligenceService) GetPerformanceMetrics(ctx context.Context, userID uint, period string) (*domain.PerformanceMetrics, error) {
	dashboardStats, err := s.analyticsService.GetDashboardStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	metrics := s.generatePerformanceMetrics(ctx, userID, dashboardStats)
	return &metrics, nil
}

func (s *businessIntelligenceService) GetAudienceInsights(ctx context.Context, userID uint, period string) (*domain.AudienceInsights, error) {
	insights := s.generateAudienceInsights(ctx, userID, period)
	return &insights, nil
}

func (s *businessIntelligenceService) GetContentAnalytics(ctx context.Context, userID uint, period string) (*domain.ContentAnalytics, error) {
	analytics := s.generateContentAnalytics(ctx, userID, period)
	return &analytics, nil
}

// Competitive analysis

func (s *businessIntelligenceService) GetCompetitiveAnalysis(ctx context.Context, userID uint) (*domain.CompetitiveAnalysis, error) {
	if s.competitiveService == nil {
		return s.generateMockCompetitiveAnalysis(userID), nil
	}

	marketPosition, err := s.competitiveService.AnalyzeMarketPosition(ctx, userID)
	if err != nil {
		return nil, err
	}

	opportunities, _ := s.competitiveService.IdentifyMarketOpportunities(ctx, userID)

	analysis := &domain.CompetitiveAnalysis{
		MarketPosition:  *marketPosition,
		OpportunityGaps: opportunities,
	}

	return analysis, nil
}

func (s *businessIntelligenceService) GetMarketPosition(ctx context.Context, userID uint) (*domain.MarketPosition, error) {
	if s.competitiveService == nil {
		return s.generateMockMarketPosition(), nil
	}

	return s.competitiveService.AnalyzeMarketPosition(ctx, userID)
}

func (s *businessIntelligenceService) GetBenchmarkData(ctx context.Context, userID uint, metric string) (*domain.BenchmarkData, error) {
	// Generate benchmark data based on industry standards
	benchmarks := &domain.BenchmarkData{
		IndustryAverage: map[string]float64{
			"ctr":             2.5,
			"conversion_rate": 3.2,
			"bounce_rate":     45.0,
			"engagement_score": 7.2,
		},
		TopPerformers: map[string]float64{
			"ctr":             8.5,
			"conversion_rate": 12.8,
			"bounce_rate":     25.0,
			"engagement_score": 9.5,
		},
		YourPerformance: map[string]float64{
			"ctr":             3.8,
			"conversion_rate": 4.2,
			"bounce_rate":     42.0,
			"engagement_score": 7.8,
		},
		Percentile: map[string]float64{
			"ctr":             65.0,
			"conversion_rate": 58.0,
			"bounce_rate":     48.0,
			"engagement_score": 72.0,
		},
	}

	return benchmarks, nil
}

// Predictive insights

func (s *businessIntelligenceService) GetPredictiveInsights(ctx context.Context, userID uint) (*domain.PredictiveInsights, error) {
	insights := &domain.PredictiveInsights{
		ForecastData: domain.ForecastData{
			Period: "30d",
			Predictions: map[string]float64{
				"clicks":       1250.0,
				"conversions":  45.0,
				"revenue":      890.0,
			},
			Confidence: map[string]float64{
				"clicks":       85.5,
				"conversions":  78.2,
				"revenue":      72.8,
			},
			Scenarios: map[string]domain.Scenario{
				"best": {
					Name:        "Best Case",
					Probability: 0.15,
					Predictions: map[string]float64{
						"clicks":       1850.0,
						"conversions":  68.0,
						"revenue":      1320.0,
					},
					Impact: "high",
				},
				"most_likely": {
					Name:        "Most Likely",
					Probability: 0.70,
					Predictions: map[string]float64{
						"clicks":       1250.0,
						"conversions":  45.0,
						"revenue":      890.0,
					},
					Impact: "medium",
				},
				"worst": {
					Name:        "Worst Case",
					Probability: 0.15,
					Predictions: map[string]float64{
						"clicks":       750.0,
						"conversions":  28.0,
						"revenue":      550.0,
					},
					Impact: "low",
				},
			},
		},
		AnomalyDetection: []domain.Anomaly{
			{
				Type:        "spike",
				Metric:      "clicks",
				Severity:    "medium",
				Detected:    time.Now().AddDate(0, 0, -2),
				Value:       450.0,
				Expected:    320.0,
				Deviation:   40.6,
				Description: "Unusual spike in click activity detected",
				Impact:      "positive",
			},
		},
		TrendPrediction: domain.TrendPrediction{
			ShortTerm: domain.TrendForecast{
				Direction:  "up",
				Magnitude:  15.2,
				Confidence: 82.5,
				FactorsInfluencing: []string{
					"recent content improvements",
					"seasonal trends",
					"increased social media activity",
				},
			},
			MediumTerm: domain.TrendForecast{
				Direction:  "stable",
				Magnitude:  3.8,
				Confidence: 75.0,
				FactorsInfluencing: []string{
					"market saturation",
					"competitive activity",
				},
			},
			LongTerm: domain.TrendForecast{
				Direction:  "up",
				Magnitude:  8.5,
				Confidence: 65.2,
				FactorsInfluencing: []string{
					"market expansion",
					"technology adoption",
				},
			},
		},
		RiskAssessment: domain.RiskAssessment{
			OverallRisk: domain.RiskLevel{
				Level:       "medium",
				Score:       45.5,
				Description: "Moderate risk level with manageable challenges",
				Trend:       "stable",
			},
			RiskFactors: []domain.RiskFactor{
				{
					Factor:      "market competition",
					Impact:      "medium",
					Probability: 0.65,
					Description: "Increasing competition in the market",
				},
				{
					Factor:      "technology changes",
					Impact:      "high",
					Probability: 0.35,
					Description: "Rapid technology evolution affecting user behavior",
				},
			},
			Mitigation: []domain.Mitigation{
				{
					Strategy:      "diversify content strategy",
					Priority:      "high",
					Effort:        "medium",
					Effectiveness: 78.5,
					Description:   "Expand content portfolio to reduce dependency on single sources",
				},
			},
		},
		OptimizationSuggestions: []domain.OptimizationSuggestion{
			{
				Area:           "content",
				Suggestion:     "Optimize URL titles for better engagement",
				Impact:         "high",
				Effort:         "low",
				Priority:       "high",
				Potential:      25.5,
				Confidence:     85.2,
				Implementation: "Update URL titles with action-oriented language and relevant keywords",
			},
			{
				Area:           "timing",
				Suggestion:     "Adjust posting schedule to peak engagement hours",
				Impact:         "medium",
				Effort:         "low",
				Priority:       "medium",
				Potential:      18.3,
				Confidence:     78.8,
				Implementation: "Schedule posts between 10-11 AM and 2-3 PM on weekdays",
			},
		},
	}

	return insights, nil
}

func (s *businessIntelligenceService) GetForecastData(ctx context.Context, userID uint, metric string, period string) (*domain.ForecastData, error) {
	insights, err := s.GetPredictiveInsights(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &insights.ForecastData, nil
}

func (s *businessIntelligenceService) DetectAnomalies(ctx context.Context, userID uint, metric string) ([]domain.Anomaly, error) {
	insights, err := s.GetPredictiveInsights(ctx, userID)
	if err != nil {
		return nil, err
	}

	return insights.AnomalyDetection, nil
}

func (s *businessIntelligenceService) GetTrendPrediction(ctx context.Context, userID uint, metric string) (*domain.TrendPrediction, error) {
	insights, err := s.GetPredictiveInsights(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &insights.TrendPrediction, nil
}

// Recommendations

func (s *businessIntelligenceService) GetRecommendations(ctx context.Context, userID uint) (*domain.RecommendationEngine, error) {
	analytics, err := s.GetAdvancedAnalytics(ctx, userID, "30d")
	if err != nil {
		return nil, err
	}

	return &analytics.RecommendationEngine, nil
}

func (s *businessIntelligenceService) GetOptimizationSuggestions(ctx context.Context, userID uint) ([]domain.OptimizationSuggestion, error) {
	insights, err := s.GetPredictiveInsights(ctx, userID)
	if err != nil {
		return nil, err
	}

	return insights.OptimizationSuggestions, nil
}

func (s *businessIntelligenceService) GetContentRecommendations(ctx context.Context, userID uint) ([]domain.ContentRec, error) {
	recommendations := []domain.ContentRec{
		{
			ContentType: "link_optimization",
			Topics:      []string{"technology", "productivity", "business"},
			Timing:      "morning_peak",
			Channels:    []string{"social_media", "email", "direct"},
			Potential:   85.5,
		},
		{
			ContentType: "viral_content",
			Topics:      []string{"trending", "entertainment", "lifestyle"},
			Timing:      "evening_engagement",
			Channels:    []string{"social_media", "forums"},
			Potential:   72.3,
		},
	}

	return recommendations, nil
}

func (s *businessIntelligenceService) GetAudienceRecommendations(ctx context.Context, userID uint) ([]domain.AudienceRec, error) {
	recommendations := []domain.AudienceRec{
		{
			SegmentName:    "tech_professionals",
			Characteristics: []string{"high_income", "mobile_heavy", "social_active"},
			Size:           15000,
			Potential:      78.5,
			Approach:       "linkedin_content_marketing",
		},
		{
			SegmentName:    "content_creators",
			Characteristics: []string{"creative", "social_native", "early_adopters"},
			Size:           8500,
			Potential:      65.2,
			Approach:       "influencer_partnerships",
		},
	}

	return recommendations, nil
}

// Helper methods

func (s *businessIntelligenceService) getDefaultDashboardLayout() domain.DashboardLayout {
	return domain.DashboardLayout{
		GridSize: domain.GridSize{
			Columns: 12,
			Rows:    20,
		},
		Breakpoints: domain.Breakpoint{
			XS: 480,
			SM: 768,
			MD: 992,
			LG: 1200,
			XL: 1600,
		},
		Margin:    []int{16, 16},
		Padding:   []int{16, 16, 16, 16},
		RowHeight: 60,
	}
}

func (s *businessIntelligenceService) getDefaultWidgets(userID uint) []domain.DashboardWidget {
	now := time.Now()
	return []domain.DashboardWidget{
		{
			ID:          1,
			UserID:      userID,
			Type:        "metric",
			Title:       "Total Clicks",
			Description: "Total number of clicks across all URLs",
			Position:    domain.WidgetPosition{X: 0, Y: 0, Z: 0},
			Size:        domain.WidgetSize{Width: 3, Height: 2},
			Config: domain.WidgetConfig{
				DataSource:   "clicks",
				TimeRange:    "30d",
				RefreshRate:  60,
				ShowLegend:   false,
			},
			IsVisible: true,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:          2,
			UserID:      userID,
			Type:        "chart",
			Title:       "Click Trend",
			Description: "Click performance over time",
			Position:    domain.WidgetPosition{X: 3, Y: 0, Z: 0},
			Size:        domain.WidgetSize{Width: 6, Height: 4},
			Config: domain.WidgetConfig{
				ChartType:    "line",
				DataSource:   "clicks",
				TimeRange:    "30d",
				RefreshRate:  300,
				ShowLegend:   true,
				ShowGrid:     true,
				Colors:       []string{"#3B82F6", "#10B981", "#F59E0B"},
			},
			IsVisible: true,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:          3,
			UserID:      userID,
			Type:        "map",
			Title:       "Geographic Distribution",
			Description: "Click distribution by location",
			Position:    domain.WidgetPosition{X: 9, Y: 0, Z: 0},
			Size:        domain.WidgetSize{Width: 3, Height: 4},
			Config: domain.WidgetConfig{
				DataSource:  "geo",
				TimeRange:   "30d",
				RefreshRate: 600,
				Limit:       10,
			},
			IsVisible: true,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:          4,
			UserID:      userID,
			Type:        "funnel",
			Title:       "Conversion Funnel",
			Description: "User conversion through different stages",
			Position:    domain.WidgetPosition{X: 0, Y: 4, Z: 0},
			Size:        domain.WidgetSize{Width: 6, Height: 3},
			Config: domain.WidgetConfig{
				DataSource:  "funnel",
				TimeRange:   "30d",
				RefreshRate: 300,
			},
			IsVisible: true,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:          5,
			UserID:      userID,
			Type:        "table",
			Title:       "Top Performing URLs",
			Description: "URLs with highest engagement",
			Position:    domain.WidgetPosition{X: 6, Y: 4, Z: 0},
			Size:        domain.WidgetSize{Width: 6, Height: 3},
			Config: domain.WidgetConfig{
				DataSource:  "urls",
				TimeRange:   "30d",
				RefreshRate: 600,
				Limit:       10,
			},
			IsVisible: true,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:          6,
			UserID:      userID,
			Type:        "chart",
			Title:       "Device Analytics",
			Description: "Click distribution by device type",
			Position:    domain.WidgetPosition{X: 0, Y: 7, Z: 0},
			Size:        domain.WidgetSize{Width: 4, Height: 3},
			Config: domain.WidgetConfig{
				ChartType:   "pie",
				DataSource:  "devices",
				TimeRange:   "30d",
				RefreshRate: 600,
				ShowLegend:  true,
				Colors:      []string{"#8B5CF6", "#06B6D4", "#84CC16", "#F97316"},
			},
			IsVisible: true,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

func (s *businessIntelligenceService) generateWidgetData(ctx context.Context, widget *domain.DashboardWidget, userID uint) (interface{}, map[string]interface{}) {
	metadata := map[string]interface{}{
		"last_updated": time.Now(),
		"data_source":  widget.Config.DataSource,
		"time_range":   widget.Config.TimeRange,
	}

	switch widget.Type {
	case "metric":
		return map[string]interface{}{
			"value":     12543,
			"change":    "+15.2%",
			"trend":     "up",
			"previous":  10892,
		}, metadata

	case "chart":
		switch widget.Config.ChartType {
		case "line":
			return s.generateTimeSeriesData(), metadata
		case "pie", "donut":
			return s.generatePieChartData(), metadata
		case "bar":
			return s.generateBarChartData(), metadata
		default:
			return s.generateTimeSeriesData(), metadata
		}

	case "table":
		return s.generateTableData(), metadata

	case "map":
		return s.generateMapData(), metadata

	case "funnel":
		return s.generateFunnelData(), metadata

	default:
		return map[string]interface{}{
			"message": "Widget data not available",
		}, metadata
	}
}

func (s *businessIntelligenceService) generateTimeSeriesData() interface{} {
	data := make([]map[string]interface{}, 30)
	baseValue := 100.0
	
	for i := 0; i < 30; i++ {
		// Add some realistic variation
		variation := (math.Sin(float64(i)*0.2) * 20) + (float64(i) * 0.5)
		value := baseValue + variation + (float64(i%7) * 5) // Weekly pattern
		
		data[i] = map[string]interface{}{
			"date":  time.Now().AddDate(0, 0, -29+i).Format("2006-01-02"),
			"value": int(value),
		}
	}
	
	return data
}

func (s *businessIntelligenceService) generatePieChartData() interface{} {
	return []map[string]interface{}{
		{"label": "Desktop", "value": 45, "color": "#3B82F6"},
		{"label": "Mobile", "value": 35, "color": "#10B981"},
		{"label": "Tablet", "value": 20, "color": "#F59E0B"},
	}
}

func (s *businessIntelligenceService) generateBarChartData() interface{} {
	return []map[string]interface{}{
		{"category": "Social Media", "value": 342},
		{"category": "Direct", "value": 298},
		{"category": "Search", "value": 187},
		{"category": "Email", "value": 156},
		{"category": "Referral", "value": 89},
	}
}

func (s *businessIntelligenceService) generateTableData() interface{} {
	return []map[string]interface{}{
		{
			"url":         "/abc123",
			"title":       "Product Launch Page",
			"clicks":      1543,
			"conversions": 87,
			"ctr":         "5.6%",
		},
		{
			"url":         "/def456",
			"title":       "Blog Article: Best Practices",
			"clicks":      1287,
			"conversions": 64,
			"ctr":         "5.0%",
		},
		{
			"url":         "/ghi789",
			"title":       "Download Whitepaper",
			"clicks":      892,
			"conversions": 178,
			"ctr":         "19.9%",
		},
	}
}

func (s *businessIntelligenceService) generateMapData() interface{} {
	return []map[string]interface{}{
		{"country": "US", "value": 435, "percentage": 34.5},
		{"country": "CA", "value": 298, "percentage": 23.6},
		{"country": "GB", "value": 187, "percentage": 14.8},
		{"country": "DE", "value": 156, "percentage": 12.4},
		{"country": "FR", "value": 89, "percentage": 7.1},
		{"country": "Others", "value": 95, "percentage": 7.6},
	}
}

func (s *businessIntelligenceService) generateFunnelData() interface{} {
	return []map[string]interface{}{
		{"step": "Awareness", "value": 1000, "percentage": 100.0},
		{"step": "Interest", "value": 750, "percentage": 75.0},
		{"step": "Consideration", "value": 400, "percentage": 40.0},
		{"step": "Purchase", "value": 120, "percentage": 12.0},
		{"step": "Retention", "value": 85, "percentage": 8.5},
	}
}

func (s *businessIntelligenceService) unsetUserDefaultDashboards(ctx context.Context, userID uint) error {
	// In a real implementation, you would update all user's dashboards to set is_default=false
	return nil
}

func (s *businessIntelligenceService) generatePerformanceMetrics(ctx context.Context, userID uint, stats *domain.DashboardStats) domain.PerformanceMetrics {
	// Calculate CTR (Click-through rate)
	ctr := 0.0
	if stats.TotalURLs > 0 {
		ctr = float64(stats.TotalClicks) / float64(stats.TotalURLs) * 100
	}

	return domain.PerformanceMetrics{
		CTR:            ctr,
		ConversionRate: 4.2,
		BounceRate:     42.5,
		EngagementScore: 7.8,
		QualityScore:   85.5,
		PerformanceTrend: map[string]float64{
			"ctr_7d":            3.2,
			"ctr_30d":           ctr,
			"conversion_7d":     3.8,
			"conversion_30d":    4.2,
			"engagement_7d":     7.5,
			"engagement_30d":    7.8,
		},
		BenchmarkComparison: map[string]interface{}{
			"industry_average": map[string]float64{
				"ctr":            2.5,
				"conversion":     3.1,
				"engagement":     6.8,
			},
			"your_performance": map[string]float64{
				"ctr":            ctr,
				"conversion":     4.2,
				"engagement":     7.8,
			},
			"percentile": 72.5,
		},
	}
}

func (s *businessIntelligenceService) generateAudienceInsights(ctx context.Context, userID uint, period string) domain.AudienceInsights {
	return domain.AudienceInsights{
		Demographics: domain.Demographics{
			AgeGroups: map[string]int64{
				"18-24": 156,
				"25-34": 432,
				"35-44": 387,
				"45-54": 298,
				"55+":   187,
			},
			Gender: map[string]int64{
				"male":   789,
				"female": 671,
				"other":  45,
			},
			Locations: map[string]int64{
				"North America": 654,
				"Europe":        432,
				"Asia":          289,
				"Other":         125,
			},
			Devices: map[string]int64{
				"mobile":  765,
				"desktop": 543,
				"tablet":  192,
			},
		},
		BehaviorPatterns: domain.BehaviorPatterns{
			ClickPatterns: map[string]interface{}{
				"peak_hours":     []int{10, 11, 14, 15},
				"peak_days":      []string{"Tuesday", "Wednesday", "Thursday"},
				"session_length": 145.5, // seconds
			},
			TimeOfDay: map[int]int64{
				9:  67, 10: 123, 11: 156, 12: 89, 13: 78,
				14: 145, 15: 167, 16: 134, 17: 98, 18: 76,
			},
			DayOfWeek: map[string]int64{
				"Monday":    145,
				"Tuesday":   189,
				"Wednesday": 203,
				"Thursday":  187,
				"Friday":    167,
				"Saturday":  98,
				"Sunday":    87,
			},
			ReturnVisitors: 234,
			AverageSession: 145.5,
		},
		SegmentAnalysis: []domain.AudienceSegment{
			{
				SegmentID: "power_users",
				Name:      "Power Users",
				Size:      156,
				Characteristics: map[string]interface{}{
					"avg_clicks_per_day": 12.5,
					"engagement_score":   9.2,
					"device_preference":  "desktop",
				},
				Performance: map[string]float64{
					"conversion_rate": 8.5,
					"retention_rate":  85.2,
				},
				GrowthRate: 15.3,
			},
			{
				SegmentID: "casual_users",
				Name:      "Casual Users",
				Size:      543,
				Characteristics: map[string]interface{}{
					"avg_clicks_per_day": 2.3,
					"engagement_score":   6.8,
					"device_preference":  "mobile",
				},
				Performance: map[string]float64{
					"conversion_rate": 3.2,
					"retention_rate":  45.8,
				},
				GrowthRate: 8.7,
			},
		},
		RetentionMetrics: domain.RetentionMetrics{
			RetentionRate:    72.5,
			ChurnRate:        27.5,
			LifetimeValue:    145.50,
			ReactivationRate: 15.2,
			RetentionCohort: map[string]float64{
				"week_1":  85.2,
				"week_2":  72.8,
				"week_4":  65.3,
				"week_8":  58.7,
				"week_12": 52.1,
			},
		},
	}
}

func (s *businessIntelligenceService) generateContentAnalytics(ctx context.Context, userID uint, period string) domain.ContentAnalytics {
	return domain.ContentAnalytics{
		TopPerformingContent: []domain.ContentPerformance{
			{
				ContentID:       "content_1",
				Title:           "Ultimate Guide to Productivity",
				URL:             "/abc123",
				Clicks:          1543,
				Shares:          89,
				EngagementScore: 9.2,
				Category:        "education",
				Tags:            []string{"productivity", "guide", "business"},
			},
			{
				ContentID:       "content_2",
				Title:           "Tech Trends 2024",
				URL:             "/def456",
				Clicks:          1287,
				Shares:          156,
				EngagementScore: 8.7,
				Category:        "technology",
				Tags:            []string{"technology", "trends", "future"},
			},
		},
		ContentCategories: map[string]int64{
			"technology": 456,
			"business":   387,
			"education":  298,
			"lifestyle":  189,
			"health":     134,
		},
		EngagementMetrics: domain.ContentEngagement{
			AverageTimeOnPage: 145.5,
			ShareRate:         12.3,
			CommentRate:       8.7,
			LikeRate:          23.4,
			EngagementTrend: map[string]float64{
				"week_1": 8.2,
				"week_2": 8.7,
				"week_3": 9.1,
				"week_4": 8.9,
			},
		},
		ViralityMetrics: domain.ViralityMetrics{
			ViralCoefficient:   1.25,
			ShareVelocity:      15.3,
			ReachAmplification: 2.8,
		},
	}
}

func (s *businessIntelligenceService) generateRecommendations(ctx context.Context, userID uint, analytics *domain.AdvancedAnalytics) domain.RecommendationEngine {
	return domain.RecommendationEngine{
		PersonalizedRecommendations: []domain.Recommendation{
			{
				Type:        "optimization",
				Title:       "Optimize Peak Hour Posting",
				Description: "Your audience is most active between 10-11 AM and 2-3 PM. Schedule important content during these times.",
				Priority:    "high",
				Impact:      "medium",
				Confidence:  85.2,
				ActionItems: []string{
					"Schedule content for 10:30 AM weekdays",
					"Avoid posting during low-engagement hours (6-8 AM)",
					"Set up automated scheduling tools",
				},
			},
			{
				Type:        "content",
				Title:       "Expand Technology Content",
				Description: "Technology-related content shows 45% higher engagement than average.",
				Priority:    "medium",
				Impact:      "high",
				Confidence:  78.9,
				ActionItems: []string{
					"Create more tech trend articles",
					"Partner with technology influencers",
					"Develop tech-focused content series",
				},
			},
		},
		ContentRecommendations: []domain.ContentRec{
			{
				ContentType: "trending_topics",
				Topics:      []string{"AI", "productivity", "remote_work"},
				Timing:      "morning_peak",
				Channels:    []string{"social_media", "email"},
				Potential:   82.5,
			},
		},
		OptimizationRecommendations: []domain.OptimizationRec{
			{
				Metric:       "click_through_rate",
				CurrentValue: 3.8,
				TargetValue:  5.2,
				Strategy:     "Improve title optimization and timing",
				Timeline:     "4-6 weeks",
			},
		},
		AudienceRecommendations: []domain.AudienceRec{
			{
				SegmentName:     "tech_professionals",
				Characteristics: []string{"high_engagement", "business_focused", "mobile_first"},
				Size:            2500,
				Potential:       75.3,
				Approach:        "Professional networking and LinkedIn content",
			},
		},
	}
}

func (s *businessIntelligenceService) generateMockCompetitiveAnalysis(userID uint) *domain.CompetitiveAnalysis {
	return &domain.CompetitiveAnalysis{
		MarketPosition: domain.MarketPosition{
			Rank:           15,
			MarketShare:    2.8,
			GrowthRate:     12.5,
			CompetitiveGap: -8.3,
			Strengths:      []string{"user_experience", "content_quality", "engagement_rate"},
			Weaknesses:     []string{"market_reach", "brand_awareness", "social_presence"},
		},
		CompetitorMetrics: []domain.CompetitorData{
			{
				CompetitorID:   "competitor_1",
				Name:           "Market Leader",
				MarketShare:    15.2,
				GrowthRate:     8.7,
				PerformanceGap: 12.4,
			},
		},
		BenchmarkData: domain.BenchmarkData{
			IndustryAverage: map[string]float64{
				"engagement": 6.5,
				"retention":  45.2,
				"growth":     8.9,
			},
			TopPerformers: map[string]float64{
				"engagement": 9.8,
				"retention":  78.5,
				"growth":     25.3,
			},
			YourPerformance: map[string]float64{
				"engagement": 7.8,
				"retention":  58.2,
				"growth":     12.5,
			},
		},
		OpportunityGaps: []domain.OpportunityGap{
			{
				Area:        "social_media_engagement",
				Impact:      "high",
				Effort:      "medium",
				Priority:    "high",
				Description: "Significant opportunity to improve social media presence",
				Potential:   35.5,
			},
		},
	}
}

func (s *businessIntelligenceService) generateMockMarketPosition() *domain.MarketPosition {
	return &domain.MarketPosition{
		Rank:           15,
		MarketShare:    2.8,
		GrowthRate:     12.5,
		CompetitiveGap: -8.3,
		Strengths:      []string{"user_experience", "content_quality", "engagement_rate"},
		Weaknesses:     []string{"market_reach", "brand_awareness", "social_presence"},
	}
}