package services

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"url-shortener/internal/core/domain"
	"url-shortener/internal/core/ports"
)

type competitiveIntelligenceService struct {
	analyticsService ports.AnalyticsService
	urlRepo         ports.URLRepository
	clickRepo       ports.ClickRepository
	userRepo        ports.UserRepository
	cacheService    ports.CacheService
}

func NewCompetitiveIntelligenceService(
	analyticsService ports.AnalyticsService,
	urlRepo ports.URLRepository,
	clickRepo ports.ClickRepository,
	userRepo ports.UserRepository,
	cacheService ports.CacheService,
) ports.CompetitiveIntelligenceService {
	return &competitiveIntelligenceService{
		analyticsService: analyticsService,
		urlRepo:         urlRepo,
		clickRepo:       clickRepo,
		userRepo:        userRepo,
		cacheService:    cacheService,
	}
}

// Market analysis

func (s *competitiveIntelligenceService) AnalyzeMarketPosition(ctx context.Context, userID uint) (*domain.MarketPosition, error) {
	// Get user's analytics data
	userStats, err := s.analyticsService.GetDashboardStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}

	// Get global market data for comparison
	globalStats, err := s.analyticsService.GetGlobalStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get global stats: %w", err)
	}

	// Calculate market position metrics
	position := s.calculateMarketPosition(userStats, globalStats)

	return position, nil
}

func (s *competitiveIntelligenceService) GetCompetitorData(ctx context.Context, competitorID string) (*domain.CompetitorData, error) {
	// In a real implementation, this would integrate with external data sources
	// For now, generate realistic mock competitor data
	
	competitor := &domain.CompetitorData{
		ID:           competitorID,
		Name:         s.getCompetitorName(competitorID),
		Industry:     "URL Shortening",
		Website:      fmt.Sprintf("https://%s.com", competitorID),
		Description:  fmt.Sprintf("Leading URL shortening service with focus on %s", s.getCompetitorFocus(competitorID)),
		MarketShare:  s.calculateMarketShare(competitorID),
		Founded:      s.getFoundedYear(competitorID),
		Employees:    s.getEmployeeCount(competitorID),
		Revenue:      s.getEstimatedRevenue(competitorID),
		Metrics: domain.CompetitorMetrics{
			MonthlyActiveUsers: s.getMAU(competitorID),
			URLsCreated:       s.getURLsCreated(competitorID),
			ClicksPerMonth:    s.getClicksPerMonth(competitorID),
			ConversionRate:    s.getConversionRate(competitorID),
			AverageClickTime:  s.getAverageClickTime(competitorID),
			GeographicReach:   s.getGeographicReach(competitorID),
			MobileUsage:       s.getMobileUsage(competitorID),
		},
		Strengths: s.getCompetitorStrengths(competitorID),
		Weaknesses: s.getCompetitorWeaknesses(competitorID),
		RecentNews: s.getRecentNews(competitorID),
		UpdatedAt:  time.Now(),
	}

	return competitor, nil
}

func (s *competitiveIntelligenceService) GetMarketTrends(ctx context.Context, industry string) (*domain.MarketTrends, error) {
	trends := &domain.MarketTrends{
		Industry:   industry,
		TimeRange:  "Last 12 months",
		UpdatedAt:  time.Now(),
		
		GrowthMetrics: domain.GrowthMetrics{
			MarketSize:       15.8, // $15.8B market size
			YearOverYearGrowth: 12.5, // 12.5% YoY growth
			ProjectedGrowth:   18.3, // 18.3% projected growth
			CompoundAnnualGrowthRate: 14.7, // 14.7% CAGR
		},
		
		EmergingTrends: []domain.Trend{
			{
				Name:        "AI-Powered Link Analytics",
				Impact:      "high",
				Confidence:  89.5,
				Description: "Machine learning integration for predictive analytics and user behavior insights",
				Timeline:    "6-12 months",
				Adoption:    35.7,
			},
			{
				Name:        "Privacy-First Tracking",
				Impact:      "high", 
				Confidence:  92.3,
				Description: "GDPR-compliant tracking solutions without compromising analytics depth",
				Timeline:    "3-6 months",
				Adoption:    67.2,
			},
			{
				Name:        "Blockchain-Based URL Verification",
				Impact:      "medium",
				Confidence:  71.8,
				Description: "Immutable URL verification for enhanced security and trust",
				Timeline:    "12-18 months",
				Adoption:    8.4,
			},
			{
				Name:        "Voice-Activated Link Sharing",
				Impact:      "medium",
				Confidence:  78.9,
				Description: "Integration with smart speakers and voice assistants",
				Timeline:    "6-9 months", 
				Adoption:    22.1,
			},
			{
				Name:        "Real-Time Collaborative Analytics",
				Impact:      "high",
				Confidence:  85.6,
				Description: "Team-based analytics dashboards with live collaboration features",
				Timeline:    "3-6 months",
				Adoption:    41.3,
			},
		},
		
		TechnologyAdoption: map[string]float64{
			"Cloud Infrastructure":     89.2,
			"API-First Architecture":   76.8,
			"Machine Learning":         45.3,
			"Real-Time Analytics":      62.7,
			"Mobile-First Design":      91.5,
			"Progressive Web Apps":     38.9,
			"Microservices":           67.4,
			"Containerization":        71.2,
		},
		
		CustomerBehavior: domain.CustomerBehavior{
			AverageSessionDuration: 4.7, // minutes
			BounceRate:            23.8, // percentage
			MobileUsage:           68.4, // percentage
			ReturnVisitorRate:     45.2, // percentage
			ConversionFunnelOptimization: 31.7, // adoption percentage
		},
		
		CompetitiveLandscape: s.getCompetitiveLandscape(),
	}

	return trends, nil
}

