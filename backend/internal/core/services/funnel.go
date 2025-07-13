package services

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"time"

	"url-shortener/internal/core/domain"
	"url-shortener/internal/core/ports"
)

type funnelService struct {
	urlRepo     ports.URLRepository
	clickRepo   ports.ClickRepository
	userRepo    ports.UserRepository
	cacheRepo   ports.CacheService
	configRepo  ports.ConfigService
}

func NewFunnelService(
	urlRepo ports.URLRepository,
	clickRepo ports.ClickRepository,
	userRepo ports.UserRepository,
	cacheRepo ports.CacheService,
	configRepo ports.ConfigService,
) ports.FunnelService {
	return &funnelService{
		urlRepo:    urlRepo,
		clickRepo:  clickRepo,
		userRepo:   userRepo,
		cacheRepo:  cacheRepo,
		configRepo: configRepo,
	}
}

// Funnel management

func (s *funnelService) CreateFunnel(ctx context.Context, userID uint, req domain.CreateFunnelRequest) (*domain.ConversionFunnel, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Create funnel
	funnel := &domain.ConversionFunnel{
		ID:          uint(time.Now().Unix()), // Mock ID
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Create funnel steps
	steps := make([]domain.FunnelStep, len(req.Steps))
	for i, stepReq := range req.Steps {
		steps[i] = domain.FunnelStep{
			ID:          uint(time.Now().Unix() + int64(i)), // Mock ID
			FunnelID:    funnel.ID,
			StepNumber:  i + 1,
			Name:        stepReq.Name,
			Description: stepReq.Description,
			URLPattern:  stepReq.URLPattern,
			EventType:   stepReq.EventType,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
	}

	funnel.Steps = steps

	// In a real implementation, you would save to repository
	return funnel, nil
}

func (s *funnelService) GetFunnel(ctx context.Context, funnelID uint, userID uint) (*domain.ConversionFunnel, error) {
	// Note: In a real implementation, you would fetch from repository
	// For now, return a mock funnel
	funnel := &domain.ConversionFunnel{
		ID:          funnelID,
		UserID:      userID,
		Name:        "E-commerce Conversion Funnel",
		Description: "Track user journey from awareness to purchase",
		IsActive:    true,
		CreatedAt:   time.Now().AddDate(0, -1, 0),
		UpdatedAt:   time.Now(),
		Steps: []domain.FunnelStep{
			{
				ID:          1,
				FunnelID:    funnelID,
				StepNumber:  1,
				Name:        "Product View",
				Description: "User views product page",
				URLPattern:  "/product/.*",
				EventType:   "view",
				CreatedAt:   time.Now().AddDate(0, -1, 0),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          2,
				FunnelID:    funnelID,
				StepNumber:  2,
				Name:        "Add to Cart",
				Description: "User adds product to cart",
				URLPattern:  "/cart/add",
				EventType:   "click",
				CreatedAt:   time.Now().AddDate(0, -1, 0),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          3,
				FunnelID:    funnelID,
				StepNumber:  3,
				Name:        "Checkout Start",
				Description: "User starts checkout process",
				URLPattern:  "/checkout",
				EventType:   "view",
				CreatedAt:   time.Now().AddDate(0, -1, 0),
				UpdatedAt:   time.Now(),
			},
			{
				ID:          4,
				FunnelID:    funnelID,
				StepNumber:  4,
				Name:        "Payment",
				Description: "User completes payment",
				URLPattern:  "/payment/success",
				EventType:   "conversion",
				CreatedAt:   time.Now().AddDate(0, -1, 0),
				UpdatedAt:   time.Now(),
			},
		},
	}

	// Verify ownership
	if funnel.UserID != userID {
		return nil, domain.ErrUnauthorized
	}

	return funnel, nil
}

func (s *funnelService) GetUserFunnels(ctx context.Context, userID uint) ([]*domain.ConversionFunnel, error) {
	// Note: In a real implementation, you would fetch from repository
	// For now, return mock funnels
	funnels := []*domain.ConversionFunnel{
		{
			ID:          1,
			UserID:      userID,
			Name:        "E-commerce Conversion",
			Description: "Product view to purchase conversion",
			IsActive:    true,
			CreatedAt:   time.Now().AddDate(0, -2, 0),
			UpdatedAt:   time.Now().AddDate(0, 0, -1),
		},
		{
			ID:          2,
			UserID:      userID,
			Name:        "Newsletter Signup",
			Description: "Landing page to newsletter subscription",
			IsActive:    true,
			CreatedAt:   time.Now().AddDate(0, -1, 0),
			UpdatedAt:   time.Now().AddDate(0, 0, -3),
		},
		{
			ID:          3,
			UserID:      userID,
			Name:        "Content Engagement",
			Description: "Blog view to resource download",
			IsActive:    false,
			CreatedAt:   time.Now().AddDate(0, -3, 0),
			UpdatedAt:   time.Now().AddDate(0, -1, 0),
		},
	}

	return funnels, nil
}

func (s *funnelService) UpdateFunnel(ctx context.Context, funnelID uint, userID uint, req domain.CreateFunnelRequest) (*domain.ConversionFunnel, error) {
	// Get existing funnel
	funnel, err := s.GetFunnel(ctx, funnelID, userID)
	if err != nil {
		return nil, err
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Update funnel
	funnel.Name = req.Name
	funnel.Description = req.Description
	funnel.UpdatedAt = time.Now()

	// Update steps
	steps := make([]domain.FunnelStep, len(req.Steps))
	for i, stepReq := range req.Steps {
		steps[i] = domain.FunnelStep{
			ID:          uint(time.Now().Unix() + int64(i)), // Mock ID
			FunnelID:    funnel.ID,
			StepNumber:  i + 1,
			Name:        stepReq.Name,
			Description: stepReq.Description,
			URLPattern:  stepReq.URLPattern,
			EventType:   stepReq.EventType,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
	}

	funnel.Steps = steps

	return funnel, nil
}

func (s *funnelService) DeleteFunnel(ctx context.Context, funnelID uint, userID uint) error {
	// Verify ownership
	funnel, err := s.GetFunnel(ctx, funnelID, userID)
	if err != nil {
		return err
	}
	if funnel.UserID != userID {
		return domain.ErrUnauthorized
	}

	// In a real implementation, you would delete from repository
	return nil
}

// Funnel analytics

func (s *funnelService) GetFunnelAnalytics(ctx context.Context, funnelID uint, userID uint, period string) (*domain.FunnelAnalyticsResponse, error) {
	// Get funnel
	funnel, err := s.GetFunnel(ctx, funnelID, userID)
	if err != nil {
		return nil, err
	}

	// Calculate analytics
	analytics := s.calculateFunnelAnalytics(ctx, funnel, period)

	// Generate optimization suggestions
	optimizations := s.generateFunnelOptimizations(ctx, funnel, analytics)

	response := &domain.FunnelAnalyticsResponse{
		FunnelAnalytics:          analytics,
		RecommendedOptimizations: optimizations,
	}

	return response, nil
}

func (s *funnelService) GetFunnelStepAnalytics(ctx context.Context, funnelID uint, stepID uint, userID uint) (*domain.FunnelStepAnalytics, error) {
	// Get funnel
	funnel, err := s.GetFunnel(ctx, funnelID, userID)
	if err != nil {
		return nil, err
	}

	// Find the step
	var step *domain.FunnelStep
	for _, s := range funnel.Steps {
		if s.ID == stepID {
			step = &s
			break
		}
	}

	if step == nil {
		return nil, domain.ErrNotFound
	}

	// Calculate step analytics
	analytics := s.calculateStepAnalytics(ctx, step)

	return analytics, nil
}

func (s *funnelService) GetConversionTrend(ctx context.Context, funnelID uint, userID uint, period string) (map[string]int64, error) {
	// Verify ownership
	funnel, err := s.GetFunnel(ctx, funnelID, userID)
	if err != nil {
		return nil, err
	}

	// Generate mock conversion trend data
	trend := s.generateConversionTrend(period)

	return trend, nil
}

// Funnel optimization

func (s *funnelService) GetFunnelOptimizations(ctx context.Context, funnelID uint, userID uint) ([]domain.OptimizationSuggestion, error) {
	// Get funnel analytics
	analyticsResponse, err := s.GetFunnelAnalytics(ctx, funnelID, userID, "30d")
	if err != nil {
		return nil, err
	}

	return analyticsResponse.RecommendedOptimizations, nil
}

func (s *funnelService) AnalyzeFunnelDropOffs(ctx context.Context, funnelID uint, userID uint) ([]domain.FunnelStepAnalytics, error) {
	// Get funnel
	funnel, err := s.GetFunnel(ctx, funnelID, userID)
	if err != nil {
		return nil, err
	}

	// Calculate drop-offs for each step
	dropOffs := make([]domain.FunnelStepAnalytics, len(funnel.Steps))
	
	for i, step := range funnel.Steps {
		stepAnalytics := s.calculateStepAnalytics(ctx, &step)
		dropOffs[i] = *stepAnalytics
	}

	// Sort by drop-off rate (highest first)
	sort.Slice(dropOffs, func(i, j int) bool {
		return dropOffs[i].DropOffRate > dropOffs[j].DropOffRate
	})

	return dropOffs, nil
}

// Helper methods

func (s *funnelService) calculateFunnelAnalytics(ctx context.Context, funnel *domain.ConversionFunnel, period string) *domain.FunnelAnalytics {
	// In a real implementation, you would query the database for actual data
	// For now, generate realistic mock data
	
	totalEntries := int64(10000)
	conversions := int64(850)
	conversionRate := float64(conversions) / float64(totalEntries) * 100
	dropOffRate := 100.0 - conversionRate

	// Calculate step analytics
	stepAnalytics := make([]domain.FunnelStepAnalytics, len(funnel.Steps))
	currentEntries := totalEntries

	for i, step := range funnel.Steps {
		// Simulate drop-off at each step
		var dropOffPercent float64
		switch i {
		case 0: // First step (Product View)
			dropOffPercent = 0.0 // Everyone who enters sees this
		case 1: // Add to Cart
			dropOffPercent = 25.0
		case 2: // Checkout Start
			dropOffPercent = 40.0
		case 3: // Payment
			dropOffPercent = 15.0
		default:
			dropOffPercent = 20.0
		}

		exits := int64(float64(currentEntries) * dropOffPercent / 100)
		conversionsAtStep := currentEntries - exits
		
		stepAnalytics[i] = domain.FunnelStepAnalytics{
			StepNumber:     step.StepNumber,
			StepName:       step.Name,
			Entries:        currentEntries,
			Exits:          exits,
			Conversions:    conversionsAtStep,
			ConversionRate: float64(conversionsAtStep) / float64(currentEntries) * 100,
			DropOffRate:    dropOffPercent,
			AvgTimeOnStep:  s.calculateAverageTimeOnStep(step.EventType),
		}

		currentEntries = conversionsAtStep
	}

	// Generate conversion trend
	conversionTrend := s.generateConversionTrend(period)

	return &domain.FunnelAnalytics{
		FunnelID:        funnel.ID,
		TotalEntries:    totalEntries,
		Conversions:     conversions,
		ConversionRate:  conversionRate,
		DropOffRate:     dropOffRate,
		StepAnalytics:   stepAnalytics,
		TimeToConvert:   s.generateTimeToConvert(),
		ConversionTrend: conversionTrend,
	}
}

func (s *funnelService) calculateStepAnalytics(ctx context.Context, step *domain.FunnelStep) *domain.FunnelStepAnalytics {
	// In a real implementation, you would query actual click data
	// For now, generate mock data based on step characteristics
	
	var entries, exits, conversions int64
	var conversionRate, dropOffRate float64

	switch step.StepNumber {
	case 1:
		entries = 10000
		exits = 0
		conversions = 10000
		conversionRate = 100.0
		dropOffRate = 0.0
	case 2:
		entries = 10000
		exits = 2500
		conversions = 7500
		conversionRate = 75.0
		dropOffRate = 25.0
	case 3:
		entries = 7500
		exits = 3000
		conversions = 4500
		conversionRate = 60.0
		dropOffRate = 40.0
	case 4:
		entries = 4500
		exits = 675
		conversions = 3825
		conversionRate = 85.0
		dropOffRate = 15.0
	default:
		entries = 1000
		exits = 200
		conversions = 800
		conversionRate = 80.0
		dropOffRate = 20.0
	}

	return &domain.FunnelStepAnalytics{
		StepNumber:     step.StepNumber,
		StepName:       step.Name,
		Entries:        entries,
		Exits:          exits,
		Conversions:    conversions,
		ConversionRate: conversionRate,
		DropOffRate:    dropOffRate,
		AvgTimeOnStep:  s.calculateAverageTimeOnStep(step.EventType),
	}
}

func (s *funnelService) calculateAverageTimeOnStep(eventType string) float64 {
	// Return realistic average times based on event type
	switch eventType {
	case "view":
		return 45.5 // 45.5 seconds average viewing time
	case "click":
		return 2.3 // 2.3 seconds to make click decision
	case "conversion":
		return 120.8 // 2 minutes to complete conversion action
	default:
		return 30.0
	}
}

func (s *funnelService) generateTimeToConvert() map[string]interface{} {
	return map[string]interface{}{
		"average_seconds": 1245.5,
		"median_seconds":  987.2,
		"percentiles": map[string]float64{
			"25th": 345.6,
			"50th": 987.2,
			"75th": 1456.8,
			"90th": 2134.5,
			"95th": 2890.1,
		},
		"distribution": map[string]int64{
			"0-5min":    2450,
			"5-15min":   3200,
			"15-30min":  2100,
			"30-60min":  1500,
			"60min+":    750,
		},
	}
}

func (s *funnelService) generateConversionTrend(period string) map[string]int64 {
	trend := make(map[string]int64)
	
	// Generate trend data based on period
	var days int
	switch period {
	case "7d":
		days = 7
	case "30d":
		days = 30
	case "90d":
		days = 90
	default:
		days = 30
	}

	baseConversions := int64(25)
	
	for i := 0; i < days; i++ {
		date := time.Now().AddDate(0, 0, -days+i+1).Format("2006-01-02")
		
		// Add some realistic variation
		variation := int64(float64(i%7) * 2.5) // Weekly pattern
		if i%7 == 0 || i%7 == 6 { // Lower on weekends
			variation -= 8
		}
		
		conversions := baseConversions + variation + int64(i/7) // Slight growth over time
		if conversions < 0 {
			conversions = 5
		}
		
		trend[date] = conversions
	}
	
	return trend
}

func (s *funnelService) generateFunnelOptimizations(ctx context.Context, funnel *domain.ConversionFunnel, analytics *domain.FunnelAnalytics) []domain.OptimizationSuggestion {
	suggestions := []domain.OptimizationSuggestion{}

	// Analyze each step for optimization opportunities
	for _, stepAnalytics := range analytics.StepAnalytics {
		if stepAnalytics.DropOffRate > 30.0 {
			suggestions = append(suggestions, domain.OptimizationSuggestion{
				Area:           "funnel_step",
				Suggestion:     fmt.Sprintf("Optimize step '%s' - high drop-off rate of %.1f%%", stepAnalytics.StepName, stepAnalytics.DropOffRate),
				Impact:         "high",
				Effort:         "medium",
				Priority:       "high",
				Potential:      stepAnalytics.DropOffRate * 0.3, // Potential 30% improvement
				Confidence:     85.5,
				Implementation: fmt.Sprintf("Review UX/UI of step %d, simplify form fields, add progress indicators, and test different messaging", stepAnalytics.StepNumber),
			})
		}

		if stepAnalytics.AvgTimeOnStep > 60.0 && stepAnalytics.StepNumber > 1 {
			suggestions = append(suggestions, domain.OptimizationSuggestion{
				Area:           "user_experience",
				Suggestion:     fmt.Sprintf("Reduce time on step '%s' - users spend %.1f seconds", stepAnalytics.StepName, stepAnalytics.AvgTimeOnStep),
				Impact:         "medium",
				Effort:         "low",
				Priority:       "medium",
				Potential:      15.2,
				Confidence:     72.8,
				Implementation: "Streamline the user interface, reduce form complexity, and add helpful tooltips or guidance",
			})
		}
	}

	// Overall conversion rate optimization
	if analytics.ConversionRate < 10.0 {
		suggestions = append(suggestions, domain.OptimizationSuggestion{
			Area:           "conversion_optimization",
			Suggestion:     fmt.Sprintf("Overall conversion rate of %.1f%% is below industry average", analytics.ConversionRate),
			Impact:         "high",
			Effort:         "high",
			Priority:       "high",
			Potential:      25.5,
			Confidence:     78.9,
			Implementation: "Implement A/B testing, improve value proposition, add social proof, and optimize mobile experience",
		})
	}

	// Step-specific optimizations
	for i, step := range funnel.Steps {
		if i < len(analytics.StepAnalytics) {
			stepAnalytics := analytics.StepAnalytics[i]
			
			switch step.EventType {
			case "view":
				if stepAnalytics.AvgTimeOnStep < 10.0 {
					suggestions = append(suggestions, domain.OptimizationSuggestion{
						Area:           "content_engagement",
						Suggestion:     fmt.Sprintf("Low engagement on '%s' - average time only %.1f seconds", step.Name, stepAnalytics.AvgTimeOnStep),
						Impact:         "medium",
						Effort:         "medium",
						Priority:       "medium",
						Potential:      18.3,
						Confidence:     68.5,
						Implementation: "Improve content quality, add interactive elements, and optimize page loading speed",
					})
				}
			case "click":
				if stepAnalytics.ConversionRate < 50.0 {
					suggestions = append(suggestions, domain.OptimizationSuggestion{
						Area:           "call_to_action",
						Suggestion:     fmt.Sprintf("Low click-through rate on '%s' - only %.1f%% conversion", step.Name, stepAnalytics.ConversionRate),
						Impact:         "high",
						Effort:         "low",
						Priority:       "high",
						Potential:      22.7,
						Confidence:     82.1,
						Implementation: "Test different button colors, text, and positioning. Add urgency or scarcity messaging",
					})
				}
			case "conversion":
				if stepAnalytics.AvgTimeOnStep > 120.0 {
					suggestions = append(suggestions, domain.OptimizationSuggestion{
						Area:           "checkout_optimization",
						Suggestion:     fmt.Sprintf("Long conversion time on '%s' - %.1f seconds average", step.Name, stepAnalytics.AvgTimeOnStep),
						Impact:         "high",
						Effort:         "medium",
						Priority:       "high",
						Potential:      28.9,
						Confidence:     85.7,
						Implementation: "Simplify checkout process, add guest checkout option, and optimize payment flow",
					})
				}
			}
		}
	}

	// Sort suggestions by priority and potential impact
	sort.Slice(suggestions, func(i, j int) bool {
		if suggestions[i].Priority != suggestions[j].Priority {
			priorityOrder := map[string]int{"high": 3, "medium": 2, "low": 1}
			return priorityOrder[suggestions[i].Priority] > priorityOrder[suggestions[j].Priority]
		}
		return suggestions[i].Potential > suggestions[j].Potential
	})

	return suggestions
}

func (s *funnelService) matchesURLPattern(url, pattern string) bool {
	// Simple regex matching for URL patterns
	matched, err := regexp.MatchString(pattern, url)
	if err != nil {
		return false
	}
	return matched
}

func (s *funnelService) trackFunnelEvent(ctx context.Context, funnel *domain.ConversionFunnel, url string, eventType string, userData map[string]interface{}) error {
	// In a real implementation, this would record funnel events
	// and update analytics in real-time
	
	// Find matching step
	for _, step := range funnel.Steps {
		if s.matchesURLPattern(url, step.URLPattern) && step.EventType == eventType {
			// Record event
			// This would typically insert into a funnel_events table
			break
		}
	}
	
	return nil
}

func (s *funnelService) calculateFunnelROI(ctx context.Context, funnelID uint, conversionValue float64) map[string]interface{} {
	// Calculate return on investment for funnel optimization
	return map[string]interface{}{
		"conversion_value":     conversionValue,
		"optimization_cost":    2500.0,
		"projected_improvement": 25.5,
		"projected_additional_revenue": conversionValue * 0.255,
		"roi_percentage": ((conversionValue * 0.255) - 2500.0) / 2500.0 * 100,
		"payback_period_days": 45,
	}
}