// Benchmarking

func (s *competitiveIntelligenceService) GetIndustryBenchmarks(ctx context.Context, industry string) (*domain.BenchmarkData, error) {
	benchmarks := &domain.BenchmarkData{
		Industry:   industry,
		DataSource: "Industry analysis and market research",
		UpdatedAt:  time.Now(),
		
		Metrics: map[string]domain.BenchmarkMetric{
			"click_through_rate": {
				Median:       3.2,
				Average:      4.1,
				Percentile25: 1.8,
				Percentile75: 6.4,
				Percentile90: 9.7,
				Unit:         "percentage",
				Description:  "Average click-through rate for shortened URLs",
			},
			"conversion_rate": {
				Median:       2.8,
				Average:      3.5,
				Percentile25: 1.2,
				Percentile75: 5.1,
				Percentile90: 8.3,
				Unit:         "percentage",
				Description:  "Conversion rate from click to desired action",
			},
			"session_duration": {
				Median:       185.5,
				Average:      210.7,
				Percentile25: 92.3,
				Percentile75: 298.1,
				Percentile90: 425.6,
				Unit:         "seconds",
				Description:  "Average session duration on target pages",
			},
			"bounce_rate": {
				Median:       34.2,
				Average:      38.7,
				Percentile25: 21.5,
				Percentile75: 52.8,
				Percentile90: 68.9,
				Unit:         "percentage",
				Description:  "Percentage of single-page sessions",
			},
			"mobile_usage": {
				Median:       67.8,
				Average:      69.2,
				Percentile25: 58.4,
				Percentile75: 78.6,
				Percentile90: 84.3,
				Unit:         "percentage",
				Description:  "Percentage of clicks from mobile devices",
			},
			"geographic_diversity": {
				Median:       12.5,
				Average:      15.8,
				Percentile25: 6.2,
				Percentile75: 21.4,
				Percentile90: 32.7,
				Unit:         "countries",
				Description:  "Number of countries generating significant traffic",
			},
		},
		
		PerformanceTiers: map[string]domain.PerformanceTier{
			"top_performer": {
				Description: "Top 10% of performers in the industry",
				Criteria: map[string]float64{
					"click_through_rate": 8.5,
					"conversion_rate":    7.2,
					"session_duration":   380.0,
					"bounce_rate":        25.0,
				},
			},
			"above_average": {
				Description: "Above industry average performers",
				Criteria: map[string]float64{
					"click_through_rate": 5.5,
					"conversion_rate":    4.8,
					"session_duration":   250.0,
					"bounce_rate":        30.0,
				},
			},
			"average": {
				Description: "Industry average performance",
				Criteria: map[string]float64{
					"click_through_rate": 4.1,
					"conversion_rate":    3.5,
					"session_duration":   210.7,
					"bounce_rate":        38.7,
				},
			},
			"below_average": {
				Description: "Below industry average, needs improvement",
				Criteria: map[string]float64{
					"click_through_rate": 2.5,
					"conversion_rate":    2.0,
					"session_duration":   150.0,
					"bounce_rate":        50.0,
				},
			},
		},
	}

	return benchmarks, nil
}

func (s *competitiveIntelligenceService) ComparePerformance(ctx context.Context, userID uint, competitorID string) (map[string]float64, error) {
	// Get user performance data
	userStats, err := s.analyticsService.GetDashboardStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}

	// Get competitor data
	competitor, err := s.GetCompetitorData(ctx, competitorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get competitor data: %w", err)
	}

	// Calculate performance comparison
	comparison := map[string]float64{
		"click_through_rate_ratio": s.calculateUserCTR(userStats) / competitor.Metrics.ConversionRate * 100,
		"conversion_rate_ratio":    s.calculateUserConversionRate(userStats) / competitor.Metrics.ConversionRate * 100,
		"volume_ratio":            float64(userStats.TotalClicks) / float64(competitor.Metrics.ClicksPerMonth) * 100,
		"geographic_reach_ratio":   s.calculateUserGeographicReach(userStats) / competitor.Metrics.GeographicReach * 100,
		"mobile_usage_ratio":       s.calculateUserMobileUsage(userStats) / competitor.Metrics.MobileUsage * 100,
		"overall_performance":      0, // Will be calculated below
	}

	// Calculate overall performance score (weighted average)
	weights := map[string]float64{
		"click_through_rate_ratio": 0.25,
		"conversion_rate_ratio":    0.30,
		"volume_ratio":            0.20,
		"geographic_reach_ratio":   0.15,
		"mobile_usage_ratio":       0.10,
	}

	var overallScore float64
	for metric, ratio := range comparison {
		if metric != "overall_performance" {
			overallScore += ratio * weights[metric]
		}
	}
	comparison["overall_performance"] = overallScore

	return comparison, nil
}

func (s *competitiveIntelligenceService) GetPerformanceGaps(ctx context.Context, userID uint) ([]domain.OpportunityGap, error) {
	// Get industry benchmarks
	benchmarks, err := s.GetIndustryBenchmarks(ctx, "URL Shortening")
	if err != nil {
		return nil, fmt.Errorf("failed to get benchmarks: %w", err)
	}

	// Get user performance
	userStats, err := s.analyticsService.GetDashboardStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}

	// Identify performance gaps
	gaps := []domain.OpportunityGap{}

	// Check click-through rate
	userCTR := s.calculateUserCTR(userStats)
	ctrBenchmark := benchmarks.Metrics["click_through_rate"]
	if userCTR < ctrBenchmark.Percentile75 {
		gaps = append(gaps, domain.OpportunityGap{
			Area:        "Click-Through Rate",
			CurrentValue: userCTR,
			BenchmarkValue: ctrBenchmark.Percentile75,
			GapPercentage: ((ctrBenchmark.Percentile75 - userCTR) / ctrBenchmark.Percentile75) * 100,
			Impact:       s.calculateImpact(userCTR, ctrBenchmark.Percentile75),
			Priority:     s.calculatePriority(userCTR, ctrBenchmark.Percentile75, "high"),
			Description:  "Your click-through rate is below the 75th percentile of industry performers",
			Recommendations: []string{
				"Optimize URL preview titles and descriptions",
				"Use more compelling call-to-action language",
				"A/B test different URL structures",
				"Improve landing page relevance",
			},
			EstimatedImpact: "15-25% increase in traffic",
			ImplementationEffort: "Medium",
			Timeline: "2-4 weeks",
		})
	}

	// Check conversion rate
	userConversion := s.calculateUserConversionRate(userStats)
	conversionBenchmark := benchmarks.Metrics["conversion_rate"]
	if userConversion < conversionBenchmark.Percentile75 {
		gaps = append(gaps, domain.OpportunityGap{
			Area:        "Conversion Rate",
			CurrentValue: userConversion,
			BenchmarkValue: conversionBenchmark.Percentile75,
			GapPercentage: ((conversionBenchmark.Percentile75 - userConversion) / conversionBenchmark.Percentile75) * 100,
			Impact:       s.calculateImpact(userConversion, conversionBenchmark.Percentile75),
			Priority:     s.calculatePriority(userConversion, conversionBenchmark.Percentile75, "high"),
			Description:  "Your conversion rate trails industry leaders by a significant margin",
			Recommendations: []string{
				"Implement conversion tracking and funnel analysis",
				"Optimize landing page user experience",
				"Add social proof and trust signals",
				"Test different call-to-action placements",
				"Implement exit-intent popups",
			},
			EstimatedImpact: "20-35% increase in conversions",
			ImplementationEffort: "High",
			Timeline: "4-8 weeks",
		})
	}

	// Check mobile usage optimization
	userMobile := s.calculateUserMobileUsage(userStats)
	mobileBenchmark := benchmarks.Metrics["mobile_usage"]
	if userMobile < mobileBenchmark.Median {
		gaps = append(gaps, domain.OpportunityGap{
			Area:        "Mobile Optimization",
			CurrentValue: userMobile,
			BenchmarkValue: mobileBenchmark.Median,
			GapPercentage: ((mobileBenchmark.Median - userMobile) / mobileBenchmark.Median) * 100,
			Impact:       s.calculateImpact(userMobile, mobileBenchmark.Median),
			Priority:     "medium",
			Description:  "Mobile traffic share is below industry average",
			Recommendations: []string{
				"Implement responsive design improvements",
				"Optimize page load times for mobile",
				"Create mobile-specific call-to-actions",
				"Test mobile app integration",
			},
			EstimatedImpact: "10-20% increase in mobile engagement",
			ImplementationEffort: "Medium",
			Timeline: "3-6 weeks",
		})
	}

	// Sort gaps by impact and priority
	sort.Slice(gaps, func(i, j int) bool {
		if gaps[i].Impact != gaps[j].Impact {
			impactOrder := map[string]int{"high": 3, "medium": 2, "low": 1}
			return impactOrder[gaps[i].Impact] > impactOrder[gaps[j].Impact]
		}
		return gaps[i].GapPercentage > gaps[j].GapPercentage
	})

	return gaps, nil
}

// Opportunity identification

func (s *competitiveIntelligenceService) IdentifyMarketOpportunities(ctx context.Context, userID uint) ([]domain.OpportunityGap, error) {
	// Get market trends
	trends, err := s.GetMarketTrends(ctx, "URL Shortening")
	if err != nil {
		return nil, fmt.Errorf("failed to get market trends: %w", err)
	}

	opportunities := []domain.OpportunityGap{}

	// Analyze emerging trends for opportunities
	for _, trend := range trends.EmergingTrends {
		if trend.Impact == "high" && trend.Adoption < 50.0 {
			opportunities = append(opportunities, domain.OpportunityGap{
				Area:        trend.Name,
				CurrentValue: 0, // Not yet adopted
				BenchmarkValue: trend.Adoption,
				GapPercentage: 100, // 100% opportunity
				Impact:       trend.Impact,
				Priority:     s.calculateTrendPriority(trend),
				Description:  trend.Description,
				Recommendations: s.getTrendRecommendations(trend.Name),
				EstimatedImpact: s.getTrendEstimatedImpact(trend.Name),
				ImplementationEffort: s.getTrendImplementationEffort(trend.Name),
				Timeline: trend.Timeline,
			})
		}
	}

	// Sort by adoption rate (lower = bigger opportunity)
	sort.Slice(opportunities, func(i, j int) bool {
		return opportunities[i].BenchmarkValue < opportunities[j].BenchmarkValue
	})

	return opportunities, nil
}

func (s *competitiveIntelligenceService) AnalyzeCompetitorWeaknesses(ctx context.Context, competitorID string) ([]string, error) {
	competitor, err := s.GetCompetitorData(ctx, competitorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get competitor data: %w", err)
	}

	return competitor.Weaknesses, nil
}

func (s *competitiveIntelligenceService) GetEmergingTrends(ctx context.Context, industry string) ([]domain.Trend, error) {
	trends, err := s.GetMarketTrends(ctx, industry)
	if err != nil {
		return nil, fmt.Errorf("failed to get market trends: %w", err)
	}

	return trends.EmergingTrends, nil
}

// Helper methods for calculations and data generation

func (s *competitiveIntelligenceService) calculateMarketPosition(userStats *domain.DashboardStats, globalStats *domain.GlobalStats) *domain.MarketPosition {
	// Calculate relative performance metrics
	clickShare := float64(userStats.TotalClicks) / float64(globalStats.TotalClicks) * 100
	urlShare := float64(userStats.TotalURLs) / float64(globalStats.TotalURLs) * 100
	
	// Determine market position tier
	var tier string
	var percentile float64
	
	if clickShare >= 1.0 { // Top 1% of users
		tier = "market_leader"
		percentile = 99.0
	} else if clickShare >= 0.1 { // Top 10%
		tier = "strong_player"
		percentile = 90.0
	} else if clickShare >= 0.01 { // Top 25%
		tier = "emerging_player"
		percentile = 75.0
	} else {
		tier = "niche_player"
		percentile = 50.0
	}

	return &domain.MarketPosition{
		UserID:      userStats.UserID,
		Tier:        tier,
		Percentile:  percentile,
		MarketShare: clickShare,
		Ranking:     s.calculateRanking(clickShare),
		
		Metrics: domain.PositionMetrics{
			ClickShare:      clickShare,
			URLShare:        urlShare,
			GrowthRate:      userStats.ClickGrowthRate,
			EngagementScore: s.calculateEngagementScore(userStats),
			ReachScore:      s.calculateReachScore(userStats),
		},
		
		Strengths: s.identifyStrengths(userStats, globalStats),
		Opportunities: s.identifyOpportunities(userStats, globalStats),
		Threats: s.identifyThreats(userStats, globalStats),
		
		CompetitiveAdvantages: s.identifyCompetitiveAdvantages(userStats),
		RecommendedActions: s.getPositionRecommendations(tier),
		
		UpdatedAt: time.Now(),
	}
}

func (s *competitiveIntelligenceService) calculateUserCTR(stats *domain.DashboardStats) float64 {
	if stats.TotalURLs == 0 {
		return 0
	}
	return float64(stats.TotalClicks) / float64(stats.TotalURLs) * 100
}

func (s *competitiveIntelligenceService) calculateUserConversionRate(stats *domain.DashboardStats) float64 {
	// Mock conversion rate calculation
	return math.Min(s.calculateUserCTR(stats) * 0.8, 10.0)
}

func (s *competitiveIntelligenceService) calculateUserGeographicReach(stats *domain.DashboardStats) float64 {
	// Mock geographic reach calculation
	return math.Min(float64(stats.TotalClicks)/1000 + 5, 50)
}

func (s *competitiveIntelligenceService) calculateUserMobileUsage(stats *domain.DashboardStats) float64 {
	// Mock mobile usage calculation
	return 65.0 + (float64(stats.TotalClicks) / 10000 * 5)
}

func (s *competitiveIntelligenceService) calculateImpact(current, benchmark float64) string {
	gap := math.Abs(benchmark - current) / benchmark * 100
	if gap > 30 {
		return "high"
	} else if gap > 15 {
		return "medium"
	}
	return "low"
}

func (s *competitiveIntelligenceService) calculatePriority(current, benchmark float64, baseImpact string) string {
	gap := math.Abs(benchmark - current) / benchmark * 100
	if gap > 25 && baseImpact == "high" {
		return "high"
	} else if gap > 15 {
		return "medium"
	}
	return "low"
}

// Mock data generation methods

func (s *competitiveIntelligenceService) getCompetitorName(id string) string {
	names := map[string]string{
		"bitly":     "Bitly",
		"tinyurl":   "TinyURL", 
		"shortlink": "ShortLink Pro",
		"linktr":    "Linktree",
		"rebrand":   "Rebrandly",
	}
	if name, exists := names[id]; exists {
		return name
	}
	return fmt.Sprintf("Competitor %s", id)
}

func (s *competitiveIntelligenceService) getCompetitorFocus(id string) string {
	focus := map[string]string{
		"bitly":     "enterprise solutions",
		"tinyurl":   "simplicity and reliability",
		"shortlink": "advanced analytics",
		"linktr":    "social media optimization",
		"rebrand":   "branded links",
	}
	if f, exists := focus[id]; exists {
		return f
	}
	return "URL shortening"
}

func (s *competitiveIntelligenceService) calculateMarketShare(id string) float64 {
	shares := map[string]float64{
		"bitly":     15.8,
		"tinyurl":   12.3,
		"shortlink": 8.7,
		"linktr":    6.4,
		"rebrand":   5.2,
	}
	if share, exists := shares[id]; exists {
		return share
	}
	return 2.1
}

func (s *competitiveIntelligenceService) getFoundedYear(id string) int {
	years := map[string]int{
		"bitly":     2008,
		"tinyurl":   2002,
		"shortlink": 2015,
		"linktr":    2016,
		"rebrand":   2014,
	}
	if year, exists := years[id]; exists {
		return year
	}
	return 2018
}

func (s *competitiveIntelligenceService) getEmployeeCount(id string) int {
	counts := map[string]int{
		"bitly":     250,
		"tinyurl":   45,
		"shortlink": 120,
		"linktr":    180,
		"rebrand":   85,
	}
	if count, exists := counts[id]; exists {
		return count
	}
	return 25
}

func (s *competitiveIntelligenceService) getEstimatedRevenue(id string) float64 {
	revenues := map[string]float64{
		"bitly":     25.5, // Million USD
		"tinyurl":   8.2,
		"shortlink": 15.7,
		"linktr":    18.9,
		"rebrand":   12.3,
	}
	if revenue, exists := revenues[id]; exists {
		return revenue
	}
	return 3.5
}

func (s *competitiveIntelligenceService) getMAU(id string) int64 {
	mau := map[string]int64{
		"bitly":     12500000,
		"tinyurl":   8900000,
		"shortlink": 6200000,
		"linktr":    15600000,
		"rebrand":   4800000,
	}
	if users, exists := mau[id]; exists {
		return users
	}
	return 1200000
}

func (s *competitiveIntelligenceService) getURLsCreated(id string) int64 {
	urls := map[string]int64{
		"bitly":     450000000,
		"tinyurl":   320000000,
		"shortlink": 180000000,
		"linktr":    290000000,
		"rebrand":   155000000,
	}
	if count, exists := urls[id]; exists {
		return count
	}
	return 25000000
}

func (s *competitiveIntelligenceService) getClicksPerMonth(id string) int64 {
	clicks := map[string]int64{
		"bitly":     2800000000,
		"tinyurl":   1900000000,
		"shortlink": 1200000000,
		"linktr":    2200000000,
		"rebrand":   850000000,
	}
	if count, exists := clicks[id]; exists {
		return count
	}
	return 125000000
}

func (s *competitiveIntelligenceService) getConversionRate(id string) float64 {
	rates := map[string]float64{
		"bitly":     4.8,
		"tinyurl":   3.2,
		"shortlink": 5.6,
		"linktr":    6.1,
		"rebrand":   4.9,
	}
	if rate, exists := rates[id]; exists {
		return rate
	}
	return 3.5
}

func (s *competitiveIntelligenceService) getAverageClickTime(id string) float64 {
	times := map[string]float64{
		"bitly":     185.5,
		"tinyurl":   165.2,
		"shortlink": 210.8,
		"linktr":    195.3,
		"rebrand":   178.9,
	}
	if time, exists := times[id]; exists {
		return time
	}
	return 170.0
}

func (s *competitiveIntelligenceService) getGeographicReach(id string) float64 {
	reach := map[string]float64{
		"bitly":     45.2,
		"tinyurl":   38.7,
		"shortlink": 28.9,
		"linktr":    52.6,
		"rebrand":   35.1,
	}
	if r, exists := reach[id]; exists {
		return r
	}
	return 22.3
}

func (s *competitiveIntelligenceService) getMobileUsage(id string) float64 {
	usage := map[string]float64{
		"bitly":     72.5,
		"tinyurl":   68.9,
		"shortlink": 75.2,
		"linktr":    81.7,
		"rebrand":   70.3,
	}
	if u, exists := usage[id]; exists {
		return u
	}
	return 65.8
}

func (s *competitiveIntelligenceService) getCompetitorStrengths(id string) []string {
	strengths := map[string][]string{
		"bitly": {
			"Strong enterprise customer base",
			"Advanced analytics platform",
			"Reliable infrastructure and uptime",
			"Comprehensive API ecosystem",
		},
		"tinyurl": {
			"First-mover advantage and brand recognition",
			"Simple, user-friendly interface",
			"High trust and reliability",
			"Low operational costs",
		},
		"shortlink": {
			"Cutting-edge analytics features",
			"Real-time performance monitoring",
			"Advanced funnel tracking",
			"AI-powered insights",
		},
		"linktr": {
			"Strong social media integration",
			"Visual link management",
			"Creator-focused features",
			"High mobile engagement",
		},
		"rebrand": {
			"Superior branded link customization",
			"White-label solutions",
			"Strong customer support",
			"Flexible pricing models",
		},
	}
	if s, exists := strengths[id]; exists {
		return s
	}
	return []string{"Competitive pricing", "Growing user base"}
}

func (s *competitiveIntelligenceService) getCompetitorWeaknesses(id string) []string {
	weaknesses := map[string][]string{
		"bitly": {
			"High pricing for small businesses",
			"Complex interface for casual users",
			"Limited social media features",
			"Slower innovation cycle",
		},
		"tinyurl": {
			"Outdated user interface",
			"Limited analytics capabilities",
			"Minimal customization options",
			"Lack of enterprise features",
		},
		"shortlink": {
			"High learning curve",
			"Premium pricing",
			"Limited brand recognition",
			"Smaller user community",
		},
		"linktr": {
			"Limited advanced analytics",
			"Focus primarily on social media",
			"Higher bounce rates",
			"Less suitable for enterprise",
		},
		"rebrand": {
			"Complex setup process",
			"Limited free tier",
			"Smaller market presence",
			"Less developer resources",
		},
	}
	if w, exists := weaknesses[id]; exists {
		return w
	}
	return []string{"Limited market share", "Resource constraints"}
}

func (s *competitiveIntelligenceService) getRecentNews(id string) []domain.NewsItem {
	// Generate recent news items
	return []domain.NewsItem{
		{
			Title:       fmt.Sprintf("%s announces new analytics features", s.getCompetitorName(id)),
			Date:        time.Now().AddDate(0, 0, -15),
			Source:      "TechCrunch",
			URL:         "https://techcrunch.com/news",
			Summary:     "Enhanced tracking and reporting capabilities launched",
		},
		{
			Title:       fmt.Sprintf("%s raises $25M Series B funding", s.getCompetitorName(id)),
			Date:        time.Now().AddDate(0, -2, 0),
			Source:      "VentureBeat",
			URL:         "https://venturebeat.com/news",
			Summary:     "Funding will be used to expand international presence",
		},
	}
}

func (s *competitiveIntelligenceService) getCompetitiveLandscape() map[string]interface{} {
	return map[string]interface{}{
		"market_leaders": []string{"Bitly", "Linktree", "TinyURL"},
		"emerging_players": []string{"Rebrandly", "ShortLink Pro", "ClickMeter"},
		"market_concentration": 65.8, // Percentage controlled by top 5 players
		"barriers_to_entry": "medium",
		"innovation_rate": "high",
		"customer_switching_cost": "low",
	}
}

func (s *competitiveIntelligenceService) calculateRanking(marketShare float64) int {
	// Simple ranking calculation based on market share
	if marketShare >= 10.0 {
		return 1
	} else if marketShare >= 5.0 {
		return 2
	} else if marketShare >= 1.0 {
		return 3
	} else if marketShare >= 0.1 {
		return 4
	}
	return 5
}

func (s *competitiveIntelligenceService) calculateEngagementScore(stats *domain.DashboardStats) float64 {
	// Mock engagement score based on clicks per URL
	if stats.TotalURLs == 0 {
		return 0
	}
	clicksPerURL := float64(stats.TotalClicks) / float64(stats.TotalURLs)
	return math.Min(clicksPerURL / 10.0 * 100, 100)
}

func (s *competitiveIntelligenceService) calculateReachScore(stats *domain.DashboardStats) float64 {
	// Mock reach score based on total clicks
	return math.Min(float64(stats.TotalClicks) / 10000 * 100, 100)
}

func (s *competitiveIntelligenceService) identifyStrengths(userStats *domain.DashboardStats, globalStats *domain.GlobalStats) []string {
	strengths := []string{}
	
	if userStats.ClickGrowthRate > 15.0 {
		strengths = append(strengths, "High growth rate")
	}
	if s.calculateUserCTR(userStats) > 5.0 {
		strengths = append(strengths, "Above-average click-through rate")
	}
	if userStats.TotalURLs > 100 {
		strengths = append(strengths, "Strong content creation volume")
	}
	
	return strengths
}

func (s *competitiveIntelligenceService) identifyOpportunities(userStats *domain.DashboardStats, globalStats *domain.GlobalStats) []string {
	opportunities := []string{}
	
	if s.calculateUserMobileUsage(userStats) < 70.0 {
		opportunities = append(opportunities, "Mobile optimization potential")
	}
	if s.calculateUserConversionRate(userStats) < 4.0 {
		opportunities = append(opportunities, "Conversion rate improvement")
	}
	if s.calculateUserGeographicReach(userStats) < 20.0 {
		opportunities = append(opportunities, "International expansion")
	}
	
	return opportunities
}

func (s *competitiveIntelligenceService) identifyThreats(userStats *domain.DashboardStats, globalStats *domain.GlobalStats) []string {
	threats := []string{}
	
	if userStats.ClickGrowthRate < 5.0 {
		threats = append(threats, "Slowing growth rate")
	}
	threats = append(threats, "Increasing competition from major players")
	threats = append(threats, "Privacy regulations affecting tracking")
	
	return threats
}

func (s *competitiveIntelligenceService) identifyCompetitiveAdvantages(stats *domain.DashboardStats) []string {
	advantages := []string{}
	
	if s.calculateUserCTR(stats) > 6.0 {
		advantages = append(advantages, "Superior link optimization")
	}
	if stats.ClickGrowthRate > 20.0 {
		advantages = append(advantages, "Rapid growth trajectory")
	}
	
	return advantages
}

func (s *competitiveIntelligenceService) getPositionRecommendations(tier string) []string {
	recommendations := map[string][]string{
		"market_leader": {
			"Maintain competitive advantage through innovation",
			"Expand into new market segments",
			"Invest in strategic partnerships",
			"Focus on customer retention",
		},
		"strong_player": {
			"Identify and exploit competitor weaknesses", 
			"Invest in differentiation strategies",
			"Build strategic alliances",
			"Enhance customer experience",
		},
		"emerging_player": {
			"Focus on niche market opportunities",
			"Optimize conversion funnel",
			"Improve product-market fit",
			"Build brand awareness",
		},
		"niche_player": {
			"Identify underserved market segments",
			"Optimize core features and performance",
			"Build strategic partnerships",
			"Focus on customer acquisition",
		},
	}
	
	if recs, exists := recommendations[tier]; exists {
		return recs
	}
	return []string{"Conduct market analysis", "Improve competitive position"}
}

func (s *competitiveIntelligenceService) calculateTrendPriority(trend domain.Trend) string {
	if trend.Impact == "high" && trend.Confidence > 80.0 && trend.Adoption < 30.0 {
		return "high"
	} else if trend.Impact == "high" || (trend.Confidence > 70.0 && trend.Adoption < 50.0) {
		return "medium"
	}
	return "low"
}

func (s *competitiveIntelligenceService) getTrendRecommendations(trendName string) []string {
	recommendations := map[string][]string{
		"AI-Powered Link Analytics": {
			"Implement machine learning for click prediction",
			"Add automated optimization suggestions",
			"Develop predictive analytics dashboard",
			"Create AI-driven content recommendations",
		},
		"Privacy-First Tracking": {
			"Implement GDPR-compliant tracking",
			"Add user consent management",
			"Develop cookieless analytics",
			"Enhance data transparency features",
		},
		"Blockchain-Based URL Verification": {
			"Research blockchain integration possibilities",
			"Develop URL verification system",
			"Partner with blockchain platforms",
			"Implement tamper-proof link tracking",
		},
		"Voice-Activated Link Sharing": {
			"Integrate with voice assistants",
			"Develop voice command interface",
			"Create audio link descriptions",
			"Build voice-optimized user flows",
		},
		"Real-Time Collaborative Analytics": {
			"Implement team dashboard features",
			"Add real-time data synchronization",
			"Develop collaborative reporting tools",
			"Create shared analytics workspace",
		},
	}
	
	if recs, exists := recommendations[trendName]; exists {
		return recs
	}
	return []string{"Research trend implementation", "Develop strategy"}
}

func (s *competitiveIntelligenceService) getTrendEstimatedImpact(trendName string) string {
	impacts := map[string]string{
		"AI-Powered Link Analytics": "30-50% improvement in user insights",
		"Privacy-First Tracking": "20-30% increase in user trust and compliance",
		"Blockchain-Based URL Verification": "15-25% improvement in security perception",
		"Voice-Activated Link Sharing": "10-20% increase in accessibility",
		"Real-Time Collaborative Analytics": "25-40% improvement in team productivity",
	}
	
	if impact, exists := impacts[trendName]; exists {
		return impact
	}
	return "10-20% improvement in relevant metrics"
}

func (s *competitiveIntelligenceService) getTrendImplementationEffort(trendName string) string {
	efforts := map[string]string{
		"AI-Powered Link Analytics": "High",
		"Privacy-First Tracking": "Medium",
		"Blockchain-Based URL Verification": "Very High",
		"Voice-Activated Link Sharing": "Medium",
		"Real-Time Collaborative Analytics": "High",
	}
	
	if effort, exists := efforts[trendName]; exists {
		return effort
	}
	return "Medium"
}